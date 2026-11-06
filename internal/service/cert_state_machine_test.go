package service_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/service"
)

const day = 24 * time.Hour

// certRecorder lê os avisos gravados no outbox (os pendentes, na ordem de gravação).
type certRecorder struct {
	t  *testing.T
	db *sql.DB
}

func (r *certRecorder) events(eventType string) []domain.OutboxEvent {
	return pendingNotices(r.t, r.db, eventType)
}

func (r *certRecorder) alerts() []domain.CertAlert {
	out := []domain.CertAlert{}
	for _, ev := range r.events(domain.EventCertExpiring) {
		out = append(out, ev.Snapshot.AsCertAlert())
	}
	return out
}

// at é o relógio simulado no momento de cada aviso.
func (r *certRecorder) at() []time.Time {
	out := []time.Time{}
	for _, ev := range r.events(domain.EventCertExpiring) {
		out = append(out, ev.CreatedAt)
	}
	return out
}

func (r *certRecorder) downs() int { return len(r.events(domain.EventDown)) }
func (r *certRecorder) up() int    { return len(r.events(domain.EventUp)) }

func (r *certRecorder) thresholds() []int {
	out := []int{}
	for _, a := range r.alerts() {
		out = append(out, a.ThresholdDays)
	}
	return out
}

// certEnv é um banco SQLite novo com um monitor https e um relógio simulado.
type certEnv struct {
	t     *testing.T
	db    *sql.DB
	repo  *repository.MonitorRepo
	now   time.Time
	rec   *certRecorder
	uowFn func() service.UnitOfWork
}

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func newCertEnv(t *testing.T, alertEnabled bool, threshold int) *certEnv {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	e := &certEnv{t: t, db: db, repo: repository.NewMonitorRepo(db), now: t0}
	e.rec = &certRecorder{t: t, db: db}
	e.uowFn = func() service.UnitOfWork {
		return newUoW(db)
	}
	m := &domain.Monitor{
		ID: "m1", Name: "API", Type: domain.TypeHTTP, URL: "https://api.test", IntervalSeconds: 3600, TimeoutSeconds: 5,
		FailureThreshold: threshold, Enabled: true, CurrentState: domain.StateUnknown, CertAlertEnabled: alertEnabled,
		CreatedAt: t0, UpdatedAt: t0,
	}
	if err := e.repo.Create(context.Background(), m); err != nil {
		t.Fatal(err)
	}
	return e
}

// monitor recarrega o monitor do banco, como o withMonitor do scheduler.
func (e *certEnv) monitor() *domain.Monitor {
	e.t.Helper()
	m, err := e.repo.FindByID(context.Background(), "m1")
	if err != nil || m == nil {
		e.t.Fatalf("monitor: %v %v", m, err)
	}
	return m
}

func (e *certEnv) process(res domain.CheckResult) error {
	e.t.Helper()
	sm := service.NewStateMachine(e.uowFn(), testChannels, nil).WithClock(func() time.Time { return e.now })
	return sm.Process(context.Background(), e.monitor(), res)
}

func (e *certEnv) mustProcess(res domain.CheckResult) {
	e.t.Helper()
	if err := e.process(res); err != nil {
		e.t.Fatal(err)
	}
}

// ok é um check com sucesso que viu um certificado vencendo em notAfter.
func ok(notAfter time.Time) domain.CheckResult {
	code := 200
	return domain.CheckResult{Success: true, StatusCode: &code, CertNotAfter: &notAfter, CertURL: "https://api.test"}
}

func fail() domain.CheckResult {
	msg := "connection refused"
	return domain.CheckResult{ErrorMessage: &msg}
}

// 90 dias de checks (1 por hora) com um certificado que vence no dia 45:
// exatamente 1 aviso de 30 dias e 1 de 7, na hora certa.
func TestCertAlertSimulation90Days(t *testing.T) {
	e := newCertEnv(t, true, 3)
	expires := t0.Add(45 * day)

	for e.now.Before(t0.Add(90 * day)) {
		if e.now.Before(expires) {
			e.mustProcess(ok(expires))
		} else {
			e.mustProcess(fail()) // vencido: o Go recusa o handshake, sem certificado
		}
		e.now = e.now.Add(time.Hour)
	}

	got := e.rec.thresholds()
	if len(got) != 2 || got[0] != 30 || got[1] != 7 {
		t.Fatalf("avisos %v, want [30 7]", got)
	}
	if want := expires.Add(-30 * day); !e.rec.at()[0].Equal(want) {
		t.Fatalf("aviso de 30 dias em %v, want %v", e.rec.at()[0], want)
	}
	if want := expires.Add(-7 * day); !e.rec.at()[1].Equal(want) {
		t.Fatalf("aviso de 7 dias em %v, want %v", e.rec.at()[1], want)
	}
	if a := e.rec.alerts()[0]; a.DaysLeft != 30 || !a.ExpiresAt.Equal(expires) {
		t.Fatalf("aviso de 30: %+v", a)
	}
	if m := e.monitor(); m.CertExpiresAt == nil || !m.CertExpiresAt.Equal(expires) {
		t.Fatalf("data perdida após checks sem certificado: %v", m.CertExpiresAt)
	}
	// quedas seguem o fluxo normal: 1 incidente quando o certificado venceu
	if e.rec.downs() != 1 {
		t.Fatalf("downs = %d, want 1", e.rec.downs())
	}
}

