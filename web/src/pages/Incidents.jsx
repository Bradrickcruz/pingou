import { useState, useEffect } from "react";
import { incidentsApi } from "../api/incidents";
import { Spinner } from "../components/ui/Spinner";
import { tokens as t } from "../theme/tokens";

export function Incidents() {
  const [incidents, setIncidents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [onlyOpen, setOnlyOpen] = useState(false);

  const handleToggleOnlyOpen = (event) => {
    setLoading(true);
    setOnlyOpen(event.target.checked);
  };

  useEffect(() => {
    incidentsApi
      .list({ open: onlyOpen, limit: 50 })
      .then((r) => setIncidents(r.data))
      .finally(() => setLoading(false));
  }, [onlyOpen]);

  return (
    <div>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: "24px",
        }}
      >
        <h1 style={{ fontSize: "20px", fontWeight: 700 }}>Incidents</h1>
        <label
          style={{
            display: "flex",
            alignItems: "center",
            gap: "8px",
            color: t.colors.textMuted,
            fontSize: "13px",
          }}
        >
          <input
            type="checkbox"
            checked={onlyOpen}
            onChange={handleToggleOnlyOpen}
            style={{ width: "auto" }}
          />
          Open only
        </label>
      </div>

      {loading && (
        <div
          style={{ display: "flex", justifyContent: "center", padding: "48px" }}
        >
          <Spinner />
        </div>
      )}

      {!loading && incidents.length === 0 && (
        <p
          style={{
            color: t.colors.textMuted,
            textAlign: "center",
            padding: "48px",
          }}
        >
          No incidents found 🎉
        </p>
      )}

      <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
        {incidents.map((i) => (
          <div
            key={i.id}
            style={{
              background: t.colors.surface,
              border: `1px solid ${i.open ? t.colors.danger : t.colors.border}`,
              borderRadius: t.radius.md,
              padding: "14px 18px",
            }}
          >
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
              }}
            >
              <span style={{ fontWeight: 600, fontSize: "13px" }}>
                {i.open ? (
                  <span style={{ color: t.colors.danger }}>● OPEN</span>
                ) : (
                  <span style={{ color: t.colors.success }}>✓ RESOLVED</span>
                )}
                <span
                  style={{
                    color: t.colors.textMuted,
                    fontWeight: 400,
                    marginLeft: "10px",
                  }}
                >
                  {i.monitor_id}
                </span>
              </span>
              <span style={{ fontSize: "11px", color: t.colors.textMuted }}>
                {new Date(i.started_at).toLocaleString()}
              </span>
            </div>
            {i.last_error && (
              <p
                style={{
                  color: t.colors.textMuted,
                  fontSize: "12px",
                  marginTop: "6px",
                  fontFamily: t.font.mono,
                }}
              >
                {i.last_error}
              </p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
