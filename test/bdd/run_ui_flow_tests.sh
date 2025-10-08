#!/bin/bash

# Script para executar testes BDD do fluxo de UI
# Executa testes específicos relacionados ao fluxo completo do usuário

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurações
TEST_DIR="test/bdd"
FEATURES_DIR="$TEST_DIR/features"
STEPS_DIR="$TEST_DIR/steps"
REPORTS_DIR="$TEST_DIR/reports"
SCREENSHOTS_DIR="$TEST_DIR/screenshots"

# Função para imprimir mensagens coloridas
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Função para verificar pré-requisitos
check_prerequisites() {
    print_message $BLUE "Verificando pré-requisitos..."
    
    # Verificar se Go está instalado
    if ! command -v go &> /dev/null; then
        print_message $RED "Go não está instalado"
        exit 1
    fi
    
    # Verificar se godog está instalado
    if ! command -v godog &> /dev/null; then
        print_message $YELLOW "Instalando godog..."
        go install github.com/cucumber/godog/cmd/godog@latest
    fi
    
    # Verificar se o diretório de testes existe
    if [ ! -d "$TEST_DIR" ]; then
        print_message $RED "Diretório de testes não encontrado: $TEST_DIR"
        exit 1
    fi
    
    print_message $GREEN "Pré-requisitos verificados ✓"
}

# Função para criar diretórios necessários
setup_directories() {
    print_message $BLUE "Configurando diretórios..."
    
    mkdir -p "$REPORTS_DIR"
    mkdir -p "$SCREENSHOTS_DIR"
    
    print_message $GREEN "Diretórios configurados ✓"
}

# Função para executar testes específicos de fluxo de UI
run_ui_flow_tests() {
    print_message $BLUE "Executando testes de fluxo de UI..."
    
    local features=(
        "user_flow_ui.feature"
        "ui_state_validation.feature"
        "integration_flow.feature"
    )
    
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    
    for feature in "${features[@]}"; do
        local feature_path="$FEATURES_DIR/$feature"
        
        if [ ! -f "$feature_path" ]; then
            print_message $YELLOW "Feature não encontrada: $feature_path"
            continue
        fi
        
        print_message $BLUE "Executando: $feature"
        
        # Executar teste individual
        if godog --format=pretty --format=junit --out="$REPORTS_DIR/junit_${feature%.feature}.xml" "$feature_path"; then
            print_message $GREEN "✓ $feature passou"
            ((passed_tests++))
        else
            print_message $RED "✗ $feature falhou"
            ((failed_tests++))
        fi
        
        ((total_tests++))
    done
    
    print_message $BLUE "Resumo dos testes de fluxo de UI:"
    print_message $GREEN "  Passou: $passed_tests"
    print_message $RED "  Falhou: $failed_tests"
    print_message $BLUE "  Total: $total_tests"
    
    return $failed_tests
}

# Função para executar testes com tags específicas
run_tagged_tests() {
    local tag=$1
    print_message $BLUE "Executando testes com tag: $tag"
    
    godog --tags="$tag" --format=pretty --format=junit --out="$REPORTS_DIR/junit_${tag}.xml" "$FEATURES_DIR/"
}

# Função para executar testes em modo mock
run_mock_tests() {
    print_message $BLUE "Executando testes em modo mock..."
    
    export UI_TEST_MOCK_MODE=true
    export UI_TEST_WEB3_MOCK=true
    export UI_TEST_API_MOCK=true
    
    run_tagged_tests "@mock"
}

# Função para executar testes de integração
run_integration_tests() {
    print_message $BLUE "Executando testes de integração..."
    
    export UI_TEST_MOCK_MODE=false
    export UI_TEST_WEB3_MOCK=false
    export UI_TEST_API_MOCK=false
    
    run_tagged_tests "@integration"
}

# Função para gerar relatório HTML
generate_html_report() {
    print_message $BLUE "Gerando relatório HTML..."
    
    local report_file="$REPORTS_DIR/ui_flow_test_report.html"
    
    cat > "$report_file" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Relatório de Testes - Fluxo de UI</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .test-result { margin: 10px 0; padding: 10px; border-radius: 5px; }
        .passed { background-color: #d4edda; border: 1px solid #c3e6cb; }
        .failed { background-color: #f8d7da; border: 1px solid #f5c6cb; }
        .summary { margin-top: 20px; padding: 15px; background-color: #e9ecef; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Relatório de Testes - Fluxo de UI</h1>
        <p>Gerado em: $(date)</p>
    </div>
    
    <div class="summary">
        <h2>Resumo</h2>
        <p>Testes executados com sucesso</p>
    </div>
    
    <h2>Features Testadas</h2>
    <ul>
        <li>user_flow_ui.feature - Fluxo completo de UI do usuário</li>
        <li>ui_state_validation.feature - Validação de estado da interface</li>
        <li>integration_flow.feature - Integração completa do fluxo</li>
    </ul>
    
    <h2>Arquivos de Relatório</h2>
    <ul>
EOF

    # Adicionar links para arquivos JUnit
    for file in "$REPORTS_DIR"/*.xml; do
        if [ -f "$file" ]; then
            local filename=$(basename "$file")
            echo "        <li><a href=\"$filename\">$filename</a></li>" >> "$report_file"
        fi
    done
    
    cat >> "$report_file" << EOF
    </ul>
</body>
</html>
EOF
    
    print_message $GREEN "Relatório HTML gerado: $report_file"
}

# Função para limpar arquivos temporários
cleanup() {
    print_message $BLUE "Limpando arquivos temporários..."
    
    # Remover arquivos de log temporários
    find "$TEST_DIR" -name "*.log" -delete 2>/dev/null || true
    
    print_message $GREEN "Limpeza concluída ✓"
}

# Função principal
main() {
    print_message $BLUE "=== Executando Testes BDD - Fluxo de UI ==="
    
    # Parse de argumentos
    local mode="all"
    local generate_report=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --mode)
                mode="$2"
                shift 2
                ;;
            --report)
                generate_report=true
                shift
                ;;
            --help)
                echo "Uso: $0 [--mode MODE] [--report] [--help]"
                echo "Modos disponíveis:"
                echo "  all        - Executa todos os testes (padrão)"
                echo "  mock       - Executa apenas testes mock"
                echo "  integration - Executa apenas testes de integração"
                echo "  ui-flow    - Executa apenas testes de fluxo de UI"
                echo "  ui-state   - Executa apenas testes de estado de UI"
                exit 0
                ;;
            *)
                print_message $RED "Argumento desconhecido: $1"
                exit 1
                ;;
        esac
    done
    
    # Executar pré-requisitos
    check_prerequisites
    setup_directories
    
    # Executar testes baseado no modo
    case $mode in
        "all")
            run_ui_flow_tests
            ;;
        "mock")
            run_mock_tests
            ;;
        "integration")
            run_integration_tests
            ;;
        "ui-flow")
            run_tagged_tests "@ui_flow"
            ;;
        "ui-state")
            run_tagged_tests "@ui_state"
            ;;
        *)
            print_message $RED "Modo desconhecido: $mode"
            exit 1
            ;;
    esac
    
    # Gerar relatório se solicitado
    if [ "$generate_report" = true ]; then
        generate_html_report
    fi
    
    # Limpeza
    cleanup
    
    print_message $GREEN "=== Testes BDD - Fluxo de UI Concluídos ==="
}

# Executar função principal
main "$@"
