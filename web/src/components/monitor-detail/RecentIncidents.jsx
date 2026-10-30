import { Check, X } from "lucide-react";
import { formatDuration, whenLabel } from "../../lib/format";

// Os 3 incidentes mais recentes do monitor, independentes do período escolhido.
// No mobile (< md), só quando começou e a duração (ou "em aberto").
export function RecentIncidents({ incidents, now }) {
  if (incidents.length === 0) {
    return (
      <p className="m-0 py-8 px-4 text-center text-pg-ui text-pg-neblina rounded-pg-card border border-dashed border-pg-borda">
        Nenhum incidente registrado.
      </p>
    );
  }
  return (
    <ul className="list-none m-0 p-0 rounded-pg-card border border-pg-borda bg-pg-branco">
      {incidents.map((i) => (
        <li
          key={i.id}
          className="flex gap-3 px-4 py-3 border-t border-pg-borda first:border-t-0"
        >
          <span className="sr-only">
            {i.open ? "Em andamento" : "Resolvido"}
          </span>
          <span
            aria-hidden="true"
            className={`hidden md:inline-flex flex-shrink-0 items-center justify-center w-6 h-6 rounded-full ${
              i.open ? "bg-pg-down-bg text-pg-down" : "bg-pg-up-bg text-pg-up"
            }`}
          >
            {i.open ? <X size={14} /> : <Check size={14} />}
          </span>
          <div className="flex-1 min-w-0">
            <div className="flex justify-between gap-3 text-pg-ui text-pg-ardosia">
              <span className="font-medium">
                {whenLabel(i.started_at, now)}
              </span>
              <span
                className={`whitespace-nowrap ${i.open ? "text-pg-down" : ""}`}
              >
                {i.open ? "em aberto" : formatDuration(i.duration_seconds)}
              </span>
            </div>
            {i.last_error && (
              <div
                className="hidden md:block font-pg-mono text-pg-caption text-pg-neblina truncate"
                title={i.last_error}
              >
                {i.last_error}
              </div>
            )}
          </div>
        </li>
      ))}
    </ul>
  );
}
