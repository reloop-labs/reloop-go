package reloop

type MailService struct {
	client *Client
}

func (s *MailService) Send(params SendMailParams) (*SendMailResponse, error) {
	var response SendMailResponse
	err := s.client.request("POST", "/api/mail/v1/send", params, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
