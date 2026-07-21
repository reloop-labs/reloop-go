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

func TestChannelsCreateHappyPath(t *testing.T) {
	var method, path, apiKey string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(ContactChannelResponse{
			ContactChannel: ContactChannel{Object: "channel", ID: "chn_1", Name: "Newsletter"},
			Event:          "channel.created",
		})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Contacts.Channels.Create(CreateChannelParams{Name: "Newsletter"})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/contacts/v1/channels/create" || apiKey != "rl_test" {
		t.Fatalf("request mismatch: %s %s key=%s", method, path, apiKey)
	}
	if body["name"] != "Newsletter" {
		t.Fatalf("body = %#v", body)
	}
	if res.ID != "chn_1" {
		t.Fatalf("response = %+v", res)
	}
}

func TestChannelsCreateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.Create(CreateChannelParams{Name: "Newsletter"})
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestChannelsCreateValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.Create(CreateChannelParams{Name: ""})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestChannelsListHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(ChannelListResponse{Object: "channel"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if method != "GET" || path != "/api/contacts/v1/channels/list" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestChannelsGetHappyPath(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_ = json.NewEncoder(w).Encode(ContactChannel{Object: "channel", ID: "chn_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.Get("chn_1")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/api/contacts/v1/channels/chn_1" {
		t.Fatalf("path=%s", path)
	}
}

func TestChannelsUpdateHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(ContactChannelResponse{ContactChannel: ContactChannel{ID: "chn_1"}})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.Update("chn_1", UpdateChannelParams{Name: String("Updated")})
	if err != nil {
		t.Fatal(err)
	}
	if method != "PATCH" || path != "/api/contacts/v1/channels/chn_1" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestChannelsDeleteHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(DeleteChannelResponse{Success: true, ID: "chn_1"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.Delete("chn_1")
	if err != nil {
		t.Fatal(err)
	}
	if method != "DELETE" || path != "/api/contacts/v1/channels/chn_1" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestChannelsAddContactHappyPath(t *testing.T) {
	var method, path string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(AddContactToChannelResponse{SubscriptionID: "sub_1"})
	}))
	defer srv.Close()

	sub := ChannelSubscriptionOptIn
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.AddContact("chn_1", AddContactToChannelParams{
		Email:        String("a@x.com"),
		Subscription: &sub,
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/contacts/channel/chn_1" || body["email"] != "a@x.com" {
		t.Fatalf("method=%s path=%s body=%v", method, path, body)
	}
}

func TestChannelsUpdateSubscriptionHappyPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(UpdateContactChannelResponse{Success: true, Status: "enrolled"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.UpdateSubscription("chn_1", UpdateContactChannelParams{
		Email:        String("a@x.com"),
		Subscription: ChannelSubscriptionOptOut,
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "PATCH" || path != "/api/contacts/channel/chn_1" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestChannelsUpdateSubscriptionValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Contacts.Channels.UpdateSubscription("chn_1", UpdateContactChannelParams{
		Subscription: ChannelSubscriptionOptIn,
	})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if hits != 0 {
		t.Fatalf("expected 0 HTTP requests, got %d", hits)
	}
}

func TestChannelsServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&ChannelsService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{"AddContact", "Create", "Delete", "Get", "List", "Update", "UpdateSubscription"}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}
