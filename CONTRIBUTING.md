# Contributing to the Reloop Go SDK

Module: **`github.com/reloop-labs/reloop-go`**.

**License:** [Apache License 2.0](./LICENSE) with additional use restrictions from Reloop Labs.

**API reference:** [reloop.sh/docs](https://reloop.sh/docs)

Port new endpoints from the [Node.js SDK](https://github.com/reloop-labs/reloop-node) reference.

---

## Development setup

```bash
git clone git@github.com:reloop-labs/reloop-go.git
cd reloop-go
go test ./...
```

Requires **Go 1.20+**.

---

## Project layout

```
client.go
mail.go
domain.go
models.go
version.go            # const Version
*_test.go             # httptest route tests
go.mod
```

---

## Conventions

| Topic | Rule |
|-------|------|
| Mail & domain structs | `json` tags in **snake_case** |
| Optional fields | Pointers (`*string`, `*bool`) with `omitempty` |
| Tests | `httptest.Server`; assert path and JSON body |
| Version | Update `version.go` and tag together |

---

## Pull request checklist

- [ ] `go test ./...` passes
- [ ] `version.go` updated only when releasing

---

## Releasing

Version: **`version.go`** → `const Version = "1.9.0"`.

```bash
git commit -am "chore: release v1.9.0"
git push origin main
git tag v1.9.0
git push origin v1.9.0
```

[`.github/workflows/release.yml`](./.github/workflows/release.yml) creates a GitHub Release with a source zip.

Consumers install via: `go get github.com/reloop-labs/reloop-go@v1.9.0`
