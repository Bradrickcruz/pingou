import { NavLink, useLocation } from "react-router-dom";

// Barra de abas fixa do mobile (< md). O desktop usa a sidebar do Shell.
export function BottomNav({ items, openIncidents }) {
  const { pathname } = useLocation();

  return (
    <nav
      aria-label="Principal"
      className="md:hidden fixed inset-x-0 bottom-0 z-40 bg-pg-branco border-t border-pg-borda pb-[env(safe-area-inset-bottom)]"
    >
      <ul className="grid grid-cols-3">
        {items.map(({ to, label, short, icon: Icon }) => (
          <li key={to}>
            <NavLink
              to={to}
              end={to === "/"}
              className={({ isActive }) =>
                `flex flex-col items-center justify-center gap-1 min-h-[56px] py-2 text-pg-caption focus-visible:outline-2 focus-visible:outline focus-visible:-outline-offset-2 focus-visible:outline-pg-ciano ${
                  isActive || (to === "/" && pathname.startsWith("/monitors/"))
                    ? "text-pg-ciano-fundo font-medium"
                    : "text-pg-neblina"
                }`
              }
            >
              <span className="relative">
                <Icon size={22} aria-hidden="true" />
                {to === "/incidents" && openIncidents > 0 && (
                  <span className="absolute -top-1.5 -right-2.5 min-w-[18px] h-[18px] px-1 rounded-full bg-pg-down text-pg-branco text-[11px] leading-[18px] font-medium text-center tabular-nums">
                    {openIncidents}
                    <span className="sr-only">
                      {" "}
                      {openIncidents === 1
                        ? "incidente aberto"
                        : "incidentes abertos"}
                    </span>
                  </span>
                )}
              </span>
              {short ?? label}
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  );
}
