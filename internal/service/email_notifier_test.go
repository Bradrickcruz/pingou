package service

import (
	"bufio"
	"io"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Bradrickcruz/pingou/internal/config"
	"github.com/Bradrickcruz/pingou/internal/domain"
)

func TestBuildEmail(t *testing.T) {
	now := time.Date(2026, 11, 1, 14, 3, 12, 0, time.UTC)
	raw := buildEmail("Pingou <alerta@exemplo.com>", []string{"a@x.com", "b@x.com"},
		"[Pingou] 🔴 Ação \r\nBcc: x@y", "Linha um\nLinha ção\n", now)

	if strings.Contains(string(raw), "\nBcc:") {
		t.Fatalf("injeção de cabeçalho:\n%s", raw)
	}
	head, _, ok := strings.Cut(string(raw), "\r\n\r\n")
	if !ok {
		t.Fatal("sem separador CRLF entre cabeçalho e corpo")
	}
	for _, line := range strings.Split(head, "\r\n") {
		if strings.ContainsAny(line, "\r\n") {
			t.Fatalf("linha com CR/LF solto: %q", line)
		}
	}

	msg, err := mail.ReadMessage(strings.NewReader(string(raw)))
	if err != nil {
		t.Fatal(err)
	}
	h := msg.Header
	for _, k := range []string{"From", "To", "Subject", "Date", "Message-Id", "Mime-Version", "Content-Type", "Content-Transfer-Encoding"} {
		if h.Get(k) == "" {
			t.Fatalf("cabeçalho %s ausente:\n%s", k, head)
		}
	}
	if !strings.HasPrefix(h.Get("Subject"), "=?utf-8?q?") {
		t.Fatalf("assunto não codificado: %q", h.Get("Subject"))
	}
	subject, err := new(mime.WordDecoder).DecodeHeader(h.Get("Subject"))
	if err != nil || subject != "[Pingou] 🔴 Ação Bcc: x@y" {
		t.Fatalf("assunto decodificado = %q (%v)", subject, err)
	}
	if h.Get("To") != "a@x.com, b@x.com" || h.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("cabeçalhos: %v", h)
	}
	if d, err := h.Date(); err != nil || !d.Equal(now) {
		t.Fatalf("Date = %v (%v)", d, err)
	}
	if !strings.HasSuffix(h.Get("Message-Id"), "@exemplo.com>") {
		t.Fatalf("Message-ID = %q", h.Get("Message-Id"))
	}
	body, _ := io.ReadAll(quotedprintable.NewReader(msg.Body))
	if string(body) != "Linha um\r\nLinha ção\r\n" {
		t.Fatalf("corpo = %q", body)
	}
}

func TestBodyFor(t *testing.T) {
	token := "segredo-do-token"
	errMsg := "connection refused"
	started := time.Date(2026, 11, 1, 14, 3, 12, 0, time.UTC)
	ended := started.Add(5*time.Minute + 12*time.Second)
	httpMon := &domain.Monitor{ID: "abc", Name: "API de pagamentos", Type: domain.TypeHTTP, URL: "https://api.exemplo.com/health"}
	hb := &domain.Monitor{ID: "hb", Name: "Backup", Type: domain.TypeHeartbeat, PingToken: &token}

	down := bodyFor(httpMon, &domain.Incident{StartedAt: started, LastError: &errMsg}, false, "", time.UTC)
	for _, want := range []string{"Monitor: API de pagamentos", "Tipo: http", "Alvo: https://api.exemplo.com/health", "Caiu em: 2026-11-01 14:03:12 UTC", "Erro: connection refused"} {
		if !strings.Contains(down, want) {
			t.Fatalf("down sem %q:\n%s", want, down)
		}
	}
	if strings.Contains(down, "Detalhes") || strings.Contains(down, "Fora do ar") {
		t.Fatalf("down com campo a mais:\n%s", down)
	}

	up := bodyFor(httpMon, &domain.Incident{StartedAt: started, EndedAt: &ended}, true, "https://pingou.exemplo.com", time.UTC)
	for _, want := range []string{"Fora do ar por: 5 min 12 s", "Detalhes: https://pingou.exemplo.com/monitors/abc"} {
		if !strings.Contains(up, want) {
			t.Fatalf("up sem %q:\n%s", want, up)
		}
	}
	if strings.Contains(up, "Erro:") {
		t.Fatalf("up com erro:\n%s", up)
	}

	beat := bodyFor(hb, &domain.Incident{StartedAt: started}, false, "", time.UTC)
	if strings.Contains(beat, "Alvo") || strings.Contains(beat, token) {
		t.Fatalf("heartbeat:\n%s", beat)
	}

	if s := subjectUp(httpMon, &domain.Incident{StartedAt: started, EndedAt: &ended}); s != "[Pingou] 🟢 API de pagamentos voltou (fora por 5 min)" {
		t.Fatalf("subjectUp = %q", s)
	}
	if s := subjectDown(&domain.Monitor{Name: "API\r\nBcc: x"}); s != "[Pingou] 🔴 API Bcc: x está fora do ar" {
		t.Fatalf("subjectDown = %q", s)
	}
}

