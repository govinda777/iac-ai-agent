# 🎨 Integração com Nation.fun

## 📋 Overview

O **IaC AI Agent** requer que você possua um **NFT da [Nation.fun](https://nation.fun/)** para executar a aplicação. Isso garante que apenas membros verificados da comunidade Nation.fun possam operar o agente.

## 🎯 Por Que Nation.fun?

[Nation.fun](https://nation.fun/) é uma plataforma na Base Network que permite criar e gerenciar comunidades on-chain através de NFTs. Ao integrar com Nation.fun, garantimos:

- ✅ **Acesso Exclusivo**: Apenas holders de Nation.fun NFT podem rodar o bot
- ✅ **Comunidade Verificada**: Membros ativos e engajados
- ✅ **Descentralização**: Verificação on-chain via smart contracts
- ✅ **Governança**: Potencial para DAO governance futura

---

## 🔑 Requisitos Obrigatórios

Para executar o IaC AI Agent, você precisa:

### 1. **Nation.fun NFT**
- Você DEVE possuir um NFT da Nation.fun
- O NFT deve estar na sua wallet configurada
- A validação é feita no startup da aplicação

### 2. **Wallet Configurada**
- Endereço da wallet (WALLET_ADDRESS)
- Private key para validação (WALLET_PRIVATE_KEY)
- Saldo de ETH na Base Network para gas

### 3. **Credenciais Privy**
- PRIVY_APP_ID
- PRIVY_APP_SECRET

### 4. **LLM API Key**
- OpenAI ou Anthropic API key funcional

---

## 🚀 Setup Completo

### Passo 1: Obter Nation.fun NFT

```bash
# 1. Acesse https://nation.fun/
# 2. Conecte sua wallet (MetaMask, Coinbase Wallet, etc)
# 3. Escolha uma Nation para participar
# 4. Compre o NFT de membership
# 5. Anote o endereço do contrato do NFT
```

### Passo 2: Configurar Variáveis de Ambiente

Crie o arquivo `.env` na raiz do projeto:

```bash
# =====================================
# NATION.FUN CONFIGURATION (REQUIRED)
# =====================================

# Sua wallet que possui o Nation.fun NFT
WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

# Private key da wallet (apenas para validação de startup)
WALLET_PRIVATE_KEY=0x...

# Endereço do contrato Nation.fun NFT
NATION_NFT_CONTRACT=0x...  # Obter de https://nation.fun/

# Requer NFT para rodar? (true/false)
NATION_NFT_REQUIRED=true

# =====================================
# PRIVY.IO CONFIGURATION (REQUIRED)
# =====================================

PRIVY_APP_ID=your-privy-app-id
PRIVY_APP_SECRET=your-privy-app-secret

# =====================================
# LLM CONFIGURATION (REQUIRED)
# =====================================

LLM_PROVIDER=openai
LLM_API_KEY=sk-...
LLM_MODEL=gpt-4

# =====================================
# BASE NETWORK CONFIGURATION
# =====================================

BASE_RPC_URL=https://mainnet.base.org
BASE_CHAIN_ID=8453

# =====================================
# SMART CONTRACTS (after deployment)
# =====================================

NFT_ACCESS_CONTRACT_ADDRESS=0x...
TOKEN_CONTRACT_ADDRESS=0x...

# =====================================
# FEATURES
# =====================================

ENABLE_NFT_ACCESS=true
ENABLE_TOKEN_PAYMENTS=true
ENABLE_STARTUP_VALIDATION=true
```

### Passo 3: Executar Aplicação

```bash
# 1. Instalar dependências
go mod download

# 2. Executar aplicação
go run cmd/agent/main.go
```

---

## ✅ Validação de Startup

Quando a aplicação inicia, ela executa as seguintes validações **obrigatórias**:

### 1. ✅ Validação de LLM

```
🤖 Validando conexão com LLM...
- Provider: openai
- Model: gpt-4
- Status: Enviando mensagem de teste...
- Resposta recebida: OK
- Latência: 2.3s
- Tokens usados: 8
✅ LLM validado com sucesso
```

**Se falhar**: Aplicação NÃO inicia

### 2. ✅ Validação de Privy.io

```
🔐 Validando credenciais Privy.io...
- App ID: app_abc...
- Status: Credenciais válidas
✅ Privy.io validado com sucesso
```

**Se falhar**: Aplicação NÃO inicia

### 3. ✅ Validação de Base Network

```
🌐 Validando conexão com Base Network...
- RPC URL: https://mainnet.base.org
- Chain ID: 8453 (Base Mainnet)
- Latest Block: 12345678
✅ Base Network validado com sucesso
```

**Se falhar**: Warning, mas continua

### 4. ✅ Validação de Nation.fun NFT

```
🎨 Validando posse do NFT Nation.fun...
- Wallet: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
- NFT Contract: 0x...
- Verificando balance...
- Balance: 1 NFT
✅ NFT Nation.fun validado com sucesso
```

**Se falhar**: Aplicação NÃO inicia

---

## 📊 Relatório de Validação

Após as validações, você verá um relatório:

```
============================================================
📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP
============================================================
✅ Status: PASSOU

📋 Checklist de Validações:
  ✅ LLM Connection
  ✅ Privy.io Credentials
  ✅ Base Network
  ✅ Nation.fun NFT

============================================================
✅ Validação completa - Aplicação iniciando...

🚀 Servidor HTTP iniciado
Address: 0.0.0.0:8080
Environment: production

📚 Documentação: http://localhost:8080/api/docs
❤️ Health Check: http://localhost:8080/health

✨ Aplicação pronta para receber requisições!
Press Ctrl+C to shutdown gracefully
```

---

## ❌ Tratamento de Erros

### Erro: NFT não encontrado

```
❌ Nation.fun NFT validation failed: NFT not found in wallet

💡 Solução:
1. Verifique se você possui um NFT Nation.fun
2. Confirme que WALLET_ADDRESS está correto
3. Confirme que NATION_NFT_CONTRACT está correto
4. Verifique no Basescan: https://basescan.org/address/YOUR_WALLET
```

### Erro: LLM não responde

```
❌ LLM validation failed: request timeout

💡 Solução:
1. Verifique sua API key: echo $LLM_API_KEY
2. Verifique seu saldo de créditos na OpenAI/Anthropic
3. Teste manualmente:
   curl https://api.openai.com/v1/chat/completions \
     -H "Authorization: Bearer $LLM_API_KEY" \
     -d '{"model":"gpt-4","messages":[{"role":"user","content":"test"}]}'
```

### Erro: Privy credentials inválidas

```
❌ Privy validation failed: invalid credentials

💡 Solução:
1. Obtenha novas credenciais em https://privy.io
2. Verifique no dashboard do Privy
3. Confirme que APP_ID e APP_SECRET estão corretos
```

---

## 🔒 Segurança

### ⚠️ IMPORTANTE: Private Key

A `WALLET_PRIVATE_KEY` é usada **APENAS** para validação de startup. Ela:

- ✅ NÃO é usada para assinar transações
- ✅ NÃO é exposta em logs
- ✅ NÃO é enviada para APIs externas
- ✅ Fica apenas em memória durante validação

**Proteção:**

```bash
# Adicione .env ao .gitignore
echo ".env" >> .gitignore

# Nunca comite private keys
git add .gitignore
git commit -m "Protect environment variables"

# Use secrets manager em produção
# AWS Secrets Manager, HashiCorp Vault, etc
```

---

## 🐳 Docker

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# Build
RUN CGO_ENABLED=1 go build -o /iacai-agent ./cmd/agent

# Runtime
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /iacai-agent .

EXPOSE 8080

CMD ["./iacai-agent"]
```

### Docker Compose com Secrets

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      # Nation.fun
      - WALLET_ADDRESS=${WALLET_ADDRESS}
      - NATION_NFT_CONTRACT=${NATION_NFT_CONTRACT}
      
      # Privy
      - PRIVY_APP_ID=${PRIVY_APP_ID}
      
      # LLM
      - LLM_API_KEY=${LLM_API_KEY}
      - LLM_PROVIDER=openai
      
      # Base Network
      - BASE_RPC_URL=https://mainnet.base.org
    
    secrets:
      - wallet_private_key
      - privy_app_secret
    
    restart: unless-stopped

secrets:
  wallet_private_key:
    external: true
  privy_app_secret:
    external: true
```

---

## 🧪 Modo de Desenvolvimento

Para desenvolvimento local, você pode **desabilitar** a validação de Nation.fun NFT:

```bash
# .env.development
NATION_NFT_REQUIRED=false
ENABLE_STARTUP_VALIDATION=false
```

**⚠️ CUIDADO**: Isso deve ser usado **APENAS** em desenvolvimento local!

---

## 📞 Suporte

### Nation.fun
- Website: https://nation.fun/
- Docs: https://docs.nation.fun/
- Discord: https://discord.gg/nation-fun

### IaC AI Agent
- GitHub Issues: https://github.com/gosouza/iac-ai-agent/issues
- Docs: `docs/`

---

## 🗺️ Roadmap

### Futuras Integrações Nation.fun

- [ ] **Token Gating**: Diferentes níveis de acesso baseados em quantidade de NFTs
- [ ] **Staking**: Stake Nation.fun NFTs para benefícios adicionais
- [ ] **Governance**: Participar de decisões sobre o bot via DAO
- [ ] **Rewards**: Ganhar tokens IACAI por usar o bot e contribuir
- [ ] **Marketplace**: Comprar/vender análises customizadas
- [ ] **Social**: Compartilhar análises com a comunidade Nation.fun

---

**Status**: ✅ Pronto para uso  
**Versão**: 1.0.0  
**Última Atualização**: 2025-01-15
