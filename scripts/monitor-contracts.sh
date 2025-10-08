#!/bin/bash

# ============================================
# 📊 Smart Contracts Monitoring Script
# ============================================
# 
# Script para monitorar contratos inteligentes em produção
# 
# Uso: ./monitor-contracts.sh [network] [options]
# 
# Exemplos:
#   ./monitor-contracts.sh base-mainnet --alerts
#   ./monitor-contracts.sh base-sepolia --health-check
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
MONITORING_DIR="$PROJECT_ROOT/monitoring"

# Parâmetros
NETWORK="${1:-base-mainnet}"
ENABLE_ALERTS="${2:-false}"
HEALTH_CHECK_ONLY="${3:-false}"

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

# Função para configurar ambiente
setup_environment() {
    log "Configurando ambiente de monitoramento para rede: $NETWORK"
    
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
    
    # Criar diretório de monitoramento
    mkdir -p "$MONITORING_DIR/logs"
    mkdir -p "$MONITORING_DIR/alerts"
    mkdir -p "$MONITORING_DIR/reports"
    
    log_success "Ambiente de monitoramento configurado"
}

# Função para carregar endereços dos contratos
load_contract_addresses() {
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    
    if [[ ! -f "$deployment_file" ]]; then
        log_error "Arquivo de deployment não encontrado: $deployment_file"
        log "Execute o deploy primeiro"
        exit 1
    fi
    
    # Extrair endereços do arquivo JSON
    TOKEN_ADDRESS=$(jq -r '.contracts.IACaiToken' "$deployment_file")
    NFT_ADDRESS=$(jq -r '.contracts.NationPassNFT' "$deployment_file")
    AGENT_ADDRESS=$(jq -r '.contracts.AgentContract' "$deployment_file")
    
    if [[ "$TOKEN_ADDRESS" == "null" || "$TOKEN_ADDRESS" == "" ]]; then
        log_error "Endereço do token não encontrado"
        exit 1
    fi
    
    log_success "Endereços dos contratos carregados"
    log "Token: $TOKEN_ADDRESS"
    log "NFT: $NFT_ADDRESS"
    log "Agent: $AGENT_ADDRESS"
}

# Função para verificar saúde dos contratos
check_contract_health() {
    log "🔍 Verificando saúde dos contratos..."
    
    local health_status="healthy"
    local issues=()
    
    # Verificar Token Contract
    if ! check_token_contract_health; then
        health_status="unhealthy"
        issues+=("Token contract issues detected")
    fi
    
    # Verificar NFT Contract
    if ! check_nft_contract_health; then
        health_status="unhealthy"
        issues+=("NFT contract issues detected")
    fi
    
    # Verificar Agent Contract
    if ! check_agent_contract_health; then
        health_status="unhealthy"
        issues+=("Agent contract issues detected")
    fi
    
    # Gerar relatório de saúde
    generate_health_report "$health_status" "${issues[@]}"
    
    if [[ "$health_status" == "healthy" ]]; then
        log_success "Todos os contratos estão saudáveis"
    else
        log_error "Problemas detectados nos contratos"
        for issue in "${issues[@]}"; do
            log_error "  - $issue"
        done
    fi
    
    return $([ "$health_status" == "healthy" ] && echo 0 || echo 1)
}

# Função para verificar saúde do contrato de token
check_token_contract_health() {
    log "Verificando saúde do IACaiToken..."
    
    local is_healthy=true
    
    # Verificar se o contrato responde
    if ! cast call "$TOKEN_ADDRESS" "name()" --rpc-url "$RPC_URL" > /dev/null 2>&1; then
        log_error "Token contract não responde"
        is_healthy=false
    fi
    
    # Verificar supply total
    local total_supply=$(cast call "$TOKEN_ADDRESS" "totalSupply()" --rpc-url "$RPC_URL")
    if [[ "$total_supply" == "0" ]]; then
        log_warning "Supply total é zero"
        is_healthy=false
    fi
    
    # Verificar se o contrato não está pausado
    local is_paused=$(cast call "$TOKEN_ADDRESS" "paused()" --rpc-url "$RPC_URL")
    if [[ "$is_paused" == "true" ]]; then
        log_warning "Token contract está pausado"
        is_healthy=false
    fi
    
    # Verificar saldo do contrato
    local contract_balance=$(cast call "$TOKEN_ADDRESS" "balanceOf(address)" "$TOKEN_ADDRESS" --rpc-url "$RPC_URL")
    if [[ "$contract_balance" == "0" ]]; then
        log_warning "Contrato não tem tokens disponíveis"
        is_healthy=false
    fi
    
    if [[ "$is_healthy" == true ]]; then
        log_success "Token contract está saudável"
    fi
    
    return $([ "$is_healthy" == true ] && echo 0 || echo 1)
}

