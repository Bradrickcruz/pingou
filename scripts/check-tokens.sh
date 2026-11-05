#!/usr/bin/env bash
# Trava do tema escuro: componentes e páginas usam só tokens de PAPEL
# (bg-pg-superficie, text-pg-texto, ...), nunca a cor crua (bg-pg-branco,
# text-pg-ardosia, ...) nem cor solta. Cor crua não troca no escuro e deixa a
# tela "furada". Ver DESIGN.md 3.6 e tasks/plano-dark-mode.md.
#
#   scripts/check-tokens.sh          # confere web/src
#
# Cores que podem aparecer cruas, por serem iguais nos dois temas:
#   pg-ciano  (símbolo, ícones, anel de foco) e pg-nuvem em texto sobre
#   superfície escura fixa (coluna lateral do login).
set -euo pipefail
cd "$(dirname "$0")/../web/src"

DIRS=(components pages)
FAIL=0

check() { # descrição, regex
  local out
  out=$(grep -rnP --include='*.jsx' --include='*.js' --exclude='*.test.js' "$2" "${DIRS[@]}" \
    | grep -v 'pages/Login.jsx:.*text-pg-nuvem' || true)
  if [ -n "$out" ]; then
    echo "✘ $1"; echo "$out" | sed 's/^/    /'; echo
    FAIL=1
  fi
}

# classe de cor crua (ignora pg-ardosia-clara, pg-up-claro etc. e os papéis)
RAW='(up|down|unknown)(-bg)?|branco|ardosia|ciano-fundo|ciano-nevoa|neblina|nuvem|borda'
UTIL='(bg|text|border(-[trblxy])?|accent|fill|stroke|ring|outline|divide|from|via|to|decoration|caret|shadow)'
check "classe de cor crua. Use o token de papel (veja tailwind.config.js)" \
  "(?<![\w-])((?:[\w\[\]-]+:)*)${UTIL}-pg-(${RAW})(?![\w-])"

check "cor hexadecimal solta" "#[0-9a-fA-F]{6}\b|#[0-9a-fA-F]{3}\b(?![\w-])|\[#[0-9a-fA-F]{3,8}\]"
check "rgb()/rgba()/hsl() solto" "\b(rgba?|hsla?)\("
check "variável de cor crua. Use var(--pg-cor-*)" "var\(--pg-(branco|ardosia|ciano-fundo|ciano-nevoa|neblina|nuvem|borda|up|down|unknown)(-bg)?\)"

if [ "$FAIL" -ne 0 ]; then
  echo "Troque pelos tokens de papel (DESIGN.md 3.6.2 e 3.6.6)."
  echo "Exceções: text-pg-nuvem na coluna lateral do login e pg-ciano (iguais nos dois temas)."
  exit 1
fi
echo "✔ tokens: nenhuma cor crua em components/ e pages/"
