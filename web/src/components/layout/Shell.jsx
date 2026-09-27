import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  TriangleAlert,
  Settings as SettingsIcon,
  LogOut,
} from "lucide-react";
import { tokens as t } from "../../theme/tokens";
import { Button } from "../ui/Button";
import logo from "../../assets/brand/pingou-simbolo.svg";

const nav = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/incidents", label: "Incidents", icon: TriangleAlert },
  { to: "/settings", label: "Settings", icon: SettingsIcon },
];

export function Shell({ children, onLogout }) {
  return (
    <div className="flex min-h-screen">
      <aside
        className="w-[220px] flex-shrink-0 flex flex-col py-6 px-0"
        style={{
          background: t.colors.surface,
          borderRight: `1px solid ${t.colors.border}`,
        }}
      >
        <div
          className="px-5 pb-6 border-b"
          style={{
            borderColor: t.colors.border,
          }}
        >
          <div className="flex items-center gap-2">
            <img src={logo} alt="" width={28} height={28} />
            <span
              className="font-bold text-base"
              style={{
                color: t.colors.primary,
              }}
            >
              Pingou
            </span>
          </div>
          <div
            className="text-[11px] mt-0.5"
            style={{
              color: t.colors.textMuted,
            }}
          >
            health checker
          </div>
        </div>

        <nav className="py-4 px-3 flex-1">
          {nav.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
              className={({ isActive }) =>
                `flex items-center gap-2 py-2 px-3 rounded mb-1 text-[13px] transition-all duration-150 ${
                  isActive
                    ? "bg-[var(--social-bg)] font-semibold"
                    : "text-[var(--text)] font-normal"
                }`
              }
              style={({ isActive }) => ({
                color: isActive ? t.colors.textPrimary : t.colors.textMuted,
                background: isActive ? t.colors.surfaceAlt : "transparent",
                fontWeight: isActive ? 600 : 400,
              })}
            >
              <Icon size={16} aria-hidden="true" />
              {label}
            </NavLink>
          ))}
        </nav>

        <div
          className="py-4 px-5 border-t flex flex-col gap-2"
          style={{
            borderColor: t.colors.border,
          }}
        >
          <Button
            variant="ghost"
            onClick={onLogout}
            className="text-[12px] py-1.5 px-3 text-left inline-flex items-center gap-2"
          >
            <LogOut size={16} aria-hidden="true" />
            Logout
          </Button>
          <span
            className="text-[11px]"
            style={{
              color: t.colors.textMuted,
            }}
          >
            Pingou v1.0
          </span>
        </div>
      </aside>

      <main className="flex-1 p-8 overflow-y-auto">
        {children}
      </main>
    </div>
  );
}