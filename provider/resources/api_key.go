package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ApiKey struct{}

type ApiKeyArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string  `pulumi:"name" json:"name"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
	ExpiresAt   string  `pulumi:"expiresAt" json:"expiresAt"`
}

type ApiKeyState struct {
	ApiKeyArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	ApiKey     string `pulumi:"apiKey,optional" json:"apiKey,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ApiKey)(nil)
var _ infer.ExplicitDependencies[ApiKeyArgs, ApiKeyState] = (*ApiKey)(nil)

func (k *ApiKey) Annotate(a infer.Annotator) {
	a.Describe(k, "Manages a OneUptime API Key. The apiKey output field is a secret and only available on create.")
}

// WireDependencies marks the generated apiKey value as always secret.
func (k *ApiKey) WireDependencies(f infer.FieldSelector, args *ApiKeyArgs, state *ApiKeyState) {
	f.OutputField(&state.ApiKey).AlwaysSecret()
}

func (k *ApiKey) Create(ctx context.Context, req infer.CreateRequest[ApiKeyArgs]) (infer.CreateResponse[ApiKeyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ApiKeyState]{
			ID:     "preview-id",
			Output: ApiKeyState{ApiKeyArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ApiKeyState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ApiKeyState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "api-key", data)
	if err != nil {
		return infer.CreateResponse[ApiKeyState]{}, err
	}

	var state ApiKeyState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ApiKeyState]{}, err
	}
	state.ApiKeyArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ApiKeyState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (k *ApiKey) Read(ctx context.Context, req infer.ReadRequest[ApiKeyArgs, ApiKeyState]) (infer.ReadResponse[ApiKeyArgs, ApiKeyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "api-key", req.ID, SelectFields(ApiKeyState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ApiKeyArgs, ApiKeyState]{}, nil
		}
		return infer.ReadResponse[ApiKeyArgs, ApiKeyState]{}, err
	}

	var state ApiKeyState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ApiKeyArgs, ApiKeyState]{}, err
	}
	// Preserve the secret apiKey from prior state (Read endpoint doesn't return it)
	state.ApiKey = req.State.ApiKey

	return infer.ReadResponse[ApiKeyArgs, ApiKeyState]{
		ID:     state.ResourceID,
		Inputs: state.ApiKeyArgs,
		State:  state,
	}, nil
}

func (k *ApiKey) Update(ctx context.Context, req infer.UpdateRequest[ApiKeyArgs, ApiKeyState]) (infer.UpdateResponse[ApiKeyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ApiKeyState]{
			Output: ApiKeyState{ApiKeyArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ApiKeyState]{}, err
	}

	if err := c.UpdateResource(ctx, "api-key", req.ID, data); err != nil {
		return infer.UpdateResponse[ApiKeyState]{}, err
	}

	result, err := c.ReadResource(ctx, "api-key", req.ID, SelectFields(ApiKeyState{}))
	if err != nil {
		return infer.UpdateResponse[ApiKeyState]{}, err
	}

	var state ApiKeyState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ApiKeyState]{}, err
	}
	// Preserve the secret apiKey from prior state (Read endpoint doesn't return it)
	state.ApiKey = req.State.ApiKey

	return infer.UpdateResponse[ApiKeyState]{
		Output: state,
	}, nil
}

func (k *ApiKey) Delete(ctx context.Context, req infer.DeleteRequest[ApiKeyState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "api-key", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
