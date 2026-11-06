package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/domain"
)

// smtpTimeout cobre a conversa SMTP inteira (net/smtp não aceita context).
// É variável só para os testes rodarem rápido.
var (
	smtpTimeout     = 30 * time.Second
	smtpDialTimeout = 10 * time.Second
)

// EmailNotifier envia um e-mail em texto simples para os destinatários salvos
// nas settings. A config do servidor vem toda do ambiente (PINGOU_SMTP_*).
type EmailNotifier struct {
	cfg           config.SMTPConfig
	getRecipients func() string         // lista normalizada, separada por vírgula; vazio = desligado
	publicURL     string                // PINGOU_PUBLIC_URL sem "/" final; vazio = sem link
	location      func() *time.Location // fuso das mensagens; nil = UTC
}

func NewEmailNotifier(cfg config.SMTPConfig, getRecipients func() string, publicURL string) *EmailNotifier {
	return &EmailNotifier{cfg: cfg, getRecipients: getRecipients, publicURL: publicURL}
}

// WithLocation define de onde vem o fuso das datas da mensagem (lido a cada envio).
func (n *EmailNotifier) WithLocation(fn func() *time.Location) *EmailNotifier {
	n.location = fn
	return n
}

func (n *EmailNotifier) loc() *time.Location {
	if n.location == nil {
		return time.UTC
	}
	return n.location()
}

func (n *EmailNotifier) NotifyDown(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	n.send(ctx, eventMonitorDown, subjectDown(m), bodyFor(m, i, false, n.publicURL, n.loc()))
}

func (n *EmailNotifier) NotifyRecovery(ctx context.Context, m *domain.Monitor, i *domain.Incident) {
	n.send(ctx, eventMonitorUp, subjectUp(m, i), bodyFor(m, i, true, n.publicURL, n.loc()))
}

func (n *EmailNotifier) NotifyCertExpiring(ctx context.Context, m *domain.Monitor, a domain.CertAlert) {
	n.send(ctx, eventCertExpiring, subjectCert(m, a), certBody(m, a, n.publicURL, n.loc()))
}

func (n *EmailNotifier) send(ctx context.Context, event, subject, body string) {
	to := splitRecipients(n.getRecipients())
	if len(to) == 0 {
		return
	}
	msg := buildEmail(n.cfg.From, to, subject, body, time.Now())
	// o log nunca leva os endereços nem a senha
	if err := sendMail(ctx, n.cfg, to, msg); err != nil {
		slog.Error("email send error", "event", event, "recipients_count", len(to), "err", smtpErrorMessage(err))
		return
	}
	slog.Info("email sent", "event", event, "recipients_count", len(to))
}

// splitRecipients separa a lista salva (já validada e normalizada).
func splitRecipients(s string) []string {
	var out []string
	for _, r := range strings.Split(s, ",") {
		if r = strings.TrimSpace(r); r != "" {
			out = append(out, r)
		}
	}
	return out
}

// ── mensagem ──────────────────────────────────────────────────

// oneLine tira CR/LF de textos que vêm do usuário (nome do monitor, erro), para
// não virarem cabeçalho novo no e-mail.
func oneLine(s string) string {
	return strings.Join(strings.Fields(strings.NewReplacer("\r", " ", "\n", " ").Replace(s)), " ")
}

func subjectDown(m *domain.Monitor) string {
	return fmt.Sprintf("[Pingou] 🔴 %s está fora do ar", oneLine(m.Name))
}

func subjectUp(m *domain.Monitor, i *domain.Incident) string {
	d := downtimeSeconds(i, recoveredAt(i))
	return fmt.Sprintf("[Pingou] 🟢 %s voltou (fora por %s)", oneLine(m.Name), shortDuration(d))
}

const subjectTest = "[Pingou] Teste de e-mail"

// certPhrase: "expira em 7 dias", "expira em 1 dia", "expira hoje", "expirou".
func certPhrase(daysLeft int) string {
	switch {
	case daysLeft < 0:
		return "expirou"
	case daysLeft == 0:
		return "expira hoje"
	case daysLeft == 1:
		return "expira em 1 dia"
	}
	return fmt.Sprintf("expira em %d dias", daysLeft)
}

