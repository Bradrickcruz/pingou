import { tokens as t } from "../../theme/tokens";

const stateColors = {
  UP: { bg: "#14532d", color: t.colors.success },
  DOWN: { bg: "#450a0a", color: t.colors.danger },
  UNKNOWN: { bg: "#1f2937", color: t.colors.unknown },
};

export function Badge({ state }) {
  const c = stateColors[state] ?? stateColors.UNKNOWN;
  return (
    <span
      style={{
        background: c.bg,
        color: c.color,
        borderRadius: t.radius.sm,
        padding: "2px 10px",
        fontSize: "11px",
        fontWeight: 600,
        letterSpacing: "0.05em",
        textTransform: "uppercase",
      }}
    >
      {state}
    </span>
  );
}
