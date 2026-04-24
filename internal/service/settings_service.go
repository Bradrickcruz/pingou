package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Bradrickcruz/pingou/internal/repository"
)

type SettingsService struct {
	repo *repository.SettingsRepo
}

func NewSettingsService(repo *repository.SettingsRepo) *SettingsService {
	return &SettingsService{repo: repo}
}

type Settings struct {
	WebhookURL    string `json:"webhook_url"`
	RetentionDays int    `json:"retention_days"`
}

type UpdateSettingsInput struct {
	WebhookURL    *string
	RetentionDays *int
}

func (s *SettingsService) Get(ctx context.Context) (*Settings, error) {
	webhookURL, err := s.repo.Get(ctx, "webhook_url")
	if err != nil {
		return nil, err
	}

	retentionStr, err := s.repo.Get(ctx, "retention_days")
	if err != nil {
		return nil, err
	}

	retention := 30
	fmt.Sscanf(retentionStr, "%d", &retention)

	return &Settings{
		WebhookURL:    webhookURL,
		RetentionDays: retention,
	}, nil
}

func (s *SettingsService) Update(ctx context.Context, in UpdateSettingsInput) (*Settings, error) {
	if in.WebhookURL != nil {
		if *in.WebhookURL != "" {
			if _, err := url.ParseRequestURI(*in.WebhookURL); err != nil {
				return nil, fmt.Errorf("%w: webhook_url must be a valid URL", ErrValidation)
			}
		}
		if err := s.repo.Set(ctx, "webhook_url", *in.WebhookURL); err != nil {
			return nil, err
		}
	}

	if in.RetentionDays != nil {
		if *in.RetentionDays < 7 || *in.RetentionDays > 90 {
			return nil, fmt.Errorf("%w: retention_days must be between 7 and 90", ErrValidation)
		}
		if err := s.repo.Set(ctx, "retention_days", fmt.Sprintf("%d", *in.RetentionDays)); err != nil {
			return nil, err
		}
	}

	return s.Get(ctx)
}
