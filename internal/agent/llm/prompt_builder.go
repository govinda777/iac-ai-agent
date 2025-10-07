package llm

// NOTE: This file contains comments and strings in Portuguese, which may be flagged by spell checkers.
// This is intentional and not a bug.

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// PromptBuilder é responsável por construir prompts estruturados para o LLM
type PromptBuilder struct {
	logger *logger.Logger
}

// NewPromptBuilder cria um novo builder de prompts
func NewPromptBuilder(log *logger.Logger) *PromptBuilder {
	return &PromptBuilder{
		logger: log,
	}
}

// PromptData contém os dados para construir um prompt
type PromptData struct {
	TerraformCode  string
	CheckovResults *models.CheckovResult
	IAMPolicies    []string
	KnowledgeBase  map[string]interface{}
	Context        map[string]interface{}
}

// BuildAnalysisPrompt constrói um prompt para análise de código
func (b *PromptBuilder) BuildAnalysisPrompt(data *PromptData) (*models.LLMRequest, error) {
	// Inicializa template
	tmpl, err := template.New("analysis").Parse(analysisPromptTemplate)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear template: %w", err)
	}

	// Executa template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar template: %w", err)
	}

	// Constrói request
	req := &models.LLMRequest{
		SystemPrompt: systemPrompt,
		Prompt:       buf.String(),
		Temperature:  0.2,
		MaxTokens:    4000,
		ResponseFormat: "json",
	}

	b.logger.Debug("Prompt construído", 
		"prompt_length", len(req.Prompt),
		"has_terraform", len(data.TerraformCode) > 0,
		"has_checkov", data.CheckovResults != nil,
		"has_iam", len(data.IAMPolicies) > 0)

	return req, nil
}

// BuildSecurityPrompt constrói um prompt focado em segurança
func (b *PromptBuilder) BuildSecurityPrompt(data *PromptData) (*models.LLMRequest, error) {
	// Inicializa template
	tmpl, err := template.New("security").Parse(securityPromptTemplate)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear template: %w", err)
	}

	// Executa template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar template: %w", err)
	}

	// Constrói request
	req := &models.LLMRequest{
		SystemPrompt: securitySystemPrompt,
		Prompt:       buf.String(),
		Temperature:  0.1, // Menor temperatura para análise de segurança
		MaxTokens:    4000,
		ResponseFormat: "json",
	}

	return req, nil
}

// BuildCostPrompt constrói um prompt focado em otimização de custos
func (b *PromptBuilder) BuildCostPrompt(data *PromptData) (*models.LLMRequest, error) {
	// Inicializa template
	tmpl, err := template.New("cost").Parse(costPromptTemplate)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear template: %w", err)
	}

	// Executa template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar template: %w", err)
	}

	// Constrói request
	req := &models.LLMRequest{
		SystemPrompt: costSystemPrompt,
		Prompt:       buf.String(),
		Temperature:  0.2,
		MaxTokens:    4000,
		ResponseFormat: "json",
	}

	return req, nil
}

// BuildArchitecturePrompt constrói um prompt focado em arquitetura
func (b *PromptBuilder) BuildArchitecturePrompt(data *PromptData) (*models.LLMRequest, error) {
	// Inicializa template
	tmpl, err := template.New("architecture").Parse(architecturePromptTemplate)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear template: %w", err)
	}

	// Executa template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar template: %w", err)
	}

	// Constrói request
	req := &models.LLMRequest{
		SystemPrompt: architectureSystemPrompt,
		Prompt:       buf.String(),
		Temperature:  0.3, // Maior temperatura para sugestões de arquitetura
		MaxTokens:    4000,
		ResponseFormat: "json",
	}

	return req, nil
}

