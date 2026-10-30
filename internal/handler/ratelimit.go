package handler

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// pingRateLimit: pings com token inválido por IP por minuto em /ping/{token}.
// Ping válido nunca conta nem é bloqueado: atrás de um proxy todos dividem o
// mesmo IP, e um atacante não pode derrubar os heartbeats de verdade.
const pingRateLimit = 60

// webhookTestRateLimit: testes de webhook por IP por minuto. Cada teste é uma
// requisição de saída para uma URL qualquer, então o limite é baixo.
const webhookTestRateLimit = 5

// rateLimiter é uma janela fixa por IP, em memória (um binário só, sem Redis).
// O map é zerado a cada janela, então a memória fica limitada aos IPs de um minuto.
// ponytail: usa RemoteAddr; atrás de proxy reverso todos os clientes contam
// como um IP só. Suporte a X-Forwarded-For confiável fica para depois.
type rateLimiter struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	windowStart time.Time
	counts      map[string]int
	now         func() time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, counts: map[string]int{}, now: time.Now}
}

// allow registra uma requisição do IP e diz se ela cabe no limite. Se não couber,
// devolve quanto falta para a janela reiniciar.
func (l *rateLimiter) allow(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if now.Sub(l.windowStart) >= l.window {
		l.windowStart = now
		clear(l.counts)
	}
	if l.counts[ip] >= l.limit {
		return false, l.window - now.Sub(l.windowStart)
	}
	l.counts[ip]++
	return true, 0
}

// reject responde 429 com Retry-After.
func rejectRateLimited(w http.ResponseWriter, retry time.Duration) {
	w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())+1))
	writeError(w, http.StatusTooManyRequests, "too many requests", "RATE_LIMITED")
}

// middleware conta toda requisição.
func (l *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok, retry := l.allow(clientIP(r)); !ok {
			rejectRateLimited(w, retry)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
