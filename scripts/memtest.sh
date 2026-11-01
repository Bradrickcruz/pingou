#!/usr/bin/env bash
# Prova da promessa do PRD: 100 monitores a cada 10s com < 50MB de RAM.
# Uso: scripts/memtest.sh [duração_em_segundos=3600]
# Requer o binário em bin/pingou (make build) e nada rodando na porta 8956.
set -euo pipefail
cd "$(dirname "$0")/.."

DUR=${1:-3600}
KEY=memtest-key
API=http://localhost:8956/api
OUT=$(mktemp -d)
LOG=$OUT/mem.log

cleanup() { kill "${PINGOU_PID:-}" "${TARGET_PID:-}" 2>/dev/null || true; }
trap cleanup EXIT

go build -o "$OUT/target" ./scripts/memtest/target.go
"$OUT/target" & TARGET_PID=$!
PINGOU_API_KEY=$KEY PINGOU_DATABASE_URL=$OUT/pingou.db ./bin/pingou >"$OUT/pingou.log" 2>&1 & PINGOU_PID=$!
until curl -sf -H "X-API-Key: $KEY" $API/info >/dev/null; do sleep 0.5; done

post() { curl -sf -X POST $API/monitors -H "X-API-Key: $KEY" -H 'Content-Type: application/json' -d "$1" >/dev/null; }
HTTPS_LEFT=${HTTPS_N:-30}   # monitores https://google.com no lugar de parte dos "ok" locais
for i in $(seq 1 100); do
  case $((i % 10)) in
    0)    post "{\"name\":\"tcp-$i\",\"type\":\"tcp\",\"config\":{\"host\":\"127.0.0.1\",\"port\":9000},\"interval_seconds\":10,\"timeout_seconds\":5,\"failure_threshold\":3}" ;;
    1)    post "{\"name\":\"slow-$i\",\"url\":\"http://127.0.0.1:9000/slow/$i\",\"interval_seconds\":10,\"timeout_seconds\":5,\"failure_threshold\":3}" ;;
    2)    post "{\"name\":\"fail-$i\",\"url\":\"http://127.0.0.1:9000/fail/$i\",\"interval_seconds\":10,\"timeout_seconds\":5,\"failure_threshold\":3}" ;;
    *)    if [ "$HTTPS_LEFT" -gt 0 ]; then
            HTTPS_LEFT=$((HTTPS_LEFT - 1))
            post "{\"name\":\"google-$i\",\"url\":\"https://google.com\",\"interval_seconds\":10,\"timeout_seconds\":5,\"failure_threshold\":3}"
          else
            post "{\"name\":\"ok-$i\",\"url\":\"http://127.0.0.1:9000/ok/$i\",\"interval_seconds\":10,\"timeout_seconds\":5,\"failure_threshold\":3}"
          fi ;;
  esac
done
# o 101º deve ser recusado
code=$(curl -s -o /dev/null -w '%{http_code}' -X POST $API/monitors -H "X-API-Key: $KEY" -H 'Content-Type: application/json' \
  -d '{"name":"extra","url":"http://127.0.0.1:9000/ok/x","interval_seconds":10,"timeout_seconds":5,"failure_threshold":3}')
curl -s -H "X-API-Key: $KEY" "$API/monitors?limit=100" | grep -o '"url":"https://google.com"' | wc -l | xargs echo "monitores google.com:"
echo "101º monitor -> HTTP $code (esperado 4xx)"

echo "medindo por ${DUR}s; log em $LOG"
end=$((SECONDS + DUR))
while [ $SECONDS -lt $end ]; do
  # carga extra a cada ciclo: dashboard (polling) e stats de 30d, como um usuário com a tela aberta
  curl -s -H "X-API-Key: $KEY" $API/monitors >/dev/null
  curl -s -H "X-API-Key: $KEY" "$API/incidents?limit=20" >/dev/null
  rss=$(awk '/VmRSS/{print $2}' /proc/$PINGOU_PID/status)
  hwm=$(awk '/VmHWM/{print $2}' /proc/$PINGOU_PID/status)
  thr=$(awk '/Threads/{print $2}' /proc/$PINGOU_PID/status)
  echo "$(date +%T) rss_kb=$rss peak_kb=$hwm threads=$thr" | tee -a "$LOG"
  sleep 10
done

peak=$(awk '/VmHWM/{print $2}' /proc/$PINGOU_PID/status)
echo "pico: $((peak / 1024)) MB (limite 50)"
[ "$peak" -lt $((50 * 1024)) ] && echo PASSOU || { echo FALHOU; exit 1; }
