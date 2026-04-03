#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../app"
go run ./cmd/uptime-api
