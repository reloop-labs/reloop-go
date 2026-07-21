# Contributing to the Reloop Go SDK

Module: **`github.com/reloop-labs/reloop-go/v2`**.

Keep the package **`Version`** in lockstep with Node (`reloop`) and Python (`reloop-email`) major/minor releases.

**License:** [Apache License 2.0](./LICENSE) with additional use restrictions from Reloop Labs.

**API reference:** [reloop.sh/docs](https://reloop.sh/docs)

Port new endpoints from the [Node.js SDK](https://github.com/reloop-labs/reloop-node) — Node wins on paths, bodies, and validation messages.

---

## Development setup

```bash
git clone git@github.com:reloop-labs/reloop-go.git
cd reloop-go
go test ./...
```

Requires **Go 1.20+**. Stdlib only — do not add third-party dependencies.

---

## Project layout

```
client.go / client_test.go   # NewClient, do/request, headers
errors.go                    # ValidationError, APIError
validate.go                  # shared field validators
version.go                   # const Version = "2.0.0"
mail.go / api_key.go / domain.go
contacts*.go                 # Contacts + Properties / Groups / Channels
webhook.go / webhook_verify.go
inbox*.go                    # Inbox + Mailboxes / Messages / Threads
*_test.go                    # httptest route + surface-lock tests
go.mod                       # module .../reloop-go/v2
```

Single package `reloop` — service files with colocated tests.

---

## Conventions

| Topic | Rule |
|-------|------|
| Return shape | `(*T, error)` — `*ValidationError` / `*APIError` |
| Wire JSON | Match Node (`snake_case` vs camelCase per endpoint) |
| Optional fields | Pointers (`*string`, `*bool`) + `omitempty` helpers `String`/`Bool`/`Int` |
| Validation | Fail before HTTP; tests assert **0** requests |
| Tests | `httptest.Server`: method, path, `x-api-key`, body; surface-lock method sets |
| Endpoints | Do not invent — copy from Node `paths.ts` |
| Version | Same semver as Node/Python for the parity release line |

---

## Pull request checklist

- [ ] `go test ./...` passes
- [ ] New ops have happy path + API error + validation (0 HTTP) + surface test
- [ ] `version.go` updated only when releasing (match other SDKs)

---

## Releasing

1. Set `version.go` → `const Version = "X.Y.Z"` (same as Node/Python)
2. Tag `vX.Y.Z` (module path includes `/v2` for major 2)
3. Push tag; consumers `go get github.com/reloop-labs/reloop-go/v2@vX.Y.Z`
