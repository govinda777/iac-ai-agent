# 🐳 Instalação com Docker - IaC AI Agent

Este guia explica como instalar e executar o IaC AI Agent usando Docker.

## 📋 Pré-requisitos

- Docker instalado e funcionando
- NFT da Nation.fun
- Conta no Privy.io
- API key da OpenAI

## 🔄 Passo a Passo

### 1️⃣ Clone o repositório

```bash
git clone https://github.com/gosouza/iac-ai-agent.git
cd iac-ai-agent
```

### 2️⃣ Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```bash
# Crie o arquivo .env
touch .env

# Adicione as variáveis necessárias
cat << 'EOF' > .env
# PRIVY.IO
PRIVY_APP_ID=app_xxxxxxxxxxxxxx
PRIVY_APP_SECRET=privy_secret_xxxxxxxxxxxxxx

# NATION.FUN NFT
WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
WALLET_PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb...
NATION_NFT_CONTRACT=0x1234567890123456789012345678901234567890
NATION_NFT_REQUIRED=true

# LLM (Nation.fun)
LLM_PROVIDER=nation.fun
LLM_MODEL=nation-1
# Não é necessária chave de API - acesso via NFT Nation.fun

# BASE NETWORK
BASE_RPC_URL=https://mainnet.base.org
BASE_CHAIN_ID=8453

# FEATURES
ENABLE_NFT_ACCESS=true
ENABLE_TOKEN_PAYMENTS=true
ENABLE_STARTUP_VALIDATION=true

# SERVER
PORT=8080
HOST=0.0.0.0
ENVIRONMENT=production

# LOGGING
LOG_LEVEL=info
LOG_FORMAT=json
EOF

# Edite o arquivo .env com seus valores reais
nano .env
```

### 3️⃣ Build da imagem Docker

```bash
docker build -t iacai-agent -f deployments/Dockerfile .
```

### 4️⃣ Execute o container

```bash
# Executando com variáveis do arquivo .env
docker run -p 8080:8080 \
  --env-file .env \
  --name iacai-container \
  iacai-agent
```

### 5️⃣ Verificar se a aplicação está funcionando

```bash
# Teste o health check
curl http://localhost:8080/health
```

## 🔁 Operações Docker Comuns

### Parar o container

```bash
docker stop iacai-container
```

### Iniciar novamente

```bash
docker start iacai-container
```

### Ver logs

```bash
docker logs iacai-container
```

### Ver logs em tempo real

```bash
docker logs -f iacai-container
```

### Remover o container

```bash
docker rm iacai-container
```

## 🐙 Docker Compose (Opcional)

Se preferir usar Docker Compose, crie um arquivo `docker-compose.yml`:

```bash
cat << 'EOF' > docker-compose.yml
version: '3'

services:
  iac-ai-agent:
    build:
      context: .
      dockerfile: deployments/Dockerfile
    ports:
      - "8080:8080"
    env_file:
      - .env
    restart: unless-stopped
    volumes:
      - ./configs:/app/configs
EOF
```

E execute com:

```bash
docker-compose up -d
```

## ⚙️ Configuração de Recursos

Se necessário, você pode configurar os recursos do container:

```bash
docker run -p 8080:8080 \
  --env-file .env \
  --name iacai-container \
  --cpus=2 \
  --memory=2g \
  iacai-agent
```

## 🔐 Segurança

**IMPORTANTE**: Ao usar Docker em produção:

1. **NÃO inclua credenciais sensíveis na imagem**
2. Use **secrets** do Docker ou serviços de gerenciamento de segredos
3. Considere usar **non-root users** dentro do container

## 🔄 Atualizações

Para atualizar para uma nova versão:

```bash
# Pare o container atual
docker stop iacai-container

# Remova o container
docker rm iacai-container

# Atualize o repositório
git pull

# Reconstrua a imagem
docker build -t iacai-agent -f deployments/Dockerfile .

# Execute novamente
docker run -p 8080:8080 \
  --env-file .env \
  --name iacai-container \
  iacai-agent
```

## 🧪 Validação

A aplicação realizará as validações de startup automaticamente ao iniciar.
