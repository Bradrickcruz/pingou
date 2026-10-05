import { useState, useEffect } from "react";
import { CircleAlert, Check } from "lucide-react";
import { incidentsApi } from "../api/incidents";
import { Spinner } from "../components/ui/Spinner";

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
        <h1 className="text-pg-h1 font-medium text-pg-ardosia my-0">Incidents</h1>
        <label className="flex items-center gap-2 text-pg-ui text-pg-neblina">
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
        <p className="text-center py-12 text-pg-neblina">No incidents found</p>
      )}

      <div className="flex flex-col gap-2">
        {incidents.map((i) => (
          <div
            key={i.id}
            className={`p-3.5 rounded-pg-card border bg-pg-branco ${
              i.open ? "border-pg-down" : "border-pg-borda"
            }`}
          >
            <div className="flex justify-between items-center">
              <span className="font-medium text-pg-ui">
                {i.open ? (
                  <span className="inline-flex items-center gap-1 text-pg-down">
                    <CircleAlert size={14} aria-hidden="true" />
                    OPEN
                  </span>
                ) : (
                  <span className="inline-flex items-center gap-1 text-pg-up">
                    <Check size={14} aria-hidden="true" />
                    RESOLVED
                  </span>
                )}
                <span className="font-normal ml-2.5 text-pg-neblina">
                  {i.monitor_id}
                </span>
              </span>
              <span className="text-pg-caption text-pg-neblina">
                {new Date(i.started_at).toLocaleString()}
              </span>
            </div>
            {i.last_error && (
              <p className="text-pg-caption mt-1.5 font-mono text-pg-neblina">
                {i.last_error}
              </p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
