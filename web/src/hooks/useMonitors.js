import { useState, useEffect, useCallback } from "react";
import { monitorsApi } from "../api/monitors";

export function useMonitors() {
  const [monitors, setMonitors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isFetching, setIsFetching] = useState(false);
  const [updatedAt, setUpdatedAt] = useState(null);

  const loadMonitors = useCallback(async () => {
    try {
      const res = await monitorsApi.list({ limit: 100 });
      setMonitors(res.data);
      setError(null);
      setUpdatedAt(Date.now());
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
      setIsFetching(false);
    }
  }, []);

  const refetch = useCallback(async () => {
    setIsFetching(true);
    await loadMonitors();
  }, [loadMonitors]);

  useEffect(() => {
    const initialLoad = setTimeout(() => {
      void loadMonitors();
    }, 0);
    const interval = setInterval(() => {
      void loadMonitors();
    }, 15_000); // polling a cada 15s
    return () => {
      clearTimeout(initialLoad);
      clearInterval(interval);
    };
  }, [loadMonitors]);

  return { monitors, loading, error, refetch, isFetching, updatedAt };
}
