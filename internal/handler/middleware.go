package handler

import (
	"log/slog"
	"net/http"
	"time"
)

// authMiddleware valida o header X-API-Key
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key != s.cfg.APIKey {
			writeError(w, http.StatusUnauthorized, "invalid API key", "UNAUTHORIZED")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware loga método, path, status e latência
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"latency_ms", time.Since(start).Milliseconds(),
		)
	})
}

// responseWriter captura o status code para o logger
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
