# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Native Pulumi provider for OneUptime, built with `pulumi-go-provider` v1.3.1 (infer pattern). Exposes 41 resources and 10 lookup functions for managing OneUptime monitoring infrastructure.

## Build & Test

```bash
make provider         # Build the provider binary to bin/
make test             # Run all unit tests
make schema           # Generate schema.json (requires pulumi CLI)
make sdk              # Generate TypeScript, Python, Go, .NET SDKs
make install          # Copy binary to $GOPATH/bin for local testing
go test ./provider/resources/ -run TestTeamResource_LifeCycle -v  # Run a single test
```

Debug mode: set `ONEUPTIME_DEBUG=1` to log all HTTP request/response bodies to stderr.

## Commit Conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/) and [release-please](https://github.com/googleapis/release-please). All commits must follow `<type>(<scope>): <description>`.

- `feat:` — new resource or function (minor bump)
- `fix:` — bug fix (patch bump)
- `feat!:` or `BREAKING CHANGE:` in body — major bump
- `chore:`, `docs:`, `ci:`, `refactor:`, `test:` — no version bump

Scope is optional: `monitor`, `client`, `config`, `ci`, `sdk`

## Architecture

```
provider/
├── cmd/pulumi-resource-oneuptime/main.go  # Binary entry point
├── provider.go          # Provider builder: registers all resources + functions
├── client/
│   └── client.go        # Generic HTTP client for OneUptime REST API
└── resources/
    ├── config.go        # Provider config (apiKey, baseUrl, projectId, timeout)
    ├── helpers.go       # ToMap, FromMap, SelectFields, ResolveProjectID, unwrapTypedValues
    ├── lookups.go       # 10 invoke functions (getMonitorStatus, getTeam, etc.)
    ├── team.go          # Canonical resource example
    ├── api_key.go       # Example with secret output (WireDependencies)
    └── *.go             # ~40 resource files, all following the same pattern
```

### OneUptime API conventions

- **Auth header:** `APIKey` (capital I, not `ApiKey`)
- **Project scoping:** `tenantid` header (set from provider config `projectId`)
- **Create/Update body:** wrapped in `{"data": {...}}`
- **Read:** `POST /api/{resource}/{id}/get-item` with `{"select": {...}}`
- **List:** `POST /api/{resource}/get-list` with `{"query": {...}, "select": {...}, "limit": N}`
- **Typed response values:** API returns `{"_type": "DateTime", "value": "..."}` — `FromMap` unwraps these automatically via `unwrapTypedValues`
- **Deleted resource reads:** API returns `200 {}` not `404` — client treats `{}` as nil

### Struct tag rules

Every Args/State field needs dual tags. **The `pulumi` and `json` tag values must match** — the Pulumi mapper checks `json` tags first, and mismatches cause silent field drops.

```go
Name  string  `pulumi:"name" json:"name"`              // ✓ matching tags
Color string  `pulumi:"color" json:"color"`             // ✓ matching tags
```

The only exception is `ResourceID` on State structs, which maps `_id` from the API:
```go
ResourceID string `pulumi:"resourceId" json:"_id"`      // ✓ output-only, not affected by input mapping
```

## Adding a Resource

1. Create `provider/resources/<name>.go` — copy `team.go` as template
2. Define `XxxArgs` (inputs) and `XxxState` (embeds Args + `ResourceID`, `CreatedAt`, `UpdatedAt`)
3. Implement `Create`, `Read`, `Update`, `Delete` methods
4. The output ID field must be `resourceId` (Pulumi reserves `id`)
5. API path is kebab-case: `on-call-duty-policy`, `status-page-group`
6. Register in `provider/provider.go` under `WithResources(...)`
7. For secret outputs, implement `infer.ExplicitDependencies` with `WireDependencies` — see `api_key.go`

## Adding a Lookup Function

1. Add to `provider/resources/lookups.go`
2. If it's a standard name-based lookup, reuse `lookupByName()` helper and `LookupByNameArgs`/`LookupResult` types
3. For custom result shapes (like `GetMonitorStatus` with `color`), define separate Args/Result structs
4. Register in `provider/provider.go` under `WithFunctions(...)`

## Release flow

Push to `main` → release-please creates PR → merge PR → release-please tags → same workflow runs GoReleaser + npm publish (OIDC trusted publishing, Node 24) + PyPI publish (OIDC) + attaches artifacts to GitHub Release.

The npm package publishes compiled JS from `sdk/nodejs/nodejs/bin/` (not raw `.ts`). The `jq` step in the workflow sets package name, version, `main`, `types`, and `repository.url`.

Manual release: use workflow_dispatch on the Release workflow, or tag directly.
