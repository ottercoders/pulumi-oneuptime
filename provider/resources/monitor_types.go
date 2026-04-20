package resources

// Shared nested types used by the Monitor resource and its sub-configurations.
// Field shapes mirror the OneUptime master branch TypeScript interfaces as of
// 2026-04-20. Values that OneUptime wraps in _type envelopes on the wire
// (MonitorSteps, MonitorStep, MonitorCriteria, MonitorCriteriaInstance, plus
// URL/IP/Hostname/Port) are enveloped by the helpers in
// monitor_steps_envelope.go rather than being expressed as typed envelopes on
// the Go structs — keeping the Go/SDK surface flat for DX while still producing
// the shape the API demands.

// MonitorStep is one step of a monitor probe run. Most fields apply only to a
// subset of monitor types; populate the sub-config (logMonitor, snmpMonitor,
// ...) that matches the parent Monitor's monitorType.
type MonitorStep struct {
	ID                        string                               `pulumi:"id,optional" json:"id,omitempty"`
	MonitorDestination        *string                              `pulumi:"monitorDestination,optional" json:"-"`
	MonitorDestinationPort    *int                                 `pulumi:"monitorDestinationPort,optional" json:"-"`
	MonitorCriteria           MonitorCriteria                      `pulumi:"monitorCriteria" json:"-"`
	RequestType               *string                              `pulumi:"requestType,optional" json:"requestType,omitempty"`
	RequestHeaders            map[string]string                    `pulumi:"requestHeaders,optional" json:"requestHeaders,omitempty"`
	RequestBody               *string                              `pulumi:"requestBody,optional" json:"requestBody,omitempty"`
	DoNotFollowRedirects      *bool                                `pulumi:"doNotFollowRedirects,optional" json:"doNotFollowRedirects,omitempty"`
	CustomCode                *string                              `pulumi:"customCode,optional" json:"customCode,omitempty"`
	ScreenSizeTypes           []string                             `pulumi:"screenSizeTypes,optional" json:"screenSizeTypes,omitempty"`
	BrowserTypes              []string                             `pulumi:"browserTypes,optional" json:"browserTypes,omitempty"`
	RetryCountOnError         *int                                 `pulumi:"retryCountOnError,optional" json:"retryCountOnError,omitempty"`
	LogMonitor                *LogMonitorConfig                    `pulumi:"logMonitor,optional" json:"logMonitor,omitempty"`
	TraceMonitor              *TraceMonitorConfig                  `pulumi:"traceMonitor,optional" json:"traceMonitor,omitempty"`
	MetricMonitor             *MetricMonitorConfig                 `pulumi:"metricMonitor,optional" json:"metricMonitor,omitempty"`
	ExceptionMonitor          *ExceptionMonitorConfig              `pulumi:"exceptionMonitor,optional" json:"exceptionMonitor,omitempty"`
	ProfileMonitor            *ProfileMonitorConfig                `pulumi:"profileMonitor,optional" json:"profileMonitor,omitempty"`
	SnmpMonitor               *SnmpMonitorConfig                   `pulumi:"snmpMonitor,optional" json:"snmpMonitor,omitempty"`
	DnsMonitor                *DnsMonitorConfig                    `pulumi:"dnsMonitor,optional" json:"dnsMonitor,omitempty"`
	DomainMonitor             *DomainMonitorConfig                 `pulumi:"domainMonitor,optional" json:"domainMonitor,omitempty"`
	ExternalStatusPageMonitor *ExternalStatusPageMonitorConfig     `pulumi:"externalStatusPageMonitor,optional" json:"externalStatusPageMonitor,omitempty"`
	KubernetesMonitor         *KubernetesMonitorConfig             `pulumi:"kubernetesMonitor,optional" json:"kubernetesMonitor,omitempty"`
	DockerMonitor             *DockerMonitorConfig                 `pulumi:"dockerMonitor,optional" json:"dockerMonitor,omitempty"`
}

// MonitorCriteria is the decision block attached to a step — an array of
// criteria instances, any (or all) of which may trigger a status change,
// incident, or alert.
type MonitorCriteria struct {
	CriteriaInstances []MonitorCriteriaInstance `pulumi:"criteriaInstances" json:"-"`
}

