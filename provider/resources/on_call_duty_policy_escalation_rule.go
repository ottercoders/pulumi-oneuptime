package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicyEscalationRule struct{}

type OnCallDutyPolicyEscalationRuleArgs struct {
	ProjectID              *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyID     string  `pulumi:"onCallDutyPolicyId" json:"onCallDutyPolicyId"`
	Name                   *string `pulumi:"name,optional" json:"name,omitempty"`
	Description            *string `pulumi:"description,optional" json:"description,omitempty"`
	Order                  int     `pulumi:"order" json:"order"`
	EscalateAfterInMinutes *int    `pulumi:"escalateAfterInMinutes,optional" json:"escalateAfterInMinutes,omitempty"`
}

type OnCallDutyPolicyEscalationRuleState struct {
	OnCallDutyPolicyEscalationRuleArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicyEscalationRule)(nil)

func (e *OnCallDutyPolicyEscalationRule) Annotate(a infer.Annotator) {
	a.Describe(e, "Manages a OneUptime On-Call Duty Policy Escalation Rule. Defines escalation timing and order within an on-call policy.")
}

func (e *OnCallDutyPolicyEscalationRule) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyEscalationRuleArgs]) (infer.CreateResponse[OnCallDutyPolicyEscalationRuleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyEscalationRuleState{OnCallDutyPolicyEscalationRuleArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "on-call-duty-policy-escalation-rule", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}
	state.OnCallDutyPolicyEscalationRuleArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyEscalationRuleState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (e *OnCallDutyPolicyEscalationRule) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyEscalationRuleArgs, OnCallDutyPolicyEscalationRuleState]) (infer.ReadResponse[OnCallDutyPolicyEscalationRuleArgs, OnCallDutyPolicyEscalationRuleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyEscalationRuleArgs, OnCallDutyPolicyEscalationRuleState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleArgs, OnCallDutyPolicyEscalationRuleState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleArgs, OnCallDutyPolicyEscalationRuleState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyEscalationRuleArgs, OnCallDutyPolicyEscalationRuleState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyEscalationRuleArgs,
		State:  state,
	}, nil
}

func (e *OnCallDutyPolicyEscalationRule) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyEscalationRuleArgs, OnCallDutyPolicyEscalationRuleState]) (infer.UpdateResponse[OnCallDutyPolicyEscalationRuleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleState]{
			Output: OnCallDutyPolicyEscalationRuleState{OnCallDutyPolicyEscalationRuleArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy-escalation-rule", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleState]{}, err
	}

	return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleState]{
		Output: state,
	}, nil
}

func (e *OnCallDutyPolicyEscalationRule) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyEscalationRuleState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-escalation-rule", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
