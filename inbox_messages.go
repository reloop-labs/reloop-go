package reloop

import (
	"net/url"
	"strconv"
)

// InboxMessagesService manages inbox messages.
type InboxMessagesService struct {
	client *Client
}

// MessageAttachmentItem is an attachment on a message.
type MessageAttachmentItem struct {
	ID                 string  `json:"id"`
	InboundEmailID     string  `json:"inboundEmailId"`
	Filename           string  `json:"filename"`
	ContentType        string  `json:"contentType"`
	Size               int     `json:"size"`
	StoragePath        string  `json:"storagePath"`
	ContentDisposition *string `json:"contentDisposition"`
	ContentID          *string `json:"contentId"`
	CreatedAt          string  `json:"createdAt"`
}

// Message is an inbox message resource.
type Message struct {
	ID             string                  `json:"id"`
	MailboxID      string                  `json:"mailboxId"`
	OrganizationID string                  `json:"organizationId"`
	FromEmail      string                  `json:"fromEmail"`
	FromName       *string                 `json:"fromName"`
	ToEmails       []string                `json:"toEmails"`
	CcEmails       []string                `json:"ccEmails,omitempty"`
	BccEmails      []string                `json:"bccEmails,omitempty"`
	ReplyTo        *string                 `json:"replyTo"`
	Subject        *string                 `json:"subject"`
	TextBody       *string                 `json:"textBody"`
	HTMLBody       *string                 `json:"htmlBody"`
	Snippet        *string                 `json:"snippet"`
	Size           int                     `json:"size"`
	Status         string                  `json:"status"`
	IsRead         bool                    `json:"isRead"`
	IsStarred      bool                    `json:"isStarred"`
	IsSpam         bool                    `json:"isSpam"`
	SpamScore      *float64                `json:"spamScore"`
	MessageID      *string                 `json:"messageId"`
	ThreadID       *string                 `json:"threadId"`
	InReplyTo      *string                 `json:"inReplyTo"`
	References     []string                `json:"references,omitempty"`
	Headers        map[string]string       `json:"headers,omitempty"`
	Date           *string                 `json:"date"`
	CreatedAt      string                  `json:"createdAt"`
	Attachments    []MessageAttachmentItem `json:"attachments,omitempty"`
}

// MessageRaw is the raw RFC822 content of a message.
type MessageRaw struct {
	ID        string  `json:"id"`
	MessageID *string `json:"messageId"`
	Raw       string  `json:"raw"`
}

// ListMessagesParams are query params for List.
type ListMessagesParams struct {
	MailboxID *string
	Limit     *int
	Offset    *int
	Q         *string
	IsSpam    *bool
}

// ListSentMessagesParams are query params for ListSent.
type ListSentMessagesParams struct {
	MailboxID *string
}

// BatchMessagesParams batch-fetches messages by ID.
type BatchMessagesParams struct {
	IDs []string `json:"ids"`
}

// UpdateMessageParams updates message flags.
type UpdateMessageParams struct {
	IsRead    *bool `json:"isRead,omitempty"`
	IsStarred *bool `json:"isStarred,omitempty"`
	IsSpam    *bool `json:"isSpam,omitempty"`
}

// SetMessageReadParams sets read state.
type SetMessageReadParams struct {
	IsRead bool `json:"isRead"`
}

// SetMessageStarParams sets starred state.
type SetMessageStarParams struct {
	IsStarred bool `json:"isStarred"`
}

// SendMessageParams sends a new message.
type SendMessageParams struct {
	MailboxID         string            `json:"mailboxId"`
	To                interface{}       `json:"to"`
	Subject           string            `json:"subject"`
	Text              *string           `json:"text,omitempty"`
	HTML              *string           `json:"html,omitempty"`
	Cc                interface{}       `json:"cc,omitempty"`
	Bcc               interface{}       `json:"bcc,omitempty"`
	Attachments       []AttachmentInput `json:"attachments,omitempty"`
	ScheduledAt       *string           `json:"scheduledAt,omitempty"`
	UndoWindowSeconds *float64          `json:"undoWindowSeconds,omitempty"`
}

