package pagerduty

import "time"

type AlertBody struct {
	Type     string                 `json:"type"`
	Contexts []APIReference           `json:"contexts"`
	Details  map[string]interface{} `json:"details"`
}

type Alert struct {
	APIReference
	CreatedAt   time.Time    `json:"created_at"`
	Status      string       `json:"status"`
	AlertKey    string       `json:"alert_key"`
	Service     APIReference `json:"service"`
	Body        AlertBody    `json:"body"`
	Incident    APIReference `json:"incident"`
	Suppressed  bool         `json:"suppressed"`
	Severity    string       `json:"severity"`
	Integration APIReference `json:"integration"`
}
