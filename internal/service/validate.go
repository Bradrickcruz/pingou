package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

var urlRegexp = regexp.MustCompile(`^https?://` +
	`([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*` +
	`|[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}` +
	`|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}` +
	`|([0-9a-fA-F]{1,4}:){1,7}:|` +
	`|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}` +
	`|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}` +
	`|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}` +
	`|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}` +
	`|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}` +
	`|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})` +
	`|([0-9a-fA-F]{1,4}:){1,1}(:[0-9a-fA-F]{1,4}){1,7}` +
	`|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])` +
	`|::([0-9a-fA-F]{1,4}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])` +
	`|fe80:(:[0-9a-fA-F]{0,4}){0,1}%[0-9a-zA-Z]{1,}` +
	`|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])` +
	`|::([0-9a-fA-F]{1,4}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])` +
	`|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}` +
	`|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}` +
	`|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}` +
	`|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}` +
	`|([0-9a-fA-F]{1,4}:){1,6}(:[0-9a-fA-F]{1,4}){1,1}` +
	`|([0-9a-fA-F]{1,4}:){1,7}(:[0-9a-fA-F]{1,4})` +
	`|:` +
	`)` +
	`(:[0-9]{1,5})?` +
	`(/[a-zA-Z0-9._~:/?#\[\]@!$&'()*+,;=%-]*)?$`)

func validateURL(url string) error {
	if len(url) < 1 || len(url) > 2048 {
		return fmt.Errorf("%w: url must be between 1 and 2048 characters", ErrValidation)
	}
	if !urlRegexp.MatchString(url) {
		return fmt.Errorf("%w: invalid url format", ErrValidation)
	}
	return nil
}

func validateName(name string) error {
	if len(name) < 1 || len(name) > 100 {
		return fmt.Errorf("%w: name must be between 1 and 100 characters", ErrValidation)
	}
	return nil
}

// Limites de intervalo: http/tcp até 1 dia; heartbeat de 1 min até 1 semana.
const (
	maxCheckInterval     = 86400
	minHeartbeatInterval = 60
	maxHeartbeatInterval = 604800
	maxGraceSeconds      = 86400
	defaultGraceSeconds  = 60
)

func validateInterval(interval int) error {
	if interval < 10 || interval > maxCheckInterval {
		return fmt.Errorf("%w: interval_seconds must be between 10 and 86400", ErrValidation)
	}
	return nil
}

func validateHeartbeatInterval(interval int) error {
	if interval < minHeartbeatInterval || interval > maxHeartbeatInterval {
		return fmt.Errorf("%w: interval_seconds must be between 60 and 604800 for heartbeat", ErrValidation)
	}
	return nil
}

func validateTimeout(timeout int) error {
	if timeout < 1 || timeout > 60 {
		return fmt.Errorf("%w: timeout_seconds must be between 1 and 60", ErrValidation)
	}
	return nil
}

func validateThreshold(threshold int) error {
	if threshold < 1 || threshold > 10 {
		return fmt.Errorf("%w: failure_threshold must be between 1 and 10", ErrValidation)
	}
	return nil
}

func validateType(t domain.MonitorType) error {
	switch t {
	case domain.TypeHTTP, domain.TypeTCP, domain.TypeHeartbeat:
		return nil
	}
	return fmt.Errorf("%w: type must be http, tcp or heartbeat", ErrValidation)
}

// validateHost aceita nome ou IP, sem esquema, caminho ou porta embutida.
func validateHost(host string) error {
	if len(host) < 1 || len(host) > 253 {
		return fmt.Errorf("%w: host must be between 1 and 253 characters", ErrValidation)
	}
	if strings.Contains(host, "://") || strings.ContainsAny(host, "/ \t[]") {
		return fmt.Errorf("%w: host must be a hostname or IP, without scheme or path", ErrValidation)
	}
	return nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("%w: port must be between 1 and 65535", ErrValidation)
	}
	return nil
}

func validateGrace(grace int) error {
	if grace < 0 || grace > maxGraceSeconds {
		return fmt.Errorf("%w: grace_seconds must be between 0 and 86400", ErrValidation)
	}
	return nil
}

// validateMonitor valida o monitor completo, com as regras do seu tipo.
func validateMonitor(m *domain.Monitor) error {
	if err := validateName(m.Name); err != nil {
		return err
	}
	if err := validateType(m.Type); err != nil {
		return err
	}

	switch m.Type {
	case domain.TypeHTTP:
		if err := validateURL(m.URL); err != nil {
			return err
		}
	case domain.TypeTCP:
		if err := validateHost(m.Config.Host); err != nil {
			return err
		}
		if err := validatePort(m.Config.Port); err != nil {
			return err
		}
	case domain.TypeHeartbeat:
		if err := validateHeartbeatInterval(m.IntervalSeconds); err != nil {
			return err
		}
		return validateGrace(m.Config.GraceSeconds)
	}

	if err := validateInterval(m.IntervalSeconds); err != nil {
		return err
	}
	if err := validateTimeout(m.TimeoutSeconds); err != nil {
		return err
	}
	return validateThreshold(m.FailureThreshold)
}
