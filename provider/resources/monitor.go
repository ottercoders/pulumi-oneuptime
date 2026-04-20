package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Monitor struct{}

type MonitorArgs struct {
	ProjectID                      *string                `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name                           string                 `pulumi:"name" json:"name"`
	Description                    *string                `pulumi:"description,optional" json:"description,omitempty"`
	MonitorType                    string                 `pulumi:"monitorType" json:"monitorType"`
	CurrentMonitorStatusID         string                 `pulumi:"currentMonitorStatusId" json:"currentMonitorStatusId"`
	DisableActiveMonitoring        *bool                  `pulumi:"disableActiveMonitoring,optional" json:"disableActiveMonitoring,omitempty"`
	MonitoringInterval             *string                `pulumi:"monitoringInterval,optional" json:"monitoringInterval,omitempty"`

	// MonitorSteps + DefaultMonitorStatusID write together into the
	// monitorSteps typed envelope. Both carry json:"-" so the default ToMap
	// skips them; attachMonitorSteps injects the correct wire shape after.
	MonitorSteps                   []MonitorStep          `pulumi:"monitorSteps,optional" json:"-"`
	DefaultMonitorStatusID         *string                `pulumi:"defaultMonitorStatusId,optional" json:"-"`

	// Labels is a list of label resource IDs. The API wants `[{_id: ...}]`
	// on write; attachLabels handles the transform.
	Labels                         []string               `pulumi:"labels,optional" json:"-"`

	CustomFields                   map[string]interface{} `pulumi:"customFields,optional" json:"customFields,omitempty"`
	PostUpdatesToWorkspaceChannels []WorkspaceChannelRef  `pulumi:"postUpdatesToWorkspaceChannels,optional" json:"postUpdatesToWorkspaceChannels,omitempty"`
}

type MonitorState struct {
	MonitorArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`

	// Server-generated secret keys used to construct webhook/agent URLs for
	// IncomingRequest, Server, and IncomingEmail monitor types respectively.
	// Wired as AlwaysSecret via WireDependencies.
	IncomingRequestSecretKey string `pulumi:"incomingRequestSecretKey,optional" json:"incomingRequestSecretKey,omitempty"`
	ServerMonitorSecretKey   string `pulumi:"serverMonitorSecretKey,optional" json:"serverMonitorSecretKey,omitempty"`
	IncomingEmailSecretKey   string `pulumi:"incomingEmailSecretKey,optional" json:"incomingEmailSecretKey,omitempty"`

	// Most-recent probe/push payloads, populated by the server once monitoring
	// has run at least once. Output-only.
	IncomingMonitorRequest      *IncomingMonitorRequest      `pulumi:"incomingMonitorRequest,optional" json:"incomingMonitorRequest,omitempty"`
	IncomingEmailMonitorRequest *IncomingEmailMonitorRequest `pulumi:"incomingEmailMonitorRequest,optional" json:"incomingEmailMonitorRequest,omitempty"`
	ServerMonitorResponse       *ServerMonitorResponse       `pulumi:"serverMonitorResponse,optional" json:"serverMonitorResponse,omitempty"`
}

var _ infer.Annotated = (*Monitor)(nil)
var _ infer.ExplicitDependencies[MonitorArgs, MonitorState] = (*Monitor)(nil)

func (m *Monitor) Annotate(a infer.Annotator) {
	a.Describe(m, "Manages a OneUptime Monitor resource. Set monitorSteps to configure probes, criteria, and incident/alert triggers; the correct sub-monitor field (logMonitor, snmpMonitor, dnsMonitor, ...) on each step must match the parent monitorType.")
}

// WireDependencies marks the three server-generated secret keys AlwaysSecret
// so Pulumi stores them encrypted and omits them from unencrypted output.
func (m *Monitor) WireDependencies(f infer.FieldSelector, args *MonitorArgs, state *MonitorState) {
	f.OutputField(&state.IncomingRequestSecretKey).AlwaysSecret()
	f.OutputField(&state.ServerMonitorSecretKey).AlwaysSecret()
	f.OutputField(&state.IncomingEmailSecretKey).AlwaysSecret()
}

