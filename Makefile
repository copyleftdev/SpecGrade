# SpecGrade Makefile
# Development and build automation for SpecGrade OpenAPI validator

# Variables
BINARY_NAME=specgrade
DOCKER_IMAGE=specgrade
DOCKER_TAG=latest
GO_VERSION=1.21
MAIN_PACKAGE=.
BUILD_DIR=build
COVERAGE_FILE=coverage.out
LDFLAGS=-s -w -X main.version=$(shell git describe --tags --always --dirty)
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

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

.PHONY: build-linux-amd64
build-linux-amd64: ## Build for Linux AMD64
	@printf "$(BLUE)Building $(BINARY_NAME) for Linux AMD64...$(NC)\n"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

.PHONY: build-linux-arm64
build-linux-arm64: ## Build for Linux ARM64
	@printf "$(BLUE)Building $(BINARY_NAME) for Linux ARM64...$(NC)\n"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)

.PHONY: build-windows-amd64
build-windows-amd64: ## Build for Windows AMD64
	@printf "$(BLUE)Building $(BINARY_NAME) for Windows AMD64...$(NC)\n"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

.PHONY: build-darwin-amd64
build-darwin-amd64: ## Build for macOS AMD64
	@printf "$(BLUE)Building $(BINARY_NAME) for macOS AMD64...$(NC)\n"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Build for macOS ARM64 (Apple Silicon)
	@printf "$(BLUE)Building $(BINARY_NAME) for macOS ARM64...$(NC)\n"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)

.PHONY: build-all
build-all: build-linux-amd64 build-linux-arm64 build-windows-amd64 build-darwin-amd64 build-darwin-arm64 ## Build for all platforms

# Enhanced Developer Reporting Demo
.PHONY: demo-developer
demo-developer: ## Demo enhanced developer reporting format
	@echo "Running SpecGrade with enhanced developer reporting..."
	@echo "=========================================="
	go run main.go --target-dir test/sample-spec --output-format developer

.PHONY: demo-all-formats
demo-all-formats: ## Demo all output formats
	@echo "Testing all output formats..."
	@echo "\n📋 CLI Format:"
	go run main.go --target-dir test/sample-spec --output-format cli
	@echo "\n🎯 Developer Format:"
	go run main.go --target-dir test/sample-spec --output-format developer
	@echo "\n📊 JSON Format:"
	go run main.go --target-dir test/sample-spec --output-format json

# Testing
.PHONY: test
test: ## Run all tests
	@printf "$(BLUE)Running tests...$(NC)\n"
	go test -v ./...

.PHONY: test-short
test-short: ## Run tests with short flag
	@printf "$(BLUE)Running short tests...$(NC)\n"
	go test -short ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@printf "$(BLUE)Running tests with coverage...$(NC)\n"
	go test -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@printf "$(GREEN)Coverage report generated: coverage.html$(NC)\n"

.PHONY: test-coverage-func
test-coverage-func: ## Show test coverage by function
	@printf "$(BLUE)Running coverage analysis...$(NC)\n"
	go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	go tool cover -func=$(COVERAGE_FILE)

.PHONY: validate-enhanced-reporting
validate-enhanced-reporting: ## Validate enhanced developer reporting works correctly
	@echo "Validating enhanced developer reporting..."
	@echo "Testing compilation..."
	go build -o /tmp/specgrade-test .
	@echo "✅ Compilation successful"
	@echo "Testing enhanced developer format..."
	go run main.go --target-dir test/sample-spec --output-format developer > /tmp/test-developer.out
	@echo "✅ Developer format works"
	@echo "Testing all formats..."
	go run main.go --target-dir test/sample-spec --output-format json > /tmp/test-json.out
	go run main.go --target-dir test/sample-spec --output-format cli > /tmp/test-cli.out
	@echo "✅ All output formats validated successfully!"
	@rm -f /tmp/specgrade-test /tmp/test-*.out

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
	@printf "$(BLUE)Running linter...$(NC)\n"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		printf "$(RED)golangci-lint not installed.$(NC)\n"; \
		printf "$(YELLOW)Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)\n"; \
		exit 1; \
	fi

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@printf "$(BLUE)Running linter with auto-fix...$(NC)\n"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix --timeout=5m; \
	else \
		printf "$(RED)golangci-lint not installed.$(NC)\n"; \
		exit 1; \
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
	@printf "$(BLUE)Verifying go modules...$(NC)\n"
	go mod verify

