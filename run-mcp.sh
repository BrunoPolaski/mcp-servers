#!/usr/bin/env bash
# Sobe os 4 MCP servers nativos (go run), usando o Postgres do compose em 127.0.0.1:5433.
# Uso: ./run-mcp.sh [start|stop]
set -euo pipefail
cd "$(dirname "$0")"

SERVICES="bureau:8081 open-finance:8082 internal-registry:8083 registration-validation:8084"
RUN=/tmp/mcp-servers-run
mkdir -p "$RUN"

case "${1:-start}" in
stop)
  for f in "$RUN"/*.pid; do
    [ -e "$f" ] || continue
    pkill -P "$(cat "$f")" 2>/dev/null || true   # go run deixa o binário como filho
    kill "$(cat "$f")" 2>/dev/null || true
    echo "parado: $(basename "$f" .pid)"
    rm -f "$f"
  done
  ;;
start)
  for pair in $SERVICES; do
    svc=${pair%:*}; port=${pair#*:}
    ( set -a; . "./$svc/.env"; set +a
      # o .env aponta para o host 'postgres' da rede do compose; no host é a porta exposta
      export DATABASE_URL=${DATABASE_URL/@postgres:5432/@127.0.0.1:5433}
      export MCP_PORT=$port
      cd "$svc" && exec go run ./cmd/mcp
    ) >"$RUN/$svc.log" 2>&1 &
    echo $! > "$RUN/$svc.pid"
    echo "iniciando $svc na :$port (log: $RUN/$svc.log)"
  done
  ;;
esac
