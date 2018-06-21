package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"time"
)

type APIResourceType string

const (
	apiEndpoint = "https://api.pagerduty.com"

	// Resource Types
	AbilityResourceType           APIResourceType = "ability"
	AddonResourceType             APIResourceType = "addon"
	EscalationPolicyResourceType  APIResourceType = "escalation_policy"
	EventResourceType             APIResourceType = "event"
	IncidentResourceType          APIResourceType = "incident"
	LogEntryResourceType          APIResourceType = "log_entry"
	MaintenanceWindowResourceType APIResourceType = "maintenance_window"
	NotificationResourceType      APIResourceType = "notification"
	OnCallResourceType            APIResourceType = "on_call"
	ResponsePLayResourceType      APIResourceType = "response_play"
	ScheduleResourceType          APIResourceType = "schedule"
	ServiceResourceType           APIResourceType = "service"
	TeamResourceType              APIResourceType = "team"
	UserResourceType              APIResourceType = "user"
	VendorResourceType            APIResourceType = "vendor"
	WebhookResourceType           APIResourceType = "webhook"
)

type Resource interface {
	GetID() string
	GetType() APIResourceType
	GetSummary() string
	GetSelf() string
	GetHTMLURL() string
}

type ResourceList interface {
	GetLimit() uint
	GetOffset() uint
	GetMore() bool
	GetTotal() uint
}

// APIObject represents generic api json response that is shared by most
// domain object (like escalation
type APIObject struct {
	ID      string          `json:"id,omitempty"`
	Type    APIResourceType `json:"type,omitempty"`
	Summary string          `json:"summary,omitempty"`
	Self    string          `json:"self,omitempty"`
	HTMLURL string          `json:"html_url,omitempty"`
}

func (apiObj APIObject) GetID() string {
	return apiObj.ID
}

func (apiObj APIObject) GetType() APIResourceType {
	return apiObj.Type
}

func (apiObj APIObject) GetSummary() string {
	return apiObj.Summary
}

func (apiObj APIObject) GetSelf() string {
	return apiObj.Self
}

func (apiObj APIObject) GetHTMLURL() string {
	return apiObj.HTMLURL
}

// APIListObject are the fields used to control pagination when listing objects.
type APIListObject struct {
	Limit  uint `url:"limit,omitempty"`
	Offset uint `url:"offset,omitempty"`
	More   bool `url:"more,omitempty"`
	Total  uint `url:"total,omitempty"`
}

func (list APIListObject) GetLimit() uint {
	return list.Limit
}

func (list APIListObject) GetOffset() uint {
	return list.Offset
}

func (list APIListObject) GetMore() bool {
	return list.More
}

func (list APIListObject) GetTotal() uint {
	return list.Total
}

// APIReference are the fields required to reference another API object.
type APIReference struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type ErrorObject struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

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
	authToken string

	// HTTPClient is the HTTP client used for making requests against the
	// PagerDuty API. You can use either *http.Client here, or your own
	// implementation.
	HTTPClient HTTPClient
}

type NewClientOptionFunc func(*Client)

func WithCustomClient(c HTTPClient) NewClientOptionFunc {
	return func(client *Client) {
		client.HTTPClient = c
	}
}

// NewClient creates an API client
func NewClient(authToken string, opts ...NewClientOptionFunc) *Client {
	c := &Client{
		authToken:  authToken,
		HTTPClient: defaultHTTPClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) setDefaultHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token token="+c.authToken)
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.do(http.MethodDelete, path, nil, nil)
}

func (c *Client) put(path string, payload interface{}, headers *map[string]string) (*http.Response, error) {

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		return c.do(http.MethodPut, path, bytes.NewBuffer(data), headers)
	}
	return c.do(http.MethodPut, path, nil, headers)
}

func (c *Client) post(path string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.do(http.MethodPost, path, bytes.NewBuffer(data), nil)
}

func (c *Client) get(path string) (*http.Response, error) {
	return c.do(http.MethodGet, path, nil, nil)
}

func (c *Client) do(method, path string, body io.Reader, headers *map[string]string) (*http.Response, error) {
	endpoint := apiEndpoint + path
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
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
	if resp.StatusCode <= 199  || resp.StatusCode >= http.StatusMultipleChoices {
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
