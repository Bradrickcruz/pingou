import { Link } from "react-router-dom";
import { ChevronLeft } from "lucide-react";

// Barra do topo do detalhe no mobile (< md): "‹ Monitores" sempre volta ao
// painel (quem chega pelo link do alerta não tem histórico) e o menu à direita.
export function DetailTopBar({ menu }) {
  return (
    <div className="md:hidden sticky top-0 z-30 -mx-4 -mt-4 mb-4 flex items-center justify-between gap-2 px-1 bg-pg-branco border-b border-pg-borda">
      <Link
        to="/"
        className="inline-flex items-center gap-1 min-h-[48px] px-3 text-pg-ui font-medium text-pg-ardosia no-underline focus-visible:outline-2 focus-visible:outline focus-visible:-outline-offset-2 focus-visible:outline-pg-ciano"
      >
        <ChevronLeft size={20} aria-hidden="true" />
        Monitores
      </Link>
      {menu}
    </div>
  );
}
