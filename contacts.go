package reloop

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const contactsBase = "/api/contacts"

var contactEmailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
var createPropertyKeyPattern = regexp.MustCompile(`^[a-z0-9_]+$`)
var updatePropertyKeyPattern = regexp.MustCompile(`^[a-z_]+$`)

// ContactStatus is a contact subscription status.
type ContactStatus string

const (
	ContactStatusSubscribed   ContactStatus = "subscribed"
	ContactStatusUnsubscribed ContactStatus = "unsubscribed"
	ContactStatusBlocked      ContactStatus = "blocked"
)

// ChannelSubscription is a channel opt-in state.
type ChannelSubscription string

const (
	ChannelSubscriptionOptIn  ChannelSubscription = "opt_in"
	ChannelSubscriptionOptOut ChannelSubscription = "opt_out"
)

// ContactChannelInput assigns a channel on contact create.
type ContactChannelInput struct {
	ChannelID    string              `json:"channelId"`
	Subscription ChannelSubscription `json:"subscription"`
}

// ContactGroupRef is a group on a contact.
type ContactGroupRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ContactChannelRef is a channel on a contact.
type ContactChannelRef struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Subscription ChannelSubscription `json:"subscription"`
}

// Contact is a contact resource.
type Contact struct {
	Object            string                 `json:"object"`
	ID                string                 `json:"id"`
	Email             string                 `json:"email"`
	FirstName         *string                `json:"firstName"`
	LastName          *string                `json:"lastName"`
	Status            ContactStatus          `json:"status"`
	Properties        map[string]interface{} `json:"properties"`
	Groups            []ContactGroupRef      `json:"groups"`
	Channels          []ContactChannelRef    `json:"channels"`
	SuppressionReason *string                `json:"suppressionReason"`
	SuppressedAt      *string                `json:"suppressedAt"`
	CreatedAt         string                 `json:"createdAt"`
	UpdatedAt         string                 `json:"updatedAt"`
}

// ContactResponse is returned by contact mutations.
type ContactResponse struct {
	Contact
	Event string `json:"event"`
}

// CreateContactParams creates a contact.
type CreateContactParams struct {
	Email      string
	FirstName  *string
	LastName   *string
	Status     *ContactStatus
	Properties map[string]interface{}
	GroupIDs   []string
	Channels   []ContactChannelInput
}

// UpdateContactParams updates a contact.
type UpdateContactParams struct {
	Email      *string
	FirstName  *string
	LastName   *string
	Status     *ContactStatus
	Properties map[string]interface{}
}

// ListContactsParams are query params for List.
type ListContactsParams struct {
	Page   *int
	Limit  *int
	Search *string
	Status *ContactStatus
}

// ContactListResponse is a paginated contact list.
type ContactListResponse struct {
	Object               string    `json:"object"`
	Contacts             []Contact `json:"contacts"`
	Total                int       `json:"total"`
	Page                 int       `json:"page"`
	Limit                int       `json:"limit"`
	TotalContacts        int       `json:"totalContacts"`
	SubscribedContacts   int       `json:"subscribedContacts"`
	UnsubscribedContacts int       `json:"unsubscribedContacts"`
	Event                string    `json:"event"`
}

// DeleteContactResponse is returned by Delete.
type DeleteContactResponse struct {
	Success bool   `json:"success"`
	Object  string `json:"object"`
	ID      string `json:"id"`
	Event   string `json:"event"`
}

// ContactsService manages contacts.
type ContactsService struct {
	client *Client

	Properties *PropertiesService
	Groups     *GroupsService
	Channels   *ChannelsService
}

func newContactsService(c *Client) *ContactsService {
	s := &ContactsService{client: c}
	s.Properties = &PropertiesService{client: c}
	s.Groups = &GroupsService{client: c}
	s.Channels = &ChannelsService{client: c}
	return s
}

func requireContactID(id, field string) (string, error) {
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Contact %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireContactEmail(email, field string) (string, error) {
	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return "", validationErrorf(field, "Contact %s is required and must be a string.", field)
	}
	if !contactEmailPattern.MatchString(trimmed) {
		return "", validationErrorf(field, "Contact %s must be a valid email address.", field)
	}
	return trimmed, nil
}

func optionalContactStatus(status *ContactStatus, field string) (*ContactStatus, error) {
	if status == nil {
		return nil, nil
	}
	switch *status {
	case ContactStatusSubscribed, ContactStatusUnsubscribed, ContactStatusBlocked:
		return status, nil
	default:
		return nil, validationErrorf(field, `%s must be "subscribed", "unsubscribed", or "blocked".`, field)
	}
}

func validateContactProperties(properties map[string]interface{}, keyPattern *regexp.Regexp, field string) (map[string]interface{}, error) {
	if properties == nil {
		return nil, nil
	}
	out := make(map[string]interface{}, len(properties))
	for key, value := range properties {
		if !keyPattern.MatchString(key) {
			return nil, validationErrorf(field, "%s keys must match %s.", field, keyPattern.String())
		}
		switch value.(type) {
		case string, float64, int, int64:
			out[key] = value
		default:
			return nil, validationErrorf(field, "%s.%s must be a string or number.", field, key)
		}
	}
	return out, nil
}

