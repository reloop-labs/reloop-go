package reloop

import (
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	groupsBase      = "/api/contacts/v1/groups"
	groupMembership = "/api/contacts/group"
)

// ContactGroup is a contact group.
type ContactGroup struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ContactGroupResponse is returned by group mutations.
type ContactGroupResponse struct {
	ContactGroup
	Event string `json:"event"`
}

// ContactGroupListItem is a group in a list response.
type ContactGroupListItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateGroupParams creates a contact group.
type CreateGroupParams struct {
	Name string
}

// UpdateGroupParams updates a contact group.
type UpdateGroupParams struct {
	Name string
}

// ListGroupsParams are query params for List.
type ListGroupsParams struct {
	Page   *int
	Limit  *int
	Search *string
}

// GroupListResponse is a paginated group list.
type GroupListResponse struct {
	Object string                 `json:"object"`
	Groups []ContactGroupListItem `json:"groups"`
	Total  int                    `json:"total"`
	Page   int                    `json:"page"`
	Limit  int                    `json:"limit"`
	Event  string                 `json:"event"`
}

// DeleteGroupResponse is returned by Delete.
type DeleteGroupResponse struct {
	Object  string `json:"object"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Event   string `json:"event"`
}

// ListGroupContactsParams are query params for ListContacts.
type ListGroupContactsParams struct {
	Page   *int
	Limit  *int
	Search *string
	Status *ContactStatus
}

// GroupContactItem is a contact in a group list.
type GroupContactItem struct {
	ID         string                 `json:"id"`
	Email      string                 `json:"email"`
	FirstName  *string                `json:"firstName"`
	LastName   *string                `json:"lastName"`
	Status     ContactStatus          `json:"status"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  string                 `json:"updatedAt"`
}

// GroupContactListResponse is a paginated group contact list.
type GroupContactListResponse struct {
	Object   string             `json:"object"`
	Group    ContactGroupRef    `json:"group"`
	Contacts []GroupContactItem `json:"contacts"`
	Total    int                `json:"total"`
	Page     int                `json:"page"`
	Limit    int                `json:"limit"`
	Event    string             `json:"event"`
}

// GroupMembershipParams identifies a contact for group membership changes.
type GroupMembershipParams struct {
	ContactID *string
	Email     *string
}

// AddContactToGroupResponse is returned by AddContact.
type AddContactToGroupResponse struct {
	Success bool   `json:"success"`
	Object  string `json:"object"`
	ID      string `json:"id"`
	Event   string `json:"event"`
}

// RemoveContactFromGroupResponse is returned by RemoveContact.
type RemoveContactFromGroupResponse struct {
	Success bool   `json:"success"`
	Object  string `json:"object"`
	ID      string `json:"id"`
	Event   string `json:"event"`
}

// GroupsService manages contact groups.
type GroupsService struct {
	client *Client
}

func requireGroupID(id, field string) (string, error) {
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Group %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireCreateGroupName(name, field string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", validationErrorf(field, "Group %s is required and must be a string.", field)
	}
	if utf8.RuneCountInString(trimmed) < 1 {
		return "", validationErrorf(field, "Group %s must be at least 1 character.", field)
	}
	if utf8.RuneCountInString(trimmed) > 50 {
		return "", validationErrorf(field, "Group %s must be at most 50 characters.", field)
	}
	return trimmed, nil
}

func requireUpdateGroupName(name, field string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", validationErrorf(field, "Group %s is required and must be a string.", field)
	}
	if utf8.RuneCountInString(trimmed) < 1 {
		return "", validationErrorf(field, "Group %s must be at least 1 character.", field)
	}
	if utf8.RuneCountInString(trimmed) > 255 {
		return "", validationErrorf(field, "Group %s must be at most 255 characters.", field)
	}
	return trimmed, nil
}

func requireGroupListContactStatus(status ContactStatus, field string) error {
	switch status {
	case ContactStatusSubscribed, ContactStatusUnsubscribed, ContactStatusBlocked:
		return nil
	default:
		return validationErrorf(field, `list %s must be "subscribed", "unsubscribed", or "blocked".`, field)
	}
}

func validateGroupMembershipParams(params GroupMembershipParams) (map[string]interface{}, error) {
	body := map[string]interface{}{}

	if params.ContactID != nil {
		contactID, err := requireNonEmptyString(*params.ContactID, "contact_id")
		if err != nil {
			return nil, validationErrorf("contact_id", "contact_id must be a non-empty string when provided.")
		}
		body["contact_id"] = contactID
	}
	if params.Email != nil {
		email, err := requireNonEmptyString(*params.Email, "email")
		if err != nil {
			return nil, validationErrorf("email", "email must be a non-empty string when provided.")
		}
		body["email"] = email
	}

	if len(body) == 0 {
		return nil, validationErrorf("params", "Either contact_id or email is required.")
	}
	return body, nil
}

