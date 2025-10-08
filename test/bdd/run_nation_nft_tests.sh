#!/bin/bash

# Script para executar testes BDD específicos de validação de NFT Pass do Nation
# Executa apenas os cenários relacionados à validação de NFT do Nation.fun

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🎨 Executando testes BDD de validação de NFT Pass do Nation${NC}"
echo "=================================================="

# Verificar se estamos no diretório correto
if [ ! -f "test/bdd/features/nation_nft_validation.feature" ]; then
    echo -e "${RED}❌ Arquivo de feature não encontrado${NC}"
    echo "Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Verificar se godog está instalado
if ! command -v godog &> /dev/null; then
    echo -e "${YELLOW}⚠️ Godog não encontrado. Instalando...${NC}"
    go install github.com/cucumber/godog/cmd/godog@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Configurar variáveis de ambiente para teste
export NATION_NFT_REQUIRED=true
export WALLET_ADDRESS=0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5
export NATION_NFT_CONTRACT=0x1234567890123456789012345678901234567890
export LLM_PROVIDER=nation.fun
export LLM_MODEL=nation-1
export PRIVY_APP_ID=cmgh6un8w007bl10ci0tgitwp

echo -e "${GREEN}✅ Variáveis de ambiente configuradas${NC}"
echo "  NATION_NFT_REQUIRED: $NATION_NFT_REQUIRED"
echo "  WALLET_ADDRESS: $WALLET_ADDRESS"
echo "  NATION_NFT_CONTRACT: $NATION_NFT_CONTRACT"
echo ""

# Executar testes BDD específicos do Nation NFT
echo -e "${BLUE}🧪 Executando testes BDD de validação de NFT Pass do Nation...${NC}"
echo ""

# Executar apenas os cenários marcados com @nation_nft
godog run test/bdd/features/nation_nft_validation.feature \
    --format=pretty \
    --tags="@nation_nft" \
    --concurrency=1 \
    --stop-on-failure

# Verificar resultado
if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 Todos os testes de validação de NFT Pass do Nation passaram!${NC}"
    echo ""
    echo -e "${GREEN}✅ Fluxo de validação implementado com sucesso:${NC}"
    echo "  • Validação de carteira padrão autorizada"
    echo "  • Verificação de NFT Pass via API do Nation.fun"
    echo "  • Teste de conectividade com Nation.fun"
    echo "  • Tratamento de erros e casos edge"
    echo "  • Validação em tempo de execução"
    echo ""
    echo -e "${BLUE}📋 Próximos passos:${NC}"
    echo "  1. Configurar NATION_NFT_CONTRACT com endereço real"
    echo "  2. Testar em ambiente de produção"
    echo "  3. Monitorar logs de validação"
    echo "  4. Implementar cache de validação se necessário"
else
    echo ""
    echo -e "${RED}❌ Alguns testes falharam${NC}"
    echo ""
    echo -e "${YELLOW}🔍 Verifique:${NC}"
    echo "  • Configuração das variáveis de ambiente"
    echo "  • Conectividade com API do Nation.fun"
    echo "  • Implementação dos steps BDD"
    echo "  • Logs de erro detalhados"
    exit 1
fi

echo ""
echo -e "${BLUE}📊 Resumo dos testes executados:${NC}"
echo "  • Cenários de validação na inicialização"
echo "  • Cenários de validação de carteira"
echo "  • Cenários de teste de conectividade"
echo "  • Cenários de tratamento de erro"
echo "  • Cenários de configuração"
echo ""
echo -e "${GREEN}✨ Validação de NFT Pass do Nation implementada e testada!${NC}"
