import { useState } from "react";
import { Activity, Globe, Server } from "lucide-react";
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

// Modais de editar e excluir um monitor, usados no painel e na página de
// detalhe. `modal` é "edit", "delete" ou null; criar fica no painel.
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
          <Button variant="ghost" onClick={onClose}>
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
