package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ProjectSmtpConfig struct{}

type ProjectSmtpConfigArgs struct {
	ProjectID         *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name              string  `pulumi:"name" json:"name"`
	Description       *string `pulumi:"description,optional" json:"description,omitempty"`
	Hostname          string  `pulumi:"hostname" json:"hostname"`
	Port              int     `pulumi:"port" json:"port"`
	Username          *string `pulumi:"username,optional" json:"username,omitempty"`
	Password          *string `pulumi:"password,optional" json:"password,omitempty"`
	FromEmail         string  `pulumi:"fromEmail" json:"fromEmail"`
	FromName          string  `pulumi:"fromName" json:"fromName"`
	Secure            *bool   `pulumi:"secure,optional" json:"secure,omitempty"`
	AuthType          *string `pulumi:"authType,optional" json:"authType,omitempty"`
	ClientID          *string `pulumi:"clientId,optional" json:"clientId,omitempty"`
	ClientSecret      *string `pulumi:"clientSecret,optional" json:"clientSecret,omitempty"`
	TokenURL          *string `pulumi:"tokenUrl,optional" json:"tokenUrl,omitempty"`
	Scope             *string `pulumi:"scope,optional" json:"scope,omitempty"`
	OauthProviderType *string `pulumi:"oauthProviderType,optional" json:"oauthProviderType,omitempty"`
}

type ProjectSmtpConfigState struct {
	ProjectSmtpConfigArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ProjectSmtpConfig)(nil)
var _ infer.ExplicitDependencies[ProjectSmtpConfigArgs, ProjectSmtpConfigState] = (*ProjectSmtpConfig)(nil)

func (s *ProjectSmtpConfig) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime project SMTP config used for outbound email. "+
		"authType accepts: \"Username and Password\", \"OAuth\", or \"None\". "+
		"password and clientSecret are secret outputs that the API does not return on read; "+
		"the provider preserves them from state.")
}

func (s *ProjectSmtpConfig) WireDependencies(f infer.FieldSelector, args *ProjectSmtpConfigArgs, state *ProjectSmtpConfigState) {
	f.OutputField(&state.Password).AlwaysSecret()
	f.OutputField(&state.ClientSecret).AlwaysSecret()
}

func (s *ProjectSmtpConfig) Create(ctx context.Context, req infer.CreateRequest[ProjectSmtpConfigArgs]) (infer.CreateResponse[ProjectSmtpConfigState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ProjectSmtpConfigState]{
			ID:     "preview-id",
			Output: ProjectSmtpConfigState{ProjectSmtpConfigArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ProjectSmtpConfigState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ProjectSmtpConfigState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "smtp-config", data)
	if err != nil {
		return infer.CreateResponse[ProjectSmtpConfigState]{}, err
	}

	var state ProjectSmtpConfigState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ProjectSmtpConfigState]{}, err
	}
	state.ProjectSmtpConfigArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ProjectSmtpConfigState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *ProjectSmtpConfig) Read(ctx context.Context, req infer.ReadRequest[ProjectSmtpConfigArgs, ProjectSmtpConfigState]) (infer.ReadResponse[ProjectSmtpConfigArgs, ProjectSmtpConfigState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "smtp-config", req.ID, SelectFields(ProjectSmtpConfigState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ProjectSmtpConfigArgs, ProjectSmtpConfigState]{}, nil
		}
		return infer.ReadResponse[ProjectSmtpConfigArgs, ProjectSmtpConfigState]{}, err
	}

	var state ProjectSmtpConfigState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ProjectSmtpConfigArgs, ProjectSmtpConfigState]{}, err
	}
	// Server's read ACL omits secrets; preserve them from prior state to avoid drift.
	state.Password = req.State.Password
	state.ClientSecret = req.State.ClientSecret

	return infer.ReadResponse[ProjectSmtpConfigArgs, ProjectSmtpConfigState]{
		ID:     state.ResourceID,
		Inputs: state.ProjectSmtpConfigArgs,
		State:  state,
	}, nil
}

func (s *ProjectSmtpConfig) Update(ctx context.Context, req infer.UpdateRequest[ProjectSmtpConfigArgs, ProjectSmtpConfigState]) (infer.UpdateResponse[ProjectSmtpConfigState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ProjectSmtpConfigState]{
			Output: ProjectSmtpConfigState{ProjectSmtpConfigArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ProjectSmtpConfigState]{}, err
	}

	if err := c.UpdateResource(ctx, "smtp-config", req.ID, data); err != nil {
		return infer.UpdateResponse[ProjectSmtpConfigState]{}, err
	}

	result, err := c.ReadResource(ctx, "smtp-config", req.ID, SelectFields(ProjectSmtpConfigState{}))
	if err != nil {
		return infer.UpdateResponse[ProjectSmtpConfigState]{}, err
	}

	var state ProjectSmtpConfigState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ProjectSmtpConfigState]{}, err
	}
	state.Password = req.Inputs.Password
	if state.Password == nil {
		state.Password = req.State.Password
	}
	state.ClientSecret = req.Inputs.ClientSecret
	if state.ClientSecret == nil {
		state.ClientSecret = req.State.ClientSecret
	}

	return infer.UpdateResponse[ProjectSmtpConfigState]{
		Output: state,
	}, nil
}

func (s *ProjectSmtpConfig) Delete(ctx context.Context, req infer.DeleteRequest[ProjectSmtpConfigState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "smtp-config", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
