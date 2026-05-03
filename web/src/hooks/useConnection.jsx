import { createContext, useContext, useState, useEffect } from "react";
import { client } from "../api/client";

const ConnectionContext = createContext();

export function ConnectionProvider({ children }) {
  const [online, setOnline] = useState(true);
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    let mounted = true;

    const checkConnection = async () => {
      if (!mounted) return;
      setChecking(true);
      try {
        await client.get("/monitors");
        if (mounted) setOnline(true);
      } catch {
        if (mounted) setOnline(false);
      } finally {
        if (mounted) setChecking(false);
      }
    };

    checkConnection();

    const interval = setInterval(checkConnection, 10000);

    return () => {
      mounted = false;
      clearInterval(interval);
    };
  }, []);

  return (
    <ConnectionContext.Provider value={{ online, checking }}>
      {children}
    </ConnectionContext.Provider>
  );
}

export function useConnection() {
  return useContext(ConnectionContext);
}