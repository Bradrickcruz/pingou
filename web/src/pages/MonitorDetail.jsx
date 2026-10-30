import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useParams, useSearchParams } from "react-router-dom";
import { Activity, ChevronLeft, Globe, Pause, Pencil, Play, Server } from "lucide-react";
import { toast } from "react-toastify";
import { monitorsApi } from "../api/monitors";
import { useNow } from "../hooks/useNow";
import { useSettings } from "../hooks/useSettings";
import { Badge } from "../components/ui/Badge";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { RowMenu } from "../components/monitors/RowMenu";
import { MonitorModals } from "../components/monitors/MonitorModals";
import { RangeTabs } from "../components/monitor-detail/RangeTabs";
import { RANGES, clampRange } from "../lib/ranges";
import { ChartLegend, ResponseChart } from "../components/monitor-detail/ResponseChart";
import { ChecksTable } from "../components/monitor-detail/ChecksTable";
import { RecentIncidents } from "../components/monitor-detail/RecentIncidents";
import { ConfigCard } from "../components/monitor-detail/ConfigCard";
import { DetailTopBar } from "../components/monitor-detail/DetailTopBar";
import { DownAlert } from "../components/monitor-detail/DownAlert";
import { formatDuration, formatUptime, timeAgo } from "../lib/format";
import { statusOf, targetOf, typeLabel } from "../lib/monitors";

const typeIcon = { http: Globe, tcp: Server, heartbeat: Activity };
const FIRST_CHECKS = 8;
const MORE_CHECKS = 20;
const MAX_CHECKS = 200; // limite da API por request

// Monitor, últimos checks e incidentes recentes, a cada 15 s (ritmo do painel).
// `alive` descarta respostas atrasadas ao trocar de monitor.
function useMonitorDetail(id, checkLimit) {
  const [state, setState] = useState({ monitor: null, checks: [], checksTotal: 0, incidents: [], loading: true, error: null, notFound: false });

  const load = useCallback(
    async (alive = () => true) => {
      try {
        const [monitor, checks, incidents] = await Promise.all([
          monitorsApi.get(id),
          monitorsApi.checks(id, { limit: checkLimit }),
          monitorsApi.incidents(id, { limit: 3 }),
        ]);
        if (!alive()) return;
        setState({
          monitor,
          checks: checks.data,
          checksTotal: checks.total,
          incidents: incidents.data,
          loading: false,
          error: null,
          notFound: false,
        });
      } catch (e) {
        if (!alive()) return;
        setState((s) => ({ ...s, loading: false, error: e.message, notFound: e.status === 404 }));
      }
    },
    [id, checkLimit],
  );

  useEffect(() => {
    let alive = true;
    const isAlive = () => alive;
    void load(isAlive);
    const interval = setInterval(() => void load(isAlive), 15_000);
    return () => {
      alive = false;
      clearInterval(interval);
    };
  }, [load]);

  return { ...state, reload: load };
}

// Estatísticas da janela: ao trocar de período e a cada 60 s.
function useStats(id, range) {
  const [state, setState] = useState({ stats: null, error: null });

  const load = useCallback(
    async (alive = () => true) => {
      if (!range) return;
      try {
        const stats = await monitorsApi.stats(id, range);
        if (alive()) setState({ stats, error: null });
      } catch (e) {
        if (alive()) setState((s) => ({ ...s, error: e.message }));
      }
    },
    [id, range],
  );

  useEffect(() => {
    let alive = true;
    const isAlive = () => alive;
    void load(isAlive);
    const interval = setInterval(() => void load(isAlive), 60_000);
    return () => {
      alive = false;
      clearInterval(interval);
    };
  }, [load]);

  // dados de outra janela não servem enquanto a nova carrega
  const stats = state.stats?.range === range ? state.stats : null;
  return { stats, error: state.error, reload: load };
}

// short/shortValue: rótulo e valor do mobile (< md), que não mostra a nota;
// alert: valor em vermelho no mobile.
function StatCard({ label, short, value, shortValue, note, alert = false, className = "" }) {
  return (
    <div className={`p-4 rounded-pg-card border border-pg-borda bg-pg-branco min-w-0 ${className}`}>
      <div className="text-pg-caption text-pg-neblina">
        <span className="md:hidden">{short ?? label}</span>
        <span className="hidden md:inline">{label}</span>
      </div>
      <div
        className={`mt-1 font-pg-heading text-pg-h3 md:text-pg-h2 font-medium tabular-nums break-words md:truncate ${
          alert ? "text-pg-down md:text-pg-ardosia" : "text-pg-ardosia"
        }`}
      >
        {shortValue != null ? (
          <>
            <span className="md:hidden">{shortValue}</span>
            <span className="hidden md:inline">{value}</span>
          </>
        ) : (
          value
        )}
      </div>
      <div className="hidden md:block mt-1 text-pg-caption text-pg-neblina">{note}</div>
    </div>
  );
}

