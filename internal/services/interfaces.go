package services

import "github.com/gosouza/iac-ai-agent/internal/models"

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

// PreviewAnalyzerInterface defines the interface for a Terraform plan preview analyzer.
type PreviewAnalyzerInterface interface {
	AnalyzePreview(planJSON []byte) (*models.PreviewAnalysis, error)
}

// SecretsAnalyzerInterface defines the interface for a secrets scanner.
type SecretsAnalyzerInterface interface {
	ScanContent(content string, filename string) []models.SecretFinding
}

// LLMClientInterface defines the interface for an LLM client.
type LLMClientInterface interface {
	Generate(req *models.LLMRequest) (*models.LLMResponse, error)
}

// KnowledgeBaseInterface defines the interface for the knowledge base.
type KnowledgeBaseInterface interface {
	GetRelevantPractices(analysis *models.AnalysisDetails) []models.BestPractice
	GetSecurityPolicies() []models.SecurityPolicy
	GetPlatformContext() models.PlatformContext
}

// ModuleRegistryInterface defines the interface for a Terraform module registry.
type ModuleRegistryInterface interface {
	FindApplicableModules(resources []models.TerraformResource) []models.ApprovedModule
}

// PromptBuilderInterface defines the interface for building LLM prompts.
type PromptBuilderInterface interface {
	BuildEnrichmentPrompt(
		analysis *models.AnalysisDetails,
		baseSuggestions []models.Suggestion,
		relevantPractices []models.BestPractice,
		relevantModules []models.ApprovedModule,
	) (string, error)
}
