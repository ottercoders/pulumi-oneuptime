package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Project struct{}

type ProjectArgs struct {
	Name        string  `pulumi:"name" json:"name"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
}

type ProjectState struct {
	ProjectArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Project)(nil)

func (p *Project) Annotate(a infer.Annotator) {
	a.Describe(p, "Manages a OneUptime Project resource.")
}

func (p *Project) Create(ctx context.Context, req infer.CreateRequest[ProjectArgs]) (infer.CreateResponse[ProjectState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ProjectState]{}, err
	}

	if req.DryRun {
		return infer.CreateResponse[ProjectState]{
			ID:     "preview-id",
			Output: ProjectState{ProjectArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "project", data)
	if err != nil {
		return infer.CreateResponse[ProjectState]{}, err
	}

	var state ProjectState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ProjectState]{}, err
	}
	state.ProjectArgs = req.Inputs

	return infer.CreateResponse[ProjectState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (p *Project) Read(ctx context.Context, req infer.ReadRequest[ProjectArgs, ProjectState]) (infer.ReadResponse[ProjectArgs, ProjectState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "project", req.ID, SelectFields(ProjectState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ProjectArgs, ProjectState]{}, nil
		}
		return infer.ReadResponse[ProjectArgs, ProjectState]{}, err
	}

	var state ProjectState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ProjectArgs, ProjectState]{}, err
	}

	return infer.ReadResponse[ProjectArgs, ProjectState]{
		ID:     state.ResourceID,
		Inputs: state.ProjectArgs,
		State:  state,
	}, nil
}

func (p *Project) Update(ctx context.Context, req infer.UpdateRequest[ProjectArgs, ProjectState]) (infer.UpdateResponse[ProjectState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ProjectState]{
			Output: ProjectState{ProjectArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ProjectState]{}, err
	}

	if err := c.UpdateResource(ctx, "project", req.ID, data); err != nil {
		return infer.UpdateResponse[ProjectState]{}, err
	}

	result, err := c.ReadResource(ctx, "project", req.ID, SelectFields(ProjectState{}))
	if err != nil {
		return infer.UpdateResponse[ProjectState]{}, err
	}

	var state ProjectState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ProjectState]{}, err
	}

	return infer.UpdateResponse[ProjectState]{
		Output: state,
	}, nil
}

func (p *Project) Delete(ctx context.Context, req infer.DeleteRequest[ProjectState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "project", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
