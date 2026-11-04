// node --test src/lib
import { test, afterEach } from "node:test";
import assert from "node:assert/strict";
import {
  defaultRange,
  toApiDate,
  validateWindow,
  windowState,
} from "./maintenance.js";
import { setTimeZone } from "./timezone.js";

afterEach(() => setTimeZone("UTC"));

const now = Date.parse("2026-11-10T12:00:00Z");

test("validateWindow: válido e cada erro", () => {
  setTimeZone("America/Sao_Paulo"); // 12:00 UTC = 09:00 em São Paulo
  const ok = {
    ids: ["a"],
    starts: "2026-11-10T10:00",
    ends: "2026-11-10T11:00",
  };
  assert.deepEqual(validateWindow(ok, now), {});
  assert.ok(validateWindow({ ...ok, ids: [] }, now).monitors);
  // 08:30 em São Paulo já passou (seria futuro se lido em UTC)
  assert.ok(validateWindow({ ...ok, starts: "2026-11-10T08:30" }, now).starts);
  assert.ok(validateWindow({ ...ok, ends: "2026-11-10T10:00" }, now).ends);
  assert.ok(validateWindow({ ...ok, ends: "2026-11-17T10:01" }, now).ends);
  assert.deepEqual(
    validateWindow({ ...ok, ends: "2026-11-17T10:00" }, now),
    {},
  );
  assert.ok(validateWindow({ ...ok, starts: "" }, now).starts);
});

test("toApiDate converte do fuso das Configurações para UTC", () => {
  setTimeZone("Asia/Kolkata");
  assert.equal(toApiDate("2026-11-10T02:00"), "2026-11-09T20:30:00Z");
});

test("defaultRange: próximo múltiplo de 5 min com folga, 1 h de duração", () => {
  assert.deepEqual(defaultRange(Date.parse("2026-11-10T12:03:10Z")), {
    starts: "2026-11-10T12:10",
    ends: "2026-11-10T13:10",
  });
});

test("validateWindow em edição", () => {
  // 12:00 UTC; a janela começou a ser editada com início 11:00 (já passou)
  const base = {
    edit: true,
    starts: "2026-11-10T11:00",
    initialStarts: "2026-11-10T11:00",
    ends: "2026-11-10T13:00",
  };
  // sem lista de monitores e sem mexer no início: não acusa "início no passado"
  assert.deepEqual(validateWindow(base, now), {});
  // mudou o início para o passado: acusa
  assert.ok(
    validateWindow({ ...base, starts: "2026-11-10T11:30" }, now).starts,
  );
  // fim antes do início continua valendo
  assert.ok(validateWindow({ ...base, ends: "2026-11-10T10:00" }, now).ends);
  // janela iniciada: horário não é validado (só o motivo vai)
  assert.deepEqual(
    validateWindow({ ...base, ends: "2026-11-10T10:00", locked: true }, now),
    {},
  );
});

test("windowState vira em curso na virada do início", () => {
  const w = { starts_at: "2026-11-10T12:00:00Z", state: "scheduled" };
  assert.equal(windowState(w, Date.parse("2026-11-10T11:59:59Z")), "scheduled");
  assert.equal(windowState(w, Date.parse("2026-11-10T12:00:00Z")), "active");
  assert.equal(windowState({ ...w, state: "active" }, 0), "active");
  // bloco maintenance do monitor (detalhe) usa active
  assert.equal(
    windowState({ starts_at: w.starts_at, active: true }, 0),
    "active",
  );
});
