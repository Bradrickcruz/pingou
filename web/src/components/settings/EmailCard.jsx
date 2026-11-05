import { useState } from "react";
import { settingsApi } from "../../api/settings";
import { Card } from "./Card";
import { TestButton, TestResult } from "./TestControls";
import { testErrorText } from "../../lib/testError";

const EMAIL_DOCS = "https://github.com/Bradrickcruz/pingou#e-mail";

// "Teste enviado em 812 ms." / "Teste falhou. timeout após 30s"
function resultText(r) {
  if (r.ok) return `Teste enviado em ${r.duration_ms} ms.`;
  return ["Teste falhou.", r.error].filter(Boolean).join(" ");
}

// Os destinatários ficam na tela; o servidor SMTP fica todo no ambiente
// (PINGOU_SMTP_*). Sem ele (available=false), o card fica bloqueado.
// O teste usa os destinatários digitados, sem salvar.
export function EmailCard({ value, onChange, available }) {
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
      const r = await settingsApi.testEmail(value);
      setResult({ ok: r.ok, text: resultText(r) });
    } catch (e) {
      setResult({ ok: false, text: testErrorText(e) });
    } finally {
      setTesting(false);
    }
  };

  return (
    <Card
      title="E-mail"
      description="O Pingou envia um e-mail em texto simples quando um monitor cai e quando ele volta."
    >
      <label
        htmlFor="email-recipients"
        className="block mb-1.5 text-pg-ui font-medium text-pg-texto"
      >
        Destinatários
      </label>
      <div className="flex flex-col md:flex-row gap-2">
        <input
          id="email-recipients"
          type="text"
          inputMode="email"
          autoComplete="off"
          value={value}
          onChange={(e) => change(e.target.value)}
          disabled={!available}
          placeholder="ops@exemplo.com, voce@exemplo.com"
          aria-describedby="email-help email-result"
          className="flex-1 min-w-0 px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-ui text-pg-texto bg-pg-superficie border border-pg-linha disabled:bg-pg-fundo disabled:text-pg-texto-2 disabled:cursor-not-allowed"
        />
        {available && (
          <TestButton
            onClick={test}
            disabled={!value.trim()}
            testing={testing}
          />
        )}
      </div>
      {available ? (
        <p
          id="email-help"
          className="mt-1.5 mb-0 text-pg-caption text-pg-texto-2"
        >
          Separe por vírgula (máximo de 10 destinatários, limite do SMTP). Vazio
          desliga o e-mail.
        </p>
      ) : (
        <p
          id="email-help"
          className="mt-1.5 mb-0 text-pg-caption text-pg-texto-2"
        >
          Para ativar, configure{" "}
          <code className="font-pg-mono">PINGOU_SMTP_HOST</code> e{" "}
          <code className="font-pg-mono">PINGOU_SMTP_FROM</code> no servidor.{" "}
          <a
            href={EMAIL_DOCS}
            target="_blank"
            rel="noreferrer"
            className="text-pg-link underline"
          >
            Como configurar
          </a>
        </p>
      )}
      <TestResult id="email-result" result={result} />
    </Card>
  );
}