// MonitorCriteriaInstance is one decision rule. When its filters (combined by
// filterCondition) match, its incidents/alerts are created and/or the monitor
// status is flipped to monitorStatusId.
type MonitorCriteriaInstance struct {
	ID                  string           `pulumi:"id,optional" json:"id,omitempty"`
	Name                string           `pulumi:"name" json:"name"`
	Description         string           `pulumi:"description" json:"description"`
	FilterCondition     string           `pulumi:"filterCondition" json:"filterCondition"`
	Filters             []CriteriaFilter `pulumi:"filters" json:"filters"`
	Incidents           []IncidentRef    `pulumi:"incidents,optional" json:"incidents,omitempty"`
	Alerts              []AlertRef       `pulumi:"alerts,optional" json:"alerts,omitempty"`
	MonitorStatusID     *string          `pulumi:"monitorStatusId,optional" json:"monitorStatusId,omitempty"`
	ChangeMonitorStatus *bool            `pulumi:"changeMonitorStatus,optional" json:"changeMonitorStatus,omitempty"`
	CreateIncidents     *bool            `pulumi:"createIncidents,optional" json:"createIncidents,omitempty"`
	CreateAlerts        *bool            `pulumi:"createAlerts,optional" json:"createAlerts,omitempty"`
}

// CriteriaFilter expresses one condition such as "Response Status Code >= 400".
// CheckOn and FilterType use OneUptime's display strings verbatim (e.g.
// "Response Status Code", "Greater Than"); see the provider README for the
// full enum list.
type CriteriaFilter struct {
	CheckOn                 string                   `pulumi:"checkOn" json:"checkOn"`
	FilterType              *string                  `pulumi:"filterType,optional" json:"filterType,omitempty"`
	Value                   *string                  `pulumi:"value,optional" json:"value,omitempty"`
	EvaluateOverTime        *bool                    `pulumi:"evaluateOverTime,optional" json:"evaluateOverTime,omitempty"`
	EvaluateOverTimeOptions *EvaluateOverTimeOptions `pulumi:"evaluateOverTimeOptions,optional" json:"evaluateOverTimeOptions,omitempty"`
	ServerMonitorOptions    *ServerMonitorOptions    `pulumi:"serverMonitorOptions,optional" json:"serverMonitorOptions,omitempty"`
	MetricMonitorOptions    *MetricMonitorOptions    `pulumi:"metricMonitorOptions,optional" json:"metricMonitorOptions,omitempty"`
	SnmpMonitorOptions      *SnmpMonitorOptions      `pulumi:"snmpMonitorOptions,optional" json:"snmpMonitorOptions,omitempty"`
}

type EvaluateOverTimeOptions struct {
	TimeValueInMinutes   *int    `pulumi:"timeValueInMinutes,optional" json:"timeValueInMinutes,omitempty"`
	EvaluateOverTimeType *string `pulumi:"evaluateOverTimeType,optional" json:"evaluateOverTimeType,omitempty"`
}

type ServerMonitorOptions struct {
	DiskPath *string `pulumi:"diskPath,optional" json:"diskPath,omitempty"`
}

type MetricMonitorOptions struct {
	MetricAlias           *string `pulumi:"metricAlias,optional" json:"metricAlias,omitempty"`
	MetricAggregationType *string `pulumi:"metricAggregationType,optional" json:"metricAggregationType,omitempty"`
	OnNoDataPolicy        *string `pulumi:"onNoDataPolicy,optional" json:"onNoDataPolicy,omitempty"`
	ThresholdUnit         *string `pulumi:"thresholdUnit,optional" json:"thresholdUnit,omitempty"`
}

type SnmpMonitorOptions struct {
	Oid *string `pulumi:"oid,optional" json:"oid,omitempty"`
}

// IncidentRef is the body used to auto-open an incident when a criteria
// instance fires. ID is a user-supplied client id (not the server ObjectID),
// commonly a UUID generated at config time.
type IncidentRef struct {
	ID                        string                          `pulumi:"id" json:"id"`
	Title                     string                          `pulumi:"title" json:"title"`
	Description               string                          `pulumi:"description" json:"description"`
	IncidentSeverityID        *string                         `pulumi:"incidentSeverityId,optional" json:"incidentSeverityId,omitempty"`
	AutoResolveIncident       *bool                           `pulumi:"autoResolveIncident,optional" json:"autoResolveIncident,omitempty"`
	RemediationNotes          *string                         `pulumi:"remediationNotes,optional" json:"remediationNotes,omitempty"`
	OnCallPolicyIDs           []string                        `pulumi:"onCallPolicyIds,optional" json:"onCallPolicyIds,omitempty"`
	LabelIDs                  []string                        `pulumi:"labelIds,optional" json:"labelIds,omitempty"`
	OwnerTeamIDs              []string                        `pulumi:"ownerTeamIds,optional" json:"ownerTeamIds,omitempty"`
	OwnerUserIDs              []string                        `pulumi:"ownerUserIds,optional" json:"ownerUserIds,omitempty"`
	IncidentMemberRoles       []IncidentMemberRoleAssignment  `pulumi:"incidentMemberRoles,optional" json:"incidentMemberRoles,omitempty"`
	ShowIncidentOnStatusPage  *bool                           `pulumi:"showIncidentOnStatusPage,optional" json:"showIncidentOnStatusPage,omitempty"`
}

