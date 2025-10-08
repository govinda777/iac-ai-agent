# 🚀 Guia de Instalação Completo - IaC AI Agent

## 📋 Pré-requisitos

Antes de começar, certifique-se de ter instalado:

- Go 1.21+ (`go version` para verificar)
- Git
- Make (opcional, para usar os comandos do Makefile)
- Uma conta na Nation.fun com NFT
- Uma conta no Privy.io
- Uma API key da OpenAI

## 🔄 Passo a Passo Detalhado

### 1️⃣ Clone o repositório

```bash
git clone https://github.com/gosouza/iac-ai-agent.git
cd iac-ai-agent
```

### 2️⃣ Obtenha as dependências necessárias

#### NFT da Nation.fun

```bash
# Acesse o site da Nation.fun
open https://nation.fun/

# Conecte sua wallet (MetaMask/Coinbase Wallet)
# Compre um NFT de membership de qualquer Nation
# Anote o endereço do contrato NFT
```

#### Conta no Privy.io

```bash
# Acesse o site do Privy.io
open https://privy.io

# Crie uma conta
# Crie um novo app
# Copie o App ID e App Secret
```

#### API Key da OpenAI

```bash
# Acesse o site da OpenAI
open https://platform.openai.com/api-keys

# Crie uma nova API key
# Adicione créditos (mínimo $5)
# Copie a chave
```

### 3️⃣ Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com base no exemplo:

```bash
# Crie o arquivo .env
touch .env

# Adicione as variáveis necessárias
cat << 'EOF' > .env
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
NATION_NFT_REQUIRED=true                                       # ← Deixe true

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
EOF

# Edite o arquivo .env com seus valores reais
nano .env
```

### 4️⃣ Configure o arquivo de configuração

```bash
# Copie o arquivo de configuração exemplo
cp configs/app.yaml.example configs/app.yaml

# Edite conforme necessário
nano configs/app.yaml
```

### 5️⃣ Instale as dependências do Go

```bash
go mod download
```

### 6️⃣ Verifique que você tem todas as configurações

Antes de executar, verifique:

- [x] Você tem NFT da Nation.fun na wallet
- [x] `PRIVY_APP_ID` e `PRIVY_APP_SECRET` estão configurados
- [x] `WALLET_ADDRESS` e `WALLET_PRIVATE_KEY` estão configurados
- [x] `NATION_NFT_CONTRACT` está configurado
- [x] `LLM_API_KEY` está configurado
- [x] Arquivo `configs/app.yaml` está configurado

### 7️⃣ Compilar e executar a aplicação

#### Método 1: Usando Go diretamente

```bash
go run cmd/agent/main.go
```

#### Método 2: Usando Make (se disponível)

```bash
make run
```

#### Método 3: Compilar e executar o binário

```bash
go build -o bin/iac-ai-agent cmd/agent/main.go
./bin/iac-ai-agent
```

#### Método 4: Docker (se Docker estiver instalado)

```bash
# Construir a imagem
docker build -t iacai-agent .

# Executar o container
docker run -p 8080:8080 \
  --env-file .env \
  iacai-agent
```

### 8️⃣ Verificar se a aplicação está funcionando

```bash
# Teste o health check
curl http://localhost:8080/health

# Veja a documentação da API
open http://localhost:8080/api/docs
```

## 🔍 Verificando o status da aplicação

Ao iniciar, a aplicação vai executar validações de startup e mostrar:

- ✅ Validação da configuração básica
- ✅ Validação da conexão com LLM
- ✅ Validação das credenciais do Privy.io
- ✅ Validação da conexão com Base Network
- ✅ Validação da posse do NFT Nation.fun

## ❌ Troubleshooting

### Erro: Variáveis de ambiente não configuradas

```
💡 Solução:
1. Verifique se o arquivo .env existe na raiz
2. Verifique se todas as variáveis OBRIGATÓRIAS estão preenchidas
3. Não deixe espaços antes/depois do =
```

### Erro: Validação do LLM falhou

```
💡 Solução:
1. Verifique sua API key em https://platform.openai.com/api-keys
2. Confirme que tem créditos disponíveis
3. Teste manualmente:
   curl https://api.openai.com/v1/models -H "Authorization: Bearer SEU_LLM_API_KEY"
```

### Erro: NFT Nation.fun não encontrado

```
💡 Solução:
1. Confirme que você possui o NFT:
   Acesse: https://basescan.org/address/SEU_WALLET_ADDRESS
2. Verifique se WALLET_ADDRESS está correto
3. Verifique se NATION_NFT_CONTRACT está correto
4. Se não tem NFT, compre em: https://nation.fun/
```

### Erro: Validação do Privy falhou

```
💡 Solução:
1. Verifique no dashboard do Privy: https://dashboard.privy.io/
2. Confirme que o App ID começa com "app_"
3. Gere novo App Secret se necessário
```

## 📚 Documentação Adicional

Para mais informações, consulte:

- 📖 [Variáveis de Ambiente](./ENVIRONMENT_VARIABLES.md)
- 📖 [Integração Nation.fun](./NATION_FUN_INTEGRATION.md)
- 📖 [Integração Web3](./WEB3_INTEGRATION_GUIDE.md)
- 📖 [Quick Start](./QUICKSTART.md)
- 📖 [Arquitetura](./ARCHITECTURE.md)
- 📖 [Sistema de Agentes](./AGENT_SYSTEM.md)

## 🔒 Segurança

⚠️ **NUNCA faça commit do arquivo .env!**

```bash
# Verifique se .env está no .gitignore:
cat .gitignore | grep .env

# Se não estiver, adicione:
echo ".env" >> .gitignore
echo ".env.*" >> .gitignore
echo "!.env.example" >> .gitignore
```

🔐 **Proteja sua WALLET_PRIVATE_KEY**

```
⚠️ Esta chave dá ACESSO TOTAL à sua wallet!

✅ Use APENAS para validação de startup
✅ NUNCA compartilhe
✅ NUNCA comite no git
✅ Em produção, use secrets manager (AWS/Vault/K8s)
✅ Considere usar uma wallet separada apenas para isto
```

## 📝 Notas Finais

- A aplicação está pronta para uso em produção (versão 1.0.0)
- O projeto utiliza Go 1.21+ e requer conexão com internet
- Para suporte, use o GitHub Issues ou entre em contato via support@iacai.com
