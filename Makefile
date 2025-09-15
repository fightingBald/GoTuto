# Makefile 用于简化常用命令
.PHONY: gen build run fmt

gen:
	./scripts/generate-api.sh

build: gen
	go build -o bin/server ./cmd/server

run: build
	./bin/server

fmt:
	gofmt -w .

