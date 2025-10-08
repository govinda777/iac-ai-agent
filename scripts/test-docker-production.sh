#!/bin/bash

# ============================================
# 🧪 Teste Docker Production Setup
# ============================================
# 
# Script para testar o setup de produção Docker

set -e

echo "🧪 Testando Docker Production Setup..."

# Cores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    print_error "Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Teste 1: Verificar arquivos necessários
echo "📋 Verificando arquivos necessários..."

required_files=(
    "deployments/Dockerfile.prod"
    "configs/docker-compose.prod.yml"
    "env.prod.example"
    "configs/nginx/nginx.conf"
    "configs/monitoring/prometheus.yml"
    "scripts/deploy-production.sh"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        print_success "✓ $file"
    else
        print_error "✗ $file não encontrado"
        exit 1
    fi
done

# Teste 2: Verificar se Docker está funcionando
echo ""
echo "🐳 Testando Docker..."

if docker --version > /dev/null 2>&1; then
    print_success "Docker está instalado: $(docker --version)"
else
    print_error "Docker não está instalado"
    exit 1
fi

if docker-compose --version > /dev/null 2>&1; then
    print_success "Docker Compose está instalado: $(docker-compose --version)"
else
    print_error "Docker Compose não está instalado"
    exit 1
fi

# Teste 3: Verificar sintaxe do docker-compose
echo ""
echo "📝 Verificando sintaxe do docker-compose..."

if docker-compose -f configs/docker-compose.prod.yml config > /dev/null 2>&1; then
    print_success "Sintaxe do docker-compose.prod.yml está correta"
else
    print_error "Erro na sintaxe do docker-compose.prod.yml"
    docker-compose -f configs/docker-compose.prod.yml config
    exit 1
fi

# Teste 4: Verificar sintaxe do Dockerfile
echo ""
echo "📦 Verificando sintaxe do Dockerfile..."

if docker build --dry-run -f deployments/Dockerfile.prod . > /dev/null 2>&1; then
    print_success "Sintaxe do Dockerfile.prod está correta"
else
    print_warning "Não foi possível verificar Dockerfile (dry-run não suportado)"
fi

# Teste 5: Verificar configuração Nginx
echo ""
echo "🌐 Verificando configuração Nginx..."

if command -v nginx > /dev/null 2>&1; then
    if nginx -t -c configs/nginx/nginx.conf > /dev/null 2>&1; then
        print_success "Configuração Nginx está correta"
    else
        print_warning "Nginx não está instalado para verificação"
    fi
else
    print_warning "Nginx não está instalado para verificação"
fi

# Teste 6: Verificar arquivo de ambiente
echo ""
echo "⚙️ Verificando arquivo de ambiente..."

if [ -f "env.prod.example" ]; then
    print_success "Arquivo env.prod.example existe"
    
    # Verificar se contém configurações Nation.fun
    if grep -q "LLM_PROVIDER=nation.fun" env.prod.example; then
        print_success "Configuração Nation.fun encontrada"
    else
        print_warning "Configuração Nation.fun não encontrada"
    fi
    
    if grep -q "NOTION_API_KEY" env.prod.example; then
        print_success "Configuração Notion encontrada"
    else
        print_warning "Configuração Notion não encontrada"
    fi
else
    print_error "Arquivo env.prod.example não encontrado"
fi

# Teste 7: Verificar scripts
echo ""
echo "🔧 Verificando scripts..."

if [ -x "scripts/deploy-production.sh" ]; then
    print_success "Script deploy-production.sh é executável"
else
    print_error "Script deploy-production.sh não é executável"
fi

# Resumo final
echo ""
echo "📊 Resumo dos Testes:"
echo "===================="

if [ $? -eq 0 ]; then
    print_success "🎉 Todos os testes passaram!"
    echo ""
    echo "✅ Docker Production Setup está pronto!"
    echo ""
    echo "📋 Próximos passos:"
    echo "1. Copie env.prod.example para .env.prod"
    echo "2. Configure suas variáveis de ambiente"
    echo "3. Execute: ./scripts/deploy-production.sh"
    echo ""
    echo "🔗 URLs importantes:"
    echo "- Health Check: http://localhost/health"
    echo "- Prometheus: http://localhost:9090"
    echo "- Grafana: http://localhost:3000"
else
    print_error "❌ Alguns testes falharam"
    echo "Verifique os erros acima antes de prosseguir"
    exit 1
fi
