package services

import "github.com/gosouza/iac-ai-agent/internal/models"

// mockLLMClient is a mock implementation of LLMClientInterface.
type mockLLMClient struct{}

func (m *mockLLMClient) Generate(req *models.LLMRequest) (*models.LLMResponse, error) {
	return &models.LLMResponse{Content: "mock llm response"}, nil
}

// mockKnowledgeBase is a mock implementation of KnowledgeBaseInterface.
type mockKnowledgeBase struct{}

func (m *mockKnowledgeBase) GetRelevantPractices(analysis *models.AnalysisDetails) []models.BestPractice {
	return nil
}
func (m *mockKnowledgeBase) GetSecurityPolicies() []models.SecurityPolicy {
	return nil
}
func (m *mockKnowledgeBase) GetPlatformContext() models.PlatformContext {
	return models.PlatformContext{}
}

// mockModuleRegistry is a mock implementation of ModuleRegistryInterface.
type mockModuleRegistry struct{}

func (m *mockModuleRegistry) FindApplicableModules(resources []models.TerraformResource) []models.ApprovedModule {
	return nil
}

// mockPromptBuilder is a mock implementation of PromptBuilderInterface.
type mockPromptBuilder struct{}

func (m *mockPromptBuilder) BuildEnrichmentPrompt(
	analysis *models.AnalysisDetails,
	baseSuggestions []models.Suggestion,
	relevantPractices []models.BestPractice,
	relevantModules []models.ApprovedModule,
) (string, error) {
	return "mock prompt", nil
}