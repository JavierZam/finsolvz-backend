# Makefile
# This line MUST be the very first non-comment line in your Makefile
# It loads environment variables from .env file directly into Makefile's scope.
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build run test clean lint migrate-up migrate-down docker-build docker-run-db docker-stop-db help docker-reset-db

APP_NAME := finsolvz-backend
BINARY_NAME := $(APP_NAME)
BUILD_DIR := bin

# Go related variables
GO_VERSION := 1.22
GO_LINTER := golangci-lint
GO_MIGRATE := migrate

# Docker variables for local PostgreSQL
DOCKER_IMAGE := $(APP_NAME)
DOCKER_CONTAINER_NAME := $(APP_NAME)-db
DOCKER_DB_PORT := 5432:5432
DOCKER_DB_PASSWORD := $(DB_PASSWORD)
DOCKER_DB_USER := $(DB_USER)
DOCKER_DB_NAME := $(DB_NAME)

# Default target
all: build

# Build the application binary
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Run linters (installs if not present)
lint:
	@echo "Running linters..."
	@if ! command -v $(GO_LINTER) &> /dev/null; then \
		echo "$(GO_LINTER) not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	$(GO_LINTER) run ./...

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run-db: docker-stop-db
	@echo "Starting PostgreSQL Docker container..."
	docker run --name $(DOCKER_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DOCKER_DB_PASSWORD) -e POSTGRES_USER=$(DOCKER_DB_USER) -e POSTGRES_DB=$(DOCKER_DB_NAME) -p $(DOCKER_DB_PORT) -d postgres:15-alpine
	@echo "Waiting for database to be responsive..."
	# Loop until the database is ready to accept connections
	until docker exec $(DOCKER_CONTAINER_NAME) pg_isready -U $(DOCKER_DB_USER) -d $(DOCKER_DB_NAME); do \
		sleep 1; \
	done;
	@echo "Database is ready."
	# Removed direct psql execution here. Migrations will be handled by migrate-up/down.

docker-stop-db:
	@echo "Stopping PostgreSQL Docker container..."
	docker stop $(DOCKER_CONTAINER_NAME) || true
	docker rm $(DOCKER_CONTAINER_NAME) || true

docker-reset-db: docker-stop-db
	@echo "Pruning unused Docker volumes (this might remove volumes from other projects if not careful!)..."
	docker volume prune -f || true
	@echo "Database has been completely reset."
	make docker-run-db # Start a fresh database container after pruning
	make migrate-up # Apply migrations after a clean database start


# Database Migrations (golang-migrate CLI)
migrate-install:
	@echo "Installing migrate CLI tool..."
	@if ! command -v $(GO_MIGRATE) &> /dev/null; then \
		echo "$(GO_MIGRATE) not found. Installing..."; \
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz; \
		sudo mv migrate /usr/local/bin/; \
	fi

migrate-up: migrate-install
	@echo "Running database migrations up..."
	migrate -path migrations -database "postgresql://$(DOCKER_DB_USER):$(DOCKER_DB_PASSWORD)@localhost:$(shell echo $(DOCKER_DB_PORT) | cut -d':' -f1)/$(DOCKER_DB_NAME)?sslmode=disable" up

migrate-down: migrate-install
	@echo "Running database migrations down..."
	migrate -path migrations -database "postgresql://$(DOCKER_DB_USER):$(DOCKER_DB_PASSWORD)@localhost:$(shell echo $(DOCKER_DB_PORT) | cut -d':' -f1)/$(DOCKER_DB_NAME)?sslmode=disable" down 1

# Help message
help:
	@echo "Usage:"
	@echo "  make build          Builds the application binary"
	@echo "  make run            Runs the application"
	@echo "  make test           Runs all tests"
	@echo "  make clean          Removes build artifacts"
	@echo "  make lint           Runs linters"
	@echo "  make docker-build   Builds the Docker image"
	@echo "  make docker-run-db  Starts a local PostgreSQL database in Docker (stops existing one first)."
	@echo "  make docker-stop-db Stops and removes the local PostgreSQL database"
	@echo "  make docker-reset-db Resets the database completely (stops, removes, prunes) and starts a new one with migrations."
	@echo "  make migrate-up     Applies database migrations using golang-migrate CLI."
	@echo "  make migrate-down   Reverts the last database migration using golang-migrate CLI."
	@echo "  make help           Displays this help message"
