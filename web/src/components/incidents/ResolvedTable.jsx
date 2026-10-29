import { Link } from "react-router-dom";
import { Check } from "lucide-react";
import { Button } from "../ui/Button";
import { dayKey, dayLabel, formatDuration, timeRange } from "../../lib/format";

// Agrupa pelo dia local do fim: uma queda que vira a meia-noite fica no dia
// em que terminou. A lista já vem ordenada por ended_at DESC.
function groupByDay(incidents) {
  const groups = [];
  for (const i of incidents) {
    const key = dayKey(i.ended_at);
    const last = groups.at(-1);
    if (last?.key === key) last.items.push(i);
    else groups.push({ key, label: dayLabel(i.ended_at), items: [i] });
  }
  return groups;
}

function MonitorCell({ incident }) {
  const m = incident.monitor;
  return (
    <div className="flex items-center gap-3 min-w-0">
      <span className="flex-shrink-0 inline-flex items-center justify-center w-6 h-6 rounded-full bg-pg-up-bg text-pg-up">
        <Check size={14} aria-hidden="true" />
      </span>
      <div className="min-w-0">
        {m ? (
          <Link
            to={`/monitors/${m.id}`}
            className="block text-pg-ui font-medium text-pg-ardosia truncate no-underline hover:underline"
          >
            {m.name}
          </Link>
        ) : (
          <div className="text-pg-ui font-medium text-pg-ardosia truncate">Monitor removido</div>
        )}
        {m?.target && (
          <div className="font-pg-mono text-pg-caption text-pg-neblina truncate" title={m.target}>
            {m.target}
          </div>
        )}
      </div>
    </div>
  );
}

const cell = "px-5 py-3 align-middle";

export function ResolvedTable({ incidents, total, onLoadMore, loadingMore }) {
  const groups = groupByDay(incidents);

  return (
    <div className="rounded-pg-card border border-pg-borda bg-pg-branco overflow-hidden">
      {/* desktop: tabela */}
      <table className="hidden md:table w-full table-fixed border-collapse text-left">
        <caption className="sr-only">Incidentes resolvidos, agrupados pelo dia em que terminaram</caption>
        <colgroup>
          <col className="w-[34%]" />
          <col className="w-[200px]" />
          <col className="w-[100px]" />
          <col />
        </colgroup>
        <thead>
          <tr className="bg-pg-ciano-nevoa text-pg-caption text-pg-ciano-fundo">
            <th scope="col" className="px-5 py-3 font-medium">Monitor</th>
            <th scope="col" className="px-5 py-3 font-medium">Período</th>
            <th scope="col" className="px-5 py-3 font-medium text-right">Duração</th>
            <th scope="col" className="px-5 py-3 font-medium">Motivo da queda</th>
          </tr>
        </thead>
        {groups.map((g) => (
          <tbody key={g.key}>
            <tr className="border-t border-pg-borda bg-pg-nuvem">
              <th scope="rowgroup" colSpan={4} className="px-5 py-2 text-pg-caption font-normal text-pg-neblina">
                {g.label}
              </th>
            </tr>
            {g.items.map((i) => (
              <tr key={i.id} className="border-t border-pg-borda">
                <td className={cell}>
                  <MonitorCell incident={i} />
                </td>
                <td className={`${cell} text-pg-ui text-pg-ardosia whitespace-nowrap tabular-nums`}>
                  {timeRange(i.started_at, i.ended_at)}
                </td>
                <td className={`${cell} text-right text-pg-ui font-medium text-pg-ardosia whitespace-nowrap`}>
                  {formatDuration(i.duration_seconds)}
                </td>
                <td className={`${cell} font-pg-mono text-pg-caption text-pg-neblina`}>
                  <div className="truncate" title={i.last_error ?? undefined}>
                    {i.last_error ?? "—"}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        ))}
      </table>

      {/* mobile: lista, sem scroll horizontal */}
      <div className="md:hidden">
        {groups.map((g) => (
          <section key={g.key} aria-label={g.label}>
            <h3 className="m-0 px-4 py-2 bg-pg-nuvem border-b border-pg-borda text-pg-caption font-normal text-pg-neblina">
              {g.label}
            </h3>
            <ul className="list-none m-0 p-0">
              {g.items.map((i) => (
                <li key={i.id} className="px-4 py-3 border-b border-pg-borda">
                  <MonitorCell incident={i} />
                  <div className="mt-1.5 pl-9 text-pg-ui text-pg-ardosia tabular-nums">
                    {timeRange(i.started_at, i.ended_at)} · <strong className="font-medium">{formatDuration(i.duration_seconds)}</strong>
                  </div>
                  {i.last_error && (
                    <p className="mt-1 mb-0 pl-9 font-pg-mono text-pg-caption text-pg-neblina break-words line-clamp-2">
                      {i.last_error}
                    </p>
                  )}
                </li>
              ))}
            </ul>
          </section>
        ))}
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3 px-5 py-4 border-pg-borda md:border-t">
        <span className="text-pg-caption text-pg-neblina" aria-live="polite">
          Mostrando {incidents.length} de {total} {total === 1 ? "incidente resolvido" : "incidentes resolvidos"}
        </span>
        {incidents.length < total && (
          <Button variant="ghost" onClick={onLoadMore} disabled={loadingMore}>
            {loadingMore ? "Carregando…" : "Carregar mais"}
          </Button>
        )}
      </div>
    </div>
  );
}
