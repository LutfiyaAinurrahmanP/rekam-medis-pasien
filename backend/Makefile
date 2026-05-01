.PHONY: help build run seed test test-users test-departments test-patients test-doctors test-rooms test-type-tests test-unit test-unit-handler test-unit-service test-unit-repository test-unit-verbose test-unit-handler-verbose test-unit-service-verbose test-unit-repository-verbose test-unit-race test-unit-failfast test-unit-timeout test-unit-run test-unit-list test-unit-coverage test-unit-coverage-html test-unit-coverage-handler test-unit-coverage-service test-unit-coverage-repository test-coverage test-ci test-all-cleanup clean clean-coverage docker-build docker-up docker-up-deps docker-down docker-logs docker-clean migrate lint fmt mod-download test-all mod-tidy dev prod-build check-docker wait-db stop-app-container

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o main cmd/api/main.go

run: check-docker docker-up-deps stop-app-container wait-db ## Run local app with auto dependency setup
	@echo "Running application..."
	@go run cmd/api/main.go

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@go run cmd/seed/main.go

test: ## Run all tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./..." -Filter ""; \
	else \
		bash ./scripts/test-runner.sh -Path "./..." -Filter ""; \
	fi

test-users: ## Run user API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_.*User"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests" -Filter "^Test_.*User"; \
	fi

test-departments: ## Run department API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_.*Department"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests" -Filter "^Test_.*Department"; \
	fi

test-patients: ## Run patient API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_.*Patient"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests" -Filter "^Test_.*Patient"; \
	fi

test-doctors: ## Run doctor API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_.*Doctor"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests" -Filter "^Test_.*Doctor"; \
	fi

test-rooms: ## Run room API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_.*Room"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests" -Filter "^Test_.*Room"; \
	fi

test-type-tests: ## Run type test API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests" -Filter "^Test_TypeTest"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests" -Filter "^Test_TypeTest"; \
	fi

# ============================================================================
# UNIT TEST COMMANDS (NEW - Comprehensive Unit Testing)
# ============================================================================

test-unit: ## Run all unit tests (handler, service, repository)
	@echo "Running all unit tests..."
	@go test ./tests/unit/... --count=1

test-unit-handler: ## Run only handler layer unit tests
	@echo "Running handler unit tests..."
	@go test ./tests/unit/handler --count=1

test-unit-service: ## Run only service layer unit tests
	@echo "Running service unit tests..."
	@go test ./tests/unit/service --count=1

test-unit-repository: ## Run only repository layer unit tests
	@echo "Running repository unit tests..."
	@go test ./tests/unit/repository --count=1

test-unit-verbose: ## Run all unit tests with verbose output
	@echo "Running all unit tests with verbose output..."
	@go test ./tests/unit/... -v --count=1

test-unit-handler-verbose: ## Run handler unit tests with verbose output
	@echo "Running handler unit tests with verbose output..."
	@go test ./tests/unit/handler -v --count=1

test-unit-service-verbose: ## Run service unit tests with verbose output
	@echo "Running service unit tests with verbose output..."
	@go test ./tests/unit/service -v --count=1

test-unit-repository-verbose: ## Run repository unit tests with verbose output
	@echo "Running repository unit tests with verbose output..."
	@go test ./tests/unit/repository -v --count=1

test-unit-race: ## Run unit tests with race detector
	@echo "Running unit tests with race detector..."
	@go test ./tests/unit/... -race --count=1

test-unit-failfast: ## Run unit tests and stop at first failure
	@echo "Running unit tests in failfast mode..."
	@go test ./tests/unit/... -failfast --count=1

test-unit-timeout: ## Run unit tests with 60s timeout
	@echo "Running unit tests with custom timeout..."
	@go test ./tests/unit/... --timeout=60s --count=1

test-unit-run: ## Run specific unit test (usage: make test-unit-run TEST=TestGetMyProfile)
	@echo "Running specific test: $(TEST)"
	@go test ./tests/unit/... -run $(TEST) -v --count=1

