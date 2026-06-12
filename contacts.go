package reloop

import (
	"fmt"
	"net/url"
)

type ContactsService struct {
	client   *Client
	Groups   *ContactGroupsService
	Channels *ContactChannelsService
}

func newContactsService(client *Client) *ContactsService {
	service := &ContactsService{client: client}
	service.Groups = &ContactGroupsService{client: client}
	service.Channels = &ContactChannelsService{client: client}
	return service
}

func (s *ContactsService) Create(parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request("POST", "/api/contacts/create", forRequest(parameters), &response)
	return response, err
}

func (s *ContactsService) Get(contactID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request("GET", fmt.Sprintf("/api/contacts/retrieve/%s", contactID), nil, &response)
	return response, err
}

func (s *ContactsService) List(options map[string]interface{}) (map[string]interface{}, error) {
	if groupID, ok := options["group_id"]; ok && groupID != nil {
		filtered := make(map[string]interface{}, len(options))
		for key, value := range options {
			if key != "group_id" && key != "groupId" {
				filtered[key] = value
			}
		}

		groupIDString, ok := groupID.(string)
		if !ok {
			return nil, fmt.Errorf("group_id must be a string")
		}

		return s.Groups.ListContacts(groupIDString, filtered)
	}

	path := "/api/contacts/list"
	if query := buildQuery(forQuery(options)); query != "" {
		path += "?" + query
	}

	var response map[string]interface{}
	err := s.client.request("GET", path, nil, &response)
	return response, err
}

func (s *ContactsService) Update(contactID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"PATCH",
		fmt.Sprintf("/api/contacts/%s", contactID),
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactsService) Delete(contactID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request("DELETE", fmt.Sprintf("/api/contacts/%s", contactID), nil, &response)
	return response, err
}

func (s *ContactsService) CreateProperty(parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"POST",
		"/api/contacts/v1/properties/create",
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactsService) ListProperties(options map[string]interface{}) (map[string]interface{}, error) {
	path := "/api/contacts/v1/properties/list"
	if query := buildQuery(forQuery(options)); query != "" {
		path += "?" + query
	}

	var response map[string]interface{}
	err := s.client.request("GET", path, nil, &response)
	return response, err
}

func (s *ContactsService) UpdateProperty(propertyID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"PATCH",
		fmt.Sprintf("/api/contacts/v1/properties/%s", propertyID),
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactsService) DeleteProperty(propertyID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"DELETE",
		fmt.Sprintf("/api/contacts/v1/properties/%s", propertyID),
		nil,
		&response,
	)
	return response, err
}

func (s *ContactsService) CreateGroup(parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"POST",
		"/api/contacts/v1/groups/create",
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactsService) ListGroups(options map[string]interface{}) (map[string]interface{}, error) {
	path := "/api/contacts/v1/groups/list"
	if query := buildQuery(forQuery(options)); query != "" {
		path += "?" + query
	}

	var response map[string]interface{}
	err := s.client.request("GET", path, nil, &response)
	return response, err
}

func (s *ContactsService) GetGroup(groupID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request("GET", fmt.Sprintf("/api/contacts/v1/groups/%s", groupID), nil, &response)
	return response, err
}

func (s *ContactsService) UpdateGroup(groupID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"PATCH",
		fmt.Sprintf("/api/contacts/v1/groups/%s", groupID),
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactsService) DeleteGroup(groupID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request("DELETE", fmt.Sprintf("/api/contacts/v1/groups/%s", groupID), nil, &response)
	return response, err
}

func buildQuery(values map[string]interface{}) string {
	if len(values) == 0 {
		return ""
	}

	query := url.Values{}
	for key, value := range values {
		query.Add(key, fmt.Sprintf("%v", value))
	}

	return query.Encode()
}
