# 🚀 Quick Start - IaC AI Agent

## ⚡ Setup em 5 Minutos

### 1. Clone e Configure

```bash
# Clone o repositório
git clone https://github.com/your-org/iac-ai-agent
cd iac-ai-agent

# Copie o arquivo de configuração
cp configs/app.yaml.example configs/app.yaml

# Configure as variáveis de ambiente
cp .env.example .env
```

### 2. Configure Privy.io (2 minutos)

```bash
# 1. Acesse https://privy.io e crie uma conta
# 2. Crie um novo app
# 3. Copie APP_ID e APP_SECRET

# Adicione ao .env:
echo "PRIVY_APP_ID=seu-app-id" >> .env
echo "PRIVY_APP_SECRET=seu-secret" >> .env
```

### 3. Configure Base Network

```bash
# Adicione ao .env:
echo "BASE_RPC_URL=https://goerli.base.org" >> .env  # Testnet
echo "BASE_CHAIN_ID=84531" >> .env  # Goerli testnet
```

### 4. Deploy Smart Contracts (Testnet)

```bash
# Instale dependências
cd contracts
npm install

# Configure hardhat
npx hardhat init

# Deploy (necessita ETH na Base Goerli)
npx hardhat run scripts/deploy.ts --network baseGoerli

# Anote os endereços dos contratos e adicione ao .env:
echo "NFT_CONTRACT_ADDRESS=0x..." >> ../.env
echo "TOKEN_CONTRACT_ADDRESS=0x..." >> ../.env
```

### 5. Execute o Backend

```bash
# Volte para raiz
cd ..

# Instale dependências Go
go mod download

# Execute
go run cmd/agent/main.go

# Ou com Docker
docker-compose up
```

### 6. Execute os Testes BDD

```bash
# Instale Godog
go install github.com/cucumber/godog/cmd/godog@latest

# Execute todos os testes
godog test/bdd/features/

# Ou teste específico
godog test/bdd/features/nft_purchase.feature
```

---

## 📋 Variáveis de Ambiente Necessárias

```bash
# .env file
# =============

# Server
PORT=8080
HOST=0.0.0.0

# LLM
LLM_PROVIDER=openai
LLM_API_KEY=sk-...
LLM_MODEL=gpt-4

# GitHub (opcional)
GITHUB_TOKEN=ghp_...
GITHUB_WEBHOOK_SECRET=...

# Privy.io (obrigatório para Web3)
PRIVY_APP_ID=...
PRIVY_APP_SECRET=...

# Base Network
BASE_RPC_URL=https://goerli.base.org
BASE_CHAIN_ID=84531

# Smart Contracts
NFT_CONTRACT_ADDRESS=0x...
TOKEN_CONTRACT_ADDRESS=0x...

# Features
ENABLE_NFT_ACCESS=true
ENABLE_TOKEN_PAYMENTS=true
```

---

## 🧪 Testar Fluxos Principais

### Teste 1: Autenticação com Privy

```bash
# Execute o teste
godog test/bdd/features/user_onboarding.feature

# Ou manualmente via API
curl -X POST http://localhost:8080/api/v1/auth/verify \
  -H "Authorization: Bearer privy-token-here" \
  -H "Content-Type: application/json"
```

### Teste 2: Comprar NFT (Simulação)

```bash
# Teste BDD
godog test/bdd/features/nft_purchase.feature

# Ou via API
curl -X POST http://localhost:8080/api/v1/nft/mint \
  -H "Authorization: Bearer token" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0x...",
    "tier_id": 2
  }'
```

### Teste 3: Comprar Tokens

```bash
# Teste BDD
godog test/bdd/features/token_purchase.feature

# Ou via API
curl -X POST http://localhost:8080/api/v1/tokens/buy \
  -H "Authorization: Bearer token" \
  -H "Content-Type: application/json" \
  -d '{
    "package_id": 2
  }'
```

### Teste 4: Análise com LLM

```bash
# Teste BDD
godog test/bdd/features/bot_analysis.feature

# Ou via API
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Authorization: Bearer token" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" { ... }",
    "type": "llm_analysis"
  }'
```

---

## 📊 Endpoints Principais

