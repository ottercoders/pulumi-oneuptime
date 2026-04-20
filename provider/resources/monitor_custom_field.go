package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorCustomField struct{}

type MonitorCustomFieldArgs struct {
	ProjectID       *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name            string  `pulumi:"name" json:"name"`
	Description     *string `pulumi:"description,optional" json:"description,omitempty"`
	// CustomFieldType is one of OneUptime's CustomFieldType values
	// (commonly Text, LongText, Number, URL, Date). The exact enum isn't
	// enforced here — server rejects unknown values.
	CustomFieldType string  `pulumi:"customFieldType" json:"customFieldType"`
}

type MonitorCustomFieldState struct {
	MonitorCustomFieldArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorCustomField)(nil)

func (r *MonitorCustomField) Annotate(a infer.Annotator) {
	a.Describe(r, "Declares a custom field schema for Monitors in a project. Values are written via Monitor.customFields at create/update time.")
}

func (r *MonitorCustomField) Create(ctx context.Context, req infer.CreateRequest[MonitorCustomFieldArgs]) (infer.CreateResponse[MonitorCustomFieldState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[MonitorCustomFieldState]{
			ID:     "preview-id",
			Output: MonitorCustomFieldState{MonitorCustomFieldArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorCustomFieldState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorCustomFieldState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "monitor-custom-field", data)
	if err != nil {
		return infer.CreateResponse[MonitorCustomFieldState]{}, err
	}

	var state MonitorCustomFieldState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorCustomFieldState]{}, err
	}
	state.MonitorCustomFieldArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorCustomFieldState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *MonitorCustomField) Read(ctx context.Context, req infer.ReadRequest[MonitorCustomFieldArgs, MonitorCustomFieldState]) (infer.ReadResponse[MonitorCustomFieldArgs, MonitorCustomFieldState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor-custom-field", req.ID, SelectFields(MonitorCustomFieldState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorCustomFieldArgs, MonitorCustomFieldState]{}, nil
		}
		return infer.ReadResponse[MonitorCustomFieldArgs, MonitorCustomFieldState]{}, err
	}

	var state MonitorCustomFieldState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorCustomFieldArgs, MonitorCustomFieldState]{}, err
	}

	return infer.ReadResponse[MonitorCustomFieldArgs, MonitorCustomFieldState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorCustomFieldArgs,
		State:  state,
	}, nil
}

func (r *MonitorCustomField) Update(ctx context.Context, req infer.UpdateRequest[MonitorCustomFieldArgs, MonitorCustomFieldState]) (infer.UpdateResponse[MonitorCustomFieldState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorCustomFieldState]{
			Output: MonitorCustomFieldState{MonitorCustomFieldArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorCustomFieldState]{}, err
	}

	if err := c.UpdateResource(ctx, "monitor-custom-field", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorCustomFieldState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor-custom-field", req.ID, SelectFields(MonitorCustomFieldState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorCustomFieldState]{}, err
	}

	var state MonitorCustomFieldState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorCustomFieldState]{}, err
	}

	return infer.UpdateResponse[MonitorCustomFieldState]{
		Output: state,
	}, nil
}

func (r *MonitorCustomField) Delete(ctx context.Context, req infer.DeleteRequest[MonitorCustomFieldState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-custom-field", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
