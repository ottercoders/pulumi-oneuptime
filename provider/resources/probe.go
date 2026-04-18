package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Probe struct{}

type ProbeArgs struct {
	ProjectID       *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name            string  `pulumi:"name" json:"name"`
	Description     *string `pulumi:"description,optional" json:"description,omitempty"`
	IconFileID      *string `pulumi:"iconFileId,optional" json:"iconFileId,omitempty"`
	ProbeVersion    *string `pulumi:"probeVersion,optional" json:"probeVersion,omitempty"`
	ShouldAutoEnableProbeOnNewMonitors *bool `pulumi:"shouldAutoEnableProbeOnNewMonitors,optional" json:"shouldAutoEnableProbeOnNewMonitors,omitempty"`
}

type ProbeState struct {
	ProbeArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	Key        string `pulumi:"key,optional" json:"key,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*Probe)(nil)
var _ infer.ExplicitDependencies[ProbeArgs, ProbeState] = (*Probe)(nil)

func (p *Probe) Annotate(a infer.Annotator) {
	a.Describe(p, "Manages a OneUptime Probe (monitoring agent). The key output field is a secret.")
}

// WireDependencies marks the probe key as always secret.
func (p *Probe) WireDependencies(f infer.FieldSelector, args *ProbeArgs, state *ProbeState) {
	f.OutputField(&state.Key).AlwaysSecret()
}

func (p *Probe) Create(ctx context.Context, req infer.CreateRequest[ProbeArgs]) (infer.CreateResponse[ProbeState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ProbeState]{
			ID:     "preview-id",
			Output: ProbeState{ProbeArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ProbeState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ProbeState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "probe", data)
	if err != nil {
		return infer.CreateResponse[ProbeState]{}, err
	}

	var state ProbeState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ProbeState]{}, err
	}
	state.ProbeArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ProbeState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (p *Probe) Read(ctx context.Context, req infer.ReadRequest[ProbeArgs, ProbeState]) (infer.ReadResponse[ProbeArgs, ProbeState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "probe", req.ID, SelectFields(ProbeState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ProbeArgs, ProbeState]{}, nil
		}
		return infer.ReadResponse[ProbeArgs, ProbeState]{}, err
	}

	var state ProbeState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ProbeArgs, ProbeState]{}, err
	}
	if state.Key == "" {
		state.Key = req.State.Key
	}

	return infer.ReadResponse[ProbeArgs, ProbeState]{
		ID:     state.ResourceID,
		Inputs: state.ProbeArgs,
		State:  state,
	}, nil
}

func (p *Probe) Update(ctx context.Context, req infer.UpdateRequest[ProbeArgs, ProbeState]) (infer.UpdateResponse[ProbeState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ProbeState]{
			Output: ProbeState{ProbeArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ProbeState]{}, err
	}

	if err := c.UpdateResource(ctx, "probe", req.ID, data); err != nil {
		return infer.UpdateResponse[ProbeState]{}, err
	}

	result, err := c.ReadResource(ctx, "probe", req.ID, SelectFields(ProbeState{}))
	if err != nil {
		return infer.UpdateResponse[ProbeState]{}, err
	}

	var state ProbeState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ProbeState]{}, err
	}
	if state.Key == "" {
		state.Key = req.State.Key
	}

	return infer.UpdateResponse[ProbeState]{
		Output: state,
	}, nil
}

func (p *Probe) Delete(ctx context.Context, req infer.DeleteRequest[ProbeState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "probe", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
