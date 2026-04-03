SHELL := /bin/bash

PORT ?= 8080
TEST_DATABASE_DSN ?= postgres://uptime:uptime@localhost:5432/uptime?sslmode=disable

.PHONY: run test test-integration lint fmt build db-up db-down stack-up stack-down monitoring-up monitoring-down quickstart quickstart-down argocd-install argocd-bootstrap

run:
	cd app && APP_PORT=$(PORT) go run ./cmd/uptime-api

test:
	cd app && go test ./...

test-integration:
	cd app && TEST_DATABASE_DSN='$(TEST_DATABASE_DSN)' go test ./internal/api -run TestIntegrationAPIFlow -v

lint:
	cd app && go vet ./...

fmt:
	cd app && go fmt ./...

build:
	cd app && go build ./cmd/uptime-api

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

stack-up:
	docker compose up -d --build postgres app

stack-down:
	docker compose down

monitoring-up:
	docker compose up -d prometheus grafana

monitoring-down:
	docker compose stop grafana prometheus

quickstart:
	./scripts/quickstart.sh up

quickstart-down:
	./scripts/quickstart.sh down

argocd-install:
	kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

argocd-bootstrap:
	kubectl apply -f deploy/argocd/bootstrap/root-app.yaml