// Opção desligada: nenhum aviso, mas a data é gravada.
func TestCertAlertDisabledStillStoresDate(t *testing.T) {
	e := newCertEnv(t, false, 3)
	expires := t0.Add(3 * day)
	e.mustProcess(ok(expires))
	e.mustProcess(ok(expires))

	if n := len(e.rec.thresholds()); n != 0 {
		t.Fatalf("%d avisos com a opção desligada", n)
	}
	m := e.monitor()
	if m.CertExpiresAt == nil || !m.CertExpiresAt.Equal(expires) || m.CertAlerted30For != nil || m.CertAlerted7For != nil {
		t.Fatalf("monitor: expires=%v a30=%v a7=%v", m.CertExpiresAt, m.CertAlerted30For, m.CertAlerted7For)
	}
}

// Ligar a opção com 5 dias restantes: um aviso só, o de 7 dias.
func TestCertAlertLateActivation(t *testing.T) {
	e := newCertEnv(t, false, 3)
	expires := t0.Add(5 * day)
	e.mustProcess(ok(expires))

	m := e.monitor()
	m.CertAlertEnabled = true
	if err := e.repo.Update(context.Background(), m); err != nil {
		t.Fatal(err)
	}
	for range 10 {
		e.now = e.now.Add(time.Hour)
		e.mustProcess(ok(expires))
	}
	got := e.rec.thresholds()
	if len(got) != 1 || got[0] != 7 || e.rec.alerts()[0].DaysLeft != 4 {
		t.Fatalf("avisos %v (%+v), want um de 7 dias com 4 dias restantes", got, e.rec.alerts())
	}
}

// Renovação no meio do ciclo: os avisos voltam a funcionar para o certificado novo.
func TestCertAlertRenewal(t *testing.T) {
	e := newCertEnv(t, true, 3)
	first := t0.Add(20 * day)
	e.mustProcess(ok(first)) // aviso de 30 do primeiro

	e.now = t0.Add(5 * day)
	second := e.now.Add(25 * day) // renovado, mas curto: já entra nos 30 dias
	e.mustProcess(ok(second))     // aviso de 30 do segundo
	e.mustProcess(ok(second))     // não repete

	e.now = second.Add(-7 * day)
	e.mustProcess(ok(second)) // aviso de 7 do segundo

	got := e.rec.thresholds()
	if len(got) != 3 || got[0] != 30 || got[1] != 30 || got[2] != 7 {
		t.Fatalf("avisos %v, want [30 30 7]", got)
	}
	if !e.rec.alerts()[1].ExpiresAt.Equal(second) {
		t.Fatalf("segundo aviso com a data antiga: %v", e.rec.alerts()[1].ExpiresAt)
	}
}

// Troca de URL: campos zerados e o certificado do endereço antigo é descartado.
func TestCertAlertURLChange(t *testing.T) {
	e := newCertEnv(t, true, 3)
	expires := t0.Add(40 * day)
	e.mustProcess(ok(expires))

	m := e.monitor()
	m.URL = "https://novo.test"
	if err := e.repo.Update(context.Background(), m); err != nil {
		t.Fatal(err)
	}
	if m := e.monitor(); m.CertExpiresAt != nil {
		t.Fatalf("data não zerada ao trocar a URL: %v", m.CertExpiresAt)
	}

	// check que rodou com a URL antiga e terminou depois da edição
	e.mustProcess(ok(t0.Add(3 * day)))
	if m := e.monitor(); m.CertExpiresAt != nil || m.CertAlerted7For != nil {
		t.Fatalf("certificado da URL antiga gravado: %v %v", m.CertExpiresAt, m.CertAlerted7For)
	}
	if n := len(e.rec.thresholds()); n != 0 {
		t.Fatalf("%d avisos do certificado antigo", n)
	}
}

