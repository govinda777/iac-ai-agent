# Implementation Roadmap - IaC AI Agent

**Última Atualização:** 2025-10-07  
**Versão Atual:** 1.0.0  
**Versão Alvo:** 2.0.0  
**Status:** Em implementação - Sprint 1

## 📊 Estado Atual vs. Objetivo

| Componente | Status Atual | Objetivo | Progresso |
|------------|--------------|----------|-----------|
| Terraform Analyzer | ✅ Implementado | ✅ Funcional | 100% |
| Checkov Integration | ✅ Implementado | ✅ Funcional | 100% |
| IAM Analyzer | ✅ Básico | ✅ Completo | 70% |
| Cost Optimizer | ⚠️ Básico | ✅ Avançado | 40% |
| Security Advisor | ✅ Implementado | ✅ Funcional | 80% |
| **LLM Integration** | ❌ Não usado | ✅ **CRÍTICO** | 0% |
| **Knowledge Base** | ❌ Não usado | ✅ **CRÍTICO** | 10% |
| Preview Analyzer | ❌ Ausente | ✅ Necessário | 0% |
| Drift Detection | ❌ Ausente | ✅ Necessário | 0% |
| Secrets Scanner | ❌ Ausente | ✅ Necessário | 0% |
| Module Suggester | ❌ Ausente | ✅ Desejável | 0% |
| Best Practices | ⚠️ Parcial | ✅ Completo | 30% |

**Score Geral de Implementação:** 24% ✅ → **Objetivo: 95%**

---

## 🚀 Sprint 1: Fundação AI (Semanas 1-2) - CRÍTICO

### Objetivo
Transformar o projeto em um verdadeiro **AI Agent** integrando LLM e Knowledge Base.

### Tasks

#### 1.1. Integrar LLM ao Fluxo de Análise ⭐ CRÍTICO
**Arquivo:** `internal/services/analysis.go`

```go
// ANTES (sem LLM)
func (as *AnalysisService) generateSuggestions(...) []models.Suggestion {
    suggestions := []models.Suggestion{}
    
    // Apenas regras hardcoded
    for _, warning := range tfAnalysis.BestPracticeWarnings {
        suggestions = append(suggestions, ...)
    }
    
    return suggestions
}

// DEPOIS (com LLM)
func (as *AnalysisService) generateSuggestions(...) []models.Suggestion {
    // 1. Sugestões baseadas em regras (rápido, determinístico)
    ruleBased := as.generateRuleBasedSuggestions(...)
    
    // 2. Enriquece com LLM (contexto, inteligência)
    enriched, err := as.enrichSuggestionsWithLLM(
        analysisDetails,
        ruleBased,
    )
    
    // 3. Combina e deduplicates
    return as.mergeSuggestions(ruleBased, enriched)
}
```

**Arquivo Novo:** `internal/services/llm_enrichment.go`

```go
package services

import (
    "github.com/gosouza/iac-ai-agent/internal/agent/llm"
    "github.com/gosouza/iac-ai-agent/internal/models"
)

// enrichSuggestionsWithLLM usa LLM para adicionar contexto inteligente
func (as *AnalysisService) enrichSuggestionsWithLLM(
    analysis *models.AnalysisDetails,
    baseSuggestions []models.Suggestion,
) ([]models.Suggestion, error) {
    
    // Consulta Knowledge Base
    relevantPractices := as.knowledgeBase.GetRelevantPractices(analysis)
    relevantModules := as.moduleRegistry.FindApplicableModules(analysis.Terraform.Resources)
    
    // Constrói prompt contextualizado
    prompt := as.promptBuilder.BuildEnrichmentPrompt(
        analysis,
        baseSuggestions,
        relevantPractices,
        relevantModules,
    )
    
    // Chama LLM
    llmResp, err := as.llmClient.Generate(&models.LLMRequest{
        Prompt:      prompt,
        Temperature: 0.2,
        MaxTokens:   2000,
    })
    if err != nil {
        as.logger.Warn("LLM enrichment failed, using base suggestions", "error", err)
        return baseSuggestions, nil // Fallback gracioso
    }
    
    // Parse resposta estruturada
    var enrichedSuggestions []models.EnrichedSuggestion
    if err := as.parseLLMSuggestions(llmResp.Content, &enrichedSuggestions); err != nil {
        return baseSuggestions, nil
    }
    
    // Converte para formato padrão
    return as.convertEnrichedSuggestions(enrichedSuggestions), nil
}

// parseLLMSuggestions extrai sugestões estruturadas da resposta LLM
func (as *AnalysisService) parseLLMSuggestions(
    content string,
    target interface{},
) error {
    // Extrai JSON da resposta (LLM pode adicionar texto antes/depois)
    jsonStr := extractJSONFromText(content)
    return json.Unmarshal([]byte(jsonStr), target)
}
```

