import { useEffect, useRef, useState } from "react";
import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  TriangleAlert,
  Settings as SettingsIcon,
  LogOut,
  Menu,
  X,
} from "lucide-react";
import { Button } from "../ui/Button";
import marca from "../../assets/brand/pingou-marca-horizontal.svg";

const nav = [
  { to: "/", label: "Painel", icon: LayoutDashboard },
  { to: "/incidents", label: "Incidentes", icon: TriangleAlert },
  { to: "/settings", label: "Configurações", icon: SettingsIcon },
];

export function Shell({ children, onLogout }) {
  const [open, setOpen] = useState(false);
  const toggleRef = useRef(null);
  const asideRef = useRef(null);

  useEffect(() => {
    if (!open) return;
    asideRef.current?.querySelector("a")?.focus();
    const onKey = (e) => {
      if (e.key !== "Escape") return;
      setOpen(false);
      toggleRef.current?.focus();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open]);

  const closeMenu = () => {
    setOpen(false);
    toggleRef.current?.focus();
  };

  return (
    <div className="flex min-h-screen flex-col md:flex-row">
      <header className="md:hidden sticky top-0 z-40 flex items-center justify-between px-4 py-3 bg-pg-branco border-b border-pg-borda">
        <img src={marca} alt="Pingou" className="h-6 w-auto" />
        <button
          ref={toggleRef}
          type="button"
          onClick={() => setOpen(true)}
          aria-expanded={open}
          aria-controls="menu-lateral"
          aria-label="Abrir menu"
          className="inline-flex items-center justify-center w-11 h-11 rounded-pg-control text-pg-ardosia hover:bg-pg-nuvem focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
        >
          <Menu size={20} aria-hidden="true" />
        </button>
      </header>

      {open && (
        <div
          className="fixed inset-0 z-50 bg-[rgba(30,43,51,0.5)] md:hidden"
          onClick={closeMenu}
          aria-hidden="true"
        />
      )}

      <aside
        id="menu-lateral"
        ref={asideRef}
        className={`fixed inset-y-0 left-0 z-50 w-[260px] flex flex-col py-6 bg-pg-branco border-r border-pg-borda transition-transform duration-200 md:static md:z-auto md:w-[220px] md:flex-shrink-0 md:translate-x-0 md:visible ${
          open ? "translate-x-0" : "-translate-x-full invisible"
        }`}
      >
        <button
          type="button"
          onClick={closeMenu}
          aria-label="Fechar menu"
          className="md:hidden absolute top-3 right-3 inline-flex items-center justify-center w-11 h-11 rounded-pg-control text-pg-ardosia hover:bg-pg-nuvem focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
        >
          <X size={20} aria-hidden="true" />
        </button>

        <div className="px-5 pb-6 border-b border-pg-borda">
          <img src={marca} alt="Pingou" className="h-7 w-auto" />
          <div className="text-pg-caption text-pg-neblina mt-1">monitor de saúde</div>
        </div>

        <nav className="py-4 px-3 flex-1">
          {nav.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
              onClick={() => setOpen(false)}
              className={({ isActive }) =>
                `flex items-center gap-2 py-2 px-3 mb-1 rounded-pg-control text-pg-ui transition-colors duration-150 ${
                  isActive
                    ? "bg-pg-ciano-nevoa text-pg-ciano-fundo font-medium"
                    : "text-pg-neblina hover:bg-pg-nuvem hover:text-pg-ardosia"
                }`
              }
            >
              <Icon size={16} aria-hidden="true" />
              {label}
            </NavLink>
          ))}
        </nav>

        <div className="py-4 px-5 border-t border-pg-borda flex flex-col gap-2">
          <Button
            variant="ghost"
            onClick={onLogout}
            className="py-1.5 px-3 justify-start gap-2"
          >
            <LogOut size={16} aria-hidden="true" />
            Sair
          </Button>
          <span className="text-pg-caption text-pg-neblina">Pingou v1.0</span>
        </div>
      </aside>

      <main className="flex-1 min-w-0 p-4 md:p-8">{children}</main>
    </div>
  );
}
