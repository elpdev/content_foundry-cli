package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type ActivityService struct {
	client *Client
}

func NewActivityService(c *Client) *ActivityService {
	return &ActivityService{client: c}
}

type ActivityListParams struct {
	Page    int
	PerPage int
}

func (p ActivityListParams) Values() url.Values {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *ActivityService) List(ctx context.Context, p ActivityListParams) (*PaginatedResponse[models.ActivityEvent], error) {
	return FetchPage[models.ActivityEvent](ctx, s.client, "/api/v1/activity_events", p.Values(), "activity_events")
}
