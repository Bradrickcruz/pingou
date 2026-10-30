import { X } from "lucide-react";
import { formatDuration, startedLabel } from "../../lib/format";

// Cartão "Fora do ar há…" do detalhe no mobile (< md), a partir do incidente
// aberto. Sem incidente aberto na lista, não aparece.
export function DownAlert({ incident, now }) {
  const down = formatDuration((now - new Date(incident.started_at).getTime()) / 1000);
  const failures = incident.failure_count ?? 0;

  return (
    <div role="status" className="md:hidden mb-4 p-4 rounded-pg-card border border-[#EDC4C8] bg-pg-down-bg">
      <div className="flex gap-3">
        <span className="flex-shrink-0 inline-flex items-center justify-center w-8 h-8 rounded-full bg-pg-down text-pg-branco">
          <X size={18} aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <p className="text-pg-body font-medium text-pg-ardosia">Fora do ar há {down}</p>
          <p className="text-pg-caption text-pg-ardosia">
            Desde {startedLabel(incident.started_at, now)}
            {failures > 0 &&
              ` · ${failures.toLocaleString("pt-BR")} ${failures === 1 ? "falha seguida" : "falhas seguidas"}`}
          </p>
        </div>
      </div>
      {incident.last_error && (
        <pre className="mt-3 mb-0 px-3 py-2.5 rounded-pg-control border border-pg-borda bg-pg-branco font-pg-mono text-pg-caption text-pg-ardosia whitespace-pre-wrap break-words">
          {incident.last_error}
        </pre>
      )}
    </div>
  );
}
