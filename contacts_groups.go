package reloop

import "fmt"

type ContactGroupsService struct {
	client *Client
}

func (s *ContactGroupsService) AddContact(groupID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"POST",
		fmt.Sprintf("/api/contacts/group/%s", groupID),
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactGroupsService) RemoveContact(groupID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"DELETE",
		fmt.Sprintf("/api/contacts/group/%s", groupID),
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactGroupsService) ListContacts(groupID string, options map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/contacts/v1/groups/%s/contacts", groupID)
	if query := buildQuery(forQuery(options)); query != "" {
		path += "?" + query
	}

	var response map[string]interface{}
	err := s.client.request("GET", path, nil, &response)
	return response, err
}
