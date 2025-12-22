# MangaHub Makefile
# Development and build automation

.PHONY: all build run test clean generate-types generate-proto help

# Default target
all: build

# ==========================================
# Build Commands
# ==========================================

build: ## Build all server binaries
	@echo "Building all servers..."
	go build -o bin/api-server ./cmd/api-server
	go build -o bin/tcp-server ./cmd/tcp-server
	go build -o bin/udp-server ./cmd/udp-server
	go build -o bin/grpc-server ./cmd/grpc-server
	go build -o bin/cli ./cmd/cli
	@echo "Build complete! Binaries in ./bin/"

build-api: ## Build HTTP API server only
	go build -o bin/api-server ./cmd/api-server

build-tcp: ## Build TCP server only
	go build -o bin/tcp-server ./cmd/tcp-server

build-grpc: ## Build gRPC server only
	go build -o bin/grpc-server ./cmd/grpc-server

# ==========================================
# Run Commands
# ==========================================

run-api: ## Run HTTP API server
	go run ./cmd/api-server

run-tcp: ## Run TCP server
	go run ./cmd/tcp-server

run-udp: ## Run UDP server
	go run ./cmd/udp-server

run-grpc: ## Run gRPC server
	go run ./cmd/grpc-server

run-all: ## Run all servers (requires tmux or multiple terminals)
	@echo "Starting all servers..."
	@echo "Use 'make run-api', 'make run-tcp', etc. in separate terminals"

# ==========================================
# Code Generation
# ==========================================

generate-types: ## Generate TypeScript types from OpenAPI spec
	@echo "Generating TypeScript types..."
	yarn workspace @mangahub/types generate
	@echo "Types generated at packages/types/src/generated.ts"

generate-proto: ## Generate Go code from Protocol Buffers
	@echo "Generating gRPC code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/manga.proto
	@echo "Proto files generated"

generate: generate-proto generate-types ## Generate all code

# ==========================================
# Database Commands
# ==========================================

migrate-up: ## Run database migrations
	go run ./scripts/migrate.go up

migrate-down: ## Rollback database migrations
	go run ./scripts/migrate.go down

seed: ## Seed database with sample data
	go run ./scripts/seed.go

db-reset: migrate-down migrate-up seed ## Reset database and reseed

# ==========================================
# Testing
# ==========================================

test: ## Run all tests
	go test -v ./internal/... ./pkg/...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-integration: ## Run integration tests
	go test -v -tags=integration ./test/integration/...

# ==========================================
# Development
# ==========================================

dev: ## Run API server with hot reload (requires air)
	air -c .air.toml

lint: ## Run linter
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

tidy: ## Tidy go modules
	go mod tidy

# ==========================================
# Documentation
# ==========================================

docs-preview: ## Preview OpenAPI documentation
	yarn workspace @mangahub/spec preview

docs-validate: ## Validate OpenAPI specification
	yarn workspace @mangahub/spec validate

# ==========================================
# JavaScript/TypeScript Commands (Monorepo)
# ==========================================

js-install: ## Install JavaScript dependencies
	yarn install

js-build: ## Build all JavaScript packages
	yarn build

js-dev: ## Run Next.js web app in dev mode
	yarn workspace @mangahub/web dev

js-test: ## Run JavaScript tests
	yarn test

js-lint: ## Lint JavaScript code
	yarn lint

js-typecheck: ## Type-check JavaScript code
	yarn typecheck

js-clean: ## Clean JavaScript build artifacts
	yarn clean

# ==========================================
# Cleanup
# ==========================================

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Cleaned Go build artifacts"

clean-all: clean js-clean ## Clean all build artifacts (Go + JS)
	@echo "Cleaned all build artifacts"

# ==========================================
# Help
# ==========================================

help: ## Show this help message
	@echo "MangaHub - Manga Tracking System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
