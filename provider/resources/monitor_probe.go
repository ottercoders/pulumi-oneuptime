package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorProbe struct{}

type MonitorProbeArgs struct {
	ProjectID *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	MonitorID string  `pulumi:"monitorId" json:"monitorId"`
	ProbeID   string  `pulumi:"probeId" json:"probeId"`
	IsEnabled *bool   `pulumi:"isEnabled,optional" json:"isEnabled,omitempty"`
}

type MonitorProbeState struct {
	MonitorProbeArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorProbe)(nil)

func (r *MonitorProbe) Annotate(a infer.Annotator) {
	a.Describe(r, "Links a monitoring probe to a specific monitor.")
}

func (r *MonitorProbe) Create(ctx context.Context, req infer.CreateRequest[MonitorProbeArgs]) (infer.CreateResponse[MonitorProbeState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorProbeState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorProbeState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[MonitorProbeState]{
			ID:     "preview-id",
			Output: MonitorProbeState{MonitorProbeArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "monitor-probe", data)
	if err != nil {
		return infer.CreateResponse[MonitorProbeState]{}, err
	}

	var state MonitorProbeState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorProbeState]{}, err
	}
	state.MonitorProbeArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorProbeState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *MonitorProbe) Read(ctx context.Context, req infer.ReadRequest[MonitorProbeArgs, MonitorProbeState]) (infer.ReadResponse[MonitorProbeArgs, MonitorProbeState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor-probe", req.ID, SelectFields(MonitorProbeState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorProbeArgs, MonitorProbeState]{}, nil
		}
		return infer.ReadResponse[MonitorProbeArgs, MonitorProbeState]{}, err
	}

	var state MonitorProbeState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorProbeArgs, MonitorProbeState]{}, err
	}

	return infer.ReadResponse[MonitorProbeArgs, MonitorProbeState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorProbeArgs,
		State:  state,
	}, nil
}

func (r *MonitorProbe) Update(ctx context.Context, req infer.UpdateRequest[MonitorProbeArgs, MonitorProbeState]) (infer.UpdateResponse[MonitorProbeState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorProbeState]{
			Output: MonitorProbeState{MonitorProbeArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorProbeState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor-probe", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorProbeState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor-probe", req.ID, SelectFields(MonitorProbeState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorProbeState]{}, err
	}

	var state MonitorProbeState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorProbeState]{}, err
	}

	return infer.UpdateResponse[MonitorProbeState]{
		Output: state,
	}, nil
}

func (r *MonitorProbe) Delete(ctx context.Context, req infer.DeleteRequest[MonitorProbeState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-probe", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
