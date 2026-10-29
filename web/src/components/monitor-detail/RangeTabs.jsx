import { useId } from "react";
import { MIN_RETENTION, RANGES } from "../../lib/ranges";

// Seletor de período. Os acima da retenção ficam com aria-disabled (não
// `disabled`) para receber foco e hover e mostrar o motivo no tooltip.
// ponytail: tooltip em CSS; trocar por popover clicável se pedirem link
export function RangeTabs({ value, onChange, retentionDays = MIN_RETENTION }) {
  const baseId = useId();
  return (
    <div
      role="group"
      aria-label="Período"
      className="inline-flex p-1 gap-1 rounded-pg-card border border-pg-borda bg-pg-branco"
    >
      {RANGES.map((r, i) => {
        const on = value === r.key;
        const off = r.days > retentionDays;
        const tipId = `${baseId}-${r.key}`;
        const last = i === RANGES.length - 1;
        return (
          <div key={r.key} className="relative group">
            <button
              type="button"
              aria-pressed={on}
              aria-disabled={off || undefined}
              aria-describedby={off ? tipId : undefined}
              onClick={() => !off && onChange(r.key)}
              className={`px-3 sm:px-4 py-1.5 rounded-pg-control text-pg-ui whitespace-nowrap transition-colors duration-150 focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
                on
                  ? "bg-pg-ardosia text-pg-branco"
                  : off
                    ? "text-pg-neblina opacity-60 cursor-not-allowed"
                    : "text-pg-ardosia hover:bg-pg-nuvem"
              }`}
            >
              {r.label}
            </button>
            {off && (
              <span
                id={tipId}
                role="tooltip"
                className={`invisible opacity-0 group-hover:visible group-hover:opacity-100 group-has-[:focus-visible]:visible group-has-[:focus-visible]:opacity-100 transition-opacity absolute top-full mt-2 z-30 w-60 px-3 py-2 rounded-pg-control bg-pg-ardosia text-pg-branco text-pg-caption text-left shadow-[0_4px_16px_rgba(30,43,51,0.16)] ${
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
