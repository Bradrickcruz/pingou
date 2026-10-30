import { Link } from "react-router-dom";
import { Globe, Server, Activity } from "lucide-react";
import { Badge } from "../ui/Badge";
import { HistoryBar } from "./HistoryBar";
import { RowMenu } from "./RowMenu";
import { formatUptime, timeAgo } from "../../lib/format";
import { statusOf, targetOf } from "../../lib/monitors";

const typeIcon = { http: Globe, tcp: Server, heartbeat: Activity };

const cell = "px-5 py-3 align-middle";

// Linha da tabela do painel. Substitui o antigo MonitorCard.
export function MonitorRow({ monitor, now, onEdit, onToggle, onDelete }) {
  const status = statusOf(monitor);
  const paused = status === "PAUSED";
  const Icon = typeIcon[monitor.type] ?? Globe;
  const isHeartbeat = monitor.type === "heartbeat";
  const muted = paused ? "text-pg-neblina" : "text-pg-ardosia";

  const uptime = formatUptime(monitor.uptime_24h);
  // latência só faz sentido com o alvo respondendo; heartbeat não mede resposta
  const latency =
    !isHeartbeat && status === "UP" && monitor.last_latency_ms != null
      ? `${monitor.last_latency_ms} ms`
      : null;
  const lastCheck = paused ? "pausado" : timeAgo(monitor.last_checked_at, now);

  return (
    <tr className="border-t border-pg-borda">
      <td className={cell}>
        <Badge state={status} />
      </td>
      <td className={`${cell} max-w-0 w-full`}>
        <Link
          to={`/monitors/${monitor.id}`}
          className={`block text-pg-ui font-medium truncate no-underline hover:underline focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${muted}`}
        >
          {monitor.name}
        </Link>
        <div className="flex items-center gap-1.5 mt-0.5 text-pg-neblina min-w-0">
          <Icon size={14} aria-hidden="true" className="flex-shrink-0" />
          <span className="sr-only">{monitor.type.toUpperCase()}:</span>
          <span
            className={`truncate text-pg-caption ${isHeartbeat ? "" : "font-pg-mono"}`}
            title={targetOf(monitor)}
          >
            {targetOf(monitor)}
          </span>
        </div>
      </td>
      <td className={`${cell} hidden lg:table-cell`}>
        <HistoryBar history={monitor.history} dimmed={paused} />
      </td>
      <td
        className={`${cell} hidden md:table-cell text-right tabular-nums text-pg-ui ${
          status === "DOWN" ? "text-pg-down" : muted
        }`}
      >
        {uptime ?? <span aria-label="sem dados">—</span>}
      </td>
      <td
        className={`${cell} hidden md:table-cell text-right tabular-nums text-pg-ui whitespace-nowrap ${muted}`}
      >
        {latency ?? <span aria-label="sem dados">—</span>}
      </td>
      <td
        className={`${cell} hidden sm:table-cell text-pg-caption text-pg-neblina whitespace-nowrap`}
      >
        {lastCheck ?? "—"}
      </td>
      <td className={`${cell} pl-0 text-right`}>
        <RowMenu
          monitor={monitor}
          onEdit={onEdit}
          onToggle={onToggle}
          onDelete={onDelete}
        />
      </td>
    </tr>
  );
}
