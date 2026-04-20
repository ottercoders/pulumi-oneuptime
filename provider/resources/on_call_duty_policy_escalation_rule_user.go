package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicyEscalationRuleUser struct{}

type OnCallDutyPolicyEscalationRuleUserArgs struct {
	ProjectID                        *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyID               string  `pulumi:"onCallDutyPolicyId" json:"onCallDutyPolicyId"`
	OnCallDutyPolicyEscalationRuleID string  `pulumi:"onCallDutyPolicyEscalationRuleId" json:"onCallDutyPolicyEscalationRuleId"`
	UserID                           string  `pulumi:"userId" json:"userId"`
}

type OnCallDutyPolicyEscalationRuleUserState struct {
	OnCallDutyPolicyEscalationRuleUserArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicyEscalationRuleUser)(nil)

func (r *OnCallDutyPolicyEscalationRuleUser) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a specific user as a paging target on an on-call duty policy escalation rule. Without at least one *RuleUser / *RuleTeam / *RuleSchedule attached, an escalation rule cannot page anyone.")
}

func (r *OnCallDutyPolicyEscalationRuleUser) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyEscalationRuleUserArgs]) (infer.CreateResponse[OnCallDutyPolicyEscalationRuleUserState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleUserState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyEscalationRuleUserState{OnCallDutyPolicyEscalationRuleUserArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleUserState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleUserState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "on-call-duty-policy-escalation-rule-user", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleUserState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleUserState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleUserState]{}, err
	}
	state.OnCallDutyPolicyEscalationRuleUserArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyEscalationRuleUserState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyEscalationRuleUser) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyEscalationRuleUserArgs, OnCallDutyPolicyEscalationRuleUserState]) (infer.ReadResponse[OnCallDutyPolicyEscalationRuleUserArgs, OnCallDutyPolicyEscalationRuleUserState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule-user", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleUserState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyEscalationRuleUserArgs, OnCallDutyPolicyEscalationRuleUserState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleUserArgs, OnCallDutyPolicyEscalationRuleUserState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleUserState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleUserArgs, OnCallDutyPolicyEscalationRuleUserState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyEscalationRuleUserArgs, OnCallDutyPolicyEscalationRuleUserState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyEscalationRuleUserArgs,
		State:  state,
	}, nil
}

// The join entity has no mutable columns on the OneUptime API (all FKs are
// create-only; the update ACL is empty upstream). Update is a no-op that
// just re-reads the current server state so Pulumi's diff loop stays
// coherent if a user marks a field changed.
func (r *OnCallDutyPolicyEscalationRuleUser) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyEscalationRuleUserArgs, OnCallDutyPolicyEscalationRuleUserState]) (infer.UpdateResponse[OnCallDutyPolicyEscalationRuleUserState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleUserState]{
			Output: OnCallDutyPolicyEscalationRuleUserState{OnCallDutyPolicyEscalationRuleUserArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule-user", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleUserState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleUserState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleUserState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleUserState]{}, err
	}
	state.OnCallDutyPolicyEscalationRuleUserArgs = req.Inputs

	return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleUserState]{
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyEscalationRuleUser) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyEscalationRuleUserState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-escalation-rule-user", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
