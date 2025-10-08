#!/bin/bash

# Script para executar testes BDD com diferentes configurações
# Uso: ./run_bdd_tests.sh [mock|integration|real] [tags]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Função para mostrar ajuda
show_help() {
    echo "Uso: $0 [MODO] [TAGS] [OPÇÕES]"
    echo ""
    echo "MODOS:"
    echo "  mock        - Executa testes usando apenas mocks (padrão)"
    echo "  integration - Executa testes de integração com serviços reais"
    echo "  real        - Executa testes end-to-end com serviços reais"
    echo "  all         - Executa todos os modos sequencialmente"
    echo ""
    echo "TAGS (opcional):"
    echo "  @mock       - Apenas testes com mocks"
    echo "  @integration - Apenas testes de integração"
    echo "  @real       - Apenas testes reais"
    echo "  @unit       - Apenas testes unitários"
    echo "  @performance - Apenas testes de performance"
    echo ""
    echo "OPÇÕES:"
    echo "  --help, -h  - Mostra esta ajuda"
    echo "  --verbose   - Executa com logs verbosos"
    echo "  --parallel  - Executa testes em paralelo"
    echo "  --coverage  - Gera relatório de cobertura"
    echo "  --fail-fast - Para no primeiro erro"
    echo ""
    echo "EXEMPLOS:"
    echo "  $0 mock                    # Testes com mocks"
    echo "  $0 integration @real        # Testes de integração real"
    echo "  $0 all --verbose           # Todos os testes com logs verbosos"
    echo "  $0 mock @unit --coverage   # Testes unitários com cobertura"
}

# Configurações padrão
MODE="mock"
TAGS=""
VERBOSE=false
PARALLEL=false
COVERAGE=false
FAIL_FAST=false
EXTRA_ARGS=""

# Parse argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        mock|integration|real|all)
            MODE="$1"
            shift
            ;;
        @*)
            TAGS="$TAGS $1"
            shift
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --parallel)
            PARALLEL=true
            shift
            ;;
        --coverage)
            COVERAGE=true
            shift
            ;;
        --fail-fast)
            FAIL_FAST=true
            shift
            ;;
        *)
            EXTRA_ARGS="$EXTRA_ARGS $1"
            shift
            ;;
    esac
done

# Configurações baseadas no modo
case $MODE in
    mock)
        export TEST_MODE=mock
        export TEST_USE_REAL_SERVICES=false
        export TEST_MOCK_FAILURES=false
        DEFAULT_TAGS="@mock"
        ;;
    integration)
        export TEST_MODE=integration
        export TEST_USE_REAL_SERVICES=true
        export TEST_MOCK_FAILURES=false
        DEFAULT_TAGS="@integration"
        ;;
    real)
        export TEST_MODE=real
        export TEST_USE_REAL_SERVICES=true
        export TEST_MOCK_FAILURES=false
        DEFAULT_TAGS="@real"
        ;;
    all)
        print_message $BLUE "Executando todos os modos de teste..."
        $0 mock $TAGS $EXTRA_ARGS
        $0 integration $TAGS $EXTRA_ARGS
        $0 real $TAGS $EXTRA_ARGS
        exit 0
        ;;
esac

# Usar tags padrão se nenhuma foi especificada
if [ -z "$TAGS" ]; then
    TAGS="$DEFAULT_TAGS"
fi

# Configurações adicionais
if [ "$VERBOSE" = true ]; then
    export TEST_DEBUG=true
    export TEST_LOG_LEVEL=debug
    EXTRA_ARGS="$EXTRA_ARGS -v"
fi

if [ "$PARALLEL" = true ]; then
    EXTRA_ARGS="$EXTRA_ARGS -parallel 4"
fi

if [ "$COVERAGE" = true ]; then
    EXTRA_ARGS="$EXTRA_ARGS -coverprofile=coverage.out -covermode=atomic"
fi

if [ "$FAIL_FAST" = true ]; then
    EXTRA_ARGS="$EXTRA_ARGS -failfast"
fi

# Configurações de ambiente
export TEST_LOG_FORMAT=json
export TEST_TIMEOUT_SECONDS=30
export TEST_LOAD_COUNT=10

# Configurações específicas para integração
if [ "$MODE" = "integration" ] || [ "$MODE" = "real" ]; then
    if [ -z "$TEST_REAL_NATION_API_KEY" ]; then
        print_message $YELLOW "AVISO: TEST_REAL_NATION_API_KEY não configurado"
        print_message $YELLOW "Alguns testes de integração podem falhar"
    fi
fi

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    print_message $RED "ERRO: Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Verificar se godog está instalado
if ! command -v godog &> /dev/null; then
    print_message $YELLOW "Instalando godog..."
    go install github.com/cucumber/godog/cmd/godog@latest
fi

# Verificar se os testes BDD existem
if [ ! -d "test/bdd" ]; then
    print_message $RED "ERRO: Diretório test/bdd não encontrado"
    exit 1
fi

# Mostrar configuração atual
print_message $BLUE "=== Configuração de Teste ==="
print_message $BLUE "Modo: $MODE"
print_message $BLUE "Tags: $TAGS"
print_message $BLUE "Verbose: $VERBOSE"
print_message $BLUE "Parallel: $PARALLEL"
print_message $BLUE "Coverage: $COVERAGE"
print_message $BLUE "Fail Fast: $FAIL_FAST"
print_message $BLUE "============================="

# Executar testes
print_message $GREEN "Executando testes BDD..."

# Comando base
CMD="godog"

# Adicionar tags se especificadas
if [ -n "$TAGS" ]; then
    CMD="$CMD --tags=\"$TAGS\""
fi

# Adicionar argumentos extras
if [ -n "$EXTRA_ARGS" ]; then
    CMD="$CMD $EXTRA_ARGS"
fi

# Adicionar diretório de features
CMD="$CMD test/bdd/features/"

# Executar comando
print_message $BLUE "Comando: $CMD"
eval $CMD

# Verificar resultado
if [ $? -eq 0 ]; then
    print_message $GREEN "✅ Testes executados com sucesso!"
    
    # Gerar relatório de cobertura se solicitado
    if [ "$COVERAGE" = true ]; then
        print_message $BLUE "Gerando relatório de cobertura..."
        go tool cover -html=coverage.out -o coverage.html
        print_message $GREEN "Relatório de cobertura gerado: coverage.html"
    fi
else
    print_message $RED "❌ Testes falharam!"
    exit 1
fi

# Mostrar resumo
print_message $BLUE "=== Resumo ==="
print_message $BLUE "Modo executado: $MODE"
print_message $BLUE "Tags utilizadas: $TAGS"
print_message $BLUE "Configurações aplicadas:"
print_message $BLUE "  TEST_MODE=$TEST_MODE"
print_message $BLUE "  TEST_USE_REAL_SERVICES=$TEST_USE_REAL_SERVICES"
print_message $BLUE "  TEST_MOCK_FAILURES=$TEST_MOCK_FAILURES"
print_message $BLUE "==============="