// Escrita guardada pela URL: se a URL mudou entre a leitura do monitor e a
// transação, nada é gravado nem enviado.
func TestCertAlertGuardedByURL(t *testing.T) {
	e := newCertEnv(t, true, 3)
	m := e.monitor() // lido antes da edição
	if _, err := e.db.Exec(`UPDATE monitors SET url='https://novo.test' WHERE id='m1'`); err != nil {
		t.Fatal(err)
	}
	sm := service.NewStateMachine(e.uowFn(), testChannels, nil).WithClock(func() time.Time { return e.now })
	if err := sm.Process(context.Background(), m, ok(t0.Add(3*day))); err != nil {
		t.Fatal(err)
	}
	if n := len(e.rec.thresholds()); n != 0 {
		t.Fatalf("%d avisos com a URL trocada", n)
	}
	if m := e.monitor(); m.CertExpiresAt != nil {
		t.Fatalf("data gravada com a URL trocada: %v", m.CertExpiresAt)
	}
}

// failingCommit desfaz a transação no lugar do commit.
type failingCommit struct{ service.UnitOfWork }

func (u failingCommit) Commit() error {
	_ = u.UnitOfWork.Rollback()
	return errors.New("commit falhou")
}

// Rollback: nenhum aviso sai e as marcas não ficam gravadas; o próximo check tenta de novo.
func TestCertAlertRollback(t *testing.T) {
	e := newCertEnv(t, true, 3)
	expires := t0.Add(10 * day)
	real := e.uowFn
	e.uowFn = func() service.UnitOfWork { return failingCommit{real()} }
	if err := e.process(ok(expires)); err == nil {
		t.Fatal("Process sem erro com commit falho")
	}
	if n := len(e.rec.thresholds()); n != 0 {
		t.Fatalf("%d avisos após rollback", n)
	}
	if m := e.monitor(); m.CertExpiresAt != nil || m.CertAlerted30For != nil {
		t.Fatalf("gravado após rollback: %v %v", m.CertExpiresAt, m.CertAlerted30For)
	}

	e.uowFn = real
	e.mustProcess(ok(expires))
	if got := e.rec.thresholds(); len(got) != 1 || got[0] != 30 {
		t.Fatalf("avisos após nova tentativa %v, want [30]", got)
	}
}

// Monitor em DOWN (status 500) com certificado visível ainda envia o aviso.
func TestCertAlertWhileDown(t *testing.T) {
	e := newCertEnv(t, true, 1)
	code := 500
	expires := t0.Add(6 * day)
	msg := "unexpected status code: 500"
	e.mustProcess(domain.CheckResult{StatusCode: &code, ErrorMessage: &msg, CertNotAfter: &expires, CertURL: "https://api.test"})

	if e.rec.downs() != 1 {
		t.Fatalf("downs = %d, want 1", e.rec.downs())
	}
	if got := e.rec.thresholds(); len(got) != 1 || got[0] != 7 {
		t.Fatalf("avisos %v, want [7]", got)
	}
	if m := e.monitor(); m.CurrentState != domain.StateDown {
		t.Fatalf("state = %s, want DOWN", m.CurrentState)
	}
}

// Regressão: UNKNOWN→UP, UP→DOWN e DOWN→UP avisam como antes, e checks sem
// mudança de certificado não escrevem nos campos de certificado.
func TestCertDoesNotChangeStateFlow(t *testing.T) {
	e := newCertEnv(t, true, 1)
	// conta cada UPDATE que grava as colunas de certificado
	if _, err := e.db.Exec(`CREATE TABLE cert_writes (n INTEGER);
		CREATE TRIGGER count_cert_writes AFTER UPDATE OF cert_expires_at, cert_alerted_30_for, cert_alerted_7_for ON monitors
		BEGIN INSERT INTO cert_writes VALUES (1); END;`); err != nil {
		t.Fatal(err)
	}
	expires := t0.Add(200 * day)
	e.mustProcess(ok(expires)) // UNKNOWN→UP: sem aviso, grava a data (1 escrita)
	e.mustProcess(fail())      // UP→DOWN
	e.mustProcess(ok(expires)) // DOWN→UP
	for range 20 {
		e.mustProcess(ok(expires))
	}
	if e.rec.downs() != 1 || e.rec.up() != 1 || len(e.rec.thresholds()) != 0 {
		t.Fatalf("downs=%d up=%d avisos=%v", e.rec.downs(), e.rec.up(), e.rec.thresholds())
	}
	var writes int
	if err := e.db.QueryRow(`SELECT COUNT(*) FROM cert_writes`).Scan(&writes); err != nil {
		t.Fatal(err)
	}
	if writes != 1 {
		t.Fatalf("%d escritas de certificado em 23 checks, want 1 (só a data nova)", writes)
	}
}
