package analyzer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

type PreviewAnalyzer struct {
	logger *logger.Logger
}

func NewPreviewAnalyzer(log *logger.Logger) *PreviewAnalyzer {
	return &PreviewAnalyzer{logger: log}
}

// AnalyzePreview analisa resultado de terraform plan em formato JSON
func (pa *PreviewAnalyzer) AnalyzePreview(planJSON []byte) (*models.PreviewAnalysis, error) {
	var plan models.TerraformPlan
	if err := json.Unmarshal(planJSON, &plan); err != nil {
		return nil, fmt.Errorf("invalid terraform plan JSON: %w", err)
	}

	analysis := &models.PreviewAnalysis{
		PlannedChanges:    []models.PlannedChange{},
		Errors:            []string{},
		Warnings:          []string{},
		ResourcesAffected: 0,
		RiskLevel:         "low",
	}

	// Analisa resource_changes
	for _, rc := range plan.ResourceChanges {
		change := pa.analyzeResourceChange(&rc)
		analysis.PlannedChanges = append(analysis.PlannedChanges, change)

		// Contabiliza ações
		for _, action := range rc.Change.Actions {
			switch action {
			case "create":
				analysis.CreateCount++
			case "update":
				analysis.UpdateCount++
			case "delete":
				analysis.DestroyCount++
			case "replace":
				analysis.ReplaceCount++
			}
		}
	}

	analysis.ResourcesAffected = len(plan.ResourceChanges)
	analysis.RiskLevel = pa.calculateRiskLevel(analysis)

	// Detecta mudanças arriscadas
	riskyChanges := pa.detectRiskyChanges(analysis.PlannedChanges)
	for _, risk := range riskyChanges {
		analysis.Warnings = append(analysis.Warnings, risk.Message)
	}

	return analysis, nil
}

func (pa *PreviewAnalyzer) analyzeResourceChange(rc *struct {
	Address string `json:"address"`
	Mode    string `json:"mode"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Change  struct {
		Actions []string               `json:"actions"`
		Before  map[string]interface{} `json:"before"`
		After   map[string]interface{} `json:"after"`
	} `json:"change"`
}) models.PlannedChange {
	action := "no-op"
	if len(rc.Change.Actions) > 0 {
		action = rc.Change.Actions[0] // Typically the primary action
		if action == "update" && len(rc.Change.Actions) > 1 && rc.Change.Actions[0] == "delete" && rc.Change.Actions[1] == "create" {
			action = "replace"
		}
	}

	return models.PlannedChange{
		Resource: rc.Address,
		Action:   action,
		Changes:  rc.Change.After, // Simplified for now
	}
}


// detectRiskyChanges identifica mudanças de alto risco
func (pa *PreviewAnalyzer) detectRiskyChanges(changes []models.PlannedChange) []models.RiskWarning {
	warnings := []models.RiskWarning{}

	for _, change := range changes {
		// Destruição de banco de dados
		if isDatabase(change.Resource) && change.Action == "delete" {
			warnings = append(warnings, models.RiskWarning{
				Severity: "critical",
				Resource: change.Resource,
				Message:  "⚠️ Database will be DESTROYED. Ensure backups exist!",
				Action:   "delete",
			})
		}

		// Replace de recursos stateful
		if isStateful(change.Resource) && change.Action == "replace" {
			warnings = append(warnings, models.RiskWarning{
				Severity: "high",
				Resource: change.Resource,
				Message:  "Resource replacement will cause downtime",
				Action:   "replace",
			})
		}

		// Mudança de VPC/Network
		if isNetworking(change.Resource) {
			warnings = append(warnings, models.RiskWarning{
				Severity: "high",
				Resource: change.Resource,
				Message:  "Network change may affect connectivity",
				Action:   change.Action,
			})
		}
	}

	return warnings
}

// calculateRiskLevel calcula nível de risco geral
func (pa *PreviewAnalyzer) calculateRiskLevel(analysis *models.PreviewAnalysis) string {
	// Destruições são sempre high risk
	if analysis.DestroyCount > 0 {
		return "high"
	}

	// Replace de muitos recursos é medium risk
	if analysis.ReplaceCount > 5 {
		return "medium"
	}

	// Muitas mudanças é medium risk
	if analysis.ResourcesAffected > 20 {
		return "medium"
	}

	return "low"
}

// isDatabase checks if a resource is a database.
func isDatabase(resource string) bool {
	dbTypes := []string{"aws_rds_instance", "aws_db_instance", "aws_dynamodb_table", "google_sql_database_instance", "azurerm_sql_database"}
	for _, t := range dbTypes {
		if strings.Contains(resource, t) {
			return true
		}
	}
	return false
}

// isStateful checks if a resource is stateful.
func isStateful(resource string) bool {
	if isDatabase(resource) {
		return true
	}
	statefulTypes := []string{"aws_s3_bucket", "aws_ebs_volume", "aws_efs_file_system", "google_storage_bucket", "azurerm_storage_account"}
	for _, t := range statefulTypes {
		if strings.Contains(resource, t) {
			return true
		}
	}
	return false
}

// isNetworking checks if a resource is a core networking component.
func isNetworking(resource string) bool {
	networkTypes := []string{"aws_vpc", "aws_subnet", "aws_security_group", "aws_route_table", "aws_nat_gateway", "aws_internet_gateway", "google_compute_network", "azurerm_virtual_network"}
	for _, t := range networkTypes {
		if strings.Contains(resource, t) {
			return true
		}
	}
	return false
}