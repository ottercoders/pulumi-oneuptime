package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Label struct{}

type LabelArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string  `pulumi:"name" json:"name"`
	Color       string  `pulumi:"color" json:"color"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
}

type LabelState struct {
	LabelArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Label)(nil)

func (l *Label) Annotate(a infer.Annotator) {
	a.Describe(l, "Manages a OneUptime Label resource.")
}

func (l *Label) Create(ctx context.Context, req infer.CreateRequest[LabelArgs]) (infer.CreateResponse[LabelState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[LabelState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[LabelState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[LabelState]{
			ID:     "preview-id",
			Output: LabelState{LabelArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "label", data)
	if err != nil {
		return infer.CreateResponse[LabelState]{}, err
	}

	var state LabelState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[LabelState]{}, err
	}
	state.LabelArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[LabelState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (l *Label) Read(ctx context.Context, req infer.ReadRequest[LabelArgs, LabelState]) (infer.ReadResponse[LabelArgs, LabelState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "label", req.ID, SelectFields(LabelState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[LabelArgs, LabelState]{}, nil
		}
		return infer.ReadResponse[LabelArgs, LabelState]{}, err
	}

	var state LabelState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[LabelArgs, LabelState]{}, err
	}

	return infer.ReadResponse[LabelArgs, LabelState]{
		ID:     state.ResourceID,
		Inputs: state.LabelArgs,
		State:  state,
	}, nil
}

func (l *Label) Update(ctx context.Context, req infer.UpdateRequest[LabelArgs, LabelState]) (infer.UpdateResponse[LabelState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[LabelState]{
			Output: LabelState{LabelArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[LabelState]{}, err
	}

	if err := c.UpdateResource(ctx, "label", req.ID, data); err != nil {
		return infer.UpdateResponse[LabelState]{}, err
	}

	result, err := c.ReadResource(ctx, "label", req.ID, SelectFields(LabelState{}))
	if err != nil {
		return infer.UpdateResponse[LabelState]{}, err
	}

	var state LabelState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[LabelState]{}, err
	}

	return infer.UpdateResponse[LabelState]{
		Output: state,
	}, nil
}

func (l *Label) Delete(ctx context.Context, req infer.DeleteRequest[LabelState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "label", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
