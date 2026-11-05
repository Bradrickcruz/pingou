package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/domain"
)

// ntfyTimeout cobre cada publicação. É variável só para os testes rodarem rápido.
var ntfyTimeout = 10 * time.Second

// ntfyBodyLimit: quanto do corpo da resposta é lido (e descartado).
const ntfyBodyLimit = 4 << 10

// ntfyTopicPattern é o que o ntfy aceita como nome de tópico.
var ntfyTopicPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

// validNtfyTopic diz se o tópico é aceito pelo ntfy.
func validNtfyTopic(topic string) bool { return ntfyTopicPattern.MatchString(topic) }

// ntfyMessage é uma publicação: texto no corpo e metadados nos headers.
type ntfyMessage struct {
	title    string
	body     string
	priority string  // "high", "default"...
	tags     string  // emojis do ntfy, separados por vírgula
	click    *string // link aberto ao tocar na notificação; nil = sem link
}

// NtfyNotifier publica no tópico salvo nas settings. Servidor e token vêm do
// ambiente (PINGOU_NTFY_*). O log nunca leva o token, a URL completa nem o
// tópico (em servidor aberto, o tópico funciona como senha).
type NtfyNotifier struct {
	cfg       config.NtfyConfig
	getTopic  func() string // vazio = desligado
	publicURL string        // PINGOU_PUBLIC_URL sem "/" final; vazio = sem link
	client    *http.Client
	location  func() *time.Location // fuso das mensagens; nil = UTC
}

func NewNtfyNotifier(cfg config.NtfyConfig, getTopic func() string, publicURL string) *NtfyNotifier {
	return &NtfyNotifier{cfg: cfg, getTopic: getTopic, publicURL: publicURL, client: newNtfyClient()}
}

// newNtfyClient não segue redirects: o POST viraria GET e não entregaria.
func newNtfyClient() *http.Client {
	return &http.Client{
		Timeout: ntfyTimeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// WithLocation define de onde vem o fuso das datas da mensagem (lido a cada envio).
func (n *NtfyNotifier) WithLocation(fn func() *time.Location) *NtfyNotifier {
	n.location = fn
	return n
}

func (n *NtfyNotifier) loc() *time.Location {
	if n.location == nil {
		return time.UTC
	}
	return n.location()
}

func (n *NtfyNotifier) Name() string { return domain.ChannelNtfy }

// Send monta a mensagem a partir do snapshot, na hora do envio, e publica no
// tópico salvo nas settings. Cert: prioridade normal no aviso de 30 dias e alta
// no de 7.
func (n *NtfyNotifier) Send(ctx context.Context, ev domain.OutboxEvent) error {
	topic := n.getTopic()
	if topic == "" {
		return nil
	}
	m := ev.Snapshot.AsMonitor()
	click := dashboardLink(n.publicURL, m)
	var msg ntfyMessage
	switch ev.EventType {
	case domain.EventDown:
		i := ev.Snapshot.AsIncident()
		msg = ntfyMessage{title: ntfyTitle(subjectDown(m)), body: bodyFor(m, i, false, n.publicURL, n.loc()),
			priority: "high", tags: "rotating_light", click: click}
	case domain.EventUp:
		i := ev.Snapshot.AsIncident()
		msg = ntfyMessage{title: ntfyTitle(subjectUp(m, i)), body: bodyFor(m, i, true, n.publicURL, n.loc()),
			priority: "default", tags: "white_check_mark", click: click}
	case domain.EventCertExpiring:
		a := ev.Snapshot.AsCertAlert()
		priority := "default"
		if a.ThresholdDays <= domain.CertAlertFinalDays {
			priority = "high"
		}
		msg = ntfyMessage{title: ntfyTitle(subjectCert(m, a)), body: certBody(m, a, n.publicURL, n.loc()),
			priority: priority, tags: "lock,warning", click: click}
	default:
		return fmt.Errorf("unknown event type %q", ev.EventType)
	}
	if _, err := publishNtfy(ctx, n.client, n.cfg, topic, msg); err != nil {
		return errors.New(ntfyErrorMessage(err))
	}
	return nil
}

// ntfyTitle tira o "[Pingou]" e o emoji do assunto do e-mail: no ntfy o app já
// mostra a origem, e o emoji vai em Tags (header não leva emoji cru).
func ntfyTitle(subject string) string {
	s := strings.TrimPrefix(subject, "[Pingou] ")
	for _, emoji := range []string{"🔴 ", "🟢 ", "⚠️ "} {
		s = strings.TrimPrefix(s, emoji)
	}
	return s
}

// errNtfyStatus é a resposta fora de 2xx do servidor.
type errNtfyStatus struct{ code int }

func (e errNtfyStatus) Error() string {
	if e.code == http.StatusUnauthorized || e.code == http.StatusForbidden {
		return fmt.Sprintf("unexpected status code: %d (token inválido ou sem permissão no tópico)", e.code)
	}
	return fmt.Sprintf("unexpected status code: %d", e.code)
}

// publishNtfy faz POST <URL>/<tópico> com o texto no corpo. Devolve o status
// HTTP (0 se não houve resposta).
func publishNtfy(ctx context.Context, client *http.Client, cfg config.NtfyConfig, topic string, msg ntfyMessage) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.URL+"/"+url.PathEscape(topic), strings.NewReader(msg.body))
	if err != nil {
		return 0, errors.New("invalid ntfy url")
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// oneLine impede injeção de header; a codificação RFC 2047 (que o ntfy
	// entende) leva acentos sem mandar bytes fora do ASCII no header.
	req.Header.Set("Title", mime.BEncoding.Encode("utf-8", oneLine(msg.title)))
	if msg.priority != "" {
		req.Header.Set("Priority", msg.priority)
	}
	if msg.tags != "" {
		req.Header.Set("Tags", msg.tags)
	}
	if msg.click != nil {
		req.Header.Set("Click", oneLine(*msg.click))
	}
	if cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, ntfyBodyLimit))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp.StatusCode, errNtfyStatus{resp.StatusCode}
	}
	return resp.StatusCode, nil
}

