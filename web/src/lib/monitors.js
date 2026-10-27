// Regras de apresentação dos monitores no painel.
import { formatDuration } from "./format";

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
  if (m.type === "heartbeat") return `Espera um ping a cada ${formatDuration(m.interval_seconds)}`;
  return m.url;
}

export function pingUrl(token) {
  return `${window.location.origin}/ping/${token}`;
}
