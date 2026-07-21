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

func TestInboxMailboxesCreateHappyPath(t *testing.T) {
	var method, path, apiKey string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(CreateMailboxResponse{
			ID: "mbx_1", Email: "user@example.com", Status: "active",
		})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Inbox.Mailboxes.Create(CreateMailboxParams{
		DomainID: "dom_1",
		Email:    "user@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/inbox/v1/mailboxes/create" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if body["domainId"] != "dom_1" || body["email"] != "user@example.com" {
		t.Fatalf("body = %#v", body)
	}
	if res.ID != "mbx_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestInboxMailboxesAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Inbox.Mailboxes.Get("mbx_missing")
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestInboxMailboxesValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Inbox.Mailboxes.Get("  ")
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Field != "id" {
		t.Fatalf("field = %q", ve.Field)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestInboxMailboxesServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&InboxMailboxesService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{"Create", "Delete", "Get", "List", "Update"}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
