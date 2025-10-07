package llm

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// PromptBuilder constrói prompts para o LLM
type PromptBuilder struct {
	logger    *logger.Logger
	templates map[string]*template.Template
}

// NewPromptBuilder cria um novo construtor de prompts
func NewPromptBuilder(log *logger.Logger) *PromptBuilder {
	pb := &PromptBuilder{
		logger:    log,
		templates: make(map[string]*template.Template),
	}

	// Carrega templates
	pb.loadTemplates()

	return pb
}

// loadTemplates carrega todos os templates de prompt
func (pb *PromptBuilder) loadTemplates() {
	// Template de enriquecimento
	enrichmentTmpl, err := template.New("enrichment").Parse(EnrichmentPromptTemplate)
	if err != nil {
		pb.logger.Error("Erro ao carregar template de enriquecimento", "error", err)
	} else {
		pb.templates["enrichment"] = enrichmentTmpl
	}

	// Template de análise de preview
	previewTmpl, err := template.New("preview").Parse(PreviewAnalysisPromptTemplate)
	if err != nil {
		pb.logger.Error("Erro ao carregar template de análise de preview", "error", err)
	} else {
		pb.templates["preview"] = previewTmpl
	}

	// Template de análise de segurança
	securityTmpl, err := template.New("security").Parse(SecurityAnalysisPromptTemplate)
	if err != nil {
		pb.logger.Error("Erro ao carregar template de análise de segurança", "error", err)
	} else {
		pb.templates["security"] = securityTmpl
	}

	// Template de otimização de custo
	costTmpl, err := template.New("cost").Parse(CostOptimizationPromptTemplate)
	if err != nil {
		pb.logger.Error("Erro ao carregar template de otimização de custo", "error", err)
	} else {
		pb.templates["cost"] = costTmpl
	}
}

// BuildEnrichmentPrompt constrói prompt para enriquecimento de sugestões
func (pb *PromptBuilder) BuildEnrichmentPrompt(
	analysis *models.AnalysisDetails,
	baseSuggestions []models.Suggestion,
	relevantPractices []models.BestPractice,
	relevantModules []models.Module,
) string {
	// Verifica se o template existe
	tmpl, ok := pb.templates["enrichment"]
	if !ok {
		pb.logger.Error("Template de enriquecimento não encontrado")
		return buildFallbackEnrichmentPrompt(analysis, baseSuggestions)
	}

	// Prepara dados para o template
	resourcesWithLine := []map[string]interface{}{}

	for _, res := range analysis.Terraform.Resources {
		resourcesWithLine = append(resourcesWithLine, map[string]interface{}{
			"Type": res.Type,
			"Name": res.Name,
			"File": res.File,
			"Line": res.LineStart,
		})
	}

	data := map[string]interface{}{
		"Resources":       resourcesWithLine,
		"BaseSuggestions": baseSuggestions,
		"BestPractices":   relevantPractices,
		"Modules":         relevantModules,
	}

	// Executa o template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		pb.logger.Error("Erro ao executar template de enriquecimento", "error", err)
		return buildFallbackEnrichmentPrompt(analysis, baseSuggestions)
	}

	return buf.String()
}

// BuildPreviewAnalysisPrompt constrói prompt para análise de preview
func (pb *PromptBuilder) BuildPreviewAnalysisPrompt(
	previewAnalysis *models.PreviewAnalysis,
) string {
	// Verifica se o template existe
	tmpl, ok := pb.templates["preview"]
	if !ok {
		pb.logger.Error("Template de análise de preview não encontrado")
		return buildFallbackPreviewPrompt(previewAnalysis)
	}

	// Executa o template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, previewAnalysis); err != nil {
		pb.logger.Error("Erro ao executar template de análise de preview", "error", err)
		return buildFallbackPreviewPrompt(previewAnalysis)
	}

	return buf.String()
}

