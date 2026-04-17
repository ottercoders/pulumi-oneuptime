package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type IncidentNoteTemplate struct{}

type IncidentNoteTemplateArgs struct {
	ProjectID           *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	TemplateName        string  `pulumi:"templateName" json:"templateName"`
	TemplateDescription *string `pulumi:"templateDescription,optional" json:"templateDescription,omitempty"`
	Note                string  `pulumi:"note" json:"note"`
}

type IncidentNoteTemplateState struct {
	IncidentNoteTemplateArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*IncidentNoteTemplate)(nil)

func (t *IncidentNoteTemplate) Annotate(a infer.Annotator) {
	a.Describe(t, "Manages a OneUptime Incident Note Template for reusable incident notes.")
}

func (t *IncidentNoteTemplate) Create(ctx context.Context, req infer.CreateRequest[IncidentNoteTemplateArgs]) (infer.CreateResponse[IncidentNoteTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[IncidentNoteTemplateState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[IncidentNoteTemplateState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[IncidentNoteTemplateState]{
			ID:     "preview-id",
			Output: IncidentNoteTemplateState{IncidentNoteTemplateArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "incident-note-template", data)
	if err != nil {
		return infer.CreateResponse[IncidentNoteTemplateState]{}, err
	}

	var state IncidentNoteTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[IncidentNoteTemplateState]{}, err
	}
	state.IncidentNoteTemplateArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[IncidentNoteTemplateState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (t *IncidentNoteTemplate) Read(ctx context.Context, req infer.ReadRequest[IncidentNoteTemplateArgs, IncidentNoteTemplateState]) (infer.ReadResponse[IncidentNoteTemplateArgs, IncidentNoteTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "incident-note-template", req.ID, SelectFields(IncidentNoteTemplateState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[IncidentNoteTemplateArgs, IncidentNoteTemplateState]{}, nil
		}
		return infer.ReadResponse[IncidentNoteTemplateArgs, IncidentNoteTemplateState]{}, err
	}

	var state IncidentNoteTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[IncidentNoteTemplateArgs, IncidentNoteTemplateState]{}, err
	}

	return infer.ReadResponse[IncidentNoteTemplateArgs, IncidentNoteTemplateState]{
		ID:     state.ResourceID,
		Inputs: state.IncidentNoteTemplateArgs,
		State:  state,
	}, nil
}

func (t *IncidentNoteTemplate) Update(ctx context.Context, req infer.UpdateRequest[IncidentNoteTemplateArgs, IncidentNoteTemplateState]) (infer.UpdateResponse[IncidentNoteTemplateState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[IncidentNoteTemplateState]{
			Output: IncidentNoteTemplateState{IncidentNoteTemplateArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[IncidentNoteTemplateState]{}, err
	}

	if err := c.UpdateResource(ctx, "incident-note-template", req.ID, data); err != nil {
		return infer.UpdateResponse[IncidentNoteTemplateState]{}, err
	}

	result, err := c.ReadResource(ctx, "incident-note-template", req.ID, SelectFields(IncidentNoteTemplateState{}))
	if err != nil {
		return infer.UpdateResponse[IncidentNoteTemplateState]{}, err
	}

	var state IncidentNoteTemplateState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[IncidentNoteTemplateState]{}, err
	}

	return infer.UpdateResponse[IncidentNoteTemplateState]{
		Output: state,
	}, nil
}

func (t *IncidentNoteTemplate) Delete(ctx context.Context, req infer.DeleteRequest[IncidentNoteTemplateState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "incident-note-template", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
