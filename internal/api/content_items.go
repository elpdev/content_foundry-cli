package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type ContentItemService struct {
	client *Client
}

func NewContentItemService(c *Client) *ContentItemService {
	return &ContentItemService{client: c}
}

type ContentItemListParams struct {
	Status  string
	Page    int
	PerPage int
}

func (p ContentItemListParams) Values() url.Values {
	v := url.Values{}
	if p.Status != "" {
		v.Set("status", p.Status)
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *ContentItemService) List(ctx context.Context, p ContentItemListParams) (*PaginatedResponse[models.ContentItem], error) {
	return FetchPage[models.ContentItem](ctx, s.client, "/api/v1/content_items", p.Values(), "content_items")
}

func (s *ContentItemService) Get(ctx context.Context, id int64) (*models.ContentItem, []models.ContentItemDraft, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/content_items/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}
	var wrapper struct {
		ContentItem models.ContentItem       `json:"content_item"`
		Drafts      []models.ContentItemDraft `json:"drafts"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, nil, fmt.Errorf("parsing content item: %w", err)
	}
	return &wrapper.ContentItem, wrapper.Drafts, nil
}

func (s *ContentItemService) Process(ctx context.Context, id int64, guidance string) (string, error) {
	var payload any
	if guidance != "" {
		payload = map[string]any{"guidance": guidance}
	}
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/content_items/%d/process_item", id), payload)
	if err != nil {
		return "", err
	}
	var wrapper struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return "", fmt.Errorf("parsing process response: %w", err)
	}
	return wrapper.Message, nil
}

func (s *ContentItemService) GenerateDrafts(ctx context.Context, id int64) (string, error) {
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/content_items/%d/generate_missing_drafts", id), nil)
	if err != nil {
		return "", err
	}
	var wrapper struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return "", fmt.Errorf("parsing generate response: %w", err)
	}
	return wrapper.Message, nil
}
