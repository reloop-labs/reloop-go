package reloop

import (
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	channelsBase      = "/api/contacts/v1/channels"
	channelMembership = "/api/contacts/channel"
)

// ChannelVisibility is a channel visibility setting.
type ChannelVisibility string

const (
	ChannelVisibilityPrivate ChannelVisibility = "private"
	ChannelVisibilityPublic  ChannelVisibility = "public"
)

// ContactChannel is a contact channel.
type ContactChannel struct {
	Object              string              `json:"object"`
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         *string             `json:"description"`
	DefaultSubscription ChannelSubscription `json:"defaultSubscription"`
	Visibility          ChannelVisibility   `json:"visibility"`
	CreatedAt           string              `json:"createdAt"`
	UpdatedAt           string              `json:"updatedAt"`
}

// ContactChannelResponse is returned by channel mutations.
type ContactChannelResponse struct {
	ContactChannel
	Event string `json:"event"`
}

// ContactChannelListItem is a channel in a list response.
type ContactChannelListItem struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         *string             `json:"description"`
	DefaultSubscription ChannelSubscription `json:"defaultSubscription"`
	Visibility          ChannelVisibility   `json:"visibility"`
	CreatedAt           string              `json:"createdAt"`
	UpdatedAt           string              `json:"updatedAt"`
	SubscriberCount     *int                `json:"subscriberCount,omitempty"`
}

// CreateChannelParams creates a contact channel.
type CreateChannelParams struct {
	Name                string
	Description         *string
	DefaultSubscription *ChannelSubscription
	Visibility          *ChannelVisibility
}

// UpdateChannelParams updates a contact channel.
type UpdateChannelParams struct {
	Name        *string
	Description **string
	Visibility  *ChannelVisibility
}

// ListChannelsParams are query params for List.
type ListChannelsParams struct {
	Page  *int
	Limit *int
}

// ChannelListResponse is a paginated channel list.
type ChannelListResponse struct {
	Object   string                   `json:"object"`
	Channels []ContactChannelListItem `json:"channels"`
	Total    int                      `json:"total"`
	Page     int                      `json:"page"`
	Limit    int                      `json:"limit"`
	Event    string                   `json:"event"`
}

