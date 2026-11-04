import { useState } from "react";
import { Wrench } from "lucide-react";
import { toast } from "react-toastify";
import { maintenanceApi } from "../../api/maintenance";
import { Button } from "../ui/Button";
import { Modal } from "../ui/Modal";
import { fullDate } from "../../lib/format";
import { zoneOffsetLabel } from "../../lib/timezone";

// Janela viva do monitor (agendada ou em curso): motivo, início, fim e cancelar.
// O estado vem do servidor (maintenance.active), não do relógio do navegador.
export function MaintenanceCard({ maintenance, onChanged }) {
  const [confirming, setConfirming] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const running = maintenance.active;
  const zone = zoneOffsetLabel(new Date(maintenance.starts_at));

  const cancel = async () => {
    setSaving(true);
    setError(null);
    try {
      await maintenanceApi.cancel(maintenance.id);
      toast.success(running ? "Manutenção encerrada" : "Manutenção cancelada");
      setConfirming(false);
      await onChanged?.();
    } catch (e) {
      setError(e.message);
    } finally {
      setSaving(false);
    }
  };

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
              <span className="text-pg-neblina">({zone})</span>
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
        <Button
          variant="ghost"
          onClick={() => setConfirming(true)}
          className="flex-shrink-0 bg-pg-branco"
        >
          {running ? "Encerrar agora" : "Cancelar"}
        </Button>
      </div>
      {error && !confirming && (
        <p role="alert" className="mt-3 mb-0 text-pg-caption text-pg-down">
          {error}
        </p>
      )}

      {confirming && (
        <Modal
          title={running ? "Encerrar manutenção" : "Cancelar manutenção"}
          onClose={() => setConfirming(false)}
        >
          <p className="mb-5 text-pg-neblina">
            {running
              ? "O monitor volta a ser conferido agora, com o primeiro check na hora."
              : "A janela agendada é apagada. Esta ação não pode ser desfeita."}
          </p>
          {error && (
            <p role="alert" className="-mt-2 mb-4 text-pg-caption text-pg-down">
              {error}
            </p>
          )}
          <div className="flex gap-2.5 justify-end">
            <Button
              variant="ghost"
              onClick={() => setConfirming(false)}
              data-autofocus
            >
              Voltar
            </Button>
            <Button variant="danger" onClick={cancel} disabled={saving}>
              {saving
                ? "Aguarde..."
                : running
                  ? "Encerrar agora"
                  : "Cancelar manutenção"}
            </Button>
          </div>
        </Modal>
      )}
    </section>
  );
}
