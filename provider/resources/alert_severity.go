package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type AlertSeverity struct{}

type AlertSeverityArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string  `pulumi:"name" json:"name"`
	Color       string  `pulumi:"color" json:"color"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
	Order       *int    `pulumi:"order,optional" json:"order,omitempty"`
}

type AlertSeverityState struct {
	AlertSeverityArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*AlertSeverity)(nil)

func (a *AlertSeverity) Annotate(ann infer.Annotator) {
	ann.Describe(a, "Manages a OneUptime Alert Severity resource.")
}

func (a *AlertSeverity) Create(ctx context.Context, req infer.CreateRequest[AlertSeverityArgs]) (infer.CreateResponse[AlertSeverityState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[AlertSeverityState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[AlertSeverityState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[AlertSeverityState]{
			ID:     "preview-id",
			Output: AlertSeverityState{AlertSeverityArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "alert-severity", data)
	if err != nil {
		return infer.CreateResponse[AlertSeverityState]{}, err
	}

	var state AlertSeverityState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[AlertSeverityState]{}, err
	}
	state.AlertSeverityArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[AlertSeverityState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (a *AlertSeverity) Read(ctx context.Context, req infer.ReadRequest[AlertSeverityArgs, AlertSeverityState]) (infer.ReadResponse[AlertSeverityArgs, AlertSeverityState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "alert-severity", req.ID, SelectFields(AlertSeverityState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[AlertSeverityArgs, AlertSeverityState]{}, nil
		}
		return infer.ReadResponse[AlertSeverityArgs, AlertSeverityState]{}, err
	}

	var state AlertSeverityState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[AlertSeverityArgs, AlertSeverityState]{}, err
	}

	return infer.ReadResponse[AlertSeverityArgs, AlertSeverityState]{
		ID:     state.ResourceID,
		Inputs: state.AlertSeverityArgs,
		State:  state,
	}, nil
}

func (a *AlertSeverity) Update(ctx context.Context, req infer.UpdateRequest[AlertSeverityArgs, AlertSeverityState]) (infer.UpdateResponse[AlertSeverityState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[AlertSeverityState]{
			Output: AlertSeverityState{AlertSeverityArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[AlertSeverityState]{}, err
	}

	if err := c.UpdateResource(ctx, "alert-severity", req.ID, data); err != nil {
		return infer.UpdateResponse[AlertSeverityState]{}, err
	}

	result, err := c.ReadResource(ctx, "alert-severity", req.ID, SelectFields(AlertSeverityState{}))
	if err != nil {
		return infer.UpdateResponse[AlertSeverityState]{}, err
	}

	var state AlertSeverityState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[AlertSeverityState]{}, err
	}

	return infer.UpdateResponse[AlertSeverityState]{
		Output: state,
	}, nil
}

func (a *AlertSeverity) Delete(ctx context.Context, req infer.DeleteRequest[AlertSeverityState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "alert-severity", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
