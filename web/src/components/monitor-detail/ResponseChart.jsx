import { useEffect, useRef, useState } from "react";
import {
  bucketLabel,
  clock,
  dayMonth,
  formatDuration,
  plural,
} from "../../lib/format";

const W = 1000; // viewBox; o SVG estica na largura (preserveAspectRatio="none")
const H = 200;

// Escala Y em 3 passos "redondos" de 100 ms: pico 402 ms → 0 / 200 / 400 / 600.
function scaleOf(values) {
  const peak = Math.max(0, ...values);
  const step = Math.max(100, Math.ceil(peak / 300) * 100);
  return { max: step * 3, ticks: [0, step, step * 2, step * 3] };
}

// Linha da série; balde sem média quebra o traço (novo "M") em vez de ligar
// por cima do buraco.
function linePath(series, xOf, yOf) {
  let d = "";
  let pen = false;
  series.forEach((b, i) => {
    if (b.avg_latency_ms == null) {
      pen = false;
      return;
    }
    d += `${pen ? "L" : "M"}${xOf(i).toFixed(1)},${yOf(b.avg_latency_ms).toFixed(1)}`;
    pen = true;
  });
  return d;
}

// Marcas do eixo X em horários "redondos" no fuso local: a cada 4 h na janela
// de 24 h; à meia-noite a cada 1, 2 ou 5 dias nas maiores (~6 a 7 rótulos).
function axisTicks(since, until) {
  const start = new Date(since);
  const end = new Date(until);
  const days = (end - start) / 86_400_000;
  const ticks = [];
  if (days <= 1) {
    const t = new Date(start);
    t.setMinutes(0, 0, 0);
    t.setHours(
      Math.ceil((t < start ? t.getHours() + 1 : t.getHours()) / 4) * 4,
    );
    for (; t < end; t.setHours(t.getHours() + 4))
      ticks.push({ t: new Date(t), label: clock(t) });
    return ticks;
  }
  const step = days <= 7 ? 1 : days <= 14 ? 2 : 5;
  const t = new Date(start);
  t.setHours(0, 0, 0, 0);
  if (t < start) t.setDate(t.getDate() + 1);
  for (; t < end; t.setDate(t.getDate() + step))
    ticks.push({ t: new Date(t), label: dayMonth(t) });
  return ticks;
}

function isDown(bucket, bucketSeconds, down) {
  const start = new Date(bucket.t).getTime();
  const end = start + bucketSeconds * 1000;
  return down.some(
    (w) =>
      new Date(w.start).getTime() < end && new Date(w.end).getTime() > start,
  );
}

// Índices dos baldes em trechos "sem dados": sequência de baldes sem média
// (fora de queda) mais longa que 3 intervalos de check.
const GAP_INTERVALS = 3;
function gapBuckets(series, bucketSeconds, intervalSeconds, down) {
  const gaps = new Set();
  let run = [];
  const flush = () => {
    if (run.length * bucketSeconds > GAP_INTERVALS * intervalSeconds)
      run.forEach((i) => gaps.add(i));
    run = [];
  };
  series.forEach((b, i) => {
    if (b.avg_latency_ms == null && !isDown(b, bucketSeconds, down))
      run.push(i);
    else flush();
  });
  flush();
  return gaps;
}

function summary(series, down) {
  const values = series.map((b) => b.avg_latency_ms).filter((v) => v != null);
  if (values.length === 0) return "Tempo de resposta: sem dados neste período.";
  const avg = Math.round(values.reduce((a, b) => a + b, 0) / values.length);
  const quedas = down.length
    ? plural(down.length, "período fora do ar", "períodos fora do ar")
    : "nenhum período fora do ar";
  return `Tempo de resposta médio ${avg} ms, pico ${Math.max(...values)} ms, ${quedas}.`;
}

