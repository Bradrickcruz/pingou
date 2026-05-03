import { Badge } from "../ui/Badge";
import { Button } from "../ui/Button";
import { tokens as t } from "../../theme/tokens";

export function MonitorCard({ monitor, onEdit, onDelete }) {
  const {
    name,
    url,
    current_state,
    last_checked_at,
    interval_seconds,
    enabled,
  } = monitor;

  const stateColor =
    current_state === "UP"
      ? t.colors.success
      : current_state === "DOWN"
        ? t.colors.danger
        : t.colors.unknown;

  const stateGlow =
    current_state === "UP"
      ? `0 0 6px ${t.colors.success}`
      : current_state === "DOWN"
        ? `0 0 6px ${t.colors.danger}`
        : "none";

  return (
    <div
      className="flex items-center gap-4 p-4 rounded-md border"
      style={{
        background: t.colors.surface,
        borderColor: t.colors.border,
        borderRadius: t.radius.md,
      }}
    >
      {/* state dot */}
      <div
        className="w-2.5 h-2.5 rounded-full flex-shrink-0"
        style={{
          background: stateColor,
          boxShadow: stateGlow,
        }}
      />

      {/* info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-0.5">
          <span className="font-semibold text-sm">{name}</span>
          {!enabled && (
            <span
              className="text-[10px] px-1.5 py-0.5 rounded"
              style={{
                color: t.colors.textMuted,
                background: t.colors.surfaceAlt,
              }}
            >
              PAUSED
            </span>
          )}
        </div>
        <div
          className="text-xs overflow-hidden text-ellipsis whitespace-nowrap"
          style={{
            color: t.colors.textMuted,
          }}
        >
          {url}
        </div>
        <div
          className="text-[11px] mt-1"
          style={{
            color: t.colors.textMuted,
          }}
        >
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
          className="py-1.5 px-3 text-xs"
        >
          Edit
        </Button>
        <Button
          variant="danger"
          onClick={() => onDelete(monitor)}
          className="py-1.5 px-3 text-xs"
        >
          Delete
        </Button>
      </div>
    </div>
  );
}