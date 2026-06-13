package reloop_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	reloop "github.com/reloop-labs/reloop-go"
)

func TestApiKeyCreateUsesApiPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/api-key/v1/" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), `"name":"Production Key"`) {
			t.Fatalf("unexpected body: %s", body)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        "key_1",
			"name":      "Production Key",
			"key":       "rl_live_secret",
			"enabled":   true,
			"createdAt": "2026-01-01T00:00:00.000Z",
			"updatedAt": "2026-01-01T00:00:00.000Z",
			"object":    "api_key",
			"event":     "evt_1",
		})
	}))
	defer server.Close()

	client, err := reloop.NewClient(reloop.ClientOptions{
		APIKey:  "rl_test",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatal(err)
	}

	key, err := client.ApiKeys.Create(reloop.CreateApiKeyParams{Name: "Production Key"})
	if err != nil {
		t.Fatal(err)
	}

	if key.ID != "key_1" {
		t.Fatalf("expected key_1, got %s", key.ID)
	}
}

func TestApiKeyPauseUsesDisableRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/api-key/v1/disable/key_1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      "key_1",
			"enabled": false,
			"object":  "api_key",
			"event":   "evt_1",
		})
	}))
	defer server.Close()

	client, err := reloop.NewClient(reloop.ClientOptions{
		APIKey:  "rl_test",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ApiKeys.Pause("key_1")
	if err != nil {
		t.Fatal(err)
	}
}
