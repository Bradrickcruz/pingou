package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// MultiNotifier repassa cada evento para todos os canais (webhook, e-mail), cada
// um na sua goroutine.
//
// Por que assíncrono: o StateMachine.Process chama o notifier dentro do
// withMonitor do scheduler, com o lock global dbMu seguro. Um destino lento
// (webhook com timeout de 10 s, SMTP com 30 s) pausaria a gravação de todos os
// monitores. Aqui o envio sai do lock, e a falha ou o panic de um canal não
// atinge o outro.
//
// m e i são seguros entre goroutines: o withMonitor recarrega o monitor do banco
// a cada chamada, e ninguém mais escreve nesses ponteiros depois do Process.
type MultiNotifier struct {
	notifiers []Notifier
	wg        sync.WaitGroup
}

func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Add inclui um canal. Use só na montagem, antes do scheduler começar.
func (n *MultiNotifier) Add(notifier Notifier) {
	n.notifiers = append(n.notifiers, notifier)
}

func (n *MultiNotifier) NotifyDown(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	n.fanOut(ctx, "down", func(ctx context.Context, c Notifier) { c.NotifyDown(ctx, m, i) })
}

func (n *MultiNotifier) NotifyRecovery(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	n.fanOut(ctx, "up", func(ctx context.Context, c Notifier) { c.NotifyRecovery(ctx, m, i) })
}

func (n *MultiNotifier) fanOut(ctx context.Context, event string, send func(context.Context, Notifier)) {
	// O ctx do job é cancelado quando o monitor é editado ou pausado; o aviso
	// já decidido precisa sair mesmo assim.
	ctx = context.WithoutCancel(ctx)
	for _, c := range n.notifiers {
		n.wg.Add(1)
		go func() {
			defer n.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					slog.Error("notifier panic", "event", event, "panic", r)
				}
			}()
			send(ctx, c)
		}()
	}
}

// Wait espera os envios em andamento. Devolve false se o timeout acabar antes.
func (n *MultiNotifier) Wait(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		n.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}
