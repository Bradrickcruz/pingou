package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// ── DTOs ──────────────────────────────────────────────────────────────────

type monitorConfigRequest struct {
	Host         *string `json:"host"`
	Port         *int    `json:"port"`
	GraceSeconds *int    `json:"grace_seconds"`
}

type createMonitorRequest struct {
	Name             string               `json:"name"`
	Type             string               `json:"type"`
	URL              string               `json:"url"`
	Config           monitorConfigRequest `json:"config"`
	IntervalSeconds  int                  `json:"interval_seconds"`
	TimeoutSeconds   int                  `json:"timeout_seconds"`
	FailureThreshold int                  `json:"failure_threshold"`
	Enabled          *bool                `json:"enabled"` // ausente = true
}

type updateMonitorRequest struct {
	Name             *string              `json:"name"`
	Type             *string              `json:"type"`
	URL              *string              `json:"url"`
	Config           monitorConfigRequest `json:"config"`
	IntervalSeconds  *int                 `json:"interval_seconds"`
	TimeoutSeconds   *int                 `json:"timeout_seconds"`
	FailureThreshold *int                 `json:"failure_threshold"`
	Enabled          *bool                `json:"enabled"`
}

type monitorResponse struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Type             string         `json:"type"`
	URL              *string        `json:"url"` // null fora de http
	Config           map[string]any `json:"config"`
	PingToken        *string        `json:"ping_token,omitempty"` // só heartbeat
	IntervalSeconds  int            `json:"interval_seconds"`
	TimeoutSeconds   int            `json:"timeout_seconds"`
	FailureThreshold int            `json:"failure_threshold"`
	Enabled          bool           `json:"enabled"`
	CurrentState     string         `json:"current_state"`
	LastCheckedAt    *string        `json:"last_checked_at"`
	LastPingAt       *string        `json:"last_ping_at"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
}

// monitorListItem acrescenta o resumo do painel na listagem.
type monitorListItem struct {
	monitorResponse
	History       []bool   `json:"history"` // últimos 30 checks, mais recente por último
	Uptime24h     *float64 `json:"uptime_24h"`
	LastLatencyMs *int64   `json:"last_latency_ms"`
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

	enabled := req.Enabled == nil || *req.Enabled
	in := service.CreateMonitorInput{
		Name:             req.Name,
		Type:             req.Type,
		URL:              req.URL,
		GraceSeconds:     req.Config.GraceSeconds,
		IntervalSeconds:  req.IntervalSeconds,
		TimeoutSeconds:   req.TimeoutSeconds,
		FailureThreshold: req.FailureThreshold,
		Enabled:          enabled,
	}
	if req.Config.Host != nil {
		in.Host = *req.Config.Host
	}
	if req.Config.Port != nil {
		in.Port = *req.Config.Port
	}

	m, err := s.monitorService.Create(r.Context(), in)
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

	summaries, err := s.monitorService.Summaries(r.Context(), monitors)
	if err != nil {
		slog.Error("monitor summaries", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error", "INTERNAL_ERROR")
		return
	}

	data := make([]monitorListItem, len(monitors))
	for i, m := range monitors {
		item := monitorListItem{monitorResponse: toMonitorResponse(m), History: []bool{}}
		if sum := summaries[m.ID]; sum != nil {
			item.History = sum.History
			item.Uptime24h = sum.Uptime24h
			item.LastLatencyMs = sum.LastLatencyMs
		}
		data[i] = item
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
		Type:             req.Type,
		URL:              req.URL,
		Host:             req.Config.Host,
		Port:             req.Config.Port,
		GraceSeconds:     req.Config.GraceSeconds,
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

func (s *Server) handleRegenerateToken(w http.ResponseWriter, r *http.Request) {
	m, err := s.monitorService.RegenerateToken(r.Context(), r.PathValue("id"))
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
		Type:             string(m.Type),
		Config:           map[string]any{},
		IntervalSeconds:  m.IntervalSeconds,
		TimeoutSeconds:   m.TimeoutSeconds,
		FailureThreshold: m.FailureThreshold,
		Enabled:          m.Enabled,
		CurrentState:     string(m.CurrentState),
		LastCheckedAt:    formatTime(m.LastCheckedAt),
		LastPingAt:       formatTime(m.LastPingAt),
		CreatedAt:        m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	switch m.Type {
	case domain.TypeHTTP:
		url := m.URL
		res.URL = &url
	case domain.TypeTCP:
		res.Config["host"] = m.Config.Host
		res.Config["port"] = m.Config.Port
	case domain.TypeHeartbeat:
		res.Config["grace_seconds"] = m.Config.GraceSeconds
		res.PingToken = m.PingToken
	}
	return res
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format("2006-01-02T15:04:05Z")
	return &s
}

func queryInt(v string, def int) int {
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return n
	}
	return def
}
