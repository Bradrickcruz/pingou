import { useId, useState } from "react";
import { MIN_RETENTION, RANGES } from "../../lib/ranges";

// Seletor de período. Os acima da retenção ficam com aria-disabled (não
// `disabled`) para receber foco e hover e mostrar o motivo no tooltip.
// ponytail: tooltip em CSS; trocar por popover clicável se pedirem link
// No toque não há hover nem :focus-visible: tocar no período bloqueado
// mostra o motivo até o botão perder o foco.
// No mobile (< md) ocupa a largura toda, com rótulos curtos ("7 d").
export function RangeTabs({ value, onChange, retentionDays = MIN_RETENTION }) {
  const baseId = useId();
  const [tip, setTip] = useState(null);
  return (
    <div
      role="group"
      aria-label="Período"
      className="flex w-full md:inline-flex md:w-auto p-1 gap-1 rounded-pg-card border border-pg-borda bg-pg-branco"
    >
      {RANGES.map((r, i) => {
        const on = value === r.key;
        const off = r.days > retentionDays;
        const tipId = `${baseId}-${r.key}`;
        const last = i === RANGES.length - 1;
        return (
          <div key={r.key} className="relative group flex-1 md:flex-none">
            <button
              type="button"
              aria-pressed={on}
              aria-disabled={off || undefined}
              aria-describedby={off ? tipId : undefined}
              onClick={() => (off ? setTip((t) => (t === r.key ? null : r.key)) : onChange(r.key))}
              onBlur={() => setTip(null)}
              className={`w-full md:w-auto min-h-[40px] md:min-h-0 px-2 md:px-4 py-1.5 rounded-pg-control text-pg-ui whitespace-nowrap transition-colors duration-150 focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
                on
                  ? "bg-pg-ardosia text-pg-branco"
                  : off
                    ? "text-pg-neblina opacity-60 cursor-not-allowed line-through md:no-underline"
                    : "text-pg-ardosia hover:bg-pg-nuvem"
              }`}
            >
              <span className="md:hidden">{r.short}</span>
              <span className="hidden md:inline">{r.label}</span>
            </button>
            {off && (
              <span
                id={tipId}
                role="tooltip"
                className={`invisible opacity-0 group-hover:visible group-hover:opacity-100 group-has-[:focus-visible]:visible group-has-[:focus-visible]:opacity-100 transition-opacity absolute top-full mt-2 z-30 w-60 ${tip === r.key ? "!visible !opacity-100" : ""} px-3 py-2 rounded-pg-control bg-pg-ardosia text-pg-branco text-pg-caption text-left shadow-[0_4px_16px_rgba(30,43,51,0.16)] ${
                  last ? "right-0" : "left-1/2 -translate-x-1/2"
                }`}
              >
                Os checks são guardados por {retentionDays} dias. Aumente a retenção em Configurações para ver{" "}
                {r.label}.
              </span>
            )}
          </div>
        );
      })}
    </div>
  );
}
