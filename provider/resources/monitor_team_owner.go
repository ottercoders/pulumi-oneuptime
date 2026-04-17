package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorTeamOwner struct{}

type MonitorTeamOwnerArgs struct {
	ProjectID *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	MonitorID string  `pulumi:"monitorId" json:"monitorId"`
	TeamID    string  `pulumi:"teamId" json:"teamId"`
}

type MonitorTeamOwnerState struct {
	MonitorTeamOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorTeamOwner)(nil)

func (o *MonitorTeamOwner) Annotate(a infer.Annotator) {
	a.Describe(o, "Assigns a OneUptime Team as owner of a Monitor.")
}

func (o *MonitorTeamOwner) Create(ctx context.Context, req infer.CreateRequest[MonitorTeamOwnerArgs]) (infer.CreateResponse[MonitorTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorTeamOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorTeamOwnerState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[MonitorTeamOwnerState]{
			ID:     "preview-id",
			Output: MonitorTeamOwnerState{MonitorTeamOwnerArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "monitor-owner-team", data)
	if err != nil {
		return infer.CreateResponse[MonitorTeamOwnerState]{}, err
	}

	var state MonitorTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorTeamOwnerState]{}, err
	}
	state.MonitorTeamOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorTeamOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (o *MonitorTeamOwner) Read(ctx context.Context, req infer.ReadRequest[MonitorTeamOwnerArgs, MonitorTeamOwnerState]) (infer.ReadResponse[MonitorTeamOwnerArgs, MonitorTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor-owner-team", req.ID, SelectFields(MonitorTeamOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorTeamOwnerArgs, MonitorTeamOwnerState]{}, nil
		}
		return infer.ReadResponse[MonitorTeamOwnerArgs, MonitorTeamOwnerState]{}, err
	}

	var state MonitorTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorTeamOwnerArgs, MonitorTeamOwnerState]{}, err
	}

	return infer.ReadResponse[MonitorTeamOwnerArgs, MonitorTeamOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorTeamOwnerArgs,
		State:  state,
	}, nil
}

func (o *MonitorTeamOwner) Update(ctx context.Context, req infer.UpdateRequest[MonitorTeamOwnerArgs, MonitorTeamOwnerState]) (infer.UpdateResponse[MonitorTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorTeamOwnerState]{
			Output: MonitorTeamOwnerState{MonitorTeamOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorTeamOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor-owner-team", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorTeamOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor-owner-team", req.ID, SelectFields(MonitorTeamOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorTeamOwnerState]{}, err
	}

	var state MonitorTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorTeamOwnerState]{}, err
	}

	return infer.UpdateResponse[MonitorTeamOwnerState]{
		Output: state,
	}, nil
}

func (o *MonitorTeamOwner) Delete(ctx context.Context, req infer.DeleteRequest[MonitorTeamOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-owner-team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
