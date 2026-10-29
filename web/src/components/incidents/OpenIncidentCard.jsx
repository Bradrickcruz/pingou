import { useId } from "react";
import { ExternalLink, X } from "lucide-react";
import { formatDuration, startedLabel } from "../../lib/format";

// Card de incidente aberto: há quanto tempo está fora, alvo, erro inteiro e,
// para http, um atalho para abrir a URL.
export function OpenIncidentCard({ incident, now }) {
  const titleId = useId();
  const monitor = incident.monitor;
  const name = monitor?.name ?? "Monitor removido";
  const down = formatDuration((now - new Date(incident.started_at).getTime()) / 1000);

  return (
    <article
      aria-labelledby={titleId}
      className="flex flex-col gap-4 md:flex-row p-5 rounded-pg-card border border-[#EDC4C8] bg-pg-branco"
    >
      <div className="flex gap-4 flex-1 min-w-0">
        <span className="flex-shrink-0 inline-flex items-center justify-center w-10 h-10 rounded-full bg-pg-down-bg text-pg-down">
          <X size={18} aria-hidden="true" />
        </span>
        <div className="flex-1 min-w-0">
          <h3 id={titleId} className="flex flex-wrap items-baseline gap-x-3 my-0 text-pg-ui">
            <span className="text-[17px] font-medium text-pg-ardosia break-words">{name}</span>
            <span className="text-pg-down">
              <span className="sr-only">, </span>Fora do ar há {down}
            </span>
          </h3>
          {monitor?.target && (
            <div className="mt-0.5 font-pg-mono text-pg-caption text-pg-neblina break-all">
              {monitor.target}
            </div>
          )}
          {incident.last_error && (
            <pre className="mt-3 mb-0 px-3 py-2.5 rounded-pg-control border border-pg-borda bg-pg-nuvem font-pg-mono text-pg-caption text-pg-ardosia whitespace-pre-wrap break-words">
              {incident.last_error}
            </pre>
          )}
        </div>
      </div>

      <div className="flex flex-col gap-3 pl-14 md:pl-0 md:items-end md:w-[200px] md:flex-shrink-0">
        <dl className="grid grid-cols-[auto_auto] gap-x-6 gap-y-0.5 my-0 text-pg-caption w-fit md:w-full">
          <dt className="text-pg-neblina">Começou</dt>
          <dd className="m-0 text-right text-pg-ardosia whitespace-nowrap">{startedLabel(incident.started_at, now)}</dd>
          <dt className="text-pg-neblina">Falhas seguidas</dt>
          <dd className="m-0 text-right text-pg-ardosia tabular-nums">
            {incident.failure_count.toLocaleString("pt-BR")}
          </dd>
        </dl>
        {monitor?.type === "http" && monitor.target && (
          <a
            href={monitor.target}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 w-fit px-4 py-2 rounded-pg-control border border-pg-borda text-pg-ui font-medium text-pg-ardosia no-underline hover:bg-pg-nuvem focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
          >
            <ExternalLink size={16} aria-hidden="true" />
            Abrir URL
            <span className="sr-only"> de {name} (abre em nova aba)</span>
          </a>
        )}
      </div>
    </article>
  );
}
