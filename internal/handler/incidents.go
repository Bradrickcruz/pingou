package handler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type incidentResponse struct {
	ID               string                   `json:"id"`
	MonitorID        string                   `json:"monitor_id"`
	Monitor          *incidentMonitorResponse `json:"monitor"`
	StartedAt        string                   `json:"started_at"`
	EndedAt          *string                  `json:"ended_at"`
	LastError        *string                  `json:"last_error"`
	FailureCount     int64                    `json:"failure_count"`
	NotificationSent bool                     `json:"notification_sent"`
	DurationSeconds  int64                    `json:"duration_seconds"`
	Open             bool                     `json:"open"`
}

type incidentMonitorResponse struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Target *string `json:"target"` // URL (http), host:porta (tcp), null (heartbeat)
}

func (s *Server) handleListIncidents(w http.ResponseWriter, r *http.Request) {
	s.listIncidents(w, r, r.URL.Query().Get("monitor_id"))
}

func (s *Server) handleListMonitorIncidents(w http.ResponseWriter, r *http.Request) {
	s.listIncidents(w, r, r.PathValue("id"))
}

// listIncidents atende as duas rotas: ?status=open|resolved (ou o antigo
// ?open=true), monitor_id, limit e offset.
func (s *Server) listIncidents(w http.ResponseWriter, r *http.Request, monitorID string) {
	q := r.URL.Query()
	status, ok := incidentStatus(q)
	if !ok {
		writeError(w, http.StatusBadRequest, "status must be open or resolved", "VALIDATION_ERROR")
		return
	}
	filter := domain.IncidentFilter{
		MonitorID: monitorID,
		Status:    status,
		Limit:     min(queryInt(q.Get("limit"), 20), 100),
		Offset:    queryInt(q.Get("offset"), 0),
	}

	incidents, total, err := s.incidentService.List(r.Context(), filter)
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
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

func incidentStatus(q url.Values) (string, bool) {
	switch status := q.Get("status"); status {
	case "":
		if q.Get("open") == "true" {
			return domain.IncidentStatusOpen, true
		}
		return "", true
	case domain.IncidentStatusOpen, domain.IncidentStatusResolved:
		return status, true
	default:
		return "", false
	}
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
		FailureCount:     i.FailureCount,
		NotificationSent: i.NotificationSent,
		DurationSeconds:  durationSeconds,
		Open:             i.EndedAt == nil,
	}
	if i.MonitorName != "" {
		res.Monitor = &incidentMonitorResponse{
			ID:     i.MonitorID,
			Name:   i.MonitorName,
			Type:   string(i.MonitorType),
			Target: i.MonitorTarget,
		}
	}
	if i.EndedAt != nil {
		s := i.EndedAt.Format(time.RFC3339)
		res.EndedAt = &s
	}
	return res
}
