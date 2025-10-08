#!/bin/bash

# Script de validação de pré-requisitos para IaC AI Agent
# Verifica se todos os requisitos estão instalados e configurados

set -e

echo "🔍 Validando pré-requisitos do IaC AI Agent..."
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Contador de erros
ERRORS=0

echo "📋 Verificando dependências do sistema..."

# 1. Verificar Go
echo -n "Verificando Go... "
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | cut -d' ' -f3)
    if [[ $GO_VERSION =~ go1\.(2[1-9]|[3-9][0-9]) ]]; then
        print_status 0 "Go $GO_VERSION instalado"
    else
        print_status 1 "Go $GO_VERSION instalado (requer 1.21+)"
        ERRORS=$((ERRORS + 1))
    fi
else
    print_status 1 "Go não instalado"
    ERRORS=$((ERRORS + 1))
fi

# 2. Verificar Git
echo -n "Verificando Git... "
if command -v git &> /dev/null; then
    print_status 0 "Git $(git --version | cut -d' ' -f3) instalado"
else
    print_status 1 "Git não instalado"
    ERRORS=$((ERRORS + 1))
fi

# 3. Verificar Make
echo -n "Verificando Make... "
if command -v make &> /dev/null; then
    print_status 0 "Make instalado"
else
    print_warning "Make não instalado (opcional, mas recomendado)"
fi

# 4. Verificar Python3 (para Checkov)
echo -n "Verificando Python3... "
if command -v python3 &> /dev/null; then
    print_status 0 "Python3 $(python3 --version | cut -d' ' -f2) instalado"
else
    print_warning "Python3 não instalado (Checkov não estará disponível)"
fi

# 5. Verificar Terraform
echo -n "Verificando Terraform... "
if command -v terraform &> /dev/null; then
    print_status 0 "Terraform $(terraform version | head -n1 | cut -d' ' -f2) instalado"
else
    print_warning "Terraform não instalado (algumas funcionalidades podem não funcionar)"
fi

echo ""
echo "📋 Verificando arquivos de configuração..."

# 6. Verificar arquivo .env
echo -n "Verificando arquivo .env... "
if [ -f .env ]; then
    print_status 0 "Arquivo .env encontrado"
    
    # Verificar variáveis obrigatórias para NFT Pass do Nation
    echo "  Verificando variáveis obrigatórias..."
    
    # Verificar WALLET_ADDRESS
    if grep -q "^WALLET_ADDRESS=" .env && ! grep -q "^WALLET_ADDRESS=$" .env; then
        print_status 0 "WALLET_ADDRESS configurado"
    else
        print_status 1 "WALLET_ADDRESS não configurado"
        ERRORS=$((ERRORS + 1))
    fi
    
    # Verificar NATION_NFT_REQUIRED
    if grep -q "^NATION_NFT_REQUIRED=true" .env; then
        print_status 0 "NATION_NFT_REQUIRED configurado"
    else
        print_status 1 "NATION_NFT_REQUIRED não configurado"
        ERRORS=$((ERRORS + 1))
    fi
    
    # Verificar LLM_PROVIDER
    if grep -q "^LLM_PROVIDER=nation.fun" .env; then
        print_status 0 "LLM_PROVIDER configurado para nation.fun"
    else
        print_status 1 "LLM_PROVIDER não configurado para nation.fun"
        ERRORS=$((ERRORS + 1))
    fi
    
    # Verificar configurações hardcoded
    echo "  Verificando configurações hardcoded..."
    
    # Verificar se Privy App ID está hardcoded
    if grep -q "cmgh6un8w007bl10ci0tgitwp" configs/app.yaml; then
        print_status 0 "Privy App ID hardcoded"
    else
        print_status 1 "Privy App ID não encontrado no app.yaml"
        ERRORS=$((ERRORS + 1))
    fi
    
    # Verificar se Wallet Address está hardcoded
    if grep -q "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5" configs/app.yaml; then
        print_status 0 "Wallet Address hardcoded"
    else
        print_status 1 "Wallet Address não encontrado no app.yaml"
        ERRORS=$((ERRORS + 1))
    fi
    
