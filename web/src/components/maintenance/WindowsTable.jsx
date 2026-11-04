import { Link } from "react-router-dom";
import { Pencil, Power, X } from "lucide-react";
import { whenLabel } from "../../lib/format";
import { windowState } from "../../lib/maintenance";

const stateBadge = {
  scheduled: {
    label: "Agendada",
    cls: "bg-pg-branco border border-pg-borda text-pg-neblina",
  },
  active: { label: "Em curso", cls: "bg-pg-ciano-nevoa text-pg-ciano-fundo" },
};

function StateBadge({ state }) {
  const s = stateBadge[state];
  return (
    <span
      className={`inline-flex rounded-full px-2.5 py-0.5 text-pg-caption font-medium whitespace-nowrap ${s.cls}`}
    >
      {s.label}
    </span>
  );
}

const actionCls =
  "inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-pg-control text-pg-ui font-medium focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano";

// Editar e Cancelar/Desligar. aria-label leva o nome do monitor.
function Actions({ w, state, onEdit, onCancel, full = false }) {
  const active = state === "active";
  const CancelIcon = active ? Power : X;
  return (
    <div className={`flex gap-2 ${full ? "" : "justify-end"}`}>
      <button
        type="button"
        onClick={() => onEdit(w)}
        aria-label={`Editar manutenção de ${w.monitor_name}`}
        className={`${actionCls} text-pg-ciano-fundo hover:bg-pg-ciano-nevoa ${full ? "flex-1 justify-center min-h-[44px] border border-pg-borda" : ""}`}
      >
        <Pencil size={16} aria-hidden="true" />
        Editar
      </button>
      <button
        type="button"
        onClick={() => onCancel(w)}
        aria-label={`${active ? "Desligar" : "Cancelar"} manutenção de ${w.monitor_name}`}
        className={`${actionCls} text-pg-down hover:bg-pg-down-bg ${full ? "flex-1 justify-center min-h-[44px] border border-pg-borda" : ""}`}
      >
        <CancelIcon size={16} aria-hidden="true" />
        {active ? "Desligar" : "Cancelar"}
      </button>
    </div>
  );
}

const cell = "px-5 py-3 align-middle";

// Janelas vivas: tabela no desktop, cartões no celular. O estado é recalculado
// com o relógio local (now), para a linha virar "Em curso" sem recarregar.
export function WindowsTable({ windows, now, onEdit, onCancel }) {
  return (
    <div className="rounded-pg-card border border-pg-borda bg-pg-branco overflow-hidden">
      <table className="hidden md:table w-full table-fixed border-collapse text-left">
        <caption className="sr-only">
          Manutenções agendadas e em curso, ordenadas pelo início
        </caption>
        <colgroup>
          <col className="w-[24%]" />
          <col className="w-[120px]" />
          <col className="w-[150px]" />
          <col className="w-[150px]" />
          <col />
          <col className="w-[230px]" />
        </colgroup>
        <thead>
          <tr className="bg-pg-ciano-nevoa text-pg-caption text-pg-ciano-fundo">
            <th scope="col" className="px-5 py-3 font-medium">
              Monitor
            </th>
            <th scope="col" className="px-5 py-3 font-medium">
              Estado
            </th>
            <th scope="col" className="px-5 py-3 font-medium">
              Início
            </th>
            <th scope="col" className="px-5 py-3 font-medium">
              Fim
            </th>
            <th scope="col" className="px-5 py-3 font-medium">
              Motivo
            </th>
            <th scope="col" className="px-5 py-3">
              <span className="sr-only">Ações</span>
            </th>
          </tr>
        </thead>
        <tbody>
          {windows.map((w) => {
            const state = windowState(w, now);
            return (
              <tr key={w.id} className="border-t border-pg-borda">
                <td className={cell}>
                  <Link
                    to={`/monitors/${w.monitor_id}`}
                    className="block text-pg-ui font-medium text-pg-ardosia truncate no-underline hover:underline"
                  >
                    {w.monitor_name}
                  </Link>
                </td>
                <td className={cell}>
                  <StateBadge state={state} />
                </td>
                <td
                  className={`${cell} text-pg-ui text-pg-ardosia whitespace-nowrap tabular-nums`}
                >
                  {whenLabel(w.starts_at, now)}
                </td>
                <td
                  className={`${cell} text-pg-ui text-pg-ardosia whitespace-nowrap tabular-nums`}
                >
                  {whenLabel(w.ends_at, now)}
                </td>
                <td className={`${cell} text-pg-ui text-pg-neblina`}>
                  <div className="truncate" title={w.reason ?? undefined}>
                    {w.reason ?? "—"}
                  </div>
                </td>
                <td className={`${cell} pl-0`}>
                  <Actions
                    w={w}
                    state={state}
                    onEdit={onEdit}
                    onCancel={onCancel}
                  />
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>

      <ul className="md:hidden list-none m-0 p-0">
        {windows.map((w) => {
          const state = windowState(w, now);
          return (
            <li
              key={w.id}
              className="px-4 py-4 border-b border-pg-borda last:border-b-0"
            >
              <div className="flex items-center justify-between gap-3">
                <Link
                  to={`/monitors/${w.monitor_id}`}
                  className="text-pg-body font-medium text-pg-ardosia truncate no-underline"
                >
                  {w.monitor_name}
                </Link>
                <StateBadge state={state} />
              </div>
              <p className="mt-1 mb-0 text-pg-caption text-pg-ardosia tabular-nums">
                {whenLabel(w.starts_at, now)} → {whenLabel(w.ends_at, now)}
              </p>
              {w.reason && (
                <p className="mt-1 mb-0 text-pg-caption text-pg-neblina break-words">
                  {w.reason}
                </p>
              )}
              <div className="mt-3">
                <Actions
                  w={w}
                  state={state}
                  onEdit={onEdit}
                  onCancel={onCancel}
                  full
                />
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
