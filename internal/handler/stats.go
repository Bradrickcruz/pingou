package handler

import (
	"math"
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type statsBucketResponse struct {
	T            string `json:"t"`
	AvgLatencyMs *int64 `json:"avg_latency_ms"`
	Checks       int    `json:"checks"`
	Failures     int    `json:"failures"`
}

type downWindowResponse struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type statsResponse struct {
	Range           string                `json:"range"`
	Since           string                `json:"since"`
	Until           string                `json:"until"`
	BucketSeconds   int                   `json:"bucket_seconds"`
	Uptime          *float64              `json:"uptime"`
	AvgLatencyMs    *int64                `json:"avg_latency_ms"`
	Checks          int                   `json:"checks"`
	Incidents       int                   `json:"incidents"`
	DowntimeSeconds int64                 `json:"downtime_seconds"`
	Series          []statsBucketResponse `json:"series"`
	Down            []downWindowResponse  `json:"down"`
}

func (s *Server) handleMonitorStats(w http.ResponseWriter, r *http.Request) {
	rng := r.URL.Query().Get("range")
	if rng == "" {
		rng = "24h"
	}
	st, err := s.monitorService.Stats(r.Context(), r.PathValue("id"), rng, time.Now())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toStatsResponse(st))
}

func toStatsResponse(st *domain.MonitorStats) statsResponse {
	res := statsResponse{
		Range:           st.Range,
		Since:           st.Since.UTC().Format(time.RFC3339),
		Until:           st.Until.UTC().Format(time.RFC3339),
		BucketSeconds:   st.BucketSeconds,
		Uptime:          st.Uptime,
		AvgLatencyMs:    roundMs(st.AvgLatencyMs),
		Checks:          st.Checks,
		Incidents:       st.Incidents,
		DowntimeSeconds: st.DowntimeSeconds,
		Series:          make([]statsBucketResponse, len(st.Series)),
		Down:            make([]downWindowResponse, len(st.Down)),
	}
	for i, b := range st.Series {
		res.Series[i] = statsBucketResponse{
			T:            b.Start.UTC().Format(time.RFC3339),
			AvgLatencyMs: roundMs(b.AvgLatencyMs),
			Checks:       b.Checks,
			Failures:     b.Failures,
		}
	}
	for i, d := range st.Down {
		res.Down[i] = downWindowResponse{
			Start: d.Start.UTC().Format(time.RFC3339),
			End:   d.End.UTC().Format(time.RFC3339),
		}
	}
	return res
}

func roundMs(v *float64) *int64 {
	if v == nil {
		return nil
	}
	n := int64(math.Round(*v))
	return &n
}
