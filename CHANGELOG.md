# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.3.4...v0.4.0) (2026-04-15)


### Features

* fix ProjectSSO fields, add Domain and ProjectSsoTeam resources ([ae5ebe1](https://github.com/ottercoders/pulumi-oneuptime/commit/ae5ebe1bfb95d76be8967e33d9c66c5479359f64))

## [0.3.4](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.3.3...v0.3.4) (2026-04-15)


### Bug Fixes

* **client:** use correct API key header name (APIKey not ApiKey) ([b9a1bef](https://github.com/ottercoders/pulumi-oneuptime/commit/b9a1befa5b9abacb06723a79a1a0d1407b10486e))
* **client:** wrap create/update bodies in {"data": ...} envelope ([d0cc729](https://github.com/ottercoders/pulumi-oneuptime/commit/d0cc729b3543af8565acd5480a0b4da6886eab66))
* unwrap OneUptime typed objects (DateTime, ObjectID) in API responses ([6d17f23](https://github.com/ottercoders/pulumi-oneuptime/commit/6d17f23114e306fc21c39f656a5650c0c7ed29b2))
* use tenantid header for project scoping and add name to escalation rules ([7e8ea2e](https://github.com/ottercoders/pulumi-oneuptime/commit/7e8ea2e2cd9820fe8396c59237091e1bdbb13b99))

## [0.3.3](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.3.2...v0.3.3) (2026-04-15)


### Bug Fixes

* add debug logging to ToMap for diagnosing empty request bodies ([fd011ef](https://github.com/ottercoders/pulumi-oneuptime/commit/fd011ef355a8246feadc68652914f74e0049e257))

## [0.3.2](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.3.1...v0.3.2) (2026-04-15)


### Bug Fixes

* align json tags with pulumi tags in ProjectSSO ([be0c3a9](https://github.com/ottercoders/pulumi-oneuptime/commit/be0c3a9be7f17079588ef69424976ac365d1c560))

## [0.3.1](https://github.com/ottercoders/pulumi-oneuptime/compare/v0.3.0...v0.3.1) (2026-04-15)


### Bug Fixes

* **ci:** add repository.url to npm package.json for provenance ([cb718b1](https://github.com/ottercoders/pulumi-oneuptime/commit/cb718b1d0c8d66d5450fbf208b2315b393715d11))
* **ci:** inline release jobs into release-please workflow ([fdb98f9](https://github.com/ottercoders/pulumi-oneuptime/commit/fdb98f9536833ecde548debce88da55f60782bb3))

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