// ComposeMessageParams are shared by reply, reply-all, and forward body fields.
type ComposeMessageParams struct {
	Text        *string           `json:"text,omitempty"`
	HTML        *string           `json:"html,omitempty"`
	Cc          interface{}       `json:"cc,omitempty"`
	Bcc         interface{}       `json:"bcc,omitempty"`
	Attachments []AttachmentInput `json:"attachments,omitempty"`
}

// ForwardMessageParams forwards a message.
type ForwardMessageParams struct {
	ComposeMessageParams
	To interface{} `json:"to"`
}

// List lists inbox messages.
func (s *InboxMessagesService) List(params *ListMessagesParams) ([]Message, error) {
	q, err := buildListMessagesQuery(params)
	if err != nil {
		return nil, err
	}
	var out []Message
	if err := s.client.do("GET", messagesV1, nil, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListSent lists sent messages.
func (s *InboxMessagesService) ListSent(params *ListSentMessagesParams) ([]Message, error) {
	q := url.Values{}
	if params != nil && params.MailboxID != nil {
		q.Set("mailboxId", *params.MailboxID)
	}
	var out []Message
	if err := s.client.do("GET", messagesV1+"/sent", nil, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a message by ID.
func (s *InboxMessagesService) Get(id string) (*Message, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	var out Message
	if err := s.client.do("GET", messagesV1+"/"+messageID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Batch batch-fetches messages.
func (s *InboxMessagesService) Batch(params BatchMessagesParams) ([]Message, error) {
	ids, err := requireInboxIDArray(params.IDs, "ids", inboxBatchIDsMax)
	if err != nil {
		return nil, err
	}
	body := BatchMessagesParams{IDs: ids}
	var out []Message
	if err := s.client.do("POST", messagesV1+"/batch", body, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetRaw fetches raw RFC822 content.
func (s *InboxMessagesService) GetRaw(id string) (*MessageRaw, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	var out MessageRaw
	if err := s.client.do("GET", messagesV1+"/"+messageID+"/raw", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAttachment fetches a message attachment.
func (s *InboxMessagesService) GetAttachment(id, attachmentID string) (*MessageAttachment, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	attachment, err := requireInboxAttachmentID(attachmentID, "attachmentId")
	if err != nil {
		return nil, err
	}
	var out MessageAttachment
	path := messagesV1 + "/" + messageID + "/attachments/" + attachment
	if err := s.client.do("GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates message flags.
func (s *InboxMessagesService) Update(id string, params UpdateMessageParams) (*InboxSuccessResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	if params.IsRead == nil && params.IsStarred == nil && params.IsSpam == nil {
		return nil, validationErrorf("params", "update requires at least one of isRead, isStarred, or isSpam.")
	}
	var out InboxSuccessResponse
	if err := s.client.do("PATCH", messagesV1+"/"+messageID, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetRead sets the read state of a message.
func (s *InboxMessagesService) SetRead(id string, params SetMessageReadParams) (*InboxSuccessResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	body := SetMessageReadParams{IsRead: params.IsRead}
	var out InboxSuccessResponse
	if err := s.client.do("PATCH", messagesV1+"/"+messageID+"/read", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetStar sets the starred state of a message.
func (s *InboxMessagesService) SetStar(id string, params SetMessageStarParams) (*InboxSuccessResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	body := SetMessageStarParams{IsStarred: params.IsStarred}
	var out InboxSuccessResponse
	if err := s.client.do("PATCH", messagesV1+"/"+messageID+"/star", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a message.
func (s *InboxMessagesService) Delete(id string) (*InboxSuccessResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	var out InboxSuccessResponse
	if err := s.client.do("DELETE", messagesV1+"/"+messageID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Send sends a new message.
func (s *InboxMessagesService) Send(params SendMessageParams) (*SendEmailOrPendingResponse, error) {
	mailboxID, err := requireNonEmptyString(params.MailboxID, "mailboxId")
	if err != nil {
		return nil, err
	}
	to, err := requireRecipient(params.To, "to")
	if err != nil {
		return nil, err
	}
	subject, err := requireNonEmptyString(params.Subject, "subject")
	if err != nil {
		return nil, err
	}
	if params.UndoWindowSeconds != nil {
		if err := requireFiniteNumber(*params.UndoWindowSeconds, "undoWindowSeconds"); err != nil {
			return nil, err
		}
	}
	body := SendMessageParams{
		MailboxID:         mailboxID,
		To:                to,
		Subject:           subject,
		Text:              params.Text,
		HTML:              params.HTML,
		ScheduledAt:       params.ScheduledAt,
		UndoWindowSeconds: params.UndoWindowSeconds,
		Attachments:       params.Attachments,
	}
	if params.Cc != nil {
		cc, err := requireRecipient(params.Cc, "cc")
		if err != nil {
			return nil, err
		}
		body.Cc = cc
	}
	if params.Bcc != nil {
		bcc, err := requireRecipient(params.Bcc, "bcc")
		if err != nil {
			return nil, err
		}
		body.Bcc = bcc
	}
	var out SendEmailOrPendingResponse
	if err := s.client.do("POST", messagesV1+"/send", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelPending cancels a pending send.
func (s *InboxMessagesService) CancelPending(id string) (*InboxSuccessResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	var out InboxSuccessResponse
	path := messagesV1 + "/pending/" + messageID + "/cancel"
	if err := s.client.do("POST", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reply replies to a message.
func (s *InboxMessagesService) Reply(id string, params ComposeMessageParams) (*SendEmailResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := buildComposeBody(params)
	if err != nil {
		return nil, err
	}
	var out SendEmailResponse
	if err := s.client.do("POST", messagesV1+"/"+messageID+"/reply", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReplyAll reply-all to a message.
func (s *InboxMessagesService) ReplyAll(id string, params ComposeMessageParams) (*SendEmailResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	body, err := buildComposeBody(params)
	if err != nil {
		return nil, err
	}
	var out SendEmailResponse
	if err := s.client.do("POST", messagesV1+"/"+messageID+"/reply-all", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Forward forwards a message.
func (s *InboxMessagesService) Forward(id string, params ForwardMessageParams) (*SendEmailResponse, error) {
	messageID, err := requireMessageID(id, "id")
	if err != nil {
		return nil, err
	}
	to, err := requireRecipient(params.To, "to")
	if err != nil {
		return nil, err
	}
	compose, err := buildComposeBody(params.ComposeMessageParams)
	if err != nil {
		return nil, err
	}
	compose["to"] = to
	var out SendEmailResponse
	if err := s.client.do("POST", messagesV1+"/"+messageID+"/forward", compose, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func buildListMessagesQuery(params *ListMessagesParams) (url.Values, error) {
	q := url.Values{}
	if params == nil {
		return q, nil
	}
	if params.MailboxID != nil {
		q.Set("mailboxId", *params.MailboxID)
	}
	if params.Limit != nil {
		if err := requireInboxLimit(*params.Limit, "limit"); err != nil {
			return nil, err
		}
		q.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		if err := requireInboxOffset(*params.Offset, "offset"); err != nil {
			return nil, err
		}
		q.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.Q != nil {
		q.Set("q", *params.Q)
	}
	if params.IsSpam != nil {
		q.Set("isSpam", strconv.FormatBool(*params.IsSpam))
	}
	return q, nil
}

func buildComposeBody(params ComposeMessageParams) (map[string]interface{}, error) {
	if err := requireComposeBody(params.Text, params.HTML); err != nil {
		return nil, err
	}
	body := map[string]interface{}{}
	if params.Text != nil {
		body["text"] = *params.Text
	}
	if params.HTML != nil {
		body["html"] = *params.HTML
	}
	if params.Cc != nil {
		cc, err := requireRecipient(params.Cc, "cc")
		if err != nil {
			return nil, err
		}
		body["cc"] = cc
	}
	if params.Bcc != nil {
		bcc, err := requireRecipient(params.Bcc, "bcc")
		if err != nil {
			return nil, err
		}
		body["bcc"] = bcc
	}
	if params.Attachments != nil {
		body["attachments"] = params.Attachments
	}
	return body, nil
}
