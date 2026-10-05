import { useConnection } from "../../hooks/connectionContext";

export function ConnectionOverlay({ children }) {
  const { online } = useConnection();

  if (online) {
    return children;
  }

  return (
    <div className="relative min-h-screen">
      <div
        className="absolute inset-0 pointer-events-none z-50 bg-[rgba(30,43,51,0.5)] backdrop-blur-sm"
      />
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none z-[51]">
        <div className="p-6 rounded-pg-card border border-pg-down bg-pg-branco text-center">
          <p className="text-pg-h3 font-medium text-pg-down">connection lost</p>
          <p className="text-pg-ui text-pg-neblina mt-1">attempting to reconnect...</p>
        </div>
      </div>
      {children}
    </div>
  );
}
