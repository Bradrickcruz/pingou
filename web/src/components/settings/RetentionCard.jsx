import { CircleAlert } from "lucide-react";
import { formatBytes } from "../../lib/format";
import { Card } from "./Card";

const OPTIONS = [7, 14, 30, 60, 90];

// Rádios escondidos dentro de chips: setas do teclado trocam a opção.
// savedValue: valor salvo; escolher um menor avisa o que será apagado.
export function RetentionCard({ value, savedValue, onChange, databaseSize }) {
  return (
    <Card
      title="Retenção de dados"
      description="Por quanto tempo o Pingou guarda os checks. Checks mais antigos que isso são apagados automaticamente; os incidentes são mantidos. Prazos maiores deixam o banco maior."
    >
      <fieldset className="min-w-0 border-0 p-0 m-0">
        <legend className="sr-only">Retenção (dias)</legend>
        <div className="flex md:flex-wrap gap-2">
          {OPTIONS.map((days) => {
            const selected = value === days;
            return (
              <label
                key={days}
                className={`flex-1 md:flex-none inline-flex items-center justify-center min-h-[44px] md:min-h-0 px-2 md:px-4 py-2 rounded-full border text-pg-ui cursor-pointer transition-colors duration-150 has-[:focus-visible]:outline-2 has-[:focus-visible]:outline has-[:focus-visible]:outline-pg-ciano ${
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
                {days}
                <span className="hidden md:inline">&nbsp;dias</span>
              </label>
            );
          })}
        </div>
        <p
          aria-hidden="true"
          className="md:hidden mt-1.5 mb-0 text-pg-caption text-pg-neblina"
        >
          dias
        </p>
      </fieldset>
      {savedValue != null && value < savedValue && (
        <p
          role="status"
          className="flex gap-2.5 items-start mt-4 mb-0 px-4 py-3 rounded-pg-control border border-pg-neblina bg-pg-nuvem text-pg-caption text-pg-ardosia"
        >
          <CircleAlert
            size={16}
            aria-hidden="true"
            className="shrink-0 mt-px"
          />
          <span>
            Ao salvar, checks com mais de{" "}
            <strong className="font-medium">{value} dias</strong> serão
            apagados. Isso não pode ser desfeito.
          </span>
        </p>
      )}
      {databaseSize != null && (
        <p className="mt-4 mb-0 text-pg-caption text-pg-neblina">
          Tamanho atual do banco: {formatBytes(databaseSize)}
        </p>
      )}
    </Card>
  );
}
