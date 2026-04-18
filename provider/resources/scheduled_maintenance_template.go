package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ScheduledMaintenanceTemplate struct{}

type ScheduledMaintenanceTemplateArgs struct {
	ProjectID           *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	TemplateName        string   `pulumi:"templateName" json:"templateName"`
	TemplateDescription *string  `pulumi:"templateDescription,optional" json:"templateDescription,omitempty"`
	Title               *string  `pulumi:"title,optional" json:"title,omitempty"`
	Description         *string  `pulumi:"description,optional" json:"description,omitempty"`
	Monitors            []string `pulumi:"monitors,optional" json:"monitors,omitempty"`
	StatusPages         []string `pulumi:"statusPages,optional" json:"statusPages,omitempty"`
	Labels              []string `pulumi:"labels,optional" json:"labels,omitempty"`
	ChangeMonitorStatusToID *string `pulumi:"changeMonitorStatusToId,optional" json:"changeMonitorStatusToId,omitempty"`
}

type ScheduledMaintenanceTemplateState struct {
	ScheduledMaintenanceTemplateArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ScheduledMaintenanceTemplate)(nil)

func (t *ScheduledMaintenanceTemplate) Annotate(a infer.Annotator) {
	a.Describe(t, "Manages a OneUptime Scheduled Maintenance Template.")
}

func (t *ScheduledMaintenanceTemplate) Create(ctx context.Context, req infer.CreateRequest[ScheduledMaintenanceTemplateArgs]) (infer.CreateResponse[ScheduledMaintenanceTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ScheduledMaintenanceTemplateState]{
			ID:     "preview-id",
			Output: ScheduledMaintenanceTemplateState{ScheduledMaintenanceTemplateArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTemplateState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTemplateState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "scheduled-maintenance-template", data)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTemplateState]{}, err
	}

	var state ScheduledMaintenanceTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ScheduledMaintenanceTemplateState]{}, err
	}
	state.ScheduledMaintenanceTemplateArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ScheduledMaintenanceTemplateState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (t *ScheduledMaintenanceTemplate) Read(ctx context.Context, req infer.ReadRequest[ScheduledMaintenanceTemplateArgs, ScheduledMaintenanceTemplateState]) (infer.ReadResponse[ScheduledMaintenanceTemplateArgs, ScheduledMaintenanceTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "scheduled-maintenance-template", req.ID, SelectFields(ScheduledMaintenanceTemplateState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ScheduledMaintenanceTemplateArgs, ScheduledMaintenanceTemplateState]{}, nil
		}
		return infer.ReadResponse[ScheduledMaintenanceTemplateArgs, ScheduledMaintenanceTemplateState]{}, err
	}

	var state ScheduledMaintenanceTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ScheduledMaintenanceTemplateArgs, ScheduledMaintenanceTemplateState]{}, err
	}

	return infer.ReadResponse[ScheduledMaintenanceTemplateArgs, ScheduledMaintenanceTemplateState]{
		ID:     state.ResourceID,
		Inputs: state.ScheduledMaintenanceTemplateArgs,
		State:  state,
	}, nil
}

func (t *ScheduledMaintenanceTemplate) Update(ctx context.Context, req infer.UpdateRequest[ScheduledMaintenanceTemplateArgs, ScheduledMaintenanceTemplateState]) (infer.UpdateResponse[ScheduledMaintenanceTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ScheduledMaintenanceTemplateState]{
			Output: ScheduledMaintenanceTemplateState{ScheduledMaintenanceTemplateArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTemplateState]{}, err
	}

	if err := c.UpdateResource(ctx, "scheduled-maintenance-template", req.ID, data); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTemplateState]{}, err
	}

	result, err := c.ReadResource(ctx, "scheduled-maintenance-template", req.ID, SelectFields(ScheduledMaintenanceTemplateState{}))
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTemplateState]{}, err
	}

	var state ScheduledMaintenanceTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceTemplateState]{}, err
	}

	return infer.UpdateResponse[ScheduledMaintenanceTemplateState]{
		Output: state,
	}, nil
}

func (t *ScheduledMaintenanceTemplate) Delete(ctx context.Context, req infer.DeleteRequest[ScheduledMaintenanceTemplateState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "scheduled-maintenance-template", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
