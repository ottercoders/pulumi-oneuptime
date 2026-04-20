package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type MonitorSecret struct{}

type MonitorSecretArgs struct {
	ProjectID   *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name        string   `pulumi:"name" json:"name"`
	Description *string  `pulumi:"description,optional" json:"description,omitempty"`
	// SecretValue is the credential itself (API key, token, password).
	// AlwaysSecret is wired on the input so Pulumi encrypts it at rest.
	SecretValue string   `pulumi:"secretValue" json:"secretValue"`
	// MonitorIDs is the list of monitor IDs this secret may be referenced
	// from. The API accepts the ManyToMany shape `[{_id: ...}]` — attached
	// at Create/Update via attachIDRefs.
	MonitorIDs  []string `pulumi:"monitorIds,optional" json:"-"`
}

type MonitorSecretState struct {
	MonitorSecretArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*MonitorSecret)(nil)
var _ infer.ExplicitDependencies[MonitorSecretArgs, MonitorSecretState] = (*MonitorSecret)(nil)

func (s *MonitorSecret) Annotate(a infer.Annotator) {
	a.Describe(s, "A named secret value that can be referenced from Monitor step request bodies and headers (e.g. API keys, bearer tokens). SecretValue is stored encrypted.")
}

func (s *MonitorSecret) WireDependencies(f infer.FieldSelector, args *MonitorSecretArgs, state *MonitorSecretState) {
	f.InputField(&args.SecretValue).Secret()
	f.OutputField(&state.SecretValue).AlwaysSecret()
}

func (s *MonitorSecret) Create(ctx context.Context, req infer.CreateRequest[MonitorSecretArgs]) (infer.CreateResponse[MonitorSecretState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[MonitorSecretState]{
			ID:     "preview-id",
			Output: MonitorSecretState{MonitorSecretArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorSecretState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorSecretState]{}, err
	}
	data["projectId"] = projectID
	attachIDRefs(data, "monitors", req.Inputs.MonitorIDs)

	result, err := c.CreateResource(ctx, "monitor-secret", data)
	if err != nil {
		return infer.CreateResponse[MonitorSecretState]{}, err
	}

	var state MonitorSecretState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorSecretState]{}, err
	}
	state.MonitorSecretArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorSecretState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *MonitorSecret) Read(ctx context.Context, req infer.ReadRequest[MonitorSecretArgs, MonitorSecretState]) (infer.ReadResponse[MonitorSecretArgs, MonitorSecretState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	// Don't ask the server for secretValue — OneUptime doesn't return it on
	// get-item and selecting it has caused schema-drift 400s historically.
	// Preserve the value from prior state.
	sel := SelectFields(MonitorSecretState{})
	delete(sel, "secretValue")

	result, err := c.ReadResource(ctx, "monitor-secret", req.ID, sel)
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorSecretArgs, MonitorSecretState]{}, nil
		}
		return infer.ReadResponse[MonitorSecretArgs, MonitorSecretState]{}, err
	}

	var state MonitorSecretState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorSecretArgs, MonitorSecretState]{}, err
	}
	// Carry secretValue forward from prior state; the server never returns it.
	state.SecretValue = req.State.SecretValue

	return infer.ReadResponse[MonitorSecretArgs, MonitorSecretState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorSecretArgs,
		State:  state,
	}, nil
}

func (s *MonitorSecret) Update(ctx context.Context, req infer.UpdateRequest[MonitorSecretArgs, MonitorSecretState]) (infer.UpdateResponse[MonitorSecretState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorSecretState]{
			Output: MonitorSecretState{MonitorSecretArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorSecretState]{}, err
	}
	attachIDRefs(data, "monitors", req.Inputs.MonitorIDs)

	if err := c.UpdateResource(ctx, "monitor-secret", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorSecretState]{}, err
	}

	sel := SelectFields(MonitorSecretState{})
	delete(sel, "secretValue")
	result, err := c.ReadResource(ctx, "monitor-secret", req.ID, sel)
	if err != nil {
		return infer.UpdateResponse[MonitorSecretState]{}, err
	}

	var state MonitorSecretState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorSecretState]{}, err
	}
	state.MonitorSecretArgs = req.Inputs

	return infer.UpdateResponse[MonitorSecretState]{
		Output: state,
	}, nil
}

func (s *MonitorSecret) Delete(ctx context.Context, req infer.DeleteRequest[MonitorSecretState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor-secret", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
