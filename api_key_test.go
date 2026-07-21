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

func TestAPIKeyCreateHappyPath(t *testing.T) {
	var method, path string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(APIKeyWithKey{ID: "key_1", Key: "rl_secret", Object: "api_key"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.APIKey.Create(CreateAPIKeyParams{Name: "prod"})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/api-key/v1/" || body["name"] != "prod" {
		t.Fatalf("mismatch method=%s path=%s body=%v", method, path, body)
	}
	if res.Key != "rl_secret" {
		t.Fatalf("key = %q", res.Key)
	}
}

func TestAPIKeyValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.APIKey.Create(CreateAPIKeyParams{Name: ""})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatal("expected no HTTP")
	}
}

func TestAPIKeyAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.APIKey.Get("key_1")
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestAPIKeyServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&APIKeyService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{"Create", "Delete", "Disable", "Enable", "Get", "List", "Rotate", "Update"}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
