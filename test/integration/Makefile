# Makefile para Testes de Integração NFT Access

.PHONY: test-integration test-nation-pass test-contracts test-startup test-all help

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Configurações padrão
TEST_TIMEOUT=30s
VERBOSE=-v

help: ## Mostra esta ajuda
	@echo "$(GREEN)Testes de Integração NFT Access - Nation Pass$(NC)"
	@echo ""
	@echo "$(YELLOW)Comandos disponíveis:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Variáveis de ambiente necessárias:$(NC)"
	@echo "  WALLET_ADDRESS        - Wallet com NFT Nation Pass"
	@echo "  NATION_NFT_CONTRACT   - Endereço do contrato Nation Pass"
	@echo "  BASE_RPC_URL          - RPC da Base Network"
	@echo ""
	@echo "$(YELLOW)Exemplo de uso:$(NC)"
	@echo "  make test-integration"
	@echo "  make test-nation-pass"
	@echo "  make test-all"

test-integration: ## Executa testes básicos de integração
	@echo "$(GREEN)🧪 Executando testes básicos de integração...$(NC)"
	@INTEGRATION_TESTS=true go test ./test/integration/ $(VERBOSE) -run "TestNFTAccessIntegration" -timeout $(TEST_TIMEOUT)

test-contracts: ## Executa testes com contratos reais
	@echo "$(GREEN)🔗 Executando testes com contratos reais...$(NC)"
	@REAL_CONTRACT_TESTS=true go test ./test/integration/ $(VERBOSE) -run "TestNFTAccessManagerRealContract" -timeout $(TEST_TIMEOUT)

test-nation-pass: ## Executa testes específicos Nation Pass
	@echo "$(GREEN)🎨 Executando testes específicos Nation Pass...$(NC)"
	@NATION_PASS_TESTS=true go test ./test/integration/ $(VERBOSE) -run "TestNationPassAccessValidation" -timeout $(TEST_TIMEOUT)

test-startup: ## Executa testes de validação de startup
	@echo "$(GREEN)🚀 Executando testes de validação de startup...$(NC)"
	@STARTUP_VALIDATION_TESTS=true go test ./test/integration/ $(VERBOSE) -run "TestStartupValidationIntegration" -timeout $(TEST_TIMEOUT)

test-tiers: ## Executa testes de validação de tiers
	@echo "$(GREEN)📊 Executando testes de validação de tiers...$(NC)"
	@TIER_VALIDATION_TESTS=true go test ./test/integration/ $(VERBOSE) -run "TestNationPassTierValidation" -timeout $(TEST_TIMEOUT)

test-flow: ## Executa testes de fluxo completo
	@echo "$(GREEN)🔄 Executando testes de fluxo completo...$(NC)"
	@NATION_PASS_FLOW_TESTS=true go test ./test/integration/ $(VERBOSE) -run "TestNationPassIntegrationFlow" -timeout $(TEST_TIMEOUT)

test-all: ## Executa todos os testes de integração
	@echo "$(GREEN)🎯 Executando todos os testes de integração...$(NC)"
	@INTEGRATION_TESTS=true REAL_CONTRACT_TESTS=true NATION_PASS_TESTS=true STARTUP_VALIDATION_TESTS=true TIER_VALIDATION_TESTS=true NATION_PASS_FLOW_TESTS=true go test ./test/integration/ $(VERBOSE) -timeout $(TEST_TIMEOUT)

test-quick: ## Executa apenas testes básicos (mais rápidos)
	@echo "$(GREEN)⚡ Executando testes rápidos...$(NC)"
	@INTEGRATION_TESTS=true go test ./test/integration/ $(VERBOSE) -run "TestNFTAccessIntegration" -timeout 10s

test-verbose: ## Executa testes com output detalhado
	@echo "$(GREEN)📝 Executando testes com output detalhado...$(NC)"
	@INTEGRATION_TESTS=true REAL_CONTRACT_TESTS=true NATION_PASS_TESTS=true go test ./test/integration/ -v -timeout $(TEST_TIMEOUT)

check-env: ## Verifica se as variáveis de ambiente estão configuradas
	@echo "$(GREEN)🔍 Verificando variáveis de ambiente...$(NC)"
	@if [ -z "$$WALLET_ADDRESS" ]; then echo "$(RED)❌ WALLET_ADDRESS não configurada$(NC)"; exit 1; fi
	@if [ -z "$$NATION_NFT_CONTRACT" ]; then echo "$(RED)❌ NATION_NFT_CONTRACT não configurada$(NC)"; exit 1; fi
	@if [ -z "$$BASE_RPC_URL" ]; then echo "$(RED)❌ BASE_RPC_URL não configurada$(NC)"; exit 1; fi
	@echo "$(GREEN)✅ Todas as variáveis obrigatórias estão configuradas$(NC)"

