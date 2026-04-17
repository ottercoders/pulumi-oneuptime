package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type IncidentTeamOwner struct{}

type IncidentTeamOwnerArgs struct {
	ProjectID  *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	IncidentID string  `pulumi:"incidentId" json:"incidentId"`
	TeamID     string  `pulumi:"teamId" json:"teamId"`
}

type IncidentTeamOwnerState struct {
	IncidentTeamOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*IncidentTeamOwner)(nil)

func (r *IncidentTeamOwner) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a team as owner of an incident.")
}

func (r *IncidentTeamOwner) Create(ctx context.Context, req infer.CreateRequest[IncidentTeamOwnerArgs]) (infer.CreateResponse[IncidentTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[IncidentTeamOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[IncidentTeamOwnerState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[IncidentTeamOwnerState]{
			ID:     "preview-id",
			Output: IncidentTeamOwnerState{IncidentTeamOwnerArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "incident-owner-team", data)
	if err != nil {
		return infer.CreateResponse[IncidentTeamOwnerState]{}, err
	}

	var state IncidentTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[IncidentTeamOwnerState]{}, err
	}
	state.IncidentTeamOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[IncidentTeamOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *IncidentTeamOwner) Read(ctx context.Context, req infer.ReadRequest[IncidentTeamOwnerArgs, IncidentTeamOwnerState]) (infer.ReadResponse[IncidentTeamOwnerArgs, IncidentTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "incident-owner-team", req.ID, SelectFields(IncidentTeamOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[IncidentTeamOwnerArgs, IncidentTeamOwnerState]{}, nil
		}
		return infer.ReadResponse[IncidentTeamOwnerArgs, IncidentTeamOwnerState]{}, err
	}

	var state IncidentTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[IncidentTeamOwnerArgs, IncidentTeamOwnerState]{}, err
	}

	return infer.ReadResponse[IncidentTeamOwnerArgs, IncidentTeamOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.IncidentTeamOwnerArgs,
		State:  state,
	}, nil
}

func (r *IncidentTeamOwner) Update(ctx context.Context, req infer.UpdateRequest[IncidentTeamOwnerArgs, IncidentTeamOwnerState]) (infer.UpdateResponse[IncidentTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[IncidentTeamOwnerState]{
			Output: IncidentTeamOwnerState{IncidentTeamOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[IncidentTeamOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "incident-owner-team", req.ID, data); err != nil {
		return infer.UpdateResponse[IncidentTeamOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "incident-owner-team", req.ID, SelectFields(IncidentTeamOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[IncidentTeamOwnerState]{}, err
	}

	var state IncidentTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[IncidentTeamOwnerState]{}, err
	}

	return infer.UpdateResponse[IncidentTeamOwnerState]{
		Output: state,
	}, nil
}

func (r *IncidentTeamOwner) Delete(ctx context.Context, req infer.DeleteRequest[IncidentTeamOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "incident-owner-team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
