package handler

import (
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/service"
)

type incidentResponse struct {
	ID               string  `json:"id"`
	MonitorID        string  `json:"monitor_id"`
	StartedAt        string  `json:"started_at"`
	EndedAt          *string `json:"ended_at"`
	LastError        *string `json:"last_error"`
	NotificationSent bool    `json:"notification_sent"`
	DurationSeconds  int64   `json:"duration_seconds"`
	Open             bool    `json:"open"`
}

func (s *Server) handleListIncidents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	onlyOpen := q.Get("open") == "true"
	limit := queryInt(q.Get("limit"), 20)
	offset := queryInt(q.Get("offset"), 0)

	incidents, total, err := s.incidentService.List(r.Context(), onlyOpen, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error", "INTERNAL_ERROR")
		return
	}

	data := make([]incidentResponse, len(incidents))
	for i, inc := range incidents {
		data[i] = toIncidentResponse(inc)
	}

	writeJSON(w, http.StatusOK, paginatedResponse{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (s *Server) handleListMonitorIncidents(w http.ResponseWriter, r *http.Request) {
	monitorID := r.PathValue("id")
	q := r.URL.Query()
	onlyOpen := q.Get("open") == "true"
	limit := queryInt(q.Get("limit"), 20)
	offset := queryInt(q.Get("offset"), 0)

	incidents, total, err := s.incidentService.ListByMonitor(r.Context(), monitorID, onlyOpen, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error", "INTERNAL_ERROR")
		return
	}

	data := make([]incidentResponse, len(incidents))
	for i, inc := range incidents {
		data[i] = toIncidentResponse(inc)
	}

	writeJSON(w, http.StatusOK, paginatedResponse{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func toIncidentResponse(i *domain.Incident) incidentResponse {
	durationSeconds := int64(time.Since(i.StartedAt).Seconds())
	if durationSeconds < 0 {
		durationSeconds = 0
	}
	if i.DurationSeconds != nil {
		durationSeconds = *i.DurationSeconds
	}

	res := incidentResponse{
		ID:               i.ID,
		MonitorID:        i.MonitorID,
		StartedAt:        i.StartedAt.Format(time.RFC3339),
		LastError:        i.LastError,
		NotificationSent: i.NotificationSent,
		DurationSeconds:  durationSeconds,
		Open:             i.EndedAt == nil,
	}
	if i.EndedAt != nil {
		s := i.EndedAt.Format(time.RFC3339)
		res.EndedAt = &s
	}
	return res
}

// garante que service está importado mesmo sem uso direto aqui
var _ *service.IncidentService
