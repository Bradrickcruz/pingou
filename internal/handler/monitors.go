package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// ── DTOs ──────────────────────────────────────────────────────────────────

type createMonitorRequest struct {
	Name             string `json:"name"`
	URL              string `json:"url"`
	IntervalSeconds  int    `json:"interval_seconds"`
	TimeoutSeconds   int    `json:"timeout_seconds"`
	FailureThreshold int    `json:"failure_threshold"`
	Enabled          bool   `json:"enabled"`
}

type updateMonitorRequest struct {
	Name             *string `json:"name"`
	URL              *string `json:"url"`
	IntervalSeconds  *int    `json:"interval_seconds"`
	TimeoutSeconds   *int    `json:"timeout_seconds"`
	FailureThreshold *int    `json:"failure_threshold"`
	Enabled          *bool   `json:"enabled"`
}

type monitorResponse struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	URL              string  `json:"url"`
	IntervalSeconds  int     `json:"interval_seconds"`
	TimeoutSeconds   int     `json:"timeout_seconds"`
	FailureThreshold int     `json:"failure_threshold"`
	Enabled          bool    `json:"enabled"`
	CurrentState     string  `json:"current_state"`
	LastCheckedAt    *string `json:"last_checked_at"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type paginatedResponse struct {
	Data   any `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ── Handlers ──────────────────────────────────────────────────────────────

func (s *Server) handleCreateMonitor(w http.ResponseWriter, r *http.Request) {
	var req createMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "BAD_REQUEST")
		return
	}

	m, err := s.monitorService.Create(r.Context(), service.CreateMonitorInput{
		Name:             req.Name,
		URL:              req.URL,
		IntervalSeconds:  req.IntervalSeconds,
		TimeoutSeconds:   req.TimeoutSeconds,
		FailureThreshold: req.FailureThreshold,
		Enabled:          req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toMonitorResponse(m))
}

func (s *Server) handleListMonitors(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := domain.MonitorFilter{
		Limit:  queryInt(q.Get("limit"), 20),
		Offset: queryInt(q.Get("offset"), 0),
	}

	if v := q.Get("enabled"); v != "" {
		b := v == "true"
		filter.Enabled = &b
	}
	if v := q.Get("state"); v != "" {
		state := domain.MonitorState(v)
		filter.State = &state
	}

	monitors, total, err := s.monitorService.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error", "INTERNAL_ERROR")
		return
	}

	data := make([]monitorResponse, len(monitors))
	for i, m := range monitors {
		data[i] = toMonitorResponse(m)
	}

	writeJSON(w, http.StatusOK, paginatedResponse{
		Data:   data,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

func (s *Server) handleGetMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.monitorService.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMonitorResponse(m))
}

func (s *Server) handleUpdateMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req updateMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "BAD_REQUEST")
		return
	}

	m, err := s.monitorService.Update(r.Context(), id, service.UpdateMonitorInput{
		Name:             req.Name,
		URL:              req.URL,
		IntervalSeconds:  req.IntervalSeconds,
		TimeoutSeconds:   req.TimeoutSeconds,
		FailureThreshold: req.FailureThreshold,
		Enabled:          req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toMonitorResponse(m))
}

func (s *Server) handleDeleteMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.monitorService.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Helpers ───────────────────────────────────────────────────────────────

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error(), "NOT_FOUND")
	case errors.Is(err, service.ErrLimitReached):
		writeError(w, http.StatusConflict, err.Error(), "LIMIT_REACHED")
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
	default:
		slog.Error("internal error", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error", "INTERNAL_ERROR")
	}
}

func toMonitorResponse(m *domain.Monitor) monitorResponse {
	res := monitorResponse{
		ID:               m.ID,
		Name:             m.Name,
		URL:              m.URL,
		IntervalSeconds:  m.IntervalSeconds,
		TimeoutSeconds:   m.TimeoutSeconds,
		FailureThreshold: m.FailureThreshold,
		Enabled:          m.Enabled,
		CurrentState:     string(m.CurrentState),
		CreatedAt:        m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if m.LastCheckedAt != nil {
		s := m.LastCheckedAt.Format("2006-01-02T15:04:05Z")
		res.LastCheckedAt = &s
	}
	return res
}

func queryInt(v string, def int) int {
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return n
	}
	return def
}