**Prompt Template:** `internal/agent/llm/prompts/enrichment.go`

```go
const EnrichmentPromptTemplate = `
# Infrastructure as Code Analysis Enhancement

## Current Analysis

### Resources
{{range .Resources}}
- {{.Type}}.{{.Name}} ({{.File}}:{{.Line}})
{{end}}

### Existing Suggestions (Rule-Based)
{{range .BaseSuggestions}}
- [{{.Severity}}] {{.Message}}
  Recommendation: {{.Recommendation}}
{{end}}

### Knowledge Base Context

#### Relevant Best Practices
{{range .BestPractices}}
- {{.}}
{{end}}

#### Recommended Modules
{{range .Modules}}
- {{.Name}} ({{.Source}}) - {{.Description}}
{{end}}

## Your Task

As an Infrastructure as Code expert, analyze the above and provide:

1. **Enhanced Explanations**: For each existing suggestion, provide:
   - Why it matters (business impact)
   - How to implement (specific code example)
   - What could go wrong if ignored

2. **New Insights**: Identify additional improvements not caught by rules:
   - Architectural patterns that could be improved
   - Security considerations
   - Cost optimization opportunities
   - Module suggestions from the recommended list

3. **Prioritization**: Order suggestions by:
   - Impact (critical/high/medium/low)
   - Effort (easy/medium/hard)
   - ROI (quick wins first)

## Response Format

Respond ONLY with valid JSON:

{
  "enriched_suggestions": [
    {
      "original_id": "suggestion-uuid or null if new",
      "type": "security|cost|best_practice|architecture",
      "severity": "critical|high|medium|low",
      "title": "Brief title",
      "message": "Detailed explanation with business impact",
      "code_example": "# HCL code showing fix",
      "implementation_effort": "easy|medium|hard",
      "estimated_impact": "Description of impact",
      "why_it_matters": "Business justification",
      "references": ["https://...", "https://..."]
    }
  ],
  "architectural_insights": {
    "pattern_detected": "3-tier web app | microservices | ...",
    "strengths": ["..."],
    "areas_for_improvement": ["..."]
  },
  "priority_actions": [
    "Most important action to take first",
    "Second priority",
    "..."
  ]
}
`
```

**Estimativa:** 3-4 dias  
**Prioridade:** 🔥 CRÍTICA  
**Blocker para:** Tudo relacionado a IA

---

#### 1.2. Conectar Knowledge Base ⭐ CRÍTICO
**Arquivo:** `internal/platform/cloudcontroller/knowledge_base.go` (expandir)

