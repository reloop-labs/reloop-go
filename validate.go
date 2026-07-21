package reloop

import (
	"strings"
	"unicode/utf8"
)

func requireNonEmptyString(value, field string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", validationErrorf(field, "%s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireMailString(value, field string) (string, error) {
	return requireNonEmptyString(value, field)
}

func requireRecipient(value interface{}, field string) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return requireMailString(v, field)
	case []string:
		if len(v) == 0 {
			return nil, validationErrorf(field, "%s must contain at least one address.", field)
		}
		out := make([]string, len(v))
		for i, entry := range v {
			s, err := requireMailString(entry, field)
			if err != nil {
				return nil, err
			}
			out[i] = s
		}
		return out, nil
	default:
		return nil, validationErrorf(field, "%s is required and must be a string or string array.", field)
	}
}

func requireAPIKeyName(name, field string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", validationErrorf(field, "API key %s must be at least 1 character.", field)
	}
	if utf8.RuneCountInString(trimmed) > 255 {
		return "", validationErrorf(field, "API key %s must be at most 255 characters.", field)
	}
	return trimmed, nil
}

func requireAPIKeyID(id, field string) (string, error) {
	return requireNonEmptyString(id, field)
}

func requirePage(page int, field string) error {
	if page < 1 {
		return validationErrorf(field, "list %s must be an integer >= 1.", field)
	}
	return nil
}

func requireLimit(limit, min, max int, field string) error {
	if limit < min || limit > max {
		return validationErrorf(field, "list %s must be an integer between %d and %d.", field, min, max)
	}
	return nil
}
