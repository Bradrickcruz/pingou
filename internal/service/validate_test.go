package service

import "testing"

func TestValidateURL(t *testing.T) {
	valid := []string{
		"http://example.com",
		"https://example.com/health?x=1",
		"http://127.0.0.1:8080/",
		"http://localhost:65535/status", // porta com 5 dígitos
	}
	for _, u := range valid {
		if err := validateURL(u); err != nil {
			t.Errorf("validateURL(%q) = %v, want nil", u, err)
		}
	}

	invalid := []string{"", "example.com", "ftp://example.com", "http://example.com:123456"}
	for _, u := range invalid {
		if err := validateURL(u); err == nil {
			t.Errorf("validateURL(%q) = nil, want erro", u)
		}
	}
}
