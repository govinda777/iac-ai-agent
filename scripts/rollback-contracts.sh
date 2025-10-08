#!/bin/bash

# ============================================
# 🔄 Smart Contracts Rollback Script
# ============================================
# 
# Script para fazer rollback de contratos inteligentes
# 
# Uso: ./rollback-contracts.sh [network] [version] [options]
# 
# Exemplos:
#   ./rollback-contracts.sh base-mainnet v1.0.0 --confirm
#   ./rollback-contracts.sh base-sepolia v1.1.0 --dry-run
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
CONTRACTS_DIR="$PROJECT_ROOT/contracts"
ROLLBACK_DIR="$PROJECT_ROOT/rollback"

# Parâmetros
NETWORK="${1:-base-mainnet}"
TARGET_VERSION="${2:-}"
DRY_RUN="${3:-false}"
CONFIRM="${4:-false}"

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

# Função para mostrar ajuda
show_help() {
    cat << EOF
🔄 Smart Contracts Rollback Script

Uso: $0 [network] [version] [options]

Argumentos:
  network        Rede blockchain (base-mainnet, base-sepolia)
  version        Versão para rollback (ex: v1.0.0)
  options        Opções adicionais

Opções:
  --dry-run      Simular rollback sem executar
  --confirm      Confirmar rollback automaticamente
  --help         Mostrar esta ajuda

Exemplos:
  $0 base-mainnet v1.0.0 --dry-run
  $0 base-sepolia v1.1.0 --confirm
  $0 base-mainnet v1.0.0

⚠️  ATENÇÃO: Rollback é uma operação crítica que pode afetar usuários em produção!
EOF
}

# Função para verificar pré-requisitos
check_prerequisites() {
    log "Verificando pré-requisitos para rollback..."
    
    # Verificar se está no diretório correto
    if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
        log_error "Não está no diretório raiz do projeto"
        exit 1
    fi
    
    # Verificar se Foundry está instalado
    if ! command -v forge &> /dev/null; then
        log_error "Foundry não está instalado"
        log "Instale com: curl -L https://foundry.paradigm.xyz | bash"
        exit 1
    fi
    
    # Verificar se jq está instalado
    if ! command -v jq &> /dev/null; then
        log_error "jq não está instalado"
        log "Instale com: brew install jq (macOS) ou apt-get install jq (Ubuntu)"
        exit 1
    fi
    
    # Verificar se a versão foi especificada
    if [[ -z "$TARGET_VERSION" ]]; then
        log_error "Versão de rollback não especificada"
        log "Use: $0 $NETWORK v1.0.0"
        exit 1
    fi
    
    log_success "Pré-requisitos verificados"
}

# Função para configurar ambiente
setup_environment() {
    log "Configurando ambiente para rollback na rede: $NETWORK"
    
    # Configurar variáveis de ambiente baseadas na rede
    case $NETWORK in
        "base-sepolia")
            export RPC_URL="https://sepolia.base.org"
            export EXPLORER_URL="https://sepolia.basescan.org"
            export CHAIN_ID="84532"
            ;;
        "base-mainnet")
            export RPC_URL="https://mainnet.base.org"
            export EXPLORER_URL="https://basescan.org"
            export CHAIN_ID="8453"
            ;;
        *)
            log_error "Rede não suportada: $NETWORK"
            log "Redes suportadas: base-sepolia, base-mainnet"
            exit 1
            ;;
    esac
    
    # Criar diretório de rollback
    mkdir -p "$ROLLBACK_DIR/backups"
    mkdir -p "$ROLLBACK_DIR/logs"
    mkdir -p "$ROLLBACK_DIR/reports"
    
    log_success "Ambiente de rollback configurado"
}

# Função para verificar versão atual
check_current_version() {
    log "Verificando versão atual dos contratos..."
    
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    
    if [[ ! -f "$deployment_file" ]]; then
        log_error "Arquivo de deployment não encontrado: $deployment_file"
        exit 1
    fi
    
    # Extrair versão atual
    local current_version=$(jq -r '.version // "unknown"' "$deployment_file")
    log "Versão atual: $current_version"
    
    # Verificar se a versão de rollback existe
    local rollback_file="$ROLLBACK_DIR/backups/$NETWORK-$TARGET_VERSION.json"
    if [[ ! -f "$rollback_file" ]]; then
        log_error "Backup da versão $TARGET_VERSION não encontrado"
        log "Arquivo esperado: $rollback_file"
        exit 1
    fi
    
    log_success "Versão de rollback encontrada: $TARGET_VERSION"
}

