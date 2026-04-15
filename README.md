# Pulumi OneUptime Provider

A native [Pulumi](https://www.pulumi.com/) provider for managing [OneUptime](https://oneuptime.com/) monitoring resources.

## Installation

### Go

```bash
go get github.com/ottercoders/pulumi-oneuptime/sdk/go/oneuptime
```

### Node.js

```bash
npm install @ottercoders/pulumi-oneuptime
```

### Python

```bash
pip install pulumi-oneuptime
```

## Configuration

The provider requires an API key and optionally a base URL and default project ID.

| Property    | Environment Variable     | Required | Default                  | Description                          |
|-------------|--------------------------|----------|--------------------------|--------------------------------------|
| `apiKey`    | `ONEUPTIME_API_KEY`      | Yes      |                          | API key for authentication           |
| `baseUrl`   | `ONEUPTIME_BASE_URL`     | No       | `https://oneuptime.com`  | Base URL of the OneUptime instance   |
| `projectId` | `ONEUPTIME_PROJECT_ID`   | No       |                          | Default project ID for all resources |

### Using Pulumi Config

```bash
pulumi config set oneuptime:apiKey --secret <your-api-key>
pulumi config set oneuptime:projectId <your-project-id>

# For self-hosted instances
pulumi config set oneuptime:baseUrl https://your-instance.example.com
```

### Using Environment Variables

```bash
export ONEUPTIME_API_KEY="your-api-key"
export ONEUPTIME_PROJECT_ID="your-project-id"
```

## Resources

| Resource             | Description                              |
|----------------------|------------------------------------------|
| `Team`               | Manages teams within a project           |
| `Monitor`            | Manages monitors (API, website, etc.)    |
| `StatusPage`         | Manages public/private status pages      |
| `Incident`           | Manages incidents                        |
| `OnCallDutyPolicy`   | Manages on-call duty escalation policies |

## Example

### YAML

```yaml
name: oneuptime-example
runtime: yaml
resources:
  engineering-team:
    type: oneuptime:index:Team
    properties:
      projectId: ${projectId}
      name: Engineering
      description: Engineering team managed by Pulumi

  api-monitor:
    type: oneuptime:index:Monitor
    properties:
      projectId: ${projectId}
      name: API Health Check
      monitorType: API
      currentMonitorStatusId: ${monitorStatusId}

  status-page:
    type: oneuptime:index:StatusPage
    properties:
      projectId: ${projectId}
      name: Public Status
      pageTitle: Service Status
      isPublicStatusPage: true
      showIncidentHistoryInDays: 14
```

### TypeScript

```typescript
import * as oneuptime from "@ottercoders/pulumi-oneuptime";

const team = new oneuptime.Team("engineering", {
    projectId: "your-project-id",
    name: "Engineering",
    description: "Engineering team",
});

const monitor = new oneuptime.Monitor("api-health", {
    projectId: "your-project-id",
    name: "API Health Check",
    monitorType: "API",
    currentMonitorStatusId: "your-status-id",
});
```

### Python

```python
import pulumi_oneuptime as oneuptime

team = oneuptime.Team("engineering",
    project_id="your-project-id",
    name="Engineering",
    description="Engineering team",
)

monitor = oneuptime.Monitor("api-health",
    project_id="your-project-id",
    name="API Health Check",
    monitor_type="API",
    current_monitor_status_id="your-status-id",
)
```

## Development

### Prerequisites

- [Go](https://golang.org/dl/) 1.22+
- [Pulumi CLI](https://www.pulumi.com/docs/install/)

### Build

```bash
make provider    # Build the provider binary
make install     # Install to $GOPATH/bin
make test        # Run unit tests
make schema      # Generate schema.json
make lint        # Run linter
make clean       # Remove build artifacts
```

### Testing

Unit tests use mocked HTTP servers and require no credentials:

```bash
make test
```

Acceptance tests run against a real OneUptime instance:

```bash
export ONEUPTIME_API_KEY="your-key"
export ONEUPTIME_PROJECT_ID="your-project-id"
make test-acceptance
```

### Adding a New Resource

1. Create a new file in `provider/resources/` (e.g., `workflow.go`)
2. Define `Args`, `State`, and resource structs with `pulumi:"..."` and `json:"..."` tags
3. Implement `Create`, `Read`, `Update`, `Delete` methods following the pattern in existing resources
4. Register the resource in `provider/provider.go`

## License

Apache 2.0 - See [LICENSE](LICENSE) for details.
