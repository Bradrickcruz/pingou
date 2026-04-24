package domain

import "time"

type Incident struct {
	ID               string
	MonitorID        string
	StartedAt        time.Time
	EndedAt          *time.Time
	LastError        *string
	NotificationSent bool
	DurationSeconds  *int64
}
