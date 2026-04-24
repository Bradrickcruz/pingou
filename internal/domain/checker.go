package domain

import "context"

type CheckResult struct {
	Success      bool
	StatusCode   *int
	LatencyMs    int64
	ErrorMessage *string
}

type Checker interface {
	Check(ctx context.Context, monitor *Monitor) CheckResult
}
