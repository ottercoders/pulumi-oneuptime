package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageTeamOwner struct{}

type StatusPageTeamOwnerArgs struct {
	ProjectID    *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID string  `pulumi:"statusPageId" json:"statusPageId"`
	TeamID       string  `pulumi:"teamId" json:"teamId"`
}

type StatusPageTeamOwnerState struct {
	StatusPageTeamOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageTeamOwner)(nil)

func (o *StatusPageTeamOwner) Annotate(a infer.Annotator) {
	a.Describe(o, "Assigns a OneUptime Team as owner of a Status Page.")
}

func (o *StatusPageTeamOwner) Create(ctx context.Context, req infer.CreateRequest[StatusPageTeamOwnerArgs]) (infer.CreateResponse[StatusPageTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageTeamOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageTeamOwnerState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[StatusPageTeamOwnerState]{
			ID:     "preview-id",
			Output: StatusPageTeamOwnerState{StatusPageTeamOwnerArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "status-page-owner-team", data)
	if err != nil {
		return infer.CreateResponse[StatusPageTeamOwnerState]{}, err
	}

	var state StatusPageTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageTeamOwnerState]{}, err
	}
	state.StatusPageTeamOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageTeamOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (o *StatusPageTeamOwner) Read(ctx context.Context, req infer.ReadRequest[StatusPageTeamOwnerArgs, StatusPageTeamOwnerState]) (infer.ReadResponse[StatusPageTeamOwnerArgs, StatusPageTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-owner-team", req.ID, SelectFields(StatusPageTeamOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageTeamOwnerArgs, StatusPageTeamOwnerState]{}, nil
		}
		return infer.ReadResponse[StatusPageTeamOwnerArgs, StatusPageTeamOwnerState]{}, err
	}

	var state StatusPageTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageTeamOwnerArgs, StatusPageTeamOwnerState]{}, err
	}

	return infer.ReadResponse[StatusPageTeamOwnerArgs, StatusPageTeamOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageTeamOwnerArgs,
		State:  state,
	}, nil
}

func (o *StatusPageTeamOwner) Update(ctx context.Context, req infer.UpdateRequest[StatusPageTeamOwnerArgs, StatusPageTeamOwnerState]) (infer.UpdateResponse[StatusPageTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageTeamOwnerState]{
			Output: StatusPageTeamOwnerState{StatusPageTeamOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageTeamOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-owner-team", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageTeamOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-owner-team", req.ID, SelectFields(StatusPageTeamOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageTeamOwnerState]{}, err
	}

	var state StatusPageTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageTeamOwnerState]{}, err
	}

	return infer.UpdateResponse[StatusPageTeamOwnerState]{
		Output: state,
	}, nil
}

func (o *StatusPageTeamOwner) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageTeamOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-owner-team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
