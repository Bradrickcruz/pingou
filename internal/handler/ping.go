package handler

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Bradrickcruz/pingou/internal/service"
)

const maxPingBody = 1 << 10 // o corpo é ignorado; só lê um pouco para liberar a conexão

// handlePing recebe o ping de um heartbeat (GET, POST ou HEAD).
// 200: token válido (mesmo pausado). 404: token desconhecido, sem detalhes.
func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost, http.MethodHead:
	default:
		w.Header().Set("Allow", "GET, POST, HEAD")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	io.Copy(io.Discard, http.MaxBytesReader(w, r.Body, maxPingBody))

	// grava mesmo se o cliente desistir da resposta (ex.: curl -m curto)
	err := s.monitorService.Ping(context.WithoutCancel(r.Context()), r.PathValue("token"))
	switch {
	case errors.Is(err, service.ErrNotFound):
		http.NotFound(w, r)
		return
	case err != nil:
		slog.Error("ping error", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error", "INTERNAL_ERROR")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		io.WriteString(w, "OK\n")
	}
}