function SectionTitle({ children, action }) {
  return (
    <div className="flex items-baseline justify-between gap-3 mb-3">
      <h2 className="text-pg-h3 font-medium text-pg-ardosia my-0">{children}</h2>
      {action}
    </div>
  );
}

const linkClass = "text-pg-ui text-pg-ciano-fundo underline focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano";

function Performance({ monitor, stats, statsError, range, onRange, retentionDays, now, onRetry }) {
  const heartbeat = monitor.type === "heartbeat";
  const current = RANGES.find((r) => r.key === range);
  const since = current?.since ?? "";
  const empty = <span aria-label="sem dados">—</span>;

  return (
    <section aria-labelledby="desempenho" className="mb-6 md:mb-9">
      <div className="flex flex-col md:flex-row md:flex-wrap md:items-center md:justify-between gap-3 mb-3 md:mb-4">
        <h2 id="desempenho" className="text-pg-h3 font-medium text-pg-ardosia my-0">
          Desempenho
        </h2>
        <RangeTabs value={range} onChange={onRange} retentionDays={retentionDays} />
      </div>

      {statsError && !stats && (
        <p className="text-pg-ui text-pg-down mb-4" role="alert">
          {statsError}{" "}
          <button type="button" onClick={onRetry} className={`${linkClass} p-0 border-0 bg-transparent cursor-pointer`}>
            Tentar de novo
          </button>
        </p>
      )}

      <div className={`grid grid-cols-2 gap-3 mb-3 ${heartbeat ? "lg:grid-cols-3" : "lg:grid-cols-4"}`}>
        <StatCard
          label="Uptime"
          short={`Uptime ${current?.short ?? ""}`}
          value={stats?.uptime != null ? formatUptime(stats.uptime) : empty}
          note={since}
          alert={stats?.uptime != null && stats.uptime < 100}
        />
        {!heartbeat && (
          <StatCard
            label="Resposta média"
            value={stats?.avg_latency_ms != null ? `${stats.avg_latency_ms} ms` : empty}
            note="só checks com sucesso"
          />
        )}
        <StatCard label="Incidentes" value={stats ? stats.incidents : empty} note={since} alert={stats?.incidents > 0} />
        <StatCard
          label="Tempo fora do ar"
          short="Tempo fora"
          className={heartbeat ? "col-span-2 md:col-span-1" : ""}
          value={stats ? (stats.downtime_seconds > 0 ? formatDuration(stats.downtime_seconds) : "nenhum") : empty}
          shortValue={stats?.downtime_seconds > 0 ? formatDuration(stats.downtime_seconds).replace(" e ", " ") : null}
          note={`somado ${since}`}
        />
      </div>

      {!heartbeat && (
        <div className="p-4 sm:p-5 rounded-pg-card border border-pg-borda bg-pg-branco">
          <div className="flex flex-col md:flex-row md:flex-wrap md:items-center md:justify-between gap-1 md:gap-2 mb-4">
            <h3 className="text-pg-ui font-medium text-pg-ardosia my-0">Tempo de resposta</h3>
            <ChartLegend bucketSeconds={stats?.bucket_seconds} />
          </div>
          {stats ? (
            <ResponseChart stats={stats} now={now} />
          ) : (
            <div className="flex justify-center items-center h-[200px]">{!statsError && <Spinner />}</div>
          )}
        </div>
      )}
    </section>
  );
}

