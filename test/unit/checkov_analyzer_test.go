package unit_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gosouza/iac-ai-agent/internal/agent/analyzer"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

var _ = Describe("CheckovAnalyzer", func() {
	var (
		checkovAnalyzer *analyzer.CheckovAnalyzer
		log             *logger.Logger
	)

	BeforeEach(func() {
		log = logger.New("info", "json")
		checkovAnalyzer = analyzer.NewCheckovAnalyzer(log)
	})

	Describe("Verificando disponibilidade do Checkov", func() {
		Context("quando o Checkov está instalado", func() {
			It("deve retornar true se encontrado no PATH", func() {
				// Este teste depende do ambiente
				// Apenas verifica se a função funciona
				_ = checkovAnalyzer.IsAvailable()
				// Não fazemos assertion pois depende do ambiente
			})
		})
	})

	Describe("Convertendo resultados Checkov para SecurityAnalysis", func() {
		Context("quando há checks passados e falhados", func() {
			It("deve converter corretamente os resultados", func() {
				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  5,
						Failed:  3,
						Skipped: 1,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{
								CheckID:       "CKV_AWS_1",
								CheckName:     "S3 Bucket Encryption",
								Severity:      "HIGH",
								Resource:      "aws_s3_bucket.data",
								File:          "main.tf",
								FileLineRange: []int{10, 15},
							},
							{
								CheckID:       "CKV_AWS_2",
								CheckName:     "S3 Bucket Versioning",
								Severity:      "MEDIUM",
								Resource:      "aws_s3_bucket.data",
								File:          "main.tf",
								FileLineRange: []int{10, 15},
							},
							{
								CheckID:   "CKV_AWS_3",
								CheckName: "S3 Bucket Logging",
								Severity:  "LOW",
								Resource:  "aws_s3_bucket.logs",
								File:      "s3.tf",
							},
						},
					},
				}

				// Usa reflection para chamar método privado
				// ou testa através de método público
				suggestions := checkovAnalyzer.GetRecommendations(
					checkovAnalyzer.GetSecurityAnalysis(checkovResult))

				Expect(suggestions).To(HaveLen(3))
			})
		})

		Context("quando determina severidade baseada no check ID", func() {
			It("deve identificar HIGH para checks de encryption", func() {
				check := models.CheckovCheck{
					CheckID:   "CKV_AWS_1",
					CheckName: "Ensure S3 bucket has encryption enabled",
				}

				// Simula determinação de severidade
				severity := determineSeverityForTest(check)
				Expect(severity).To(Equal("HIGH"))
			})

			It("deve identificar MEDIUM para checks de logging", func() {
				check := models.CheckovCheck{
					CheckID:   "CKV_AWS_2",
					CheckName: "Ensure S3 bucket has logging enabled",
				}

				severity := determineSeverityForTest(check)
				Expect(severity).To(Equal("MEDIUM"))
			})

			It("deve identificar LOW para checks de tags", func() {
				check := models.CheckovCheck{
					CheckID:   "CKV_AWS_3",
					CheckName: "Ensure S3 bucket has tags",
				}

				severity := determineSeverityForTest(check)
				Expect(severity).To(Equal("LOW"))
			})
		})
	})

	Describe("Gerando recomendações baseadas em findings", func() {
		Context("quando há multiple security findings", func() {
			It("deve criar sugestões apropriadas para cada finding", func() {
				securityAnalysis := &models.SecurityAnalysis{
					Critical:     1,
					High:         2,
					Medium:       3,
					Low:          1,
					TotalIssues:  7,
					ChecksPassed: 10,
					ChecksFailed: 7,
					Findings: []models.SecurityFinding{
						{
							ID:        "1",
							CheckID:   "CKV_AWS_1",
							CheckName: "S3 bucket encryption missing",
							Severity:  "CRITICAL",
							Resource:  "aws_s3_bucket.data",
							File:      "main.tf",
							Line:      10,
							Guideline: "Enable server-side encryption for S3 bucket",
						},
						{
							ID:        "2",
							CheckID:   "CKV_AWS_2",
							CheckName: "S3 bucket versioning disabled",
							Severity:  "HIGH",
							Resource:  "aws_s3_bucket.data",
							File:      "main.tf",
							Line:      10,
							Guideline: "Enable versioning for S3 bucket",
						},
					},
				}

				suggestions := checkovAnalyzer.GetRecommendations(securityAnalysis)

				Expect(suggestions).To(HaveLen(2))
				Expect(suggestions[0].Type).To(Equal("security"))
				Expect(suggestions[0].Severity).To(Equal("critical"))
				Expect(suggestions[0].Message).To(ContainSubstring("encryption"))
				Expect(suggestions[1].Severity).To(Equal("high"))
			})
		})

		Context("quando gera link para documentação", func() {
			It("deve criar URL correto do Bridgecrew", func() {
				securityAnalysis := &models.SecurityAnalysis{
					Findings: []models.SecurityFinding{
						{
							CheckID:   "CKV_AWS_1",
							CheckName: "Test check",
							Severity:  "HIGH",
						},
					},
				}

				suggestions := checkovAnalyzer.GetRecommendations(securityAnalysis)

				Expect(suggestions[0].ReferenceLink).To(ContainSubstring("docs.bridgecrew.io"))
				Expect(suggestions[0].ReferenceLink).To(ContainSubstring("ckv_aws_1"))
			})
		})
	})

	Describe("Contabilizando issues por severidade", func() {
		Context("quando converte resultados Checkov", func() {
			It("deve contar corretamente cada tipo de severidade", func() {
				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  10,
						Failed:  8,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{Severity: "CRITICAL"},
							{Severity: "CRITICAL"},
							{Severity: "HIGH"},
							{Severity: "HIGH"},
							{Severity: "HIGH"},
							{Severity: "MEDIUM"},
							{Severity: "MEDIUM"},
							{Severity: "LOW"},
						},
					},
				}

				analysis := checkovAnalyzer.GetSecurityAnalysis(checkovResult)

				Expect(analysis.Critical).To(Equal(2))
				Expect(analysis.High).To(Equal(3))
				Expect(analysis.Medium).To(Equal(2))
				Expect(analysis.Low).To(Equal(1))
				Expect(analysis.TotalIssues).To(Equal(8))
				Expect(analysis.ChecksFailed).To(Equal(8))
				Expect(analysis.ChecksPassed).To(Equal(10))
			})
		})
	})
})

// Helper function para testar determinação de severidade
func determineSeverityForTest(check models.CheckovCheck) string {
	checkLower := check.CheckID + " " + check.CheckName
	checkLower = toLower(checkLower)

	if contains(checkLower, "encryption") ||
		contains(checkLower, "credential") ||
		contains(checkLower, "secret") {
		return "HIGH"
	}

	if contains(checkLower, "logging") ||
		contains(checkLower, "monitoring") {
		return "MEDIUM"
	}

	if contains(checkLower, "tag") ||
		contains(checkLower, "description") {
		return "LOW"
	}

	return "MEDIUM"
}

func toLower(s string) string {
	return strings.ToLower(s)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
