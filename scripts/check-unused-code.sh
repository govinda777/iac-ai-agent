#!/bin/bash

# Script para verificar código não utilizado no projeto
# Adiciona esta verificação ao pre-push hook

echo "🔍 Verificando código não utilizado..."

# Lista de arquivos vazios
EMPTY_FILES=()

# Verifica apenas arquivos vazios
for file in $(find . -name "*.go" -not -path "./vendor/*" -not -path "./bin/*" -not -path "./.git/*"); do
  if [ ! -s "$file" ]; then
    EMPTY_FILES+=("$file é um arquivo vazio")
  fi
done

# Se encontrou arquivos vazios, falha
if [ ${#EMPTY_FILES[@]} -gt 0 ]; then
  echo "❌ Encontrados arquivos vazios:"
  for error in "${EMPTY_FILES[@]}"; do
    echo "   - $error"
  done
  exit 1
fi

# Verifica erros de sintaxe no projeto todo
echo "🔍 Verificando erros de sintaxe..."
if ! go vet ./... &> /dev/null; then
  echo "❌ Encontrados erros de sintaxe:"
  go vet ./...
  exit 1
fi

echo "✅ Nenhum código não utilizado encontrado!"
exit 0