export function MonitorDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [params, setParams] = useSearchParams();
  const now = useNow(1000);
  const { settings, loading: settingsLoading } = useSettings();
  const [checkLimit, setCheckLimit] = useState(FIRST_CHECKS);
  const [modal, setModal] = useState(null);
  const [selected, setSelected] = useState(null);
  const [actionError, setActionError] = useState(null);

  // período: só 24h e 7d valem antes de saber a retenção (mínimo 7 dias)
  const retentionDays = settings?.retention_days ?? 7;
  const asked = params.get("range") ?? "24h";
  const range = settingsLoading ? (["24h", "7d"].includes(asked) ? asked : null) : clampRange(asked, retentionDays);

  // link antigo ou retenção reduzida: reescreve a URL com o período permitido
  useEffect(() => {
    if (settingsLoading || !params.has("range") || asked === range) return;
    setParams((p) => {
      const next = new URLSearchParams(p);
      next.set("range", range);
      return next;
    }, { replace: true });
  }, [settingsLoading, asked, range, params, setParams]);

  const { monitor, checks, checksTotal, incidents, loading, error, notFound, reload } = useMonitorDetail(id, checkLimit);
  const { stats, error: statsError, reload: reloadStats } = useStats(id, range);

  const setRange = (key) =>
    setParams((p) => {
      const next = new URLSearchParams(p);
      if (key === "24h") next.delete("range");
      else next.set("range", key);
      return next;
    });

  const loadMore = () => setCheckLimit((l) => Math.min(MAX_CHECKS, l + MORE_CHECKS));
  // pediu mais do que chegou: a página seguinte ainda está a caminho
  const loadingMore = checks.length < Math.min(checkLimit, checksTotal);

  const refresh = async () => {
    await Promise.all([reload(), reloadStats()]);
  };

  const openModal = (kind) => (m) => {
    setSelected(m);
    setModal(kind);
  };
  const closeModal = () => {
    setModal(null);
    setSelected(null);
  };

  const handleToggle = async (m) => {
    setActionError(null);
    try {
      await monitorsApi.update(m.id, { enabled: !m.enabled });
      toast.success(m.enabled ? "Monitor pausado" : "Monitor retomado");
      await reload();
    } catch (e) {
      setActionError(e.message);
    }
  };

  const back = (
    <Link to="/" className="hidden md:inline-flex items-center gap-1 mb-5 text-pg-ui text-pg-ciano-fundo no-underline hover:underline">
      <ChevronLeft size={16} aria-hidden="true" />
      Monitores
    </Link>
  );

  if (loading) {
    return (
      <div className="max-w-[1180px]">
        <DetailTopBar />
        {back}
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      </div>
    );
  }

  if (!monitor) {
    return (
      <div className="max-w-[1180px]">
        <DetailTopBar />
        {back}
        {notFound ? (
          <div className="text-center py-16 px-4 rounded-pg-card border border-dashed border-pg-borda">
            <p className="font-medium text-pg-ardosia mb-1.5">Monitor não encontrado</p>
            <p className="text-pg-ui text-pg-neblina mb-5">Ele pode ter sido excluído.</p>
            <Link to="/" className={linkClass}>
              Voltar ao painel
            </Link>
          </div>
        ) : (
          <p className="text-pg-ui text-pg-down" role="alert">
            {error}{" "}
            <Button variant="ghost" onClick={() => reload()} className="ml-2">
              Tentar de novo
            </Button>
          </p>
        )}
      </div>
    );
  }

  const status = statusOf(monitor);
  const paused = status === "PAUSED";
  const heartbeat = monitor.type === "heartbeat";
  const Icon = typeIcon[monitor.type] ?? Globe;
  const target = targetOf(monitor);
  const lastAt = heartbeat ? monitor.last_ping_at : monitor.last_checked_at;
  const lastLabel = heartbeat ? "último ping" : "último check";
  const last = paused ? "pausado" : lastAt ? `${lastLabel} ${timeAgo(lastAt, now)}` : `${lastLabel}: ainda não houve`;
  // o incidente aberto é sempre o mais novo; fora dos 3 carregados, o cartão não aparece
  const openIncident = status === "DOWN" ? incidents.find((i) => i.open) : null;
  const menu = (
    <RowMenu monitor={monitor} onEdit={openModal("edit")} onToggle={handleToggle} onDelete={openModal("delete")} />
  );

  return (
    <div className="max-w-[1180px]">
      <DetailTopBar menu={menu} />
      {back}

      <header className="flex items-start justify-between gap-4 mb-4 md:mb-8">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <h1 className="text-pg-h1 font-medium text-pg-ardosia my-0 break-words min-w-0">{monitor.name}</h1>
            <Badge state={status} />
          </div>
          <div className="flex flex-col items-start md:flex-row md:flex-wrap md:items-center gap-x-4 gap-y-1 mt-2 text-pg-ui text-pg-neblina">
            <span className="inline-flex items-center gap-1.5 min-w-0 max-w-full">
              <Icon size={14} aria-hidden="true" className="flex-shrink-0" />
              <span className="sr-only">{monitor.type.toUpperCase()}:</span>
              {monitor.type === "http" ? (
                <a
                  href={monitor.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="relative font-pg-mono text-pg-caption text-pg-ciano-fundo underline truncate"
                  title={target}
                >
                  {target}
                  <span className="sr-only"> (abre em nova aba)</span>
                </a>
              ) : (
                <span className={`truncate ${heartbeat ? "" : "font-pg-mono text-pg-caption"}`} title={target}>
                  {target}
                </span>
              )}
            </span>
            <span className="text-pg-caption md:text-pg-ui">
              <span className="md:hidden">{typeLabel[monitor.type] ?? monitor.type} · </span>
              a cada {formatDuration(monitor.interval_seconds)} · {last}
            </span>
          </div>
        </div>
        <div className="flex flex-shrink-0 items-center gap-2">
          <Button variant="ghost" onClick={() => handleToggle(monitor)} className="hidden md:inline-flex gap-1.5 h-10 bg-pg-branco">
            {monitor.enabled ? <Pause size={16} aria-hidden="true" /> : <Play size={16} aria-hidden="true" />}
            {monitor.enabled ? "Pausar" : "Retomar"}
          </Button>
          <Button variant="ghost" onClick={() => openModal("edit")(monitor)} className="hidden md:inline-flex gap-1.5 h-10 bg-pg-branco">
            <Pencil size={16} aria-hidden="true" />
            Editar
          </Button>
          <div className="hidden md:block rounded-pg-control border border-pg-borda bg-pg-branco">{menu}</div>
        </div>
      </header>

      {openIncident && <DownAlert incident={openIncident} now={now} />}

      <div className="md:hidden grid grid-cols-2 gap-2 mb-6">
        <Button variant="ghost" onClick={() => handleToggle(monitor)} className="gap-1.5 min-h-[44px] bg-pg-branco">
          {monitor.enabled ? <Pause size={16} aria-hidden="true" /> : <Play size={16} aria-hidden="true" />}
          {monitor.enabled ? "Pausar" : "Retomar"}
        </Button>
        <Button variant="ghost" onClick={() => openModal("edit")(monitor)} className="gap-1.5 min-h-[44px] bg-pg-branco">
          <Pencil size={16} aria-hidden="true" />
          Editar
        </Button>
      </div>

      {(actionError || (error && monitor)) && (
        <p className="text-pg-ui text-pg-down mb-4" role="alert">
          {actionError ?? error}
        </p>
      )}

      <Performance
        monitor={monitor}
        stats={stats}
        statsError={statsError}
        range={range ?? "24h"}
        onRange={setRange}
        retentionDays={retentionDays}
        now={now}
        onRetry={() => reloadStats()}
      />

      <div className="grid gap-6 md:gap-8 lg:grid-cols-3">
        <section aria-labelledby="ultimos-checks" className="lg:col-span-2 min-w-0">
          <SectionTitle>
            <span id="ultimos-checks">Últimos checks</span>
          </SectionTitle>
          <ChecksTable
            checks={checks}
            total={Math.min(checksTotal, MAX_CHECKS)}
            type={monitor.type}
            now={now}
            onMore={loadMore}
            loadingMore={loadingMore}
            compact={checkLimit === FIRST_CHECKS}
          />
        </section>

        <div className="flex flex-col gap-6 md:gap-8 min-w-0">
          <section aria-labelledby="incidentes">
            <SectionTitle
              action={
                <Link to={`/incidents?monitor=${encodeURIComponent(monitor.id)}`} className={linkClass}>
                  Ver todos
                </Link>
              }
            >
              <span id="incidentes">Incidentes</span>
            </SectionTitle>
            <RecentIncidents incidents={incidents} now={now} />
          </section>

          <section aria-labelledby="configuracao">
            <SectionTitle
              action={
                <button
                  type="button"
                  onClick={() => openModal("edit")(monitor)}
                  className={`${linkClass} hidden md:inline p-0 border-0 bg-transparent cursor-pointer`}
                >
                  Editar
                </button>
              }
            >
              <span id="configuracao">Configuração</span>
            </SectionTitle>
            <ConfigCard monitor={monitor} />
          </section>
        </div>
      </div>

      <MonitorModals
        modal={modal}
        selected={selected}
        onModalChange={setModal}
        onSelectedChange={setSelected}
        onClose={closeModal}
        onSaved={refresh}
        onDeleted={() => navigate("/")}
      />
    </div>
  );
}
