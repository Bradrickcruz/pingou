import { useState } from "react";
import { useMonitors } from "../hooks/useMonitors";
import { monitorsApi } from "../api/monitors";
import { MonitorCard } from "../components/monitors/MonitorCard";
import { MonitorForm } from "../components/monitors/MonitorForm";
import { Modal } from "../components/ui/Modal";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { tokens as t } from "../theme/tokens";

export function Dashboard() {
  const { monitors, loading, error, refetch } = useMonitors();
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
      {/* header */}
      <div className="flex justify-between items-center mb-7">
        <div>
          <h1 className="text-xl font-bold">Dashboard</h1>
          <p
            className="text-sm mt-0.5"
            style={{
              color: t.colors.textMuted,
            }}
          >
            {monitors.length} monitors · {up} up · {down} down · {unknown}{" "}
            unknown
          </p>
        </div>
        <Button onClick={() => setModal("create")}>+ Add Monitor</Button>
      </div>

      {/* stats */}
      {monitors.length > 0 && (
        <div className="grid grid-cols-3 gap-3 mb-6">
          {[
            { label: "Up", value: up, color: t.colors.success },
            { label: "Down", value: down, color: t.colors.danger },
            { label: "Unknown", value: unknown, color: t.colors.unknown },
          ].map(({ label, value, color }) => (
            <div
              key={label}
              className="p-4 rounded-md border"
              style={{
                background: t.colors.surface,
                borderColor: t.colors.border,
                borderRadius: t.radius.md,
              }}
            >
              <div
                className="text-[28px] font-bold"
                style={{
                  color,
                }}
              >
                {value}
              </div>
              <div
                className="text-xs mt-0.5"
                style={{
                  color: t.colors.textMuted,
                }}
              >
                {label}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* list */}
      {loading && (
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      )}

      {error && (
        <p
          className="text-sm"
          style={{
            color: t.colors.danger,
          }}
        >
          {error}
        </p>
      )}

      {!loading && monitors.length === 0 && (
        <div
          className="text-center py-16 px-4 rounded-lg border border-dashed"
          style={{
            color: t.colors.textMuted,
            borderColor: t.colors.border,
            borderRadius: t.radius.lg,
          }}
        >
          <p className="text-[32px] mb-3">🏓</p>
          <p className="font-semibold mb-1.5">No monitors yet</p>
          <p className="text-sm mb-5">Add your first URL to start monitoring</p>
          <Button onClick={() => setModal("create")}>+ Add Monitor</Button>
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

      {/* modais */}
      {modal === "create" && (
        <Modal title="Add Monitor" onClose={closeModal}>
          <MonitorForm
            onSubmit={handleCreate}
            onCancel={closeModal}
            loading={saving}
          />
        </Modal>
      )}

      {modal === "edit" && selected && (
        <Modal title="Edit Monitor" onClose={closeModal}>
          <MonitorForm
            initial={selected}
            onSubmit={handleEdit}
            onCancel={closeModal}
            loading={saving}
          />
        </Modal>
      )}

      {modal === "delete" && selected && (
        <Modal title="Delete Monitor" onClose={closeModal}>
          <p
            className="mb-5"
            style={{
              color: t.colors.textMuted,
            }}
          >
            Delete{" "}
            <strong
              style={{
                color: t.colors.textPrimary,
              }}
            >
              {selected.name}
            </strong>
            ? This action cannot be undone.
          </p>
          <div className="flex gap-2.5 justify-end">
            <Button variant="ghost" onClick={closeModal}>
              Cancel
            </Button>
            <Button variant="danger" onClick={handleDelete} disabled={saving}>
              {saving ? "Deleting..." : "Delete"}
            </Button>
          </div>
        </Modal>
      )}
    </div>
  );
}