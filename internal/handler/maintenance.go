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

// updateMaintenanceRequest: campo ausente não muda. reason é RawMessage para
// separar "ausente" de null (null ou "" apagam o motivo).
type updateMaintenanceRequest struct {
	StartsAt  *string         `json:"starts_at"`
	EndsAt    *string         `json:"ends_at"`
	Reason    json.RawMessage `json:"reason"`
	MonitorID json.RawMessage `json:"monitor_id"` // não editável (D13)
}

// maintenanceListItem: janela da página de manutenções. state vem do relógio
// do servidor: "scheduled" ou "active" (D12).
type maintenanceListItem struct {
	ID          string  `json:"id"`
	MonitorID   string  `json:"monitor_id"`
	MonitorName string  `json:"monitor_name"`
	StartsAt    string  `json:"starts_at"`
	EndsAt      string  `json:"ends_at"`
	Reason      *string `json:"reason"`
	State       string  `json:"state"`
}

func toMaintenanceListItem(w *domain.MaintenanceWindow, monitorName string, now time.Time) maintenanceListItem {
	state := "scheduled"
	if w.Started(now) {
		state = "active"
	}
	return maintenanceListItem{
		ID:          w.ID,
		MonitorID:   w.MonitorID,
		MonitorName: monitorName,
		StartsAt:    w.StartsAt.UTC().Format(time.RFC3339),
		EndsAt:      w.EndsAt.UTC().Format(time.RFC3339),
		Reason:      w.Reason,
		State:       state,
	}
}

func (s *Server) handleListMaintenance(w http.ResponseWriter, r *http.Request) {
	views, err := s.maintenanceService.List(r.Context())
	if err != nil {
		writeMaintenanceError(w, err)
		return
	}
	now := s.maintenanceService.Now()
	data := make([]maintenanceListItem, len(views))
	for i, v := range views {
		data[i] = toMaintenanceListItem(&v.MaintenanceWindow, v.MonitorName, now)
	}
	writeJSON(w, http.StatusOK, map[string]any{"windows": data})
}

func (s *Server) handleUpdateMaintenance(w http.ResponseWriter, r *http.Request) {
	var req updateMaintenanceRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxMaintenanceBody)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "BAD_REQUEST")
		return
	}
	if req.MonitorID != nil {
		writeError(w, http.StatusBadRequest, "validation error: monitor_id cannot be changed; cancel and schedule another window", "VALIDATION_ERROR")
		return
	}
	var in service.UpdateMaintenanceInput
	for _, f := range []struct {
		raw *string
		dst **time.Time
	}{{req.StartsAt, &in.StartsAt}, {req.EndsAt, &in.EndsAt}} {
		if f.raw == nil {
			continue
		}
		t, err := time.Parse(time.RFC3339, *f.raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "validation error: starts_at and ends_at must be RFC 3339 dates", "VALIDATION_ERROR")
			return
		}
		*f.dst = &t
	}
	if req.Reason != nil {
		in.ReasonSet = true
		if string(req.Reason) != "null" {
			var reason string
			if err := json.Unmarshal(req.Reason, &reason); err != nil {
				writeError(w, http.StatusBadRequest, "validation error: reason must be a string", "VALIDATION_ERROR")
				return
			}
			in.Reason = &reason
		}
	}

	mw, err := s.maintenanceService.Update(r.Context(), r.PathValue("id"), in)
	if err != nil {
		writeMaintenanceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMaintenanceResponse(mw, s.maintenanceService.Now()))
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

// writeMaintenanceError: conflito vira 409 com os monitores em conflito,
// janela iniciada com horário alterado vira 409 MAINTENANCE_LOCKED; o resto
// segue writeServiceError.
func writeMaintenanceError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrMaintenanceLocked) {
		writeError(w, http.StatusConflict, err.Error(), "MAINTENANCE_LOCKED")
		return
	}
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
