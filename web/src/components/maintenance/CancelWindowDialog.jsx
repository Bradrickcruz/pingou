import { useState } from "react";
import { toast } from "react-toastify";
import { maintenanceApi } from "../../api/maintenance";
import { Button } from "../ui/Button";
import { Modal } from "../ui/Modal";

// Confirmação de cancelar (agendada: a janela some) ou desligar (em curso: o
// monitor volta a ser checado agora). Erro sai em toast; o diálogo fecha.
export function CancelWindowDialog({ window: w, active, onClose, onDone }) {
  const [saving, setSaving] = useState(false);
  const who = w.monitor_name ? (
    <>
      {" "}
      de{" "}
      <strong className="font-medium text-pg-ardosia">{w.monitor_name}</strong>
    </>
  ) : null;

  const confirm = async () => {
    setSaving(true);
    try {
      await maintenanceApi.cancel(w.id);
      toast.success(active ? "Manutenção desligada" : "Manutenção cancelada");
      onClose();
      await onDone?.();
    } catch (e) {
      toast.error(e.message);
      onClose();
      await onDone?.(); // a lista pode ter mudado (ex.: a janela já acabou)
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      title={active ? "Desligar manutenção" : "Cancelar manutenção"}
      onClose={onClose}
    >
      <p className="mb-5 text-pg-neblina">
        {active ? (
          <>
            A manutenção{who} termina agora. O monitor volta a ser checado na
            hora.
          </>
        ) : (
          <>
            A janela agendada{who} será removida. Esta ação não pode ser
            desfeita.
          </>
        )}
      </p>
      <div className="flex gap-2.5 justify-end">
        <Button variant="ghost" onClick={onClose} data-autofocus>
          Voltar
        </Button>
        <Button variant="danger" onClick={confirm} disabled={saving}>
          {saving ? "Aguarde..." : active ? "Desligar" : "Cancelar manutenção"}
        </Button>
      </div>
    </Modal>
  );
}
