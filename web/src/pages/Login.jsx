import { useState } from "react";
import { CircleX, Eye, EyeOff } from "lucide-react";
import { Button } from "../components/ui/Button";
import marca from "../assets/brand/pingou-marca-horizontal.svg";
import marcaNegativa from "../assets/brand/pingou-marca-horizontal-negativa.svg";
import { version } from "../../package.json";

const GITHUB_URL = "https://github.com/Bradrickcruz/pingou";

function Footer({ className = "", linkClassName = "" }) {
  return (
    <p className={`text-pg-caption ${className}`}>
      Versão {version} ·{" "}
      <a
        href={GITHUB_URL}
        target="_blank"
        rel="noopener noreferrer"
        className={`underline ${linkClassName}`}
      >
        Código no GitHub
      </a>
    </p>
  );
}

// Ondas (DESIGN.md §8): nascem do canto, 4 anéis, opacidade crescente.
function Ondas() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 250 250"
      className="pointer-events-none absolute bottom-0 left-0 h-[340px] w-[340px] text-pg-ciano"
      fill="none"
      stroke="currentColor"
      strokeWidth="6"
    >
      {[
        [248, 0.15],
        [192, 0.25],
        [136, 0.45],
        [80, 0.8],
      ].map(([r, o]) => (
        <circle key={r} cx="0" cy="250" r={r - 3} opacity={o} />
      ))}
    </svg>
  );
}

export function Login({ onLogin }) {
  const [key, setKey] = useState("");
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [visible, setVisible] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await fetch("/api/info", {
        headers: { "X-API-Key": key },
      });

      if (res.status === 401) {
        setError(
          "API Key inválida. Confira o valor de PINGOU_API_KEY no servidor.",
        );
        return;
      }

      if (!res.ok) {
        setError(
          "Não foi possível conectar à API. Verifique a URL e tente de novo.",
        );
        return;
      }

      localStorage.setItem("pingou_api_key", key);
      onLogin(key);
    } catch {
      setError(
        "Não foi possível conectar à API. Verifique a URL e tente de novo.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-pg-nuvem md:grid md:grid-cols-[45%_1fr]">
      <aside className="relative hidden flex-col justify-between overflow-hidden bg-pg-ardosia p-14 text-pg-branco md:flex">
        <img
          src={marcaNegativa}
          alt="pingou"
          className="h-9 w-auto self-start"
        />
        <div>
          <p className="font-pg-heading text-pg-h1 font-medium">
            Rodou, pingou.
          </p>
          <p className="mt-2 text-pg-body text-pg-nuvem">
            Monitor de uptime leve e self-hosted.
          </p>
        </div>
        <Footer className="self-end" linkClassName="text-pg-branco" />
        <Ondas />
      </aside>

      <main className="flex min-h-screen flex-col px-6 pb-6 pt-8 md:items-center md:justify-center md:p-0">
        <img
          src={marca}
          alt="pingou"
          className="h-9 w-auto self-start md:hidden"
        />

        <div className="mt-12 w-full max-w-[360px] flex-1 md:mt-0 md:flex-none">
          <h1 className="font-pg-heading text-pg-h2 font-medium text-pg-ardosia">
            Entrar no painel
          </h1>

          <form onSubmit={handleSubmit} aria-busy={loading} className="mt-6">
            <label
              htmlFor="api-key"
              className="mb-1.5 block text-pg-ui font-medium text-pg-ardosia"
            >
              API Key
            </label>
            <div className="relative">
              <input
                id="api-key"
                type={visible ? "text" : "password"}
                value={key}
                onChange={(e) => {
                  setKey(e.target.value);
                  setError(null);
                }}
                placeholder="Cole sua API Key"
                autoComplete="current-password"
                autoFocus
                required
                readOnly={loading}
                aria-invalid={error ? "true" : undefined}
                aria-describedby={error ? "api-key-error" : undefined}
                className={`w-full rounded-pg-control py-2.5 pl-3 pr-11 text-pg-ui text-pg-ardosia ${
                  error ? "border-pg-down" : "border-pg-borda"
                }`}
              />
              <button
                type="button"
                onClick={() => setVisible((v) => !v)}
                aria-label={visible ? "Ocultar API Key" : "Mostrar API Key"}
                aria-pressed={visible}
                className="absolute right-1 top-1/2 flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-pg-control text-pg-neblina hover:bg-pg-nuvem"
              >
                {visible ? (
                  <EyeOff size={18} aria-hidden="true" />
                ) : (
                  <Eye size={18} aria-hidden="true" />
                )}
              </button>
            </div>

            {error && (
              <p
                id="api-key-error"
                role="alert"
                className="mt-2 flex items-start gap-2 text-pg-ui text-pg-down"
              >
                <CircleX
                  size={16}
                  aria-hidden="true"
                  className="mt-0.5 shrink-0"
                />
                {error}
              </p>
            )}

            <Button
              type="submit"
              disabled={loading || !key}
              className="mt-4 w-full"
            >
              {loading ? "Verificando..." : "Entrar"}
            </Button>
          </form>

          <p className="mt-6 text-pg-caption text-pg-neblina">
            A chave fica salva neste navegador até você sair.
          </p>
        </div>

        <Footer
          className="mt-8 text-pg-neblina md:hidden"
          linkClassName="text-pg-ciano-fundo"
        />
      </main>
    </div>
  );
}
