package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ProjectSsoTeam struct{}

type ProjectSsoTeamArgs struct {
	ProjectID    *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	ProjectSsoID string  `pulumi:"projectSsoId" json:"projectSsoId"`
	TeamID       string  `pulumi:"teamId" json:"teamId"`
}

type ProjectSsoTeamState struct {
	ProjectSsoTeamArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ProjectSsoTeam)(nil)

func (t *ProjectSsoTeam) Annotate(a infer.Annotator) {
	a.Describe(t, "Links a OneUptime SSO configuration to a Team. Users logging in via this SSO will be auto-assigned to the linked team.")
}

func (t *ProjectSsoTeam) Create(ctx context.Context, req infer.CreateRequest[ProjectSsoTeamArgs]) (infer.CreateResponse[ProjectSsoTeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ProjectSsoTeamState]{
			ID:     "preview-id",
			Output: ProjectSsoTeamState{ProjectSsoTeamArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ProjectSsoTeamState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ProjectSsoTeamState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "project-sso-team", data)
	if err != nil {
		return infer.CreateResponse[ProjectSsoTeamState]{}, err
	}

	var state ProjectSsoTeamState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ProjectSsoTeamState]{}, err
	}
	state.ProjectSsoTeamArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ProjectSsoTeamState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (t *ProjectSsoTeam) Read(ctx context.Context, req infer.ReadRequest[ProjectSsoTeamArgs, ProjectSsoTeamState]) (infer.ReadResponse[ProjectSsoTeamArgs, ProjectSsoTeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "project-sso-team", req.ID, SelectFields(ProjectSsoTeamState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ProjectSsoTeamArgs, ProjectSsoTeamState]{}, nil
		}
		return infer.ReadResponse[ProjectSsoTeamArgs, ProjectSsoTeamState]{}, err
	}

	var state ProjectSsoTeamState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ProjectSsoTeamArgs, ProjectSsoTeamState]{}, err
	}

	return infer.ReadResponse[ProjectSsoTeamArgs, ProjectSsoTeamState]{
		ID:     state.ResourceID,
		Inputs: state.ProjectSsoTeamArgs,
		State:  state,
	}, nil
}

func (t *ProjectSsoTeam) Update(ctx context.Context, req infer.UpdateRequest[ProjectSsoTeamArgs, ProjectSsoTeamState]) (infer.UpdateResponse[ProjectSsoTeamState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ProjectSsoTeamState]{
			Output: ProjectSsoTeamState{ProjectSsoTeamArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ProjectSsoTeamState]{}, err
	}

	if err := c.UpdateResource(ctx, "project-sso-team", req.ID, data); err != nil {
		return infer.UpdateResponse[ProjectSsoTeamState]{}, err
	}

	result, err := c.ReadResource(ctx, "project-sso-team", req.ID, SelectFields(ProjectSsoTeamState{}))
	if err != nil {
		return infer.UpdateResponse[ProjectSsoTeamState]{}, err
	}

	var state ProjectSsoTeamState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ProjectSsoTeamState]{}, err
	}

	return infer.UpdateResponse[ProjectSsoTeamState]{
		Output: state,
	}, nil
}

func (t *ProjectSsoTeam) Delete(ctx context.Context, req infer.DeleteRequest[ProjectSsoTeamState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "project-sso-team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
