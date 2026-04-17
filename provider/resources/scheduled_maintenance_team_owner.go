package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ScheduledMaintenanceTeamOwner struct{}

type ScheduledMaintenanceTeamOwnerArgs struct {
	ProjectID              *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	ScheduledMaintenanceID string  `pulumi:"scheduledMaintenanceId" json:"scheduledMaintenanceId"`
	TeamID                 string  `pulumi:"teamId" json:"teamId"`
}

type ScheduledMaintenanceTeamOwnerState struct {
	ScheduledMaintenanceTeamOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ScheduledMaintenanceTeamOwner)(nil)

func (r *ScheduledMaintenanceTeamOwner) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a team as owner of a scheduled maintenance event.")
}

func (r *ScheduledMaintenanceTeamOwner) Create(ctx context.Context, req infer.CreateRequest[ScheduledMaintenanceTeamOwnerArgs]) (infer.CreateResponse[ScheduledMaintenanceTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[ScheduledMaintenanceTeamOwnerState]{
			ID:     "preview-id",
			Output: ScheduledMaintenanceTeamOwnerState{ScheduledMaintenanceTeamOwnerArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "scheduled-maintenance-owner-team", data)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}

	var state ScheduledMaintenanceTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}
	state.ScheduledMaintenanceTeamOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ScheduledMaintenanceTeamOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *ScheduledMaintenanceTeamOwner) Read(ctx context.Context, req infer.ReadRequest[ScheduledMaintenanceTeamOwnerArgs, ScheduledMaintenanceTeamOwnerState]) (infer.ReadResponse[ScheduledMaintenanceTeamOwnerArgs, ScheduledMaintenanceTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "scheduled-maintenance-owner-team", req.ID, SelectFields(ScheduledMaintenanceTeamOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ScheduledMaintenanceTeamOwnerArgs, ScheduledMaintenanceTeamOwnerState]{}, nil
		}
		return infer.ReadResponse[ScheduledMaintenanceTeamOwnerArgs, ScheduledMaintenanceTeamOwnerState]{}, err
	}

	var state ScheduledMaintenanceTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ScheduledMaintenanceTeamOwnerArgs, ScheduledMaintenanceTeamOwnerState]{}, err
	}

	return infer.ReadResponse[ScheduledMaintenanceTeamOwnerArgs, ScheduledMaintenanceTeamOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.ScheduledMaintenanceTeamOwnerArgs,
		State:  state,
	}, nil
}

func (r *ScheduledMaintenanceTeamOwner) Update(ctx context.Context, req infer.UpdateRequest[ScheduledMaintenanceTeamOwnerArgs, ScheduledMaintenanceTeamOwnerState]) (infer.UpdateResponse[ScheduledMaintenanceTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ScheduledMaintenanceTeamOwnerState]{
			Output: ScheduledMaintenanceTeamOwnerState{ScheduledMaintenanceTeamOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "scheduled-maintenance-owner-team", req.ID, data); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "scheduled-maintenance-owner-team", req.ID, SelectFields(ScheduledMaintenanceTeamOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}

	var state ScheduledMaintenanceTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTeamOwnerState]{}, err
	}

	return infer.UpdateResponse[ScheduledMaintenanceTeamOwnerState]{
		Output: state,
	}, nil
}

func (r *ScheduledMaintenanceTeamOwner) Delete(ctx context.Context, req infer.DeleteRequest[ScheduledMaintenanceTeamOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "scheduled-maintenance-owner-team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
