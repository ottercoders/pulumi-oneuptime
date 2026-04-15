package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageGroup struct{}

type StatusPageGroupArgs struct {
	ProjectID    *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID string  `pulumi:"statusPageId" json:"statusPageId"`
	Name         string  `pulumi:"name" json:"name"`
	Description  *string `pulumi:"description,optional" json:"description,omitempty"`
	Order        *int    `pulumi:"order,optional" json:"order,omitempty"`
	IsExpandedByDefault *bool `pulumi:"isExpandedByDefault,optional" json:"isExpandedByDefault,omitempty"`
}

type StatusPageGroupState struct {
	StatusPageGroupArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageGroup)(nil)

func (s *StatusPageGroup) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime Status Page Group resource. Groups organize monitors on a status page.")
}

func (s *StatusPageGroup) Create(ctx context.Context, req infer.CreateRequest[StatusPageGroupArgs]) (infer.CreateResponse[StatusPageGroupState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageGroupState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageGroupState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[StatusPageGroupState]{
			ID:     "preview-id",
			Output: StatusPageGroupState{StatusPageGroupArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "status-page-group", data)
	if err != nil {
		return infer.CreateResponse[StatusPageGroupState]{}, err
	}

	var state StatusPageGroupState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageGroupState]{}, err
	}
	state.StatusPageGroupArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageGroupState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *StatusPageGroup) Read(ctx context.Context, req infer.ReadRequest[StatusPageGroupArgs, StatusPageGroupState]) (infer.ReadResponse[StatusPageGroupArgs, StatusPageGroupState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-group", req.ID, SelectFields(StatusPageGroupState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageGroupArgs, StatusPageGroupState]{}, nil
		}
		return infer.ReadResponse[StatusPageGroupArgs, StatusPageGroupState]{}, err
	}

	var state StatusPageGroupState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageGroupArgs, StatusPageGroupState]{}, err
	}

	return infer.ReadResponse[StatusPageGroupArgs, StatusPageGroupState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageGroupArgs,
		State:  state,
	}, nil
}

func (s *StatusPageGroup) Update(ctx context.Context, req infer.UpdateRequest[StatusPageGroupArgs, StatusPageGroupState]) (infer.UpdateResponse[StatusPageGroupState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageGroupState]{
			Output: StatusPageGroupState{StatusPageGroupArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageGroupState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-group", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageGroupState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-group", req.ID, SelectFields(StatusPageGroupState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageGroupState]{}, err
	}

	var state StatusPageGroupState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageGroupState]{}, err
	}

	return infer.UpdateResponse[StatusPageGroupState]{
		Output: state,
	}, nil
}

func (s *StatusPageGroup) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageGroupState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-group", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
