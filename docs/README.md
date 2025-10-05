# IaC AI Agent

Um agente inteligente para análise, revisão e otimização de código Infrastructure as Code (IaC) com foco em Terraform.

## 🚀 Features

- ✅ **Análise de Terraform**: Parse e validação de código HCL
- 🔒 **Segurança**: Integração com Checkov para detecção de vulnerabilidades
- 🔑 **IAM Analysis**: Análise especializada de políticas e permissões
- 🤖 **AI-Powered**: Usa LLM para sugestões contextualizadas
- 💰 **Otimização de Custo**: Recomendações para redução de gastos
- 📊 **PR Scoring**: Atribuição de score de qualidade para pull requests
- 🔗 **GitHub Integration**: Webhooks e comentários automáticos
- 📚 **Knowledge Base**: Base de best practices e módulos recomendados

## 🏗️ Arquitetura

Consulte [ARCHITECTURE.md](./ARCHITECTURE.md) para detalhes completos da arquitetura.

```
┌─────────────┐
│   GitHub    │
│   Webhook   │
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│    REST API         │
│  (handlers.go)      │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│   Services Layer    │
│ analysis.go         │
│ review.go           │
└──────┬──────────────┘
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
┌─────────────┐    ┌──────────────┐
│  Analyzers  │    │     LLM      │
│  terraform  │    │    client    │
│  checkov    │    │prompt_builder│
│  iam        │    └──────────────┘
└─────────────┘
       │
       ▼
┌─────────────────────┐
│  Scorer/Suggester   │
│  pr_scorer.go       │
│  cost_optimizer.go  │
│  security_advisor.go│
└─────────────────────┘
```

## 🛠️ Instalação

### Pré-requisitos

- Go 1.21+
- Docker (opcional)
- Checkov instalado (`pip install checkov`)
- Terraform instalado
- Token GitHub
- API Key de LLM (OpenAI ou Anthropic)

### Setup Local

```bash
# Clone o repositório
git clone <repo-url>
cd iac-ai-agent

# Instale dependências
go mod download

# Configure variáveis de ambiente
cp .env.example .env
# Edite .env com suas credenciais

# Execute setup
make setup

# Execute o agente
make run
```

### Docker

```bash
# Build
docker build -t iac-ai-agent .

# Run
docker run -p 8080:8080 --env-file .env iac-ai-agent
```

### Docker Compose

```bash
docker-compose up -d
```

## 📝 Configuração

### Variáveis de Ambiente

```bash
# LLM Configuration
LLM_PROVIDER=openai          # openai ou anthropic
LLM_API_KEY=sk-xxx...        # Sua API key
LLM_MODEL=gpt-4              # Modelo a usar

# GitHub Configuration
GITHUB_TOKEN=ghp_xxx...      # Token do GitHub
GITHUB_WEBHOOK_SECRET=xxx    # Secret do webhook

# Analysis Configuration
CHECKOV_ENABLED=true
IAM_ANALYSIS_ENABLED=true
COST_OPTIMIZATION_ENABLED=true

# Server Configuration
PORT=8080
LOG_LEVEL=info
```

### app.yaml

Ver exemplo em `configs/app.yaml`

## 🚦 Uso

### Como API REST

```bash
# Health check
curl http://localhost:8080/health

# Analisar código Terraform
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "repository": "org/repo",
    "path": "infrastructure/",
    "content": "<terraform-code>"
  }'

# Review de PR
curl -X POST http://localhost:8080/review \
  -H "Content-Type: application/json" \
  -d '{
    "repository": "org/repo",
    "pr_number": 123
  }'
```

### Integração com GitHub

1. Vá em Settings → Webhooks → Add webhook
2. **Payload URL**: `https://your-domain.com/webhook/github`
3. **Content type**: `application/json`
4. **Secret**: Seu webhook secret
5. **Events**: Pull request, Push
6. Save

O agente comentará automaticamente nos PRs com análises e sugestões.

## 🧪 Testes

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Com coverage
make test-coverage
```

## 📊 Exemplo de Saída

```json
{
  "score": 85,
  "analysis": {
    "terraform": {
      "resources": 12,
      "modules": 3,
      "valid": true
    },
    "security": {
      "critical": 0,
      "high": 1,
      "medium": 3,
      "low": 5
    },
    "iam": {
      "overly_permissive": false,
      "recommendations": [...]
    }
  },
  "suggestions": [
    {
      "type": "security",
      "severity": "high",
      "message": "S3 bucket is publicly accessible",
      "recommendation": "Add bucket_acl = \"private\"",
      "file": "main.tf",
      "line": 45
    },
    {
      "type": "cost",
      "severity": "medium",
      "message": "Consider using spot instances",
      "recommendation": "Add spot_price parameter",
      "estimated_savings": "$450/month"
    }
  ]
}
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📄 Licença

MIT License - veja [LICENSE](../LICENSE) para detalhes.

## 🙏 Agradecimentos

- [Checkov](https://www.checkov.io/) - Security scanning
- [Terraform](https://www.terraform.io/) - IaC platform
- [OpenAI](https://openai.com/) - LLM capabilities
