import { useState } from "react";
import { useSettings } from "../hooks/useSettings";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { tokens as t } from "../theme/tokens";

export function Settings() {
  const { settings, loading, update } = useSettings();
  const [form, setForm] = useState(null);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState(null);

  if (loading)
    return (
      <div
        style={{ display: "flex", justifyContent: "center", padding: "48px" }}
      >
        <Spinner />
      </div>
    );

  const current = form ?? settings;

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError(null);
    try {
      await update(current);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (e) {
      setError(e.message);
    } finally {
      setSaving(false);
    }
  };

  const set = (k, v) => setForm((f) => ({ ...(f ?? settings), [k]: v }));

  return (
    <div style={{ maxWidth: "480px" }}>
      <h1 style={{ fontSize: "20px", fontWeight: 700, marginBottom: "24px" }}>
        Settings
      </h1>

      <form onSubmit={handleSubmit}>
        <div
          style={{
            background: t.colors.surface,
            border: `1px solid ${t.colors.border}`,
            borderRadius: t.radius.md,
            padding: "24px",
            display: "flex",
            flexDirection: "column",
            gap: "18px",
          }}
        >
          <div>
            <label
              style={{
                display: "block",
                marginBottom: "6px",
                color: t.colors.textMuted,
                fontSize: "12px",
                fontWeight: 500,
              }}
            >
              Webhook URL
            </label>
            <input
              type="url"
              value={current.webhook_url ?? ""}
              onChange={(e) => set("webhook_url", e.target.value)}
              placeholder="https://hooks.example.com/..."
            />
            <p
              style={{
                fontSize: "11px",
                color: t.colors.textMuted,
                marginTop: "4px",
              }}
            >
              Receives <code>down</code> and <code>up</code> events.
            </p>
          </div>

          <div>
            <label
              style={{
                display: "block",
                marginBottom: "6px",
                color: t.colors.textMuted,
                fontSize: "12px",
                fontWeight: 500,
              }}
            >
              Retention (days)
            </label>
            <input
              type="number"
              min={7}
              max={90}
              value={current.retention_days ?? 30}
              onChange={(e) => set("retention_days", Number(e.target.value))}
            />
            <p
              style={{
                fontSize: "11px",
                color: t.colors.textMuted,
                marginTop: "4px",
              }}
            >
              Checks older than this are automatically deleted. Min 7, max 90.
            </p>
          </div>

          {error && (
            <p style={{ color: t.colors.danger, fontSize: "13px" }}>{error}</p>
          )}

          <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Save Settings"}
            </Button>
            {saved && (
              <span style={{ color: t.colors.success, fontSize: "13px" }}>
                ✓ Saved
              </span>
            )}
          </div>
        </div>
      </form>

      <div
        style={{
          marginTop: "24px",
          background: t.colors.surface,
          border: `1px solid ${t.colors.border}`,
          borderRadius: t.radius.md,
          padding: "20px",
        }}
      >
        <p style={{ fontWeight: 600, marginBottom: "12px", fontSize: "13px" }}>
          Database Export
        </p>
        <p
          style={{
            color: t.colors.textMuted,
            fontSize: "12px",
            marginBottom: "14px",
          }}
        >
          Download a full SQLite dump of all monitors, checks and incidents.
        </p>
        <a
          href="/api/export"
          style={{
            display: "inline-block",
            background: t.colors.surfaceAlt,
            color: t.colors.textPrimary,
            padding: "8px 16px",
            borderRadius: t.radius.sm,
            fontSize: "13px",
            fontWeight: 600,
            border: `1px solid ${t.colors.border}`,
          }}
        >
          ↓ Download dump
        </a>
      </div>
    </div>
  );
}
