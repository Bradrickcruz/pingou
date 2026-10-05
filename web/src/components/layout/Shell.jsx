import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  TriangleAlert,
  Settings as SettingsIcon,
  LogOut,
} from "lucide-react";
import { Button } from "../ui/Button";
import marca from "../../assets/brand/pingou-marca-horizontal.svg";

const nav = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/incidents", label: "Incidents", icon: TriangleAlert },
  { to: "/settings", label: "Settings", icon: SettingsIcon },
];

export function Shell({ children, onLogout }) {
  return (
    <div className="flex min-h-screen">
      <aside className="w-[220px] flex-shrink-0 flex flex-col py-6 bg-pg-branco border-r border-pg-borda">
        <div className="px-5 pb-6 border-b border-pg-borda">
          <img src={marca} alt="Pingou" className="h-7 w-auto" />
          <div className="text-pg-caption text-pg-neblina mt-1">health checker</div>
        </div>

        <nav className="py-4 px-3 flex-1">
          {nav.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
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
            Logout
          </Button>
          <span className="text-pg-caption text-pg-neblina">Pingou v1.0</span>
        </div>
      </aside>

      <main className="flex-1 p-8 overflow-y-auto">
        {children}
      </main>
    </div>
  );
}
