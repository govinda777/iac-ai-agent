# 🔐 Variáveis de Ambiente - IaC AI Agent

## ⚠️ IMPORTANTE: Variáveis OBRIGATÓRIAS para Iniciar a Aplicação

A aplicação **NÃO VAI INICIAR** sem estas variáveis configuradas corretamente:

---

## 🔴 1. PRIVY.IO (OBRIGATÓRIO)

### `PRIVY_APP_ID`
- **O que é**: ID da sua aplicação no Privy.io
- **Como obter**: 
  1. Acesse https://privy.io
  2. Crie uma conta/login
  3. Crie um novo app
  4. Copie o App ID do dashboard
- **Formato**: `app_xxxxxxxxxxxxxx`
- **Exemplo**: `PRIVY_APP_ID=app_cmzuw6w6p0002mg08i6q84v3x`

### `PRIVY_APP_SECRET`
- **O que é**: Chave secreta da sua aplicação Privy
- **Como obter**: Dashboard do Privy → Settings → API Keys
- **Formato**: String alfanumérica longa
- **Exemplo**: `PRIVY_APP_SECRET=privy_secret_xxxxxxxxxxxxxx`
- **⚠️ NUNCA commite este valor!**

---

## 🔴 2. NATION.FUN NFT (OBRIGATÓRIO)

