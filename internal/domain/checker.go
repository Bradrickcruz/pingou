package domain

import (
	"context"
	"time"
)

type CheckResult struct {
	Success      bool
	StatusCode   *int
	LatencyMs    int64
	ErrorMessage *string
	// CertNotAfter é a validade (UTC) do certificado do site na resposta final
	// de um check HTTPS. nil sem resposta TLS (http://, falha de conexão ou de
	// handshake, TCP, heartbeat).
	CertNotAfter *time.Time
	// CertURL é a URL do monitor usada no check que leu o certificado. A state
	// machine só grava a data se a URL ainda for a mesma.
	CertURL string
}

type Checker interface {
	Check(ctx context.Context, monitor *Monitor) CheckResult
}
