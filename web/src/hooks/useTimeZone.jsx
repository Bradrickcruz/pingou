import { useCallback, useEffect, useState } from "react";
import { settingsApi } from "../api/settings";
import { isValidTimeZone, setTimeZone as applyTimeZone } from "../lib/timezone";
import { TimeZoneContext } from "./timeZoneContext";

// Carrega o fuso das Configurações antes de mostrar o painel, para nenhuma
// data aparecer primeiro em UTC e depois mudar. Se a API falhar, segue em UTC
// (a tela de conexão perdida cuida do resto).
// Fuso nunca configurado: salva o do navegador uma vez (e-mail e ntfy são
// montados no servidor, então precisa persistir). Sem fuso válido, fica em UTC.
async function loadTimeZone() {
  const s = await settingsApi.get();
  if (s.timezone_set) return s.timezone;
  const browser = Intl.DateTimeFormat().resolvedOptions().timeZone;
  if (!isValidTimeZone(browser)) return s.timezone;
  try {
    return (await settingsApi.update({ timezone: browser })).timezone;
  } catch {
    return s.timezone;
  }
}
export function TimeZoneProvider({ children }) {
  const [timeZone, setState] = useState(null);

  // troca o fuso dos formatadores antes do estado: quem renderizar depois já
  // usa o novo. O fuso só muda na tela de Configurações, que não mostra datas;
  // as outras telas montam de novo ao navegar e pegam o fuso novo.
  const setTimeZone = useCallback((tz) => setState(applyTimeZone(tz)), []);

  useEffect(() => {
    let mounted = true;
    loadTimeZone()
      .then((tz) => mounted && setTimeZone(tz))
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
