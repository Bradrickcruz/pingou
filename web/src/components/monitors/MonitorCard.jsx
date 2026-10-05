import { Badge } from "../ui/Badge";
import { Button } from "../ui/Button";

const dotColor = {
  UP: "bg-pg-up",
  DOWN: "bg-pg-down",
};

export function MonitorCard({ monitor, onEdit, onDelete }) {
  const {
    name,
    url,
    current_state,
    last_checked_at,
    interval_seconds,
    enabled,
  } = monitor;

  return (
    <div className="flex flex-wrap items-center gap-x-4 gap-y-3 p-4 rounded-pg-card border border-pg-borda bg-pg-branco">
      <div className="flex items-center gap-4 flex-1 min-w-[12rem]">
        <div
          className={`w-2.5 h-2.5 rounded-full flex-shrink-0 ${dotColor[current_state] ?? "bg-pg-unknown"}`}
        />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-0.5">
            <span className="font-medium text-pg-ui text-pg-ardosia truncate">{name}</span>
            {!enabled && (
              <span className="text-pg-caption px-1.5 py-0.5 rounded-full bg-pg-unknown-bg text-pg-neblina flex-shrink-0">
                PAUSADO
              </span>
            )}
          </div>
          <div className="text-pg-caption text-pg-neblina truncate">{url}</div>
          <div className="text-pg-caption text-pg-neblina mt-1">
            a cada {interval_seconds}s
            {last_checked_at &&
              ` · última verificação ${new Date(last_checked_at).toLocaleTimeString("pt-BR")}`}
          </div>
        </div>
      </div>

      <div className="flex items-center gap-2 ml-auto sm:ml-0">
        <Badge state={current_state} />
      </div>

      <div className="flex gap-2 w-full sm:w-auto justify-end">
        <Button variant="ghost" onClick={() => onEdit(monitor)} className="py-1.5 px-3 text-pg-caption">
          Editar
        </Button>
        <Button variant="destructive" onClick={() => onDelete(monitor)} className="py-1.5 px-3 text-pg-caption">
          Excluir
        </Button>
      </div>
    </div>
  );
}
