// node --test src/lib
import { test } from "node:test";
import assert from "node:assert/strict";
import { niceScale } from "./chartScale.js";

test("respeita o piso", () => {
  assert.equal(niceScale(0).max, 100);
  assert.equal(niceScale(NaN).max, 100);
  assert.equal(niceScale(70).max, 100);
  assert.equal(niceScale(100).max, 100);
  assert.deepEqual(niceScale(70).ticks, [0, 20, 40, 60, 80, 100]);
});

test("exemplos do pedido", () => {
  assert.equal(niceScale(350).max, 400);
  assert.equal(niceScale(2300).max, 2500);
  assert.equal(niceScale(402).max, 500);
});

test("passo salta para o próximo {1,2,5}×10ⁿ", () => {
  // 101/5 = 20,2 > 20 → passo 50
  assert.deepEqual(niceScale(101).ticks, [0, 50, 100, 150]);
});

test("piso 1 (contagem de pings)", () => {
  assert.deepEqual(niceScale(3, 1).ticks, [0, 1, 2, 3]);
  assert.deepEqual(niceScale(0, 1).ticks, [0, 1]);
});
