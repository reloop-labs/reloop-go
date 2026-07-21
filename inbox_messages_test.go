package reloop

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestInboxMessagesSendHappyPath(t *testing.T) {
	var method, path, apiKey string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(SendEmailResponse{
			Success: true, MessageID: "msg_1", Status: "queued", Timestamp: "t", ID: "em_1",
		})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Inbox.Messages.Send(SendMessageParams{
		MailboxID: "mbx_1",
		To:        "user@example.com",
		Subject:   "Hello",
		HTML:      String("<p>Hi</p>"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/inbox/v1/messages/send" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if body["mailboxId"] != "mbx_1" || body["to"] != "user@example.com" || body["subject"] != "Hello" {
		t.Fatalf("body = %#v", body)
	}
	if !res.Success || res.ID != "em_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestInboxMessagesAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Inbox.Messages.Get("msg_missing")
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestInboxMessagesValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Inbox.Messages.Send(SendMessageParams{
		MailboxID: "",
		To:        "user@example.com",
		Subject:   "Hello",
	})
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Field != "mailboxId" {
		t.Fatalf("field = %q", ve.Field)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestInboxMessagesServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&InboxMessagesService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{
		"Batch", "CancelPending", "Delete", "Forward", "Get", "GetAttachment",
		"GetRaw", "List", "ListSent", "Reply", "ReplyAll", "Send", "SetRead",
		"SetStar", "Update",
	}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