// BuildSecurityAnalysisPrompt constrói prompt para análise de segurança
func (pb *PromptBuilder) BuildSecurityAnalysisPrompt(
	analysis *models.AnalysisDetails,
	securityFindings []models.SecurityFinding,
) string {
	// Verifica se o template existe
	tmpl, ok := pb.templates["security"]
	if !ok {
		pb.logger.Error("Template de análise de segurança não encontrado")
		return buildFallbackSecurityPrompt(analysis, securityFindings)
	}

	// Prepara dados para o template
	resourcesWithLine := []map[string]interface{}{}

	for _, res := range analysis.Terraform.Resources {
		resourcesWithLine = append(resourcesWithLine, map[string]interface{}{
			"Type": res.Type,
			"Name": res.Name,
			"File": res.File,
			"Line": res.LineStart,
		})
	}

	securityFindingsWithType := []map[string]interface{}{}

	for _, finding := range securityFindings {
		securityFindingsWithType = append(securityFindingsWithType, map[string]interface{}{
			"Type":        finding.CheckName,
			"Severity":    finding.Severity,
			"Description": finding.Description,
			"Resource":    finding.Resource,
			"File":        finding.File,
			"Line":        finding.Line,
		})
	}

	data := map[string]interface{}{
		"Resources":        resourcesWithLine,
		"SecurityFindings": securityFindingsWithType,
	}

	// Executa o template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		pb.logger.Error("Erro ao executar template de análise de segurança", "error", err)
		return buildFallbackSecurityPrompt(analysis, securityFindings)
	}

	return buf.String()
}

// BuildCostOptimizationPrompt constrói prompt para otimização de custo
func (pb *PromptBuilder) BuildCostOptimizationPrompt(
	analysis *models.AnalysisDetails,
) string {
	// Verifica se o template existe
	tmpl, ok := pb.templates["cost"]
	if !ok {
		pb.logger.Error("Template de otimização de custo não encontrado")
		return buildFallbackCostPrompt(analysis)
	}

	// Executa o template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, analysis); err != nil {
		pb.logger.Error("Erro ao executar template de otimização de custo", "error", err)
		return buildFallbackCostPrompt(analysis)
	}

	return buf.String()
}

// Fallback prompts para caso de erro nos templates

func buildFallbackEnrichmentPrompt(
	analysis *models.AnalysisDetails,
	baseSuggestions []models.Suggestion,
) string {
	var sb strings.Builder

	sb.WriteString("# Infrastructure as Code Analysis Enhancement\n\n")
	sb.WriteString("## Current Analysis\n\n")

	// Resources
	sb.WriteString("### Resources\n")
	for _, res := range analysis.Terraform.Resources {
		sb.WriteString(fmt.Sprintf("- %s.%s (%s:%d)\n", res.Type, res.Name, res.File, res.LineStart))
	}
	sb.WriteString("\n")

	// Base suggestions
	sb.WriteString("### Existing Suggestions (Rule-Based)\n")
	for _, sugg := range baseSuggestions {
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", sugg.Severity, sugg.Message))
		sb.WriteString(fmt.Sprintf("  Recommendation: %s\n", sugg.Recommendation))
	}
	sb.WriteString("\n")

	// Task
	sb.WriteString("## Your Task\n\n")
	sb.WriteString("As an Infrastructure as Code expert, analyze the above and provide:\n\n")
	sb.WriteString("1. **Enhanced Explanations**: For each existing suggestion\n")
	sb.WriteString("2. **New Insights**: Identify additional improvements\n")
	sb.WriteString("3. **Prioritization**: Order suggestions by impact\n\n")

	// Response format
	sb.WriteString("## Response Format\n\n")
	sb.WriteString("Respond ONLY with valid JSON containing enriched_suggestions array.\n")

	return sb.String()
}

func buildFallbackPreviewPrompt(previewAnalysis *models.PreviewAnalysis) string {
	var sb strings.Builder

	sb.WriteString("# Terraform Plan Analysis\n\n")
	sb.WriteString(fmt.Sprintf("Analyze the following Terraform plan with %d resources affected.\n\n", previewAnalysis.ResourcesAffected))

	// Changes summary
	sb.WriteString("## Changes Summary\n")
	sb.WriteString(fmt.Sprintf("- Create: %d\n", previewAnalysis.CreateCount))
	sb.WriteString(fmt.Sprintf("- Update: %d\n", previewAnalysis.UpdateCount))
	sb.WriteString(fmt.Sprintf("- Replace: %d\n", previewAnalysis.ReplaceCount))
	sb.WriteString(fmt.Sprintf("- Destroy: %d\n", previewAnalysis.DestroyCount))
	sb.WriteString("\n")

	// Planned changes
	sb.WriteString("## Planned Changes\n")
	for _, change := range previewAnalysis.PlannedChanges {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", change.Action, change.Resource))
	}
	sb.WriteString("\n")

	// Task
	sb.WriteString("## Your Task\n\n")
	sb.WriteString("Analyze this plan and provide:\n\n")
	sb.WriteString("1. Risk assessment\n")
	sb.WriteString("2. Potential issues\n")
	sb.WriteString("3. Recommendations\n\n")

	// Response format
	sb.WriteString("## Response Format\n\n")
	sb.WriteString("Respond ONLY with valid JSON.\n")

	return sb.String()
}