```go
// ADICIONAR: Métodos de busca contextual
func (kb *KnowledgeBase) GetRelevantPractices(
    analysis *models.AnalysisDetails,
) []BestPractice {
    relevant := []BestPractice{}
    
    // Por tipo de recurso
    for _, resource := range analysis.Terraform.Resources {
        if practices, ok := kb.bestPractices[resource.Type]; ok {
            relevant = append(relevant, practices...)
        }
    }
    
    // Por provider
    for _, provider := range analysis.Terraform.Providers {
        if practices, ok := kb.providerPractices[provider]; ok {
            relevant = append(relevant, practices...)
        }
    }
    
    // Por padrão arquitetural detectado
    pattern := kb.detectArchitecturalPattern(analysis)
    if practices, ok := kb.architecturePractices[pattern]; ok {
        relevant = append(relevant, practices...)
    }
    
    return kb.deduplicate(relevant)
}

// ADICIONAR: Detecção de padrões
func (kb *KnowledgeBase) detectArchitecturalPattern(
    analysis *models.AnalysisDetails,
) string {
    resources := analysis.Terraform.Resources
    
    // 3-tier web app
    hasLB := containsResourceType(resources, "aws_lb", "aws_alb", "aws_elb")
    hasCompute := containsResourceType(resources, "aws_instance", "aws_ecs_service")
    hasDB := containsResourceType(resources, "aws_rds_instance", "aws_db_instance")
    if hasLB && hasCompute && hasDB {
        return "3-tier-web-app"
    }
    
    // Serverless
    hasLambda := containsResourceType(resources, "aws_lambda_function")
    hasAPIGW := containsResourceType(resources, "aws_api_gateway_rest_api")
    hasDynamoDB := containsResourceType(resources, "aws_dynamodb_table")
    if hasLambda && (hasAPIGW || hasDynamoDB) {
        return "serverless"
    }
    
    // Microservices
    hasECS := containsResourceType(resources, "aws_ecs_service")
    hasEKS := containsResourceType(resources, "aws_eks_cluster")
    if (hasECS || hasEKS) && len(resources) > 10 {
        return "microservices"
    }
    
    return "general"
}
```

**Novo Arquivo:** `internal/platform/cloudcontroller/knowledge_data.go`

```go
package cloudcontroller

// loadPlatformContext carrega contexto específico da plataforma Nation.fun
func (kb *KnowledgeBase) loadPlatformContext() {
    kb.platformContext = PlatformContext{
        Name: "Nation.fun",
        Standards: Standards{
            TaggingPolicy: map[string]string{
                "Environment":  "required (dev|staging|prod)",
                "Owner":        "required (email)",
                "CostCenter":   "required",
                "ManagedBy":    "required (terraform)",
                "Project":      "required",
            },
            NamingConvention: "{project}-{environment}-{resource}-{random}",
            RequiredOutputs: []string{
                "vpc_id",
                "resource_arns",
                "endpoints",
            },
        },
        SupportedVersions: SupportedVersions{
            Terraform: []string{">= 1.5.0"},
            OpenTofu:  []string{">= 1.6.0"},
            Providers: map[string]string{
                "aws":   ">= 5.0",
                "azure": ">= 3.0",
                "gcp":   ">= 4.0",
            },
        },
        ApprovedModules: []ApprovedModule{
            {
                Source:      "terraform-aws-modules/vpc/aws",
                Version:     "~> 5.0",
                UseCase:     "VPC creation",
                Recommended: true,
            },
            {
                Source:      "terraform-aws-modules/eks/aws",
                Version:     "~> 19.0",
                UseCase:     "EKS cluster",
                Recommended: true,
            },
            // ... mais módulos
        },
    }
}

// loadSecurityPolicies carrega políticas de segurança
func (kb *KnowledgeBase) loadSecurityPolicies() {
    kb.securityPolicies = []SecurityPolicy{
        {
            ID:          "SEC-001",
            Title:       "No Public S3 Buckets",
            Description: "S3 buckets must not be publicly accessible unless explicitly approved",
            Severity:    "critical",
            AutoFix:     true,
            FixCode: `
resource "aws_s3_bucket_public_access_block" "example" {
  bucket = aws_s3_bucket.example.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}`,
        },
        {
            ID:          "SEC-002",
            Title:       "Encryption at Rest Required",
            Description: "All storage resources must be encrypted at rest",
            Severity:    "high",
            Resources:   []string{"aws_s3_bucket", "aws_ebs_volume", "aws_rds_instance"},
        },
        // ... mais políticas
    }
}
```

