// Formatação de tempo e números para o painel (pt-BR).

const units = [
  { s: 604800, one: "semana", many: "semanas" },
  { s: 86400, one: "dia", many: "dias" },
  { s: 3600, one: "h", many: "h" },
  { s: 60, one: "min", many: "min" },
  { s: 1, one: "s", many: "s" },
];

// formatDuration(90000) → "1 dia e 1 h". Usa no máximo as duas maiores unidades.
export function formatDuration(seconds) {
  let rest = Math.max(0, Math.round(seconds));
  if (rest === 0) return "0 s";
  const parts = [];
  for (const u of units) {
    if (parts.length === 2) break;
    const n = Math.floor(rest / u.s);
    if (n > 0) {
      parts.push(`${n} ${n === 1 ? u.one : u.many}`);
      rest -= n * u.s;
    } else if (parts.length > 0) {
      break; // "1 dia e 30 s" confunde mais do que ajuda
    }
  }
  return parts.join(" e ");
}

// timeAgo("2026-…Z", Date.now()) → "há 4 s"
export function timeAgo(iso, now = Date.now()) {
  if (!iso) return null;
  const diff = Math.max(0, Math.floor((now - new Date(iso).getTime()) / 1000));
  if (diff < 60) return `há ${diff} s`;
  if (diff < 3600) return `há ${Math.floor(diff / 60)} min`;
  if (diff < 86400) return `há ${Math.floor(diff / 3600)} h`;
  const days = Math.floor(diff / 86400);
  return `há ${days} ${days === 1 ? "dia" : "dias"}`;
}

// formatUptime(76.666) → "76,7%"; 100 → "100%"
export function formatUptime(value) {
  if (value == null) return null;
  if (value >= 100) return "100%";
  return `${value.toLocaleString("pt-BR", { minimumFractionDigits: 1, maximumFractionDigits: 1 })}%`;
}

export function plural(n, one, many) {
  return `${n} ${n === 1 ? one : many}`;
}
