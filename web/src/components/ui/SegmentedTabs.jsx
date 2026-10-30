// Filtro de opções exclusivas: no mobile (< md) é um controle segmentado de
// largura total; no desktop, a linha de chips de sempre. Só muda a classe.
export function SegmentedTabs({ label, options, value, onChange }) {
  return (
    <div
      role="group"
      aria-label={label}
      className="flex w-full gap-1 p-1 rounded-pg-control border border-pg-borda bg-pg-branco md:w-auto md:flex-wrap md:gap-2 md:p-0 md:border-0 md:bg-transparent"
    >
      {options.map((o) => {
        const on = value === o.key;
        return (
          <button
            key={o.key}
            type="button"
            aria-pressed={on}
            onClick={() => onChange(o.key)}
            className={`flex-1 md:flex-none min-h-[44px] md:min-h-0 px-2 md:px-4 py-2 rounded-[6px] md:rounded-full md:border text-pg-ui whitespace-nowrap transition-colors duration-150 focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
              on
                ? "bg-pg-ardosia md:border-pg-ardosia text-pg-branco"
                : "text-pg-ardosia md:bg-pg-branco md:border-pg-borda hover:bg-pg-nuvem"
            }`}
          >
            {o.label}
            {o.count != null && (
              <>
                <span className="md:hidden"> {o.count}</span>
                <span className="hidden md:inline"> ({o.count})</span>
              </>
            )}
          </button>
        );
      })}
    </div>
  );
}
