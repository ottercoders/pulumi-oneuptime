package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ProjectSSO struct{}

type ProjectSSOArgs struct {
	ProjectID          *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name               string  `pulumi:"name" json:"name"`
	Description        *string `pulumi:"description,optional" json:"description,omitempty"`
	SignOnURL          *string `pulumi:"signOnURL,optional" json:"signOnURL,omitempty"`
	IssuerURL          *string `pulumi:"issuerURL,optional" json:"issuerURL,omitempty"`
	SsoType            *string `pulumi:"ssoType,optional" json:"ssoType,omitempty"`
	IsEnabled          *bool   `pulumi:"isEnabled,optional" json:"isEnabled,omitempty"`
	IsSaml             *bool   `pulumi:"isSaml,optional" json:"isSaml,omitempty"`
	OidcDiscoveryURL   *string `pulumi:"oidcDiscoveryUrl,optional" json:"oidcDiscoveryUrl,omitempty"`
	OidcClientID       *string `pulumi:"oidcClientId,optional" json:"oidcClientId,omitempty"`
	OidcClientSecret   *string `pulumi:"oidcClientSecret,optional" json:"oidcClientSecret,omitempty"`
}

type ProjectSSOState struct {
	ProjectSSOArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ProjectSSO)(nil)

func (s *ProjectSSO) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime Project SSO/OIDC configuration.")
}

func (s *ProjectSSO) Create(ctx context.Context, req infer.CreateRequest[ProjectSSOArgs]) (infer.CreateResponse[ProjectSSOState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ProjectSSOState]{
			ID:     "preview-id",
			Output: ProjectSSOState{ProjectSSOArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ProjectSSOState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ProjectSSOState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "project-sso", data)
	if err != nil {
		return infer.CreateResponse[ProjectSSOState]{}, err
	}

	var state ProjectSSOState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ProjectSSOState]{}, err
	}
	state.ProjectSSOArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ProjectSSOState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *ProjectSSO) Read(ctx context.Context, req infer.ReadRequest[ProjectSSOArgs, ProjectSSOState]) (infer.ReadResponse[ProjectSSOArgs, ProjectSSOState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "project-sso", req.ID, SelectFields(ProjectSSOState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ProjectSSOArgs, ProjectSSOState]{}, nil
		}
		return infer.ReadResponse[ProjectSSOArgs, ProjectSSOState]{}, err
	}

	var state ProjectSSOState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ProjectSSOArgs, ProjectSSOState]{}, err
	}

	return infer.ReadResponse[ProjectSSOArgs, ProjectSSOState]{
		ID:     state.ResourceID,
		Inputs: state.ProjectSSOArgs,
		State:  state,
	}, nil
}

func (s *ProjectSSO) Update(ctx context.Context, req infer.UpdateRequest[ProjectSSOArgs, ProjectSSOState]) (infer.UpdateResponse[ProjectSSOState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ProjectSSOState]{
			Output: ProjectSSOState{ProjectSSOArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ProjectSSOState]{}, err
	}

	if err := c.UpdateResource(ctx, "project-sso", req.ID, data); err != nil {
		return infer.UpdateResponse[ProjectSSOState]{}, err
	}

	result, err := c.ReadResource(ctx, "project-sso", req.ID, SelectFields(ProjectSSOState{}))
	if err != nil {
		return infer.UpdateResponse[ProjectSSOState]{}, err
	}

	var state ProjectSSOState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ProjectSSOState]{}, err
	}

	return infer.UpdateResponse[ProjectSSOState]{
		Output: state,
	}, nil
}

func (s *ProjectSSO) Delete(ctx context.Context, req infer.DeleteRequest[ProjectSSOState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "project-sso", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