setup-test-env: ## Configura ambiente de teste
	@echo "$(GREEN)⚙️ Configurando ambiente de teste...$(NC)"
	@if [ ! -f ".env" ]; then \
		echo "$(YELLOW)⚠️ Arquivo .env não encontrado, copiando env.example...$(NC)"; \
		cp env.example .env; \
		echo "$(YELLOW)⚠️ Configure suas variáveis no arquivo .env$(NC)"; \
	else \
		echo "$(GREEN)✅ Arquivo .env encontrado$(NC)"; \
	fi

clean-test-cache: ## Limpa cache de testes
	@echo "$(GREEN)🧹 Limpando cache de testes...$(NC)"
	@go clean -testcache
	@go clean -cache

test-coverage: ## Executa testes com cobertura
	@echo "$(GREEN)📊 Executando testes com cobertura...$(NC)"
	@INTEGRATION_TESTS=true REAL_CONTRACT_TESTS=true NATION_PASS_TESTS=true go test ./test/integration/ -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Relatório de cobertura gerado: coverage.html$(NC)"

test-benchmark: ## Executa benchmarks de performance
	@echo "$(GREEN)🏃 Executando benchmarks...$(NC)"
	@INTEGRATION_TESTS=true go test ./test/integration/ -bench=. -benchmem

test-race: ## Executa testes com detecção de race conditions
	@echo "$(GREEN)🏁 Executando testes com detecção de race conditions...$(NC)"
	@INTEGRATION_TESTS=true go test ./test/integration/ -race -timeout $(TEST_TIMEOUT)

# Comandos de desenvolvimento
dev-setup: ## Configura ambiente de desenvolvimento
	@echo "$(GREEN)🛠️ Configurando ambiente de desenvolvimento...$(NC)"
	@make setup-test-env
	@make check-env
	@echo "$(GREEN)✅ Ambiente de desenvolvimento configurado$(NC)"

dev-test: ## Executa testes de desenvolvimento (rápidos)
	@echo "$(GREEN)🔧 Executando testes de desenvolvimento...$(NC)"
	@make test-quick

# Comandos de CI/CD
ci-test: ## Executa testes para CI/CD
	@echo "$(GREEN)🤖 Executando testes para CI/CD...$(NC)"
	@make test-all

ci-test-quick: ## Executa testes rápidos para CI/CD
	@echo "$(GREEN)⚡ Executando testes rápidos para CI/CD...$(NC)"
	@make test-quick

# Comandos de produção
prod-test: ## Executa testes de produção (completos)
	@echo "$(GREEN)🏭 Executando testes de produção...$(NC)"
	@make test-all
	@make test-coverage

# Comandos de troubleshooting
debug-test: ## Executa testes em modo debug
	@echo "$(GREEN)🐛 Executando testes em modo debug...$(NC)"
	@LOG_LEVEL=debug INTEGRATION_TESTS=true go test ./test/integration/ -v -timeout $(TEST_TIMEOUT)

test-specific: ## Executa teste específico (use TEST_NAME=TestName)
	@echo "$(GREEN)🎯 Executando teste específico: $(TEST_NAME)$(NC)"
	@INTEGRATION_TESTS=true go test ./test/integration/ $(VERBOSE) -run "$(TEST_NAME)" -timeout $(TEST_TIMEOUT)

# Comandos de documentação
test-docs: ## Gera documentação dos testes
	@echo "$(GREEN)📚 Gerando documentação dos testes...$(NC)"
	@echo "Testes de Integração NFT Access - Nation Pass" > test-results.md
	@echo "=============================================" >> test-results.md
	@echo "" >> test-results.md
	@echo "Executado em: $$(date)" >> test-results.md
	@echo "" >> test-results.md
	@make test-all >> test-results.md 2>&1
	@echo "$(GREEN)✅ Documentação gerada: test-results.md$(NC)"

# Comandos de limpeza
clean: ## Limpa arquivos temporários
	@echo "$(GREEN)🧹 Limpando arquivos temporários...$(NC)"
	@rm -f coverage.out coverage.html test-results.md
	@make clean-test-cache

# Comandos de verificação
verify: ## Verifica se tudo está funcionando
	@echo "$(GREEN)✅ Verificando sistema...$(NC)"
	@make check-env
	@make test-quick
	@echo "$(GREEN)🎉 Sistema funcionando corretamente!$(NC)"

# Comandos de status
status: ## Mostra status do sistema
	@echo "$(GREEN)📊 Status do sistema:$(NC)"
	@echo "  Wallet Address: $$WALLET_ADDRESS"
	@echo "  Nation Contract: $$NATION_NFT_CONTRACT"
	@echo "  Base RPC: $$BASE_RPC_URL"
	@echo "  Go Version: $$(go version)"
	@echo "  Test Timeout: $(TEST_TIMEOUT)"
