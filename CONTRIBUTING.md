# Contributing to pulumi-oneuptime

Thank you for your interest in contributing to the Pulumi OneUptime provider!

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/<your-username>/pulumi-oneuptime.git
   cd pulumi-oneuptime
   ```
3. Build the provider:
   ```bash
   make provider
   ```
4. Run tests:
   ```bash
   make test
   ```

## Development Workflow

### Prerequisites

- Go 1.22+
- Pulumi CLI
- A OneUptime account (for acceptance tests)

### Building

```bash
make provider    # Build the provider binary to bin/
make install     # Install to $GOPATH/bin for local testing
```

### Testing

**Unit tests** (no credentials needed):

```bash
make test
```

**Acceptance tests** (requires a OneUptime instance):

```bash
export ONEUPTIME_API_KEY="your-key"
export ONEUPTIME_PROJECT_ID="your-project-id"
make test-acceptance
```

### Adding a New Resource

1. Create `provider/resources/<resource_name>.go`
2. Define the struct types:
   - `<Resource>` — empty struct, the resource controller
   - `<Resource>Args` — input properties with `pulumi:"..."` and `json:"..."` tags
   - `<Resource>State` — output properties, embedding `Args` plus computed fields
3. Implement CRUD methods: `Create`, `Read`, `Update`, `Delete`
4. Register the resource in `provider/provider.go`
5. Add unit tests
6. Update `README.md` with the new resource

Use existing resources (e.g., `team.go`) as a template.

## Pull Request Process

1. Create a feature branch from `main`
2. Make your changes
3. Ensure `make test` passes
4. Ensure `make provider` builds successfully
5. Open a pull request with a clear description of the changes

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use dual struct tags: `pulumi:"fieldName"` for Pulumi schema, `json:"fieldName"` for API serialization
- Keep resource implementations consistent with the existing pattern
- Add tests for new functionality

## Reporting Issues

Please open an issue on GitHub with:

- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Provider version and Pulumi CLI version
