package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Incident struct{}

type IncidentArgs struct {
	ProjectID              *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Title                  string  `pulumi:"title" json:"title"`
	Description            *string `pulumi:"description,optional" json:"description,omitempty"`
	CurrentIncidentStateID string  `pulumi:"currentIncidentStateId" json:"currentIncidentStateId"`
	IncidentSeverityID     string  `pulumi:"incidentSeverityId" json:"incidentSeverityId"`
	DeclaredAt             *string `pulumi:"declaredAt,optional" json:"declaredAt,omitempty"`
}

type IncidentState struct {
	IncidentArgs
	ID             string `pulumi:"id" json:"_id"`
	Slug           string `pulumi:"slug,optional" json:"slug,omitempty"`
	IncidentNumber *int   `pulumi:"incidentNumber,optional" json:"incidentNumber,omitempty"`
	CreatedAt      string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt      string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Incident)(nil)

func (i *Incident) Annotate(a infer.Annotator) {
	a.Describe(i, "Manages a OneUptime Incident resource.")
}

func (i *Incident) Create(ctx context.Context, req infer.CreateRequest[IncidentArgs]) (infer.CreateResponse[IncidentState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[IncidentState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[IncidentState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[IncidentState]{
			ID:     "preview-id",
			Output: IncidentState{IncidentArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "incident", data)
	if err != nil {
		return infer.CreateResponse[IncidentState]{}, err
	}

	var state IncidentState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[IncidentState]{}, err
	}
	state.IncidentArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[IncidentState]{
		ID:     state.ID,
		Output: state,
	}, nil
}

func (i *Incident) Read(ctx context.Context, req infer.ReadRequest[IncidentArgs, IncidentState]) (infer.ReadResponse[IncidentArgs, IncidentState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "incident", req.ID, SelectFields(IncidentState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[IncidentArgs, IncidentState]{}, nil
		}
		return infer.ReadResponse[IncidentArgs, IncidentState]{}, err
	}

	var state IncidentState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[IncidentArgs, IncidentState]{}, err
	}

	return infer.ReadResponse[IncidentArgs, IncidentState]{
		ID:     state.ID,
		Inputs: state.IncidentArgs,
		State:  state,
	}, nil
}

func (i *Incident) Update(ctx context.Context, req infer.UpdateRequest[IncidentArgs, IncidentState]) (infer.UpdateResponse[IncidentState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[IncidentState]{
			Output: IncidentState{IncidentArgs: req.Inputs, ID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[IncidentState]{}, err
	}

	if err := c.UpdateResource(ctx, "incident", req.ID, data); err != nil {
		return infer.UpdateResponse[IncidentState]{}, err
	}

	result, err := c.ReadResource(ctx, "incident", req.ID, SelectFields(IncidentState{}))
	if err != nil {
		return infer.UpdateResponse[IncidentState]{}, err
	}

	var state IncidentState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[IncidentState]{}, err
	}

	return infer.UpdateResponse[IncidentState]{
		Output: state,
	}, nil
}

func (i *Incident) Delete(ctx context.Context, req infer.DeleteRequest[IncidentState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "incident", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
