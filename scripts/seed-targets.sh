#!/usr/bin/env bash
set -euo pipefail

API_URL="${API_URL:-http://localhost:8080}"

curl -sS -X POST "$API_URL/api/v1/targets" \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","check_interval_seconds":30,"timeout_seconds":10,"enabled":true}' || true
