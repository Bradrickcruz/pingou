import { useId, useState } from "react";

// Unidades do campo "Outro…"; o valor é sempre em segundos.
const UNITS = [
  { label: "segundos", factor: 1 },
  { label: "minutos", factor: 60 },
  { label: "horas", factor: 3600 },
  { label: "dias", factor: 86400 },
];

function bestUnit(seconds) {
  return [...UNITS].reverse().find((u) => seconds % u.factor === 0 && seconds >= u.factor) ?? UNITS[0];
}

const chip = (active) =>
  `px-4 py-2 rounded-full border text-pg-ui transition-colors duration-150 focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
    active
      ? "bg-pg-ciano-nevoa border-pg-ciano-fundo text-pg-ciano-fundo"
      : "bg-pg-branco border-pg-borda text-pg-ardosia hover:bg-pg-nuvem"
  }`;

// ChipGroup: valores prontos + "Outro…" com número e unidade.
// min/max (segundos) só orientam o input; quem valida de verdade é a API.
export function ChipGroup({ label, options, value, onChange, min, max }) {
  const labelId = useId();
  const inputId = useId();
  const isPreset = options.some((o) => o.value === value);
  const [custom, setCustom] = useState(!isPreset);
  const [unit, setUnit] = useState(() => bestUnit(value));

  const pick = (v) => {
    setCustom(false);
    onChange(v);
  };

  const openCustom = () => {
    setCustom(true);
    setUnit(bestUnit(value));
  };

  return (
    <div className="mb-5">
      <div id={labelId} className="mb-2 text-pg-ui font-medium text-pg-ardosia">
        {label}
      </div>
      <div role="group" aria-labelledby={labelId} className="flex flex-wrap gap-2">
        {options.map((o) => (
          <button
            key={o.value}
            type="button"
            aria-pressed={!custom && value === o.value}
            onClick={() => pick(o.value)}
            className={chip(!custom && value === o.value)}
          >
            {o.label}
          </button>
        ))}
        <button type="button" aria-pressed={custom} onClick={openCustom} className={chip(custom)}>
          Outro…
        </button>
      </div>

      {custom && (
        <div className="flex gap-2 mt-2.5 items-center">
          <label htmlFor={inputId} className="sr-only">
            {label} (valor)
          </label>
          <input
            id={inputId}
            type="number"
            inputMode="numeric"
            min={Math.max(0, Math.ceil((min ?? 0) / unit.factor))}
            max={max ? Math.floor(max / unit.factor) : undefined}
            value={Math.round(value / unit.factor)}
            onChange={(e) => onChange(Math.max(0, Number(e.target.value)) * unit.factor)}
            className="w-28 px-3 py-2 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda"
          />
          <select
            aria-label={`${label} (unidade)`}
            value={unit.factor}
            onChange={(e) => {
              const next = UNITS.find((u) => u.factor === Number(e.target.value));
              onChange(Math.round(value / unit.factor) * next.factor);
              setUnit(next);
            }}
            className="px-3 py-2 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda"
          >
            {UNITS.map((u) => (
              <option key={u.factor} value={u.factor}>
                {u.label}
              </option>
            ))}
          </select>
        </div>
      )}
    </div>
  );
}
