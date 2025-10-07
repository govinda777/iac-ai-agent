package cloudcontroller

import (
	"github.com/gosouza/iac-ai-agent/internal/models"
)

// KnowledgeBase contém conhecimento sobre best practices, padrões e políticas.
type KnowledgeBase struct {
	platformContext       models.PlatformContext
	securityPolicies      []models.SecurityPolicy
	bestPractices         map[string][]models.BestPractice
	architecturePractices map[string][]models.BestPractice
	providerPractices     map[string][]models.BestPractice
}

// NewKnowledgeBase cria e carrega uma nova base de conhecimento.
func NewKnowledgeBase() *KnowledgeBase {
	kb := &KnowledgeBase{
		bestPractices:         make(map[string][]models.BestPractice),
		architecturePractices: make(map[string][]models.BestPractice),
		providerPractices:     make(map[string][]models.BestPractice),
	}

	kb.loadPlatformContext()
	kb.loadSecurityPolicies()
	// Future loaders for best practices can be added here
	return kb
}

// GetRelevantPractices coleta best practices relevantes para uma dada análise.
func (kb *KnowledgeBase) GetRelevantPractices(
	analysis *models.AnalysisDetails,
) []models.BestPractice {
	relevant := []models.BestPractice{}

	// Por tipo de recurso
	for _, resource := range analysis.Terraform.Resources {
		if practices, ok := kb.bestPractices[resource.Type]; ok {
			relevant = append(relevant, practices...)
		}
	}

	// Por provider
	for _, provider := range analysis.Terraform.Providers {
		if practices, ok := kb.providerPractices[provider]; ok {
			relevant = append(relevant, practices...)
		}
	}

	// Por padrão arquitetural detectado
	pattern := kb.detectArchitecturalPattern(analysis)
	if practices, ok := kb.architecturePractices[pattern]; ok {
		relevant = append(relevant, practices...)
	}

	return kb.deduplicate(relevant)
}

// detectArchitecturalPattern infere o padrão de arquitetura a partir dos recursos.
func (kb *KnowledgeBase) detectArchitecturalPattern(
	analysis *models.AnalysisDetails,
) string {
	resources := analysis.Terraform.Resources

	// 3-tier web app
	hasLB := containsResourceType(resources, "aws_lb", "aws_alb", "aws_elb")
	hasCompute := containsResourceType(resources, "aws_instance", "aws_ecs_service")
	hasDB := containsResourceType(resources, "aws_rds_instance", "aws_db_instance")
	if hasLB && hasCompute && hasDB {
		return "3-tier-web-app"
	}

	// Serverless
	hasLambda := containsResourceType(resources, "aws_lambda_function")
	hasAPIGW := containsResourceType(resources, "aws_api_gateway_rest_api")
	hasDynamoDB := containsResourceType(resources, "aws_dynamodb_table")
	if hasLambda && (hasAPIGW || hasDynamoDB) {
		return "serverless"
	}

	// Microservices
	hasECS := containsResourceType(resources, "aws_ecs_service")
	hasEKS := containsResourceType(resources, "aws_eks_cluster")
	if (hasECS || hasEKS) && len(resources) > 10 {
		return "microservices"
	}

	return "general"
}

// containsResourceType é um helper para verificar a existência de tipos de recursos.
func containsResourceType(resources []models.TerraformResource, types ...string) bool {
	typeSet := make(map[string]struct{}, len(types))
	for _, t := range types {
		typeSet[t] = struct{}{}
	}

	for _, r := range resources {
		if _, ok := typeSet[r.Type]; ok {
			return true
		}
	}
	return false
}

// deduplicate remove práticas duplicadas de uma lista.
func (kb *KnowledgeBase) deduplicate(practices []models.BestPractice) []models.BestPractice {
	seen := make(map[string]bool)
	result := []models.BestPractice{}
	for _, p := range practices {
		if p.ID != "" {
			if _, ok := seen[p.ID]; !ok {
				seen[p.ID] = true
				result = append(result, p)
			}
		}
	}
	return result
}

// GetSecurityPolicies retorna todas as políticas de segurança carregadas.
func (kb *KnowledgeBase) GetSecurityPolicies() []models.SecurityPolicy {
	return kb.securityPolicies
}

// GetPlatformContext retorna o contexto da plataforma carregado.
func (kb *KnowledgeBase) GetPlatformContext() models.PlatformContext {
	return kb.platformContext
}