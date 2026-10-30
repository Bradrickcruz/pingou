import { Check, X } from "lucide-react";
import { checkTime } from "../../lib/format";

// Detalhe técnico do check: erro, código HTTP ou o que o tipo verifica.
function detailOf(c, type) {
  if (!c.success && c.error_message) return c.error_message;
  if (type === "heartbeat") return "ping recebido";
  if (type === "tcp") return c.success ? "conexão aberta" : "falha";
  return c.status_code ? `HTTP ${c.status_code}` : "—";
}

const cell = "px-4 py-2.5 align-middle";

// No mobile (< md), a primeira página mostra só 5 dos checks carregados.
const COMPACT = 5;

// Últimos checks do monitor, do mais recente ao mais antigo. "ver mais"
// acrescenta mais 20 na própria página. No mobile vira lista (hora, status,
// resposta e o erro embaixo); compact = ainda na primeira página.
export function ChecksTable({ checks, total, type, now, onMore, loadingMore, compact = false }) {
  const heartbeat = type === "heartbeat";
  const mobile = compact ? checks.slice(0, COMPACT) : checks;

  if (checks.length === 0) {
    return (
      <p className="m-0 py-10 px-4 text-center text-pg-ui text-pg-neblina rounded-pg-card border border-dashed border-pg-borda">
        Aguardando o primeiro check.
      </p>
    );
  }

  return (
    <div className="rounded-pg-card border border-pg-borda bg-pg-branco overflow-hidden">
      <ul className="md:hidden list-none m-0 p-0">
        {mobile.map((c) => (
          <li key={c.id} className="px-4 py-2.5 border-t border-pg-borda first:border-t-0">
            <div className="flex items-center gap-3">
              <span className="text-pg-ui text-pg-ardosia tabular-nums">{checkTime(c.checked_at, now)}</span>
              <span className={`inline-flex items-center gap-1 text-pg-caption font-medium ${c.success ? "text-pg-up" : "text-pg-down"}`}>
                {c.success ? <Check size={12} aria-hidden="true" /> : <X size={12} aria-hidden="true" />}
                {c.success ? "OK" : "Falha"}
              </span>
              <span className="ml-auto text-pg-ui text-pg-ardosia tabular-nums whitespace-nowrap">
                {heartbeat || !c.success ? <span aria-label="sem dados">—</span> : `${c.latency_ms} ms`}
              </span>
            </div>
            {!c.success && (
              <div className="font-pg-mono text-pg-caption text-pg-neblina truncate" title={detailOf(c, type)}>
                {detailOf(c, type)}
              </div>
            )}
          </li>
        ))}
      </ul>
      {(mobile.length < checks.length || checks.length < total) && (
        <button
          type="button"
          onClick={onMore}
          disabled={loadingMore}
          className="md:hidden w-full min-h-[48px] border-t border-pg-borda text-pg-ui font-medium text-pg-ciano-fundo disabled:opacity-50 focus-visible:outline-2 focus-visible:outline focus-visible:-outline-offset-2 focus-visible:outline-pg-ciano"
        >
          {loadingMore ? "Carregando…" : "Ver mais checks"}
        </button>
      )}

      <table className="hidden md:table w-full table-fixed border-collapse text-left">
        <caption className="sr-only">Últimos checks, do mais recente ao mais antigo</caption>
        <colgroup>
          <col className="w-[112px] sm:w-[148px]" />
          <col className="sm:w-[150px]" />
          <col className="w-[84px]" />
          <col className="hidden sm:table-column" />
        </colgroup>
        <thead>
          <tr className="bg-pg-ciano-nevoa text-pg-caption text-pg-ciano-fundo">
            <th scope="col" className="px-4 py-3 font-medium">Horário</th>
            <th scope="col" className="px-4 py-3 font-medium">Status</th>
            <th scope="col" className="px-4 py-3 font-medium text-right">Resposta</th>
            <th scope="col" className="px-4 py-3 font-medium hidden sm:table-cell">Detalhe</th>
          </tr>
        </thead>
        <tbody>
          {checks.map((c) => {
            const detail = detailOf(c, type);
            return (
              <tr key={c.id} className="border-t border-pg-borda">
                <td className={`${cell} text-pg-ui text-pg-ardosia tabular-nums`}>{checkTime(c.checked_at, now)}</td>
                <td className={cell}>
                  <span className={`inline-flex items-center gap-1.5 text-pg-caption font-medium ${c.success ? "text-pg-up" : "text-pg-down"}`}>
                    {c.success ? <Check size={14} aria-hidden="true" /> : <X size={14} aria-hidden="true" />}
                    {c.success ? "Sucesso" : "Falha"}
                  </span>
                  {/* mobile: sem a coluna Detalhe, o erro aparece sob o status */}
                  {!c.success && (
                    <div className="sm:hidden mt-0.5 font-pg-mono text-pg-caption text-pg-neblina truncate" title={detail}>
                      {detail}
                    </div>
                  )}
                </td>
                <td className={`${cell} text-right text-pg-ui text-pg-ardosia tabular-nums whitespace-nowrap`}>
                  {heartbeat || !c.success ? <span aria-label="sem dados">—</span> : `${c.latency_ms} ms`}
                </td>
                <td className={`${cell} hidden sm:table-cell font-pg-mono text-pg-caption text-pg-neblina`}>
                  <div className="truncate" title={detail}>
                    {detail}
                  </div>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
      <div className="hidden md:block px-4 py-3 border-t border-pg-borda text-pg-caption text-pg-neblina" aria-live="polite">
        Mostrando {checks.length === 1 ? "o mais recente" : `os ${checks.length} mais recentes`}
        {checks.length < total && (
          <>
            {" · "}
            <button
              type="button"
              onClick={onMore}
              disabled={loadingMore}
              className="p-0 border-0 bg-transparent text-pg-caption text-pg-ciano-fundo underline cursor-pointer disabled:opacity-50 focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
            >
              {loadingMore ? "carregando…" : "ver mais"}
            </button>
          </>
        )}
      </div>
    </div>
  );
}
