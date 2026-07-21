package reloop

// MailService sends transactional email.
type MailService struct {
	client *Client
}

// SendMailTag is a custom tag on a send.
type SendMailTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// SendMailAttachment is an email attachment.
type SendMailAttachment struct {
	Content     interface{} `json:"content,omitempty"`
	Filename    *string     `json:"filename,omitempty"`
	Path        *string     `json:"path,omitempty"`
	ContentType *string     `json:"content_type,omitempty"`
	ContentID   *string     `json:"content_id,omitempty"`
}

// SendMailTemplate references a template.
type SendMailTemplate struct {
	ID        string                 `json:"id"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// SendMailParams are parameters for Mail.Send.
type SendMailParams struct {
	From        string               `json:"from"`
	To          interface{}          `json:"to"`
	Subject     string               `json:"subject"`
	Cc          interface{}          `json:"cc,omitempty"`
	Bcc         interface{}          `json:"bcc,omitempty"`
	Text        *string              `json:"text,omitempty"`
	HTML        *string              `json:"html,omitempty"`
	ReplyTo     interface{}          `json:"reply_to,omitempty"`
	ScheduledAt *string              `json:"scheduled_at,omitempty"`
	Headers     map[string]string    `json:"headers,omitempty"`
	ChannelID   *string              `json:"channel_id,omitempty"`
	Attachments []SendMailAttachment `json:"attachments,omitempty"`
	Tags        []SendMailTag        `json:"tags,omitempty"`
	Template    *SendMailTemplate    `json:"template,omitempty"`
	ThreadID    *string              `json:"thread_id,omitempty"`
}

// SendMailResponse is returned by Mail.Send.
type SendMailResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	ID        string `json:"id"`
}

// Send sends an email.
func (s *MailService) Send(params SendMailParams) (*SendMailResponse, error) {
	from, err := requireMailString(params.From, "from")
	if err != nil {
		return nil, err
	}
	to, err := requireRecipient(params.To, "to")
	if err != nil {
		return nil, err
	}
	subject, err := requireMailString(params.Subject, "subject")
	if err != nil {
		return nil, err
	}

	body := SendMailParams{
		From:        from,
		To:          to,
		Subject:     subject,
		Cc:          params.Cc,
		Bcc:         params.Bcc,
		Text:        params.Text,
		HTML:        params.HTML,
		ReplyTo:     params.ReplyTo,
		ScheduledAt: params.ScheduledAt,
		Headers:     params.Headers,
		ChannelID:   params.ChannelID,
		Attachments: params.Attachments,
		Tags:        params.Tags,
		Template:    params.Template,
		ThreadID:    params.ThreadID,
	}

	var out SendMailResponse
	if err := s.client.do("POST", "/api/mail/v1/send", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
