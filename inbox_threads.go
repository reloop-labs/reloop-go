package reloop

import (
	"net/url"
	"strconv"
)

// InboxThreadsService manages inbox threads.
type InboxThreadsService struct {
	client *Client
}

// ThreadStatus is a thread lifecycle status.
type ThreadStatus string

const (
	ThreadStatusActive   ThreadStatus = "active"
	ThreadStatusArchived ThreadStatus = "archived"
	ThreadStatusClosed   ThreadStatus = "closed"
	ThreadStatusTrash    ThreadStatus = "trash"
)

// ThreadFilter filters thread list results.
type ThreadFilter string

const (
	ThreadFilterPrimary ThreadFilter = "primary"
	ThreadFilterAlerts  ThreadFilter = "alerts"
	ThreadFilterPerson  ThreadFilter = "person"
	ThreadFilterTag     ThreadFilter = "tag"
)

// ThreadBatchAction is a batch operation on threads.
type ThreadBatchAction string

const (
	ThreadBatchActionArchive     ThreadBatchAction = "archive"
	ThreadBatchActionTrash       ThreadBatchAction = "trash"
	ThreadBatchActionRestore     ThreadBatchAction = "restore"
	ThreadBatchActionStar        ThreadBatchAction = "star"
	ThreadBatchActionUnstar      ThreadBatchAction = "unstar"
	ThreadBatchActionRead        ThreadBatchAction = "read"
	ThreadBatchActionUnread      ThreadBatchAction = "unread"
	ThreadBatchActionImportant   ThreadBatchAction = "important"
	ThreadBatchActionUnimportant ThreadBatchAction = "unimportant"
	ThreadBatchActionSpam        ThreadBatchAction = "spam"
	ThreadBatchActionUnspam      ThreadBatchAction = "unspam"
	ThreadBatchActionPin         ThreadBatchAction = "pin"
	ThreadBatchActionUnpin       ThreadBatchAction = "unpin"
)

// ThreadLabel is a label on a thread.
type ThreadLabel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Thread is an inbox thread resource.
type Thread struct {
	ID                 string        `json:"id"`
	MailboxID          *string       `json:"mailboxId"`
	OrganizationID     string        `json:"organizationId"`
	Subject            *string       `json:"subject"`
	LastMessagePreview *string       `json:"lastMessagePreview"`
	LastMessageAt      string        `json:"lastMessageAt"`
	Status             string        `json:"status"`
	MessageCount       int           `json:"messageCount"`
	Participants       []string      `json:"participants"`
	IsRead             bool          `json:"isRead"`
	IsStarred          bool          `json:"isStarred"`
	IsImportant        *bool         `json:"isImportant,omitempty"`
	IsPinned           *bool         `json:"isPinned,omitempty"`
	PinnedAt           *string       `json:"pinnedAt,omitempty"`
	Labels             []ThreadLabel `json:"labels,omitempty"`
	DeletedAt          *string       `json:"deletedAt,omitempty"`
	CreatedAt          string        `json:"createdAt"`
	UpdatedAt          string        `json:"updatedAt"`
}

// ThreadMessage is a message within a thread detail response.
type ThreadMessage struct {
	ID              string      `json:"id"`
	ThreadID        string      `json:"threadId"`
	Direction       string      `json:"direction"`
	InboundEmailID  *string     `json:"inboundEmailId"`
	EmailLogID      *string     `json:"emailLogId"`
	FromEmail       string      `json:"fromEmail"`
	FromName        *string     `json:"fromName"`
	Subject         *string     `json:"subject"`
	Preview         *string     `json:"preview"`
	MessageAt       string      `json:"messageAt"`
	RFC822MessageID *string     `json:"rfc822MessageId"`
	InReplyTo       *string     `json:"inReplyTo"`
	CreatedAt       string      `json:"createdAt"`
	Email           interface{} `json:"email"`
}

// ThreadDetail is a thread with its messages.
type ThreadDetail struct {
	Thread
	Messages []ThreadMessage `json:"messages"`
}

// ListThreadsParams are query params for List.
type ListThreadsParams struct {
	MailboxID *string
	Limit     *int
	Offset    *int
	Folder    *string
	Q         *string
	IsPinned  *bool
	Filter    *ThreadFilter
}

// BatchThreadsParams batch-updates threads.
type BatchThreadsParams struct {
	IDs    []string          `json:"ids"`
	Action ThreadBatchAction `json:"action"`
}

// UpdateThreadParams updates thread flags.
type UpdateThreadParams struct {
	IsRead      *bool         `json:"isRead,omitempty"`
	IsStarred   *bool         `json:"isStarred,omitempty"`
	IsImportant *bool         `json:"isImportant,omitempty"`
	IsPinned    *bool         `json:"isPinned,omitempty"`
	Status      *ThreadStatus `json:"status,omitempty"`
}

