package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicy struct{}

type OnCallDutyPolicyArgs struct {
	ProjectID                                *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name                                     string  `pulumi:"name" json:"name"`
	Description                              *string `pulumi:"description,optional" json:"description,omitempty"`
	RepeatPolicyIfNoOneAcknowledges          *bool   `pulumi:"repeatPolicyIfNoOneAcknowledges,optional" json:"repeatPolicyIfNoOneAcknowledges,omitempty"`
	RepeatPolicyIfNoOneAcknowledgesNoOfTimes *int    `pulumi:"repeatPolicyIfNoOneAcknowledgesNoOfTimes,optional" json:"repeatPolicyIfNoOneAcknowledgesNoOfTimes,omitempty"`
}

type OnCallDutyPolicyState struct {
	OnCallDutyPolicyArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicy)(nil)

func (o *OnCallDutyPolicy) Annotate(a infer.Annotator) {
	a.Describe(o, "Manages a OneUptime On-Call Duty Policy resource.")
}

func (o *OnCallDutyPolicy) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyArgs]) (infer.CreateResponse[OnCallDutyPolicyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyState{OnCallDutyPolicyArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "on-call-duty-policy", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyState]{}, err
	}

	var state OnCallDutyPolicyState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyState]{}, err
	}
	state.OnCallDutyPolicyArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (o *OnCallDutyPolicy) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyArgs, OnCallDutyPolicyState]) (infer.ReadResponse[OnCallDutyPolicyArgs, OnCallDutyPolicyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy", req.ID, SelectFields(OnCallDutyPolicyState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyArgs, OnCallDutyPolicyState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyArgs, OnCallDutyPolicyState]{}, err
	}

	var state OnCallDutyPolicyState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyArgs, OnCallDutyPolicyState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyArgs, OnCallDutyPolicyState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyArgs,
		State:  state,
	}, nil
}

func (o *OnCallDutyPolicy) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyArgs, OnCallDutyPolicyState]) (infer.UpdateResponse[OnCallDutyPolicyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyState]{
			Output: OnCallDutyPolicyState{OnCallDutyPolicyArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy", req.ID, SelectFields(OnCallDutyPolicyState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyState]{}, err
	}

	var state OnCallDutyPolicyState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyState]{}, err
	}

	return infer.UpdateResponse[OnCallDutyPolicyState]{
		Output: state,
	}, nil
}

func (o *OnCallDutyPolicy) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
