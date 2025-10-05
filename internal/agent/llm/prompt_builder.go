package llm

import (
	"fmt"
	"strings"

	"github.com/gosouza/iac-ai-agent/internal/models"
)

// PromptBuilder constrói prompts contextualizados para o LLM
type PromptBuilder struct{}

// NewPromptBuilder cria um novo construtor de prompts
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// BuildAnalysisPrompt constrói prompt para análise de código
func (pb *PromptBuilder) BuildAnalysisPrompt(analysis *models.AnalysisDetails, code string) string {
	var prompt strings.Builder

	prompt.WriteString("# Análise de Infrastructure as Code\n\n")
	prompt.WriteString("## Código Terraform\n")
	prompt.WriteString("```hcl\n")
	prompt.WriteString(code)
	prompt.WriteString("\n```\n\n")

	// Adiciona contexto da análise
	prompt.WriteString("## Contexto da Análise\n\n")

	// Terraform Analysis
	prompt.WriteString("### Análise Terraform\n")
	prompt.WriteString(fmt.Sprintf("- Total de recursos: %d\n", analysis.Terraform.TotalResources))
	prompt.WriteString(fmt.Sprintf("- Total de módulos: %d\n", analysis.Terraform.TotalModules))
	prompt.WriteString(fmt.Sprintf("- Providers: %s\n", strings.Join(analysis.Terraform.Providers, ", ")))
	
	if len(analysis.Terraform.SyntaxErrors) > 0 {
		prompt.WriteString(fmt.Sprintf("- ⚠️ Erros de sintaxe: %d\n", len(analysis.Terraform.SyntaxErrors)))
	}

	// Security Analysis
	if analysis.Security.TotalIssues > 0 {
		prompt.WriteString("\n### Análise de Segurança (Checkov)\n")
		prompt.WriteString(fmt.Sprintf("- 🔴 Crítico: %d\n", analysis.Security.Critical))
		prompt.WriteString(fmt.Sprintf("- 🟠 Alto: %d\n", analysis.Security.High))
		prompt.WriteString(fmt.Sprintf("- 🟡 Médio: %d\n", analysis.Security.Medium))
		prompt.WriteString(fmt.Sprintf("- 🔵 Baixo: %d\n", analysis.Security.Low))
		
		// Top findings
		if len(analysis.Security.Findings) > 0 {
			prompt.WriteString("\nPrincipais problemas:\n")
			for i, finding := range analysis.Security.Findings {
				if i >= 5 {
					break // Limita a 5 findings
				}
				prompt.WriteString(fmt.Sprintf("- [%s] %s (arquivo: %s)\n",
					finding.Severity, finding.CheckName, finding.File))
			}
		}
	}

	// IAM Analysis
	if analysis.IAM.TotalPolicies > 0 || analysis.IAM.TotalRoles > 0 {
		prompt.WriteString("\n### Análise IAM\n")
		prompt.WriteString(fmt.Sprintf("- Políticas: %d\n", analysis.IAM.TotalPolicies))
		prompt.WriteString(fmt.Sprintf("- Roles: %d\n", analysis.IAM.TotalRoles))
		
		if analysis.IAM.AdminAccessDetected {
			prompt.WriteString("- ⚠️ Acesso administrativo detectado\n")
		}
		if analysis.IAM.OverlyPermissive {
			prompt.WriteString("- ⚠️ Políticas excessivamente permissivas\n")
		}
		if len(analysis.IAM.PublicAccess) > 0 {
			prompt.WriteString(fmt.Sprintf("- ⚠️ Recursos com acesso público: %d\n", len(analysis.IAM.PublicAccess)))
		}
	}

	// Tarefa
	prompt.WriteString("\n## Tarefa\n\n")
	prompt.WriteString("Com base na análise acima, forneça:\n\n")
	prompt.WriteString("1. **Resumo Executivo**: Um parágrafo resumindo a qualidade geral do código\n")
	prompt.WriteString("2. **Principais Problemas**: Lista dos 3-5 problemas mais críticos\n")
	prompt.WriteString("3. **Recomendações Prioritárias**: Ações concretas para melhorar a infraestrutura\n")
	prompt.WriteString("4. **Otimizações de Custo**: Sugestões para reduzir gastos (se aplicável)\n")
	prompt.WriteString("5. **Best Practices**: Recomendações de melhores práticas não implementadas\n\n")
	prompt.WriteString("Seja específico, prático e focado em ações que podem ser tomadas imediatamente.\n")

	return prompt.String()
}

