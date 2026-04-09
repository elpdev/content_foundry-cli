package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type BrandService struct {
	client *Client
}

func NewBrandService(c *Client) *BrandService {
	return &BrandService{client: c}
}

type BrandListParams struct {
	Page    int
	PerPage int
}

func (p BrandListParams) Values() url.Values {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *BrandService) List(ctx context.Context, p BrandListParams) (*PaginatedResponse[models.Brand], error) {
	return FetchPage[models.Brand](ctx, s.client, "/api/v1/brands", p.Values(), "brands")
}

func brandRefPath(ref string) string {
	return "/api/v1/brands/" + url.PathEscape(ref)
}

func (s *BrandService) GetByRef(ctx context.Context, ref string) (*models.BrandDetail, error) {
	body, _, err := s.client.Get(ctx, brandRefPath(ref), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Brand models.BrandDetail `json:"brand"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing brand: %w", err)
	}
	return &wrapper.Brand, nil
}

func (s *BrandService) Get(ctx context.Context, id int64) (*models.BrandDetail, error) {
	return s.GetByRef(ctx, fmt.Sprintf("%d", id))
}

func (s *BrandService) Create(ctx context.Context, name, slug, description string) (*models.Brand, error) {
	fields := map[string]any{
		"name": name,
	}
	if slug != "" {
		fields["slug"] = slug
	}
	if description != "" {
		fields["description"] = description
	}
	payload := map[string]any{"brand": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/brands", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Brand models.Brand `json:"brand"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing brand: %w", err)
	}
	return &wrapper.Brand, nil
}

func (s *BrandService) UpdateByRef(ctx context.Context, ref string, fields map[string]any) (*models.BrandDetail, error) {
	payload := map[string]any{"brand": fields}
	body, _, err := s.client.Patch(ctx, brandRefPath(ref), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Brand models.BrandDetail `json:"brand"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing brand: %w", err)
	}
	return &wrapper.Brand, nil
}

func (s *BrandService) Update(ctx context.Context, id int64, fields map[string]any) (*models.BrandDetail, error) {
	return s.UpdateByRef(ctx, fmt.Sprintf("%d", id), fields)
}

func (s *BrandService) DeleteByRef(ctx context.Context, ref string) error {
	_, err := s.client.Delete(ctx, brandRefPath(ref))
	return err
}

func (s *BrandService) Delete(ctx context.Context, id int64) error {
	return s.DeleteByRef(ctx, fmt.Sprintf("%d", id))
}
