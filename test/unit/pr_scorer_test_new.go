package unit

import (
	"github.com/gosouza/iac-ai-agent/internal/agent/scorer"
	"github.com/gosouza/iac-ai-agent/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PRScorer (Nova API)", func() {
	var (
		prScorer *scorer.PRScorer
	)

	BeforeEach(func() {
		prScorer = scorer.NewPRScorer()
	})

	Describe("Criando uma nova instância do PRScorer", func() {
		Context("quando criado", func() {
			It("deve criar o scorer com sucesso", func() {
				scorer := scorer.NewPRScorer()
				Expect(scorer).NotTo(BeNil())
			})
		})
	})

	Describe("Calculando score completo", func() {
		Context("quando tudo está perfeito", func() {
			It("deve retornar score próximo de 100", func() {
				analysisDetails := &models.AnalysisDetails{
					Terraform: models.TerraformAnalysis{
						Valid:          true,
						TotalResources: 10,
						TotalModules:   2,
						TotalVariables: 5,
						TotalOutputs:   3,
						Variables: []models.TerraformVariable{
							{Name: "var1", Description: "Description 1"},
							{Name: "var2", Description: "Description 2"},
						},
						Outputs: []models.TerraformOutput{
							{Name: "out1", Description: "Output 1"},
						},
						BestPracticeWarnings: []string{},
						SyntaxErrors:         []models.SyntaxError{},
					},
					Security: models.SecurityAnalysis{
						ChecksPassed: 20,
						ChecksFailed: 0,
						TotalIssues:  0,
					},
				}

				score := prScorer.CalculateScore(analysisDetails)
				Expect(score.Total).To(BeNumerically(">=", 90))
				Expect(score.Security).To(Equal(100))
			})
		})

		Context("quando há problemas de segurança", func() {
			It("deve penalizar o score de segurança", func() {
				analysisDetails := &models.AnalysisDetails{
					Terraform: models.TerraformAnalysis{
						Valid: true,
					},
					Security: models.SecurityAnalysis{
						ChecksPassed: 10,
						ChecksFailed: 5,
						TotalIssues:  5,
						Critical:     2,
						High:         1,
						Medium:       1,
						Low:          1,
					},
				}

				score := prScorer.CalculateScore(analysisDetails)
				// Critical=2 (2*20=40) + High=1 (10) + Medium=1 (5) + Low=1 (2) = 57
				// 100 - 57 = 43
				Expect(score.Security).To(Equal(43))
				Expect(score.Total).To(BeNumerically("<", 100))
			})
		})

		Context("quando há muitos warnings de best practices", func() {
			It("deve penalizar o score de best practices", func() {
				analysisDetails := &models.AnalysisDetails{
					Terraform: models.TerraformAnalysis{
						Valid: true,
						BestPracticeWarnings: []string{
							"warning1", "warning2", "warning3",
							"warning4", "warning5",
						},
						SyntaxErrors: []models.SyntaxError{},
					},
					Security: models.SecurityAnalysis{},
				}

				score := prScorer.CalculateScore(analysisDetails)
				// 100 - (5 * 3) = 85
				Expect(score.BestPractices).To(Equal(85))
			})
		})

		Context("quando falta documentação", func() {
			It("deve penalizar o score de documentação", func() {
				analysisDetails := &models.AnalysisDetails{
					Terraform: models.TerraformAnalysis{
						Valid:          true,
						TotalVariables: 3,
						TotalOutputs:   2,
						Variables: []models.TerraformVariable{
							{Name: "var1", Description: ""},
							{Name: "var2", Description: ""},
							{Name: "var3", Description: "Documented"},
						},
						Outputs: []models.TerraformOutput{
							{Name: "out1", Description: ""},
							{Name: "out2", Description: "Documented"},
						},
					},
					Security: models.SecurityAnalysis{},
				}

				score := prScorer.CalculateScore(analysisDetails)
				// 1/3 vars documented = 33% of 50 = 16.5 ≈ 16
				// 1/2 outputs documented = 50% of 50 = 25
				// Total: 16 + 25 = 41
				Expect(score.Documentation).To(BeNumerically(">=", 40))
				Expect(score.Documentation).To(BeNumerically("<=", 42))
			})
		})

		Context("quando há muitos recursos sem modularização", func() {
			It("deve penalizar maintainability e performance", func() {
				analysisDetails := &models.AnalysisDetails{
					Terraform: models.TerraformAnalysis{
						Valid:          true,
						TotalResources: 60,
						TotalModules:   0,
					},
					Security: models.SecurityAnalysis{},
				}

				score := prScorer.CalculateScore(analysisDetails)
				// Performance: 100 - 10 (muitos recursos) - 15 (sem módulos) = 75
				Expect(score.Performance).To(BeNumerically("<=", 80))
				// Maintainability: também penalizado
				Expect(score.Maintainability).To(BeNumerically("<", 100))
			})
		})
	})

	Describe("GetScoreLevel", func() {
		Context("para diferentes scores", func() {
			It("deve retornar o nível correto", func() {
				Expect(prScorer.GetScoreLevel(95)).To(Equal("Excelente"))
				Expect(prScorer.GetScoreLevel(80)).To(Equal("Bom"))
				Expect(prScorer.GetScoreLevel(65)).To(Equal("Regular"))
				Expect(prScorer.GetScoreLevel(45)).To(Equal("Ruim"))
				Expect(prScorer.GetScoreLevel(30)).To(Equal("Crítico"))
			})
		})
	})

	Describe("ShouldApprove", func() {
		Context("quando score é alto e segurança é boa", func() {
			It("deve aprovar", func() {
				score := &models.PRScore{
					Total:    85,
					Security: 90,
				}
				Expect(prScorer.ShouldApprove(score, 70)).To(BeTrue())
			})
		})

		Context("quando score está abaixo do mínimo", func() {
			It("não deve aprovar", func() {
				score := &models.PRScore{
					Total:    65,
					Security: 90,
				}
				Expect(prScorer.ShouldApprove(score, 70)).To(BeFalse())
			})
		})

		Context("quando segurança é crítica", func() {
			It("não deve aprovar mesmo com score alto", func() {
				score := &models.PRScore{
					Total:    90,
					Security: 45,
				}
				Expect(prScorer.ShouldApprove(score, 70)).To(BeFalse())
			})
		})
	})

	Describe("GenerateScoreSummary", func() {
		Context("quando gerado", func() {
			It("deve conter as informações principais", func() {
				score := &models.PRScore{
					Total:           85,
					Security:        90,
					BestPractices:   85,
					Performance:     80,
					Maintainability: 85,
					Documentation:   75,
				}

				summary := prScorer.GenerateScoreSummary(score)
				Expect(summary).To(ContainSubstring("Score de Qualidade"))
				Expect(summary).To(ContainSubstring("Segurança"))
				Expect(summary).To(ContainSubstring("Best Practices"))
			})
		})
	})
})
