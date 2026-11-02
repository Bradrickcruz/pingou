// Formatação de tempo e números para o painel (pt-BR).
import { zonedParts } from "./timezone.js";

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

// formatBytes(13002342) → "12,4 MB"
const BYTE_UNITS = ["B", "KB", "MB", "GB", "TB"];
const bytesFormat = new Intl.NumberFormat("pt-BR", {
  maximumFractionDigits: 1,
});
export function formatBytes(bytes) {
  if (bytes == null) return null;
  let value = bytes;
  let unit = 0;
  while (value >= 1024 && unit < BYTE_UNITS.length - 1) {
    value /= 1024;
    unit++;
  }
  return `${bytesFormat.format(value)} ${BYTE_UNITS[unit]}`;
}

// backupFilename() → "pingou-2026-10-29.db" (data UTC, igual ao backend)
export function backupFilename(now = new Date()) {
  return `pingou-${now.toISOString().slice(0, 10)}.db`;
}

export function plural(n, one, many) {
  return `${n} ${n === 1 ? one : many}`;
}

// ── datas dos incidentes ───────────────────────────────────────────────────
// Meses e dias à mão: o Intl em pt-BR escreve "28 de out.", e a tela usa "28 out".
const MONTHS = [
  "jan",
  "fev",
  "mar",
  "abr",
  "mai",
  "jun",
  "jul",
  "ago",
  "set",
  "out",
  "nov",
  "dez",
];
const WEEKDAYS = [
  "Domingo",
  "Segunda",
  "Terça",
  "Quarta",
  "Quinta",
  "Sexta",
  "Sábado",
];

// Todas as datas abaixo saem no fuso das Configurações (lib/timezone.js).
const pad = (n) => String(n).padStart(2, "0");

// dias de calendário (no fuso escolhido) entre a data e hoje: 0 = hoje, 1 = ontem
function daysAgo(date, now) {
  const day = (d) => {
    const p = zonedParts(d);
    return Date.UTC(p.year, p.month, p.day);
  };
  return Math.round((day(new Date(now)) - day(date)) / 86_400_000);
}

function shortDate(date, now) {
  const p = zonedParts(date);
  const year = p.year !== zonedParts(new Date(now)).year ? ` ${p.year}` : "";
  return `${p.day} ${MONTHS[p.month]}${year}`;
}

// clock(iso) → "08:03"
export function clock(iso) {
  const p = zonedParts(new Date(iso));
  return `${pad(p.hour)}:${pad(p.minute)}`;
}

// clockSeconds(iso) → "08:03:41"
function clockSeconds(iso) {
  const p = zonedParts(new Date(iso));
  return `${pad(p.hour)}:${pad(p.minute)}:${pad(p.second)}`;
}

// dayKey(iso) → "2026-10-28" no fuso escolhido, para agrupar por dia
export function dayKey(iso) {
  const p = zonedParts(new Date(iso));
  return `${p.year}-${p.month + 1}-${p.day}`;
}

// dayLabel → "Hoje · 28 out", "Ontem · 27 out", "Domingo · 26 out" (até 6 dias), "20 out"
export function dayLabel(iso, now = Date.now()) {
  const date = new Date(iso);
  const diff = daysAgo(date, now);
  const short = shortDate(date, now);
  if (diff === 0) return `Hoje · ${short}`;
  if (diff === 1) return `Ontem · ${short}`;
  if (diff > 1 && diff < 7)
    return `${WEEKDAYS[zonedParts(date).weekday]} · ${short}`;
  return short;
}

// startedLabel → "hoje, 08:03", "ontem, 11:43", "26 out, 19:05"
export function startedLabel(iso, now = Date.now()) {
  const date = new Date(iso);
  const diff = daysAgo(date, now);
  const day = diff === 0 ? "hoje" : diff === 1 ? "ontem" : shortDate(date, now);
  return `${day}, ${clock(iso)}`;
}

// timeRange → "23:58 → 00:04". Quedas de 1 dia ou mais levam a data do início,
// senão "08:00 → 09:00" esconderia o dia que passou.
export function timeRange(startIso, endIso, now = Date.now()) {
  const long = new Date(endIso) - new Date(startIso) >= 86_400_000;
  const start = long
    ? `${shortDate(new Date(startIso), now)}, ${clock(startIso)}`
    : clock(startIso);
  return `${start} → ${clock(endIso)}`;
}

// ── datas da página de detalhe ─────────────────────────────────────────────

// fullDate → "20 out 2026, 14:32" (sempre com o ano e a hora)
export function fullDate(iso) {
  const p = zonedParts(new Date(iso));
  return `${p.day} ${MONTHS[p.month]} ${p.year}, ${clock(iso)}`;
}

// dateTimeSeconds → "20 out 2026, 14:32:05" (aviso de conexão perdida)
export function dateTimeSeconds(iso) {
  const p = zonedParts(new Date(iso));
  return `${p.day} ${MONTHS[p.month]} ${p.year}, ${clockSeconds(iso)}`;
}

// dayMonth → "27 out" (rótulo do eixo X do gráfico)
export function dayMonth(date) {
  const p = zonedParts(date);
  return `${p.day} ${MONTHS[p.month]}`;
}

function dayWord(date, now) {
  const diff = daysAgo(date, now);
  if (diff === 0) return "Hoje";
  if (diff === 1) return "Ontem";
  if (diff > 1 && diff < 7) return WEEKDAYS[zonedParts(date).weekday];
  return shortDate(date, now);
}

// whenLabel → "Hoje, 10:12", "Ontem, 08:00", "Domingo, 19:05" (até 6 dias), "24 out, 03:20"
export function whenLabel(iso, now = Date.now()) {
  return `${dayWord(new Date(iso), now)}, ${clock(iso)}`;
}

// bucketLabel → "Hoje, 14:10 – 14:20", "27 out, 06:00 – 12:00". Balde que vira
// a meia-noite leva só a data do início.
export function bucketLabel(startIso, bucketSeconds, now = Date.now()) {
  const end = new Date(
    new Date(startIso).getTime() + bucketSeconds * 1000,
  ).toISOString();
  return `${whenLabel(startIso, now)} – ${clock(end)}`;
}

// checkTime → "14:02:41" hoje; "Ontem, 14:02", "24 out, 03:20" em outros dias
export function checkTime(iso, now = Date.now()) {
  if (daysAgo(new Date(iso), now) === 0) return clockSeconds(iso);
  return whenLabel(iso, now);
}
