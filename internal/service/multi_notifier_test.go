package service

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

type funcNotifier struct {
	fn func()
}

func (f funcNotifier) NotifyDown(context.Context, *domain.Monitor, *domain.Incident)     { f.fn() }
func (f funcNotifier) NotifyRecovery(context.Context, *domain.Monitor, *domain.Incident) { f.fn() }

func TestMultiNotifierDoesNotBlock(t *testing.T) {
	release := make(chan struct{})
	n := NewMultiNotifier(funcNotifier{func() { <-release }})
	defer close(release)

	start := time.Now()
	n.NotifyDown(t.Context(), &domain.Monitor{}, &domain.Incident{})
	if d := time.Since(start); d > 50*time.Millisecond {
		t.Fatalf("NotifyDown levou %s, want < 50ms", d)
	}
}

func TestMultiNotifierIsolatesFailures(t *testing.T) {
	var calls atomic.Int64
	n := NewMultiNotifier(
		funcNotifier{func() { panic("canal quebrado") }},
		funcNotifier{func() { calls.Add(1) }},
	)
	n.NotifyRecovery(t.Context(), &domain.Monitor{}, &domain.Incident{})
	if !n.Wait(time.Second) {
		t.Fatal("Wait expirou")
	}
	if calls.Load() != 1 {
		t.Fatalf("segundo canal recebeu %d chamadas, want 1", calls.Load())
	}
}

func TestMultiNotifierWait(t *testing.T) {
	fast := NewMultiNotifier(funcNotifier{func() {}})
	fast.NotifyDown(t.Context(), &domain.Monitor{}, &domain.Incident{})
	if !fast.Wait(time.Second) {
		t.Fatal("envio rápido: Wait = false")
	}

	release := make(chan struct{})
	defer close(release)
	stuck := NewMultiNotifier(funcNotifier{func() { <-release }})
	stuck.NotifyDown(t.Context(), &domain.Monitor{}, &domain.Incident{})
	if stuck.Wait(100 * time.Millisecond) {
		t.Fatal("envio preso: Wait = true")
	}
}

func TestMultiNotifierIgnoresCanceledContext(t *testing.T) {
	got := make(chan error, 1)
	n := NewMultiNotifier(ctxNotifier{got})
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	n.NotifyDown(ctx, &domain.Monitor{}, &domain.Incident{})
	if err := <-got; err != nil {
		t.Fatalf("ctx do envio cancelado: %v", err)
	}
}

type ctxNotifier struct{ got chan error }

func (c ctxNotifier) NotifyDown(ctx context.Context, _ *domain.Monitor, _ *domain.Incident) {
	c.got <- ctx.Err()
}

func (c ctxNotifier) NotifyRecovery(ctx context.Context, _ *domain.Monitor, _ *domain.Incident) {
	c.got <- ctx.Err()
}
