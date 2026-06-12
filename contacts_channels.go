package reloop

import "fmt"

type ContactChannelsService struct {
	client *Client
}

func (s *ContactChannelsService) Create(parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"POST",
		"/api/contacts/v1/channels/create",
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactChannelsService) List(options map[string]interface{}) (map[string]interface{}, error) {
	path := "/api/contacts/v1/channels/list"
	if query := buildQuery(forQuery(options)); query != "" {
		path += "?" + query
	}

	var response map[string]interface{}
	err := s.client.request("GET", path, nil, &response)
	return response, err
}

func (s *ContactChannelsService) Get(channelID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request("GET", fmt.Sprintf("/api/contacts/v1/channels/%s", channelID), nil, &response)
	return response, err
}

func (s *ContactChannelsService) Update(channelID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"PATCH",
		fmt.Sprintf("/api/contacts/v1/channels/%s", channelID),
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactChannelsService) Delete(channelID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"DELETE",
		fmt.Sprintf("/api/contacts/v1/channels/%s", channelID),
		nil,
		&response,
	)
	return response, err
}

func (s *ContactChannelsService) AddContact(channelID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"POST",
		fmt.Sprintf("/api/contacts/channel/%s", channelID),
		forRequest(parameters),
		&response,
	)
	return response, err
}

func (s *ContactChannelsService) UpdateSubscription(channelID string, parameters map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.request(
		"PATCH",
		fmt.Sprintf("/api/contacts/channel/%s", channelID),
		forRequest(parameters),
		&response,
	)
	return response, err
}
