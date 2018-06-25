package pagerduty

import (
	"encoding/json"
	"io"
	"net/http"
)

// IncidentDetail contains a representation of the incident associated with the action that caused this webhook message.
type IncidentDetail struct {
	ID                    string           `json:"id"`
	IncidentNumber        uint             `json:"incident_number"`
	CreatedOn             string           `json:"created_on"`
	Status                string           `json:"status"`
	HTMLUrl               string           `json:"html_url"`
	Service               string           `json:"service"`
	AssignedToUser        *json.RawMessage `json:"assigned_to_user"`
	AssignedTo            []string         `json:"assigned_to"`
	TriggerSummaryData    *json.RawMessage `json:"trigger_summary_data"`
	TriggerDetailsHTMLUrl string           `json:"trigger_details_html_url"`
}

// WebhookPayload is a single message array for a webhook.
type WebhookPayload struct {
	ID        string           `json:"id"`
	Type      string           `json:"type"`
	CreatedOn string           `json:"created_on"`
	Data      *json.RawMessage `json:"data"`
}

// DecodeWebhook decodes a webhook from a response object.
func DecodeWebhook(r io.Reader) (*WebhookPayload, error) {
	var payload WebhookPayload
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

type Webhook struct {
	APIObject
}

type WebhookResponse struct {
	APIResponse
}

func (r WebhookResponse) GetResource() (Resource, error) {
	var dest Webhook
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewWebhookResponse(resp *http.Response) WebhookResponse {
	return WebhookResponse{APIResponse{raw: resp, apiType: WebhookResourceType}}
}

type ListWebhooksResponse struct {
	APIListObject
}
