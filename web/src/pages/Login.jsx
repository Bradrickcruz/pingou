import { useState } from "react";
import { Button } from "../components/ui/Button";
import marca from "../assets/brand/pingou-marca-horizontal.svg";

export function Login({ onLogin }) {
  const [key, setKey] = useState("");
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await fetch("/api/monitors?limit=1", {
        headers: { "X-API-Key": key },
      });

      if (res.status === 401) {
        setError("Chave de API inválida.");
        return;
      }

      if (!res.ok) {
        setError("Não foi possível conectar à API.");
        return;
      }

      localStorage.setItem("pingou_api_key", key);
      onLogin(key);
    } catch {
      setError("Não foi possível conectar à API.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-pg-nuvem">
      <div className="p-10 rounded-pg-card border border-pg-borda bg-pg-branco w-full max-w-[380px]">
        <div className="text-center mb-8">
          <h1 className="my-0">
            <img src={marca} alt="Pingou" className="h-10 w-auto mx-auto" />
          </h1>
          <p className="text-pg-caption text-pg-neblina mt-2">monitor de saúde</p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label
              htmlFor="api-key"
              className="block mb-1.5 text-pg-caption font-medium text-pg-neblina"
            >
              Chave de API
            </label>
            <input
              id="api-key"
              type="password"
              value={key}
              onChange={(e) => setKey(e.target.value)}
              placeholder="Digite sua chave de API"
              required
              autoFocus
              className="w-full px-3 py-2 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda"
            />
          </div>

          {error && <p role="alert" className="text-pg-ui text-pg-down mb-3.5">{error}</p>}

          <Button type="submit" disabled={loading || !key} className="w-full">
            {loading ? "Verificando..." : "Entrar"}
          </Button>
        </form>
      </div>
    </div>
  );
}
