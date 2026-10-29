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

func (s *IncidentService) List(ctx context.Context, filter domain.IncidentFilter) ([]*domain.Incident, int, error) {
	return s.incidents.FindAll(ctx, filter)
}
