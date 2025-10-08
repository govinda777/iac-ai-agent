#!/bin/bash

# Script para executar testes de integração NFT Access
# Valida efetivamente o acesso a NFTs Nation Pass

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Função para verificar se uma variável está configurada
check_env_var() {
    if [ -z "${!1}" ]; then
        print_error "$1 não está configurada"
        return 1
    else
        print_success "$1 está configurada"
        return 0
    fi
}

# Função para verificar pré-requisitos
check_prerequisites() {
    print_info "Verificando pré-requisitos..."
    
    # Verificar se Go está instalado
    if ! command -v go &> /dev/null; then
        print_error "Go não está instalado"
        exit 1
    fi
    print_success "Go está instalado"
    
    # Verificar se o arquivo .env existe
    if [ ! -f ".env" ]; then
        print_warning "Arquivo .env não encontrado"
        print_info "Copiando env.example para .env..."
        cp env.example .env
        print_warning "Configure suas variáveis no arquivo .env antes de continuar"
        exit 1
    fi
    print_success "Arquivo .env encontrado"
    
    # Carregar variáveis do .env
    source .env
    
    # Verificar variáveis obrigatórias
    print_info "Verificando variáveis obrigatórias..."
    
    local missing_vars=()
    
    if ! check_env_var "WALLET_ADDRESS"; then
        missing_vars+=("WALLET_ADDRESS")
    fi
    
    if ! check_env_var "NATION_NFT_CONTRACT"; then
        missing_vars+=("NATION_NFT_CONTRACT")
    fi
    
    if ! check_env_var "BASE_RPC_URL"; then
        missing_vars+=("BASE_RPC_URL")
    fi
    
    if [ ${#missing_vars[@]} -gt 0 ]; then
        print_error "Variáveis obrigatórias não configuradas:"
        for var in "${missing_vars[@]}"; do
            echo "  - $var"
        done
        print_info "Configure essas variáveis no arquivo .env"
        exit 1
    fi
    
    print_success "Todos os pré-requisitos estão OK"
}

# Função para executar testes básicos
run_basic_tests() {
    print_info "Executando testes básicos de integração..."
    
    export INTEGRATION_TESTS=true
    
    if go test ./test/integration/ -v -run "TestNFTAccessIntegration"; then
        print_success "Testes básicos passaram"
    else
        print_error "Testes básicos falharam"
        return 1
    fi
}

# Função para executar testes com contratos reais
run_contract_tests() {
    print_info "Executando testes com contratos reais..."
    
    export REAL_CONTRACT_TESTS=true
    
    if go test ./test/integration/ -v -run "TestNFTAccessManagerRealContract"; then
        print_success "Testes com contratos reais passaram"
    else
        print_error "Testes com contratos reais falharam"
        return 1
    fi
}

# Função para executar testes específicos Nation Pass
run_nation_pass_tests() {
    print_info "Executando testes específicos Nation Pass..."
    
    export NATION_PASS_TESTS=true
    
    if go test ./test/integration/ -v -run "TestNationPassAccessValidation"; then
        print_success "Testes Nation Pass passaram"
    else
        print_error "Testes Nation Pass falharam"
        return 1
    fi
}

# Função para executar todos os testes
run_all_tests() {
    print_info "Executando todos os testes de integração..."
    
    export INTEGRATION_TESTS=true
    export REAL_CONTRACT_TESTS=true
    export NATION_PASS_TESTS=true
    
    if go test ./test/integration/ -v; then
        print_success "Todos os testes passaram"
    else
        print_error "Alguns testes falharam"
        return 1
    fi
}

# Função para mostrar ajuda
show_help() {
    echo "Uso: $0 [OPÇÃO]"
    echo ""
    echo "Opções:"
    echo "  basic     Executa apenas testes básicos de integração"
    echo "  contract  Executa testes com contratos reais"
    echo "  nation    Executa testes específicos Nation Pass"
    echo "  all       Executa todos os testes (padrão)"
    echo "  check     Verifica apenas pré-requisitos"
    echo "  help      Mostra esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  $0 basic     # Apenas testes básicos"
    echo "  $0 all       # Todos os testes"
    echo "  $0 check     # Verificar configuração"
}

# Função principal
main() {
    local command=${1:-all}
    
    case $command in
        "basic")
            check_prerequisites
            run_basic_tests
            ;;
        "contract")
            check_prerequisites
            run_contract_tests
            ;;
        "nation")
            check_prerequisites
            run_nation_pass_tests
            ;;
        "all")
            check_prerequisites
            run_all_tests
            ;;
        "check")
            check_prerequisites
            print_success "Verificação de pré-requisitos concluída"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Comando inválido: $command"
            show_help
            exit 1
            ;;
    esac
}

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    print_error "Execute este script na raiz do projeto (onde está o go.mod)"
    exit 1
fi

# Executar função principal
main "$@"
