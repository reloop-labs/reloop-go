package reloop

import (
	"fmt"
	"net/url"
)

type APIKeyService struct {
	client *Client
}

func (s *APIKeyService) Create(params CreateApiKeyParams) (*ApiKeyWithKey, error) {
	var response ApiKeyWithKey
	err := s.client.request("POST", "/api-key/v1/", params, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *APIKeyService) List(params *ApiKeyListParams) (*ApiKeyListResponse, error) {
	query := url.Values{}
	if params != nil {
		if params.Page != nil {
			query.Add("page", fmt.Sprintf("%d", *params.Page))
		}
		if params.Limit != nil {
			query.Add("limit", fmt.Sprintf("%d", *params.Limit))
		}
		if params.Enabled != nil {
			query.Add("enabled", fmt.Sprintf("%t", *params.Enabled))
		}
		if params.UserID != nil {
			query.Add("userId", *params.UserID)
		}
		if params.Q != nil {
			query.Add("q", *params.Q)
		}
	}

	path := "/api-key/v1/"
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}

	var response ApiKeyListResponse
	err := s.client.request("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *APIKeyService) Get(id string) (*ApiKey, error) {
	var response ApiKey
	err := s.client.request("GET", fmt.Sprintf("/api-key/v1/%s", id), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *APIKeyService) Update(id string, params UpdateApiKeyParams) (*ApiKey, error) {
	var response ApiKey
	err := s.client.request("PATCH", fmt.Sprintf("/api-key/v1/%s", id), params, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *APIKeyService) Delete(id string) (*DeleteApiKeyResponse, error) {
	var response DeleteApiKeyResponse
	err := s.client.request("DELETE", fmt.Sprintf("/api-key/v1/%s", id), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *APIKeyService) Rotate(id string) (*ApiKeyWithKey, error) {
	var response ApiKeyWithKey
	err := s.client.request("POST", fmt.Sprintf("/api-key/v1/rotate/%s", id), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *APIKeyService) Enable(id string) (*ApiKey, error) {
	var response ApiKey
	err := s.client.request("POST", fmt.Sprintf("/api-key/v1/enable/%s", id), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *APIKeyService) Disable(id string) (*ApiKey, error) {
	var response ApiKey
	err := s.client.request("POST", fmt.Sprintf("/api-key/v1/disable/%s", id), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
