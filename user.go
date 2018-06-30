package pagerduty

import (
	"net/http"
)

// ContactMethod is a way of contacting the user.
type ContactMethod struct {
	ID             string
	Label          string
	Address        string
	Type           string
	SendShortEmail bool `json:"send_short_email"`
}

// NotificationRule is a rule for notifying the user.
type NotificationRule struct {
	ID                  string
	StartDelayInMinutes uint          `json:"start_delay_in_minutes"`
	CreatedAt           string        `json:"created_at"`
	ContactMethod       ContactMethod `json:"contact_method"`
	Urgency             string
	Type                string
}

// User is a member of a PagerDuty account that has the ability to interact with incidents and other data on the account.
type User struct {
	APIObject
	Name              string `json:"name"`
	Email             string `json:"email"`
	Timezone          string `json:"timezone,omitempty"`
	Color             string `json:"color,omitempty"`
	Role              string `json:"role,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	Description       string `json:"description,omitempty"`
	InvitationSent    bool
	ContactMethods    []ContactMethod    `json:"contact_methods"`
	NotificationRules []NotificationRule `json:"notification_rules"`
	JobTitle          string             `json:"job_title,omitempty"`
	Teams             []Team
}

type UserResponse struct {
	APIResponse
}

func (r UserResponse) GetResource() (Resource, error) {
	var dest User
	err := r.getResourceFromResponse(&dest)
	return dest, err
}

func NewUserResponse(resp *http.Response) UserResponse {
	return UserResponse{APIResponse{raw: resp, apiType: UserResourceType}}
}

// ListUsersResponse is the data structure returned from calling the ListUsers API endpoint.
type ListUsersResponse struct {
	APIListObject
	Users []User
}

// ListUsersOptions is the data structure used when calling the ListUsers API endpoint.
type ListUsersOptions struct {
	APIListObject
	Query    string   `url:"query,omitempty"`
	TeamIDs  []string `url:"team_ids,omitempty,brackets"`
	Includes []string `url:"include,omitempty,brackets"`
}

// GetUserOptions is the data structure used when calling the GetUser API endpoint.
type GetUserOptions struct {
	Includes []string `url:"include,omitempty,brackets"`
}

// ListUsers lists users of your PagerDuty account, optionally filtered by a search query.
func (c *Client) ListUsers(opts ...ResourceRequestOptionFunc) (*ListUsersResponse, error) {
	resp, err := c.ListResources(UserResourceType, opts...)
	if err != nil {
		return nil, err
	}
	var result ListUsersResponse
	return &result, deserialize(resp, &result)
}

// CreateUser creates a new user.
func (c *Client) CreateUser(u User) (*User, error) {
	resp, err := c.CreateResource(u)
	if err != nil {
		return nil, err
	}
	user := resp.(User)
	return &user, nil
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(id string) error {
	return c.DeleteResource(UserResourceType, id)
}

// GetUser gets details about an existing user.
func (c *Client) GetUser(id string, opts ...ResourceRequestOptionFunc) (*User, error) {
	res, err := c.GetResource(UserResourceType, id, opts...)
	if err != nil {
	    return nil, err
	}
	obj := res.(User)
	return &obj, nil
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(u User) (*User, error) {
	resp, err := c.UpdateResource(u)
	if err != nil {
		return nil, err
	}
	user := resp.(User)
	return &user, nil
}
