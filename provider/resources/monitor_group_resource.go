package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorGroupResource struct{}

type MonitorGroupResourceArgs struct {
	ProjectID      *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	MonitorGroupID string  `pulumi:"monitorGroupId" json:"monitorGroupId"`
	MonitorID      string  `pulumi:"monitorId" json:"monitorId"`
}

type MonitorGroupResourceState struct {
	MonitorGroupResourceArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorGroupResource)(nil)

func (r *MonitorGroupResource) Annotate(a infer.Annotator) {
	a.Describe(r, "Adds a Monitor to a Monitor Group.")
}

func (r *MonitorGroupResource) Create(ctx context.Context, req infer.CreateRequest[MonitorGroupResourceArgs]) (infer.CreateResponse[MonitorGroupResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[MonitorGroupResourceState]{
			ID:     "preview-id",
			Output: MonitorGroupResourceState{MonitorGroupResourceArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorGroupResourceState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorGroupResourceState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "monitor-group-resource", data)
	if err != nil {
		return infer.CreateResponse[MonitorGroupResourceState]{}, err
	}

	var state MonitorGroupResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorGroupResourceState]{}, err
	}
	state.MonitorGroupResourceArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorGroupResourceState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *MonitorGroupResource) Read(ctx context.Context, req infer.ReadRequest[MonitorGroupResourceArgs, MonitorGroupResourceState]) (infer.ReadResponse[MonitorGroupResourceArgs, MonitorGroupResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor-group-resource", req.ID, SelectFields(MonitorGroupResourceState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorGroupResourceArgs, MonitorGroupResourceState]{}, nil
		}
		return infer.ReadResponse[MonitorGroupResourceArgs, MonitorGroupResourceState]{}, err
	}

	var state MonitorGroupResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorGroupResourceArgs, MonitorGroupResourceState]{}, err
	}

	return infer.ReadResponse[MonitorGroupResourceArgs, MonitorGroupResourceState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorGroupResourceArgs,
		State:  state,
	}, nil
}

func (r *MonitorGroupResource) Update(ctx context.Context, req infer.UpdateRequest[MonitorGroupResourceArgs, MonitorGroupResourceState]) (infer.UpdateResponse[MonitorGroupResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorGroupResourceState]{
			Output: MonitorGroupResourceState{MonitorGroupResourceArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorGroupResourceState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor-group-resource", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorGroupResourceState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor-group-resource", req.ID, SelectFields(MonitorGroupResourceState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorGroupResourceState]{}, err
	}

	var state MonitorGroupResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorGroupResourceState]{}, err
	}

	return infer.UpdateResponse[MonitorGroupResourceState]{
		Output: state,
	}, nil
}

func (r *MonitorGroupResource) Delete(ctx context.Context, req infer.DeleteRequest[MonitorGroupResourceState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-group-resource", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
