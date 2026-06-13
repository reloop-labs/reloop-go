package reloop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client

	ApiKeys       *ApiKeysService
	APIKeyService *ApiKeysService
	Contacts      *ContactsService
	Domain        *DomainService
}

type ClientOptions struct {
	APIKey  string
	BaseURL string
}

func NewClient(options ClientOptions) (*Client, error) {
	if options.APIKey == "" {
		return nil, fmt.Errorf("Reloop SDK requires an APIKey")
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "https://reloop.sh"
	}

	c := &Client{
		APIKey:  options.APIKey,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	c.ApiKeys = &ApiKeysService{client: c}
	c.APIKeyService = c.ApiKeys
	c.Contacts = newContactsService(c)
	c.Domain = &DomainService{client: c}

	return c, nil
}

type APIError struct {
	StatusCode int
	Message    string
	Cause      interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Reloop API Error: %d %s", e.StatusCode, e.Message)
}

func (c *Client) request(method, path string, body interface{}, responseObj interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Reloop Network Error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var errBody map[string]interface{}
		json.NewDecoder(res.Body).Decode(&errBody)
		return &APIError{
			StatusCode: res.StatusCode,
			Message:    res.Status,
			Cause:      errBody,
		}
	}

	if res.StatusCode == 204 || responseObj == nil {
		return nil
	}

	return json.NewDecoder(res.Body).Decode(responseObj)
}