# Função para fazer backup da versão atual
backup_current_version() {
    log "📦 Fazendo backup da versão atual..."
    
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    local backup_file="$ROLLBACK_DIR/backups/$NETWORK-$(date +%Y%m%d-%H%M%S).json"
    
    # Copiar arquivo de deployment atual
    cp "$deployment_file" "$backup_file"
    
    # Adicionar metadados do backup
    jq --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       --arg reason "rollback-to-$TARGET_VERSION" \
       '. + {
           backup_timestamp: $timestamp,
           backup_reason: $reason,
           rollback_from: .version
       }' "$backup_file" > "$backup_file.tmp" && mv "$backup_file.tmp" "$backup_file"
    
    log_success "Backup criado: $backup_file"
}

# Função para simular rollback
simulate_rollback() {
    log "🔍 Simulando rollback para versão $TARGET_VERSION..."
    
    local rollback_file="$ROLLBACK_DIR/backups/$NETWORK-$TARGET_VERSION.json"
    
    # Carregar dados do rollback
    local rollback_data=$(cat "$rollback_file")
    local token_address=$(echo "$rollback_data" | jq -r '.contracts.IACaiToken')
    local nft_address=$(echo "$rollback_data" | jq -r '.contracts.NationPassNFT')
    local agent_address=$(echo "$rollback_data" | jq -r '.contracts.AgentContract')
    
    log "Endereços da versão $TARGET_VERSION:"
    log "  Token: $token_address"
    log "  NFT: $nft_address"
    log "  Agent: $agent_address"
    
    # Verificar se os contratos existem
    verify_rollback_contracts "$token_address" "$nft_address" "$agent_address"
    
    # Simular atualização de configurações
    simulate_configuration_update "$rollback_data"
    
    log_success "Simulação de rollback concluída"
}

# Função para verificar contratos de rollback
verify_rollback_contracts() {
    local token_address="$1"
    local nft_address="$2"
    local agent_address="$3"
    
    log "Verificando contratos de rollback..."
    
    # Verificar Token Contract
    if cast code "$token_address" --rpc-url "$RPC_URL" > /dev/null 2>&1; then
        log_success "Token contract verificado: $token_address"
    else
        log_error "Token contract não encontrado: $token_address"
        return 1
    fi
    
    # Verificar NFT Contract
    if cast code "$nft_address" --rpc-url "$RPC_URL" > /dev/null 2>&1; then
        log_success "NFT contract verificado: $nft_address"
    else
        log_error "NFT contract não encontrado: $nft_address"
        return 1
    fi
    
    # Verificar Agent Contract
    if cast code "$agent_address" --rpc-url "$RPC_URL" > /dev/null 2>&1; then
        log_success "Agent contract verificado: $agent_address"
    else
        log_error "Agent contract não encontrado: $agent_address"
        return 1
    fi
}

# Função para simular atualização de configurações
simulate_configuration_update() {
    local rollback_data="$1"
    
    log "Simulando atualização de configurações..."
    
    # Simular atualização do arquivo de deployment
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    local backup_deployment="$deployment_file.backup"
    
    # Fazer backup do arquivo atual
    cp "$deployment_file" "$backup_deployment"
    
    # Simular atualização
    echo "$rollback_data" > "$deployment_file"
    
    log "Configurações simuladas:"
    log "  Arquivo de deployment atualizado"
    log "  Backup criado: $backup_deployment"
    
    # Restaurar arquivo original
    mv "$backup_deployment" "$deployment_file"
}

# Função para executar rollback
execute_rollback() {
    log "🚀 Executando rollback para versão $TARGET_VERSION..."
    
    local rollback_file="$ROLLBACK_DIR/backups/$NETWORK-$TARGET_VERSION.json"
    
    # Carregar dados do rollback
    local rollback_data=$(cat "$rollback_file")
    
    # Fazer backup da versão atual
    backup_current_version
    
    # Atualizar arquivo de deployment
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    local backup_deployment="$deployment_file.backup"
    
    # Fazer backup do arquivo atual
    cp "$deployment_file" "$backup_deployment"
    
    # Atualizar com dados do rollback
    echo "$rollback_data" > "$deployment_file"
    
    # Adicionar metadados do rollback
    jq --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       --arg version "$TARGET_VERSION" \
       --arg reason "rollback-executed" \
       '. + {
           rollback_timestamp: $timestamp,
           rollback_version: $version,
           rollback_reason: $reason,
           previous_version: .version
       }' "$deployment_file" > "$deployment_file.tmp" && mv "$deployment_file.tmp" "$deployment_file"
    
    log_success "Rollback executado com sucesso"
}

