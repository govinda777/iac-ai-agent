package analyzer

import (
	"testing"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

func TestPreviewAnalyzer_AnalyzePreview(t *testing.T) {
	log := logger.New("debug", "text")
	pa := NewPreviewAnalyzer(log)

	testCases := []struct {
		name                 string
		planJSON             string
		expectedErr          bool
		expectedRiskLevel    string
		expectedDestroyCount int
		expectedCreateCount  int
		expectedUpdateCount  int
		expectedWarnings     int
	}{
		{
			name: "Should detect database destruction as high risk",
			planJSON: `{
				"resource_changes": [
					{
						"address": "aws_db_instance.default",
						"type": "aws_db_instance",
						"change": {
							"actions": ["delete"]
						}
					}
				]
			}`,
			expectedErr:          false,
			expectedRiskLevel:    "high",
			expectedDestroyCount: 1,
			expectedWarnings:     1,
		},
		{
			name: "Should identify simple resource creation as low risk",
			planJSON: `{
				"resource_changes": [
					{
						"address": "aws_s3_bucket.my_bucket",
						"type": "aws_s3_bucket",
						"change": {
							"actions": ["create"]
						}
					}
				]
			}`,
			expectedErr:         false,
			expectedRiskLevel:   "low",
			expectedCreateCount: 1,
			expectedWarnings:    0,
		},
		{
			name: "Should identify many changes as medium risk",
			planJSON: `{
				"resource_changes": [
					{"address": "res.1", "change": {"actions": ["create"]}},
					{"address": "res.2", "change": {"actions": ["create"]}},
					{"address": "res.3", "change": {"actions": ["create"]}},
					{"address": "res.4", "change": {"actions": ["create"]}},
					{"address": "res.5", "change": {"actions": ["create"]}},
					{"address": "res.6", "change": {"actions": ["create"]}},
					{"address": "res.7", "change": {"actions": ["create"]}},
					{"address": "res.8", "change": {"actions": ["create"]}},
					{"address": "res.9", "change": {"actions": ["create"]}},
					{"address": "res.10", "change": {"actions": ["create"]}},
					{"address": "res.11", "change": {"actions": ["create"]}},
					{"address": "res.12", "change": {"actions": ["create"]}},
					{"address": "res.13", "change": {"actions": ["create"]}},
					{"address": "res.14", "change": {"actions": ["create"]}},
					{"address": "res.15", "change": {"actions": ["create"]}},
					{"address": "res.16", "change": {"actions": ["create"]}},
					{"address": "res.17", "change": {"actions": ["create"]}},
					{"address": "res.18", "change": {"actions": ["create"]}},
					{"address": "res.19", "change": {"actions": ["create"]}},
					{"address": "res.20", "change": {"actions": ["create"]}},
					{"address": "res.21", "change": {"actions": ["create"]}}
				]
			}`,
			expectedErr:       false,
			expectedRiskLevel: "medium",
			expectedCreateCount: 21,
			expectedWarnings:    0,
		},
		{
			name:                 "Should return error for invalid JSON",
			planJSON:             `{ "invalid_json": `,
			expectedErr:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analysis, err := pa.AnalyzePreview([]byte(tc.planJSON))

			if tc.expectedErr {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Did not expect an error but got: %v", err)
			}

			if analysis.RiskLevel != tc.expectedRiskLevel {
				t.Errorf("Expected risk level '%s', but got '%s'", tc.expectedRiskLevel, analysis.RiskLevel)
			}
			if analysis.DestroyCount != tc.expectedDestroyCount {
				t.Errorf("Expected destroy count %d, but got %d", tc.expectedDestroyCount, analysis.DestroyCount)
			}
			if analysis.CreateCount != tc.expectedCreateCount {
				t.Errorf("Expected create count %d, but got %d", tc.expectedCreateCount, analysis.CreateCount)
			}
			if analysis.UpdateCount != tc.expectedUpdateCount {
				t.Errorf("Expected update count %d, but got %d", tc.expectedUpdateCount, analysis.UpdateCount)
			}
			if len(analysis.Warnings) != tc.expectedWarnings {
				t.Errorf("Expected %d warnings, but got %d", tc.expectedWarnings, len(analysis.Warnings))
			}
		})
	}
}

func TestPreviewAnalyzer_detectRiskyChanges(t *testing.T) {
	pa := &PreviewAnalyzer{}

	testCases := []struct {
		name     string
		changes  []models.PlannedChange
		expected int
	}{
		{
			name: "Detects database destruction",
			changes: []models.PlannedChange{
				{Resource: "aws_rds_instance.main", Action: "delete"},
			},
			expected: 1,
		},
		{
			name: "Detects stateful resource replacement",
			changes: []models.PlannedChange{
				{Resource: "aws_s3_bucket.data", Action: "replace"},
			},
			expected: 1,
		},
		{
			name: "Detects networking change",
			changes: []models.PlannedChange{
				{Resource: "aws_vpc.main", Action: "update"},
			},
			expected: 1,
		},
		{
			name:     "No risky changes",
			changes:  []models.PlannedChange{{Resource: "null_resource.foo", Action: "create"}},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			warnings := pa.detectRiskyChanges(tc.changes)
			if len(warnings) != tc.expected {
				t.Errorf("Expected %d warnings, got %d", tc.expected, len(warnings))
			}
		})
	}
}