func buildFallbackSecurityPrompt(
	analysis *models.AnalysisDetails,
	securityFindings []models.SecurityFinding,
) string {
	var sb strings.Builder

	sb.WriteString("# Security Analysis for Terraform Code\n\n")

	// Resources
	sb.WriteString("## Resources\n")
	for _, res := range analysis.Terraform.Resources {
		sb.WriteString(fmt.Sprintf("- %s.%s\n", res.Type, res.Name))
	}
	sb.WriteString("\n")

	// Security findings
	sb.WriteString("## Security Findings\n")
	for _, finding := range securityFindings {
		sb.WriteString(fmt.Sprintf("- [%s] %s: %s\n", finding.Severity, finding.CheckName, finding.Description))
	}
	sb.WriteString("\n")

	// Task
	sb.WriteString("## Your Task\n\n")
	sb.WriteString("Analyze these security findings and provide:\n\n")
	sb.WriteString("1. Detailed explanation of each issue\n")
	sb.WriteString("2. Remediation steps with code examples\n")
	sb.WriteString("3. Risk assessment\n\n")

	// Response format
	sb.WriteString("## Response Format\n\n")
	sb.WriteString("Respond ONLY with valid JSON.\n")

	return sb.String()
}

func buildFallbackCostPrompt(analysis *models.AnalysisDetails) string {
	var sb strings.Builder

	sb.WriteString("# Cost Optimization Analysis for Terraform Code\n\n")

	// Resources
	sb.WriteString("## Resources\n")
	for _, res := range analysis.Terraform.Resources {
		sb.WriteString(fmt.Sprintf("- %s.%s\n", res.Type, res.Name))
	}
	sb.WriteString("\n")

	// Task
	sb.WriteString("## Your Task\n\n")
	sb.WriteString("Analyze these resources for cost optimization and provide:\n\n")
	sb.WriteString("1. Cost optimization opportunities\n")
	sb.WriteString("2. Estimated savings\n")
	sb.WriteString("3. Implementation recommendations\n\n")

	// Response format
	sb.WriteString("## Response Format\n\n")
	sb.WriteString("Respond ONLY with valid JSON.\n")

	return sb.String()
}

// Templates de prompt

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
- {{.Title}}: {{.Description}}
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

const PreviewAnalysisPromptTemplate = `
# Terraform Plan Analysis

## Plan Summary

- **Resources Affected**: {{.ResourcesAffected}}
- **Create**: {{.CreateCount}}
- **Update**: {{.UpdateCount}}
- **Replace**: {{.ReplaceCount}}
- **Destroy**: {{.DestroyCount}}
- **Risk Level**: {{.RiskLevel}}

## Planned Changes
{{range .PlannedChanges}}
- [{{.Action}}] {{.Resource}}
{{end}}

## Warnings
{{range .Warnings}}
- {{.}}
{{end}}

## Your Task

As an Infrastructure as Code expert, analyze this Terraform plan and provide:

1. **Risk Assessment**:
   - Evaluate the overall risk of this plan
   - Identify high-risk changes (data loss, downtime, etc)
   - Assess potential blast radius

2. **Impact Analysis**:
   - Estimate downtime for services
   - Identify dependencies that might be affected
   - Highlight potential performance impacts

3. **Recommendations**:
   - Suggest safer ways to apply these changes
   - Recommend testing strategies
   - Propose rollback plan

4. **Best Practices**:
   - Identify any best practices violations in the plan
   - Suggest improvements for future changes

## Response Format

Respond ONLY with valid JSON:

{
  "risk_assessment": {
    "overall_risk": "low|medium|high|critical",
    "high_risk_changes": [
      {
        "resource": "resource_name",
        "action": "create|update|replace|destroy",
        "risk_details": "Description of risk",
        "mitigation": "How to mitigate"
      }
    ],
    "blast_radius": "Description of potential impact"
  },
  "impact_analysis": {
    "estimated_downtime": "None|Minimal|Significant|Extended",
    "affected_dependencies": ["service1", "service2"],
    "performance_impact": "Description of performance impact"
  },
  "recommendations": {
    "apply_strategy": "Recommended approach to apply",
    "testing_strategy": "Recommended testing approach",
    "rollback_plan": "Steps to rollback if needed"
  },
  "best_practices": [
    {
      "issue": "Description of issue",
      "recommendation": "How to improve"
    }
  ]
}
`

