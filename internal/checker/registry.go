package checker

import (
	"context"
	"fmt"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// Registry implementa domain.Checker escolhendo o checker pelo tipo do monitor.
// O scheduler continua recebendo um único domain.Checker.
type Registry struct {
	checkers map[domain.MonitorType]domain.Checker
}

func NewRegistry(checkers map[domain.MonitorType]domain.Checker) *Registry {
	return &Registry{checkers: checkers}
}

func (r *Registry) Check(ctx context.Context, m *domain.Monitor) domain.CheckResult {
	c, ok := r.checkers[m.Type]
	if !ok {
		msg := fmt.Sprintf("no checker for monitor type %q", m.Type)
		return domain.CheckResult{Success: false, ErrorMessage: &msg}
	}
	return c.Check(ctx, m)
}