// ntfyErrorMessage tira o prefixo `Post "<url>": ` do erro, para a URL do
// servidor e o tópico não voltarem na resposta nem no log.
func ntfyErrorMessage(err error) string {
	var ue *url.Error
	if errors.As(err, &ue) {
		if ue.Timeout() {
			return fmt.Sprintf("timeout após %s", ntfyTimeout)
		}
		return ue.Err.Error()
	}
	return err.Error()
}

// ── teste ─────────────────────────────────────────────────────

// NtfyTestResult é a resposta de POST /api/settings/ntfy/test.
type NtfyTestResult struct {
	OK         bool    `json:"ok"`
	StatusCode *int    `json:"status_code"`
	DurationMS int64   `json:"duration_ms"`
	Error      *string `json:"error"`
}

// SendTestNtfy publica uma mensagem de teste no tópico. Não grava nada no banco.
func SendTestNtfy(ctx context.Context, cfg config.NtfyConfig, topic string) NtfyTestResult {
	msg := ntfyMessage{
		title:    "Teste do Pingou",
		body:     "Esta é uma notificação de teste do Pingou.\n\nSe ela chegou, as notificações pelo ntfy estão funcionando.\n",
		priority: "default",
		tags:     "test_tube",
	}
	start := time.Now()
	status, err := publishNtfy(ctx, newNtfyClient(), cfg, topic, msg)
	result := NtfyTestResult{DurationMS: time.Since(start).Milliseconds()}
	if status != 0 {
		result.StatusCode = &status
	}
	if err != nil {
		errMsg := ntfyErrorMessage(err)
		slog.Info("ntfy test", "host", cfg.Host(), "status", status, "error", errMsg)
		result.Error = &errMsg
		return result
	}
	slog.Info("ntfy test", "host", cfg.Host(), "status", status)
	result.OK = true
	return result
}
