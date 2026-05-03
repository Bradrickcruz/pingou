import { useState } from "react";
import { tokens as t } from "../theme/tokens";
import { Button } from "../components/ui/Button";

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
        setError("Invalid API key.");
        return;
      }

      if (!res.ok) {
        setError("Could not connect to the API.");
        return;
      }

      localStorage.setItem("pingou_api_key", key);
      onLogin(key);
    } catch {
      setError("Could not connect to the API.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen flex items-center justify-center"
      style={{
        background: t.colors.bg,
      }}
    >
      <div
        className="p-10 rounded-lg border w-full max-w-[380px]"
        style={{
          background: t.colors.surface,
          borderColor: t.colors.border,
          borderRadius: t.radius.lg,
        }}
      >
        <div className="text-center mb-8">
          <div className="text-[40px] mb-2">🏓</div>
          <h1
            className="text-[22px] font-bold"
            style={{
              color: t.colors.primary,
            }}
          >
            Pingou
          </h1>
          <p
            className="text-sm mt-1"
            style={{
              color: t.colors.textMuted,
            }}
          >
            health checker
          </p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label
              className="block mb-1.5 text-xs font-medium"
              style={{
                color: t.colors.textMuted,
              }}
            >
              API Key
            </label>
            <input
              type="password"
              value={key}
              onChange={(e) => setKey(e.target.value)}
              placeholder="Enter your API key"
              required
              autoFocus
              className="w-full px-3 py-2 rounded text-sm bg-[var(--bg)] border border-[var(--border)] text-[var(--text-h)] focus:outline-none focus:border-[var(--accent)]"
            />
          </div>

          {error && (
            <p
              className="text-sm mb-3.5"
              style={{
                color: t.colors.danger,
              }}
            >
              {error}
            </p>
          )}

          <Button
            type="submit"
            disabled={loading || !key}
            className="w-full justify-center"
          >
            {loading ? "Verifying..." : "Enter"}
          </Button>
        </form>
      </div>
    </div>
  );
}