# Função para verificar saúde do contrato de NFT
check_nft_contract_health() {
    log "Verificando saúde do NationPassNFT..."
    
    local is_healthy=true
    
    # Verificar se o contrato responde
    if ! cast call "$NFT_ADDRESS" "name()" --rpc-url "$RPC_URL" > /dev/null 2>&1; then
        log_error "NFT contract não responde"
        is_healthy=false
    fi
    
    # Verificar se o contrato não está pausado
    local is_paused=$(cast call "$NFT_ADDRESS" "paused()" --rpc-url "$RPC_URL")
    if [[ "$is_paused" == "true" ]]; then
        log_warning "NFT contract está pausado"
        is_healthy=false
    fi
    
    # Verificar tiers
    for i in {0..2}; do
        local tier_info=$(cast call "$NFT_ADDRESS" "tierInfo(uint8)" "$i" --rpc-url "$RPC_URL")
        if [[ -z "$tier_info" ]]; then
            log_warning "Tier $i não configurado"
            is_healthy=false
        fi
    done
    
    if [[ "$is_healthy" == true ]]; then
        log_success "NFT contract está saudável"
    fi
    
    return $([ "$is_healthy" == true ] && echo 0 || echo 1)
}

# Função para verificar saúde do contrato do agente
check_agent_contract_health() {
    log "Verificando saúde do AgentContract..."
    
    local is_healthy=true
    
    # Verificar se o contrato responde
    if ! cast call "$AGENT_ADDRESS" "config()" --rpc-url "$RPC_URL" > /dev/null 2>&1; then
        log_error "Agent contract não responde"
        is_healthy=false
    fi
    
    # Verificar se o contrato não está pausado
    local is_paused=$(cast call "$AGENT_ADDRESS" "paused()" --rpc-url "$RPC_URL")
    if [[ "$is_paused" == "true" ]]; then
        log_warning "Agent contract está pausado"
        is_healthy=false
    fi
    
    # Verificar contratos vinculados
    local token_contract=$(cast call "$AGENT_ADDRESS" "tokenContract()" --rpc-url "$RPC_URL")
    if [[ "$token_contract" != "$TOKEN_ADDRESS" ]]; then
        log_error "Token contract não vinculado corretamente"
        is_healthy=false
    fi
    
    local nft_contract=$(cast call "$AGENT_ADDRESS" "nftContract()" --rpc-url "$RPC_URL")
    if [[ "$nft_contract" != "$NFT_ADDRESS" ]]; then
        log_error "NFT contract não vinculado corretamente"
        is_healthy=false
    fi
    
    if [[ "$is_healthy" == true ]]; then
        log_success "Agent contract está saudável"
    fi
    
    return $([ "$is_healthy" == true ] && echo 0 || echo 1)
}

# Função para monitorar eventos
monitor_events() {
    log "📡 Monitorando eventos dos contratos..."
    
    # Monitorar eventos do Token Contract
    monitor_token_events &
    local token_pid=$!
    
    # Monitorar eventos do NFT Contract
    monitor_nft_events &
    local nft_pid=$!
    
    # Monitorar eventos do Agent Contract
    monitor_agent_events &
    local agent_pid=$!
    
    # Aguardar processos
    wait $token_pid
    wait $nft_pid
    wait $agent_pid
    
    log_success "Monitoramento de eventos concluído"
}

# Função para monitorar eventos do token
monitor_token_events() {
    log "Monitorando eventos do IACaiToken..."
    
    # Monitorar eventos de compra
    cast logs \
        --from-block latest \
        --to-block latest \
        --address "$TOKEN_ADDRESS" \
        --topic "0x..." \
        --rpc-url "$RPC_URL" > "$MONITORING_DIR/logs/token-events-$(date +%Y%m%d-%H%M%S).log" 2>/dev/null || true
}

