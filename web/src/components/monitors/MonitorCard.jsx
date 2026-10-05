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
    <div className="flex items-center gap-4 p-4 rounded-pg-card border border-pg-borda bg-pg-branco">
      <div
        className={`w-2.5 h-2.5 rounded-full flex-shrink-0 ${dotColor[current_state] ?? "bg-pg-unknown"}`}
      />

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-0.5">
          <span className="font-medium text-pg-ui text-pg-ardosia">{name}</span>
          {!enabled && (
            <span className="text-pg-caption px-1.5 py-0.5 rounded-full bg-pg-unknown-bg text-pg-neblina">
              PAUSED
            </span>
          )}
        </div>
        <div className="text-pg-caption text-pg-neblina overflow-hidden text-ellipsis whitespace-nowrap">
          {url}
        </div>
        <div className="text-pg-caption text-pg-neblina mt-1">
          every {interval_seconds}s
          {last_checked_at &&
            ` · last check ${new Date(last_checked_at).toLocaleTimeString()}`}
        </div>
      </div>

      <Badge state={current_state} />

      <div className="flex gap-2">
        <Button
          variant="ghost"
          onClick={() => onEdit(monitor)}
          className="py-1.5 px-3 text-pg-caption"
        >
          Edit
        </Button>
        <Button
          variant="destructive"
          onClick={() => onDelete(monitor)}
          className="py-1.5 px-3 text-pg-caption"
        >
          Delete
        </Button>
      </div>
    </div>
  );
}
