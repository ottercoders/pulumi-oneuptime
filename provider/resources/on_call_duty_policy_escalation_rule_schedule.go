package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicyEscalationRuleSchedule struct{}

type OnCallDutyPolicyEscalationRuleScheduleArgs struct {
	ProjectID                        *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyID               string  `pulumi:"onCallDutyPolicyId" json:"onCallDutyPolicyId"`
	OnCallDutyPolicyEscalationRuleID string  `pulumi:"onCallDutyPolicyEscalationRuleId" json:"onCallDutyPolicyEscalationRuleId"`
	OnCallDutyPolicyScheduleID       string  `pulumi:"onCallDutyPolicyScheduleId" json:"onCallDutyPolicyScheduleId"`
}

type OnCallDutyPolicyEscalationRuleScheduleState struct {
	OnCallDutyPolicyEscalationRuleScheduleArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicyEscalationRuleSchedule)(nil)

func (l *OnCallDutyPolicyEscalationRuleSchedule) Annotate(a infer.Annotator) {
	a.Describe(l, "Links a OneUptime On-Call Duty Policy Schedule to an Escalation Rule.")
}

func (l *OnCallDutyPolicyEscalationRuleSchedule) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyEscalationRuleScheduleArgs]) (infer.CreateResponse[OnCallDutyPolicyEscalationRuleScheduleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyEscalationRuleScheduleState{OnCallDutyPolicyEscalationRuleScheduleArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "on-call-duty-policy-escalation-rule-schedule", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleScheduleState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}
	state.OnCallDutyPolicyEscalationRuleScheduleArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (l *OnCallDutyPolicyEscalationRuleSchedule) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyEscalationRuleScheduleArgs, OnCallDutyPolicyEscalationRuleScheduleState]) (infer.ReadResponse[OnCallDutyPolicyEscalationRuleScheduleArgs, OnCallDutyPolicyEscalationRuleScheduleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule-schedule", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleScheduleState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyEscalationRuleScheduleArgs, OnCallDutyPolicyEscalationRuleScheduleState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleScheduleArgs, OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleScheduleState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyEscalationRuleScheduleArgs, OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyEscalationRuleScheduleArgs, OnCallDutyPolicyEscalationRuleScheduleState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyEscalationRuleScheduleArgs,
		State:  state,
	}, nil
}

func (l *OnCallDutyPolicyEscalationRuleSchedule) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyEscalationRuleScheduleArgs, OnCallDutyPolicyEscalationRuleScheduleState]) (infer.UpdateResponse[OnCallDutyPolicyEscalationRuleScheduleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{
			Output: OnCallDutyPolicyEscalationRuleScheduleState{OnCallDutyPolicyEscalationRuleScheduleArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy-escalation-rule-schedule", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-escalation-rule-schedule", req.ID, SelectFields(OnCallDutyPolicyEscalationRuleScheduleState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	var state OnCallDutyPolicyEscalationRuleScheduleState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{}, err
	}

	return infer.UpdateResponse[OnCallDutyPolicyEscalationRuleScheduleState]{
		Output: state,
	}, nil
}

func (l *OnCallDutyPolicyEscalationRuleSchedule) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyEscalationRuleScheduleState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-escalation-rule-schedule", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