// BuildReviewPrompt constrói prompt para review de PR
func (pb *PromptBuilder) BuildReviewPrompt(review *models.ReviewResponse, filesChanged []string) string {
	var prompt strings.Builder

	prompt.WriteString("# Review de Pull Request - Infrastructure as Code\n\n")
	prompt.WriteString(fmt.Sprintf("**Repositório**: %s\n", review.Repository))
	prompt.WriteString(fmt.Sprintf("**PR**: #%d\n", review.PRNumber))
	prompt.WriteString(fmt.Sprintf("**Arquivos modificados**: %d\n\n", review.FilesAnalyzed))

	// Lista arquivos modificados
	if len(filesChanged) > 0 {
		prompt.WriteString("## Arquivos Modificados\n")
		for _, file := range filesChanged {
			prompt.WriteString(fmt.Sprintf("- %s\n", file))
		}
		prompt.WriteString("\n")
	}

	// Análise agregada
	prompt.WriteString("## Resultado da Análise\n\n")
	
	if review.Analysis.Security.TotalIssues > 0 {
		prompt.WriteString(fmt.Sprintf("**Segurança**: %d issues (%d críticos, %d altos)\n",
			review.Analysis.Security.TotalIssues,
			review.Analysis.Security.Critical,
			review.Analysis.Security.High))
	}

	if review.TotalSuggestions > 0 {
		prompt.WriteString(fmt.Sprintf("**Sugestões**: %d recomendações\n", review.TotalSuggestions))
	}

	// File reviews
	if len(review.FileReviews) > 0 {
		prompt.WriteString("\n## Análise por Arquivo\n\n")
		for _, fr := range review.FileReviews {
			prompt.WriteString(fmt.Sprintf("### %s\n", fr.Filename))
			prompt.WriteString(fmt.Sprintf("- Status: %s (+%d -%d)\n", fr.Status, fr.Additions, fr.Deletions))
			
			if len(fr.Suggestions) > 0 {
				prompt.WriteString(fmt.Sprintf("- Sugestões: %d\n", len(fr.Suggestions)))
				for _, sug := range fr.Suggestions {
					if sug.Severity == "critical" || sug.Severity == "high" {
						prompt.WriteString(fmt.Sprintf("  - [%s] %s\n", sug.Severity, sug.Message))
					}
				}
			}
			prompt.WriteString("\n")
		}
	}

	// Tarefa
	prompt.WriteString("\n## Tarefa\n\n")
	prompt.WriteString("Gere um comentário de review para o PR com:\n\n")
	prompt.WriteString("1. **Resumo**: Avaliação geral das mudanças (2-3 frases)\n")
	prompt.WriteString("2. **Pontos Positivos**: O que está bem implementado\n")
	prompt.WriteString("3. **Problemas Críticos**: Issues que DEVEM ser corrigidos antes do merge\n")
	prompt.WriteString("4. **Melhorias Sugeridas**: Recomendações não-bloqueantes\n")
	prompt.WriteString("5. **Veredicto**: APPROVED, CHANGES_REQUESTED, ou COMMENTED\n\n")
	prompt.WriteString("Use tom profissional e construtivo. Seja específico e cite linhas/arquivos quando relevante.\n")

	return prompt.String()
}

// BuildSuggestionPrompt constrói prompt para gerar sugestões específicas
func (pb *PromptBuilder) BuildSuggestionPrompt(finding *models.SecurityFinding, context string) string {
	var prompt strings.Builder

	prompt.WriteString("# Geração de Sugestão de Correção\n\n")
	prompt.WriteString("## Problema Detectado\n\n")
	prompt.WriteString(fmt.Sprintf("**Check**: %s\n", finding.CheckName))
	prompt.WriteString(fmt.Sprintf("**Severidade**: %s\n", finding.Severity))
	prompt.WriteString(fmt.Sprintf("**Recurso**: %s\n", finding.Resource))
	prompt.WriteString(fmt.Sprintf("**Arquivo**: %s (linha %d)\n\n", finding.File, finding.Line))
	
	if finding.Description != "" {
		prompt.WriteString(fmt.Sprintf("**Descrição**: %s\n\n", finding.Description))
	}

	// Contexto do código
	if context != "" {
		prompt.WriteString("## Código Atual\n")
		prompt.WriteString("```hcl\n")
		prompt.WriteString(context)
		prompt.WriteString("\n```\n\n")
	}

	// Tarefa
	prompt.WriteString("## Tarefa\n\n")
	prompt.WriteString("Forneça:\n\n")
	prompt.WriteString("1. **Explicação**: Por que isso é um problema (1-2 frases)\n")
	prompt.WriteString("2. **Código Corrigido**: Exemplo concreto de como corrigir\n")
	prompt.WriteString("3. **Impacto**: O que muda ao implementar a correção\n")
	prompt.WriteString("4. **Referências**: Links úteis (se houver)\n\n")
	prompt.WriteString("Seja prático e direto ao ponto.\n")

	return prompt.String()
}

