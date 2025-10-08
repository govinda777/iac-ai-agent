package scorer

import (
	"math"

	"github.com/govinda777/iac-ai-agent/internal/models"
)

// PRScorer calcula scores de qualidade para PRs
type PRScorer struct{}

// NewPRScorer cria uma nova instância do scorer
func NewPRScorer() *PRScorer {
	return &PRScorer{}
}

// CalculateScore calcula o score total de um PR baseado na análise
func (ps *PRScorer) CalculateScore(analysis *models.AnalysisDetails) *models.PRScore {
	score := &models.PRScore{
		Breakdown: make(map[string]int),
	}

	// Calcula scores individuais
	score.Security = ps.calculateSecurityScore(&analysis.Security)
	score.BestPractices = ps.calculateBestPracticesScore(&analysis.Terraform)
	score.Performance = ps.calculatePerformanceScore(&analysis.Terraform)
	score.Maintainability = ps.calculateMaintainabilityScore(&analysis.Terraform)
	score.Documentation = ps.calculateDocumentationScore(&analysis.Terraform)

	// Pesos para cada categoria
	weights := map[string]float64{
		"security":        0.35, // 35%
		"best_practices":  0.25, // 25%
		"performance":     0.15, // 15%
		"maintainability": 0.15, // 15%
		"documentation":   0.10, // 10%
	}

	// Calcula score ponderado
	totalScore := float64(score.Security)*weights["security"] +
		float64(score.BestPractices)*weights["best_practices"] +
		float64(score.Performance)*weights["performance"] +
		float64(score.Maintainability)*weights["maintainability"] +
		float64(score.Documentation)*weights["documentation"]

	score.Total = int(math.Round(totalScore))

	// Preenche breakdown
	score.Breakdown["security"] = score.Security
	score.Breakdown["best_practices"] = score.BestPractices
	score.Breakdown["performance"] = score.Performance
	score.Breakdown["maintainability"] = score.Maintainability
	score.Breakdown["documentation"] = score.Documentation

	return score
}

// calculateSecurityScore calcula score de segurança (0-100)
func (ps *PRScorer) calculateSecurityScore(security *models.SecurityAnalysis) int {
	if security.TotalIssues == 0 {
		return 100
	}

	// Penalidades por severidade
	penalties := security.Critical*20 +
		security.High*10 +
		security.Medium*5 +
		security.Low*2

	score := 100 - penalties

	// Mínimo de 0
	if score < 0 {
		score = 0
	}

	return score
}

// calculateBestPracticesScore calcula score de best practices (0-100)
func (ps *PRScorer) calculateBestPracticesScore(terraform *models.TerraformAnalysis) int {
	score := 100

	// Penaliza por warnings
	score -= len(terraform.BestPracticeWarnings) * 3

	// Penaliza por erros de sintaxe
	score -= len(terraform.SyntaxErrors) * 10

	// Bonus por ter outputs
	if terraform.TotalOutputs > 0 {
		score += 5
	}

	// Bonus por ter módulos (reutilização)
	if terraform.TotalModules > 0 {
		score += 5
	}

	// Verifica variáveis sem descrição
	varsWithoutDesc := 0
	for _, v := range terraform.Variables {
		if v.Description == "" {
			varsWithoutDesc++
		}
	}
	score -= varsWithoutDesc * 2

	// Normaliza entre 0-100
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// calculatePerformanceScore calcula score de performance (0-100)
func (ps *PRScorer) calculatePerformanceScore(terraform *models.TerraformAnalysis) int {
	score := 100

	// Analisa complexidade
	if terraform.TotalResources > 50 {
		score -= 10 // Muitos recursos em um único lugar
	}

	// Penaliza por não usar módulos quando há muitos recursos
	if terraform.TotalResources > 20 && terraform.TotalModules == 0 {
		score -= 15
	}

	// Bonus por modularização apropriada
	if terraform.TotalModules > 0 && terraform.TotalResources < 30 {
		score += 10
	}

	// Normaliza
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// calculateMaintainabilityScore calcula score de manutenibilidade (0-100)
func (ps *PRScorer) calculateMaintainabilityScore(terraform *models.TerraformAnalysis) int {
	score := 100

	// Penaliza por falta de modularização
	if terraform.TotalResources > 15 && terraform.TotalModules == 0 {
		score -= 20
	}

	// Penaliza por não ter variáveis parametrizadas
	if terraform.TotalResources > 5 && terraform.TotalVariables < 3 {
		score -= 15
	}

	// Verifica uso de valores hardcoded (heurística)
	// Em uma implementação real, isso seria detectado no parser

	// Bonus por boa organização
	if terraform.TotalModules > 2 {
		score += 10
	}

	// Normaliza
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// calculateDocumentationScore calcula score de documentação (0-100)
func (ps *PRScorer) calculateDocumentationScore(terraform *models.TerraformAnalysis) int {
	var varScore, outputScore int

	// Calcula score das variáveis (até 50 pontos)
	totalVars := terraform.TotalVariables
	if totalVars > 0 {
		varsWithDesc := 0
		for _, v := range terraform.Variables {
			if v.Description != "" {
				varsWithDesc++
			}
		}
		docPercentage := float64(varsWithDesc) / float64(totalVars)
		varScore = int(docPercentage * 50)
	} else {
		varScore = 50 // Score neutro se não há variáveis
	}

	// Calcula score dos outputs (até 50 pontos)
	totalOutputs := terraform.TotalOutputs
	if totalOutputs > 0 {
		outputsWithDesc := 0
		for _, o := range terraform.Outputs {
			if o.Description != "" {
				outputsWithDesc++
			}
		}
		docPercentage := float64(outputsWithDesc) / float64(totalOutputs)
		outputScore = int(docPercentage * 50)
	} else {
		outputScore = 50 // Score neutro se não há outputs
	}

	score := varScore + outputScore

	// Normaliza para garantir que está entre 0 e 100
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// GetScoreLevel retorna o nível do score (Excelente, Bom, Regular, Ruim)
func (ps *PRScorer) GetScoreLevel(score int) string {
	switch {
	case score >= 90:
		return "Excelente"
	case score >= 75:
		return "Bom"
	case score >= 60:
		return "Regular"
	case score >= 40:
		return "Ruim"
	default:
		return "Crítico"
	}
}

// ShouldApprove determina se o PR deveria ser aprovado baseado no score
func (ps *PRScorer) ShouldApprove(score *models.PRScore, minScore int) bool {
	// Verifica score mínimo
	if score.Total < minScore {
		return false
	}

	// Não aprova se há problemas críticos de segurança
	if score.Security < 50 {
		return false
	}

	return true
}

// GenerateScoreSummary gera um resumo textual do score
func (ps *PRScorer) GenerateScoreSummary(score *models.PRScore) string {
	level := ps.GetScoreLevel(score.Total)

	summary := "📊 **Score de Qualidade**: " + string(rune(score.Total)) + "/100 - " + level + "\n\n"
	summary += "**Breakdown**:\n"
	summary += "- 🔒 Segurança: " + string(rune(score.Security)) + "/100\n"
	summary += "- ✅ Best Practices: " + string(rune(score.BestPractices)) + "/100\n"
	summary += "- ⚡ Performance: " + string(rune(score.Performance)) + "/100\n"
	summary += "- 🔧 Manutenibilidade: " + string(rune(score.Maintainability)) + "/100\n"
	summary += "- 📚 Documentação: " + string(rune(score.Documentation)) + "/100\n"

	return summary
}
