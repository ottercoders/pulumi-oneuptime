package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorStatus struct{}

type MonitorStatusArgs struct {
	ProjectID        *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name             string  `pulumi:"name" json:"name"`
	Description      *string `pulumi:"description,optional" json:"description,omitempty"`
	Color            string  `pulumi:"color" json:"color"`
	Priority         *int    `pulumi:"priority,optional" json:"priority,omitempty"`
	IsOperationalState *bool `pulumi:"isOperationalState,optional" json:"isOperationalState,omitempty"`
	IsDownState      *bool   `pulumi:"isDownState,optional" json:"isDownState,omitempty"`
	IsDegradedState  *bool   `pulumi:"isDegradedState,optional" json:"isDegradedState,omitempty"`
}

type MonitorStatusState struct {
	MonitorStatusArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorStatus)(nil)

func (r *MonitorStatus) Annotate(a infer.Annotator) {
	a.Describe(r, "Manages a custom monitor status (e.g., Operational, Degraded, Offline).")
}

func (r *MonitorStatus) Create(ctx context.Context, req infer.CreateRequest[MonitorStatusArgs]) (infer.CreateResponse[MonitorStatusState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[MonitorStatusState]{
			ID:     "preview-id",
			Output: MonitorStatusState{MonitorStatusArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorStatusState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorStatusState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "monitor-status", data)
	if err != nil {
		return infer.CreateResponse[MonitorStatusState]{}, err
	}

	var state MonitorStatusState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorStatusState]{}, err
	}
	state.MonitorStatusArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorStatusState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *MonitorStatus) Read(ctx context.Context, req infer.ReadRequest[MonitorStatusArgs, MonitorStatusState]) (infer.ReadResponse[MonitorStatusArgs, MonitorStatusState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor-status", req.ID, SelectFields(MonitorStatusState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorStatusArgs, MonitorStatusState]{}, nil
		}
		return infer.ReadResponse[MonitorStatusArgs, MonitorStatusState]{}, err
	}

	var state MonitorStatusState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorStatusArgs, MonitorStatusState]{}, err
	}

	return infer.ReadResponse[MonitorStatusArgs, MonitorStatusState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorStatusArgs,
		State:  state,
	}, nil
}

func (r *MonitorStatus) Update(ctx context.Context, req infer.UpdateRequest[MonitorStatusArgs, MonitorStatusState]) (infer.UpdateResponse[MonitorStatusState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorStatusState]{
			Output: MonitorStatusState{MonitorStatusArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorStatusState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor-status", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorStatusState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor-status", req.ID, SelectFields(MonitorStatusState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorStatusState]{}, err
	}

	var state MonitorStatusState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorStatusState]{}, err
	}

	return infer.UpdateResponse[MonitorStatusState]{
		Output: state,
	}, nil
}

func (r *MonitorStatus) Delete(ctx context.Context, req infer.DeleteRequest[MonitorStatusState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-status", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
