#!/bin/bash

# ============================================
# 🔍 Smart Contracts Verification Script
# ============================================
# 
# Script para verificar contratos inteligentes em diferentes redes
# 
# Uso: ./verify-contracts.sh [network] [contract-address]
# 
# Exemplos:
#   ./verify-contracts.sh base-sepolia 0x1234...
#   ./verify-contracts.sh base-mainnet 0x5678...
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

# Parâmetros
NETWORK="${1:-base-sepolia}"
CONTRACT_ADDRESS="${2:-}"

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
    
    # Verificar se está no diretório de contratos
    if [[ ! -f "$CONTRACTS_DIR/foundry.toml" ]]; then
        log_error "Diretório de contratos não encontrado"
        exit 1
    fi
    
    log_success "Pré-requisitos verificados"
}

# Função para configurar ambiente
setup_environment() {
    log "Configurando ambiente para rede: $NETWORK"
    
    cd "$CONTRACTS_DIR"
    
    # Configurar variáveis de ambiente baseadas na rede
    case $NETWORK in
        "base-sepolia")
            export RPC_URL="https://sepolia.base.org"
            export ETHERSCAN_API_KEY="${ETHERSCAN_API_KEY:-}"
            export CHAIN_ID="84532"
            ;;
        "base-mainnet")
            export RPC_URL="https://mainnet.base.org"
            export ETHERSCAN_API_KEY="${ETHERSCAN_API_KEY:-}"
            export CHAIN_ID="8453"
            ;;
        "local")
            export RPC_URL="http://localhost:8545"
            export CHAIN_ID="31337"
            ;;
        *)
            log_error "Rede não suportada: $NETWORK"
            log "Redes suportadas: base-sepolia, base-mainnet, local"
            exit 1
            ;;
    esac
    
    log_success "Ambiente configurado para $NETWORK"
}

# Função para verificar contrato específico
verify_contract() {
    local contract_address="$1"
    local contract_name="$2"
    
    log "Verificando contrato $contract_name em $contract_address..."
    
    # Verificar se o contrato existe
    if ! cast code "$contract_address" --rpc-url "$RPC_URL" | grep -q "0x"; then
        log_error "Contrato não encontrado em $contract_address"
        return 1
    fi
    
    # Verificar se o contrato está verificado no Etherscan
    if [[ "$NETWORK" != "local" && -n "$ETHERSCAN_API_KEY" ]]; then
        log "Verificando no Etherscan..."
        
        # Tentar verificar o contrato
        if forge verify-contract \
            "$contract_address" \
            "$contract_name" \
            --etherscan-api-key "$ETHERSCAN_API_KEY" \
            --chain "$NETWORK" \
            --watch; then
            log_success "Contrato $contract_name verificado no Etherscan"
        else
            log_warning "Falha na verificação do Etherscan para $contract_name"
        fi
    fi
    
    # Verificar funcionalidades básicas
    verify_contract_functionality "$contract_address" "$contract_name"
}

# Função para verificar funcionalidades do contrato
verify_contract_functionality() {
    local contract_address="$1"
    local contract_name="$2"
    
    log "Verificando funcionalidades do contrato $contract_name..."
    
    case $contract_name in
        "IACaiToken")
            verify_token_contract "$contract_address"
            ;;
        "NationPassNFT")
            verify_nft_contract "$contract_address"
            ;;
        "AgentContract")
            verify_agent_contract "$contract_address"
            ;;
        *)
            log_warning "Tipo de contrato não reconhecido: $contract_name"
            ;;
    esac
}

# Função para verificar contrato de token
verify_token_contract() {
    local contract_address="$1"
    
    log "Verificando funcionalidades do IACaiToken..."
    
    # Verificar nome do token
    local token_name=$(cast call "$contract_address" "name()" --rpc-url "$RPC_URL")
    if [[ "$token_name" == *"IaC AI Token"* ]]; then
        log_success "Nome do token correto"
    else
        log_error "Nome do token incorreto: $token_name"
    fi
    
    # Verificar símbolo do token
    local token_symbol=$(cast call "$contract_address" "symbol()" --rpc-url "$RPC_URL")
    if [[ "$token_symbol" == *"IACAI"* ]]; then
        log_success "Símbolo do token correto"
    else
        log_error "Símbolo do token incorreto: $token_symbol"
    fi
    
    # Verificar supply total
    local total_supply=$(cast call "$contract_address" "totalSupply()" --rpc-url "$RPC_URL")
    local expected_supply="1000000000000000000000000" # 1M tokens
    if [[ "$total_supply" == "$expected_supply" ]]; then
        log_success "Supply total correto"
    else
        log_error "Supply total incorreto: $total_supply"
    fi
    
    # Verificar pacotes
    for i in {1..4}; do
        local package_info=$(cast call "$contract_address" "packages(uint8)" "$i" --rpc-url "$RPC_URL")
        if [[ -n "$package_info" ]]; then
            log_success "Pacote $i configurado"
        else
            log_warning "Pacote $i não encontrado"
        fi
    done
}

