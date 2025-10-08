# ⚡ Setup Rápido - IaC AI Agent

## 🚨 ATENÇÃO: A aplicação NÃO VAI INICIAR sem estas configurações!

---

## 📋 PASSO A PASSO OBRIGATÓRIO

### 1️⃣ Obtenha um NFT da Nation.fun

```
🎨 Acesse: https://nation.fun/
👛 Conecte sua wallet (MetaMask/Coinbase Wallet)
💰 Compre um NFT de membership de qualquer Nation
✅ Anote o endereço do contrato NFT
```

### 2️⃣ Crie conta no Privy.io

```
🔐 Acesse: https://privy.io
📝 Crie uma conta
🆕 Crie um novo app
📋 Copie:
   - App ID (começa com app_...)
   - App Secret (Settings → API Keys)
```

### 3️⃣ Obtenha API Key do OpenAI

```
🤖 Acesse: https://platform.openai.com/api-keys
🔑 Create new secret key
💳 Adicione créditos (mínimo $5)
📋 Configure o provedor LLM para Nation.fun
```

### 4️⃣ Configure o arquivo .env

```bash
# Na raiz do projeto, crie o arquivo .env:
cp .env.example .env

# Edite o arquivo .env com seus valores:
nano .env
```

---

## 🔴 VARIÁVEIS OBRIGATÓRIAS (Copie e Cole no seu .env)

```bash
# ============================================
# 🔴 OBRIGATÓRIAS - APP NÃO INICIA SEM ELAS
# ============================================

# 1. PRIVY.IO
PRIVY_APP_ID=app_xxxxxxxxxxxxxx              # ← Copie do dashboard Privy
PRIVY_APP_SECRET=privy_secret_xxxxxxxxxxxxxx # ← Copie do dashboard Privy

# 2. NATION.FUN NFT
WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb      # ← Sua wallet com NFT
WALLET_PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb... # ← Private key (CUIDADO!)
NATION_NFT_CONTRACT=0x1234567890123456789012345678901234567890 # ← Contrato da Nation
NATION_NFT_REQUIRED=true                                        # ← Deixe true

# 3. LLM (Nation.fun)
LLM_PROVIDER=nation.fun
LLM_MODEL=nation-1
# Não é necessária chave de API - acesso via NFT Nation.fun

# ============================================
# 🟡 RECOMENDADAS
# ============================================

# 4. BASE NETWORK
BASE_RPC_URL=https://mainnet.base.org        # ← Base Mainnet
BASE_CHAIN_ID=8453                           # ← 8453 = Mainnet

# 5. FEATURES
ENABLE_NFT_ACCESS=true
ENABLE_TOKEN_PAYMENTS=true
ENABLE_STARTUP_VALIDATION=true

# ============================================
# 🟢 OPCIONAIS
# ============================================

# 6. SERVER
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=production

# 7. LOGGING
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## ✅ Checklist de Verificação

Antes de executar, confirme:

- [ ] ✅ Você possui um NFT da Nation.fun na sua wallet
- [ ] ✅ `PRIVY_APP_ID` está preenchido (app_...)
- [ ] ✅ `PRIVY_APP_SECRET` está preenchido
- [ ] ✅ `WALLET_ADDRESS` é o endereço da wallet com o NFT
- [ ] ✅ `WALLET_PRIVATE_KEY` está preenchida (começa com 0x)
- [ ] ✅ `NATION_NFT_CONTRACT` é o endereço do contrato Nation.fun
- [ ] ✅ `LLM_PROVIDER` está configurado como `nation.fun`
- [ ] ✅ Você possui um NFT Nation.fun
- [ ] ✅ Arquivo `.env` está na raiz do projeto
- [ ] ✅ `.env` está no `.gitignore` (NUNCA comite!)

---

## 🚀 Executar a Aplicação

```bash
# 1. Instalar dependências
go mod download

# 2. Executar
go run cmd/agent/main.go
```

### O que vai acontecer:

```
	██╗ █████╗  ██████╗     █████╗ ██╗     █████╗  ██████╗ ███████╗███╗   ██╗████████╗
	██║██╔══██╗██╔════╝    ██╔══██╗██║    ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝
	██║███████║██║         ███████║██║    ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   
	██║██╔══██║██║         ██╔══██║██║    ██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   
	██║██║  ██║╚██████╗    ██║  ██║██║    ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   
	╚═╝╚═╝  ╚═╝ ╚═════╝    ╚═╝  ╚═╝╚═╝    ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   
	
	Infrastructure as Code AI Agent v1.0.0
	Powered by LLM | Secured by Privy.io | Running on Base Network

🚀 Starting IaC AI Agent...

🔍 Executando validações de startup...

📋 Validando configuração básica...
✅ Configuração básica validada

🤖 Validando conexão com LLM...
   Provider: openai
   Model: gpt-4
   Testando conexão...
   Resposta recebida: OK
   Latência: 2.3s
   Tokens usados: 8
