import { Bell, Clock, ShieldAlert, ShieldCheck, ShieldX } from "lucide-react";
import { certStatus } from "../../lib/cert";

// Nível → pílula (ícone + texto, nunca só cor) e cor do texto da data.
// DESIGN.md não tem token de "atenção": o nível de 30 dias usa o neutro
// (unknown) com texto "Atenção"; vermelho fica para 7 dias ou vencido.
const levels = {
  ok: {
    icon: ShieldCheck,
    iconColor: "text-pg-status-up",
    bg: "bg-pg-status-up-bg",
    label: "Válido",
  },
  warn: {
    icon: ShieldAlert,
    iconColor: "text-pg-texto",
    bg: "bg-pg-status-unknown-bg",
    label: "Atenção",
  },
  alert: {
    icon: ShieldAlert,
    iconColor: "text-pg-status-down",
    bg: "bg-pg-status-down-bg",
    label: "Renove já",
    text: "text-pg-status-down font-medium",
  },
  expired: {
    icon: ShieldX,
    iconColor: "text-pg-status-down",
    bg: "bg-pg-status-down-bg",
    label: "Expirado",
    text: "text-pg-status-down font-medium",
  },
};

// Certificado SSL do monitor HTTP com URL https: validade, dias restantes e se
// o aviso está ligado. Fora disso, não aparece.
export function CertCard({ monitor, now }) {
  const status = certStatus(monitor, now);
  if (!status) return null;
  const level = levels[status.level];
  // site fora do ar não devolve certificado: a data é a do último check que viu
  const stale = monitor.current_state === "DOWN" && status.level !== "pending";

  return (
    <div className="p-4 rounded-pg-card border border-pg-linha bg-pg-superficie text-pg-ui">
      <div className="flex flex-wrap items-center justify-between gap-2">
        {level ? (
          <span
            className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-pg-caption font-medium uppercase tracking-wider whitespace-nowrap text-pg-texto ${level.bg}`}
          >
            <level.icon
              size={12}
              aria-hidden="true"
              className={level.iconColor}
            />
            {level.label}
          </span>
        ) : (
          <span className="inline-flex items-center gap-1.5 text-pg-texto-2">
            <Clock size={14} aria-hidden="true" />
            {status.text}
          </span>
        )}
        {stale && (
          <span className="text-pg-caption text-pg-texto-2">
            último valor conhecido
          </span>
        )}
      </div>
      {level && (
        <p className={`mt-2 mb-0 ${level.text ?? "text-pg-texto"}`}>
          {status.text}
        </p>
      )}
      {monitor.cert_alert_enabled && (
        <p className="inline-flex items-center gap-1.5 mt-2 mb-0 text-pg-caption text-pg-texto-2">
          <Bell size={12} aria-hidden="true" />
          Aviso ativo (30 e 7 dias antes)
        </p>
      )}
    </div>
  );
}
