package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallScheduleLayer struct{}

type OnCallScheduleLayerArgs struct {
	ProjectID                  *string                `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyScheduleID string                 `pulumi:"onCallDutyPolicyScheduleId" json:"onCallDutyPolicyScheduleId"`
	Name                       string                 `pulumi:"name" json:"name"`
	Description                *string                `pulumi:"description,optional" json:"description,omitempty"`
	Order                      *int                   `pulumi:"order,optional" json:"order,omitempty"`
	StartsAt                   *string                `pulumi:"startsAt,optional" json:"startsAt,omitempty"`
	Rotation                   map[string]interface{} `pulumi:"rotation,optional" json:"rotation,omitempty"`
	RestrictionTimes           map[string]interface{} `pulumi:"restrictionTimes,optional" json:"restrictionTimes,omitempty"`
}

type OnCallScheduleLayerState struct {
	OnCallScheduleLayerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallScheduleLayer)(nil)

func (r *OnCallScheduleLayer) Annotate(a infer.Annotator) {
	a.Describe(r, "Manages an on-call schedule layer (rotation level within a schedule).")
}

func (r *OnCallScheduleLayer) Create(ctx context.Context, req infer.CreateRequest[OnCallScheduleLayerArgs]) (infer.CreateResponse[OnCallScheduleLayerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallScheduleLayerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallScheduleLayerState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[OnCallScheduleLayerState]{
			ID:     "preview-id",
			Output: OnCallScheduleLayerState{OnCallScheduleLayerArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "on-call-duty-policy-schedule-layer", data)
	if err != nil {
		return infer.CreateResponse[OnCallScheduleLayerState]{}, err
	}

	var state OnCallScheduleLayerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallScheduleLayerState]{}, err
	}
	state.OnCallScheduleLayerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallScheduleLayerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *OnCallScheduleLayer) Read(ctx context.Context, req infer.ReadRequest[OnCallScheduleLayerArgs, OnCallScheduleLayerState]) (infer.ReadResponse[OnCallScheduleLayerArgs, OnCallScheduleLayerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-schedule-layer", req.ID, SelectFields(OnCallScheduleLayerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallScheduleLayerArgs, OnCallScheduleLayerState]{}, nil
		}
		return infer.ReadResponse[OnCallScheduleLayerArgs, OnCallScheduleLayerState]{}, err
	}

	var state OnCallScheduleLayerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallScheduleLayerArgs, OnCallScheduleLayerState]{}, err
	}

	return infer.ReadResponse[OnCallScheduleLayerArgs, OnCallScheduleLayerState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallScheduleLayerArgs,
		State:  state,
	}, nil
}

func (r *OnCallScheduleLayer) Update(ctx context.Context, req infer.UpdateRequest[OnCallScheduleLayerArgs, OnCallScheduleLayerState]) (infer.UpdateResponse[OnCallScheduleLayerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallScheduleLayerState]{
			Output: OnCallScheduleLayerState{OnCallScheduleLayerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy-schedule-layer", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-schedule-layer", req.ID, SelectFields(OnCallScheduleLayerState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerState]{}, err
	}

	var state OnCallScheduleLayerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallScheduleLayerState]{}, err
	}

	return infer.UpdateResponse[OnCallScheduleLayerState]{
		Output: state,
	}, nil
}

func (r *OnCallScheduleLayer) Delete(ctx context.Context, req infer.DeleteRequest[OnCallScheduleLayerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-schedule-layer", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
