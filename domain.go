package reloop

import (
	"fmt"
	"net/url"
)

type DomainService struct {
	client *Client
}

func (s *DomainService) Create(params CreateDomainParams) (*Domain, error) {
	var response Domain
	err := s.client.request("POST", "/api/domain/v1/create", params, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *DomainService) List(params *ListDomainsParams) (*DomainListResponse, error) {
	query := url.Values{}
	if params != nil {
		if params.Page != nil {
			query.Add("page", fmt.Sprintf("%d", *params.Page))
		}
		if params.Limit != nil {
			query.Add("limit", fmt.Sprintf("%d", *params.Limit))
		}
		if params.Q != nil {
			query.Add("q", *params.Q)
		}
		if params.Status != nil {
			query.Add("status", string(*params.Status))
		}
	}

	path := "/api/domain/v1/list"
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}

	var response DomainListResponse
	err := s.client.request("GET", path, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *DomainService) Get(domainID string) (*Domain, error) {
	var response Domain
	err := s.client.request("GET", fmt.Sprintf("/api/domain/v1/%s", domainID), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *DomainService) GetNameservers(domainID string) (*DomainNameserversResponse, error) {
	var response DomainNameserversResponse
	err := s.client.request(
		"GET",
		fmt.Sprintf("/api/domain/v1/nameservers/%s", domainID),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *DomainService) Update(domainID string, params UpdateDomainParams) (*Domain, error) {
	var response Domain
	err := s.client.request(
		"PATCH",
		fmt.Sprintf("/api/domain/v1/%s", domainID),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *DomainService) Delete(domainID string) (*Domain, error) {
	var response Domain
	err := s.client.request("DELETE", fmt.Sprintf("/api/domain/v1/%s", domainID), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *DomainService) Verify(domainID string) (*DomainStatusResponse, error) {
	var response DomainStatusResponse
	err := s.client.request("POST", fmt.Sprintf("/api/domain/v1/verify/%s", domainID), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *DomainService) ForwardDNS(domainID string, params ForwardDNSParams) (*ForwardDNSResponse, error) {
	var response ForwardDNSResponse
	err := s.client.request(
		"POST",
		fmt.Sprintf("/api/domain/v1/verify/%s/forward-dns", domainID),
		params,
		&response,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
