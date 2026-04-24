import { useState, useEffect } from "react";
import { settingsApi } from "../api/settings";

export function useSettings() {
  const [settings, setSettings] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadSettings = async () => {
    try {
      setSettings(await settingsApi.get());
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const update = async (data) => {
    const updated = await settingsApi.update(data);
    setSettings(updated);
    return updated;
  };

  useEffect(() => {
    const initialLoad = setTimeout(() => {
      void loadSettings();
    }, 0);

    return () => clearTimeout(initialLoad);
  }, []);

  return { settings, loading, error, update };
}
