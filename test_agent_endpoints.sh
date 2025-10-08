#!/bin/bash

# Script para testar os novos endpoints do agente
# ================================================

echo "🚀 Testando novos endpoints do IaC AI Agent"
echo "=========================================="

BASE_URL="http://localhost:8080"

# Função para fazer requisição e mostrar resultado
test_endpoint() {
    local endpoint=$1
    local description=$2
    
    echo ""
    echo "📡 Testando: $description"
    echo "URL: $BASE_URL$endpoint"
    echo "----------------------------------------"
    
    response=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$BASE_URL$endpoint")
    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
    body=$(echo "$response" | sed '/HTTP_CODE:/d')
    
    if [ "$http_code" = "200" ]; then
        echo "✅ Status: OK ($http_code)"
        echo "📄 Resposta:"
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo "❌ Status: ERRO ($http_code)"
        echo "📄 Resposta:"
        echo "$body"
    fi
}

# Verificar se o servidor está rodando
echo "🔍 Verificando se o servidor está rodando..."
if curl -s "$BASE_URL/health" > /dev/null; then
    echo "✅ Servidor está rodando!"
else
    echo "❌ Servidor não está rodando. Inicie o servidor primeiro:"
    echo "   go run cmd/agent/main.go"
    exit 1
fi

# Testar endpoints
test_endpoint "/health" "Health Check Básico"
test_endpoint "/agent/health" "Health Check Detalhado do Agente"
test_endpoint "/agent/template" "Comparação Agente vs Template"

echo ""
echo "🎉 Testes concluídos!"
echo ""
echo "📋 Resumo dos endpoints implementados:"
echo "   • GET /health - Health check básico"
echo "   • GET /agent/health - Health check detalhado com dados do agente, NFT Nation e teste de conectividade"
echo "   • GET /agent/template - Comparação entre dados do agente atual e template de referência"
echo ""
echo "🔍 Principais funcionalidades implementadas:"
echo "   ✅ Verificação de NFT da Nation na conta padrão"
echo "   ✅ Teste de request no agente da Nation"
echo "   ✅ Detalhes completos do agente (endereço, configurações)"
echo "   ✅ Comparação com dados do template do agente"
echo "   ✅ Informações de conformidade e recomendações"
