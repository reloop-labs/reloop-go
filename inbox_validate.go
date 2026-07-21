package reloop

import (
	"math"
)

const (
	inboxLimitMin    = 1
	inboxLimitMax    = 200
	inboxOffsetMin   = 0
	inboxBatchIDsMax = 100
)

func requireMailboxID(id, field string) (string, error) {
	if field == "" {
		field = "id"
	}
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Mailbox %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireMessageID(id, field string) (string, error) {
	if field == "" {
		field = "id"
	}
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Message %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireThreadID(id, field string) (string, error) {
	if field == "" {
		field = "id"
	}
	trimmed, err := requireNonEmptyString(id, field)
	if err != nil {
		return "", validationErrorf(field, "Thread %s is required and must be a non-empty string.", field)
	}
	return trimmed, nil
}

func requireInboxAttachmentID(id, field string) (string, error) {
	return requireNonEmptyString(id, field)
}

func requireInboxLimit(limit int, field string) error {
	if field == "" {
		field = "limit"
	}
	return requireLimit(limit, inboxLimitMin, inboxLimitMax, field)
}

func requireInboxOffset(offset int, field string) error {
	if field == "" {
		field = "offset"
	}
	if offset < inboxOffsetMin {
		return validationErrorf(field, "list %s must be an integer >= %d.", field, inboxOffsetMin)
	}
	return nil
}

func requireInboxIDArray(ids []string, field string, max int) ([]string, error) {
	if field == "" {
		field = "ids"
	}
	if len(ids) == 0 {
		return nil, validationErrorf(field, "%s is required and must be a non-empty array.", field)
	}
	if len(ids) > max {
		return nil, validationErrorf(field, "%s must contain at most %d items.", field, max)
	}
	out := make([]string, len(ids))
	for i, id := range ids {
		trimmed, err := requireNonEmptyString(id, field)
		if err != nil {
			return nil, validationErrorf(field, "%s[%d] must be a non-empty string.", field, i)
		}
		out[i] = trimmed
	}
	return out, nil
}

func requireMailboxStatus(status MailboxStatus, field string) error {
	if field == "" {
		field = "status"
	}
	switch status {
	case MailboxStatusActive, MailboxStatusDisabled:
		return nil
	default:
		return validationErrorf(field, `%s must be "active" or "disabled".`, field)
	}
}

func requireThreadStatusValue(status ThreadStatus, field string) error {
	if field == "" {
		field = "status"
	}
	switch status {
	case ThreadStatusActive, ThreadStatusArchived, ThreadStatusClosed, ThreadStatusTrash:
		return nil
	default:
		return validationErrorf(field, `%s must be "active", "archived", "closed", or "trash".`, field)
	}
}

func requireThreadFilterValue(filter ThreadFilter, field string) error {
	if field == "" {
		field = "filter"
	}
	switch filter {
	case ThreadFilterPrimary, ThreadFilterAlerts, ThreadFilterPerson, ThreadFilterTag:
		return nil
	default:
		return validationErrorf(field, `%s must be "primary", "alerts", "person", or "tag".`, field)
	}
}

func requireThreadBatchAction(action ThreadBatchAction) error {
	switch action {
	case ThreadBatchActionArchive, ThreadBatchActionTrash, ThreadBatchActionRestore,
		ThreadBatchActionStar, ThreadBatchActionUnstar, ThreadBatchActionRead, ThreadBatchActionUnread,
		ThreadBatchActionImportant, ThreadBatchActionUnimportant, ThreadBatchActionSpam, ThreadBatchActionUnspam,
		ThreadBatchActionPin, ThreadBatchActionUnpin:
		return nil
	default:
		return validationErrorf("action", "batch action must be a valid thread batch action.")
	}
}

func requireComposeBody(text, html *string) error {
	hasText := text != nil
	hasHTML := html != nil
	if !hasText && !hasHTML {
		return validationErrorf("params", "at least one of text or html is required.")
	}
	return nil
}

func requireFiniteNumber(v float64, field string) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return validationErrorf(field, "%s must be a number when provided.", field)
	}
	return nil
}
