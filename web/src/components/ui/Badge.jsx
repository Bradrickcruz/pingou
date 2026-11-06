import { Check, X, CircleHelp, Pause, Wrench } from "lucide-react";

const stateStyles = {
  UP: {
    bg: "bg-pg-status-up-bg",
    iconColor: "text-pg-status-up",
    icon: Check,
    label: "UP",
  },
  DOWN: {
    bg: "bg-pg-status-down-bg",
    iconColor: "text-pg-status-down",
    icon: X,
    label: "DOWN",
  },
  UNKNOWN: {
    bg: "bg-pg-status-unknown-bg",
    iconColor: "text-pg-status-unknown",
    icon: CircleHelp,
    label: "UNKNOWN",
  },
  // neutro e distinto do pausado: o monitor não está sendo conferido de propósito
  MAINTENANCE: {
    bg: "bg-pg-destaque",
    iconColor: "text-pg-link",
    icon: Wrench,
    label: "MANUTENÇÃO",
    text: "text-pg-link",
  },
  PAUSED: {
    bg: "bg-pg-superficie border border-pg-linha",
    iconColor: "text-pg-texto-2",
    icon: Pause,
    label: "PAUSADO",
    text: "text-pg-texto-2",
  },
};

export function Badge({ state }) {
  const s = stateStyles[state] ?? stateStyles.UNKNOWN;
  const Icon = s.icon;
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-pg-caption font-medium uppercase tracking-wider whitespace-nowrap ${s.text ?? "text-pg-texto"} ${s.bg}`}
    >
      <Icon size={12} aria-hidden="true" className={s.iconColor} />
      {s.label}
    </span>
  );
}
