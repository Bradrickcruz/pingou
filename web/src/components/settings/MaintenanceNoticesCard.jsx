import { useId } from "react";
import { Card } from "./Card";

// Interruptor global do aviso de manutenção. Entra no formulário de
// Configurações como os outros cards (salva com "Salvar alterações").
export function MaintenanceNoticesCard({ checked, onChange }) {
  const id = useId();
  return (
    <Card
      title="Aviso de manutenção"
      description="Avisa nos canais ligados quando uma janela agendada começa e termina. Pausa manual não avisa."
    >
      <label
        htmlFor={id}
        className="inline-flex items-center gap-3 min-h-[44px] md:min-h-0 text-pg-ui text-pg-texto cursor-pointer"
      >
        <input
          id={id}
          type="checkbox"
          role="switch"
          checked={checked}
          onChange={(e) => onChange(e.target.checked)}
          className="w-5 h-5 shrink-0 accent-pg-link"
        />
        Avisar quando uma manutenção começar e terminar
      </label>
    </Card>
  );
}
