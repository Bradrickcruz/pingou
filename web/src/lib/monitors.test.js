// node --test src/lib
import { test } from "node:test";
import assert from "node:assert/strict";
import { cardFooter, heartbeatExpectation, matchesQuery } from "./monitors.js";

const now = Date.parse("2026-10-30T12:00:00Z");
const base = {
  type: "http",
  enabled: true,
  uptime_24h: 99.2,
  last_latency_ms: 214,
  last_checked_at: "2026-10-30T11:59:19Z",
};

test("cardFooter UP com latência", () => {
  assert.equal(
    cardFooter({ ...base, current_state: "UP" }, now),
    "99,2% em 24 h · 214 ms · há 41 s",
  );
});

test("cardFooter UP sem latência nem uptime", () => {
  assert.equal(
    cardFooter(
      { ...base, current_state: "UP", last_latency_ms: null, uptime_24h: null },
      now,
    ),
    "há 41 s",
  );
});

test("cardFooter DOWN", () => {
  const m = {
    ...base,
    current_state: "DOWN",
    uptime_24h: 90,
    down_since: "2026-10-29T08:00:00Z",
  };
  assert.equal(cardFooter(m, now), "Fora há 1 dia e 4 h · 90,0% em 24 h");
  assert.equal(
    cardFooter({ ...m, down_since: null }, now),
    "Fora do ar · 90,0% em 24 h",
  );
});

test("cardFooter UNKNOWN", () => {
  assert.equal(
    cardFooter({ ...base, current_state: "UNKNOWN" }, now),
    "Aguardando o primeiro check",
  );
});

test("cardFooter heartbeat", () => {
  const hb = {
    type: "heartbeat",
    enabled: true,
    current_state: "UP",
    uptime_24h: 100,
    last_latency_ms: 0,
    last_ping_at: "2026-10-30T11:00:00Z",
    interval_seconds: 86400,
  };
  assert.equal(cardFooter(hb, now), "100% em 24 h · ping há 1 h");
  assert.equal(heartbeatExpectation(86400), "espera um ping por dia");
  assert.equal(heartbeatExpectation(7200), "espera um ping a cada 2 h");
});

test("cardFooter PAUSED", () => {
  assert.equal(
    cardFooter({ ...base, enabled: false, current_state: "DOWN" }, now),
    "Pausado",
  );
});

const http = {
  type: "http",
  name: "Serviço de Pagamentos",
  url: "https://api.exemplo.com/v1",
};
const tcp = {
  type: "tcp",
  name: "Banco",
  config: { host: "db.interno", port: 5432 },
};
const beat = {
  type: "heartbeat",
  name: "Backup noturno",
  interval_seconds: 86400,
};

test("matchesQuery ignora acento e caixa", () => {
  assert.equal(matchesQuery(http, "servico"), true);
  assert.equal(matchesQuery(http, "PAGAMENTOS"), true);
});

test("matchesQuery acha por parte do nome", () => {
  assert.equal(matchesQuery(http, "gam"), true);
  assert.equal(matchesQuery(http, "xyz"), false);
});

test("matchesQuery acha por URL e host:porta", () => {
  assert.equal(matchesQuery(http, "exemplo.com/v1"), true);
  assert.equal(matchesQuery(tcp, "db.interno:5432"), true);
  assert.equal(matchesQuery(tcp, "5432"), true);
});

test("matchesQuery em heartbeat usa só o nome", () => {
  assert.equal(matchesQuery(beat, "backup"), true);
  assert.equal(matchesQuery(beat, "espera"), false);
});

test("matchesQuery com busca vazia casa tudo", () => {
  assert.equal(matchesQuery(http, ""), true);
  assert.equal(matchesQuery(http, "   "), true);
});
