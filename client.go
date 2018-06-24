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

type APIResourceType string

var vowels = map[rune]bool{
	'a': true,
	'e': true,
	'i': true,
	'o': true,
	'u': true,
}

//IV. Nouns ending in "y"
//a. If the common noun ends with a consonant + "y" or "qu" + "y" , remove the "y" and add "ies".
//The vowels are the letters a, e, i, o, and u. All other letters are consonants.
//b. Common nouns with a vowel + "y", just add "s"
//Exception: To form the plural of proper nouns ending in "y" preceded by a consonant, just add an "s".
func (r APIResourceType) Plural() APIResourceType {
	l := len(r)
	penultimate, ultimate := r[l-2], r[l-1]
	if ultimate == 'y' {
		if _, ok := vowels[rune(penultimate)]; !ok {
			return APIResourceType(r.String()[:l-1] + "ies")
		}
	}
	return APIResourceType(r.String() + "s")
}

func (r APIResourceType) String() string {
	return string(r)
}

const (
	apiEndpoint = "https://api.pagerduty.com"

	// Resource Types
	AbilityResourceType           APIResourceType = "ability"
	AddonResourceType             APIResourceType = "addon"
	EscalationPolicyResourceType  APIResourceType = "escalation_policy"
	EventResourceType             APIResourceType = "event"
	ExtensionResourceType         APIResourceType = "extension"
	IncidentResourceType          APIResourceType = "incident"
	LogEntryResourceType          APIResourceType = "log_entry"
	MaintenanceWindowResourceType APIResourceType = "maintenance_window"
	NotificationResourceType      APIResourceType = "notification"
	OnCallResourceType            APIResourceType = "on_call"
	ResponsePlayResourceType      APIResourceType = "response_play"
	ScheduleResourceType          APIResourceType = "schedule"
	ServiceResourceType           APIResourceType = "service"
	TeamResourceType              APIResourceType = "team"
	UserResourceType              APIResourceType = "user"
	VendorResourceType            APIResourceType = "vendor"
	WebhookResourceType           APIResourceType = "webhook"
)

type ResourceTypeFunc func() Resource
type ResponseTypeFunc func(response *http.Response) Response
type apiResourceTypes map[APIResourceType]ResponseTypeFunc

func (rt apiResourceTypes) Get(typ APIResourceType, resp *http.Response) (Response, error) {
	r, ok := rt[typ]
	if !ok {
		return nil, NewInvalidResourceTypeError(typ)
	}
	return r(resp), nil
}

//var APIResources = apiResourceTypes{
//	AbilityResourceType:           func() Resource { return new(Ability) },
//	AddonResourceType:             func() Resource { return new(Addon) },
//	EscalationPolicyResourceType:  func() Resource { return new(EscalationPolicy) },
//	IncidentResourceType:          func() Resource { return new(Incident) },
//	LogEntryResourceType:          func() Resource { return new(LogEntry) },
//	MaintenanceWindowResourceType: func() Resource { return new(MaintenanceWindow) },
//	NotificationResourceType:      func() Resource { return new(Notification) },
//	ResponsePlayResourceType:      func() Resource { return new(ResponsePlay) },
//	ScheduleResourceType:          func() Resource { return new(Schedule) },
//	ServiceResourceType:           func() Resource { return new(Service) },
//	TeamResourceType:              func() Resource { return new(Team) },
//	UserResourceType:              func() Resource { return new(User) },
//	VendorResourceType:            func() Resource { return new(Vendor) },
//	WebhookResourceType:           func() Resource { return new(Webhook) },
//}

var APIResponses = apiResourceTypes{
	AbilityResourceType:           func(response *http.Response) Response { return NewAbilityResponse(response) },
	AddonResourceType:             func(response *http.Response) Response { return NewAddonResponse(response) },
	EscalationPolicyResourceType:  func(response *http.Response) Response { return NewEscalationPolicyResponse(response) },
	IncidentResourceType:          func(response *http.Response) Response { return NewIncidentResponse(response) },
	LogEntryResourceType:          func(response *http.Response) Response { return NewLogEntryResponse(response) },
	MaintenanceWindowResourceType: func(response *http.Response) Response { return NewMaintenanceWindowResponse(response) },
	NotificationResourceType:      func(response *http.Response) Response { return NewNotificationResponse(response) },
	ResponsePlayResourceType:      func(response *http.Response) Response { return NewResponsePlayResponse(response) },
	ScheduleResourceType:          func(response *http.Response) Response { return NewScheduleResponse(response) },
	ServiceResourceType:           func(response *http.Response) Response { return NewServiceResponse(response) },
	TeamResourceType:              func(response *http.Response) Response { return NewTeamResponse(response) },
	UserResourceType:              func(response *http.Response) Response { return NewUserResponse(response) },
	VendorResourceType:            func(response *http.Response) Response { return NewVendorResponse(response) },
	WebhookResourceType:           func(response *http.Response) Response { return NewWebhookResponse(response) },
}

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
	HasMore() bool
	GetTotal() uint
}

type Response interface {
	GetResource() (Resource, error)
}

type APIResponse struct {
	raw     *http.Response
	apiType APIResourceType
}

func (ar APIResponse) getResourceFromResponse(target Resource) error {
	var dest map[string]json.RawMessage
	body, err := ioutil.ReadAll(ar.raw.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("could not decode JSON response: %v", err)
	}
	t, nodeOK := dest[string(ar.apiType)]
	if !nodeOK {
		return fmt.Errorf("JSON response does not have %s field", ar.apiType)
	}
	if err := json.Unmarshal(t, target); err != nil {
		return nil
	}
	return nil
}

func (ar APIResponse) GetResource(tgt Resource) error {
	return ar.getResourceFromResponse(tgt)
}

func NewAPIResponse(res *http.Response, typ APIResourceType) APIResponse {
	return APIResponse{raw: res, apiType: typ}
}

type ListResponse interface {
	GetResources() []Resource
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

func (apiObj APIObject) String() string {
	return fmt.Sprintf("<%s %s: %s>", apiObj.GetType(), apiObj.GetID(), apiObj.GetSummary())
}

// APIListObject are the fields used to control pagination when listing objects.
type APIListObject struct {
	Limit  uint `url:"limit,omitempty" json:"limit"`
	Offset uint `url:"offset,omitempty" json:"offset"`
	More   bool `url:"more,omitempty" json:"more"`
	Total  uint `url:"total,omitempty" json:"total"`
}

func (list APIListObject) GetLimit() uint {
	return list.Limit
}

func (list APIListObject) GetOffset() uint {
	return list.Offset
}

func (list APIListObject) HasMore() bool {
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
func (c *Client) GetResource(typ APIResourceType, id string) (Resource, error) {
	path := fmt.Sprintf("/%s/%s", typ.Plural(), id)
	res, err := c.get(path)
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
	res, err := c.post(path, resource)
	if err != nil {
		return nil, err
	}
	apiRes, err := APIResponses.Get(resource.GetType(), res)
	if err != nil {
		return nil, err
	}
	return apiRes.GetResource()
}

func (c *Client) ListResources(typ APIResourceType) (ResourceList, error) {

}
