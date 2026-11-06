import { useEffect, useId, useRef, useState } from "react";
import { useMatch, useNavigate } from "react-router-dom";
import {
  Copy,
  EllipsisVertical,
  Eye,
  Pencil,
  Pause,
  Play,
  Trash2,
  Wrench,
} from "lucide-react";
import { toast } from "react-toastify";
import { copyText, pingUrl } from "../../lib/monitors";

// Menu ⋮ da linha. "Ver detalhes" abre a página do monitor (some nela mesma).
// onSchedule (opcional) agenda uma manutenção; some se já há janela viva.
export function RowMenu({ monitor, onEdit, onToggle, onDelete, onSchedule }) {
  const navigate = useNavigate();
  const onDetailPage = useMatch("/monitors/:id") != null;
  const [open, setOpen] = useState(false);
  const rootRef = useRef(null);
  const buttonRef = useRef(null);
  const menuId = useId();

  useEffect(() => {
    if (!open) return;
    rootRef.current?.querySelector('[role="menuitem"]')?.focus();
    const onDown = (e) => {
      if (!rootRef.current?.contains(e.target)) setOpen(false);
    };
    const onKey = (e) => {
      const items = [
        ...(rootRef.current?.querySelectorAll('[role="menuitem"]') ?? []),
      ];
      const i = items.indexOf(document.activeElement);
      if (e.key === "Escape") {
        setOpen(false);
        buttonRef.current?.focus();
      } else if (e.key === "ArrowDown") {
        e.preventDefault();
        items[(i + 1) % items.length]?.focus();
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        items[(i - 1 + items.length) % items.length]?.focus();
      } else if (e.key === "Tab") {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", onDown);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onDown);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  const run = (fn) => () => {
    setOpen(false);
    fn(monitor);
  };

  const item =
    "flex w-full items-center gap-2.5 px-4 py-2.5 text-left text-pg-ui hover:bg-pg-hover focus:bg-pg-hover focus:outline-none";

  return (
    <div ref={rootRef} className="relative inline-block">
      <button
        ref={buttonRef}
        type="button"
        aria-haspopup="menu"
        aria-expanded={open}
        aria-controls={open ? menuId : undefined}
        aria-label={`Ações de ${monitor.name}`}
        onClick={() => setOpen((o) => !o)}
        className={`inline-flex items-center justify-center w-10 h-10 rounded-pg-control text-pg-texto focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
          open ? "bg-pg-destaque" : "hover:bg-pg-hover"
        }`}
      >
        <EllipsisVertical size={18} aria-hidden="true" />
      </button>

      {open && (
        <div
          id={menuId}
          role="menu"
          aria-label={`Ações de ${monitor.name}`}
          className="absolute right-0 top-full mt-1 z-30 w-52 py-1.5 rounded-pg-card border border-pg-linha bg-pg-superficie shadow-pg-flutuante"
        >
          <button
            type="button"
            role="menuitem"
            onClick={run(onEdit)}
            className={`${item} text-pg-texto`}
          >
            <Pencil size={16} aria-hidden="true" />
            Editar
          </button>
          <button
            type="button"
            role="menuitem"
            onClick={run(onToggle)}
            className={`${item} text-pg-texto`}
          >
            {monitor.enabled ? (
              <Pause size={16} aria-hidden="true" />
            ) : (
              <Play size={16} aria-hidden="true" />
            )}
            {monitor.enabled ? "Pausar" : "Retomar"}
          </button>
          {onSchedule && !monitor.maintenance && (
            <button
              type="button"
              role="menuitem"
              onClick={run(onSchedule)}
              className={`${item} text-pg-texto`}
            >
              <Wrench size={16} aria-hidden="true" />
              Agendar manutenção
            </button>
          )}
          {!onDetailPage && (
            <button
              type="button"
              role="menuitem"
              onClick={run((m) =>
                navigate(`/monitors/${encodeURIComponent(m.id)}`),
              )}
              className={`${item} text-pg-texto`}
            >
              <Eye size={16} aria-hidden="true" />
              Ver detalhes
            </button>
          )}
          {monitor.ping_token && (
            <button
              type="button"
              role="menuitem"
              onClick={run(async (m) => {
                if (await copyText(pingUrl(m.ping_token)))
                  toast.success("URL de ping copiada");
                else toast.error("Não foi possível copiar a URL");
              })}
              className={`${item} text-pg-texto`}
            >
              <Copy size={16} aria-hidden="true" />
              Copiar URL de ping
            </button>
          )}
          <div className="my-1.5 border-t border-pg-linha" role="separator" />
          <button
            type="button"
            role="menuitem"
            onClick={run(onDelete)}
            className={`${item} text-pg-status-down`}
          >
            <Trash2 size={16} aria-hidden="true" />
            Excluir
          </button>
        </div>
      )}
    </div>
  );
}
