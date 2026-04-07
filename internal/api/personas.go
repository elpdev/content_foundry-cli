package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/leo/content-foundry-cli/internal/models"
)

type PersonaService struct {
	client *Client
}

func NewPersonaService(c *Client) *PersonaService {
	return &PersonaService{client: c}
}

func (s *PersonaService) List(ctx context.Context) ([]models.Persona, error) {
	body, _, err := s.client.Get(ctx, "/api/v1/personas", nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Personas []models.Persona `json:"personas"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing personas: %w", err)
	}
	return wrapper.Personas, nil
}

func (s *PersonaService) Get(ctx context.Context, id int64) (*models.Persona, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/personas/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Persona models.Persona `json:"persona"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing persona: %w", err)
	}
	return &wrapper.Persona, nil
}

func (s *PersonaService) Create(ctx context.Context, fields map[string]any) (*models.Persona, error) {
	payload := map[string]any{"persona": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/personas", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Persona models.Persona `json:"persona"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing persona: %w", err)
	}
	return &wrapper.Persona, nil
}

func (s *PersonaService) Update(ctx context.Context, id int64, fields map[string]any) (*models.Persona, error) {
	payload := map[string]any{"persona": fields}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/personas/%d", id), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Persona models.Persona `json:"persona"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing persona: %w", err)
	}
	return &wrapper.Persona, nil
}

func (s *PersonaService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/personas/%d", id))
	return err
}
