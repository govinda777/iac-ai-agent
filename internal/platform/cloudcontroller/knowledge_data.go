package cloudcontroller

import "github.com/gosouza/iac-ai-agent/internal/models"

// loadPlatformContext carrega contexto específico da plataforma Nation.fun
func (kb *KnowledgeBase) loadPlatformContext() {
    kb.platformContext = models.PlatformContext{
        Name: "Nation.fun",
        Standards: models.Standards{
            TaggingPolicy: map[string]string{
                "Environment":  "required (dev|staging|prod)",
                "Owner":        "required (email)",
                "CostCenter":   "required",
                "ManagedBy":    "required (terraform)",
                "Project":      "required",
            },
            NamingConvention: "{project}-{environment}-{resource}-{random}",
            RequiredOutputs: []string{
                "vpc_id",
                "resource_arns",
                "endpoints",
            },
        },
        SupportedVersions: models.SupportedVersions{
            Terraform: []string{">= 1.5.0"},
            OpenTofu:  []string{">= 1.6.0"},
            Providers: map[string]string{
                "aws":   ">= 5.0",
                "azure": ">= 3.0",
                "gcp":   ">= 4.0",
            },
        },
        ApprovedModules: []models.ApprovedModule{
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
    kb.securityPolicies = []models.SecurityPolicy{
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