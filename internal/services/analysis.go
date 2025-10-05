package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gosouza/iac-ai-agent/internal/agent/analyzer"
	"github.com/gosouza/iac-ai-agent/internal/agent/scorer"
	"github.com/gosouza/iac-ai-agent/internal/agent/suggester"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// AnalysisService orquestra análise completa de código IaC
type AnalysisService struct {
	tfAnalyzer      *analyzer.TerraformAnalyzer
	checkovAnalyzer *analyzer.CheckovAnalyzer
	iamAnalyzer     *analyzer.IAMAnalyzer
	prScorer        *scorer.PRScorer
	costOptimizer   *suggester.CostOptimizer
	securityAdvisor *suggester.SecurityAdvisor
	logger          *logger.Logger
}

// NewAnalysisService cria uma nova instância do serviço de análise
func NewAnalysisService(log *logger.Logger, minPassScore int) *AnalysisService {
	return &AnalysisService{
		tfAnalyzer:      analyzer.NewTerraformAnalyzer(),
		checkovAnalyzer: analyzer.NewCheckovAnalyzer(log),
		iamAnalyzer:     analyzer.NewIAMAnalyzer(log),
		prScorer:        scorer.NewPRScorer(minPassScore),
		costOptimizer:   suggester.NewCostOptimizer(log),
		securityAdvisor: suggester.NewSecurityAdvisor(log),
		logger:          log,
	}
}

// AnalyzeContent analisa conteúdo Terraform
func (as *AnalysisService) AnalyzeContent(content string, filename string) (*models.AnalysisResponse, error) {
	as.logger.Info("Iniciando análise de conteúdo", "filename", filename)

	// 1. Análise Terraform
	tfAnalysis, err := as.tfAnalyzer.AnalyzeContent(content, filename)
	if err != nil {
		return nil, fmt.Errorf("erro na análise Terraform: %w", err)
	}

	// 2. Análise IAM
	iamAnalysis, err := as.iamAnalyzer.AnalyzeTerraform(tfAnalysis)
	if err != nil {
		return nil, fmt.Errorf("erro na análise IAM: %w", err)
	}

	// 3. Análise de segurança (Checkov é opcional)
	var securityAnalysis *models.SecurityAnalysis
	if as.checkovAnalyzer.IsAvailable() {
		as.logger.Info("Checkov disponível, executando análise de segurança")
		// Para análise de conteúdo, não executamos Checkov diretamente
		securityAnalysis = &models.SecurityAnalysis{}
	} else {
		as.logger.Warn("Checkov não disponível, pulando análise de segurança")
		securityAnalysis = &models.SecurityAnalysis{}
	}

	// 4. Calcula score
	var checkovResult *models.CheckovResult
	score := as.prScorer.CalculateScore(tfAnalysis, checkovResult)

	// 5. Gera sugestões
	suggestions := as.generateSuggestions(tfAnalysis, securityAnalysis, iamAnalysis)

	// 6. Monta resposta
	response := &models.AnalysisResponse{
		ID:    uuid.New().String(),
		Score: score.Total,
		Analysis: models.AnalysisDetails{
			Terraform: *tfAnalysis,
			Security:  *securityAnalysis,
			IAM:       *iamAnalysis,
		},
		Suggestions: suggestions,
		Metadata: map[string]interface{}{
			"pr_score":       score,
			"is_approved":    as.prScorer.IsApproved(score),
			"recommendation": as.prScorer.GetRecommendation(score),
		},
		Timestamp: time.Now(),
	}

	as.logger.Info("Análise concluída", "score", score.Total, "suggestions", len(suggestions))
	return response, nil
}

