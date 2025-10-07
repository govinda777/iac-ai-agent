# 📦 Dependências do Projeto

## Go Dependencies

### Instalar Todas as Dependências

```bash
# Adicionar dependências Ethereum
go get github.com/ethereum/go-ethereum@latest

# Outras dependências necessárias
go get gopkg.in/yaml.v3
go get github.com/gorilla/mux
go get github.com/rs/cors

# Atualizar go.mod e go.sum
go mod tidy
```

### Dependências Principais

```go
// go.mod
module github.com/gosouza/iac-ai-agent

go 1.21

require (
    // Ethereum & Web3
    github.com/ethereum/go-ethereum v1.13.8
    
    // Configuration
    gopkg.in/yaml.v3 v3.0.1
    
    // HTTP Router
    github.com/gorilla/mux v1.8.1
    
    // CORS
    github.com/rs/cors v1.10.1
    
    // Testing
    github.com/stretchr/testify v1.8.4
    github.com/cucumber/godog v0.14.0
)
```

---

## Smart Contracts Dependencies

```bash
cd contracts

# Inicializar projeto
npm init -y

# Hardhat & Plugins
npm install --save-dev \
  hardhat \
  @nomiclabs/hardhat-ethers \
  @nomiclabs/hardhat-etherscan \
  hardhat-gas-reporter \
  hardhat-deploy

# OpenZeppelin Contracts
npm install @openzeppelin/contracts

# Ethers.js
npm install ethers

# TypeScript (opcional mas recomendado)
npm install --save-dev \
  typescript \
  @types/node \
  ts-node
```

### package.json

```json
{
  "name": "iacai-contracts",
  "version": "1.0.0",
  "scripts": {
    "compile": "hardhat compile",
    "test": "hardhat test",
    "deploy:testnet": "hardhat run scripts/deploy.ts --network baseGoerli",
    "deploy:mainnet": "hardhat run scripts/deploy.ts --network base",
    "verify": "hardhat verify"
  },
  "devDependencies": {
    "@nomiclabs/hardhat-ethers": "^2.2.3",
    "@nomiclabs/hardhat-etherscan": "^3.1.7",
    "@openzeppelin/contracts": "^5.0.1",
    "hardhat": "^2.19.4",
    "ethers": "^6.9.0",
    "typescript": "^5.3.3"
  }
}
```

---

## Frontend Dependencies (Next.js)

```bash
# Criar projeto Next.js
npx create-next-app@latest frontend --typescript --tailwind --app

cd frontend

# Privy SDK
npm install @privy-io/react-auth

# Wagmi & Viem (Web3)
npm install wagmi viem

# UI Components
npm install @radix-ui/react-dialog @radix-ui/react-dropdown-menu
npm install lucide-react
npm install class-variance-authority clsx tailwind-merge

# State Management (opcional)
npm install zustand

# Forms (opcional)
npm install react-hook-form zod
```

---

## Docker Dependencies

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências de build
RUN apk add --no-cache git gcc musl-dev

# Copiar go.mod e go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /iacai-agent ./cmd/agent

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /iacai-agent .
COPY configs/app.yaml configs/app.yaml

EXPOSE 8080

CMD ["./iacai-agent"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - LLM_API_KEY=${LLM_API_KEY}
      - PRIVY_APP_ID=${PRIVY_APP_ID}
      - PRIVY_APP_SECRET=${PRIVY_APP_SECRET}
      - BASE_RPC_URL=${BASE_RPC_URL}
      - NFT_CONTRACT_ADDRESS=${NFT_CONTRACT_ADDRESS}
      - TOKEN_CONTRACT_ADDRESS=${TOKEN_CONTRACT_ADDRESS}
    volumes:
      - ./configs:/root/configs
    restart: unless-stopped

  # Opcional: Redis para cache
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    restart: unless-stopped
```

---

## System Dependencies

### macOS

```bash
# Homebrew
brew install go
brew install node
brew install yarn
brew install git

# Ethereum tools (opcional)
brew install geth
```

### Ubuntu/Debian

```bash
# Go
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Node.js
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Build tools
sudo apt-get install -y build-essential git
```

---

## Verificar Instalação

### Go

```bash
go version
# Esperado: go version go1.21.x

go env GOPATH
# Esperado: /home/user/go ou similar
```

### Node & NPM

```bash
node --version
# Esperado: v20.x.x

npm --version
# Esperado: 10.x.x
```

### Ethereum Tools

```bash
# Verificar se go-ethereum foi instalado
go list -m github.com/ethereum/go-ethereum
# Esperado: github.com/ethereum/go-ethereum v1.13.8
```

---

## Troubleshooting

### Erro: "package github.com/ethereum/go-ethereum not found"

```bash
# Limpar cache
go clean -modcache

# Reinstalar
go get github.com/ethereum/go-ethereum@latest
go mod tidy
```

### Erro: "CGO_ENABLED required"

```bash
# Instalar gcc
# macOS
xcode-select --install

# Ubuntu
sudo apt-get install build-essential

# Habilitar CGO
export CGO_ENABLED=1
go build
```

### Erro: "npm ERR! 404 Not Found"

```bash
# Limpar cache npm
npm cache clean --force

# Reinstalar
rm -rf node_modules package-lock.json
npm install
```

---

## Comandos Úteis

```bash
# Atualizar todas as dependências Go
go get -u ./...
go mod tidy

# Verificar dependências desatualizadas
go list -u -m all

# Adicionar dependência específica
go get package@version

# Remover dependências não usadas
go mod tidy
```

---

**Status**: ✅ Completo  
**Última Atualização**: 2025-01-15
