import { tokens as t } from "./tokens";

export function injectGlobalStyles() {
  const style = document.createElement("style");
  style.textContent = `
    @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap');

    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    body {
      background: ${t.colors.bg};
      color: ${t.colors.textPrimary};
      font-family: ${t.font.sans};
      font-size: 14px;
      line-height: 1.6;
      -webkit-font-smoothing: antialiased;
    }

    ::-webkit-scrollbar { width: 6px; }
    ::-webkit-scrollbar-track { background: ${t.colors.bg}; }
    ::-webkit-scrollbar-thumb { background: ${t.colors.border}; border-radius: 3px; }

    a { color: inherit; text-decoration: none; }
    button { cursor: pointer; border: none; background: none; font: inherit; }
    input, textarea, select {
      font: inherit;
      background: ${t.colors.surfaceAlt};
      color: ${t.colors.textPrimary};
      border: 1px solid ${t.colors.border};
      border-radius: ${t.radius.sm};
      padding: 8px 12px;
      outline: none;
      width: 100%;
    }
    input:focus, textarea:focus, select:focus {
      border-color: ${t.colors.primary};
    }
  `;
  document.head.appendChild(style);
}
