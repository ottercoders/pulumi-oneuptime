package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageFooterLink struct{}

type StatusPageFooterLinkArgs struct {
	ProjectID    *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID string  `pulumi:"statusPageId" json:"statusPageId"`
	Title        string  `pulumi:"title" json:"title"`
	Link         string  `pulumi:"link" json:"link"`
	Order        *int    `pulumi:"order,optional" json:"order,omitempty"`
	OpenInNewTab *bool   `pulumi:"openInNewTab,optional" json:"openInNewTab,omitempty"`
}

type StatusPageFooterLinkState struct {
	StatusPageFooterLinkArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageFooterLink)(nil)

func (l *StatusPageFooterLink) Annotate(a infer.Annotator) {
	a.Describe(l, "Adds a link to the footer of a OneUptime Status Page.")
}

func (l *StatusPageFooterLink) Create(ctx context.Context, req infer.CreateRequest[StatusPageFooterLinkArgs]) (infer.CreateResponse[StatusPageFooterLinkState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[StatusPageFooterLinkState]{
			ID:     "preview-id",
			Output: StatusPageFooterLinkState{StatusPageFooterLinkArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageFooterLinkState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageFooterLinkState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "status-page-footer-link", data)
	if err != nil {
		return infer.CreateResponse[StatusPageFooterLinkState]{}, err
	}

	var state StatusPageFooterLinkState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageFooterLinkState]{}, err
	}
	state.StatusPageFooterLinkArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageFooterLinkState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (l *StatusPageFooterLink) Read(ctx context.Context, req infer.ReadRequest[StatusPageFooterLinkArgs, StatusPageFooterLinkState]) (infer.ReadResponse[StatusPageFooterLinkArgs, StatusPageFooterLinkState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-footer-link", req.ID, SelectFields(StatusPageFooterLinkState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageFooterLinkArgs, StatusPageFooterLinkState]{}, nil
		}
		return infer.ReadResponse[StatusPageFooterLinkArgs, StatusPageFooterLinkState]{}, err
	}

	var state StatusPageFooterLinkState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageFooterLinkArgs, StatusPageFooterLinkState]{}, err
	}

	return infer.ReadResponse[StatusPageFooterLinkArgs, StatusPageFooterLinkState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageFooterLinkArgs,
		State:  state,
	}, nil
}

func (l *StatusPageFooterLink) Update(ctx context.Context, req infer.UpdateRequest[StatusPageFooterLinkArgs, StatusPageFooterLinkState]) (infer.UpdateResponse[StatusPageFooterLinkState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageFooterLinkState]{
			Output: StatusPageFooterLinkState{StatusPageFooterLinkArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageFooterLinkState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-footer-link", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageFooterLinkState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-footer-link", req.ID, SelectFields(StatusPageFooterLinkState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageFooterLinkState]{}, err
	}

	var state StatusPageFooterLinkState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageFooterLinkState]{}, err
	}

	return infer.UpdateResponse[StatusPageFooterLinkState]{
		Output: state,
	}, nil
}

func (l *StatusPageFooterLink) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageFooterLinkState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-footer-link", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
