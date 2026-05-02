package service

import (
	"fmt"
	"regexp"

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
	`(:[0-9]{1,4})?` +
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

func validateInterval(interval int) error {
	if interval < 10 || interval > 86400 {
		return fmt.Errorf("%w: interval_seconds must be between 10 and 86400", ErrValidation)
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

func validateMonitor(m *domain.Monitor) error {
	if err := validateName(m.Name); err != nil {
		return err
	}
	if err := validateURL(m.URL); err != nil {
		return err
	}
	if err := validateInterval(m.IntervalSeconds); err != nil {
		return err
	}
	if err := validateTimeout(m.TimeoutSeconds); err != nil {
		return err
	}
	if err := validateThreshold(m.FailureThreshold); err != nil {
		return err
	}
	return nil
}

func validateCreateInput(in CreateMonitorInput) error {
	if err := validateName(in.Name); err != nil {
		return err
	}
	if err := validateURL(in.URL); err != nil {
		return err
	}
	if err := validateInterval(in.IntervalSeconds); err != nil {
		return err
	}
	if err := validateTimeout(in.TimeoutSeconds); err != nil {
		return err
	}
	if err := validateThreshold(in.FailureThreshold); err != nil {
		return err
	}
	return nil
}

func validateUpdateInput(in UpdateMonitorInput) error {
	if in.Name != nil {
		if err := validateName(*in.Name); err != nil {
			return err
		}
	}
	if in.URL != nil {
		if err := validateURL(*in.URL); err != nil {
			return err
		}
	}
	if in.IntervalSeconds != nil {
		if err := validateInterval(*in.IntervalSeconds); err != nil {
			return err
		}
	}
	if in.TimeoutSeconds != nil {
		if err := validateTimeout(*in.TimeoutSeconds); err != nil {
			return err
		}
	}
	if in.FailureThreshold != nil {
		if err := validateThreshold(*in.FailureThreshold); err != nil {
			return err
		}
	}
	return nil
}
