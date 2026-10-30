import { useId, useState } from "react";
import {
  Globe,
  Server,
  Activity,
  Copy,
  Check,
  RefreshCw,
  Trash2,
} from "lucide-react";
import { Button } from "../ui/Button";
import { ChipGroup } from "../ui/ChipGroup";
import { formatDuration, plural, timeAgo } from "../../lib/format";
import { copyText, pingUrl } from "../../lib/monitors";

const TYPES = [
  {
    value: "http",
    label: "HTTP",
    hint: "Um site ou API responde",
    icon: Globe,
  },
  {
    value: "tcp",
    label: "Porta (TCP)",
    hint: "Uma porta aceita conexão",
    icon: Server,
  },
  {
    value: "heartbeat",
    label: "Heartbeat",
    hint: "Seu job avisa que rodou",
    icon: Activity,
  },
];

const CHECK_EVERY = [
  { value: 30, label: "30 s" },
  { value: 60, label: "1 min" },
  { value: 300, label: "5 min" },
  { value: 900, label: "15 min" },
  { value: 3600, label: "1 h" },
];

const PING_EVERY = [
  { value: 3600, label: "1 h" },
  { value: 21600, label: "6 h" },
  { value: 86400, label: "1 dia" },
  { value: 604800, label: "1 semana" },
];

const GRACE = [
  { value: 60, label: "1 min" },
  { value: 900, label: "15 min" },
  { value: 3600, label: "1 h" },
];

const blank = {
  type: "http",
  name: "",
  url: "",
  host: "",
  port: "",
  interval_seconds: 60,
  timeout_seconds: 10,
  failure_threshold: 3,
  grace_seconds: 3600,
};

function fromMonitor(m) {
  return {
    ...blank,
    type: m.type,
    name: m.name,
    url: m.url ?? "",
    host: m.config?.host ?? "",
    port: m.config?.port ?? "",
    interval_seconds: m.interval_seconds,
    timeout_seconds: m.timeout_seconds,
    failure_threshold: m.failure_threshold,
    grace_seconds: m.config?.grace_seconds ?? blank.grace_seconds,
  };
}

const inputCls =
  "w-full px-3 py-2.5 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda placeholder:text-pg-neblina";

function Field({ label, children, id }) {
  return (
    <div className="mb-5 min-w-0">
      <label
        htmlFor={id}
        className="block mb-2 text-pg-ui font-medium text-pg-ardosia"
      >
        {label}
      </label>
      {children}
    </div>
  );
}

function NumberWithUnit({ id, label, unit, value, onChange, min, max }) {
  return (
    <div className="min-w-0">
      <label
        htmlFor={id}
        className="block mb-2 text-pg-ui font-medium text-pg-ardosia"
      >
        {label}
      </label>
      <div className="flex items-center rounded-pg-control border border-pg-borda bg-pg-branco pr-3 focus-within:outline-2 focus-within:outline focus-within:outline-pg-ciano">
        <input
          id={id}
          type="number"
          inputMode="numeric"
          min={min}
          max={max}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="w-full min-w-0 px-3 py-2.5 border-0 rounded-pg-control text-pg-ui text-pg-ardosia bg-transparent focus:outline-none"
        />
        <span className="text-pg-caption text-pg-neblina">{unit}</span>
      </div>
    </div>
  );
}

