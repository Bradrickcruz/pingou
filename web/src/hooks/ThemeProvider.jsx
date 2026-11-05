import { useCallback, useEffect, useMemo, useState } from "react";
import {
  THEME_KEY,
  apply,
  getPreference,
  resolve,
  savePreference,
  systemIsDark,
} from "../lib/theme";
import { ThemeContext } from "./themeContext";

// Fica na raiz do app (inclusive no Login) para que, em "Sistema", o painel
// acompanhe o sistema operacional ao vivo em qualquer tela. Também acompanha
// trocas feitas em outra aba do mesmo navegador (evento `storage`).
export function ThemeProvider({ children }) {
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

  const value = useMemo(
    () => ({ preference, resolved, setPreference }),
    [preference, resolved, setPreference],
  );

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
}
