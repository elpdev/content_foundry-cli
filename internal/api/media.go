package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/leo/content-foundry-cli/internal/models"
)

type MediaSessionService struct {
	client *Client
}

func NewMediaSessionService(c *Client) *MediaSessionService {
	return &MediaSessionService{client: c}
}

type MediaListParams struct {
	Type    string
	Page    int
	PerPage int
}

func (p MediaListParams) Values() url.Values {
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

func (s *MediaSessionService) List(ctx context.Context, p MediaListParams) (*PaginatedResponse[models.MediaSession], error) {
	return FetchPage[models.MediaSession](ctx, s.client, "/api/v1/media_sessions", p.Values(), "media_sessions")
}

type MediaSessionDetail struct {
	MediaSession models.MediaSession `json:"media_session"`
	MediaTurns   []models.MediaTurn  `json:"media_turns"`
}

func (s *MediaSessionService) Get(ctx context.Context, id int64) (*MediaSessionDetail, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/media_sessions/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var detail MediaSessionDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("parsing media session: %w", err)
	}
	return &detail, nil
}

func (s *MediaSessionService) Create(ctx context.Context, mediaType, prompt, aspectRatio string, durationSeconds, captureSeconds int) (*models.MediaSession, error) {
	fields := map[string]any{
		"prompt": prompt,
	}
	if aspectRatio != "" {
		fields["aspect_ratio"] = aspectRatio
	}
	if durationSeconds > 0 {
		fields["duration_seconds"] = durationSeconds
	}
	if captureSeconds > 0 {
		fields["capture_seconds"] = captureSeconds
	}
	payload := map[string]any{
		"media_type":    mediaType,
		"media_session": fields,
	}
	body, _, err := s.client.Post(ctx, "/api/v1/media_sessions", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		MediaSession models.MediaSession `json:"media_session"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing media session: %w", err)
	}
	return &wrapper.MediaSession, nil
}

func (s *MediaSessionService) Delete(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/media_sessions/%d", id))
	return err
}

// Image turns

func (s *MediaSessionService) CreateImageTurn(ctx context.Context, sessionID int64, prompt, aspectRatio string) (*models.ImageTurn, error) {
	fields := map[string]any{
		"user_prompt": prompt,
	}
	if aspectRatio != "" {
		fields["aspect_ratio"] = aspectRatio
	}
	payload := map[string]any{"image_turn": fields}
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/image_turns", sessionID), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		ImageTurn models.ImageTurn `json:"image_turn"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing image turn: %w", err)
	}
	return &wrapper.ImageTurn, nil
}

func (s *MediaSessionService) DeleteImageTurn(ctx context.Context, sessionID, turnID int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/image_turns/%d", sessionID, turnID))
	return err
}

func (s *MediaSessionService) ConvertImageTurn(ctx context.Context, sessionID, turnID int64) (*models.ImageTurn, error) {
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/image_turns/%d/convert", sessionID, turnID), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		ImageTurn models.ImageTurn `json:"image_turn"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing image turn: %w", err)
	}
	return &wrapper.ImageTurn, nil
}

// Video turns

func (s *MediaSessionService) CreateVideoTurn(ctx context.Context, sessionID int64, prompt, aspectRatio string, durationSeconds int) (*models.VideoTurn, error) {
	fields := map[string]any{
		"user_prompt": prompt,
	}
	if aspectRatio != "" {
		fields["aspect_ratio"] = aspectRatio
	}
	if durationSeconds > 0 {
		fields["duration_seconds"] = durationSeconds
	}
	payload := map[string]any{"video_turn": fields}
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/video_turns", sessionID), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		VideoTurn models.VideoTurn `json:"video_turn"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing video turn: %w", err)
	}
	return &wrapper.VideoTurn, nil
}

func (s *MediaSessionService) DeleteVideoTurn(ctx context.Context, sessionID, turnID int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/video_turns/%d", sessionID, turnID))
	return err
}

func (s *MediaSessionService) ExtendVideoTurn(ctx context.Context, sessionID, turnID int64, targetDuration int, extensionPrompt string) (*models.VideoTurn, error) {
	payload := map[string]any{
		"target_duration": targetDuration,
	}
	if extensionPrompt != "" {
		payload["extension_prompt"] = extensionPrompt
	}
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/video_turns/%d/extend_video", sessionID, turnID), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		VideoTurn models.VideoTurn `json:"video_turn"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing video turn: %w", err)
	}
	return &wrapper.VideoTurn, nil
}

// Audio turns

func (s *MediaSessionService) CreateAudioTurn(ctx context.Context, sessionID int64, prompt string, captureSeconds int) (*models.AudioTurn, error) {
	fields := map[string]any{
		"user_prompt": prompt,
	}
	if captureSeconds > 0 {
		fields["capture_seconds"] = captureSeconds
	}
	payload := map[string]any{"audio_turn": fields}
	body, _, err := s.client.Post(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/audio_turns", sessionID), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		AudioTurn models.AudioTurn `json:"audio_turn"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing audio turn: %w", err)
	}
	return &wrapper.AudioTurn, nil
}

func (s *MediaSessionService) DeleteAudioTurn(ctx context.Context, sessionID, turnID int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/media_sessions/%d/audio_turns/%d", sessionID, turnID))
	return err
}
