package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/service"
)

type Server struct {
	cfg            *config.Config
	router         *http.ServeMux
	monitorService *service.MonitorService
}

func NewServer(cfg *config.Config, monitorService *service.MonitorService) *Server {
	s := &Server{
		cfg:            cfg,
		router:         http.NewServeMux(),
		monitorService: monitorService,
	}
	s.registerRoutes()
	return s
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      loggingMiddleware(s.router),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	slog.Info("server listening", "port", s.cfg.Port)
	return srv.ListenAndServe()
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("GET /healthz", s.handleHealthz)

	api := http.NewServeMux()
	api.HandleFunc("GET /monitors", s.handleListMonitors)
	api.HandleFunc("POST /monitors", s.handleCreateMonitor)
	api.HandleFunc("GET /monitors/{id}", s.handleGetMonitor)
	api.HandleFunc("PATCH /monitors/{id}", s.handleUpdateMonitor)
	api.HandleFunc("DELETE /monitors/{id}", s.handleDeleteMonitor)

	s.router.Handle("/api/", s.authMiddleware(http.StripPrefix("/api", api)))
}

func (s *Server) shutdown(ctx context.Context) error {
	return nil // expandir na fase de Docker
}
