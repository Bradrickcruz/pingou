import { useState } from "react";
import { settingsApi } from "../../api/settings";
import { Card } from "./Card";
import { TestButton, TestResult } from "./TestControls";
import { testErrorText } from "../../lib/testError";

const JSON_DOCS = "https://github.com/Bradrickcruz/pingou#webhook";

// "Teste enviado. Resposta 200 OK em 142 ms." / "Teste falhou. Resposta 500."
function resultText(r) {
  const status = r.status_code === 200 ? "200 OK" : r.status_code;
  if (r.ok) return `Teste enviado. Resposta ${status} em ${r.duration_ms} ms.`;
  const parts = ["Teste falhou."];
  if (r.status_code != null) parts.push(`Resposta ${status}.`);
  if (r.error) parts.push(r.error);
  return parts.join(" ");
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
      setResult({ ok: false, text: testErrorText(e) });
    } finally {
      setTesting(false);
    }
  };

  return (
    <Card
      title="Webhook"
      description={
        <>
          O Pingou envia um POST com JSON para esta URL quando um monitor cai e
          quando ele volta. Funciona com Slack, Discord, n8n ou qualquer serviço
          que receba webhooks.{" "}
          <a
            href={JSON_DOCS}
            target="_blank"
            rel="noreferrer"
            className="text-pg-link underline"
          >
            Ver formato do JSON
          </a>
        </>
      }
    >
      <label
        htmlFor="webhook-url"
        className="block mb-1.5 text-pg-ui font-medium text-pg-texto"
      >
        URL do webhook
      </label>
      <div className="flex flex-col md:flex-row gap-2">
        <input
          id="webhook-url"
          type="url"
          value={value}
          onChange={(e) => change(e.target.value)}
          placeholder="https://hooks.exemplo.com/..."
          aria-describedby="webhook-help webhook-result"
          className="flex-1 min-w-0 px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-ui text-pg-texto bg-pg-superficie border border-pg-linha"
        />
        <TestButton onClick={test} disabled={!value} testing={testing} />
      </div>
      <p
        id="webhook-help"
        className="mt-1.5 mb-0 text-pg-caption text-pg-texto-2"
      >
        Deixe vazio para não enviar notificações.
      </p>
      <TestResult id="webhook-result" result={result} />
    </Card>
  );
}
