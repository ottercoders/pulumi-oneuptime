package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Monitor struct{}

type MonitorArgs struct {
	ProjectID               *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name                    string  `pulumi:"name" json:"name"`
	Description             *string `pulumi:"description,optional" json:"description,omitempty"`
	MonitorType             string  `pulumi:"monitorType" json:"monitorType"`
	CurrentMonitorStatusID  string  `pulumi:"currentMonitorStatusId" json:"currentMonitorStatusId"`
	DisableActiveMonitoring *bool   `pulumi:"disableActiveMonitoring,optional" json:"disableActiveMonitoring,omitempty"`
	MonitoringInterval      *string `pulumi:"monitoringInterval,optional" json:"monitoringInterval,omitempty"`
}

type MonitorState struct {
	MonitorArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Monitor)(nil)

func (m *Monitor) Annotate(a infer.Annotator) {
	a.Describe(m, "Manages a OneUptime Monitor resource.")
}

func (m *Monitor) Create(ctx context.Context, req infer.CreateRequest[MonitorArgs]) (infer.CreateResponse[MonitorState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[MonitorState]{
			ID:     "preview-id",
			Output: MonitorState{MonitorArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "monitor", data)
	if err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}

	var state MonitorState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}
	state.MonitorArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (m *Monitor) Read(ctx context.Context, req infer.ReadRequest[MonitorArgs, MonitorState]) (infer.ReadResponse[MonitorArgs, MonitorState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor", req.ID, SelectFields(MonitorState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorArgs, MonitorState]{}, nil
		}
		return infer.ReadResponse[MonitorArgs, MonitorState]{}, err
	}

	var state MonitorState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorArgs, MonitorState]{}, err
	}

	return infer.ReadResponse[MonitorArgs, MonitorState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorArgs,
		State:  state,
	}, nil
}

func (m *Monitor) Update(ctx context.Context, req infer.UpdateRequest[MonitorArgs, MonitorState]) (infer.UpdateResponse[MonitorState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorState]{
			Output: MonitorState{MonitorArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor", req.ID, SelectFields(MonitorState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}

	var state MonitorState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}

	return infer.UpdateResponse[MonitorState]{
		Output: state,
	}, nil
}

func (m *Monitor) Delete(ctx context.Context, req infer.DeleteRequest[MonitorState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
