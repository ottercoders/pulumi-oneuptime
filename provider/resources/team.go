package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Team struct{}

type TeamArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string  `pulumi:"name" json:"name"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
}

type TeamState struct {
	TeamArgs
	ID        string `pulumi:"id" json:"_id"`
	Slug      string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Team)(nil)

func (t *Team) Annotate(a infer.Annotator) {
	a.Describe(t, "Manages a OneUptime Team resource.")
}

func (t *Team) Create(ctx context.Context, req infer.CreateRequest[TeamArgs]) (infer.CreateResponse[TeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[TeamState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[TeamState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[TeamState]{
			ID:     "preview-id",
			Output: TeamState{TeamArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "team", data)
	if err != nil {
		return infer.CreateResponse[TeamState]{}, err
	}

	var state TeamState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[TeamState]{}, err
	}
	state.TeamArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[TeamState]{
		ID:     state.ID,
		Output: state,
	}, nil
}

func (t *Team) Read(ctx context.Context, req infer.ReadRequest[TeamArgs, TeamState]) (infer.ReadResponse[TeamArgs, TeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "team", req.ID, SelectFields(TeamState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[TeamArgs, TeamState]{}, nil
		}
		return infer.ReadResponse[TeamArgs, TeamState]{}, err
	}

	var state TeamState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[TeamArgs, TeamState]{}, err
	}

	return infer.ReadResponse[TeamArgs, TeamState]{
		ID:     state.ID,
		Inputs: state.TeamArgs,
		State:  state,
	}, nil
}

func (t *Team) Update(ctx context.Context, req infer.UpdateRequest[TeamArgs, TeamState]) (infer.UpdateResponse[TeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[TeamState]{
			Output: TeamState{TeamArgs: req.Inputs, ID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[TeamState]{}, err
	}

	if err := c.UpdateResource(ctx, "team", req.ID, data); err != nil {
		return infer.UpdateResponse[TeamState]{}, err
	}

	// Read back to get server-computed fields
	result, err := c.ReadResource(ctx, "team", req.ID, SelectFields(TeamState{}))
	if err != nil {
		return infer.UpdateResponse[TeamState]{}, err
	}

	var state TeamState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[TeamState]{}, err
	}

	return infer.UpdateResponse[TeamState]{
		Output: state,
	}, nil
}

func (t *Team) Delete(ctx context.Context, req infer.DeleteRequest[TeamState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
