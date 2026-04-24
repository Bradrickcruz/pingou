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
  const [modal, setModal] = useState(null); // null | 'create' | 'edit' | 'delete'
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
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: "28px",
        }}
      >
        <div>
          <h1 style={{ fontSize: "20px", fontWeight: 700 }}>Dashboard</h1>
          <p
            style={{
              color: t.colors.textMuted,
              fontSize: "13px",
              marginTop: "2px",
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
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(3, 1fr)",
            gap: "12px",
            marginBottom: "24px",
          }}
        >
          {[
            { label: "Up", value: up, color: t.colors.success },
            { label: "Down", value: down, color: t.colors.danger },
            { label: "Unknown", value: unknown, color: t.colors.unknown },
          ].map(({ label, value, color }) => (
            <div
              key={label}
              style={{
                background: t.colors.surface,
                border: `1px solid ${t.colors.border}`,
                borderRadius: t.radius.md,
                padding: "16px 20px",
              }}
            >
              <div style={{ fontSize: "28px", fontWeight: 700, color }}>
                {value}
              </div>
              <div
                style={{
                  fontSize: "12px",
                  color: t.colors.textMuted,
                  marginTop: "2px",
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
        <div
          style={{ display: "flex", justifyContent: "center", padding: "48px" }}
        >
          <Spinner />
        </div>
      )}

      {error && <p style={{ color: t.colors.danger }}>{error}</p>}

      {!loading && monitors.length === 0 && (
        <div
          style={{
            textAlign: "center",
            padding: "64px",
            color: t.colors.textMuted,
            border: `1px dashed ${t.colors.border}`,
            borderRadius: t.radius.lg,
          }}
        >
          <p style={{ fontSize: "32px", marginBottom: "12px" }}>🏓</p>
          <p style={{ fontWeight: 600, marginBottom: "6px" }}>
            No monitors yet
          </p>
          <p style={{ fontSize: "13px", marginBottom: "20px" }}>
            Add your first URL to start monitoring
          </p>
          <Button onClick={() => setModal("create")}>+ Add Monitor</Button>
        </div>
      )}

      <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
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
          <p style={{ color: t.colors.textMuted, marginBottom: "20px" }}>
            Delete{" "}
            <strong style={{ color: t.colors.textPrimary }}>
              {selected.name}
            </strong>
            ? This action cannot be undone.
          </p>
          <div
            style={{ display: "flex", gap: "10px", justifyContent: "flex-end" }}
          >
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
