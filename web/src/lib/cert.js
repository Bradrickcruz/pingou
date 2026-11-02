// Regras do certificado SSL para o formulário e a página de detalhe.
import { zonedParts } from "./timezone.js";

const DAY = 86_400_000;
// mesmos limites do backend (30 e 7 dias, medidos em horas)
export const CERT_WARN_DAYS = 30;
export const CERT_ALERT_DAYS = 7;

// isHttps("HTTPS://x") → true. Só https tem certificado para avisar.
export function isHttps(url) {
  return /^https:\/\//i.test(String(url ?? "").trim());
}

// O bloco de certificado só existe para monitor HTTP com URL https.
export function hasCert(m) {
  return m?.type === "http" && isHttps(m.url);
}

const pad = (n) => String(n).padStart(2, "0");

// certDate(iso) → "23/12/2026", no fuso das Configurações
export function certDate(iso) {
  const p = zonedParts(new Date(iso));
  return `${pad(p.day)}/${pad(p.month + 1)}/${p.year}`;
}

function relative(days) {
  if (days === 0) return "hoje";
  return `daqui a ${days} ${days === 1 ? "dia" : "dias"}`;
}

// certStatus → { level, text } ou null (sem bloco de certificado).
// level: "pending" (ainda não visto), "ok" (> 30 dias), "warn" (≤ 30),
// "alert" (≤ 7) ou "expired".
export function certStatus(m, now = Date.now()) {
  if (!hasCert(m)) return null;
  if (!m.cert_expires_at)
    return { level: "pending", text: "Aguardando o primeiro check" };

  const expires = new Date(m.cert_expires_at).getTime();
  const remaining = expires - now;
  const date = certDate(m.cert_expires_at);
  if (remaining <= 0) return { level: "expired", text: `Expirou em ${date}` };

  const days = Math.floor(remaining / DAY);
  const text = `Expira em ${date} (${relative(days)})`;
  if (remaining <= CERT_ALERT_DAYS * DAY) return { level: "alert", text };
  if (remaining <= CERT_WARN_DAYS * DAY) return { level: "warn", text };
  return { level: "ok", text };
}
