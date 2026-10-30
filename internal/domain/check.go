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

// StatsBucket agrega os checks de um intervalo fixo da janela de estatísticas.
// AvgLatencyMs considera só checks com sucesso; nil quando não há nenhum.
type StatsBucket struct {
	Start        time.Time
	Checks       int
	Failures     int
	Successes    int
	SumLatencyMs int64
	AvgLatencyMs *float64
}

// DownWindow é um incidente recortado aos limites da janela.
type DownWindow struct {
	Start time.Time
	End   time.Time
}

// MonitorStats resume um monitor numa janela (24h, 7d, 14d, 30d).
type MonitorStats struct {
	Range           string
	Since, Until    time.Time
	BucketSeconds   int
	Uptime          *float64 // nil sem checks na janela
	AvgLatencyMs    *float64 // nil sem checks com sucesso
	Checks          int
	Incidents       int
	DowntimeSeconds int64
	Series          []StatsBucket
	Down            []DownWindow
}
