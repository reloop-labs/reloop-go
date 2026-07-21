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

func TestWebhookCreateHappyPath(t *testing.T) {
	var method, path string
	var body map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_ = json.NewEncoder(w).Encode(Webhook{ID: "wh_1", Name: "hooks", URL: "https://example.com/hook", Status: WebhookStatusActive})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	res, err := c.Webhook.Create(CreateWebhookParams{
		Description: "prod hooks",
		URL:         "https://example.com/hook",
		Events:      []string{"domain.created"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/webhook/v1/" {
		t.Fatalf("mismatch method=%s path=%s", method, path)
	}
	if body["url"] != "https://example.com/hook" {
		t.Fatalf("body = %#v", body)
	}
	if res.ID != "wh_1" {
		t.Fatalf("id = %q", res.ID)
	}
}

func TestWebhookValidationNoHTTP(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})

	_, err := c.Webhook.Get("  ")
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	_, err = c.Webhook.Update("wh_1", UpdateWebhookParams{})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError for empty update, got %T", err)
	}

	_, err = c.Webhook.Trigger(TriggerWebhookParams{Event: "", Payload: map[string]interface{}{"a": 1}})
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError for trigger event, got %T", err)
	}

	if hits != 0 {
		t.Fatal("expected no HTTP")
	}
}

func TestWebhookAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "missing"})
	}))
	defer srv.Close()
	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Webhook.Get("wh_missing")
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestWebhookServiceSurface(t *testing.T) {
	methods := []string{}
	typ := reflect.TypeOf(&WebhookService{})
	for i := 0; i < typ.NumMethod(); i++ {
		methods = append(methods, typ.Method(i).Name)
	}
	sort.Strings(methods)
	expected := []string{
		"Create", "Delete", "Disable", "Enable", "Get", "List",
		"ListDeliveries", "Pause", "RetryDelivery", "Trigger", "Update", "Verify",
	}
	if !reflect.DeepEqual(methods, expected) {
		t.Fatalf("surface = %v want %v", methods, expected)
	}
}

func TestWebhookPauseEnableDisable(t *testing.T) {
	var bodies []map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		var body map[string]interface{}
		_ = json.Unmarshal(data, &body)
		bodies = append(bodies, body)
		_ = json.NewEncoder(w).Encode(Webhook{ID: "wh_1", Status: WebhookStatusActive})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	if _, err := c.Webhook.Pause("wh_1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Webhook.Enable("wh_1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Webhook.Disable("wh_1"); err != nil {
		t.Fatal(err)
	}

	if len(bodies) != 3 {
		t.Fatalf("bodies = %d", len(bodies))
	}
	if bodies[0]["status"] != "paused" || bodies[1]["status"] != "active" || bodies[2]["status"] != "disabled" {
		t.Fatalf("status bodies = %#v", bodies)
	}
}

func TestWebhookTriggerPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(TriggerWebhookResponse{Success: true, Message: "queued"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Webhook.Trigger(TriggerWebhookParams{
		Event:   "domain.created",
		Payload: map[string]interface{}{"domainId": "dom_1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/webhook/v1/trigger" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}

func TestWebhookListDeliveriesPath(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_ = json.NewEncoder(w).Encode(WebhookDeliveryListResponse{Deliveries: []WebhookDelivery{}, Total: 0})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Webhook.ListDeliveries("wh_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if path != "/api/webhook/v1/wh_1/deliveries" {
		t.Fatalf("path = %q", path)
	}
}

func TestWebhookRetryDeliveryPath(t *testing.T) {
	var method, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(RetryWebhookDeliveryResponse{Success: true, Message: "retry queued"})
	}))
	defer srv.Close()

	c, _ := NewClient(ClientOptions{APIKey: "rl_test", BaseURL: srv.URL})
	_, err := c.Webhook.RetryDelivery("del_1")
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/api/webhook/deliveries/del_1/retry" {
		t.Fatalf("method=%s path=%s", method, path)
	}
}
