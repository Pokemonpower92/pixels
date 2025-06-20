SERVICE_NAME ?= api
GO=go
GOFLAGS=-mod=vendor
BIN_DIR=bin

.PHONY: build run test clean dev build-service deploy-staging deploy-production teardown-staging teardown-production k8s-status k8s-check-env k8s-create-namespaces k8s-create-secrets

# Build the service
build:
	$(GO) build -o $(BIN_DIR)/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

# Run the service
run:
	$(GO) run ./cmd/$(SERVICE_NAME)

# Run tests
test:
	$(GO) test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)

# Build for specific service
build-service:
	$(GO) build -o $(BIN_DIR)/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

# ==================== Kubernetes Commands ====================

# Check required environment variables for K8s
k8s-check-env:
	@echo "Checking required environment variables..."
	@test -n "$(GITHUB_USERNAME)" || (echo "❌ GITHUB_USERNAME is not set" && exit 1)
	@test -n "$(GITHUB_PAT)" || (echo "❌ GITHUB_PAT is not set" && exit 1)
	@test -n "$(GITHUB_EMAIL)" || (echo "❌ GITHUB_EMAIL is not set" && exit 1)
	@echo "✅ All required environment variables are set"

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
	@echo "   test               - Run tests"
	@echo "   clean              - Clean build artifacts"
	@echo ""
	@echo "☸️  Kubernetes:"
	@echo "   k8s-status         - Show cluster status"
