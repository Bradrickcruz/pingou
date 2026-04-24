package service

import (
	"fmt"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

func validate(m *domain.Monitor) error {
	if len(m.Name) < 1 || len(m.Name) > 100 {
		return fmt.Errorf("%w: name must be between 1 and 100 characters", ErrValidation)
	}
	if len(m.URL) < 1 || len(m.URL) > 2048 {
		return fmt.Errorf("%w: url must be between 1 and 2048 characters", ErrValidation)
	}
	if m.IntervalSeconds < 10 || m.IntervalSeconds > 86400 {
		return fmt.Errorf("%w: interval_seconds must be between 10 and 86400", ErrValidation)
	}
	if m.TimeoutSeconds < 1 || m.TimeoutSeconds > 60 {
		return fmt.Errorf("%w: timeout_seconds must be between 1 and 60", ErrValidation)
	}
	if m.FailureThreshold < 1 || m.FailureThreshold > 10 {
		return fmt.Errorf("%w: failure_threshold must be between 1 and 10", ErrValidation)
	}
	return nil
}
