import { useMemo } from "react";
import { Card } from "./Card";
import { useNow } from "../../hooks/useNow";
import {
  isValidTimeZone,
  timeZoneOptions,
  zonedParts,
  zoneOffsetLabel,
} from "../../lib/timezone";

const pad = (n) => String(n).padStart(2, "0");

// Select com a lista de fusos; UTC fica no topo.
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
        className="block mb-1.5 text-pg-ui font-medium text-pg-ardosia"
      >
        Fuso
      </label>
      <select
        id="timezone"
        value={tz}
        onChange={(e) => onChange(e.target.value)}
        aria-describedby="timezone-help"
        className="w-full md:max-w-[360px] px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda focus-visible:outline-2 focus-visible:outline focus-visible:outline-pg-ciano"
      >
        {options.map((z) => (
          <option key={z} value={z}>
            {z}
          </option>
        ))}
      </select>
      <p
        id="timezone-help"
        className={`mt-1.5 mb-0 text-pg-caption ${invalid ? "text-pg-down" : "text-pg-neblina"}`}
      >
        {invalid
          ? "Fuso desconhecido neste navegador. Escolha outro da lista."
          : preview}
      </p>
    </Card>
  );
}