// Gráfico de tempo de resposta (SVG à mão) com tooltip por ponto: mouse,
// toque (toque fixa o ponto, tocar fora limpa) e setas do teclado.
export function ResponseChart({ stats, now, intervalSeconds }) {
  const { series, down, since, until, bucket_seconds: bucketSeconds } = stats;
  const n = series.length;
  const gaps = gapBuckets(series, bucketSeconds, intervalSeconds, down);
  const [active, setActive] = useState(null);
  const plotRef = useRef(null);

  const span = new Date(until) - new Date(since);
  const xOf = (i) => ((i + 0.5) / n) * W;
  const xAt = (iso) => ((new Date(iso) - new Date(since)) / span) * W;
  const { max, ticks } = scaleOf(series.map((b) => b.avg_latency_ms ?? 0));
  const yOf = (ms) => H - (ms / max) * H;

  // toque: tocar fora do gráfico limpa o ponto fixado
  useEffect(() => {
    if (active == null) return;
    const onDown = (e) => {
      if (!plotRef.current?.contains(e.target)) setActive(null);
    };
    document.addEventListener("pointerdown", onDown);
    return () => document.removeEventListener("pointerdown", onDown);
  }, [active]);

  const hasData = series.some((b) => b.avg_latency_ms != null);
  if (!hasData) {
    return (
      <p className="m-0 py-16 text-center text-pg-ui text-pg-neblina">
        Sem dados de resposta neste período.
      </p>
    );
  }

  const pick = (e) => {
    const rect = plotRef.current.getBoundingClientRect();
    const i = Math.floor(((e.clientX - rect.left) / rect.width) * n);
    setActive(Math.min(n - 1, Math.max(0, i)));
  };

  const onKeyDown = (e) => {
    const cur = active ?? n - 1;
    const moves = {
      ArrowLeft: cur - 1,
      ArrowRight: cur + 1,
      Home: 0,
      End: n - 1,
    };
    if (e.key in moves) {
      e.preventDefault();
      setActive(Math.min(n - 1, Math.max(0, moves[e.key])));
    } else if (e.key === "Escape") {
      setActive(null);
    }
  };

  const b = active != null ? series[active] : null;
  const pointDown = b && isDown(b, bucketSeconds, down);
  const leftPct = active != null ? (xOf(active) / W) * 100 : 0;
  const tipText = b
    ? [
        bucketLabel(b.t, bucketSeconds, now),
        b.avg_latency_ms != null
          ? `média ${b.avg_latency_ms} ms`
          : "sem resposta",
        `${plural(b.checks, "check", "checks")}, ${plural(b.failures, "falha", "falhas")}`,
        pointDown ? "fora do ar" : null,
      ]
        .filter(Boolean)
        .join(", ")
    : "";

  return (
    <div className="flex gap-3 pb-6">
      {/* rótulos do eixo Y em HTML: texto não deforma com o SVG */}
      <div
        className="relative w-12 flex-shrink-0 h-[200px] text-pg-caption text-pg-neblina tabular-nums"
        aria-hidden="true"
      >
        {ticks.map((t) => (
          <span
            key={t}
            className="absolute right-0 -translate-y-1/2 whitespace-nowrap"
            style={{ top: `${(yOf(t) / H) * 100}%` }}
          >
            {t === 0 ? "0" : `${t} ms`}
          </span>
        ))}
      </div>

      <div
        ref={plotRef}
        tabIndex={0}
        role="group"
        aria-label={`${summary(series, down)} Use as setas para ler cada intervalo.`}
        onPointerMove={(e) => e.pointerType === "mouse" && pick(e)}
        onPointerDown={pick}
        onPointerLeave={(e) => e.pointerType === "mouse" && setActive(null)}
        onKeyDown={onKeyDown}
        onFocus={() => setActive((a) => a ?? n - 1)}
        onBlur={() => setActive(null)}
        className="relative flex-1 min-w-0 h-[200px] touch-pan-y rounded-sm focus-visible:outline-2 focus-visible:outline focus-visible:outline-offset-4 focus-visible:outline-pg-ciano"
      >
        <svg
          viewBox={`0 0 ${W} ${H}`}
          preserveAspectRatio="none"
          className="absolute inset-0 w-full h-full overflow-visible"
          aria-hidden="true"
        >
          {ticks.map((t) => (
            <line
              key={t}
              x1="0"
              x2={W}
              y1={yOf(t)}
              y2={yOf(t)}
              stroke="var(--pg-borda)"
              strokeWidth="1"
              vectorEffect="non-scaling-stroke"
            />
          ))}
          {/* faixa cinza nos trechos sem dados (legenda "sem dados") */}
          {series.map((s, i) =>
            gaps.has(i) ? (
              <rect
                key={s.t}
                x={(i / n) * W}
                y="0"
                width={W / n}
                height={H}
                fill="var(--pg-unknown-bg)"
              />
            ) : null,
          )}
          {down.map((w) => {
            const x = Math.max(0, xAt(w.start));
            // mínimo de ~2 px (num gráfico de ~1000 px) para quedas curtas aparecerem em 30 d
            const width = Math.max(2, Math.min(W, xAt(w.end)) - x);
            return (
              <g key={w.start}>
                <rect
                  x={x}
                  y="0"
                  width={width}
                  height={H}
                  fill="var(--pg-down-bg)"
                />
                <line
                  x1={x}
                  x2={x}
                  y1="0"
                  y2={H}
                  stroke="var(--pg-down)"
                  strokeWidth="1"
                  strokeDasharray="3 3"
                  vectorEffect="non-scaling-stroke"
                />
              </g>
            );
          })}
          <path
            d={linePath(series, xOf, yOf)}
            fill="none"
            stroke="var(--pg-ciano)"
            strokeWidth="2"
            strokeLinejoin="round"
            vectorEffect="non-scaling-stroke"
          />
        </svg>

        {b && (
          <>
            <div
              className="absolute top-0 bottom-0 w-px bg-pg-borda pointer-events-none"
              style={{ left: `${leftPct}%` }}
            />
            {b.avg_latency_ms != null && (
              <div
                className="absolute w-2.5 h-2.5 -ml-[5px] -mt-[5px] rounded-full bg-pg-ciano border-2 border-pg-branco pointer-events-none"
                style={{
                  left: `${leftPct}%`,
                  top: `${(yOf(b.avg_latency_ms) / H) * 100}%`,
                }}
              />
            )}
            <div
              aria-hidden="true"
              className={`absolute top-2 z-10 w-max max-w-[220px] px-3 py-2 rounded-pg-control border border-pg-borda bg-pg-branco shadow-[0_4px_16px_rgba(30,43,51,0.08)] text-pg-caption text-pg-ardosia pointer-events-none ${
                leftPct > 50 ? "-translate-x-[calc(100%+8px)]" : "translate-x-2"
              }`}
              style={{ left: `${leftPct}%` }}
            >
              <div className="font-medium mb-1">
                {bucketLabel(b.t, bucketSeconds, now)}
              </div>
              <div>
                {b.avg_latency_ms != null ? (
                  <>
                    Média{" "}
                    <strong className="font-medium tabular-nums">
                      {b.avg_latency_ms} ms
                    </strong>
                  </>
                ) : (
                  <span className="text-pg-neblina">Sem resposta</span>
                )}
              </div>
              <div className="text-pg-neblina tabular-nums">
                {plural(b.checks, "check", "checks")} ·{" "}
                <span className={b.failures > 0 ? "text-pg-down" : ""}>
                  {plural(b.failures, "falha", "falhas")}
                </span>
              </div>
              {pointDown && (
                <div className="mt-1 font-medium text-pg-down">Fora do ar</div>
              )}
            </div>
          </>
        )}
        {/* eixo X: rótulos que encostariam nas bordas do card ficam de fora */}
        <div
          className="absolute left-0 right-0 top-full mt-2 h-4 text-pg-caption text-pg-neblina tabular-nums"
          aria-hidden="true"
        >
          {axisTicks(since, until)
            .map(({ t, label }) => ({ label, pct: (xAt(t) / W) * 100 }))
            .filter(({ pct }) => pct > 4 && pct < 96)
            .map(({ label, pct }, i) => (
              <span
                key={label}
                className={`absolute -translate-x-1/2 whitespace-nowrap ${i % 2 ? "hidden sm:inline" : ""}`}
                style={{ left: `${pct}%` }}
              >
                {label}
              </span>
            ))}
        </div>
        <div role="status" aria-live="polite" className="sr-only">
          {tipText}
        </div>
      </div>
    </div>
  );
}

export function ChartLegend({ bucketSeconds }) {
  return (
    <div
      className="flex gap-3 md:gap-4 text-pg-caption text-pg-neblina"
      aria-hidden="true"
    >
      <span className="inline-flex items-center gap-1.5">
        <span className="w-3.5 h-0.5 rounded bg-pg-ciano" />{" "}
        <span className="md:hidden">média</span>
        <span className="hidden md:inline">
          {bucketSeconds
            ? `média a cada ${formatDuration(bucketSeconds)}`
            : "média por intervalo"}
        </span>
      </span>
      <span className="inline-flex items-center gap-1.5">
        <span className="w-3 h-3 rounded-[2px] border border-pg-down bg-pg-down-bg" />{" "}
        fora do ar
      </span>
      <span className="inline-flex items-center gap-1.5">
        <span className="w-3 h-3 rounded-[2px] border border-pg-neblina bg-pg-unknown-bg" />{" "}
        sem dados
      </span>
    </div>
  );
}
