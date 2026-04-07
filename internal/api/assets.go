package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/leo/content-foundry-cli/internal/models"
)

type AssetService struct {
	client *Client
}

func NewAssetService(c *Client) *AssetService {
	return &AssetService{client: c}
}

type AssetListParams struct {
	Type    string
	Page    int
	PerPage int
}

func (p AssetListParams) Values() url.Values {
	v := url.Values{}
	if p.Type != "" {
		v.Set("type", p.Type)
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *AssetService) List(ctx context.Context, p AssetListParams) (*PaginatedResponse[models.Asset], error) {
	return FetchPage[models.Asset](ctx, s.client, "/api/v1/assets", p.Values(), "assets")
}

func (s *AssetService) Get(ctx context.Context, id int64) (*models.Asset, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/assets/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Asset models.Asset `json:"asset"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing asset: %w", err)
	}
	return &wrapper.Asset, nil
}

func (s *AssetService) Upload(ctx context.Context, filePath string) (*models.Asset, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("copying file: %w", err)
	}
	writer.Close()

	if err := s.client.ensureAuth(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.client.BaseURL+"/api/v1/assets", &buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.client.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if s.client.BrandID > 0 {
		req.Header.Set("X-Brand-Id", fmt.Sprintf("%d", s.client.BrandID))
	}

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("uploading: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, body)
	}

	var wrapper struct {
		Asset models.Asset `json:"asset"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing asset: %w", err)
	}
	return &wrapper.Asset, nil
}

func (s *AssetService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/assets/%d", id))
	return err
}
