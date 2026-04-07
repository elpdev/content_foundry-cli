package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/leo/content-foundry-cli/internal/models"
)

type AccountService struct {
	client *Client
}

func NewAccountService(c *Client) *AccountService {
	return &AccountService{client: c}
}

type AccountDetail struct {
	Account models.Account     `json:"account"`
	Users   []AccountUserBrief `json:"users"`
}

type AccountUserBrief struct {
	ID     int64    `json:"id"`
	UserID int64    `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

func (s *AccountService) Get(ctx context.Context) (*AccountDetail, error) {
	body, _, err := s.client.Get(ctx, "/api/v1/account", nil)
	if err != nil {
		return nil, err
	}
	var detail AccountDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("parsing account: %w", err)
	}
	return &detail, nil
}

func (s *AccountService) Update(ctx context.Context, name, slug string) (*models.Account, error) {
	fields := map[string]any{}
	if name != "" {
		fields["name"] = name
	}
	if slug != "" {
		fields["slug"] = slug
	}
	payload := map[string]any{"account": fields}
	body, _, err := s.client.Patch(ctx, "/api/v1/account", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Account models.Account `json:"account"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing account: %w", err)
	}
	return &wrapper.Account, nil
}

// Members

func (s *AccountService) ListMembers(ctx context.Context) ([]models.AccountUser, error) {
	body, _, err := s.client.Get(ctx, "/api/v1/account_users", nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		AccountUsers []models.AccountUser `json:"account_users"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing account users: %w", err)
	}
	return wrapper.AccountUsers, nil
}

func (s *AccountService) GetMember(ctx context.Context, id int64) (*models.AccountUser, []int64, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/account_users/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}
	var wrapper struct {
		AccountUser     models.AccountUser `json:"account_user"`
		GrantedBrandIDs []int64            `json:"granted_brand_ids"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, nil, fmt.Errorf("parsing account user: %w", err)
	}
	return &wrapper.AccountUser, wrapper.GrantedBrandIDs, nil
}

func (s *AccountService) UpdateMember(ctx context.Context, id int64, fields map[string]any) (*models.AccountUser, error) {
	payload := map[string]any{"account_user": fields}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/account_users/%d", id), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		AccountUser models.AccountUser `json:"account_user"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing account user: %w", err)
	}
	return &wrapper.AccountUser, nil
}

func (s *AccountService) DeleteMember(ctx context.Context, id int64) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/account_users/%d", id))
	return err
}

func (s *AccountService) UpdateBrandAccess(ctx context.Context, id int64, brandIDs []int64) (*models.AccountUser, error) {
	payload := map[string]any{"brand_ids": brandIDs}
	body, _, err := s.client.Patch(ctx, fmt.Sprintf("/api/v1/account_users/%d/update_brand_access", id), payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		AccountUser models.AccountUser `json:"account_user"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing account user: %w", err)
	}
	return &wrapper.AccountUser, nil
}

// Invitations

func (s *AccountService) ListInvitations(ctx context.Context) ([]models.AccountInvitation, error) {
	body, _, err := s.client.Get(ctx, "/api/v1/account_invitations", nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Invitations []models.AccountInvitation `json:"account_invitations"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing invitations: %w", err)
	}
	return wrapper.Invitations, nil
}

func (s *AccountService) GetInvitation(ctx context.Context, token string) (*models.AccountInvitation, error) {
	body, _, err := s.client.Get(ctx, fmt.Sprintf("/api/v1/account_invitations/%s", token), nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Invitation models.AccountInvitation `json:"account_invitation"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing invitation: %w", err)
	}
	return &wrapper.Invitation, nil
}

func (s *AccountService) CreateInvitation(ctx context.Context, fields map[string]any) (*models.AccountInvitation, error) {
	payload := map[string]any{"account_invitation": fields}
	body, _, err := s.client.Post(ctx, "/api/v1/account_invitations", payload)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Invitation models.AccountInvitation `json:"account_invitation"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing invitation: %w", err)
	}
	return &wrapper.Invitation, nil
}

func (s *AccountService) DeleteInvitation(ctx context.Context, token string) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("/api/v1/account_invitations/%s", token))
	return err
}
