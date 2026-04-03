# Runbook

## Startup Checks

1. Confirm API health:

```bash
curl -sS http://localhost:8080/healthz
curl -sS http://localhost:8080/readyz
```

2. Confirm metrics are exposed:

```bash
curl -sS http://localhost:8080/metrics | head
```

3. Confirm DB is reachable (local stack):

```bash
docker compose ps
```

## Failure Drill

Local drill:

```bash
./scripts/failure-drill.sh local
```

Kubernetes drill (staging):

```bash
./scripts/failure-drill.sh k8s
```

## Alert: High Failure Ratio

1. Review status API:

```bash
curl -sS http://localhost:8080/api/v1/status
```

2. Check whether failures are concentrated on one target or global.
3. Verify outbound network and DNS from app environment.
4. If failure started after deploy, run rollback.

## Alert: High Latency

1. Check p95 latency panel in Grafana.
2. Compare latency increase with CPU/memory pressure.
3. Validate target endpoints from outside cluster.
4. Revert last deployment if latency regression is release-linked.

## Alert: No Checks

1. Verify scheduler process is running (app logs).
2. Verify enabled target count is greater than zero.
3. Confirm DB writes are succeeding.
4. Restart app deployment only after logs and DB checks.

## Rollback Test

Quick rollback procedure:

```bash
kubectl -n uptime-staging rollout undo deployment/uptime-api
kubectl -n uptime-staging rollout status deployment/uptime-api --timeout=120s
```

GitOps rollback (preferred):

1. Revert the commit that changed image tag or overlay values.
2. Push to `main`.
3. Confirm Argo CD syncs and deployment becomes healthy.

## Post-Incident Checklist

- Record timeline with UTC timestamps.
- Capture root cause and exact command output.
- Add at least one preventive action item.
- Link evidence screenshots in `docs/screenshots/`.
