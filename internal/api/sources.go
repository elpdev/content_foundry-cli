package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type SourceService struct {
	client *Client
}

func NewSourceService(c *Client) *SourceService {
	return &SourceService{client: c}
}

type SourceListParams struct {
	Active  string
	Page    int
	PerPage int
}

func (p SourceListParams) Values() url.Values {
	v := url.Values{}
	if p.Active != "" {
		v.Set("active", p.Active)
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *SourceService) List(ctx context.Context, p SourceListParams) (*PaginatedResponse[models.Source], error) {
	return FetchPage[models.Source](ctx, s.client, "/api/v1/sources", p.Values(), "sources")
}

func (s *SourceService) GetRaw(ctx context.Context, id int64) ([]byte, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/sources/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s *SourceService) Get(ctx context.Context, id int64) (*models.Source, error) {
	body, err := s.GetRaw(ctx, id)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Source models.Source `json:"source"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}
	return &wrapper.Source, nil
}

func (s *SourceService) Create(ctx context.Context, name, sourceType, pollingSchedule, prompt string, active bool, cfg map[string]any) (*models.Source, error) {
	fields := map[string]any{
		"name":   name,
		"type":   sourceType,
		"active": active,
	}
	if pollingSchedule != "" {
		fields["polling_schedule"] = pollingSchedule
	}
	if prompt != "" {
		fields["prompt"] = prompt
	}
	if cfg != nil {
		fields["config"] = cfg
	}
	payload := map[string]any{"source": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/sources", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Source models.Source `json:"source"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}
	return &wrapper.Source, nil
}

func (s *SourceService) Update(ctx context.Context, id int64, fields map[string]any) (*models.Source, error) {
	payload := map[string]any{"source": fields}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/sources/%d", id), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Source models.Source `json:"source"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}
	return &wrapper.Source, nil
}

func (s *SourceService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/sources/%d", id))
	return err
}

func (s *SourceService) Fetch(ctx context.Context, id int64) (*models.Source, string, error) {
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/sources/%d/fetch", id), nil)
	if err != nil {
		return nil, "", err
	}
	var wrapper struct {
		Source  models.Source `json:"source"`
		Message string        `json:"message"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, "", fmt.Errorf("parsing fetch response: %w", err)
	}
	return &wrapper.Source, wrapper.Message, nil
}
