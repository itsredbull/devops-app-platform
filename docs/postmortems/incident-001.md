# Incident 001 - Staging Rollout Failure and Recovery

## Summary

- Date: 2026-04-03
- Environment: `uptime-staging`
- Impact: Uptime API rollout failed during a controlled drill. Existing pods stayed serving traffic, but deploy progress was blocked until rollback.

## Trigger

A failure drill intentionally set a non-existent container tag on `uptime-api` to verify detection, rollback speed, and runbook quality.

## Timeline (UTC)

- 09:20 - Drill started with bad image tag update.
- 09:21 - `kubectl rollout status` reported progress deadline issues.
- 09:22 - Rollback initiated with `kubectl rollout undo deployment/uptime-api`.
- 09:24 - Deployment returned to healthy state; alerts returned to normal.

## Evidence

Command log excerpt:

```text
[2026-04-03T09:20:49Z] Injecting bad image to force rollout failure
error: timed out waiting for the condition
[2026-04-03T09:22:11Z] Rollout failed as expected; executing rollback
deployment.apps/uptime-api rolled back
[2026-04-03T09:24:03Z] deployment "uptime-api" successfully rolled out
```

## Root Cause

Image tag referenced an unavailable artifact in registry. Kubernetes could not pull the new image, so new ReplicaSet pods never became ready.

## What Worked

- Rollout status check detected the problem immediately.
- Rollback command restored known-good ReplicaSet quickly.
- Existing service stayed available during rollback due to rolling strategy.

## Corrective Actions

- Add pre-deploy image existence check in CI before manifest update.
- Keep rollback steps as a standard release gate.
- Keep one known-good immutable tag pinned for emergency rollback.

## Preventive Follow-up

- Add release checklist item: verify image digest exists in registry.
- Add dashboard panel for rollout errors/events in next iteration.
