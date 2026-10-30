import simbolo from "../../assets/brand/pingou-simbolo.svg";

// Cabeçalho das páginas no mobile (< md): símbolo, título e ações da página.
// No desktop some; a página mostra o próprio cabeçalho.
export function MobileHeader({ title, subtitle, actions }) {
  return (
    <header className="md:hidden sticky top-0 z-30 -mx-4 -mt-4 mb-4 flex items-center gap-3 px-4 py-3 bg-pg-branco border-b border-pg-borda">
      <img
        src={simbolo}
        alt=""
        aria-hidden="true"
        className="w-8 h-8 shrink-0"
      />
      <div className="flex-1 min-w-0">
        <h1 className="text-pg-h3 truncate">{title}</h1>
        {subtitle && (
          <p className="text-pg-caption text-pg-neblina truncate">{subtitle}</p>
        )}
      </div>
      {actions && (
        <div className="flex items-center gap-2 shrink-0">{actions}</div>
      )}
    </header>
  );
}
