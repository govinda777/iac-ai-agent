package analyzer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
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
			case "destroy":
				analysis.DestroyCount++
			case "replace":
				analysis.ReplaceCount++
			}
		}
	}

	analysis.ResourcesAffected = len(plan.ResourceChanges)
	analysis.RiskLevel = pa.calculateRiskLevel(analysis)

	// Detecta mudanças arriscadas
	analysis.RiskWarnings = pa.detectRiskyChanges(analysis.PlannedChanges)

	return analysis, nil
}

// analyzeResourceChange analisa uma mudança de recurso individual
func (pa *PreviewAnalyzer) analyzeResourceChange(rc *models.TerraformResourceChange) models.PlannedChange {
	change := models.PlannedChange{
		Resource:  rc.Address,
		Action:    strings.Join(rc.Change.Actions, ","),
		Changes:   make(map[string]interface{}),
		RiskScore: pa.calculateResourceRiskScore(rc),
		Impact:    pa.determineImpact(rc),
	}

	// Analisa mudanças específicas
	if rc.Change.Before != nil && rc.Change.After != nil {
		change.Changes = pa.compareResourceStates(rc.Change.Before, rc.Change.After)
	}

	return change
}

// calculateResourceRiskScore calcula score de risco para um recurso
func (pa *PreviewAnalyzer) calculateResourceRiskScore(rc *models.TerraformResourceChange) int {
	score := 0

	// Ações de alto risco
	for _, action := range rc.Change.Actions {
		switch action {
		case "destroy":
			score += 50
		case "replace":
			score += 30
		case "update":
			score += 10
		case "create":
			score += 5
		}
	}

	// Recursos críticos
	if pa.isCriticalResource(rc.Type) {
		score += 20
	}

	// Recursos stateful
	if pa.isStatefulResource(rc.Type) {
		score += 15
	}

	// Recursos de rede
	if pa.isNetworkResource(rc.Type) {
		score += 10
	}

	return score
}

// determineImpact determina o impacto da mudança
func (pa *PreviewAnalyzer) determineImpact(rc *models.TerraformResourceChange) string {
	for _, action := range rc.Change.Actions {
		switch action {
		case "destroy":
			if pa.isCriticalResource(rc.Type) {
				return "critical"
			}
			return "high"
		case "replace":
			if pa.isStatefulResource(rc.Type) {
				return "high"
			}
			return "medium"
		case "update":
			if pa.isNetworkResource(rc.Type) {
				return "medium"
			}
			return "low"
		case "create":
			return "low"
		}
	}
	return "unknown"
}

// compareResourceStates compara estados antes e depois
func (pa *PreviewAnalyzer) compareResourceStates(before, after map[string]interface{}) map[string]interface{} {
	changes := make(map[string]interface{})

	// Compara campos comuns
	commonFields := []string{"name", "instance_type", "size", "count", "tags"}
	for _, field := range commonFields {
		if beforeVal, beforeExists := before[field]; beforeExists {
			if afterVal, afterExists := after[field]; afterExists {
				if beforeVal != afterVal {
					changes[field] = map[string]interface{}{
						"before": beforeVal,
						"after":  afterVal,
					}
				}
			}
		}
	}

	return changes
}

// detectRiskyChanges identifica mudanças de alto risco
func (pa *PreviewAnalyzer) detectRiskyChanges(changes []models.PlannedChange) []models.RiskWarning {
	warnings := []models.RiskWarning{}

	for _, change := range changes {
		// Destruição de banco de dados
		if pa.isDatabaseResource(change.Resource) && strings.Contains(change.Action, "destroy") {
			warnings = append(warnings, models.RiskWarning{
				Severity: "critical",
				Resource: change.Resource,
				Message:  "⚠️ Database will be DESTROYED. Ensure backups exist!",
				Action:   "destroy",
			})
		}

		// Replace de recursos stateful
		if pa.isStatefulResourceFromAddress(change.Resource) && strings.Contains(change.Action, "replace") {
			warnings = append(warnings, models.RiskWarning{
				Severity: "high",
				Resource: change.Resource,
				Message:  "Resource replacement will cause downtime",
				Action:   "replace",
			})
		}

		// Mudanças em recursos de rede críticos
		if pa.isNetworkResourceFromAddress(change.Resource) && change.RiskScore > 20 {
			warnings = append(warnings, models.RiskWarning{
				Severity: "high",
				Resource: change.Resource,
				Message:  "Network resource changes may affect connectivity",
				Action:   change.Action,
			})
		}

		// Score de risco muito alto
		if change.RiskScore > 50 {
			warnings = append(warnings, models.RiskWarning{
				Severity: "critical",
				Resource: change.Resource,
				Message:  fmt.Sprintf("High risk change detected (score: %d)", change.RiskScore),
				Action:   change.Action,
			})
		}
	}

	return warnings
}

