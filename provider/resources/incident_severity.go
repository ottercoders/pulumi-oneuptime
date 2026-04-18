package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type IncidentSeverity struct{}

type IncidentSeverityArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string  `pulumi:"name" json:"name"`
	Color       string  `pulumi:"color" json:"color"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
	Order       *int    `pulumi:"order,optional" json:"order,omitempty"`
}

type IncidentSeverityState struct {
	IncidentSeverityArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*IncidentSeverity)(nil)

func (i *IncidentSeverity) Annotate(a infer.Annotator) {
	a.Describe(i, "Manages a OneUptime Incident Severity resource.")
}

func (i *IncidentSeverity) Create(ctx context.Context, req infer.CreateRequest[IncidentSeverityArgs]) (infer.CreateResponse[IncidentSeverityState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[IncidentSeverityState]{
			ID:     "preview-id",
			Output: IncidentSeverityState{IncidentSeverityArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[IncidentSeverityState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[IncidentSeverityState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "incident-severity", data)
	if err != nil {
		return infer.CreateResponse[IncidentSeverityState]{}, err
	}

	var state IncidentSeverityState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[IncidentSeverityState]{}, err
	}
	state.IncidentSeverityArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[IncidentSeverityState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (i *IncidentSeverity) Read(ctx context.Context, req infer.ReadRequest[IncidentSeverityArgs, IncidentSeverityState]) (infer.ReadResponse[IncidentSeverityArgs, IncidentSeverityState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "incident-severity", req.ID, SelectFields(IncidentSeverityState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[IncidentSeverityArgs, IncidentSeverityState]{}, nil
		}
		return infer.ReadResponse[IncidentSeverityArgs, IncidentSeverityState]{}, err
	}

	var state IncidentSeverityState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[IncidentSeverityArgs, IncidentSeverityState]{}, err
	}

	return infer.ReadResponse[IncidentSeverityArgs, IncidentSeverityState]{
		ID:     state.ResourceID,
		Inputs: state.IncidentSeverityArgs,
		State:  state,
	}, nil
}

func (i *IncidentSeverity) Update(ctx context.Context, req infer.UpdateRequest[IncidentSeverityArgs, IncidentSeverityState]) (infer.UpdateResponse[IncidentSeverityState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[IncidentSeverityState]{
			Output: IncidentSeverityState{IncidentSeverityArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[IncidentSeverityState]{}, err
	}

	if err := c.UpdateResource(ctx, "incident-severity", req.ID, data); err != nil {
		return infer.UpdateResponse[IncidentSeverityState]{}, err
	}

	result, err := c.ReadResource(ctx, "incident-severity", req.ID, SelectFields(IncidentSeverityState{}))
	if err != nil {
		return infer.UpdateResponse[IncidentSeverityState]{}, err
	}

	var state IncidentSeverityState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[IncidentSeverityState]{}, err
	}

	return infer.UpdateResponse[IncidentSeverityState]{
		Output: state,
	}, nil
}

func (i *IncidentSeverity) Delete(ctx context.Context, req infer.DeleteRequest[IncidentSeverityState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "incident-severity", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
