package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type BrandDocumentService struct {
	client *Client
}

func NewBrandDocumentService(c *Client) *BrandDocumentService {
	return &BrandDocumentService{client: c}
}

type BrandDocumentListParams struct {
	Page    int
	PerPage int
}

func (p BrandDocumentListParams) Values() url.Values {
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *BrandDocumentService) List(ctx context.Context, p BrandDocumentListParams) (*PaginatedResponse[models.BrandDocument], error) {
	return FetchPage[models.BrandDocument](ctx, s.client, "/api/v1/brand_documents", p.Values(), "brand_documents")
}

func (s *BrandDocumentService) Get(ctx context.Context, id int64) (*models.BrandDocument, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/brand_documents/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		BrandDocument models.BrandDocument `json:"brand_document"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing brand document: %w", err)
	}
	return &wrapper.BrandDocument, nil
}

func (s *BrandDocumentService) CreateFromURL(ctx context.Context, docURL string) (*models.BrandDocument, string, error) {
	payload := map[string]any{"url": docURL}
	body, _, err := s.client.Post(ctx, "/api/v1/brand_documents", payload)
	if err != nil {
		return nil, "", err
	}
	var wrapper struct {
		BrandDocument models.BrandDocument `json:"brand_document"`
		Message       string               `json:"message"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, "", fmt.Errorf("parsing brand document: %w", err)
	}
	return &wrapper.BrandDocument, wrapper.Message, nil
}

func (s *BrandDocumentService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/brand_documents/%d", id))
	return err
}

func (s *BrandDocumentService) IndexContent(ctx context.Context) (string, error) {
	body, _, err := s.client.Post(ctx, "/api/v1/brand_documents/index_content", nil)
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
