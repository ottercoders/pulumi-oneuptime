package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageResource struct{}

type StatusPageResourceArgs struct {
	ProjectID          *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID       string  `pulumi:"statusPageId" json:"statusPageId"`
	StatusPageGroupID  *string `pulumi:"statusPageGroupId,optional" json:"statusPageGroupId,omitempty"`
	MonitorID          *string `pulumi:"monitorId,optional" json:"monitorId,omitempty"`
	MonitorGroupID     *string `pulumi:"monitorGroupId,optional" json:"monitorGroupId,omitempty"`
	DisplayName        *string `pulumi:"displayName,optional" json:"displayName,omitempty"`
	DisplayDescription *string `pulumi:"displayDescription,optional" json:"displayDescription,omitempty"`
	DisplayTooltip     *string `pulumi:"displayTooltip,optional" json:"displayTooltip,omitempty"`
	Order              *int    `pulumi:"order,optional" json:"order,omitempty"`
	ShowCurrentStatus  *bool   `pulumi:"showCurrentStatus,optional" json:"showCurrentStatus,omitempty"`
	ShowUptimePercent  *bool   `pulumi:"showUptimePercent,optional" json:"showUptimePercent,omitempty"`
}

type StatusPageResourceState struct {
	StatusPageResourceArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageResource)(nil)

func (s *StatusPageResource) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime Status Page Resource. Links a monitor or monitor group to a status page, optionally within a group.")
}

func (s *StatusPageResource) Create(ctx context.Context, req infer.CreateRequest[StatusPageResourceArgs]) (infer.CreateResponse[StatusPageResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.CreateResponse[StatusPageResourceState]{
			ID:     "preview-id",
			Output: StatusPageResourceState{StatusPageResourceArgs: req.Inputs},
		}, nil
	}

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageResourceState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageResourceState]{}, err
	}
	data["projectId"] = projectID

	result, err := c.CreateResource(ctx, "status-page-resource", data)
	if err != nil {
		return infer.CreateResponse[StatusPageResourceState]{}, err
	}

	var state StatusPageResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageResourceState]{}, err
	}
	state.StatusPageResourceArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageResourceState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *StatusPageResource) Read(ctx context.Context, req infer.ReadRequest[StatusPageResourceArgs, StatusPageResourceState]) (infer.ReadResponse[StatusPageResourceArgs, StatusPageResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-resource", req.ID, SelectFields(StatusPageResourceState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageResourceArgs, StatusPageResourceState]{}, nil
		}
		return infer.ReadResponse[StatusPageResourceArgs, StatusPageResourceState]{}, err
	}

	var state StatusPageResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageResourceArgs, StatusPageResourceState]{}, err
	}

	return infer.ReadResponse[StatusPageResourceArgs, StatusPageResourceState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageResourceArgs,
		State:  state,
	}, nil
}

func (s *StatusPageResource) Update(ctx context.Context, req infer.UpdateRequest[StatusPageResourceArgs, StatusPageResourceState]) (infer.UpdateResponse[StatusPageResourceState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageResourceState]{
			Output: StatusPageResourceState{StatusPageResourceArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageResourceState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-resource", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageResourceState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-resource", req.ID, SelectFields(StatusPageResourceState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageResourceState]{}, err
	}

	var state StatusPageResourceState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageResourceState]{}, err
	}

	return infer.UpdateResponse[StatusPageResourceState]{
		Output: state,
	}, nil
}

func (s *StatusPageResource) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageResourceState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-resource", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
