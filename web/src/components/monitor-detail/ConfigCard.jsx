import { formatDuration, fullDate, plural } from "../../lib/format";

const typeLabel = { http: "HTTP", tcp: "TCP", heartbeat: "Heartbeat" };

export function ConfigCard({ monitor: m }) {
  const heartbeat = m.type === "heartbeat";
  const rows = [
    ["Tipo", typeLabel[m.type] ?? m.type],
    ["Verifica a cada", formatDuration(m.interval_seconds)],
    heartbeat
      ? ["Tolerância", formatDuration(m.config?.grace_seconds ?? 0)]
      : ["Tempo limite", formatDuration(m.timeout_seconds)],
    !heartbeat && ["Cai depois de", plural(m.failure_threshold, "falha", "falhas seguidas")],
    ["Criado em", fullDate(m.created_at)],
  ].filter(Boolean);

  return (
    <dl className="m-0 px-4 rounded-pg-card border border-pg-borda bg-pg-branco">
      {rows.map(([k, v]) => (
        <div key={k} className="flex justify-between gap-4 py-2.5 border-t border-pg-borda first:border-t-0 text-pg-ui">
          <dt className="text-pg-neblina">{k}</dt>
          <dd className="m-0 text-right text-pg-ardosia">{v}</dd>
        </div>
      ))}
    </dl>
  );
}
