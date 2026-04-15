# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.5](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.1.4...v0.1.5) (2026-04-15)


### Bug Fixes

* **ci:** switch npm back to OIDC trusted publishing ([c665efd](https://github.com/ottercoders/pulumi-oneuptime/commit/c665efd6d7e621c170cb302b465e84f76cb25f37))

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
