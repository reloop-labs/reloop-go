package reloop

import (
	"net/url"
	"strconv"
)

const webhookV1 = "/api/webhook/v1"

// WebhookStatus is a webhook lifecycle status.
type WebhookStatus string

const (
	WebhookStatusActive   WebhookStatus = "active"
	WebhookStatusPaused   WebhookStatus = "paused"
	WebhookStatusDisabled WebhookStatus = "disabled"
	WebhookStatusFailed   WebhookStatus = "failed"
)

// WebhookDeliveryStatus is a delivery attempt status.
type WebhookDeliveryStatus string

const (
	WebhookDeliveryStatusPending  WebhookDeliveryStatus = "pending"
	WebhookDeliveryStatusSuccess  WebhookDeliveryStatus = "success"
	WebhookDeliveryStatusFailed   WebhookDeliveryStatus = "failed"
	WebhookDeliveryStatusRetrying WebhookDeliveryStatus = "retrying"
)

// WebhookService manages webhooks.
type WebhookService struct {
	client *Client
}

// Webhook is a webhook resource.
type Webhook struct {
	ID                     string                 `json:"id"`
	Name                   string                 `json:"name"`
	URL                    string                 `json:"url"`
	Secret                 string                 `json:"secret"`
	Status                 WebhookStatus          `json:"status"`
	CustomHeaders          map[string]string      `json:"customHeaders"`
	RateLimitEnabled       bool                   `json:"rateLimitEnabled"`
	MaxRequestsPerMinute   int                    `json:"maxRequestsPerMinute"`
	MaxRetries             int                    `json:"maxRetries"`
	RetryBackoffMultiplier float64                `json:"retryBackoffMultiplier"`
	FilteringOptions       map[string]interface{} `json:"filteringOptions"`
	LastTriggeredAt        *string                `json:"lastTriggeredAt"`
	SuccessCount           int                    `json:"successCount"`
	FailureCount           int                    `json:"failureCount"`
	ConsecutiveFailures    int                    `json:"consecutiveFailures"`
	Events                 []string               `json:"events"`
	CreatedAt              string                 `json:"createdAt"`
	UpdatedAt              string                 `json:"updatedAt"`
}

// CreateWebhookParams creates a webhook.
type CreateWebhookParams struct {
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Events      []string `json:"events"`
}

// UpdateWebhookParams updates a webhook.
type UpdateWebhookParams struct {
	Description            *string                `json:"description,omitempty"`
	Name                   *string                `json:"name,omitempty"`
	URL                    *string                `json:"url,omitempty"`
	Secret                 *string                `json:"secret,omitempty"`
	Status                 *WebhookStatus         `json:"status,omitempty"`
	CustomHeaders          map[string]string      `json:"customHeaders,omitempty"`
	RateLimitEnabled       *bool                  `json:"rateLimitEnabled,omitempty"`
	MaxRequestsPerMinute   *float64               `json:"maxRequestsPerMinute,omitempty"`
	MaxRetries             *float64               `json:"maxRetries,omitempty"`
	RetryBackoffMultiplier *float64               `json:"retryBackoffMultiplier,omitempty"`
	FilteringOptions       map[string]interface{} `json:"filteringOptions,omitempty"`
}

// ListWebhooksParams are query params for List.
type ListWebhooksParams struct {
	Page           *int
	Limit          *int
	Status         *WebhookStatus
	OrganizationID *string
	UserID         *string
}

