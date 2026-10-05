import { useState } from "react";
import { Button } from "../ui/Button";
import { TimeRangeSlider } from "../ui/TimeRangeSlider";

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

  const field = (label, key, type = "text", extra = {}) => {
    const id = `monitor-${key}`;
    return (
      <div className="mb-3.5">
        <label htmlFor={id} className="block mb-1.5 text-pg-caption font-medium text-pg-neblina">
          {label}
        </label>
        <input
          id={id}
          type={type}
          value={form[key]}
          onChange={(e) => set(key, e.target.value)}
          className="w-full px-3 py-2 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda"
          {...extra}
        />
      </div>
    );
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="flex items-center justify-between mb-4">
        <span id="monitor-enabled-label" className="text-pg-ui font-medium text-pg-ardosia">
          Ativado
        </span>
        <button
          type="button"
          role="switch"
          aria-checked={form.enabled}
          aria-labelledby="monitor-enabled-label"
          onClick={() => set("enabled", !form.enabled)}
          className={`relative w-11 h-6 rounded-full transition-colors focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
            form.enabled ? "bg-pg-ciano-fundo" : "bg-pg-borda"
          }`}
        >
          <span
            className={`absolute top-0.5 left-0.5 w-5 h-5 rounded-full bg-pg-branco transition-transform ${
              form.enabled ? "translate-x-5" : "translate-x-0"
            }`}
          />
        </button>
      </div>

      {field("URL", "url", "url", { required: true, placeholder: "https://exemplo.com" })}
      {field("Nome", "name", "text", { required: true, placeholder: "Minha API" })}
      <TimeRangeSlider
        label="Intervalo (segundos)"
        value={form.interval_seconds}
        onChange={(v) => set("interval_seconds", v)}
        min={10}
        max={3600}
        step={10}
      />
      <TimeRangeSlider
        label="Tempo limite (segundos)"
        value={form.timeout_seconds}
        onChange={(v) => set("timeout_seconds", v)}
        min={5}
        max={60}
        step={5}
      />
      {field("Limite de falhas", "failure_threshold", "number", { min: 1, max: 10 })}

      {error && (
        <p role="alert" className="text-pg-ui text-pg-down mb-3.5">
          {error}
        </p>
      )}

      <div className="flex gap-2.5 justify-end">
        <Button variant="ghost" onClick={onCancel} type="button">
          Cancelar
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? "Salvando..." : "Salvar"}
        </Button>
      </div>
    </form>
  );
}
