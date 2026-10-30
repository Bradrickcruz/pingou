// Barra fixa do mobile (< md) com Descartar/Salvar, acima da barra de abas.
// Só aparece com alterações; o Salvar envia o formulário pelo atributo form.
export function UnsavedBar({ formId, onDiscard, saving, error }) {
  return (
    <div
      role="status"
      className="md:hidden fixed inset-x-0 bottom-[calc(57px+env(safe-area-inset-bottom))] z-40 px-4 py-3 bg-pg-ardosia text-pg-branco"
    >
      {error && <p className="mb-2 text-pg-caption text-[#F3B5BB]">{error}</p>}
      <div className="flex items-center gap-2">
        <span className="flex-1 min-w-0 text-pg-ui">Alterações não salvas</span>
        <button
          type="button"
          onClick={onDiscard}
          disabled={saving}
          className="min-h-[44px] px-4 rounded-pg-control border border-pg-branco text-pg-ui font-medium disabled:opacity-50 focus-visible:outline-2 focus-visible:outline focus-visible:outline-offset-2 focus-visible:outline-pg-ciano"
        >
          Descartar
        </button>
        <button
          type="submit"
          form={formId}
          disabled={saving}
          className="min-h-[44px] px-4 rounded-pg-control bg-pg-branco text-pg-ardosia text-pg-ui font-medium disabled:opacity-50 focus-visible:outline-2 focus-visible:outline focus-visible:outline-offset-2 focus-visible:outline-pg-ciano"
        >
          {saving ? "Salvando…" : "Salvar"}
        </button>
      </div>
    </div>
  );
}
