package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Domain struct{}

type DomainArgs struct {
	ProjectID  *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Domain     string  `pulumi:"domain" json:"domain"`
	IsVerified *bool   `pulumi:"isVerified,optional" json:"isVerified,omitempty"`
}

type DomainState struct {
	DomainArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Domain)(nil)

func (d *Domain) Annotate(a infer.Annotator) {
	a.Describe(d, "Manages a OneUptime Domain resource. Required for SSO to match users by email domain.")
}

func (d *Domain) Create(ctx context.Context, req infer.CreateRequest[DomainArgs]) (infer.CreateResponse[DomainState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[DomainState]{
			ID:     "preview-id",
			Output: DomainState{DomainArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[DomainState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[DomainState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "domain", data)
	if err != nil {
		return infer.CreateResponse[DomainState]{}, err
	}

	var state DomainState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[DomainState]{}, err
	}
	state.DomainArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[DomainState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (d *Domain) Read(ctx context.Context, req infer.ReadRequest[DomainArgs, DomainState]) (infer.ReadResponse[DomainArgs, DomainState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "domain", req.ID, SelectFields(DomainState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[DomainArgs, DomainState]{}, nil
		}
		return infer.ReadResponse[DomainArgs, DomainState]{}, err
	}

	var state DomainState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[DomainArgs, DomainState]{}, err
	}

	return infer.ReadResponse[DomainArgs, DomainState]{
		ID:     state.ResourceID,
		Inputs: state.DomainArgs,
		State:  state,
	}, nil
}

func (d *Domain) Update(ctx context.Context, req infer.UpdateRequest[DomainArgs, DomainState]) (infer.UpdateResponse[DomainState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[DomainState]{
			Output: DomainState{DomainArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[DomainState]{}, err
	}

	if err := c.UpdateResource(ctx, "domain", req.ID, data); err != nil {
		return infer.UpdateResponse[DomainState]{}, err
	}

	result, err := c.ReadResource(ctx, "domain", req.ID, SelectFields(DomainState{}))
	if err != nil {
		return infer.UpdateResponse[DomainState]{}, err
	}

	var state DomainState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[DomainState]{}, err
	}

	return infer.UpdateResponse[DomainState]{
		Output: state,
	}, nil
}

func (d *Domain) Delete(ctx context.Context, req infer.DeleteRequest[DomainState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "domain", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
