package scheduler

import (
	"context"
	"testing"
	"time"
)

type fakePurger struct {
	before time.Time
	calls  int
}

func (p *fakePurger) PurgeDoneBefore(_ context.Context, before time.Time) (int64, error) {
	p.calls++
	p.before = before
	return 0, nil
}

// A limpeza do outbox apaga o que foi encerrado há mais de 7 dias.
func TestRetentionPurgesOutboxAfterSevenDays(t *testing.T) {
	p := &fakePurger{}
	w := (&RetentionWorker{}).WithOutbox(p)
	w.purgeOutbox(context.Background())
	if p.calls != 1 {
		t.Fatalf("chamadas = %d, want 1", p.calls)
	}
	if age := time.Since(p.before); age < 7*24*time.Hour-time.Minute || age > 7*24*time.Hour+time.Minute {
		t.Fatalf("corte a %v, want 7 dias", age)
	}
}

// Sem outbox ligado, a rotina segue como antes.
func TestRetentionWithoutOutboxIsNoop(t *testing.T) {
	(&RetentionWorker{}).purgeOutbox(context.Background())
}
