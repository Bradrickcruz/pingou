// node --test src/lib
import { test, afterEach } from "node:test";
import assert from "node:assert/strict";
import {
  setTimeZone,
  getTimeZone,
  timeZoneOptions,
  wallTicks,
  zoneOffsetLabel,
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
