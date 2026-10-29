import { useState } from "react";
import { Check, Send, X } from "lucide-react";
import { settingsApi } from "../../api/settings";
import { Card } from "./Card";

const JSON_DOCS = "https://github.com/Bradrickcruz/pingou#webhook";

function resultText(r) {
  const parts = [r.ok ? "Entregue" : "Falhou"];
  if (r.status_code != null) parts.push(r.status_code);
  if (r.ok) parts.push(`${r.duration_ms} ms`);
  else if (r.error) parts.push(r.error);
  return parts.join(" · ");
}

// O teste usa a URL digitada, sem salvar.
export function WebhookCard({ value, onChange }) {
  const [testing, setTesting] = useState(false);
  const [result, setResult] = useState(null); // {ok, text}

  const change = (v) => {
    setResult(null);
    onChange(v);
  };

  const test = async () => {
    setTesting(true);
    setResult(null);
    try {
      const r = await settingsApi.testWebhook(value);
      setResult({ ok: r.ok, text: resultText(r) });
    } catch (e) {
      const text = e.status === 429 ? "Muitos testes seguidos. Tente de novo em 1 minuto." : e.message;
      setResult({ ok: false, text });
    } finally {
      setTesting(false);
    }
  };

  return (
    <Card
      title="Notificações"
      description={
        <>
          O Pingou envia um POST com JSON para esta URL quando um monitor cai e quando ele volta. Funciona com Slack,
          Discord, n8n ou qualquer serviço que receba webhooks.{" "}
          <a href={JSON_DOCS} target="_blank" rel="noreferrer" className="text-pg-ciano-fundo underline">
            Ver formato do JSON
          </a>
        </>
      }
    >
      <label htmlFor="webhook-url" className="block mb-1.5 text-pg-ui font-medium text-pg-ardosia">
        URL do webhook
      </label>
      <div className="flex flex-col sm:flex-row gap-2">
        <input
          id="webhook-url"
          type="url"
          value={value}
          onChange={(e) => change(e.target.value)}
          placeholder="https://hooks.exemplo.com/..."
          aria-describedby="webhook-help"
          className="flex-1 min-w-0 px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda"
        />
        <button
          type="button"
          onClick={test}
          disabled={!value || testing}
          className="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-pg-control border border-pg-borda bg-pg-branco text-pg-ui font-medium text-pg-ardosia hover:bg-pg-nuvem disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
        >
          {testing ? (
            <span className="w-4 h-4 border-2 border-pg-borda border-t-pg-ciano-fundo rounded-full animate-spin" aria-hidden="true" />
          ) : (
            <Send size={16} aria-hidden="true" />
          )}
          {testing ? "Enviando…" : "Enviar teste"}
        </button>
      </div>
      <p id="webhook-help" className="mt-1.5 mb-0 text-pg-caption text-pg-neblina">
        Deixe vazio para não enviar notificações.
      </p>
      <p role="status" aria-live="polite" className="mt-2 mb-0 text-pg-ui empty:hidden">
        {result && (
          <span className={`inline-flex items-center gap-1 ${result.ok ? "text-pg-up" : "text-pg-down"}`}>
            {result.ok ? <Check size={14} aria-hidden="true" /> : <X size={14} aria-hidden="true" />}
            {result.text}
          </span>
        )}
      </p>
    </Card>
  );
}
