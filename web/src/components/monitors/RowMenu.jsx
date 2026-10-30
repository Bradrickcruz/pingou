import { useEffect, useId, useRef, useState } from "react";
import { useMatch, useNavigate } from "react-router-dom";
import { EllipsisVertical, Eye, Pencil, Pause, Play, Trash2 } from "lucide-react";

// Menu ⋮ da linha. "Ver detalhes" abre a página do monitor (some nela mesma).
export function RowMenu({ monitor, onEdit, onToggle, onDelete }) {
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
      const items = [...(rootRef.current?.querySelectorAll('[role="menuitem"]') ?? [])];
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
    "flex w-full items-center gap-2.5 px-4 py-2.5 text-left text-pg-ui hover:bg-pg-nuvem focus:bg-pg-nuvem focus:outline-none";

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
        className={`inline-flex items-center justify-center w-10 h-10 rounded-pg-control text-pg-ardosia focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano ${
          open ? "bg-pg-ciano-nevoa" : "hover:bg-pg-nuvem"
        }`}
      >
        <EllipsisVertical size={18} aria-hidden="true" />
      </button>

      {open && (
        <div
          id={menuId}
          role="menu"
          aria-label={`Ações de ${monitor.name}`}
          className="absolute right-0 top-full mt-1 z-30 w-52 py-1.5 rounded-pg-card border border-pg-borda bg-pg-branco shadow-[0_4px_16px_rgba(30,43,51,0.08)]"
        >
          <button type="button" role="menuitem" onClick={run(onEdit)} className={`${item} text-pg-ardosia`}>
            <Pencil size={16} aria-hidden="true" />
            Editar
          </button>
          <button type="button" role="menuitem" onClick={run(onToggle)} className={`${item} text-pg-ardosia`}>
            {monitor.enabled ? <Pause size={16} aria-hidden="true" /> : <Play size={16} aria-hidden="true" />}
            {monitor.enabled ? "Pausar" : "Retomar"}
          </button>
          {!onDetailPage && (
            <button
              type="button"
              role="menuitem"
              onClick={run((m) => navigate(`/monitors/${encodeURIComponent(m.id)}`))}
              className={`${item} text-pg-ardosia`}
            >
              <Eye size={16} aria-hidden="true" />
              Ver detalhes
            </button>
          )}
          <div className="my-1.5 border-t border-pg-borda" role="separator" />
          <button type="button" role="menuitem" onClick={run(onDelete)} className={`${item} text-pg-down`}>
            <Trash2 size={16} aria-hidden="true" />
            Excluir
          </button>
        </div>
      )}
    </div>
  );
}
