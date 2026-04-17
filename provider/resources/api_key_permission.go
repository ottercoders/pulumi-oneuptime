package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ApiKeyPermission struct{}

type ApiKeyPermissionArgs struct {
	ProjectID         *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	ApiKeyID          string   `pulumi:"apiKeyId" json:"apiKeyId"`
	Permission        string   `pulumi:"permission" json:"permission"`
	Labels            []string `pulumi:"labels,optional" json:"labels,omitempty"`
	IsBlockPermission *bool    `pulumi:"isBlockPermission,optional" json:"isBlockPermission,omitempty"`
}

type ApiKeyPermissionState struct {
	ApiKeyPermissionArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ApiKeyPermission)(nil)

func (p *ApiKeyPermission) Annotate(a infer.Annotator) {
	a.Describe(p, "Manages a permission grant on a OneUptime API Key.")
}

func (p *ApiKeyPermission) Create(ctx context.Context, req infer.CreateRequest[ApiKeyPermissionArgs]) (infer.CreateResponse[ApiKeyPermissionState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ApiKeyPermissionState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ApiKeyPermissionState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[ApiKeyPermissionState]{
			ID:     "preview-id",
			Output: ApiKeyPermissionState{ApiKeyPermissionArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "api-key-permission", data)
	if err != nil {
		return infer.CreateResponse[ApiKeyPermissionState]{}, err
	}

	var state ApiKeyPermissionState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ApiKeyPermissionState]{}, err
	}
	state.ApiKeyPermissionArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ApiKeyPermissionState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (p *ApiKeyPermission) Read(ctx context.Context, req infer.ReadRequest[ApiKeyPermissionArgs, ApiKeyPermissionState]) (infer.ReadResponse[ApiKeyPermissionArgs, ApiKeyPermissionState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "api-key-permission", req.ID, SelectFields(ApiKeyPermissionState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ApiKeyPermissionArgs, ApiKeyPermissionState]{}, nil
		}
		return infer.ReadResponse[ApiKeyPermissionArgs, ApiKeyPermissionState]{}, err
	}

	var state ApiKeyPermissionState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ApiKeyPermissionArgs, ApiKeyPermissionState]{}, err
	}

	return infer.ReadResponse[ApiKeyPermissionArgs, ApiKeyPermissionState]{
		ID:     state.ResourceID,
		Inputs: state.ApiKeyPermissionArgs,
		State:  state,
	}, nil
}

func (p *ApiKeyPermission) Update(ctx context.Context, req infer.UpdateRequest[ApiKeyPermissionArgs, ApiKeyPermissionState]) (infer.UpdateResponse[ApiKeyPermissionState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ApiKeyPermissionState]{
			Output: ApiKeyPermissionState{ApiKeyPermissionArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ApiKeyPermissionState]{}, err
	}

	if err := c.UpdateResource(ctx, "api-key-permission", req.ID, data); err != nil {
		return infer.UpdateResponse[ApiKeyPermissionState]{}, err
	}

	result, err := c.ReadResource(ctx, "api-key-permission", req.ID, SelectFields(ApiKeyPermissionState{}))
	if err != nil {
		return infer.UpdateResponse[ApiKeyPermissionState]{}, err
	}

	var state ApiKeyPermissionState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ApiKeyPermissionState]{}, err
	}

	return infer.UpdateResponse[ApiKeyPermissionState]{
		Output: state,
	}, nil
}

func (p *ApiKeyPermission) Delete(ctx context.Context, req infer.DeleteRequest[ApiKeyPermissionState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "api-key-permission", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