// BuildCostOptimizationPrompt constrói prompt para análise de custos
func (pb *PromptBuilder) BuildCostOptimizationPrompt(resources []models.TerraformResource) string {
	var prompt strings.Builder

	prompt.WriteString("# Análise de Otimização de Custos\n\n")
	prompt.WriteString("## Recursos Analisados\n\n")

	// Lista recursos por tipo
	resourcesByType := make(map[string][]models.TerraformResource)
	for _, r := range resources {
		resourcesByType[r.Type] = append(resourcesByType[r.Type], r)
	}

	for rType, rList := range resourcesByType {
		prompt.WriteString(fmt.Sprintf("### %s (%d)\n", rType, len(rList)))
		for _, r := range rList {
			prompt.WriteString(fmt.Sprintf("- %s\n", r.Name))
		}
		prompt.WriteString("\n")
	}

	// Tarefa
	prompt.WriteString("## Tarefa\n\n")
	prompt.WriteString("Analise os recursos acima e forneça:\n\n")
	prompt.WriteString("1. **Oportunidades de Economia**: Recursos que podem ser otimizados\n")
	prompt.WriteString("2. **Estimativa de Custo**: Custo mensal aproximado da infraestrutura atual\n")
	prompt.WriteString("3. **Recomendações Específicas**: Para cada recurso otimizável:\n")
	prompt.WriteString("   - Mudança sugerida\n")
	prompt.WriteString("   - Economia estimada\n")
	prompt.WriteString("   - Impacto na operação\n")
	prompt.WriteString("4. **Quick Wins**: Otimizações fáceis de implementar\n\n")
	prompt.WriteString("Use valores realistas baseados em pricing típico de cloud providers.\n")

	return prompt.String()
}

// BuildSecurityAdvisoryPrompt constrói prompt para consultoria de segurança
func (pb *PromptBuilder) BuildSecurityAdvisoryPrompt(analysis *models.SecurityAnalysis, iamAnalysis *models.IAMAnalysis) string {
	var prompt strings.Builder

	prompt.WriteString("# Consultoria de Segurança - Infrastructure as Code\n\n")
	
	// Security findings
	prompt.WriteString("## Findings de Segurança\n\n")
	prompt.WriteString(fmt.Sprintf("- Total de issues: %d\n", analysis.TotalIssues))
	prompt.WriteString(fmt.Sprintf("- Críticos: %d\n", analysis.Critical))
	prompt.WriteString(fmt.Sprintf("- Altos: %d\n", analysis.High))
	prompt.WriteString(fmt.Sprintf("- Médios: %d\n", analysis.Medium))
	prompt.WriteString(fmt.Sprintf("- Baixos: %d\n\n", analysis.Low))

	// Top findings
	if len(analysis.Findings) > 0 {
		prompt.WriteString("### Principais Problemas\n\n")
		for i, finding := range analysis.Findings {
			if i >= 10 {
				break
			}
			prompt.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, finding.Severity, finding.CheckName))
			prompt.WriteString(fmt.Sprintf("   - Recurso: %s\n", finding.Resource))
			if finding.Description != "" {
				prompt.WriteString(fmt.Sprintf("   - Descrição: %s\n", finding.Description))
			}
			prompt.WriteString("\n")
		}
	}

	// IAM concerns
	if iamAnalysis != nil {
		prompt.WriteString("## Análise IAM\n\n")
		if iamAnalysis.AdminAccessDetected {
			prompt.WriteString("⚠️ **CRÍTICO**: Acesso administrativo detectado\n\n")
		}
		if len(iamAnalysis.PublicAccess) > 0 {
			prompt.WriteString("⚠️ **Acesso Público Detectado**:\n")
			for _, pa := range iamAnalysis.PublicAccess {
				prompt.WriteString(fmt.Sprintf("- %s\n", pa))
			}
			prompt.WriteString("\n")
		}
	}

	// Tarefa
	prompt.WriteString("## Tarefa\n\n")
	prompt.WriteString("Como especialista em segurança cloud, forneça:\n\n")
	prompt.WriteString("1. **Avaliação de Risco**: Classificação geral (Crítico/Alto/Médio/Baixo)\n")
	prompt.WriteString("2. **Top 5 Prioridades**: Issues que devem ser corrigidos IMEDIATAMENTE\n")
	prompt.WriteString("3. **Roadmap de Segurança**: Plano em 3 fases (Urgente/Curto Prazo/Médio Prazo)\n")
	prompt.WriteString("4. **Compliance**: Frameworks que podem ser impactados (GDPR, SOC2, PCI-DSS, etc)\n")
	prompt.WriteString("5. **Automação**: Como prevenir esses problemas no futuro\n\n")
	prompt.WriteString("Seja direto e focado em ação. Use linguagem clara para não-especialistas.\n")

	return prompt.String()
}
