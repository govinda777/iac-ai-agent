#!/bin/bash

# Script para gerar documentação Swagger

set -e

echo "🔧 Gerando documentação Swagger..."

# Verifica se swag está instalado
if ! command -v swag &> /dev/null; then
    echo "⚠️  swag não encontrado. Instalando..."
    go install github.com/swaggo/swag/cmd/swag@v1.8.12
fi

# Gera documentação
swag init -g cmd/agent/main.go -o docs --parseDependency --parseInternal

echo "✅ Documentação Swagger gerada em ./docs"
echo ""
echo "📚 Acesse: http://localhost:8080/swagger/"
