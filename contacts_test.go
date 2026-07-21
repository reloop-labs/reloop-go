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

func TestContactsCreateHappyPath(t *testing.T) {
	var method, path, apiKey string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(ContactResponse{
			Contact: Contact{Object: "contact", ID: "con_1", Email: "john@example.com", Status: ContactStatusSubscribed},
			Event:   "contact.created",
		})
	}))
	defer srv.Close()

	status := ContactStatusSubscribed
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Contacts.Create(CreateContactParams{
		Email:      "john@example.com",
		FirstName:  String("John"),
		LastName:   String("Doe"),
		Status:     &status,
		Properties: map[string]interface{}{"company": "Reloop"},
		GroupIDs:   []string{"grp_1"},
		Channels:   []ContactChannelInput{{ChannelID: "chn_1", Subscription: ChannelSubscriptionOptIn}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/contacts/create" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if body["email"] != "john@example.com" || body["firstName"] != "John" {
		t.Fatalf("body = %#v", body)
	}
	if res.ID != "con_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestContactsCreateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "exists"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Create(CreateContactParams{Email: "john@example.com"})
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestContactsCreateValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Create(CreateContactParams{Email: "not-an-email"})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestContactsGetHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(Contact{Object: "contact", ID: "con_1", Email: "a@x.com", Status: ContactStatusSubscribed})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Contacts.Get("con_1")
	if err != nil {
		t.Fatal(err)
	}
	if method != "GET" || path != "/api/contacts/retrieve/con_1" {
		t.Fatalf("request mismatch: %s %s", method, path)
	}
	if res.ID != "con_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestContactsListHappyPath(t *testing.T) {
	var path, query string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, query = r.URL.Path, r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(ContactListResponse{Object: "contact", Event: "contacts.listed"})
	}))
	defer srv.Close()

	page, limit := 2, 10
	status := ContactStatusSubscribed
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.List(&ListContactsParams{Page: &page, Limit: &limit, Search: String("john"), Status: &status})
	if err != nil {
		t.Fatal(err)
	}
	if path != "/api/contacts/list" || query == "" {
		t.Fatalf("path=%s query=%s", path, query)
	}
}

func TestContactsUpdateHappyPath(t *testing.T) {
	var method, path string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(ContactResponse{Contact: Contact{ID: "con_1", Email: "new@x.com"}, Event: "contact.updated"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Update("con_1", UpdateContactParams{Email: String("new@x.com")})
	if err != nil {
		t.Fatal(err)
	}
	if method != "PATCH" || path != "/api/contacts/con_1" || body["email"] != "new@x.com" {
		t.Fatalf("method=%s path=%s body=%v", method, path, body)
	}
}

func TestContactsUpdateValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Update("con_1", UpdateContactParams{})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestContactsDeleteHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(DeleteContactResponse{Success: true, Object: "contact", ID: "con_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Contacts.Delete("con_1")
	if err != nil {
		t.Fatal(err)
	}
	if method != "DELETE" || path != "/api/contacts/con_1" || !res.Success {
		t.Fatalf("method=%s path=%s res=%+v", method, path, res)
	}
}

func TestContactsServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&ContactsService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{"Create", "Delete", "Get", "List", "Update"}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
