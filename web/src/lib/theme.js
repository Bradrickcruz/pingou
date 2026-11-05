// Tema do painel (Claro, Escuro, Sistema). Vale por navegador: fica só no
// localStorage, nunca vai para a API. Funções puras e sem React, para o hook
// e os testes usarem a mesma regra.
//
// ATENÇÃO: o index.html tem um script inline com a mesma regra (roda antes da
// primeira pintura, para não piscar o claro). Mexeu aqui, mexa lá. O teste em
// theme.test.js confere se os dois concordam.

export const THEME_KEY = "pingou_theme";
export const PREFERENCES = ["light", "dark", "system"];
export const DEFAULT_PREFERENCE = "system";

// Trava de segurança. Com "Sistema" como padrão, ligar o escuro antes de todas
// as telas estarem migradas (T4 a T8 do plano) quebraria o painel de quem usa
// o sistema no escuro, sem a pessoa ter escolhido nada. Com false, o tema
// efetivo é sempre o claro. Virar true junto com o bloco "Aparência", e trocar
// o `ON` do script do index.html no mesmo commit.
export const THEME_ENABLED = false;

export function normalize(value) {
  return PREFERENCES.includes(value) ? value : DEFAULT_PREFERENCE;
}

// localStorage pode lançar erro (aba anônima, dados bloqueados). Nesse caso o
// tema segue funcionando, só não lembra a escolha.
function defaultStorage() {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

export function getPreference(storage = defaultStorage()) {
  try {
    return normalize(storage?.getItem(THEME_KEY));
  } catch {
    return DEFAULT_PREFERENCE;
  }
}

export function savePreference(pref, storage = defaultStorage()) {
  const value = normalize(pref);
  try {
    storage?.setItem(THEME_KEY, value);
  } catch {
    // sem armazenamento: vale só até recarregar
  }
  return value;
}

export function systemIsDark(matchMedia = globalThis.matchMedia) {
  try {
    return !!matchMedia?.("(prefers-color-scheme: dark)").matches;
  } catch {
    return false;
  }
}

// "light" | "dark" que de fato vai para a tela.
export function resolve(pref, isSystemDark, enabled = THEME_ENABLED) {
  if (!enabled) return "light";
  const p = normalize(pref);
  if (p === "system") return isSystemDark ? "dark" : "light";
  return p;
}

export function apply(resolved, root = document.documentElement) {
  const theme = resolved === "dark" ? "dark" : "light";
  root.setAttribute("data-theme", theme);
  root.style.colorScheme = theme;
  return theme;
}
