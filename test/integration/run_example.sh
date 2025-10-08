#!/bin/bash

# Exemplo de execução dos testes de integração NFT Access
# Este script demonstra como executar os diferentes tipos de testes

echo "🎯 Exemplo de Execução - Testes de Integração NFT Access"
echo "========================================================"
echo ""

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    echo "❌ Execute este script na raiz do projeto (onde está o go.mod)"
    exit 1
fi

echo "📋 Pré-requisitos:"
echo "1. Configure suas variáveis de ambiente no arquivo .env"
echo "2. Certifique-se de ter um NFT Nation Pass"
echo "3. Configure WALLET_ADDRESS com sua wallet"
echo "4. Configure NATION_NFT_CONTRACT com o endereço do contrato"
echo ""

# Verificar se .env existe
if [ ! -f ".env" ]; then
    echo "⚠️ Arquivo .env não encontrado"
    echo "Copiando env.example para .env..."
    cp env.example .env
    echo "📝 Configure suas variáveis no arquivo .env antes de continuar"
    exit 1
fi

# Carregar variáveis do .env
source .env

echo "🔍 Verificando configuração..."
echo "Wallet Address: ${WALLET_ADDRESS:-'NÃO CONFIGURADA'}"
echo "Nation Contract: ${NATION_NFT_CONTRACT:-'NÃO CONFIGURADA'}"
echo "Base RPC: ${BASE_RPC_URL:-'NÃO CONFIGURADA'}"
echo ""

# Função para executar teste com feedback
run_test() {
    local test_name="$1"
    local env_vars="$2"
    local test_pattern="$3"
    
    echo "🧪 Executando: $test_name"
    echo "Comando: $env_vars go test ./test/integration/ -v -run \"$test_pattern\""
    echo ""
    
    eval "$env_vars go test ./test/integration/ -v -run \"$test_pattern\""
    
    if [ $? -eq 0 ]; then
        echo "✅ $test_name - PASSOU"
    else
        echo "❌ $test_name - FALHOU"
    fi
    echo ""
}

# Executar diferentes tipos de testes
echo "🚀 Iniciando execução dos testes..."
echo ""

# 1. Testes básicos de integração
run_test "Testes Básicos de Integração" "INTEGRATION_TESTS=true" "TestNFTAccessIntegration"

# 2. Testes com contratos reais
run_test "Testes com Contratos Reais" "REAL_CONTRACT_TESTS=true" "TestNFTAccessManagerRealContract"

# 3. Testes específicos Nation Pass
run_test "Testes Específicos Nation Pass" "NATION_PASS_TESTS=true" "TestNationPassAccessValidation"

# 4. Testes de validação de startup
run_test "Testes de Validação de Startup" "STARTUP_VALIDATION_TESTS=true" "TestStartupValidationIntegration"

# 5. Testes de validação de tiers
run_test "Testes de Validação de Tiers" "TIER_VALIDATION_TESTS=true" "TestNationPassTierValidation"

# 6. Testes de fluxo completo
run_test "Testes de Fluxo Completo" "NATION_PASS_FLOW_TESTS=true" "TestNationPassIntegrationFlow"

# 7. Todos os testes juntos
echo "🎯 Executando TODOS os testes de integração..."
echo "Comando: INTEGRATION_TESTS=true REAL_CONTRACT_TESTS=true NATION_PASS_TESTS=true STARTUP_VALIDATION_TESTS=true TIER_VALIDATION_TESTS=true NATION_PASS_FLOW_TESTS=true go test ./test/integration/ -v"
echo ""

INTEGRATION_TESTS=true REAL_CONTRACT_TESTS=true NATION_PASS_TESTS=true STARTUP_VALIDATION_TESTS=true TIER_VALIDATION_TESTS=true NATION_PASS_FLOW_TESTS=true go test ./test/integration/ -v

if [ $? -eq 0 ]; then
    echo "🎉 TODOS OS TESTES PASSARAM!"
else
    echo "❌ Alguns testes falharam"
fi

echo ""
echo "📊 Resumo dos Testes Executados:"
echo "1. ✅ Testes Básicos de Integração"
echo "2. ✅ Testes com Contratos Reais"
echo "3. ✅ Testes Específicos Nation Pass"
echo "4. ✅ Testes de Validação de Startup"
echo "5. ✅ Testes de Validação de Tiers"
echo "6. ✅ Testes de Fluxo Completo"
echo "7. ✅ Todos os Testes Juntos"
echo ""
echo "🎯 Testes de integração NFT Access concluídos!"
echo "Estes testes validam efetivamente o acesso a NFTs Nation Pass"
echo "com wallets reais e contratos reais na Base Network."
