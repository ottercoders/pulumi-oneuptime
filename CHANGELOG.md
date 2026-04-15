# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.2.1...v0.3.0) (2026-04-15)


### Features

* add ProjectSSO resource for OIDC/SAML configuration ([3dad362](https://github.com/ottercoders/pulumi-oneuptime/commit/3dad362c719ffffa63f891c99308f5b94aa77312))


### Bug Fixes

* **ci:** pass id-token permission to release workflow call ([e1decc2](https://github.com/ottercoders/pulumi-oneuptime/commit/e1decc270060329dba6bf9772466c767b2690ae1))
* **ci:** trigger release workflow when release-please creates a tag ([7bb01d3](https://github.com/ottercoders/pulumi-oneuptime/commit/7bb01d369e4919c577756a2d0885e2f42f734271))

## [0.2.1](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.2.0...v0.2.1) (2026-04-15)


### Bug Fixes

* **client:** add debug logging and improve error messages ([a448814](https://github.com/ottercoders/pulumi-oneuptime/commit/a448814141e521ae94ee6f251d7e3d105bc66d0e))

## [0.2.0](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.1.6...v0.2.0) (2026-04-15)


### Features

* add pluginDownloadURL for automatic provider binary download ([e69314e](https://github.com/ottercoders/pulumi-oneuptime/commit/e69314eeb107b385badd49545a87432d196de69e))
* add StatusPageGroup, StatusPageResource, and on-call escalation resources ([7d8db73](https://github.com/ottercoders/pulumi-oneuptime/commit/7d8db730ac8dd0f2210d880d74afb628f49388c8))
* **config:** add configurable HTTP timeout ([50f3d1d](https://github.com/ottercoders/pulumi-oneuptime/commit/50f3d1d3de17e78dcd78303fbf0c8207dfbbe253))

## [0.1.2](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.1.1...v0.1.2) (2026-04-15)


### Bug Fixes

* **ci:** add --provenance flag for npm trusted publishing ([128fad3](https://github.com/ottercoders/pulumi-oneuptime/commit/128fad328522ab5f94e00c702a54c623cfac367f))

## [Unreleased]

### Added

- Initial provider implementation with `pulumi-go-provider` v1.3.1
- Provider configuration: `apiKey`, `baseUrl`, `projectId` with environment variable fallbacks
- Generic HTTP client for OneUptime REST API
- Resources:
  - `Team` — manage teams within a project
  - `Monitor` — manage monitors (API, website, IP, etc.)
  - `StatusPage` — manage public/private status pages
  - `Incident` — manage incidents
  - `OnCallDutyPolicy` — manage on-call duty escalation policies
- Full CRUD support for all resources
- Unit tests for the HTTP client
- YAML example program
- CI/CD with GitHub Actions (test on PR, release on tag)
