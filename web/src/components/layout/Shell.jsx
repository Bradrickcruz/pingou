import { useEffect, useState } from "react";
import { NavLink, useLocation } from "react-router-dom";
import {
  LayoutDashboard,
  TriangleAlert,
  Settings as SettingsIcon,
  LogOut,
} from "lucide-react";
import { Button } from "../ui/Button";
import { BottomNav } from "./BottomNav";
import { incidentsApi } from "../../api/incidents";
import marca from "../../assets/brand/pingou-marca-horizontal.svg";

const nav = [
  { to: "/", label: "Painel", icon: LayoutDashboard },
  { to: "/incidents", label: "Incidentes", icon: TriangleAlert },
  {
    to: "/settings",
    label: "Configurações",
    short: "Ajustes",
    icon: SettingsIcon,
  },
];

// Contagem de incidentes abertos para o badge da sidebar e da barra de abas. Atualiza junto com o
// polling do painel (15 s) e ao trocar de página.
function useOpenIncidents() {
  const [count, setCount] = useState(0);
  const { pathname } = useLocation();

  useEffect(() => {
    let alive = true;
    const load = () =>
      incidentsApi
        .list({ status: "open", limit: 1 })
        .then((res) => alive && setCount(res.total ?? 0))
        .catch(() => {}); // o overlay de conexão já avisa se a API caiu
    load();
    const id = setInterval(load, 15_000);
    return () => {
      alive = false;
      clearInterval(id);
    };
  }, [pathname]);

  return count;
}

export function Shell({ children, onLogout }) {
  const openIncidents = useOpenIncidents();
  const { pathname } = useLocation();

  return (
    <div className="flex min-h-screen flex-col md:flex-row">
      <aside className="hidden md:flex md:w-[220px] md:flex-shrink-0 flex-col py-6 bg-pg-branco border-r border-pg-borda">
        <div className="px-5 pb-6 border-b border-pg-borda">
          <img src={marca} alt="Pingou" className="h-7 w-auto" />
          <div className="text-pg-caption text-pg-neblina mt-1">
            monitor de saúde
          </div>
        </div>

        <nav className="py-4 px-3 flex-1">
          {nav.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
              className={({ isActive }) =>
                `flex items-center gap-2 py-2 px-3 mb-1 rounded-pg-control text-pg-ui transition-colors duration-150 ${
                  isActive || (to === "/" && pathname.startsWith("/monitors/"))
                    ? "bg-pg-ciano-nevoa text-pg-ciano-fundo font-medium"
                    : "text-pg-neblina hover:bg-pg-nuvem hover:text-pg-ardosia"
                }`
              }
            >
              <Icon size={16} aria-hidden="true" />
              {label}
              {to === "/incidents" && openIncidents > 0 && (
                <span className="ml-auto min-w-[22px] px-1.5 rounded-full bg-pg-down-bg text-pg-down text-pg-caption font-medium text-center tabular-nums">
                  {openIncidents}
                  <span className="sr-only">
                    {" "}
                    {openIncidents === 1
                      ? "incidente aberto"
                      : "incidentes abertos"}
                  </span>
                </span>
              )}
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

      <main className="flex-1 min-w-0 p-4 pb-24 md:p-8">{children}</main>

      <BottomNav items={nav} openIncidents={openIncidents} />
    </div>
  );
}
