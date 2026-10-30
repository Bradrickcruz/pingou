// Regras de apresentação dos monitores no painel.
import { formatDuration, formatUptime, timeAgo } from "./format.js";

export const typeLabel = { http: "HTTP", tcp: "TCP", heartbeat: "Heartbeat" };

// Situação do monitor para o painel: pausado vence o estado do check.
export function statusOf(m) {
  return m.enabled ? m.current_state : "PAUSED";
}

// Ordem por gravidade: fora do ar, desconhecido, no ar, pausado.
const severity = { DOWN: 0, UNKNOWN: 1, UP: 2, PAUSED: 3 };

export function bySeverity(a, b) {
  return (
    severity[statusOf(a)] - severity[statusOf(b)] ||
    a.name.localeCompare(b.name, "pt-BR")
  );
}

export function countByStatus(monitors) {
  const counts = { UP: 0, DOWN: 0, UNKNOWN: 0, PAUSED: 0 };
  for (const m of monitors) counts[statusOf(m)]++;
  return counts;
}

// Texto do alvo: URL (http), host:porta (tcp) ou a expectativa do heartbeat.
export function targetOf(m) {
  if (m.type === "tcp") return `${m.config.host}:${m.config.port}`;
  if (m.type === "heartbeat")
    return `Espera um ping a cada ${formatDuration(m.interval_seconds)}`;
  return m.url;
}

export function pingUrl(token) {
  return `${window.location.origin}/ping/${token}`;
}

// heartbeatExpectation(86400) → "espera um ping por dia"; 7200 → "espera um ping a cada 2 h"
const PER = { 3600: "por hora", 86400: "por dia", 604800: "por semana" };
export function heartbeatExpectation(seconds) {
  return `espera um ping ${PER[seconds] ?? `a cada ${formatDuration(seconds)}`}`;
}

// cardFooter → linha de rodapé do card do painel (mobile), por estado:
// DOWN "Fora há 1 dia e 4 h · 90,0% em 24 h"; UP "99,2% em 24 h · 214 ms · há 41 s";
// UNKNOWN "Aguardando o primeiro check"; PAUSED "Pausado".
export function cardFooter(m, now = Date.now()) {
  const status = statusOf(m);
  const heartbeat = m.type === "heartbeat";
  if (status === "PAUSED") return "Pausado";
  if (status === "UNKNOWN")
    return heartbeat
      ? "Aguardando o primeiro ping"
      : "Aguardando o primeiro check";

  const uptime = formatUptime(m.uptime_24h);
  const parts = [];
  if (status === "DOWN") {
    parts.push(
      m.down_since
        ? `Fora há ${formatDuration((now - new Date(m.down_since)) / 1000)}`
        : "Fora do ar",
    );
    if (uptime) parts.push(`${uptime} em 24 h`);
    return parts.join(" · ");
  }
  if (uptime) parts.push(`${uptime} em 24 h`);
  if (!heartbeat && m.last_latency_ms != null)
    parts.push(`${m.last_latency_ms} ms`);
  const ago = timeAgo(
    heartbeat ? (m.last_ping_at ?? m.last_checked_at) : m.last_checked_at,
    now,
  );
  if (ago) parts.push(heartbeat ? `ping ${ago}` : ago);
  return parts.join(" · ");
}
