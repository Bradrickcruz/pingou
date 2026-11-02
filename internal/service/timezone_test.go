package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

func TestNormalizeTimezone(t *testing.T) {
	for in, want := range map[string]string{
		"":                    "UTC",
		"   ":                 "UTC",
		"UTC":                 "UTC",
		" America/Sao_Paulo ": "America/Sao_Paulo",
		"Asia/Tokyo":          "Asia/Tokyo",
	} {
		got, err := normalizeTimezone(in)
		if err != nil || got != want {
			t.Fatalf("normalizeTimezone(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"Marte/Olimpo", "Local", "../../etc/passwd", strings.Repeat("a", 65)} {
		if _, err := normalizeTimezone(in); !errors.Is(err, ErrValidation) {
			t.Fatalf("normalizeTimezone(%q) err = %v, want ErrValidation", in, err)
		}
	}
}

func TestLocationLoaderFallsBackToUTC(t *testing.T) {
	cases := map[string]func(context.Context, string) (string, error){
		"vazio":    func(context.Context, string) (string, error) { return "", nil },
		"inválido": func(context.Context, string) (string, error) { return "Marte/Olimpo", nil },
		"erro":     func(context.Context, string) (string, error) { return "", errors.New("db fechado") },
	}
	for name, get := range cases {
		if loc := LocationLoader(get)(); loc != time.UTC {
			t.Fatalf("%s: loc = %v, want UTC", name, loc)
		}
	}
	loc := LocationLoader(func(context.Context, string) (string, error) { return "Asia/Tokyo", nil })()
	if loc.String() != "Asia/Tokyo" {
		t.Fatalf("loc = %v, want Asia/Tokyo", loc)
	}
}

// 14:03:12 UTC aparece como 11:03:12 -03 em São Paulo; "Fora do ar por" não muda.
func TestBodyForUsesTimezone(t *testing.T) {
	sp, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		t.Fatal(err)
	}
	started := time.Date(2026, 11, 1, 14, 3, 12, 0, time.UTC)
	ended := started.Add(5*time.Minute + 12*time.Second)
	m := &domain.Monitor{ID: "abc", Name: "API", Type: domain.TypeHTTP, URL: "https://api.exemplo.com"}

	body := bodyFor(m, &domain.Incident{StartedAt: started, EndedAt: &ended}, true, "", sp)
	for _, want := range []string{"Caiu em: 2026-11-01 11:03:12 -03", "Voltou em: 2026-11-01 11:08:24 -03", "Fora do ar por: 5 min 12 s"} {
		if !strings.Contains(body, want) {
			t.Fatalf("corpo sem %q:\n%s", want, body)
		}
	}

	// virada de dia: 23:30 UTC já é dia seguinte em Tóquio
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	late := time.Date(2026, 11, 1, 23, 30, 0, 0, time.UTC)
	if got := formatMessageTime(late, tokyo); got != "2026-11-02 08:30:00 JST" {
		t.Fatalf("Tóquio = %q", got)
	}
	if got := formatMessageTime(late, nil); got != "2026-11-01 23:30:00 UTC" {
		t.Fatalf("nil = %q", got)
	}
}
