// Card com título e descrição, igual nos três blocos de Configurações.
export function Card({ title, description, children }) {
  return (
    <section className="p-6 rounded-pg-card border border-pg-borda bg-pg-branco">
      <h2 className="text-pg-h3 my-0">{title}</h2>
      <p className="mt-1 mb-4 text-pg-ui text-pg-neblina">{description}</p>
      {children}
    </section>
  );
}
