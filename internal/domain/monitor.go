package domain

import "time"

type (
	MonitorState   string
	MonitorEnabled bool
)

const (
	StateUnknown MonitorState = "UNKNOWN"
	StateUp      MonitorState = "UP"
	StateDown    MonitorState = "DOWN"
)

type Monitor struct {
	ID               string
	Name             string
	URL              string
	IntervalSeconds  int
	TimeoutSeconds   int
	FailureThreshold int
	Enabled          bool
	CurrentState     MonitorState
	LastCheckedAt    *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
