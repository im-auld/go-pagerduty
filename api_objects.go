package pagerduty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	HasMore() bool
	GetTotal() uint
	//GetResources() []Resource
}

type Response interface {
	GetResource() (Resource, error)
}

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

type ResponseTypeFunc func(response *http.Response) Response
type apiResourceTypes map[APIResourceType]ResponseTypeFunc

func (rt apiResourceTypes) Get(typ APIResourceType, resp *http.Response) (Response, error) {
	r, ok := rt[typ]
	if !ok {
		return nil, NewInvalidResourceTypeError(typ)
	}
	return r(resp), nil
}

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

type ListResponseTypeFunc func(response *http.Response) ResourceList
type apiListResourceTypes map[APIResourceType]ListResponseTypeFunc

func (rt apiListResourceTypes) Get(typ APIResourceType, resp *http.Response) (ResourceList, error) {
	r, ok := rt[typ]
	if !ok {
		return nil, NewInvalidResourceTypeError(typ)
	}
	return r(resp), nil
}

var APIListResponses = apiListResourceTypes{
	AbilityResourceType:           func(response *http.Response) ResourceList { return new(ListAbilityResponse) },
	AddonResourceType:             func(response *http.Response) ResourceList { return new(ListAddonResponse) },
	EscalationPolicyResourceType:  func(response *http.Response) ResourceList { return new(ListEscalationPoliciesResponse) },
	IncidentResourceType:          func(response *http.Response) ResourceList { return new(ListIncidentsResponse) },
	LogEntryResourceType:          func(response *http.Response) ResourceList { return new(ListLogEntryResponse) },
	MaintenanceWindowResourceType: func(response *http.Response) ResourceList { return new(ListMaintenanceWindowsResponse) },
	NotificationResourceType:      func(response *http.Response) ResourceList { return new(ListNotificationsResponse) },
	ResponsePlayResourceType:      func(response *http.Response) ResourceList { return new(ListResponsePlaysResponse) },
	ScheduleResourceType:          func(response *http.Response) ResourceList { return new(ListSchedulesResponse) },
	ServiceResourceType:           func(response *http.Response) ResourceList { return new(ListServiceResponse) },
	TeamResourceType:              func(response *http.Response) ResourceList { return new(ListTeamResponse) },
	UserResourceType:              func(response *http.Response) ResourceList { return new(ListUsersResponse) },
	VendorResourceType:            func(response *http.Response) ResourceList { return new(ListVendorResponse) },
	WebhookResourceType:           func(response *http.Response) ResourceList { return new(ListWebhooksResponse) },
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
