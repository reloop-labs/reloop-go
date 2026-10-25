package reloop

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	// WebhookSignatureHeader is the header Reloop signs deliveries with: t=<unix>,v1=<hex>.
	WebhookSignatureHeader = "Reloop-Signature"
	// WebhookTimestampHeader carries the Unix timestamp used in the signature.
	WebhookTimestampHeader = "Reloop-Timestamp"
	// WebhookIDHeader carries the event id, the same as the body "id".
	WebhookIDHeader = "Reloop-Id"

	// legacyWebhookSignatureHeader was used before 2.0.1. Still accepted as a fallback.
	legacyWebhookSignatureHeader = "x-webhook-signature"

	defaultWebhookToleranceSeconds = 300
)

// WebhookSignatureVerificationError is returned when signature verification fails.
type WebhookSignatureVerificationError struct {
	Message string
}

func (e *WebhookSignatureVerificationError) Error() string {
	return e.Message
}

func webhookSigErr(message string) *WebhookSignatureVerificationError {
	return &WebhookSignatureVerificationError{Message: message}
}

// WebhookEvent is a verified webhook delivery: {id, type, created_at, data}.
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	CreatedAt string                 `json:"created_at"`
	Data      map[string]interface{} `json:"data"`

	// Deprecated: use Type.
	Event string `json:"-"`
	// Deprecated: use Data.
	Payload map[string]interface{} `json:"-"`
	// Deprecated: use CreatedAt. Unix seconds.
	Timestamp float64 `json:"-"`
}

// VerifyWebhookParams configures webhook signature verification.
type VerifyWebhookParams struct {
	Payload   []byte
	Headers   map[string]string
	Secret    string
	Tolerance int
}

func getWebhookHeader(headers map[string]string, name string) (string, bool) {
	lower := strings.ToLower(name)
	for key, value := range headers {
		if strings.ToLower(key) == lower {
			return value, true
		}
	}
	return "", false
}

func parseSignatureHeader(header string) (timestamp string, signatures []string, err error) {
	for _, element := range strings.Split(header, ",") {
		trimmed := strings.TrimSpace(element)
		eqIndex := strings.Index(trimmed, "=")
		if eqIndex == -1 {
			continue
		}
		prefix := trimmed[:eqIndex]
		value := trimmed[eqIndex+1:]
		switch prefix {
		case "t":
			timestamp = value
		case "v1":
			signatures = append(signatures, value)
		}
	}
	if timestamp == "" || len(signatures) == 0 {
		return "", nil, webhookSigErr("Invalid " + WebhookSignatureHeader + " header: expected t= and v1= values")
	}
	return timestamp, signatures, nil
}

func computeExpectedWebhookSignature(secret, timestamp string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func verifyWebhookSignature(expected, received string) bool {
	expectedBuf, err1 := hex.DecodeString(expected)
	receivedBuf, err2 := hex.DecodeString(received)
	if err1 != nil || err2 != nil || len(expectedBuf) != len(receivedBuf) {
		return false
	}
	return hmac.Equal(expectedBuf, receivedBuf)
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func verifyWebhookTimestamp(timestamp string, tolerance int) error {
	ts, err := strconv.ParseFloat(timestamp, 64)
	if err != nil || !isFinite(ts) {
		return webhookSigErr("Invalid timestamp in " + WebhookSignatureHeader + " header")
	}
	age := math.Abs(float64(time.Now().Unix()) - math.Floor(ts))
	if int(age) > tolerance {
		return webhookSigErr("Timestamp outside tolerance: allowed drift is " + strconv.Itoa(tolerance) + " seconds")
	}
	return nil
}

func parseWebhookEvent(raw []byte) (*WebhookEvent, error) {
	var record map[string]interface{}
	if err := json.Unmarshal(raw, &record); err != nil {
		return nil, webhookSigErr("Webhook payload is not valid JSON")
	}
	if record == nil {
		return nil, webhookSigErr("Webhook payload must be a JSON object")
	}

	id, _ := record["id"].(string)
	if id == "" {
		return nil, webhookSigErr("Webhook payload missing required field: id")
	}

	// Production envelope: {id, type, created_at, data}.
	// Legacy envelope (pre-2.0.1 fixtures): {id, event, payload, timestamp}.
	eventType, _ := record["type"].(string)
	if eventType == "" {
		eventType, _ = record["event"].(string)
	}
	if eventType == "" {
		return nil, webhookSigErr("Webhook payload missing required field: type")
	}

	data, ok := record["data"].(map[string]interface{})
	if !ok {
		data, ok = record["payload"].(map[string]interface{})
	}
	if !ok || data == nil {
		return nil, webhookSigErr("Webhook payload missing required field: data")
	}

	var createdAt string
	var unix float64
	if s, _ := record["created_at"].(string); s != "" {
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return nil, webhookSigErr("Webhook payload has an invalid created_at")
		}
		createdAt = s
		unix = float64(t.UnixMilli()) / 1000
	} else if ts, ok := record["timestamp"].(float64); ok && isFinite(ts) {
		unix = ts
		createdAt = time.UnixMilli(int64(ts * 1000)).UTC().Format("2006-01-02T15:04:05.000Z07:00")
	} else {
		return nil, webhookSigErr("Webhook payload missing required field: created_at")
	}

	return &WebhookEvent{
		ID:        id,
		Type:      eventType,
		CreatedAt: createdAt,
		Data:      data,
		Event:     eventType,
		Payload:   data,
		Timestamp: unix,
	}, nil
}

// VerifyWebhook verifies a webhook signature and parses the event payload.
func VerifyWebhook(params VerifyWebhookParams) (*WebhookEvent, error) {
	if params.Secret == "" {
		return nil, webhookSigErr("Webhook secret is required")
	}

	tolerance := params.Tolerance
	if tolerance == 0 {
		tolerance = defaultWebhookToleranceSeconds
	}

	signatureHeader, ok := getWebhookHeader(params.Headers, WebhookSignatureHeader)
	if !ok || signatureHeader == "" {
		signatureHeader, ok = getWebhookHeader(params.Headers, legacyWebhookSignatureHeader)
	}
	if !ok || signatureHeader == "" {
		return nil, webhookSigErr("Missing " + WebhookSignatureHeader + " header")
	}

	timestamp, signatures, err := parseSignatureHeader(signatureHeader)
	if err != nil {
		return nil, err
	}
	if err := verifyWebhookTimestamp(timestamp, tolerance); err != nil {
		return nil, err
	}

	expected := computeExpectedWebhookSignature(params.Secret, timestamp, params.Payload)
	valid := false
	for _, sig := range signatures {
		if verifyWebhookSignature(expected, sig) {
			valid = true
			break
		}
	}
	if !valid {
		return nil, webhookSigErr("Webhook signature verification failed")
	}

	return parseWebhookEvent(params.Payload)
}

// ConstructEvent verifies a webhook signature from raw payload bytes and header value.
func ConstructEvent(payload []byte, signature string, secret string, tolerance int) (*WebhookEvent, error) {
	headers := map[string]string{}
	if signature != "" {
		headers[WebhookSignatureHeader] = signature
	}
	return VerifyWebhook(VerifyWebhookParams{
		Payload:   payload,
		Headers:   headers,
		Secret:    secret,
		Tolerance: tolerance,
	})
}

// Verify verifies a webhook signature (instance helper).
func (s *WebhookService) Verify(params VerifyWebhookParams) (*WebhookEvent, error) {
	return VerifyWebhook(params)
}
