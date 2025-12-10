.PHONY: help start stop restart build rebuild logs clean status shell test

# Default target
.DEFAULT_GOAL := help

# Variables
DOCKER_COMPOSE = docker-compose
CONTAINER_NAME = testing-studio
IMAGE_NAME = testing-studio:latest
PORT = 8888

##@ General

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Docker Operations

start: ## Start the application (builds if needed)
	@echo "🚀 Starting Testing Studio..."
	@if ! docker images -q $(IMAGE_NAME) 2> /dev/null | grep -q .; then \
		echo "📦 Image not found, building..."; \
		$(MAKE) build; \
	fi
	@$(DOCKER_COMPOSE) up -d
	@echo "✅ Testing Studio started on http://localhost:$(PORT)"
	@echo "📊 View logs: make logs"
	@echo "🛑 Stop: make stop"

stop: ## Stop the application
	@echo "🛑 Stopping Testing Studio..."
	@$(DOCKER_COMPOSE) down
	@echo "✅ Testing Studio stopped"

restart: ## Restart the application
	@echo "🔄 Restarting Testing Studio..."
	@$(DOCKER_COMPOSE) restart
	@echo "✅ Testing Studio restarted"

##@ Build Operations

build: ## Build the Docker image
	@echo "🔨 Building Testing Studio Docker image..."
	@$(DOCKER_COMPOSE) build --no-cache
	@echo "✅ Build complete"

rebuild: stop build start ## Rebuild and restart the application
	@echo "✅ Rebuild complete"

##@ Monitoring & Debug

logs: ## Show application logs (follow mode)
	@$(DOCKER_COMPOSE) logs -f testing-studio

logs-tail: ## Show last 50 lines of logs
	@$(DOCKER_COMPOSE) logs --tail=50 testing-studio

status: ## Show container status
	@echo "📊 Testing Studio Status:"
	@docker ps -a --filter "name=$(CONTAINER_NAME)" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
	@echo ""
	@echo "🐳 Image Info:"
	@docker images $(IMAGE_NAME) --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

shell: ## Open shell in running container
	@echo "🐚 Opening shell in $(CONTAINER_NAME)..."
	@docker exec -it $(CONTAINER_NAME) /bin/sh

##@ Cleanup

clean: stop ## Stop and remove containers, networks, images
	@echo "🧹 Cleaning up Testing Studio..."
	@$(DOCKER_COMPOSE) down -v --rmi local
	@echo "✅ Cleanup complete"

clean-all: ## Remove everything including volumes
	@echo "🧹 Deep cleaning Testing Studio..."
	@$(DOCKER_COMPOSE) down -v --rmi all
	@rm -f configs.json
	@echo "✅ Deep cleanup complete"

##@ Development

dev: ## Run in development mode (without Docker)
	@echo "🔧 Starting in development mode..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test ./... -v

fmt: ## Format code
	@echo "✨ Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run || echo "⚠️  Install golangci-lint: https://golangci-lint.run/usage/install/"

##@ Quick Actions

quick-start: build start ## Quick start (build + run)
	@echo "✅ Quick start complete"

open: ## Open Testing Studio in browser
	@echo "🌐 Opening Testing Studio..."
	@open http://localhost:$(PORT) || xdg-open http://localhost:$(PORT) || echo "Please open http://localhost:$(PORT) in your browser"

health: ## Check application health
	@echo "🏥 Checking health..."
	@curl -sf http://localhost:$(PORT)/api/configs > /dev/null && echo "✅ Healthy" || echo "❌ Unhealthy"

##@ Information

info: ## Show project information
	@echo "📋 Testing Studio Information"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Name:          Testing Studio (CloudEvents Explorer)"
	@echo "Version:       2.0.0"
	@echo "Port:          $(PORT)"
	@echo "Container:     $(CONTAINER_NAME)"
	@echo "Image:         $(IMAGE_NAME)"
	@echo "URL:           http://localhost:$(PORT)"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