const SecurityAnalysisPromptTemplate = `
# Security Analysis for Terraform Code

## Resources
{{range .Resources}}
- {{.Type}}.{{.Name}} ({{.File}}:{{.Line}})
{{end}}

## Security Findings
{{range .SecurityFindings}}
- [{{.Severity}}] {{.Type}}: {{.Description}}
  Resource: {{.Resource}}
  File: {{.File}}:{{.Line}}
{{end}}

## Your Task

As a Security Expert for Infrastructure as Code, analyze these findings and provide:

1. **Detailed Analysis**:
   - Explain each security issue in detail
   - Assess the real-world risk and potential exploit scenarios
   - Identify any false positives

2. **Remediation**:
   - Provide specific code fixes for each issue
   - Explain the security principles behind each fix
   - Suggest additional security controls where appropriate

3. **Prioritization**:
   - Rank issues by severity and exploitability
   - Identify quick wins vs. complex fixes
   - Suggest implementation order

4. **Security Posture Improvement**:
   - Recommend additional security controls
   - Suggest security testing strategies
   - Propose security best practices

## Response Format

Respond ONLY with valid JSON:

{
  "security_analysis": {
    "critical_findings": [
      {
        "finding_id": "finding reference",
        "title": "Concise title",
        "description": "Detailed explanation",
        "exploit_scenario": "How this could be exploited",
        "is_false_positive": false,
        "remediation_code": "# Code fix example",
        "remediation_explanation": "Why this fixes the issue",
        "additional_controls": ["control1", "control2"]
      }
    ],
    "high_findings": [...],
    "medium_findings": [...],
    "low_findings": [...]
  },
  "implementation_plan": {
    "immediate_actions": ["action1", "action2"],
    "short_term_actions": ["action1", "action2"],
    "long_term_actions": ["action1", "action2"]
  },
  "security_posture_recommendations": [
    {
      "area": "Area of improvement",
      "recommendation": "Detailed recommendation",
      "implementation": "How to implement"
    }
  ]
}
`

const CostOptimizationPromptTemplate = `
# Cost Optimization Analysis for Terraform Code

## Resources
{{range .Terraform.Resources}}
- {{.Type}}.{{.Name}} ({{.File}}:{{.Line}})
{{end}}

## Your Task

As a Cloud Cost Optimization Expert, analyze these resources and provide:

1. **Cost Optimization Opportunities**:
   - Identify resources that could be optimized
   - Suggest right-sizing opportunities
   - Identify unused or underutilized resources
   - Recommend reserved instances or savings plans
   - Suggest architecture changes for cost efficiency

2. **Estimated Savings**:
   - Provide estimated monthly savings for each recommendation
   - Calculate ROI for implementation effort
   - Prioritize by savings potential

3. **Implementation Recommendations**:
   - Provide specific code changes
   - Suggest implementation approach
   - Identify any risks or trade-offs

## Response Format

Respond ONLY with valid JSON:

{
  "cost_optimization": {
    "total_estimated_savings": "Estimated monthly savings",
    "opportunities": [
      {
        "resource": "resource_name",
        "current_configuration": "Current setup",
        "recommendation": "Recommended change",
        "estimated_savings": "Monthly savings estimate",
        "implementation_effort": "low|medium|high",
        "roi": "high|medium|low",
        "code_example": "# Code example",
        "risks": ["risk1", "risk2"]
      }
    ]
  },
  "architecture_recommendations": [
    {
      "area": "Area of improvement",
      "current_approach": "Current architecture",
      "recommended_approach": "Recommended architecture",
      "estimated_savings": "Monthly savings estimate",
      "implementation_complexity": "low|medium|high"
    }
  ],
  "quick_wins": [
    {
      "action": "Action to take",
      "savings": "Savings estimate",
      "effort": "Implementation effort"
    }
  ]
}
`