func (m *Monitor) Create(ctx context.Context, req infer.CreateRequest[MonitorArgs]) (infer.CreateResponse[MonitorState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[MonitorState]{
			ID:     "preview-id",
			Output: MonitorState{MonitorArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}
	data["projectId"] = projectID

	if err := attachMonitorSteps(data, req.Inputs.MonitorSteps, req.Inputs.DefaultMonitorStatusID); err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}
	attachLabels(data, req.Inputs.Labels)

	result, err := c.CreateResource(ctx, "monitor", data)
	if err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}

	var state MonitorState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[MonitorState]{}, err
	}
	state.MonitorArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[MonitorState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (m *Monitor) Read(ctx context.Context, req infer.ReadRequest[MonitorArgs, MonitorState]) (infer.ReadResponse[MonitorArgs, MonitorState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "monitor", req.ID, SelectFields(MonitorState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[MonitorArgs, MonitorState]{}, nil
		}
		return infer.ReadResponse[MonitorArgs, MonitorState]{}, err
	}

	var state MonitorState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[MonitorArgs, MonitorState]{}, err
	}
	// FromMap via unwrapTypedValues flattens _type/value envelopes, so the
	// nested monitorSteps/criteria shapes come back as plain nested maps.
	// They aren't re-hydrated into the MonitorStep struct on Read today
	// because Pulumi's refresh loop only needs the _id/createdAt/updatedAt
	// to detect drift; if callers explicitly compare inputs, they should
	// use ignoreChanges: ["monitorSteps"] until we add a round-trip
	// hydrator. Kept intentionally minimal.
	// Preserve the secret keys from state across refreshes (the API only
	// exposes them on get-item endpoints already — this is defensive).
	state.IncomingRequestSecretKey = reqStateSecret(req.State.IncomingRequestSecretKey, state.IncomingRequestSecretKey)
	state.ServerMonitorSecretKey = reqStateSecret(req.State.ServerMonitorSecretKey, state.ServerMonitorSecretKey)
	state.IncomingEmailSecretKey = reqStateSecret(req.State.IncomingEmailSecretKey, state.IncomingEmailSecretKey)

	return infer.ReadResponse[MonitorArgs, MonitorState]{
		ID:     state.ResourceID,
		Inputs: state.MonitorArgs,
		State:  state,
	}, nil
}

// reqStateSecret keeps a previously-seen secret value when the server stops
// returning it (e.g. after a secret rotation endpoint), while preferring a
// fresh value when the server provides one.
func reqStateSecret(prior, fresh string) string {
	if fresh != "" {
		return fresh
	}
	return prior
}

func (m *Monitor) Update(ctx context.Context, req infer.UpdateRequest[MonitorArgs, MonitorState]) (infer.UpdateResponse[MonitorState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[MonitorState]{
			Output: MonitorState{MonitorArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}

	if err := attachMonitorSteps(data, req.Inputs.MonitorSteps, req.Inputs.DefaultMonitorStatusID); err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}
	attachLabels(data, req.Inputs.Labels)

	if err := c.UpdateResource(ctx, "monitor", req.ID, data); err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}

	result, err := c.ReadResource(ctx, "monitor", req.ID, SelectFields(MonitorState{}))
	if err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}

	var state MonitorState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[MonitorState]{}, err
	}
	state.MonitorArgs = req.Inputs
	// Preserve secret keys across updates — they don't change and the read
	// select may drop them silently on schema-drifted builds.
	state.IncomingRequestSecretKey = reqStateSecret(req.State.IncomingRequestSecretKey, state.IncomingRequestSecretKey)
	state.ServerMonitorSecretKey = reqStateSecret(req.State.ServerMonitorSecretKey, state.ServerMonitorSecretKey)
	state.IncomingEmailSecretKey = reqStateSecret(req.State.IncomingEmailSecretKey, state.IncomingEmailSecretKey)

	return infer.UpdateResponse[MonitorState]{
		Output: state,
	}, nil
}

func (m *Monitor) Delete(ctx context.Context, req infer.DeleteRequest[MonitorState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "monitor", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
