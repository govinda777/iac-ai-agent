package unit

import (
	"encoding/json"
	"testing"

	"github.com/gosouza/iac-ai-agent/internal/agent/analyzer"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/internal/services"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestValidation(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Validation Suite")
}

var _ = ginkgo.Describe("Validação de Resultados Pré-existentes", func() {
	var (
		log             *logger.Logger
		analysisService *services.AnalysisService
		checkovAnalyzer *analyzer.CheckovAnalyzer
	)

	ginkgo.BeforeEach(func() {
		log = logger.New("debug", "json")
		analysisService = services.NewAnalysisService(log, 70)
		checkovAnalyzer = analyzer.NewCheckovAnalyzer(log)
	})

	ginkgo.Describe("ValidateAndParseResult do CheckovAnalyzer", func() {
		ginkgo.Context("quando fornecido um resultado Checkov válido", func() {
			ginkgo.It("deve validar e converter para SecurityAnalysis", func() {
				// Resultado Checkov de exemplo
				checkovResult := models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  10,
						Failed:  3,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{
								CheckID:       "CKV_AWS_1",
								CheckName:     "Ensure S3 bucket has encryption enabled",
								Description:   "S3 bucket should have encryption enabled",
								Severity:      "HIGH",
								Resource:      "aws_s3_bucket.example",
								File:          "main.tf",
								FileLineRange: []int{10, 15},
								Guideline:     "https://docs.bridgecrew.io/docs/s3_14",
							},
							{
								CheckID:       "CKV_AWS_2",
								CheckName:     "Ensure S3 bucket has logging enabled",
								Description:   "S3 bucket should have logging enabled",
								Severity:      "MEDIUM",
								Resource:      "aws_s3_bucket.example",
								File:          "main.tf",
								FileLineRange: []int{10, 15},
								Guideline:     "https://docs.bridgecrew.io/docs/s3_13",
							},
							{
								CheckID:       "CKV_AWS_3",
								CheckName:     "Ensure S3 bucket has versioning enabled",
								Description:   "S3 bucket should have versioning enabled",
								Severity:      "LOW",
								Resource:      "aws_s3_bucket.example",
								File:          "main.tf",
								FileLineRange: []int{10, 15},
								Guideline:     "https://docs.bridgecrew.io/docs/s3_16",
							},
						},
					},
				}

				// Converte para JSON
				checkovJSON, err := json.Marshal(checkovResult)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				// Valida o resultado
				securityAnalysis, err := checkovAnalyzer.ValidateAndParseResult(checkovJSON)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(securityAnalysis).ToNot(gomega.BeNil())

				// Verifica os resultados
				gomega.Expect(securityAnalysis.ChecksPassed).To(gomega.Equal(10))
				gomega.Expect(securityAnalysis.ChecksFailed).To(gomega.Equal(3))
				gomega.Expect(securityAnalysis.TotalIssues).To(gomega.Equal(3))
				gomega.Expect(len(securityAnalysis.Findings)).To(gomega.Equal(3))

				// Verifica contagem por severidade
				gomega.Expect(securityAnalysis.High).To(gomega.Equal(1))
				gomega.Expect(securityAnalysis.Medium).To(gomega.Equal(1))
				gomega.Expect(securityAnalysis.Low).To(gomega.Equal(1))
			})
		})

		ginkgo.Context("quando fornecido um JSON inválido", func() {
			ginkgo.It("deve retornar erro", func() {
				invalidJSON := []byte(`{"invalid": "json"`)

				_, err := checkovAnalyzer.ValidateAndParseResult(invalidJSON)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("quando o resultado Checkov tem dados inválidos", func() {
			ginkgo.It("deve retornar erro para checks sem ID", func() {
				checkovResult := models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed: 10,
						Failed: 1,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{
								CheckID:   "", // ID vazio - inválido
								CheckName: "Some check",
							},
						},
					},
				}

				checkovJSON, _ := json.Marshal(checkovResult)

				_, err := checkovAnalyzer.ValidateAndParseResult(checkovJSON)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("ValidatePreExistingResults do AnalysisService", func() {
		ginkgo.Context("quando fornecido resultados válidos", func() {
			ginkgo.It("deve criar AnalysisResponse completa sem executar ferramentas", func() {
				// Resultado Checkov de exemplo
				checkovResult := models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed: 15,
						Failed: 2,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{
								CheckID:   "CKV_AWS_1",
								CheckName: "Security check 1",
								Severity:  "HIGH",
								Resource:  "aws_s3_bucket.test",
								File:      "main.tf",
							},
							{
								CheckID:   "CKV_AWS_2",
								CheckName: "Security check 2",
								Severity:  "MEDIUM",
								Resource:  "aws_s3_bucket.test",
								File:      "main.tf",
							},
						},
					},
				}
				checkovJSON, _ := json.Marshal(checkovResult)

				// Análise Terraform de exemplo
				tfAnalysis := &models.TerraformAnalysis{
					TotalResources: 3,
					TotalModules:   1,
					TotalVariables: 2,
					TotalOutputs:   1,
					Resources: []models.TerraformResource{
						{Type: "aws_s3_bucket", Name: "test"},
						{Type: "aws_s3_bucket_policy", Name: "test"},
						{Type: "aws_iam_role", Name: "test"},
					},
					Variables: []models.TerraformVariable{
						{Name: "bucket_name", Description: "Name of the S3 bucket"},
						{Name: "region", Description: "AWS region"},
					},
					Outputs: []models.TerraformOutput{
						{Name: "bucket_id", Description: "ID of the S3 bucket"},
					},
					BestPracticeWarnings: []string{},
					SyntaxErrors:         []models.SyntaxError{},
				}

				// Valida os resultados pré-existentes
				response, err := analysisService.ValidatePreExistingResults(checkovJSON, tfAnalysis)

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(response).ToNot(gomega.BeNil())
				gomega.Expect(response.ID).ToNot(gomega.BeEmpty())

				// Verifica análise de segurança
				gomega.Expect(response.Analysis.Security.ChecksPassed).To(gomega.Equal(15))
				gomega.Expect(response.Analysis.Security.ChecksFailed).To(gomega.Equal(2))
				gomega.Expect(response.Analysis.Security.TotalIssues).To(gomega.Equal(2))

				// Verifica análise Terraform
				gomega.Expect(response.Analysis.Terraform.TotalResources).To(gomega.Equal(3))

				// Verifica metadata
				gomega.Expect(response.Metadata["validation_mode"]).To(gomega.Equal("pre_existing_results"))
				gomega.Expect(response.Score).To(gomega.BeNumerically(">=", 0))
				gomega.Expect(response.Score).To(gomega.BeNumerically("<=", 100))
			})
		})

		ginkgo.Context("quando não fornecido resultados Checkov", func() {
			ginkgo.It("deve processar apenas análise Terraform", func() {
				tfAnalysis := &models.TerraformAnalysis{
					TotalResources: 5,
					TotalModules:   0,
					Resources: []models.TerraformResource{
						{Type: "aws_instance", Name: "web"},
					},
				}

				response, err := analysisService.ValidatePreExistingResults(nil, tfAnalysis)

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(response).ToNot(gomega.BeNil())
				gomega.Expect(response.Analysis.Security.TotalIssues).To(gomega.Equal(0))
				gomega.Expect(response.Analysis.Terraform.TotalResources).To(gomega.Equal(5))
			})
		})
	})
})
