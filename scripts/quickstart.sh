#!/usr/bin/env bash
set -euo pipefail

ACTION="${1:-up}"

up() {
  echo "[quickstart] starting app, db, prometheus, grafana"
  docker compose up -d --build postgres app prometheus grafana

  echo "[quickstart] waiting for API health"
  for _ in {1..30}; do
    if curl -fsS http://localhost:8080/healthz >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done

  echo "[quickstart] ready"
  echo "- API:        http://localhost:8080/healthz"
  echo "- Prometheus: http://localhost:9090"
  echo "- Grafana:    http://localhost:3000 (admin/admin)"
}

down() {
  echo "[quickstart] stopping local stack"
  docker compose down
}

case "${ACTION}" in
  up)
    up
    ;;
  down)
    down
    ;;
  *)
    echo "Usage: $0 [up|down]" >&2
    exit 1
    ;;
esac
