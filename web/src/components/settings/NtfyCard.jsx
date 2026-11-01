import { useState } from "react";
import { settingsApi } from "../../api/settings";
import { Card } from "./Card";
import { TestButton, TestResult } from "./TestControls";
import { testErrorText } from "../../lib/testError";

const NTFY_DOCS = "https://github.com/Bradrickcruz/pingou#ntfy";

// o mesmo que o ntfy aceita como nome de tópico
const TOPIC_PATTERN = /^[A-Za-z0-9_-]{1,64}$/;

// "Teste enviado em 120 ms." / "Teste falhou. unexpected status code: 403 (...)"
function resultText(r) {
  if (r.ok) return `Teste enviado em ${r.duration_ms} ms.`;
  return ["Teste falhou.", r.error].filter(Boolean).join(" ");
}

// Só o tópico fica na tela; servidor e token ficam no ambiente
// (PINGOU_NTFY_URL e PINGOU_NTFY_TOKEN). Sem servidor (available=false), o
// card fica bloqueado. O teste usa o tópico digitado, sem salvar.
export function NtfyCard({ value, onChange, available }) {
  const [testing, setTesting] = useState(false);
  const [result, setResult] = useState(null); // {ok, text}

  const topic = value.trim();
  const invalid = topic !== "" && !TOPIC_PATTERN.test(topic);

  const change = (v) => {
    setResult(null);
    onChange(v);
  };

  const test = async () => {
    setTesting(true);
    setResult(null);
    try {
      const r = await settingsApi.testNtfy(topic);
      setResult({ ok: r.ok, text: resultText(r) });
    } catch (e) {
      setResult({ ok: false, text: testErrorText(e) });
    } finally {
      setTesting(false);
    }
  };

  return (
    <Card
      title="ntfy"
      description="O Pingou publica uma notificação no seu tópico do ntfy quando um monitor cai e quando ele volta."
    >
      <label
        htmlFor="ntfy-topic"
        className="block mb-1.5 text-pg-ui font-medium text-pg-ardosia"
      >
        Tópico
      </label>
      <div className="flex flex-col md:flex-row gap-2">
        <input
          id="ntfy-topic"
          type="text"
          autoComplete="off"
          spellCheck={false}
          value={value}
          onChange={(e) => change(e.target.value)}
          disabled={!available}
          placeholder="pingou-alertas"
          aria-invalid={invalid || undefined}
          aria-describedby="ntfy-help ntfy-result"
          className="flex-1 min-w-0 px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda aria-[invalid=true]:border-pg-down disabled:bg-pg-nuvem disabled:text-pg-neblina disabled:cursor-not-allowed"
        />
        {available && (
          <TestButton
            onClick={test}
            disabled={!topic || invalid}
            testing={testing}
          />
        )}
      </div>
      {available ? (
        <p
          id="ntfy-help"
          className={`mt-1.5 mb-0 text-pg-caption ${invalid ? "text-pg-down" : "text-pg-neblina"}`}
        >
          {invalid
            ? "Use de 1 a 64 caracteres: letras sem acento, números, _ ou -."
            : "Vazio desliga o ntfy. Em servidor aberto, quem souber o tópico lê os alertas: escolha um nome difícil de adivinhar."}
        </p>
      ) : (
        <p
          id="ntfy-help"
          className="mt-1.5 mb-0 text-pg-caption text-pg-neblina"
        >
          Para ativar, configure{" "}
          <code className="font-pg-mono">PINGOU_NTFY_URL</code> (e{" "}
          <code className="font-pg-mono">PINGOU_NTFY_TOKEN</code>, se o servidor
          exigir) no servidor.{" "}
          <a
            href={NTFY_DOCS}
            target="_blank"
            rel="noreferrer"
            className="text-pg-ciano-fundo underline"
          >
            Como configurar
          </a>
        </p>
      )}
      <TestResult id="ntfy-result" result={result} />
    </Card>
  );
}
