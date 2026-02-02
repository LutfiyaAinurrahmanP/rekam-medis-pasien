.PHONY: help build run test test-users test-coverage clean docker-build docker-up docker-down migrate lint fmt mod-download test-all mod-tidy dev prod-build

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o main cmd/api/main.go

run: ## Run the application
	@echo "Running application..."
	@go run cmd/api/main.go

test: ## Run all tests with Jest-like output
	@powershell -File scripts/test-runner.ps1 -Path "./..." -Filter ""

test-users: ## Run user API tests with Jest-like output
	@powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_.*User"

test-departments: ## Run department API tests with Jest-like output
	@powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_.*Department"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

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

docker-down: ## Stop docker containers
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

test-all: ## Test all API
	@echo "Test all API"
	@go test -v ./internal/routes/... -count=1

mod-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

dev: docker-up ## Start development environment
	@echo "Development environment started!"

prod-build: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go
