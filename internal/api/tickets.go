package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type TicketService struct {
	client *Client
}

func NewTicketService(c *Client) *TicketService {
	return &TicketService{client: c}
}

type TicketListParams struct {
	Category   string
	Priority   string
	Unresolved bool
	Page       int
	PerPage    int
}

func (p TicketListParams) Values() url.Values {
	v := url.Values{}
	if p.Category != "" {
		v.Set("category", p.Category)
	}
	if p.Priority != "" {
		v.Set("priority", p.Priority)
	}
	if p.Unresolved {
		v.Set("unresolved", "true")
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *TicketService) List(ctx context.Context, p TicketListParams) (*PaginatedResponse[models.Ticket], error) {
	return FetchPage[models.Ticket](ctx, s.client, "/api/v1/tickets", p.Values(), "tickets")
}

func (s *TicketService) Get(ctx context.Context, id int64) (*models.Ticket, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/tickets/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Ticket models.Ticket `json:"ticket"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing ticket: %w", err)
	}
	return &wrapper.Ticket, nil
}

func (s *TicketService) Create(ctx context.Context, subject, description, category, priority string) (*models.Ticket, error) {
	fields := map[string]any{
		"subject":     subject,
		"description": description,
		"category":    category,
		"priority":    priority,
	}
	payload := map[string]any{"ticket": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/tickets", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Ticket models.Ticket `json:"ticket"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing ticket: %w", err)
	}
	return &wrapper.Ticket, nil
}

func (s *TicketService) Update(ctx context.Context, id int64, fields map[string]any) (*models.Ticket, error) {
	payload := map[string]any{"ticket": fields}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/tickets/%d", id), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Ticket models.Ticket `json:"ticket"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing ticket: %w", err)
	}
	return &wrapper.Ticket, nil
}
