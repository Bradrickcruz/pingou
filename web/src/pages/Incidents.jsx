import { useCallback, useEffect, useRef, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { incidentsApi } from "../api/incidents";
import { monitorsApi } from "../api/monitors";
import { useNow } from "../hooks/useNow";
import { useSettings } from "../hooks/useSettings";
import { OpenIncidentCard } from "../components/incidents/OpenIncidentCard";
import { ResolvedTable } from "../components/incidents/ResolvedTable";
import { Button } from "../components/ui/Button";
import { SegmentedTabs } from "../components/ui/SegmentedTabs";
import { MobileHeader } from "../components/layout/MobileHeader";
import { Spinner } from "../components/ui/Spinner";

const PAGE_SIZE = 20;

// Abertos: uma busca só (até 100), recarregada a cada 15 s como o painel.
function useOpenIncidents(monitorId) {
  const [state, setState] = useState({ items: [], total: 0, loading: true, error: null });
  const [attempt, setAttempt] = useState(0);

  useEffect(() => {
    let alive = true; // descarta resposta atrasada de um filtro anterior
    const load = () =>
      incidentsApi
        .list({ status: "open", monitor_id: monitorId || undefined, limit: 100 })
        .then((res) => alive && setState({ items: res.data, total: res.total, loading: false, error: null }))
        .catch((e) => alive && setState((s) => ({ ...s, loading: false, error: e.message })));
    const first = setTimeout(() => {
      setState((s) => ({ ...s, loading: true, error: null }));
      void load();
    }, 0);
    const id = setInterval(load, 15_000);
    return () => {
      alive = false;
      clearTimeout(first);
      clearInterval(id);
    };
  }, [monitorId, attempt]);

  return { ...state, retry: () => setAttempt((n) => n + 1) };
}

// Resolvidos: paginados com "Carregar mais"; não fazem polling.
function useResolvedIncidents(monitorId) {
  const [state, setState] = useState({ items: [], total: 0, loading: true, loadingMore: false, error: null });
  const generation = useRef(0); // muda a cada filtro; respostas de outra geração são ignoradas

  const load = useCallback(
    async (offset) => {
      const gen = generation.current;
      setState((s) => (offset === 0 ? { ...s, loading: true, error: null } : { ...s, loadingMore: true, error: null }));
      try {
        const res = await incidentsApi.list({
          status: "resolved",
          monitor_id: monitorId || undefined,
          limit: PAGE_SIZE,
          offset,
        });
        if (gen !== generation.current) return;
        setState((s) => {
          // um incidente resolvido entre as páginas desloca o offset: evita repetir
          const seen = new Set(offset === 0 ? [] : s.items.map((i) => i.id));
          const items = [...(offset === 0 ? [] : s.items), ...res.data.filter((i) => !seen.has(i.id))];
          return { items, total: res.total, loading: false, loadingMore: false, error: null };
        });
      } catch (e) {
        if (gen !== generation.current) return;
        setState((s) => ({ ...s, loading: false, loadingMore: false, error: e.message }));
      }
    },
    [monitorId],
  );

  useEffect(() => {
    generation.current++;
    const first = setTimeout(() => void load(0), 0);
    return () => clearTimeout(first);
  }, [load]);

  return {
    ...state,
    loadMore: () => load(state.items.length),
    retry: () => load(0),
  };
}

const TABS = [
  { key: "all", label: "Todos" },
  { key: "open", label: "Abertos" },
  { key: "resolved", label: "Resolvidos" },
];

function LoadError({ message, onRetry }) {
  return (
    <div role="alert" className="flex flex-wrap items-center gap-3 p-4 rounded-pg-card border border-pg-borda bg-pg-branco">
      <span className="text-pg-ui text-pg-down">Não foi possível carregar: {message}</span>
      <Button variant="ghost" onClick={onRetry}>
        Tentar de novo
      </Button>
    </div>
  );
}

function Empty({ children }) {
  return (
    <p className="m-0 py-8 px-4 text-center text-pg-ui text-pg-neblina rounded-pg-card border border-dashed border-pg-borda">
      {children}
    </p>
  );
}

export function Incidents() {
  const [params, setParams] = useSearchParams();
  const tab = TABS.some((t) => t.key === params.get("tab")) ? params.get("tab") : "all";
  const monitorId = params.get("monitor") ?? "";

  const now = useNow(1000);
  const { settings } = useSettings();
  const [monitors, setMonitors] = useState([]);
  const open = useOpenIncidents(monitorId);
  const resolved = useResolvedIncidents(monitorId);

  useEffect(() => {
    monitorsApi
      .list({ limit: 100 })
      .then((res) => setMonitors([...res.data].sort((a, b) => a.name.localeCompare(b.name, "pt-BR"))))
      .catch(() => {}); // sem a lista, o filtro só mostra "Todos os monitores"
  }, []);

  const setParam = (key, value) =>
    setParams(
      (p) => {
        const next = new URLSearchParams(p);
        if (value) next.set(key, value);
        else next.delete(key);
        return next;
      },
      { replace: true },
    );

  const days = settings?.retention_days;
  const selected = monitors.find((m) => m.id === monitorId);
  const showOpen = tab !== "resolved";
  const showResolved = tab !== "open";

  return (
    <div className="max-w-[1180px]">
      <MobileHeader title="Incidentes" />
      <h1 className="hidden md:block text-pg-h1 font-medium text-pg-ardosia my-0">Incidentes</h1>
      <p className="hidden md:block mt-2 mb-6 text-pg-ui text-pg-neblina">
        Cada vez que um monitor caiu e quando voltou.
        {days && (
          <>
            {" "}Guardados por {days} dias (
            <Link to="/settings" className="text-pg-ciano-fundo underline">
              mudar em Configurações
            </Link>
            ).
          </>
        )}
      </p>

      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between mb-6 md:mb-7">
        <SegmentedTabs
          label="Filtrar incidentes"
          options={TABS.map((t) => ({
            ...t,
            count: t.key === "open" && !open.loading && !open.error ? open.total : null,
          }))}
          value={tab}
          onChange={(key) => setParam("tab", key === "all" ? "" : key)}
        />
        <div className="flex flex-col gap-1.5 md:flex-row md:items-center md:gap-3">
          <label htmlFor="incident-monitor" className="text-pg-caption md:text-pg-ui text-pg-neblina">
            Monitor
          </label>
          <select
            id="incident-monitor"
            value={monitorId}
            onChange={(e) => setParam("monitor", e.target.value)}
            className="w-full md:w-[220px] px-3 py-2.5 min-h-[44px] md:min-h-0 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
          >
            <option value="">Todos os monitores</option>
            {monitors.map((m) => (
              <option key={m.id} value={m.id}>
                {m.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      {showOpen && (
        <section aria-labelledby="open-title" className="mb-7 md:mb-9">
          <h2 id="open-title" className="flex items-center gap-2.5 text-pg-h3 font-medium text-pg-ardosia mt-0 mb-3">
            Abertos agora
            {open.total > 0 && (
              <span className="min-w-[22px] px-1.5 rounded-full bg-pg-down-bg text-pg-down text-pg-caption font-medium text-center tabular-nums">
                {open.total}
              </span>
            )}
          </h2>
          {open.loading ? (
            <div className="flex justify-center py-8">
              <Spinner />
            </div>
          ) : open.error ? (
            <LoadError message={open.error} onRetry={open.retry} />
          ) : open.items.length === 0 ? (
            <Empty>
              {selected ? `${selected.name} não tem incidente aberto.` : "Nenhum incidente aberto. Tudo no ar."}
            </Empty>
          ) : (
            <div className="flex flex-col gap-3">
              {open.items.map((i) => (
                <OpenIncidentCard key={i.id} incident={i} now={now} />
              ))}
            </div>
          )}
        </section>
      )}

      {showResolved && (
        <section aria-labelledby="resolved-title">
          <h2 id="resolved-title" className="text-pg-h3 font-medium text-pg-ardosia mt-0 mb-3">
            Resolvidos
          </h2>
          {resolved.loading ? (
            <div className="flex justify-center py-8">
              <Spinner />
            </div>
          ) : resolved.error && resolved.items.length === 0 ? (
            <LoadError message={resolved.error} onRetry={resolved.retry} />
          ) : resolved.items.length === 0 ? (
            <Empty>
              {selected ? `${selected.name} não teve incidentes resolvidos` : "Nenhum incidente resolvido"}
              {days ? ` nos últimos ${days} dias.` : "."}
            </Empty>
          ) : (
            <>
              <ResolvedTable
                incidents={resolved.items}
                total={resolved.total}
                onLoadMore={resolved.loadMore}
                loadingMore={resolved.loadingMore}
              />
              {resolved.error && (
                <p role="alert" className="mt-3 text-pg-ui text-pg-down">
                  Não foi possível carregar mais: {resolved.error}
                </p>
              )}
            </>
          )}
        </section>
      )}
    </div>
  );
}
