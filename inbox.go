package reloop

const (
	inboxV1     = "/api/inbox/v1"
	mailboxesV1 = inboxV1 + "/mailboxes"
	messagesV1  = inboxV1 + "/messages"
	threadsV1   = inboxV1 + "/threads"
)

// InboxService provides inbox mailboxes, messages, and threads.
type InboxService struct {
	Mailboxes *InboxMailboxesService
	Messages  *InboxMessagesService
	Threads   *InboxThreadsService
}

func newInboxService(c *Client) *InboxService {
	return &InboxService{
		Mailboxes: &InboxMailboxesService{client: c},
		Messages:  &InboxMessagesService{client: c},
		Threads:   &InboxThreadsService{client: c},
	}
}

// InboxSuccessResponse is returned by many inbox mutation endpoints.
type InboxSuccessResponse struct {
	Success     bool    `json:"success"`
	ID          *string `json:"id,omitempty"`
	Message     *string `json:"message,omitempty"`
	IsRead      *bool   `json:"isRead,omitempty"`
	IsStarred   *bool   `json:"isStarred,omitempty"`
	IsSpam      *bool   `json:"isSpam,omitempty"`
	IsImportant *bool   `json:"isImportant,omitempty"`
	IsPinned    *bool   `json:"isPinned,omitempty"`
	Status      *string `json:"status,omitempty"`
	DeletedAt   *string `json:"deletedAt,omitempty"`
}

// SendEmailResponse is returned when a message is sent immediately.
type SendEmailResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	ID        string `json:"id"`
}

// PendingSendResponse is returned when a send is scheduled with undo window.
type PendingSendResponse struct {
	Success   bool   `json:"success"`
	Pending   bool   `json:"pending"`
	ID        string `json:"id"`
	SendAt    string `json:"sendAt"`
	MessageID string `json:"messageId"`
}

// SendEmailOrPendingResponse is the union of immediate and pending send responses.
type SendEmailOrPendingResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	ID        string `json:"id"`
	Pending   *bool  `json:"pending,omitempty"`
	SendAt    string `json:"sendAt,omitempty"`
}

// ThreadBatchResponse is returned by thread batch operations.
type ThreadBatchResponse struct {
	Success bool     `json:"success"`
	IDs     []string `json:"ids"`
	Action  string   `json:"action"`
}

// MessageAttachment is an inbox attachment resource.
type MessageAttachment struct {
	ID                 string  `json:"id"`
	Filename           string  `json:"filename"`
	ContentType        string  `json:"contentType"`
	Size               int     `json:"size"`
	StoragePath        string  `json:"storagePath"`
	ContentDisposition *string `json:"contentDisposition"`
	ContentID          *string `json:"contentId"`
	CreatedAt          string  `json:"createdAt"`
}

// AttachmentInput is an attachment on a compose/send request.
type AttachmentInput struct {
	Content     *string `json:"content,omitempty"`
	Filename    *string `json:"filename,omitempty"`
	Path        *string `json:"path,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
	ContentID   *string `json:"content_id,omitempty"`
}
