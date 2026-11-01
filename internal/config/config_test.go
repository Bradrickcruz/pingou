package config

import (
	"strings"
	"testing"
)

func TestLoadSMTP(t *testing.T) {
	vars := []string{"HOST", "PORT", "USERNAME", "PASSWORD", "FROM", "TLS"}
	from := "Pingou <alerta@exemplo.com>"

	cases := []struct {
		name    string
		env     map[string]string
		want    SMTPConfig
		wantErr string
	}{
		{name: "vazio desliga", env: map[string]string{"PORT": "abc", "TLS": "foo"}, want: SMTPConfig{}},
		{
			name: "completo",
			env:  map[string]string{"HOST": "smtp.x.com", "USERNAME": "u", "PASSWORD": "p", "FROM": from},
			want: SMTPConfig{Host: "smtp.x.com", Port: 587, Username: "u", Password: "p", From: from, TLS: "starttls"},
		},
		{
			name: "465 vira tls",
			env:  map[string]string{"HOST": "smtp.x.com", "PORT": "465", "USERNAME": "u", "PASSWORD": "p", "FROM": from},
			want: SMTPConfig{Host: "smtp.x.com", Port: 465, Username: "u", Password: "p", From: from, TLS: "tls"},
		},
		{
			name: "relay none",
			env:  map[string]string{"HOST": "localhost", "PORT": "1025", "TLS": "NONE", "FROM": from},
			want: SMTPConfig{Host: "localhost", Port: 1025, From: from, TLS: "none"},
		},
		{name: "sem from", env: map[string]string{"HOST": "x"}, wantErr: "PINGOU_SMTP_FROM is required when PINGOU_SMTP_HOST is set"},
		{name: "from inválido", env: map[string]string{"HOST": "x", "FROM": "não é email"}, wantErr: "PINGOU_SMTP_FROM must be a valid address"},
		{name: "porta abc", env: map[string]string{"HOST": "x", "FROM": from, "PORT": "abc"}, wantErr: "PINGOU_SMTP_PORT"},
		{name: "porta 70000", env: map[string]string{"HOST": "x", "FROM": from, "PORT": "70000"}, wantErr: "PINGOU_SMTP_PORT"},
		{name: "tls foo", env: map[string]string{"HOST": "x", "FROM": from, "TLS": "foo"}, wantErr: "PINGOU_SMTP_TLS"},
		{name: "usuário sem senha", env: map[string]string{"HOST": "x", "FROM": from, "USERNAME": "u"}, wantErr: "PINGOU_SMTP_PASSWORD is required"},
		{name: "none com usuário", env: map[string]string{"HOST": "x", "FROM": from, "USERNAME": "u", "PASSWORD": "p", "TLS": "none"}, wantErr: "PINGOU_SMTP_TLS=none"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, v := range vars {
				t.Setenv("PINGOU_SMTP_"+v, c.env[v])
			}
			got, err := LoadSMTP()
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("err = %v, want %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Fatalf("got %+v, want %+v", got, c.want)
			}
			if got.Enabled() != (c.want.Host != "") {
				t.Fatalf("Enabled = %v", got.Enabled())
			}
		})
	}
}

func TestLoadFailsOnInvalidSMTP(t *testing.T) {
	t.Setenv("PINGOU_API_KEY", "k")
	t.Setenv("PINGOU_SMTP_HOST", "x")
	t.Setenv("PINGOU_SMTP_FROM", "")
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "PINGOU_SMTP_FROM is required") {
		t.Fatalf("Load: %v", err)
	}
}

func TestLoadNtfy(t *testing.T) {
	cases := []struct {
		name, url, token string
		want             NtfyConfig
		wantErr          string
	}{
		{name: "vazio desliga"},
		{name: "token sem url é ignorado", token: "tk_x"},
		{name: "url válida", url: "https://ntfy.sh", want: NtfyConfig{URL: "https://ntfy.sh"}},
		{name: "barra final", url: " http://ntfy.local:8080/ ", token: " tk_x ", want: NtfyConfig{URL: "http://ntfy.local:8080", Token: "tk_x"}},
		{name: "com caminho", url: "https://x.com/ntfy/", want: NtfyConfig{URL: "https://x.com/ntfy"}},
		{name: "sem esquema", url: "ntfy.sh", wantErr: "PINGOU_NTFY_URL must be a valid"},
		{name: "ftp", url: "ftp://ntfy.sh", wantErr: "PINGOU_NTFY_URL must be a valid"},
		{name: "sem host", url: "https://", wantErr: "PINGOU_NTFY_URL must be a valid"},
		{name: "com query", url: "https://ntfy.sh?auth=x", wantErr: "query"},
		{name: "com usuário", url: "https://u:p@ntfy.sh", wantErr: "user info"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("PINGOU_NTFY_URL", c.url)
			t.Setenv("PINGOU_NTFY_TOKEN", c.token)
			got, err := LoadNtfy()
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("err = %v, want %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Fatalf("got %+v, want %+v", got, c.want)
			}
			if got.Enabled() != (c.want.URL != "") {
				t.Fatalf("Enabled = %v", got.Enabled())
			}
		})
	}
}

func TestNtfyHost(t *testing.T) {
	if h := (NtfyConfig{URL: "https://ntfy.sh/base"}).Host(); h != "ntfy.sh" {
		t.Fatalf("Host = %q", h)
	}
}
