package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageDomain struct{}

type StatusPageDomainArgs struct {
	ProjectID    *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID string  `pulumi:"statusPageId" json:"statusPageId"`
	DomainID     string  `pulumi:"domainId" json:"domainId"`
	Subdomain    *string `pulumi:"subdomain,optional" json:"subdomain,omitempty"`
	FullDomain   *string `pulumi:"fullDomain,optional" json:"fullDomain,omitempty"`
	CnameVerificationToken *string `pulumi:"cnameVerificationToken,optional" json:"cnameVerificationToken,omitempty"`
	IsCnameVerified *bool `pulumi:"isCnameVerified,optional" json:"isCnameVerified,omitempty"`
	IsSslOrdered  *bool   `pulumi:"isSslOrdered,optional" json:"isSslOrdered,omitempty"`
}

type StatusPageDomainState struct {
	StatusPageDomainArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageDomain)(nil)

func (d *StatusPageDomain) Annotate(a infer.Annotator) {
	a.Describe(d, "Binds a custom Domain to a OneUptime Status Page.")
}

func (d *StatusPageDomain) Create(ctx context.Context, req infer.CreateRequest[StatusPageDomainArgs]) (infer.CreateResponse[StatusPageDomainState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageDomainState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageDomainState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[StatusPageDomainState]{
			ID:     "preview-id",
			Output: StatusPageDomainState{StatusPageDomainArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "status-page-domain", data)
	if err != nil {
		return infer.CreateResponse[StatusPageDomainState]{}, err
	}

	var state StatusPageDomainState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageDomainState]{}, err
	}
	state.StatusPageDomainArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageDomainState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (d *StatusPageDomain) Read(ctx context.Context, req infer.ReadRequest[StatusPageDomainArgs, StatusPageDomainState]) (infer.ReadResponse[StatusPageDomainArgs, StatusPageDomainState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-domain", req.ID, SelectFields(StatusPageDomainState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageDomainArgs, StatusPageDomainState]{}, nil
		}
		return infer.ReadResponse[StatusPageDomainArgs, StatusPageDomainState]{}, err
	}

	var state StatusPageDomainState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageDomainArgs, StatusPageDomainState]{}, err
	}

	return infer.ReadResponse[StatusPageDomainArgs, StatusPageDomainState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageDomainArgs,
		State:  state,
	}, nil
}

func (d *StatusPageDomain) Update(ctx context.Context, req infer.UpdateRequest[StatusPageDomainArgs, StatusPageDomainState]) (infer.UpdateResponse[StatusPageDomainState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageDomainState]{
			Output: StatusPageDomainState{StatusPageDomainArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageDomainState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-domain", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageDomainState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-domain", req.ID, SelectFields(StatusPageDomainState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageDomainState]{}, err
	}

	var state StatusPageDomainState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageDomainState]{}, err
	}

	return infer.UpdateResponse[StatusPageDomainState]{
		Output: state,
	}, nil
}

func (d *StatusPageDomain) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageDomainState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-domain", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
