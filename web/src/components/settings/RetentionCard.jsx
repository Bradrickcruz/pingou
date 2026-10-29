import { formatBytes } from "../../lib/format";
import { Card } from "./Card";

const OPTIONS = [7, 14, 30, 60, 90];

// Rádios escondidos dentro de chips: setas do teclado trocam a opção.
export function RetentionCard({ value, onChange, databaseSize }) {
  return (
    <Card
      title="Retenção de dados"
      description="Por quanto tempo o Pingou guarda checks e incidentes. Checks e incidentes resolvidos mais antigos que isso são apagados automaticamente. Incidentes abertos nunca. Prazos maiores deixam o banco maior."
    >
      <fieldset className="min-w-0 border-0 p-0 m-0">
        <legend className="sr-only">Retenção (dias)</legend>
        <div className="flex flex-wrap gap-2">
          {OPTIONS.map((days) => {
            const selected = value === days;
            return (
              <label
                key={days}
                className={`px-4 py-2 rounded-full border text-pg-ui cursor-pointer transition-colors duration-150 has-[:focus-visible]:outline-2 has-[:focus-visible]:outline has-[:focus-visible]:outline-pg-ciano ${
                  selected
                    ? "bg-pg-ciano-nevoa border-pg-ciano-fundo text-pg-ciano-fundo"
                    : "bg-pg-branco border-pg-borda text-pg-ardosia hover:bg-pg-nuvem"
                }`}
              >
                <input
                  type="radio"
                  name="retention_days"
                  value={days}
                  checked={selected}
                  onChange={() => onChange(days)}
                  className="sr-only"
                />
                {days} dias
              </label>
            );
          })}
        </div>
      </fieldset>
      {databaseSize != null && (
        <p className="mt-4 mb-0 text-pg-caption text-pg-neblina">
          Tamanho atual do banco: {formatBytes(databaseSize)}
        </p>
      )}
    </Card>
  );
}
