package reloop_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	reloop "github.com/reloop-labs/reloop-go"
)

func TestMailSendUsesSnakeCaseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/mail/v1/send" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}

		if payload["from"] != "Reloop <hello@send.example.com>" {
			t.Fatalf("expected from in body, got %s", body)
		}
		if payload["reply_to"] != "support@example.com" {
			t.Fatalf("expected reply_to in body, got %s", body)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"messageId": "msg_123456789",
			"status":    "sent",
			"timestamp": "2026-01-01T00:00:00.000Z",
			"id":        "log_123456789",
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

	result, err := client.Mail.Send(reloop.SendMailParams{
		From:    "Reloop <hello@send.example.com>",
		To:      "user@example.com",
		Subject: "Welcome to Reloop",
		HTML:    reloop.String("<p>Thanks for signing up.</p>"),
		Text:    reloop.String("Thanks for signing up."),
		ReplyTo: "support@example.com",
		Tags: []reloop.SendMailTag{
			{Name: "campaign", Value: "welcome"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.MessageID != "msg_123456789" || result.ID != "log_123456789" {
		t.Fatalf("unexpected response: %+v", result)
	}
}
