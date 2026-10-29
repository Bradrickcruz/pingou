package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/database"
	"github.com/Bradrickcruz/pingou/internal/domain"
	"github.com/Bradrickcruz/pingou/internal/repository"
	"github.com/Bradrickcruz/pingou/internal/scheduler"
	"github.com/Bradrickcruz/pingou/internal/service"
)

const testKey = "test-key"

type noopNotifier struct{}

func (noopNotifier) NotifyDown(context.Context, *domain.Monitor, *domain.Incident)     {}
func (noopNotifier) NotifyRecovery(context.Context, *domain.Monitor, *domain.Incident) {}

func newTestServer(t *testing.T) (*Server, http.Handler, func(id string) int) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "t.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	monitors := repository.NewMonitorRepo(db)
	checks := repository.NewCheckRepo(db)
	incidents := repository.NewIncidentRepo(db)
	svc := service.NewMonitorService(monitors, checks, incidents)
	// scheduler sem Start: só grava os pings (Record), sem jobs de check
	sched := scheduler.NewScheduler(monitors, nil, db,
		repository.NewCheckRepoTx(db), repository.NewMonitorRepoTx(db), repository.NewIncidentRepoTx(db),
		noopNotifier{})
	svc.SetRecorder(sched)

	s := NewServer(&config.Config{APIKey: testKey, DatabaseURL: dbPath, Version: "v1.2.3", Commit: "abc1234", BuildDate: "2026-10-29"}, db, svc,
		service.NewIncidentService(incidents), service.NewSettingsService(repository.NewSettingsRepo(db)))

	countChecks := func(id string) int {
		var n int
		if err := db.QueryRow(`SELECT COUNT(*) FROM checks WHERE monitor_id = ?`, id).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	return s, loggingMiddleware(s.router), countChecks
}

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("X-API-Key", testKey)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("json: %v (%s)", err, rec.Body.String())
	}
	return v
}

func TestCreateMonitorTypes(t *testing.T) {
	_, h, _ := newTestServer(t)

	// sem type → http
	rec := do(t, h, "POST", "/api/monitors", `{"name":"API","url":"https://x.test","interval_seconds":60,"timeout_seconds":5,"failure_threshold":3}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("http: %d %s", rec.Code, rec.Body)
	}
	httpMon := decode[monitorResponse](t, rec)
	if httpMon.Type != "http" || httpMon.URL == nil || !httpMon.Enabled || httpMon.PingToken != nil {
		t.Fatalf("http: %+v", httpMon)
	}

	rec = do(t, h, "POST", "/api/monitors", `{"name":"Postgres","type":"tcp","config":{"host":"db.interno","port":5432},"interval_seconds":60,"timeout_seconds":5,"failure_threshold":3}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("tcp: %d %s", rec.Code, rec.Body)
	}
	tcp := decode[monitorResponse](t, rec)
	if tcp.URL != nil || tcp.Config["host"] != "db.interno" || tcp.Config["port"] != float64(5432) {
		t.Fatalf("tcp: %+v", tcp)
	}

	rec = do(t, h, "POST", "/api/monitors", `{"name":"Backup","type":"heartbeat","interval_seconds":604800,"config":{"grace_seconds":3600}}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("heartbeat: %d %s", rec.Code, rec.Body)
	}
	hb := decode[monitorResponse](t, rec)
	if hb.PingToken == nil || hb.Config["grace_seconds"] != float64(3600) || hb.FailureThreshold != 1 {
		t.Fatalf("heartbeat: %+v", hb)
	}

	rec = do(t, h, "POST", "/api/monitors", `{"name":"x","type":"tcp","config":{"host":"db","port":0},"interval_seconds":60,"timeout_seconds":5,"failure_threshold":3}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("tcp porta 0: %d, want 400", rec.Code)
	}

	// tipo não pode mudar
	rec = do(t, h, "PATCH", "/api/monitors/"+tcp.ID, `{"type":"http"}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("trocar tipo: %d, want 400", rec.Code)
	}
}

func createHeartbeat(t *testing.T, h http.Handler) monitorResponse {
	t.Helper()
	rec := do(t, h, "POST", "/api/monitors", `{"name":"Backup","type":"heartbeat","interval_seconds":86400}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("heartbeat: %d %s", rec.Code, rec.Body)
	}
	return decode[monitorResponse](t, rec)
}

func ping(h http.Handler, method, token, ip string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, "/ping/"+token, strings.NewReader("ignorado"))
	req.RemoteAddr = ip + ":5555"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPingRoute(t *testing.T) {
	_, h, countChecks := newTestServer(t)
	hb := createHeartbeat(t, h)
	token := *hb.PingToken

	for _, method := range []string{"GET", "POST", "HEAD"} {
		if rec := ping(h, method, token, "10.0.0.1"); rec.Code != http.StatusOK {
			t.Fatalf("%s válido: %d, want 200", method, rec.Code)
		}
	}
	if n := countChecks(hb.ID); n != 3 {
		t.Fatalf("checks = %d, want 3", n)
	}
	got := decode[monitorResponse](t, do(t, h, "GET", "/api/monitors/"+hb.ID, ""))
	if got.CurrentState != "UP" || got.LastPingAt == nil {
		t.Fatalf("após ping: %+v", got)
	}

	if rec := ping(h, "GET", "nao-existe", "10.0.0.1"); rec.Code != http.StatusNotFound {
		t.Fatalf("token inválido: %d, want 404", rec.Code)
	}
	if rec := ping(h, "DELETE", token, "10.0.0.1"); rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("DELETE: %d, want 405", rec.Code)
	}

	// pausado: 200, sem registrar
	do(t, h, "PATCH", "/api/monitors/"+hb.ID, `{"enabled":false}`)
	if rec := ping(h, "GET", token, "10.0.0.1"); rec.Code != http.StatusOK {
		t.Fatalf("pausado: %d, want 200", rec.Code)
	}
	if n := countChecks(hb.ID); n != 3 {
		t.Fatalf("pausado registrou check: %d", n)
	}

	// ping não exige API Key
	req := httptest.NewRequest("GET", "/api/monitors", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("sanidade: /api sem chave = %d", rec.Code)
	}
}

