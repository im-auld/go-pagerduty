package pagerduty

import (
	"net/http"
)

// Notification is a message containing the details of the incident.
type Notification struct {
	ID        string `json:"id"`
	Type      string
	StartedAt string `json:"started_at"`
	Address   string
	User      APIObject
}

func (n Notification) GetID() string {
	return n.ID
}

func (n Notification) GetType() APIResourceType {
	return NotificationResourceType
}

func (n Notification) GetSummary() string {
	return n.ID
}

func (n Notification) GetSelf() string {
	return ""
}

func (n Notification) GetHTMLURL() string {
	return ""
}

type NotificationResponse struct {
	APIResponse
}

func (r NotificationResponse) GetResource() (Resource, error) {
	var dest Notification
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewNotificationResponse(resp *http.Response) NotificationResponse {
	return NotificationResponse{APIResponse{raw: resp, apiType: NotificationResourceType}}
}

// ListNotificationOptions is the data structure used when calling the ListNotifications API endpoint.
type ListNotificationOptions struct {
	APIListObject
	TimeZone string   `url:"time_zone,omitempty"`
	Since    string   `url:"since,omitempty"`
	Until    string   `url:"until,omitempty"`
	Filter   string   `url:"filter,omitempty"`
	Includes []string `url:"include,omitempty"`
}

// ListNotificationsResponse is the data structure returned from the ListNotifications API endpoint.
type ListNotificationsResponse struct {
	APIListObject
	Notifications []Notification
}

// ListNotifications lists notifications for a given time range, optionally filtered by type (sms_notification, email_notification, phone_notification, or push_notification).
func (c *Client) ListNotifications(opts ...ResourceRequestOptionFunc) (*ListNotificationsResponse, error) {
	resp, err := c.ListResources(NotificationResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListNotificationsResponse
	return &result, deserialize(resp, &result)
}