// WebhookListResponse is a paginated webhook list.
type WebhookListResponse struct {
	Webhooks []Webhook `json:"webhooks"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

// DeleteWebhookResponse is returned by Delete.
type DeleteWebhookResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// TriggerWebhookParams triggers a webhook event.
type TriggerWebhookParams struct {
	Event          string                 `json:"event"`
	Payload        map[string]interface{} `json:"payload"`
	OrganizationID *string                `json:"organizationId,omitempty"`
	UserID         *string                `json:"userId,omitempty"`
}

// TriggerWebhookResponse is returned by Trigger.
type TriggerWebhookResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	JobID   *string `json:"jobId,omitempty"`
}

// WebhookDelivery is a webhook delivery attempt.
type WebhookDelivery struct {
	ID              string                 `json:"id"`
	WebhookID       string                 `json:"webhookId"`
	WebhookEventID  *string                `json:"webhookEventId"`
	EventType       string                 `json:"eventType"`
	EventData       map[string]interface{} `json:"eventData"`
	Status          WebhookDeliveryStatus  `json:"status"`
	RequestURL      string                 `json:"requestUrl"`
	RequestHeaders  map[string]string      `json:"requestHeaders"`
	RequestBody     map[string]interface{} `json:"requestBody"`
	ResponseStatus  *int                   `json:"responseStatus"`
	ResponseBody    *string                `json:"responseBody"`
	ResponseHeaders map[string]string      `json:"responseHeaders"`
	AttemptNumber   int                    `json:"attemptNumber"`
	MaxAttempts     int                    `json:"maxAttempts"`
	NextRetryAt     *string                `json:"nextRetryAt"`
	LastAttemptAt   *string                `json:"lastAttemptAt"`
	ErrorMessage    *string                `json:"errorMessage"`
	ErrorDetails    map[string]interface{} `json:"errorDetails"`
	CompletedAt     *string                `json:"completedAt"`
	DurationMs      *int                   `json:"durationMs"`
	CreatedAt       string                 `json:"createdAt"`
}

// ListWebhookDeliveriesParams are query params for ListDeliveries.
type ListWebhookDeliveriesParams struct {
	Page   *int
	Limit  *int
	Status *WebhookDeliveryStatus
}

// WebhookDeliveryListResponse is a paginated delivery list.
type WebhookDeliveryListResponse struct {
	Deliveries []WebhookDelivery `json:"deliveries"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// RetryWebhookDeliveryResponse is returned by RetryDelivery.
type RetryWebhookDeliveryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func requireWebhookID(id, field string) (string, error) {
	if field == "" {
		field = "webhookId"
	}
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Webhook %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireDeliveryID(id, field string) (string, error) {
	if field == "" {
		field = "deliveryId"
	}
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Webhook delivery %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireWebhookStatus(status WebhookStatus, field string) error {
	switch status {
	case WebhookStatusActive, WebhookStatusPaused, WebhookStatusDisabled, WebhookStatusFailed:
		return nil
	default:
		return validationErrorf(field, "update status must be a valid webhook status.")
	}
}

func validateUpdateWebhookParams(params UpdateWebhookParams) error {
	if params.Description == nil && params.Name == nil && params.URL == nil &&
		params.Secret == nil && params.Status == nil && params.CustomHeaders == nil &&
		params.RateLimitEnabled == nil && params.MaxRequestsPerMinute == nil &&
		params.MaxRetries == nil && params.RetryBackoffMultiplier == nil &&
		params.FilteringOptions == nil {
		return validationErrorf("params", "update requires at least one field to change.")
	}
	if params.Status != nil {
		if err := requireWebhookStatus(*params.Status, "status"); err != nil {
			return err
		}
	}
	if params.MaxRequestsPerMinute != nil {
		if err := requireFiniteNumber(*params.MaxRequestsPerMinute, "maxRequestsPerMinute"); err != nil {
			return err
		}
	}
	if params.MaxRetries != nil {
		if err := requireFiniteNumber(*params.MaxRetries, "maxRetries"); err != nil {
			return err
		}
	}
	if params.RetryBackoffMultiplier != nil {
		if err := requireFiniteNumber(*params.RetryBackoffMultiplier, "retryBackoffMultiplier"); err != nil {
			return err
		}
	}
	return nil
}

func buildListWebhooksQuery(params *ListWebhooksParams) (url.Values, error) {
	q := url.Values{}
	if params == nil {
		return q, nil
	}
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
	if params.Status != nil {
		if err := requireWebhookStatus(*params.Status, "status"); err != nil {
			return nil, err
		}
		q.Set("status", string(*params.Status))
	}
	if params.OrganizationID != nil {
		orgID, err := requireNonEmptyString(*params.OrganizationID, "organizationId")
		if err != nil {
			return nil, validationErrorf("organizationId", "list organizationId must be a non-empty string when provided.")
		}
		q.Set("organizationId", orgID)
	}
	if params.UserID != nil {
		userID, err := requireNonEmptyString(*params.UserID, "userId")
		if err != nil {
			return nil, validationErrorf("userId", "list userId must be a non-empty string when provided.")
		}
		q.Set("userId", userID)
	}
	return q, nil
}

func buildListWebhookDeliveriesQuery(params *ListWebhookDeliveriesParams) (url.Values, error) {
	q := url.Values{}
	if params == nil {
		return q, nil
	}
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
	if params.Status != nil {
		switch *params.Status {
		case WebhookDeliveryStatusPending, WebhookDeliveryStatusSuccess,
			WebhookDeliveryStatusFailed, WebhookDeliveryStatusRetrying:
			q.Set("status", string(*params.Status))
		default:
			return nil, validationErrorf("status", "listDeliveries status must be a valid delivery status.")
		}
	}
	return q, nil
}

// Create creates a webhook.
func (s *WebhookService) Create(params CreateWebhookParams) (*Webhook, error) {
	var out Webhook
	if err := s.client.do("POST", webhookV1+"/", params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List lists webhooks.
func (s *WebhookService) List(params *ListWebhooksParams) (*WebhookListResponse, error) {
	q, err := buildListWebhooksQuery(params)
	if err != nil {
		return nil, err
	}
	var out WebhookListResponse
	if err := s.client.do("GET", webhookV1, nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get fetches a webhook by ID.
func (s *WebhookService) Get(webhookID string) (*Webhook, error) {
	id, err := requireWebhookID(webhookID, "webhookId")
	if err != nil {
		return nil, err
	}
	var out Webhook
	if err := s.client.do("GET", webhookV1+"/"+id, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a webhook.
func (s *WebhookService) Update(webhookID string, params UpdateWebhookParams) (*Webhook, error) {
	id, err := requireWebhookID(webhookID, "webhookId")
	if err != nil {
		return nil, err
	}
	if err := validateUpdateWebhookParams(params); err != nil {
		return nil, err
	}
	var out Webhook
	if err := s.client.do("PATCH", webhookV1+"/"+id, params, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a webhook.
func (s *WebhookService) Delete(webhookID string) (*DeleteWebhookResponse, error) {
	id, err := requireWebhookID(webhookID, "webhookId")
	if err != nil {
		return nil, err
	}
	var out DeleteWebhookResponse
	if err := s.client.do("DELETE", webhookV1+"/"+id, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Pause pauses a webhook.
func (s *WebhookService) Pause(webhookID string) (*Webhook, error) {
	status := WebhookStatusPaused
	return s.Update(webhookID, UpdateWebhookParams{Status: &status})
}

// Enable enables a webhook.
func (s *WebhookService) Enable(webhookID string) (*Webhook, error) {
	status := WebhookStatusActive
	return s.Update(webhookID, UpdateWebhookParams{Status: &status})
}

// Disable disables a webhook.
func (s *WebhookService) Disable(webhookID string) (*Webhook, error) {
	status := WebhookStatusDisabled
	return s.Update(webhookID, UpdateWebhookParams{Status: &status})
}

// Trigger triggers a webhook event.
func (s *WebhookService) Trigger(params TriggerWebhookParams) (*TriggerWebhookResponse, error) {
	event, err := requireNonEmptyString(params.Event, "event")
	if err != nil {
		return nil, validationErrorf("event", "trigger event is required and must be a non-empty string.")
	}
	if params.Payload == nil {
		return nil, validationErrorf("payload", "trigger payload is required and must be an object.")
	}
	body := TriggerWebhookParams{
		Event:   event,
		Payload: params.Payload,
	}
	if params.OrganizationID != nil {
		orgID, err := requireNonEmptyString(*params.OrganizationID, "organizationId")
		if err != nil {
			return nil, validationErrorf("organizationId", "trigger organizationId must be a non-empty string when provided.")
		}
		body.OrganizationID = &orgID
	}
	if params.UserID != nil {
		userID, err := requireNonEmptyString(*params.UserID, "userId")
		if err != nil {
			return nil, validationErrorf("userId", "trigger userId must be a non-empty string when provided.")
		}
		body.UserID = &userID
	}
	var out TriggerWebhookResponse
	if err := s.client.do("POST", webhookV1+"/trigger", body, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListDeliveries lists delivery attempts for a webhook.
func (s *WebhookService) ListDeliveries(webhookID string, params *ListWebhookDeliveriesParams) (*WebhookDeliveryListResponse, error) {
	id, err := requireWebhookID(webhookID, "webhookId")
	if err != nil {
		return nil, err
	}
	q, err := buildListWebhookDeliveriesQuery(params)
	if err != nil {
		return nil, err
	}
	var out WebhookDeliveryListResponse
	path := webhookV1 + "/" + id + "/deliveries"
	if err := s.client.do("GET", path, nil, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RetryDelivery retries a failed delivery.
func (s *WebhookService) RetryDelivery(deliveryID string) (*RetryWebhookDeliveryResponse, error) {
	id, err := requireDeliveryID(deliveryID, "deliveryId")
	if err != nil {
		return nil, err
	}
	var out RetryWebhookDeliveryResponse
	path := "/api/webhook/deliveries/" + id + "/retry"
	if err := s.client.do("POST", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