else
    print_status 1 "Arquivo .env não encontrado"
    print_info "Execute: cp env.example .env"
    ERRORS=$((ERRORS + 1))
fi

# 7. Verificar arquivo app.yaml
echo -n "Verificando arquivo app.yaml... "
if [ -f configs/app.yaml ]; then
    print_status 0 "Arquivo configs/app.yaml encontrado"
else
    print_status 1 "Arquivo configs/app.yaml não encontrado"
    print_info "Execute: cp configs/app.yaml.example configs/app.yaml"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "📋 Verificando dependências do Go..."

# 8. Verificar go.mod
echo -n "Verificando go.mod... "
if [ -f go.mod ]; then
    print_status 0 "go.mod encontrado"
else
    print_status 1 "go.mod não encontrado"
    ERRORS=$((ERRORS + 1))
fi

# 9. Verificar se dependências estão baixadas
echo -n "Verificando dependências Go... "
if [ -d vendor ] || go mod download &> /dev/null; then
    print_status 0 "Dependências Go OK"
else
    print_status 1 "Erro ao baixar dependências Go"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "📋 Verificando conectividade..."

# 10. Verificar conectividade com Nation.fun (apenas se validação estiver habilitada)
echo -n "Verificando conectividade Nation.fun... "
if [ -f .env ] && grep -q "^ENABLE_STARTUP_VALIDATION=false" .env; then
    print_warning "Validação de startup desabilitada - pulando verificação Nation.fun"
elif [ -f .env ] && grep -q "^WALLET_ADDRESS=" .env; then
    WALLET_ADDRESS=$(grep "^WALLET_ADDRESS=" .env | cut -d'=' -f2)
    if curl -s "https://api.nation.fun/v1/nft/check/$WALLET_ADDRESS" > /dev/null 2>&1; then
        print_status 0 "Conectividade Nation.fun OK"
    else
        print_status 1 "Erro de conectividade Nation.fun"
        ERRORS=$((ERRORS + 1))
    fi
else
    print_warning "Não é possível verificar Nation.fun (WALLET_ADDRESS não configurado)"
fi

# 11. Verificar Git Secrets
echo -n "Verificando Git Secrets... "
if command -v git-secret &> /dev/null; then
    if git secret list &> /dev/null; then
        print_status 0 "Git Secrets configurado"
    else
        print_warning "Git Secrets não inicializado"
    fi
else
    print_warning "Git Secrets não instalado (opcional)"
fi

echo ""
echo "📊 RELATÓRIO DE VALIDAÇÃO"
echo "=========================="

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ Todos os pré-requisitos estão OK!${NC}"
    echo ""
    echo "🚀 Próximos passos:"
    echo "   1. Execute: make run"
    echo "   2. Teste: curl http://localhost:8080/health"
    echo "   3. Acesse: http://localhost:8080"
else
    echo -e "${RED}❌ $ERRORS erro(s) encontrado(s)${NC}"
    echo ""
    echo "🔧 Para corrigir:"
    echo "   1. Instale as dependências faltantes"
    echo "   2. Configure WALLET_ADDRESS e NATION_NFT_REQUIRED no arquivo .env"
    echo "   3. Execute: cp configs/app.yaml.example configs/app.yaml"
    echo "   4. Execute: make setup"
    echo ""
    echo "📚 Consulte a documentação:"
    echo "   - docs/GUIA_INSTALACAO.md"
    echo "   - docs/QUICKSTART_ATUALIZADO.md"
fi

echo ""
echo "📚 Documentação disponível:"
echo "   - README.md - Visão geral e quick start"
echo "   - docs/GUIA_INSTALACAO.md - Instalação completa"
echo "   - docs/EXEMPLOS_PRATICOS.md - Exemplos de uso"
echo "   - docs/CONFIGURACAO_VARIAVEIS.md - Configuração detalhada"

exit $ERRORS
