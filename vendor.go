package pagerduty

import (
	"net/http"
)

// Vendor represents a specific type of integration. AWS Cloudwatch, Splunk, Datadog,
// etc are all examples of vendors that can be integrated in PagerDuty by making an integration.
type Vendor struct {
	APIObject
	Name                string `json:"name,omitempty"`
	LogoURL             string `json:"logo_url,omitempty"`
	LongName            string `json:"long_name,omitempty"`
	WebsiteURL          string `json:"website_url,omitempty"`
	Description         string `json:"description,omitempty"`
	Connectable         bool   `json:"connectable,omitempty"`
	ThumbnailURL        string `json:"thumbnail_url,omitempty"`
	GenericServiceType  string `json:"generic_service_type,omitempty"`
	IntegrationGuideURL string `json:"integration_guide_url,omitempty"`
}

type VendorResponse struct {
	APIResponse
}

func (r VendorResponse) GetResource() (Resource, error) {
	var dest Vendor
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewVendorResponse(resp *http.Response) VendorResponse {
	return VendorResponse{APIResponse{raw: resp, apiType: VendorResourceType}}
}

// ListVendorResponse is the data structure returned from calling the ListVendors API endpoint.
type ListVendorResponse struct {
	APIListObject
	Vendors []Vendor
}

// ListVendorOptions is the data structure used when calling the ListVendors API endpoint.
type ListVendorOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

// ListVendors lists existing vendors.
func (c *Client) ListVendors(opts ...ResourceRequestOptionFunc) (*ListVendorResponse, error) {
	resp, err := c.ListResources(VendorResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListVendorResponse
	return &result, deserialize(resp, &result)
}

// GetVendor gets details about an existing vendor.
func (c *Client) GetVendor(id string, opts ...ResourceRequestOptionFunc) (*Vendor, error) {
	res, err := c.GetResource(VendorResourceType, id, opts...)
	if err != nil {
	    return nil, err
	}
	obj := res.(Vendor)
	return &obj, nil
}
