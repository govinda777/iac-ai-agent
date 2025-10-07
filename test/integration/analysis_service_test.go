package integration_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gosouza/iac-ai-agent/internal/agent/analyzer"
	"github.com/gosouza/iac-ai-agent/internal/agent/llm"
	"github.com/gosouza/iac-ai-agent/internal/agent/scorer"
	"github.com/gosouza/iac-ai-agent/internal/agent/suggester"
	"github.com/gosouza/iac-ai-agent/internal/platform/cloudcontroller"
	"github.com/gosouza/iac-ai-agent/internal/services"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
	"github.com/gosouza/iac-ai-agent/test/mocks"
)

var _ = Describe("AnalysisService Integration", func() {
	var (
		analysisService *services.AnalysisService
		log             *logger.Logger
		tempDir         string
	)

	BeforeEach(func() {
		log = logger.New("debug", "text")

		// Instantiate concrete analyzers and the mock
		tfAnalyzer := analyzer.NewTerraformAnalyzer()
		mockCheckovAnalyzer := &mocks.MockCheckovAnalyzer{}
		iamAnalyzer := analyzer.NewIAMAnalyzer(log)
		prScorer := scorer.NewPRScorer()
		costOptimizer := suggester.NewCostOptimizer(log)
		securityAdvisor := suggester.NewSecurityAdvisor(log)

		// Setup the mock to always be "available"
		mockCheckovAnalyzer.IsAvailableFunc = func() bool {
			return true
		}

		cfg := &config.Config{}
		previewAnalyzer := analyzer.NewPreviewAnalyzer(log)
		secretsAnalyzer := analyzer.NewSecretsAnalyzer(log)
		llmClient := llm.NewClient(cfg, log)
		knowledgeBase := cloudcontroller.NewKnowledgeBase()
		moduleRegistry := cloudcontroller.NewModuleRegistry()
		promptBuilder := llm.NewPromptBuilder(log)

		analysisService = services.NewAnalysisService(
			log,
			70, // minPassScore
			tfAnalyzer,
			mockCheckovAnalyzer,
			iamAnalyzer,
			prScorer,
			costOptimizer,
			securityAdvisor,
			previewAnalyzer,
			secretsAnalyzer,
			llmClient,
			knowledgeBase,
			moduleRegistry,
			promptBuilder,
		)

		var err error
		tempDir, err = os.MkdirTemp("", "analysis-integration-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("Analisando conteúdo Terraform completo", func() {
		Context("quando analisa código Terraform válido com recursos AWS", func() {
			It("deve retornar análise completa com score", func() {
				content := `
resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
    Environment = "Production"
  }
}

resource "aws_s3_bucket" "data" {
  bucket = "my-data-bucket"
  
  tags = {
    Name = "DataBucket"
  }
}

variable "environment" {
  type        = string
  description = "Environment name"
  default     = "dev"
}

output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.web.id
}
`
				response, err := analysisService.AnalyzeContent(content, "main.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())
				Expect(response.ID).NotTo(BeEmpty())
				Expect(response.Score).To(BeNumerically(">=", 0))
				Expect(response.Score).To(BeNumerically("<=", 100))
				Expect(response.Analysis.Terraform.Valid).To(BeTrue())
				Expect(response.Analysis.Terraform.TotalResources).To(Equal(2))
				Expect(response.Analysis.Terraform.TotalVariables).To(Equal(1))
				Expect(response.Analysis.Terraform.TotalOutputs).To(Equal(1))
			})

			It("deve incluir metadados com score detalhado", func() {
				content := `
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  
  tags = {
    Name = "MainVPC"
  }
}
`
				err := os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(content), 0644)
				Expect(err).NotTo(HaveOccurred())

				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Metadata).To(HaveKey("pr_score"))
				Expect(response.Metadata).To(HaveKey("is_approved"))
			})
		})

		Context("quando analisa código com problemas de segurança IAM", func() {
			It("deve detectar políticas permissivas e gerar sugestões", func() {
				content := `
resource "aws_iam_policy" "admin_policy" {
  name        = "admin_policy"
  description = "Admin policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "*"
        Resource = "*"
      }
    ]
  })
}
`
				err := os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(content), 0644)
				Expect(err).NotTo(HaveOccurred())

				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Analysis.IAM.AdminAccessDetected).To(BeTrue())
				Expect(response.Analysis.IAM.OverlyPermissive).To(BeTrue())
				Expect(response.Suggestions).NotTo(BeEmpty())

				// Deve ter sugestão de segurança
				hasSecuritySuggestion := false
				for _, sugg := range response.Suggestions {
					if sugg.Type == "security" {
						hasSecuritySuggestion = true
						break
					}
				}
				Expect(hasSecuritySuggestion).To(BeTrue())
			})

			It("deve detectar recursos com acesso público", func() {
				content := `
resource "aws_s3_bucket" "public_bucket" {
  bucket = "public-bucket"
  acl    = "public-read"
}

resource "aws_db_instance" "public_db" {
  identifier           = "mydb"
  allocated_storage    = 20
  engine               = "postgres"
  instance_class       = "db.t2.micro"
  publicly_accessible  = true
}
`
				err := os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(content), 0644)
				Expect(err).NotTo(HaveOccurred())

				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Analysis.IAM.PublicAccess).NotTo(BeEmpty())
				Expect(len(response.Analysis.IAM.PublicAccess)).To(Equal(2))
			})
		})

		Context("quando analisa código com problemas de best practices", func() {
			It("deve gerar warnings e reduzir score", func() {
				content := `
resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

variable "region" {
  type = string
}
`
				response, err := analysisService.AnalyzeContent(content, "main.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Analysis.Terraform.BestPracticeWarnings).NotTo(BeEmpty())

				// Deve ter warnings sobre falta de tags e descrição
				warnings := response.Analysis.Terraform.BestPracticeWarnings
				hasTagWarning := false
				hasDescriptionWarning := false

				for _, w := range warnings {
					if Contains(w, "tags") {
						hasTagWarning = true
					}
					if Contains(w, "descrição") {
						hasDescriptionWarning = true
					}
				}

				Expect(hasTagWarning || hasDescriptionWarning).To(BeTrue())
			})
		})

		Context("quando analisa código com oportunidades de otimização de custo", func() {
			It("deve gerar sugestões de custo", func() {
				content := `
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.xlarge"
}

resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public.id
}
`
				response, err := analysisService.AnalyzeContent(content, "main.tf")

				Expect(err).NotTo(HaveOccurred())

				// Deve ter sugestões de custo
				hasCostSuggestion := false
				for _, sugg := range response.Suggestions {
					if sugg.Type == "cost" {
						hasCostSuggestion = true
						Expect(sugg.EstimatedSavings).NotTo(BeEmpty())
						break
					}
				}
				Expect(hasCostSuggestion).To(BeTrue())
			})
		})

		Context("quando analisa código inválido", func() {
			It("deve retornar análise com erros de sintaxe", func() {
				content := `
resource "aws_instance" "web" {
  ami = "ami-12345678
  instance_type = "t2.micro"
}
`
				response, err := analysisService.AnalyzeContent(content, "invalid.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Analysis.Terraform.Valid).To(BeFalse())
				Expect(response.Analysis.Terraform.SyntaxErrors).NotTo(BeEmpty())
			})
		})
	})

	Describe("Analisando diretório completo", func() {
		Context("quando analisa projeto Terraform multi-arquivo", func() {
			BeforeEach(func() {
				// Cria estrutura de projeto Terraform
				mainTf := `
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  
  tags = {
    Name = "MainVPC"
  }
}

resource "aws_subnet" "public" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
  
  tags = {
    Name = "PublicSubnet"
  }
}
`
				variablesTf := `
variable "environment" {
  type        = string
  description = "Environment name"
  default     = "production"
}

variable "region" {
  type        = string
  description = "AWS region"
  default     = "us-east-1"
}
`
				outputsTf := `
output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "subnet_id" {
  description = "Subnet ID"
  value       = aws_subnet.public.id
}
`
				err := os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(mainTf), 0644)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "variables.tf"), []byte(variablesTf), 0644)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "outputs.tf"), []byte(outputsTf), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			It("deve analisar todos os arquivos e consolidar resultados", func() {
				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response).NotTo(BeNil())
				Expect(response.Analysis.Terraform.TotalResources).To(Equal(2))
				Expect(response.Analysis.Terraform.TotalVariables).To(Equal(2))
				Expect(response.Analysis.Terraform.TotalOutputs).To(Equal(2))
				Expect(response.Analysis.Terraform.Valid).To(BeTrue())
			})

			It("deve incluir análise de custos", func() {
				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Analysis.Cost).NotTo(BeNil())
			})

			It("deve calcular score baseado em múltiplos fatores", func() {
				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Score).To(BeNumerically(">", 0))

				// Score deve refletir boa qualidade (tem descrições, outputs, tags)
				Expect(response.Score).To(BeNumerically(">=", 70))
			})
		})

		Context("quando diretório contém mix de arquivos válidos e inválidos", func() {
			BeforeEach(func() {
				validTf := `
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
`
				invalidTf := `
resource "aws_instance" "web" {
  ami = "ami-12345678
}
`
				err := os.WriteFile(filepath.Join(tempDir, "valid.tf"), []byte(validTf), 0644)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "invalid.tf"), []byte(invalidTf), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			It("deve processar arquivos válidos e reportar erros", func() {
				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Analysis.Terraform.Valid).To(BeFalse())
				Expect(response.Analysis.Terraform.SyntaxErrors).NotTo(BeEmpty())
				Expect(response.Analysis.Terraform.TotalResources).To(Equal(1))
			})
		})

		Context("quando diretório está vazio", func() {
			It("deve retornar análise sem recursos", func() {
				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(response.Analysis.Terraform.TotalResources).To(Equal(0))
				Expect(response.Analysis.Terraform.Valid).To(BeTrue())
			})
		})
	})

	Describe("Validando análises", func() {
		Context("quando análise é válida", func() {
			It("deve passar na validação", func() {
				content := `resource "aws_vpc" "main" { cidr_block = "10.0.0.0/16" }`
				response, err := analysisService.AnalyzeContent(content, "main.tf")

				Expect(err).NotTo(HaveOccurred())

				validationErr := analysisService.ValidateAnalysis(response)
				Expect(validationErr).NotTo(HaveOccurred())
			})
		})

		Context("quando análise é nula", func() {
			It("deve falhar na validação", func() {
				err := analysisService.ValidateAnalysis(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nula"))
			})
		})
	})

	Describe("Fluxo completo end-to-end", func() {
		Context("quando analisa projeto real com múltiplos componentes", func() {
			BeforeEach(func() {
				// Projeto completo com VPC, EC2, S3, IAM
				mainTf := `
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  
  tags = {
    Name        = "MainVPC"
    Environment = var.environment
  }
}

resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.public.id
  
  tags = {
    Name = "WebServer"
  }
}

resource "aws_subnet" "public" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
  
  tags = {
    Name = "PublicSubnet"
  }
}

resource "aws_s3_bucket" "data" {
  bucket = "my-data-bucket"
  
  tags = {
    Name = "DataBucket"
  }
}
`
				iamTf := `
resource "aws_iam_role" "ec2_role" {
  name = "ec2_role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "s3_read_policy" {
  name = "s3_read_policy"
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:ListBucket"
      ]
      Resource = "*"
    }]
  })
}
`
				variablesTf := `
variable "environment" {
  type        = string
  description = "Environment name"
  default     = "production"
}
`
				outputsTf := `
output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "instance_id" {
  description = "EC2 Instance ID"
  value       = aws_instance.web.id
}
`
				os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(mainTf), 0644)
				os.WriteFile(filepath.Join(tempDir, "iam.tf"), []byte(iamTf), 0644)
				os.WriteFile(filepath.Join(tempDir, "variables.tf"), []byte(variablesTf), 0644)
				os.WriteFile(filepath.Join(tempDir, "outputs.tf"), []byte(outputsTf), 0644)
			})

			It("deve realizar análise completa e gerar relatório abrangente", func() {
				response, err := analysisService.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())

				// Verifica estrutura básica
				Expect(response.ID).NotTo(BeEmpty())
				Expect(response.Timestamp).NotTo(BeZero())

				// Verifica análise Terraform
				Expect(response.Analysis.Terraform.Valid).To(BeTrue())
				Expect(response.Analysis.Terraform.TotalResources).To(Equal(6))
				Expect(response.Analysis.Terraform.TotalVariables).To(Equal(1))
				Expect(response.Analysis.Terraform.TotalOutputs).To(Equal(2))

				// Verifica análise IAM
				Expect(response.Analysis.IAM.TotalPolicies).To(Equal(1))
				Expect(response.Analysis.IAM.TotalRoles).To(Equal(1))

				// Verifica sugestões foram geradas
				Expect(response.Suggestions).NotTo(BeEmpty())

				// Verifica score foi calculado
				Expect(response.Score).To(BeNumerically(">=", 0))
				Expect(response.Score).To(BeNumerically("<=", 100))

				// Verifica metadados
				Expect(response.Metadata["is_approved"]).NotTo(BeNil())
			})
		})
	})
})

// Helper function
func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}
