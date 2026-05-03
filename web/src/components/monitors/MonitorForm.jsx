import { useState } from "react";
import { Button } from "../ui/Button";
import { tokens as t } from "../../theme/tokens";

const defaults = {
  name: "",
  url: "",
  interval_seconds: 60,
  timeout_seconds: 10,
  failure_threshold: 3,
  enabled: true,
};

export function MonitorForm({ initial = {}, onSubmit, onCancel, loading }) {
  const [form, setForm] = useState({ ...defaults, ...initial });
  const [error, setError] = useState(null);

  const set = (k, v) => setForm((f) => ({ ...f, [k]: v }));

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    try {
      await onSubmit({
        ...form,
        interval_seconds: Number(form.interval_seconds),
        timeout_seconds: Number(form.timeout_seconds),
        failure_threshold: Number(form.failure_threshold),
      });
    } catch (e) {
      setError(e.message);
    }
  };

  const field = (label, key, type = "text", extra = {}) => (
    <div className="mb-3.5">
      <label
        className="block mb-1.5 text-xs font-medium"
        style={{
          color: t.colors.textMuted,
        }}
      >
        {label}
      </label>
      <input
        type={type}
        value={form[key]}
        onChange={(e) =>
          set(key, type === "number" ? e.target.value : e.target.value)
        }
        className="w-full px-3 py-2 rounded text-sm bg-[var(--bg)] border border-[var(--border)] text-[var(--text-h)] focus:outline-none focus:border-[var(--accent)]"
        {...extra}
      />
    </div>
  );

  return (
    <form onSubmit={handleSubmit}>
      {field("Name", "name", "text", { required: true, placeholder: "My API" })}
      {field("URL", "url", "url", {
        required: true,
        placeholder: "https://example.com",
      })}
      {field("Interval (seconds)", "interval_seconds", "number", {
        min: 10,
        max: 3600,
      })}
      {field("Timeout (seconds)", "timeout_seconds", "number", {
        min: 1,
        max: 60,
      })}
      {field("Failure threshold", "failure_threshold", "number", {
        min: 1,
        max: 10,
      })}

      <div className="flex items-center gap-2.5 mb-5">
        <input
          type="checkbox"
          id="enabled"
          checked={form.enabled}
          onChange={(e) => set("enabled", e.target.checked)}
          className="w-auto"
        />
        <label
          htmlFor="enabled"
          className="text-sm"
          style={{
            color: t.colors.textMuted,
          }}
        >
          Enabled
        </label>
      </div>

      {error && (
        <p
          className="text-sm mb-3.5"
          style={{
            color: t.colors.danger,
          }}
        >
          {error}
        </p>
      )}

      <div className="flex gap-2.5 justify-end">
        <Button variant="ghost" onClick={onCancel} type="button">
          Cancel
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? "Saving..." : "Save"}
        </Button>
      </div>
    </form>
  );
}