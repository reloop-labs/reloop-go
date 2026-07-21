package reloop

// InboxMailboxesService manages inbox mailboxes.
type InboxMailboxesService struct {
	client *Client
}

// MailboxStatus is a mailbox lifecycle status.
type MailboxStatus string

const (
	MailboxStatusActive   MailboxStatus = "active"
	MailboxStatusDisabled MailboxStatus = "disabled"
)

// Mailbox is a mailbox resource.
type Mailbox struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Quota       string  `json:"quota"`
	Status      string  `json:"status"`
	DisplayName *string `json:"displayName"`
	CreatedAt   string  `json:"createdAt"`
}

// MailboxDetail is a mailbox with extended fields.
type MailboxDetail struct {
	Mailbox
	DomainID  string `json:"domainId"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateMailboxResponse is returned by Create.
type CreateMailboxResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// CreateMailboxParams creates a mailbox.
type CreateMailboxParams struct {
	DomainID    string  `json:"domainId"`
	Email       string  `json:"email"`
	Password    *string `json:"password,omitempty"`
	Quota       *string `json:"quota,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
}

// UpdateMailboxParams updates a mailbox.
type UpdateMailboxParams struct {
	DisplayName *string        `json:"displayName,omitempty"`
	Status      *MailboxStatus `json:"status,omitempty"`
	Quota       *string        `json:"quota,omitempty"`
}

// List returns all mailboxes.
func (s *InboxMailboxesService) List() ([]Mailbox, error) {
	var out []Mailbox
	if err := s.client.do("GET", mailboxesV1+"/list", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a mailbox by ID.
func (s *InboxMailboxesService) Get(id string) (*MailboxDetail, error) {
	mailboxID, err := requireMailboxID(id, "id")
	if err != nil {
		return nil, err
	}
	var out MailboxDetail
	if err := s.client.do("GET", mailboxesV1+"/"+mailboxID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a mailbox.
func (s *InboxMailboxesService) Create(params CreateMailboxParams) (*CreateMailboxResponse, error) {
	domainID, err := requireNonEmptyString(params.DomainID, "domainId")
	if err != nil {
		return nil, err
	}
	email, err := requireNonEmptyString(params.Email, "email")
	if err != nil {
		return nil, err
	}
	body := CreateMailboxParams{
		DomainID:    domainID,
		Email:       email,
		Password:    params.Password,
		Quota:       params.Quota,
		DisplayName: params.DisplayName,
	}
	var out CreateMailboxResponse
	if err := s.client.do("POST", mailboxesV1+"/create", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a mailbox.
func (s *InboxMailboxesService) Update(id string, params UpdateMailboxParams) (*InboxSuccessResponse, error) {
	mailboxID, err := requireMailboxID(id, "id")
	if err != nil {
		return nil, err
	}
	if params.DisplayName == nil && params.Status == nil && params.Quota == nil {
		return nil, validationErrorf("params", "update requires at least one of displayName, status, or quota.")
	}
	if params.Status != nil {
		if err := requireMailboxStatus(*params.Status, "status"); err != nil {
			return nil, err
		}
	}
	var out InboxSuccessResponse
	if err := s.client.do("PATCH", mailboxesV1+"/"+mailboxID, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a mailbox.
func (s *InboxMailboxesService) Delete(id string) (*InboxSuccessResponse, error) {
	mailboxID, err := requireMailboxID(id, "id")
	if err != nil {
		return nil, err
	}
	var out InboxSuccessResponse
	if err := s.client.do("DELETE", mailboxesV1+"/"+mailboxID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
