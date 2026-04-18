package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type AlertState struct{}

type AlertStateArgs struct {
	ProjectID           *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name                string  `pulumi:"name" json:"name"`
	Color               string  `pulumi:"color" json:"color"`
	Description         *string `pulumi:"description,optional" json:"description,omitempty"`
	IsCreatedState      *bool   `pulumi:"isCreatedState,optional" json:"isCreatedState,omitempty"`
	IsAcknowledgedState *bool   `pulumi:"isAcknowledgedState,optional" json:"isAcknowledgedState,omitempty"`
	IsResolvedState     *bool   `pulumi:"isResolvedState,optional" json:"isResolvedState,omitempty"`
	Order               *int    `pulumi:"order,optional" json:"order,omitempty"`
}

type AlertStateState struct {
	AlertStateArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*AlertState)(nil)

func (a *AlertState) Annotate(ann infer.Annotator) {
	ann.Describe(a, "Manages a OneUptime Alert State resource.")
}

func (a *AlertState) Create(ctx context.Context, req infer.CreateRequest[AlertStateArgs]) (infer.CreateResponse[AlertStateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[AlertStateState]{
			ID:     "preview-id",
			Output: AlertStateState{AlertStateArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[AlertStateState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[AlertStateState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "alert-state", data)
	if err != nil {
		return infer.CreateResponse[AlertStateState]{}, err
	}

	var state AlertStateState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[AlertStateState]{}, err
	}
	state.AlertStateArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[AlertStateState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (a *AlertState) Read(ctx context.Context, req infer.ReadRequest[AlertStateArgs, AlertStateState]) (infer.ReadResponse[AlertStateArgs, AlertStateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "alert-state", req.ID, SelectFields(AlertStateState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[AlertStateArgs, AlertStateState]{}, nil
		}
		return infer.ReadResponse[AlertStateArgs, AlertStateState]{}, err
	}

	var state AlertStateState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[AlertStateArgs, AlertStateState]{}, err
	}

	return infer.ReadResponse[AlertStateArgs, AlertStateState]{
		ID:     state.ResourceID,
		Inputs: state.AlertStateArgs,
		State:  state,
	}, nil
}

func (a *AlertState) Update(ctx context.Context, req infer.UpdateRequest[AlertStateArgs, AlertStateState]) (infer.UpdateResponse[AlertStateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[AlertStateState]{
			Output: AlertStateState{AlertStateArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[AlertStateState]{}, err
	}

	if err := c.UpdateResource(ctx, "alert-state", req.ID, data); err != nil {
		return infer.UpdateResponse[AlertStateState]{}, err
	}

	result, err := c.ReadResource(ctx, "alert-state", req.ID, SelectFields(AlertStateState{}))
	if err != nil {
		return infer.UpdateResponse[AlertStateState]{}, err
	}

	var state AlertStateState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[AlertStateState]{}, err
	}

	return infer.UpdateResponse[AlertStateState]{
		Output: state,
	}, nil
}

func (a *AlertState) Delete(ctx context.Context, req infer.DeleteRequest[AlertStateState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "alert-state", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
