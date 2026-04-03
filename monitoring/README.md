# Observability Stack

## Prometheus

Main config:

- `monitoring/prometheus/prometheus.yml`
- alert rules: `monitoring/prometheus/alerts/uptime-rules.yaml`

Scrapes:

- Prometheus (`localhost:9090`)
- Uptime API metrics (`uptime-api:8080/metrics`)

## Grafana

Provisioning:

- datasource: `monitoring/grafana/provisioning/datasources/prometheus.yaml`
- dashboards provider: `monitoring/grafana/provisioning/dashboards/dashboards.yaml`
- dashboard: `monitoring/grafana/dashboards/uptime-overview.json`

## Included Alert Rules

- `UptimeHighFailureRatio`: failure ratio > 20% for 10m
- `UptimeHighLatencyP95`: p95 latency > 1000ms for 10m
- `UptimeNoChecks`: no checks for 10m while enabled targets > 0
