// node --test src/lib
import { test, afterEach } from "node:test";
import assert from "node:assert/strict";
import { certStatus, hasCert, isHttps } from "./cert.js";
import { setTimeZone } from "./timezone.js";

afterEach(() => setTimeZone("UTC"));

const now = Date.parse("2026-11-01T12:00:00Z");
const day = 86_400_000;
const at = (ms) => new Date(now + ms).toISOString();
const http = (extra) => ({
  type: "http",
  url: "https://api.test",
  cert_expires_at: null,
  ...extra,
});

test("isHttps e hasCert", () => {
  assert.equal(isHttps("https://x"), true);
  assert.equal(isHttps(" HTTPS://x"), true);
  assert.equal(isHttps("http://x"), false);
  assert.equal(isHttps(""), false);
  assert.equal(hasCert(http()), true);
  assert.equal(hasCert(http({ url: "http://api.test" })), false);
  assert.equal(hasCert({ type: "tcp", config: {} }), false);
  assert.equal(hasCert({ type: "heartbeat" }), false);
});

test("sem certificado: null fora de http https", () => {
  assert.equal(certStatus({ type: "tcp" }, now), null);
  assert.equal(certStatus(http({ url: "http://api.test" }), now), null);
});

test("aguardando o primeiro check", () => {
  assert.deepEqual(certStatus(http(), now), {
    level: "pending",
    text: "Aguardando o primeiro check",
  });
});

test("níveis por dias restantes", () => {
  const cases = [
    [52 * day, "ok", "Expira em 23/12/2026 (daqui a 52 dias)"],
    [30 * day + 1000, "ok", "Expira em 01/12/2026 (daqui a 30 dias)"],
    [30 * day, "warn", "Expira em 01/12/2026 (daqui a 30 dias)"],
    [8 * day, "warn", "Expira em 09/11/2026 (daqui a 8 dias)"],
    [7 * day, "alert", "Expira em 08/11/2026 (daqui a 7 dias)"],
    [1 * day + 1000, "alert", "Expira em 02/11/2026 (daqui a 1 dia)"],
    [3600_000, "alert", "Expira em 01/11/2026 (hoje)"],
    [-1000, "expired", "Expirou em 01/11/2026"],
  ];
  for (const [ms, level, text] of cases) {
    assert.deepEqual(
      certStatus(http({ cert_expires_at: at(ms) }), now),
      { level, text },
      `${ms} ms`,
    );
  }
});

test("data no fuso das Configurações", () => {
  setTimeZone("America/Sao_Paulo");
  // 02:00 UTC do dia 23 ainda é dia 22 em São Paulo
  const m = http({ cert_expires_at: "2026-12-23T02:00:00Z" });
  assert.match(certStatus(m, now).text, /^Expira em 22\/12\/2026/);
});
