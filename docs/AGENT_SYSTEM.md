
# 🤖 Sistema de Agentes - IaC AI Agent

## 📋 O Que São Agentes?

Um **Agente** no IaC AI Agent é uma instância configurada de IA com personalidade, habilidades e conhecimento específicos. Cada usuário pode ter múltiplos agentes customizados para diferentes propósitos.

## 🎯 Criação Automática

**Quando você inicia a aplicação pela primeira vez, um agente é automaticamente criado para você!**

```
🚀 Starting IaC AI Agent...
✅ Validação completa

🤖 Verificando agente padrão...
ℹ️  Nenhum agente encontrado
✨ Criando novo agente automaticamente...
✅ Novo agente criado automaticamente
   Agent ID: agent-abc123
   Name: IaC Agent - 0x742d35
   Template: general-purpose

📦 Agente pronto para uso!
```

---

## 📐 Anatomia de um Agente

Um agente é composto de 7 componentes principais:

### 1. **Configuração Técnica** (`config`)

Define como o agente se comunica com o LLM e quais análises executar.

```yaml
config:
  # LLM Settings
  llm_provider: openai          # ou anthropic
  llm_model: gpt-4              # Modelo específico
  temperature: 0.2              # Criatividade (0-1)
  max_tokens: 4000              # Tokens por request
  
  # Analysis Features
  enable_checkov: true          # Análise de segurança
  enable_iam_analysis: true     # Análise IAM
  enable_cost_analysis: true    # Análise de custos
  enable_drift_detection: true  # Detecção de drift
  enable_preview_analysis: true # Análise de terraform plan
  enable_secrets_scanning: true # Scan de secrets
  
  # Response Settings
  response_format: json         # json, markdown, text
  include_code_examples: true   # Incluir exemplos de código
  include_references: true      # Incluir referências
  detail_level: standard        # brief, standard, detailed
  language: pt-br               # Idioma das respostas
```

### 2. **Habilidades** (`capabilities`)

Define o que o agente pode fazer.

```yaml
capabilities:
  # Analysis
  can_analyze_terraform: true
  can_analyze_checkov: true
  can_analyze_iam: true
  can_analyze_costs: true
  can_detect_drift: true
  can_analyze_preview: true
  can_scan_secrets: true
  
  # Generation
  can_generate_code: true
  can_generate_tests: false
  can_generate_documentation: true
  can_refactor_code: true
  
  # Advisory
  can_suggest_architecture: true
  can_suggest_modules: true
  can_suggest_optimizations: true
  can_suggest_security: true
  
  # Integrations
  can_integrate_github: true
  can_integrate_ci: false
  can_integrate_slack: false
  
  # Learning
  can_learn_from_feedback: false
  can_adapt_to_context: true
  can_remember_preferences: false
```

### 3. **Personalidade** (`personality`)

Define como o agente se comunica.

```yaml
personality:
  style: professional           # professional, casual, friendly, technical
  tone: encouraging             # formal, informal, encouraging, direct
  verbosity: balanced           # concise, balanced, verbose
  use_emojis: true              # Usar emojis?
  use_humor: false              # Usar humor?
  be_encouraging: true          # Ser encorajador?
  be_directive: false           # Dar comandos diretos?
  
  # Communication
  explain_reasoning: true       # Explicar o porquê
  provide_examples: true        # Dar exemplos
  compare_alternatives: true    # Comparar alternativas
  highlight_risks: true         # Destacar riscos
  
  # Interaction
  ask_clarifying_questions: false  # Fazer perguntas?
  offer_alternatives: true         # Oferecer alternativas?
  suggest_best_practices: true     # Sugerir boas práticas?
```

### 4. **Conhecimento** (`knowledge`)

Define a expertise e regras customizadas do agente.

```yaml
knowledge:
  # Expertise Levels (beginner, intermediate, expert)
  terraform_expertise: expert
  aws_expertise: expert
  azure_expertise: intermediate
  gcp_expertise: intermediate
  security_expertise: expert
  networking_expertise: intermediate
  kubernetes_expertise: intermediate
  database_expertise: intermediate
  
  # Compliance & Industry
  compliance_frameworks:
    - GDPR
    - SOC2
    - HIPAA
    - PCI-DSS
  
  industry_focus:
    - general
  
  architecture_patterns:
    - 3-tier
    - microservices
    - serverless
  
  # Custom Rules
  custom_rules:
    - id: rule-001
      name: "Require backup tags"
      severity: high
      pattern: "aws_db_instance.*"
      message: "Database must have backup enabled"
      enabled: true
  
  preferred_modules:
    - "terraform-aws-modules/vpc/aws"
    - "terraform-aws-modules/eks/aws"
  
  banned_resources:
    - "aws_instance.t1.*"  # Banir instâncias antigas
  
  required_tags:
    - Environment
    - Owner
    - ManagedBy
    - Project
```

