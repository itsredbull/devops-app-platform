# Runbook

## Startup Checks

- API health endpoint returns `200`.
- DB connectivity is healthy.
- Metrics endpoint is exposed.

## Alert: High Failure Ratio

1. Check recent target status via `/api/v1/status`.
2. Verify DNS and outbound network from app pod.
3. Check recent deploys and image tags.
4. Roll back to previous known good release if failures started after deploy.

## Alert: High Latency

1. Validate target endpoint performance externally.
2. Compare app resource usage (CPU/memory) and DB load.
3. Increase timeout only if target behavior justifies it.

## Rollback

- GitOps: revert image tag commit in deployment manifests.
- Confirm Argo CD sync and watch alert clear trend.
