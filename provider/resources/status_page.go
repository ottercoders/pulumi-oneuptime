package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPage struct{}

type StatusPageArgs struct {
	ProjectID                       *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	Name                            string  `pulumi:"name" json:"name"`
	Description                     *string `pulumi:"description,optional" json:"description,omitempty"`
	PageTitle                       *string `pulumi:"pageTitle,optional" json:"pageTitle,omitempty"`
	PageDescription                 *string `pulumi:"pageDescription,optional" json:"pageDescription,omitempty"`
	IsPublicStatusPage              *bool   `pulumi:"isPublicStatusPage,optional" json:"isPublicStatusPage,omitempty"`
	ShowIncidentHistoryInDays       *int    `pulumi:"showIncidentHistoryInDays,optional" json:"showIncidentHistoryInDays,omitempty"`
	ShowAnnouncementHistoryInDays   *int    `pulumi:"showAnnouncementHistoryInDays,optional" json:"showAnnouncementHistoryInDays,omitempty"`
	ShowScheduledEventHistoryInDays *int    `pulumi:"showScheduledEventHistoryInDays,optional" json:"showScheduledEventHistoryInDays,omitempty"`
}

type StatusPageState struct {
	StatusPageArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	Slug       string `pulumi:"slug,optional" json:"slug,omitempty"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPage)(nil)

func (s *StatusPage) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a OneUptime Status Page resource.")
}

func (s *StatusPage) Create(ctx context.Context, req infer.CreateRequest[StatusPageArgs]) (infer.CreateResponse[StatusPageState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[StatusPageState]{
			ID:     "preview-id",
			Output: StatusPageState{StatusPageArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "status-page", data)
	if err != nil {
		return infer.CreateResponse[StatusPageState]{}, err
	}

	var state StatusPageState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageState]{}, err
	}
	state.StatusPageArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *StatusPage) Read(ctx context.Context, req infer.ReadRequest[StatusPageArgs, StatusPageState]) (infer.ReadResponse[StatusPageArgs, StatusPageState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page", req.ID, SelectFields(StatusPageState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageArgs, StatusPageState]{}, nil
		}
		return infer.ReadResponse[StatusPageArgs, StatusPageState]{}, err
	}

	var state StatusPageState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageArgs, StatusPageState]{}, err
	}

	return infer.ReadResponse[StatusPageArgs, StatusPageState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageArgs,
		State:  state,
	}, nil
}

func (s *StatusPage) Update(ctx context.Context, req infer.UpdateRequest[StatusPageArgs, StatusPageState]) (infer.UpdateResponse[StatusPageState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageState]{
			Output: StatusPageState{StatusPageArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page", req.ID, SelectFields(StatusPageState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageState]{}, err
	}

	var state StatusPageState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageState]{}, err
	}

	return infer.UpdateResponse[StatusPageState]{
		Output: state,
	}, nil
}

func (s *StatusPage) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