### 5. **Limites** (`limits`)

Define restrições de uso do agente.

```yaml
limits:
  # Rate Limits
  max_requests_per_hour: 100
  max_requests_per_day: 1000
  max_concurrent_requests: 5
  
  # Token Limits
  max_tokens_per_request: 4000
  max_tokens_per_day: 100000
  
  # Analysis Limits
  max_files_per_analysis: 50
  max_file_size_mb: 10
  max_resources_per_file: 200
  
  # Cost Limits (USD)
  max_cost_per_request: 0.50
  max_cost_per_day: 10.00
  max_cost_per_month: 200.00
  
  # Time Limits (seconds)
  max_analysis_time_seconds: 300
  request_timeout_seconds: 60
```

### 6. **Métricas** (`metrics`)

Estatísticas de uso e performance do agente.

```yaml
metrics:
  # Usage
  total_requests: 1543
  successful_requests: 1520
  failed_requests: 23
  total_tokens_used: 2450000
  total_cost_usd: 48.50
  
  # Performance
  average_response_time: 3.2        # seconds
  average_tokens_per_request: 1588
  average_cost_per_request: 0.031   # USD
  
  # Quality
  average_user_rating: 4.7          # 0-5
  positive_feedback_rate: 0.94      # 0-1
  issues_detected: 15432
  issues_resolved: 12847
  
  # Time
  last_used: "2025-01-15T10:30:00Z"
  total_uptime: 2592000             # seconds (30 days)
```

### 7. **Metadados**

Informações gerais do agente.

```yaml
id: agent-abc123
name: "My IaC Agent"
version: "1.0.0"
description: "Specialized security agent for production"
owner: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
created_at: "2025-01-01T00:00:00Z"
updated_at: "2025-01-15T10:00:00Z"
status: active                      # active, inactive, training
```

---

## 🎨 Templates Pré-Definidos

O sistema vem com 4 templates prontos:

### 1. **General Purpose** (Recomendado)

Agente versátil para análise completa.

```yaml
Template ID: general-purpose
Uso: Análise geral de Terraform
Features:
  ✅ Checkov Security
  ✅ IAM Analysis
  ✅ Cost Analysis
  ✅ Drift Detection
  ✅ Preview Analysis
  ✅ Secrets Scanning
  ✅ Architecture Suggestions
  ✅ Module Recommendations
```

### 2. **Security Specialist**

Focado em segurança e compliance.

```yaml
Template ID: security-specialist
Uso: Auditoria de segurança
Features:
  ✅ Deep Security Analysis
  ✅ Compliance Frameworks (GDPR, SOC2, HIPAA, PCI-DSS, ISO27001)
  ✅ Secrets Detection
  ✅ IAM Deep Dive
  ❌ Cost Analysis (desabilitado)
Personality:
  - Formal
  - Direto
  - Detalhado
  - Destaca riscos
```

### 3. **Cost Optimizer**

Especializado em otimização de custos.

```yaml
Template ID: cost-optimizer
Uso: Reduzir custos de infraestrutura
Features:
  ✅ Cost Analysis
  ✅ Savings Recommendations
  ✅ Resource Rightsizing
  ✅ Reserved Instances Suggestions
  ✅ Alternative Resources
  ❌ Security Analysis (desabilitado)
Personality:
  - Prático
  - Foco em números
  - Comparações detalhadas
```

### 4. **Architecture Advisor**

Focado em arquitetura e design.

```yaml
Template ID: architecture-advisor
Uso: Melhorar arquitetura
Features:
  ✅ Pattern Detection
  ✅ Architecture Suggestions
  ✅ Module Recommendations
  ✅ Code Refactoring
  ✅ Best Practices
Personality:
  - Verbose
  - Explicativo
  - Educacional
  - Compara alternativas
```

---

## 🚀 Como Usar

### Opção 1: Usar o Agente Automático

Quando você inicia a aplicação, um agente é criado automaticamente para você. Basta usar!

```bash
# Execute a aplicação
go run cmd/agent/main.go

# Faça sua primeira análise
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" { ... }"
  }'

# O agente padrão será usado automaticamente!
```

### Opção 2: Criar Agente Customizado

```bash
# Listar templates disponíveis
curl http://localhost:8080/api/v1/agents/templates

# Criar agente a partir de template
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "security-specialist",
    "name": "My Security Agent",
    "description": "Specialized for production security audits"
  }'

# Usar agente específico em análise
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent-xyz789",
    "code": "..."
  }'
```