**Estimativa:** 2-3 dias  
**Prioridade:** 🔥 CRÍTICA  
**Depende de:** Estrutura de dados definida

---

#### 1.3. Implementar Preview Analyzer ⭐ ALTA
**Novo Arquivo:** `internal/agent/analyzer/preview.go`

```go
package analyzer

import (
    "encoding/json"
    "fmt"
    
    "github.com/gosouza/iac-ai-agent/internal/models"
    "github.com/gosouza/iac-ai-agent/pkg/logger"
)

type PreviewAnalyzer struct {
    logger *logger.Logger
}

func NewPreviewAnalyzer(log *logger.Logger) *PreviewAnalyzer {
    return &PreviewAnalyzer{logger: log}
}

// AnalyzePreview analisa resultado de terraform plan em formato JSON
func (pa *PreviewAnalyzer) AnalyzePreview(planJSON []byte) (*models.PreviewAnalysis, error) {
    var plan models.TerraformPlan
    if err := json.Unmarshal(planJSON, &plan); err != nil {
        return nil, fmt.Errorf("invalid terraform plan JSON: %w", err)
    }
    
    analysis := &models.PreviewAnalysis{
        PlannedChanges:    []models.PlannedChange{},
        Errors:            []string{},
        Warnings:          []string{},
        ResourcesAffected: 0,
        RiskLevel:         "low",
    }
    
    // Analisa resource_changes
    for _, rc := range plan.ResourceChanges {
        change := pa.analyzeResourceChange(&rc)
        analysis.PlannedChanges = append(analysis.PlannedChanges, change)
        
        // Contabiliza ações
        for _, action := range rc.Change.Actions {
            switch action {
            case "create":
                analysis.CreateCount++
            case "update":
                analysis.UpdateCount++
            case "delete":
                analysis.DestroyCount++
            case "replace":
                analysis.ReplaceCount++
            }
        }
    }
    
    analysis.ResourcesAffected = len(plan.ResourceChanges)
    analysis.RiskLevel = pa.calculateRiskLevel(analysis)
    
    // Detecta mudanças arriscadas
    riskyChanges := pa.detectRiskyChanges(analysis.PlannedChanges)
    for _, risk := range riskyChanges {
        analysis.Warnings = append(analysis.Warnings, risk.Message)
    }
    
    return analysis, nil
}

// detectRiskyChanges identifica mudanças de alto risco
func (pa *PreviewAnalyzer) detectRiskyChanges(changes []models.PlannedChange) []models.RiskWarning {
    warnings := []models.RiskWarning{}
    
    for _, change := range changes {
        // Destruição de banco de dados
        if isDatabase(change.Resource) && change.Action == "destroy" {
            warnings = append(warnings, models.RiskWarning{
                Severity: "critical",
                Resource: change.Resource,
                Message:  "⚠️ Database will be DESTROYED. Ensure backups exist!",
                Action:   "destroy",
            })
        }
        
        // Replace de recursos stateful
        if isStateful(change.Resource) && change.Action == "replace" {
            warnings = append(warnings, models.RiskWarning{
                Severity: "high",
                Resource: change.Resource,
                Message:  "Resource replacement will cause downtime",
                Action:   "replace",
            })
        }
        
        // Mudança de VPC/Network
        if isNetworking(change.Resource) {
            warnings = append(warnings, models.RiskWarning{
                Severity: "high",
                Resource: change.Resource,
                Message:  "Network change may affect connectivity",
                Action:   change.Action,
            })
        }
    }
    
    return warnings
}

// calculateRiskLevel calcula nível de risco geral
func (pa *PreviewAnalyzer) calculateRiskLevel(analysis *models.PreviewAnalysis) string {
    // Destruições são sempre high risk
    if analysis.DestroyCount > 0 {
        return "high"
    }
    
    // Replace de muitos recursos é medium risk
    if analysis.ReplaceCount > 5 {
        return "medium"
    }
    
    // Muitas mudanças é medium risk
    if analysis.ResourcesAffected > 20 {
        return "medium"
    }
    
    return "low"
}
```

