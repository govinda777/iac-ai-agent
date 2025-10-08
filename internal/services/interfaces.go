package services

import "github.com/govinda777/iac-ai-agent/internal/models"

// TerraformAnalyzerInterface defines the interface for a Terraform analyzer.
type TerraformAnalyzerInterface interface {
	AnalyzeDirectory(dir string) (*models.TerraformAnalysis, error)
	AnalyzeContent(content string, filename string) (*models.TerraformAnalysis, error)
}

// CheckovAnalyzerInterface defines the interface for a Checkov analyzer.
type CheckovAnalyzerInterface interface {
	IsAvailable() bool
	AnalyzeDirectory(dir string, config *models.CheckovConfig) (*models.SecurityAnalysis, error)
	ValidateAndParseResult(jsonResult []byte) (*models.SecurityAnalysis, error)
}

// IAMAnalyzerInterface defines the interface for an IAM analyzer.
type IAMAnalyzerInterface interface {
	AnalyzeTerraform(tfAnalysis *models.TerraformAnalysis) (*models.IAMAnalysis, error)
}

// PRScorerInterface defines the interface for a pull request scorer.
type PRScorerInterface interface {
	CalculateScore(details *models.AnalysisDetails) *models.PRScore
	ShouldApprove(score *models.PRScore, minPassScore int) bool
	GetScoreLevel(score int) string
	GenerateScoreSummary(score *models.PRScore) string
}

// CostOptimizerInterface defines the interface for a cost optimizer.
type CostOptimizerInterface interface {
	AnalyzeCosts(tfAnalysis *models.TerraformAnalysis) *models.CostAnalysis
	GenerateSuggestions(tfAnalysis *models.TerraformAnalysis) []models.Suggestion
}

// SecurityAdvisorInterface defines the interface for a security advisor.
type SecurityAdvisorInterface interface {
	GenerateSuggestions(
		securityAnalysis *models.SecurityAnalysis,
		iamAnalysis *models.IAMAnalysis,
	) []models.Suggestion
}
