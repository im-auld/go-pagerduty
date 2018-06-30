package pagerduty

import (
	"net/http"
)

// Team is a collection of users and escalation policies that represent a group of people within an organization.
type Team struct {
	APIObject
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type TeamResponse struct {
	APIResponse
}

func (r TeamResponse) GetResource() (Resource, error) {
	var dest Team
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewTeamResponse(resp *http.Response) TeamResponse {
	return TeamResponse{APIResponse{raw: resp, apiType: TeamResourceType}}
}

// ListTeamResponse is the structure used when calling the ListTeams API endpoint.
type ListTeamResponse struct {
	APIListObject
	Teams []Team
}

// ListTeamOptions are the input parameters used when calling the ListTeams API endpoint.
type ListTeamOptions struct {
	APIListObject
	Query string `url:"query,omitempty"`
}

// ListTeams lists teams of your PagerDuty account, optionally filtered by a search query.
func (c *Client) ListTeams(opts ...ResourceRequestOptionFunc) (*ListTeamResponse, error) {
	resp, err := c.ListResources(TeamResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListTeamResponse
	return &result, deserialize(resp, &result)
}

// CreateTeam creates a new team.
func (c *Client) CreateTeam(t *Team) (*Team, error) {
	resp, err := c.CreateResource(t)
	if err != nil {
		return nil, err
	}
	team := resp.(Team)
	return &team, nil
}

// DeleteTeam removes an existing team.
func (c *Client) DeleteTeam(id string) error {
	return c.DeleteResource(TeamResourceType, id)
}

// GetTeam gets details about an existing team.
func (c *Client) GetTeam(id string) (*Team, error) {
	res, err := c.GetResource(TeamResourceType, id)
	if err != nil {
	    return nil, err
	}
	obj := res.(Team)
	return &obj, nil
}

// UpdateTeam updates an existing team.
func (c *Client) UpdateTeam(id string, t *Team) (*Team, error) {
	resp, err := c.UpdateResource(t)
	if err != nil {
		return nil, err
	}
	team := resp.(Team)
	return &team, nil
}

// RemoveEscalationPolicyFromTeam removes an escalation policy from a team.
func (c *Client) RemoveEscalationPolicyFromTeam(teamID, epID string) error {
	_, err := c.delete("/teams/" + teamID + "/escalation_policies/" + epID)
	return err
}

// AddEscalationPolicyToTeam adds an escalation policy to a team.
func (c *Client) AddEscalationPolicyToTeam(teamID, epID string) error {
	_, err := c.put("/teams/"+teamID+"/escalation_policies/"+epID, nil, nil)
	return err
}

// RemoveUserFromTeam removes a user from a team.
func (c *Client) RemoveUserFromTeam(teamID, userID string) error {
	_, err := c.delete("/teams/" + teamID + "/users/" + userID)
	return err
}

// AddUserToTeam adds a user to a team.
func (c *Client) AddUserToTeam(teamID, userID string) error {
	_, err := c.put("/teams/"+teamID+"/users/"+userID, nil, nil)
	return err
}
