package pagerduty

import (
	"net/http"
)

// Addon is a third-party add-on to PagerDuty's UI.
type Addon struct {
	APIObject
	Name     string      `json:"name,omitempty"`
	Src      string      `json:"src,omitempty"`
	Services []APIObject `json:"services,omitempty"`
}

type AddonResponse struct {
	APIResponse
}

func (r AddonResponse) GetResource() (Resource, error) {
	var dest Addon
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewAddonResponse(resp *http.Response) AddonResponse {
	return AddonResponse{APIResponse{raw: resp, apiType: AddonResourceType}}
}

// ListAddonOptions are the options available when calling the ListAddons API endpoint.
type ListAddonOptions struct {
	APIListObject
	Includes   []string `url:"include,omitempty,brackets"`
	ServiceIDs []string `url:"service_ids,omitempty,brackets"`
	Filter     string   `url:"filter,omitempty"`
}

// ListAddonResponse is the response when calling the ListAddons API endpoint.
type ListAddonResponse struct {
	APIListObject
	Addons []Addon `json:"addons"`
}

// ListAddons lists all of the add-ons installed on your account.
func (c *Client) ListAddons(opts ...ResourceRequestOptionFunc) (*ListAddonResponse, error) {
	resp, err := c.ListResources(AddonResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListAddonResponse
	return &result, deserialize(resp, &result)
}

// GetAddon gets details about an existing add-on.
func (c *Client) GetAddon(id string) (*Addon, error) {
	res, err := c.GetResource(AddonResourceType, id)
	if err != nil {
	    return nil, err
	}
	obj := res.(Addon)
	return &obj, nil
}

// InstallAddon installs an add-on for your account.
func (c *Client) InstallAddon(a Addon) (*Addon, error) {
	resp, err := c.CreateResource(a)
	if err != nil {
		return nil, err
	}
	addon := resp.(Addon)
	return &addon, nil
}