func TestPingRateLimit(t *testing.T) {
	_, h, _ := newTestServer(t)

	// 404 também conta: limita varredura de token
	for i := range pingRateLimit {
		if rec := ping(h, "GET", fmt.Sprintf("x%d", i), "10.0.0.2"); rec.Code != http.StatusNotFound {
			t.Fatalf("req %d: %d, want 404", i, rec.Code)
		}
	}
	rec := ping(h, "GET", "x", "10.0.0.2")
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") == "" {
		t.Fatalf("excesso: %d Retry-After=%q, want 429 com Retry-After", rec.Code, rec.Header().Get("Retry-After"))
	}
	// outro IP não é afetado
	if rec := ping(h, "GET", "x", "10.0.0.3"); rec.Code != http.StatusNotFound {
		t.Fatalf("outro IP: %d, want 404", rec.Code)
	}
}

func TestRateLimiterWindowResets(t *testing.T) {
	now := time.Unix(0, 0)
	l := newRateLimiter(1, time.Minute)
	l.now = func() time.Time { return now }
	if ok, _ := l.allow("a"); !ok {
		t.Fatal("1ª requisição barrada")
	}
	if ok, _ := l.allow("a"); ok {
		t.Fatal("2ª requisição passou")
	}
	now = now.Add(time.Minute)
	if ok, _ := l.allow("a"); !ok {
		t.Fatal("janela nova continuou barrando")
	}
}

func TestRegenerateTokenEndpoint(t *testing.T) {
	_, h, _ := newTestServer(t)
	hb := createHeartbeat(t, h)

	rec := do(t, h, "POST", "/api/monitors/"+hb.ID+"/token", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("regenerar: %d %s", rec.Code, rec.Body)
	}
	updated := decode[monitorResponse](t, rec)
	if *updated.PingToken == *hb.PingToken {
		t.Fatal("token não mudou")
	}
	if rec := ping(h, "GET", *hb.PingToken, "10.0.0.4"); rec.Code != http.StatusNotFound {
		t.Fatalf("token antigo: %d, want 404", rec.Code)
	}
	if rec := ping(h, "GET", *updated.PingToken, "10.0.0.4"); rec.Code != http.StatusOK {
		t.Fatalf("token novo: %d, want 200", rec.Code)
	}
}

func TestListMonitorsSummary(t *testing.T) {
	s, h, _ := newTestServer(t)
	rec := do(t, h, "POST", "/api/monitors", `{"name":"API","url":"https://x.test","interval_seconds":60,"timeout_seconds":5,"failure_threshold":3,"enabled":false}`)
	m := decode[monitorResponse](t, rec)
	empty := decode[monitorResponse](t, do(t, h, "POST", "/api/monitors", `{"name":"Novo","url":"https://y.test","interval_seconds":60,"timeout_seconds":5,"failure_threshold":3,"enabled":false}`))

	// 35 checks, um por minuto; os 4 mais recentes falham; o último tem latência 99
	base := time.Now().UTC().Add(-35 * time.Minute)
	for i := range 35 {
		success := 1
		if i >= 31 {
			success = 0
		}
		latency := 10
		if i == 34 {
			latency = 99
		}
		if _, err := s.db.Exec(`INSERT INTO checks (monitor_id, success, latency_ms, checked_at) VALUES (?,?,?,?)`,
			m.ID, success, latency, base.Add(time.Duration(i)*time.Minute).Format(time.RFC3339)); err != nil {
			t.Fatal(err)
		}
	}

	page := decode[struct {
		Data []monitorListItem `json:"data"`
	}](t, do(t, h, "GET", "/api/monitors", ""))
	byID := map[string]monitorListItem{}
	for _, it := range page.Data {
		byID[it.ID] = it
	}

	got := byID[m.ID]
	if len(got.History) != 30 {
		t.Fatalf("history = %d itens, want 30", len(got.History))
	}
	if got.History[0] != true || got.History[29] != false || got.History[26] != false || got.History[25] != true {
		t.Fatalf("ordem do history errada (mais recente por último): %v", got.History)
	}
	if got.LastLatencyMs == nil || *got.LastLatencyMs != 99 {
		t.Fatalf("last_latency_ms = %v, want 99", got.LastLatencyMs)
	}
	if got.Uptime24h == nil || int(*got.Uptime24h) != 88 { // 31/35
		t.Fatalf("uptime_24h = %v, want ~88.57", got.Uptime24h)
	}

	if e := byID[empty.ID]; e.History == nil || len(e.History) != 0 || e.Uptime24h != nil || e.LastLatencyMs != nil {
		t.Fatalf("monitor sem checks: %+v", e)
	}
}

// O token é o segredo do heartbeat: não pode aparecer no log de acesso.
func TestPingTokenNotLogged(t *testing.T) {
	_, h, _ := newTestServer(t)
	hb := createHeartbeat(t, h)

	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	ping(h, "GET", *hb.PingToken, "10.0.0.5")
	ping(h, "GET", "token-inexistente", "10.0.0.5")

	if strings.Contains(buf.String(), *hb.PingToken) || strings.Contains(buf.String(), "token-inexistente") {
		t.Fatalf("token no log: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "/ping/[redacted]") {
		t.Fatalf("log sem o caminho redigido: %s", buf.String())
	}
}