func TestFormatDuration(t *testing.T) {
	cases := map[int64]string{0: "0 s", 45: "45 s", 60: "1 min", 312: "5 min 12 s", 7380: "2 h 3 min", 3600: "1 h", 100800: "1 d 4 h"}
	for sec, want := range cases {
		if got := formatDuration(sec); got != want {
			t.Fatalf("formatDuration(%d) = %q, want %q", sec, got, want)
		}
	}
}

// fakeSMTP é um servidor SMTP mínimo em memória (modo none, sem TLS).
type fakeSMTP struct {
	ln     net.Listener
	silent bool // aceita a conexão e não responde

	mu    sync.Mutex
	from  string
	rcpts []string
	data  string
	done  chan struct{}
}

// newFakeSMTP sobe o servidor; silent = aceita a conexão e não responde.
func newFakeSMTP(t *testing.T, silent ...bool) *fakeSMTP {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	f := &fakeSMTP{ln: ln, done: make(chan struct{}, 10), silent: len(silent) > 0 && silent[0]}
	t.Cleanup(func() { ln.Close() })
	go f.serve()
	return f
}

func (f *fakeSMTP) cfg(tlsMode string) config.SMTPConfig {
	_, port, _ := net.SplitHostPort(f.ln.Addr().String())
	p, _ := strconv.Atoi(port)
	return config.SMTPConfig{Host: "127.0.0.1", Port: p, From: "Pingou <pingou@localhost>", TLS: tlsMode}
}

func (f *fakeSMTP) serve() {
	for {
		conn, err := f.ln.Accept()
		if err != nil {
			return
		}
		go f.handle(conn)
	}
}

func (f *fakeSMTP) handle(conn net.Conn) {
	defer conn.Close()
	if f.silent {
		_, _ = io.Copy(io.Discard, conn)
		return
	}
	r := bufio.NewReader(conn)
	reply := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
	reply("220 fake ESMTP")
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		cmd := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(cmd, "EHLO"):
			reply("250 fake") // nunca anuncia STARTTLS
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			f.mu.Lock()
			f.from = line[len("MAIL FROM:"):]
			f.mu.Unlock()
			reply("250 ok")
		case strings.HasPrefix(cmd, "RCPT TO:"):
			f.mu.Lock()
			f.rcpts = append(f.rcpts, line[len("RCPT TO:"):])
			f.mu.Unlock()
			reply("250 ok")
		case cmd == "DATA":
			reply("354 go ahead")
			var b strings.Builder
			for {
				l, err := r.ReadString('\n')
				if err != nil {
					return
				}
				if l == ".\r\n" {
					break
				}
				b.WriteString(l)
			}
			f.mu.Lock()
			f.data = b.String()
			f.mu.Unlock()
			reply("250 queued")
		case cmd == "QUIT":
			reply("221 bye")
			f.done <- struct{}{}
			return
		default:
			reply("250 ok")
		}
	}
}

