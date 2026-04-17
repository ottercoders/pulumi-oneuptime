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

All resources (except `Project`) accept an optional `projectId` property. When omitted, the provider-level `projectId` config is used. Every resource exposes `resourceId`, `slug`, `createdAt`, and `updatedAt` as outputs.

| Resource             | Description                              |
|----------------------|------------------------------------------|
| `Team`               | Manages teams within a project           |
| `Monitor`            | Manages monitors (API, website, etc.)    |
| `StatusPage`         | Manages public/private status pages      |
| `Incident`           | Manages incidents                        |
| `OnCallDutyPolicy`   | Manages on-call duty escalation policies |
| `Label`              | Manages labels for categorizing resources |
| `MonitorGroup`       | Manages logical groupings of monitors    |
| `IncidentState`      | Manages incident workflow states (e.g. Created, Acknowledged, Resolved) |
| `IncidentSeverity`   | Manages incident severity levels (e.g. Critical, Major, Minor) |
| `AlertState`         | Manages alert workflow states            |
| `AlertSeverity`      | Manages alert severity levels            |
| `Project`            | Manages OneUptime projects               |
| `StatusPageGroup`    | Groups monitors on a status page          |
| `StatusPageResource` | Links a monitor or monitor group to a status page |
| `OnCallDutyPolicyEscalationRule` | Escalation timing/order within an on-call policy |
| `OnCallDutyPolicySchedule` | Rotation schedule within an on-call policy |
| `OnCallDutyPolicyEscalationRuleSchedule` | Links a schedule to an escalation rule |
| `ProjectSSO`         | OIDC/SAML SSO configuration              |
| `Domain`             | Verified email domain for SSO user matching |
| `ProjectSsoTeam`     | Auto-assigns SSO users to a team on first login |
| `ApiKey`             | API key lifecycle (the `apiKey` output is a secret) |
| `ApiKeyPermission`   | Fine-grained permissions on an API key   |
| `TeamMember`         | Adds an existing user to a team          |
| `TeamPermission`     | Assigns permissions to a team            |
| `MonitorTeamOwner`   | Team owns a monitor                      |
| `MonitorUserOwner`   | User owns a monitor                      |
| `MonitorGroupResource` | Adds a monitor to a monitor group     |
| `Probe`              | Monitoring probe registration (the `key` output is a secret) |
| `StatusPageAnnouncement` | Announcement banner on status pages  |
| `StatusPageDomain`   | Binds a custom domain to a status page   |
| `StatusPageHeaderLink` | Link in the status page header         |
| `StatusPageFooterLink` | Link in the status page footer         |
| `StatusPageSubscriber` | Email/SMS/webhook subscriber           |
| `StatusPageTeamOwner` | Team owns a status page                 |
| `StatusPageUserOwner` | User owns a status page                 |
| `IncidentTemplate`   | Standardized template for incident creation |
| `IncidentNoteTemplate` | Reusable incident note template        |
| `ScheduledMaintenanceEvent` | Planned maintenance window         |
| `ScheduledMaintenanceState` | Workflow state for maintenance events |
| `ScheduledMaintenanceTemplate` | Template for maintenance events |
| `Workflow`           | Automation workflow                      |

### Bootstrap flow

The provider cannot create the very first OneUptime user or the master API key; both must be set up via the OneUptime UI:

1. Install OneUptime and register an admin user via the UI
2. Create a project in the UI
3. Mint the master API key in **Project Settings → API Keys**
4. Set `oneuptime:apiKey` in your Pulumi config (or export `ONEUPTIME_API_KEY`)
5. Everything else — additional API keys, team membership, permissions, monitors, etc. — is declarative from here on

### Resource Properties

#### Team

| Property      | Type   | Required | Description          |
|---------------|--------|----------|----------------------|
| `name`        | string | Yes      | Team name            |
| `description` | string | No       | Team description     |

#### Monitor

| Property                    | Type   | Required | Description                        |
|-----------------------------|--------|----------|------------------------------------|
| `name`                      | string | Yes      | Monitor name                       |
| `monitorType`               | string | Yes      | Type of monitor (e.g. API, Website)|
| `currentMonitorStatusId`    | string | Yes      | ID of the current monitor status   |
| `description`               | string | No       | Monitor description                |
| `disableActiveMonitoring`   | bool   | No       | Disable active monitoring          |
| `monitoringInterval`        | string | No       | Monitoring check interval          |

#### StatusPage

| Property                            | Type   | Required | Description                       |
|-------------------------------------|--------|----------|-----------------------------------|
| `name`                              | string | Yes      | Status page name                  |
| `description`                       | string | No       | Status page description           |
| `pageTitle`                         | string | No       | HTML page title                   |
| `pageDescription`                   | string | No       | Page description text             |
| `isPublicStatusPage`                | bool   | No       | Whether the page is public        |
| `showIncidentHistoryInDays`         | int    | No       | Days of incident history to show  |
| `showAnnouncementHistoryInDays`     | int    | No       | Days of announcement history      |
| `showScheduledEventHistoryInDays`   | int    | No       | Days of scheduled event history   |

#### Incident

| Property                  | Type   | Required | Description                     |
|---------------------------|--------|----------|---------------------------------|
| `title`                   | string | Yes      | Incident title                  |
| `currentIncidentStateId`  | string | Yes      | ID of the current state         |
| `incidentSeverityId`      | string | Yes      | ID of the severity level        |
| `description`             | string | No       | Incident description            |
| `declaredAt`              | string | No       | When the incident was declared  |

