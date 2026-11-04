// Regras do formulário de manutenção (espelham o servidor).
import { utcToZonedInput, zonedToUtc } from "./timezone.js";

export const MAX_MAINTENANCE_MS = 7 * 24 * 3600 * 1000; // teto de uma janela
const STEP_MS = 5 * 60 * 1000;

// Início padrão: próximo múltiplo de 5 min com pelo menos 5 min de folga; 1 h de duração.
export function defaultRange(now = Date.now()) {
  const start = Math.ceil((now + STEP_MS) / STEP_MS) * STEP_MS;
  return {
    starts: utcToZonedInput(new Date(start)),
    ends: utcToZonedInput(new Date(start + 3600 * 1000)),
  };
}

// validateWindow → { campo: mensagem }; vazio = válido. starts/ends são valores
// de <input type="datetime-local"> no fuso das Configurações.
export function validateWindow({ ids, starts, ends }, now = Date.now()) {
  const errors = {};
  if (ids.length === 0) errors.monitors = "Escolha pelo menos um monitor.";
  const s = zonedToUtc(starts);
  const e = zonedToUtc(ends);
  if (!s) errors.starts = "Informe a data e a hora do início.";
  else if (s.getTime() <= now)
    errors.starts = "O início precisa ser no futuro.";
  if (!e) errors.ends = "Informe a data e a hora do fim.";
  else if (s && e <= s) errors.ends = "O fim precisa ser depois do início.";
  else if (s && e - s > MAX_MAINTENANCE_MS)
    errors.ends = "Uma manutenção dura no máximo 7 dias.";
  return errors;
}

// toApiDate("2026-11-10T02:00") → "2026-11-10T05:00:00Z" (RFC 3339, UTC, sem ms)
export function toApiDate(value) {
  return zonedToUtc(value)
    .toISOString()
    .replace(/\.\d{3}Z$/, "Z");
}
