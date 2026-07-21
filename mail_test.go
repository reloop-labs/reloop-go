package reloop

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMailSendHappyPath(t *testing.T) {
	var method, path, apiKey string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(SendMailResponse{
			Success: true, MessageID: "msg_1", Status: "queued", Timestamp: "t", ID: "em_1",
		})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Mail.Send(SendMailParams{
		From: "a@example.com", To: "b@example.com", Subject: "Hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/mail/v1/send" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if body["from"] != "a@example.com" || body["to"] != "b@example.com" || body["subject"] != "Hi" {
		t.Fatalf("body = %#v", body)
	}
	if !res.Success || res.ID != "em_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestMailSendAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Mail.Send(SendMailParams{From: "a@x.com", To: "b@x.com", Subject: "Hi"})
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestMailSendValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Mail.Send(SendMailParams{From: "", To: "b@x.com", Subject: "Hi"})
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Field != "from" {
		t.Fatalf("field = %q", ve.Field)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestMailServiceSurface(t *testing.T) {
	methods := map[string]bool{}
	typ := reflect.TypeOf(&MailService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods[typ.Method(i).Name] = true
	}
	expected := []string{"Send"}
	if len(methods) != len(expected) {
		t.Fatalf("methods = %v", methods)
	}
	for _, name := range expected {
		if !methods[name] {
			t.Fatalf("missing method %s", name)
		}
	}
}
