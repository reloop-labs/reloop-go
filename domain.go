package reloop

import (
	"net/url"
	"strconv"
)

// DomainService manages sending/receiving domains.
type DomainService struct {
	client *Client
}

// DomainStatus is a domain lifecycle status.
type DomainStatus string

const (
	DomainStatusPending   DomainStatus = "pending"
	DomainStatusVerifying DomainStatus = "verifying"
	DomainStatusActive    DomainStatus = "active"
	DomainStatusSuspended DomainStatus = "suspended"
	DomainStatusFailed    DomainStatus = "failed"
)

// DomainTLSMode is TLS enforcement mode.
type DomainTLSMode string

const (
	DomainTLSOpportunistic DomainTLSMode = "opportunistic"
	DomainTLSEnforced      DomainTLSMode = "enforced"
)

// DNSRecordPurpose classifies a DNS record.
type DNSRecordPurpose string

const (
	DNSRecordPurposeSending   DNSRecordPurpose = "sending"
	DNSRecordPurposeReceiving DNSRecordPurpose = "receiving"
	DNSRecordPurposeTracking  DNSRecordPurpose = "tracking"
)

// DNSRecord is a domain DNS record.
type DNSRecord struct {
	ID                string           `json:"id"`
	RecordType        string           `json:"recordType"`
	RecordTypeName    string           `json:"recordTypeName"`
	Domain            string           `json:"domain"`
	Name              string           `json:"name"`
	Value             string           `json:"value"`
	TTL               string           `json:"ttl"`
	Priority          *int             `json:"priority"`
	VerificationError *string          `json:"verificationError"`
	Purpose           DNSRecordPurpose `json:"purpose,omitempty"`
	CreatedAt         string           `json:"createdAt"`
	Status            DomainStatus     `json:"status"`
	UpdatedAt         string           `json:"updatedAt"`
}

// Domain is a domain resource.
type Domain struct {
	Object                   string        `json:"object"`
	ID                       string        `json:"id"`
	Domain                   string        `json:"domain"`
	Status                   DomainStatus  `json:"status"`
	UserVerifiedDomain       bool          `json:"userVerifiedDomain"`
	SystemVerified           bool          `json:"systemVerified"`
	CustomReturnPath         string        `json:"customReturnPath"`
	TrackingSubdomain        string        `json:"trackingSubdomain"`
	IsClickTrackingEnabled   bool          `json:"isClickTrackingEnabled"`
	IsOpenTrackingEnabled    bool          `json:"isOpenTrackingEnabled"`
	TLS                      DomainTLSMode `json:"tls"`
	IsTrackingDomain         bool          `json:"isTrackingDomain"`
	IsSendingEmailEnabled    bool          `json:"isSendingEmailEnabled"`
	IsReceivingEmailEnabled  bool          `json:"isReceivingEmailEnabled"`
	VerificationFailedReason *string       `json:"verificationFailedReason"`
	DNSRecords               []DNSRecord   `json:"dnsRecords"`
	LastVerifiedAt           *string       `json:"lastVerifiedAt"`
	CreatedAt                string        `json:"createdAt"`
	UpdatedAt                string        `json:"updatedAt"`
	Event                    string        `json:"event,omitempty"`
}

// CreateDomainParams creates a domain.
type CreateDomainParams struct {
	Domain         string         `json:"domain"`
	ClickTracking  *bool          `json:"click_tracking,omitempty"`
	OpenTracking   *bool          `json:"open_tracking,omitempty"`
	TLS            *DomainTLSMode `json:"tls,omitempty"`
	SendingEmail   *bool          `json:"sending_email,omitempty"`
	ReceivingEmail *bool          `json:"receiving_email,omitempty"`
}

// UpdateDomainParams updates a domain.
type UpdateDomainParams struct {
	ClickTracking  *bool          `json:"click_tracking,omitempty"`
	OpenTracking   *bool          `json:"open_tracking,omitempty"`
	SendingEmail   *bool          `json:"sending_email,omitempty"`
	ReceivingEmail *bool          `json:"receiving_email,omitempty"`
	TLS            *DomainTLSMode `json:"tls,omitempty"`
}

// ListDomainsParams are query params for List.
type ListDomainsParams struct {
	Page   *int
	Limit  *int
	Q      *string
	Status *DomainStatus
}

// DomainListResponse is a paginated domain list.
type DomainListResponse struct {
	Object  string   `json:"object"`
	Domains []Domain `json:"domains"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
	Event   string   `json:"event"`
}

// DomainStatusResponse is returned by Verify.
type DomainStatusResponse struct {
	ID     string       `json:"id"`
	Status DomainStatus `json:"status"`
	Event  string       `json:"event,omitempty"`
}

func requireDomainID(id, field string) (string, error) {
	return requireNonEmptyString(id, field)
}

// Create creates a domain.
func (s *DomainService) Create(params CreateDomainParams) (*Domain, error) {
	domain, err := requireNonEmptyString(params.Domain, "domain")
	if err != nil {
		return nil, err
	}
	body := CreateDomainParams{
		Domain:         domain,
		ClickTracking:  params.ClickTracking,
		OpenTracking:   params.OpenTracking,
		TLS:            params.TLS,
		SendingEmail:   params.SendingEmail,
		ReceivingEmail: params.ReceivingEmail,
	}
	var out Domain
	if err := s.client.do("POST", domainV1+"/create", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

const domainV1 = "/api/domain/v1"

// List lists domains.
func (s *DomainService) List(params *ListDomainsParams) (*DomainListResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Page != nil {
			if err := requirePage(*params.Page, "page"); err != nil {
				return nil, err
			}
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.Limit != nil {
			if err := requireLimit(*params.Limit, 1, 100, "limit"); err != nil {
				return nil, err
			}
			q.Set("limit", strconv.Itoa(*params.Limit))
		}
		if params.Q != nil {
			q.Set("q", *params.Q)
		}
		if params.Status != nil {
			q.Set("status", string(*params.Status))
		}
	}
	var out DomainListResponse
	if err := s.client.do("GET", domainV1+"/list", nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get fetches a domain.
func (s *DomainService) Get(domainID string) (*Domain, error) {
	id, err := requireDomainID(domainID, "domainId")
	if err != nil {
		return nil, err
	}
	var out Domain
	if err := s.client.do("GET", domainV1+"/"+id, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a domain.
func (s *DomainService) Update(domainID string, params UpdateDomainParams) (*Domain, error) {
	id, err := requireDomainID(domainID, "domainId")
	if err != nil {
		return nil, err
	}
	var out Domain
	if err := s.client.do("PATCH", domainV1+"/"+id, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a domain.
func (s *DomainService) Delete(domainID string) error {
	id, err := requireDomainID(domainID, "domainId")
	if err != nil {
		return err
	}
	return s.client.do("DELETE", domainV1+"/"+id, nil, nil, nil)
}

// Verify triggers domain verification.
func (s *DomainService) Verify(domainID string) (*DomainStatusResponse, error) {
	id, err := requireDomainID(domainID, "domainId")
	if err != nil {
		return nil, err
	}
	var out DomainStatusResponse
	if err := s.client.do("POST", domainV1+"/verify/"+id, map[string]interface{}{}, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
