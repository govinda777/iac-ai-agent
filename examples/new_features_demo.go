package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/govinda777/iac-ai-agent/internal/agent/analyzer"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

func demoNewFeatures() {
	// Inicializa logger
	logger := logger.New("info", "json")

	// Exemplo 1: Preview Analyzer
	fmt.Println("=== Preview Analyzer Example ===")
	previewAnalyzer := analyzer.NewPreviewAnalyzer(logger)

	// Simula um terraform plan JSON
	terraformPlanJSON := `{
		"format_version": "1.0",
		"terraform_version": "1.5.0",
		"resource_changes": [
			{
				"address": "aws_s3_bucket.example",
				"mode": "managed",
				"type": "aws_s3_bucket",
				"name": "example",
				"change": {
					"actions": ["create"],
					"before": null,
					"after": {
						"bucket": "my-example-bucket",
						"versioning": {
							"enabled": true
						}
					}
				}
			},
			{
				"address": "aws_rds_instance.database",
				"mode": "managed",
				"type": "aws_rds_instance",
				"name": "database",
				"change": {
					"actions": ["destroy"],
					"before": {
						"identifier": "my-database",
						"engine": "postgres"
					},
					"after": null
				}
			}
		]
	}`

	// Analisa o preview
	previewAnalysis, err := previewAnalyzer.AnalyzePreview([]byte(terraformPlanJSON))
	if err != nil {
		log.Fatalf("Erro ao analisar preview: %v", err)
	}

	// Exibe resultados
	fmt.Printf("Recursos afetados: %d\n", previewAnalysis.ResourcesAffected)
	fmt.Printf("Criar: %d, Atualizar: %d, Destruir: %d, Substituir: %d\n",
		previewAnalysis.CreateCount,
		previewAnalysis.UpdateCount,
		previewAnalysis.DestroyCount,
		previewAnalysis.ReplaceCount)
	fmt.Printf("Nível de risco: %s\n", previewAnalysis.RiskLevel)
	fmt.Printf("Tempo estimado: %s\n", previewAnalysis.EstimatedApplyTime)

	// Exibe operações perigosas
	if len(previewAnalysis.DangerousOperations) > 0 {
		fmt.Println("\nOperações perigosas detectadas:")
		for _, op := range previewAnalysis.DangerousOperations {
			fmt.Printf("- %s (%s): %s\n", op.Resource, op.Action, op.Reason)
			fmt.Printf("  Mitigação: %v\n", op.Mitigation)
		}
	}

	// Exemplo 2: Secrets Scanner
	fmt.Println("\n=== Secrets Scanner Example ===")
	secretsAnalyzer := analyzer.NewSecretsAnalyzer(logger)

	// Simula código Terraform com secrets
	terraformCode := `
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"
  
  # SECRET DETECTADO: AWS Access Key
  access_key = "AKIA1234567890ABCDEF"
  
  # SECRET DETECTADO: Password
  password = "super_secret_password123"
  
  # SECRET DETECTADO: API Key
  api_key = "sk-1234567890abcdef1234567890abcdef"
  
  tags = {
    Name = "example-instance"
  }
}

# SECRET DETECTADO: Private Key
resource "aws_key_pair" "example" {
  key_name   = "example-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC..."
}

# SECRET DETECTADO: Database Password
resource "aws_db_instance" "example" {
  engine         = "postgres"
  instance_class = "db.t3.micro"
  db_password    = "my_secret_db_password"
}
`

	// Escaneia o código em busca de secrets
	secretFindings := secretsAnalyzer.ScanContent(terraformCode, "main.tf")

	// Exibe resultados
	fmt.Printf("Secrets detectados: %d\n", len(secretFindings))

	if len(secretFindings) > 0 {
		fmt.Println("\nDetalhes dos secrets:")
		for _, finding := range secretFindings {
			fmt.Printf("- Tipo: %s (Severidade: %s)\n", finding.Type, finding.Severity)
			fmt.Printf("  Arquivo: %s, Linha: %d\n", finding.File, finding.Line)
			fmt.Printf("  Descrição: %s\n", finding.Description)
			fmt.Printf("  Sugestão: %s\n", finding.Suggestion)
			fmt.Printf("  Valor mascarado: %s\n", finding.MaskedValue)
			fmt.Println()
		}

		// Resumo por severidade
		summary := secretsAnalyzer.GetSeveritySummary(secretFindings)
		fmt.Println("Resumo por severidade:")
		for severity, count := range summary {
			if count > 0 {
				fmt.Printf("- %s: %d\n", severity, count)
			}
		}
	}

	// Exemplo 3: Demonstração de JSON output
	fmt.Println("\n=== JSON Output Example ===")

	// Converte análise para JSON
	previewJSON, err := json.MarshalIndent(previewAnalysis, "", "  ")
	if err != nil {
		log.Fatalf("Erro ao converter para JSON: %v", err)
	}

	fmt.Println("Preview Analysis JSON:")
	fmt.Println(string(previewJSON))

	// Converte findings para JSON
	secretsJSON, err := json.MarshalIndent(secretFindings, "", "  ")
	if err != nil {
		log.Fatalf("Erro ao converter secrets para JSON: %v", err)
	}

	fmt.Println("\nSecrets Findings JSON:")
	fmt.Println(string(secretsJSON))

	fmt.Println("\n=== Exemplo concluído ===")
	fmt.Println("As funcionalidades Preview Analyzer e Secrets Scanner foram implementadas com sucesso!")
	fmt.Println("Agora o agente pode:")
	fmt.Println("- Analisar terraform plans e detectar operações perigosas")
	fmt.Println("- Escanear código em busca de secrets e credenciais")
	fmt.Println("- Integrar com LLM para sugestões inteligentes")
	fmt.Println("- Usar Knowledge Base para práticas recomendadas")
}

func main() {
	demoNewFeatures()
}
