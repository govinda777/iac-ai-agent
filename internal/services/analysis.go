package services

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// AnalysisService orquestra análise completa de código IaC
type AnalysisService struct {
	tfAnalyzer      TerraformAnalyzerInterface
	checkovAnalyzer CheckovAnalyzerInterface
	iamAnalyzer     IAMAnalyzerInterface
	prScorer        PRScorerInterface
	costOptimizer   CostOptimizerInterface
	securityAdvisor SecurityAdvisorInterface
	previewAnalyzer PreviewAnalyzerInterface
	secretsAnalyzer SecretsAnalyzerInterface
	llmClient       LLMClientInterface
	knowledgeBase   KnowledgeBaseInterface
	moduleRegistry  ModuleRegistryInterface
	promptBuilder   PromptBuilderInterface
	logger          *logger.Logger
	minPassScore    int
}

// NewAnalysisService cria uma nova instância do serviço de análise com injeção de dependência
func NewAnalysisService(
	log *logger.Logger,
	minPassScore int,
	tfAnalyzer TerraformAnalyzerInterface,
	checkovAnalyzer CheckovAnalyzerInterface,
	iamAnalyzer IAMAnalyzerInterface,
	prScorer PRScorerInterface,
	costOptimizer CostOptimizerInterface,
	securityAdvisor SecurityAdvisorInterface,
	previewAnalyzer PreviewAnalyzerInterface,
	secretsAnalyzer SecretsAnalyzerInterface,
	llmClient LLMClientInterface,
	knowledgeBase KnowledgeBaseInterface,
	moduleRegistry ModuleRegistryInterface,
	promptBuilder PromptBuilderInterface,
) *AnalysisService {
	return &AnalysisService{
		tfAnalyzer:      tfAnalyzer,
		checkovAnalyzer: checkovAnalyzer,
		iamAnalyzer:     iamAnalyzer,
		prScorer:        prScorer,
		costOptimizer:   costOptimizer,
		securityAdvisor: securityAdvisor,
		previewAnalyzer: previewAnalyzer,
		secretsAnalyzer: secretsAnalyzer,
		llmClient:       llmClient,
		knowledgeBase:   knowledgeBase,
		moduleRegistry:  moduleRegistry,
		promptBuilder:   promptBuilder,
		logger:          log,
		minPassScore:    minPassScore,
	}
}

