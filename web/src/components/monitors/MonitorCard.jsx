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

  return (
    <div
      style={{
        background: t.colors.surface,
        border: `1px solid ${t.colors.border}`,
        borderRadius: t.radius.md,
        padding: "18px 20px",
        display: "flex",
        alignItems: "center",
        gap: "16px",
      }}
    >
      {/* state dot */}
      <div
        style={{
          width: "10px",
          height: "10px",
          borderRadius: "50%",
          flexShrink: 0,
          background:
            current_state === "UP"
              ? t.colors.success
              : current_state === "DOWN"
                ? t.colors.danger
                : t.colors.unknown,
          boxShadow:
            current_state === "UP"
              ? `0 0 6px ${t.colors.success}`
              : current_state === "DOWN"
                ? `0 0 6px ${t.colors.danger}`
                : "none",
        }}
      />

      {/* info */}
      <div style={{ flex: 1, minWidth: 0 }}>
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "8px",
            marginBottom: "2px",
          }}
        >
          <span style={{ fontWeight: 600, fontSize: "14px" }}>{name}</span>
          {!enabled && (
            <span
              style={{
                fontSize: "10px",
                color: t.colors.textMuted,
                background: t.colors.surfaceAlt,
                padding: "1px 6px",
                borderRadius: "4px",
              }}
            >
              PAUSED
            </span>
          )}
        </div>
        <div
          style={{
            color: t.colors.textMuted,
            fontSize: "12px",
            overflow: "hidden",
            textOverflow: "ellipsis",
            whiteSpace: "nowrap",
          }}
        >
          {url}
        </div>
        <div
          style={{
            color: t.colors.textMuted,
            fontSize: "11px",
            marginTop: "4px",
          }}
        >
          every {interval_seconds}s
          {last_checked_at &&
            ` · last check ${new Date(last_checked_at).toLocaleTimeString()}`}
        </div>
      </div>

      <Badge state={current_state} />

      <div style={{ display: "flex", gap: "8px" }}>
        <Button
          variant="ghost"
          onClick={() => onEdit(monitor)}
          style={{ padding: "6px 12px", fontSize: "12px" }}
        >
          Edit
        </Button>
        <Button
          variant="danger"
          onClick={() => onDelete(monitor)}
          style={{ padding: "6px 12px", fontSize: "12px" }}
        >
          Delete
        </Button>
      </div>
    </div>
  );
}