func subjectCert(m *domain.Monitor, a domain.CertAlert) string {
	return fmt.Sprintf("[Pingou] ⚠️ Certificado SSL de %s %s", oneLine(m.Name), certPhrase(a.DaysLeft))
}

// certBody monta o corpo do aviso de certificado (e-mail e ntfy), com a data
// no fuso loc.
func certBody(m *domain.Monitor, a domain.CertAlert, publicURL string, loc *time.Location) string {
	var b strings.Builder
	line := func(label, value string) { fmt.Fprintf(&b, "%s: %s\n", label, value) }

	line("Monitor", oneLine(m.Name))
	line("URL", oneLine(m.URL))
	line("Expira em", formatMessageTime(a.ExpiresAt, loc))
	switch {
	case a.DaysLeft < 0:
		line("Situação", "o certificado já expirou")
	case a.DaysLeft == 1:
		line("Faltam", "1 dia")
	default:
		line("Faltam", fmt.Sprintf("%d dias", a.DaysLeft))
	}
	line("Aviso", fmt.Sprintf("%d dias antes de expirar", a.ThresholdDays))
	if u := dashboardLink(publicURL, m); u != nil {
		line("Detalhes", *u)
	}
	b.WriteString("\nRenove o certificado antes dessa data, ou o site fica fora do ar para os visitantes.\n")
	return b.String()
}

// bodyFor monta o corpo no formato do README, com as datas no fuso loc. Nunca
// leva o ping_token.
func bodyFor(m *domain.Monitor, i *domain.Incident, recovered bool, publicURL string, loc *time.Location) string {
	var b strings.Builder
	line := func(label, value string) { fmt.Fprintf(&b, "%s: %s\n", label, value) }

	line("Monitor", oneLine(m.Name))
	line("Tipo", string(m.Type))
	if t := m.Target(); t != nil {
		line("Alvo", oneLine(*t))
	}
	line("Caiu em", formatMessageTime(i.StartedAt, loc))
	if recovered {
		end := recoveredAt(i)
		line("Voltou em", formatMessageTime(end, loc))
		line("Fora do ar por", formatDuration(downtimeSeconds(i, end)))
	} else if i.LastError != nil && *i.LastError != "" {
		line("Erro", oneLine(*i.LastError))
	}
	if u := dashboardLink(publicURL, m); u != nil {
		line("Detalhes", *u)
	}
	return b.String()
}

// formatDuration: "45 s", "5 min 12 s", "2 h 3 min", "1 d 4 h".
func formatDuration(sec int64) string {
	d, h, mi, s := sec/86400, sec%86400/3600, sec%3600/60, sec%60
	switch {
	case d > 0:
		return joinUnits(d, "d", h, "h")
	case h > 0:
		return joinUnits(h, "h", mi, "min")
	case mi > 0:
		return joinUnits(mi, "min", s, "s")
	}
	return fmt.Sprintf("%d s", s)
}

func joinUnits(a int64, ua string, b int64, ub string) string {
	if b == 0 {
		return fmt.Sprintf("%d %s", a, ua)
	}
	return fmt.Sprintf("%d %s %d %s", a, ua, b, ub)
}

// shortDuration: só a maior unidade, para o assunto ("5 min").
func shortDuration(sec int64) string {
	full := formatDuration(sec)
	parts := strings.Fields(full)
	if len(parts) > 2 {
		return parts[0] + " " + parts[1]
	}
	return full
}

