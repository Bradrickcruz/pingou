import { WifiOff } from "lucide-react";
import { useConnection } from "../../hooks/connectionContext";
import { useNow } from "../../hooks/useNow";

export function ConnectionOverlay({ children }) {
  const { online, offlineSince, nextCheckAt } = useConnection();
  const now = useNow();

  if (online) {
    return children;
  }

  const since = new Date(offlineSince ?? now).toLocaleString("pt-BR");
  const seconds = Math.max(0, Math.ceil(((nextCheckAt ?? now) - now) / 1000));

  return (
    <div className="relative min-h-screen">
      <div className="absolute inset-0 pointer-events-none z-50 bg-[rgba(30,43,51,0.5)] backdrop-blur-sm" />
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none z-[51]">
        <div
          role="alert"
          className="flex items-center gap-4 p-8 rounded-pg-card border border-pg-neblina bg-pg-ardosia text-pg-branco"
        >
          <WifiOff
            size={48}
            aria-hidden="true"
            className="shrink-0 text-pg-down"
          />
          <div>
            <p className="m-0 text-2xl font-medium">
              Sem conexão com o servidor
            </p>
            <p className="text-lg text-pg-borda mt-2 mb-0">
              Sem conexão desde {since}.
            </p>
            <p className="text-lg text-pg-borda m-0">
              Nova tentativa em {seconds} s
            </p>
          </div>
        </div>
      </div>
      {children}
    </div>
  );
}