```yaml
# Autenticação
POST   /api/v1/auth/verify          # Verifica token Privy
GET    /api/v1/auth/user            # Info do usuário

# NFT de Acesso
GET    /api/v1/nft/tiers            # Lista tiers
POST   /api/v1/nft/mint             # Minta NFT
GET    /api/v1/nft/my-access        # Meu acesso
POST   /api/v1/nft/upgrade          # Upgrade tier

# Tokens (IACAI)
GET    /api/v1/tokens/packages      # Pacotes disponíveis
POST   /api/v1/tokens/buy           # Compra tokens
GET    /api/v1/tokens/balance       # Meu saldo
GET    /api/v1/tokens/history       # Histórico

# Onramp (Privy)
POST   /api/v1/onramp/session       # Cria sessão de compra
POST   /api/v1/onramp/initiate      # Inicia pagamento
GET    /api/v1/onramp/status/:id    # Status da transação

# Análise
POST   /api/v1/analyze              # Análise de código
GET    /api/v1/analyze/history      # Histórico
GET    /api/v1/analyze/costs        # Tabela de custos

# Health
GET    /health                      # Health check geral
GET    /health/privy                # Status Privy
GET    /health/base                 # Status Base Network
```

---

## 🐛 Troubleshooting

### Erro: "Privy API retornou erro 401"

```bash
# Verifique as credenciais
echo $PRIVY_APP_ID
echo $PRIVY_APP_SECRET

# Teste manualmente
curl -u "$PRIVY_APP_ID:$PRIVY_APP_SECRET" \
  https://auth.privy.io/api/v1/verification_keys
```

### Erro: "Base Network connection failed"

```bash
# Teste RPC
curl -X POST $BASE_RPC_URL \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Use RPC público se necessário
export BASE_RPC_URL=https://goerli.base.org
```

### Erro: "Smart contract not found"

```bash
# Verifique se os contratos foram deployados
echo $NFT_CONTRACT_ADDRESS
echo $TOKEN_CONTRACT_ADDRESS

# Verifique no explorer
open https://goerli.basescan.org/address/$NFT_CONTRACT_ADDRESS
```

### Erro: "Insufficient balance"

```bash
# Para testnet, obtenha ETH grátis:
# 1. Bridge Goerli ETH: https://goerli.base.org/bridge
# 2. Faucet: https://www.alchemy.com/faucets/base-goerli
```

---

## 📚 Próximos Passos

### Para Desenvolvimento

1. Leia `docs/WEB3_INTEGRATION_GUIDE.md` (guia completo)
2. Leia `docs/IMPLEMENTATION_SUMMARY.md` (resumo técnico)
3. Execute todos os testes BDD
4. Customize os tiers e preços conforme necessário

### Para Produção

1. Deploy contratos na Base Mainnet
2. Configure domínio e SSL
3. Setup monitoring (Sentry, DataDog)
4. Configure backup de banco de dados
5. Teste tudo em staging primeiro

### Para Frontend

```bash
# Exemplo Next.js com Privy
npx create-next-app@latest frontend
cd frontend
npm install @privy-io/react-auth wagmi viem

# Ver exemplos em docs/WEB3_INTEGRATION_GUIDE.md
```

---

## 🎯 Casos de Uso

### 1. Desenvolvedor Individual

```
- Compra NFT Basic ($25)
- Faz 10 análises/hora
- Compra tokens quando necessário
- Usa análise básica + Checkov
```

### 2. Time DevOps

```
- Compra NFTs Pro ($125 cada)
- Integra com CI/CD
- Usa análise com LLM
- Compra pacotes Enterprise de tokens
```

### 3. Empresa

```
- Compra NFTs Enterprise ($500)
- API dedicada com rate limit alto
- Custom knowledge base
- SLA e suporte 24/7
```

---

## 💡 Dicas

1. **Testnet primeiro**: Use Base Goerli para testar antes de ir para mainnet
2. **Cache de análises**: Implemente cache para economizar tokens
3. **Rate limiting**: Ajuste os limites conforme necessário
4. **Monitoring**: Configure alerts para saldo baixo de tokens
5. **Backup**: Sempre faça backup dos private keys dos contratos

---

## 📞 Suporte

- **Docs**: `docs/` (todos os guias)
- **Issues**: GitHub Issues
- **Privy**: https://docs.privy.io
- **Base**: https://docs.base.org
- **Discord**: (criar servidor da comunidade)

---

**Status**: ✅ Pronto para uso  
**Versão**: 1.0.0  
**Última atualização**: 2025-01-15
