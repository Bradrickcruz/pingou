import { useMemo, useState } from "react";
import { Plus, RefreshCw, Search } from "lucide-react";
import { useMonitors } from "../hooks/useMonitors";
import { useNow } from "../hooks/useNow";
import { monitorsApi } from "../api/monitors";
import { MonitorRow } from "../components/monitors/MonitorRow";
import { MonitorForm } from "../components/monitors/MonitorForm";
import { MonitorModals } from "../components/monitors/MonitorModals";
import { StatusBanner } from "../components/monitors/StatusBanner";
import { Modal } from "../components/ui/Modal";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { timeAgo } from "../lib/format";
import { bySeverity, countByStatus, statusOf } from "../lib/monitors";
import logo from "../assets/brand/pingou-simbolo.svg";

// a busca só aparece quando a lista fica longa (decisão do plano)
const SEARCH_FROM = 10;

const FILTERS = [
  { key: "all", label: "Todos", match: () => true },
  { key: "down", label: "Fora do ar", match: (s) => s === "DOWN" },
  { key: "up", label: "No ar", match: (s) => s === "UP" },
  { key: "paused", label: "Pausados", match: (s) => s === "PAUSED" },
];

export function Dashboard() {
  const { monitors, loading, error, refetch, isFetching, updatedAt } = useMonitors();
  const now = useNow(1000);
  const [modal, setModal] = useState(null);
  const [selected, setSelected] = useState(null);
  const [saving, setSaving] = useState(false);
  const [filter, setFilter] = useState("all");
  const [query, setQuery] = useState("");
  const [actionError, setActionError] = useState(null);

  const counts = useMemo(() => countByStatus(monitors), [monitors]);
  const sorted = useMemo(() => [...monitors].sort(bySeverity), [monitors]);

  const active = FILTERS.find((f) => f.key === filter);
  const q = query.trim().toLocaleLowerCase("pt-BR");
  const visible = sorted.filter(
    (m) => active.match(statusOf(m)) && (!q || m.name.toLocaleLowerCase("pt-BR").includes(q)),
  );
  const filterCount = (f) => monitors.filter((m) => f.match(statusOf(m))).length;

  const openEdit = (m) => {
    setSelected(m);
    setModal("edit");
  };
  const openDelete = (m) => {
    setSelected(m);
    setModal("delete");
  };
  const closeModal = () => {
    setModal(null);
    setSelected(null);
  };

  const run = async (fn) => {
    setSaving(true);
    try {
      await fn();
      await refetch();
      closeModal();
    } finally {
      setSaving(false);
    }
  };

  const handleCreate = (data) => run(() => monitorsApi.create(data));

  const handleToggle = async (m) => {
    setActionError(null);
    try {
      await monitorsApi.update(m.id, { enabled: !m.enabled });
      await refetch();
    } catch (e) {
      setActionError(e.message);
    }
  };

  return (
    <div className="max-w-[1180px]">
      <div className="flex flex-col gap-4 sm:flex-row sm:justify-between sm:items-center mb-7">
        <h1 className="text-pg-h1 font-medium text-pg-ardosia my-0">Monitores</h1>
        <div className="flex items-center gap-3">
          {updatedAt && (
            <span className="text-pg-caption text-pg-neblina whitespace-nowrap" aria-live="off">
              Atualizado {timeAgo(updatedAt, now)}
            </span>
          )}
          <button
            type="button"
            onClick={refetch}
            disabled={isFetching}
            aria-label="Atualizar lista de monitores"
            aria-busy={isFetching}
            className={`flex items-center justify-center w-11 h-11 rounded-pg-control border border-pg-borda bg-pg-branco transition-colors ${
              isFetching ? "opacity-50 cursor-not-allowed" : "hover:bg-pg-nuvem"
            }`}
          >
            <RefreshCw
              size={16}
              aria-hidden="true"
              className={isFetching ? "animate-spin motion-reduce:animate-none" : ""}
            />
          </button>
          <Button onClick={() => setModal("create")} className="inline-flex items-center gap-1.5 h-11">
            <Plus size={16} aria-hidden="true" />
            Novo monitor
          </Button>
        </div>
      </div>

      {loading && (
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      )}

      {error && (
        <p className="text-pg-ui text-pg-down mb-4" role="alert">
          {error}
        </p>
      )}
      {actionError && (
        <p className="text-pg-ui text-pg-down mb-4" role="alert">
          {actionError}
        </p>
      )}

      {!loading && monitors.length === 0 && (
        <div className="text-center py-16 px-4 rounded-pg-card border border-dashed border-pg-borda text-pg-neblina">
          <img src={logo} alt="" width={40} height={40} className="mx-auto mb-3" />
          <p className="font-medium text-pg-ardosia mb-1.5">Nenhum monitor ainda</p>
          <p className="text-pg-ui mb-5">
            Monitore um site, uma porta TCP ou um job agendado.
          </p>
          <Button onClick={() => setModal("create")} className="inline-flex items-center gap-1.5">
            <Plus size={16} aria-hidden="true" />
            Novo monitor
          </Button>
        </div>
      )}

      {monitors.length > 0 && (
        <>
          <StatusBanner monitors={monitors} counts={counts} />

          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between mb-5">
            <div role="group" aria-label="Filtrar monitores" className="flex flex-wrap gap-2">
              {FILTERS.map((f) => {
                const on = filter === f.key;
                return (
                  <button
                    key={f.key}
                    type="button"
                    aria-pressed={on}
                    onClick={() => setFilter(f.key)}
                    className={`px-4 py-2 rounded-full border text-pg-ui transition-colors duration-150 focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
                      on
                        ? "bg-pg-ardosia border-pg-ardosia text-pg-branco"
                        : "bg-pg-branco border-pg-borda text-pg-ardosia hover:bg-pg-nuvem"
                    }`}
                  >
                    {f.label} ({filterCount(f)})
                  </button>
                );
              })}
            </div>
            {monitors.length >= SEARCH_FROM && (
              <div className="relative sm:w-[280px]">
                <Search
                  size={16}
                  aria-hidden="true"
                  className="absolute left-3 top-1/2 -translate-y-1/2 text-pg-neblina"
                />
                <input
                  type="search"
                  aria-label="Buscar por nome"
                  placeholder="Buscar por nome"
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  className="w-full pl-9 pr-3 py-2.5 rounded-pg-control text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda placeholder:text-pg-neblina"
                />
              </div>
            )}
          </div>

          <div className="rounded-pg-card border border-pg-borda bg-pg-branco">
            <table className="w-full border-collapse text-left">
              <caption className="sr-only">Monitores, ordenados por gravidade</caption>
              <thead>
                <tr className="bg-pg-ciano-nevoa text-pg-caption text-pg-ciano-fundo">
                  <th scope="col" className="px-5 py-3 font-medium rounded-tl-pg-card">Status</th>
                  <th scope="col" className="px-5 py-3 font-medium">Monitor</th>
                  <th scope="col" className="px-5 py-3 font-medium hidden lg:table-cell whitespace-nowrap">
                    Últimos 30 checks
                  </th>
                  <th scope="col" className="px-5 py-3 font-medium hidden md:table-cell text-right whitespace-nowrap">
                    Uptime 24 h
                  </th>
                  <th scope="col" className="px-5 py-3 font-medium hidden md:table-cell text-right">Resposta</th>
                  <th scope="col" className="px-5 py-3 font-medium hidden sm:table-cell whitespace-nowrap">
                    Último check
                  </th>
                  <th scope="col" className="px-5 py-3 rounded-tr-pg-card">
                    <span className="sr-only">Ações</span>
                  </th>
                </tr>
              </thead>
              <tbody>
                {visible.map((m) => (
                  <MonitorRow
                    key={m.id}
                    monitor={m}
                    now={now}
                    onEdit={openEdit}
                    onToggle={handleToggle}
                    onDelete={openDelete}
                  />
                ))}
                {visible.length === 0 && (
                  <tr className="border-t border-pg-borda">
                    <td colSpan={7} className="px-5 py-10 text-center text-pg-ui text-pg-neblina">
                      Nenhum monitor neste filtro.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
            <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between px-5 py-3 border-t border-pg-borda text-pg-caption text-pg-neblina">
              <span>Ordenado por gravidade: fora do ar primeiro.</span>
              <span className="flex gap-4" aria-hidden="true">
                <span className="inline-flex items-center gap-1.5">
                  <span className="w-1.5 h-3 rounded-[1px] bg-pg-up" /> no ar
                </span>
                <span className="inline-flex items-center gap-1.5">
                  <span className="w-1.5 h-3 rounded-[1px] bg-pg-down" /> fora do ar
                </span>
                <span className="inline-flex items-center gap-1.5">
                  <span className="w-1.5 h-3 rounded-[1px] bg-pg-borda" /> sem dados
                </span>
              </span>
            </div>
          </div>
        </>
      )}

      {modal === "create" && (
        <Modal title="Novo monitor" onClose={closeModal}>
          <MonitorForm onSubmit={handleCreate} onCancel={closeModal} loading={saving} />
        </Modal>
      )}

      <MonitorModals
        modal={modal}
        selected={selected}
        onModalChange={setModal}
        onSelectedChange={setSelected}
        onClose={closeModal}
        onSaved={refetch}
        onDeleted={refetch}
      />
    </div>
  );
}
