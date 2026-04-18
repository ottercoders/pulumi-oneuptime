package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicyTeamOwner struct{}

type OnCallDutyPolicyTeamOwnerArgs struct {
	ProjectID        *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyID string  `pulumi:"onCallDutyPolicyId" json:"onCallDutyPolicyId"`
	TeamID           string  `pulumi:"teamId" json:"teamId"`
}

type OnCallDutyPolicyTeamOwnerState struct {
	OnCallDutyPolicyTeamOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicyTeamOwner)(nil)

func (r *OnCallDutyPolicyTeamOwner) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a team as owner of an on-call duty policy.")
}

func (r *OnCallDutyPolicyTeamOwner) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyTeamOwnerArgs]) (infer.CreateResponse[OnCallDutyPolicyTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyTeamOwnerState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyTeamOwnerState{OnCallDutyPolicyTeamOwnerArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "on-call-duty-policy-owner-team", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}

	var state OnCallDutyPolicyTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}
	state.OnCallDutyPolicyTeamOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyTeamOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyTeamOwner) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyTeamOwnerArgs, OnCallDutyPolicyTeamOwnerState]) (infer.ReadResponse[OnCallDutyPolicyTeamOwnerArgs, OnCallDutyPolicyTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-owner-team", req.ID, SelectFields(OnCallDutyPolicyTeamOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyTeamOwnerArgs, OnCallDutyPolicyTeamOwnerState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyTeamOwnerArgs, OnCallDutyPolicyTeamOwnerState]{}, err
	}

	var state OnCallDutyPolicyTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyTeamOwnerArgs, OnCallDutyPolicyTeamOwnerState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyTeamOwnerArgs, OnCallDutyPolicyTeamOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyTeamOwnerArgs,
		State:  state,
	}, nil
}

func (r *OnCallDutyPolicyTeamOwner) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyTeamOwnerArgs, OnCallDutyPolicyTeamOwnerState]) (infer.UpdateResponse[OnCallDutyPolicyTeamOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyTeamOwnerState]{
			Output: OnCallDutyPolicyTeamOwnerState{OnCallDutyPolicyTeamOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy-owner-team", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-owner-team", req.ID, SelectFields(OnCallDutyPolicyTeamOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}

	var state OnCallDutyPolicyTeamOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyTeamOwnerState]{}, err
	}

	return infer.UpdateResponse[OnCallDutyPolicyTeamOwnerState]{
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyTeamOwner) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyTeamOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-owner-team", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
