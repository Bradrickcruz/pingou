import { useState } from "react";
import { Button } from "../ui/Button";
import { TimeRangeSlider } from "../ui/TimeRangeSlider";
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
      <div className="flex items-center justify-between mb-4">
        <span
          className="text-sm font-medium"
          style={{ color: t.colors.textPrimary }}
        >
          Enabled
        </span>
        <button
          type="button"
          onClick={() => set("enabled", !form.enabled)}
          className="relative w-11 h-6 rounded-full transition-colors"
          style={{
            background: form.enabled ? t.colors.primary : t.colors.surface,
          }}
        >
          <span
            className="absolute top-0.5 left-0.5 w-5 h-5 rounded-full bg-white shadow transition-transform"
            style={{
              transform: form.enabled ? "translateX(20px)" : "translateX(0)",
            }}
          />
        </button>
      </div>

      {field("URL", "url", "url", {
        required: true,
        placeholder: "https://example.com",
      })}
      {field("Name", "name", "text", { required: true, placeholder: "My API" })}
      <TimeRangeSlider
        label="Interval (seconds)"
        value={form.interval_seconds}
        onChange={(v) => set("interval_seconds", v)}
        min={10}
        max={3600}
        step={10}
      />
      <TimeRangeSlider
        label="Timeout (seconds)"
        value={form.timeout_seconds}
        onChange={(v) => set("timeout_seconds", v)}
        min={5}
        max={60}
        step={5}
      />
      {field("Failure threshold", "failure_threshold", "number", {
        min: 1,
        max: 10,
      })}

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
