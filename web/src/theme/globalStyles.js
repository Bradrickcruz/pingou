export function injectGlobalStyles() {
  const style = document.createElement("style");
  style.textContent = `
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

    body {
      background: var(--pg-nuvem);
      color: var(--pg-ardosia);
      font-family: Inter, system-ui, sans-serif;
      font-size: 14px;
      line-height: 1.6;
      -webkit-font-smoothing: antialiased;
    }

    ::-webkit-scrollbar { width: 6px; }
    ::-webkit-scrollbar-track { background: var(--pg-nuvem); }
    ::-webkit-scrollbar-thumb { background: var(--pg-borda); border-radius: 3px; }

    a { color: inherit; text-decoration: none; }
    button { cursor: pointer; font: inherit; }
    input, textarea, select {
      font: inherit;
      background: var(--pg-branco);
      color: var(--pg-ardosia);
      border: 1px solid var(--pg-borda);
      border-radius: 8px;
      padding: 8px 12px;
      width: 100%;
    }
    input:focus-visible, textarea:focus-visible, select:focus-visible {
      border-color: var(--pg-ciano);
      outline: 2px solid var(--pg-ciano);
      outline-offset: 1px;
    }
  `;
  document.head.appendChild(style);
}
