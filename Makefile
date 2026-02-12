.PHONY: help test build lint clean docker-build run

# Variables
APP_NAME := bookcommunity
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)

help: ## Show this help message
	@echo "BookCommunity - Makefile commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy

test: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

test-short: ## Run short tests
	go test -v -short ./...

build: ## Build binary
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/$(APP_NAME) main.go

build-linux: ## Build for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(APP_NAME)-linux main.go

lint: ## Run linter
	golangci-lint run --timeout=5m

lint-fix: ## Run linter with auto-fix
	golangci-lint run --fix --timeout=5m

fmt: ## Format code
	gofmt -s -w .
	goimports -w -local github.com/sylvia-ymlin/Coconut-book-community .

vet: ## Run go vet
	go vet ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f $(APP_NAME)

docker-build: ## Build Docker image
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

docker-run: ## Run Docker container
	docker-compose up -d

docker-stop: ## Stop Docker container
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

run: ## Run application locally
	go run main.go

dev: ## Run with hot reload (requires air)
	air

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	# TODO: Add migration command

migrate-down: ## Run database migrations down
	@echo "Rolling back migrations..."
	# TODO: Add migration rollback command

ci: lint test build ## Run CI pipeline locally

release: ## Create a new release (requires goreleaser)
	goreleaser release --clean

release-snapshot: ## Create a snapshot release
	goreleaser release --snapshot --clean

security: ## Run security scan
	gosec -no-fail ./...

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

.DEFAULT_GOAL := help
