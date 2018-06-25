package pagerduty

import "fmt"

// InvalidResourceTypeError is a custom error type.
type InvalidResourceTypeError struct {
	Message string
}

func (e InvalidResourceTypeError) Error() string {
	return e.Message
}

// NewInvalidResourceTypeError creates a new `InvalidResourceTypeError`.
func NewInvalidResourceTypeError(typ APIResourceType) InvalidResourceTypeError {
	msg := fmt.Sprintf("%s is not a known resource type.", typ)
	return InvalidResourceTypeError{Message: msg}
}


var ErrorCode_Message = map[int]string{
	1001: "Incident Already Resolved",
	1002: "Incident Already Acknowledged",
	1003: "Invalid Status Provided",
	1004: "Invalid Id Provided",
	1005: "Data Updated Since Last Request",
	1006: "Cannot Escalate",
	1007: "Assigned To User Not Found",
	1008: "Requester User Not Found",
	1009: "Error Parsing before_time Parameter",
	1010: "Before Time or Before Incident Number Not Found",
	1011: "Cannot process reassign action",
	1012: "Incident Is Not Acknowledged",
	1013: "Snooze Duration is Invalid",
	1014: "Escalation Policy Not Found",
	2000: "An internal error has occurred.PagerDuty administrators have been notified.",
	2001: "Cannot Delete Only Rule For Escalation Policy, Cannot create more than 50 users at one time, Escalation Policy has already been taken, Incident could not be created, Invalid Input Provided, Users must be an array",
	2002: "Arguments Caused Error",
	2003: "Missing Arguments",
	2004: "Invalid time values in 'since' and/or 'until' parameters",
	2005: "Invalid Query Date Range",
	2006: "Authentication failed",
	2007: "Account Not Found",
	2008: "Account Locked",
	2009: "Only HTTPS Allowed For This Call",
	2010: "Access Denied",
	2011: "You must specify a requester_id to perform this action",
	2012: "Your account is expired and cannot use the API.",
	2013: "User not found",
	2015: "Invalid operation",
	2016: "The request took too long to process",
	2100: "Not Found",
	3001: "Invalid Schedule",
	3004: "Schedule Not Found, Unknown Schedule",
	4004: "Invalid Override, Override Not Found",
	5001: "Affiliate Partner Not Found, Color Not Found, Escalation Policy Not Found, Escalation Rule Not Found, Incident Not Found, Log Entry Not Found, Maintenance Window Not Found, User Next On-call Not Found, Vendor Not Found, Webhook Not Found",
	5002: "Service Not Found, Simple Log Entry Limit Reached",
	5003: "Cannot cancel a maintenance window from the past",
	6001: "Team Not Found",
	6002: "Team has Incidents but Reassignment Team is not valid",
	6003: "Team has Incidents but Reassignment Team is not found",
	6004: "Team could not be deleted because there are too many incidents, please contact support.",
	6006: "Team has existing associations",
	6007: "The request failed to complete",
	7001: "saml is not enabled",
	7002: "SAMLResponse is not provided",
	7003: "failed to parse SAMLResponse",
	7004: "SAMLResponse is invalid",
	7005: "unable to find user by name_id",
	7006: "invalid name_id found in SAMLResponse",
	7007: "unexpected validation error",
	8001: "google auth is not enabled, pre_auth.authorizable? returned false",
	8002: "redirect_uri is invalid",
	8004: "unable to find account",
	8005: "unable to find user by name_id",
	8006: "user has denied google auth",
	8007: "unable to log in using google auth",
	8008: "user has tried to access an invalid host domain",
}

// APIError is a custom error type.
type APIError struct {
	Message string
	Code int
}

func (e APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new `APIError`.
func NewAPIError(code int) APIError {
	return APIError{Code: code, Message: ErrorCode_Message[code]}
}
