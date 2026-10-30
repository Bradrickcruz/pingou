package checker

import (
	"context"
	"errors"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// TCPChecker considera UP quando host:porta aceita a conexão dentro do timeout.
// A latência registrada é o tempo para conectar.
type TCPChecker struct {
	dialer net.Dialer
}

func NewTCPChecker() *TCPChecker {
	return &TCPChecker{}
}

func (c *TCPChecker) Check(ctx context.Context, m *domain.Monitor) domain.CheckResult {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(m.TimeoutSeconds)*time.Second)
	defer cancel()

	addr := net.JoinHostPort(m.Config.Host, strconv.Itoa(m.Config.Port))

	start := time.Now()
	conn, err := c.dialer.DialContext(ctx, "tcp", addr)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		msg := tcpErrorMessage(err)
		return domain.CheckResult{Success: false, LatencyMs: latency, ErrorMessage: &msg}
	}
	conn.Close()

	return domain.CheckResult{Success: true, LatencyMs: latency}
}

// tcpErrorMessage traduz o erro de conexão para um texto curto em last_error.
func tcpErrorMessage(err error) string {
	var dnsErr *net.DNSError
	switch {
	case errors.Is(err, syscall.ECONNREFUSED):
		return "connection refused"
	case errors.As(err, &dnsErr) && dnsErr.IsNotFound:
		return "no such host"
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, os.ErrDeadlineExceeded):
		return "timeout"
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout"
	}
	return err.Error()
}
