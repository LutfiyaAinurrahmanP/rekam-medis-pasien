. PHONY: help build run test clean docker-build docker-up docker-down migrate test-all

help:  ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o main cmd/api/main.go

run: ## Run the application
	@echo "Running application..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -f main
	@go clean

docker-build: ## Build docker images
	@echo "Building docker images..."
	@docker-compose build

docker-up: ## Start docker containers
	@echo "Starting docker containers..."
	@docker-compose up -d

docker-down:  ## Stop docker containers
	@echo "Stopping docker containers..."
	@docker-compose down

docker-logs: ## View docker logs
	@docker-compose logs -f

docker-clean: ## Remove docker containers and volumes
	@echo "Removing docker containers and volumes..."
	@docker-compose down -v

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/api/main.go

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

mod-download: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

test-all:
	@echo "Test all API"
	@go test -v ./internal/routes/... -count=1

mod-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

dev:  docker-up ## Start development environment
	@echo "Development environment started!"

prod-build: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go