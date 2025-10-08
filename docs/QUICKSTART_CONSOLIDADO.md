# 🚀 Quick Start Consolidado - IaC AI Agent

<div align="center">

![IaC AI Agent Banner](../img/logo.svg)

<h3>Guia de início rápido para análise inteligente de código Infrastructure as Code</h3>
<h4>Com autenticação Web3 e sistema de agentes IA</h4>

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Privy](https://img.shields.io/badge/Auth-Privy.io-6366F1?style=flat&logo=ethereum)](https://privy.io)
[![Base Network](https://img.shields.io/badge/L2-Base-0052FF?style=flat&logo=coinbase)](https://base.org)
[![Nation.fun](https://img.shields.io/badge/Community-Nation.fun-FF6B6B)](https://nation.fun)

</div>

---

## 📋 Pré-requisitos

### ✅ Obrigatórios
- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Git** - Para clonar o repositório
- **OpenAI API Key** - [Obter aqui](https://platform.openai.com/api-keys)

### 🟡 Recomendados
- **Docker** - Para execução em containers
- **Make** - Para comandos automatizados
- **Node.js 18+** - Para desenvolvimento frontend

---

## ⚡ Setup em 5 Minutos

### 1️⃣ Clone e Prepare

```bash
# Clone o repositório
git clone https://github.com/gosouza/iac-ai-agent.git
cd iac-ai-agent

# Setup automático (instala dependências)
make setup

# Verifique se tudo está OK
make check-env
```

### 2️⃣ Configure suas Credenciais

```bash
# Copie o arquivo de exemplo
cp env.example .env

# Edite com suas credenciais
nano .env
```

**Configure APENAS estas variáveis obrigatórias:**

```bash
# ============================================
# 🔴 ÚNICA VARIÁVEL OBRIGATÓRIA
# ============================================

# LLM (Nation.fun)
LLM_PROVIDER=nation.fun
LLM_MODEL=nation-1
# Não é necessária chave de API - acesso via NFT Nation.fun

# ============================================
# 🟢 JÁ CONFIGURADO AUTOMATICAMENTE
# ============================================

# ✅ Privy.io App ID: cmgh6un8w007bl10ci0tgitwp (hardcoded)
# ✅ Wallet Address: 0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5 (hardcoded)
# ✅ Nation.fun NFT: Validação automática via carteira padrão
# ✅ Secrets: Gerenciados via Git Secrets + Lit Protocol
```

### 3️⃣ Execute a Aplicação

```bash
# Execute com validação automática
make run

# Ou diretamente
go run cmd/agent/main.go
```

### 4️⃣ Verifique o Funcionamento

```bash
# Teste o health check
curl http://localhost:8080/agent/health

# Verifique o status do agente
curl http://localhost:8080/agent/status

# Acesse a documentação da API
open http://localhost:8080/api/docs

# Teste uma análise simples
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" { bucket = \"my-bucket\" }",
    "type": "terraform_analysis"
  }'
```

---

## 🔍 O Que Esperar Durante a Inicialização

A aplicação executa validações rigorosas e mostra:

```bash
🚀 Starting IaC AI Agent v1.0.0

📋 Validando configuração básica...
✅ Verificando variáveis obrigatórias:
  - LLM_PROVIDER: nation.fun (configurado)
  - PRIVY_APP_ID: cmgh6un8w007bl10ci0tgitwp (hardcoded)

🤖 Validando conexão com Nation.fun LLM...
✅ Testando conexão com Nation.fun...
✅ Nation.fun LLM autenticado e funcionando

🔐 Validando credenciais Privy.io...
✅ Privy credentials configuradas
  - App ID: cmgh6un8w...

🌐 Validando conexão com Base Network...
✅ Base Network conectado
  - Chain ID: 8453
  - Latest Block: 12345678

🎨 Validando posse do NFT Nation.fun...
✅ Verificando carteira autorizada...
✅ Consultando API Nation.fun...
✅ NFT Pass válido encontrado
✅ Teste de conectividade bem-sucedido

📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP
✅ Status: PASSOU

🤖 Inicializando sistema de agentes...
✅ Agente padrão criado automaticamente
  - ID: agent-abc123-def456
  - Name: IaC Agent - 0x17eDfB
  - Template: General Purpose

🚀 Servidor HTTP iniciado
Address: 0.0.0.0:8080
```

---

## 🎯 Primeiros Passos

### 1. Teste uma Análise Simples

```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n  versioning {\n    enabled = true\n  }\n}",
    "type": "security_analysis"
  }'
```

### 2. Verifique seu Agente Automático

```bash
curl http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer your-token"
```

### 3. Explore a Documentação da API

Acesse: `http://localhost:8080/api/docs` para ver todos os endpoints disponíveis.

---

## 🧪 Exemplos Práticos

### Análise de Segurança

```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n  versioning {\n    enabled = true\n  }\n}",
    "type": "security_analysis"
  }'
```

### Otimização de Custos

```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_instance\" \"web\" {\n  instance_type = \"t3.large\"\n  ami = \"ami-0c02fb55956c7d316\"\n}",
    "type": "cost_optimization"
  }'
```

### Análise Completa com LLM

```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n}",
    "type": "full_analysis",
    "include_llm": true
  }'
```

### Geração de Código

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Crie um bucket S3 com versionamento e criptografia habilitados",
    "type": "terraform_code"
  }'
```

---

## 🛠️ Comandos Úteis

### Desenvolvimento

```bash
# Executar em modo desenvolvimento
make dev

# Executar testes
make test

# Executar testes BDD
make test-bdd

# Verificar configuração
make check-env

# Formatar código
make fmt

# Executar linter
make lint
```

### Docker

```bash
# Construir imagem Docker
make docker-build

# Executar container
make docker-run

# Usar docker-compose
make docker-compose-up
```

### Smart Contracts

```bash
# Setup Foundry
make contracts-setup

# Executar testes de contratos
make contracts-test

# Deploy em testnet
make contracts-deploy-testnet

# Deploy em mainnet
make contracts-deploy-mainnet
```

---

## ❌ Troubleshooting Rápido

### Problemas Comuns

| Erro | Causa | Solução |
|------|-------|---------|
| **"LLM validation failed"** | API key inválida ou sem créditos | Verificar `LLM_API_KEY` e créditos OpenAI |
| **"Privy validation failed"** | App ID não configurado | App ID já está hardcoded, verificar logs |
| **"Nation.fun NFT validation failed"** | Carteira não possui NFT | NFT já validado automaticamente |
| **"Base Network validation failed"** | Problema de conectividade | Verificar `BASE_RPC_URL` e internet |
| **"Agent creation failed"** | Erro na criação do agente | Verificar configurações de wallet |
| **"Erro de conectividade Nation.fun"** | API Nation.fun indisponível | Configurar `ENABLE_STARTUP_VALIDATION=false` |

### Comandos de Diagnóstico

```bash
# Verificar configuração
make check-env

# Ver logs detalhados
LOG_LEVEL=debug make run

# Testar conectividade
curl http://localhost:8080/agent/health

# Verificar variáveis de ambiente
env | grep -E "(LLM|PRIVY|BASE|NATION)"

# Testar API do Nation.fun
curl -X GET "https://api.nation.fun/v1/nft/check/0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"
```

### 🔧 Solução Específica: Erro de Conectividade Nation.fun

Se você receber o erro **"Erro de conectividade Nation.fun"**, isso significa que a API do Nation.fun está temporariamente indisponível ou não existe (modo desenvolvimento). Para resolver:

#### Opção 1: Desabilitar Validação (Recomendado para Desenvolvimento)

```bash
# Editar arquivo .env
nano .env

# Alterar esta linha:
ENABLE_STARTUP_VALIDATION=false

# Testar novamente
make check-env
```

#### Opção 2: Verificar Conectividade

```bash
# Testar conectividade manual
curl -v "https://api.nation.fun/v1/nft/check/0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"

# Se retornar erro de DNS ou timeout, use Opção 1
```

#### Opção 3: Usar Modo Mock (Para Testes)

```bash
# Configurar modo mock no .env
echo "MOCK_MODE=true" >> .env
echo "ENABLE_STARTUP_VALIDATION=false" >> .env
```

### Logs de Debug

```bash
# Executar com logs detalhados
LOG_LEVEL=debug make run

# Ver logs em tempo real
docker logs -f iac-ai-agent | grep -E "(NFT|Nation|validação)"

# Verificar logs de validação
make run 2>&1 | grep -E "(NFT|Nation|validação)"
```

---

## 🎨 Sistema de Agentes Automático

### ✨ Criação Automática

Quando você inicia a aplicação, **um agente é criado automaticamente** para você:

- **ID**: `agent-abc123-def456`
- **Nome**: `IaC Agent - 0x17eDfB`
- **Template**: General Purpose
- **Capabilities**: Todas as análises habilitadas
- **Personalidade**: Professional e encorajador

### 🧠 O Que Seu Agente Pode Fazer

- ✅ Análise de Terraform
- ✅ Checkov Security Scanning
- ✅ IAM Policy Analysis
- ✅ Cost Analysis & Optimization
- ✅ Drift Detection
- ✅ Preview Analysis
- ✅ Secrets Scanning
- ✅ Geração de código
- ✅ Documentação automática
- ✅ Sugestões de arquitetura

### 📊 Monitorar Seu Agente

```bash
# Ver métricas do agente
curl http://localhost:8080/api/v1/agents/agent-abc123/metrics \
  -H "Authorization: Bearer your-token"

# Response:
{
  "total_requests": 150,
  "successful_requests": 145,
  "total_cost_usd": 4.85,
  "average_response_time": 3.2,
  "average_user_rating": 4.7
}
```

---

## 🔐 Segurança e Web3

### Autenticação Automática

- **Privy.io**: App ID já configurado
- **Base Network**: Conectividade automática
- **Nation.fun NFT**: Validação automática via carteira padrão
- **Secrets**: Gerenciados via Git Secrets + Lit Protocol

### Validação de NFT Pass

O sistema valida automaticamente:
1. ✅ Carteira padrão autorizada (`0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5`)
2. ✅ Posse de NFT Pass válido via API Nation.fun
3. ✅ Teste de conectividade com Nation.fun
4. ✅ Criação automática do agente

---

## 📚 Próximos Passos

### Para Desenvolvimento

1. 📖 **[Guia de Instalação Completo](GUIA_INSTALACAO.md)** - Setup detalhado
2. 📖 **[Sistema de Agentes](AGENT_SYSTEM.md)** - Como funciona o sistema de IA
3. 📖 **[Exemplos Práticos](EXEMPLOS_PRATICOS.md)** - Casos de uso reais
4. 📖 **[Integração Web3](WEB3_INTEGRATION_GUIDE.md)** - Guia completo Web3

### Para Produção

1. 📖 **[Arquitetura](ARCHITECTURE.md)** - Design técnico
2. 📖 **[Deploy Guide](DEPLOY_GUIDE.md)** - Guia de deploy
3. 📖 **[Configuração de Variáveis](CONFIGURACAO_VARIAVEIS.md)** - Detalhamento completo
4. 📖 **[Segurança](SECURE_TOKEN_USAGE.md)** - Uso seguro de tokens

### Para Testes

1. 📖 **[Testes](TESTING.md)** - Estratégia e execução
2. 📖 **[Relatório BDD](BDD_TEST_REPORT.md)** - Cobertura de testes
3. 📖 **[Modo Validação](VALIDATION_MODE.md)** - Debug e testes

---

## 🎯 Casos de Uso

### 👨‍💻 Desenvolvedor Individual

```bash
# Análise rápida de código
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" { bucket = \"my-bucket\" }",
    "type": "terraform_analysis"
  }'
```

### 🏢 Time DevOps

```bash
# Integração com CI/CD
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "...",
    "type": "full_analysis",
    "include_llm": true,
    "ci_cd_mode": true
  }'
```

### 🏭 Empresa

```bash
# Análise enterprise com métricas
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "...",
    "type": "enterprise_analysis",
    "include_cost_optimization": true,
    "include_security_audit": true,
    "include_compliance_check": true
  }'
```

---

## 📊 Endpoints Principais

### Autenticação
- `POST /api/v1/auth/verify` - Verifica token Privy
- `GET /api/v1/auth/user` - Informações do usuário

### Análise
- `POST /api/v1/analyze` - Análise de código
- `GET /api/v1/analyze/history` - Histórico de análises
- `GET /api/v1/analyze/costs` - Tabela de custos

### Agentes
- `GET /api/v1/agents` - Lista agentes
- `POST /api/v1/agents` - Criar agente
- `GET /api/v1/agents/{id}/metrics` - Métricas do agente

### Health
- `GET /agent/health` - Health check geral do agente
- `GET /agent/status` - Status detalhado do agente
- `GET /agent/capabilities` - Lista de capabilities disponíveis

---

## 🆘 Suporte

- **Issues**: [GitHub Issues](https://github.com/gosouza/iac-ai-agent/issues)
- **Email**: support@iacai.com
- **Discord**: (em breve)
- **Twitter**: [@iacaiagent](https://twitter.com/iacaiagent)

---

<div align="center">
  <p>Made with ❤️ by the IaC AI Agent Team</p>
  <p>
    <strong>Status</strong>: 🚀 Pronto para produção<br>
    <strong>Versão</strong>: 1.0.0<br>
    <strong>Última Atualização</strong>: 2025-01-15
  </p>
</div>

---

## 🎉 Parabéns!

Você configurou com sucesso o IaC AI Agent! 

**Próximos passos:**
1. ✅ Teste uma análise simples
2. 📖 Explore a documentação da API
3. 🧪 Execute os testes BDD
4. 🎨 Customize seu agente conforme necessário

**Seu agente está pronto para analisar código Terraform!** 🚀