#### OnCallDutyPolicy

| Property                                      | Type   | Required | Description                              |
|-----------------------------------------------|--------|----------|------------------------------------------|
| `name`                                        | string | Yes      | Policy name                              |
| `description`                                 | string | No       | Policy description                       |
| `repeatPolicyIfNoOneAcknowledges`             | bool   | No       | Repeat if no acknowledgement             |
| `repeatPolicyIfNoOneAcknowledgesNoOfTimes`    | int    | No       | Number of times to repeat                |

#### Label

| Property      | Type   | Required | Description                          |
|---------------|--------|----------|--------------------------------------|
| `name`        | string | Yes      | Label name                           |
| `color`       | string | Yes      | Hex color code (e.g. `#e74c3c`)     |
| `description` | string | No       | Label description                    |

#### MonitorGroup

| Property      | Type   | Required | Description               |
|---------------|--------|----------|---------------------------|
| `name`        | string | Yes      | Group name                |
| `description` | string | No       | Group description         |

#### IncidentState

| Property              | Type   | Required | Description                        |
|-----------------------|--------|----------|------------------------------------|
| `name`                | string | Yes      | State name                         |
| `color`               | string | Yes      | Hex color code                     |
| `description`         | string | No       | State description                  |
| `isCreatedState`      | bool   | No       | Mark as the initial created state  |
| `isAcknowledgedState` | bool   | No       | Mark as the acknowledged state     |
| `isResolvedState`     | bool   | No       | Mark as the resolved state         |
| `order`               | int    | No       | Display order                      |

#### IncidentSeverity

| Property      | Type   | Required | Description                    |
|---------------|--------|----------|--------------------------------|
| `name`        | string | Yes      | Severity name                  |
| `color`       | string | Yes      | Hex color code                 |
| `description` | string | No       | Severity description           |
| `order`       | int    | No       | Display order (1 = most severe)|

#### AlertState

| Property              | Type   | Required | Description                        |
|-----------------------|--------|----------|------------------------------------|
| `name`                | string | Yes      | State name                         |
| `color`               | string | Yes      | Hex color code                     |
| `description`         | string | No       | State description                  |
| `isCreatedState`      | bool   | No       | Mark as the initial created state  |
| `isAcknowledgedState` | bool   | No       | Mark as the acknowledged state     |
| `isResolvedState`     | bool   | No       | Mark as the resolved state         |
| `order`               | int    | No       | Display order                      |

#### AlertSeverity

| Property      | Type   | Required | Description                    |
|---------------|--------|----------|--------------------------------|
| `name`        | string | Yes      | Severity name                  |
| `color`       | string | Yes      | Hex color code                 |
| `description` | string | No       | Severity description           |
| `order`       | int    | No       | Display order (1 = most severe)|

#### Project

| Property      | Type   | Required | Description            |
|---------------|--------|----------|------------------------|
| `name`        | string | Yes      | Project name           |
| `description` | string | No       | Project description    |

> **Note:** The `Project` resource does not accept a `projectId` property since it represents the project itself.

## Example

### YAML

```yaml
name: oneuptime-example
runtime: yaml
resources:
  engineering-team:
    type: oneuptime:resources:Team
    properties:
      name: Engineering
      description: Engineering team managed by Pulumi

  api-monitor:
    type: oneuptime:resources:Monitor
    properties:
      name: API Health Check
      monitorType: API
      currentMonitorStatusId: ${monitorStatusId}

  status-page:
    type: oneuptime:resources:StatusPage
    properties:
      name: Public Status
      pageTitle: Service Status
      isPublicStatusPage: true
      showIncidentHistoryInDays: 14

  env-label:
    type: oneuptime:resources:Label
    properties:
      name: production
      color: "#e74c3c"
      description: Production environment

  api-monitors:
    type: oneuptime:resources:MonitorGroup
    properties:
      name: API Monitors
      description: All API health check monitors

  critical-severity:
    type: oneuptime:resources:IncidentSeverity
    properties:
      name: Critical
      color: "#e74c3c"
      description: Service is completely down
      order: 1
```

### TypeScript

```typescript
import * as oneuptime from "@ottercoders/pulumi-oneuptime";

const team = new oneuptime.resources.Team("engineering", {
    name: "Engineering",
    description: "Engineering team",
});

const label = new oneuptime.resources.Label("production", {
    name: "production",
    color: "#e74c3c",
    description: "Production environment",
});

const severity = new oneuptime.resources.IncidentSeverity("critical", {
    name: "Critical",
    color: "#e74c3c",
    description: "Service is completely down",
    order: 1,
});
```

### Python

```python
from pulumi_oneuptime import resources

team = resources.Team("engineering",
    name="Engineering",
    description="Engineering team",
)

label = resources.Label("production",
    name="production",
    color="#e74c3c",
    description="Production environment",
)

severity = resources.IncidentSeverity("critical",
    name="Critical",
    color="#e74c3c",
    description="Service is completely down",
    order=1,
)
```

## Development

### Prerequisites

- [Go](https://golang.org/dl/) 1.25+
- [Pulumi CLI](https://www.pulumi.com/docs/install/)

### Build

```bash
make provider    # Build the provider binary
make install     # Install to $GOPATH/bin
make schema      # Generate schema.json
make sdk         # Generate SDKs (Go, Node.js, Python, .NET)
make test        # Run unit tests
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
