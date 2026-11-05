/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["selector", '[data-theme="dark"]'],
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        // Light mode colors
        text: "var(--text)",
        "text-h": "var(--text-h)",
        bg: "var(--bg)",
        border: "var(--border)",
        "code-bg": "var(--code-bg)",
        accent: "var(--accent)",
        "accent-bg": "var(--accent-bg)",
        "accent-border": "var(--accent-border)",
        "social-bg": "var(--social-bg)",

        // Pingou design system (DESIGN.md) — prefixed to avoid colliding
        // with the tokens above during the gradual migration
        pg: {
          ciano: "var(--pg-ciano)",
          "ciano-fundo": "var(--pg-ciano-fundo)",
          "ciano-nevoa": "var(--pg-ciano-nevoa)",
          ardosia: "var(--pg-ardosia)",
          neblina: "var(--pg-neblina)",
          nuvem: "var(--pg-nuvem)",
          branco: "var(--pg-branco)",
          borda: "var(--pg-borda)",
          up: "var(--pg-up)",
          "up-bg": "var(--pg-up-bg)",
          down: "var(--pg-down)",
          "down-bg": "var(--pg-down-bg)",
          unknown: "var(--pg-unknown)",
          "unknown-bg": "var(--pg-unknown-bg)",

          // Tokens de papel (DESIGN.md 3.6.2 / 3.6.6): trocam com o tema.
          // Componentes devem usar estes, não as cores cruas acima.
          fundo: "var(--pg-cor-fundo)",
          superficie: "var(--pg-cor-superficie)",
          linha: "var(--pg-cor-linha)",
          texto: "var(--pg-cor-texto)",
          "texto-2": "var(--pg-cor-texto-2)",
          link: "var(--pg-cor-link)",
          "borda-input": "var(--pg-cor-borda-input)",
          primario: "var(--pg-cor-primario)",
          "primario-texto": "var(--pg-cor-primario-texto)",
          "tabela-cab": "var(--pg-cor-tabela-cab)",
          "tabela-cab-texto": "var(--pg-cor-tabela-cab-texto)",
          destaque: "var(--pg-cor-destaque)",
          hover: "var(--pg-cor-hover)",
          ativo: "var(--pg-cor-ativo)",
          "ativo-texto": "var(--pg-cor-ativo-texto)",
          "sobre-escuro": "var(--pg-cor-sobre-escuro)",
          "sobre-escuro-2": "var(--pg-cor-sobre-escuro-2)",
          tooltip: "var(--pg-cor-tooltip)",
          flutuante: "var(--pg-cor-flutuante)",
          "flutuante-borda": "var(--pg-cor-flutuante-borda)",
          "flutuante-erro": "var(--pg-cor-flutuante-erro)",
          codigo: "var(--pg-cor-codigo)",
          "login-lateral": "var(--pg-cor-login-lateral)",
          veu: "var(--pg-cor-veu)",
          "down-borda": "var(--pg-cor-down-borda)",
          destrutivo: "var(--pg-cor-destrutivo)",
          "destrutivo-texto": "var(--pg-cor-destrutivo-texto)",
          "status-up": "var(--pg-cor-status-up)",
          "status-up-bg": "var(--pg-cor-status-up-bg)",
          "status-down": "var(--pg-cor-status-down)",
          "status-down-bg": "var(--pg-cor-status-down-bg)",
          "status-unknown": "var(--pg-cor-status-unknown)",
          "status-unknown-bg": "var(--pg-cor-status-unknown-bg)",
        },
      },
      fontFamily: {
        sans: ["system-ui", "Segoe UI", "Roboto", "sans-serif"],
        heading: ["system-ui", "Segoe UI", "Roboto", "sans-serif"],
        mono: ["ui-monospace", "Consolas", "monospace"],

        "pg-heading": ["Rubik", "system-ui", "sans-serif"],
        "pg-sans": ["Inter", "system-ui", "sans-serif"],
        "pg-mono": [
          "ui-monospace",
          "SFMono-Regular",
          "Menlo",
          "Consolas",
          "monospace",
        ],
      },
      fontSize: {
        base: "18px",
        h1: "56px",
        h2: "24px",

        "pg-h1": ["32px", { lineHeight: "1.2" }],
        "pg-h2": ["24px", { lineHeight: "1.25" }],
        "pg-h3": ["20px", { lineHeight: "1.3" }],
        "pg-body": ["16px", { lineHeight: "1.6" }],
        "pg-ui": ["14px", { lineHeight: "1.5" }],
        "pg-caption": ["12px", { lineHeight: "1.4" }],
      },
      spacing: {
        spacer: "88px",
      },
      boxShadow: {
        custom: "var(--shadow)",
        "pg-flutuante": "var(--pg-sombra-flutuante)",
        "pg-flutuante-forte": "var(--pg-sombra-flutuante-forte)",
      },
      borderRadius: {
        sm: "4px",
        md: "5px",

        "pg-control": "8px",
        "pg-card": "12px",
      },
    },
  },
  plugins: [],
};
