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
	// WebhookSignatureHeader is the HTTP header carrying the signed payload.
	WebhookSignatureHeader = "x-webhook-signature"
	// WebhookTimestampHeader is the HTTP header carrying the event timestamp.
	WebhookTimestampHeader = "x-webhook-timestamp"

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

// WebhookEvent is a verified webhook payload.
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Event     string                 `json:"event"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp float64                `json:"timestamp"`
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
		return "", nil, webhookSigErr("Invalid X-Webhook-Signature header: expected t= and v1= values")
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
		return webhookSigErr("Invalid timestamp in X-Webhook-Signature header")
	}
	age := math.Abs(float64(time.Now().Unix()) - math.Floor(ts))
	if int(age) > tolerance {
		return webhookSigErr("Timestamp outside tolerance: allowed drift is " + strconv.Itoa(tolerance) + " seconds")
	}
	return nil
}

func parseWebhookEvent(raw string) (*WebhookEvent, error) {
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return nil, webhookSigErr("Webhook payload is not valid JSON")
	}
	if record == nil {
		return nil, webhookSigErr("Webhook payload must be a JSON object")
	}

	id, _ := record["id"].(string)
	if id == "" {
		return nil, webhookSigErr("Webhook payload missing required field: id")
	}
	event, _ := record["event"].(string)
	if event == "" {
		return nil, webhookSigErr("Webhook payload missing required field: event")
	}
	payload, ok := record["payload"].(map[string]interface{})
	if !ok || payload == nil {
		return nil, webhookSigErr("Webhook payload missing required field: payload")
	}
	ts, ok := record["timestamp"].(float64)
	if !ok || !isFinite(ts) {
		return nil, webhookSigErr("Webhook payload missing required field: timestamp")
	}

	return &WebhookEvent{
		ID:        id,
		Event:     event,
		Payload:   payload,
		Timestamp: ts,
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

	return parseWebhookEvent(string(params.Payload))
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
