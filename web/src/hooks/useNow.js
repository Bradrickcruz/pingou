import { useEffect, useState } from "react";

// Relógio que atualiza a cada `ms`, para textos como "há 4 s".
export function useNow(ms = 1000) {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    const id = setInterval(() => setNow(Date.now()), ms);
    return () => clearInterval(id);
  }, [ms]);
  return now;
}
