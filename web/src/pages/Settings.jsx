import { useState } from "react";
import { useSettings } from "../hooks/useSettings";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { tokens as t } from "../theme/tokens";
import { client } from "../api/client";

export function Settings() {
  const { settings, loading, update } = useSettings();
  const [form, setForm] = useState(null);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState(null);

  if (loading)
    return (
      <div className="flex justify-center py-12">
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

  const handleExport = async () => {
    try {
      const res = await client.get("/export", { responseType: "blob" });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute(
        "download",
        `pingou_bkp_${new Date().toISOString().slice(0, 19).replace(/[-:]/g, "")}.db`,
      );
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div className="max-w-[480px]">
      <h1 className="text-xl font-bold mb-6">Settings</h1>

      <form onSubmit={handleSubmit}>
        <div
          className="p-6 rounded-md border flex flex-col gap-4.5"
          style={{
            background: t.colors.surface,
            borderColor: t.colors.border,
            borderRadius: t.radius.md,
          }}
        >
          <div>
            <label
              className="block mb-1.5 text-xs font-medium"
              style={{
                color: t.colors.textMuted,
              }}
            >
              Webhook URL
            </label>
            <input
              type="url"
              value={current.webhook_url ?? ""}
              onChange={(e) => set("webhook_url", e.target.value)}
              placeholder="https://hooks.example.com/..."
              className="w-full px-3 py-2 rounded text-sm bg-[var(--bg)] border border-[var(--border)] text-[var(--text-h)] focus:outline-none focus:border-[var(--accent)]"
            />
            <p
              className="text-[11px] mt-1"
              style={{
                color: t.colors.textMuted,
              }}
            >
              Receives <code>down</code> and <code>up</code> events.
            </p>
          </div>

          <div>
            <label
              className="block mb-1.5 text-xs font-medium"
              style={{
                color: t.colors.textMuted,
              }}
            >
              Retention (days)
            </label>
            <div className="grid grid-cols-5 gap-2">
              {[7, 14, 30, 60, 90].map((days) => (
                <label
                  key={days}
                  className="flex items-center justify-center gap-2 py-2 rounded cursor-pointer border transition-colors"
                  style={{
                    background:
                      current.retention_days === days
                        ? t.colors.primary
                        : t.colors.surface,
                    borderColor:
                      current.retention_days === days
                        ? t.colors.primary
                        : t.colors.border,
                    color:
                      current.retention_days === days
                        ? "#fff"
                        : t.colors.textPrimary,
                  }}
                >
                  <input
                    type="radio"
                    name="retention_days"
                    value={days}
                    checked={current.retention_days === days}
                    onChange={() => set("retention_days", days)}
                    className="sr-only"
                  />
                  <span className="text-sm font-medium">{days}</span>
                </label>
              ))}
            </div>
            <p
              className="text-[11px] mt-1.5"
              style={{
                color: t.colors.textMuted,
              }}
            >
              Checks older than this are automatically deleted.
            </p>
          </div>

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

          <div className="flex items-center gap-3">
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Save Settings"}
            </Button>
            {saved && (
              <span
                className="text-sm"
                style={{
                  color: t.colors.success,
                }}
              >
                ✓ Saved
              </span>
            )}
          </div>
        </div>
      </form>

      <div
        className="mt-6 p-5 rounded-md border"
        style={{
          background: t.colors.surface,
          borderColor: t.colors.border,
          borderRadius: t.radius.md,
        }}
      >
        <p className="font-semibold mb-3 text-sm">Database Export</p>
        <p
          className="text-xs mb-3.5"
          style={{
            color: t.colors.textMuted,
          }}
        >
          Download a full SQLite dump of all monitors, checks and incidents.
        </p>
        <button
          onClick={handleExport}
          className="inline-block px-4 py-2 rounded text-sm font-semibold border cursor-pointer"
          style={{
            background: t.colors.surfaceAlt,
            color: t.colors.textPrimary,
            borderColor: t.colors.border,
          }}
        >
          ↓ Download dump
        </button>
      </div>
    </div>
  );
}