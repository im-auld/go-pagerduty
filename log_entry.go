package pagerduty

import (
	"net/http"
)

// Agent is the actor who carried out the action.
type Agent APIObject

// Channel is the means by which the action was carried out.
type Channel struct {
	Type string
}

// Context are to be included with the trigger such as links to graphs or images.
type Context struct {
	Alt  string
	Href string
	Src  string
	Text string
	Type string
}

// LogEntry is a list of all of the events that happened to an incident.
type LogEntry struct {
	APIObject
	CreatedAt              string `json:"created_at"`
	Agent                  Agent
	Channel                Channel
	Incident               Incident
	Teams                  []Team
	Contexts               []Context
	AcknowledgementTimeout int `json:"acknowledgement_timeout"`
	EventDetails           map[string]string
}

type LogEntryResponse struct {
	APIResponse
}

func (r LogEntryResponse) GetResource() (Resource, error) {
	var dest LogEntry
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewLogEntryResponse(resp *http.Response) LogEntryResponse {
	return LogEntryResponse{APIResponse{raw: resp, apiType: LogEntryResourceType}}
}

// ListLogEntryResponse is the response data when calling the ListLogEntry API endpoint.
type ListLogEntryResponse struct {
	APIListObject
	LogEntries []LogEntry `json:"log_entries"`
}

// ListLogEntriesOptions is the data structure used when calling the ListLogEntry API endpoint.
type ListLogEntriesOptions struct {
	APIListObject
	TimeZone   string   `url:"time_zone"`
	Since      string   `url:"since,omitempty"`
	Until      string   `url:"until,omitempty"`
	IsOverview bool     `url:"is_overview,omitempty"`
	Includes   []string `url:"include,omitempty,brackets"`
}

// ListLogEntries lists all of the incident log entries across the entire account.
func (c *Client) ListLogEntries(opts ...ResourceRequestOptionFunc) (*ListLogEntryResponse, error) {
	resp, err := c.ListResources(LogEntryResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListLogEntryResponse
	return &result, deserialize(resp, &result)
}
