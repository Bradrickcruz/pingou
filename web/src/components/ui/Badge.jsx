import { Check, X, CircleHelp } from "lucide-react";
import { tokens as t } from "../../theme/tokens";

const stateColors = {
  UP: { bg: "#14532d", color: t.colors.success, icon: Check },
  DOWN: { bg: "#450a0a", color: t.colors.danger, icon: X },
  UNKNOWN: { bg: "#1f2937", color: t.colors.unknown, icon: CircleHelp },
};

export function Badge({ state }) {
  const c = stateColors[state] ?? stateColors.UNKNOWN;
  const Icon = c.icon;
  return (
    <span
      className="inline-flex items-center gap-1 px-2.5 py-0.5 text-[11px] font-semibold uppercase tracking-wider rounded"
      style={{
        background: c.bg,
        color: c.color,
        borderRadius: t.radius.sm,
      }}
    >
      <Icon size={12} aria-hidden="true" />
      {state}
    </span>
  );
}