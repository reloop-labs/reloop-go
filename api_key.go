package reloop

import (
	"net/url"
	"strconv"
)

// APIKeyService manages organization API keys.
type APIKeyService struct {
	client *Client
}

// User is an API key creator.
type User struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Image *string `json:"image"`
	Email string  `json:"email"`
}

// APIKey is an API key resource.
type APIKey struct {
	ID                  string  `json:"id"`
	Name                *string `json:"name"`
	Start               *string `json:"start"`
	Prefix              *string `json:"prefix"`
	OrganizationID      string  `json:"organizationId"`
	UserID              string  `json:"userId"`
	RefillInterval      *int    `json:"refillInterval"`
	RefillAmount        *int    `json:"refillAmount"`
	LastRefillAt        *string `json:"lastRefillAt"`
	Enabled             bool    `json:"enabled"`
	RateLimitEnabled    bool    `json:"rateLimitEnabled"`
	RateLimitTimeWindow int     `json:"rateLimitTimeWindow"`
	RateLimitMax        int     `json:"rateLimitMax"`
	RequestCount        int     `json:"requestCount"`
	Remaining           *int    `json:"remaining"`
	LastRequest         *string `json:"lastRequest"`
	ExpiresAt           *string `json:"expiresAt"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
	Permissions         *string `json:"permissions"`
	Metadata            *string `json:"metadata"`
	CreatedBy           *User   `json:"createdBy,omitempty"`
	Object              string  `json:"object"`
	Event               string  `json:"event"`
}

// APIKeyWithKey includes the raw secret (create/rotate).
type APIKeyWithKey struct {
	ID          string  `json:"id"`
	Name        *string `json:"name"`
	Key         string  `json:"key"`
	Enabled     bool    `json:"enabled"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Permissions *string `json:"permissions"`
	Object      string  `json:"object"`
	Event       string  `json:"event"`
}

// APIKeyListResponse is a paginated API key list.
type APIKeyListResponse struct {
	Object  string   `json:"object"`
	APIKeys []APIKey `json:"apiKeys"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
	Event   string   `json:"event"`
}

// APIKeyListParams are query params for List.
type APIKeyListParams struct {
	Page    *int
	Limit   *int
	Enabled *bool
	UserID  *string
	Q       *string
}

// DeleteAPIKeyResponse is returned by Delete.
type DeleteAPIKeyResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Object  string `json:"object"`
	Event   string `json:"event"`
}

// CreateAPIKeyParams creates an API key.
type CreateAPIKeyParams struct {
	Name string `json:"name"`
}

// UpdateAPIKeyParams updates an API key.
type UpdateAPIKeyParams struct {
	Name string `json:"name"`
}

// Create creates an API key.
func (s *APIKeyService) Create(params CreateAPIKeyParams) (*APIKeyWithKey, error) {
	name, err := requireAPIKeyName(params.Name, "name")
	if err != nil {
		return nil, err
	}
	var out APIKeyWithKey
	if err := s.client.do("POST", "/api/api-key/v1/", map[string]string{"name": name}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

const apiKeyV1 = "/api/api-key/v1"

// List lists API keys.
func (s *APIKeyService) List(params *APIKeyListParams) (*APIKeyListResponse, error) {
	q := url.Values{}
	if params != nil {
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
		if params.Enabled != nil {
			q.Set("enabled", strconv.FormatBool(*params.Enabled))
		}
		if params.UserID != nil {
			q.Set("userId", *params.UserID)
		}
		if params.Q != nil {
			q.Set("q", *params.Q)
		}
	}
	var out APIKeyListResponse
	if err := s.client.do("GET", apiKeyV1+"/", nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get fetches an API key by ID.
func (s *APIKeyService) Get(apiKeyID string) (*APIKey, error) {
	id, err := requireAPIKeyID(apiKeyID, "apiKeyId")
	if err != nil {
		return nil, err
	}
	var out APIKey
	if err := s.client.do("GET", apiKeyV1+"/"+id, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates an API key.
func (s *APIKeyService) Update(apiKeyID string, params UpdateAPIKeyParams) (*APIKey, error) {
	id, err := requireAPIKeyID(apiKeyID, "apiKeyId")
	if err != nil {
		return nil, err
	}
	name, err := requireAPIKeyName(params.Name, "name")
	if err != nil {
		return nil, err
	}
	var out APIKey
	if err := s.client.do("PATCH", apiKeyV1+"/"+id, map[string]string{"name": name}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes an API key.
func (s *APIKeyService) Delete(apiKeyID string) (*DeleteAPIKeyResponse, error) {
	id, err := requireAPIKeyID(apiKeyID, "apiKeyId")
	if err != nil {
		return nil, err
	}
	var out DeleteAPIKeyResponse
	if err := s.client.do("DELETE", apiKeyV1+"/"+id, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Rotate rotates an API key secret.
func (s *APIKeyService) Rotate(apiKeyID string) (*APIKeyWithKey, error) {
	id, err := requireAPIKeyID(apiKeyID, "apiKeyId")
	if err != nil {
		return nil, err
	}
	var out APIKeyWithKey
	if err := s.client.do("POST", apiKeyV1+"/rotate/"+id, map[string]interface{}{}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Enable enables an API key.
func (s *APIKeyService) Enable(apiKeyID string) (*APIKey, error) {
	id, err := requireAPIKeyID(apiKeyID, "apiKeyId")
	if err != nil {
		return nil, err
	}
	var out APIKey
	if err := s.client.do("POST", apiKeyV1+"/enable/"+id, map[string]interface{}{}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Disable disables an API key.
func (s *APIKeyService) Disable(apiKeyID string) (*APIKey, error) {
	id, err := requireAPIKeyID(apiKeyID, "apiKeyId")
	if err != nil {
		return nil, err
	}
	var out APIKey
	if err := s.client.do("POST", apiKeyV1+"/disable/"+id, map[string]interface{}{}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
