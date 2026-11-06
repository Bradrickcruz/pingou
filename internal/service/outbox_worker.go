package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

const (
	// outboxBatch é quantos avisos vencidos o worker lê por vez.
	outboxBatch = 20
	// outboxPoll cobre os next_attempt_at futuros: sem wake, o worker acorda a cada 30 s.
	outboxPoll = 30 * time.Second
	// outboxSendTimeout limita cada envio (o SMTP chega a 30 s sozinho).
	outboxSendTimeout = 35 * time.Second
)

// outboxBackoff é a espera depois da 1ª, 2ª, 3ª e 4ª falha; da 4ª em diante
// repete a cada 30 min até o prazo do aviso.
var outboxBackoff = []time.Duration{30 * time.Second, 2 * time.Minute, 10 * time.Minute, 30 * time.Minute}

// OutboxStore é o que o worker usa do banco. Cada chamada solta a conexão ao
// voltar: o envio de rede acontece entre elas.
type OutboxStore interface {
	Due(ctx context.Context, now time.Time, limit int) ([]domain.OutboxEvent, error)
	MarkSent(ctx context.Context, id int64, attempts int, at time.Time) error
	MarkRetry(ctx context.Context, id int64, attempts int, next time.Time, lastError string) error
	MarkExpired(ctx context.Context, id int64, attempts int, at time.Time, lastError string) error
	CountPending(ctx context.Context) (int, error)
}

// OutboxWorker envia os avisos do outbox e repete os que falharam até dar
// certo ou expirar. Entrega "pelo menos uma vez": um crash entre enviar e
// marcar reenvia (o webhook leva o mesmo id para o receptor ignorar).
//
// Existe um único worker por banco (há lock de arquivo), então não há estado
// "sending": um crash nunca deixa linha presa.
// ponytail: envio sequencial; um webhook lento atrasa o e-mail da fila. Se virar
// problema, um worker por canal.
type OutboxWorker struct {
	store    OutboxStore
	channels map[string]Channel
	now      func() time.Time
	wake     chan struct{}
	done     chan struct{} // fechado quando o loop termina (nil antes do Start)
}

func NewOutboxWorker(store OutboxStore, channels ...Channel) *OutboxWorker {
	byName := make(map[string]Channel, len(channels))
	for _, c := range channels {
		byName[c.Name()] = c
	}
	return &OutboxWorker{store: store, channels: byName, now: time.Now, wake: make(chan struct{}, 1)}
}

// WithClock troca o relógio (testes simulam esperas sem dormir).
func (w *OutboxWorker) WithClock(now func() time.Time) *OutboxWorker {
	w.now = now
	return w
}

// Wake pede uma rodada agora. Nunca bloqueia: se já há um pedido na fila, basta.
func (w *OutboxWorker) Wake() {
	select {
	case w.wake <- struct{}{}:
	default:
	}
}

// Start roda o loop até o ctx ser cancelado. Retoma os pendentes de antes.
func (w *OutboxWorker) Start(ctx context.Context) {
	w.done = make(chan struct{})
	if n, err := w.store.CountPending(ctx); err != nil {
		slog.Error("outbox: could not count pending", "err", err)
	} else {
		slog.Info("outbox worker started", "pending", n)
	}
	go func() {
		defer close(w.done)
		w.run(ctx)
	}()
}

// Wait espera o loop terminar depois que o ctx de Start for cancelado. O envio
// em curso termina; o resto fica pending no banco.
func (w *OutboxWorker) Wait(ctx context.Context) error {
	if w.done == nil {
		return nil
	}
	select {
	case <-w.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *OutboxWorker) run(ctx context.Context) {
	ticker := time.NewTicker(outboxPoll)
	defer ticker.Stop()
	for {
		w.drain(ctx)
		select {
		case <-ctx.Done():
			slog.Info("outbox worker stopped")
			return
		case <-w.wake:
		case <-ticker.C:
		}
	}
}

// drain roda rodadas até não sobrar aviso vencido.
func (w *OutboxWorker) drain(ctx context.Context) {
	for ctx.Err() == nil {
		n, ok := w.runOnce(ctx)
		if !ok || n < outboxBatch {
			return
		}
	}
}

// runOnce processa um lote de avisos vencidos. Devolve quantos processou e
// false se o banco falhou (para o drain não girar em falso).
func (w *OutboxWorker) runOnce(ctx context.Context) (int, bool) {
	events, err := w.store.Due(ctx, w.now(), outboxBatch)
	if err != nil {
		if ctx.Err() == nil {
			slog.Error("outbox: could not read due notices", "err", err)
		}
		return 0, false
	}
	// o resultado de um envio que já saiu precisa ser gravado mesmo no
	// desligamento; senão o aviso sairia de novo na próxima subida
	save := context.WithoutCancel(ctx)
	for n, ev := range events {
		if ctx.Err() != nil {
			return n, true
		}
		if !w.process(ctx, save, ev) {
			return n + 1, false
		}
	}
	return len(events), true
}

// process trata uma linha. Devolve false se não conseguiu gravar o resultado.
func (w *OutboxWorker) process(ctx, save context.Context, ev domain.OutboxEvent) bool {
	now := w.now()
	if now.After(ev.ExpiresAt) {
		return w.expire(save, ev, ev.Attempts, now, "expired before sending")
	}
	ch, ok := w.channels[ev.Channel]
	if !ok {
		// ex.: SMTP desligado depois do aviso entrar na fila
		return w.expire(save, ev, ev.Attempts, now, "channel not available")
	}

	err := w.send(ctx, ch, ev)
	attempts := ev.Attempts + 1
	now = w.now()
	if err == nil {
		slog.Info("notice sent", "event_id", ev.EventID, "event_type", ev.EventType, "channel", ev.Channel, "attempt", attempts)
		return w.save(ev, w.store.MarkSent(save, ev.ID, attempts, now))
	}

	wait := outboxBackoff[min(attempts-1, len(outboxBackoff)-1)]
	next := now.Add(wait)
	if next.After(ev.ExpiresAt) {
		return w.expire(save, ev, attempts, now, err.Error())
	}
	slog.Warn("notice failed, will retry", "event_id", ev.EventID, "event_type", ev.EventType, "channel", ev.Channel,
		"attempt", attempts, "retry_in", wait.String(), "err", err)
	return w.save(ev, w.store.MarkRetry(save, ev.ID, attempts, next, err.Error()))
}

func (w *OutboxWorker) expire(save context.Context, ev domain.OutboxEvent, attempts int, now time.Time, reason string) bool {
	slog.Warn("notice expired", "event_id", ev.EventID, "event_type", ev.EventType, "channel", ev.Channel,
		"attempts", attempts, "reason", reason)
	return w.save(ev, w.store.MarkExpired(save, ev.ID, attempts, now, reason))
}

func (w *OutboxWorker) save(ev domain.OutboxEvent, err error) bool {
	if err != nil {
		slog.Error("outbox: could not save result", "event_id", ev.EventID, "channel", ev.Channel, "err", err)
		return false
	}
	return true
}

// send chama o canal com prazo próprio. Pânico vira erro (e retry), sem derrubar o worker.
// O prazo parte de um ctx que sobrevive ao desligamento: o envio em curso termina.
func (w *OutboxWorker) send(ctx context.Context, ch Channel, ev domain.OutboxEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("channel panic", "channel", ev.Channel, "event_type", ev.EventType, "panic", r)
			err = fmt.Errorf("channel panic")
		}
	}()
	sendCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), outboxSendTimeout)
	defer cancel()
	return ch.Send(sendCtx, ev)
}
