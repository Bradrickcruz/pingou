#!/usr/bin/env bash
# Pico de RAM do Pingou com um banco grande aberto e o painel pesado sendo
# consultado: /stats?range=30d e a lista de monitores, sem parar.
#
#   scripts/dbmem.sh .bench/limite.db [segundos]   # padrão: 120 s de consultas
#
# Trabalha numa CÓPIA do banco (em .bench/tmp). Se a cópia estiver num schema
# antigo, o app aplica as migrations ao subir, e o tempo aparece no log.
# Requisitos: bin/pingou compilado na versão atual, porta 8956 livre, Linux.
set -euo pipefail
cd "$(dirname "$0")/.."

SRC="${1:?uso: scripts/dbmem.sh <banco.db> [segundos]}"
SECS="${2:-120}"
KEY=dbmem-key
API=http://localhost:8956/api
OUT=.bench/tmp/dbmem
mkdir -p "$OUT"
rm -f "$OUT"/pingou.db*
cp "$SRC" "$OUT/pingou.db"

cleanup() { kill "${PID:-}" 2>/dev/null || true; }
trap cleanup EXIT

start=$(date +%s)
PINGOU_API_KEY=$KEY PINGOU_DATABASE_URL=$OUT/pingou.db ./bin/pingou >"$OUT/pingou.log" 2>&1 &
PID=$!
until curl -sf http://localhost:8956/healthz >/dev/null; do
  kill -0 $PID 2>/dev/null || { cat "$OUT/pingou.log"; exit 1; }
  sleep 0.5
done
echo "subiu em $(($(date +%s) - start)) s"
grep -i "migration" "$OUT/pingou.log" || true

ids=$(curl -sf -H "X-API-Key: $KEY" "$API/monitors?limit=3" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
end=$(($(date +%s) + SECS))
n=0
while [ "$(date +%s)" -lt "$end" ]; do
  for id in $ids; do
    curl -sf -o /dev/null -H "X-API-Key: $KEY" "$API/monitors/$id/stats?range=30d" || echo "falhou: /stats $id"
  done
  curl -sf -o /dev/null -H "X-API-Key: $KEY" "$API/monitors?limit=100" || echo "falhou: lista"
  n=$((n + 1))
done

rss=$(awk '/VmRSS/{print $2}' /proc/$PID/status)
hwm=$(awk '/VmHWM/{print $2}' /proc/$PID/status)
echo "ciclos: $n (3× /stats 30d + lista) em $SECS s"
echo "rss: $((rss / 1024)) MB · pico: $((hwm / 1024)) MB (meta < 50 MB)"
