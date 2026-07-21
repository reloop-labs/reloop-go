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
	if !strings.Contains(sigErr.Message, "Invalid X-Webhook-Signature") {
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
	if !strings.Contains(strings.ToLower(sigErr.Message), "missing x-webhook-signature") {
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
