import { useConnection } from "../../hooks/useConnection";
import { tokens as t } from "../../theme/tokens";

export function ConnectionOverlay({ children }) {
  const { online } = useConnection();

  if (online) {
    return children;
  }

  return (
    <div className="relative min-h-screen">
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          backdropFilter: "blur(8px)",
          background: "rgba(0, 0, 0, 0.5)",
          zIndex: 50,
        }}
      />
      <div
        className="absolute inset-0 flex items-center justify-center pointer-events-none"
        style={{ zIndex: 51 }}
      >
        <div
          className="p-6 rounded-lg border text-center"
          style={{
            background: t.colors.surface,
            borderColor: t.colors.danger,
          }}
        >
          <p className="text-lg font-semibold" style={{ color: t.colors.danger }}>
            connection lost
          </p>
          <p className="text-sm mt-1" style={{ color: t.colors.textMuted }}>
            attempting to reconnect...
          </p>
        </div>
      </div>
      {children}
    </div>
  );
}