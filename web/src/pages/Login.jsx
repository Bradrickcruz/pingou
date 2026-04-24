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
      // valida a key fazendo uma request real pra API
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

      // persiste no sessionStorage e notifica o App
      sessionStorage.setItem("pingou_api_key", key);
      onLogin(key);
    } catch {
      setError("Could not connect to the API.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: "100vh",
        background: t.colors.bg,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <div
        style={{
          background: t.colors.surface,
          border: `1px solid ${t.colors.border}`,
          borderRadius: t.radius.lg,
          padding: "40px",
          width: "100%",
          maxWidth: "380px",
        }}
      >
        <div style={{ textAlign: "center", marginBottom: "32px" }}>
          <div style={{ fontSize: "40px", marginBottom: "8px" }}>🏓</div>
          <h1
            style={{
              fontSize: "22px",
              fontWeight: 700,
              color: t.colors.primary,
            }}
          >
            Pingou
          </h1>
          <p
            style={{
              color: t.colors.textMuted,
              fontSize: "13px",
              marginTop: "4px",
            }}
          >
            health checker
          </p>
        </div>

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: "16px" }}>
            <label
              style={{
                display: "block",
                marginBottom: "6px",
                color: t.colors.textMuted,
                fontSize: "12px",
                fontWeight: 500,
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
            />
          </div>

          {error && (
            <p
              style={{
                color: t.colors.danger,
                fontSize: "13px",
                marginBottom: "14px",
              }}
            >
              {error}
            </p>
          )}

          <Button
            type="submit"
            disabled={loading || !key}
            style={{ width: "100%", justifyContent: "center" }}
          >
            {loading ? "Verifying..." : "Enter"}
          </Button>
        </form>
      </div>
    </div>
  );
}
