SHELL := /bin/bash

PORT ?= 8080
TEST_DATABASE_DSN ?= postgres://uptime:uptime@localhost:5432/uptime?sslmode=disable

.PHONY: run test test-integration lint fmt build db-up db-down stack-up stack-down

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
