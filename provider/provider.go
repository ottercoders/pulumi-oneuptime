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
		).
		Build()
	if err != nil {
		panic(err)
	}
	return prov
}
