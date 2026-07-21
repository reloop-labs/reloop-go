# Changelog

## 2.0.0

Breaking rewrite for Node SDK parity (version aligned with Node/Python `2.0.0`):

- Module path `github.com/reloop-labs/reloop-go/v2`
- Idiomatic `(T, error)` returns: `*ValidationError` (bad args, no HTTP) and `*APIError` (HTTP/network)
- Facade: `client.APIKey` (was `ApiKeys` / `APIKeyService` alias), `Mail`, `Domain`, `Contacts`, `Webhook`, `Inbox`
- Typed Contacts with nested `Properties`, `Groups`, `Channels` (no flat helpers)
- Full Webhook CRUD + pause/enable/disable/trigger/deliveries + local HMAC `Verify` / `ConstructEvent`
- Full Inbox: `Mailboxes`, `Messages`, `Threads`
- Removed: API key `Pause`, domain `GetNameservers` / `ForwardDNS`
- Paths aligned with Node (`/api/mail/v1/send`, `/api/api-key/v1/`, `/api/domain/v1/`, …)
- Stdlib-only; httptest tests with surface locks
