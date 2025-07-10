# Finsolvz Backend Makefile
# Comprehensive testing and development commands

.PHONY: help test test-unit test-integration test-e2e test-all test-coverage test-performance build run clean lint format docker-build docker-run setup-test-db

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

# Default target
help: ## Show this help message
	@echo "$(BLUE)Finsolvz Backend - Available Commands$(NC)"
	@echo "======================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development commands
build: ## Build the application
	@echo "$(BLUE)Building Finsolvz Backend...$(NC)"
	go build -o bin/finsolvz-backend cmd/server/main.go
	@echo "$(GREEN)✅ Build completed$(NC)"

run: ## Run the application locally
	@echo "$(BLUE)Starting Finsolvz Backend...$(NC)"
	go run cmd/server/main.go

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	go clean
	@echo "$(GREEN)✅ Clean completed$(NC)"

# Code quality
lint: ## Run linting
	@echo "$(BLUE)Running linters...$(NC)"
	golangci-lint run --timeout=5m
	@echo "$(GREEN)✅ Linting completed$(NC)"

format: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	goimports -w .
	@echo "$(GREEN)✅ Formatting completed$(NC)"

# Testing commands
test: test-unit ## Run unit tests (default)

test-unit: ## Run unit tests only
	@echo "$(BLUE)Running unit tests...$(NC)"
	go test -v -race -timeout=30s ./internal/app/...
	@echo "$(GREEN)✅ Unit tests completed$(NC)"

test-integration: setup-test-db ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(NC)"
	@echo "$(YELLOW)Note: Requires MongoDB running on localhost:27017$(NC)"
	go test -v -race -timeout=60s ./tests -run "TestIntegration"
	@echo "$(GREEN)✅ Integration tests completed$(NC)"

test-e2e: ## Run E2E tests against live server
	@echo "$(BLUE)Running E2E tests...$(NC)"
	@if [ -z "$(FINSOLVZ_E2E_URL)" ]; then \
		echo "$(YELLOW)Set FINSOLVZ_E2E_URL environment variable to run E2E tests$(NC)"; \
		echo "Example: make test-e2e FINSOLVZ_E2E_URL=https://your-service.a.run.app"; \
	else \
		echo "Testing against: $(FINSOLVZ_E2E_URL)"; \
		FINSOLVZ_E2E_URL=$(FINSOLVZ_E2E_URL) go test -v -timeout=120s ./tests -run "TestE2E"; \
	fi
	@echo "$(GREEN)✅ E2E tests completed$(NC)"

test-all: test-unit test-integration ## Run all tests (unit + integration)
	@echo "$(GREEN)✅ All tests completed$(NC)"

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	go test -v -race -coverprofile=coverage.out ./internal/app/...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

test-performance: ## Run performance/benchmark tests
	@echo "$(BLUE)Running performance tests...$(NC)"
	go test -v -bench=. -benchmem ./internal/app/...
	@if [ ! -z "$(FINSOLVZ_E2E_URL)" ]; then \
		echo "$(BLUE)Running E2E benchmarks...$(NC)"; \
		FINSOLVZ_E2E_URL=$(FINSOLVZ_E2E_URL) go test -bench=BenchmarkE2E -benchmem ./tests; \
	fi
	@echo "$(GREEN)✅ Performance tests completed$(NC)"

# Database setup
setup-test-db: ## Setup test database (MongoDB)
	@echo "$(BLUE)Setting up test database...$(NC)"
	@if command -v mongod >/dev/null 2>&1; then \
		echo "$(GREEN)✅ MongoDB found$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  MongoDB not found. Install MongoDB or use Docker:$(NC)"; \
		echo "   docker run -d --name mongo-test -p 27017:27017 mongo:7.0"; \
	fi

# Docker commands
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t finsolvz-backend:latest .
	@echo "$(GREEN)✅ Docker image built$(NC)"

docker-run: ## Run Docker container
	@echo "$(BLUE)Running Docker container...$(NC)"
	docker run -d --name finsolvz-backend -p 8787:8787 \
		-e MONGO_URI=mongodb://host.docker.internal:27017/Finsolvz \
		-e JWT_SECRET=docker-secret-key \
		finsolvz-backend:latest
	@echo "$(GREEN)✅ Docker container started on http://localhost:8787$(NC)"

