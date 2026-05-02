package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker(maxRedirects, globalTimeout int) *HTTPChecker {
	if maxRedirects <= 0 {
		maxRedirects = 5
	}
	if globalTimeout <= 0 {
		globalTimeout = 60
	}
	return &HTTPChecker{
		client: &http.Client{
			Timeout: time.Duration(globalTimeout) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirects {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

func (c *HTTPChecker) Check(ctx context.Context, m *domain.Monitor) domain.CheckResult {
	// timeout por monitor sobrescreve o contexto
	ctx, cancel := context.WithTimeout(ctx, time.Duration(m.TimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		msg := fmt.Sprintf("failed to build request: %s", err)
		return domain.CheckResult{Success: false, ErrorMessage: &msg}
	}
	req.Header.Set("User-Agent", "Pingou/1.0")

	start := time.Now()
	resp, err := c.client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		msg := err.Error()
		return domain.CheckResult{Success: false, LatencyMs: latency, ErrorMessage: &msg}
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	code := resp.StatusCode

	var errMsg *string
	if !success {
		s := fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		errMsg = &s
	}

	return domain.CheckResult{
		Success:      success,
		StatusCode:   &code,
		LatencyMs:    latency,
		ErrorMessage: errMsg,
	}
}