### Opção 3: Customizar Agente Existente

```bash
# Atualizar configuração do agente
curl -X PATCH http://localhost:8080/api/v1/agents/agent-abc123 \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "temperature": 0.3,
      "detail_level": "detailed"
    },
    "personality": {
      "use_emojis": false,
      "tone": "formal"
    }
  }'
```

---

## 🎯 Casos de Uso

### Caso 1: Desenvolvedor Individual

```yaml
Template: general-purpose
Configuração:
  - LLM: GPT-4
  - Idioma: pt-br
  - Emojis: true
  - Detail: standard
  
Uso:
  - Análise durante desenvolvimento
  - Review de código próprio
  - Aprendizado de best practices
```

### Caso 2: Time DevOps

```yaml
Múltiplos Agentes:
  1. Agent: general-purpose (desenvolvimento)
  2. Agent: security-specialist (pre-prod)
  3. Agent: cost-optimizer (produção)
  
Workflow:
  1. Dev → Usa general-purpose
  2. PR Review → Usa security-specialist
  3. Deploy → Usa cost-optimizer para validar custos
```

### Caso 3: Empresa Regulada

```yaml
Template: security-specialist
Customizações:
  compliance_frameworks:
    - GDPR
    - HIPAA
    - SOC2
  
  custom_rules:
    - Encryption at rest obrigatório
    - Logging habilitado
    - Multi-AZ mandatório
    - Backup com retenção 90 dias
  
  banned_resources:
    - Recursos públicos
    - Instâncias antigas
```

---

## 📖 API Endpoints

### Listar Templates

```http
GET /api/v1/agents/templates

Response:
[
  {
    "id": "general-purpose",
    "name": "General Purpose IaC Agent",
    "description": "...",
    "is_recommended": true,
    "use_cases": [...]
  }
]
```

### Criar Agente

```http
POST /api/v1/agents
Content-Type: application/json

{
  "template_id": "general-purpose",
  "name": "My Agent",
  "description": "...",
  "overrides": {
    "config": {
      "temperature": 0.3
    }
  }
}
```

### Listar Meus Agentes

```http
GET /api/v1/agents

Response:
{
  "agents": [...],
  "total_count": 3
}
```

### Obter Agente

```http
GET /api/v1/agents/{agent_id}

Response:
{
  "id": "agent-abc123",
  "name": "...",
  "config": {...},
  "capabilities": {...},
  ...
}
```

### Atualizar Agente

```http
PATCH /api/v1/agents/{agent_id}
Content-Type: application/json

{
  "name": "New Name",
  "config": {...},
  "personality": {...}
}
```

### Deletar Agente

```http
DELETE /api/v1/agents/{agent_id}
```

### Analisar com Agente Específico

```http
POST /api/v1/analyze
Content-Type: application/json

{
  "agent_id": "agent-abc123",  # Opcional: usa default se omitido
  "code": "...",
  "analysis_type": "full"
}
```

---

## 💡 Dicas e Boas Práticas

### 1. Comece com o Template

Sempre comece com um template e customize depois. Não crie do zero.

### 2. Um Agente por Propósito

- Desenvolvimento → `general-purpose`
- Segurança → `security-specialist`
- Custos → `cost-optimizer`
- Arquitetura → `architecture-advisor`

### 3. Customize Gradualmente

Não tente customizar tudo de uma vez. Ajuste conforme necessidade:

1. Use padrão 1 semana
2. Ajuste personalidade se necessário
3. Adicione regras customizadas
4. Ajuste limites de custo

### 4. Monitore Métricas

Verifique regularmente:
- Custo por request
- Taxa de sucesso
- Feedback dos usuários
- Problemas detectados vs resolvidos

### 5. Multiple Agents = Multiple Perspectives

Use agentes diferentes para obter perspectivas diferentes sobre o mesmo código.

---

## 🔒 Segurança e Ownership

- ✅ Cada agente pertence a uma wallet (owner)
- ✅ Apenas o owner pode modificar/deletar
- ✅ Agentes são isolados por owner
- ✅ Métricas são privadas do owner
- ✅ Regras customizadas são privadas

---

## 🗂️ Armazenamento

**Desenvolvimento**: In-memory (dados perdidos ao reiniciar)

**Produção**: Recomendado usar banco de dados:
- PostgreSQL para agentes
- Redis para cache de respostas
- S3 para histórico de análises

---

## 📞 Suporte

- **Docs Completas**: `docs/`
- **API Reference**: `docs/API.md`
- **GitHub Issues**: Issues do projeto

---

**Status**: ✅ Sistema completo de agentes implementado  
**Versão**: 1.0.0  
**Última Atualização**: 2025-01-15
