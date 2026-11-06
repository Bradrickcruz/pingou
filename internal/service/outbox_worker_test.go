package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
)

var owT0 = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

// owClock é um relógio que só anda quando o teste manda.
type owClock struct{ t time.Time }

func (c *owClock) now() time.Time { return c.t }

// owChannel é um canal falso: falha enquanto fail devolver erro.
type owChannel struct {
	name  string
	mu    sync.Mutex
	calls []domain.OutboxEvent
	fail  func(call int) error // call começa em 1; nil = sempre ok
	boom  bool                 // panic em vez de devolver erro
}

func (c *owChannel) Name() string { return c.name }

func (c *owChannel) Send(_ context.Context, ev domain.OutboxEvent) error {
	c.mu.Lock()
	c.calls = append(c.calls, ev)
	n := len(c.calls)
	c.mu.Unlock()
	if c.boom {
		panic("canal quebrado")
	}
	if c.fail != nil {
		return c.fail(n)
	}
	return nil
}

func (c *owChannel) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.calls)
}

func owAlways(msg string) func(int) error { return func(int) error { return errors.New(msg) } }

func owDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// owEnqueue grava um aviso por canal, como o StateMachine faz.
func owEnqueue(t *testing.T, db *sql.DB, eventType string, ttl time.Duration, channels ...string) {
	t.Helper()
	m := &domain.Monitor{ID: "m1", Name: "API", Type: domain.TypeHTTP, URL: "https://api.test"}
	snap := domain.NewDownSnapshot(m, &domain.Incident{StartedAt: owT0})
	switch eventType {
	case domain.EventUp:
		end := owT0.Add(time.Minute)
		snap = domain.NewRecoverySnapshot(m, &domain.Incident{StartedAt: owT0, EndedAt: &end}, end)
	case domain.EventCertExpiring:
		snap = domain.NewCertSnapshot(m, domain.CertAlert{ExpiresAt: owT0.Add(7 * 24 * time.Hour), DaysLeft: 7, ThresholdDays: 7}, owT0)
	}
	id, err := newUUIDv7()
	if err != nil {
		t.Fatal(err)
	}
	var rows []domain.OutboxEvent
	for _, ch := range channels {
		rows = append(rows, domain.OutboxEvent{
			EventID: id, EventType: eventType, Channel: ch, MonitorID: "m1", Snapshot: snap,
			CreatedAt: owT0, ExpiresAt: owT0.Add(ttl), NextAttemptAt: owT0,
		})
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.NewOutboxRepoTx().EnqueueWithTx(context.Background(), tx, rows); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

type owRow struct {
	status      string
	attempts    int
	next        time.Time
	done        sql.NullString
	lastError   sql.NullString
	channel, id string
}

func owRows(t *testing.T, db *sql.DB) []owRow {
	t.Helper()
	rows, err := db.Query(`SELECT channel, event_id, status, attempts, next_attempt_at, done_at, last_error FROM notification_outbox ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []owRow
	for rows.Next() {
		var r owRow
		var next string
		if err := rows.Scan(&r.channel, &r.id, &r.status, &r.attempts, &next, &r.done, &r.lastError); err != nil {
			t.Fatal(err)
		}
		r.next, _ = time.Parse(time.RFC3339, next)
		out = append(out, r)
	}
	return out
}

func owWorker(db *sql.DB, clock *owClock, channels ...Channel) *OutboxWorker {
	return NewOutboxWorker(repository.NewOutboxRepo(db), channels...).WithClock(clock.now)
}

// Queda e recuperação: esperas de 30 s, 2 min, 10 min e 30 min, e expira aos
// 60 min do evento (o retry que passaria do prazo não é agendado).
func TestWorkerBackoffAndExpiryForDown(t *testing.T) {
	db := owDB(t)
	clock := &owClock{t: owT0}
	ch := &owChannel{name: domain.ChannelWebhook, fail: owAlways("HTTP 500")}
	w := owWorker(db, clock, ch)
	owEnqueue(t, db, domain.EventDown, 60*time.Minute, domain.ChannelWebhook)

	var waits []time.Duration
	for range 10 {
		w.runOnce(t.Context())
		r := owRows(t, db)[0]
		if r.status != "pending" {
			break
		}
		waits = append(waits, r.next.Sub(clock.t))
		clock.t = r.next
	}
	want := []time.Duration{30 * time.Second, 2 * time.Minute, 10 * time.Minute, 30 * time.Minute}
	if len(waits) != len(want) {
		t.Fatalf("esperas = %v, want %v", waits, want)
	}
	for i := range want {
		if waits[i] != want[i] {
			t.Fatalf("esperas = %v, want %v", waits, want)
		}
	}
	r := owRows(t, db)[0]
	if r.status != "expired" || r.attempts != 5 || !r.done.Valid || r.lastError.String != "HTTP 500" {
		t.Fatalf("final: %+v, want expired com 5 tentativas e o último erro", r)
	}
	if ch.count() != 5 {
		t.Fatalf("envios = %d, want 5", ch.count())
	}
	if clock.t.Sub(owT0) > 60*time.Minute {
		t.Fatalf("tentou depois do prazo: %v", clock.t.Sub(owT0))
	}
}

// SSL: segue repetindo a cada 30 min depois da 4ª espera, até 24 h.
func TestWorkerCertNoticeLivesFor24Hours(t *testing.T) {
	db := owDB(t)
	clock := &owClock{t: owT0}
	ch := &owChannel{name: domain.ChannelWebhook, fail: owAlways("fora")}
	w := owWorker(db, clock, ch)
	owEnqueue(t, db, domain.EventCertExpiring, 24*time.Hour, domain.ChannelWebhook)

	for range 200 {
		w.runOnce(t.Context())
		r := owRows(t, db)[0]
		if r.status != "pending" {
			break
		}
		// depois dos 60 min de uma queda, o aviso de SSL continua vivo
		if clock.t.Sub(owT0) > 90*time.Minute && r.next.Sub(clock.t) != 30*time.Minute {
			t.Fatalf("espera %v em %v, want 30 min", r.next.Sub(clock.t), clock.t.Sub(owT0))
		}
		clock.t = r.next
	}
	r := owRows(t, db)[0]
	if r.status != "expired" {
		t.Fatalf("status = %s, want expired", r.status)
	}
	if last := clock.t.Sub(owT0); last <= 23*time.Hour || last > 24*time.Hour {
		t.Fatalf("última tentativa em %v, want entre 23 h e 24 h", last)
	}
}

// Só o canal que falhou repete; os outros não saem de novo.
func TestWorkerOnlyFailingChannelRetries(t *testing.T) {
	db := owDB(t)
	clock := &owClock{t: owT0}
	hook := &owChannel{name: domain.ChannelWebhook, fail: func(n int) error {
		if n == 1 {
			return errors.New("fora do ar")
		}
		return nil
	}}
	mail := &owChannel{name: domain.ChannelEmail}
	ntfy := &owChannel{name: domain.ChannelNtfy}
	w := owWorker(db, clock, hook, mail, ntfy)
	owEnqueue(t, db, domain.EventDown, time.Hour, domain.ChannelWebhook, domain.ChannelEmail, domain.ChannelNtfy)

	w.runOnce(t.Context())
	got := map[string]string{}
	for _, r := range owRows(t, db) {
		got[r.channel] = r.status
	}
	if got[domain.ChannelWebhook] != "pending" || got[domain.ChannelEmail] != "sent" || got[domain.ChannelNtfy] != "sent" {
		t.Fatalf("após a 1ª rodada: %v", got)
	}

	// antes do prazo do retry nada acontece
	clock.t = owT0.Add(29 * time.Second)
	if n, _ := w.runOnce(t.Context()); n != 0 {
		t.Fatalf("processou %d linhas antes do prazo", n)
	}
	clock.t = owT0.Add(30 * time.Second)
	w.runOnce(t.Context())
	for _, r := range owRows(t, db) {
		if r.status != "sent" {
			t.Fatalf("%s: %s, want sent", r.channel, r.status)
		}
	}
	if hook.count() != 2 || mail.count() != 1 || ntfy.count() != 1 {
		t.Fatalf("envios: webhook=%d email=%d ntfy=%d, want 2, 1 e 1", hook.count(), mail.count(), ntfy.count())
	}
	if hook.calls[0].EventID != hook.calls[1].EventID || hook.calls[0].EventID != mail.calls[0].EventID {
		t.Fatal("event_id mudou entre tentativas ou canais")
	}
}

// Canal que não existe mais (SMTP desligado depois) expira, sem pânico nem retry.
func TestWorkerUnknownChannelExpires(t *testing.T) {
	db := owDB(t)
	clock := &owClock{t: owT0}
	w := owWorker(db, clock, &owChannel{name: domain.ChannelWebhook})
	owEnqueue(t, db, domain.EventDown, time.Hour, domain.ChannelEmail)
	w.runOnce(t.Context())
	if r := owRows(t, db)[0]; r.status != "expired" || !strings.Contains(r.lastError.String, "channel") {
		t.Fatalf("%+v, want expired por canal indisponível", r)
	}
}

// Pânico no canal vira falha (e retry); o worker segue vivo.
func TestWorkerChannelPanicBecomesRetry(t *testing.T) {
	db := owDB(t)
	clock := &owClock{t: owT0}
	bad := &owChannel{name: domain.ChannelWebhook, boom: true}
	good := &owChannel{name: domain.ChannelEmail}
	w := owWorker(db, clock, bad, good)
	owEnqueue(t, db, domain.EventDown, time.Hour, domain.ChannelWebhook, domain.ChannelEmail)
	w.runOnce(t.Context())
	got := map[string]string{}
	for _, r := range owRows(t, db) {
		got[r.channel] = r.status
	}
	if got[domain.ChannelWebhook] != "pending" || got[domain.ChannelEmail] != "sent" {
		t.Fatalf("%v, want webhook pending e email sent", got)
	}
}

// Pingou ficou fora mais que o prazo: o aviso velho expira sem ser enviado.
func TestWorkerExpiresStaleNoticeWithoutSending(t *testing.T) {
	db := owDB(t)
	clock := &owClock{t: owT0.Add(2 * time.Hour)}
	ch := &owChannel{name: domain.ChannelWebhook}
	w := owWorker(db, clock, ch)
	owEnqueue(t, db, domain.EventDown, 60*time.Minute, domain.ChannelWebhook)
	w.runOnce(t.Context())
	if r := owRows(t, db)[0]; r.status != "expired" {
		t.Fatalf("status = %s, want expired", r.status)
	}
	if ch.count() != 0 {
		t.Fatalf("enviou um aviso vencido (%d envios)", ch.count())
	}
}

// Reinício: o que ficou pendente sai uma vez só, por um worker novo sobre o
// mesmo banco, sem reenviar o que já tinha saído.
func TestWorkerResumesPendingAfterRestart(t *testing.T) {
	db := owDB(t)
	clock := &owClock{t: owT0}
	hook1 := &owChannel{name: domain.ChannelWebhook, fail: owAlways("fora")}
	mail1 := &owChannel{name: domain.ChannelEmail}
	owEnqueue(t, db, domain.EventDown, time.Hour, domain.ChannelWebhook, domain.ChannelEmail)
	owWorker(db, clock, hook1, mail1).runOnce(t.Context()) // webhook falha, e-mail sai

	// "reinício": worker descartado, outro sobre o mesmo banco
	clock.t = owT0.Add(5 * time.Minute)
	hook2 := &owChannel{name: domain.ChannelWebhook}
	mail2 := &owChannel{name: domain.ChannelEmail}
	owWorker(db, clock, hook2, mail2).runOnce(t.Context())

	if hook2.count() != 1 || mail2.count() != 0 || mail1.count() != 1 {
		t.Fatalf("webhook novo=%d, email novo=%d, email antigo=%d, want 1, 0 e 1", hook2.count(), mail2.count(), mail1.count())
	}
	for _, r := range owRows(t, db) {
		if r.status != "sent" {
			t.Fatalf("%s: %s, want sent", r.channel, r.status)
		}
	}
	// segunda rodada: nada a fazer
	owWorker(db, clock, hook2, mail2).runOnce(t.Context())
	if hook2.count() != 1 {
		t.Fatal("reenviou o que já tinha saído")
	}
}

// Webhook de verdade: o receptor falha duas vezes (500) e depois responde 200.
// O aviso chega, o id é o mesmo em todas as tentativas e a linha termina sent.
func TestWorkerWebhookRetriesWithSameID(t *testing.T) {
	var mu sync.Mutex
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		b := make([]byte, 4096)
		n, _ := r.Body.Read(b)
		buf.Write(b[:n])
		mu.Lock()
		bodies = append(bodies, buf.String())
		call := len(bodies)
		mu.Unlock()
		if call <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	db := owDB(t)
	clock := &owClock{t: owT0}
	hook := NewWebhookNotifier(func() string { return srv.URL }, "")
	w := owWorker(db, clock, hook)
	owEnqueue(t, db, domain.EventDown, time.Hour, domain.ChannelWebhook)

	for range 6 {
		w.runOnce(t.Context())
		r := owRows(t, db)[0]
		if r.status != "pending" {
			break
		}
		clock.t = r.next
	}
	r := owRows(t, db)[0]
	if r.status != "sent" || r.attempts != 3 {
		t.Fatalf("%+v, want sent na 3ª tentativa", r)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(bodies) != 3 {
		t.Fatalf("receptor viu %d chamadas, want 3", len(bodies))
	}
	if !strings.Contains(bodies[0], `"id":"`+r.id+`"`) {
		t.Fatalf("payload sem o id do evento %s: %s", r.id, bodies[0])
	}
	for _, b := range bodies[1:] {
		if b != bodies[0] {
			t.Fatalf("payload mudou entre tentativas:\n%s\n%s", bodies[0], b)
		}
	}
	// o erro guardado nunca leva a URL
	if strings.Contains(r.lastError.String, srv.URL) {
		t.Fatalf("last_error com a URL: %s", r.lastError.String)
	}
}

// Wake nunca bloqueia, mesmo sem ninguém ouvindo.
func TestWorkerWakeNeverBlocks(t *testing.T) {
	w := NewOutboxWorker(nil)
	done := make(chan struct{})
	go func() {
		for range 100 {
			w.Wake()
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wake bloqueou")
	}
}

// Start/Wake/Wait de ponta a ponta com relógio real: o aviso sai logo depois do
// Wake, e Wait só volta depois que o loop para.
func TestWorkerStartWakeWait(t *testing.T) {
	db := owDB(t)
	ch := &owChannel{name: domain.ChannelWebhook}
	w := NewOutboxWorker(repository.NewOutboxRepo(db), ch)
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)

	// grava com o relógio real (o aviso nasce vencido) e acorda o worker
	m := &domain.Monitor{ID: "m1", Name: "API", Type: domain.TypeHTTP, URL: "https://api.test"}
	now := time.Now().UTC()
	snap := domain.NewDownSnapshot(m, &domain.Incident{StartedAt: now})
	tx, _ := db.Begin()
	err := repository.NewOutboxRepoTx().EnqueueWithTx(context.Background(), tx, []domain.OutboxEvent{{
		EventID: "e1", EventType: domain.EventDown, Channel: domain.ChannelWebhook, MonitorID: "m1", Snapshot: snap,
		CreatedAt: now, ExpiresAt: now.Add(time.Hour), NextAttemptAt: now,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	w.Wake()

	deadline := time.Now().Add(3 * time.Second)
	for ch.count() == 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if ch.count() != 1 {
		t.Fatalf("envios = %d, want 1", ch.count())
	}
	cancel()
	waitCtx, stop := context.WithTimeout(context.Background(), 2*time.Second)
	defer stop()
	if err := w.Wait(waitCtx); err != nil {
		t.Fatalf("Wait: %v", err)
	}
}
