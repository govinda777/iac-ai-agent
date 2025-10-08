# ⚡ Quick Start Atualizado - IaC AI Agent

Este guia contém os passos essenciais para começar a usar o IaC AI Agent rapidamente.

## 📋 Passo 1: Pré-requisitos

Você vai precisar de:

- Go 1.21+ instalado
- Git instalado
- NFT da Nation.fun
- Conta no Privy.io
- API key da OpenAI

## 📋 Passo 2: Clone e Prepare o Projeto

```bash
# Clone o repositório
git clone https://github.com/gosouza/iac-ai-agent.git
cd iac-ai-agent

# Configure as variáveis de ambiente
touch .env
```

## 📋 Passo 3: Configure as Variáveis Obrigatórias

Edite o arquivo `.env` e adicione (substitua pelos seus valores reais):

```bash
# PRIVY.IO
PRIVY_APP_ID=app_xxxxxxxxxxxxxx
PRIVY_APP_SECRET=privy_secret_xxxxxxxxxxxxxx

# NATION.FUN NFT
WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
WALLET_PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb...
NATION_NFT_CONTRACT=0x1234567890123456789012345678901234567890
NATION_NFT_REQUIRED=true

# LLM (OpenAI)
LLM_PROVIDER=openai
LLM_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxx
LLM_MODEL=gpt-4
```

## 📋 Passo 4: Configure o Arquivo YAML

```bash
# Copie o arquivo de configuração exemplo
cp configs/app.yaml.example configs/app.yaml
```

## 📋 Passo 5: Execute a Aplicação

```bash
# Instale dependências
go mod download

# Execute
go run cmd/agent/main.go
```

## 📋 Passo 6: Verifique o Funcionamento

```bash
# Teste o health check
curl http://localhost:8080/health

# Acesse a documentação
open http://localhost:8080/api/docs
```

## 🔍 O Que Esperar Durante a Inicialização

A aplicação vai executar validações e mostrar:

```
📋 Validando configuração básica... ✅
🤖 Validando conexão com LLM... ✅
🔐 Validando credenciais Privy.io... ✅
🌐 Validando conexão com Base Network... ✅
🎨 Validando posse do NFT Nation.fun... ✅

📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP
✅ Status: PASSOU

🚀 Servidor HTTP iniciado
Address: 0.0.0.0:8080
```

## ❌ Resolução Rápida de Problemas

| Erro | Solução |
|------|---------|
| "Variável obrigatória não configurada" | Verifique todas as variáveis no `.env` |
| "LLM validation failed" | Verifique sua API key e créditos da OpenAI |
| "Nation.fun NFT not found" | Confirme que possui o NFT e o endereço do contrato |
| "Privy validation failed" | Verifique suas credenciais no dashboard do Privy |

## 🔒 Segurança

⚠️ **NUNCA faça commit do arquivo .env!**

🔐 **Proteja sua WALLET_PRIVATE_KEY** - use apenas para validação de startup!

## 📚 Próximos Passos

Para informações mais detalhadas, consulte:
- 📖 [Guia de Instalação Completo](./GUIA_INSTALACAO.md)
- 📖 [Variáveis de Ambiente](./ENVIRONMENT_VARIABLES.md)
- 📖 [Integração Web3](./WEB3_INTEGRATION_GUIDE.md)
