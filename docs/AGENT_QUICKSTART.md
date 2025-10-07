# 🚀 Quick Start - Sistema de Agentes

## 🤖 O Que Você Precisa Saber

Quando você inicia o IaC AI Agent, **um agente é criado automaticamente** para você. Você não precisa fazer nada!

---

## ✨ Criação Automática

### O que acontece no startup:

```
🚀 Starting IaC AI Agent v1.0.0

🔍 Executando validações de startup...
✅ LLM Connection
✅ Privy.io Credentials
✅ Base Network
✅ Nation.fun NFT

📦 Inicializando serviços...
✅ Agent Service inicializado

🤖 Verificando agente padrão...
   Owner: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

ℹ️  Nenhum agente encontrado
✨ Criando novo agente automaticamente...

🎨 Template: general-purpose
   ✓ Configurando LLM (GPT-4)
   ✓ Habilitando todas análises
   ✓ Configurando personalidade
   ✓ Carregando knowledge base
   ✓ Definindo limites

✅ Novo agente criado automaticamente!
   ID: agent-abc123-def456
   Name: IaC Agent - 0x742d35
   Template: General Purpose

🤖 Agente pronto!
```

---

## 📋 O Que Seu Agente Pode Fazer

Seu agente automático vem configurado com:

### ✅ Análises Habilitadas

- ✅ Terraform Analysis
- ✅ Checkov Security Scanning
- ✅ IAM Policy Analysis
- ✅ Cost Analysis & Optimization
- ✅ Drift Detection
- ✅ Preview Analysis (terraform plan)
- ✅ Secrets Scanning

### 🧠 Habilidades

- ✅ Gerar código Terraform
- ✅ Gerar documentação
- ✅ Refatorar código
- ✅ Sugerir arquitetura
- ✅ Recomendar módulos
- ✅ Sugerir otimizações
- ✅ Alertas de segurança

### 💬 Personalidade

- **Estilo**: Professional
- **Tom**: Encorajador
- **Verbosidade**: Balanceada
- **Usa Emojis**: ✅ Sim
- **Explica Raciocínio**: ✅ Sim
- **Dá Exemplos**: ✅ Sim

### 📊 Limites Padrão

- **Requests/hora**: 100
- **Requests/dia**: 1000
- **Tokens/request**: 4000
- **Custo máximo/request**: $0.50
- **Custo máximo/dia**: $10.00

---

## 🎯 Como Usar Seu Agente

### 1. Análise Simples (Usa Agente Automático)

```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" { bucket = \"my-bucket\" }"
  }'

# Seu agente automático será usado! ✨
```

### 2. Ver Informações do Seu Agente

```bash
curl http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer your-token"

# Response:
{
  "agents": [
    {
      "id": "agent-abc123-def456",
      "name": "IaC Agent - 0x742d35",
      "version": "1.0.0",
      "status": "active",
      "config": { ... },
      "capabilities": { ... },
      "metrics": {
        "total_requests": 0,
        "successful_requests": 0,
        "total_cost_usd": 0.00
      }
    }
  ]
}
```

### 3. Customizar Seu Agente (Opcional)

```bash
curl -X PATCH http://localhost:8080/api/v1/agents/agent-abc123-def456 \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "personality": {
      "use_emojis": false,
      "tone": "formal",
      "verbosity": "concise"
    },
    "config": {
      "detail_level": "detailed"
    }
  }'
```

---

## 🎨 Criar Agentes Adicionais (Opcional)

Se você quiser agentes especializados:

### Agente de Segurança

```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "security-specialist",
    "name": "Security Auditor",
    "description": "Specialized in security audits"
  }'
```

### Agente de Custos

```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "cost-optimizer",
    "name": "Cost Optimizer",
    "description": "Find savings opportunities"
  }'
```

### Usar Agente Específico

```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent-security-xyz",
    "code": "..."
  }'
```

---

## 📊 Monitorar Seu Agente

```bash
# Ver métricas
curl http://localhost:8080/api/v1/agents/agent-abc123/metrics \
  -H "Authorization: Bearer your-token"

# Response:
{
  "total_requests": 150,
  "successful_requests": 145,
  "failed_requests": 5,
  "total_tokens_used": 245000,
  "total_cost_usd": 4.85,
  "average_response_time": 3.2,
  "average_cost_per_request": 0.032,
  "average_user_rating": 4.7
}
```

---

## ❓ FAQ

### P: Preciso criar um agente manualmente?
**R**: Não! Um agente é criado automaticamente para você no startup.

### P: Posso ter múltiplos agentes?
**R**: Sim! Você pode criar quantos quiser. Cada um com configuração diferente.

### P: O agente automático é bom o suficiente?
**R**: Sim! O template "General Purpose" é ótimo para 90% dos casos.

### P: Quando devo criar agentes customizados?
**R**: Quando você precisar de:
- Foco específico (só segurança, só custos)
- Personalidade diferente (formal vs casual)
- Limites diferentes (mais/menos requests)
- Regras customizadas

### P: Quanto custa rodar um agente?
**R**: Depende do uso:
- ~$0.03 por análise simples
- ~$0.10 por análise com LLM
- ~$0.25 por análise completa

### P: Os dados do agente são privados?
**R**: Sim! Cada agente pertence exclusivamente ao owner (wallet).

---

## 🎯 Próximos Passos

1. ✅ **Use o agente automático** - Já está pronto!
2. 📖 **Leia a doc completa** - [`docs/AGENT_SYSTEM.md`](AGENT_SYSTEM.md)
3. 🎨 **Explore templates** - Crie agentes especializados
4. 📊 **Monitore métricas** - Acompanhe uso e custos
5. ⚙️ **Customize** - Ajuste conforme necessidade

---

**Status**: ✅ Sistema pronto para uso  
**Agente Automático**: ✅ Criado no startup  
**Templates Disponíveis**: 4 (general, security, cost, architecture)  
**Documentação**: `docs/AGENT_SYSTEM.md`

---

🎉 **Seu agente está pronto! Comece a analisar Terraform agora!**
