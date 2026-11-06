import { useMemo } from "react";
import { Card } from "./Card";
import { Combobox } from "../ui/Combobox";
import { useNow } from "../../hooks/useNow";
import {
  filterTimeZones,
  isValidTimeZone,
  timeZoneOptions,
  zonedParts,
  zoneOffsetLabel,
} from "../../lib/timezone";

const pad = (n) => String(n).padStart(2, "0");

// Campo digitável que filtra a lista de fusos; UTC fica no topo.
// Mostra a hora de agora no fuso digitado, para conferir antes de salvar.
export function TimezoneCard({ value, onChange }) {
  const now = useNow();
  const tz = value.trim() || "UTC";
  // ponytail: fuso salvo fora da lista do navegador continua selecionável
  const options = useMemo(() => {
    const list = timeZoneOptions();
    return list.includes(tz) ? list : [tz, ...list];
  }, [tz]);
  const invalid = !isValidTimeZone(tz);

  let preview = null;
  if (!invalid) {
    const p = zonedParts(new Date(now), tz);
    preview = `Agora: ${pad(p.hour)}:${pad(p.minute)} (${zoneOffsetLabel(new Date(now), tz)})`;
  }

  return (
    <Card
      title="Fuso horário"
      description="As datas do painel e dos alertas por e-mail e ntfy usam este fuso. A API e o webhook continuam em UTC."
    >
      <label
        htmlFor="timezone"
        className="block mb-1.5 text-pg-ui font-medium text-pg-texto"
      >
        Fuso
      </label>
      <Combobox
        id="timezone"
        value={tz}
        options={options}
        onChange={onChange}
        filter={filterTimeZones}
        emptyText="Nenhum fuso encontrado"
        aria-describedby="timezone-help"
      />
      <p
        id="timezone-help"
        className={`mt-1.5 mb-0 text-pg-caption ${invalid ? "text-pg-status-down" : "text-pg-texto-2"}`}
      >
        {invalid
          ? "Fuso desconhecido neste navegador. Escolha outro da lista."
          : preview}
      </p>
    </Card>
  );
}
