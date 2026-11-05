import { useId } from "react";
import { Link } from "react-router-dom";
import { ExternalLink, X } from "lucide-react";
import { formatDuration, startedLabel } from "../../lib/format";

// Card de incidente aberto: há quanto tempo está fora, alvo, erro inteiro e,
// para http, um atalho para abrir a URL. No mobile (< md), "Ver monitor" e
// "Abrir URL" ficam lado a lado no rodapé.
const mobileAction =
  "flex-1 inline-flex items-center justify-center gap-2 min-h-[44px] px-3 rounded-pg-control border border-pg-linha text-pg-ui font-medium text-pg-texto no-underline hover:bg-pg-hover focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano";

export function OpenIncidentCard({ incident, now }) {
  const titleId = useId();
  const monitor = incident.monitor;
  const name = monitor?.name ?? "Monitor removido";
  const down = formatDuration(
    (now - new Date(incident.started_at).getTime()) / 1000,
  );

  return (
    <article
      aria-labelledby={titleId}
      className="flex flex-col gap-4 md:flex-row p-4 md:p-5 rounded-pg-card border border-pg-down-borda bg-pg-superficie"
    >
      <div className="flex gap-4 flex-1 min-w-0">
        <span className="flex-shrink-0 inline-flex items-center justify-center w-10 h-10 rounded-full bg-pg-status-down-bg text-pg-status-down">
          <X size={18} aria-hidden="true" />
        </span>
        <div className="flex-1 min-w-0">
          <h3
            id={titleId}
            className="flex flex-col md:flex-row md:flex-wrap md:items-baseline gap-x-3 my-0 text-pg-ui"
          >
            {monitor ? (
              <Link
                to={`/monitors/${monitor.id}`}
                className="text-[17px] font-medium text-pg-texto break-words no-underline hover:underline"
              >
                {name}
              </Link>
            ) : (
              <span className="text-[17px] font-medium text-pg-texto break-words">
                {name}
              </span>
            )}
            <span className="text-pg-status-down">
              <span className="sr-only">, </span>Fora do ar há {down}
            </span>
          </h3>
          <p className="md:hidden text-pg-caption text-pg-texto-2">
            Desde {startedLabel(incident.started_at, now)}
          </p>
          {monitor?.target && (
            <div className="hidden md:block mt-0.5 font-pg-mono text-pg-caption text-pg-texto-2 break-all">
              {monitor.target}
            </div>
          )}
          {incident.last_error && (
            <pre className="-ml-14 md:ml-0 mt-3 mb-0 px-3 py-2.5 rounded-pg-control border border-pg-linha bg-pg-fundo font-pg-mono text-pg-caption text-pg-texto whitespace-pre-wrap break-words">
              {incident.last_error}
            </pre>
          )}
        </div>
      </div>

      <div className="hidden md:flex flex-col gap-3 md:items-end md:w-[200px] md:flex-shrink-0">
        <dl className="grid grid-cols-[auto_auto] gap-x-6 gap-y-0.5 my-0 text-pg-caption w-fit md:w-full">
          <dt className="text-pg-texto-2">Começou</dt>
          <dd className="m-0 text-right text-pg-texto whitespace-nowrap">
            {startedLabel(incident.started_at, now)}
          </dd>
          <dt className="text-pg-texto-2">Falhas seguidas</dt>
          <dd className="m-0 text-right text-pg-texto tabular-nums">
            {incident.failure_count.toLocaleString("pt-BR")}
          </dd>
        </dl>
        {monitor?.type === "http" && monitor.target && (
          <a
            href={monitor.target}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 w-fit px-4 py-2 rounded-pg-control border border-pg-linha text-pg-ui font-medium text-pg-texto no-underline hover:bg-pg-hover focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
          >
            <ExternalLink size={16} aria-hidden="true" />
            Abrir URL
            <span className="sr-only"> de {name} (abre em nova aba)</span>
          </a>
        )}
      </div>

      {monitor && (
        <div className="md:hidden flex gap-2">
          <Link to={`/monitors/${monitor.id}`} className={mobileAction}>
            Ver monitor
          </Link>
          {monitor.type === "http" && monitor.target && (
            <a
              href={monitor.target}
              target="_blank"
              rel="noopener noreferrer"
              className={mobileAction}
            >
              <ExternalLink size={16} aria-hidden="true" />
              Abrir URL
              <span className="sr-only"> de {name} (abre em nova aba)</span>
            </a>
          )}
        </div>
      )}
    </article>
  );
}