func validateCreateContactChannels(channels []ContactChannelInput) ([]ContactChannelInput, error) {
	if channels == nil {
		return nil, nil
	}
	out := make([]ContactChannelInput, len(channels))
	for i, entry := range channels {
		channelID, err := requireNonEmptyString(entry.ChannelID, "channels")
		if err != nil {
			return nil, validationErrorf("channels", "channels[%d].channelId must be a non-empty string.", i)
		}
		switch entry.Subscription {
		case ChannelSubscriptionOptIn, ChannelSubscriptionOptOut:
		default:
			return nil, validationErrorf("channels", `channels[%d].subscription must be "opt_in" or "opt_out".`, i)
		}
		out[i] = ContactChannelInput{
			ChannelID:    channelID,
			Subscription: entry.Subscription,
		}
	}
	return out, nil
}

func validateCreateContactGroupIDs(groupIDs []string) ([]string, error) {
	if groupIDs == nil {
		return nil, nil
	}
	out := make([]string, len(groupIDs))
	for i, id := range groupIDs {
		trimmed, err := requireNonEmptyString(id, "groupIds")
		if err != nil {
			return nil, validationErrorf("groupIds", "groupIds[%d] must be a non-empty string.", i)
		}
		out[i] = trimmed
	}
	return out, nil
}

func validateCreateContactParams(params CreateContactParams) (map[string]interface{}, error) {
	email, err := requireContactEmail(params.Email, "email")
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{"email": email}

	if params.FirstName != nil {
		body["firstName"] = *params.FirstName
	}
	if params.LastName != nil {
		body["lastName"] = *params.LastName
	}

	status, err := optionalContactStatus(params.Status, "status")
	if err != nil {
		return nil, err
	}
	if status != nil {
		body["status"] = string(*status)
	}

	properties, err := validateContactProperties(params.Properties, createPropertyKeyPattern, "properties")
	if err != nil {
		return nil, err
	}
	if properties != nil {
		body["properties"] = properties
	}

	groupIDs, err := validateCreateContactGroupIDs(params.GroupIDs)
	if err != nil {
		return nil, err
	}
	if groupIDs != nil {
		body["groupIds"] = groupIDs
	}

	channels, err := validateCreateContactChannels(params.Channels)
	if err != nil {
		return nil, err
	}
	if channels != nil {
		body["channels"] = channels
	}

	return body, nil
}

func validateUpdateContactParams(params UpdateContactParams) (map[string]interface{}, error) {
	body := map[string]interface{}{}

	if params.Email != nil {
		email, err := requireContactEmail(*params.Email, "email")
		if err != nil {
			return nil, err
		}
		body["email"] = email
	}
	if params.FirstName != nil {
		body["firstName"] = *params.FirstName
	}
	if params.LastName != nil {
		body["lastName"] = *params.LastName
	}

	status, err := optionalContactStatus(params.Status, "status")
	if err != nil {
		return nil, err
	}
	if status != nil {
		body["status"] = string(*status)
	}

	properties, err := validateContactProperties(params.Properties, updatePropertyKeyPattern, "properties")
	if err != nil {
		return nil, err
	}
	if properties != nil {
		body["properties"] = properties
	}

	if len(body) == 0 {
		return nil, validationErrorf("params", "update requires at least one of email, firstName, lastName, status, or properties.")
	}

	return body, nil
}

func validateListContactsParams(params *ListContactsParams) (url.Values, error) {
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
		if _, err := optionalContactStatus(params.Status, "status"); err != nil {
			return nil, err
		}
		q.Set("status", string(*params.Status))
	}
	return q, nil
}

// Create creates a contact.
func (s *ContactsService) Create(params CreateContactParams) (*ContactResponse, error) {
	body, err := validateCreateContactParams(params)
	if err != nil {
		return nil, err
	}
	var out ContactResponse
	if err := s.client.do("POST", contactsBase+"/create", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a contact by ID.
func (s *ContactsService) Get(id string) (*Contact, error) {
	contactID, err := requireContactID(id, "id")
	if err != nil {
		return nil, err
	}
	var out Contact
	if err := s.client.do("GET", contactsBase+"/retrieve/"+contactID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List lists contacts.
func (s *ContactsService) List(params *ListContactsParams) (*ContactListResponse, error) {
	q, err := validateListContactsParams(params)
	if err != nil {
		return nil, err
	}
	var out ContactListResponse
	if err := s.client.do("GET", contactsBase+"/list", nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a contact.
func (s *ContactsService) Update(id string, params UpdateContactParams) (*ContactResponse, error) {
	contactID, err := requireContactID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := validateUpdateContactParams(params)
	if err != nil {
		return nil, err
	}
	var out ContactResponse
	if err := s.client.do("PATCH", contactsBase+"/"+contactID, body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a contact.
func (s *ContactsService) Delete(id string) (*DeleteContactResponse, error) {
	contactID, err := requireContactID(id, "id")
	if err != nil {
		return nil, err
	}
	var out DeleteContactResponse
	if err := s.client.do("DELETE", contactsBase+"/"+contactID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