type IncidentMemberRoleAssignment struct {
	RoleID string `pulumi:"roleId" json:"roleId"`
	UserID string `pulumi:"userId" json:"userId"`
}

// AlertRef is the body used to auto-open an alert. Shape is a subset of
// IncidentRef — alerts lack member-role assignment and status-page flags.
type AlertRef struct {
	ID               string   `pulumi:"id" json:"id"`
	Title            string   `pulumi:"title" json:"title"`
	Description      string   `pulumi:"description" json:"description"`
	AlertSeverityID  *string  `pulumi:"alertSeverityId,optional" json:"alertSeverityId,omitempty"`
	AutoResolveAlert *bool    `pulumi:"autoResolveAlert,optional" json:"autoResolveAlert,omitempty"`
	RemediationNotes *string  `pulumi:"remediationNotes,optional" json:"remediationNotes,omitempty"`
	OnCallPolicyIDs  []string `pulumi:"onCallPolicyIds,optional" json:"onCallPolicyIds,omitempty"`
	LabelIDs         []string `pulumi:"labelIds,optional" json:"labelIds,omitempty"`
	OwnerTeamIDs     []string `pulumi:"ownerTeamIds,optional" json:"ownerTeamIds,omitempty"`
	OwnerUserIDs     []string `pulumi:"ownerUserIds,optional" json:"ownerUserIds,omitempty"`
}

// WorkspaceChannelRef targets a Slack/Teams channel for monitor update
// notifications. Used by Monitor.postUpdatesToWorkspaceChannels.
type WorkspaceChannelRef struct {
	WorkspaceType string  `pulumi:"workspaceType" json:"workspaceType"`
	ChannelName   string  `pulumi:"channelName" json:"channelName"`
	ChannelID     *string `pulumi:"channelId,optional" json:"channelId,omitempty"`
}

// ── Sub-monitor configs (one-per-step, matched to Monitor.monitorType) ──

type LogMonitorConfig struct {
	Attributes          map[string]interface{} `pulumi:"attributes" json:"attributes"`
	Body                string                 `pulumi:"body" json:"body"`
	SeverityTexts       []string               `pulumi:"severityTexts" json:"severityTexts"`
	TelemetryServiceIDs []string               `pulumi:"telemetryServiceIds" json:"telemetryServiceIds"`
	LastXSecondsOfLogs  int                    `pulumi:"lastXSecondsOfLogs" json:"lastXSecondsOfLogs"`
}

type TraceMonitorConfig struct {
	Attributes          map[string]interface{} `pulumi:"attributes" json:"attributes"`
	SpanStatuses        []string               `pulumi:"spanStatuses" json:"spanStatuses"`
	TelemetryServiceIDs []string               `pulumi:"telemetryServiceIds" json:"telemetryServiceIds"`
	LastXSecondsOfSpans int                    `pulumi:"lastXSecondsOfSpans" json:"lastXSecondsOfSpans"`
	SpanName            string                 `pulumi:"spanName" json:"spanName"`
}

// MetricMonitorConfig's metricViewConfig is a deeply nested OneUptime-internal
// shape (queryConfigs/formulaConfigs) we intentionally don't model — pass the
// raw map matching whatever the OneUptime UI builds.
type MetricMonitorConfig struct {
	MetricViewConfig map[string]interface{} `pulumi:"metricViewConfig" json:"metricViewConfig"`
	RollingTime      string                 `pulumi:"rollingTime" json:"rollingTime"`
}

type ExceptionMonitorConfig struct {
	TelemetryServiceIDs      []string `pulumi:"telemetryServiceIds" json:"telemetryServiceIds"`
	ExceptionTypes           []string `pulumi:"exceptionTypes" json:"exceptionTypes"`
	Message                  string   `pulumi:"message" json:"message"`
	IncludeResolved          bool     `pulumi:"includeResolved" json:"includeResolved"`
	IncludeArchived          bool     `pulumi:"includeArchived" json:"includeArchived"`
	LastXSecondsOfExceptions int      `pulumi:"lastXSecondsOfExceptions" json:"lastXSecondsOfExceptions"`
}

