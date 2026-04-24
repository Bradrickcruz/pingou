package service

import (
	"context"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type IncidentService struct {
	incidents domain.IncidentRepository
}

func NewIncidentService(incidents domain.IncidentRepository) *IncidentService {
	return &IncidentService{incidents: incidents}
}

func (s *IncidentService) List(ctx context.Context, onlyOpen bool, limit, offset int) ([]*domain.Incident, int, error) {
	return s.incidents.FindAll(ctx, onlyOpen, limit, offset)
}

func (s *IncidentService) ListByMonitor(ctx context.Context, monitorID string, onlyOpen bool, limit, offset int) ([]*domain.Incident, int, error) {
	return s.incidents.FindByMonitor(ctx, monitorID, onlyOpen, limit, offset)
}
