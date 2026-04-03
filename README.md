# DevOps App Platform

Production-style DevOps portfolio project built around a URL Uptime Monitor service.

## Purpose

This project demonstrates end-to-end DevOps work:

- infrastructure as code
- CI/CD with quality and security gates
- Kubernetes deployment with GitOps
- observability and alerting
- incident handling and rollback

## Current Phase

`Phase 10: Portfolio Polish (completed)`

- PostgreSQL-backed target/check storage
- target CRUD APIs implemented
- status and checks APIs implemented
- scheduler + checker worker implemented
- Prometheus metrics endpoint implemented
- integration API flow test added
- checker retry/backoff and error classification added
- containerized local stack for app + db added
- GitHub Actions quality gates implemented
- Trivy image scan added with critical-fail policy
- Kubernetes base manifests for app + postgres
- dev and staging overlays via Kustomize
- config/secret management through generators
- Argo CD app-of-apps bootstrap added
- auto-sync for dev/staging from Git
- Terraform modules for network, Kubernetes, database, monitoring
- Remote state bootstrap (S3 + DynamoDB)
- Dedicated Terraform env stacks for dev and staging
- Prometheus scrape configuration added
- Grafana dashboard for uptime checks and latency
- Alert rules for failure ratio, high latency, and no-checks condition
- Failure drill script added for local and Kubernetes scenarios
- Rollback test procedure documented and drill-backed
- Incident postmortem and demo screenshot checklist added
- One-command quickstart added (`make quickstart`)
- Architecture diagram added in docs
- README now includes build/learning/trade-off summary

## CI Quality Gates

Workflow: `.github/workflows/ci.yaml`

On every push/PR to `main`, CI runs:

- format check (`gofmt -l`)
- lint (`go vet ./...`)
- tests (`go test ./...`)
- build (`go build ./cmd/uptime-api`)
- Docker image build (`docker build`)
- Trivy image scan

Security gate policy:

- pipeline fails if Trivy finds `CRITICAL` vulnerabilities
- SARIF report is uploaded to GitHub Security tab

## Kubernetes Deploy

This phase uses **Kustomize** overlays (not Helm):

- base: `deploy/k8s/base`
- dev overlay: `deploy/k8s/overlays/dev`
- staging overlay: `deploy/k8s/overlays/staging`

Deployment guide:

- `deploy/k8s/README.md`

## GitOps (Argo CD)

Argo CD manifests are in:

- `deploy/argocd/bootstrap/root-app.yaml`
- `deploy/argocd/apps/project.yaml`
- `deploy/argocd/apps/app-dev.yaml`
- `deploy/argocd/apps/app-staging.yaml`

Install and bootstrap:

```bash
make argocd-install
make argocd-bootstrap
```

Detailed guide:

- `deploy/argocd/README.md`

## Terraform Infra

Terraform layout:

- `infra/terraform/modules/network`
- `infra/terraform/modules/kubernetes`
- `infra/terraform/modules/database`
- `infra/terraform/modules/monitoring`
- `infra/terraform/envs/dev`
- `infra/terraform/envs/staging`
- `infra/terraform/global`

Quick start:

```bash
cd infra/terraform/global
cp terraform.tfvars.example terraform.tfvars
terraform init && terraform apply

cd ../envs/dev
cp backend.hcl.example backend.hcl
cp terraform.tfvars.example terraform.tfvars
terraform init -backend-config=backend.hcl
terraform plan
```

## Observability

Monitoring assets:

- `monitoring/prometheus/prometheus.yml`
- `monitoring/prometheus/alerts/uptime-rules.yaml`
- `monitoring/grafana/dashboards/uptime-overview.json`

Start local monitoring stack:

```bash
make stack-up
make monitoring-up
```

Access:

- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (admin/admin)

## Project Structure

- `app/`: Go service source, migrations, Dockerfile
- `infra/`: Terraform and Ansible layout
- `deploy/`: Kubernetes, Helm, Argo CD manifests
- `monitoring/`: Prometheus rules and Grafana dashboards
- `docs/`: architecture, API, runbook, postmortems
- `.github/workflows/`: CI/CD workflows

## One-Command Quickstart

```bash
make quickstart
```

This starts app, database, Prometheus, and Grafana together.

Stop everything:

```bash
make quickstart-down
```

## Architecture Diagram

- See: `docs/architecture.md`

## Local Run (host app + docker db)

1. Start PostgreSQL:

```bash
make db-up
```

2. Run API service:

```bash
make run PORT=8081
```

3. Verify endpoints:

```bash
curl -sS http://localhost:8081/healthz
curl -sS http://localhost:8081/readyz
curl -sS http://localhost:8081/metrics
```

4. Create target and view status:

```bash
curl -sS -X POST http://localhost:8081/api/v1/targets \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","check_interval_seconds":30,"timeout_seconds":10,"enabled":true}'

curl -sS http://localhost:8081/api/v1/status
```

## Local Run (fully containerized)

```bash
make stack-up
curl -sS http://localhost:8080/healthz
make stack-down
```

## Tests

```bash
make test
make test-integration
```

## Checker Runtime Settings

- `CHECK_MAX_ATTEMPTS` (default `3`)
- `CHECK_RETRY_BACKOFF_MS` (default `200`)

Retries are used for transient failures (network errors, timeouts, and HTTP 5xx).

## Troubleshooting

### `listen tcp :8080: bind: address already in use`

Run on a different port:

```bash
make run PORT=8081
```

Or find and stop process on `8080`:

```bash
lsof -nP -iTCP:8080 -sTCP:LISTEN
kill -9 <PID>
```

## Ops Credibility Artifacts

- Failure drill script: `scripts/failure-drill.sh`
- Runbook: `docs/runbook.md`
- Postmortem: `docs/postmortems/incident-001.md`
- Demo screenshots checklist: `docs/screenshots/README.md`

## What I Built

- A Go-based uptime monitoring API with scheduler and checker workers
- PostgreSQL-backed target and check history persistence
- CI quality gates with lint/test/build and container vulnerability scanning
- Kubernetes deployment with dev/staging overlays and GitOps sync via Argo CD
- Terraform stacks for infra modules and remote state
- Observability with Prometheus metrics, alert rules, and Grafana dashboard
- Ops workflows with failure drills, rollback tests, runbook, and postmortem

## What I Learned

- How to connect app reliability signals to actionable alerting
- How to structure GitOps so rollback is fast and low-risk
- How to design Terraform modules usable across environments
- How to write operational docs that support incident response

## Trade-Offs

- Postgres is deployed in-cluster for simplicity; managed DB would be better for production
- Single app service keeps architecture clear; fewer components but less specialization
- Local quickstart favors speed over strict production parity
- Alert thresholds are practical defaults and should be tuned with real traffic
