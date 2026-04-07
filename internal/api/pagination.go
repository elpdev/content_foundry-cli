package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

// PaginatedResponse wraps a page of items with pagination metadata.
type PaginatedResponse[T any] struct {
	Items      []T
	Pagination models.Pagination
}

// FetchPage fetches a single page of items from a paginated endpoint.
// The key parameter is the JSON key containing the items array (e.g., "brands").
func FetchPage[T any](ctx context.Context, c *Client, path string, params url.Values, key string) (*PaginatedResponse[T], error) {
	body, _, err := c.Get(ctx, path, params)
	if err != nil {
		return nil, err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing paginated response: %w", err)
	}

	result := &PaginatedResponse[T]{}

	if itemsRaw, ok := raw[key]; ok {
		if err := json.Unmarshal(itemsRaw, &result.Items); err != nil {
			return nil, fmt.Errorf("parsing %s items: %w", key, err)
		}
	}

	if pagRaw, ok := raw["pagination"]; ok {
		if err := json.Unmarshal(pagRaw, &result.Pagination); err != nil {
			return nil, fmt.Errorf("parsing pagination: %w", err)
		}
	}

	return result, nil
}

// FetchAll iterates through all pages and collects every item.
func FetchAll[T any](ctx context.Context, c *Client, path string, params url.Values, key string) ([]T, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("items", "100")

	var all []T
	page := 1

	for {
		params.Set("page", fmt.Sprintf("%d", page))
		resp, err := FetchPage[T](ctx, c, path, params, key)
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Items...)
		if !resp.Pagination.HasNext() {
			break
		}
		page++
	}

	return all, nil
}