// SetThreadReadParams sets read state.
type SetThreadReadParams struct {
	IsRead bool `json:"isRead"`
}

// SetThreadStarParams sets starred state.
type SetThreadStarParams struct {
	IsStarred bool `json:"isStarred"`
}

// List lists inbox threads.
func (s *InboxThreadsService) List(params *ListThreadsParams) ([]Thread, error) {
	q, err := buildListThreadsQuery(params)
	if err != nil {
		return nil, err
	}
	var out []Thread
	if err := s.client.do("GET", threadsV1, nil, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Batch performs a batch action on threads.
func (s *InboxThreadsService) Batch(params BatchThreadsParams) (*ThreadBatchResponse, error) {
	ids, err := requireInboxIDArray(params.IDs, "ids", inboxBatchIDsMax)
	if err != nil {
		return nil, err
	}
	if err := requireThreadBatchAction(params.Action); err != nil {
		return nil, err
	}
	body := BatchThreadsParams{IDs: ids, Action: params.Action}
	var out ThreadBatchResponse
	if err := s.client.do("POST", threadsV1+"/batch", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get fetches a thread by ID.
func (s *InboxThreadsService) Get(id string) (*ThreadDetail, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	var out ThreadDetail
	if err := s.client.do("GET", threadsV1+"/"+threadID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAttachment fetches a thread attachment.
func (s *InboxThreadsService) GetAttachment(id, attachmentID string) (*MessageAttachment, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	attachment, err := requireInboxAttachmentID(attachmentID, "attachmentId")
	if err != nil {
		return nil, err
	}
	var out MessageAttachment
	path := threadsV1 + "/" + threadID + "/attachments/" + attachment
	if err := s.client.do("GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates thread flags.
func (s *InboxThreadsService) Update(id string, params UpdateThreadParams) (*InboxSuccessResponse, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	if params.IsRead == nil && params.IsStarred == nil && params.IsImportant == nil &&
		params.IsPinned == nil && params.Status == nil {
		return nil, validationErrorf("params", "update requires at least one of isRead, isStarred, isImportant, isPinned, or status.")
	}
	if params.Status != nil {
		if err := requireThreadStatusValue(*params.Status, "status"); err != nil {
			return nil, err
		}
	}
	var out InboxSuccessResponse
	if err := s.client.do("PATCH", threadsV1+"/"+threadID, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetRead sets the read state of a thread.
func (s *InboxThreadsService) SetRead(id string, params SetThreadReadParams) (*InboxSuccessResponse, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	body := SetThreadReadParams{IsRead: params.IsRead}
	var out InboxSuccessResponse
	if err := s.client.do("PATCH", threadsV1+"/"+threadID+"/read", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetStar sets the starred state of a thread.
func (s *InboxThreadsService) SetStar(id string, params SetThreadStarParams) (*InboxSuccessResponse, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	body := SetThreadStarParams{IsStarred: params.IsStarred}
	var out InboxSuccessResponse
	if err := s.client.do("PATCH", threadsV1+"/"+threadID+"/star", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Archive archives a thread.
func (s *InboxThreadsService) Archive(id string) (*InboxSuccessResponse, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	var out InboxSuccessResponse
	if err := s.client.do("POST", threadsV1+"/"+threadID+"/archive", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Trash moves a thread to trash.
func (s *InboxThreadsService) Trash(id string) (*InboxSuccessResponse, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	var out InboxSuccessResponse
	if err := s.client.do("POST", threadsV1+"/"+threadID+"/trash", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Restore restores a thread from trash or archive.
func (s *InboxThreadsService) Restore(id string) (*InboxSuccessResponse, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	var out InboxSuccessResponse
	if err := s.client.do("POST", threadsV1+"/"+threadID+"/restore", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete permanently deletes a thread.
func (s *InboxThreadsService) Delete(id string) (*InboxSuccessResponse, error) {
	threadID, err := requireThreadID(id, "id")
	if err != nil {
		return nil, err
	}
	var out InboxSuccessResponse
	if err := s.client.do("DELETE", threadsV1+"/"+threadID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func buildListThreadsQuery(params *ListThreadsParams) (url.Values, error) {
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
	if params.Folder != nil {
		q.Set("folder", *params.Folder)
	}
	if params.Q != nil {
		q.Set("q", *params.Q)
	}
	if params.IsPinned != nil {
		q.Set("isPinned", strconv.FormatBool(*params.IsPinned))
	}
	if params.Filter != nil {
		if err := requireThreadFilterValue(*params.Filter, "filter"); err != nil {
			return nil, err
		}
		q.Set("filter", string(*params.Filter))
	}
	return q, nil
}
