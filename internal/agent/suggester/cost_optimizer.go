package suggester

import (
	"fmt"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// CostOptimizer gera sugestões de otimização de custo
type CostOptimizer struct {
	logger *logger.Logger
}

// NewCostOptimizer cria uma nova instância
func NewCostOptimizer(log *logger.Logger) *CostOptimizer {
	return &CostOptimizer{
		logger: log,
	}
}

// AnalyzeCosts analisa custos potenciais
func (co *CostOptimizer) AnalyzeCosts(tfAnalysis *models.TerraformAnalysis) *models.CostAnalysis {
	analysis := &models.CostAnalysis{
		Currency:              "USD",
		EstimatedMonthlyCost:  0,
		OptimizationPotential: 0,
		Recommendations:       []models.CostRecommendation{},
	}

	// Estima custos básicos baseado em tipos de recursos
	for _, resource := range tfAnalysis.Resources {
		cost := co.estimateResourceCost(resource)
		analysis.EstimatedMonthlyCost += cost

		// Verifica oportunidades de otimização
		if rec := co.findOptimization(resource); rec != nil {
			analysis.Recommendations = append(analysis.Recommendations, *rec)
			analysis.OptimizationPotential += rec.PotentialSavings
		}
	}

	return analysis
}

// GenerateSuggestions gera sugestões de otimização de custo
func (co *CostOptimizer) GenerateSuggestions(tfAnalysis *models.TerraformAnalysis) []models.Suggestion {
	suggestions := []models.Suggestion{}

	for _, resource := range tfAnalysis.Resources {
		if sugg := co.suggestOptimization(resource); sugg != nil {
			suggestions = append(suggestions, *sugg)
		}
	}

	return suggestions
}

// estimateResourceCost estima custo mensal de um recurso
func (co *CostOptimizer) estimateResourceCost(resource models.TerraformResource) float64 {
	// Estimativas simplificadas
	costMap := map[string]float64{
		"aws_instance":        50.0,
		"aws_rds_instance":    100.0,
		"aws_s3_bucket":       5.0,
		"aws_lambda_function": 10.0,
		"aws_nat_gateway":     45.0,
		"aws_elb":             25.0,
		"aws_alb":             25.0,
	}

	if cost, ok := costMap[resource.Type]; ok {
		return cost
	}

	return 0
}

// findOptimization encontra oportunidades de otimização
func (co *CostOptimizer) findOptimization(resource models.TerraformResource) *models.CostRecommendation {
	switch resource.Type {
	case "aws_instance":
		if instanceType, ok := resource.Attributes["instance_type"].(string); ok {
			if instanceType == "t2.large" || instanceType == "t2.xlarge" {
				return &models.CostRecommendation{
					Resource:                 fmt.Sprintf("%s.%s", resource.Type, resource.Name),
					CurrentCost:              100.0,
					PotentialSavings:         30.0,
					Recommendation:           "Considere usar t3 instances que são mais eficientes e baratas",
					ImplementationDifficulty: "easy",
				}
			}
		}
	case "aws_nat_gateway":
		return &models.CostRecommendation{
			Resource:                 fmt.Sprintf("%s.%s", resource.Type, resource.Name),
			CurrentCost:              45.0,
			PotentialSavings:         15.0,
			Recommendation:           "Considere usar NAT instance para ambientes não-produção",
			ImplementationDifficulty: "medium",
		}
	}

	return nil
}

// suggestOptimization gera sugestão de otimização
func (co *CostOptimizer) suggestOptimization(resource models.TerraformResource) *models.Suggestion {
	rec := co.findOptimization(resource)
	if rec == nil {
		return nil
	}

	return &models.Suggestion{
		Type:             "cost",
		Severity:         "info",
		Message:          fmt.Sprintf("Oportunidade de otimização de custo em %s", rec.Resource),
		Recommendation:   rec.Recommendation,
		File:             resource.File,
		Line:             resource.LineStart,
		EstimatedSavings: fmt.Sprintf("$%.2f/mês", rec.PotentialSavings),
		AutoFixAvailable: false,
	}
}
