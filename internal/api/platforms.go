package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type PlatformService struct {
	client *Client
}

func NewPlatformService(c *Client) *PlatformService {
	return &PlatformService{client: c}
}

type PlatformListParams struct {
	Active  string
	Page    int
	PerPage int
}

func (p PlatformListParams) Values() url.Values {
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

func (s *PlatformService) List(ctx context.Context, p PlatformListParams) (*PaginatedResponse[models.Platform], error) {
	return FetchPage[models.Platform](ctx, s.client, "/api/v1/platforms", p.Values(), "platforms")
}

func (s *PlatformService) Get(ctx context.Context, id int64) (*models.Platform, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/platforms/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Platform models.Platform `json:"platform"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing platform: %w", err)
	}
	return &wrapper.Platform, nil
}

func (s *PlatformService) Create(ctx context.Context, name, platformType, slug string, active bool, promptTemplate, model string) (*models.Platform, error) {
	fields := map[string]any{
		"name":   name,
		"type":   platformType,
		"active": active,
	}
	if slug != "" {
		fields["slug"] = slug
	}
	if promptTemplate != "" {
		fields["prompt_template"] = promptTemplate
	}
	if model != "" {
		fields["model_id"] = model
	}
	payload := map[string]any{"platform": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/platforms", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Platform models.Platform `json:"platform"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing platform: %w", err)
	}
	return &wrapper.Platform, nil
}

func (s *PlatformService) Update(ctx context.Context, id int64, fields map[string]any) (*models.Platform, error) {
	payload := map[string]any{"platform": fields}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/platforms/%d", id), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Platform models.Platform `json:"platform"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing platform: %w", err)
	}
	return &wrapper.Platform, nil
}

func (s *PlatformService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/platforms/%d", id))
	return err
}

// Label methods

func (s *PlatformService) ListLabels(ctx context.Context, platformID int64, p PlatformListParams) (*PaginatedResponse[models.Label], error) {
	path := fmt.Sprintf("/api/v1/platforms/%d/labels", platformID)
	return FetchPage[models.Label](ctx, s.client, path, p.Values(), "labels")
}

func (s *PlatformService) GetLabel(ctx context.Context, platformID, labelID int64) (*models.Label, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/platforms/%d/labels/%d", platformID, labelID), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Label models.Label `json:"label"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing label: %w", err)
	}
	return &wrapper.Label, nil
}

func (s *PlatformService) CreateLabel(ctx context.Context, platformID int64, fields map[string]any) (*models.Label, error) {
	payload := map[string]any{"label": fields}
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/platforms/%d/labels", platformID), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Label models.Label `json:"label"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing label: %w", err)
	}
	return &wrapper.Label, nil
}

func (s *PlatformService) UpdateLabel(ctx context.Context, platformID, labelID int64, fields map[string]any) (*models.Label, error) {
	payload := map[string]any{"label": fields}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/platforms/%d/labels/%d", platformID, labelID), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Label models.Label `json:"label"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing label: %w", err)
	}
	return &wrapper.Label, nil
}

func (s *PlatformService) DeleteLabel(ctx context.Context, platformID, labelID int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/platforms/%d/labels/%d", platformID, labelID))
	return err
}
