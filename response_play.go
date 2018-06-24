package pagerduty

import "net/http"

type ResponsePlay struct {
	APIObject
}

type ResponsePlayResponse struct {
	APIResponse
}

func (r ResponsePlayResponse) GetResource() (Resource, error) {
	var dest ResponsePlay
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewResponsePlayResponse(resp *http.Response) ResponsePlayResponse {
	return ResponsePlayResponse{APIResponse{raw: resp, apiType: ResponsePlayResourceType}}
}
