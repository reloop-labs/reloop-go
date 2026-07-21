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

func TestPropertiesCreateHappyPath(t *testing.T) {
	var method, path, apiKey string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(ContactPropertyResponse{
			ContactProperty: ContactProperty{Object: "contact_property", ID: "prop_1", PropertyName: "company", PropertyType: PropertyTypeString},
			Event:           "property.created",
		})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Contacts.Properties.Create(CreatePropertyParams{Name: "company", Type: PropertyTypeString})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/contacts/v1/properties/create" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if body["name"] != "company" || body["type"] != "string" {
		t.Fatalf("body = %#v", body)
	}
	if res.ID != "prop_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestPropertiesCreateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Properties.Create(CreatePropertyParams{Name: "company", Type: PropertyTypeString})
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestPropertiesCreateValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Properties.Create(CreatePropertyParams{Name: "", Type: PropertyTypeString})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestPropertiesListHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(PropertyListResponse{Object: "contact_property", Event: "properties.listed"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Properties.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if method != "GET" || path != "/api/contacts/v1/properties/list" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestPropertiesUpdateHappyPath(t *testing.T) {
	var method, path string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(ContactPropertyResponse{ContactProperty: ContactProperty{ID: "prop_1"}})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Properties.Update("prop_1", UpdatePropertyParams{FallbackValue: String("default")})
	if err != nil {
		t.Fatal(err)
	}
	if method != "PATCH" || path != "/api/contacts/v1/properties/prop_1" || body["fallbackValue"] != "default" {
		t.Fatalf("method=%s path=%s body=%v", method, path, body)
	}
}

func TestPropertiesDeleteHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(DeletePropertyResponse{Success: true, ID: "prop_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Properties.Delete("prop_1")
	if err != nil {
		t.Fatal(err)
	}
	if method != "DELETE" || path != "/api/contacts/v1/properties/prop_1" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestPropertiesServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&PropertiesService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{"Create", "Delete", "List", "Update"}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