type ProfileMonitorConfig struct {
	Attributes             map[string]interface{} `pulumi:"attributes" json:"attributes"`
	ProfileTypes           []string               `pulumi:"profileTypes" json:"profileTypes"`
	TelemetryServiceIDs    []string               `pulumi:"telemetryServiceIds" json:"telemetryServiceIds"`
	LastXSecondsOfProfiles int                    `pulumi:"lastXSecondsOfProfiles" json:"lastXSecondsOfProfiles"`
	ProfileType            string                 `pulumi:"profileType" json:"profileType"`
}

type SnmpMonitorConfig struct {
	SnmpVersion     string       `pulumi:"snmpVersion" json:"snmpVersion"`
	Hostname        string       `pulumi:"hostname" json:"hostname"`
	Port            int          `pulumi:"port" json:"port"`
	CommunityString *string      `pulumi:"communityString,optional" json:"communityString,omitempty"`
	SnmpV3Auth      *SnmpV3Auth  `pulumi:"snmpV3Auth,optional" json:"snmpV3Auth,omitempty"`
	Oids            []SnmpOid    `pulumi:"oids" json:"oids"`
	Timeout         int          `pulumi:"timeout" json:"timeout"`
	Retries         int          `pulumi:"retries" json:"retries"`
}

type SnmpOid struct {
	Oid         string  `pulumi:"oid" json:"oid"`
	Name        *string `pulumi:"name,optional" json:"name,omitempty"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
}

type SnmpV3Auth struct {
	SecurityLevel string  `pulumi:"securityLevel" json:"securityLevel"`
	Username      string  `pulumi:"username" json:"username"`
	AuthProtocol  *string `pulumi:"authProtocol,optional" json:"authProtocol,omitempty"`
	AuthKey       *string `pulumi:"authKey,optional" json:"authKey,omitempty"`
	PrivProtocol  *string `pulumi:"privProtocol,optional" json:"privProtocol,omitempty"`
	PrivKey       *string `pulumi:"privKey,optional" json:"privKey,omitempty"`
}

type DnsMonitorConfig struct {
	QueryName  string  `pulumi:"queryName" json:"queryName"`
	RecordType string  `pulumi:"recordType" json:"recordType"`
	Hostname   *string `pulumi:"hostname,optional" json:"hostname,omitempty"`
	Port       int     `pulumi:"port" json:"port"`
	Timeout    int     `pulumi:"timeout" json:"timeout"`
	Retries    int     `pulumi:"retries" json:"retries"`
}

type DomainMonitorConfig struct {
	DomainName string `pulumi:"domainName" json:"domainName"`
	Timeout    int    `pulumi:"timeout" json:"timeout"`
	Retries    int    `pulumi:"retries" json:"retries"`
}

type ExternalStatusPageMonitorConfig struct {
	StatusPageURL string  `pulumi:"statusPageUrl" json:"statusPageUrl"`
	Provider      string  `pulumi:"provider" json:"provider"`
	ComponentName *string `pulumi:"componentName,optional" json:"componentName,omitempty"`
	Timeout       int     `pulumi:"timeout" json:"timeout"`
	Retries       int     `pulumi:"retries" json:"retries"`
}

type KubernetesMonitorConfig struct {
	ClusterIdentifier string                      `pulumi:"clusterIdentifier" json:"clusterIdentifier"`
	ResourceScope     string                      `pulumi:"resourceScope" json:"resourceScope"`
	ResourceFilters   KubernetesResourceFilters   `pulumi:"resourceFilters" json:"resourceFilters"`
	MetricViewConfig  map[string]interface{}      `pulumi:"metricViewConfig" json:"metricViewConfig"`
	RollingTime       string                      `pulumi:"rollingTime" json:"rollingTime"`
}

type KubernetesResourceFilters struct {
	Namespace    *string `pulumi:"namespace,optional" json:"namespace,omitempty"`
	WorkloadType *string `pulumi:"workloadType,optional" json:"workloadType,omitempty"`
	WorkloadName *string `pulumi:"workloadName,optional" json:"workloadName,omitempty"`
	NodeName     *string `pulumi:"nodeName,optional" json:"nodeName,omitempty"`
	PodName      *string `pulumi:"podName,optional" json:"podName,omitempty"`
}

type DockerMonitorConfig struct {
	HostIdentifier   string                   `pulumi:"hostIdentifier" json:"hostIdentifier"`
	ContainerFilters DockerContainerFilters   `pulumi:"containerFilters" json:"containerFilters"`
	MetricViewConfig map[string]interface{}   `pulumi:"metricViewConfig" json:"metricViewConfig"`
	RollingTime      string                   `pulumi:"rollingTime" json:"rollingTime"`
}

type DockerContainerFilters struct {
	ContainerName  *string `pulumi:"containerName,optional" json:"containerName,omitempty"`
	ContainerImage *string `pulumi:"containerImage,optional" json:"containerImage,omitempty"`
	HostName       *string `pulumi:"hostName,optional" json:"hostName,omitempty"`
}

// ── Output-only types on MonitorState ──

type IncomingMonitorRequest struct {
	ProjectID                             string                 `pulumi:"projectId,optional" json:"projectId,omitempty"`
	MonitorID                             string                 `pulumi:"monitorId,optional" json:"monitorId,omitempty"`
	RequestHeaders                        map[string]string      `pulumi:"requestHeaders,optional" json:"requestHeaders,omitempty"`
	RequestBody                           map[string]interface{} `pulumi:"requestBody,optional" json:"requestBody,omitempty"`
	RequestMethod                         *string                `pulumi:"requestMethod,optional" json:"requestMethod,omitempty"`
	IncomingRequestReceivedAt             *string                `pulumi:"incomingRequestReceivedAt,optional" json:"incomingRequestReceivedAt,omitempty"`
	OnlyCheckForIncomingRequestReceivedAt *bool                  `pulumi:"onlyCheckForIncomingRequestReceivedAt,optional" json:"onlyCheckForIncomingRequestReceivedAt,omitempty"`
	CheckedAt                             *string                `pulumi:"checkedAt,optional" json:"checkedAt,omitempty"`
}

type IncomingEmailMonitorRequest struct {
	ProjectID                           string                 `pulumi:"projectId,optional" json:"projectId,omitempty"`
	MonitorID                           string                 `pulumi:"monitorId,optional" json:"monitorId,omitempty"`
	EmailFrom                           *string                `pulumi:"emailFrom,optional" json:"emailFrom,omitempty"`
	EmailTo                             *string                `pulumi:"emailTo,optional" json:"emailTo,omitempty"`
	EmailSubject                        *string                `pulumi:"emailSubject,optional" json:"emailSubject,omitempty"`
	EmailBody                           *string                `pulumi:"emailBody,optional" json:"emailBody,omitempty"`
	EmailBodyHtml                       *string                `pulumi:"emailBodyHtml,optional" json:"emailBodyHtml,omitempty"`
	EmailHeaders                        map[string]string      `pulumi:"emailHeaders,optional" json:"emailHeaders,omitempty"`
	EmailReceivedAt                     *string                `pulumi:"emailReceivedAt,optional" json:"emailReceivedAt,omitempty"`
	CheckedAt                           *string                `pulumi:"checkedAt,optional" json:"checkedAt,omitempty"`
	Attachments                         []EmailAttachment      `pulumi:"attachments,optional" json:"attachments,omitempty"`
	OnlyCheckForIncomingEmailReceivedAt *bool                  `pulumi:"onlyCheckForIncomingEmailReceivedAt,optional" json:"onlyCheckForIncomingEmailReceivedAt,omitempty"`
}

type EmailAttachment struct {
	Filename    string `pulumi:"filename,optional" json:"filename,omitempty"`
	ContentType string `pulumi:"contentType,optional" json:"contentType,omitempty"`
	Size        int    `pulumi:"size,optional" json:"size,omitempty"`
}

type ServerMonitorResponse struct {
	ProjectID                  string                 `pulumi:"projectId,optional" json:"projectId,omitempty"`
	MonitorID                  string                 `pulumi:"monitorId,optional" json:"monitorId,omitempty"`
	Hostname                   *string                `pulumi:"hostname,optional" json:"hostname,omitempty"`
	BasicInfrastructureMetrics map[string]interface{} `pulumi:"basicInfrastructureMetrics,optional" json:"basicInfrastructureMetrics,omitempty"`
	RequestReceivedAt          *string                `pulumi:"requestReceivedAt,optional" json:"requestReceivedAt,omitempty"`
	OnlyCheckRequestReceivedAt *bool                  `pulumi:"onlyCheckRequestReceivedAt,optional" json:"onlyCheckRequestReceivedAt,omitempty"`
	Processes                  []ServerProcess        `pulumi:"processes,optional" json:"processes,omitempty"`
	FailureCause               *string                `pulumi:"failureCause,optional" json:"failureCause,omitempty"`
	TimeNow                    *string                `pulumi:"timeNow,optional" json:"timeNow,omitempty"`
}

type ServerProcess struct {
	PID     int    `pulumi:"pid,optional" json:"pid,omitempty"`
	Name    string `pulumi:"name,optional" json:"name,omitempty"`
	Command string `pulumi:"command,optional" json:"command,omitempty"`
}
