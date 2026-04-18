package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ScheduledMaintenanceState struct{}

type ScheduledMaintenanceStateArgs struct {
	ProjectID          *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name               string  `pulumi:"name" json:"name"`
	Description        *string `pulumi:"description,optional" json:"description,omitempty"`
	Color              string  `pulumi:"color" json:"color"`
	IsScheduledState   *bool   `pulumi:"isScheduledState,optional" json:"isScheduledState,omitempty"`
	IsOngoingState     *bool   `pulumi:"isOngoingState,optional" json:"isOngoingState,omitempty"`
	IsEndedState       *bool   `pulumi:"isEndedState,optional" json:"isEndedState,omitempty"`
	Order              int     `pulumi:"order" json:"order"`
}

type ScheduledMaintenanceStateState struct {
	ScheduledMaintenanceStateArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ScheduledMaintenanceState)(nil)

func (s *ScheduledMaintenanceState) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime Scheduled Maintenance State (workflow state for maintenance events).")
}

func (s *ScheduledMaintenanceState) Create(ctx context.Context, req infer.CreateRequest[ScheduledMaintenanceStateArgs]) (infer.CreateResponse[ScheduledMaintenanceStateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ScheduledMaintenanceStateState]{
			ID:     "preview-id",
			Output: ScheduledMaintenanceStateState{ScheduledMaintenanceStateArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceStateState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceStateState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "scheduled-maintenance-state", data)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceStateState]{}, err
	}

	var state ScheduledMaintenanceStateState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ScheduledMaintenanceStateState]{}, err
	}
	state.ScheduledMaintenanceStateArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ScheduledMaintenanceStateState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *ScheduledMaintenanceState) Read(ctx context.Context, req infer.ReadRequest[ScheduledMaintenanceStateArgs, ScheduledMaintenanceStateState]) (infer.ReadResponse[ScheduledMaintenanceStateArgs, ScheduledMaintenanceStateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "scheduled-maintenance-state", req.ID, SelectFields(ScheduledMaintenanceStateState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ScheduledMaintenanceStateArgs, ScheduledMaintenanceStateState]{}, nil
		}
		return infer.ReadResponse[ScheduledMaintenanceStateArgs, ScheduledMaintenanceStateState]{}, err
	}

	var state ScheduledMaintenanceStateState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ScheduledMaintenanceStateArgs, ScheduledMaintenanceStateState]{}, err
	}

	return infer.ReadResponse[ScheduledMaintenanceStateArgs, ScheduledMaintenanceStateState]{
		ID:     state.ResourceID,
		Inputs: state.ScheduledMaintenanceStateArgs,
		State:  state,
	}, nil
}

func (s *ScheduledMaintenanceState) Update(ctx context.Context, req infer.UpdateRequest[ScheduledMaintenanceStateArgs, ScheduledMaintenanceStateState]) (infer.UpdateResponse[ScheduledMaintenanceStateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ScheduledMaintenanceStateState]{
			Output: ScheduledMaintenanceStateState{ScheduledMaintenanceStateArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceStateState]{}, err
	}

	if err := c.UpdateResource(ctx, "scheduled-maintenance-state", req.ID, data); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceStateState]{}, err
	}

	result, err := c.ReadResource(ctx, "scheduled-maintenance-state", req.ID, SelectFields(ScheduledMaintenanceStateState{}))
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceStateState]{}, err
	}

	var state ScheduledMaintenanceStateState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceStateState]{}, err
	}

	return infer.UpdateResponse[ScheduledMaintenanceStateState]{
		Output: state,
	}, nil
}

func (s *ScheduledMaintenanceState) Delete(ctx context.Context, req infer.DeleteRequest[ScheduledMaintenanceStateState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "scheduled-maintenance-state", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
