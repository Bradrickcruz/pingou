package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/repository"
)

type SettingsService struct {
	repo *repository.SettingsRepo
	smtp config.SMTPConfig // vem do ambiente; Enabled() decide o email_available
	ntfy config.NtfyConfig // vem do ambiente; Enabled() decide o ntfy_available
}

func NewSettingsService(repo *repository.SettingsRepo, smtp config.SMTPConfig, ntfy config.NtfyConfig) *SettingsService {
	return &SettingsService{repo: repo, smtp: smtp, ntfy: ntfy}
}

type Settings struct {
	WebhookURL      string `json:"webhook_url"`
	RetentionDays   int    `json:"retention_days"`
	EmailRecipients string `json:"email_recipients"` // "" = e-mail desligado
	EmailAvailable  bool   `json:"email_available"`  // false = SMTP não configurado no ambiente
	NtfyTopic       string `json:"ntfy_topic"`       // "" = ntfy desligado
	NtfyAvailable   bool   `json:"ntfy_available"`   // false = PINGOU_NTFY_URL não configurado
	Timezone        string `json:"timezone"`         // nome IANA; painel, e-mail e ntfy. API e banco seguem em UTC
}

type UpdateSettingsInput struct {
	WebhookURL      *string
	RetentionDays   *int
	EmailRecipients *string
	NtfyTopic       *string
	Timezone        *string
}

const (
	settingEmailRecipients = "email_recipients"
	settingNtfyTopic       = "ntfy_topic"
	settingTimezone        = "timezone"
)

var (
	errEmailUnavailable = fmt.Errorf("%w: email is not configured on the server (PINGOU_SMTP_*)", ErrValidation)
	errNtfyUnavailable  = fmt.Errorf("%w: ntfy is not configured on the server (PINGOU_NTFY_URL)", ErrValidation)
	errNtfyTopicInvalid = fmt.Errorf("%w: ntfy_topic must have 1 to 64 characters: letters, numbers, _ or -", ErrValidation)
	errTimezoneInvalid  = fmt.Errorf("%w: timezone must be a valid IANA time zone name (ex.: America/Sao_Paulo)", ErrValidation)
)

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

	recipients, err := s.repo.Get(ctx, settingEmailRecipients)
	if err != nil {
		return nil, err
	}

	topic, err := s.repo.Get(ctx, settingNtfyTopic)
	if err != nil {
		return nil, err
	}

	tz, err := s.repo.Get(ctx, settingTimezone)
	if err != nil {
		return nil, err
	}
	if tz == "" {
		tz = DefaultTimezone // instalação antiga: chave ainda não existe
	}

	return &Settings{
		WebhookURL:      webhookURL,
		RetentionDays:   retention,
		EmailRecipients: recipients,
		EmailAvailable:  s.smtp.Enabled(),
		NtfyTopic:       topic,
		NtfyAvailable:   s.ntfy.Enabled(),
		Timezone:        tz,
	}, nil
}

func (s *SettingsService) Update(ctx context.Context, in UpdateSettingsInput) (*Settings, error) {
	if in.WebhookURL != nil {
		if *in.WebhookURL != "" {
			if err := validateWebhookURL(*in.WebhookURL); err != nil {
				return nil, err
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

	if in.EmailRecipients != nil {
		list, err := parseRecipients(*in.EmailRecipients)
		if err != nil {
			return nil, err
		}
		// "" sempre aceito: desliga o e-mail mesmo sem SMTP no ambiente
		if len(list) > 0 && !s.smtp.Enabled() {
			return nil, errEmailUnavailable
		}
		if err := s.repo.Set(ctx, settingEmailRecipients, strings.Join(list, ",")); err != nil {
			return nil, err
		}
	}

	if in.NtfyTopic != nil {
		topic := strings.TrimSpace(*in.NtfyTopic)
		// "" sempre aceito: desliga o ntfy mesmo sem PINGOU_NTFY_URL
		if topic != "" {
			if !s.ntfy.Enabled() {
				return nil, errNtfyUnavailable
			}
			if !validNtfyTopic(topic) {
				return nil, errNtfyTopicInvalid
			}
		}
		if err := s.repo.Set(ctx, settingNtfyTopic, topic); err != nil {
			return nil, err
		}
	}

	if in.Timezone != nil {
		tz, err := normalizeTimezone(*in.Timezone)
		if err != nil {
			return nil, err
		}
		if err := s.repo.Set(ctx, settingTimezone, tz); err != nil {
			return nil, err
		}
	}

	return s.Get(ctx)
}

// validateWebhookURL exige http/https e host. Vale para salvar e para testar.
func validateWebhookURL(s string) error {
	u, err := url.ParseRequestURI(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("%w: webhook_url must be a valid http(s) URL", ErrValidation)
	}
	return nil
}

// TestWebhook envia o payload de teste para rawURL ou, se vier vazia, para a URL
// salva. Não salva nada.
func (s *SettingsService) TestWebhook(ctx context.Context, rawURL *string) (*WebhookTestResult, error) {
	target := ""
	if rawURL != nil {
		target = *rawURL
	}
	if target == "" {
		saved, err := s.repo.Get(ctx, "webhook_url")
		if err != nil {
			return nil, fmt.Errorf("get webhook_url: %w", err)
		}
		target = saved
	}
	if target == "" {
		return nil, fmt.Errorf("%w: no webhook url to test", ErrValidation)
	}
	if err := validateWebhookURL(target); err != nil {
		return nil, err
	}

	result := SendTestWebhook(ctx, target)
	return &result, nil
}

// TestEmail envia o e-mail de teste para raw ou, se vier vazio, para os
// destinatários salvos. Não salva nada.
func (s *SettingsService) TestEmail(ctx context.Context, raw *string) (*EmailTestResult, error) {
	if !s.smtp.Enabled() {
		return nil, errEmailUnavailable
	}
	target := ""
	if raw != nil {
		target = *raw
	}
	if strings.TrimSpace(target) == "" {
		saved, err := s.repo.Get(ctx, settingEmailRecipients)
		if err != nil {
			return nil, fmt.Errorf("get email_recipients: %w", err)
		}
		target = saved
	}
	to, err := parseRecipients(target)
	if err != nil {
		return nil, err
	}
	if len(to) == 0 {
		return nil, fmt.Errorf("%w: no email recipients to test", ErrValidation)
	}

	result := SendTestEmail(ctx, s.smtp, to)
	return &result, nil
}

// TestNtfy publica a mensagem de teste em raw ou, se vier vazio, no tópico
// salvo. Não salva nada.
func (s *SettingsService) TestNtfy(ctx context.Context, raw *string) (*NtfyTestResult, error) {
	if !s.ntfy.Enabled() {
		return nil, errNtfyUnavailable
	}
	topic := ""
	if raw != nil {
		topic = strings.TrimSpace(*raw)
	}
	if topic == "" {
		saved, err := s.repo.Get(ctx, settingNtfyTopic)
		if err != nil {
			return nil, fmt.Errorf("get ntfy_topic: %w", err)
		}
		topic = saved
	}
	if topic == "" {
		return nil, fmt.Errorf("%w: no ntfy topic to test", ErrValidation)
	}
	if !validNtfyTopic(topic) {
		return nil, errNtfyTopicInvalid
	}

	result := SendTestNtfy(ctx, s.ntfy, topic)
	return &result, nil
}
