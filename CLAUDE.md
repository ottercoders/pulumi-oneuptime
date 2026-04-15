# CLAUDE.md

## Project

Native Pulumi provider for OneUptime, built with `pulumi-go-provider` (infer pattern).

## Build & Test

```bash
make provider    # Build the provider binary
make test        # Run unit tests
make schema      # Generate schema.json (requires pulumi CLI)
make sdk         # Generate language SDKs
```

## Commit Conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/) and [release-please](https://github.com/googleapis/release-please) for automated releases. **All commits must follow this format:**

```
<type>(<scope>): <description>
```

### Types

- `feat:` — new feature or resource (triggers **minor** version bump)
- `fix:` — bug fix (triggers **patch** version bump)
- `docs:` — documentation only
- `chore:` — maintenance, CI, dependencies
- `refactor:` — code change that neither fixes a bug nor adds a feature
- `test:` — adding or updating tests
- `ci:` — CI/CD workflow changes

### Breaking Changes

Append `!` after the type or add `BREAKING CHANGE:` in the commit body to trigger a **major** version bump:

```
feat!: rename resourceId output field to oneuptimeId
```

### Scope (optional)

Use the resource or area name: `monitor`, `client`, `config`, `ci`, `sdk`

### Examples

```
feat(monitor): add monitorSteps field for multi-step monitors
fix(client): handle 429 rate limit responses with retry
chore(ci): update goreleaser to v2
docs: add Go SDK usage example to README
feat!: require projectId on all resources, remove provider-level default
```

## Adding a Resource

1. Create `provider/resources/<name>.go` with `Args`, `State`, and CRUD methods
2. Use dual struct tags: `pulumi:"fieldName"` and `json:"apiFieldName"`
3. Register in `provider/provider.go`
4. The output ID field must be named `resourceId` (not `id` — Pulumi reserves that)
5. Use the `ToMap`/`FromMap`/`SelectFields` helpers from `helpers.go`
6. API path is the kebab-case resource name (e.g., `on-call-duty-policy`)

## Architecture

- `provider/client/` — Generic HTTP client for OneUptime REST API
- `provider/resources/` — Resource definitions, config, and helpers
- `provider/provider.go` — Provider registration and wiring
- All resources share the same CRUD pattern via the generic client
