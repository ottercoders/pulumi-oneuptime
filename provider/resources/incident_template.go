package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type IncidentTemplate struct{}

type IncidentTemplateArgs struct {
	ProjectID              *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	TemplateName           string   `pulumi:"templateName" json:"templateName"`
	TemplateDescription    *string  `pulumi:"templateDescription,optional" json:"templateDescription,omitempty"`
	Title                  *string  `pulumi:"title,optional" json:"title,omitempty"`
	Description            *string  `pulumi:"description,optional" json:"description,omitempty"`
	IncidentSeverityID     *string  `pulumi:"incidentSeverityId,optional" json:"incidentSeverityId,omitempty"`
	Monitors               []string `pulumi:"monitors,optional" json:"monitors,omitempty"`
	OnCallDutyPolicies     []string `pulumi:"onCallDutyPolicies,optional" json:"onCallDutyPolicies,omitempty"`
	Labels                 []string `pulumi:"labels,optional" json:"labels,omitempty"`
	ChangeMonitorStatusToID *string `pulumi:"changeMonitorStatusToId,optional" json:"changeMonitorStatusToId,omitempty"`
}

type IncidentTemplateState struct {
	IncidentTemplateArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*IncidentTemplate)(nil)

func (t *IncidentTemplate) Annotate(a infer.Annotator) {
	a.Describe(t, "Manages a OneUptime Incident Template for standardizing incident creation.")
}

func (t *IncidentTemplate) Create(ctx context.Context, req infer.CreateRequest[IncidentTemplateArgs]) (infer.CreateResponse[IncidentTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[IncidentTemplateState]{
			ID:     "preview-id",
			Output: IncidentTemplateState{IncidentTemplateArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[IncidentTemplateState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[IncidentTemplateState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "incident-template", data)
	if err != nil {
		return infer.CreateResponse[IncidentTemplateState]{}, err
	}

	var state IncidentTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[IncidentTemplateState]{}, err
	}
	state.IncidentTemplateArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[IncidentTemplateState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (t *IncidentTemplate) Read(ctx context.Context, req infer.ReadRequest[IncidentTemplateArgs, IncidentTemplateState]) (infer.ReadResponse[IncidentTemplateArgs, IncidentTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "incident-template", req.ID, SelectFields(IncidentTemplateState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[IncidentTemplateArgs, IncidentTemplateState]{}, nil
		}
		return infer.ReadResponse[IncidentTemplateArgs, IncidentTemplateState]{}, err
	}

	var state IncidentTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[IncidentTemplateArgs, IncidentTemplateState]{}, err
	}

	return infer.ReadResponse[IncidentTemplateArgs, IncidentTemplateState]{
		ID:     state.ResourceID,
		Inputs: state.IncidentTemplateArgs,
		State:  state,
	}, nil
}

func (t *IncidentTemplate) Update(ctx context.Context, req infer.UpdateRequest[IncidentTemplateArgs, IncidentTemplateState]) (infer.UpdateResponse[IncidentTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[IncidentTemplateState]{
			Output: IncidentTemplateState{IncidentTemplateArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[IncidentTemplateState]{}, err
	}

	if err := c.UpdateResource(ctx, "incident-template", req.ID, data); err != nil {
		return infer.UpdateResponse[IncidentTemplateState]{}, err
	}

	result, err := c.ReadResource(ctx, "incident-template", req.ID, SelectFields(IncidentTemplateState{}))
	if err != nil {
		return infer.UpdateResponse[IncidentTemplateState]{}, err
	}

	var state IncidentTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[IncidentTemplateState]{}, err
	}

	return infer.UpdateResponse[IncidentTemplateState]{
		Output: state,
	}, nil
}

func (t *IncidentTemplate) Delete(ctx context.Context, req infer.DeleteRequest[IncidentTemplateState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "incident-template", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
