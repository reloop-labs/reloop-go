package reloop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://reloop.sh"

// Client is the Reloop SDK entry point.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	APIKey   *APIKeyService
	Contacts *ContactsService
	Domain   *DomainService
	Mail     *MailService
	Webhook  *WebhookService
	Inbox    *InboxService
}

// ClientOptions configures NewClient.
type ClientOptions struct {
	APIKey  string
	BaseURL string
}

// NewClient creates a Reloop client.
func NewClient(options ClientOptions) (*Client, error) {
	apiKey := strings.TrimSpace(options.APIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("Reloop SDK requires an APIKey")
	}

	baseURL := strings.TrimRight(options.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	c := &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	c.APIKey = &APIKeyService{client: c}
	c.Contacts = newContactsService(c)
	c.Domain = &DomainService{client: c}
	c.Mail = &MailService{client: c}
	c.Webhook = &WebhookService{client: c}
	c.Inbox = newInboxService(c)

	return c, nil
}

// HTTPClient returns the underlying HTTP client (for tests / custom transport).
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// SetHTTPClient replaces the HTTP client (useful in tests).
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	if httpClient != nil {
		c.httpClient = httpClient
	}
}

// BaseURL returns the configured API base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) do(method, path string, body interface{}, query url.Values, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	fullURL := c.baseURL + path
	if len(query) > 0 {
		fullURL = fullURL + "?" + query.Encode()
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return &APIError{Status: 0, StatusText: "Network Error", Body: APIErrorBody{Message: err.Error()}}
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var errBody APIErrorBody
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		return &APIError{
			Status:     res.StatusCode,
			StatusText: res.Status,
			Body:       errBody,
		}
	}

	if res.StatusCode == 204 || out == nil {
		return nil
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, out)
}

// Pointer helpers for optional fields.
func String(v string) *string { return &v }
func Bool(v bool) *bool       { return &v }
func Int(v int) *int          { return &v }
