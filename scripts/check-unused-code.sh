#!/bin/bash

# Script para verificar código não utilizado no projeto
# Adiciona esta verificação ao pre-push hook

echo "🔍 Verificando código não utilizado..."

# Verifica se o staticcheck está instalado
if ! command -v staticcheck &> /dev/null; then
  echo "⚠️  staticcheck não está instalado. Instalando..."
  go install honnef.co/go/tools/cmd/staticcheck@latest
fi

# Lista de arquivos vazios ou com erro de sintaxe
EMPTY_FILES=()

# Verifica arquivos vazios ou com problemas de sintaxe
for file in $(find . -name "*.go" -not -path "./vendor/*" -not -path "./bin/*" -not -path "./.git/*"); do
  if [ ! -s "$file" ]; then
    EMPTY_FILES+=("$file é um arquivo vazio")
  elif ! go vet "$file" &> /dev/null; then
    EMPTY_FILES+=("$file tem erro de sintaxe")
  fi
done

# Executa staticcheck para encontrar código não utilizado
staticcheck_result=$(staticcheck -checks=U1000 ./...)
staticcheck_status=$?

# Se encontrou arquivos vazios ou com erro de sintaxe, falha
if [ ${#EMPTY_FILES[@]} -gt 0 ]; then
  echo "❌ Encontrados arquivos vazios ou com erro de sintaxe:"
  for error in "${EMPTY_FILES[@]}"; do
    echo "   - $error"
  done
  exit 1
fi

# Se encontrou código não utilizado, falha
if [ $staticcheck_status -ne 0 ]; then
  echo "❌ Encontrado código não utilizado:"
  echo "$staticcheck_result"
  exit 1
fi

echo "✅ Nenhum código não utilizado encontrado!"
exit 0