// Bloco da URL de ping do heartbeat (edição e página de detalhe).
export function PingBlock({ monitor, onRegenerate }) {
  const [copied, setCopied] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState(null);
  const inputId = useId();
  const url = pingUrl(monitor.ping_token);
  const cron = `0 3 * * * backup.sh && curl -fsS -m 10 --retry 3 ${url}`;

  const copy = async () => {
    if (!(await copyText(url))) {
      // nem o fallback copiou: deixa selecionado para Ctrl+C
      document.getElementById(inputId)?.select();
      return;
    }
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const regenerate = async () => {
    setBusy(true);
    setError(null);
    try {
      await onRegenerate();
      setConfirming(false);
    } catch (e) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  };

  return (
    <section className="mb-6 p-4 sm:p-5 rounded-pg-card border border-pg-borda bg-pg-nuvem">
      <h3 className="text-pg-ui font-medium text-pg-ardosia my-0">
        URL de ping
      </h3>
      <p className="mt-1 mb-3 text-pg-caption text-pg-neblina">
        Chame esta URL sempre que o job terminar com sucesso.
      </p>
      <div className="flex gap-2">
        <label htmlFor={inputId} className="sr-only">
          URL de ping
        </label>
        <input
          id={inputId}
          readOnly
          value={url}
          onFocus={(e) => e.target.select()}
          className="flex-1 min-w-0 px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-caption text-pg-ciano-fundo bg-pg-branco border border-pg-borda"
        />
        <Button
          variant="ghost"
          onClick={copy}
          className="gap-1.5 bg-pg-branco flex-shrink-0"
        >
          {copied ? (
            <Check size={16} aria-hidden="true" />
          ) : (
            <Copy size={16} aria-hidden="true" />
          )}
          {copied ? "Copiado" : "Copiar"}
        </Button>
      </div>
      <span className="sr-only" aria-live="polite">
        {copied ? "URL copiada" : ""}
      </span>

      <div className="mt-4 mb-1.5 text-pg-caption text-pg-neblina">
        Exemplo com cron
      </div>
      <pre className="m-0 p-3 rounded-pg-control bg-pg-ardosia text-pg-branco font-pg-mono text-pg-caption whitespace-pre-wrap break-all">
        {cron}
      </pre>

      {confirming ? (
        <div
          role="alertdialog"
          aria-label="Gerar nova URL"
          className="mt-4 p-3 rounded-pg-control border border-pg-borda bg-pg-branco"
        >
          <p className="text-pg-ui text-pg-ardosia mb-3">
            Gerar uma nova URL? A atual para de funcionar na hora, e o cron
            precisa ser atualizado.
          </p>
          <div className="flex gap-2 justify-end">
            <Button
              variant="ghost"
              onClick={() => setConfirming(false)}
              disabled={busy}
            >
              Cancelar
            </Button>
            <Button variant="danger" onClick={regenerate} disabled={busy}>
              {busy ? "Gerando..." : "Gerar nova URL"}
            </Button>
          </div>
        </div>
      ) : (
        <div className="flex flex-wrap items-center justify-between gap-2 mt-4">
          <span className="text-pg-caption text-pg-neblina">
            {monitor.last_ping_at
              ? `Último ping recebido ${timeAgo(monitor.last_ping_at)}`
              : "Nenhum ping recebido ainda"}
          </span>
          <button
            type="button"
            onClick={() => setConfirming(true)}
            className="inline-flex items-center gap-1.5 px-2 py-1 -mr-2 rounded-pg-control text-pg-ui font-medium text-pg-ciano-fundo hover:bg-pg-ciano-nevoa focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
          >
            <RefreshCw size={16} aria-hidden="true" />
            Gerar nova URL
          </button>
        </div>
      )}
      {error && (
        <p role="alert" className="mt-2 text-pg-caption text-pg-down">
          {error}
        </p>
      )}
    </section>
  );
}

// Formulário de Novo/Editar monitor. Com `initial`, edita (tipo travado).
export function MonitorForm({
  initial,
  onSubmit,
  onCancel,
  onDelete,
  onRegenerate,
  loading,
}) {
  const editing = Boolean(initial);
  const [form, setForm] = useState(() =>
    editing ? fromMonitor(initial) : blank,
  );
  const [error, setError] = useState(null);
  const ids = {
    name: useId(),
    url: useId(),
    host: useId(),
    port: useId(),
    timeout: useId(),
    threshold: useId(),
  };

  const set = (k, v) => setForm((f) => ({ ...f, [k]: v }));
  const isHeartbeat = form.type === "heartbeat";

  const chooseType = (type) =>
    setForm((f) => ({
      ...f,
      type,
      // intervalos de heartbeat e de check têm faixas diferentes
      interval_seconds:
        type === "heartbeat"
          ? 86400
          : f.type === "heartbeat"
            ? 60
            : f.interval_seconds,
    }));

  const payload = () => {
    const base = {
      name: form.name.trim(),
      interval_seconds: Number(form.interval_seconds),
    };
    if (!editing) {
      base.type = form.type;
      base.enabled = true;
    }
    if (isHeartbeat)
      return { ...base, config: { grace_seconds: Number(form.grace_seconds) } };

    base.timeout_seconds = Number(form.timeout_seconds);
    base.failure_threshold = Number(form.failure_threshold);
    if (form.type === "tcp") {
      return {
        ...base,
        config: { host: form.host.trim(), port: Number(form.port) },
      };
    }
    return { ...base, url: form.url.trim() };
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    // erro mais comum do TCP, explicado antes de ir à API
    if (form.type === "tcp" && /:\/\/|\//.test(form.host)) {
      setError("Informe só o host (ex.: db.interno), sem http:// nem caminho.");
      return;
    }
    try {
      await onSubmit(payload());
    } catch (err) {
      setError(
        err.message.replace(/^validation error: /, "Verifique os campos: "),
      );
    }
  };

  const threshold = Number(form.failure_threshold) || 0;
  const downAfter = formatDuration(threshold * Number(form.interval_seconds));
  const deadline = formatDuration(
    Number(form.interval_seconds) + Number(form.grace_seconds),
  );

  return (
    <form onSubmit={handleSubmit}>
      {!editing && (
        <fieldset className="mb-5 border-0 p-0 m-0">
          <legend className="mb-2 text-pg-ui font-medium text-pg-ardosia">
            O que você quer monitorar?
          </legend>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-2">
            {TYPES.map(({ value, label, hint, icon: Icon }) => {
              const active = form.type === value;
              return (
                <label
                  key={value}
                  className={`relative flex flex-col gap-1 p-3 rounded-pg-control border cursor-pointer transition-colors focus-within:outline-2 focus-within:outline focus-within:outline-pg-ciano ${
                    active
                      ? "bg-pg-ciano-nevoa border-pg-ciano-fundo"
                      : "bg-pg-branco border-pg-borda hover:bg-pg-nuvem"
                  }`}
                >
                  <input
                    type="radio"
                    name="monitor-type"
                    value={value}
                    checked={active}
                    onChange={() => chooseType(value)}
                    className="sr-only"
                  />
                  <span className="flex items-center gap-2 text-pg-ui font-medium text-pg-ardosia">
                    <Icon
                      size={16}
                      aria-hidden="true"
                      className="text-pg-ciano-fundo"
                    />
                    {label}
                  </span>
                  <span className="text-pg-caption text-pg-neblina">
                    {hint}
                  </span>
                </label>
              );
            })}
          </div>
        </fieldset>
      )}

      {editing && isHeartbeat && initial.ping_token && (
        <PingBlock monitor={initial} onRegenerate={onRegenerate} />
      )}

      <Field label="Nome" id={ids.name}>
        <input
          id={ids.name}
          required
          maxLength={100}
          value={form.name}
          onChange={(e) => set("name", e.target.value)}
          placeholder="Ex.: API de pagamentos"
          className={inputCls}
        />
      </Field>

      {form.type === "http" && (
        <Field label="URL" id={ids.url}>
          <input
            id={ids.url}
            type="url"
            required
            value={form.url}
            onChange={(e) => set("url", e.target.value)}
            placeholder="https://api.exemplo.com/health"
            className={`${inputCls} font-pg-mono`}
          />
        </Field>
      )}

      {form.type === "tcp" && (
        <div className="grid grid-cols-[1fr_7rem] gap-3">
          <Field label="Host" id={ids.host}>
            <input
              id={ids.host}
              required
              value={form.host}
              onChange={(e) => set("host", e.target.value)}
              placeholder="db.interno"
              autoCapitalize="none"
              spellCheck={false}
              className={`${inputCls} font-pg-mono`}
            />
          </Field>
          <Field label="Porta" id={ids.port}>
            <input
              id={ids.port}
              type="number"
              inputMode="numeric"
              required
              min={1}
              max={65535}
              value={form.port}
              onChange={(e) => set("port", e.target.value)}
              placeholder="5432"
              className={`${inputCls} font-pg-mono`}
            />
          </Field>
        </div>
      )}

      {isHeartbeat ? (
        <>
          <ChipGroup
            key="ping-every"
            label="Esperar um ping a cada"
            options={PING_EVERY}
            value={form.interval_seconds}
            onChange={(v) => set("interval_seconds", v)}
            min={60}
            max={604800}
          />
          <ChipGroup
            label="Tolerância para atraso"
            options={GRACE}
            value={form.grace_seconds}
            onChange={(v) => set("grace_seconds", v)}
            min={0}
            max={86400}
          />
          <p className="-mt-2 mb-5 text-pg-caption text-pg-neblina">
            Se nenhum ping chegar em {deadline}, o monitor fica fora do ar.
          </p>
        </>
      ) : (
        <>
          <ChipGroup
            key="check-every"
            label="Verificar a cada"
            options={CHECK_EVERY}
            value={form.interval_seconds}
            onChange={(v) => set("interval_seconds", v)}
            min={10}
            max={86400}
          />
          <fieldset className="mb-5 px-4 pt-2 pb-4 rounded-pg-card border border-pg-borda">
            <legend className="px-1.5 text-pg-caption text-pg-neblina">
              Avançado
            </legend>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <NumberWithUnit
                id={ids.timeout}
                label="Tempo limite"
                unit="segundos"
                min={1}
                max={60}
                value={form.timeout_seconds}
                onChange={(v) => set("timeout_seconds", v)}
              />
              <NumberWithUnit
                id={ids.threshold}
                label="Falhas até cair"
                unit="seguidas"
                min={1}
                max={10}
                value={form.failure_threshold}
                onChange={(v) => set("failure_threshold", v)}
              />
            </div>
            {threshold > 0 && (
              <p className="mt-3 text-pg-caption text-pg-neblina">
                Fica fora do ar depois de{" "}
                {plural(threshold, "falha", "falhas seguidas")}, ou seja, em
                cerca de {downAfter}.
              </p>
            )}
          </fieldset>
        </>
      )}

      {error && (
        <p role="alert" className="text-pg-ui text-pg-down mb-4">
          {error}
        </p>
      )}

      <div className="flex items-center gap-2.5 pt-5 -mx-6 sm:-mx-7 px-6 sm:px-7 border-t border-pg-borda">
        {editing && onDelete && (
          <button
            type="button"
            onClick={onDelete}
            className="inline-flex items-center gap-1.5 px-2 py-2 -ml-2 rounded-pg-control text-pg-ui font-medium text-pg-down hover:bg-pg-down-bg focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
          >
            <Trash2 size={16} aria-hidden="true" />
            Excluir
          </button>
        )}
        <div className="flex gap-2.5 ml-auto">
          <Button variant="ghost" onClick={onCancel} type="button">
            Cancelar
          </Button>
          <Button type="submit" disabled={loading}>
            {loading ? "Salvando..." : editing ? "Salvar" : "Criar monitor"}
          </Button>
        </div>
      </div>
    </form>
  );
}
