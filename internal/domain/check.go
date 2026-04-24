package domain

import "time"

type Check struct {
	ID           int64
	MonitorID    string
	Success      bool
	StatusCode   *int
	LatencyMs    int64
	ErrorMessage *string
	CheckedAt    time.Time
}