// AnalyzeDirectory analisa um diretório completo
func (as *AnalysisService) AnalyzeDirectory(dir string) (*models.AnalysisResponse, error) {
	as.logger.Info("Iniciando análise de diretório", "directory", dir)

	// 1. Análise Terraform
	tfAnalysis, err := as.tfAnalyzer.AnalyzeDirectory(dir)
	if err != nil {
		return nil, fmt.Errorf("erro na análise Terraform: %w", err)
	}

	// 2. Análise IAM
	iamAnalysis, err := as.iamAnalyzer.AnalyzeTerraform(tfAnalysis)
	if err != nil {
		return nil, fmt.Errorf("erro na análise IAM: %w", err)
	}

	// 3. Análise de segurança (Checkov)
	var securityAnalysis *models.SecurityAnalysis
	var checkovResult *models.CheckovResult

	if as.checkovAnalyzer.IsAvailable() {
		as.logger.Info("Executando análise Checkov")
		config := &models.CheckovConfig{
			Directory:     dir,
			Framework:     "terraform",
			CompactOutput: true,
			Quiet:         true,
		}

		secAnalysis, err := as.checkovAnalyzer.AnalyzeDirectory(dir, config)
		if err != nil {
			as.logger.Warn("Erro na análise Checkov", "error", err)
			securityAnalysis = &models.SecurityAnalysis{}
		} else {
			securityAnalysis = secAnalysis
		}
	} else {
		as.logger.Warn("Checkov não disponível")
		securityAnalysis = &models.SecurityAnalysis{}
	}

	// 4. Calcula score
	score := as.prScorer.CalculateScore(tfAnalysis, checkovResult)

	// 5. Gera sugestões
	suggestions := as.generateSuggestions(tfAnalysis, securityAnalysis, iamAnalysis)

	// 6. Análise de custo (se habilitada)
	costAnalysis := as.costOptimizer.AnalyzeCosts(tfAnalysis)

	// 7. Monta resposta
	response := &models.AnalysisResponse{
		ID:    uuid.New().String(),
		Score: score.Total,
		Analysis: models.AnalysisDetails{
			Terraform: *tfAnalysis,
			Security:  *securityAnalysis,
			IAM:       *iamAnalysis,
			Cost:      *costAnalysis,
		},
		Suggestions: suggestions,
		Metadata: map[string]interface{}{
			"pr_score":       score,
			"is_approved":    as.prScorer.IsApproved(score),
			"recommendation": as.prScorer.GetRecommendation(score),
		},
		Timestamp: time.Now(),
	}

	as.logger.Info("Análise de diretório concluída",
		"score", score.Total,
		"resources", tfAnalysis.TotalResources,
		"suggestions", len(suggestions))

	return response, nil
}

// generateSuggestions gera sugestões baseadas nas análises
func (as *AnalysisService) generateSuggestions(
	tfAnalysis *models.TerraformAnalysis,
	securityAnalysis *models.SecurityAnalysis,
	iamAnalysis *models.IAMAnalysis,
) []models.Suggestion {
	suggestions := []models.Suggestion{}

	// Sugestões de best practices
	for _, warning := range tfAnalysis.BestPracticeWarnings {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "best_practice",
			Severity:       "medium",
			Message:        warning,
			Recommendation: "Siga as best practices do Terraform",
		})
	}

	// Sugestões de segurança
	secSuggestions := as.securityAdvisor.GenerateSuggestions(securityAnalysis, iamAnalysis)
	suggestions = append(suggestions, secSuggestions...)

	// Sugestões de custo
	costSuggestions := as.costOptimizer.GenerateSuggestions(tfAnalysis)
	suggestions = append(suggestions, costSuggestions...)

	return suggestions
}

// ValidateAnalysis valida se uma análise está completa e correta
func (as *AnalysisService) ValidateAnalysis(analysis *models.AnalysisResponse) error {
	if analysis == nil {
		return fmt.Errorf("análise é nula")
	}

	if analysis.ID == "" {
		return fmt.Errorf("análise sem ID")
	}

	if analysis.Score < 0 || analysis.Score > 100 {
		return fmt.Errorf("score inválido: %d", analysis.Score)
	}

	return nil
}
