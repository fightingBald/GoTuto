# Makefile for common tasks
.PHONY: gen build run fmt tidy migrate-up migrate-down migrate-create db-init test-repo-docker

SERVICE_PKG=./backend/cmd/product-query-svc
BIN_DIR=bin
BIN=$(BIN_DIR)/product-query-svc

gen:
	go generate ./api

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(SERVICE_PKG)

run:
	go run $(SERVICE_PKG)

fmt:
	go fmt ./...

tidy:
	go mod tidy

MIGRATE_PATH=apps/product-query-svc/adapters/outbound/postgres/migrations
# Example; override on CLI: make migrate-up MIGRATE_URL=postgres://...
MIGRATE_URL?=postgres://app:app_password@localhost:5432/productdb?sslmode=disable

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) $(name)

db-init:
	bash scripts/db-init.sh

.PHONY: test-integration-docker
test-integration-docker:
	bash scripts/test-integration-docker.sh ./test -run Postgres

test-repo-docker:
	go test -tags docker ./apps/product-query-svc/adapters/outbound/postgres -run TestCommentRepository_WithDocker -count=1


# Notes:
# - Requires golang-migrate installed to use migrate-* targets
# - Use gen only if codegen toolchain is available
