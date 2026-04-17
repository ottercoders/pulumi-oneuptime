package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageAnnouncement struct{}

type StatusPageAnnouncementArgs struct {
	ProjectID         *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageIDs     []string `pulumi:"statusPageIds" json:"statusPages"`
	Title             string   `pulumi:"title" json:"title"`
	Description       *string  `pulumi:"description,optional" json:"description,omitempty"`
	ShowAnnouncementAt *string `pulumi:"showAnnouncementAt,optional" json:"showAnnouncementAt,omitempty"`
	EndAnnouncementAt *string  `pulumi:"endAnnouncementAt,optional" json:"endAnnouncementAt,omitempty"`
}

type StatusPageAnnouncementState struct {
	StatusPageAnnouncementArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageAnnouncement)(nil)

func (s *StatusPageAnnouncement) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages an announcement on one or more OneUptime Status Pages.")
}

func (s *StatusPageAnnouncement) Create(ctx context.Context, req infer.CreateRequest[StatusPageAnnouncementArgs]) (infer.CreateResponse[StatusPageAnnouncementState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageAnnouncementState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageAnnouncementState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[StatusPageAnnouncementState]{
			ID:     "preview-id",
			Output: StatusPageAnnouncementState{StatusPageAnnouncementArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "status-page-announcement", data)
	if err != nil {
		return infer.CreateResponse[StatusPageAnnouncementState]{}, err
	}

	var state StatusPageAnnouncementState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageAnnouncementState]{}, err
	}
	state.StatusPageAnnouncementArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageAnnouncementState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *StatusPageAnnouncement) Read(ctx context.Context, req infer.ReadRequest[StatusPageAnnouncementArgs, StatusPageAnnouncementState]) (infer.ReadResponse[StatusPageAnnouncementArgs, StatusPageAnnouncementState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-announcement", req.ID, SelectFields(StatusPageAnnouncementState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageAnnouncementArgs, StatusPageAnnouncementState]{}, nil
		}
		return infer.ReadResponse[StatusPageAnnouncementArgs, StatusPageAnnouncementState]{}, err
	}

	var state StatusPageAnnouncementState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageAnnouncementArgs, StatusPageAnnouncementState]{}, err
	}

	return infer.ReadResponse[StatusPageAnnouncementArgs, StatusPageAnnouncementState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageAnnouncementArgs,
		State:  state,
	}, nil
}

func (s *StatusPageAnnouncement) Update(ctx context.Context, req infer.UpdateRequest[StatusPageAnnouncementArgs, StatusPageAnnouncementState]) (infer.UpdateResponse[StatusPageAnnouncementState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageAnnouncementState]{
			Output: StatusPageAnnouncementState{StatusPageAnnouncementArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageAnnouncementState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-announcement", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageAnnouncementState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-announcement", req.ID, SelectFields(StatusPageAnnouncementState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageAnnouncementState]{}, err
	}

	var state StatusPageAnnouncementState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageAnnouncementState]{}, err
	}

	return infer.UpdateResponse[StatusPageAnnouncementState]{
		Output: state,
	}, nil
}

func (s *StatusPageAnnouncement) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageAnnouncementState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-announcement", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