func validateListGroupsParams(params *ListGroupsParams) (url.Values, error) {
	if params == nil {
		return url.Values{}, nil
	}

	q := url.Values{}
	if params.Page != nil {
		if err := requirePage(*params.Page, "page"); err != nil {
			return nil, err
		}
		q.Set("page", strconv.Itoa(*params.Page))
	}
	if params.Limit != nil {
		if err := requireLimit(*params.Limit, 1, 100, "limit"); err != nil {
			return nil, err
		}
		q.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Search != nil {
		q.Set("search", *params.Search)
	}
	return q, nil
}

func validateListGroupContactsParams(params *ListGroupContactsParams) (url.Values, error) {
	if params == nil {
		return url.Values{}, nil
	}

	q := url.Values{}
	if params.Page != nil {
		if err := requirePage(*params.Page, "page"); err != nil {
			return nil, err
		}
		q.Set("page", strconv.Itoa(*params.Page))
	}
	if params.Limit != nil {
		if err := requireLimit(*params.Limit, 1, 100, "limit"); err != nil {
			return nil, err
		}
		q.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Search != nil {
		q.Set("search", *params.Search)
	}
	if params.Status != nil {
		if err := requireGroupListContactStatus(*params.Status, "status"); err != nil {
			return nil, err
		}
		q.Set("status", string(*params.Status))
	}
	return q, nil
}

// Create creates a contact group.
func (s *GroupsService) Create(params CreateGroupParams) (*ContactGroupResponse, error) {
	name, err := requireCreateGroupName(params.Name, "name")
	if err != nil {
		return nil, err
	}
	var out ContactGroupResponse
	if err := s.client.do("POST", groupsBase+"/create", map[string]string{"name": name}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List lists contact groups.
func (s *GroupsService) List(params *ListGroupsParams) (*GroupListResponse, error) {
	q, err := validateListGroupsParams(params)
	if err != nil {
		return nil, err
	}
	var out GroupListResponse
	if err := s.client.do("GET", groupsBase+"/list", nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a contact group by ID.
func (s *GroupsService) Get(id string) (*ContactGroup, error) {
	groupID, err := requireGroupID(id, "id")
	if err != nil {
		return nil, err
	}
	var out ContactGroup
	if err := s.client.do("GET", groupsBase+"/"+groupID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a contact group.
func (s *GroupsService) Update(id string, params UpdateGroupParams) (*ContactGroupResponse, error) {
	groupID, err := requireGroupID(id, "id")
	if err != nil {
		return nil, err
	}
	name, err := requireUpdateGroupName(params.Name, "name")
	if err != nil {
		return nil, err
	}
	var out ContactGroupResponse
	if err := s.client.do("PATCH", groupsBase+"/"+groupID, map[string]string{"name": name}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a contact group.
func (s *GroupsService) Delete(id string) (*DeleteGroupResponse, error) {
	groupID, err := requireGroupID(id, "id")
	if err != nil {
		return nil, err
	}
	var out DeleteGroupResponse
	if err := s.client.do("DELETE", groupsBase+"/"+groupID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListContacts lists contacts in a group.
func (s *GroupsService) ListContacts(id string, params *ListGroupContactsParams) (*GroupContactListResponse, error) {
	groupID, err := requireGroupID(id, "id")
	if err != nil {
		return nil, err
	}
	q, err := validateListGroupContactsParams(params)
	if err != nil {
		return nil, err
	}
	var out GroupContactListResponse
	if err := s.client.do("GET", groupsBase+"/"+groupID+"/contacts", nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddContact adds a contact to a group.
func (s *GroupsService) AddContact(id string, params GroupMembershipParams) (*AddContactToGroupResponse, error) {
	groupID, err := requireGroupID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := validateGroupMembershipParams(params)
	if err != nil {
		return nil, err
	}
	var out AddContactToGroupResponse
	if err := s.client.do("POST", groupMembership+"/"+groupID, body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveContact removes a contact from a group.
func (s *GroupsService) RemoveContact(id string, params GroupMembershipParams) (*RemoveContactFromGroupResponse, error) {
	groupID, err := requireGroupID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := validateGroupMembershipParams(params)
	if err != nil {
		return nil, err
	}
	var out RemoveContactFromGroupResponse
	if err := s.client.do("DELETE", groupMembership+"/"+groupID, body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