// Analyze é um wrapper que decide entre AnalyzeContent ou AnalyzeDirectory
func (as *AnalysisService) Analyze(req *models.AnalysisRequest) (*models.AnalysisResponse, error) {
	if req.Content != "" {
		filename := "main.tf"
		if req.Path != "" {
			filename = req.Path
		}
		return as.AnalyzeContent(req.Content, filename)
	}
	if req.Path != "" {
		return as.AnalyzeDirectory(req.Path)
	}
	return nil, fmt.Errorf("nenhum conteúdo ou caminho fornecido")
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

	// 4. Análise de Secrets
	var secretFindings []models.SecretFinding
	if as.secretsAnalyzer != nil {
		as.logger.Info("Executando análise de secrets")
		secretFindings = as.secretsAnalyzer.ScanContent(content, filename)
	}

	// 5. Gera sugestões
	suggestions := as.generateSuggestions(tfAnalysis, securityAnalysis, iamAnalysis, secretFindings)

	// 6. Monta análise completa
	analysisDetails := models.AnalysisDetails{
		Terraform:      *tfAnalysis,
		Security:       *securityAnalysis,
		IAM:            *iamAnalysis,
		SecretFindings: secretFindings,
	}

	// 7. Calcula score
	score := as.prScorer.CalculateScore(&analysisDetails)

	// 7. Monta resposta
	response := &models.AnalysisResponse{
		ID:          uuid.New().String(),
		Score:       score.Total,
		Analysis:    analysisDetails,
		Suggestions: suggestions,
		Metadata: map[string]interface{}{
			"pr_score":       score,
			"is_approved":    as.prScorer.ShouldApprove(score, as.minPassScore),
			"score_level":    as.prScorer.GetScoreLevel(score.Total),
			"score_summary":  as.prScorer.GenerateScoreSummary(score),
			"recommendation": as.prScorer.GenerateScoreSummary(score),
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

	// 4. Análise de Secrets
	var allSecretFindings []models.SecretFinding
	if as.secretsAnalyzer != nil {
		as.logger.Info("Executando análise de secrets no diretório")
		// tfAnalysis.Files contains the list of files parsed by terraform
		for _, file := range tfAnalysis.Files {
			content, err := os.ReadFile(file)
			if err != nil {
				as.logger.Warn("Falha ao ler arquivo para análise de secret", "file", file, "error", err)
				continue
			}
			findings := as.secretsAnalyzer.ScanContent(string(content), file)
			allSecretFindings = append(allSecretFindings, findings...)
		}
	}

	// 5. Gera sugestões
	suggestions := as.generateSuggestions(tfAnalysis, securityAnalysis, iamAnalysis, allSecretFindings)

	// 6. Análise de custo (se habilitada)
	costAnalysis := as.costOptimizer.AnalyzeCosts(tfAnalysis)

	// 7. Monta análise completa
	analysisDetails := models.AnalysisDetails{
		Terraform:      *tfAnalysis,
		Security:       *securityAnalysis,
		IAM:            *iamAnalysis,
		Cost:           *costAnalysis,
		SecretFindings: allSecretFindings,
	}

	// 8. Calcula score
	score := as.prScorer.CalculateScore(&analysisDetails)

	// 9. Monta resposta
	response := &models.AnalysisResponse{
		ID:          uuid.New().String(),
		Score:       score.Total,
		Analysis:    analysisDetails,
		Suggestions: suggestions,
		Metadata: map[string]interface{}{
			"pr_score":       score,
			"is_approved":    as.prScorer.ShouldApprove(score, as.minPassScore),
			"score_level":    as.prScorer.GetScoreLevel(score.Total),
			"score_summary":  as.prScorer.GenerateScoreSummary(score),
			"recommendation": as.prScorer.GenerateScoreSummary(score),
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
	secretFindings []models.SecretFinding,
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

	// Sugestões de secrets
	for _, finding := range secretFindings {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "security",
			Severity:       finding.Severity,
			Message:        finding.Description,
			Recommendation: finding.Suggestion,
			File:           finding.File,
			Line:           finding.Line,
		})
	}

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

// AnalyzePreviewPlan analyzes a Terraform plan JSON and returns a structured analysis.
func (as *AnalysisService) AnalyzePreviewPlan(planJSON []byte) (*models.PreviewAnalysis, error) {
	as.logger.Info("Starting Terraform plan preview analysis")

	if as.previewAnalyzer == nil {
		as.logger.Error("PreviewAnalyzer is not initialized")
		return nil, fmt.Errorf("preview analyzer is not available")
	}

	previewAnalysis, err := as.previewAnalyzer.AnalyzePreview(planJSON)
	if err != nil {
		as.logger.Error("Error during preview analysis", "error", err)
		return nil, fmt.Errorf("failed to analyze Terraform plan: %w", err)
	}

	as.logger.Info("Terraform plan preview analysis completed", "risk_level", previewAnalysis.RiskLevel, "changes", previewAnalysis.ResourcesAffected)
	return previewAnalysis, nil
}

// ValidatePreExistingResults valida resultados de análises já executadas externamente
// Este método NÃO executa nenhuma ferramenta, apenas valida os resultados fornecidos
func (as *AnalysisService) ValidatePreExistingResults(
	checkovJSON []byte,
	tfAnalysis *models.TerraformAnalysis,
) (*models.AnalysisResponse, error) {
	as.logger.Info("Validando resultados pré-existentes (sem execução)")

	// 1. Valida resultado Checkov se fornecido
	var securityAnalysis *models.SecurityAnalysis
	if len(checkovJSON) > 0 {
		as.logger.Info("Validando resultado Checkov fornecido")
		secAnalysis, err := as.checkovAnalyzer.ValidateAndParseResult(checkovJSON)
		if err != nil {
			return nil, fmt.Errorf("erro ao validar resultado Checkov: %w", err)
		}
		securityAnalysis = secAnalysis
	} else {
		as.logger.Info("Nenhum resultado Checkov fornecido")
		securityAnalysis = &models.SecurityAnalysis{}
	}

	// 2. Valida análise Terraform se fornecida
	if tfAnalysis == nil {
		as.logger.Warn("Nenhuma análise Terraform fornecida")
		tfAnalysis = &models.TerraformAnalysis{}
	} else {
		if err := as.validateTerraformAnalysis(tfAnalysis); err != nil {
			return nil, fmt.Errorf("análise Terraform inválida: %w", err)
		}
	}

	// 3. Análise IAM baseada no Terraform fornecido
	iamAnalysis, err := as.iamAnalyzer.AnalyzeTerraform(tfAnalysis)
	if err != nil {
		as.logger.Warn("Erro na análise IAM", "error", err)
		iamAnalysis = &models.IAMAnalysis{}
	}

	// 4. Gera sugestões
	// No secret scan is performed in this mode, so we pass an empty slice.
	suggestions := as.generateSuggestions(tfAnalysis, securityAnalysis, iamAnalysis, []models.SecretFinding{})

	// 5. Análise de custo
	costAnalysis := as.costOptimizer.AnalyzeCosts(tfAnalysis)

	// 6. Monta análise completa
	analysisDetails := models.AnalysisDetails{
		Terraform: *tfAnalysis,
		Security:  *securityAnalysis,
		IAM:       *iamAnalysis,
		Cost:      *costAnalysis,
	}

	// 7. Calcula score baseado nos resultados validados
	score := as.prScorer.CalculateScore(&analysisDetails)

	// 8. Monta resposta
	response := &models.AnalysisResponse{
		ID:          uuid.New().String(),
		Score:       score.Total,
		Analysis:    analysisDetails,
		Suggestions: suggestions,
		Metadata: map[string]interface{}{
			"pr_score":        score,
			"is_approved":     as.prScorer.ShouldApprove(score, as.minPassScore),
			"score_level":     as.prScorer.GetScoreLevel(score.Total),
			"score_summary":   as.prScorer.GenerateScoreSummary(score),
			"validation_mode": "pre_existing_results",
		},
		Timestamp: time.Now(),
	}

	as.logger.Info("Validação de resultados concluída",
		"score", score.Total,
		"security_issues", securityAnalysis.TotalIssues,
		"suggestions", len(suggestions))

	return response, nil
}

// validateTerraformAnalysis valida a estrutura de uma análise Terraform
func (as *AnalysisService) validateTerraformAnalysis(analysis *models.TerraformAnalysis) error {
	if analysis == nil {
		return fmt.Errorf("análise Terraform é nula")
	}

	if analysis.TotalResources < 0 {
		return fmt.Errorf("número total de recursos é negativo: %d", analysis.TotalResources)
	}

	if analysis.TotalModules < 0 {
		return fmt.Errorf("número total de módulos é negativo: %d", analysis.TotalModules)
	}

	// Valida consistência dos recursos
	resourceCount := len(analysis.Resources)
	if resourceCount != analysis.TotalResources {
		as.logger.Warn("Inconsistência no número de recursos",
			"total_resources", analysis.TotalResources,
			"actual_count", resourceCount)
	}

	return nil
}
