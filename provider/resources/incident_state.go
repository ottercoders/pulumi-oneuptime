package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// IncidentStateResource manages OneUptime Incident State resources.
// The Go type is suffixed with "Resource" to avoid collision with the
// IncidentState struct in incident.go (state output for the Incident resource).
type IncidentStateResource struct{}

type IncidentStateResourceArgs struct {
	ProjectID           *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name                string  `pulumi:"name" json:"name"`
	Color               string  `pulumi:"color" json:"color"`
	Description         *string `pulumi:"description,optional" json:"description,omitempty"`
	IsCreatedState      *bool   `pulumi:"isCreatedState,optional" json:"isCreatedState,omitempty"`
	IsAcknowledgedState *bool   `pulumi:"isAcknowledgedState,optional" json:"isAcknowledgedState,omitempty"`
	IsResolvedState     *bool   `pulumi:"isResolvedState,optional" json:"isResolvedState,omitempty"`
	Order               *int    `pulumi:"order,optional" json:"order,omitempty"`
}

type IncidentStateResourceState struct {
	IncidentStateResourceArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*IncidentStateResource)(nil)

func (i *IncidentStateResource) Annotate(a infer.Annotator) {
	a.Describe(i, "Manages a OneUptime Incident State resource.")
	a.SetToken("resources", "IncidentState")
}

func (i *IncidentStateResource) Create(ctx context.Context, req infer.CreateRequest[IncidentStateResourceArgs]) (infer.CreateResponse[IncidentStateResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[IncidentStateResourceState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[IncidentStateResourceState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[IncidentStateResourceState]{
			ID:     "preview-id",
			Output: IncidentStateResourceState{IncidentStateResourceArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "incident-state", data)
	if err != nil {
		return infer.CreateResponse[IncidentStateResourceState]{}, err
	}

	var state IncidentStateResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[IncidentStateResourceState]{}, err
	}
	state.IncidentStateResourceArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[IncidentStateResourceState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (i *IncidentStateResource) Read(ctx context.Context, req infer.ReadRequest[IncidentStateResourceArgs, IncidentStateResourceState]) (infer.ReadResponse[IncidentStateResourceArgs, IncidentStateResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "incident-state", req.ID, SelectFields(IncidentStateResourceState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[IncidentStateResourceArgs, IncidentStateResourceState]{}, nil
		}
		return infer.ReadResponse[IncidentStateResourceArgs, IncidentStateResourceState]{}, err
	}

	var state IncidentStateResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[IncidentStateResourceArgs, IncidentStateResourceState]{}, err
	}

	return infer.ReadResponse[IncidentStateResourceArgs, IncidentStateResourceState]{
		ID:     state.ResourceID,
		Inputs: state.IncidentStateResourceArgs,
		State:  state,
	}, nil
}

func (i *IncidentStateResource) Update(ctx context.Context, req infer.UpdateRequest[IncidentStateResourceArgs, IncidentStateResourceState]) (infer.UpdateResponse[IncidentStateResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[IncidentStateResourceState]{
			Output: IncidentStateResourceState{IncidentStateResourceArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[IncidentStateResourceState]{}, err
	}

	if err := c.UpdateResource(ctx, "incident-state", req.ID, data); err != nil {
		return infer.UpdateResponse[IncidentStateResourceState]{}, err
	}

	result, err := c.ReadResource(ctx, "incident-state", req.ID, SelectFields(IncidentStateResourceState{}))
	if err != nil {
		return infer.UpdateResponse[IncidentStateResourceState]{}, err
	}

	var state IncidentStateResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[IncidentStateResourceState]{}, err
	}

	return infer.UpdateResponse[IncidentStateResourceState]{
		Output: state,
	}, nil
}

func (i *IncidentStateResource) Delete(ctx context.Context, req infer.DeleteRequest[IncidentStateResourceState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "incident-state", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
