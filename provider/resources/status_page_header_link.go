package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageHeaderLink struct{}

type StatusPageHeaderLinkArgs struct {
	ProjectID    *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID string  `pulumi:"statusPageId" json:"statusPageId"`
	Title        string  `pulumi:"title" json:"title"`
	Link         string  `pulumi:"link" json:"link"`
	Order        *int    `pulumi:"order,optional" json:"order,omitempty"`
	OpenInNewTab *bool   `pulumi:"openInNewTab,optional" json:"openInNewTab,omitempty"`
}

type StatusPageHeaderLinkState struct {
	StatusPageHeaderLinkArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageHeaderLink)(nil)

func (l *StatusPageHeaderLink) Annotate(a infer.Annotator) {
	a.Describe(l, "Adds a link to the header of a OneUptime Status Page.")
}

func (l *StatusPageHeaderLink) Create(ctx context.Context, req infer.CreateRequest[StatusPageHeaderLinkArgs]) (infer.CreateResponse[StatusPageHeaderLinkState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[StatusPageHeaderLinkState]{
			ID:     "preview-id",
			Output: StatusPageHeaderLinkState{StatusPageHeaderLinkArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageHeaderLinkState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageHeaderLinkState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "status-page-header-link", data)
	if err != nil {
		return infer.CreateResponse[StatusPageHeaderLinkState]{}, err
	}

	var state StatusPageHeaderLinkState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageHeaderLinkState]{}, err
	}
	state.StatusPageHeaderLinkArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageHeaderLinkState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (l *StatusPageHeaderLink) Read(ctx context.Context, req infer.ReadRequest[StatusPageHeaderLinkArgs, StatusPageHeaderLinkState]) (infer.ReadResponse[StatusPageHeaderLinkArgs, StatusPageHeaderLinkState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-header-link", req.ID, SelectFields(StatusPageHeaderLinkState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageHeaderLinkArgs, StatusPageHeaderLinkState]{}, nil
		}
		return infer.ReadResponse[StatusPageHeaderLinkArgs, StatusPageHeaderLinkState]{}, err
	}

	var state StatusPageHeaderLinkState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageHeaderLinkArgs, StatusPageHeaderLinkState]{}, err
	}

	return infer.ReadResponse[StatusPageHeaderLinkArgs, StatusPageHeaderLinkState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageHeaderLinkArgs,
		State:  state,
	}, nil
}

func (l *StatusPageHeaderLink) Update(ctx context.Context, req infer.UpdateRequest[StatusPageHeaderLinkArgs, StatusPageHeaderLinkState]) (infer.UpdateResponse[StatusPageHeaderLinkState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageHeaderLinkState]{
			Output: StatusPageHeaderLinkState{StatusPageHeaderLinkArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageHeaderLinkState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-header-link", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageHeaderLinkState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-header-link", req.ID, SelectFields(StatusPageHeaderLinkState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageHeaderLinkState]{}, err
	}

	var state StatusPageHeaderLinkState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageHeaderLinkState]{}, err
	}

	return infer.UpdateResponse[StatusPageHeaderLinkState]{
		Output: state,
	}, nil
}

func (l *StatusPageHeaderLink) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageHeaderLinkState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-header-link", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
