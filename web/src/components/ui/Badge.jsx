import { Check, X, CircleHelp } from "lucide-react";

const stateStyles = {
  UP: { bg: "bg-pg-up-bg", iconColor: "text-pg-up", icon: Check },
  DOWN: { bg: "bg-pg-down-bg", iconColor: "text-pg-down", icon: X },
  UNKNOWN: { bg: "bg-pg-unknown-bg", iconColor: "text-pg-unknown", icon: CircleHelp },
};

export function Badge({ state }) {
  const s = stateStyles[state] ?? stateStyles.UNKNOWN;
  const Icon = s.icon;
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-pg-caption font-medium uppercase tracking-wider text-pg-ardosia ${s.bg}`}
    >
      <Icon size={12} aria-hidden="true" className={s.iconColor} />
      {state}
    </span>
  );
}
