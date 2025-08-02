# SpecGrade Makefile
# Development and build automation for SpecGrade OpenAPI validator

# Variables
BINARY_NAME=specgrade
DOCKER_IMAGE=specgrade
DOCKER_TAG=latest
GO_VERSION=1.24.1
MAIN_PACKAGE=.
BUILD_DIR=build
COVERAGE_FILE=coverage.out

# Default target
.PHONY: help
help: ## Show this help message
	@echo "SpecGrade Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
.PHONY: dev
dev: ## Start development environment with hot reload
	@echo "Starting development environment..."
	docker-compose up --build

.PHONY: dev-down
dev-down: ## Stop development environment
	docker-compose down

.PHONY: dev-logs
dev-logs: ## Show development environment logs
	docker-compose logs -f

# Building
.PHONY: build
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building $(BINARY_NAME) for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

.PHONY: build-darwin
build-darwin: ## Build for macOS
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## Build for all platforms

# Testing
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -race ./...

# Code Quality
.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

.PHONY: mod-verify
mod-verify: ## Verify go modules
	@echo "Verifying go modules..."
	go mod verify

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -v $(PWD)/test/sample-spec:/specs $(DOCKER_IMAGE):$(DOCKER_TAG) --target-dir=/specs

.PHONY: docker-shell
docker-shell: ## Get shell in Docker container
	docker run --rm -it -v $(PWD):/app -w /app $(DOCKER_IMAGE):$(DOCKER_TAG) sh

.PHONY: docker-clean
docker-clean: ## Clean Docker images
	@echo "Cleaning Docker images..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	docker system prune -f

# Installation
.PHONY: install
install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(MAIN_PACKAGE)

.PHONY: uninstall
uninstall: ## Uninstall binary from GOPATH/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Demo and Examples
.PHONY: demo
demo: build ## Run demo with sample spec
	@echo "Running SpecGrade demo..."
	./$(BUILD_DIR)/$(BINARY_NAME) --target-dir=./test/sample-spec --spec-version=3.1.0

.PHONY: demo-json
demo-json: build ## Run demo with JSON output
	@echo "Running SpecGrade demo (JSON output)..."
	./$(BUILD_DIR)/$(BINARY_NAME) --target-dir=./test/sample-spec --spec-version=3.1.0 --output-format=json

.PHONY: demo-config
demo-config: build ## Run demo with config file
	@echo "Running SpecGrade demo (config file)..."
	./$(BUILD_DIR)/$(BINARY_NAME) --config=specgrade.yaml

.PHONY: rules-list
rules-list: build ## List all available rules
	@echo "Listing available rules..."
	./$(BUILD_DIR)/$(BINARY_NAME) rules ls

.PHONY: docs-generate
docs-generate: build ## Generate rule documentation
	@echo "Generating rule documentation..."
	./$(BUILD_DIR)/$(BINARY_NAME) --docs > RULES.md

# Cleanup
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html

.PHONY: clean-all
clean-all: clean docker-clean ## Clean everything including Docker

# CI/CD helpers
.PHONY: ci-test
ci-test: mod-verify vet test-race test-coverage ## Run all CI tests

.PHONY: ci-build
ci-build: build-all ## Build all platforms for CI

# Release
.PHONY: release-check
release-check: ci-test ci-build ## Check if ready for release
	@echo "Release check completed successfully!"

# Quick development workflow
.PHONY: quick
quick: fmt vet test build ## Quick development cycle: format, vet, test, build

.PHONY: all
all: clean mod-tidy fmt vet test-coverage build-all ## Full build pipeline
