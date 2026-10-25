package reloop

import (
	"encoding/json"
	"fmt"
)

// ValidationError is returned when client arguments are invalid (no HTTP call).
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// APIErrorBody is the decoded Reloop error payload: message, why, fix, link.
type APIErrorBody struct {
	Message string `json:"message,omitempty"`
	Why     string `json:"why,omitempty"`
	Fix     string `json:"fix,omitempty"`
	Link    string `json:"link,omitempty"`
	// Deprecated: the API returns "fix". Tip mirrors Fix.
	Tip string `json:"tip,omitempty"`
}

// UnmarshalJSON reads "fix" (and the legacy "tip") into both Fix and Tip.
func (b *APIErrorBody) UnmarshalJSON(data []byte) error {
	type rawBody APIErrorBody
	var raw rawBody
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*b = APIErrorBody(raw)
	if b.Fix == "" {
		b.Fix = b.Tip
	}
	b.Tip = b.Fix
	return nil
}

// APIError is returned for non-2xx HTTP responses and network failures.
type APIError struct {
	Status     int
	StatusText string
	Body       APIErrorBody
}

func (e *APIError) Error() string {
	if e.Body.Message != "" {
		return e.Body.Message
	}
	if e.Status == 0 {
		return fmt.Sprintf("Reloop network error: %s", e.StatusText)
	}
	return fmt.Sprintf("Reloop API Error: %d %s", e.Status, e.StatusText)
}

func validationErrorf(field, format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		Message: fmt.Sprintf(format, args...),
		Field:   field,
	}
}
