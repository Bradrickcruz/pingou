import { Check, Send, X } from "lucide-react";

// Botão "Enviar teste", igual nos cards de notificação.
export function TestButton({ onClick, disabled, testing }) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled || testing}
      className="inline-flex items-center justify-center gap-2 min-h-[44px] md:min-h-0 px-4 py-2 rounded-pg-control border border-pg-linha bg-pg-superficie text-pg-ui font-medium text-pg-texto hover:bg-pg-hover disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
    >
      {testing ? (
        <span
          className="w-4 h-4 border-2 border-pg-linha border-t-pg-link rounded-full animate-spin"
          aria-hidden="true"
        />
      ) : (
        <Send size={16} aria-hidden="true" />
      )}
      {testing ? "Enviando…" : "Enviar teste"}
    </button>
  );
}

// Resultado do teste ({ok, text}), lido por leitores de tela.
export function TestResult({ id, result }) {
  return (
    <p
      id={id}
      role="status"
      aria-live="polite"
      className="mt-2 mb-0 text-pg-caption md:text-pg-ui empty:hidden"
    >
      {result && (
        <span
          className={`flex md:inline-flex items-start md:items-center gap-2 md:gap-1 p-3 md:p-0 rounded-pg-control md:bg-transparent ${
            result.ok
              ? "text-pg-status-up bg-pg-status-up-bg"
              : "text-pg-status-down bg-pg-status-down-bg"
          }`}
        >
          {result.ok ? (
            <Check
              size={14}
              aria-hidden="true"
              className="shrink-0 mt-1 md:mt-0"
            />
          ) : (
            <X size={14} aria-hidden="true" className="shrink-0 mt-1 md:mt-0" />
          )}
          <span className="md:text-inherit text-pg-texto">{result.text}</span>
        </span>
      )}
    </p>
  );
}
