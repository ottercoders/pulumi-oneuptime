package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type WorkflowVariable struct{}

type WorkflowVariableArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	WorkflowID  string  `pulumi:"workflowId" json:"workflowId"`
	Name        string  `pulumi:"name" json:"name"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
	Content     *string `pulumi:"content,optional" json:"content,omitempty"`
	IsSecret    *bool   `pulumi:"isSecret,optional" json:"isSecret,omitempty"`
}

type WorkflowVariableState struct {
	WorkflowVariableArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*WorkflowVariable)(nil)

func (r *WorkflowVariable) Annotate(a infer.Annotator) {
	a.Describe(r, "Manages a variable within a OneUptime Workflow.")
}

func (r *WorkflowVariable) Create(ctx context.Context, req infer.CreateRequest[WorkflowVariableArgs]) (infer.CreateResponse[WorkflowVariableState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[WorkflowVariableState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[WorkflowVariableState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[WorkflowVariableState]{
			ID:     "preview-id",
			Output: WorkflowVariableState{WorkflowVariableArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "workflow-variable", data)
	if err != nil {
		return infer.CreateResponse[WorkflowVariableState]{}, err
	}

	var state WorkflowVariableState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[WorkflowVariableState]{}, err
	}
	state.WorkflowVariableArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[WorkflowVariableState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *WorkflowVariable) Read(ctx context.Context, req infer.ReadRequest[WorkflowVariableArgs, WorkflowVariableState]) (infer.ReadResponse[WorkflowVariableArgs, WorkflowVariableState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "workflow-variable", req.ID, SelectFields(WorkflowVariableState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[WorkflowVariableArgs, WorkflowVariableState]{}, nil
		}
		return infer.ReadResponse[WorkflowVariableArgs, WorkflowVariableState]{}, err
	}

	var state WorkflowVariableState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[WorkflowVariableArgs, WorkflowVariableState]{}, err
	}

	return infer.ReadResponse[WorkflowVariableArgs, WorkflowVariableState]{
		ID:     state.ResourceID,
		Inputs: state.WorkflowVariableArgs,
		State:  state,
	}, nil
}

func (r *WorkflowVariable) Update(ctx context.Context, req infer.UpdateRequest[WorkflowVariableArgs, WorkflowVariableState]) (infer.UpdateResponse[WorkflowVariableState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[WorkflowVariableState]{
			Output: WorkflowVariableState{WorkflowVariableArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[WorkflowVariableState]{}, err
	}

	if err := c.UpdateResource(ctx, "workflow-variable", req.ID, data); err != nil {
		return infer.UpdateResponse[WorkflowVariableState]{}, err
	}

	result, err := c.ReadResource(ctx, "workflow-variable", req.ID, SelectFields(WorkflowVariableState{}))
	if err != nil {
		return infer.UpdateResponse[WorkflowVariableState]{}, err
	}

	var state WorkflowVariableState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[WorkflowVariableState]{}, err
	}

	return infer.UpdateResponse[WorkflowVariableState]{
		Output: state,
	}, nil
}

func (r *WorkflowVariable) Delete(ctx context.Context, req infer.DeleteRequest[WorkflowVariableState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "workflow-variable", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