// DeleteChannelResponse is returned by Delete.
type DeleteChannelResponse struct {
	Object  string `json:"object"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Event   string `json:"event"`
}

// ChannelMembershipParams identifies a contact for channel membership changes.
type ChannelMembershipParams struct {
	ContactID *string
	Email     *string
}

// AddContactToChannelParams adds a contact to a channel.
type AddContactToChannelParams struct {
	ContactID    *string
	Email        *string
	Subscription *ChannelSubscription
}

// AddContactToChannelResponse is returned by AddContact.
type AddContactToChannelResponse struct {
	Contact        ContactResponse `json:"contact"`
	SubscriptionID string          `json:"subscriptionId"`
	Event          string          `json:"event"`
}

// UpdateContactChannelParams updates a contact channel subscription.
type UpdateContactChannelParams struct {
	ContactID    *string
	Email        *string
	Subscription ChannelSubscription
}

// UpdateContactChannelResponse is returned by UpdateSubscription.
type UpdateContactChannelResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Event   string `json:"event"`
}

// ChannelsService manages contact channels.
type ChannelsService struct {
	client *Client
}

func requireChannelID(id, field string) (string, error) {
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Channel %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireChannelName(name, field string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", validationErrorf(field, "Channel %s is required and must be a string.", field)
	}
	if utf8.RuneCountInString(trimmed) < 1 {
		return "", validationErrorf(field, "Channel %s must be at least 1 character.", field)
	}
	if utf8.RuneCountInString(trimmed) > 255 {
		return "", validationErrorf(field, "Channel %s must be at most 255 characters.", field)
	}
	return trimmed, nil
}

func optionalChannelDescription(description *string, field string) (*string, error) {
	if description == nil {
		return nil, nil
	}
	if utf8.RuneCountInString(*description) > 1000 {
		return nil, validationErrorf(field, "Channel %s must be at most 1000 characters.", field)
	}
	return description, nil
}

func optionalChannelSubscription(subscription *ChannelSubscription, field string) (*ChannelSubscription, error) {
	if subscription == nil {
		return nil, nil
	}
	switch *subscription {
	case ChannelSubscriptionOptIn, ChannelSubscriptionOptOut:
		return subscription, nil
	default:
		return nil, validationErrorf(field, `Channel %s must be "opt_in" or "opt_out".`, field)
	}
}

func requireChannelSubscription(subscription ChannelSubscription, field string) (ChannelSubscription, error) {
	switch subscription {
	case ChannelSubscriptionOptIn, ChannelSubscriptionOptOut:
		return subscription, nil
	default:
		return "", validationErrorf(field, `Channel %s is required and must be "opt_in" or "opt_out".`, field)
	}
}

func optionalChannelVisibility(visibility *ChannelVisibility, field string) (*ChannelVisibility, error) {
	if visibility == nil {
		return nil, nil
	}
	switch *visibility {
	case ChannelVisibilityPrivate, ChannelVisibilityPublic:
		return visibility, nil
	default:
		return nil, validationErrorf(field, `Channel %s must be "private" or "public".`, field)
	}
}

func validateChannelMembershipParams(params ChannelMembershipParams) (map[string]interface{}, error) {
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

func validateCreateChannelParams(params CreateChannelParams) (map[string]interface{}, error) {
	name, err := requireChannelName(params.Name, "name")
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{"name": name}

	description, err := optionalChannelDescription(params.Description, "description")
	if err != nil {
		return nil, err
	}
	if description != nil {
		body["description"] = *description
	}

	defaultSubscription, err := optionalChannelSubscription(params.DefaultSubscription, "defaultSubscription")
	if err != nil {
		return nil, err
	}
	if defaultSubscription != nil {
		body["defaultSubscription"] = string(*defaultSubscription)
	}

	visibility, err := optionalChannelVisibility(params.Visibility, "visibility")
	if err != nil {
		return nil, err
	}
	if visibility != nil {
		body["visibility"] = string(*visibility)
	}

	return body, nil
}

func validateUpdateChannelParams(params UpdateChannelParams) (map[string]interface{}, error) {
	body := map[string]interface{}{}

	if params.Name != nil {
		name, err := requireChannelName(*params.Name, "name")
		if err != nil {
			return nil, err
		}
		body["name"] = name
	}
	if params.Description != nil {
		if *params.Description == nil {
			body["description"] = nil
		} else {
			description, err := optionalChannelDescription(*params.Description, "description")
			if err != nil {
				return nil, err
			}
			body["description"] = *description
		}
	}
	if params.Visibility != nil {
		visibility, err := optionalChannelVisibility(params.Visibility, "visibility")
		if err != nil {
			return nil, err
		}
		body["visibility"] = string(*visibility)
	}

	if len(body) == 0 {
		return nil, validationErrorf("params", "update requires at least one of name, description, or visibility.")
	}

	return body, nil
}

func validateAddContactToChannelParams(params AddContactToChannelParams) (map[string]interface{}, error) {
	body, err := validateChannelMembershipParams(ChannelMembershipParams{
		ContactID: params.ContactID,
		Email:     params.Email,
	})
	if err != nil {
		return nil, err
	}

	subscription, err := optionalChannelSubscription(params.Subscription, "subscription")
	if err != nil {
		return nil, err
	}
	if subscription != nil {
		body["subscription"] = string(*subscription)
	}

	return body, nil
}

func validateUpdateContactChannelParams(params UpdateContactChannelParams) (map[string]interface{}, error) {
	body, err := validateChannelMembershipParams(ChannelMembershipParams{
		ContactID: params.ContactID,
		Email:     params.Email,
	})
	if err != nil {
		return nil, err
	}

	subscription, err := requireChannelSubscription(params.Subscription, "subscription")
	if err != nil {
		return nil, err
	}
	body["subscription"] = string(subscription)

	return body, nil
}

func validateListChannelsParams(params *ListChannelsParams) (url.Values, error) {
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
	return q, nil
}

// Create creates a contact channel.
func (s *ChannelsService) Create(params CreateChannelParams) (*ContactChannelResponse, error) {
	body, err := validateCreateChannelParams(params)
	if err != nil {
		return nil, err
	}
	var out ContactChannelResponse
	if err := s.client.do("POST", channelsBase+"/create", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List lists contact channels.
func (s *ChannelsService) List(params *ListChannelsParams) (*ChannelListResponse, error) {
	q, err := validateListChannelsParams(params)
	if err != nil {
		return nil, err
	}
	var out ChannelListResponse
	if err := s.client.do("GET", channelsBase+"/list", nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a contact channel by ID.
func (s *ChannelsService) Get(id string) (*ContactChannel, error) {
	channelID, err := requireChannelID(id, "id")
	if err != nil {
		return nil, err
	}
	var out ContactChannel
	if err := s.client.do("GET", channelsBase+"/"+channelID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a contact channel.
func (s *ChannelsService) Update(id string, params UpdateChannelParams) (*ContactChannelResponse, error) {
	channelID, err := requireChannelID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := validateUpdateChannelParams(params)
	if err != nil {
		return nil, err
	}
	var out ContactChannelResponse
	if err := s.client.do("PATCH", channelsBase+"/"+channelID, body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a contact channel.
func (s *ChannelsService) Delete(id string) (*DeleteChannelResponse, error) {
	channelID, err := requireChannelID(id, "id")
	if err != nil {
		return nil, err
	}
	var out DeleteChannelResponse
	if err := s.client.do("DELETE", channelsBase+"/"+channelID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddContact adds a contact to a channel.
func (s *ChannelsService) AddContact(id string, params AddContactToChannelParams) (*AddContactToChannelResponse, error) {
	channelID, err := requireChannelID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := validateAddContactToChannelParams(params)
	if err != nil {
		return nil, err
	}
	var out AddContactToChannelResponse
	if err := s.client.do("POST", channelMembership+"/"+channelID, body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSubscription updates a contact's channel subscription.
func (s *ChannelsService) UpdateSubscription(id string, params UpdateContactChannelParams) (*UpdateContactChannelResponse, error) {
	channelID, err := requireChannelID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := validateUpdateContactChannelParams(params)
	if err != nil {
		return nil, err
	}
	var out UpdateContactChannelResponse
	if err := s.client.do("PATCH", channelMembership+"/"+channelID, body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
