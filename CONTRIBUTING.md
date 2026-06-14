# Contributing & releasing Reloop SDKs

This guide applies to all official Reloop language SDKs:

| Repository | Package | Version file |
|------------|---------|--------------|
| `reloop-node` | `reloop-email` (npm) | `package.json` |
| `reloop-python` | `reloop-email` (PyPI) | `pyproject.toml` |
| `reloop-php` | `reloop/reloop-email` (Composer) | `composer.json` |
| `reloop-go` | `github.com/reloop-labs/reloop-go` | `version.go` |
| `reloop-java` | `sh.reloop:reloop-java` (Maven) | `pom.xml` |
| `reloop-dotnet` | `Reloop` (NuGet) | `Reloop.csproj` |
| `reloop-rust` | `reloop` (crates.io) | `Cargo.toml` |
| `reloop-ruby` | `reloop` (RubyGems) | `reloop.gemspec` |
| `reloop-elixir` | `reloop` (Hex) | `mix.exs` |

The same guide is maintained in every Reloop SDK repository.

All SDKs are licensed under the **Apache License 2.0** with additional use restrictions from Reloop Labs. See [`LICENSE`](./LICENSE) in each repository (same text as the [main Reloop project](https://github.com/reloop-labs/reloop/blob/main/LICENSE)).

---

## Contributing

### Before you start

1. Read the API reference at [reloop.sh/docs](https://reloop.sh/docs).
2. Use **Node.js** (`reloop-node`) as the reference implementation for new endpoints and behaviour.
3. Keep **the same version number** across all SDKs when shipping a coordinated release (see [Versioning](#versioning)).

### What to change

When adding or updating an API (for example mail send):

1. Implement the service in **Node.js** first.
2. Port the same routes, request shape, and response fields to the other SDKs.
3. Add **route tests** (mock HTTP client; assert path, method, and body).
4. Update the SDK **README** with prerequisites (API key + verified domain), a minimal send example, and a link to [reloop.sh/docs](https://reloop.sh/docs).
5. Update backend **x-codeSamples** and docs sync if the public API docs should show SDK snippets.

### Conventions

| Area | Convention |
|------|------------|
| Mail, domain | **snake_case** in JSON request bodies (`reply_to`, `click_tracking`) |
| Contacts, API keys | camelCase in JSON via each SDK’s `for_request` helper |
| Node.js | `{ response, error }` result type — no throws for HTTP errors |
| Python | `from_` instead of `from` for mail sender |
| PHP / Ruby / Elixir | snake_case hash keys in user-facing APIs |
| Go / Java / .NET / Rust | Typed params or maps matching OpenAPI field names |

### Pull request checklist

- [ ] Tests pass locally for the SDK you changed.
- [ ] New endpoints have route/body tests consistent with existing SDK tests.
- [ ] README stays minimal (send email example + link to docs only).
- [ ] Version bumped in the package manifest **only when preparing a release** (see below).
- [ ] No secrets, API keys, or `.env` files committed.

### Running tests

```bash
# Node
cd reloop-node && npm test

# Python
cd reloop-python && python -m unittest discover -s tests -v

# Go
cd reloop-go && go test ./...

# PHP
cd reloop-php && composer install && vendor/bin/phpunit

# .NET
cd reloop-dotnet && dotnet test tests/Reloop.Tests/Reloop.Tests.csproj

# Rust
cd reloop-rust && cargo test

# Ruby
cd reloop-ruby && ruby -Itest test/*_test.rb

# Java
cd reloop-java && mvn verify

# Elixir
cd reloop-elixir && mix test
```

---

## Releasing

Reloop SDKs follow common industry practice:

- **[Semantic Versioning](https://semver.org)** (`MAJOR.MINOR.PATCH`)
- **Git tags** drive GitHub Releases (`v1.8.0`)
- **Package manifest version** must match the tag
- **Registry publish** (npm, PyPI, etc.) is separate from but usually done alongside the GitHub Release

### Versioning

| Change type | Version bump | Example |
|-------------|--------------|---------|
| Breaking API change | MAJOR | `1.8.0` → `2.0.0` |
| New feature, backward compatible | MINOR | `1.8.0` → `1.9.0` |
| Bug fix only | PATCH | `1.8.0` → `1.8.1` |

When shipping a feature across all languages (e.g. mail send), bump **every SDK** to the same version before tagging each repository.

**Version locations:**

```
reloop-node/package.json          → "version"
reloop-python/pyproject.toml      → project.version
reloop-php/composer.json          → "version"
reloop-go/version.go              → const Version
reloop-java/pom.xml               → <version>
reloop-dotnet/Reloop.csproj       → <Version>
reloop-rust/Cargo.toml            → version
reloop-ruby/reloop.gemspec        → spec.version
reloop-elixir/mix.exs             → version:
```

Also update `reloop-node/package-lock.json` and `reloop-rust/Cargo.lock` (reloop crate entry only) when those versions change.

---

### Release workflow (per repository)

Each SDK has `.github/workflows/release.yml` that runs on:

- Push of a tag matching `v*` (primary)
- Manual **workflow_dispatch** (backup)

The workflow will:

1. Read the version from the package manifest
2. Verify the git tag matches that version (e.g. tag `v1.8.0` ↔ manifest `1.8.0`)
3. Run tests and build
4. Create a **source zip** (e.g. `reloop-node-1.8.0-source.zip`)
5. Attach **built artifacts** where applicable (npm tarball, wheel, gem, nupkg, jar, etc.)
6. Create a **GitHub Release** with auto-generated notes

Registry publishing (npm, PyPI, NuGet, …) is handled by separate `publish.yml` workflows where configured. Run those after the release succeeds or wire them into the same pipeline when ready.

---

### Step-by-step: cut a release

Do this **in each SDK repository** you want to ship.

#### 1. Bump the version

Set the new version in that repo’s manifest (example: `1.9.0`). Commit:

```bash
git add package.json package-lock.json   # example for Node
git commit -m "chore: release v1.9.0"
git push origin main
```

#### 2. Create and push the tag

The tag **must** match the manifest version with a `v` prefix:

```bash
git tag v1.9.0
git push origin v1.9.0
```

#### 3. Confirm the GitHub Action

Open **Actions → Release** in the repository. The workflow should:

- Pass tests
- Upload assets to **Releases**
- Show a new entry such as `v1.9.0`

#### 4. Publish to the package registry (if applicable)

| SDK | Registry | Typical command / workflow |
|-----|----------|----------------------------|
| Node | npm | `publish.yml` or `npm publish` |
| Python | PyPI | `publish.yml` / `hatch publish` |
| PHP | Packagist | Tag on GitHub (Packagist mirrors tags) |
| Go | — | Tag only; `go get` uses module tags |
| Java | GitHub Packages / Maven Central | `publish.yml` |
| .NET | NuGet | `publish.yml` |
| Rust | crates.io | `cargo publish` |
| Ruby | RubyGems | `gem push` |
| Elixir | Hex | `mix hex.publish` |

Ensure registry secrets (`NPM_TOKEN`, `PYPI`, `NUGET_API_KEY`, etc.) are set in **Repository → Settings → Secrets**.

---

### Coordinated multi-SDK release

When releasing one API version across all languages:

1. Merge all SDK changes to each repo’s `main`.
2. Bump every manifest to the **same version** (e.g. `1.9.0`).
3. Tag and push **`v1.9.0`** in **each** repository:

   ```bash
   for repo in reloop-node reloop-python reloop-php reloop-go reloop-java \
               reloop-dotnet reloop-rust reloop-ruby reloop-elixir; do
     cd "$repo"
     git tag v1.9.0
     git push origin v1.9.0
     cd ..
   done
   ```

4. Verify GitHub Releases and registry packages for each language.

---

### Pre-release and hotfix tags

| Type | Tag example | When |
|------|-------------|------|
| Beta | `v1.9.0-beta.1` | Early testers; avoid for production docs |
| Release candidate | `v1.9.0-rc.1` | Final testing before stable |
| Patch hotfix | `v1.8.1` | Bug fix on latest stable line |

Do **not** delete or move tags that are already published. Ship a new patch version instead.

---

### Troubleshooting

| Problem | Fix |
|---------|-----|
| Workflow fails: tag ≠ manifest version | Align tag and manifest, delete local tag, retag after fix |
| Release already exists | Bump patch version or delete the draft release only if nothing was published |
| Tests pass locally, fail in CI | Match CI runtime versions (Node 22, Go 1.20, Python 3.10, etc.) |
| npm publish works, no GitHub Release | Ensure `release.yml` ran on the same tag push |

---

## Questions

- **API behaviour:** [reloop.sh/docs](https://reloop.sh/docs)
- **Backend / OpenAPI:** `reloop` monorepo (`apps/backend`)
- **SDK samples in docs:** `generate_code_samples.py` + `sync:code-samples`

For larger changes (new services, breaking renames), open an issue or discuss before implementing across all nine SDKs.
