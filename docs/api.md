# API Contract (v1)

## Health

- `GET /healthz`
- `GET /readyz`
- `GET /metrics` (Prometheus)

## Targets

- `POST /api/v1/targets`
- `GET /api/v1/targets`
- `GET /api/v1/targets/{id}`
- `DELETE /api/v1/targets/{id}`

### POST /api/v1/targets request

```json
{
  "url": "https://example.com",
  "check_interval_seconds": 30,
  "timeout_seconds": 10,
  "enabled": true
}
```

### Target object

```json
{
  "id": "uuid",
  "url": "https://example.com",
  "check_interval_seconds": 30,
  "timeout_seconds": 10,
  "enabled": true,
  "created_at": "2026-04-03T10:00:00Z",
  "updated_at": "2026-04-03T10:00:00Z"
}
```

## Checks

- `GET /api/v1/checks?target_id=<uuid>&limit=100`
- `GET /api/v1/status` (latest result for each target)

### Check object

```json
{
  "id": "uuid",
  "target_id": "uuid",
  "checked_at": "2026-04-03T10:00:30Z",
  "status_code": 200,
  "latency_ms": 123,
  "success": true,
  "error": ""
}
```

## Expected responses

- `201` on target creation
- `204` on successful delete
- `404` for missing target id
- `409` when URL already exists
