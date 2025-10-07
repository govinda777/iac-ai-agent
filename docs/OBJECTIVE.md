# Objetivo do Projeto - IaC AI Agent

## 🎯 Visão Geral

Um agente de inteligência artificial responsável por analisar resultados de **IAC Preview** (terraform plan/apply) e **Checkov Policies** para propor sugestões de melhorias baseadas em sua base de conhecimento.

**Plataforma Alvo:** [Nation.fun Platform](https://nation.fun/)

## 📥 Inputs Esperados

### 1. Preview de IAC
- **Formato**: JSON do `terraform plan` ou saída de preview do Spacelift/Terraform Cloud
- **Conteúdo**:
  - Recursos que serão criados/modificados/destruídos
  - Mudanças planejadas
  - Erros de validação
  - Warnings

**Exemplo:**
```json
{
  "resource_changes": [
    {
      "address": "aws_s3_bucket.example",
      "mode": "managed",
      "type": "aws_s3_bucket",
      "name": "example",
      "change": {
        "actions": ["create"],
        "before": null,
        "after": {...}
      }
    }
  ],
  "configuration": {...},
  "planned_values": {...}
}
```

### 2. Checkov Policies Result
- **Formato**: JSON do output do Checkov
- **Conteúdo**:
  - Checks passed/failed/skipped
  - Severidades
  - Guias de correção

**Exemplo:**
```json
{
  "summary": {
    "passed": 15,
    "failed": 3,
    "skipped": 2
  },
  "results": {
    "failed_checks": [...],
    "passed_checks": [...],
    "skipped_checks": [...]
  }
}
```

### 3. Base de Conhecimento (Contexto)

#### 3.1. Contexto da Plataforma
- Nome da plataforma: Nation.fun
- Arquitetura padrão
- Recursos comumente utilizados
- Padrões de naming
- Estrutura de tags obrigatórias

#### 3.2. Gerenciamento de Providers
- Providers permitidos (AWS, Azure, GCP, etc.)
- Versões aprovadas de providers
- Configurações de provider recomendadas
- Restrictions e políticas

#### 3.3. Versões Suportadas
- **Terraform**: >= 1.5.0
- **OpenTofu**: >= 1.6.0
- **Providers**:
  - AWS: >= 5.0
  - Azure: >= 3.0
  - GCP: >= 4.0

#### 3.4. Módulos Aprovados
- Registry interno de módulos
- Módulos Terraform Registry recomendados
- Templates de infraestrutura

---

## 🐛 Categoria 1: Bugs / Correções

### 1.1. Correção de Preview
**Objetivo**: Identificar e sugerir correções para erros no preview

**Casos de Uso**:
- ❌ Erro: `Resource 'aws_instance.web' depends on 'aws_vpc.main' which does not exist`
  - ✅ Sugestão: "Crie o recurso aws_vpc.main ou corrija a referência"
  
- ❌ Erro: `Invalid value for variable "instance_type"`
  - ✅ Sugestão: "Use um tipo de instância válido. Valores aceitos: t3.micro, t3.small..."

### 1.2. Correção de Policies (Checkov)
**Objetivo**: Analisar resultados Checkov e categorizar ações

#### Failed Checks
- Identificar checks críticos que bloqueiam deploy
- Priorizar por severidade
- Gerar fix automatizado quando possível

#### Passed Checks
- Validar que best practices estão sendo seguidas
- Registrar para histórico de qualidade

#### Skipped Checks
**Por Exception**: Validar se exception está documentada e aprovada
```yaml
# Exemplo de exception válida
skip_check:
  - CKV_AWS_20  # S3 bucket precisa ser público para website
  reason: "Bucket usado para static website hosting"
  approved_by: "security-team"
  expires: "2025-12-31"
```

**Por Change Window**: Verificar se mudança está dentro da janela permitida
```yaml
change_window:
  start: "2025-10-07T02:00:00Z"
  end: "2025-10-07T06:00:00Z"
  reason: "Maintenance window"
```

**Por Account Test**: Validar que está em ambiente de testes
```yaml
environment: test
account_id: "123456789012"  # Test account
```

### 1.3. Correção de Acesso IAM
**Objetivo**: Identificar e corrigir problemas de permissões

**Detecções**:
- ❌ Admin access (`Action: "*"`)
  - ✅ Sugestão: Especificar apenas actions necessárias
  
- ❌ Resource: `"*"` sem justificativa
  - ✅ Sugestão: Limitar a recursos específicos com ARNs
  
- ❌ Principal: `"*"` (público)
  - ✅ Sugestão: Restringir a principals específicos

### 1.4. Correção de Timeout
**Objetivo**: Identificar recursos que costumam dar timeout

**Análise**:
```yaml
resource: aws_rds_instance.main
problem: "Timeout após 20 minutos em 3 dos últimos 5 applies"
suggestion: |
  - Aumentar timeout para 30 minutos
  - Considerar usar snapshot para faster recovery
  - Verificar se instance_class não é muito grande
```

### 1.5. Recurso Travado
**Objetivo**: Detectar recursos em estado inconsistente

**Cenários**:
- Recurso marcado como "tainted"
- Recurso com lock pendente
- Recurso com apply em andamento há muito tempo

**Sugestões**:
```bash
# Recurso tainted
terraform untaint aws_instance.web

# Lock pendente
terraform force-unlock <lock-id>

# Reimportar recurso
terraform import aws_instance.web i-1234567890
```

### 1.6. Drifted Resources
**Objetivo**: Detectar drift entre código e estado real

**Análise**:
```yaml
resource: aws_s3_bucket.data
drift_detected: true
changes:
  versioning:
    in_code: enabled
    in_cloud: disabled
  tags:
    in_code: {Environment: prod}
    in_cloud: {Environment: prod, ManagedBy: terraform, Owner: john}

suggestions:
  - "Executar 'terraform apply' para alinhar configuração"
  - "Ou atualizar código para incluir tags adicionadas manualmente"
  - "Considerar importar configuração atual com 'terraform import'"
```

---

## 💡 Categoria 2: Melhorias

### 2.1. Diminuição de Custo
**Objetivo**: Identificar oportunidades de economia

**Análises**:

#### Instâncias EC2
```yaml
current: t2.xlarge ($0.1856/hour = $135/month)
suggestion: t3.xlarge ($0.1664/hour = $121/month)
savings: $14/month (10%)
reason: "t3 instances have better price/performance ratio"
```

#### RDS
```yaml
current: db.m5.large (24x7)
suggestion: db.m5.large com Reserved Instance (1 year)
savings: $500/year (30%)
```

#### NAT Gateway
```yaml
current: NAT Gateway ($45/month)
suggestion: "Use NAT Instance para ambientes não-prod"
savings: $30/month
caveat: "Menos HA, ok para dev/test"
```

### 2.2. Refatoração de Arquitetura
**Objetivo**: Sugerir melhorias arquiteturais

**Padrões Detectados**:

#### Anti-pattern: Monolith Stack
```hcl
# ❌ Todos os recursos em um único arquivo
# main.tf com 500 linhas, 80 recursos

# ✅ Sugestão: Separar em módulos
modules/
  networking/
  compute/
  database/
  security/
```

#### Anti-pattern: Hardcoded Values
```hcl
# ❌ Valores hardcoded
resource "aws_instance" "web" {
  instance_type = "t3.large"
  ami           = "ami-12345"
}

# ✅ Sugestão: Usar variáveis e data sources
variable "instance_type" {
  default = "t3.large"
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
}
```

### 2.3. Importar Recurso
**Objetivo**: Identificar recursos criados manualmente que deveriam estar no Terraform

**Detecção**:
```yaml
detected_resources:
  - type: aws_s3_bucket
    id: my-manual-bucket
    status: not_in_state
    suggestion: |
      terraform import aws_s3_bucket.manual my-manual-bucket
      
      # Depois adicione no código:
      resource "aws_s3_bucket" "manual" {
        bucket = "my-manual-bucket"
        # ... outras configurações
      }
```

### 2.4. Remover Dados Sensíveis
**Objetivo**: Detectar secrets/passwords/keys expostos

**Padrões Detectados**:

```hcl
# ❌ Password em plaintext
resource "aws_db_instance" "main" {
  password = "MySecretPassword123!"
}

# ✅ Sugestão: Usar AWS Secrets Manager
data "aws_secretsmanager_secret_version" "db_password" {
  secret_id = "prod/db/password"
}

resource "aws_db_instance" "main" {
  password = data.aws_secretsmanager_secret_version.db_password.secret_string
}
```

```hcl
# ❌ API Key exposta
variable "api_key" {
  default = "sk-abc123xyz789"
}

# ✅ Sugestão: Usar variável de ambiente
variable "api_key" {
  type      = string
  sensitive = true
  # Set via TF_VAR_api_key environment variable
}
```

### 2.5. Sugestão de Módulos Community
**Objetivo**: Recomendar módulos Terraform Registry quando apropriado

**Exemplo**:
```yaml
detected_pattern: "VPC com 3 subnets públicas e 3 privadas"
suggestion:
  module: "terraform-aws-modules/vpc/aws"
  version: "~> 5.0"
  benefits:
    - Mantido pela community
    - 500+ contribuidores
    - Best practices embutidas
    - Suporte a NAT Gateway, VPN, etc.
  example: |
    module "vpc" {
      source  = "terraform-aws-modules/vpc/aws"
      version = "~> 5.0"
      
      name = "my-vpc"
      cidr = "10.0.0.0/16"
      
      azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
      private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
      public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
      
      enable_nat_gateway = true
    }
```

### 2.6. Sugestão de Provider
**Objetivo**: Recomendar updates ou alternativas de provider

**Cenários**:

#### Provider Desatualizado
```yaml
current_version: "hashicorp/aws ~> 4.0"
latest_version: "hashicorp/aws ~> 5.0"
breaking_changes: false
suggestion: "Atualizar para ~> 5.0 para novos recursos e bugfixes"
```

#### Provider Deprecado
```yaml
current: "hashicorp/template"
status: deprecated
alternative: "hashicorp/terraform (templatefile function)"
migration_guide: "https://..."
```

---

## ✅ Categoria 3: Boas Práticas

### 3.1. Stack com Menos de 100 Recursos
**Validação**:
```yaml
current_resources: 120
threshold: 100
status: ❌ FAIL
recommendation: |
  Stack muito grande dificulta:
  - Debugging
  - Review de código
  - Apply performance
  - Blast radius de erros
  
  Sugestões de refatoração:
  1. Separar por ambiente (dev/staging/prod)
  2. Separar por componente (networking/compute/data)
  3. Usar workspace do Terraform
```

### 3.2. Stack com Tamanho Menor que 100 MB
**Validação**:
```yaml
stack_size: 150 MB
threshold: 100 MB
status: ❌ FAIL
issues:
  - "terraform.tfstate muito grande (80 MB)"
  - "Muitos data sources inline"
suggestion: |
  - Usar remote state
  - Remover data sources grandes
  - Considerar separar em múltiplos stacks
```

### 3.3. Stack com Documentação (README)
**Validação**:
```yaml
readme_exists: false
status: ❌ FAIL
suggestion: |
  Criar README.md com:
  
  # Infrastructure - [Component Name]
  
  ## Overview
  Brief description of what this stack creates
  
  ## Requirements
  - Terraform >= 1.5.0
  - AWS Provider >= 5.0
  
  ## Usage
  ```hcl
  terraform init
  terraform plan
  terraform apply
  ```
  
  ## Variables
  | Name | Description | Type | Default |
  |------|-------------|------|---------|
  | ... | ... | ... | ... |
  
  ## Outputs
  | Name | Description |
  |------|-------------|
  | ... | ... |
```

### 3.4. Stack com Testes de Unidade
**Validação**:
```yaml
tests_found: false
test_framework: null
status: ❌ FAIL
suggestion: |
  Implementar testes com Terratest:
  
  # tests/terraform_test.go
  package test
  
  import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
  )
  
  func TestTerraformVPC(t *testing.T) {
    opts := &terraform.Options{
      TerraformDir: "../",
    }
    
    defer terraform.Destroy(t, opts)
    terraform.InitAndApply(t, opts)
    
    vpcId := terraform.Output(t, opts, "vpc_id")
    assert.NotEmpty(t, vpcId)
  }
```

---

## 🤖 Integração com LLM

O agente usa LLM para:

### 1. Análise Contextual
- Entender o propósito da infraestrutura
- Identificar padrões arquiteturais
- Sugerir melhorias específicas ao contexto

### 2. Geração de Correções
- Criar código Terraform para corrigir problemas
- Gerar exemplos de uso
- Explicar o impacto de mudanças

### 3. Consulta à Knowledge Base
- Buscar best practices relevantes
- Encontrar módulos aplicáveis
- Verificar políticas da plataforma

### 4. Priorização Inteligente
- Ordenar sugestões por impacto
- Identificar quick wins
- Agrupar mudanças relacionadas

---

## 📊 Output Esperado

```json
{
  "id": "analysis-uuid",
  "timestamp": "2025-10-07T10:30:00Z",
  "score": 75,
  "score_breakdown": {
    "security": 70,
    "cost": 85,
    "best_practices": 60,
    "architecture": 80
  },
  "status": "PASSED",
  "summary": "Infraestrutura está em boa forma com algumas melhorias sugeridas",
  
  "bugs": [
    {
      "category": "preview_error",
      "severity": "high",
      "message": "Dependência circular detectada",
      "resource": "aws_instance.web",
      "fix": "terraform code...",
      "explanation": "LLM generated explanation..."
    }
  ],
  
  "improvements": [
    {
      "category": "cost_optimization",
      "impact": "medium",
      "savings": "$50/month",
      "message": "Usar t3.large ao invés de t2.large",
      "implementation": "terraform code..."
    }
  ],
  
  "best_practices": [
    {
      "check": "stack_size",
      "status": "warning",
      "current": 85,
      "threshold": 100,
      "message": "Stack approaching size limit"
    }
  ],
  
  "llm_insights": {
    "architectural_pattern": "3-tier web application",
    "recommendations": [
      "Consider using Auto Scaling Groups",
      "Add Application Load Balancer for HA"
    ],
    "knowledge_base_matches": [
      {
        "topic": "AWS Web Application Best Practices",
        "relevance": 0.95
      }
    ]
  }
}
```

---

## 🎯 Success Criteria

O agente será considerado bem-sucedido quando:

1. ✅ Detectar 95%+ dos problemas de preview
2. ✅ Classificar corretamente skipped checks (exception/window/test)
3. ✅ Identificar 80%+ das oportunidades de economia de custo
4. ✅ Sugerir módulos community relevantes em 90% dos casos
5. ✅ Detectar 99% dos dados sensíveis expostos
6. ✅ Validar todas as best practices definidas
7. ✅ Gerar correções aplicáveis em 70%+ dos casos
8. ✅ Usar LLM para enriquecer 100% das análises

---

## 📚 Referências

- [Terraform Best Practices](https://www.terraform-best-practices.com/)
- [Checkov Documentation](https://www.checkov.io/)
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)
- [Terraform Registry](https://registry.terraform.io/)
- [Nation.fun Platform Docs](https://nation.fun/) (placeholder)
