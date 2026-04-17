package resources

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// Common lookup result fields shared by most lookup functions.
type LookupResult struct {
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Name       string `pulumi:"name" json:"name"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
}

type LookupByNameArgs struct {
	Name string `pulumi:"name"`
}

// lookupByName is a generic helper for all name-based lookups.
func lookupByName(ctx context.Context, apiPath string, name string) (LookupResult, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	results, err := c.ListResources(ctx, apiPath, map[string]interface{}{
		"name": name,
	}, map[string]bool{
		"_id":  true,
		"name": true,
		"slug": true,
	}, 1)
	if err != nil {
		return LookupResult{}, fmt.Errorf("looking up %s by name %q: %w", apiPath, name, err)
	}
	if len(results) == 0 {
		return LookupResult{}, fmt.Errorf("no %s found with name %q", apiPath, name)
	}

	var result LookupResult
	if err := FromMap(results[0], &result); err != nil {
		return LookupResult{}, err
	}
	return result, nil
}

// ── GetMonitorStatus ──

type GetMonitorStatusArgs struct {
	Name string `pulumi:"name"`
}

type GetMonitorStatusResult struct {
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Name       string `pulumi:"name" json:"name"`
	Color      string `pulumi:"color,optional" json:"color,omitempty"`
}

type GetMonitorStatus struct{}

var _ infer.Annotated = (*GetMonitorStatus)(nil)

func (f *GetMonitorStatus) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Monitor Status by name (e.g., 'Operational', 'Degraded', 'Offline').")
}

func (f *GetMonitorStatus) Invoke(ctx context.Context, req infer.FunctionRequest[GetMonitorStatusArgs]) (infer.FunctionResponse[GetMonitorStatusResult], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	results, err := c.ListResources(ctx, "monitor-status", map[string]interface{}{
		"name": req.Input.Name,
	}, map[string]bool{
		"_id":   true,
		"name":  true,
		"color": true,
	}, 1)
	if err != nil {
		return infer.FunctionResponse[GetMonitorStatusResult]{}, fmt.Errorf("looking up monitor status %q: %w", req.Input.Name, err)
	}
	if len(results) == 0 {
		return infer.FunctionResponse[GetMonitorStatusResult]{}, fmt.Errorf("no monitor status found with name %q", req.Input.Name)
	}

	var result GetMonitorStatusResult
	if err := FromMap(results[0], &result); err != nil {
		return infer.FunctionResponse[GetMonitorStatusResult]{}, err
	}

	return infer.FunctionResponse[GetMonitorStatusResult]{Output: result}, nil
}

// ── GetMonitor ──

type GetMonitor struct{}

var _ infer.Annotated = (*GetMonitor)(nil)

func (f *GetMonitor) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Monitor by name.")
}

func (f *GetMonitor) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "monitor", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetTeam ──

type GetTeam struct{}

var _ infer.Annotated = (*GetTeam)(nil)

func (f *GetTeam) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Team by name.")
}

func (f *GetTeam) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "team", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetLabel ──

type GetLabel struct{}

var _ infer.Annotated = (*GetLabel)(nil)

func (f *GetLabel) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Label by name.")
}

func (f *GetLabel) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "label", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetStatusPage ──

type GetStatusPage struct{}

var _ infer.Annotated = (*GetStatusPage)(nil)

func (f *GetStatusPage) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Status Page by name.")
}

func (f *GetStatusPage) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "status-page", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetIncidentState ──

type GetIncidentState struct{}

var _ infer.Annotated = (*GetIncidentState)(nil)

func (f *GetIncidentState) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Incident State by name (e.g., 'Investigating', 'Identified', 'Resolved').")
}

func (f *GetIncidentState) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "incident-state", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetIncidentSeverity ──

type GetIncidentSeverity struct{}

var _ infer.Annotated = (*GetIncidentSeverity)(nil)

func (f *GetIncidentSeverity) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Incident Severity by name (e.g., 'Critical', 'Major', 'Minor').")
}

func (f *GetIncidentSeverity) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "incident-severity", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetAlertState ──

type GetAlertState struct{}

var _ infer.Annotated = (*GetAlertState)(nil)

func (f *GetAlertState) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Alert State by name (e.g., 'Triggered', 'Acknowledged', 'Resolved').")
}

func (f *GetAlertState) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "alert-state", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetAlertSeverity ──

type GetAlertSeverity struct{}

var _ infer.Annotated = (*GetAlertSeverity)(nil)

func (f *GetAlertSeverity) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime Alert Severity by name.")
}

func (f *GetAlertSeverity) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "alert-severity", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}

// ── GetOnCallDutyPolicy ──

type GetOnCallDutyPolicy struct{}

var _ infer.Annotated = (*GetOnCallDutyPolicy)(nil)

func (f *GetOnCallDutyPolicy) Annotate(a infer.Annotator) {
	a.Describe(f, "Look up a OneUptime On-Call Duty Policy by name.")
}

func (f *GetOnCallDutyPolicy) Invoke(ctx context.Context, req infer.FunctionRequest[LookupByNameArgs]) (infer.FunctionResponse[LookupResult], error) {
	result, err := lookupByName(ctx, "on-call-duty-policy", req.Input.Name)
	if err != nil {
		return infer.FunctionResponse[LookupResult]{}, err
	}
	return infer.FunctionResponse[LookupResult]{Output: result}, nil
}
