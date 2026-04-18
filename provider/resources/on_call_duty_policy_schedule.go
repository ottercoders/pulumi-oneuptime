package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicySchedule struct{}

type OnCallDutyPolicyScheduleArgs struct {
	ProjectID          *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyID string  `pulumi:"onCallDutyPolicyId" json:"onCallDutyPolicyId"`
	Name               string  `pulumi:"name" json:"name"`
	Description        *string `pulumi:"description,optional" json:"description,omitempty"`
}

type OnCallDutyPolicyScheduleState struct {
	OnCallDutyPolicyScheduleArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicySchedule)(nil)

func (s *OnCallDutyPolicySchedule) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime On-Call Duty Policy Schedule. Defines rotation schedules within an on-call policy.")
}

func (s *OnCallDutyPolicySchedule) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyScheduleArgs]) (infer.CreateResponse[OnCallDutyPolicyScheduleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyScheduleState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyScheduleState{OnCallDutyPolicyScheduleArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyScheduleState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyScheduleState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "on-call-duty-policy-schedule", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyScheduleState]{}, err
	}

	var state OnCallDutyPolicyScheduleState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyScheduleState]{}, err
	}
	state.OnCallDutyPolicyScheduleArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyScheduleState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *OnCallDutyPolicySchedule) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyScheduleArgs, OnCallDutyPolicyScheduleState]) (infer.ReadResponse[OnCallDutyPolicyScheduleArgs, OnCallDutyPolicyScheduleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-schedule", req.ID, SelectFields(OnCallDutyPolicyScheduleState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyScheduleArgs, OnCallDutyPolicyScheduleState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyScheduleArgs, OnCallDutyPolicyScheduleState]{}, err
	}

	var state OnCallDutyPolicyScheduleState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyScheduleArgs, OnCallDutyPolicyScheduleState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyScheduleArgs, OnCallDutyPolicyScheduleState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyScheduleArgs,
		State:  state,
	}, nil
}

func (s *OnCallDutyPolicySchedule) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyScheduleArgs, OnCallDutyPolicyScheduleState]) (infer.UpdateResponse[OnCallDutyPolicyScheduleState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyScheduleState]{
			Output: OnCallDutyPolicyScheduleState{OnCallDutyPolicyScheduleArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyScheduleState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy-schedule", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyScheduleState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-schedule", req.ID, SelectFields(OnCallDutyPolicyScheduleState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyScheduleState]{}, err
	}

	var state OnCallDutyPolicyScheduleState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyScheduleState]{}, err
	}

	return infer.UpdateResponse[OnCallDutyPolicyScheduleState]{
		Output: state,
	}, nil
}

func (s *OnCallDutyPolicySchedule) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyScheduleState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-schedule", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
