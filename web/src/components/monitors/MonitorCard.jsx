import { Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { Badge } from "../ui/Badge";
import { HistoryBar } from "./HistoryBar";
import {
  cardFooter,
  heartbeatExpectation,
  statusOf,
  targetOf,
} from "../../lib/monitors";

const stateLabel = {
  UP: "no ar",
  DOWN: "fora do ar",
  UNKNOWN: "desconhecido",
  PAUSED: "pausado",
};

// Card do painel no mobile (< md). O card todo leva ao detalhe, onde ficam
// editar, pausar e excluir; a tabela do desktop é o MonitorRow.
export function MonitorCard({ monitor, now }) {
  const status = statusOf(monitor);
  const down = status === "DOWN";
  const paused = status === "PAUSED";
  const isHeartbeat = monitor.type === "heartbeat";

  return (
    <li
      className={`relative p-4 rounded-pg-card border bg-pg-branco has-[:focus-visible]:outline-2 has-[:focus-visible]:outline has-[:focus-visible]:outline-pg-ciano ${
        down ? "border-[#EDC4C8]" : "border-pg-borda"
      }`}
    >
      <div className="flex items-center gap-2 min-w-0">
        <Badge state={status} />
        <Link
          to={`/monitors/${monitor.id}`}
          aria-label={`${monitor.name}, ${stateLabel[status]}`}
          className={`flex-1 min-w-0 truncate text-pg-body font-medium no-underline outline-none after:absolute after:inset-0 after:content-[''] ${
            paused ? "text-pg-neblina" : "text-pg-ardosia"
          }`}
        >
          {monitor.name}
        </Link>
        <ChevronRight
          size={18}
          aria-hidden="true"
          className="shrink-0 text-pg-neblina"
        />
      </div>

      <p
        className={`mt-1 truncate text-pg-caption text-pg-neblina ${isHeartbeat ? "" : "font-pg-mono"}`}
      >
        {isHeartbeat
          ? `Heartbeat · ${heartbeatExpectation(monitor.interval_seconds)}`
          : targetOf(monitor)}
      </p>

      {down && monitor.last_error && (
        <p className="mt-3 p-3 rounded-pg-control bg-pg-nuvem font-pg-mono text-pg-caption text-pg-ardosia break-words">
          {monitor.last_error}
        </p>
      )}

      <div className="mt-3">
        <HistoryBar history={monitor.history} dimmed={paused} fluid />
      </div>

      <p className="mt-2 text-pg-caption text-pg-neblina">
        {cardFooter(monitor, now)}
      </p>
    </li>
  );
}
