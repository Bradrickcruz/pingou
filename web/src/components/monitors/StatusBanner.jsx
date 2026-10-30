import { X, Check, CircleHelp } from "lucide-react";
import { plural } from "../../lib/format";

function joinNames(names) {
  if (names.length <= 1) return names[0] ?? "";
  if (names.length > 3) return `${names.slice(0, 2).join(", ")} e mais ${names.length - 2}`;
  return `${names.slice(0, -1).join(", ")} e ${names.at(-1)}`;
}

// Faixa de situação do topo do painel: frase principal + contagens.
export function StatusBanner({ monitors, counts }) {
  const down = monitors.filter((m) => m.enabled && m.current_state === "DOWN");
  const active = monitors.length - counts.PAUSED;

  let tone, Icon, title, detail;
  if (down.length > 0) {
    tone = { box: "bg-pg-down-bg border-pg-down", icon: "bg-pg-down" };
    Icon = X;
    title = down.length === 1 ? "1 monitor fora do ar" : `${down.length} monitores fora do ar`;
    detail = `${joinNames(down.map((m) => m.name))} ${down.length === 1 ? "não está respondendo" : "não estão respondendo"}.`;
  } else if (counts.UNKNOWN > 0 && counts.UP === 0) {
    tone = { box: "bg-pg-branco border-pg-borda", icon: "bg-pg-unknown" };
    Icon = CircleHelp;
    title = "Aguardando os primeiros checks";
    detail = "Os monitores aparecem como no ar ou fora do ar depois do primeiro resultado.";
  } else {
    tone = { box: "bg-pg-up-bg border-pg-borda", icon: "bg-pg-up" };
    Icon = Check;
    title = active === 0 ? "Todos os monitores estão pausados" : "Tudo no ar";
    detail =
      active === 0
        ? "Retome um monitor para voltar a verificar."
        : `${plural(active, "monitor ativo respondendo", "monitores ativos respondendo")} normalmente.`;
  }

  const stats = [
    { n: counts.UP, label: "no ar" },
    { n: counts.DOWN, label: "fora do ar", strong: counts.DOWN > 0 },
    { n: counts.UNKNOWN, label: counts.UNKNOWN === 1 ? "desconhecido" : "desconhecidos" },
    { n: counts.PAUSED, label: counts.PAUSED === 1 ? "pausado" : "pausados" },
  ];

  // linha curta do mobile: "3 no ar · 2 fora · 1 desconhecido · 1 pausado"
  const short = [
    { n: counts.UP, label: "no ar" },
    { n: counts.DOWN, label: "fora", strong: true },
    { n: counts.UNKNOWN, label: counts.UNKNOWN === 1 ? "desconhecido" : "desconhecidos" },
    { n: counts.PAUSED, label: counts.PAUSED === 1 ? "pausado" : "pausados" },
  ].filter((s) => s.n > 0);

  return (
    <section
      aria-live="polite"
      className={`flex flex-col gap-3 md:gap-4 lg:flex-row lg:items-center lg:justify-between p-4 md:p-5 mb-4 md:mb-7 rounded-pg-card border ${tone.box}`}
    >
      <div className="flex items-center gap-4 min-w-0">
        <span className={`flex-shrink-0 inline-flex items-center justify-center w-10 h-10 rounded-full text-pg-branco ${tone.icon}`}>
          <Icon size={20} aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <h2 className="text-pg-h3 font-medium text-pg-ardosia my-0">{title}</h2>
          <p className="hidden md:block text-pg-ui text-pg-ardosia">{detail}</p>
        </div>
      </div>
      <p className="md:hidden text-pg-caption text-pg-ardosia">
        {short.map((s, i) => (
          <span key={s.label}>
            {i > 0 && " · "}
            <span className={s.strong ? "text-pg-down" : undefined}>
              <span className="tabular-nums">{s.n}</span> {s.label}
            </span>
          </span>
        ))}
      </p>
      <ul className="hidden md:flex flex-wrap gap-x-5 gap-y-1 text-pg-ui text-pg-ardosia list-none m-0 p-0">
        {stats.map((s) => (
          <li key={s.label} className={s.strong ? "text-pg-down" : undefined}>
            <span className="tabular-nums">{s.n}</span> {s.label}
          </li>
        ))}
      </ul>
    </section>
  );
}
