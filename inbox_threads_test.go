package reloop

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestInboxThreadsListHappyPath(t *testing.T) {
	var method, path, apiKey, limit string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		limit = r.URL.Query().Get("limit")
		_ = json.NewEncoder(w).Encode([]Thread{{
			ID: "thr_1", OrganizationID: "org_1", LastMessageAt: "t", Status: "active",
		}})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Inbox.Threads.List(&ListThreadsParams{Limit: Int(50)})
	if err != nil {
		t.Fatal(err)
	}
	if method != "GET" || path != "/api/inbox/v1/threads" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if limit != "50" {
		t.Fatalf("limit = %q", limit)
	}
	if len(res) != 1 || res[0].ID != "thr_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestInboxThreadsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Inbox.Threads.Get("thr_missing")
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestInboxThreadsValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Inbox.Threads.Batch(BatchThreadsParams{
		IDs:    []string{},
		Action: ThreadBatchActionArchive,
	})
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Field != "ids" {
		t.Fatalf("field = %q", ve.Field)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestInboxThreadsServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&InboxThreadsService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{
		"Archive", "Batch", "Delete", "Get", "GetAttachment", "List",
		"Restore", "SetRead", "SetStar", "Trash", "Update",
	}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
