#!/bin/bash

# Demonstração: Testes de Integração SEM CONFIGURAÇÃO
# Este script mostra como os testes funcionam automaticamente

set -e

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🎯 Demonstração: Testes de Integração SEM CONFIGURAÇÃO${NC}"
echo "=============================================================="
echo ""

echo -e "${YELLOW}📋 O que este script demonstra:${NC}"
echo "✅ WALLET_ADDRESS - Usa padrão automaticamente"
echo "✅ NATION_NFT_CONTRACT - Descobre automaticamente"
echo "✅ BASE_RPC_URL - Testa múltiplos RPCs automaticamente"
echo ""

echo -e "${GREEN}🚀 Executando testes com configuração automática...${NC}"
echo ""

# Limpar variáveis de ambiente para garantir que usa configuração automática
unset WALLET_ADDRESS
unset NATION_NFT_CONTRACT
unset BASE_RPC_URL

echo -e "${BLUE}🔍 Configuração descoberta automaticamente:${NC}"
echo "Wallet Address: $(go run -c 'package main; import "github.com/govinda777/iac-ai-agent/pkg/config"; import "fmt"; func main() { fmt.Println(config.GetDefaultWalletAddress()) }')"
echo "Base RPC: $(go run -c 'package main; import "github.com/govinda777/iac-ai-agent/pkg/config"; import "fmt"; func main() { fmt.Println(config.GetDefaultBaseRPC()) }')"
echo "Nation Contract: $(go run -c 'package main; import "github.com/govinda777/iac-ai-agent/pkg/config"; import "fmt"; func main() { fmt.Println(config.GetDefaultNationPassContract()) }')"
echo ""

echo -e "${GREEN}🧪 Executando teste de configuração automática...${NC}"
go test ./test/integration/ -v -run "TestAutoConfiguration"

echo ""
echo -e "${GREEN}🎯 Executando teste de integração sem configuração...${NC}"
INTEGRATION_TESTS=true go test ./test/integration/ -v -run "TestNFTAccessIntegration" -timeout 30s

echo ""
echo -e "${GREEN}🎉 Demonstração concluída!${NC}"
echo ""
echo -e "${YELLOW}📝 Resumo:${NC}"
echo "✅ Os testes funcionam SEM configuração manual"
echo "✅ O sistema descobre automaticamente todas as configurações"
echo "✅ Fallbacks robustos garantem que sempre funciona"
echo "✅ Facilidade de uso máxima"
echo ""
echo -e "${BLUE}💡 Para usar em seus próprios testes:${NC}"
echo "INTEGRATION_TESTS=true go test ./test/integration/ -v"
echo ""
echo -e "${BLUE}💡 Para sobrescrever configurações (opcional):${NC}"
echo "WALLET_ADDRESS=0x... INTEGRATION_TESTS=true go test ./test/integration/ -v"
