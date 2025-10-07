.PHONY: help build run test clean docker-build docker-run setup

# Variables
APP_NAME=iac-ai-agent
DOCKER_IMAGE=iac-ai-agent:latest
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify
	@pip3 install checkov || echo "Warning: checkov not installed"

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) ./cmd/agent

run: build ## Run the application
	@echo "Running $(APP_NAME)..."
	@./bin/$(APP_NAME)

run-swagger: swagger build ## Generate Swagger docs and run the application
	@echo "Running $(APP_NAME) with Swagger UI..."
	@./bin/$(APP_NAME)

test: ## Run all tests
	@echo "Running all tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-unit: ## Run unit tests (BDD)
	@echo "Running unit tests..."
	@go test -v -race ./test/unit/

test-integration: ## Run integration tests (BDD)
	@echo "Running integration tests..."
	@go test -v -race ./test/integration/

test-bdd: ## Run all BDD tests with Ginkgo
	@echo "Running BDD tests with Ginkgo..."
	@which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@ginkgo -r ./test/

test-bdd-verbose: ## Run BDD tests with verbose output
	@echo "Running BDD tests (verbose)..."
	@which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@ginkgo -v -r ./test/

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out

test-coverage-bdd: ## Run BDD tests with coverage
	@echo "Running BDD tests with coverage..."
	@which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@ginkgo -cover -coverprofile=coverage.out -r ./test/
	@go tool cover -html=coverage.out

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed"

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -w $(GO_FILES)
	@go mod tidy

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@v1.8.12
	@swag init -g cmd/agent/main.go -o docs --parseDependency --parseInternal
	@echo "✅ Swagger docs generated! Access: http://localhost:8080/swagger/"

swagger-install: ## Install Swagger CLI
	@echo "Installing Swagger CLI..."
	@go install github.com/swaggo/swag/cmd/swag@v1.8.12
	@echo "✅ Swagger CLI installed!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf docs/
	@rm -f coverage.out

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -f deployments/Dockerfile -t $(DOCKER_IMAGE) .

docker-run: docker-build ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services..."
	@docker-compose -f configs/docker-compose.yml up -d

docker-compose-down: ## Stop services
	@echo "Stopping services..."
	@docker-compose -f configs/docker-compose.yml down

install: build ## Install binary to $GOPATH/bin
	@echo "Installing $(APP_NAME)..."
	@cp bin/$(APP_NAME) $(GOPATH)/bin/

dev: ## Run in development mode with hot reload
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air

.DEFAULT_GOAL := help
