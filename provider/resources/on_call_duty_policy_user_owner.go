package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type OnCallDutyPolicyUserOwner struct{}

type OnCallDutyPolicyUserOwnerArgs struct {
	ProjectID        *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	OnCallDutyPolicyID string  `pulumi:"onCallDutyPolicyId" json:"onCallDutyPolicyId"`
	UserID           string  `pulumi:"userId" json:"userId"`
}

type OnCallDutyPolicyUserOwnerState struct {
	OnCallDutyPolicyUserOwnerArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*OnCallDutyPolicyUserOwner)(nil)

func (r *OnCallDutyPolicyUserOwner) Annotate(a infer.Annotator) {
	a.Describe(r, "Assigns a user as owner of an on-call duty policy.")
}

func (r *OnCallDutyPolicyUserOwner) Create(ctx context.Context, req infer.CreateRequest[OnCallDutyPolicyUserOwnerArgs]) (infer.CreateResponse[OnCallDutyPolicyUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[OnCallDutyPolicyUserOwnerState]{
			ID:     "preview-id",
			Output: OnCallDutyPolicyUserOwnerState{OnCallDutyPolicyUserOwnerArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "on-call-duty-policy-owner-user", data)
	if err != nil {
		return infer.CreateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}

	var state OnCallDutyPolicyUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}
	state.OnCallDutyPolicyUserOwnerArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[OnCallDutyPolicyUserOwnerState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyUserOwner) Read(ctx context.Context, req infer.ReadRequest[OnCallDutyPolicyUserOwnerArgs, OnCallDutyPolicyUserOwnerState]) (infer.ReadResponse[OnCallDutyPolicyUserOwnerArgs, OnCallDutyPolicyUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "on-call-duty-policy-owner-user", req.ID, SelectFields(OnCallDutyPolicyUserOwnerState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[OnCallDutyPolicyUserOwnerArgs, OnCallDutyPolicyUserOwnerState]{}, nil
		}
		return infer.ReadResponse[OnCallDutyPolicyUserOwnerArgs, OnCallDutyPolicyUserOwnerState]{}, err
	}

	var state OnCallDutyPolicyUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[OnCallDutyPolicyUserOwnerArgs, OnCallDutyPolicyUserOwnerState]{}, err
	}

	return infer.ReadResponse[OnCallDutyPolicyUserOwnerArgs, OnCallDutyPolicyUserOwnerState]{
		ID:     state.ResourceID,
		Inputs: state.OnCallDutyPolicyUserOwnerArgs,
		State:  state,
	}, nil
}

func (r *OnCallDutyPolicyUserOwner) Update(ctx context.Context, req infer.UpdateRequest[OnCallDutyPolicyUserOwnerArgs, OnCallDutyPolicyUserOwnerState]) (infer.UpdateResponse[OnCallDutyPolicyUserOwnerState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[OnCallDutyPolicyUserOwnerState]{
			Output: OnCallDutyPolicyUserOwnerState{OnCallDutyPolicyUserOwnerArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}

	if err := c.UpdateResource(ctx, "on-call-duty-policy-owner-user", req.ID, data); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}

	result, err := c.ReadResource(ctx, "on-call-duty-policy-owner-user", req.ID, SelectFields(OnCallDutyPolicyUserOwnerState{}))
	if err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}

	var state OnCallDutyPolicyUserOwnerState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[OnCallDutyPolicyUserOwnerState]{}, err
	}

	return infer.UpdateResponse[OnCallDutyPolicyUserOwnerState]{
		Output: state,
	}, nil
}

func (r *OnCallDutyPolicyUserOwner) Delete(ctx context.Context, req infer.DeleteRequest[OnCallDutyPolicyUserOwnerState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "on-call-duty-policy-owner-user", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
