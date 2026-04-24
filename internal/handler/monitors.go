package handler

import "net/http"

func (s *Server) handleListMonitors(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data":   []any{},
		"total":  0,
		"limit":  20,
		"offset": 0,
	})
}
