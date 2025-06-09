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

# Run air for development
dev:
	air -c .air.$(SERVICE_NAME).toml

# ==================== Kubernetes Commands ====================

# Check required environment variables for K8s
k8s-check-env:
	@echo "Checking required environment variables..."
	@test -n "$(GITHUB_USERNAME)" || (echo "❌ GITHUB_USERNAME is not set" && exit 1)
	@test -n "$(GITHUB_PAT)" || (echo "❌ GITHUB_PAT is not set" && exit 1)
	@test -n "$(GITHUB_EMAIL)" || (echo "❌ GITHUB_EMAIL is not set" && exit 1)
	@echo "✅ All required environment variables are set"

# Create K8s namespaces
k8s-create-namespaces:
	@echo "Creating namespaces..."
	kubectl create namespace pixels-staging --dry-run=client -o yaml | kubectl apply -f -
	kubectl create namespace pixels-production --dry-run=client -o yaml | kubectl apply -f -
	kubectl label namespace pixels-staging name=pixels-staging --overwrite
	kubectl label namespace pixels-production name=pixels-production --overwrite
	@echo "✅ Namespaces created and labeled"

# Create image pull secrets
k8s-create-secrets: k8s-check-env
	@echo "Creating image pull secrets..."
	kubectl create secret docker-registry ghcr-secret \
		--docker-server=ghcr.io \
		--docker-username=$(GITHUB_USERNAME) \
		--docker-password=$(GITHUB_PAT) \
		--docker-email=$(GITHUB_EMAIL) \
		--namespace=pixels-staging \
		--dry-run=client -o yaml | kubectl apply -f -
	kubectl create secret docker-registry ghcr-secret \
		--docker-server=ghcr.io \
		--docker-username=$(GITHUB_USERNAME) \
		--docker-password=$(GITHUB_PAT) \
		--docker-email=$(GITHUB_EMAIL) \
		--namespace=pixels-production \
		--dry-run=client -o yaml | kubectl apply -f -
	@echo "✅ Image pull secrets created"

# Deploy to staging environment
deploy-staging: k8s-create-namespaces k8s-create-secrets
	@echo "Deploying to staging..."
	kubectl apply -k k8s/overlays/staging/
	@echo "✅ Staging deployment complete"
	@echo ""
	@echo "🔗 Access staging at:"
	@echo "   http://10.0.0.130:30080"
	@echo "   http://10.0.0.99:30080"
	@echo ""
	@echo "📋 To check status: make k8s-status"
	@echo "📋 To view logs:    kubectl logs -n pixels-staging deployment/pixels -f"

# Deploy to production environment
deploy-production: k8s-create-namespaces k8s-create-secrets
	@echo "Deploying to production..."
	kubectl apply -k k8s/overlays/production/
	@echo "✅ Production deployment complete"
	@echo ""
	@echo "🔗 Access production at:"
	@echo "   http://10.0.0.130:30080"
	@echo "   http://10.0.0.99:30080"
	@echo ""
	@echo "📋 To check status: make k8s-status"
	@echo "📋 To view logs:    kubectl logs -n pixels-production deployment/pixels -f"

# Tear down staging environment
teardown-staging:
	@echo "Tearing down staging environment..."
	kubectl delete namespace pixels-staging --ignore-not-found=true
	@echo "✅ Staging environment torn down"

# Tear down production environment  
teardown-production:
	@echo "Tearing down production environment..."
	kubectl delete namespace pixels-production --ignore-not-found=true
	@echo "✅ Production environment torn down"

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
	@echo "🚀 Production Pods:"
	@kubectl get pods -n pixels-production 2>/dev/null || echo "   No production namespace"
	@echo ""
	@echo "🌐 Services:"
	@kubectl get services -n pixels-staging 2>/dev/null || echo "   No staging services"
	@kubectl get services -n pixels-production 2>/dev/null || echo "   No production services"
	@echo ""
	@echo "🔗 Access URLs:"
	@echo "   Staging:    http://10.0.0.130:30080 | http://10.0.0.99:30080"
	@echo "   Production: http://10.0.0.130:30080 | http://10.0.0.99:30080"

# Help command
help:
	@echo "Available commands:"
	@echo ""
	@echo "🔨 Development:"
	@echo "   build              - Build the service"
	@echo "   run                - Run the service locally"
	@echo "   test               - Run tests"
	@echo "   clean              - Clean build artifacts"
	@echo "   dev                - Run with air for development"
	@echo ""
	@echo "☸️  Kubernetes:"
	@echo "   deploy-staging     - Deploy to staging environment"
	@echo "   deploy-production  - Deploy to production environment"
	@echo "   teardown-staging   - Tear down staging environment"
	@echo "   teardown-production - Tear down production environment"
	@echo "   k8s-status         - Show cluster status"
	@echo ""
	@echo "📋 Required environment variables for K8s:"
	@echo "   GITHUB_USERNAME    - Your GitHub username"
	@echo "   GITHUB_PAT         - Your GitHub Personal Access Token"
	@echo "   GITHUB_EMAIL       - Your GitHub email"
	@echo ""
	@echo "💡 Example usage:"
	@echo "   source .env && make deploy-staging"
