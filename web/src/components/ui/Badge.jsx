import { Check, X, CircleHelp } from "lucide-react";

const stateStyles = {
  UP: { className: "bg-pg-up-bg text-pg-up", icon: Check },
  DOWN: { className: "bg-pg-down-bg text-pg-down", icon: X },
  UNKNOWN: { className: "bg-pg-unknown-bg text-pg-unknown", icon: CircleHelp },
};

export function Badge({ state }) {
  const s = stateStyles[state] ?? stateStyles.UNKNOWN;
  const Icon = s.icon;
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[12px] font-medium uppercase tracking-wider ${s.className}`}
    >
      <Icon size={12} aria-hidden="true" />
      {state}
    </span>
  );
}
