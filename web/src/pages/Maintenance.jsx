import { useCallback, useEffect, useRef, useState } from "react";
import { Plus, Wrench } from "lucide-react";
import { maintenanceApi } from "../api/maintenance";
import { useNow } from "../hooks/useNow";
import { MobileHeader } from "../components/layout/MobileHeader";
import { MaintenanceModal } from "../components/monitors/MaintenanceModal";
import { WindowsTable } from "../components/maintenance/WindowsTable";
import { CancelWindowDialog } from "../components/maintenance/CancelWindowDialog";
import { Button } from "../components/ui/Button";
import { Spinner } from "../components/ui/Spinner";
import { windowState } from "../lib/maintenance";
import { getTimeZone } from "../lib/timezone";

// Janelas vivas, a cada 15 s (ritmo do painel). `alive` descarta resposta atrasada.
function useWindows() {
  const [state, setState] = useState({
    windows: [],
    loading: true,
    error: null,
  });

  const load = useCallback(async (alive = () => true) => {
    try {
      const res = await maintenanceApi.list();
      if (alive())
        setState({ windows: res.windows, loading: false, error: null });
    } catch (e) {
      if (alive())
        setState((s) => ({ ...s, loading: false, error: e.message }));
    }
  }, []);

  useEffect(() => {
    let alive = true;
    const isAlive = () => alive;
    void load(isAlive);
    const id = setInterval(() => void load(isAlive), 15_000);
    return () => {
      alive = false;
      clearInterval(id);
    };
  }, [load]);

  return { ...state, reload: load };
}

// Recarrega quando o relógio passa do início ou do fim de alguma janela (uma
// vez por virada), para a lista acompanhar o servidor sem esperar o polling.
function useReloadOnBoundary(windows, now, reload) {
  const seen = useRef(new Set());
  useEffect(() => {
    for (const w of windows) {
      for (const at of [w.starts_at, w.ends_at]) {
        const key = `${w.id}@${at}`;
        if (new Date(at).getTime() <= now && !seen.current.has(key)) {
          seen.current.add(key);
          void reload();
          return;
        }
      }
    }
  }, [windows, now, reload]);
}

export function Maintenance() {
  const now = useNow(1000);
  const { windows, loading, error, reload } = useWindows();
  const [modal, setModal] = useState(null); // { kind: "create" | "edit" | "cancel", window? }
  useReloadOnBoundary(windows, now, reload);

  // a janela que terminou sai da lista na hora; o servidor confirma no próximo
  // recarregamento (o fim é aplicado em até 10 s)
  const visible = windows.filter((w) => new Date(w.ends_at).getTime() > now);
  const close = () => setModal(null);
  const refresh = () => reload();

  const scheduleButton = (
    <Button
      onClick={() => setModal({ kind: "create" })}
      className="inline-flex items-center gap-1.5 h-11"
    >
      <Plus size={16} aria-hidden="true" />
      Agendar manutenção
    </Button>
  );

  return (
    <div className="max-w-[1680px]">
      <MobileHeader
        title="Manutenções"
        subtitle={`Horários em ${getTimeZone()}`}
        actions={
          <Button
            onClick={() => setModal({ kind: "create" })}
            className="w-11 h-11 !p-0"
          >
            <Plus size={20} aria-hidden="true" />
            <span className="sr-only">Agendar manutenção</span>
          </Button>
        }
      />
      <div className="hidden md:flex md:items-start md:justify-between gap-4 mb-7">
        <div>
          <h1 className="text-pg-h1 font-medium text-pg-texto my-0">
            Manutenções
          </h1>
          <p className="mt-2 mb-0 text-pg-ui text-pg-texto-2">
            Janelas agendadas e em curso. Horários em {getTimeZone()}.
          </p>
        </div>
        {scheduleButton}
      </div>

      {loading && (
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      )}

      {error && (
        <p className="text-pg-ui text-pg-status-down mb-4" role="alert">
          {error}{" "}
          <Button variant="ghost" onClick={() => reload()} className="ml-2">
            Tentar de novo
          </Button>
        </p>
      )}

      {!loading && !error && visible.length === 0 && (
        <div className="text-center py-16 px-4 rounded-pg-card border border-dashed border-pg-linha text-pg-texto-2">
          <Wrench
            size={32}
            aria-hidden="true"
            className="mx-auto mb-3 text-pg-link"
          />
          <p className="font-medium text-pg-texto mb-1.5">
            Nenhuma manutenção agendada
          </p>
          <p className="text-pg-ui mb-5">
            Durante uma janela, o monitor não é checado e ninguém é avisado.
          </p>
          {scheduleButton}
        </div>
      )}

      {visible.length > 0 && (
        <div aria-live="polite">
          <WindowsTable
            windows={visible}
            now={now}
            onEdit={(w) => setModal({ kind: "edit", window: w })}
            onCancel={(w) => setModal({ kind: "cancel", window: w })}
          />
        </div>
      )}

      {modal?.kind === "create" && (
        <MaintenanceModal onClose={close} onCreated={refresh} />
      )}
      {modal?.kind === "edit" && (
        <MaintenanceModal
          window={modal.window}
          onClose={close}
          onSaved={refresh}
          onStale={refresh}
        />
      )}
      {modal?.kind === "cancel" && (
        <CancelWindowDialog
          window={modal.window}
          active={windowState(modal.window, now) === "active"}
          onClose={close}
          onDone={refresh}
        />
      )}
    </div>
  );
}
