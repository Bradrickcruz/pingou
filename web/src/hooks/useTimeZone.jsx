import { useCallback, useEffect, useState } from "react";
import { settingsApi } from "../api/settings";
import { setTimeZone as applyTimeZone } from "../lib/timezone";
import { TimeZoneContext } from "./timeZoneContext";

// Carrega o fuso das Configurações antes de mostrar o painel, para nenhuma
// data aparecer primeiro em UTC e depois mudar. Se a API falhar, segue em UTC
// (a tela de conexão perdida cuida do resto).
export function TimeZoneProvider({ children }) {
  const [timeZone, setState] = useState(null);

  // troca o fuso dos formatadores antes do estado: quem renderizar depois já
  // usa o novo. O fuso só muda na tela de Configurações, que não mostra datas;
  // as outras telas montam de novo ao navegar e pegam o fuso novo.
  const setTimeZone = useCallback((tz) => setState(applyTimeZone(tz)), []);

  useEffect(() => {
    let mounted = true;
    settingsApi
      .get()
      .then((s) => mounted && setTimeZone(s.timezone))
      .catch(() => mounted && setTimeZone(null));
    return () => {
      mounted = false;
    };
  }, [setTimeZone]);

  if (timeZone == null) return null;

  return (
    <TimeZoneContext.Provider value={{ timeZone, setTimeZone }}>
      {children}
    </TimeZoneContext.Provider>
  );
}
