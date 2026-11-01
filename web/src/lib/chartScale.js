// Escala Y "redonda": passo em {1,2,5}×10ⁿ, no máximo 5 intervalos.
// peak 350 → 0..400 (passo 100); 2300 → 0..2500 (passo 500); piso evita eixo achatado.
export function niceScale(peak, floor = 100) {
  const p = Math.max(Number.isFinite(peak) ? peak : 0, floor);
  let step;
  for (let mag = 1; !step; mag *= 10)
    step = [1, 2, 5].map((k) => k * mag).find((v) => v * 5 >= p);
  const max = Math.ceil(p / step) * step;
  const ticks = [];
  for (let t = 0; t <= max; t += step) ticks.push(t);
  return { max, ticks };
}
