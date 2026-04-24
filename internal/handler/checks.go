package handler

import (
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type checkResponse struct {
	ID           int64   `json:"id"`
	MonitorID    string  `json:"monitor_id"`
	Success      bool    `json:"success"`
	StatusCode   *int    `json:"status_code"`
	LatencyMs    int64   `json:"latency_ms"`
	ErrorMessage *string `json:"error_message"`
	CheckedAt    string  `json:"checked_at"`
}

func (s *Server) handleListMonitorChecks(w http.ResponseWriter, r *http.Request) {
	monitorID := r.PathValue("id")
	q := r.URL.Query()
	limit := queryInt(q.Get("limit"), 50)
	offset := queryInt(q.Get("offset"), 0)

	checks, total, err := s.monitorService.ListChecks(r.Context(), monitorID, limit, offset)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	data := make([]checkResponse, len(checks))
	for i, c := range checks {
		data[i] = toCheckResponse(c)
	}

	writeJSON(w, http.StatusOK, paginatedResponse{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func toCheckResponse(c *domain.Check) checkResponse {
	return checkResponse{
		ID:           c.ID,
		MonitorID:    c.MonitorID,
		Success:      c.Success,
		StatusCode:   c.StatusCode,
		LatencyMs:    c.LatencyMs,
		ErrorMessage: c.ErrorMessage,
		CheckedAt:    c.CheckedAt.Format(time.RFC3339),
	}
}
