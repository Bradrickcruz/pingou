package service

import (
	"testing"
	"time"
)

func TestEvaluateCertAlert(t *testing.T) {
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	day := 24 * time.Hour
	exp := func(d time.Duration) time.Time { return now.Add(d) }
	ptr := func(t time.Time) *time.Time { return &t }

	cases := []struct {
		name          string
		expiresAt     time.Time
		a30, a7       *time.Time
		wantThreshold int // 0 = sem aviso
		wantDays      int
		mark30, mark7 bool
	}{
		{name: "mais de 30 dias", expiresAt: exp(30*day + time.Second)},
		{name: "exatamente 30 dias", expiresAt: exp(30 * day), wantThreshold: 30, wantDays: 30, mark30: true},
		{name: "30 dias menos 1 s", expiresAt: exp(30*day - time.Second), wantThreshold: 30, wantDays: 29, mark30: true},
		{name: "entre 30 e 7, não avisado", expiresAt: exp(20 * day), wantThreshold: 30, wantDays: 20, mark30: true},
		{name: "entre 30 e 7, já avisado", expiresAt: exp(20 * day), a30: ptr(exp(20 * day))},
		{name: "exatamente 7 dias", expiresAt: exp(7 * day), a30: ptr(exp(7 * day)), wantThreshold: 7, wantDays: 7, mark30: true, mark7: true},
		{name: "7 dias e 1 s, já avisado de 30", expiresAt: exp(7*day + time.Second), a30: ptr(exp(7*day + time.Second))},
		{name: "ativação tardia com 5 dias", expiresAt: exp(5 * day), wantThreshold: 7, wantDays: 5, mark30: true, mark7: true},
		{name: "7 dias já avisado", expiresAt: exp(3 * day), a30: ptr(exp(3 * day)), a7: ptr(exp(3 * day))},
		{name: "já vencido", expiresAt: exp(-2 * time.Hour), wantThreshold: 7, wantDays: -1, mark30: true, mark7: true},
		{name: "renovado: marcas de outra data", expiresAt: exp(25 * day), a30: ptr(exp(-60 * day)), a7: ptr(exp(-60 * day)), wantThreshold: 30, wantDays: 25, mark30: true},
		{name: "trocado por um mais curto", expiresAt: exp(4 * day), a30: ptr(exp(20 * day)), a7: ptr(exp(20 * day)), wantThreshold: 7, wantDays: 4, mark30: true, mark7: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := EvaluateCertAlert(now, c.expiresAt, c.a30, c.a7)
			if c.wantThreshold == 0 {
				if d.Alert != nil || d.Mark30 || d.Mark7 {
					t.Fatalf("got %+v (alert %+v), want nada", d, d.Alert)
				}
				return
			}
			if d.Alert == nil {
				t.Fatalf("sem aviso, want limite %d", c.wantThreshold)
			}
			if d.Alert.ThresholdDays != c.wantThreshold || d.Alert.DaysLeft != c.wantDays || !d.Alert.ExpiresAt.Equal(c.expiresAt) {
				t.Fatalf("alert=%+v, want limite %d e %d dias", *d.Alert, c.wantThreshold, c.wantDays)
			}
			if d.Mark30 != c.mark30 || d.Mark7 != c.mark7 {
				t.Fatalf("marcas 30=%v 7=%v, want %v e %v", d.Mark30, d.Mark7, c.mark30, c.mark7)
			}
		})
	}
}

// 100 avaliações seguidas, aplicando as marcas como a state machine faz: um
// aviso só.
func TestEvaluateCertAlertRepeated(t *testing.T) {
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	expiresAt := now.Add(20 * 24 * time.Hour)
	var a30, a7 *time.Time
	alerts := 0
	for i := range 100 {
		d := EvaluateCertAlert(now.Add(time.Duration(i)*time.Minute), expiresAt, a30, a7)
		if d.Alert != nil {
			alerts++
		}
		if d.Mark30 {
			a30 = &expiresAt
		}
		if d.Mark7 {
			a7 = &expiresAt
		}
	}
	if alerts != 1 {
		t.Fatalf("%d avisos, want 1", alerts)
	}
}