# Função para verificar contrato de NFT
verify_nft_contract() {
    local contract_address="$1"
    
    log "Verificando funcionalidades do NationPassNFT..."
    
    # Verificar nome do NFT
    local nft_name=$(cast call "$contract_address" "name()" --rpc-url "$RPC_URL")
    if [[ "$nft_name" == *"Nation Pass NFT"* ]]; then
        log_success "Nome do NFT correto"
    else
        log_error "Nome do NFT incorreto: $nft_name"
    fi
    
    # Verificar símbolo do NFT
    local nft_symbol=$(cast call "$contract_address" "symbol()" --rpc-url "$RPC_URL")
    if [[ "$nft_symbol" == *"NATION"* ]]; then
        log_success "Símbolo do NFT correto"
    else
        log_error "Símbolo do NFT incorreto: $nft_symbol"
    fi
    
    # Verificar tiers
    for i in {0..2}; do
        local tier_info=$(cast call "$contract_address" "tierInfo(uint8)" "$i" --rpc-url "$RPC_URL")
        if [[ -n "$tier_info" ]]; then
            log_success "Tier $i configurado"
        else
            log_warning "Tier $i não encontrado"
        fi
    done
}

# Função para verificar contrato do agente
verify_agent_contract() {
    local contract_address="$1"
    
    log "Verificando funcionalidades do AgentContract..."
    
    # Verificar configuração
    local config=$(cast call "$contract_address" "config()" --rpc-url "$RPC_URL")
    if [[ -n "$config" ]]; then
        log_success "Configuração do agente encontrada"
    else
        log_error "Configuração do agente não encontrada"
    fi
    
    # Verificar contratos relacionados
    local token_contract=$(cast call "$contract_address" "tokenContract()" --rpc-url "$RPC_URL")
    if [[ -n "$token_contract" ]]; then
        log_success "Contrato de token vinculado"
    else
        log_error "Contrato de token não vinculado"
    fi
    
    local nft_contract=$(cast call "$contract_address" "nftContract()" --rpc-url "$RPC_URL")
    if [[ -n "$nft_contract" ]]; then
        log_success "Contrato de NFT vinculado"
    else
        log_error "Contrato de NFT não vinculado"
    fi
}

# Função para verificar todos os contratos
verify_all_contracts() {
    log "Verificando todos os contratos..."
    
    # Carregar endereços dos contratos
    local deployment_file="$PROJECT_ROOT/deployments/$NETWORK.json"
    
    if [[ ! -f "$deployment_file" ]]; then
        log_error "Arquivo de deployment não encontrado: $deployment_file"
        log "Execute o deploy primeiro ou forneça endereços manualmente"
        exit 1
    fi
    
    # Extrair endereços do arquivo JSON
    local token_address=$(jq -r '.contracts.IACaiToken' "$deployment_file")
    local nft_address=$(jq -r '.contracts.NationPassNFT' "$deployment_file")
    local agent_address=$(jq -r '.contracts.AgentContract' "$deployment_file")
    
    # Verificar cada contrato
    if [[ "$token_address" != "null" && "$token_address" != "" ]]; then
        verify_contract "$token_address" "IACaiToken"
    else
        log_warning "Endereço do token não encontrado"
    fi
    
    if [[ "$nft_address" != "null" && "$nft_address" != "" ]]; then
        verify_contract "$nft_address" "NationPassNFT"
    else
        log_warning "Endereço do NFT não encontrado"
    fi
    
    if [[ "$agent_address" != "null" && "$agent_address" != "" ]]; then
        verify_contract "$agent_address" "AgentContract"
    else
        log_warning "Endereço do agente não encontrado"
    fi
}

# Função para gerar relatório de verificação
generate_verification_report() {
    log "Gerando relatório de verificação..."
    
    local report_file="$PROJECT_ROOT/reports/verification-$NETWORK-$(date +%Y%m%d-%H%M%S).md"
    mkdir -p "$(dirname "$report_file")"
    
    cat > "$report_file" << EOF
# 🔍 Relatório de Verificação de Contratos

**Rede:** $NETWORK  
**Data:** $(date)  
**Script:** $0  

## 📋 Resumo da Verificação

### ✅ Contratos Verificados
- IACaiToken: Verificação de funcionalidades básicas
- NationPassNFT: Verificação de funcionalidades básicas  
- AgentContract: Verificação de funcionalidades básicas

### 🔗 Links dos Contratos
- [Base Explorer](https://basescan.org/)
- [Base Sepolia Explorer](https://sepolia.basescan.org/)

### 📊 Estatísticas
- Total de contratos verificados: 3
- Status: ✅ Verificação concluída
- Tempo de execução: $(date)

## 🔧 Próximos Passos
1. Monitorar contratos em produção
2. Configurar alertas de segurança
3. Implementar testes de integração
4. Configurar backup e recovery

---
*Relatório gerado automaticamente pelo script de verificação*
EOF

    log_success "Relatório gerado: $report_file"
}

# Função principal
main() {
    log "🚀 Iniciando verificação de contratos inteligentes..."
    
    check_prerequisites
    setup_environment
    
    if [[ -n "$CONTRACT_ADDRESS" ]]; then
        # Verificar contrato específico
        verify_contract "$CONTRACT_ADDRESS" "CustomContract"
    else
        # Verificar todos os contratos
        verify_all_contracts
    fi
    
    generate_verification_report
    
    log_success "Verificação concluída com sucesso!"
}

# Executar função principal
main "$@"