✅ LLM validado com sucesso

🔐 Validando credenciais Privy.io...
   App ID: app_cmzu...
✅ Privy.io validado com sucesso

🌐 Validando conexão com Base Network...
   RPC URL: https://mainnet.base.org
   Chain ID: 8453
   Latest Block: 12345678
✅ Base Network validado com sucesso

🎨 Validando posse do NFT Nation.fun...
   Wallet: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
   NFT Contract: 0x...
   Verificando balance...
   Balance: 1 NFT
✅ NFT Nation.fun validado com sucesso

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

✅ Todas as validações passaram!
✅ Validação completa - Aplicação iniciando...

📦 Inicializando serviços...
✅ Analysis Service inicializado
✅ Review Service inicializado

🌐 Configurando servidor HTTP...

🚀 Servidor HTTP iniciado
Address: 0.0.0.0:8080
Environment: production

📚 Documentação: http://localhost:8080/api/docs
❤️ Health Check: http://localhost:8080/health

✨ Aplicação pronta para receber requisições!
Press Ctrl+C to shutdown gracefully
```

---

## ❌ Se algo der errado...

### Erro: "variável obrigatória não configurada"

```
💡 Solução:
1. Verifique se o arquivo .env existe na raiz
2. Verifique se todas as variáveis OBRIGATÓRIAS estão preenchidas
3. Não deixe espaços antes/depois do =
   ✅ Correto:   LLM_PROVIDER=nation.fun
   ❌ Errado:    LLM_PROVIDER = nation.fun
```

### Erro: "LLM validation failed"

```
💡 Solução:
1. Verifique sua API key em https://platform.openai.com/api-keys
2. Confirme que tem créditos disponíveis
3. Teste manualmente:
   curl https://api.openai.com/v1/models \
     -H "Authorization: Bearer SEU_LLM_API_KEY"
```

### Erro: "Nation.fun NFT not found"

```
💡 Solução:
1. Confirme que você possui o NFT:
   Acesse: https://basescan.org/address/SEU_WALLET_ADDRESS
2. Verifique se WALLET_ADDRESS está correto
3. Verifique se NATION_NFT_CONTRACT está correto
4. Se não tem NFT, compre em: https://nation.fun/
```

### Erro: "Privy validation failed"

```
💡 Solução:
1. Verifique no dashboard do Privy: https://dashboard.privy.io/
2. Confirme que o App ID começa com "app_"
3. Gere novo App Secret se necessário
```

---

## 🔒 Segurança IMPORTANTE

### ⚠️ NUNCA faça commit do arquivo .env!

```bash
# Verifique se .env está no .gitignore:
cat .gitignore | grep .env

# Se não estiver, adicione:
echo ".env" >> .gitignore
echo ".env.*" >> .gitignore
echo "!.env.example" >> .gitignore

# Confirme que está ignorado:
git status
# .env NÃO deve aparecer na lista
```

### 🔐 Proteja sua WALLET_PRIVATE_KEY

```
⚠️ Esta chave dá ACESSO TOTAL à sua wallet!

✅ Use APENAS para validação de startup
✅ NUNCA compartilhe
✅ NUNCA comite no git
✅ Em produção, use secrets manager (AWS/Vault/K8s)
✅ Considere usar uma wallet separada apenas para isto
```

---

## 📚 Documentação Completa

- 📖 **Variáveis de Ambiente Detalhadas**: [`docs/ENVIRONMENT_VARIABLES.md`](docs/ENVIRONMENT_VARIABLES.md)
- 📖 **Integração Nation.fun**: [`docs/NATION_FUN_INTEGRATION.md`](docs/NATION_FUN_INTEGRATION.md)
- 📖 **Integração Web3**: [`docs/WEB3_INTEGRATION_GUIDE.md`](docs/WEB3_INTEGRATION_GUIDE.md)
- 📖 **Quick Start**: [`docs/QUICKSTART.md`](docs/QUICKSTART.md)

---

## 🎯 Próximos Passos

Após a aplicação iniciar com sucesso:

1. ✅ Teste o health check: `curl http://localhost:8080/health`
2. ✅ Leia a documentação completa em `docs/`
3. ✅ Execute os testes BDD: `godog test/bdd/features/`
4. ✅ Deploy dos smart contracts (se necessário)
5. ✅ Configure frontend com Privy SDK

---

## 💬 Precisa de Ajuda?

- **Nation.fun**: https://nation.fun/
- **Privy.io**: https://docs.privy.io
- **OpenAI**: https://platform.openai.com/docs
- **GitHub Issues**: https://github.com/gosouza/iac-ai-agent/issues
- **Documentação**: `docs/`

---

**🎉 Boa sorte com seu IaC AI Agent!**

Status: ✅ Pronto para uso  
Versão: 1.0.0
