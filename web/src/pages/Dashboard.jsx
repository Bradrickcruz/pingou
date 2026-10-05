import { useState } from "react";
import { Plus, RefreshCw } from "lucide-react";
import { useMonitors } from "../hooks/useMonitors";
import { monitorsApi } from "../api/monitors";
import { MonitorCard } from "../components/monitors/MonitorCard";
import { MonitorForm } from "../components/monitors/MonitorForm";
import { Modal } from "../components/ui/Modal";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import logo from "../assets/brand/pingou-simbolo.svg";

export function Dashboard() {
  const { monitors, loading, error, refetch, isFetching } = useMonitors();
  const [modal, setModal] = useState(null);
  const [selected, setSelected] = useState(null);
  const [saving, setSaving] = useState(false);

  const up = monitors.filter((m) => m.current_state === "UP").length;
  const down = monitors.filter((m) => m.current_state === "DOWN").length;
  const unknown = monitors.filter((m) => m.current_state === "UNKNOWN").length;

  const openEdit = (m) => {
    setSelected(m);
    setModal("edit");
  };
  const openDelete = (m) => {
    setSelected(m);
    setModal("delete");
  };
  const closeModal = () => {
    setModal(null);
    setSelected(null);
  };

  const handleCreate = async (data) => {
    setSaving(true);
    try {
      await monitorsApi.create(data);
      await refetch();
      closeModal();
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = async (data) => {
    setSaving(true);
    try {
      await monitorsApi.update(selected.id, data);
      await refetch();
      closeModal();
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setSaving(true);
    try {
      await monitorsApi.delete(selected.id);
      await refetch();
      closeModal();
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <div className="flex flex-col gap-4 sm:flex-row sm:justify-between sm:items-center mb-7">
        <div>
          <h1 className="text-pg-h1 font-medium text-pg-ardosia my-0">Painel</h1>
          <p className="text-pg-ui text-pg-neblina mt-0.5">
            {monitors.length} monitores · {up} no ar · {down} fora do ar · {unknown}{" "}
            desconhecidos
          </p>
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={refetch}
            disabled={isFetching}
            aria-label="Atualizar lista de monitores"
            aria-busy={isFetching}
            className={`flex items-center justify-center w-11 h-11 rounded-pg-control border border-pg-borda bg-pg-branco transition-colors ${
              isFetching ? "opacity-50 cursor-not-allowed" : "hover:bg-pg-nuvem"
            }`}
          >
            <RefreshCw
              size={16}
              aria-hidden="true"
              className={isFetching ? "animate-spin motion-reduce:animate-none" : ""}
            />
          </button>
          <Button
            onClick={() => setModal("create")}
            className="inline-flex items-center gap-1.5"
          >
            <Plus size={16} aria-hidden="true" />
            Criar monitor
          </Button>
        </div>
      </div>

      {monitors.length > 0 && (
        <div className="grid grid-cols-3 gap-2 sm:gap-3 mb-6">
          {[
            { label: "No ar", value: up, tone: "text-pg-up" },
            { label: "Fora do ar", value: down, tone: "text-pg-down" },
            { label: "Desconhecido", value: unknown, tone: "text-pg-unknown" },
          ].map(({ label, value, tone }) => (
            <div
              key={label}
              className="p-3 sm:p-4 rounded-pg-card border border-pg-borda bg-pg-branco min-w-0"
            >
              <div className={`text-pg-h1 font-medium tabular-nums ${tone}`}>
                {value}
              </div>
              <div className="text-pg-caption text-pg-neblina mt-0.5">{label}</div>
            </div>
          ))}
        </div>
      )}

      {loading && (
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      )}

      {error && <p className="text-pg-ui text-pg-down" role="alert">{error}</p>}

      {!loading && monitors.length === 0 && (
        <div className="text-center py-16 px-4 rounded-pg-card border border-dashed border-pg-borda text-pg-neblina">
          <img src={logo} alt="" width={40} height={40} className="mx-auto mb-3" />
          <p className="font-medium text-pg-ardosia mb-1.5">Nenhum monitor ainda</p>
          <p className="text-pg-ui mb-5">Adicione a primeira URL para começar o monitoramento.</p>
          <Button
            onClick={() => setModal("create")}
            className="inline-flex items-center gap-1.5"
          >
            <Plus size={16} aria-hidden="true" />
            Criar monitor
          </Button>
        </div>
      )}

      <div className="flex flex-col gap-2.5">
        {monitors.map((m) => (
          <MonitorCard
            key={m.id}
            monitor={m}
            onEdit={openEdit}
            onDelete={openDelete}
          />
        ))}
      </div>

      {modal === "create" && (
        <Modal title="Criar monitor" onClose={closeModal}>
          <MonitorForm
            onSubmit={handleCreate}
            onCancel={closeModal}
            loading={saving}
          />
        </Modal>
      )}

      {modal === "edit" && selected && (
        <Modal title="Editar monitor" onClose={closeModal}>
          <MonitorForm
            initial={selected}
            onSubmit={handleEdit}
            onCancel={closeModal}
            loading={saving}
          />
        </Modal>
      )}

      {modal === "delete" && selected && (
        <Modal title="Excluir monitor" onClose={closeModal}>
          <p className="mb-5 text-pg-neblina">
            Excluir <strong className="font-medium text-pg-ardosia">{selected.name}</strong>?
            Esta ação não pode ser desfeita.
          </p>
          <div className="flex gap-2.5 justify-end">
            <Button variant="ghost" onClick={closeModal}>
              Cancelar
            </Button>
            <Button variant="danger" onClick={handleDelete} disabled={saving}>
              {saving ? "Excluindo..." : "Excluir"}
            </Button>
          </div>
        </Modal>
      )}
    </div>
  );
}
