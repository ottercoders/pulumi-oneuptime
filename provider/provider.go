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
		WithConfig(infer.Config[*resources.Config](&resources.Config{})).
		WithResources(
			infer.Resource[*resources.Team, resources.TeamArgs, resources.TeamState](&resources.Team{}),
			infer.Resource[*resources.Monitor, resources.MonitorArgs, resources.MonitorState](&resources.Monitor{}),
			infer.Resource[*resources.StatusPage, resources.StatusPageArgs, resources.StatusPageState](&resources.StatusPage{}),
			infer.Resource[*resources.Incident, resources.IncidentArgs, resources.IncidentState](&resources.Incident{}),
			infer.Resource[*resources.OnCallDutyPolicy, resources.OnCallDutyPolicyArgs, resources.OnCallDutyPolicyState](&resources.OnCallDutyPolicy{}),
		).
		Build()
	if err != nil {
		panic(err)
	}
	return prov
}