func TestSendMailNone(t *testing.T) {
	f := newFakeSMTP(t)
	to := []string{"ops@exemplo.com", "voce@exemplo.com"}
	msg := buildEmail("Pingou <pingou@localhost>", to, "[Pingou] Teste", "oi\n", time.Now())
	if err := sendMail(t.Context(), f.cfg(config.SMTPTLSNone), to, msg); err != nil {
		t.Fatal(err)
	}
	<-f.done
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.from != "<pingou@localhost>" {
		t.Fatalf("MAIL FROM = %q", f.from)
	}
	if strings.Join(f.rcpts, ",") != "<ops@exemplo.com>,<voce@exemplo.com>" {
		t.Fatalf("RCPT TO = %v", f.rcpts)
	}
	if !strings.Contains(f.data, "Subject: [Pingou] Teste\r\n") || !strings.Contains(f.data, "\r\n\r\noi\r\n") {
		t.Fatalf("DATA:\n%s", f.data)
	}
}

func TestSendMailTimeout(t *testing.T) {
	prev := smtpTimeout
	smtpTimeout = 200 * time.Millisecond
	defer func() { smtpTimeout = prev }()

	f := newFakeSMTP(t, true)
	start := time.Now()
	err := sendMail(t.Context(), f.cfg(config.SMTPTLSNone), []string{"a@x.com"}, []byte("x"))
	if err == nil {
		t.Fatal("sem erro com servidor mudo")
	}
	if d := time.Since(start); d > 2*time.Second {
		t.Fatalf("demorou %s", d)
	}
	if msg := smtpErrorMessage(err); !strings.HasPrefix(msg, "timeout após") {
		t.Fatalf("mensagem = %q", msg)
	}
}

func TestStartTLSRequired(t *testing.T) {
	f := newFakeSMTP(t) // não anuncia STARTTLS
	err := sendMail(t.Context(), f.cfg(config.SMTPTLSStartTLS), []string{"a@x.com"}, []byte("x"))
	if err != errStartTLSUnsupported {
		t.Fatalf("err = %v, want %v", err, errStartTLSUnsupported)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.from != "" {
		t.Fatal("enviou MAIL FROM sem TLS")
	}
}

func TestEmailNotifierSkipsWithoutRecipients(t *testing.T) {
	f := newFakeSMTP(t)
	n := NewEmailNotifier(f.cfg(config.SMTPTLSNone), func() string { return "" }, "")
	n.NotifyDown(t.Context(), &domain.Monitor{Name: "x"}, &domain.Incident{StartedAt: time.Now()})
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.from != "" {
		t.Fatal("enviou sem destinatários")
	}
}

func TestEmailNotifierSendsDown(t *testing.T) {
	f := newFakeSMTP(t)
	n := NewEmailNotifier(f.cfg(config.SMTPTLSNone), func() string { return "ops@exemplo.com" }, "")
	n.NotifyDown(t.Context(), &domain.Monitor{Name: "API", Type: domain.TypeHTTP, URL: "https://x.test"}, &domain.Incident{StartedAt: time.Now()})
	<-f.done
	f.mu.Lock()
	defer f.mu.Unlock()
	if !strings.Contains(f.data, "Subject: =?utf-8?q?") || len(f.rcpts) != 1 {
		t.Fatalf("rcpts=%v DATA:\n%s", f.rcpts, f.data)
	}
}

func TestSendTestEmail(t *testing.T) {
	f := newFakeSMTP(t)
	res := SendTestEmail(t.Context(), f.cfg(config.SMTPTLSNone), []string{"ops@exemplo.com"})
	if !res.OK || res.Error != nil {
		t.Fatalf("resultado: %+v", res)
	}
	<-f.done

	res = SendTestEmail(t.Context(), config.SMTPConfig{Host: "127.0.0.1", Port: 1, From: "a@b.c", TLS: "none"}, []string{"ops@exemplo.com"})
	if res.OK || res.Error == nil {
		t.Fatalf("porta fechada: %+v", res)
	}
}
