import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import {
  THEME_KEY,
  THEME_ENABLED,
  normalize,
  getPreference,
  savePreference,
  systemIsDark,
  resolve,
  apply,
} from "./theme.js";

const memory = (initial = {}) => {
  const data = { ...initial };
  return {
    getItem: (k) => (k in data ? data[k] : null),
    setItem: (k, v) => {
      data[k] = String(v);
    },
    data,
  };
};
const broken = {
  getItem() {
    throw new Error("blocked");
  },
  setItem() {
    throw new Error("blocked");
  },
};

test("normalize: valor inválido ou ausente vira system", () => {
  for (const v of [null, undefined, "", "auto", "DARK", 1, {}]) {
    assert.equal(normalize(v), "system");
  }
  for (const v of ["light", "dark", "system"]) assert.equal(normalize(v), v);
});

test("getPreference: lê o salvo; sem nada, system", () => {
  assert.equal(getPreference(memory()), "system");
  assert.equal(getPreference(memory({ [THEME_KEY]: "dark" })), "dark");
  assert.equal(getPreference(memory({ [THEME_KEY]: "lixo" })), "system");
});

test("getPreference/savePreference: storage ausente ou com erro não quebra", () => {
  assert.equal(getPreference(null), "system");
  assert.equal(getPreference(broken), "system");
  assert.equal(savePreference("dark", null), "dark");
  assert.equal(savePreference("dark", broken), "dark");
});

test("savePreference: grava só valor válido", () => {
  const s = memory();
  assert.equal(savePreference("light", s), "light");
  assert.equal(s.data[THEME_KEY], "light");
  assert.equal(savePreference("xyz", s), "system");
  assert.equal(s.data[THEME_KEY], "system");
});

test("systemIsDark: segue o matchMedia e tolera falta dele", () => {
  assert.equal(
    systemIsDark(() => ({ matches: true })),
    true,
  );
  assert.equal(
    systemIsDark(() => ({ matches: false })),
    false,
  );
  assert.equal(systemIsDark(undefined), false);
  assert.equal(
    systemIsDark(() => {
      throw new Error("x");
    }),
    false,
  );
});

test("resolve (ligado): os 3 casos", () => {
  assert.equal(resolve("light", true, true), "light");
  assert.equal(resolve("dark", false, true), "dark");
  assert.equal(resolve("system", true, true), "dark");
  assert.equal(resolve("system", false, true), "light");
  assert.equal(resolve("lixo", true, true), "dark"); // inválido = system
});

test("resolve (desligado): sempre claro", () => {
  for (const p of ["light", "dark", "system"]) {
    for (const sys of [true, false])
      assert.equal(resolve(p, sys, false), "light");
  }
});

test("apply: data-theme e color-scheme", () => {
  const attrs = {};
  const root = {
    setAttribute: (k, v) => (attrs[k] = v),
    style: {},
  };
  assert.equal(apply("dark", root), "dark");
  assert.equal(attrs["data-theme"], "dark");
  assert.equal(root.style.colorScheme, "dark");
  assert.equal(apply("qualquer", root), "light");
  assert.equal(attrs["data-theme"], "light");
});

// O script inline do index.html repete a regra. Garante que não divergem.
test("index.html: mesma chave, mesma trava e roda antes de tudo", () => {
  const html = readFileSync(
    new URL("../../index.html", import.meta.url),
    "utf8",
  );
  assert.ok(html.includes(`"${THEME_KEY}"`), "chave do localStorage diferente");
  assert.ok(
    html.includes(`var ON = ${THEME_ENABLED};`),
    "trava `ON` do index.html difere de THEME_ENABLED",
  );
  const script = html.indexOf("<script>");
  assert.ok(
    script > -1 && script < html.indexOf("<title>"),
    "script inline fora do lugar",
  );
  assert.ok(
    script < html.indexOf('type="module"'),
    "script inline depois do bundle",
  );
});

// Executa o script do index.html e compara com resolve().
test("index.html: o script produz o mesmo tema que resolve()", () => {
  const html = readFileSync(
    new URL("../../index.html", import.meta.url),
    "utf8",
  );
  const code = html.match(/<script>([\s\S]*?)<\/script>/)[1];
  for (const on of [true, false]) {
    for (const stored of [null, "light", "dark", "system", "lixo"]) {
      for (const sys of [true, false]) {
        const attrs = {};
        const el = { setAttribute: (k, v) => (attrs[k] = v), style: {} };
        const win = {
          localStorage: { getItem: () => stored },
          matchMedia: () => ({ matches: sys }),
        };
        const src = code.replace(/var ON = (true|false);/, `var ON = ${on};`);
        new Function("window", "document", src)(win, { documentElement: el });
        assert.equal(
          attrs["data-theme"],
          resolve(stored, sys, on),
          `ON=${on} salvo=${stored} sistema_escuro=${sys}`,
        );
        assert.equal(el.style.colorScheme, attrs["data-theme"]);
      }
    }
  }
});
