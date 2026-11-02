// Fuso das datas do painel. Vem das Configurações (settings.timezone); a API
// sempre devolve UTC e a conversão acontece só aqui, na hora de mostrar.

export const DEFAULT_TIME_ZONE = "UTC";

let current = DEFAULT_TIME_ZONE;
const formatters = new Map();

// Lista curta de reserva para navegadores sem Intl.supportedValuesOf.
const FALLBACK_ZONES = [
  "America/Sao_Paulo",
  "America/Manaus",
  "America/Fortaleza",
  "America/Recife",
  "America/Belem",
  "America/Cuiaba",
  "America/Rio_Branco",
  "America/Noronha",
  "America/Argentina/Buenos_Aires",
  "America/Santiago",
  "America/Bogota",
  "America/Lima",
  "America/Mexico_City",
  "America/New_York",
  "America/Chicago",
  "America/Denver",
  "America/Los_Angeles",
  "Europe/Lisbon",
  "Europe/London",
  "Europe/Madrid",
  "Europe/Paris",
  "Europe/Berlin",
  "Africa/Luanda",
  "Africa/Maputo",
  "Asia/Kolkata",
  "Asia/Shanghai",
  "Asia/Tokyo",
  "Australia/Sydney",
];

export function isValidTimeZone(tz) {
  if (!tz || typeof tz !== "string") return false;
  try {
    new Intl.DateTimeFormat("en-US", { timeZone: tz });
    return true;
  } catch {
    return false;
  }
}

// setTimeZone troca o fuso usado por todos os formatadores. Nome inválido
// (ex.: navegador antigo sem esse fuso) cai em UTC.
export function setTimeZone(tz) {
  current = isValidTimeZone(tz) ? tz : DEFAULT_TIME_ZONE;
  return current;
}

export function getTimeZone() {
  return current;
}

// Lista para o seletor: UTC no topo, depois os nomes IANA em ordem.
export function timeZoneOptions() {
  let zones;
  try {
    zones = Intl.supportedValuesOf?.("timeZone") ?? [];
  } catch {
    zones = [];
  }
  if (zones.length === 0) zones = FALLBACK_ZONES;
  const rest = zones.filter((z) => z !== "UTC" && z !== "Etc/UTC").sort();
  return ["UTC", ...rest];
}

// Normaliza para comparar: sem acento, minúsculas, "/" e "_" viram espaço.
const norm = (s) =>
  s
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[/_]+/g, " ");

// filterTimeZones: cada palavra digitada precisa aparecer no nome, em qualquer
// ordem. "sao paulo" acha America/Sao_Paulo; "tokyo" acha Asia/Tokyo.
export function filterTimeZones(zones, query) {
  const words = norm(query).split(/\s+/).filter(Boolean);
  if (words.length === 0) return zones;
  return zones.filter((z) => {
    const n = norm(z);
    return words.every((w) => n.includes(w));
  });
}

function formatter(tz) {
  let f = formatters.get(tz);
  if (!f) {
    f = new Intl.DateTimeFormat("en-US", {
      timeZone: tz,
      hourCycle: "h23",
      year: "numeric",
      month: "numeric",
      day: "numeric",
      hour: "numeric",
      minute: "numeric",
      second: "numeric",
      weekday: "short",
    });
    formatters.set(tz, f);
  }
  return f;
}

const WEEKDAY_INDEX = {
  Sun: 0,
  Mon: 1,
  Tue: 2,
  Wed: 3,
  Thu: 4,
  Fri: 5,
  Sat: 6,
};

// zonedParts(date) → { year, month (0–11), day, hour, minute, second, weekday (0 = domingo) }
// no fuso atual (ou em tz).
export function zonedParts(date, tz = current) {
  const d = date instanceof Date ? date : new Date(date);
  const p = {};
  for (const { type, value } of formatter(tz).formatToParts(d)) p[type] = value;
  return {
    year: Number(p.year),
    month: Number(p.month) - 1,
    day: Number(p.day),
    hour: Number(p.hour) % 24,
    minute: Number(p.minute),
    second: Number(p.second),
    weekday: WEEKDAY_INDEX[p.weekday],
  };
}

// "Hora de parede" em ms: a data e hora do fuso escritas como se fossem UTC.
// Diferença entre dois desses valores conta dias e horas de calendário.
export function wallMs(date, tz = current) {
  const p = zonedParts(date, tz);
  return Date.UTC(p.year, p.month, p.day, p.hour, p.minute, p.second);
}

// Instante real de uma hora de parede. Duas passadas acertam a troca de horário
// de verão na maioria dos casos.
function instantOf(wall, tz) {
  let t = wall;
  for (let i = 0; i < 2; i++) {
    const offset = wallMs(new Date(t), tz) - t;
    t = wall - offset;
  }
  return new Date(t);
}

// wallTicks: marcas em horários "redondos" do fuso entre since e until.
// stepMs múltiplo de 1 h (ex.: 4 h, 1 dia); as marcas caem em múltiplos de
// alignMs contados da meia-noite (padrão: o próprio passo, até 1 dia).
export function wallTicks(
  since,
  until,
  stepMs,
  alignMs = Math.min(stepMs, 86_400_000),
  tz = current,
) {
  const start = new Date(since).getTime();
  const end = new Date(until).getTime();
  // o wallMs ignora os ms: arredonda o início para cima no segundo
  let wall = Math.ceil(wallMs(new Date(start), tz) / alignMs) * alignMs;
  const ticks = [];
  for (let guard = 0; guard < 1000; guard++, wall += stepMs) {
    const t = instantOf(wall, tz);
    if (t.getTime() >= end) break;
    if (t.getTime() >= start) ticks.push(t);
  }
  return ticks;
}

// zoneOffsetLabel(agora, "America/Sao_Paulo") → "-03"; "UTC" → "UTC" (igual às mensagens)
export function zoneOffsetLabel(date = new Date(), tz = current) {
  const offsetMin = Math.round(
    (wallMs(date, tz) - Math.floor(date.getTime() / 1000) * 1000) / 60_000,
  );
  if (offsetMin === 0) return "UTC";
  const sign = offsetMin < 0 ? "-" : "+";
  const abs = Math.abs(offsetMin);
  const h = String(Math.floor(abs / 60)).padStart(2, "0");
  const m = abs % 60;
  return `${sign}${h}${m ? `:${String(m).padStart(2, "0")}` : ""}`;
}
