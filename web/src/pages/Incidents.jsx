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
      <div className="flex flex-col gap-3 sm:flex-row sm:justify-between sm:items-center mb-6">
        <h1 className="text-pg-h1 font-medium text-pg-ardosia my-0">Incidentes</h1>
        <label className="inline-flex items-center gap-2 text-pg-ui text-pg-neblina">
          <input
            type="checkbox"
            checked={onlyOpen}
            onChange={handleToggleOnlyOpen}
            className="w-auto"
          />
          Somente abertos
        </label>
      </div>

      {loading && (
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      )}

      {!loading && incidents.length === 0 && (
        <p className="text-center py-12 text-pg-neblina">Nenhum incidente encontrado</p>
      )}

      <div className="flex flex-col gap-2">
        {incidents.map((i) => (
          <div
            key={i.id}
            className={`p-3.5 rounded-pg-card border bg-pg-branco ${
              i.open ? "border-pg-down" : "border-pg-borda"
            }`}
          >
            <div className="flex flex-wrap justify-between items-center gap-x-3 gap-y-1">
              <span className="font-medium text-pg-ui">
                {i.open ? (
                  <span className="inline-flex items-center gap-1 text-pg-down">
                    <CircleAlert size={14} aria-hidden="true" />
                    ABERTO
                  </span>
                ) : (
                  <span className="inline-flex items-center gap-1 text-pg-up">
                    <Check size={14} aria-hidden="true" />
                    RESOLVIDO
                  </span>
                )}
                <span className="font-normal ml-2.5 text-pg-neblina">{i.monitor_id}</span>
              </span>
              <span className="text-pg-caption text-pg-neblina">
                {new Date(i.started_at).toLocaleString("pt-BR")}
              </span>
            </div>
            {i.last_error && (
              <p className="text-pg-caption mt-1.5 font-mono text-pg-neblina break-words">
                {i.last_error}
              </p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
