// Card com título e descrição, igual nos três blocos de Configurações.
export function Card({ title, description, children }) {
  return (
    <section className="p-4 md:p-6 rounded-pg-card border border-pg-linha bg-pg-superficie">
      <h2 className="text-pg-h3 my-0">{title}</h2>
      <p className="mt-1 mb-4 text-pg-ui text-pg-texto-2">{description}</p>
      {children}
    </section>
  );
}
