package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorGroup struct{}

type MonitorGroupArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string  `pulumi:"name" json:"name"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
}

type MonitorGroupState struct {
	MonitorGroupArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorGroup)(nil)

func (m *MonitorGroup) Annotate(a infer.Annotator) {
	a.Describe(m, "Manages a OneUptime Monitor Group resource.")
}

func (m *MonitorGroup) Create(ctx context.Context, req infer.CreateRequest[MonitorGroupArgs]) (infer.CreateResponse[MonitorGroupState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[MonitorGroupState]{
			ID:     "preview-id",
			Output: MonitorGroupState{MonitorGroupArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorGroupState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorGroupState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "monitor-group", data)
	if err != nil {
		return infer.CreateResponse[MonitorGroupState]{}, err
	}

	var state MonitorGroupState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorGroupState]{}, err
	}
	state.MonitorGroupArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorGroupState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (m *MonitorGroup) Read(ctx context.Context, req infer.ReadRequest[MonitorGroupArgs, MonitorGroupState]) (infer.ReadResponse[MonitorGroupArgs, MonitorGroupState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor-group", req.ID, SelectFields(MonitorGroupState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorGroupArgs, MonitorGroupState]{}, nil
		}
		return infer.ReadResponse[MonitorGroupArgs, MonitorGroupState]{}, err
	}

	var state MonitorGroupState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorGroupArgs, MonitorGroupState]{}, err
	}

	return infer.ReadResponse[MonitorGroupArgs, MonitorGroupState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorGroupArgs,
		State:  state,
	}, nil
}

func (m *MonitorGroup) Update(ctx context.Context, req infer.UpdateRequest[MonitorGroupArgs, MonitorGroupState]) (infer.UpdateResponse[MonitorGroupState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorGroupState]{
			Output: MonitorGroupState{MonitorGroupArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorGroupState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor-group", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorGroupState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor-group", req.ID, SelectFields(MonitorGroupState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorGroupState]{}, err
	}

	var state MonitorGroupState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorGroupState]{}, err
	}

	return infer.UpdateResponse[MonitorGroupState]{
		Output: state,
	}, nil
}

func (m *MonitorGroup) Delete(ctx context.Context, req infer.DeleteRequest[MonitorGroupState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-group", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