docker-stop: ## Stop Docker container
	@echo "$(BLUE)Stopping Docker container...$(NC)"
	docker stop finsolvz-backend || true
	docker rm finsolvz-backend || true
	@echo "$(GREEN)✅ Docker container stopped$(NC)"

# GCP commands
gcp-deploy: ## Deploy to Google Cloud Run
	@echo "$(BLUE)Deploying to Google Cloud Run...$(NC)"
	@if [ -z "$(PROJECT_ID)" ]; then \
		echo "$(RED)❌ PROJECT_ID environment variable is required$(NC)"; \
		echo "Usage: make gcp-deploy PROJECT_ID=your-project-id"; \
		exit 1; \
	fi
	gcloud builds submit --project=$(PROJECT_ID)
	@echo "$(GREEN)✅ Deployment completed$(NC)"

gcp-setup: ## Setup GCP environment
	@echo "$(BLUE)Setting up GCP environment...$(NC)"
	@if [ -z "$(PROJECT_ID)" ] || [ -z "$(GITHUB_USER)" ]; then \
		echo "$(RED)❌ PROJECT_ID and GITHUB_USER are required$(NC)"; \
		echo "Usage: make gcp-setup PROJECT_ID=your-project GITHUB_USER=your-username"; \
		exit 1; \
	fi
	./setup-gcp-environment.sh $(PROJECT_ID) $(GITHUB_USER)
	@echo "$(GREEN)✅ GCP setup completed$(NC)"

# Testing workflows
test-quick: ## Quick test (unit tests only, no race detection)
	@echo "$(BLUE)Running quick tests...$(NC)"
	go test -timeout=15s ./internal/app/...
	@echo "$(GREEN)✅ Quick tests completed$(NC)"

test-ci: ## CI/CD test pipeline
	@echo "$(BLUE)Running CI/CD test pipeline...$(NC)"
	@echo "$(YELLOW)1. Formatting check...$(NC)"
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "$(RED)❌ Code not formatted. Run: make format$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Code formatting OK$(NC)"
	
	@echo "$(YELLOW)2. Linting...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=3m; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found, skipping$(NC)"; \
	fi
	
	@echo "$(YELLOW)3. Unit tests with race detection...$(NC)"
	go test -v -race -timeout=60s ./internal/app/...
	
	@echo "$(YELLOW)4. Build test...$(NC)"
	go build -o /tmp/finsolvz-test cmd/server/main.go
	rm -f /tmp/finsolvz-test
	
	@echo "$(GREEN)✅ CI/CD pipeline completed$(NC)"

# Performance testing
perf-test: ## Run performance tests against live service
	@echo "$(BLUE)Running performance tests...$(NC)"
	@if [ -z "$(SERVICE_URL)" ]; then \
		echo "$(YELLOW)Set SERVICE_URL to test against live service$(NC)"; \
		echo "Example: make perf-test SERVICE_URL=https://your-service.a.run.app"; \
	else \
		./performance-test.sh $(SERVICE_URL); \
	fi

# Development helpers
dev-setup: ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@echo "$(YELLOW)1. Installing Go dependencies...$(NC)"
	go mod download
	go mod tidy
	
	@echo "$(YELLOW)2. Installing development tools...$(NC)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	fi
	
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	
	@echo "$(YELLOW)3. Creating .env file...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.example .env || true; \
		echo "✅ .env file created from .env.example"; \
	fi
	
	@echo "$(GREEN)✅ Development environment setup completed$(NC)"

# Show test status
test-status: ## Show testing status and coverage
	@echo "$(BLUE)Testing Status$(NC)"
	@echo "=============="
	@echo "$(YELLOW)Unit Tests:$(NC)"
	@find ./internal/app -name '*_test.go' | wc -l | xargs -I {} echo "  Test files: {}"
	@echo "$(YELLOW)Integration Tests:$(NC)"
	@find ./tests -name '*integration*' | wc -l | xargs -I {} echo "  Test files: {}"
	@echo "$(YELLOW)E2E Tests:$(NC)"
	@find ./tests -name '*e2e*' | wc -l | xargs -I {} echo "  Test files: {}"
	@if [ -f coverage.out ]; then \
		echo "$(YELLOW)Coverage:$(NC)"; \
		go tool cover -func=coverage.out | tail -1; \
	else \
		echo "$(YELLOW)Coverage:$(NC) Run 'make test-coverage' to generate"; \
	fi
