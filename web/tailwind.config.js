/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Light mode colors
        text: 'var(--text)',
        'text-h': 'var(--text-h)',
        bg: 'var(--bg)',
        border: 'var(--border)',
        'code-bg': 'var(--code-bg)',
        accent: 'var(--accent)',
        'accent-bg': 'var(--accent-bg)',
        'accent-border': 'var(--accent-border)',
        'social-bg': 'var(--social-bg)',

        // Pingou design system (DESIGN.md) — prefixed to avoid colliding
        // with the tokens above during the gradual migration
        pg: {
          ciano: 'var(--pg-ciano)',
          'ciano-fundo': 'var(--pg-ciano-fundo)',
          'ciano-nevoa': 'var(--pg-ciano-nevoa)',
          ardosia: 'var(--pg-ardosia)',
          neblina: 'var(--pg-neblina)',
          nuvem: 'var(--pg-nuvem)',
          branco: 'var(--pg-branco)',
          borda: 'var(--pg-borda)',
          up: 'var(--pg-up)',
          'up-bg': 'var(--pg-up-bg)',
          down: 'var(--pg-down)',
          'down-bg': 'var(--pg-down-bg)',
          unknown: 'var(--pg-unknown)',
          'unknown-bg': 'var(--pg-unknown-bg)',
        },
      },
      fontFamily: {
        sans: ['system-ui', 'Segoe UI', 'Roboto', 'sans-serif'],
        heading: ['system-ui', 'Segoe UI', 'Roboto', 'sans-serif'],
        mono: ['ui-monospace', 'Consolas', 'monospace'],

        'pg-heading': ['Rubik', 'system-ui', 'sans-serif'],
        'pg-sans': ['Inter', 'system-ui', 'sans-serif'],
        'pg-mono': ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace'],
      },
      fontSize: {
        'base': '18px',
        'h1': '56px',
        'h2': '24px',

        'pg-h1': ['32px', { lineHeight: '1.2' }],
        'pg-h2': ['24px', { lineHeight: '1.25' }],
        'pg-h3': ['20px', { lineHeight: '1.3' }],
        'pg-body': ['16px', { lineHeight: '1.6' }],
        'pg-ui': ['14px', { lineHeight: '1.5' }],
        'pg-caption': ['12px', { lineHeight: '1.4' }],
      },
      spacing: {
        'spacer': '88px',
      },
      boxShadow: {
        'custom': 'var(--shadow)',
      },
      borderRadius: {
        'sm': '4px',
        'md': '5px',

        'pg-control': '8px',
        'pg-card': '12px',
      },
    },
  },
  plugins: [],
}