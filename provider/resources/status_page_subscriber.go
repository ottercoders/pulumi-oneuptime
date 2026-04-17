package resources

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type StatusPageSubscriber struct{}

type StatusPageSubscriberArgs struct {
	ProjectID           *string  `pulumi:"projectId,optional" json:"projectId,omitempty"`
	StatusPageID        string   `pulumi:"statusPageId" json:"statusPageId"`
	SubscriberEmail     *string  `pulumi:"subscriberEmail,optional" json:"subscriberEmail,omitempty"`
	SubscriberPhone     *string  `pulumi:"subscriberPhone,optional" json:"subscriberPhone,omitempty"`
	SubscriberWebhook   *string  `pulumi:"subscriberWebhook,optional" json:"subscriberWebhook,omitempty"`
	SubscriberSlackWebhookURL *string `pulumi:"subscriberSlackWebhookUrl,optional" json:"subscriberSlackWebhookUrl,omitempty"`
	IsUnsubscribed      *bool    `pulumi:"isUnsubscribed,optional" json:"isUnsubscribed,omitempty"`
	IsSubscribedToAllResources *bool `pulumi:"isSubscribedToAllResources,optional" json:"isSubscribedToAllResources,omitempty"`
	StatusPageResources []string `pulumi:"statusPageResources,optional" json:"statusPageResources,omitempty"`
}

type StatusPageSubscriberState struct {
	StatusPageSubscriberArgs
	ResourceID string `pulumi:"resourceId" json:"_id"`
	CreatedAt  string `pulumi:"createdAt,optional" json:"createdAt,omitempty"`
	UpdatedAt  string `pulumi:"updatedAt,optional" json:"updatedAt,omitempty"`
}

var _ infer.Annotated = (*StatusPageSubscriber)(nil)

func (s *StatusPageSubscriber) Annotate(a infer.Annotator) {
	a.Describe(s, "Manages a subscriber to a OneUptime Status Page (email, SMS, webhook, or Slack).")
}

func (s *StatusPageSubscriber) Create(ctx context.Context, req infer.CreateRequest[StatusPageSubscriberArgs]) (infer.CreateResponse[StatusPageSubscriberState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	projectID, err := ResolveProjectID(req.Inputs.ProjectID, cfg.ProjectID)
	if err != nil {
		return infer.CreateResponse[StatusPageSubscriberState]{}, err
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.CreateResponse[StatusPageSubscriberState]{}, err
	}
	data["projectId"] = projectID

	if req.DryRun {
		return infer.CreateResponse[StatusPageSubscriberState]{
			ID:     "preview-id",
			Output: StatusPageSubscriberState{StatusPageSubscriberArgs: req.Inputs},
		}, nil
	}

	result, err := c.CreateResource(ctx, "status-page-subscriber", data)
	if err != nil {
		return infer.CreateResponse[StatusPageSubscriberState]{}, err
	}

	var state StatusPageSubscriberState
	if err := FromMap(result, &state); err != nil {
		return infer.CreateResponse[StatusPageSubscriberState]{}, err
	}
	state.StatusPageSubscriberArgs = req.Inputs
	if state.ProjectID == nil {
		state.ProjectID = &projectID
	}

	return infer.CreateResponse[StatusPageSubscriberState]{
		ID:     state.ResourceID,
		Output: state,
	}, nil
}

func (s *StatusPageSubscriber) Read(ctx context.Context, req infer.ReadRequest[StatusPageSubscriberArgs, StatusPageSubscriberState]) (infer.ReadResponse[StatusPageSubscriberArgs, StatusPageSubscriberState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	result, err := c.ReadResource(ctx, "status-page-subscriber", req.ID, SelectFields(StatusPageSubscriberState{}))
	if err != nil {
		if IsNotFound(err) {
			return infer.ReadResponse[StatusPageSubscriberArgs, StatusPageSubscriberState]{}, nil
		}
		return infer.ReadResponse[StatusPageSubscriberArgs, StatusPageSubscriberState]{}, err
	}

	var state StatusPageSubscriberState
	if err := FromMap(result, &state); err != nil {
		return infer.ReadResponse[StatusPageSubscriberArgs, StatusPageSubscriberState]{}, err
	}

	return infer.ReadResponse[StatusPageSubscriberArgs, StatusPageSubscriberState]{
		ID:     state.ResourceID,
		Inputs: state.StatusPageSubscriberArgs,
		State:  state,
	}, nil
}

func (s *StatusPageSubscriber) Update(ctx context.Context, req infer.UpdateRequest[StatusPageSubscriberArgs, StatusPageSubscriberState]) (infer.UpdateResponse[StatusPageSubscriberState], error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if req.DryRun {
		return infer.UpdateResponse[StatusPageSubscriberState]{
			Output: StatusPageSubscriberState{StatusPageSubscriberArgs: req.Inputs, ResourceID: req.ID},
		}, nil
	}

	data, err := ToMap(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[StatusPageSubscriberState]{}, err
	}

	if err := c.UpdateResource(ctx, "status-page-subscriber", req.ID, data); err != nil {
		return infer.UpdateResponse[StatusPageSubscriberState]{}, err
	}

	result, err := c.ReadResource(ctx, "status-page-subscriber", req.ID, SelectFields(StatusPageSubscriberState{}))
	if err != nil {
		return infer.UpdateResponse[StatusPageSubscriberState]{}, err
	}

	var state StatusPageSubscriberState
	if err := FromMap(result, &state); err != nil {
		return infer.UpdateResponse[StatusPageSubscriberState]{}, err
	}

	return infer.UpdateResponse[StatusPageSubscriberState]{
		Output: state,
	}, nil
}

func (s *StatusPageSubscriber) Delete(ctx context.Context, req infer.DeleteRequest[StatusPageSubscriberState]) (infer.DeleteResponse, error) {
	cfg := infer.GetConfig[*Config](ctx)
	c := cfg.GetClient()

	if err := c.DeleteResource(ctx, "status-page-subscriber", req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
