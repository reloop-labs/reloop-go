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

func TestDomainCreateHappyPath(t *testing.T) {
	var method, path string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(Domain{ID: "dom_1", Domain: "example.com", Object: "domain"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Domain.Create(CreateDomainParams{Domain: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/domain/v1/create" || body["domain"] != "example.com" {
		t.Fatalf("mismatch %s %s %#v", method, path, body)
	}
	if res.ID != "dom_1" {
		t.Fatalf("id = %q", res.ID)
	}
}

func TestDomainValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Domain.Get("  ")
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatal("expected no HTTP")
	}
}

func TestDomainAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "missing"})
	}))
	defer srv.Close()
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Domain.Get("dom_missing")
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestDomainServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&DomainService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{"Create", "Delete", "Get", "List", "Update", "Verify"}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
