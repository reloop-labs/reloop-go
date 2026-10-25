package reloop

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
	"time"
)

const testWebhookSecret = "whsec_test_secret"

var testWebhookPayload = []byte(`{"id":"evt_123456789","event":"domain.created","payload":{"domainId":"dom_1"},"timestamp":1735689600}`)

func signWebhookPayload(secret, timestamp string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))
	return "t=" + timestamp + ",v1=" + signature
}

func TestVerifyWebhookAcceptsValidSignature(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, testWebhookPayload)

	event, err := VerifyWebhook(VerifyWebhookParams{
		Payload: testWebhookPayload,
		Headers: map[string]string{WebhookSignatureHeader: header},
		Secret:  testWebhookSecret,
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.ID != "evt_123456789" || event.Event != "domain.created" {
		t.Fatalf("event = %#v", event)
	}
	if event.Payload["domainId"] != "dom_1" {
		t.Fatalf("payload = %#v", event.Payload)
	}
}

func TestConstructEventVerifiesSignature(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, testWebhookPayload)

	event, err := ConstructEvent(testWebhookPayload, header, testWebhookSecret, 0)
	if err != nil {
		t.Fatal(err)
	}
	if event.Event != "domain.created" {
		t.Fatalf("event = %q", event.Event)
	}
}

func TestWebhookServiceVerify(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, testWebhookPayload)

	c, _ := NewClient(ClientOptions{APIKey: "rl_test"})
	event, err := c.Webhook.Verify(VerifyWebhookParams{
		Payload: testWebhookPayload,
		Headers: map[string]string{WebhookSignatureHeader: header},
		Secret:  testWebhookSecret,
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.ID != "evt_123456789" {
		t.Fatalf("id = %q", event.ID)
	}
}

func TestVerifyWebhookRejectsWrongSecret(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, testWebhookPayload)

	_, err := VerifyWebhook(VerifyWebhookParams{
		Payload: testWebhookPayload,
		Headers: map[string]string{WebhookSignatureHeader: header},
		Secret:  "wrong_secret",
	})
	if _, ok := err.(*WebhookSignatureVerificationError); !ok {
		t.Fatalf("expected WebhookSignatureVerificationError, got %T", err)
	}
}

func TestVerifyWebhookRejectsExpiredTimestamp(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix()-600, 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, testWebhookPayload)

	_, err := VerifyWebhook(VerifyWebhookParams{
		Payload:   testWebhookPayload,
		Headers:   map[string]string{WebhookSignatureHeader: header},
		Secret:    testWebhookSecret,
		Tolerance: 300,
	})
	sigErr, ok := err.(*WebhookSignatureVerificationError)
	if !ok {
		t.Fatalf("expected WebhookSignatureVerificationError, got %T", err)
	}
	if !strings.Contains(strings.ToLower(sigErr.Message), "tolerance") {
		t.Fatalf("message = %q", sigErr.Message)
	}
}

func TestVerifyWebhookRejectsMalformedSignatureHeader(t *testing.T) {
	_, err := VerifyWebhook(VerifyWebhookParams{
		Payload: testWebhookPayload,
		Headers: map[string]string{WebhookSignatureHeader: "invalid-header"},
		Secret:  testWebhookSecret,
	})
	sigErr, ok := err.(*WebhookSignatureVerificationError)
	if !ok {
		t.Fatalf("expected WebhookSignatureVerificationError, got %T", err)
	}
	if !strings.Contains(sigErr.Message, "Invalid Reloop-Signature") {
		t.Fatalf("message = %q", sigErr.Message)
	}
}

func TestVerifyWebhookRejectsInvalidJSONPayload(t *testing.T) {
	badPayload := []byte("not-json")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, badPayload)

	_, err := VerifyWebhook(VerifyWebhookParams{
		Payload: badPayload,
		Headers: map[string]string{WebhookSignatureHeader: header},
		Secret:  testWebhookSecret,
	})
	sigErr, ok := err.(*WebhookSignatureVerificationError)
	if !ok {
		t.Fatalf("expected WebhookSignatureVerificationError, got %T", err)
	}
	if !strings.Contains(strings.ToLower(sigErr.Message), "not valid json") {
		t.Fatalf("message = %q", sigErr.Message)
	}
}

func TestVerifyWebhookRejectsMissingSignatureHeader(t *testing.T) {
	_, err := VerifyWebhook(VerifyWebhookParams{
		Payload: testWebhookPayload,
		Headers: map[string]string{},
		Secret:  testWebhookSecret,
	})
	sigErr, ok := err.(*WebhookSignatureVerificationError)
	if !ok {
		t.Fatalf("expected WebhookSignatureVerificationError, got %T", err)
	}
	if !strings.Contains(strings.ToLower(sigErr.Message), "missing reloop-signature") {
		t.Fatalf("message = %q", sigErr.Message)
	}
}

func TestVerifyWebhookRejectsMissingSecret(t *testing.T) {
	_, err := VerifyWebhook(VerifyWebhookParams{
		Payload: testWebhookPayload,
		Headers: map[string]string{WebhookSignatureHeader: "t=1,v1=abc"},
		Secret:  "",
	})
	if _, ok := err.(*WebhookSignatureVerificationError); !ok {
		t.Fatalf("expected WebhookSignatureVerificationError, got %T", err)
	}
}

// productionWebhookPayload matches the envelope Reloop actually sends.
var productionWebhookPayload = []byte(`{"id":"whev_01abc","type":"email.received","created_at":"2026-10-26T12:00:00.000Z","data":{"email_id":"in_1","subject":"Hello"}}`)

func TestVerifyWebhookAcceptsProductionDelivery(t *testing.T) {
	if WebhookSignatureHeader != "Reloop-Signature" {
		t.Fatalf("WebhookSignatureHeader = %q", WebhookSignatureHeader)
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, productionWebhookPayload)

	event, err := VerifyWebhook(VerifyWebhookParams{
		Payload: productionWebhookPayload,
		Headers: map[string]string{
			"reloop-signature": header,
			"reloop-timestamp": timestamp,
			"reloop-id":        "whev_01abc",
		},
		Secret: testWebhookSecret,
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.ID != "whev_01abc" || event.Type != "email.received" || event.CreatedAt != "2026-10-26T12:00:00.000Z" {
		t.Fatalf("event = %#v", event)
	}
	if event.Data["subject"] != "Hello" || event.Event != event.Type || event.Payload["subject"] != "Hello" {
		t.Fatalf("event data = %#v", event)
	}
}

func TestVerifyWebhookAcceptsLegacyHeader(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, productionWebhookPayload)

	event, err := VerifyWebhook(VerifyWebhookParams{
		Payload: productionWebhookPayload,
		Headers: map[string]string{"x-webhook-signature": header},
		Secret:  testWebhookSecret,
	})
	if err != nil || event.Type != "email.received" {
		t.Fatalf("event = %#v, err = %v", event, err)
	}
}

func TestVerifyWebhookRejectsTamperedProductionBody(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	header := signWebhookPayload(testWebhookSecret, timestamp, productionWebhookPayload)
	tampered := []byte(strings.Replace(string(productionWebhookPayload), "Hello", "Hacked", 1))

	if _, err := ConstructEvent(tampered, header, testWebhookSecret, 0); err == nil {
		t.Fatal("expected verification error")
	}
}
