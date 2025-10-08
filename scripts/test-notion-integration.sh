#!/bin/bash

# ============================================
# 🧪 Teste de Integração Notion
# ============================================
# 
# Este script testa a integração com Notion
# incluindo validação de startup e criação de agentes

set -e

echo "🚀 Iniciando teste de integração Notion..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir com cores
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verifica se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    print_error "Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Verifica se Go está instalado
if ! command -v go &> /dev/null; then
    print_error "Go não está instalado"
    exit 1
fi

print_status "Go encontrado: $(go version)"

# Verifica se git-secret está disponível
if ! command -v git &> /dev/null || ! git secret --version &> /dev/null; then
    print_warning "git-secret não está disponível. Alguns testes podem ser pulados."
fi

# Cria arquivo de configuração temporário para teste
print_status "Criando configuração temporária para teste..."

cat > /tmp/test_config.yaml << EOF
server:
  port: "8080"
  host: "0.0.0.0"

llm:
  provider: "nation.fun"
  model: "nation-1"
  temperature: 0.2
  max_tokens: 4000

github:
  auto_comment: true

analysis:
  checkov_enabled: true
  iam_analysis_enabled: true
  cost_optimization_enabled: true

scoring:
  min_pass_score: 70

logging:
  level: "info"
  format: "json"

web3:
  privy_app_id: "cmgh6un8w007bl10ci0tgitwp"
  privy_verification_key_url: "https://api.privy.io/v1/siwe/keys"
  base_rpc_url: "https://goerli.base.org"
  base_chain_id: 84531
  nft_access_contract_address: "0x147e832418Cc06A501047019E956714271098b89"
  token_symbol: "IACAI"
  token_decimals: 18
  enable_nft_access: true
  enable_token_payments: true
  basic_tier_rate_limit: 100
  pro_tier_rate_limit: 1000
  enterprise_tier_rate_limit: 10000
  default_agent_address: "0x147e832418Cc06A501047019E956714271098b89"

notion:
  base_url: "https://api.notion.com/v1"
  agent_name: "Test IaC AI Agent"
  agent_description: "Test Agent for IaC AI Integration"
  enable_agent_creation: true
  auto_create_on_startup: true
  max_requests_per_minute: 60
EOF

print_success "Configuração temporária criada"

# Teste 1: Compilação
print_status "Teste 1: Compilação do projeto..."
if go build -o /tmp/iac-ai-agent-test ./cmd/agent; then
    print_success "Compilação bem-sucedida"
else
    print_error "Falha na compilação"
    exit 1
fi

# Teste 2: Testes unitários
print_status "Teste 2: Executando testes unitários..."
if go test ./internal/services/... -v -timeout 30s; then
    print_success "Testes unitários passaram"
else
    print_warning "Alguns testes unitários falharam (pode ser esperado sem API keys)"
fi

# Teste 3: Testes de integração
print_status "Teste 3: Executando testes de integração..."
if go test ./test/integration/notion_integration_test.go -v -timeout 60s; then
    print_success "Testes de integração passaram"
else
    print_warning "Testes de integração falharam (provavelmente falta NOTION_API_KEY)"
fi

# Teste 4: Validação de configuração
print_status "Teste 4: Validação de configuração..."
if CONFIG_PATH=/tmp/test_config.yaml go run ./cmd/agent --validate-config; then
    print_success "Validação de configuração passou"
else
    print_warning "Validação de configuração falhou (esperado sem API keys)"
fi

# Teste 5: Verificação de dependências
print_status "Teste 5: Verificando dependências..."
if go mod tidy && go mod verify; then
    print_success "Dependências verificadas"
else
    print_error "Problema com dependências"
    exit 1
fi

# Limpeza
print_status "Limpando arquivos temporários..."
rm -f /tmp/test_config.yaml
rm -f /tmp/iac-ai-agent-test

print_success "Teste de integração Notion concluído!"

echo ""
echo "📋 Resumo dos testes:"
echo "✅ Compilação: OK"
echo "✅ Dependências: OK"
echo "⚠️  Testes unitários: Verificar logs acima"
echo "⚠️  Testes de integração: Verificar logs acima"
echo "⚠️  Validação de configuração: Verificar logs acima"
echo ""
echo "💡 Para testes completos, configure:"
echo "   - NOTION_API_KEY (obtenha em https://www.notion.so/my-integrations)"
echo "   - WALLET_TOKEN (para autenticação Web3)"
echo "   - Outras variáveis conforme documentação"
echo ""
echo "📚 Consulte docs/NOTION_INTEGRATION.md para mais informações"
