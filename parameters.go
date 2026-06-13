package reloop

import (
	"regexp"
	"strings"
)

var camelCasePattern = regexp.MustCompile(`([a-z0-9])([A-Z])`)

var requestKeyMap = map[string]string{
	"first_name":           "firstName",
	"last_name":            "lastName",
	"group_ids":            "groupIds",
	"group_id":             "groupId",
	"fallback_value":       "fallbackValue",
	"default_subscription": "defaultSubscription",
	"channel_id":           "channelId",
	"property_name":        "propertyName",
	"property_type":        "propertyType",
	"contact_id":           "contactId",
	"rate_limit_enabled":   "rateLimitEnabled",
	"user_id":              "userId",
}

func Int(value int) *int {
	return &value
}

func Bool(value bool) *bool {
	return &value
}

func String(value string) *string {
	return &value
}

func DomainTLS(value DomainTLSMode) *DomainTLSMode {
	return &value
}

func forRequest(parameters map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(parameters))

	for key, value := range parameters {
		if key == "unsubscribed" {
			if _, ok := parameters["status"]; !ok {
				if unsubscribed, ok := value.(bool); ok {
					if unsubscribed {
						normalized["status"] = "unsubscribed"
					} else {
						normalized["status"] = "subscribed"
					}
				}
			}
			continue
		}

		apiKey := requestKeyMap[key]
		if apiKey == "" {
			apiKey = toCamelCase(key)
		}

		normalized[apiKey] = normalizeValue(value, true)
	}

	return normalized
}

func forQuery(options map[string]interface{}) map[string]interface{} {
	return forRequest(options)
}

func normalizeValue(value interface{}, isRequest bool) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		if isRequest {
			return forRequest(typed)
		}
		return typed
	case []map[string]interface{}:
		items := make([]interface{}, len(typed))
		for index, item := range typed {
			items[index] = normalizeValue(item, isRequest)
		}
		return items
	case []interface{}:
		items := make([]interface{}, len(typed))
		for index, item := range typed {
			items[index] = normalizeValue(item, isRequest)
		}
		return items
	default:
		return value
	}
}

func toCamelCase(key string) string {
	if mapped, ok := requestKeyMap[key]; ok {
		return mapped
	}
	if !strings.Contains(key, "_") {
		return key
	}

	parts := strings.Split(key, "_")
	for index := 1; index < len(parts); index++ {
		if parts[index] == "" {
			continue
		}
		parts[index] = strings.ToUpper(parts[index][:1]) + parts[index][1:]
	}

	return strings.Join(parts, "")
}

func toSnakeCase(key string) string {
	if strings.Contains(key, "_") {
		return key
	}
	return camelCasePattern.ReplaceAllString(key, "${1}_${2}")
}
