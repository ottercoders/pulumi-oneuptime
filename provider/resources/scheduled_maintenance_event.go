package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ScheduledMaintenanceEvent struct{}

type ScheduledMaintenanceEventArgs struct {
	ProjectID                  *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Title                      string   `pulumi:"title" json:"title"`
	Description                *string  `pulumi:"description,optional" json:"description,omitempty"`
	StartsAt                   string   `pulumi:"startsAt" json:"startsAt"`
	EndsAt                     string   `pulumi:"endsAt" json:"endsAt"`
	CurrentScheduledMaintenanceStateID *string `pulumi:"currentScheduledMaintenanceStateId,optional" json:"currentScheduledMaintenanceStateId,omitempty"`
	Monitors                   []string `pulumi:"monitors,optional" json:"monitors,omitempty"`
	StatusPages                []string `pulumi:"statusPages,optional" json:"statusPages,omitempty"`
	Labels                     []string `pulumi:"labels,optional" json:"labels,omitempty"`
	ChangeMonitorStatusToID    *string  `pulumi:"changeMonitorStatusToId,optional" json:"changeMonitorStatusToId,omitempty"`
	ShouldStatusPageSubscribersBeNotifiedOnEventCreated *bool `pulumi:"shouldStatusPageSubscribersBeNotifiedOnEventCreated,optional" json:"shouldStatusPageSubscribersBeNotifiedOnEventCreated,omitempty"`
	ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing *bool `pulumi:"shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing,optional" json:"shouldStatusPageSubscribersBeNotifiedWhenEventChangedToOngoing,omitempty"`
	ShouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded *bool `pulumi:"shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded,optional" json:"shouldStatusPageSubscribersBeNotifiedWhenEventChangedToEnded,omitempty"`
}

type ScheduledMaintenanceEventState struct {
	ScheduledMaintenanceEventArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*ScheduledMaintenanceEvent)(nil)

func (s *ScheduledMaintenanceEvent) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime Scheduled Maintenance event (planned downtime window).")
}

func (s *ScheduledMaintenanceEvent) Create(ctx context.Context, req infer.CreateRequest[ScheduledMaintenanceEventArgs]) (infer.CreateResponse[ScheduledMaintenanceEventState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[ScheduledMaintenanceEventState]{
			ID:     "preview-id",
			Output: ScheduledMaintenanceEventState{ScheduledMaintenanceEventArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceEventState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceEventState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "scheduled-maintenance", data)
	if err != nil {
		return infer.CreateResponse[ScheduledMaintenanceEventState]{}, err
	}

	var state ScheduledMaintenanceEventState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[ScheduledMaintenanceEventState]{}, err
	}
	state.ScheduledMaintenanceEventArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[ScheduledMaintenanceEventState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *ScheduledMaintenanceEvent) Read(ctx context.Context, req infer.ReadRequest[ScheduledMaintenanceEventArgs, ScheduledMaintenanceEventState]) (infer.ReadResponse[ScheduledMaintenanceEventArgs, ScheduledMaintenanceEventState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "scheduled-maintenance", req.ID, SelectFields(ScheduledMaintenanceEventState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[ScheduledMaintenanceEventArgs, ScheduledMaintenanceEventState]{}, nil
		}
		return infer.ReadResponse[ScheduledMaintenanceEventArgs, ScheduledMaintenanceEventState]{}, err
	}

	var state ScheduledMaintenanceEventState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[ScheduledMaintenanceEventArgs, ScheduledMaintenanceEventState]{}, err
	}

	return infer.ReadResponse[ScheduledMaintenanceEventArgs, ScheduledMaintenanceEventState]{
		ID:     state.ResourceID,
		Inputs: state.ScheduledMaintenanceEventArgs,
		State:  state,
	}, nil
}

func (s *ScheduledMaintenanceEvent) Update(ctx context.Context, req infer.UpdateRequest[ScheduledMaintenanceEventArgs, ScheduledMaintenanceEventState]) (infer.UpdateResponse[ScheduledMaintenanceEventState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[ScheduledMaintenanceEventState]{
			Output: ScheduledMaintenanceEventState{ScheduledMaintenanceEventArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceEventState]{}, err
	}

	if err := c.UpdateResource(ctx, "scheduled-maintenance", req.ID, data); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceEventState]{}, err
	}

	result, err := c.ReadResource(ctx, "scheduled-maintenance", req.ID, SelectFields(ScheduledMaintenanceEventState{}))
	if err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceEventState]{}, err
	}

	var state ScheduledMaintenanceEventState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[ScheduledMaintenanceEventState]{}, err
	}

	return infer.UpdateResponse[ScheduledMaintenanceEventState]{
		Output: state,
	}, nil
}

func (s *ScheduledMaintenanceEvent) Delete(ctx context.Context, req infer.DeleteRequest[ScheduledMaintenanceEventState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "scheduled-maintenance", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
