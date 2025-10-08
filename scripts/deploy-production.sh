#!/bin/bash

# ============================================
# 🚀 IaC AI Agent - Script de Deploy Produção
# ============================================
# 
# Script automatizado para deploy em produção
# 
# Uso: ./deploy-production.sh [environment]
# 
# Exemplos:
#   ./deploy-production.sh staging
#   ./deploy-production.sh production
#

set -euo pipefail

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurações
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENVIRONMENT="${1:-production}"
DOCKER_COMPOSE_FILE="configs/docker-compose.prod.yml"
ENV_FILE=".env.prod"

# Função para logging
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] ✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] ⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ❌ $1${NC}"
}

# Função para verificar pré-requisitos
check_prerequisites() {
    log "Verificando pré-requisitos..."
    
    # Verificar Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker não está instalado"
        exit 1
    fi
    
    # Verificar Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose não está instalado"
        exit 1
    fi
    
    # Verificar arquivo de ambiente
    if [[ ! -f "$PROJECT_ROOT/$ENV_FILE" ]]; then
        log_error "Arquivo de ambiente não encontrado: $ENV_FILE"
        log "Execute: cp env.prod.example $ENV_FILE"
        log "E configure as variáveis necessárias"
        exit 1
    fi
    
    # Verificar se está no diretório correto
    if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
        log_error "Não está no diretório raiz do projeto"
        exit 1
    fi
    
    log_success "Pré-requisitos verificados"
}

# Função para validar configuração
validate_config() {
    log "Validando configuração..."
    
    # Verificar variáveis obrigatórias
    source "$PROJECT_ROOT/$ENV_FILE"
    
    local missing_vars=()
    
    # Variáveis obrigatórias para produção
    if [[ -z "${LLM_PROVIDER:-}" ]]; then
        missing_vars+=("LLM_PROVIDER")
    fi
    
    if [[ -z "${LLM_MODEL:-}" ]]; then
        missing_vars+=("LLM_MODEL")
    fi
    
    if [[ -z "${WALLET_ADDRESS:-}" ]]; then
        missing_vars+=("WALLET_ADDRESS")
    fi
    
    if [[ -z "${PRIVY_APP_ID:-}" ]]; then
        missing_vars+=("PRIVY_APP_ID")
    fi
    
    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        log_error "Variáveis obrigatórias não configuradas:"
        for var in "${missing_vars[@]}"; do
            log_error "  - $var"
        done
        exit 1
    fi
    
    log_success "Configuração validada"
}

# Função para construir imagem
build_image() {
    log "Construindo imagem Docker..."
    
    cd "$PROJECT_ROOT"
    
    # Construir imagem de produção
    docker build -t iacai-agent:prod -f deployments/Dockerfile.prod .
    
    if [[ $? -eq 0 ]]; then
        log_success "Imagem construída com sucesso"
    else
        log_error "Falha ao construir imagem"
        exit 1
    fi
}

# Função para parar serviços existentes
stop_existing_services() {
    log "Parando serviços existentes..."
    
    cd "$PROJECT_ROOT"
    
    # Parar containers existentes
    docker-compose -f "$DOCKER_COMPOSE_FILE" down --remove-orphans || true
    
    log_success "Serviços parados"
}

# Função para iniciar serviços
start_services() {
    log "Iniciando serviços..."
    
    cd "$PROJECT_ROOT"
    
    # Iniciar serviços
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d
    
    if [[ $? -eq 0 ]]; then
        log_success "Serviços iniciados"
    else
        log_error "Falha ao iniciar serviços"
        exit 1
    fi
}

# Função para aguardar serviços ficarem saudáveis
wait_for_services() {
    log "Aguardando serviços ficarem saudáveis..."
    
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        log "Tentativa $attempt/$max_attempts..."
        
        # Verificar se o serviço principal está saudável
        if docker-compose -f "$DOCKER_COMPOSE_FILE" ps iac-ai-agent | grep -q "healthy"; then
            log_success "Serviços estão saudáveis"
            return 0
        fi
        
        sleep 10
        ((attempt++))
    done
    
    log_error "Timeout aguardando serviços ficarem saudáveis"
    return 1
}

# Função para verificar saúde da aplicação
health_check() {
    log "Verificando saúde da aplicação..."
    
    # Aguardar um pouco para garantir que está rodando
    sleep 5
    
    # Tentar acessar o health check
    local health_url="http://localhost:8080/health"
    
    if curl -f -s "$health_url" > /dev/null; then
        log_success "Aplicação está saudável"
        return 0
    else
        log_error "Aplicação não está respondendo no health check"
        return 1
    fi
}

# Função para mostrar status
show_status() {
    log "Status dos serviços:"
    
    cd "$PROJECT_ROOT"
    docker-compose -f "$DOCKER_COMPOSE_FILE" ps
    
    echo ""
    log "Logs da aplicação (últimas 20 linhas):"
    docker-compose -f "$DOCKER_COMPOSE_FILE" logs --tail=20 iac-ai-agent
    
    echo ""
    log "URLs importantes:"
    echo "  - Aplicação: http://localhost:8080"
    echo "  - Health Check: http://localhost:8080/health"
    echo "  - Prometheus: http://localhost:9090 (se habilitado)"
    echo "  - Grafana: http://localhost:3000 (se habilitado)"
}

# Função para limpeza em caso de erro
cleanup_on_error() {
    log_error "Erro durante o deploy. Fazendo limpeza..."
    
    cd "$PROJECT_ROOT"
    docker-compose -f "$DOCKER_COMPOSE_FILE" down --remove-orphans || true
    
    log "Limpeza concluída"
}

# Função principal
main() {
    log "🚀 Iniciando deploy de produção do IaC AI Agent"
    log "Ambiente: $ENVIRONMENT"
    
    # Configurar trap para limpeza em caso de erro
    trap cleanup_on_error ERR
    
    # Executar etapas do deploy
    check_prerequisites
    validate_config
    build_image
    stop_existing_services
    start_services
    
    if wait_for_services; then
        if health_check; then
            log_success "🎉 Deploy concluído com sucesso!"
            show_status
        else
            log_error "Deploy falhou no health check"
            exit 1
        fi
    else
        log_error "Deploy falhou - serviços não ficaram saudáveis"
        exit 1
    fi
}

# Executar função principal
main "$@"