// node --test src/lib
import { test, afterEach } from "node:test";
import assert from "node:assert/strict";
import {
  setTimeZone,
  getTimeZone,
  timeZoneOptions,
  filterTimeZones,
  wallTicks,
  zoneOffsetLabel,
  zonedToUtc,
  utcToZonedInput,
} from "./timezone.js";
import {
  clock,
  checkTime,
  dayKey,
  dayLabel,
  fullDate,
  startedLabel,
} from "./format.js";

afterEach(() => setTimeZone("UTC"));

const ISO = "2026-11-01T14:03:12Z";

test("padrão UTC e nome inválido cai em UTC", () => {
  assert.equal(getTimeZone(), "UTC");
  assert.equal(clock(ISO), "14:03");
  assert.equal(setTimeZone("Marte/Olimpo"), "UTC");
  assert.equal(setTimeZone(""), "UTC");
});

test("UTC → America/Sao_Paulo e Asia/Tokyo", () => {
  setTimeZone("America/Sao_Paulo");
  assert.equal(clock(ISO), "11:03");
  assert.equal(fullDate(ISO), "1 nov 2026, 11:03");
  setTimeZone("Asia/Tokyo");
  assert.equal(clock(ISO), "23:03");
});

test("virada de dia depende do fuso", () => {
  const late = "2026-11-01T23:30:00Z";
  const now = Date.parse("2026-11-02T01:00:00Z");
  setTimeZone("UTC");
  assert.equal(dayKey(late), "2026-11-1");
  assert.equal(startedLabel(late, now), "ontem, 23:30");
  setTimeZone("Asia/Tokyo"); // 08:30 de 2 nov; agora 10:00 de 2 nov
  assert.equal(dayKey(late), "2026-11-2");
  assert.equal(startedLabel(late, now), "hoje, 08:30");
  assert.equal(dayLabel(late, now), "Hoje · 2 nov");
  assert.equal(checkTime(late, now), "08:30:00");
  setTimeZone("America/Sao_Paulo"); // 20:30 de 1 nov; agora 22:00 de 1 nov
  assert.equal(startedLabel(late, now), "hoje, 20:30");
});

test("marcas do gráfico caem em horas redondas do fuso", () => {
  const since = "2026-11-01T13:17:00Z";
  const until = "2026-11-02T13:17:00Z";
  setTimeZone("America/Sao_Paulo");
  const labels = wallTicks(since, until, 4 * 3_600_000).map(clock);
  assert.deepEqual(labels, [
    "12:00",
    "16:00",
    "20:00",
    "00:00",
    "04:00",
    "08:00",
  ]);

  setTimeZone("Asia/Kolkata"); // +05:30
  const k = wallTicks(since, until, 4 * 3_600_000).map(clock);
  assert.ok(
    k.every((l) => l.endsWith(":00")),
    k.join(","),
  );

  setTimeZone("Asia/Tokyo");
  const days = wallTicks(
    "2026-10-26T00:00:00Z",
    "2026-11-02T00:00:00Z",
    86_400_000,
    86_400_000,
  );
  assert.equal(days.length, 7);
  assert.ok(days.every((d) => clock(d) === "00:00"));
});

test("rótulo do deslocamento", () => {
  const d = new Date(ISO);
  assert.equal(zoneOffsetLabel(d, "UTC"), "UTC");
  assert.equal(zoneOffsetLabel(d, "America/Sao_Paulo"), "-03");
  assert.equal(zoneOffsetLabel(d, "Asia/Kolkata"), "+05:30");
});

test("lista de fusos com UTC no topo", () => {
  const list = timeZoneOptions();
  assert.equal(list[0], "UTC");
  assert.ok(list.includes("America/Sao_Paulo"));
  assert.equal(list.filter((z) => z === "UTC").length, 1);
});

test("filtro: palavras em qualquer ordem, sem acento nem _ ou /", () => {
  const zones = ["UTC", "America/Sao_Paulo", "America/Santiago", "Asia/Tokyo"];
  assert.deepEqual(filterTimeZones(zones, "sao paulo"), ["America/Sao_Paulo"]);
  assert.deepEqual(filterTimeZones(zones, "SÃO"), ["America/Sao_Paulo"]);
  assert.deepEqual(filterTimeZones(zones, "paulo america"), [
    "America/Sao_Paulo",
  ]);
  assert.deepEqual(filterTimeZones(zones, "san"), ["America/Santiago"]);
  assert.deepEqual(filterTimeZones(zones, "  "), zones);
  assert.deepEqual(filterTimeZones(zones, "marte"), []);
});

test("zonedToUtc lê a hora no fuso das Configurações, não no do navegador", () => {
  const iso = (v, tz) => zonedToUtc(v, tz).toISOString();
  assert.equal(
    iso("2026-11-10T02:00", "America/Sao_Paulo"),
    "2026-11-10T05:00:00.000Z",
  );
  assert.equal(
    iso("2026-11-10T02:00", "Asia/Kolkata"),
    "2026-11-09T20:30:00.000Z",
  );
  assert.equal(iso("2026-11-10T02:00", "UTC"), "2026-11-10T02:00:00.000Z");
  // usa o fuso atual quando não recebe um
  setTimeZone("Asia/Kolkata");
  assert.equal(
    zonedToUtc("2026-11-10T02:00").toISOString(),
    "2026-11-09T20:30:00.000Z",
  );
  assert.equal(zonedToUtc(""), null);
  assert.equal(zonedToUtc("10/11/2026 02:00"), null);
});

test("zonedToUtc atravessa o horário de verão", () => {
  // Nova York: 8 mar 2026, 02:00–03:00 não existe; 1 nov, 01:00–02:00 repete
  const iso = (v) => zonedToUtc(v, "America/New_York").toISOString();
  assert.equal(iso("2026-03-07T12:00"), "2026-03-07T17:00:00.000Z"); // EST
  assert.equal(iso("2026-03-09T12:00"), "2026-03-09T16:00:00.000Z"); // EDT
  assert.equal(iso("2026-11-02T12:00"), "2026-11-02T17:00:00.000Z"); // EST de novo
});

test("utcToZonedInput é o inverso de zonedToUtc", () => {
  for (const tz of ["America/Sao_Paulo", "Asia/Kolkata", "UTC"]) {
    const v = "2026-11-10T02:00";
    assert.equal(utcToZonedInput(zonedToUtc(v, tz), tz), v);
  }
});
