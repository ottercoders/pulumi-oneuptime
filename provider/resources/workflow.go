package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Workflow struct{}

type WorkflowArgs struct {
	ProjectID   *string                `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string                 `pulumi:"name" json:"name"`
	Description *string                `pulumi:"description,optional" json:"description,omitempty"`
	Graph       map[string]interface{} `pulumi:"graph,optional" json:"graph,omitempty"`
	TriggerID   *string                `pulumi:"triggerId,optional" json:"triggerId,omitempty"`
	TriggerArguments map[string]interface{} `pulumi:"triggerArguments,optional" json:"triggerArguments,omitempty"`
	IsEnabled   *bool                  `pulumi:"isEnabled,optional" json:"isEnabled,omitempty"`
}

type WorkflowState struct {
	WorkflowArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Workflow)(nil)

func (w *Workflow) Annotate(a infer.Annotator) {
	a.Describe(w, "Manages a OneUptime Workflow (automation flow).")
}

func (w *Workflow) Create(ctx context.Context, req infer.CreateRequest[WorkflowArgs]) (infer.CreateResponse[WorkflowState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[WorkflowState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[WorkflowState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[WorkflowState]{
			ID:     "preview-id",
			Output: WorkflowState{WorkflowArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "workflow", data)
	if err != nil {
		return infer.CreateResponse[WorkflowState]{}, err
	}

	var state WorkflowState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[WorkflowState]{}, err
	}
	state.WorkflowArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[WorkflowState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (w *Workflow) Read(ctx context.Context, req infer.ReadRequest[WorkflowArgs, WorkflowState]) (infer.ReadResponse[WorkflowArgs, WorkflowState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "workflow", req.ID, SelectFields(WorkflowState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[WorkflowArgs, WorkflowState]{}, nil
		}
		return infer.ReadResponse[WorkflowArgs, WorkflowState]{}, err
	}

	var state WorkflowState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[WorkflowArgs, WorkflowState]{}, err
	}

	return infer.ReadResponse[WorkflowArgs, WorkflowState]{
		ID:     state.ResourceID,
		Inputs: state.WorkflowArgs,
		State:  state,
	}, nil
}

func (w *Workflow) Update(ctx context.Context, req infer.UpdateRequest[WorkflowArgs, WorkflowState]) (infer.UpdateResponse[WorkflowState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[WorkflowState]{
			Output: WorkflowState{WorkflowArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[WorkflowState]{}, err
	}

	if err := c.UpdateResource(ctx, "workflow", req.ID, data); err != nil {
		return infer.UpdateResponse[WorkflowState]{}, err
	}

	result, err := c.ReadResource(ctx, "workflow", req.ID, SelectFields(WorkflowState{}))
	if err != nil {
		return infer.UpdateResponse[WorkflowState]{}, err
	}

	var state WorkflowState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[WorkflowState]{}, err
	}

	return infer.UpdateResponse[WorkflowState]{
		Output: state,
	}, nil
}

func (w *Workflow) Delete(ctx context.Context, req infer.DeleteRequest[WorkflowState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "workflow", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
