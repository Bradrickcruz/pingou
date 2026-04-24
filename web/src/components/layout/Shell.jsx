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
    <div style={{ display: "flex", minHeight: "100vh" }}>
      <aside
        style={{
          width: "220px",
          flexShrink: 0,
          background: t.colors.surface,
          borderRight: `1px solid ${t.colors.border}`,
          display: "flex",
          flexDirection: "column",
          padding: "24px 0",
        }}
      >
        <div
          style={{
            padding: "0 20px 24px",
            borderBottom: `1px solid ${t.colors.border}`,
          }}
        >
          <div
            style={{
              fontWeight: 700,
              fontSize: "16px",
              color: t.colors.primary,
            }}
          >
            🏓 Pingou
          </div>
          <div
            style={{
              fontSize: "11px",
              color: t.colors.textMuted,
              marginTop: "2px",
            }}
          >
            health checker
          </div>
        </div>

        <nav style={{ padding: "16px 12px", flex: 1 }}>
          {nav.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/"}
              style={({ isActive }) => ({
                display: "block",
                padding: "9px 12px",
                borderRadius: t.radius.sm,
                marginBottom: "4px",
                color: isActive ? t.colors.textPrimary : t.colors.textMuted,
                background: isActive ? t.colors.surfaceAlt : "transparent",
                fontWeight: isActive ? 600 : 400,
                fontSize: "13px",
                transition: "all 0.15s",
              })}
            >
              {label}
            </NavLink>
          ))}
        </nav>

        <div
          style={{
            padding: "16px 20px",
            borderTop: `1px solid ${t.colors.border}`,
            display: "flex",
            flexDirection: "column",
            gap: "8px",
          }}
        >
          <Button
            variant="ghost"
            onClick={onLogout}
            style={{ fontSize: "12px", padding: "6px 12px", textAlign: "left" }}
          >
            ⎋ Logout
          </Button>
          <span style={{ fontSize: "11px", color: t.colors.textMuted }}>
            Pingou v1.0
          </span>
        </div>
      </aside>

      <main style={{ flex: 1, padding: "32px", overflowY: "auto" }}>
        {children}
      </main>
    </div>
  );
}
