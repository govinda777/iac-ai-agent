# 🎯 Exemplos Práticos de Uso - IaC AI Agent

Este guia contém exemplos práticos de como usar o IaC AI Agent para diferentes cenários.

## 🚀 Primeiro Uso

### 1. Verificar se a API está funcionando

```bash
# Health check
curl http://localhost:8080/health

# Resposta esperada:
# {"status":"ok","version":"1.0.0","timestamp":"2025-01-15T10:30:00Z"}
```

### 2. Teste básico de análise

```bash
# Análise simples de um recurso S3
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n}",
    "type": "terraform_analysis"
  }'
```

## 🔒 Análise de Segurança

### 1. Detecção de Vulnerabilidades

```bash
# Analise um bucket S3 sem criptografia
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"insecure\" {\n  bucket = \"my-insecure-bucket\"\n  versioning {\n    enabled = true\n  }\n}",
    "type": "security_analysis"
  }'
```

### 2. Análise de Políticas IAM

```bash
# Analise uma política IAM permissiva
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_iam_policy\" \"admin\" {\n  name = \"admin-policy\"\n  policy = jsonencode({\n    Version = \"2012-10-17\"\n    Statement = [\n      {\n        Effect = \"Allow\"\n        Action = \"*\"\n        Resource = \"*\"\n      }\n    ]\n  })\n}",
    "type": "iam_analysis"
  }'
```

### 3. Scan com Checkov

```bash
# Execute scan de segurança com Checkov
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n  versioning {\n    enabled = true\n  }\n}",
    "type": "checkov_scan"
  }'
```

## 💰 Otimização de Custos

### 1. Análise de Instâncias EC2

```bash
# Analise custos de instâncias EC2
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_instance\" \"web\" {\n  instance_type = \"t3.large\"\n  ami = \"ami-0c02fb55956c7d316\"\n  \n  tags = {\n    Name = \"web-server\"\n    Environment = \"production\"\n  }\n}",
    "type": "cost_optimization"
  }'
```

### 2. Sugestões de Otimização

```bash
# Peça sugestões de otimização
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_instance\" \"expensive\" {\n  instance_type = \"m5.2xlarge\"\n  ami = \"ami-0c02fb55956c7d316\"\n}",
    "type": "cost_optimization",
    "include_suggestions": true
  }'
```

## 🧠 Análise com LLM

### 1. Análise Completa com IA

```bash
# Análise completa com sugestões inteligentes
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n  versioning {\n    enabled = true\n  }\n  server_side_encryption_configuration {\n    rule {\n      apply_server_side_encryption_by_default {\n        sse_algorithm = \"AES256\"\n      }\n    }\n  }\n}",
    "type": "full_analysis",
    "include_llm": true
  }'
```

### 2. Geração de Código

```bash
# Peça para gerar código Terraform
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Crie um bucket S3 com versionamento, criptografia e logging habilitados",
    "type": "terraform_code"
  }'
```

### 3. Refatoração de Código

```bash
# Peça para refatorar código existente
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_instance\" \"old\" {\n  instance_type = \"t2.micro\"\n  ami = \"ami-0c02fb55956c7d316\"\n}",
    "type": "refactor",
    "prompt": "Modernize this EC2 instance with best practices"
  }'
```

## 🏗️ Análise de Arquitetura

### 1. Validação de Arquitetura

```bash
# Analise uma arquitetura completa
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_vpc\" \"main\" {\n  cidr_block = \"10.0.0.0/16\"\n  \n  tags = {\n    Name = \"main-vpc\"\n  }\n}\n\nresource \"aws_subnet\" \"public\" {\n  vpc_id = aws_vpc.main.id\n  cidr_block = \"10.0.1.0/24\"\n  availability_zone = \"us-west-2a\"\n  \n  tags = {\n    Name = \"public-subnet\"\n  }\n}",
    "type": "architecture_analysis"
  }'
```

### 2. Sugestões de Melhoria

```bash
# Peça sugestões de melhoria arquitetural
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n  ami = \"ami-0c02fb55956c7d316\"\n}",
    "type": "architecture_review",
    "prompt": "Suggest improvements for high availability and scalability"
  }'
```

## 📊 Análise de Drift

### 1. Detecção de Drift

```bash
# Analise diferenças entre código e infraestrutura
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n  versioning {\n    enabled = true\n  }\n}",
    "type": "drift_detection",
    "state_file": "terraform.tfstate"
  }'
```

## 🔍 Análise de Múltiplos Arquivos

### 1. Análise de Projeto Completo

```bash
# Analise um projeto Terraform completo
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "files": [
      {
        "path": "main.tf",
        "content": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n}"
      },
      {
        "path": "variables.tf", 
        "content": "variable \"bucket_name\" {\n  type = string\n  default = \"my-bucket\"\n}"
      }
    ],
    "type": "project_analysis"
  }'
```

## 🎯 Casos de Uso Específicos

### 1. Compliance e Auditoria

```bash
# Verifique compliance com padrões
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n}",
    "type": "compliance_check",
    "standards": ["SOC2", "PCI-DSS", "HIPAA"]
  }'
```

### 2. Análise de Performance

```bash
# Analise performance de recursos
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n  ami = \"ami-0c02fb55956c7d316\"\n}",
    "type": "performance_analysis"
  }'
```

## 📈 Monitoramento e Métricas

### 1. Verificar Status do Agente

```bash
# Ver informações do seu agente
curl http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer your-token"
```

### 2. Ver Métricas de Uso

```bash
# Ver métricas de uso
curl http://localhost:8080/api/v1/agents/agent-id/metrics \
  -H "Authorization: Bearer your-token"
```

### 3. Histórico de Análises

```bash
# Ver histórico de análises
curl http://localhost:8080/api/v1/analyze/history \
  -H "Authorization: Bearer your-token"
```

## 🛠️ Integração com CI/CD

### 1. GitHub Actions

```yaml
name: IaC Analysis
on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run IaC Analysis
        run: |
          curl -X POST http://localhost:8080/api/v1/analyze \
            -H "Content-Type: application/json" \
            -d @analysis-request.json
```

### 2. GitLab CI

```yaml
analyze:
  stage: test
  script:
    - curl -X POST http://localhost:8080/api/v1/analyze \
        -H "Content-Type: application/json" \
        -d @analysis-request.json
```

## 🔧 Troubleshooting

### 1. Erro de Autenticação

```bash
# Verifique se está autenticado
curl http://localhost:8080/api/v1/auth/user \
  -H "Authorization: Bearer your-token"
```

### 2. Erro de Rate Limit

```bash
# Verifique seus limites
curl http://localhost:8080/api/v1/agents/agent-id/limits \
  -H "Authorization: Bearer your-token"
```

### 3. Erro de Validação

```bash
# Valide sua configuração
curl http://localhost:8080/api/v1/validate \
  -H "Content-Type: application/json" \
  -d '{
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n}"
  }'
```

## 📚 Próximos Passos

- 📖 [Sistema de Agentes](AGENT_SYSTEM.md) - Como funciona o sistema de IA
- 📖 [Integração Web3](WEB3_INTEGRATION_GUIDE.md) - Autenticação e pagamentos
- 📖 [Arquitetura](ARCHITECTURE.md) - Design técnico da aplicação
- 📖 [Testes](TESTING.md) - Como testar a aplicação

---

**Status**: ✅ Exemplos prontos para uso  
**Versão**: 1.0.0  
**Última atualização**: 2025-01-15
