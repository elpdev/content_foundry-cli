package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type NotificationService struct {
	client *Client
}

func NewNotificationService(c *Client) *NotificationService {
	return &NotificationService{client: c}
}

type NotificationListParams struct {
	Page    int
	PerPage int
}

func (p NotificationListParams) Values() url.Values {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *NotificationService) List(ctx context.Context, p NotificationListParams) (*PaginatedResponse[models.Notification], error) {
	return FetchPage[models.Notification](ctx, s.client, "/api/v1/notifications", p.Values(), "notifications")
}

func (s *NotificationService) MarkRead(ctx context.Context, id int64) (*models.Notification, error) {
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/notifications/%d/mark_read", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Notification models.Notification `json:"notification"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing notification: %w", err)
	}
	return &wrapper.Notification, nil
}

func (s *NotificationService) MarkAllRead(ctx context.Context) (string, error) {
	body, _, err := s.client.Post(ctx, "/api/v1/notifications/mark_all_read", nil)
	if err != nil {
		return "", err
	}
	var wrapper struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}
	return wrapper.Message, nil
}
