import { useState, useEffect, useCallback } from "react";
import { settingsApi } from "../api/settings";

// withInfo: busca também /api/info (versão e tamanho do banco). Se /info falhar,
// a tela segue sem esses dados.
export function useSettings({ withInfo = false } = {}) {
  const [settings, setSettings] = useState(null);
  const [info, setInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    const [s, i] = await Promise.allSettled([settingsApi.get(), withInfo ? settingsApi.info() : null]);
    if (s.status === "fulfilled") setSettings(s.value);
    else setError(s.reason.message);
    if (i.status === "fulfilled") setInfo(i.value);
    setLoading(false);
  }, [withInfo]);

  const update = async (data) => {
    const updated = await settingsApi.update(data);
    setSettings(updated);
    return updated;
  };

  useEffect(() => {
    const initialLoad = setTimeout(() => {
      void load();
    }, 0);

    return () => clearTimeout(initialLoad);
  }, [load]);

  return { settings, info, loading, error, update, reload: load };
}