### `WALLET_ADDRESS`
- **O que é**: Endereço da sua wallet Ethereum que possui o NFT da Nation.fun
- **Requisito**: Esta wallet DEVE ter um NFT da [Nation.fun](https://nation.fun/)
- **Formato**: `0x` seguido de 40 caracteres hexadecimais
- **Exemplo**: `WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb`
- **Como obter**: Copie o endereço da sua MetaMask/Coinbase Wallet

### `WALLET_PRIVATE_KEY`
- **O que é**: Chave privada da wallet (usado APENAS para validação no startup)
- **Formato**: `0x` seguido de 64 caracteres hexadecimais
- **Exemplo**: `WALLET_PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
- **⚠️ CRÍTICO: NUNCA compartilhe ou comite esta chave!**
- **Uso**: Apenas para validar que você possui o NFT no startup

### `NATION_NFT_CONTRACT`
- **O que é**: Endereço do contrato do NFT Nation.fun na Base Network
- **Como obter**: https://nation.fun/ → Sua Nation → Ver contrato no Basescan
- **Formato**: `0x` seguido de 40 caracteres hexadecimais
- **Exemplo**: `NATION_NFT_CONTRACT=0x1234567890123456789012345678901234567890`

### `NATION_NFT_REQUIRED`
- **O que é**: Se deve validar posse do NFT no startup
- **Valor**: `true` ou `false`
- **Recomendado**: `true` (produção)
- **Exemplo**: `NATION_NFT_REQUIRED=true`

---

## 🔴 3. LLM API (OBRIGATÓRIO)

### `LLM_API_KEY`
- **O que é**: Chave de API do provedor de LLM (OpenAI ou Anthropic)
- **Como obter**:
  - **OpenAI**: https://platform.openai.com/api-keys
  - **Anthropic**: https://console.anthropic.com/
- **Formato OpenAI**: `sk-` seguido de caracteres
- **Formato Anthropic**: `sk-ant-` seguido de caracteres
- **Exemplo**: `LLM_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
- **⚠️ Esta chave gera custos! Monitore seu uso**

### `LLM_PROVIDER`
- **O que é**: Qual provedor de LLM usar
- **Valores aceitos**: `openai` ou `anthropic`
- **Padrão**: `openai`
- **Exemplo**: `LLM_PROVIDER=openai`

### `LLM_MODEL`
- **O que é**: Modelo específico do LLM a usar
- **OpenAI**: `gpt-4`, `gpt-4-turbo`, `gpt-3.5-turbo`
- **Anthropic**: `claude-3-opus`, `claude-3-sonnet`, `claude-3-haiku`
- **Recomendado**: `gpt-4` (melhor qualidade)
- **Exemplo**: `LLM_MODEL=gpt-4`

---

## 🟡 4. BASE NETWORK (Recomendado)

### `BASE_RPC_URL`
- **O que é**: URL do RPC da Base Network
- **Mainnet**: `https://mainnet.base.org`
- **Testnet (Goerli)**: `https://goerli.base.org`
- **Padrão**: Mainnet
- **Exemplo**: `BASE_RPC_URL=https://mainnet.base.org`

### `BASE_CHAIN_ID`
- **O que é**: ID da chain da Base Network
- **Mainnet**: `8453`
- **Goerli Testnet**: `84531`
- **Sepolia Testnet**: `84532`
- **Exemplo**: `BASE_CHAIN_ID=8453`

---

## 🟢 5. SMART CONTRACTS (Após Deploy)

### `NFT_ACCESS_CONTRACT_ADDRESS`
- **O que é**: Endereço do seu contrato de NFT de acesso
- **Como obter**: Deploy o contrato primeiro (ver `docs/WEB3_INTEGRATION_GUIDE.md`)
- **Exemplo**: `NFT_ACCESS_CONTRACT_ADDRESS=0x...`

### `TOKEN_CONTRACT_ADDRESS`
- **O que é**: Endereço do seu contrato de token IACAI
- **Como obter**: Deploy o contrato primeiro
- **Exemplo**: `TOKEN_CONTRACT_ADDRESS=0x...`

---

## 🟢 6. FEATURES (Opcional)

### `ENABLE_NFT_ACCESS`
- **O que é**: Habilitar sistema de NFT de acesso
- **Valores**: `true` ou `false`
- **Padrão**: `true`
- **Exemplo**: `ENABLE_NFT_ACCESS=true`

### `ENABLE_TOKEN_PAYMENTS`
- **O que é**: Habilitar pagamentos com tokens
- **Valores**: `true` ou `false`
- **Padrão**: `true`
- **Exemplo**: `ENABLE_TOKEN_PAYMENTS=true`

### `ENABLE_STARTUP_VALIDATION`
- **O que é**: Executar validações no startup
- **Valores**: `true` ou `false`
- **Recomendado**: `true`
- **Exemplo**: `ENABLE_STARTUP_VALIDATION=true`

---

## 🟢 7. SERVER (Opcional)

### `PORT`
- **O que é**: Porta do servidor HTTP
- **Padrão**: `8080`
- **Exemplo**: `PORT=8080`

### `HOST`
- **O que é**: Host do servidor
- **Padrão**: `0.0.0.0`
- **Exemplo**: `HOST=0.0.0.0`

### `ENVIRONMENT`
- **O que é**: Ambiente de execução
- **Valores**: `development`, `staging`, `production`
- **Padrão**: `development`
- **Exemplo**: `ENVIRONMENT=production`

---

## 🟢 8. LOGGING (Opcional)

### `LOG_LEVEL`
- **O que é**: Nível de logging
- **Valores**: `debug`, `info`, `warn`, `error`
- **Padrão**: `info`
- **Exemplo**: `LOG_LEVEL=info`

### `LOG_FORMAT`
- **O que é**: Formato dos logs
- **Valores**: `json` ou `text`
- **Padrão**: `json`
- **Exemplo**: `LOG_FORMAT=json`

---

## 📝 Arquivo .env Completo (Template)

Crie um arquivo `.env` na raiz do projeto com este conteúdo:

```bash
# =====================================================
# VARIÁVEIS OBRIGATÓRIAS (SEM ESTAS, APP NÃO INICIA!)
# =====================================================

# 1. PRIVY.IO
PRIVY_APP_ID=app_xxxxxxxxxxxxxx
PRIVY_APP_SECRET=privy_secret_xxxxxxxxxxxxxx

# 2. NATION.FUN NFT
WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
WALLET_PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
NATION_NFT_CONTRACT=0x1234567890123456789012345678901234567890
NATION_NFT_REQUIRED=true

# 3. LLM API
LLM_PROVIDER=openai
LLM_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
LLM_MODEL=gpt-4
LLM_TEMPERATURE=0.2
LLM_MAX_TOKENS=4000

# =====================================================
# VARIÁVEIS RECOMENDADAS
# =====================================================

# 4. BASE NETWORK
BASE_RPC_URL=https://mainnet.base.org
BASE_CHAIN_ID=8453

# 5. SMART CONTRACTS (após deploy)
NFT_ACCESS_CONTRACT_ADDRESS=0x...
TOKEN_CONTRACT_ADDRESS=0x...

# =====================================================
# VARIÁVEIS OPCIONAIS
# =====================================================

# 6. FEATURES
ENABLE_NFT_ACCESS=true
ENABLE_TOKEN_PAYMENTS=true
ENABLE_STARTUP_VALIDATION=true

# 7. SERVER
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=production

# 8. LOGGING
LOG_LEVEL=info
LOG_FORMAT=json

# 9. RATE LIMITING
BASIC_TIER_RATE_LIMIT=10
PRO_TIER_RATE_LIMIT=100
ENTERPRISE_TIER_RATE_LIMIT=1000

# 10. GITHUB (Opcional)
# GITHUB_TOKEN=ghp_...
# GITHUB_WEBHOOK_SECRET=...

# 11. MONITORING (Opcional)
# SENTRY_DSN=https://...
# DATADOG_API_KEY=...
```

---

## ✅ Validação de Startup

Quando você executa `go run cmd/agent/main.go`, a aplicação valida:

```
🔍 Executando validações de startup...

✅ 1. Configuração básica
   - LLM_API_KEY: Configurada
   - PRIVY_APP_ID: Configurada
   - PRIVY_APP_SECRET: Configurada
   - WALLET_ADDRESS: Configurada

✅ 2. LLM Connection
   - Provider: openai
   - Model: gpt-4
   - Status: Enviando mensagem de teste...
   - Resposta: OK
   - Latência: 2.3s
   - Tokens: 8

✅ 3. Privy.io Credentials
   - App ID: app_cmzu...
   - Status: Válido

✅ 4. Base Network
   - RPC: https://mainnet.base.org
   - Chain ID: 8453
   - Latest Block: 12345678

✅ 5. Nation.fun NFT
   - Wallet: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
   - NFT Contract: 0x...
   - Balance: 1 NFT
   - Status: ✅ VÁLIDO

🎉 Todas as validações passaram!
🚀 Servidor HTTP iniciado em http://0.0.0.0:8080
```

---

## ❌ Erros Comuns

### Erro 1: Variável não configurada

```
❌ Failed to load configuration: variável obrigatória não configurada: LLM_API_KEY

💡 Solução:
1. Crie arquivo .env na raiz do projeto
2. Adicione: LLM_API_KEY=sk-...
3. Execute novamente
```

### Erro 2: LLM não responde

```
❌ LLM validation failed: invalid API key

💡 Solução:
1. Verifique sua API key em:
   - OpenAI: https://platform.openai.com/api-keys
   - Anthropic: https://console.anthropic.com/
2. Verifique se tem créditos disponíveis
3. Teste manualmente:
   curl https://api.openai.com/v1/models \
     -H "Authorization: Bearer $LLM_API_KEY"
```

### Erro 3: NFT não encontrado

```
❌ Nation.fun NFT validation failed: NFT not found in wallet

💡 Solução:
1. Compre um NFT em https://nation.fun/
2. Verifique que WALLET_ADDRESS está correto
3. Verifique que NATION_NFT_CONTRACT está correto
4. Confirme no Basescan:
   https://basescan.org/address/YOUR_WALLET
```

### Erro 4: Private key inválida

```
❌ Invalid private key format

💡 Solução:
1. Private key deve começar com 0x
2. Deve ter 66 caracteres (0x + 64 hex)
3. Exemplo: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
4. NUNCA compartilhe sua private key real!
```

---

## 🔒 Segurança

### ⚠️ NUNCA faça commit destas variáveis:

```bash
# Adicione ao .gitignore
echo ".env" >> .gitignore
echo ".env.*" >> .gitignore
echo "!.env.example" >> .gitignore
```

### ✅ Use secrets manager em produção:

```bash
# AWS Secrets Manager
aws secretsmanager create-secret \
  --name iacai-agent/privy-secret \
  --secret-string "your-secret"

# HashiCorp Vault
vault kv put secret/iacai-agent \
  privy_app_secret=xxx \
  wallet_private_key=xxx

# Kubernetes Secrets
kubectl create secret generic iacai-agent \
  --from-literal=privy-app-secret=xxx \
  --from-literal=wallet-private-key=xxx
```

---

## 📞 Precisa de Ajuda?

### Nation.fun NFT
- Website: https://nation.fun/
- Como comprar NFT: https://docs.nation.fun/getting-started

### Privy.io
- Dashboard: https://dashboard.privy.io/
- Docs: https://docs.privy.io/

### OpenAI
- API Keys: https://platform.openai.com/api-keys
- Pricing: https://openai.com/pricing

### Anthropic
- Console: https://console.anthropic.com/
- Docs: https://docs.anthropic.com/

---

## 🎯 Checklist Rápido

Antes de iniciar a aplicação, verifique:

- [ ] ✅ Arquivo `.env` criado na raiz
- [ ] ✅ `PRIVY_APP_ID` configurado
- [ ] ✅ `PRIVY_APP_SECRET` configurado
- [ ] ✅ `WALLET_ADDRESS` configurado (com Nation.fun NFT)
- [ ] ✅ `WALLET_PRIVATE_KEY` configurado
- [ ] ✅ `NATION_NFT_CONTRACT` configurado
- [ ] ✅ `LLM_API_KEY` configurado (com créditos)
- [ ] ✅ `LLM_PROVIDER` configurado
- [ ] ✅ `LLM_MODEL` configurado
- [ ] ✅ `.env` adicionado ao `.gitignore`

Pronto! Execute: `go run cmd/agent/main.go`

---

**Status**: ✅ Documentação completa  
**Última Atualização**: 2025-01-15
