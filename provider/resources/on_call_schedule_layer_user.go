package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallScheduleLayerUser struct{}

type OnCallScheduleLayerUserArgs struct {
	ProjectID                      *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyScheduleLayerID string  `pulumi:"onCallDutyPolicyScheduleLayerId" json:"onCallDutyPolicyScheduleLayerId"`
	UserID                         string  `pulumi:"userId" json:"userId"`
	Order                          *int    `pulumi:"order,optional" json:"order,omitempty"`
}

type OnCallScheduleLayerUserState struct {
	OnCallScheduleLayerUserArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallScheduleLayerUser)(nil)

func (r *OnCallScheduleLayerUser) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a user to an on-call schedule layer.")
}

func (r *OnCallScheduleLayerUser) Create(ctx context.Context, req infer.CreateRequest[OnCallScheduleLayerUserArgs]) (infer.CreateResponse[OnCallScheduleLayerUserState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallScheduleLayerUserState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallScheduleLayerUserState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[OnCallScheduleLayerUserState]{
			ID:     "preview-id",
			Output: OnCallScheduleLayerUserState{OnCallScheduleLayerUserArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "on-call-duty-policy-schedule-layer-user", data)
	if err != nil {
		return infer.CreateResponse[OnCallScheduleLayerUserState]{}, err
	}

	var state OnCallScheduleLayerUserState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallScheduleLayerUserState]{}, err
	}
	state.OnCallScheduleLayerUserArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallScheduleLayerUserState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *OnCallScheduleLayerUser) Read(ctx context.Context, req infer.ReadRequest[OnCallScheduleLayerUserArgs, OnCallScheduleLayerUserState]) (infer.ReadResponse[OnCallScheduleLayerUserArgs, OnCallScheduleLayerUserState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-schedule-layer-user", req.ID, SelectFields(OnCallScheduleLayerUserState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallScheduleLayerUserArgs, OnCallScheduleLayerUserState]{}, nil
		}
		return infer.ReadResponse[OnCallScheduleLayerUserArgs, OnCallScheduleLayerUserState]{}, err
	}

	var state OnCallScheduleLayerUserState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallScheduleLayerUserArgs, OnCallScheduleLayerUserState]{}, err
	}

	return infer.ReadResponse[OnCallScheduleLayerUserArgs, OnCallScheduleLayerUserState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallScheduleLayerUserArgs,
		State:  state,
	}, nil
}

func (r *OnCallScheduleLayerUser) Update(ctx context.Context, req infer.UpdateRequest[OnCallScheduleLayerUserArgs, OnCallScheduleLayerUserState]) (infer.UpdateResponse[OnCallScheduleLayerUserState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallScheduleLayerUserState]{
			Output: OnCallScheduleLayerUserState{OnCallScheduleLayerUserArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerUserState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy-schedule-layer-user", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerUserState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-schedule-layer-user", req.ID, SelectFields(OnCallScheduleLayerUserState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerUserState]{}, err
	}

	var state OnCallScheduleLayerUserState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerUserState]{}, err
	}

	return infer.UpdateResponse[OnCallScheduleLayerUserState]{
		Output: state,
	}, nil
}

func (r *OnCallScheduleLayerUser) Delete(ctx context.Context, req infer.DeleteRequest[OnCallScheduleLayerUserState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-schedule-layer-user", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
