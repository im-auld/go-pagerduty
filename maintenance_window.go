package pagerduty

import (
	"net/http"
)

// MaintenanceWindow is used to temporarily disable one or more services for a set period of time.
type MaintenanceWindow struct {
	APIObject
	SequenceNumber uint            `json:"sequence_number,omitempty"`
	StartTime      string          `json:"start_time"`
	EndTime        string          `json:"end_time"`
	Description    string          `json:"description"`
	Services       []APIObject     `json:"services"`
	Teams          []APIListObject `json:"teams"`
	CreatedBy      APIListObject   `json:"created_by"`
}

type MaintenanceWindowResponse struct {
	APIResponse
}

func (r MaintenanceWindowResponse) GetResource() (Resource, error) {
	var dest MaintenanceWindow
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewMaintenanceWindowResponse(resp *http.Response) MaintenanceWindowResponse {
	return MaintenanceWindowResponse{APIResponse{raw: resp, apiType: MaintenanceWindowResourceType}}
}

// ListMaintenanceWindowsResponse is the data structur returned from calling the ListMaintenanceWindows API endpoint.
type ListMaintenanceWindowsResponse struct {
	APIListObject
	MaintenanceWindows []MaintenanceWindow `json:"maintenance_windows"`
}

// ListMaintenanceWindowsOptions is the data structure used when calling the ListMaintenanceWindows API endpoint.
type ListMaintenanceWindowsOptions struct {
	APIListObject
	Query      string   `url:"query,omitempty"`
	Includes   []string `url:"include,omitempty,brackets"`
	TeamIDs    []string `url:"team_ids,omitempty,brackets"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
	Filter     string   `url:"filter,omitempty,brackets"`
}

// ListMaintenanceWindows lists existing maintenance windows, optionally filtered by service and/or team, or whether they are from the past, present or future.
func (c *Client) ListMaintenanceWindows(opts ...ResourceRequestOptionFunc) (*ListMaintenanceWindowsResponse, error) {
	resp, err := c.ListResources(MaintenanceWindowResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListMaintenanceWindowsResponse
	return &result, deserialize(resp, &result)
}

// CreateMaintenanceWindows creates a new maintenance window for the specified services.
func (c *Client) CreateMaintenanceWindows(m MaintenanceWindow) (*MaintenanceWindow, error) {
	resp, err := c.CreateResource(m)
	if err != nil {
		return nil, err
	}
	mw := resp.(MaintenanceWindow)
	return &mw, nil
}

// DeleteMaintenanceWindow deletes an existing maintenance window if it's in the future, or ends it if it's currently on-going.
func (c *Client) DeleteMaintenanceWindow(id string) error {
	return c.DeleteResource(MaintenanceWindowResourceType, id)
}

// GetMaintenanceWindowOptions is the data structure used when calling the GetMaintenanceWindow API endpoint.
type GetMaintenanceWindowOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// GetMaintenanceWindow gets an existing maintenance window.
func (c *Client) GetMaintenanceWindow(id string, opts ...ResourceRequestOptionFunc) (*MaintenanceWindow, error) {
	res, err := c.GetResource(MaintenanceWindowResourceType, id, opts...)
	if err != nil {
	    return nil, err
	}
	obj := res.(MaintenanceWindow)
	return &obj, nil
}

// UpdateMaintenanceWindow updates an existing maintenance window.
func (c *Client) UpdateMaintenanceWindow(m MaintenanceWindow) (*MaintenanceWindow, error) {
	resp, err := c.UpdateResource(m)
	if err != nil {
		return nil, err
	}
	mw := resp.(MaintenanceWindow)
	return &mw, nil
}
