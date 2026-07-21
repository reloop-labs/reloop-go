package reloop

import "fmt"

// ValidationError is returned when client arguments are invalid (no HTTP call).
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// APIErrorBody is the decoded Reloop error payload.
type APIErrorBody struct {
	Message string `json:"message,omitempty"`
	Why     string `json:"why,omitempty"`
	Tip     string `json:"tip,omitempty"`
	Link    string `json:"link,omitempty"`
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
