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
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-xl font-bold">Incidents</h1>
        <label
          className="flex items-center gap-2 text-sm"
          style={{
            color: t.colors.textMuted,
          }}
        >
          <input
            type="checkbox"
            checked={onlyOpen}
            onChange={handleToggleOnlyOpen}
            className="w-auto"
          />
          Open only
        </label>
      </div>

      {loading && (
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      )}

      {!loading && incidents.length === 0 && (
        <p
          className="text-center py-12"
          style={{
            color: t.colors.textMuted,
          }}
        >
          No incidents found 🎉
        </p>
      )}

      <div className="flex flex-col gap-2">
        {incidents.map((i) => (
          <div
            key={i.id}
            className="p-3.5 rounded-md border"
            style={{
              background: t.colors.surface,
              borderColor: i.open ? t.colors.danger : t.colors.border,
              borderRadius: t.radius.md,
            }}
          >
            <div className="flex justify-between items-center">
              <span className="font-semibold text-sm">
                {i.open ? (
                  <span
                    style={{
                      color: t.colors.danger,
                    }}
                  >
                    ● OPEN
                  </span>
                ) : (
                  <span
                    style={{
                      color: t.colors.success,
                    }}
                  >
                    ✓ RESOLVED
                  </span>
                )}
                <span
                  className="font-normal ml-2.5"
                  style={{
                    color: t.colors.textMuted,
                  }}
                >
                  {i.monitor_id}
                </span>
              </span>
              <span
                className="text-[11px]"
                style={{
                  color: t.colors.textMuted,
                }}
              >
                {new Date(i.started_at).toLocaleString()}
              </span>
            </div>
            {i.last_error && (
              <p
                className="text-xs mt-1.5 font-mono"
                style={{
                  color: t.colors.textMuted,
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