package cloudcontroller

import (
	"testing"

	"github.com/gosouza/iac-ai-agent/internal/models"
)

func TestKnowledgeBase_NewKnowledgeBase(t *testing.T) {
	kb := NewKnowledgeBase()

	if kb == nil {
		t.Fatal("NewKnowledgeBase returned nil")
	}

	// Check if data is loaded
	if kb.platformContext.Name != "Nation.fun" {
		t.Errorf("Expected platform context name to be 'Nation.fun', got '%s'", kb.platformContext.Name)
	}

	if len(kb.securityPolicies) == 0 {
		t.Error("Expected security policies to be loaded, but slice is empty")
	}

	if kb.securityPolicies[0].ID != "SEC-001" {
		t.Errorf("Expected first security policy ID to be 'SEC-001', got '%s'", kb.securityPolicies[0].ID)
	}
}

func TestKnowledgeBase_detectArchitecturalPattern(t *testing.T) {
	kb := NewKnowledgeBase()

	testCases := []struct {
		name      string
		resources []models.TerraformResource
		expected  string
	}{
		{
			name: "Should detect 3-tier-web-app",
			resources: []models.TerraformResource{
				{Type: "aws_lb"},
				{Type: "aws_instance"},
				{Type: "aws_db_instance"},
			},
			expected: "3-tier-web-app",
		},
		{
			name: "Should detect serverless",
			resources: []models.TerraformResource{
				{Type: "aws_lambda_function"},
				{Type: "aws_api_gateway_rest_api"},
				{Type: "aws_dynamodb_table"},
			},
			expected: "serverless",
		},
		{
			name: "Should detect microservices from ECS",
			resources: []models.TerraformResource{
				{Type: "aws_ecs_service"}, {Type: "aws_ecs_service"}, {Type: "aws_ecs_service"},
				{Type: "aws_ecs_service"}, {Type: "aws_ecs_service"}, {Type: "aws_ecs_service"},
				{Type: "aws_ecs_service"}, {Type: "aws_ecs_service"}, {Type: "aws_ecs_service"},
				{Type: "aws_ecs_service"}, {Type: "aws_ecs_service"},
			},
			expected: "microservices",
		},
		{
			name: "Should detect microservices from EKS",
			resources: []models.TerraformResource{
				{Type: "aws_eks_cluster"}, {Type: "aws_instance"}, {Type: "aws_instance"},
				{Type: "aws_instance"}, {Type: "aws_instance"}, {Type: "aws_instance"},
				{Type: "aws_instance"}, {Type: "aws_instance"}, {Type: "aws_instance"},
				{Type: "aws_instance"}, {Type: "aws_instance"},
			},
			expected: "microservices",
		},
		{
			name:      "Should default to general",
			resources: []models.TerraformResource{{Type: "aws_s3_bucket"}},
			expected:  "general",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analysis := &models.AnalysisDetails{
				Terraform: models.TerraformAnalysis{
					Resources: tc.resources,
				},
			}
			pattern := kb.detectArchitecturalPattern(analysis)
			if pattern != tc.expected {
				t.Errorf("Expected pattern '%s', but got '%s'", tc.expected, pattern)
			}
		})
	}
}

func TestKnowledgeBase_GetRelevantPractices(t *testing.T) {
	kb := NewKnowledgeBase()
	// Manually add some practices for testing, since the loader is not implemented yet
	kb.bestPractices["aws_s3_bucket"] = []models.BestPractice{
		{ID: "S3-001", Title: "Enable Versioning"},
	}
	kb.providerPractices["aws"] = []models.BestPractice{
		{ID: "AWS-001", Title: "Use IAM Roles"},
	}
	kb.architecturePractices["serverless"] = []models.BestPractice{
		{ID: "SRV-001", Title: "Use X-Ray for tracing"},
	}

	analysis := &models.AnalysisDetails{
		Terraform: models.TerraformAnalysis{
			Resources: []models.TerraformResource{
				{Type: "aws_s3_bucket"},
				{Type: "aws_lambda_function"},      // to trigger serverless pattern
				{Type: "aws_api_gateway_rest_api"}, // to trigger serverless pattern
			},
			Providers: []string{"aws"},
		},
	}

	practices := kb.GetRelevantPractices(analysis)

	if len(practices) != 3 {
		t.Errorf("Expected 3 relevant practices, but got %d", len(practices))
	}

	foundS3 := false
	foundAWS := false
	foundSRV := false
	for _, p := range practices {
		switch p.ID {
		case "S3-001":
			foundS3 = true
		case "AWS-001":
			foundAWS = true
		case "SRV-001":
			foundSRV = true
		}
	}

	if !foundS3 {
		t.Error("Did not find expected S3 best practice")
	}
	if !foundAWS {
		t.Error("Did not find expected AWS provider practice")
	}
	if !foundSRV {
		t.Error("Did not find expected Serverless architecture practice")
	}
}