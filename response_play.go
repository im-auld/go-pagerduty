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

type ListResponsePlaysResponse struct {
	APIListObject
}

func (c *Client) GetResponsePlay(id string, opts ...ResourceRequestOptionFunc) (*ResponsePlay, error) {
	res, err := c.GetResource(ResponsePlayResourceType, id, opts...)
	if err != nil {
	    return nil, err
	}
	obj := res.(ResponsePlay)
	return &obj, nil
}
