package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
)

type Server struct {
	cfg    *config.Config
	router *http.ServeMux
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		cfg:    cfg,
		router: http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server listening", "port", s.cfg.Port)
	return srv.ListenAndServe()
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("GET /healthz", s.handleHealthz)

	// rotas autenticadas — prefixo /api
	api := http.NewServeMux()
	api.HandleFunc("GET /monitors", s.handleListMonitors)

	s.router.Handle("/api/", s.authMiddleware(http.StripPrefix("/api", api)))
}
