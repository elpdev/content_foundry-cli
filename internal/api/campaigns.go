package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type CampaignService struct {
	client *Client
}

func NewCampaignService(c *Client) *CampaignService {
	return &CampaignService{client: c}
}

type CampaignListParams struct {
	Page    int
	PerPage int
}

func (p CampaignListParams) Values() url.Values {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *CampaignService) List(ctx context.Context, p CampaignListParams) (*PaginatedResponse[models.Campaign], error) {
	return FetchPage[models.Campaign](ctx, s.client, "/api/v1/campaigns", p.Values(), "campaigns")
}

func (s *CampaignService) GetRaw(ctx context.Context, id int64) ([]byte, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/campaigns/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s *CampaignService) Get(ctx context.Context, id int64) (*models.Campaign, error) {
	body, err := s.GetRaw(ctx, id)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Campaign models.Campaign `json:"campaign"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing campaign: %w", err)
	}
	return &wrapper.Campaign, nil
}

func (s *CampaignService) Create(ctx context.Context, name, slug, description string) (*models.Campaign, error) {
	fields := map[string]any{
		"name": name,
		"slug": slug,
	}
	if description != "" {
		fields["description"] = description
	}
	payload := map[string]any{"campaign": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/campaigns", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Campaign models.Campaign `json:"campaign"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing campaign: %w", err)
	}
	return &wrapper.Campaign, nil
}

func (s *CampaignService) Update(ctx context.Context, id int64, fields map[string]any) (*models.Campaign, error) {
	payload := map[string]any{"campaign": fields}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/campaigns/%d", id), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Campaign models.Campaign `json:"campaign"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing campaign: %w", err)
	}
	return &wrapper.Campaign, nil
}

func (s *CampaignService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/campaigns/%d", id))
	return err
}
