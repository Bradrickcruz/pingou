import { useState } from "react";
import { Check, Download } from "lucide-react";
import { useSettings } from "../hooks/useSettings";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
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
      <h1 className="text-pg-h1 font-medium text-pg-ardosia mb-6 my-0">Settings</h1>

      <form onSubmit={handleSubmit}>
        <div className="p-6 rounded-pg-card border border-pg-borda bg-pg-branco flex flex-col gap-4">
          <div>
            <label className="block mb-1.5 text-pg-caption font-medium text-pg-neblina">
              Webhook URL
            </label>
            <input
              type="url"
              value={current.webhook_url ?? ""}
              onChange={(e) => set("webhook_url", e.target.value)}
              placeholder="https://hooks.example.com/..."
              className="w-full px-3 py-2 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda"
            />
            <p className="text-pg-caption text-pg-neblina mt-1">
              Receives <code>down</code> and <code>up</code> events.
            </p>
          </div>

          <div>
            <span className="block mb-1.5 text-pg-caption font-medium text-pg-neblina">
              Retention (days)
            </span>
            <div className="grid grid-cols-5 gap-2">
              {[7, 14, 30, 60, 90].map((days) => {
                const selected = current.retention_days === days;
                return (
                  <label
                    key={days}
                    className={`flex items-center justify-center py-2 rounded-pg-control border cursor-pointer transition-colors has-[:focus-visible]:outline-2 has-[:focus-visible]:outline has-[:focus-visible]:outline-pg-ciano ${
                      selected
                        ? "bg-pg-ciano-nevoa border-pg-ciano-fundo text-pg-ciano-fundo font-medium"
                        : "bg-pg-branco border-pg-borda text-pg-ardosia"
                    }`}
                  >
                    <input
                      type="radio"
                      name="retention_days"
                      value={days}
                      checked={selected}
                      onChange={() => set("retention_days", days)}
                      className="sr-only"
                    />
                    <span className="text-pg-ui font-medium">{days}</span>
                  </label>
                );
              })}
            </div>
            <p className="text-pg-caption text-pg-neblina mt-1.5">
              Checks older than this are automatically deleted.
            </p>
          </div>

          {error && <p className="text-pg-ui text-pg-down">{error}</p>}

          <div className="flex items-center gap-3">
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Save Settings"}
            </Button>
            {saved && (
              <span className="text-pg-ui inline-flex items-center gap-1 text-pg-up">
                <Check size={14} aria-hidden="true" />
                Saved
              </span>
            )}
          </div>
        </div>
      </form>

      <div className="mt-6 p-5 rounded-pg-card border border-pg-borda bg-pg-branco">
        <p className="font-medium mb-3 text-pg-ui text-pg-ardosia">Database Export</p>
        <p className="text-pg-caption text-pg-neblina mb-3.5">
          Download a full SQLite dump of all monitors, checks and incidents.
        </p>
        <Button variant="ghost" onClick={handleExport} className="gap-1.5">
          <Download size={16} aria-hidden="true" />
          Download dump
        </Button>
      </div>
    </div>
  );
}
