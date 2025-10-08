# ⚡ Quick Start Atualizado - IaC AI Agent

Este guia contém os passos essenciais para começar a usar o IaC AI Agent rapidamente.

## 🚀 Setup em 5 Minutos

### 1️⃣ Clone e Prepare

```bash
# Clone o repositório
git clone https://github.com/gosouza/iac-ai-agent.git
cd iac-ai-agent

# Setup automático (instala dependências)
make setup

# Copie o arquivo de configuração
cp configs/app.yaml.example configs/app.yaml
```

### 2️⃣ Configure suas Credenciais

Crie o arquivo `.env` com suas credenciais:

```bash
# Copie o arquivo de exemplo
cp env.example .env

# Edite com suas credenciais reais
nano .env
```

**Variáveis obrigatórias no `.env`:**

```bash
# ============================================
# 🔴 ÚNICA VARIÁVEL OBRIGATÓRIA
# ============================================

# LLM (OpenAI)
LLM_PROVIDER=openai                          # ← openai ou anthropic
LLM_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxx   # ← Sua OpenAI API key
LLM_MODEL=gpt-4                              # ← Modelo (gpt-4 recomendado)

# ============================================
# 🟢 JÁ CONFIGURADO AUTOMATICAMENTE
# ============================================

# ✅ Privy.io App ID: cmgh6un8w007bl10ci0tgitwp (hardcoded)
# ✅ Wallet Address: 0x147e832418Cc06A501047019E956714271098b89 (hardcoded)
# ✅ Secrets: Gerenciados via Git Secrets + Lit Protocol
```

### 3️⃣ Execute a Aplicação

```bash
# Execute
make run

# Ou diretamente com Go
go run cmd/agent/main.go
```

### 4️⃣ Verifique o Funcionamento

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

## 🧪 Teste Rápido

```bash
# Teste uma análise simples
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" { bucket = \"my-bucket\" }",
    "type": "terraform_analysis"
  }'
```

## ❌ Resolução Rápida de Problemas

| Erro | Solução |
|------|---------|
| "LLM validation failed" | Verifique sua API key e créditos da OpenAI |
| "Git secret não encontrado" | Execute `git secret reveal` para descriptografar secrets |
| "Lit Protocol error" | Verifique se a wallet está conectada corretamente |
| "Variável obrigatória não configurada" | Verifique se `LLM_API_KEY` está no `.env` |

## 🔒 Segurança

⚠️ **NUNCA faça commit do arquivo .env!**

🔐 **Proteja sua WALLET_PRIVATE_KEY** - use apenas para validação de startup!

## 📚 Próximos Passos

Para informações mais detalhadas, consulte:
- 📖 [Guia de Instalação Completo](./GUIA_INSTALACAO.md)
- 📖 [Variáveis de Ambiente](./ENVIRONMENT_VARIABLES.md)
- 📖 [Integração Web3](./WEB3_INTEGRATION_GUIDE.md)
