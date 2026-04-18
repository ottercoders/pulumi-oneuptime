package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type TeamPermission struct{}

type TeamPermissionArgs struct {
	ProjectID         *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	TeamID            string   `pulumi:"teamId" json:"teamId"`
	Permission        string   `pulumi:"permission" json:"permission"`
	Labels            []string `pulumi:"labels,optional" json:"labels,omitempty"`
	IsBlockPermission *bool    `pulumi:"isBlockPermission,optional" json:"isBlockPermission,omitempty"`
}

type TeamPermissionState struct {
	TeamPermissionArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*TeamPermission)(nil)

func (p *TeamPermission) Annotate(a infer.Annotator) {
	a.Describe(p, "Manages a permission grant on a OneUptime Team.")
}

func (p *TeamPermission) Create(ctx context.Context, req infer.CreateRequest[TeamPermissionArgs]) (infer.CreateResponse[TeamPermissionState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[TeamPermissionState]{
			ID:     "preview-id",
			Output: TeamPermissionState{TeamPermissionArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[TeamPermissionState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[TeamPermissionState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "team-permission", data)
	if err != nil {
		return infer.CreateResponse[TeamPermissionState]{}, err
	}

	var state TeamPermissionState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[TeamPermissionState]{}, err
	}
	state.TeamPermissionArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[TeamPermissionState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (p *TeamPermission) Read(ctx context.Context, req infer.ReadRequest[TeamPermissionArgs, TeamPermissionState]) (infer.ReadResponse[TeamPermissionArgs, TeamPermissionState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "team-permission", req.ID, SelectFields(TeamPermissionState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[TeamPermissionArgs, TeamPermissionState]{}, nil
		}
		return infer.ReadResponse[TeamPermissionArgs, TeamPermissionState]{}, err
	}

	var state TeamPermissionState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[TeamPermissionArgs, TeamPermissionState]{}, err
	}

	return infer.ReadResponse[TeamPermissionArgs, TeamPermissionState]{
		ID:     state.ResourceID,
		Inputs: state.TeamPermissionArgs,
		State:  state,
	}, nil
}

func (p *TeamPermission) Update(ctx context.Context, req infer.UpdateRequest[TeamPermissionArgs, TeamPermissionState]) (infer.UpdateResponse[TeamPermissionState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[TeamPermissionState]{
			Output: TeamPermissionState{TeamPermissionArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[TeamPermissionState]{}, err
	}

	if err := c.UpdateResource(ctx, "team-permission", req.ID, data); err != nil {
		return infer.UpdateResponse[TeamPermissionState]{}, err
	}

	result, err := c.ReadResource(ctx, "team-permission", req.ID, SelectFields(TeamPermissionState{}))
	if err != nil {
		return infer.UpdateResponse[TeamPermissionState]{}, err
	}

	var state TeamPermissionState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[TeamPermissionState]{}, err
	}

	return infer.UpdateResponse[TeamPermissionState]{
		Output: state,
	}, nil
}

func (p *TeamPermission) Delete(ctx context.Context, req infer.DeleteRequest[TeamPermissionState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "team-permission", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
