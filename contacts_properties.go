package reloop

import (
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

const propertiesBase = "/api/contacts/v1/properties"

// PropertyType is a contact property data type.
type PropertyType string

const (
	PropertyTypeString PropertyType = "string"
	PropertyTypeNumber PropertyType = "number"
)

// ContactProperty is a contact property definition.
type ContactProperty struct {
	Object       string       `json:"object"`
	ID           string       `json:"id"`
	PropertyName string       `json:"propertyName"`
	PropertyType PropertyType `json:"propertyType"`
	DefaultValue *string      `json:"defaultValue"`
	CreatedAt    string       `json:"createdAt"`
	UpdatedAt    string       `json:"updatedAt"`
}

// ContactPropertyResponse is returned by property mutations.
type ContactPropertyResponse struct {
	ContactProperty
	Event string `json:"event"`
}

// ContactPropertyListItem is a property in a list response.
type ContactPropertyListItem struct {
	ID           string       `json:"id"`
	PropertyName string       `json:"propertyName"`
	PropertyType PropertyType `json:"propertyType"`
	DefaultValue *string      `json:"defaultValue"`
	CreatedAt    string       `json:"createdAt"`
	UpdatedAt    string       `json:"updatedAt"`
}

// CreatePropertyParams creates a contact property.
type CreatePropertyParams struct {
	Name          string
	Type          PropertyType
	FallbackValue *string
}

// UpdatePropertyParams updates a contact property.
type UpdatePropertyParams struct {
	FallbackValue *string
}

// ListPropertiesParams are query params for List.
type ListPropertiesParams struct {
	Page   *int
	Limit  *int
	Search *string
	Type   *PropertyType
}

// PropertyListResponse is a paginated property list.
type PropertyListResponse struct {
	Object     string                    `json:"object"`
	Properties []ContactPropertyListItem `json:"properties"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	Event      string                    `json:"event"`
}

// DeletePropertyResponse is returned by Delete.
type DeletePropertyResponse struct {
	Object  string `json:"object"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Event   string `json:"event"`
}

// PropertiesService manages contact properties.
type PropertiesService struct {
	client *Client
}

func requirePropertyID(id, field string) (string, error) {
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Property %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requirePropertyName(name, field string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", validationErrorf(field, "Property %s is required and must be a string.", field)
	}
	if utf8.RuneCountInString(trimmed) < 1 {
		return "", validationErrorf(field, "Property %s must be at least 1 character.", field)
	}
	if utf8.RuneCountInString(trimmed) > 255 {
		return "", validationErrorf(field, "Property %s must be at most 255 characters.", field)
	}
	return trimmed, nil
}

func requirePropertyType(propertyType PropertyType, field string) (PropertyType, error) {
	switch propertyType {
	case PropertyTypeString, PropertyTypeNumber:
		return propertyType, nil
	default:
		return "", validationErrorf(field, `Property %s must be "string" or "number".`, field)
	}
}

func validateCreatePropertyParams(params CreatePropertyParams) (map[string]interface{}, error) {
	name, err := requirePropertyName(params.Name, "name")
	if err != nil {
		return nil, err
	}
	propertyType, err := requirePropertyType(params.Type, "type")
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"name": name,
		"type": string(propertyType),
	}
	if params.FallbackValue != nil {
		body["fallbackValue"] = *params.FallbackValue
	}
	return body, nil
}

func validateUpdatePropertyParams(params UpdatePropertyParams) (map[string]interface{}, error) {
	return map[string]interface{}{"fallbackValue": params.FallbackValue}, nil
}

func validateListPropertiesParams(params *ListPropertiesParams) (url.Values, error) {
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
	if params.Type != nil {
		if _, err := requirePropertyType(*params.Type, "type"); err != nil {
			return nil, err
		}
		q.Set("type", string(*params.Type))
	}
	return q, nil
}

// Create creates a contact property.
func (s *PropertiesService) Create(params CreatePropertyParams) (*ContactPropertyResponse, error) {
	body, err := validateCreatePropertyParams(params)
	if err != nil {
		return nil, err
	}
	var out ContactPropertyResponse
	if err := s.client.do("POST", propertiesBase+"/create", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List lists contact properties.
func (s *PropertiesService) List(params *ListPropertiesParams) (*PropertyListResponse, error) {
	q, err := validateListPropertiesParams(params)
	if err != nil {
		return nil, err
	}
	var out PropertyListResponse
	if err := s.client.do("GET", propertiesBase+"/list", nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a contact property.
func (s *PropertiesService) Update(id string, params UpdatePropertyParams) (*ContactPropertyResponse, error) {
	propertyID, err := requirePropertyID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := validateUpdatePropertyParams(params)
	if err != nil {
		return nil, err
	}
	var out ContactPropertyResponse
	if err := s.client.do("PATCH", propertiesBase+"/"+propertyID, body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a contact property.
func (s *PropertiesService) Delete(id string) (*DeletePropertyResponse, error) {
	propertyID, err := requirePropertyID(id, "id")
	if err != nil {
		return nil, err
	}
	var out DeletePropertyResponse
	if err := s.client.do("DELETE", propertiesBase+"/"+propertyID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
