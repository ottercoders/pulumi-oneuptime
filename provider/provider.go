package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/ottercoders/pulumi-oneuptime/provider/resources"
)

const Name = "oneuptime"

var Version = "0.0.1-dev"

func Provider() p.Provider {
	prov, err := infer.NewProviderBuilder().
		WithDisplayName("OneUptime").
		WithDescription("Manage OneUptime monitoring resources").
		WithPluginDownloadURL("github://api.github.com/ottercoders/pulumi-oneuptime").
		WithConfig(infer.Config[*resources.Config](&resources.Config{})).
		WithResources(
			infer.Resource[*resources.Team, resources.TeamArgs, resources.TeamState](&resources.Team{}),
			infer.Resource[*resources.Monitor, resources.MonitorArgs, resources.MonitorState](&resources.Monitor{}),
			infer.Resource[*resources.StatusPage, resources.StatusPageArgs, resources.StatusPageState](&resources.StatusPage{}),
			infer.Resource[*resources.Incident, resources.IncidentArgs, resources.IncidentState](&resources.Incident{}),
			infer.Resource[*resources.OnCallDutyPolicy, resources.OnCallDutyPolicyArgs, resources.OnCallDutyPolicyState](&resources.OnCallDutyPolicy{}),
			infer.Resource[*resources.Label, resources.LabelArgs, resources.LabelState](&resources.Label{}),
			infer.Resource[*resources.MonitorGroup, resources.MonitorGroupArgs, resources.MonitorGroupState](&resources.MonitorGroup{}),
			infer.Resource[*resources.IncidentStateResource, resources.IncidentStateResourceArgs, resources.IncidentStateResourceState](&resources.IncidentStateResource{}),
			infer.Resource[*resources.IncidentSeverity, resources.IncidentSeverityArgs, resources.IncidentSeverityState](&resources.IncidentSeverity{}),
			infer.Resource[*resources.AlertState, resources.AlertStateArgs, resources.AlertStateState](&resources.AlertState{}),
			infer.Resource[*resources.AlertSeverity, resources.AlertSeverityArgs, resources.AlertSeverityState](&resources.AlertSeverity{}),
			infer.Resource[*resources.Project, resources.ProjectArgs, resources.ProjectState](&resources.Project{}),
			infer.Resource[*resources.StatusPageGroup, resources.StatusPageGroupArgs, resources.StatusPageGroupState](&resources.StatusPageGroup{}),
			infer.Resource[*resources.StatusPageResource, resources.StatusPageResourceArgs, resources.StatusPageResourceState](&resources.StatusPageResource{}),
			infer.Resource[*resources.OnCallDutyPolicyEscalationRule, resources.OnCallDutyPolicyEscalationRuleArgs, resources.OnCallDutyPolicyEscalationRuleState](&resources.OnCallDutyPolicyEscalationRule{}),
			infer.Resource[*resources.OnCallDutyPolicySchedule, resources.OnCallDutyPolicyScheduleArgs, resources.OnCallDutyPolicyScheduleState](&resources.OnCallDutyPolicySchedule{}),
			infer.Resource[*resources.OnCallDutyPolicyEscalationRuleSchedule, resources.OnCallDutyPolicyEscalationRuleScheduleArgs, resources.OnCallDutyPolicyEscalationRuleScheduleState](&resources.OnCallDutyPolicyEscalationRuleSchedule{}),
			infer.Resource[*resources.ProjectSSO, resources.ProjectSSOArgs, resources.ProjectSSOState](&resources.ProjectSSO{}),
			infer.Resource[*resources.Domain, resources.DomainArgs, resources.DomainState](&resources.Domain{}),
			infer.Resource[*resources.ProjectSsoTeam, resources.ProjectSsoTeamArgs, resources.ProjectSsoTeamState](&resources.ProjectSsoTeam{}),
			// Phase 1: Auth + ownership
			infer.Resource[*resources.ApiKey, resources.ApiKeyArgs, resources.ApiKeyState](&resources.ApiKey{}),
			infer.Resource[*resources.ApiKeyPermission, resources.ApiKeyPermissionArgs, resources.ApiKeyPermissionState](&resources.ApiKeyPermission{}),
			infer.Resource[*resources.TeamMember, resources.TeamMemberArgs, resources.TeamMemberState](&resources.TeamMember{}),
			infer.Resource[*resources.TeamPermission, resources.TeamPermissionArgs, resources.TeamPermissionState](&resources.TeamPermission{}),
			infer.Resource[*resources.MonitorTeamOwner, resources.MonitorTeamOwnerArgs, resources.MonitorTeamOwnerState](&resources.MonitorTeamOwner{}),
			infer.Resource[*resources.MonitorUserOwner, resources.MonitorUserOwnerArgs, resources.MonitorUserOwnerState](&resources.MonitorUserOwner{}),
			infer.Resource[*resources.MonitorGroupResource, resources.MonitorGroupResourceArgs, resources.MonitorGroupResourceState](&resources.MonitorGroupResource{}),
			infer.Resource[*resources.Probe, resources.ProbeArgs, resources.ProbeState](&resources.Probe{}),
			// Phase 2: Status page completeness
			infer.Resource[*resources.StatusPageAnnouncement, resources.StatusPageAnnouncementArgs, resources.StatusPageAnnouncementState](&resources.StatusPageAnnouncement{}),
			infer.Resource[*resources.StatusPageDomain, resources.StatusPageDomainArgs, resources.StatusPageDomainState](&resources.StatusPageDomain{}),
			infer.Resource[*resources.StatusPageHeaderLink, resources.StatusPageHeaderLinkArgs, resources.StatusPageHeaderLinkState](&resources.StatusPageHeaderLink{}),
			infer.Resource[*resources.StatusPageFooterLink, resources.StatusPageFooterLinkArgs, resources.StatusPageFooterLinkState](&resources.StatusPageFooterLink{}),
			infer.Resource[*resources.StatusPageSubscriber, resources.StatusPageSubscriberArgs, resources.StatusPageSubscriberState](&resources.StatusPageSubscriber{}),
			infer.Resource[*resources.StatusPageTeamOwner, resources.StatusPageTeamOwnerArgs, resources.StatusPageTeamOwnerState](&resources.StatusPageTeamOwner{}),
			infer.Resource[*resources.StatusPageUserOwner, resources.StatusPageUserOwnerArgs, resources.StatusPageUserOwnerState](&resources.StatusPageUserOwner{}),
			// Phase 3: Incident + maintenance lifecycle
			infer.Resource[*resources.IncidentTemplate, resources.IncidentTemplateArgs, resources.IncidentTemplateState](&resources.IncidentTemplate{}),
			infer.Resource[*resources.IncidentNoteTemplate, resources.IncidentNoteTemplateArgs, resources.IncidentNoteTemplateState](&resources.IncidentNoteTemplate{}),
			infer.Resource[*resources.ScheduledMaintenanceEvent, resources.ScheduledMaintenanceEventArgs, resources.ScheduledMaintenanceEventState](&resources.ScheduledMaintenanceEvent{}),
			infer.Resource[*resources.ScheduledMaintenanceState, resources.ScheduledMaintenanceStateArgs, resources.ScheduledMaintenanceStateState](&resources.ScheduledMaintenanceState{}),
			infer.Resource[*resources.ScheduledMaintenanceTemplate, resources.ScheduledMaintenanceTemplateArgs, resources.ScheduledMaintenanceTemplateState](&resources.ScheduledMaintenanceTemplate{}),
			infer.Resource[*resources.Workflow, resources.WorkflowArgs, resources.WorkflowState](&resources.Workflow{}),
			// TF parity: on-call rotation, ownership, monitor status, workflow variables, telemetry
			infer.Resource[*resources.OnCallScheduleLayer, resources.OnCallScheduleLayerArgs, resources.OnCallScheduleLayerState](&resources.OnCallScheduleLayer{}),
			infer.Resource[*resources.OnCallScheduleLayerUser, resources.OnCallScheduleLayerUserArgs, resources.OnCallScheduleLayerUserState](&resources.OnCallScheduleLayerUser{}),
			infer.Resource[*resources.OnCallDutyPolicyTeamOwner, resources.OnCallDutyPolicyTeamOwnerArgs, resources.OnCallDutyPolicyTeamOwnerState](&resources.OnCallDutyPolicyTeamOwner{}),
			infer.Resource[*resources.OnCallDutyPolicyUserOwner, resources.OnCallDutyPolicyUserOwnerArgs, resources.OnCallDutyPolicyUserOwnerState](&resources.OnCallDutyPolicyUserOwner{}),
			infer.Resource[*resources.OnCallDutyPolicyEscalationRuleUser, resources.OnCallDutyPolicyEscalationRuleUserArgs, resources.OnCallDutyPolicyEscalationRuleUserState](&resources.OnCallDutyPolicyEscalationRuleUser{}),
			infer.Resource[*resources.OnCallDutyPolicyEscalationRuleTeam, resources.OnCallDutyPolicyEscalationRuleTeamArgs, resources.OnCallDutyPolicyEscalationRuleTeamState](&resources.OnCallDutyPolicyEscalationRuleTeam{}),
			infer.Resource[*resources.MonitorStatus, resources.MonitorStatusArgs, resources.MonitorStatusState](&resources.MonitorStatus{}),
			infer.Resource[*resources.MonitorProbe, resources.MonitorProbeArgs, resources.MonitorProbeState](&resources.MonitorProbe{}),
			infer.Resource[*resources.MonitorSecret, resources.MonitorSecretArgs, resources.MonitorSecretState](&resources.MonitorSecret{}),
			infer.Resource[*resources.MonitorCustomField, resources.MonitorCustomFieldArgs, resources.MonitorCustomFieldState](&resources.MonitorCustomField{}),
			infer.Resource[*resources.IncidentTeamOwner, resources.IncidentTeamOwnerArgs, resources.IncidentTeamOwnerState](&resources.IncidentTeamOwner{}),
			infer.Resource[*resources.IncidentUserOwner, resources.IncidentUserOwnerArgs, resources.IncidentUserOwnerState](&resources.IncidentUserOwner{}),
			infer.Resource[*resources.WorkflowVariable, resources.WorkflowVariableArgs, resources.WorkflowVariableState](&resources.WorkflowVariable{}),
			infer.Resource[*resources.TelemetryIngestionKey, resources.TelemetryIngestionKeyArgs, resources.TelemetryIngestionKeyState](&resources.TelemetryIngestionKey{}),
			infer.Resource[*resources.ScheduledMaintenanceTeamOwner, resources.ScheduledMaintenanceTeamOwnerArgs, resources.ScheduledMaintenanceTeamOwnerState](&resources.ScheduledMaintenanceTeamOwner{}),
			// Project-level outbound SMTP config
			infer.Resource[*resources.ProjectSmtpConfig, resources.ProjectSmtpConfigArgs, resources.ProjectSmtpConfigState](&resources.ProjectSmtpConfig{}),
		).
		WithFunctions(
			infer.Function[*resources.GetMonitorStatus, resources.GetMonitorStatusArgs, resources.GetMonitorStatusResult](&resources.GetMonitorStatus{}),
			infer.Function[*resources.GetMonitor, resources.LookupByNameArgs, resources.LookupResult](&resources.GetMonitor{}),
			infer.Function[*resources.GetTeam, resources.LookupByNameArgs, resources.LookupResult](&resources.GetTeam{}),
			infer.Function[*resources.GetLabel, resources.LookupByNameArgs, resources.LookupResult](&resources.GetLabel{}),
			infer.Function[*resources.GetStatusPage, resources.LookupByNameArgs, resources.LookupResult](&resources.GetStatusPage{}),
			infer.Function[*resources.GetIncidentState, resources.LookupByNameArgs, resources.LookupResult](&resources.GetIncidentState{}),
			infer.Function[*resources.GetIncidentSeverity, resources.LookupByNameArgs, resources.LookupResult](&resources.GetIncidentSeverity{}),
			infer.Function[*resources.GetAlertState, resources.LookupByNameArgs, resources.LookupResult](&resources.GetAlertState{}),
			infer.Function[*resources.GetAlertSeverity, resources.LookupByNameArgs, resources.LookupResult](&resources.GetAlertSeverity{}),
			infer.Function[*resources.GetOnCallDutyPolicy, resources.LookupByNameArgs, resources.LookupResult](&resources.GetOnCallDutyPolicy{}),
			infer.Function[*resources.GetMonitorSecret, resources.LookupByNameArgs, resources.LookupResult](&resources.GetMonitorSecret{}),
		).
		Build()
	if err != nil {
		panic(err)
	}
	return prov
}
