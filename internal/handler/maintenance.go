package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/service"
)

// maxMaintenanceBody: 100 ids de monitor + datas + motivo cabem com folga.
const maxMaintenanceBody = 16 << 10

type createMaintenanceRequest struct {
	MonitorIDs []string `json:"monitor_ids"`
	StartsAt   string   `json:"starts_at"` // RFC 3339
	EndsAt     string   `json:"ends_at"`
	Reason     *string  `json:"reason"`
}

type maintenanceResponse struct {
	ID        string  `json:"id"`
	MonitorID string  `json:"monitor_id"`
	StartsAt  string  `json:"starts_at"`
	EndsAt    string  `json:"ends_at"`
	Reason    *string `json:"reason"`
	// Active é calculado no servidor, para não depender do relógio do navegador.
	Active bool `json:"active"`
}

func toMaintenanceResponse(w *domain.MaintenanceWindow, now time.Time) *maintenanceResponse {
	if w == nil {
		return nil
	}
	return &maintenanceResponse{
		ID:        w.ID,
		MonitorID: w.MonitorID,
		StartsAt:  w.StartsAt.UTC().Format(time.RFC3339),
		EndsAt:    w.EndsAt.UTC().Format(time.RFC3339),
		Reason:    w.Reason,
		Active:    w.Active(now),
	}
}

func (s *Server) handleCreateMaintenance(w http.ResponseWriter, r *http.Request) {
	var req createMaintenanceRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxMaintenanceBody)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "BAD_REQUEST")
		return
	}
	startsAt, err1 := time.Parse(time.RFC3339, req.StartsAt)
	endsAt, err2 := time.Parse(time.RFC3339, req.EndsAt)
	if err1 != nil || err2 != nil {
		writeError(w, http.StatusBadRequest, "validation error: starts_at and ends_at must be RFC 3339 dates", "VALIDATION_ERROR")
		return
	}

	windows, err := s.maintenanceService.CreateBatch(r.Context(), service.CreateMaintenanceInput{
		MonitorIDs: req.MonitorIDs,
		StartsAt:   startsAt,
		EndsAt:     endsAt,
		Reason:     req.Reason,
	})
	if err != nil {
		writeMaintenanceError(w, err)
		return
	}

	now := s.maintenanceService.Now()
	data := make([]*maintenanceResponse, len(windows))
	for i, mw := range windows {
		data[i] = toMaintenanceResponse(mw, now)
	}
	writeJSON(w, http.StatusCreated, map[string]any{"windows": data})
}

func (s *Server) handleCancelMaintenance(w http.ResponseWriter, r *http.Request) {
	if err := s.maintenanceService.Cancel(r.Context(), r.PathValue("id")); err != nil {
		writeMaintenanceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeMaintenanceError: conflito vira 409 com os monitores em conflito; o
// resto segue writeServiceError.
func writeMaintenanceError(w http.ResponseWriter, err error) {
	var conflict *service.MaintenanceConflictError
	if errors.As(err, &conflict) {
		writeJSON(w, http.StatusConflict, map[string]any{
			"error":       conflict.Error(),
			"code":        "MAINTENANCE_CONFLICT",
			"monitor_ids": conflict.MonitorIDs,
		})
		return
	}
	writeServiceError(w, err)
}

// liveMaintenance busca a janela viva dos monitores numa consulta (sem N+1).
// Sem o serviço (testes antigos) ou com erro, a resposta sai sem o bloco.
func (s *Server) liveMaintenance(ctx context.Context, monitors ...*domain.Monitor) (map[string]*domain.MaintenanceWindow, time.Time) {
	if s.maintenanceService == nil || len(monitors) == 0 {
		return nil, time.Time{}
	}
	ids := make([]string, len(monitors))
	for i, m := range monitors {
		ids[i] = m.ID
	}
	live, err := s.maintenanceService.LiveByMonitors(ctx, ids)
	if err != nil {
		slog.Error("maintenance windows", "err", err)
		return nil, time.Time{}
	}
	return live, s.maintenanceService.Now()
}
