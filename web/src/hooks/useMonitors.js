import { useState, useEffect, useCallback } from "react";
import { monitorsApi } from "../api/monitors";

export function useMonitors() {
  const [monitors, setMonitors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetch = useCallback(async () => {
    try {
      setLoading(true);
      const res = await monitorsApi.list({ limit: 100 });
      setMonitors(res.data);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetch();
    const interval = setInterval(fetch, 15_000); // polling a cada 15s
    return () => clearInterval(interval);
  }, [fetch]);

  return { monitors, loading, error, refetch: fetch };
}