// FormatCheckovResults formata os resultados do Checkov para o prompt
func (b *PromptBuilder) FormatCheckovResults(results *models.CheckovResult) string {
	if results == nil || len(results.Results.FailedChecks) == 0 {
		return "Nenhum resultado Checkov disponível."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Resultados Checkov\n\n"))
	sb.WriteString(fmt.Sprintf("- Total de verificações: %d\n", results.Summary.Passed+results.Summary.Failed))
	sb.WriteString(fmt.Sprintf("- Verificações passaram: %d\n", results.Summary.Passed))
	sb.WriteString(fmt.Sprintf("- Verificações falharam: %d\n\n", results.Summary.Failed))

	sb.WriteString("### Falhas de Segurança\n\n")
	for i, check := range results.Results.FailedChecks {
		if i >= 20 { // Limita para não exceder tokens
			sb.WriteString(fmt.Sprintf("\n... mais %d falhas omitidas ...\n", len(results.Results.FailedChecks)-20))
			break
		}

		sb.WriteString(fmt.Sprintf("- **%s**\n", check.CheckID))
		sb.WriteString(fmt.Sprintf("  - **Severidade**: %s\n", check.Severity))
		sb.WriteString(fmt.Sprintf("  - **Arquivo**: %s\n", check.File))
		sb.WriteString(fmt.Sprintf("  - **Recurso**: %s\n", check.Resource))
		sb.WriteString(fmt.Sprintf("  - **Descrição**: %s\n", check.CheckName))
		if check.Guideline != "" {
			sb.WriteString(fmt.Sprintf("  - **Recomendação**: %s\n", check.Guideline))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Templates de prompts

const systemPrompt = `Você é um especialista em Infrastructure as Code (IaC), segurança e otimização de cloud.

Sua tarefa é analisar código Terraform e fornecer insights detalhados sobre:
1. Segurança e compliance
2. Otimização de custos
3. Boas práticas
4. Arquitetura e design

Forneça análises detalhadas e recomendações práticas, sempre incluindo exemplos de código quando relevante.
Suas respostas devem ser estruturadas, claras e acionáveis.`

const securitySystemPrompt = `Você é um especialista em segurança de cloud e Infrastructure as Code.

Sua tarefa é realizar uma análise profunda de segurança em código Terraform, identificando:
1. Vulnerabilidades críticas
2. Problemas de compliance
3. Configurações inseguras
4. Exposição de dados sensíveis
5. Problemas de IAM e permissões

Seja minucioso e detalhado. Priorize os problemas por severidade (Crítico, Alto, Médio, Baixo).
Forneça recomendações específicas para correção com exemplos de código.`

const costSystemPrompt = `Você é um especialista em otimização de custos para infraestrutura cloud.

Sua tarefa é analisar código Terraform e identificar oportunidades de economia:
1. Recursos superdimensionados
2. Recursos subutilizados
3. Opções de compra mais econômicas (Reserved Instances, Savings Plans)
4. Arquiteturas mais eficientes em custo
5. Recursos desnecessários ou redundantes

Quantifique as economias potenciais quando possível. Forneça recomendações específicas com exemplos de código.`

const architectureSystemPrompt = `Você é um arquiteto de soluções cloud especializado em Infrastructure as Code.

Sua tarefa é analisar código Terraform e fornecer insights arquiteturais:
1. Padrões de arquitetura identificados
2. Sugestões de melhoria arquitetural
3. Escalabilidade e resiliência
4. Modularização e reutilização
5. Integração com serviços gerenciados

Forneça recomendações específicas com exemplos de código e diagramas quando relevante.`

const analysisPromptTemplate = `# Análise de Infrastructure as Code

## Código Terraform

` + "```hcl" + `
{{.TerraformCode}}
` + "```" + `

{{if .CheckovResults}}
{{.CheckovResults}}
{{end}}

{{if .IAMPolicies}}
## Políticas IAM

{{range .IAMPolicies}}
` + "```json" + `
{{.}}
` + "```" + `
{{end}}
{{end}}

## Instruções

Analise o código Terraform fornecido e forneça:

1. **Resumo Executivo**: Visão geral da infraestrutura e principais pontos de atenção.

2. **Problemas Críticos**: Identifique vulnerabilidades de segurança, configurações incorretas ou riscos significativos.

3. **Recomendações Prioritárias**: Sugestões de melhorias mais importantes, com exemplos de código.

4. **Otimizações de Custo**: Oportunidades para reduzir custos sem comprometer funcionalidade.

5. **Boas Práticas**: Recomendações para melhorar o código seguindo as melhores práticas de IaC.

6. **Insights Arquiteturais**: Observações sobre a arquitetura e sugestões de design.

Forneça sua resposta como um JSON estruturado seguindo o formato LLMStructuredResponse.`

const securityPromptTemplate = `# Análise de Segurança - Infrastructure as Code

## Código Terraform

` + "```hcl" + `
{{.TerraformCode}}
` + "```" + `

{{if .CheckovResults}}
{{.CheckovResults}}
{{end}}

{{if .IAMPolicies}}
## Políticas IAM

{{range .IAMPolicies}}
` + "```json" + `
{{.}}
` + "```" + `
{{end}}
{{end}}

## Instruções

Realize uma análise profunda de segurança do código Terraform fornecido. Concentre-se em:

1. **Vulnerabilidades Críticas**: Identifique problemas de segurança de alta severidade.

2. **Configurações Inseguras**: Detecte configurações que violam princípios de segurança.

3. **Problemas de Compliance**: Identifique violações de compliance (GDPR, HIPAA, PCI-DSS, etc).

4. **Exposição de Dados**: Detecte potencial exposição de dados sensíveis.

5. **Problemas de IAM**: Analise permissões excessivas ou inseguras.

6. **Recomendações de Correção**: Forneça exemplos de código para corrigir cada problema.

Forneça sua resposta como um JSON estruturado seguindo o formato SecurityAuditResponse.`

const costPromptTemplate = `# Análise de Custos - Infrastructure as Code

## Código Terraform

` + "```hcl" + `
{{.TerraformCode}}
` + "```" + `

## Instruções

Realize uma análise de otimização de custos do código Terraform fornecido. Concentre-se em:

1. **Dimensionamento de Recursos**: Identifique recursos superdimensionados.

2. **Opções de Compra**: Sugira Reserved Instances, Savings Plans ou outras opções de desconto.

3. **Recursos Desnecessários**: Identifique recursos que podem ser eliminados ou consolidados.

4. **Arquitetura Econômica**: Sugira alternativas arquiteturais mais econômicas.

5. **Estimativas de Economia**: Quando possível, forneça estimativas de economia potencial.

6. **Recomendações de Implementação**: Forneça exemplos de código para cada otimização.

Forneça sua resposta como um JSON estruturado seguindo o formato CostOptimizationResponse.`

const architecturePromptTemplate = `# Análise Arquitetural - Infrastructure as Code

## Código Terraform

` + "```hcl" + `
{{.TerraformCode}}
` + "```" + `

## Instruções

Realize uma análise arquitetural do código Terraform fornecido. Concentre-se em:

1. **Padrões de Arquitetura**: Identifique os padrões arquiteturais utilizados.

2. **Escalabilidade**: Avalie a capacidade de escalar da infraestrutura.

3. **Resiliência**: Analise a resiliência a falhas e disaster recovery.

4. **Modularização**: Sugira oportunidades para melhorar a modularização e reutilização.

5. **Serviços Gerenciados**: Recomende substituição por serviços gerenciados quando apropriado.

6. **Recomendações de Implementação**: Forneça exemplos de código para cada melhoria.

Forneça sua resposta como um JSON estruturado seguindo o formato ArchitecturalInsight.`