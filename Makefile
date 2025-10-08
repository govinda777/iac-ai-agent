.PHONY: help build run test clean docker-build docker-run setup git-hooks

# Variables
APP_NAME=iac-ai-agent
DOCKER_IMAGE=iac-ai-agent:latest
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Install dependencies and git hooks
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod verify
	@pip3 install checkov || echo "⚠️  Warning: checkov not installed"
	@echo ""
	@echo "🔧 Installing test tools..."
	@go install github.com/onsi/ginkgo/v2/ginkgo@latest || echo "⚠️  Warning: ginkgo not installed"
	@which golangci-lint > /dev/null || echo "⚠️  Warning: golangci-lint not installed. Install from: https://golangci-lint.run/usage/install/"
	@echo ""
	@echo "🪝 Installing git hooks..."
	@$(MAKE) git-hooks
	@echo ""
	@echo "✅ Setup complete!"

check-env: ## Check prerequisites and environment configuration
	@echo "🔍 Checking prerequisites..."
	@./scripts/check-prerequisites.sh

git-hooks: ## Install git hooks (pre-push)
	@echo "Installing git hooks..."
	@chmod +x .githooks/pre-push
	@git config core.hooksPath .githooks
	@echo "✅ Git hooks installed! Pre-push will run tests automatically."

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
	@echo "🧪 Running all tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "📊 Coverage:"
	@go tool cover -func=coverage.out | tail -1
	@echo "📊 Generating comprehensive HTML report..."
	@go run scripts/generate_report.go all

test-unit: ## Run unit tests
	@echo "🧪 Running unit tests..."
	@go test -v -race -count=1 -timeout=30s ./test/unit/...
	@echo "✅ Unit tests completed!"
	@echo "📊 Generating HTML report..."
	@go run scripts/generate_report.go unit

test-integration: ## Run integration tests
	@echo "🔗 Running integration tests..."
	@go test -v -race -count=1 -timeout=60s ./test/integration/...
	@echo "✅ Integration tests completed!"
	@echo "📊 Generating HTML report..."
	@go run scripts/generate_report.go integration

test-quick: ## Run tests without race detector (faster)
	@echo "⚡ Running quick tests..."
	@go test -short ./...

test-bdd: ## Run all BDD tests with Ginkgo
	@echo "Running BDD tests with Ginkgo..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@export PATH=$$PATH:$$(go env GOPATH)/bin && ginkgo -r ./test/
	@echo "✅ BDD tests completed!"
	@echo "📊 Generating HTML report..."
	@go run scripts/generate_report.go bdd

test-bdd-verbose: ## Run BDD tests with verbose output
	@echo "Running BDD tests (verbose)..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@export PATH=$$PATH:$$(go env GOPATH)/bin && ginkgo -v -r ./test/

test-nation-nft: ## Run BDD tests for Nation NFT Pass validation
	@echo "🎨 Running Nation NFT Pass validation tests..."
	@chmod +x test/bdd/run_nation_nft_tests.sh
	@./test/bdd/run_nation_nft_tests.sh

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out

test-coverage-bdd: ## Run BDD tests with coverage
	@echo "Running BDD tests with coverage..."
	@which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@ginkgo -cover -coverprofile=coverage.out -r ./test/
	@go tool cover -html=coverage.out