test-unit-list: ## List all available unit tests
	@echo "Available unit tests:"
	@go test ./tests/unit/... -list .

test-unit-coverage: ## Generate unit test coverage report
	@echo "Generating unit test coverage report..."
	@go test ./tests/unit/... -coverprofile=coverage_unit.out --count=1
	@go tool cover -func=coverage_unit.out | tail -30
	@echo "Full coverage report saved to: coverage_unit.out"

test-unit-coverage-html: ## Generate HTML unit test coverage report
	@echo "Generating HTML unit test coverage report..."
	@go test ./tests/unit/... -coverprofile=coverage_unit.out --count=1
	@go tool cover -html=coverage_unit.out -o coverage_unit.html
	@echo "HTML coverage report generated: coverage_unit.html"

test-unit-coverage-handler: ## Coverage report for handler layer only
	@echo "Generating handler layer coverage..."
	@go test ./tests/unit/handler -coverprofile=coverage_handler.out --count=1
	@go tool cover -func=coverage_handler.out

test-unit-coverage-service: ## Coverage report for service layer only
	@echo "Generating service layer coverage..."
	@go test ./tests/unit/service -coverprofile=coverage_service.out --count=1
	@go tool cover -func=coverage_service.out

test-unit-coverage-repository: ## Coverage report for repository layer only
	@echo "Generating repository layer coverage..."
	@go test ./tests/unit/repository -coverprofile=coverage_repository.out --count=1
	@go tool cover -func=coverage_repository.out

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
	@docker compose build

docker-up: ## Start docker containers
	@echo "Starting docker containers..."
	@docker compose up -d

docker-up-deps: ## Start only dependency containers (without app)
	@echo "Starting dependency containers (postgres, redis, kafka, zookeeper)..."
	@docker compose up -d postgres redis zookeeper kafka kafka-ui

docker-down: ## Stop docker containers
	@echo "Stopping docker containers..."
	@docker compose down

docker-logs: ## View docker logs
	@docker compose logs -f

docker-clean: ## Remove docker containers and volumes
	@echo "Removing docker containers and volumes..."
	@docker compose down -v

# ─── Helper targets ────────────────────────────────────────────────────────

check-docker: ## Check if docker and docker-compose are available
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "❌ Docker is not installed. Please install Docker first."; \
		exit 1; \
	fi
	@if ! docker compose version >/dev/null 2>&1; then \
		echo "❌ Docker Compose is not available. Please install Docker Compose first."; \
		exit 1; \
	fi
	@echo "✅ Docker and Docker Compose are available"

wait-db: ## Wait for database to be ready
	@echo "⏳ Waiting for database to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30; do \
		if nc -z 127.0.0.1 5432 2>/dev/null; then \
			echo "✅ Database is ready!"; \
			sleep 1; \
			exit 0; \
		fi; \
		echo "Attempt $$i/30 - Database not ready yet, waiting..."; \
		sleep 2; \
	done; \
	echo "❌ Database failed to start after 60 seconds"; \
	exit 1

stop-app-container: ## Stop docker app container to avoid port conflicts with local go run
	@docker compose stop app >/dev/null 2>&1 || true
	@docker compose rm -f app >/dev/null 2>&1 || true

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

test-ci: ## Run tests in CI mode (unit tests with coverage and race detection)
	@echo "Running CI test suite..."
	@go test ./tests/unit/... -race -coverprofile=coverage_unit.out --count=1
	@echo "✅ CI tests passed!"

test-all-cleanup: clean-coverage ## Clean all test artifacts
	@echo "All test artifacts cleaned!"

clean-coverage: ## Clean coverage files
	@echo "Cleaning coverage files..."
	@rm -f coverage.out coverage.html coverage_unit.out coverage_unit.html coverage_handler.out coverage_service.out coverage_repository.out

dev: docker-up-deps ## Start development dependencies only
	@echo "Development environment started!"

prod-build: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go
