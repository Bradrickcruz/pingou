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
    <div style={{ marginBottom: "14px" }}>
      <label
        style={{
          display: "block",
          marginBottom: "6px",
          color: t.colors.textMuted,
          fontSize: "12px",
          fontWeight: 500,
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

      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: "10px",
          marginBottom: "20px",
        }}
      >
        <input
          type="checkbox"
          id="enabled"
          checked={form.enabled}
          onChange={(e) => set("enabled", e.target.checked)}
          style={{ width: "auto" }}
        />
        <label
          htmlFor="enabled"
          style={{ color: t.colors.textMuted, fontSize: "13px" }}
        >
          Enabled
        </label>
      </div>

      {error && (
        <p
          style={{
            color: t.colors.danger,
            fontSize: "13px",
            marginBottom: "14px",
          }}
        >
          {error}
        </p>
      )}

      <div style={{ display: "flex", gap: "10px", justifyContent: "flex-end" }}>
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
