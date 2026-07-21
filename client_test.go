package reloop

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientRequiresAPIKey(t *testing.T) {
	_, err := NewClient(ClientOptions{})
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
	_, err = NewClient(ClientOptions{APIKey: "   "})
	if err == nil {
		t.Fatal("expected error for whitespace API key")
	}
}

func TestNewClientDefaults(t *testing.T) {
	c, err := NewClient(ClientOptions{APIKey: "rl_test"})
	if err != nil {
		t.Fatal(err)
	}
	if c.BaseURL() != defaultBaseURL {
		t.Fatalf("base URL = %q", c.BaseURL())
	}
	if c.APIKey == nil || c.Mail == nil || c.Domain == nil || c.Contacts == nil || c.Webhook == nil || c.Inbox == nil {
		t.Fatal("expected all services to be wired")
	}
}

func TestClientAuthHeaderAndAPIError(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("x-api-key")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad key"})
	}))
	defer srv.Close()

	c, err := NewClient(ClientOptions{APIKey: "rl_secret", BaseURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}

	var out map[string]interface{}
	err = c.do("GET", "/v1/test", nil, nil, &out)
	if gotKey != "rl_secret" {
		t.Fatalf("auth header = %q", gotKey)
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T %v", err, err)
	}
	if apiErr.Status != 401 || apiErr.Body.Message != "bad key" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
}

func TestClientNetworkError(t *testing.T) {
	c, err := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	err = c.do("GET", "/v1/test", nil, nil, &out)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Status != 0 {
		t.Fatalf("expected status 0 for network error, got %d", apiErr.Status)
	}
}
