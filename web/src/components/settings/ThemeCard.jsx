import { Card } from "./Card";
import { SegmentedTabs } from "../ui/SegmentedTabs";

const OPTIONS = [
  { key: "light", label: "Claro" },
  { key: "dark", label: "Escuro" },
  { key: "system", label: "Sistema" },
];

// Tema do painel. Vale só neste navegador e aplica na hora: não entra no
// "Salvar" das Configurações, que guarda valores globais na API.
export function ThemeCard({ preference, onChange }) {
  return (
    <Card title="Aparência" description="Vale só neste navegador.">
      <SegmentedTabs
        label="Tema"
        options={OPTIONS}
        value={preference}
        onChange={onChange}
      />
      <p className="mt-3 mb-0 text-pg-caption text-pg-texto-2">
        Sistema segue o tema do seu aparelho.
      </p>
    </Card>
  );
}