.PHONY: mod-download
mod-download: ## Download go modules
	@printf "$(BLUE)Downloading go modules...$(NC)\n"
	go mod download

.PHONY: deps-check
deps-check: ## Check for dependency updates
	@printf "$(BLUE)Checking for dependency updates...$(NC)\n"
	@if command -v go-mod-outdated >/dev/null 2>&1; then \
		go list -u -m -json all | go-mod-outdated -update -direct; \
	else \
		printf "$(YELLOW)go-mod-outdated not installed. Install with: go install github.com/psampaz/go-mod-outdated@latest$(NC)\n"; \
	fi

.PHONY: deps-update
deps-update: ## Update dependencies
	@printf "$(BLUE)Updating dependencies...$(NC)\n"
	go get -u ./...
	go mod tidy

.PHONY: security-scan
security-scan: ## Run security scan
	@printf "$(BLUE)Running security scan...$(NC)\n"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		printf "$(YELLOW)gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(NC)\n"; \
	fi

.PHONY: vuln-check
vuln-check: ## Check for known vulnerabilities
	@printf "$(BLUE)Checking for vulnerabilities...$(NC)\n"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		printf "$(YELLOW)govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest$(NC)\n"; \
	fi

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
ci-test: mod-verify vet lint test-race test-coverage ## Run all CI tests

.PHONY: ci-build
ci-build: build-all ## Build all platforms for CI

.PHONY: ci-security
ci-security: security-scan vuln-check ## Run all security checks

.PHONY: ci-full
ci-full: clean mod-tidy ci-test ci-security ci-build ## Full CI pipeline

# Release
.PHONY: release-check
release-check: ci-full ## Check if ready for release
	@printf "$(GREEN)Release check completed successfully!$(NC)\n"

.PHONY: release-notes
release-notes: ## Generate release notes
	@printf "$(BLUE)Generating release notes...$(NC)\n"
	@if command -v git-chglog >/dev/null 2>&1; then \
		git-chglog --output CHANGELOG.md; \
	else \
		printf "$(YELLOW)git-chglog not installed. Using git log instead...$(NC)\n"; \
		git log --oneline --decorate --graph --since="$(shell git describe --tags --abbrev=0)" > RELEASE_NOTES.md; \
	fi

.PHONY: tag
tag: ## Create and push a new tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		printf "$(RED)VERSION is required. Usage: make tag VERSION=v1.0.0$(NC)\n"; \
		exit 1; \
	fi
	@printf "$(BLUE)Creating tag $(VERSION)...$(NC)\n"
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# Quick development workflow
.PHONY: quick
quick: fmt vet test build ## Quick development cycle: format, vet, test, build

.PHONY: dev-setup
dev-setup: ## Setup development environment
	@printf "$(BLUE)Setting up development environment...$(NC)\n"
	@printf "$(YELLOW)Installing development tools...$(NC)\n"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/psampaz/go-mod-outdated@latest
	@printf "$(GREEN)Development environment setup complete!$(NC)\n"

.PHONY: pre-commit
pre-commit: fmt vet lint test-short ## Pre-commit checks
	@printf "$(GREEN)Pre-commit checks passed!$(NC)\n"

.PHONY: all
all: clean mod-tidy fmt vet test-coverage build-all ## Full build pipeline

# Performance
.PHONY: profile-cpu
profile-cpu: build ## Run CPU profiling
	@printf "$(BLUE)Running CPU profiling...$(NC)\n"
	./$(BUILD_DIR)/$(BINARY_NAME) --target-dir=./test/sample-spec --cpuprofile=cpu.prof
	go tool pprof cpu.prof

.PHONY: profile-mem
profile-mem: build ## Run memory profiling
	@printf "$(BLUE)Running memory profiling...$(NC)\n"
	./$(BUILD_DIR)/$(BINARY_NAME) --target-dir=./test/sample-spec --memprofile=mem.prof
	go tool pprof mem.prof