# Função para monitorar eventos do NFT
monitor_nft_events() {
    log "Monitorando eventos do NationPassNFT..."
    
    # Monitorar eventos de mint
    cast logs \
        --from-block latest \
        --to-block latest \
        --address "$NFT_ADDRESS" \
        --topic "0x..." \
        --rpc-url "$RPC_URL" > "$MONITORING_DIR/logs/nft-events-$(date +%Y%m%d-%H%M%S).log" 2>/dev/null || true
}

# Função para monitorar eventos do agente
monitor_agent_events() {
    log "Monitorando eventos do AgentContract..."
    
    # Monitorar eventos de análise
    cast logs \
        --from-block latest \
        --to-block latest \
        --address "$AGENT_ADDRESS" \
        --topic "0x..." \
        --rpc-url "$RPC_URL" > "$MONITORING_DIR/logs/agent-events-$(date +%Y%m%d-%H%M%S).log" 2>/dev/null || true
}

# Função para gerar relatório de saúde
generate_health_report() {
    local health_status="$1"
    shift
    local issues=("$@")
    
    local report_file="$MONITORING_DIR/reports/health-report-$NETWORK-$(date +%Y%m%d-%H%M%S).md"
    
    cat > "$report_file" << EOF
# 📊 Relatório de Saúde dos Contratos

**Rede:** $NETWORK  
**Data:** $(date)  
**Status:** $health_status  

## 🔍 Verificações Realizadas

### ✅ Contratos Verificados
- IACaiToken: $TOKEN_ADDRESS
- NationPassNFT: $NFT_ADDRESS  
- AgentContract: $AGENT_ADDRESS

### 📋 Status dos Contratos
- Token Contract: $(check_token_contract_health && echo "✅ Saudável" || echo "❌ Problemas")
- NFT Contract: $(check_nft_contract_health && echo "✅ Saudável" || echo "❌ Problemas")
- Agent Contract: $(check_agent_contract_health && echo "✅ Saudável" || echo "❌ Problemas")

### ⚠️ Problemas Detectados
EOF

    if [[ ${#issues[@]} -eq 0 ]]; then
        echo "- Nenhum problema detectado" >> "$report_file"
    else
        for issue in "${issues[@]}"; do
            echo "- $issue" >> "$report_file"
        done
    fi
    
    cat >> "$report_file" << EOF

### 🔗 Links Úteis
- [Base Explorer]($EXPLORER_URL)
- [Token Contract]($EXPLORER_URL/address/$TOKEN_ADDRESS)
- [NFT Contract]($EXPLORER_URL/address/$NFT_ADDRESS)
- [Agent Contract]($EXPLORER_URL/address/$AGENT_ADDRESS)

### 📊 Próximas Verificações
- Monitoramento contínuo de eventos
- Verificação de segurança
- Backup de dados
- Atualização de configurações

---
*Relatório gerado automaticamente pelo script de monitoramento*
EOF

    log_success "Relatório de saúde gerado: $report_file"
}

# Função para configurar alertas
setup_alerts() {
    log "🚨 Configurando sistema de alertas..."
    
    # Criar arquivo de configuração de alertas
    local alerts_config="$MONITORING_DIR/alerts/alerts-config.json"
    
    cat > "$alerts_config" << EOF
{
  "network": "$NETWORK",
  "contracts": {
    "token": "$TOKEN_ADDRESS",
    "nft": "$NFT_ADDRESS",
    "agent": "$AGENT_ADDRESS"
  },
  "alerts": {
    "health_check_failure": {
      "enabled": true,
      "threshold": 1,
      "notification": "email"
    },
    "contract_paused": {
      "enabled": true,
      "threshold": 1,
      "notification": "email"
    },
    "unusual_activity": {
      "enabled": true,
      "threshold": 10,
      "notification": "slack"
    }
  },
  "notifications": {
    "email": {
      "enabled": true,
      "recipients": ["admin@example.com"]
    },
    "slack": {
      "enabled": true,
      "webhook_url": "\$SLACK_WEBHOOK_URL"
    }
  }
}
EOF

    log_success "Sistema de alertas configurado"
}

# Função principal
main() {
    log "🚀 Iniciando monitoramento de contratos inteligentes..."
    
    setup_environment
    load_contract_addresses
    
    if [[ "$HEALTH_CHECK_ONLY" == "true" ]]; then
        check_contract_health
    else
        check_contract_health
        
        if [[ "$ENABLE_ALERTS" == "true" ]]; then
            setup_alerts
        fi
        
        monitor_events
    fi
    
    log_success "Monitoramento concluído!"
}

# Executar função principal
main "$@"
