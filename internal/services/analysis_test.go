package services

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// Mock Implementations
type mockTerraformAnalyzer struct{}

func (m *mockTerraformAnalyzer) AnalyzeDirectory(dir string) (*models.TerraformAnalysis, error) {
	// Return a realistic path that includes the directory passed to the function.
	return &models.TerraformAnalysis{Files: []string{filepath.Join(dir, "main.tf")}}, nil
}
func (m *mockTerraformAnalyzer) AnalyzeContent(content, filename string) (*models.TerraformAnalysis, error) {
	return &models.TerraformAnalysis{Files: []string{filename}}, nil
}

type mockCheckovAnalyzer struct{}

func (m *mockCheckovAnalyzer) IsAvailable() bool                                          { return true }
func (m *mockCheckovAnalyzer) AnalyzeDirectory(d string, c *models.CheckovConfig) (*models.SecurityAnalysis, error) { return &models.SecurityAnalysis{}, nil }
func (m *mockCheckovAnalyzer) ValidateAndParseResult(j []byte) (*models.SecurityAnalysis, error) { return &models.SecurityAnalysis{}, nil }

type mockIAMAnalyzer struct{}

func (m *mockIAMAnalyzer) AnalyzeTerraform(tf *models.TerraformAnalysis) (*models.IAMAnalysis, error) {
	return &models.IAMAnalysis{}, nil
}

type mockPRScorer struct{}

func (m *mockPRScorer) CalculateScore(d *models.AnalysisDetails) *models.PRScore { return &models.PRScore{Total: 100} }
func (m *mockPRScorer) ShouldApprove(s *models.PRScore, min int) bool             { return true }
func (m *mockPRScorer) GetScoreLevel(s int) string                              { return "good" }
func (m *mockPRScorer) GenerateScoreSummary(s *models.PRScore) string           { return "summary" }

type mockCostOptimizer struct{}

func (m *mockCostOptimizer) AnalyzeCosts(tf *models.TerraformAnalysis) *models.CostAnalysis {
	return &models.CostAnalysis{}
}
func (m *mockCostOptimizer) GenerateSuggestions(tf *models.TerraformAnalysis) []models.Suggestion {
	return nil
}

type mockSecurityAdvisor struct{}

func (m *mockSecurityAdvisor) GenerateSuggestions(sec *models.SecurityAnalysis, iam *models.IAMAnalysis) []models.Suggestion {
	return nil
}

type mockPreviewAnalyzer struct {
	err bool
}

func (m *mockPreviewAnalyzer) AnalyzePreview(planJSON []byte) (*models.PreviewAnalysis, error) {
	if m.err {
		return nil, errors.New("mock preview error")
	}
	var plan models.TerraformPlan
	json.Unmarshal(planJSON, &plan)
	return &models.PreviewAnalysis{RiskLevel: "low", ResourcesAffected: len(plan.ResourceChanges)}, nil
}

type mockSecretsAnalyzer struct {
	findings []models.SecretFinding
}

func (m *mockSecretsAnalyzer) ScanContent(content, filename string) []models.SecretFinding {
	// Return the findings that were set up for this mock.
	return m.findings
}

func TestAnalysisService_AnalyzeDirectory_WithSecrets(t *testing.T) {
	log := logger.New("debug", "text")

	// 1. Create a temporary directory for the test.
	tempDir, err := os.MkdirTemp("", "analysis-secrets-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 2. Create a dummy file inside the temp directory.
	mainTfPath := filepath.Join(tempDir, "main.tf")
	secretFileContent := []byte(`password = "super-secret"`)
	if err := os.WriteFile(mainTfPath, secretFileContent, 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// 3. Set up the mock SecretsAnalyzer to return a specific finding.
	mockSecrets := &mockSecretsAnalyzer{
		findings: []models.SecretFinding{{
			Type:        "Generic Password",
			File:        mainTfPath,
			Line:        1,
			Description: "Password in plaintext",
			Suggestion:  "Use variable with sensitive=true or secrets manager",
			Severity:    "high",
		}},
	}

	as := NewAnalysisService(
		log, 80, &mockTerraformAnalyzer{}, &mockCheckovAnalyzer{}, &mockIAMAnalyzer{},
		&mockPRScorer{}, &mockCostOptimizer{}, &mockSecurityAdvisor{}, &mockPreviewAnalyzer{}, mockSecrets,
		&mockLLMClient{}, &mockKnowledgeBase{}, &mockModuleRegistry{}, &mockPromptBuilder{},
	)

	// 4. Analyze the temporary directory.
	resp, err := as.AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeDirectory failed: %v", err)
	}

	// 5. Assert the results.
	if len(resp.Analysis.SecretFindings) != 1 {
		t.Errorf("Expected 1 secret finding, got %d", len(resp.Analysis.SecretFindings))
	}

	foundSuggestion := false
	for _, s := range resp.Suggestions {
		if s.Type == "security" && s.File == mainTfPath {
			foundSuggestion = true
			break
		}
	}
	if !foundSuggestion {
		t.Error("Expected a suggestion for the found secret, but none was found")
	}
}

func TestAnalysisService_AnalyzePreviewPlan(t *testing.T) {
	log := logger.New("debug", "text")
	as := NewAnalysisService(
		log, 80, nil, nil, nil, nil, nil, nil, &mockPreviewAnalyzer{}, nil,
		nil, nil, nil, nil,
	)

	planJSON := []byte(`{"resource_changes": [{"address": "res.1"}]}`)
	resp, err := as.AnalyzePreviewPlan(planJSON)
	if err != nil {
		t.Fatalf("AnalyzePreviewPlan failed: %v", err)
	}

	if resp.ResourcesAffected != 1 {
		t.Errorf("Expected 1 resource affected, got %d", resp.ResourcesAffected)
	}
}

func TestAnalysisService_AnalyzePreviewPlan_NoAnalyzer(t *testing.T) {
	log := logger.New("debug", "text")
	as := NewAnalysisService(
		log, 80, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil,
	)

	_, err := as.AnalyzePreviewPlan([]byte(`{}`))
	if err == nil {
		t.Fatal("Expected an error when preview analyzer is not available, but got nil")
	}

	if err.Error() != "preview analyzer is not available" {
		t.Errorf("Expected error message 'preview analyzer is not available', got '%s'", err.Error())
	}
}

func TestMain(tm *testing.M) {
	// Setup can go here
	code := tm.Run()
	// Teardown can go here
	os.Exit(code)
}