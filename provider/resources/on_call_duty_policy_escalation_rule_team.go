package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicyEscalationRuleTeam struct{}

type OnCallDutyPolicyEscalationRuleTeamArgs struct {
	ProjectID                        *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyID               string  `pulumi:"onCallDutyPolicyId" json:"onCallDutyPolicyId"`
	OnCallDutyPolicyEscalationRuleID string  `pulumi:"onCallDutyPolicyEscalationRuleId" json:"onCallDutyPolicyEscalationRuleId"`
	TeamID                           string  `pulumi:"teamId" json:"teamId"`
}

type OnCallDutyPolicyEscalationRuleTeamState struct {
	OnCallDutyPolicyEscalationRuleTeamArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicyEscalationRuleTeam)(nil)

func (r *OnCallDutyPolicyEscalationRuleTeam) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a whole team as a paging target on an on-call duty policy escalation rule. Every member of the team is paged according to their notification settings.")
}

func (r *OnCallDutyPolicyEscalationRuleTeam) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyEscalationRuleTeamArgs]) (infer.CreateResponse[OnCallDutyPolicyEscalationRuleTeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleTeamState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyEscalationRuleTeamState{OnCallDutyPolicyEscalationRuleTeamArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "on-call-duty-policy-escalation-rule-team", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleTeamState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}
	state.OnCallDutyPolicyEscalationRuleTeamArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyEscalationRuleTeamState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyEscalationRuleTeam) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyEscalationRuleTeamArgs, OnCallDutyPolicyEscalationRuleTeamState]) (infer.ReadResponse[OnCallDutyPolicyEscalationRuleTeamArgs, OnCallDutyPolicyEscalationRuleTeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule-team", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleTeamState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyEscalationRuleTeamArgs, OnCallDutyPolicyEscalationRuleTeamState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleTeamArgs, OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleTeamState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleTeamArgs, OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyEscalationRuleTeamArgs, OnCallDutyPolicyEscalationRuleTeamState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyEscalationRuleTeamArgs,
		State:  state,
	}, nil
}

func (r *OnCallDutyPolicyEscalationRuleTeam) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyEscalationRuleTeamArgs, OnCallDutyPolicyEscalationRuleTeamState]) (infer.UpdateResponse[OnCallDutyPolicyEscalationRuleTeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleTeamState]{
			Output: OnCallDutyPolicyEscalationRuleTeamState{OnCallDutyPolicyEscalationRuleTeamArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule-team", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleTeamState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleTeamState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleTeamState]{}, err
	}
	state.OnCallDutyPolicyEscalationRuleTeamArgs = req.Inputs

	return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleTeamState]{
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyEscalationRuleTeam) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyEscalationRuleTeamState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-escalation-rule-team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
