package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorUserOwner struct{}

type MonitorUserOwnerArgs struct {
	ProjectID *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	MonitorID string  `pulumi:"monitorId" json:"monitorId"`
	UserID    string  `pulumi:"userId" json:"userId"`
}

type MonitorUserOwnerState struct {
	MonitorUserOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorUserOwner)(nil)

func (o *MonitorUserOwner) Annotate(a infer.Annotator) {
	a.Describe(o, "Assigns a OneUptime User as owner of a Monitor.")
}

func (o *MonitorUserOwner) Create(ctx context.Context, req infer.CreateRequest[MonitorUserOwnerArgs]) (infer.CreateResponse[MonitorUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorUserOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorUserOwnerState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[MonitorUserOwnerState]{
			ID:     "preview-id",
			Output: MonitorUserOwnerState{MonitorUserOwnerArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "monitor-owner-user", data)
	if err != nil {
		return infer.CreateResponse[MonitorUserOwnerState]{}, err
	}

	var state MonitorUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorUserOwnerState]{}, err
	}
	state.MonitorUserOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorUserOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (o *MonitorUserOwner) Read(ctx context.Context, req infer.ReadRequest[MonitorUserOwnerArgs, MonitorUserOwnerState]) (infer.ReadResponse[MonitorUserOwnerArgs, MonitorUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor-owner-user", req.ID, SelectFields(MonitorUserOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorUserOwnerArgs, MonitorUserOwnerState]{}, nil
		}
		return infer.ReadResponse[MonitorUserOwnerArgs, MonitorUserOwnerState]{}, err
	}

	var state MonitorUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorUserOwnerArgs, MonitorUserOwnerState]{}, err
	}

	return infer.ReadResponse[MonitorUserOwnerArgs, MonitorUserOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorUserOwnerArgs,
		State:  state,
	}, nil
}

func (o *MonitorUserOwner) Update(ctx context.Context, req infer.UpdateRequest[MonitorUserOwnerArgs, MonitorUserOwnerState]) (infer.UpdateResponse[MonitorUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorUserOwnerState]{
			Output: MonitorUserOwnerState{MonitorUserOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorUserOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor-owner-user", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorUserOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor-owner-user", req.ID, SelectFields(MonitorUserOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorUserOwnerState]{}, err
	}

	var state MonitorUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorUserOwnerState]{}, err
	}

	return infer.UpdateResponse[MonitorUserOwnerState]{
		Output: state,
	}, nil
}

func (o *MonitorUserOwner) Delete(ctx context.Context, req infer.DeleteRequest[MonitorUserOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-owner-user", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
