package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Bradrickcruz/pingou/internal/service"
)

type updateSettingsRequest struct {
	WebhookURL    *string `json:"webhook_url"`
	RetentionDays *int    `json:"retention_days"`
}

type testWebhookRequest struct {
	URL *string `json:"url"`
}

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.settingsService.Get(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (s *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req updateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "BAD_REQUEST")
		return
	}

	settings, err := s.settingsService.Update(r.Context(), toUpdateSettingsInput(req))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

func toUpdateSettingsInput(req updateSettingsRequest) service.UpdateSettingsInput {
	return service.UpdateSettingsInput{
		WebhookURL:    req.WebhookURL,
		RetentionDays: req.RetentionDays,
	}
}

// handleTestWebhook: corpo opcional. Sem corpo, testa a URL salva.
func (s *Server) handleTestWebhook(w http.ResponseWriter, r *http.Request) {
	var req testWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid request body", "BAD_REQUEST")
		return
	}

	result, err := s.settingsService.TestWebhook(r.Context(), req.URL)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
