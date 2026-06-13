package reloop

type User struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Image *string `json:"image"`
	Email string  `json:"email"`
}

type ApiKey struct {
	ID                  string  `json:"id"`
	Name                *string `json:"name"`
	Start               *string `json:"start"`
	Prefix              *string `json:"prefix"`
	OrganizationID      string  `json:"organizationId"`
	UserID              string  `json:"userId"`
	RefillInterval      *int    `json:"refillInterval"`
	RefillAmount        *int    `json:"refillAmount"`
	LastRefillAt        *string `json:"lastRefillAt"`
	Enabled             bool    `json:"enabled"`
	RateLimitEnabled    bool    `json:"rateLimitEnabled"`
	RateLimitTimeWindow int     `json:"rateLimitTimeWindow"`
	RateLimitMax        int     `json:"rateLimitMax"`
	RequestCount        int     `json:"requestCount"`
	Remaining           *int    `json:"remaining"`
	LastRequest         *string `json:"lastRequest"`
	ExpiresAt           *string `json:"expiresAt"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
	Permissions         *string `json:"permissions"`
	Metadata            *string `json:"metadata"`
	CreatedBy           *User   `json:"createdBy,omitempty"`
	Object              string  `json:"object"`
	Event               string  `json:"event"`
}

type ApiKeyWithKey struct {
	ID          string  `json:"id"`
	Name        *string `json:"name"`
	Key         string  `json:"key"`
	Enabled     bool    `json:"enabled"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Permissions *string `json:"permissions"`
	Object      string  `json:"object"`
	Event       string  `json:"event"`
}

type ApiKeyListResponse struct {
	Object  string   `json:"object"`
	ApiKeys []ApiKey `json:"apiKeys"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
	Event   string   `json:"event"`
}

type ApiKeyListParams struct {
	Page    *int
	Limit   *int
	Enabled *bool
	UserID  *string
	Q       *string
}

type DeleteApiKeyResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Object  string `json:"object"`
	Event   string `json:"event"`
}

type CreateApiKeyParams struct {
	Name             string `json:"name"`
	Enabled          *bool  `json:"enabled,omitempty"`
	RateLimitEnabled *bool  `json:"rateLimitEnabled,omitempty"`
}

type UpdateApiKeyParams struct {
	Name    string `json:"name,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type DomainStatus string

const (
	DomainStatusPending   DomainStatus = "pending"
	DomainStatusVerifying DomainStatus = "verifying"
	DomainStatusActive    DomainStatus = "active"
	DomainStatusSuspended DomainStatus = "suspended"
	DomainStatusFailed    DomainStatus = "failed"
)

type DomainTLSMode string

const (
	DomainTLSOpportunistic DomainTLSMode = "opportunistic"
	DomainTLSEnforced      DomainTLSMode = "enforced"
)

type DNSRecordPurpose string

const (
	DNSRecordPurposeSending   DNSRecordPurpose = "sending"
	DNSRecordPurposeReceiving DNSRecordPurpose = "receiving"
	DNSRecordPurposeTracking  DNSRecordPurpose = "tracking"
)

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

type CreateDomainParams struct {
	Domain           string         `json:"domain"`
	CustomReturnPath *string        `json:"custom_return_path,omitempty"`
	Tracking         *string        `json:"tracking,omitempty"`
	ClickTracking    *bool          `json:"click_tracking,omitempty"`
	OpenTracking     *bool          `json:"open_tracking,omitempty"`
	TLS              *DomainTLSMode `json:"tls,omitempty"`
	SendingEmail     *bool          `json:"sending_email,omitempty"`
	ReceivingEmail   *bool          `json:"receiving_email,omitempty"`
}

type UpdateDomainParams struct {
	ClickTracking  *bool          `json:"click_tracking,omitempty"`
	OpenTracking   *bool          `json:"open_tracking,omitempty"`
	SendingEmail   *bool          `json:"sending_email,omitempty"`
	ReceivingEmail *bool          `json:"receiving_email,omitempty"`
	TLS            *DomainTLSMode `json:"tls,omitempty"`
}

type ListDomainsParams struct {
	Page   *int
	Limit  *int
	Q      *string
	Status *DomainStatus
}

type DomainListResponse struct {
	Object  string   `json:"object"`
	Domains []Domain `json:"domains"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
	Event   string   `json:"event"`
}

type DomainStatusResponse struct {
	ID     string       `json:"id"`
	Status DomainStatus `json:"status"`
	Event  string       `json:"event,omitempty"`
}

type ForwardDNSParams struct {
	Email string `json:"email"`
}

type ForwardDNSResponse struct {
	Success bool `json:"success"`
}

type DomainNameserversResponse struct {
	Object      string   `json:"object"`
	DomainID    string   `json:"domainId"`
	Domain      string   `json:"domain"`
	Nameservers []string `json:"nameservers"`
	DNSProvider *string  `json:"dnsProvider"`
	Event       string   `json:"event"`
}
