# IaC AI Agent - Arquitetura

## Visão Geral

O IaC AI Agent é um sistema inteligente para análise, revisão e otimização de código Infrastructure as Code (IaC), com foco em Terraform. O agente combina análise estática, verificações de segurança e inteligência artificial para fornecer feedback detalhado sobre infraestrutura como código.

## Componentes Principais

### 1. API REST (`api/rest/`)
- **Responsabilidade**: Expor endpoints HTTP para interação com o agente
- **Endpoints principais**:
  - `POST /analyze` - Analisa código Terraform
  - `POST /review` - Realiza review completo de PR
  - `POST /webhook/github` - Recebe eventos do GitHub
  - `GET /health` - Health check

### 2. Agent (`internal/agent/`)

#### 2.1 Analyzer
- **terraform.go**: Analisa código Terraform (parsing, validação, recursos)
- **checkov.go**: Integração com Checkov para análise de segurança
- **iam.go**: Análise especializada de políticas e permissões IAM

#### 2.2 LLM
- **client.go**: Cliente para integração com LLM (OpenAI, Anthropic, etc)
- **prompt_builder.go**: Constrói prompts contextualizados para o LLM

#### 2.3 Scorer
- **pr_scorer.go**: Calcula score de qualidade e segurança de PRs

#### 2.4 Suggester
- **cost_optimizer.go**: Sugere otimizações de custo
- **security_advisor.go**: Recomendações de segurança

### 3. Platform (`internal/platform/`)

#### 3.1 Cloud Controller
- **knowledge_base.go**: Base de conhecimento sobre best practices
- **module_registry.go**: Registro de módulos Terraform recomendados

#### 3.2 Webhook
- **github_client.go**: Cliente para API do GitHub
- **handlers.go**: Processa eventos de webhook

### 4. Services (`internal/services/`)
- **analysis.go**: Orquestra análise completa de código
- **review.go**: Orquestra processo de review de PR

### 5. Models (`internal/models/`)
- Estruturas de dados compartilhadas
- DTOs para comunicação entre camadas

### 6. Infrastructure (`pkg/`)
- **config**: Gerenciamento de configurações
- **logger**: Logging estruturado
- **utils**: Utilitários gerais

## Fluxo de Dados

### Análise de Código Terraform
```
1. GitHub Webhook → API Handler
2. Handler → Analysis Service
3. Analysis Service → {Terraform Analyzer, Checkov Analyzer, IAM Analyzer}
4. Analyzers → Results
5. Results → LLM Client (contexto + prompt)
6. LLM → Suggestions
7. Suggestions → {Cost Optimizer, Security Advisor}
8. Final Report → PR Scorer
9. Score + Report → GitHub Comment
```

### Review de PR
```
1. PR Criado/Atualizado → Webhook
2. Webhook Handler → Review Service
3. Review Service → Git Diff Analysis
4. Diff → Changed Files → Analyzers
5. Analysis Results → LLM para síntese
6. LLM → Summary + Score
7. Post Comment no GitHub PR
```

## Decisões Arquiteturais

### 1. Separação em Camadas
- **API**: Entrada HTTP
- **Services**: Lógica de orquestração
- **Agent**: Lógica de negócio especializada
- **Platform**: Integrações externas
- **Models**: Contratos de dados

### 2. Análise Multi-Fonte
- Terraform nativo (parsing HCL)
- Checkov (segurança e compliance)
- IAM especializado (políticas)
- LLM (contexto e inteligência)

### 3. Extensibilidade
- Novos analyzers podem ser adicionados facilmente
- Suggesters plugáveis
- Múltiplos provedores LLM suportados

### 4. Observabilidade
- Logging estruturado em todas as camadas
- Métricas de performance
- Rastreamento de análises

## Tecnologias

- **Linguagem**: Go 1.21+
- **LLM**: OpenAI GPT-4 / Anthropic Claude
- **Security Scanner**: Checkov
- **Terraform**: HCL parser nativo
- **GitHub**: API REST + Webhooks
- **Config**: YAML
- **Container**: Docker

## Configurações

### Variáveis de Ambiente
```
LLM_PROVIDER=openai
LLM_API_KEY=sk-xxx
GITHUB_TOKEN=ghp_xxx
CHECKOV_ENABLED=true
LOG_LEVEL=info
PORT=8080
```

### Arquivo de Configuração (app.yaml)
```yaml
server:
  port: 8080
  
llm:
  provider: openai
  model: gpt-4
  temperature: 0.2
  
analysis:
  checkov_enabled: true
  iam_analysis_enabled: true
  cost_optimization_enabled: true
  
scoring:
  min_pass_score: 70
```

## Deployment

### Docker
```bash
docker build -t iac-ai-agent .
docker run -p 8080:8080 --env-file .env iac-ai-agent
```

### Docker Compose
```bash
docker-compose up -d
```

## Segurança

1. **Secrets**: Nunca commitar tokens/keys
2. **Validação**: Input sanitization em todos os endpoints
3. **Rate Limiting**: Proteção contra abuse
4. **HTTPS**: TLS obrigatório em produção
5. **Webhook Validation**: Verificar assinatura GitHub

## Próximos Passos

1. Adicionar suporte a outros IaC (CloudFormation, Pulumi)
2. Cache de análises
3. Dashboard web
4. Métricas de evolução de qualidade ao longo do tempo
5. Integração com CI/CD (GitLab, Bitbucket)
