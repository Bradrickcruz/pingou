import { useId, useState } from "react";
import { Activity, Globe, Server } from "lucide-react";
import { toast } from "react-toastify";
import { monitorsApi } from "../../api/monitors";
import { MonitorForm } from "./MonitorForm";
import { Modal } from "../ui/Modal";
import { Button } from "../ui/Button";

const typeLabel = { http: "HTTP", tcp: "TCP", heartbeat: "Heartbeat" };
const typeIcon = { http: Globe, tcp: Server, heartbeat: Activity };

function TypeNote({ type }) {
  const Icon = typeIcon[type] ?? Globe;
  return (
    <span className="inline-flex items-center gap-1.5">
      <Icon size={14} aria-hidden="true" />
      {typeLabel[type]} · o tipo não pode ser alterado
    </span>
  );
}

export const REASON_MAX = 200;

// Motivo opcional (pausa ou manutenção): textarea com contador de caracteres.
export function ReasonField({ value, onChange, label, placeholder }) {
  const id = useId();
  const counterId = useId();
  return (
    <div className="mb-5">
      <label
        htmlFor={id}
        className="block mb-2 text-pg-ui font-medium text-pg-ardosia"
      >
        {label} <span className="font-normal text-pg-neblina">(opcional)</span>
      </label>
      <textarea
        id={id}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        maxLength={REASON_MAX}
        rows={2}
        placeholder={placeholder}
        aria-describedby={counterId}
        className="w-full px-3 py-2.5 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda placeholder:text-pg-neblina resize-y"
      />
      <p
        id={counterId}
        className="mt-1 mb-0 text-right text-pg-caption text-pg-neblina tabular-nums"
      >
        {value.length}/{REASON_MAX}
      </p>
    </div>
  );
}

// Pausar com motivo opcional. Enter confirma (Shift+Enter quebra linha).
function PauseModal({ monitor, onConfirm, onClose, saving }) {
  const [reason, setReason] = useState("");
  const [error, setError] = useState(null);
  const submit = async (e) => {
    e.preventDefault();
    if (saving) return;
    setError(null);
    try {
      await onConfirm(reason.trim());
    } catch (err) {
      setError(err.message);
    }
  };
  return (
    <Modal title="Pausar monitor" onClose={onClose}>
      <form
        onSubmit={submit}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.shiftKey && !e.nativeEvent.isComposing)
            submit(e);
        }}
      >
        <p className="mb-4 text-pg-ui text-pg-neblina">
          Enquanto{" "}
          <strong className="font-medium text-pg-ardosia">
            {monitor.name}
          </strong>{" "}
          estiver pausado, nenhum check roda e ninguém é avisado. Um incidente
          aberto é fechado sem aviso.
        </p>
        <ReasonField
          label="Motivo"
          value={reason}
          onChange={setReason}
          placeholder="Ex.: troca do servidor"
        />
        {error && (
          <p role="alert" className="-mt-2 mb-4 text-pg-caption text-pg-down">
            {error}
          </p>
        )}
        <div className="flex gap-2.5 justify-end">
          <Button variant="ghost" onClick={onClose}>
            Cancelar
          </Button>
          <Button type="submit" disabled={saving}>
            {saving ? "Pausando..." : "Pausar"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}

// Modais de editar, pausar e excluir um monitor, usados no painel e na página
// de detalhe. `modal` é "edit", "pause", "delete" ou null; criar fica no painel.
// onSaved/onDeleted rodam antes de fechar (ex.: refetch ou navegar).
export function MonitorModals({
  modal,
  selected,
  onModalChange,
  onSelectedChange,
  onClose,
  onSaved,
  onDeleted,
}) {
  const [saving, setSaving] = useState(false);

  const run = async (fn, after) => {
    setSaving(true);
    try {
      await fn();
      await after?.();
      onClose();
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = (data) =>
    run(() => monitorsApi.update(selected.id, data), onSaved);
  const handleDelete = () =>
    run(() => monitorsApi.delete(selected.id), onDeleted);
  const handlePause = (reason) =>
    run(async () => {
      await monitorsApi.update(selected.id, {
        enabled: false,
        pause_reason: reason,
      });
      toast.success("Monitor pausado");
    }, onSaved);

  const handleRegenerate = async () => {
    const updated = await monitorsApi.regenerateToken(selected.id);
    onSelectedChange((s) => ({ ...s, ...updated }));
    void onSaved?.();
  };

  if (!selected) return null;

  if (modal === "edit") {
    return (
      <Modal
        title="Editar monitor"
        subtitle={<TypeNote type={selected.type} />}
        onClose={onClose}
      >
        <MonitorForm
          initial={selected}
          onSubmit={handleEdit}
          onCancel={onClose}
          onDelete={() => onModalChange("delete")}
          onRegenerate={handleRegenerate}
          loading={saving}
        />
      </Modal>
    );
  }

  if (modal === "pause") {
    return (
      <PauseModal
        monitor={selected}
        onConfirm={handlePause}
        onClose={onClose}
        saving={saving}
      />
    );
  }

  if (modal === "delete") {
    return (
      <Modal title="Excluir monitor" onClose={onClose}>
        <p className="mb-5 text-pg-neblina">
          Excluir{" "}
          <strong className="font-medium text-pg-ardosia">
            {selected.name}
          </strong>
          ? O histórico de checks e incidentes também é apagado. Esta ação não
          pode ser desfeita.
        </p>
        <div className="flex gap-2.5 justify-end">
          <Button variant="ghost" onClick={onClose} data-autofocus>
            Cancelar
          </Button>
          <Button variant="danger" onClick={handleDelete} disabled={saving}>
            {saving ? "Excluindo..." : "Excluir"}
          </Button>
        </div>
      </Modal>
    );
  }

  return null;
}