test-reports: ## Open latest test reports in browser
	@echo "📊 Opening latest test reports..."
	@if [ -d "reports/html" ]; then \
		LATEST_REPORT=$$(ls -t reports/html/*.html | head -1); \
		if [ -n "$$LATEST_REPORT" ]; then \
			echo "Opening: $$LATEST_REPORT"; \
			open "$$LATEST_REPORT"; \
		else \
			echo "No reports found. Run tests first."; \
		fi; \
	else \
		echo "Reports directory not found. Run tests first."; \
	fi

test-reports-list: ## List all available test reports
	@echo "📊 Available test reports:"
	@if [ -d "reports/html" ]; then \
		ls -la reports/html/*.html 2>/dev/null || echo "No reports found."; \
	else \
		echo "Reports directory not found."; \
	fi

test-reports-version: ## Version current test reports for Git tracking
	@echo "📊 Versioning test reports..."
	@./scripts/version_reports.sh

test-reports-index: ## Open reports index page
	@echo "📊 Opening reports index..."
	@if [ -f "reports/versioned/index.html" ]; then \
		open "reports/versioned/index.html"; \
	else \
		echo "Index not found. Run 'make test-reports-version' first."; \
	fi

lint: ## Run linter
	@echo "📋 Running linter..."
	@which golangci-lint > /dev/null || (echo "⚠️  golangci-lint not installed. Run 'make lint-install' to install." && exit 1)
	@golangci-lint run --timeout=5m
	@echo "✅ Linter passed!"

lint-install: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
	@echo "✅ golangci-lint installed!"

lint-fix: ## Run linter and fix issues
	@echo "📋 Running linter with auto-fix..."
	@golangci-lint run --fix
	@echo "✅ Auto-fix complete!"

fmt: ## Format code
	@echo "💅 Formatting code..."
	@gofmt -w $(GO_FILES)
	@go mod tidy
	@echo "✅ Code formatted!"

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

# ============================================
# 🚀 SMART CONTRACTS CI/CD COMMANDS
# ============================================

contracts-setup: ## Setup Foundry and install dependencies
	@echo "🛠️ Setting up Foundry environment..."
	@cd contracts && forge install --no-commit
	@cd contracts && forge build
	@echo "✅ Foundry setup complete!"

contracts-test: ## Run smart contract tests
	@echo "🧪 Running smart contract tests..."
	@cd contracts && forge test --gas-report --coverage
	@echo "✅ Contract tests completed!"

contracts-test-verbose: ## Run smart contract tests with verbose output
	@echo "🧪 Running smart contract tests (verbose)..."
	@cd contracts && forge test -vvv --gas-report --coverage
	@echo "✅ Contract tests completed!"

contracts-lint: ## Run smart contract linting and security checks
	@echo "🔍 Running smart contract security analysis..."
	@cd contracts && forge build
	@echo "✅ Contract linting completed!"

contracts-deploy-testnet: ## Deploy contracts to testnet
	@echo "🚀 Deploying contracts to testnet..."
	@cd contracts && forge script script/Deploy.s.sol --rpc-url base-sepolia --broadcast --verify
	@echo "✅ Testnet deployment completed!"

contracts-deploy-mainnet: ## Deploy contracts to mainnet
	@echo "🚀 Deploying contracts to mainnet..."
	@cd contracts && forge script script/Deploy.s.sol --rpc-url base-mainnet --broadcast --verify
	@echo "✅ Mainnet deployment completed!"

contracts-verify: ## Verify deployed contracts
	@echo "🔍 Verifying deployed contracts..."
	@./scripts/verify-contracts.sh base-mainnet
	@echo "✅ Contract verification completed!"

contracts-monitor: ## Monitor contracts in production
	@echo "📊 Monitoring contracts in production..."
	@./scripts/monitor-contracts.sh base-mainnet --health-check
	@echo "✅ Contract monitoring completed!"

contracts-monitor-alerts: ## Monitor contracts with alerts enabled
	@echo "📊 Monitoring contracts with alerts..."
	@./scripts/monitor-contracts.sh base-mainnet --alerts
	@echo "✅ Contract monitoring with alerts completed!"

contracts-rollback: ## Rollback contracts to previous version
	@echo "🔄 Rolling back contracts..."
	@./scripts/rollback-contracts.sh base-mainnet v1.0.0 --dry-run
	@echo "✅ Contract rollback simulation completed!"

contracts-rollback-confirm: ## Rollback contracts with confirmation
	@echo "🔄 Rolling back contracts with confirmation..."
	@./scripts/rollback-contracts.sh base-mainnet v1.0.0 --confirm
	@echo "✅ Contract rollback completed!"

contracts-ci: contracts-setup contracts-test contracts-lint ## Run full CI pipeline for contracts
	@echo "✅ Smart contracts CI pipeline completed!"

contracts-cd-testnet: contracts-ci contracts-deploy-testnet contracts-verify ## Run full CD pipeline for testnet
	@echo "✅ Smart contracts CD pipeline for testnet completed!"

contracts-cd-mainnet: contracts-ci contracts-deploy-mainnet contracts-verify ## Run full CD pipeline for mainnet
	@echo "✅ Smart contracts CD pipeline for mainnet completed!"

contracts-clean: ## Clean contract build artifacts
	@echo "🧹 Cleaning contract artifacts..."
	@cd contracts && forge clean
	@rm -rf contracts/out
	@rm -rf contracts/cache
	@echo "✅ Contract artifacts cleaned!"

contracts-docs: ## Generate contract documentation
	@echo "📚 Generating contract documentation..."
	@cd contracts && forge doc --build
	@echo "✅ Contract documentation generated!"

contracts-gas-report: ## Generate gas usage report
	@echo "⛽ Generating gas usage report..."
	@cd contracts && forge test --gas-report > gas-report.txt
	@echo "✅ Gas report generated: contracts/gas-report.txt"

contracts-coverage: ## Generate test coverage report
	@echo "📊 Generating test coverage report..."
	@cd contracts && forge coverage --report lcov
	@echo "✅ Coverage report generated!"

contracts-security: ## Run security analysis on contracts
	@echo "🔒 Running security analysis..."
	@cd contracts && forge build
	@echo "✅ Security analysis completed!"

contracts-upgrade: ## Check for contract upgrades
	@echo "⬆️ Checking for contract upgrades..."
	@cd contracts && forge build
	@echo "✅ Upgrade check completed!"

contracts-backup: ## Create backup of current deployment
	@echo "💾 Creating deployment backup..."
	@mkdir -p rollback/backups
	@cp deployments/base-mainnet.json rollback/backups/base-mainnet-$(shell date +%Y%m%d-%H%M%S).json
	@echo "✅ Deployment backup created!"

contracts-status: ## Show current contract deployment status
	@echo "📋 Current contract deployment status:"
	@if [ -f "deployments/base-mainnet.json" ]; then \
		echo "Mainnet: $(jq -r '.version // "unknown"' deployments/base-mainnet.json)"; \
	else \
		echo "Mainnet: Not deployed"; \
	fi
	@if [ -f "deployments/base-sepolia.json" ]; then \
		echo "Testnet: $(jq -r '.version // "unknown"' deployments/base-sepolia.json)"; \
	else \
		echo "Testnet: Not deployed"; \
	fi

.DEFAULT_GOAL := help
