package unit_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gosouza/iac-ai-agent/internal/agent/analyzer"
)

var _ = Describe("TerraformAnalyzer", func() {
	var (
		tfAnalyzer *analyzer.TerraformAnalyzer
		tempDir    string
	)

	BeforeEach(func() {
		tfAnalyzer = analyzer.NewTerraformAnalyzer()
		var err error
		tempDir, err = os.MkdirTemp("", "terraform-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("Analisando código Terraform válido", func() {
		Context("quando o código contém recursos básicos", func() {
			It("deve identificar corretamente os recursos", func() {
				content := `
resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "main.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Valid).To(BeTrue())
				Expect(analysis.TotalResources).To(Equal(2))
				Expect(analysis.Resources).To(HaveLen(2))
				Expect(analysis.Resources[0].Type).To(Equal("aws_instance"))
				Expect(analysis.Resources[0].Name).To(Equal("web"))
				Expect(analysis.Resources[1].Type).To(Equal("aws_s3_bucket"))
				Expect(analysis.Resources[1].Name).To(Equal("data"))
			})

			It("deve identificar o provider correto", func() {
				content := `
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "vpc.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Resources[0].Provider).To(Equal("aws"))
			})
		})

		Context("quando o código contém módulos", func() {
			It("deve identificar módulos corretamente", func() {
				content := `
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "3.0.0"
}

module "rds" {
  source = "./modules/rds"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "modules.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.TotalModules).To(Equal(2))
				Expect(analysis.Modules).To(HaveLen(2))
				Expect(analysis.Modules[0].Name).To(Equal("vpc"))
				Expect(analysis.Modules[1].Name).To(Equal("rds"))
			})
		})

		Context("quando o código contém variáveis e outputs", func() {
			It("deve identificar variáveis corretamente", func() {
				content := `
variable "environment" {
  type        = string
  description = "Environment name"
  default     = "dev"
}

variable "instance_count" {
  type = number
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "variables.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.TotalVariables).To(Equal(2))
				Expect(analysis.Variables).To(HaveLen(2))
				Expect(analysis.Variables[0].Name).To(Equal("environment"))
				Expect(analysis.Variables[1].Name).To(Equal("instance_count"))
			})

			It("deve identificar outputs corretamente", func() {
				content := `
output "vpc_id" {
  value = aws_vpc.main.id
}

output "subnet_ids" {
  value = aws_subnet.private[*].id
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "outputs.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.TotalOutputs).To(Equal(2))
				Expect(analysis.Outputs).To(HaveLen(2))
				Expect(analysis.Outputs[0].Name).To(Equal("vpc_id"))
				Expect(analysis.Outputs[1].Name).To(Equal("subnet_ids"))
			})
		})

		Context("quando o código contém providers", func() {
			It("deve identificar providers declarados", func() {
				content := `
provider "aws" {
  region = "us-east-1"
}

provider "azurerm" {
  features {}
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "providers.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Providers).To(ContainElements("aws", "azurerm"))
			})
		})
	})

	Describe("Analisando código Terraform com erros", func() {
		Context("quando o código possui erros de sintaxe", func() {
			It("deve retornar análise inválida com detalhes do erro", func() {
				content := `
resource "aws_instance" "web" {
  ami = "ami-12345678
  instance_type = "t2.micro"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "invalid.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Valid).To(BeFalse())
				Expect(analysis.SyntaxErrors).NotTo(BeEmpty())
				Expect(analysis.SyntaxErrors[0].File).To(Equal("invalid.tf"))
			})
		})

		Context("quando o código contém blocos malformados", func() {
			It("deve capturar o erro de sintaxe", func() {
				content := `
resource "aws_instance" {
  ami = "ami-12345678"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "malformed.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Valid).To(BeFalse())
			})
		})
	})

	Describe("Verificando best practices", func() {
		Context("quando recursos não possuem tags", func() {
			It("deve gerar warnings sobre falta de tags", func() {
				content := `
resource "aws_instance" "web" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "main.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.BestPracticeWarnings).To(ContainElement(
					ContainSubstring("aws_instance.web não possui tags")))
			})
		})

		Context("quando variáveis não possuem descrição", func() {
			It("deve gerar warnings sobre falta de descrição", func() {
				content := `
variable "region" {
  type = string
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "variables.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.BestPracticeWarnings).To(ContainElement(
					ContainSubstring("Variável region não possui descrição")))
			})
		})

		Context("quando há recursos mas não há outputs", func() {
			It("deve sugerir adicionar outputs", func() {
				content := `
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "main.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.BestPracticeWarnings).To(ContainElement(
					ContainSubstring("Considere adicionar outputs")))
			})
		})
	})

	Describe("Analisando diretório completo", func() {
		Context("quando o diretório contém múltiplos arquivos Terraform", func() {
			BeforeEach(func() {
				// Cria arquivos de teste
				mainTf := `
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
`
				variablesTf := `
variable "environment" {
  type        = string
  description = "Environment name"
}
`
				outputsTf := `
output "vpc_id" {
  value = aws_vpc.main.id
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
				analysis, err := tfAnalyzer.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Valid).To(BeTrue())
				Expect(analysis.TotalResources).To(Equal(1))
				Expect(analysis.TotalVariables).To(Equal(1))
				Expect(analysis.TotalOutputs).To(Equal(1))
			})
		})

		Context("quando o diretório contém arquivos com erros", func() {
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

			It("deve processar arquivos válidos e reportar erros dos inválidos", func() {
				analysis, err := tfAnalyzer.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Valid).To(BeFalse())
				Expect(analysis.TotalResources).To(Equal(1)) // Apenas o recurso válido
				Expect(analysis.SyntaxErrors).NotTo(BeEmpty())
			})
		})

		Context("quando o diretório está vazio", func() {
			It("deve retornar análise vazia mas válida", func() {
				analysis, err := tfAnalyzer.AnalyzeDirectory(tempDir)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Valid).To(BeTrue())
				Expect(analysis.TotalResources).To(Equal(0))
				Expect(analysis.TotalModules).To(Equal(0))
			})
		})
	})

	Describe("Verificando tipos de recursos que devem ter tags", func() {
		Context("quando o recurso é um tipo que deve ter tags", func() {
			It("deve identificar aws_s3_bucket como recurso que necessita tags", func() {
				content := `
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "s3.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.BestPracticeWarnings).To(ContainElement(
					ContainSubstring("aws_s3_bucket.data não possui tags")))
			})

			It("deve identificar aws_lambda_function como recurso que necessita tags", func() {
				content := `
resource "aws_lambda_function" "processor" {
  function_name = "my-function"
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "lambda.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.BestPracticeWarnings).To(ContainElement(
					ContainSubstring("aws_lambda_function.processor não possui tags")))
			})
		})

		Context("quando o recurso não precisa de tags", func() {
			It("não deve gerar warning para tipos que não necessitam tags", func() {
				content := `
resource "random_string" "suffix" {
  length = 8
}
`
				analysis, err := tfAnalyzer.AnalyzeContent(content, "random.tf")

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.BestPracticeWarnings).NotTo(ContainElement(
					ContainSubstring("random_string.suffix não possui tags")))
			})
		})
	})
})
