package analyzer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// IAMAnalyzer analisa políticas e configurações IAM
type IAMAnalyzer struct {
	logger *logger.Logger
}

// NewIAMAnalyzer cria uma nova instância do analisador IAM
func NewIAMAnalyzer(log *logger.Logger) *IAMAnalyzer {
	return &IAMAnalyzer{
		logger: log,
	}
}

// AnalyzeTerraform analisa recursos IAM no código Terraform
func (ia *IAMAnalyzer) AnalyzeTerraform(tfAnalysis *models.TerraformAnalysis) (*models.IAMAnalysis, error) {
	analysis := &models.IAMAnalysis{
		WildcardActions:     []string{},
		PublicAccess:        []string{},
		Recommendations:     []string{},
		PrincipalRisks:      []models.PrincipalRisk{},
	}

	// Analisa cada recurso
	for _, resource := range tfAnalysis.Resources {
		switch {
		case ia.isIAMPolicy(resource.Type):
			ia.analyzePolicyResource(resource, analysis)
		case ia.isIAMRole(resource.Type):
			ia.analyzeRoleResource(resource, analysis)
			analysis.TotalRoles++
		case ia.hasPublicAccess(resource):
			ia.analyzePublicAccess(resource, analysis)
		}
	}

	// Gera recomendações gerais
	ia.generateRecommendations(analysis)

	return analysis, nil
}

// isIAMPolicy verifica se é um recurso de política IAM
func (ia *IAMAnalyzer) isIAMPolicy(resourceType string) bool {
	policyTypes := []string{
		"aws_iam_policy",
		"aws_iam_role_policy",
		"aws_iam_user_policy",
		"aws_iam_group_policy",
		"azurerm_role_definition",
		"google_project_iam_custom_role",
	}

	for _, t := range policyTypes {
		if t == resourceType {
			return true
		}
	}
	return false
}

// isIAMRole verifica se é um recurso de role IAM
func (ia *IAMAnalyzer) isIAMRole(resourceType string) bool {
	roleTypes := []string{
		"aws_iam_role",
		"azurerm_role_assignment",
		"google_project_iam_member",
	}

	for _, t := range roleTypes {
		if t == resourceType {
			return true
		}
	}
	return false
}

// analyzePolicyResource analisa um recurso de política
func (ia *IAMAnalyzer) analyzePolicyResource(resource models.TerraformResource, analysis *models.IAMAnalysis) {
	analysis.TotalPolicies++

	// Procura pelo documento de política
	var policyDoc string
	if doc, ok := resource.Attributes["policy"].(string); ok {
		policyDoc = doc
	} else if doc, ok := resource.Attributes["policy_document"].(string); ok {
		policyDoc = doc
	}

	if policyDoc == "" {
		return
	}

	// Parse policy document (JSON)
	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(policyDoc), &policy); err != nil {
		ia.logger.Warn("Erro ao fazer parse de política IAM", "resource", resource.Name, "error", err)
		return
	}

	// Analisa statements
	if statements, ok := policy["Statement"].([]interface{}); ok {
		for _, stmt := range statements {
			if statement, ok := stmt.(map[string]interface{}); ok {
				ia.analyzeStatement(statement, resource, analysis)
			}
		}
	}
}

// analyzeStatement analisa um statement de política IAM
func (ia *IAMAnalyzer) analyzeStatement(statement map[string]interface{}, resource models.TerraformResource, analysis *models.IAMAnalysis) {
	// Verifica Effect
	effect, ok := statement["Effect"].(string)
	if !ok || effect != "Allow" {
		return // Só analisa permissões
	}

	// Analisa Actions
	actions := ia.extractActions(statement)
	for _, action := range actions {
		if strings.Contains(action, "*") {
			analysis.WildcardActions = append(analysis.WildcardActions,
				fmt.Sprintf("%s: %s", resource.Name, action))
			analysis.OverlyPermissive = true

			// Detecta acesso admin
			if action == "*" || action == "*:*" {
				analysis.AdminAccessDetected = true
			}
		}
	}

	// Analisa Resource
	resources := ia.extractResources(statement)
	for _, res := range resources {
		if res == "*" {
			analysis.OverlyPermissive = true
		}
	}

	// Analisa Principal
	if principal, ok := statement["Principal"].(map[string]interface{}); ok {
		ia.analyzePrincipal(principal, resource, analysis)
	}
}

// extractActions extrai actions de um statement
func (ia *IAMAnalyzer) extractActions(statement map[string]interface{}) []string {
	actions := []string{}

	if action, ok := statement["Action"].(string); ok {
		actions = append(actions, action)
	} else if actionList, ok := statement["Action"].([]interface{}); ok {
		for _, a := range actionList {
			if actionStr, ok := a.(string); ok {
				actions = append(actions, actionStr)
			}
		}
	}

	return actions
}

// extractResources extrai resources de um statement
func (ia *IAMAnalyzer) extractResources(statement map[string]interface{}) []string {
	resources := []string{}

	if resource, ok := statement["Resource"].(string); ok {
		resources = append(resources, resource)
	} else if resourceList, ok := statement["Resource"].([]interface{}); ok {
		for _, r := range resourceList {
			if resourceStr, ok := r.(string); ok {
				resources = append(resources, resourceStr)
			}
		}
	}

	return resources
}

