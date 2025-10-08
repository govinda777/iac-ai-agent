package cloudcontroller

import (
	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// KnowledgeBase armazena e gerencia conhecimento sobre IaC
type KnowledgeBase struct {
	logger                *logger.Logger
	bestPractices         map[string][]models.BestPractice
	providerPractices     map[string][]models.BestPractice
	architecturePractices map[string][]models.BestPractice
	platformContext       PlatformContext
	securityPolicies      []SecurityPolicy
}

// PlatformContext contém contexto específico da plataforma
type PlatformContext struct {
	Name              string
	Standards         Standards
	SupportedVersions SupportedVersions
	ApprovedModules   []ApprovedModule
}

// Standards define padrões da plataforma
type Standards struct {
	TaggingPolicy    map[string]string
	NamingConvention string
	RequiredOutputs  []string
}

// SupportedVersions define versões suportadas
type SupportedVersions struct {
	Terraform []string
	OpenTofu  []string
	Providers map[string]string
}

// ApprovedModule representa um módulo aprovado
type ApprovedModule struct {
	Source      string
	Version     string
	UseCase     string
	Recommended bool
}

// SecurityPolicy representa uma política de segurança
type SecurityPolicy struct {
	ID          string
	Title       string
	Description string
	Severity    string
	AutoFix     bool
	FixCode     string
	Resources   []string
}

// NewKnowledgeBase cria uma nova base de conhecimento
func NewKnowledgeBase(log *logger.Logger) *KnowledgeBase {
	kb := &KnowledgeBase{
		logger:                log,
		bestPractices:         make(map[string][]models.BestPractice),
		providerPractices:     make(map[string][]models.BestPractice),
		architecturePractices: make(map[string][]models.BestPractice),
	}

	// Carrega dados iniciais
	kb.loadPlatformContext()
	kb.loadSecurityPolicies()
	kb.loadBestPractices()

	log.Info("Knowledge Base inicializada",
		"best_practices", len(kb.bestPractices),
		"security_policies", len(kb.securityPolicies))

	return kb
}

// GetRelevantPractices retorna práticas relevantes para a análise
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

// detectArchitecturalPattern detecta o padrão arquitetural
func (kb *KnowledgeBase) detectArchitecturalPattern(
	analysis *models.AnalysisDetails,
) string {
	resources := analysis.Terraform.Resources

	// 3-tier web app
	hasLB := kb.containsResourceType(resources, "aws_lb", "aws_alb", "aws_elb")
	hasCompute := kb.containsResourceType(resources, "aws_instance", "aws_ecs_service")
	hasDB := kb.containsResourceType(resources, "aws_rds_instance", "aws_db_instance")
	if hasLB && hasCompute && hasDB {
		return "3-tier-web-app"
	}

	// Serverless
	hasLambda := kb.containsResourceType(resources, "aws_lambda_function")
	hasAPIGW := kb.containsResourceType(resources, "aws_api_gateway_rest_api")
	hasDynamoDB := kb.containsResourceType(resources, "aws_dynamodb_table")
	if hasLambda && (hasAPIGW || hasDynamoDB) {
		return "serverless"
	}

	// Microservices
	hasECS := kb.containsResourceType(resources, "aws_ecs_service")
	hasEKS := kb.containsResourceType(resources, "aws_eks_cluster")
	if (hasECS || hasEKS) && len(resources) > 10 {
		return "microservices"
	}

	return "general"
}

// containsResourceType verifica se os recursos contêm algum dos tipos especificados
func (kb *KnowledgeBase) containsResourceType(
	resources []models.TerraformResource,
	types ...string,
) bool {
	for _, res := range resources {
		for _, t := range types {
			if res.Type == t {
				return true
			}
		}
	}
	return false
}

// deduplicate remove duplicatas de práticas
func (kb *KnowledgeBase) deduplicate(
	practices []models.BestPractice,
) []models.BestPractice {
	seen := make(map[string]bool)
	result := []models.BestPractice{}

	for _, practice := range practices {
		if !seen[practice.ID] {
			seen[practice.ID] = true
			result = append(result, practice)
		}
	}

	return result
}

// loadPlatformContext carrega contexto específico da plataforma
func (kb *KnowledgeBase) loadPlatformContext() {
	kb.platformContext = PlatformContext{
		Name: "Nation.fun",
		Standards: Standards{
			TaggingPolicy: map[string]string{
				"Environment": "required (dev|staging|prod)",
				"Owner":       "required (email)",
				"CostCenter":  "required",
				"ManagedBy":   "required (terraform)",
				"Project":     "required",
			},
			NamingConvention: "{project}-{environment}-{resource}-{random}",
			RequiredOutputs: []string{
				"vpc_id",
				"resource_arns",
				"endpoints",
			},
		},
		SupportedVersions: SupportedVersions{
			Terraform: []string{">= 1.5.0"},
			OpenTofu:  []string{">= 1.6.0"},
			Providers: map[string]string{
				"aws":   ">= 5.0",
				"azure": ">= 3.0",
				"gcp":   ">= 4.0",
			},
		},
		ApprovedModules: []ApprovedModule{
			{
				Source:      "terraform-aws-modules/vpc/aws",
				Version:     "~> 5.0",
				UseCase:     "VPC creation",
				Recommended: true,
			},
			{
				Source:      "terraform-aws-modules/eks/aws",
				Version:     "~> 19.0",
				UseCase:     "EKS cluster",
				Recommended: true,
			},
			// ... mais módulos
		},
	}
}

// loadSecurityPolicies carrega políticas de segurança
func (kb *KnowledgeBase) loadSecurityPolicies() {
	kb.securityPolicies = []SecurityPolicy{
		{
			ID:          "SEC-001",
			Title:       "No Public S3 Buckets",
			Description: "S3 buckets must not be publicly accessible unless explicitly approved",
			Severity:    "critical",
			AutoFix:     true,
			FixCode: `
resource "aws_s3_bucket_public_access_block" "example" {
  bucket = aws_s3_bucket.example.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}`,
		},
		{
			ID:          "SEC-002",
			Title:       "Encryption at Rest Required",
			Description: "All storage resources must be encrypted at rest",
			Severity:    "high",
			Resources:   []string{"aws_s3_bucket", "aws_ebs_volume", "aws_rds_instance"},
		},
		// ... mais políticas
	}
}

// loadBestPractices carrega melhores práticas
func (kb *KnowledgeBase) loadBestPractices() {
	// AWS S3
	kb.bestPractices["aws_s3_bucket"] = []models.BestPractice{
		{
			ID:          "BP-S3-001",
			Title:       "Enable Versioning",
			Description: "S3 buckets should have versioning enabled",
			Category:    "data_protection",
			Severity:    "medium",
		},
		{
			ID:          "BP-S3-002",
			Title:       "Enable Access Logging",
			Description: "S3 buckets should have access logging enabled",
			Category:    "security",
			Severity:    "medium",
		},
	}

	// AWS EC2
	kb.bestPractices["aws_instance"] = []models.BestPractice{
		{
			ID:          "BP-EC2-001",
			Title:       "Use IMDSv2",
			Description: "EC2 instances should use IMDSv2",
			Category:    "security",
			Severity:    "high",
		},
		{
			ID:          "BP-EC2-002",
			Title:       "Use Latest Generation",
			Description: "Use latest generation instance types for better performance/cost",
			Category:    "cost",
			Severity:    "medium",
		},
	}

	// Provider practices
	kb.providerPractices["aws"] = []models.BestPractice{
		{
			ID:          "BP-AWS-001",
			Title:       "Use Provider Alias",
			Description: "Use provider aliases for multi-region resources",
			Category:    "structure",
			Severity:    "low",
		},
		{
			ID:          "BP-AWS-002",
			Title:       "Lock Provider Version",
			Description: "Lock AWS provider version to avoid unexpected changes",
			Category:    "stability",
			Severity:    "medium",
		},
	}

	// Architecture practices
	kb.architecturePractices["serverless"] = []models.BestPractice{
		{
			ID:          "BP-ARCH-001",
			Title:       "Lambda Concurrency",
			Description: "Set Lambda concurrency limits to avoid unexpected costs",
			Category:    "cost",
			Severity:    "medium",
		},
		{
			ID:          "BP-ARCH-002",
			Title:       "API Gateway Throttling",
			Description: "Configure API Gateway throttling to prevent abuse",
			Category:    "security",
			Severity:    "medium",
		},
	}
}
