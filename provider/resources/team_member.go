package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type TeamMember struct{}

type TeamMemberArgs struct {
	ProjectID             *string `pulumi:"projectId,optional" json:"projectId,omitempty"`
	TeamID                string  `pulumi:"teamId" json:"teamId"`
	UserID                string  `pulumi:"userId" json:"userId"`
	HasAcceptedInvitation *bool   `pulumi:"hasAcceptedInvitation,optional" json:"hasAcceptedInvitation,omitempty"`
	InvitationAcceptedAt  *string `pulumi:"invitationAcceptedAt,optional" json:"invitationAcceptedAt,omitempty"`
}

type TeamMemberState struct {
	TeamMemberArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*TeamMember)(nil)

func (m *TeamMember) Annotate(a infer.Annotator) {
	a.Describe(m, "Adds a user to a OneUptime team. The user must already exist (registered via the OneUptime UI).")
}

func (m *TeamMember) Create(ctx context.Context, req infer.CreateRequest[TeamMemberArgs]) (infer.CreateResponse[TeamMemberState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[TeamMemberState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[TeamMemberState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[TeamMemberState]{
			ID:     "preview-id",
			Output: TeamMemberState{TeamMemberArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "team-member", data)
	if err != nil {
		return infer.CreateResponse[TeamMemberState]{}, err
	}

	var state TeamMemberState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[TeamMemberState]{}, err
	}
	state.TeamMemberArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[TeamMemberState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (m *TeamMember) Read(ctx context.Context, req infer.ReadRequest[TeamMemberArgs, TeamMemberState]) (infer.ReadResponse[TeamMemberArgs, TeamMemberState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "team-member", req.ID, SelectFields(TeamMemberState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[TeamMemberArgs, TeamMemberState]{}, nil
		}
		return infer.ReadResponse[TeamMemberArgs, TeamMemberState]{}, err
	}

	var state TeamMemberState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[TeamMemberArgs, TeamMemberState]{}, err
	}

	return infer.ReadResponse[TeamMemberArgs, TeamMemberState]{
		ID:     state.ResourceID,
		Inputs: state.TeamMemberArgs,
		State:  state,
	}, nil
}

func (m *TeamMember) Update(ctx context.Context, req infer.UpdateRequest[TeamMemberArgs, TeamMemberState]) (infer.UpdateResponse[TeamMemberState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[TeamMemberState]{
			Output: TeamMemberState{TeamMemberArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[TeamMemberState]{}, err
	}

	if err := c.UpdateResource(ctx, "team-member", req.ID, data); err != nil {
		return infer.UpdateResponse[TeamMemberState]{}, err
	}

	result, err := c.ReadResource(ctx, "team-member", req.ID, SelectFields(TeamMemberState{}))
	if err != nil {
		return infer.UpdateResponse[TeamMemberState]{}, err
	}

	var state TeamMemberState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[TeamMemberState]{}, err
	}

	return infer.UpdateResponse[TeamMemberState]{
		Output: state,
	}, nil
}

func (m *TeamMember) Delete(ctx context.Context, req infer.DeleteRequest[TeamMemberState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "team-member", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
