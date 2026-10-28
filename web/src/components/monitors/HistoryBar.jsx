// Últimos N checks como barras: verde (no ar), vermelho (fora), cinza (sem dados).
// history vem do mais antigo para o mais recente; completa à direita com cinza.
export function HistoryBar({ history = [], size = 30, dimmed = false }) {
  const bars = history.slice(-size);
  const empty = size - bars.length;
  const ok = bars.filter(Boolean).length;
  const label = bars.length
    ? `Últimos ${bars.length} checks: ${ok} no ar, ${bars.length - ok} fora do ar`
    : "Sem checks ainda";

  return (
    <div
      role="img"
      aria-label={label}
      className={`flex items-stretch gap-[2px] h-6 w-[198px] ${dimmed ? "opacity-50" : ""}`}
    >
      {bars.map((success, i) => (
        <span
          key={i}
          className={`flex-1 rounded-[1px] ${success ? "bg-pg-up" : "bg-pg-down"}`}
        />
      ))}
      {Array.from({ length: empty }, (_, i) => (
        <span key={`e${i}`} className="flex-1 rounded-[1px] bg-pg-borda" />
      ))}
    </div>
  );
}
