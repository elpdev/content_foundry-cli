package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type PublicationService struct {
	client *Client
}

func NewPublicationService(c *Client) *PublicationService {
	return &PublicationService{client: c}
}

type PublicationListParams struct {
	URL     string
	Status  string
	Page    int
	PerPage int
}

func (p PublicationListParams) Values() url.Values {
	v := url.Values{}
	if p.URL != "" {
		v.Set("url", p.URL)
	}
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

func (s *PublicationService) List(ctx context.Context, p PublicationListParams) (*PaginatedResponse[models.PublicationListItem], error) {
	return FetchPage[models.PublicationListItem](ctx, s.client, "/api/v1/publications", p.Values(), "publications")
}

func (s *PublicationService) Get(ctx context.Context, id int64) (*models.PublicationDetail, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/publications/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var detail models.PublicationDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("parsing publication: %w", err)
	}
	return &detail, nil
}