// analyzePrincipal analisa principal de uma política
func (ia *IAMAnalyzer) analyzePrincipal(principal map[string]interface{}, resource models.TerraformResource, analysis *models.IAMAnalysis) {
	// Verifica se é público
	if aws, ok := principal["AWS"].(string); ok {
		if aws == "*" {
			analysis.PublicAccess = append(analysis.PublicAccess,
				fmt.Sprintf("Política %s permite acesso público (Principal: *)", resource.Name))
			
			risk := models.PrincipalRisk{
				Principal:   "*",
				Type:        "public",
				RiskLevel:   "critical",
				Permissions: ia.extractActions(principal),
				Reason:      "Política permite acesso público irrestrito",
			}
			analysis.PrincipalRisks = append(analysis.PrincipalRisks, risk)
		}
	}
}

// analyzeRoleResource analisa um recurso de role
func (ia *IAMAnalyzer) analyzeRoleResource(resource models.TerraformResource, analysis *models.IAMAnalysis) {
	// Verifica assume role policy
	if assumePolicy, ok := resource.Attributes["assume_role_policy"].(string); ok {
		var policy map[string]interface{}
		if err := json.Unmarshal([]byte(assumePolicy), &policy); err == nil {
			if statements, ok := policy["Statement"].([]interface{}); ok {
				for _, stmt := range statements {
					if statement, ok := stmt.(map[string]interface{}); ok {
						if principal, ok := statement["Principal"].(map[string]interface{}); ok {
							if service, ok := principal["Service"].(string); ok {
								// Verifica serviços de risco
								if ia.isRiskyService(service) {
									risk := models.PrincipalRisk{
										Principal:   service,
										Type:        "service",
										RiskLevel:   "medium",
										Reason:      fmt.Sprintf("Serviço %s pode assumir esta role", service),
									}
									analysis.PrincipalRisks = append(analysis.PrincipalRisks, risk)
								}
							}
						}
					}
				}
			}
		}
	}
}

// hasPublicAccess verifica se recurso tem acesso público
func (ia *IAMAnalyzer) hasPublicAccess(resource models.TerraformResource) bool {
	publicResourceTypes := map[string][]string{
		"aws_s3_bucket": {"acl"},
		"aws_s3_bucket_acl": {"acl"},
		"aws_security_group": {"ingress"},
		"aws_db_instance": {"publicly_accessible"},
		"aws_rds_cluster": {"publicly_accessible"},
	}

	if attrs, ok := publicResourceTypes[resource.Type]; ok {
		for _, attr := range attrs {
			if val, exists := resource.Attributes[attr]; exists {
				if ia.isPublicValue(val) {
					return true
				}
			}
		}
	}

	return false
}

// analyzePublicAccess analisa acesso público em recursos
func (ia *IAMAnalyzer) analyzePublicAccess(resource models.TerraformResource, analysis *models.IAMAnalysis) {
	msg := fmt.Sprintf("Recurso %s.%s pode ter acesso público", resource.Type, resource.Name)
	analysis.PublicAccess = append(analysis.PublicAccess, msg)
}

// isPublicValue verifica se um valor indica acesso público
func (ia *IAMAnalyzer) isPublicValue(value interface{}) bool {
	switch v := value.(type) {
	case string:
		publicValues := []string{"public", "public-read", "public-read-write"}
		for _, pv := range publicValues {
			if v == pv {
				return true
			}
		}
	case bool:
		return v
	case []interface{}:
		// Verifica ingress rules
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				if cidr, ok := m["cidr_blocks"].([]interface{}); ok {
					for _, c := range cidr {
						if c == "0.0.0.0/0" {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// isRiskyService verifica se um serviço é considerado de risco
func (ia *IAMAnalyzer) isRiskyService(service string) bool {
	riskyServices := []string{
		"ec2.amazonaws.com",
		"lambda.amazonaws.com",
	}

	for _, rs := range riskyServices {
		if service == rs {
			return true
		}
	}
	return false
}

// generateRecommendations gera recomendações baseadas na análise
func (ia *IAMAnalyzer) generateRecommendations(analysis *models.IAMAnalysis) {
	if analysis.AdminAccessDetected {
		analysis.Recommendations = append(analysis.Recommendations,
			"Evite usar permissões administrativas (*:*). Use princípio do menor privilégio.")
	}

	if analysis.OverlyPermissive {
		analysis.Recommendations = append(analysis.Recommendations,
			"Políticas com wildcard (*) são muito permissivas. Especifique ações e recursos explicitamente.")
	}

	if len(analysis.PublicAccess) > 0 {
		analysis.Recommendations = append(analysis.Recommendations,
			"Recursos com acesso público devem ser revisados cuidadosamente. Considere restringir acesso.")
	}

	if len(analysis.WildcardActions) > 3 {
		analysis.Recommendations = append(analysis.Recommendations,
			"Muitas ações com wildcard detectadas. Revise as permissões e aplique princípio do menor privilégio.")
	}
}
