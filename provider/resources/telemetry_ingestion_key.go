package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type TelemetryIngestionKey struct{}

type TelemetryIngestionKeyArgs struct {
	ProjectID   *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string  `pulumi:"name" json:"name"`
	Description *string `pulumi:"description,optional" json:"description,omitempty"`
}

type TelemetryIngestionKeyState struct {
	TelemetryIngestionKeyArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	SecretKey  string `pulumi:"secretKey,optional" json:"secretKey,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*TelemetryIngestionKey)(nil)
var _ infer.ExplicitDependencies[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState] = (*TelemetryIngestionKey)(nil)

func (r *TelemetryIngestionKey) Annotate(a infer.Annotator) {
	a.Describe(r, "Manages a telemetry ingestion key for pushing metrics, logs, and traces.")
}

// WireDependencies marks the generated secretKey value as always secret.
func (r *TelemetryIngestionKey) WireDependencies(f infer.FieldSelector, args *TelemetryIngestionKeyArgs, state *TelemetryIngestionKeyState) {
	f.OutputField(&state.SecretKey).AlwaysSecret()
}

func (r *TelemetryIngestionKey) Create(ctx context.Context, req infer.CreateRequest[TelemetryIngestionKeyArgs]) (infer.CreateResponse[TelemetryIngestionKeyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[TelemetryIngestionKeyState]{
			ID:     "preview-id",
			Output: TelemetryIngestionKeyState{TelemetryIngestionKeyArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[TelemetryIngestionKeyState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[TelemetryIngestionKeyState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "telemetry-ingestion-key", data)
	if err != nil {
		return infer.CreateResponse[TelemetryIngestionKeyState]{}, err
	}

	var state TelemetryIngestionKeyState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[TelemetryIngestionKeyState]{}, err
	}
	state.TelemetryIngestionKeyArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[TelemetryIngestionKeyState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (r *TelemetryIngestionKey) Read(ctx context.Context, req infer.ReadRequest[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState]) (infer.ReadResponse[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "telemetry-ingestion-key", req.ID, SelectFields(TelemetryIngestionKeyState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState]{}, nil
		}
		return infer.ReadResponse[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState]{}, err
	}

	var state TelemetryIngestionKeyState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState]{}, err
	}
	// Preserve the secret key from prior state (Read endpoint may not return it)
	state.SecretKey = req.State.SecretKey

	return infer.ReadResponse[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState]{
		ID:     state.ResourceID,
		Inputs: state.TelemetryIngestionKeyArgs,
		State:  state,
	}, nil
}

func (r *TelemetryIngestionKey) Update(ctx context.Context, req infer.UpdateRequest[TelemetryIngestionKeyArgs, TelemetryIngestionKeyState]) (infer.UpdateResponse[TelemetryIngestionKeyState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[TelemetryIngestionKeyState]{
			Output: TelemetryIngestionKeyState{TelemetryIngestionKeyArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[TelemetryIngestionKeyState]{}, err
	}

	if err := c.UpdateResource(ctx, "telemetry-ingestion-key", req.ID, data); err != nil {
		return infer.UpdateResponse[TelemetryIngestionKeyState]{}, err
	}

	result, err := c.ReadResource(ctx, "telemetry-ingestion-key", req.ID, SelectFields(TelemetryIngestionKeyState{}))
	if err != nil {
		return infer.UpdateResponse[TelemetryIngestionKeyState]{}, err
	}

	var state TelemetryIngestionKeyState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[TelemetryIngestionKeyState]{}, err
	}
	// Preserve the secret key from prior state (Read endpoint may not return it)
	state.SecretKey = req.State.SecretKey

	return infer.UpdateResponse[TelemetryIngestionKeyState]{
		Output: state,
	}, nil
}

func (r *TelemetryIngestionKey) Delete(ctx context.Context, req infer.DeleteRequest[TelemetryIngestionKeyState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "telemetry-ingestion-key", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
