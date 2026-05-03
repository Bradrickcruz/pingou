import { NavLink } from "react-router-dom";
import { tokens as t } from "../../theme/tokens";
import { Button } from "../ui/Button";

const nav = [
  { to: "/", label: "⬡  Dashboard" },
  { to: "/incidents", label: "⚠  Incidents" },
  { to: "/settings", label: "⚙  Settings" },
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
          <div
            className="font-bold text-base"
            style={{
              color: t.colors.primary,
            }}
          >
            🏓 Pingou
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
          {nav.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
              className={({ isActive }) =>
                `block py-2 px-3 rounded mb-1 text-[13px] transition-all duration-150 ${
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
            className="text-[12px] py-1.5 px-3 text-left"
          >
            ⎋ Logout
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