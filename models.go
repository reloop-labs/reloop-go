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
	Name string `json:"name"`
}

type UpdateApiKeyParams struct {
	Name string `json:"name"`
}
