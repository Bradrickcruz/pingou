package checker

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

func newTestRegistry() *Registry {
	return NewRegistry(map[domain.MonitorType]domain.Checker{
		domain.TypeHTTP: NewHTTPChecker(5, 10),
		domain.TypeTCP:  NewTCPChecker(),
	})
}

// Regressão da Fase 1: o HTTP passando pelo Registry se comporta como antes.
func TestRegistryHTTP(t *testing.T) {
	status := http.StatusOK
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	}))
	defer srv.Close()

	r := newTestRegistry()
	m := &domain.Monitor{Type: domain.TypeHTTP, URL: srv.URL, TimeoutSeconds: 2}

	res := r.Check(context.Background(), m)
	if !res.Success || res.StatusCode == nil || *res.StatusCode != 200 {
		t.Fatalf("200: %+v, want sucesso com status 200", res)
	}

	status = http.StatusInternalServerError
	res = r.Check(context.Background(), m)
	if res.Success || res.ErrorMessage == nil || *res.ErrorMessage != "unexpected status code: 500" {
		t.Fatalf("500: %+v, want falha com mensagem de status", res)
	}
}

func TestRegistryUnknownType(t *testing.T) {
	res := newTestRegistry().Check(context.Background(), &domain.Monitor{Type: "dns"})
	if res.Success || res.ErrorMessage == nil {
		t.Fatalf("tipo desconhecido: %+v, want falha com erro", res)
	}
}

func tcpMonitor(host string, port int) *domain.Monitor {
	return &domain.Monitor{
		Type:           domain.TypeTCP,
		Config:         domain.MonitorConfig{Host: host, Port: port},
		TimeoutSeconds: 1,
	}
}

func TestTCPCheckerUp(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()

	port := ln.Addr().(*net.TCPAddr).Port
	res := newTestRegistry().Check(context.Background(), tcpMonitor("127.0.0.1", port))
	if !res.Success || res.ErrorMessage != nil || res.StatusCode != nil {
		t.Fatalf("porta aberta: %+v, want sucesso sem status code", res)
	}
}

func TestTCPCheckerRefused(t *testing.T) {
	// abre e fecha para obter uma porta livre (fechada)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	res := NewTCPChecker().Check(context.Background(), tcpMonitor("127.0.0.1", port))
	if res.Success || res.ErrorMessage == nil || *res.ErrorMessage != "connection refused" {
		t.Fatalf("porta fechada: %+v, want 'connection refused'", res)
	}
}

func TestTCPCheckerTimeout(t *testing.T) {
	// 10.255.255.1 não é roteável: o SYN fica sem resposta até o timeout
	res := NewTCPChecker().Check(context.Background(), tcpMonitor("10.255.255.1", 81))
	if res.Success {
		t.Fatal("IP não roteável deu sucesso")
	}
	if *res.ErrorMessage != "timeout" {
		t.Skipf("rede do ambiente respondeu %q em vez de esperar o timeout", *res.ErrorMessage)
	}
}

func TestTCPCheckerNoSuchHost(t *testing.T) {
	res := NewTCPChecker().Check(context.Background(), tcpMonitor("pingou-nao-existe.invalid", 80))
	if res.Success {
		t.Fatal("host inexistente deu sucesso")
	}
	if *res.ErrorMessage != "no such host" {
		t.Skipf("DNS do ambiente respondeu %q", *res.ErrorMessage)
	}
}

// A mensagem não repete a URL (`Get "<url>": …`): ela já aparece ao lado do erro.
func TestHTTPCheckerErrorMessage(t *testing.T) {
	refused := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	refused.Close()

	res := NewHTTPChecker(5, 10).Check(context.Background(), &domain.Monitor{Type: domain.TypeHTTP, URL: refused.URL, TimeoutSeconds: 2})
	if res.Success || res.ErrorMessage == nil || strings.HasPrefix(*res.ErrorMessage, "Get ") ||
		!strings.Contains(*res.ErrorMessage, "connection refused") {
		t.Fatalf("connection refused: %+v", res.ErrorMessage)
	}

	release := make(chan struct{})
	slow := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	defer slow.Close()
	defer close(release)

	res = NewHTTPChecker(5, 10).Check(context.Background(), &domain.Monitor{Type: domain.TypeHTTP, URL: slow.URL, TimeoutSeconds: 1})
	if res.ErrorMessage == nil || *res.ErrorMessage != "context deadline exceeded (timeout 1s)" {
		t.Fatalf("timeout: %v, want 'context deadline exceeded (timeout 1s)'", res.ErrorMessage)
	}
}