**Novo Modelo:** `internal/models/preview.go`

```go
package models

type PreviewAnalysis struct {
    PlannedChanges    []PlannedChange `json:"planned_changes"`
    Errors            []string        `json:"errors"`
    Warnings          []string        `json:"warnings"`
    ResourcesAffected int             `json:"resources_affected"`
    CreateCount       int             `json:"create_count"`
    UpdateCount       int             `json:"update_count"`
    DestroyCount      int             `json:"destroy_count"`
    ReplaceCount      int             `json:"replace_count"`
    RiskLevel         string          `json:"risk_level"` // low, medium, high, critical
}

type PlannedChange struct {
    Resource    string                 `json:"resource"`
    Action      string                 `json:"action"` // create, update, destroy, replace
    Changes     map[string]interface{} `json:"changes"`
    RiskScore   int                    `json:"risk_score"`
    Impact      string                 `json:"impact"`
}

type TerraformPlan struct {
    FormatVersion    string `json:"format_version"`
    TerraformVersion string `json:"terraform_version"`
    ResourceChanges  []struct {
        Address string `json:"address"`
        Mode    string `json:"mode"`
        Type    string `json:"type"`
        Name    string `json:"name"`
        Change  struct {
            Actions []string               `json:"actions"`
            Before  map[string]interface{} `json:"before"`
            After   map[string]interface{} `json:"after"`
        } `json:"change"`
    } `json:"resource_changes"`
}
```

**Estimativa:** 2-3 dias  
**Prioridade:** 🔥 ALTA  
**Blocker para:** Análise completa de preview

---

#### 1.4. Implementar Secrets Scanner ⭐ ALTA
**Novo Arquivo:** `internal/agent/analyzer/secrets.go`

```go
package analyzer

import (
    "regexp"
    "strings"
    
    "github.com/gosouza/iac-ai-agent/internal/models"
    "github.com/gosouza/iac-ai-agent/pkg/logger"
)

type SecretsAnalyzer struct {
    patterns []SecretPattern
    logger   *logger.Logger
}

type SecretPattern struct {
    Name        string
    Regex       *regexp.Regexp
    Severity    string
    Description string
    Suggestion  string
}

func NewSecretsAnalyzer(log *logger.Logger) *SecretsAnalyzer {
    sa := &SecretsAnalyzer{
        logger:   log,
        patterns: []SecretPattern{},
    }
    sa.loadPatterns()
    return sa
}

func (sa *SecretsAnalyzer) loadPatterns() {
    sa.patterns = []SecretPattern{
        {
            Name:        "AWS Access Key",
            Regex:       regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
            Severity:    "critical",
            Description: "AWS Access Key ID detected in plaintext",
            Suggestion:  "Use AWS Secrets Manager or environment variables",
        },
        {
            Name:        "AWS Secret Key",
            Regex:       regexp.MustCompile(`aws_secret_access_key\s*=\s*["']([^"']+)["']`),
            Severity:    "critical",
            Description: "AWS Secret Access Key in plaintext",
            Suggestion:  "Never commit AWS credentials. Use IAM roles or AWS SSO",
        },
        {
            Name:        "Generic Password",
            Regex:       regexp.MustCompile(`password\s*=\s*["']([^"']+)["']`),
            Severity:    "high",
            Description: "Password in plaintext",
            Suggestion:  "Use variable with sensitive=true or secrets manager",
        },
        {
            Name:        "Private Key",
            Regex:       regexp.MustCompile(`-----BEGIN.*PRIVATE KEY-----`),
            Severity:    "critical",
            Description: "Private key detected in code",
            Suggestion:  "Store private keys in secure vault",
        },
        {
            Name:        "API Key",
            Regex:       regexp.MustCompile(`api[_-]?key\s*=\s*["']([A-Za-z0-9]{20,})["']`),
            Severity:    "high",
            Description: "API key in plaintext",
            Suggestion:  "Use environment variables or parameter store",
        },
    }
}

