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
      className="inline-block px-2.5 py-0.5 text-[11px] font-semibold uppercase tracking-wider rounded"
      style={{
        background: c.bg,
        color: c.color,
        borderRadius: t.radius.sm,
      }}
    >
      {state}
    </span>
  );
}