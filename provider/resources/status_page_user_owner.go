package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageUserOwner struct{}

type StatusPageUserOwnerArgs struct {
	ProjectID    *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID string  `pulumi:"statusPageId" json:"statusPageId"`
	UserID       string  `pulumi:"userId" json:"userId"`
}

type StatusPageUserOwnerState struct {
	StatusPageUserOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageUserOwner)(nil)

func (o *StatusPageUserOwner) Annotate(a infer.Annotator) {
	a.Describe(o, "Assigns a OneUptime User as owner of a Status Page.")
}

func (o *StatusPageUserOwner) Create(ctx context.Context, req infer.CreateRequest[StatusPageUserOwnerArgs]) (infer.CreateResponse[StatusPageUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[StatusPageUserOwnerState]{
			ID:     "preview-id",
			Output: StatusPageUserOwnerState{StatusPageUserOwnerArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageUserOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageUserOwnerState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "status-page-owner-user", data)
	if err != nil {
		return infer.CreateResponse[StatusPageUserOwnerState]{}, err
	}

	var state StatusPageUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageUserOwnerState]{}, err
	}
	state.StatusPageUserOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageUserOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (o *StatusPageUserOwner) Read(ctx context.Context, req infer.ReadRequest[StatusPageUserOwnerArgs, StatusPageUserOwnerState]) (infer.ReadResponse[StatusPageUserOwnerArgs, StatusPageUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-owner-user", req.ID, SelectFields(StatusPageUserOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageUserOwnerArgs, StatusPageUserOwnerState]{}, nil
		}
		return infer.ReadResponse[StatusPageUserOwnerArgs, StatusPageUserOwnerState]{}, err
	}

	var state StatusPageUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageUserOwnerArgs, StatusPageUserOwnerState]{}, err
	}

	return infer.ReadResponse[StatusPageUserOwnerArgs, StatusPageUserOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageUserOwnerArgs,
		State:  state,
	}, nil
}

func (o *StatusPageUserOwner) Update(ctx context.Context, req infer.UpdateRequest[StatusPageUserOwnerArgs, StatusPageUserOwnerState]) (infer.UpdateResponse[StatusPageUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageUserOwnerState]{
			Output: StatusPageUserOwnerState{StatusPageUserOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageUserOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-owner-user", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageUserOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-owner-user", req.ID, SelectFields(StatusPageUserOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageUserOwnerState]{}, err
	}

	var state StatusPageUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageUserOwnerState]{}, err
	}

	return infer.UpdateResponse[StatusPageUserOwnerState]{
		Output: state,
	}, nil
}

func (o *StatusPageUserOwner) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageUserOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-owner-user", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
