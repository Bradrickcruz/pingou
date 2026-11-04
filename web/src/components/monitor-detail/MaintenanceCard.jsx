import { useState } from "react";
import { Pencil, Wrench } from "lucide-react";
import { Button } from "../ui/Button";
import { MaintenanceModal } from "../monitors/MaintenanceModal";
import { CancelWindowDialog } from "../maintenance/CancelWindowDialog";
import { fullDate } from "../../lib/format";
import { getTimeZone } from "../../lib/timezone";

// Janela viva do monitor (agendada ou em curso): motivo, início, fim, editar e
// cancelar/desligar (D17). O estado vem do servidor (maintenance.active).
export function MaintenanceCard({ maintenance, monitorName, onChanged }) {
  const [open, setOpen] = useState(null); // "edit" | "cancel"
  const running = maintenance.active;
  const w = { ...maintenance, monitor_name: monitorName };
  const close = () => setOpen(null);

  return (
    <section
      aria-labelledby="manutencao"
      className="mb-6 md:mb-8 p-4 sm:p-5 rounded-pg-card border border-pg-borda bg-pg-ciano-nevoa"
    >
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
        <div className="flex items-start gap-3 min-w-0">
          <span className="flex-shrink-0 inline-flex items-center justify-center w-9 h-9 rounded-full bg-pg-branco text-pg-ciano-fundo">
            <Wrench size={18} aria-hidden="true" />
          </span>
          <div className="min-w-0">
            <h2
              id="manutencao"
              className="text-pg-ui font-medium text-pg-ardosia my-0"
            >
              {running ? "Em manutenção" : "Manutenção agendada"}
            </h2>
            <p className="mt-0.5 mb-0 text-pg-caption text-pg-ardosia tabular-nums">
              {fullDate(maintenance.starts_at)} →{" "}
              {fullDate(maintenance.ends_at)}{" "}
              <span className="text-pg-neblina">({getTimeZone()})</span>
            </p>
            {maintenance.reason && (
              <p className="mt-1 mb-0 text-pg-ui text-pg-ardosia break-words">
                <span className="font-medium">Motivo:</span>{" "}
                {maintenance.reason}
              </p>
            )}
            <p className="mt-1 mb-0 text-pg-caption text-pg-neblina">
              {running
                ? "Nenhum check roda e ninguém é avisado até o fim."
                : "Quando começar, os checks param e um incidente aberto é fechado sem aviso."}
            </p>
          </div>
        </div>
        <div className="flex flex-shrink-0 gap-2">
          <Button
            variant="ghost"
            onClick={() => setOpen("edit")}
            aria-label="Editar manutenção"
            className="flex-1 sm:flex-none gap-1.5 bg-pg-branco"
          >
            <Pencil size={16} aria-hidden="true" />
            Editar
          </Button>
          <Button
            variant="ghost"
            onClick={() => setOpen("cancel")}
            className="flex-1 sm:flex-none bg-pg-branco"
          >
            {running ? "Desligar" : "Cancelar"}
          </Button>
        </div>
      </div>

      {open === "edit" && (
        <MaintenanceModal
          window={w}
          onClose={close}
          onSaved={onChanged}
          onStale={onChanged}
        />
      )}
      {open === "cancel" && (
        <CancelWindowDialog
          window={w}
          active={running}
          onClose={close}
          onDone={onChanged}
        />
      )}
    </section>
  );
}
