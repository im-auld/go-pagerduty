package pagerduty

import "net/http"

type Ability string

func (a Ability) GetID() string {
	return string(a)
}

func (a Ability) GetType() APIResourceType {
	return AbilityResourceType
}

func (a Ability) GetSummary() string {
	return string(a)
}

func (a Ability) GetSelf() string {
	return "https://api.pagerduty.com/abilities/" + string(a)
}

func (a Ability) GetHTMLURL() string {
	return ""
}

func (a Ability) String() string {
	return string(a)
}

type AbilityResponse struct {
	APIResponse
}

func (r AbilityResponse) GetResource() (Resource, error) {
	var ability Ability
	err := r.getResourceFromResponse(&ability)
	return ability, err
}

func NewAbilityResponse(resp *http.Response) AbilityResponse {
	return AbilityResponse{APIResponse{raw: resp, apiType: AbilityResourceType}}
}

type Abilities []Ability

func (a Abilities) Len() int           { return len(a) }
func (a Abilities) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Abilities) Less(i, j int) bool { return a[i].GetID() < a[j].GetID() }

// ListAbilityResponse is the response when calling the ListAbility API endpoint.
type ListAbilityResponse struct {
	Abilities Abilities `json:"abilities"`
}

func (list ListAbilityResponse) GetLimit() uint {
	return 27
}

func (list ListAbilityResponse) GetOffset() uint {
	return 0
}

func (list ListAbilityResponse) HasMore() bool {
	return false
}

func (list ListAbilityResponse) GetTotal() uint {
	return uint(len(list.Abilities))
}

// ListAbilities lists all abilities on your account.
func (c *Client) ListAbilities() (*ListAbilityResponse, error) {
	resp, err := c.ListResources(AbilityResourceType)
	if err != nil {
		return nil, err
	}
	var result ListAbilityResponse
	return &result, c.decodeJSON(resp, &result)
}

// TestAbility Check if your account has the given ability.
func (c *Client) TestAbility(ability Ability) error {
	_, err := c.get("/abilities/" + ability.GetID())
	return err
}
