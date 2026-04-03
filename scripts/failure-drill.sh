#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-local}"
API_URL="${API_URL:-http://localhost:8080}"
NAMESPACE="${NAMESPACE:-uptime-staging}"
DEPLOYMENT="${DEPLOYMENT:-uptime-api}"
GOOD_IMAGE="${GOOD_IMAGE:-ghcr.io/itsredbull/devops-app-platform:latest}"
BAD_IMAGE="${BAD_IMAGE:-ghcr.io/itsredbull/devops-app-platform:does-not-exist}"

log() {
  printf "[%s] %s\n" "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" "$*"
}

run_local_drill() {
  log "Starting local failure drill against ${API_URL}"
  log "Injecting one failing target"

  curl -sS -X POST "${API_URL}/api/v1/targets" \
    -H 'Content-Type: application/json' \
    -d '{"url":"http://127.0.0.1:65534","check_interval_seconds":10,"timeout_seconds":2,"enabled":true}' \
    >/tmp/failure-drill-target.json

  log "Waiting 20s for scheduler/checker to process failures"
  sleep 20

  log "Current status snapshot (filtered)"
  curl -sS "${API_URL}/api/v1/status" | tr ',' '\n' | rg -E '127.0.0.1:65534|"success":false|"error"|"status_code"' || true

  log "Local drill completed"
  log "Expected result: failure ratio and latency panels should move in Grafana"
}

run_k8s_rollback_drill() {
  command -v kubectl >/dev/null 2>&1 || {
    echo "kubectl is required for k8s mode" >&2
    exit 1
  }

  log "Starting Kubernetes rollback drill"
  log "Namespace=${NAMESPACE} Deployment=${DEPLOYMENT}"

  log "Current image before drill"
  kubectl -n "${NAMESPACE}" get deploy "${DEPLOYMENT}" -o jsonpath='{.spec.template.spec.containers[0].image}' && echo

  log "Injecting bad image to force rollout failure"
  kubectl -n "${NAMESPACE}" set image "deployment/${DEPLOYMENT}" "${DEPLOYMENT}=${BAD_IMAGE}"

  set +e
  kubectl -n "${NAMESPACE}" rollout status "deployment/${DEPLOYMENT}" --timeout=90s
  rollout_rc=$?
  set -e

  if [[ ${rollout_rc} -eq 0 ]]; then
    log "WARNING: rollout unexpectedly succeeded with bad image"
  else
    log "Rollout failed as expected; executing rollback"
  fi

  kubectl -n "${NAMESPACE}" rollout undo "deployment/${DEPLOYMENT}"
  kubectl -n "${NAMESPACE}" rollout status "deployment/${DEPLOYMENT}" --timeout=120s

  log "Image after rollback"
  kubectl -n "${NAMESPACE}" get deploy "${DEPLOYMENT}" -o jsonpath='{.spec.template.spec.containers[0].image}' && echo

  log "Rollback drill completed"
  log "If image after rollback is not expected, restore manually to ${GOOD_IMAGE}"
}

case "${MODE}" in
  local)
    run_local_drill
    ;;
  k8s)
    run_k8s_rollback_drill
    ;;
  *)
    echo "Usage: $0 [local|k8s]" >&2
    exit 1
    ;;
esac
