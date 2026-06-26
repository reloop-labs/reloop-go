package reloop_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	reloop "github.com/reloop-labs/reloop-go/v2"
)

func TestDomainCreateUsesSnakeCaseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/domain/v1/create" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), `"click_tracking":true`) {
			t.Fatalf("expected snake_case click_tracking in body, got %s", body)
		}
		if strings.Contains(string(body), `"clickTracking"`) {
			t.Fatalf("did not expect camelCase clickTracking in body, got %s", body)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"object":                   "domain",
			"id":                       "dom_1",
			"domain":                   "send.example.com",
			"status":                   "pending",
			"userVerifiedDomain":       false,
			"systemVerified":           false,
			"customReturnPath":         "inbound",
			"trackingSubdomain":        "tracking",
			"isClickTrackingEnabled":   true,
			"isOpenTrackingEnabled":    false,
			"tls":                      "opportunistic",
			"isTrackingDomain":         false,
			"isSendingEmailEnabled":    true,
			"isReceivingEmailEnabled":  true,
			"verificationFailedReason": nil,
			"dnsRecords":               []interface{}{},
			"lastVerifiedAt":           nil,
			"createdAt":                "2026-01-01T00:00:00.000Z",
			"updatedAt":                "2026-01-01T00:00:00.000Z",
		})
	}))
	defer server.Close()

	client, err := reloop.NewClient(reloop.ClientOptions{
		APIKey:  "rl_test",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatal(err)
	}

	domain, err := client.Domain.Create(reloop.CreateDomainParams{
		Domain:        "send.example.com",
		ClickTracking: reloop.Bool(true),
	})
	if err != nil {
		t.Fatal(err)
	}

	if domain.ID != "dom_1" {
		t.Fatalf("expected dom_1, got %s", domain.ID)
	}
}

func TestDomainListBuildsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/domain/v1/list" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("page") != "2" || query.Get("limit") != "5" || query.Get("q") != "example" || query.Get("status") != "active" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"object":  "domain",
			"domains": []interface{}{},
			"total":   0,
			"page":    2,
			"limit":   5,
			"event":   "evt_1",
		})
	}))
	defer server.Close()

	client, err := reloop.NewClient(reloop.ClientOptions{
		APIKey:  "rl_test",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatal(err)
	}

	status := reloop.DomainStatusActive
	result, err := client.Domain.List(&reloop.ListDomainsParams{
		Page:   reloop.Int(2),
		Limit:  reloop.Int(5),
		Q:      reloop.String("example"),
		Status: &status,
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.Total != 0 || result.Page != 2 {
		t.Fatalf("unexpected list response: %+v", result)
	}
}

func TestDomainGetNameserversRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/domain/v1/nameservers/dom_1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		provider := "cloudflare"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reloop.DomainNameserversResponse{
			Object:      "domain_nameservers",
			DomainID:    "dom_1",
			Domain:      "send.example.com",
			Nameservers: []string{"ns1.example.net"},
			DNSProvider: &provider,
			Event:       "evt_1",
		})
	}))
	defer server.Close()

	client, err := reloop.NewClient(reloop.ClientOptions{
		APIKey:  "rl_test",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Domain.GetNameservers("dom_1")
	if err != nil {
		t.Fatal(err)
	}

	if result.DNSProvider == nil || *result.DNSProvider != "cloudflare" {
		t.Fatalf("unexpected dns provider: %+v", result.DNSProvider)
	}
}

func TestDomainForwardDNSRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/domain/v1/verify/dom_1/forward-dns" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), `"email":"admin@example.com"`) {
			t.Fatalf("unexpected body: %s", body)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reloop.ForwardDNSResponse{Success: true})
	}))
	defer server.Close()

	client, err := reloop.NewClient(reloop.ClientOptions{
		APIKey:  "rl_test",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.Domain.ForwardDNS("dom_1", reloop.ForwardDNSParams{
		Email: "admin@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !result.Success {
		t.Fatal("expected success response")
	}
}
