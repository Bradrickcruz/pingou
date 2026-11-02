package domain

import (
	"net"
	"strconv"
	"time"
)

type (
	MonitorState   string
	MonitorEnabled bool
	MonitorType    string
)

const (
	StateUnknown MonitorState = "UNKNOWN"
	StateUp      MonitorState = "UP"
	StateDown    MonitorState = "DOWN"
)

const (
	TypeHTTP      MonitorType = "http"
	TypeTCP       MonitorType = "tcp"
	TypeHeartbeat MonitorType = "heartbeat"
)

// MonitorConfig guarda os campos específicos de cada tipo (coluna config, JSON).
// http: nenhum (usa Monitor.URL); tcp: Host e Port; heartbeat: GraceSeconds.
type MonitorConfig struct {
	Host         string
	Port         int
	GraceSeconds int
}

type Monitor struct {
	ID               string
	Name             string
	Type             MonitorType
	URL              string // vazio fora de http
	Config           MonitorConfig
	PingToken        *string // só heartbeat
	IntervalSeconds  int
	TimeoutSeconds   int
	FailureThreshold int
	Enabled          bool
	CurrentState     MonitorState
	LastCheckedAt    *time.Time
	LastPingAt       *time.Time // só heartbeat
	// Certificado SSL (só http com URL https). CertAlertEnabled é a opção de
	// aviso; CertExpiresAt é a validade do último certificado visto (gravada
	// mesmo com o aviso desligado); CertAlerted*For é a data de expiração para
	// a qual o aviso de 30/7 dias já saiu.
	CertAlertEnabled bool
	CertExpiresAt    *time.Time
	CertAlerted30For *time.Time
	CertAlerted7For  *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Target é o alvo do check: URL (http), host:porta (tcp) ou nil (heartbeat,
// que não tem alvo: quem chama é o job monitorado).
func (m *Monitor) Target() *string {
	switch m.Type {
	case TypeHTTP:
		url := m.URL
		return &url
	case TypeTCP:
		target := net.JoinHostPort(m.Config.Host, strconv.Itoa(m.Config.Port))
		return &target
	}
	return nil
}

// MonitorSummary resume os checks recentes de um monitor para o painel.
type MonitorSummary struct {
	History       []bool   // últimos checks, mais antigo primeiro
	Uptime24h     *float64 // nil sem checks nas últimas 24 h
	LastLatencyMs *int64
	LastError     *string    // só DOWN, do incidente aberto
	DownSince     *time.Time // só DOWN, início do incidente aberto
}
