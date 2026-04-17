package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type IncidentUserOwner struct{}

type IncidentUserOwnerArgs struct {
	ProjectID  *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	IncidentID string  `pulumi:"incidentId" json:"incidentId"`
	UserID     string  `pulumi:"userId" json:"userId"`
}

type IncidentUserOwnerState struct {
	IncidentUserOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*IncidentUserOwner)(nil)

func (r *IncidentUserOwner) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a user as owner of an incident.")
}

func (r *IncidentUserOwner) Create(ctx context.Context, req infer.CreateRequest[IncidentUserOwnerArgs]) (infer.CreateResponse[IncidentUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[IncidentUserOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[IncidentUserOwnerState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[IncidentUserOwnerState]{
			ID:     "preview-id",
			Output: IncidentUserOwnerState{IncidentUserOwnerArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "incident-owner-user", data)
	if err != nil {
		return infer.CreateResponse[IncidentUserOwnerState]{}, err
	}

	var state IncidentUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[IncidentUserOwnerState]{}, err
	}
	state.IncidentUserOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[IncidentUserOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *IncidentUserOwner) Read(ctx context.Context, req infer.ReadRequest[IncidentUserOwnerArgs, IncidentUserOwnerState]) (infer.ReadResponse[IncidentUserOwnerArgs, IncidentUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "incident-owner-user", req.ID, SelectFields(IncidentUserOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[IncidentUserOwnerArgs, IncidentUserOwnerState]{}, nil
		}
		return infer.ReadResponse[IncidentUserOwnerArgs, IncidentUserOwnerState]{}, err
	}

	var state IncidentUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[IncidentUserOwnerArgs, IncidentUserOwnerState]{}, err
	}

	return infer.ReadResponse[IncidentUserOwnerArgs, IncidentUserOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.IncidentUserOwnerArgs,
		State:  state,
	}, nil
}

func (r *IncidentUserOwner) Update(ctx context.Context, req infer.UpdateRequest[IncidentUserOwnerArgs, IncidentUserOwnerState]) (infer.UpdateResponse[IncidentUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[IncidentUserOwnerState]{
			Output: IncidentUserOwnerState{IncidentUserOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[IncidentUserOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "incident-owner-user", req.ID, data); err != nil {
		return infer.UpdateResponse[IncidentUserOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "incident-owner-user", req.ID, SelectFields(IncidentUserOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[IncidentUserOwnerState]{}, err
	}

	var state IncidentUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[IncidentUserOwnerState]{}, err
	}

	return infer.UpdateResponse[IncidentUserOwnerState]{
		Output: state,
	}, nil
}

func (r *IncidentUserOwner) Delete(ctx context.Context, req infer.DeleteRequest[IncidentUserOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "incident-owner-user", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
