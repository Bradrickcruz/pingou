import { useCallback, useEffect, useState } from "react";
import {
  THEME_KEY,
  apply,
  getPreference,
  resolve,
  savePreference,
  systemIsDark,
} from "../lib/theme";

// Estado do tema para a tela de Configurações. Em "Sistema", acompanha o
// sistema operacional ao vivo; e acompanha trocas feitas em outra aba do mesmo
// navegador (evento `storage`).
export function useTheme() {
  const [preference, setPreferenceState] = useState(getPreference);
  const [isSystemDark, setIsSystemDark] = useState(() => systemIsDark());

  const resolved = resolve(preference, isSystemDark);

  useEffect(() => {
    apply(resolved);
  }, [resolved]);

  useEffect(() => {
    const mq = window.matchMedia?.("(prefers-color-scheme: dark)");
    if (!mq) return;
    const onChange = (e) => setIsSystemDark(e.matches);
    mq.addEventListener("change", onChange);
    return () => mq.removeEventListener("change", onChange);
  }, []);

  useEffect(() => {
    const onStorage = (e) => {
      if (e.key === THEME_KEY || e.key === null) {
        setPreferenceState(getPreference());
      }
    };
    window.addEventListener("storage", onStorage);
    return () => window.removeEventListener("storage", onStorage);
  }, []);

  const setPreference = useCallback((pref) => {
    setPreferenceState(savePreference(pref));
  }, []);

  return { preference, resolved, setPreference };
}
