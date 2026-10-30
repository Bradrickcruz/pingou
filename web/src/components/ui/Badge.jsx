import { Check, X, CircleHelp, Pause } from "lucide-react";

const stateStyles = {
  UP: { bg: "bg-pg-up-bg", iconColor: "text-pg-up", icon: Check, label: "UP" },
  DOWN: {
    bg: "bg-pg-down-bg",
    iconColor: "text-pg-down",
    icon: X,
    label: "DOWN",
  },
  UNKNOWN: {
    bg: "bg-pg-unknown-bg",
    iconColor: "text-pg-unknown",
    icon: CircleHelp,
    label: "UNKNOWN",
  },
  PAUSED: {
    bg: "bg-pg-branco border border-pg-borda",
    iconColor: "text-pg-neblina",
    icon: Pause,
    label: "PAUSADO",
    text: "text-pg-neblina",
  },
};

export function Badge({ state }) {
  const s = stateStyles[state] ?? stateStyles.UNKNOWN;
  const Icon = s.icon;
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-pg-caption font-medium uppercase tracking-wider whitespace-nowrap ${s.text ?? "text-pg-ardosia"} ${s.bg}`}
    >
      <Icon size={12} aria-hidden="true" className={s.iconColor} />
      {s.label}
    </span>
  );
}