// calculateRiskLevel calcula nível de risco geral
func (pa *PreviewAnalyzer) calculateRiskLevel(analysis *models.PreviewAnalysis) string {
	totalRisk := 0
	criticalCount := 0

	for _, change := range analysis.PlannedChanges {
		totalRisk += change.RiskScore
		if change.RiskScore > 50 {
			criticalCount++
		}
	}

	// Muitas mudanças críticas
	if criticalCount > 3 {
		return "critical"
	}

	// Score total muito alto
	if totalRisk > 200 {
		return "high"
	}

	// Score total médio
	if totalRisk > 100 {
		return "medium"
	}

	return "low"
}

// Helper functions para identificar tipos de recursos

func (pa *PreviewAnalyzer) isCriticalResource(resourceType string) bool {
	criticalTypes := []string{
		"aws_rds_instance",
		"aws_db_instance",
		"aws_eks_cluster",
		"aws_ecs_cluster",
		"aws_vpc",
		"aws_internet_gateway",
		"aws_nat_gateway",
	}
	return pa.containsString(criticalTypes, resourceType)
}

func (pa *PreviewAnalyzer) isStatefulResource(resourceType string) bool {
	statefulTypes := []string{
		"aws_rds_instance",
		"aws_db_instance",
		"aws_ebs_volume",
		"aws_instance",
		"aws_elasticache_cluster",
		"aws_dynamodb_table",
	}
	return pa.containsString(statefulTypes, resourceType)
}

func (pa *PreviewAnalyzer) isNetworkResource(resourceType string) bool {
	networkTypes := []string{
		"aws_vpc",
		"aws_subnet",
		"aws_internet_gateway",
		"aws_nat_gateway",
		"aws_route_table",
		"aws_security_group",
		"aws_network_acl",
		"aws_lb",
		"aws_alb",
		"aws_elb",
	}
	return pa.containsString(networkTypes, resourceType)
}

func (pa *PreviewAnalyzer) isDatabaseResource(address string) bool {
	return strings.Contains(address, "aws_rds_instance") ||
		strings.Contains(address, "aws_db_instance") ||
		strings.Contains(address, "aws_dynamodb_table") ||
		strings.Contains(address, "aws_elasticache_cluster")
}

func (pa *PreviewAnalyzer) isStatefulResourceFromAddress(address string) bool {
	return strings.Contains(address, "aws_rds_instance") ||
		strings.Contains(address, "aws_db_instance") ||
		strings.Contains(address, "aws_ebs_volume") ||
		strings.Contains(address, "aws_instance") ||
		strings.Contains(address, "aws_elasticache_cluster") ||
		strings.Contains(address, "aws_dynamodb_table")
}

func (pa *PreviewAnalyzer) isNetworkResourceFromAddress(address string) bool {
	return strings.Contains(address, "aws_vpc") ||
		strings.Contains(address, "aws_subnet") ||
		strings.Contains(address, "aws_internet_gateway") ||
		strings.Contains(address, "aws_nat_gateway") ||
		strings.Contains(address, "aws_route_table") ||
		strings.Contains(address, "aws_security_group") ||
		strings.Contains(address, "aws_network_acl") ||
		strings.Contains(address, "aws_lb") ||
		strings.Contains(address, "aws_alb") ||
		strings.Contains(address, "aws_elb")
}

func (pa *PreviewAnalyzer) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
