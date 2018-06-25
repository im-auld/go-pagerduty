package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"time"
)

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}

// HTTPClient is an interface which declares the functionality we need from an
// HTTP client. This is to allow consumers to provide their own HTTP client as
// needed, without restricting them to only using *http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// defaultHTTPClient is our own default HTTP client. We use this, instead of
// http.DefaultClient, to avoid other packages tweaks to http.DefaultClient
// causing issues with our HTTP calls. This also allows us to tweak the
// transport values to be more resilient without making changes to the
// http.DefaultClient.
//
// Keep this unexported so consumers of the package can't make changes to it.
var defaultHTTPClient HTTPClient = newDefaultHTTPClient()

// Client wraps http client
type Client struct {
	authToken   string
	apiEndpoint string
	// HTTPClient is the HTTP client used for making requests against the
	// PagerDuty API. You can use either *http.Client here, or your own
	// implementation.
	HTTPClient HTTPClient
}

// DeleteResource deletes the given Resource. The given Resource should return a valid API URL from GetSelf()
func (c *Client) DeleteResource(resource Resource) error {
	res, err := c.delete(resource.GetSelf())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("received non-OK response %d: %s", res.StatusCode, body)
	}
	return nil
}

// GetResource fetches a Resource for a given type and ID
func (c *Client) GetResource(typ APIResourceType, id string, opts ...ResourceRequestOptionFunc) (Resource, error) {
	path := fmt.Sprintf("/%s/%s", typ.Plural(), id)
	res, err := c.get(path, opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	apiRes, err := APIResponses.Get(typ, res)
	if err != nil {
		return nil, err
	}
	return apiRes.GetResource()
}

func (c *Client) CreateResource(resource Resource) (Resource, error) {
	path := fmt.Sprintf("/%s", resource.GetType().Plural())
	res, err := c.post(path, map[APIResourceType]Resource{resource.GetType(): resource})
	if err != nil {
		return nil, err
	}
	apiRes, err := APIResponses.Get(resource.GetType(), res)
	if err != nil {
		return nil, err
	}
	return apiRes.GetResource()
}

func (c *Client) ListResources(typ APIResourceType, opts ...ResourceRequestOptionFunc) (*http.Response, error) {
	path := fmt.Sprintf("/%s", typ.Plural())
	res, err := c.get(path, opts...)
	if err != nil {
		return nil, err
	}

	//apiResp, err := APIListResponses.Get(typ, res)
	return res, nil
}

func (c *Client) setDefaultHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token token="+c.authToken)
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.do(http.MethodDelete, path, nil, nil)
}

func (c *Client) put(path string, payload interface{}, opts ...ResourceRequestOptionFunc) (*http.Response, error) {
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		return c.do(http.MethodPut, path, bytes.NewBuffer(data), opts...)
	}
	return c.do(http.MethodPut, path, nil, opts...)
}

func (c *Client) post(path string, payload interface{}, opts ...ResourceRequestOptionFunc) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.do(http.MethodPost, path, bytes.NewBuffer(data), nil)
}

func (c *Client) get(path string, opts ...ResourceRequestOptionFunc) (*http.Response, error) {
	return c.do(http.MethodGet, path, nil, opts...)
}

func (c *Client) do(method, path string, body io.Reader, opts ...ResourceRequestOptionFunc) (*http.Response, error) {
	endpoint := c.apiEndpoint + path
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	for _, opt := range opts {
		err := opt(req)
		if err != nil {
			return nil, err
		}
	}
	c.setDefaultHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	return c.checkResponse(resp, err)
}

func (c *Client) decodeJSON(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}

func (c *Client) checkResponse(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, fmt.Errorf("error calling the API endpoint: %v", err)
	}
	if resp.StatusCode <= 199 || resp.StatusCode >= http.StatusMultipleChoices {
		var eo *ErrorObject
		var getErr error
		if eo, getErr = c.getErrorFromResponse(resp); getErr != nil {
			return resp, fmt.Errorf("response did not contain formatted error: %s. HTTP response code: %v. Raw response: %+v", getErr, resp.StatusCode, resp)
		}
		return resp, fmt.Errorf("failed call API endpoint. HTTP response code: %v. Error: %v", resp.StatusCode, eo)
	}
	return resp, nil
}

func (c *Client) getErrorFromResponse(resp *http.Response) (*ErrorObject, error) {
	var result map[string]ErrorObject
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %v", err)
	}
	s, ok := result["error"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have error field")
	}
	return &s, nil
}

type NewClientOptionFunc func(*Client)

func WithCustomClient(c HTTPClient) NewClientOptionFunc {
	return func(client *Client) {
		client.HTTPClient = c
	}
}

func WithCustomHost(host string) NewClientOptionFunc {
	return func(client *Client) {
		client.apiEndpoint = host
	}
}

// NewClient creates an API client
func NewClient(authToken string, opts ...NewClientOptionFunc) *Client {
	c := &Client{
		authToken:   authToken,
		apiEndpoint: apiEndpoint,
		HTTPClient:  defaultHTTPClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
