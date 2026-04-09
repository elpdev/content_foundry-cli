package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type DraftService struct {
	client *Client
}

func NewDraftService(c *Client) *DraftService {
	return &DraftService{client: c}
}

type DraftListParams struct {
	Status       string
	PlatformID   string
	AssignedToID string
	Page         int
	PerPage      int
}

func (p DraftListParams) Values() url.Values {
	v := url.Values{}
	if p.Status != "" {
		v.Set("status", p.Status)
	}
	if p.PlatformID != "" {
		v.Set("platform_id", p.PlatformID)
	}
	if p.AssignedToID != "" {
		v.Set("assigned_to_id", p.AssignedToID)
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("items", fmt.Sprintf("%d", p.PerPage))
	}
	return v
}

func (s *DraftService) List(ctx context.Context, p DraftListParams) (*PaginatedResponse[models.Draft], error) {
	return FetchPage[models.Draft](ctx, s.client, "/api/v1/drafts", p.Values(), "drafts")
}

type DraftDetail struct {
	Draft            models.Draft             `json:"draft"`
	Comments         []models.DraftComment    `json:"comments"`
	RevisionRequests []models.RevisionRequest `json:"revision_requests"`
	Publication      *models.Publication      `json:"publication"`
}

func (s *DraftService) Get(ctx context.Context, id int64) (*DraftDetail, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/drafts/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var detail DraftDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &detail, nil
}

func (s *DraftService) Create(ctx context.Context, title, content string, platformID int64) (*models.Draft, error) {
	payload := map[string]any{
		"draft": map[string]any{
			"title":       title,
			"content":     content,
			"platform_id": platformID,
		},
	}
	body, _, err := s.client.Post(ctx, "/api/v1/drafts", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Draft models.Draft `json:"draft"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &wrapper.Draft, nil
}

func (s *DraftService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/drafts/%d", id))
	return err
}

// Workflow actions

func (s *DraftService) Approve(ctx context.Context, id int64) (*models.Draft, error) {
	return s.workflowAction(ctx, id, "approve", nil)
}

func (s *DraftService) Reject(ctx context.Context, id int64) (*models.Draft, error) {
	return s.workflowAction(ctx, id, "reject", nil)
}

func (s *DraftService) RequestRevision(ctx context.Context, id int64, notes string) (*models.Draft, error) {
	return s.workflowAction(ctx, id, "request_revision", map[string]any{"notes": notes})
}

func (s *DraftService) Schedule(ctx context.Context, id int64, scheduledFor string) (*models.Draft, error) {
	return s.workflowAction(ctx, id, "schedule", map[string]any{"scheduled_for": scheduledFor})
}

func (s *DraftService) Reschedule(ctx context.Context, id int64, scheduledFor string) (*models.Draft, error) {
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/drafts/%d/reschedule", id), map[string]any{"scheduled_for": scheduledFor})
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Draft models.Draft `json:"draft"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &wrapper.Draft, nil
}

func (s *DraftService) Unschedule(ctx context.Context, id int64) (*models.Draft, error) {
	body, _, err := s.client.do(ctx, "DELETE", fmt.Sprintf("/api/v1/drafts/%d/unschedule", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Draft models.Draft `json:"draft"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &wrapper.Draft, nil
}

func (s *DraftService) Assign(ctx context.Context, id, userID int64) (*models.Draft, error) {
	return s.workflowAction(ctx, id, "assign", map[string]any{"assigned_to_id": userID})
}

func (s *DraftService) Unassign(ctx context.Context, id int64) (*models.Draft, error) {
	body, _, err := s.client.do(ctx, "DELETE", fmt.Sprintf("/api/v1/drafts/%d/unassign", id), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Draft models.Draft `json:"draft"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &wrapper.Draft, nil
}

func (s *DraftService) SaveMedia(ctx context.Context, id int64, turnIDs []int64) (*models.Draft, error) {
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/drafts/%d/save_media", id), map[string]any{"media_turn_ids": turnIDs})
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Draft models.Draft `json:"draft"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &wrapper.Draft, nil
}

func (s *DraftService) workflowAction(ctx context.Context, id int64, action string, payload map[string]any) (*models.Draft, error) {
	var body any
	if payload != nil {
		body = payload
	}
	respBody, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/drafts/%d/%s", id, action), body)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Draft models.Draft `json:"draft"`
	}
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &wrapper.Draft, nil
}

// Comments

func (s *DraftService) AddComment(ctx context.Context, draftID int64, body string) (*models.DraftComment, error) {
	payload := map[string]any{"body": body}
	respBody, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/drafts/%d/comments", draftID), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Comment models.DraftComment `json:"comment"`
	}
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing comment: %w", err)
	}
	return &wrapper.Comment, nil
}

func (s *DraftService) DeleteComment(ctx context.Context, draftID, commentID int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/drafts/%d/comments/%d", draftID, commentID))
	return err
}

// Publish

func (s *DraftService) Publish(ctx context.Context, draftID int64, turnIDs []int64) (*models.Publication, string, error) {
	var payload any
	if len(turnIDs) > 0 {
		payload = map[string]any{"media_turn_ids": turnIDs}
	}
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/drafts/%d/publications", draftID), payload)
	if err != nil {
		return nil, "", err
	}
	var wrapper struct {
		Publication models.Publication `json:"publication"`
		Message     string             `json:"message"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, "", fmt.Errorf("parsing publication: %w", err)
	}
	return &wrapper.Publication, wrapper.Message, nil
}