// ScanContent escaneia conteúdo em busca de secrets
func (sa *SecretsAnalyzer) ScanContent(
    content string,
    filename string,
) []models.SecretFinding {
    findings := []models.SecretFinding{}
    lines := strings.Split(content, "\n")
    
    for lineNum, line := range lines {
        for _, pattern := range sa.patterns {
            if pattern.Regex.MatchString(line) {
                findings = append(findings, models.SecretFinding{
                    Type:        pattern.Name,
                    File:        filename,
                    Line:        lineNum + 1,
                    Value:       sa.maskSecret(line),
                    Severity:    pattern.Severity,
                    Description: pattern.Description,
                    Suggestion:  pattern.Suggestion,
                })
            }
        }
    }
    
    return findings
}

// maskSecret mascara o valor do secret
func (sa *SecretsAnalyzer) maskSecret(line string) string {
    if len(line) > 50 {
        return line[:20] + "***REDACTED***" + line[len(line)-10:]
    }
    return "***REDACTED***"
}
```

**Estimativa:** 2 dias  
**Prioridade:** 🔥 ALTA  
**Impacto:** Segurança crítica

---

#### 1.5. Testes e Documentação
- Criar testes unitários para novos componentes
- Atualizar documentação técnica
- Criar exemplos de uso

**Estimativa:** 2 dias  
**Prioridade:** MÉDIA

### Entregáveis Sprint 1
- ✅ LLM totalmente integrado ao fluxo
- ✅ Knowledge Base consultada em todas análises
- ✅ Preview Analyzer funcional
- ✅ Secrets Scanner operacional
- ✅ Testes passando (80%+ coverage dos novos códigos)
- ✅ Documentação atualizada

**Duração Total:** 2 semanas  
**Recursos:** 1 dev full-time

---

## 🔨 Sprint 2: Features Core (Semanas 3-4)

### Objetivo
Implementar funcionalidades essenciais que ainda faltam.

### Tasks

#### 2.1. Drift Detection
**Arquivo:** `internal/agent/analyzer/drift.go`

```go
// CompareTerraformStateWithCode compara state com código
func (da *DriftAnalyzer) CompareTerraformStateWithCode(
    stateJSON []byte,
    tfAnalysis *models.TerraformAnalysis,
) (*models.DriftAnalysis, error)
```

**Estimativa:** 3 dias

#### 2.2. Best Practices Validator Completo
**Arquivo:** `internal/agent/analyzer/best_practices.go`

```go
// ValidateStackSize verifica tamanho de stack
// ValidateDocumentation verifica README
// ValidateTests verifica testes
// ValidateFileSize verifica tamanho de arquivos
```

**Estimativa:** 2 dias

#### 2.3. Module Suggester
**Arquivo:** `internal/agent/suggester/module_suggester.go`

Integrar `ModuleRegistry` para sugerir módulos community.

**Estimativa:** 2 dias

#### 2.4. Operational Analyzer (Timeout/Stuck Resources)
**Arquivo:** `internal/agent/analyzer/operational.go`

Analisar histórico de apply/destroy para detectar problemas.

**Estimativa:** 3 dias

### Entregáveis Sprint 2
- ✅ Drift detection funcional
- ✅ Best practices 100% validadas
- ✅ Sugestões de módulos community
- ✅ Detecção de problemas operacionais

**Duração:** 2 semanas

---

## 🎯 Sprint 3: Features Avançadas (Semanas 5-6)

### Objetivo
Implementar features avançadas de refatoração e importação.

### Tasks

#### 3.1. Architecture Advisor
**Arquivo:** `internal/agent/suggester/architecture_advisor.go`

Análise de padrões arquiteturais e sugestões de refatoração.

**Estimativa:** 4 dias

#### 3.2. Resource Import Suggester
**Arquivo:** `internal/agent/suggester/import_suggester.go`

Detectar recursos não gerenciados e sugerir importação.

**Estimativa:** 3 dias

#### 3.3. Provider Update Advisor
**Arquivo:** `internal/agent/suggester/provider_advisor.go`

Sugerir updates de providers e alternativos.

**Estimativa:** 2 dias

### Entregáveis Sprint 3
- ✅ Refatoração arquitetural
- ✅ Sugestões de importação de recursos
- ✅ Advisor de providers

**Duração:** 2 semanas

---

## 📚 Sprint 4: Documentação e Polish (Semana 7)

### Objetivo
Finalizar documentação e polish.

### Tasks
- Documentação completa de todas features
- Guias de integração (Spacelift, Terraform Cloud, CI/CD)
- Exemplos práticos
- Melhorias de UI/UX nas respostas
- Otimizações de performance

**Duração:** 1 semana

---

## 📊 Métricas de Sucesso

| Métrica | Atual | Objetivo Sprint 1 | Objetivo Sprint 2 | Objetivo Sprint 3 | Objetivo Sprint 4 |
|---------|-------|-------------------|-------------------|-------------------|-------------------|
| Features Implementadas | 24% | 50% | 75% | 90% | 95% |
| LLM Integration | 0% | 100% | 100% | 100% | 100% |
| Knowledge Base Usage | 10% | 100% | 100% | 100% | 100% |
| Test Coverage | 60% | 75% | 80% | 85% | 90% |
| Documentation | 40% | 60% | 75% | 85% | 100% |

---

## 🚨 Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| LLM API custos altos | Média | Alto | Implementar cache agressivo, rate limiting |
| LLM latência | Alta | Médio | Fallback para rule-based, async processing |
| Complexidade de drift | Alta | Médio | Começar com casos simples, expandir gradualmente |
| Knowledge Base desatualizada | Média | Médio | Processo de review mensal, versionamento |
| Overhead de testes | Média | Baixo | Priorizar testes de alto valor, mocks para LLM |

---

## 👥 Recursos Necessários

### Sprint 1
- 1x Senior Backend Dev (Go)
- 1x ML Engineer (part-time, prompts)

### Sprint 2-3
- 1x Senior Backend Dev (Go)
- 1x DevOps Engineer (part-time, integrations)

### Sprint 4
- 1x Technical Writer (documentação)
- 1x QA Engineer (testes end-to-end)

---

## ✅ Definition of Done

Cada sprint será considerado concluído quando:

1. ✅ Todos os testes passando (>= 80% coverage)
2. ✅ Code review aprovado
3. ✅ Documentação atualizada
4. ✅ Exemplos de uso criados
5. ✅ Demo funcional
6. ✅ Performance dentro de SLA (<5s para análises, <10s com LLM)
7. ✅ Sem linter errors
8. ✅ Security scan passing

---

## 📅 Timeline Visual

```
Semanas:  1     2     3     4     5     6     7
Sprint:  |--1--|--2--|--3--|--4--|
         
Sprint 1: [████████████] LLM + KB + Preview + Secrets
Sprint 2:             [████████] Drift + BP + Modules + Ops
Sprint 3:                     [████████] Arch + Import + Provider
Sprint 4:                             [████] Docs + Polish

Releases:
v1.5.0: Sprint 1 complete ←
v1.8.0: Sprint 2 complete
v2.0.0: Sprint 3 complete ← OBJETIVO ALCANÇADO
v2.1.0: Sprint 4 complete (polish)
```

---

## 🎉 Marcos

- **Milestone 1 (Semana 2)**: Primeiro AI Agent funcional com LLM
- **Milestone 2 (Semana 4)**: Features core completas
- **Milestone 3 (Semana 6)**: Objetivo 95% alcançado
- **Milestone 4 (Semana 7)**: Produção ready

---

**Status:** 📋 Planejamento completo  
**Próximo Passo:** 🚀 Iniciar Sprint 1, Task 1.1 (LLM Integration)
