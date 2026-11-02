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

// Campo com sugestões (datalist): digitar filtra a lista; UTC fica no topo.
// Mostra a hora de agora no fuso digitado, para conferir antes de salvar.
export function TimezoneCard({ value, onChange }) {
  const options = useMemo(() => timeZoneOptions(), []);
  const now = useNow();
  const tz = value.trim();
  const invalid = tz !== "" && !isValidTimeZone(tz);

  let preview = null;
  if (tz && !invalid) {
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
      <input
        id="timezone"
        type="text"
        list="timezone-options"
        autoComplete="off"
        spellCheck={false}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="UTC"
        aria-invalid={invalid || undefined}
        aria-describedby="timezone-help"
        className="w-full md:max-w-[360px] px-3 py-2.5 rounded-pg-control font-pg-mono text-pg-ui text-pg-ardosia bg-pg-branco border border-pg-borda aria-[invalid=true]:border-pg-down"
      />
      <datalist id="timezone-options">
        {options.map((z) => (
          <option key={z} value={z} />
        ))}
      </datalist>
      <p
        id="timezone-help"
        className={`mt-1.5 mb-0 text-pg-caption ${invalid ? "text-pg-down" : "text-pg-neblina"}`}
      >
        {invalid
          ? "Fuso desconhecido. Escolha um nome da lista, como America/Sao_Paulo."
          : (preview ?? "Vazio volta para UTC.")}
      </p>
    </Card>
  );
}
