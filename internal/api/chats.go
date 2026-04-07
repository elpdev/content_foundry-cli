package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type ChatService struct {
	client *Client
}

func NewChatService(c *Client) *ChatService {
	return &ChatService{client: c}
}

type ChatListParams struct {
	PersonaID string
	Page      int
	PerPage   int
}

func (p ChatListParams) Values() url.Values {
	v := url.Values{}
	if p.PersonaID != "" {
		v.Set("persona_id", p.PersonaID)
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *ChatService) List(ctx context.Context, p ChatListParams) (*PaginatedResponse[models.Chat], error) {
	return FetchPage[models.Chat](ctx, s.client, "/api/v1/chats", p.Values(), "chats")
}

type ChatDetail struct {
	Chat     models.Chat      `json:"chat"`
	Messages []models.Message `json:"messages"`
}

func (s *ChatService) Get(ctx context.Context, id int64) (*ChatDetail, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/chats/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var detail ChatDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("parsing chat: %w", err)
	}
	return &detail, nil
}

func (s *ChatService) Create(ctx context.Context, prompt string, personaID *int64, model string) (*models.Chat, error) {
	fields := map[string]any{
		"prompt": prompt,
	}
	if personaID != nil {
		fields["persona_id"] = *personaID
	}
	if model != "" {
		fields["model"] = model
	}
	payload := map[string]any{"chat": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/chats", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Chat models.Chat `json:"chat"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing chat: %w", err)
	}
	return &wrapper.Chat, nil
}

func (s *ChatService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/chats/%d", id))
	return err
}

func (s *ChatService) SendMessage(ctx context.Context, chatID int64, content string) error {
	payload := map[string]any{
		"message": map[string]any{
			"content": content,
		},
	}
	_, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/chats/%d/messages", chatID), payload)
	return err
}
