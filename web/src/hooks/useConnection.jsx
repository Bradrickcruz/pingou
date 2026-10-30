import { useState, useEffect } from "react";
import { client } from "../api/client";
import { ConnectionContext } from "./connectionContext";

const RETRY_MS = 10000;

export function ConnectionProvider({ children, onUnauthorized }) {
  const [online, setOnline] = useState(true);
  const [checking, setChecking] = useState(true);
  const [offlineSince, setOfflineSince] = useState(null);
  const [nextCheckAt, setNextCheckAt] = useState(null);

  useEffect(() => {
    let mounted = true;

    const checkConnection = async () => {
      if (!mounted) return;
      setChecking(true);
      setNextCheckAt(Date.now() + RETRY_MS);
      try {
        await client.get("/monitors");
        if (mounted) {
          setOnline(true);
          setOfflineSince(null);
        }
      } catch (e) {
        // chave inválida/revogada: reconectar não adianta, volta pro login
        if (e.status === 401) return onUnauthorized?.();
        if (mounted) {
          setOnline(false);
          setOfflineSince((since) => since ?? Date.now());
        }
      } finally {
        if (mounted) setChecking(false);
      }
    };

    checkConnection();

    const interval = setInterval(checkConnection, RETRY_MS);

    return () => {
      mounted = false;
      clearInterval(interval);
    };
  }, [onUnauthorized]);

  return (
    <ConnectionContext.Provider
      value={{ online, checking, offlineSince, nextCheckAt }}
    >
      {children}
    </ConnectionContext.Provider>
  );
}