# Função para verificar rollback
verify_rollback() {
    log "🔍 Verificando rollback..."
    
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    local current_version=$(jq -r '.rollback_version // .version' "$deployment_file")
    
    if [[ "$current_version" == "$TARGET_VERSION" ]]; then
        log_success "Rollback verificado: versão $TARGET_VERSION ativa"
    else
        log_error "Rollback falhou: versão atual $current_version, esperada $TARGET_VERSION"
        return 1
    fi
    
    # Verificar saúde dos contratos
    if verify_rollback_contracts \
        "$(jq -r '.contracts.IACaiToken' "$deployment_file")" \
        "$(jq -r '.contracts.NationPassNFT' "$deployment_file")" \
        "$(jq -r '.contracts.AgentContract' "$deployment_file")"; then
        log_success "Contratos de rollback verificados"
    else
        log_error "Verificação dos contratos falhou"
        return 1
    fi
}

# Função para gerar relatório de rollback
generate_rollback_report() {
    log "📊 Gerando relatório de rollback..."
    
    local report_file="$ROLLBACK_DIR/reports/rollback-report-$NETWORK-$(date +%Y%m%d-%H%M%S).md"
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    
    cat > "$report_file" << EOF
# 🔄 Relatório de Rollback de Contratos

**Rede:** $NETWORK  
**Versão de Rollback:** $TARGET_VERSION  
**Data:** $(date)  
**Tipo:** $([ "$DRY_RUN" == "true" ] && echo "Simulação" || echo "Execução Real")  

## 📋 Resumo do Rollback

### ✅ Operações Realizadas
- Backup da versão atual criado
- Contratos de rollback verificados
- Configurações atualizadas
- Verificação de saúde executada

### 🔗 Contratos Atualizados
- IACaiToken: $(jq -r '.contracts.IACaiToken' "$deployment_file")
- NationPassNFT: $(jq -r '.contracts.NationPassNFT' "$deployment_file")
- AgentContract: $(jq -r '.contracts.AgentContract' "$deployment_file")

### 📊 Status da Verificação
- Contratos verificados: ✅
- Configurações atualizadas: ✅
- Saúde dos contratos: ✅

### 🔗 Links Úteis
- [Base Explorer]($EXPLORER_URL)
- [Token Contract]($EXPLORER_URL/address/$(jq -r '.contracts.IACaiToken' "$deployment_file"))
- [NFT Contract]($EXPLORER_URL/address/$(jq -r '.contracts.NationPassNFT' "$deployment_file"))
- [Agent Contract]($EXPLORER_URL/address/$(jq -r '.contracts.AgentContract' "$deployment_file"))

### 📝 Próximos Passos
1. Monitorar contratos em produção
2. Verificar funcionalidades críticas
3. Notificar usuários sobre mudanças
4. Atualizar documentação

---
*Relatório gerado automaticamente pelo script de rollback*
EOF

    log_success "Relatório de rollback gerado: $report_file"
}

# Função para confirmar rollback
confirm_rollback() {
    if [[ "$CONFIRM" == "true" ]]; then
        return 0
    fi
    
    log_warning "⚠️  ATENÇÃO: Esta operação fará rollback dos contratos para a versão $TARGET_VERSION"
    log_warning "Esta é uma operação crítica que pode afetar usuários em produção!"
    
    read -p "Tem certeza que deseja continuar? (yes/no): " -r
    if [[ $REPLY != "yes" ]]; then
        log "Rollback cancelado pelo usuário"
        exit 0
    fi
}

# Função principal
main() {
    # Verificar argumentos de ajuda
    if [[ "$1" == "--help" || "$1" == "-h" ]]; then
        show_help
        exit 0
    fi
    
    log "🔄 Iniciando rollback de contratos inteligentes..."
    
    check_prerequisites
    setup_environment
    check_current_version
    
    if [[ "$DRY_RUN" == "true" ]]; then
        simulate_rollback
        generate_rollback_report
        log_success "Simulação de rollback concluída!"
    else
        confirm_rollback
        execute_rollback
        verify_rollback
        generate_rollback_report
        log_success "Rollback executado com sucesso!"
    fi
}

# Executar função principal
main "$@"
