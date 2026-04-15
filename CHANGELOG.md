# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.0 (2026-04-15)


### Features

* add release-please for automated releases ([1249cc2](https://github.com/ottercoders/pulumi-oneuptime/commit/1249cc2e442c9e642fac28aa5ae540915c94ea48))

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