// buildEmail monta a mensagem RFC 5322 em texto simples, UTF-8 e
// quoted-printable, com linhas terminadas em CRLF.
func buildEmail(from string, to []string, subject, body string, now time.Time) []byte {
	fromHeader, domainPart := from, "pingou.local"
	if a, err := mail.ParseAddress(from); err == nil {
		fromHeader = a.String() // codifica nomes com acento
		if at := strings.LastIndex(a.Address, "@"); at >= 0 {
			domainPart = a.Address[at+1:]
		}
	}
	msgID, err := newUUIDv7()
	if err != nil {
		msgID = strconv.FormatInt(now.UnixNano(), 36)
	}

	var buf bytes.Buffer
	header := func(k, v string) { fmt.Fprintf(&buf, "%s: %s\r\n", k, oneLine(v)) }
	header("From", fromHeader)
	header("To", strings.Join(to, ", "))
	header("Subject", mime.QEncoding.Encode("utf-8", oneLine(subject)))
	header("Date", now.Format(time.RFC1123Z))
	header("Message-ID", "<"+msgID+"@"+domainPart+">")
	header("MIME-Version", "1.0")
	header("Content-Type", "text/plain; charset=utf-8")
	header("Content-Transfer-Encoding", "quoted-printable")
	buf.WriteString("\r\n")

	qp := quotedprintable.NewWriter(&buf)
	_, _ = qp.Write([]byte(strings.ReplaceAll(body, "\r\n", "\n")))
	_ = qp.Close()
	return buf.Bytes()
}

// ── transporte ────────────────────────────────────────────────

var errStartTLSUnsupported = errors.New("server does not support STARTTLS")

// sendMail entrega msg usando só a stdlib. Modos:
//   - tls: TLS desde a conexão (465)
//   - starttls: conexão aberta e STARTTLS obrigatório; sem ele, erro (nunca
//     rebaixa para texto aberto)
//   - none: sem TLS e sem autenticação (relay local)
func sendMail(ctx context.Context, cfg config.SMTPConfig, to []string, msg []byte) error {
	from, err := mail.ParseAddress(cfg.From)
	if err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	tlsConfig := &tls.Config{ServerName: cfg.Host, MinVersion: tls.VersionTLS12}

	dialer := &net.Dialer{Timeout: smtpDialTimeout}
	var conn net.Conn
	if cfg.TLS == config.SMTPTLSImplicit {
		conn, err = (&tls.Dialer{NetDialer: dialer, Config: tlsConfig}).DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	// net/smtp não aceita context: o deadline cobre a conversa inteira, e o
	// cancelamento do ctx fecha a conexão.
	if err := conn.SetDeadline(time.Now().Add(smtpTimeout)); err != nil {
		return fmt.Errorf("set deadline: %w", err)
	}
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()

	c, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("greeting: %w", err)
	}
	defer c.Close()

	if cfg.TLS == config.SMTPTLSStartTLS {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return errStartTLSUnsupported
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	if cfg.Username != "" {
		if err := c.Auth(smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)); err != nil {
			return fmt.Errorf("auth: %w", err)
		}
	}

	if err := c.Mail(from.Address); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	for _, r := range to {
		if err := c.Rcpt(r); err != nil {
			return fmt.Errorf("rcpt to: %w", err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("data write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("data close: %w", err)
	}
	return c.Quit()
}

// smtpErrorMessage deixa o erro legível; timeout vira "timeout após 30s".
func smtpErrorMessage(err error) string {
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return fmt.Sprintf("timeout após %s", smtpTimeout)
	}
	return err.Error()
}

// ── teste ─────────────────────────────────────────────────────

// EmailTestResult é a resposta de POST /api/settings/email/test.
type EmailTestResult struct {
	OK         bool    `json:"ok"`
	DurationMS int64   `json:"duration_ms"`
	Error      *string `json:"error"`
}

// SendTestEmail envia um e-mail de teste para to. Não grava nada no banco.
func SendTestEmail(ctx context.Context, cfg config.SMTPConfig, to []string) EmailTestResult {
	body := "Este é um e-mail de teste do Pingou.\n\n" +
		"Se ele chegou, as notificações por e-mail estão funcionando.\n"
	start := time.Now()
	err := sendMail(ctx, cfg, to, buildEmail(cfg.From, to, subjectTest, body, start))
	result := EmailTestResult{DurationMS: time.Since(start).Milliseconds()}
	if err != nil {
		msg := smtpErrorMessage(err)
		slog.Info("email test", "recipients_count", len(to), "error", msg)
		result.Error = &msg
		return result
	}
	slog.Info("email test", "recipients_count", len(to), "ok", true)
	result.OK = true
	return result
}
