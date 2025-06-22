SERVICE_NAME ?= api
GO=go
GOFLAGS=-mod=vendor
BIN_DIR=bin
COVERAGE_DIR=coverage

.PHONY: build run test test-e2e test-verbose test-short test-coverage test-coverage-html clean build-service k8s-status help

# Build the service
build:
	$(GO) build -o $(BIN_DIR)/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

# Run the service
run:
	$(GO) run ./cmd/$(SERVICE_NAME)

# Run tests with clean output (only shows failures and summary)
test:
	@echo "🧪 Running tests..."
	@$(GO) test ./... -count=1

test-e2e:
	@$(GO) test ./internal/tests -tags=e2e -count=1

test-integration:
	@$(GO) test ./internal/tests -tags=integration -count=1

# Run tests with verbose output (original behavior)
test-verbose:
	$(GO) test -v ./...

# Run tests with short output (skip long-running tests)
test-short:
	@echo "🧪 Running short tests..."
	@$(GO) test -short ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@$(GO) test -cover ./...

# Run tests with coverage report
test-coverage-html:
	@echo "🧪 Generating coverage report..."
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "✅ Coverage report generated: $(COVERAGE_DIR)/coverage.html"

# Clean build artifacts
clean:
	rm -f $(BIN_DIR)/*
	rm -f $(COVERAGE_DIR)/*

# Build for specific service
build-service:
	$(GO) build -o $(BIN_DIR)/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

# ==================== Kubernetes Commands ====================

# Show K8s cluster status
k8s-status:
	@echo "==================== Kubernetes Status ===================="
	@echo ""
	@echo "📦 Namespaces:"
	@kubectl get namespaces | grep pixels || echo "   No pixels namespaces found"
	@echo ""
	@echo "🚀 Staging Pods:"
	@kubectl get pods -n pixels-staging 2>/dev/null || echo "   No staging namespace"
	@echo ""
	@echo "🌐 Services:"
	@kubectl get services -n pixels-staging 2>/dev/null || echo "   No staging services"
	@echo ""

# Help command
help:
	@echo "Available commands:"
	@echo ""
	@echo "🔨 Development:"
	@echo "   build              - Build the service"
	@echo "   test               - Run tests (clean output)"
	@echo "   test-verbose       - Run tests with verbose output"
	@echo "   test-short         - Run quick tests only"
	@echo "   test-coverage      - Run tests with coverage"
	@echo "   test-coverage-html - Generate HTML coverage report"
	@echo "   clean              - Clean build artifacts"
	@echo ""
	@echo "☸️  Kubernetes:"
	@echo "   k8s-status         - Show cluster status"