package mocks

import (
	"github.com/gosouza/iac-ai-agent/internal/models"
)

// MockCheckovAnalyzer is a mock implementation of the CheckovAnalyzerInterface.
type MockCheckovAnalyzer struct {
	IsAvailableFunc          func() bool
	AnalyzeDirectoryFunc     func(dir string, config *models.CheckovConfig) (*models.SecurityAnalysis, error)
	ValidateAndParseResultFunc func(jsonResult []byte) (*models.SecurityAnalysis, error)
}

// IsAvailable mocks the IsAvailable method.
func (m *MockCheckovAnalyzer) IsAvailable() bool {
	if m.IsAvailableFunc != nil {
		return m.IsAvailableFunc()
	}
	return true // Default to true for tests
}

// AnalyzeDirectory mocks the AnalyzeDirectory method.
func (m *MockCheckovAnalyzer) AnalyzeDirectory(dir string, config *models.CheckovConfig) (*models.SecurityAnalysis, error) {
	if m.AnalyzeDirectoryFunc != nil {
		return m.AnalyzeDirectoryFunc(dir, config)
	}
	// Default mock behavior
	return &models.SecurityAnalysis{
		ChecksPassed: 10,
		ChecksFailed: 2,
		TotalIssues:  2,
		Findings: []models.SecurityFinding{
			{
				CheckID:   "CKV_AWS_117",
				CheckName: "Ensure all data stored in the S3 bucket is securely encrypted at rest",
				Severity:  "HIGH",
				Resource:  "aws_s3_bucket.example",
				File:      "main.tf",
				Line:      1,
			},
			{
				CheckID:   "CKV_AWS_18",
				CheckName: "Ensure the S3 bucket has access logging enabled",
				Severity:  "MEDIUM",
				Resource:  "aws_s3_bucket.example",
				File:      "main.tf",
				Line:      1,
			},
		},
	}, nil
}

// ValidateAndParseResult mocks the ValidateAndParseResult method.
func (m *MockCheckovAnalyzer) ValidateAndParseResult(jsonResult []byte) (*models.SecurityAnalysis, error) {
	if m.ValidateAndParseResultFunc != nil {
		return m.ValidateAndParseResultFunc(jsonResult)
	}
	// Default mock behavior
	return &models.SecurityAnalysis{
		ChecksPassed: 5,
		ChecksFailed: 1,
		TotalIssues:  1,
	}, nil
}
