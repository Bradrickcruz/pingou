#!/usr/bin/env bash
# Mede o banco com volume real: gera os dados (se faltarem), imprime os planos
# de consulta e roda os benchmarks. Resultado em .bench/results-<perfil>-<data>.md
#
#   scripts/dbbench.sh            # perfil típico (10 monitores, 60 s)
#   scripts/dbbench.sh limite     # perfil limite (100 monitores, 10 s; ~5 GB, vários minutos)
#
# Variáveis repassadas aos benchmarks: PINGOU_BENCH_REPS, PINGOU_BENCH_WRITES,
# PINGOU_BENCH_CONCURRENCY (ver internal/repository/dbbench_test.go).
set -euo pipefail
cd "$(dirname "$0")/.."

PROFILE="${1:-tipico}"
DB=".bench/${PROFILE}.db"
OUT=".bench/results-${PROFILE}-$(date -u +%Y%m%d-%H%M%S).md"
mkdir -p .bench

if [ ! -f "$DB" ]; then
  echo "gerando $DB..."
  go run ./scripts/seed --profile "$PROFILE" --out "$DB"
fi

# a cópia do banco vai para o TMPDIR: precisa de espaço igual ao tamanho do .db
export TMPDIR="${TMPDIR:-$PWD/.bench/tmp}"
mkdir -p "$TMPDIR"

{
  echo "# Benchmark do banco: perfil $PROFILE"
  echo
  echo "- data: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "- commit: $(git rev-parse --short HEAD 2>/dev/null || echo '?')"
  echo "- CPU: $(grep -m1 'model name' /proc/cpuinfo 2>/dev/null | cut -d: -f2- | xargs || uname -m) ($(nproc 2>/dev/null || echo '?') núcleos)"
  echo "- RAM: $(free -m 2>/dev/null | awk '/Mem:/ {print $2 " MB"}')"
  echo "- disco: $(df -h .bench | awk 'NR==2 {print $1 " (" $2 ")"}')"
  echo "- banco: $DB ($(du -h "$DB" | cut -f1))"
  echo
  echo "## Planos de consulta"
  PINGOU_BENCH_DB="$DB" go test ./internal/repository -run '^TestPrintQueryPlans$' -count=1 -v \
    | grep '^|'
  echo
  echo "## Tempos"
  PINGOU_BENCH_DB="$DB" go test ./internal/repository -run '^TestDBBench$' -count=1 -v -timeout 3h \
    | grep -v -e '^=== RUN' -e '^--- ' -e '^PASS' -e '^ok' -e '^[0-9]\{4\}/'
} | tee "$OUT"

echo
echo "salvo em $OUT"
