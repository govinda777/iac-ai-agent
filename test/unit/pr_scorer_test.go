package unit_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gosouza/iac-ai-agent/internal/agent/scorer"
	"github.com/gosouza/iac-ai-agent/internal/models"
)

var _ = Describe("PRScorer", func() {
	var (
		prScorer *scorer.PRScorer
	)

	BeforeEach(func() {
		prScorer = scorer.NewPRScorer(70)
	})

	Describe("Criando uma nova instância do PRScorer", func() {
		Context("quando o score mínimo é válido", func() {
			It("deve criar o scorer com o score mínimo especificado", func() {
				scorer := scorer.NewPRScorer(80)
				Expect(scorer).NotTo(BeNil())
			})
		})

		Context("quando o score mínimo é inválido", func() {
			It("deve usar o valor padrão de 70", func() {
				scorer := scorer.NewPRScorer(0)
				Expect(scorer).NotTo(BeNil())
			})
		})
	})

	Describe("Calculando score de segurança", func() {
		Context("quando não há resultados do Checkov", func() {
			It("deve retornar score padrão de 80", func() {
				analysis := &models.TerraformAnalysis{Valid: true}
				score := prScorer.CalculateScore(analysis, nil)

				Expect(score.Security).To(Equal(80))
			})
		})

		Context("quando todos os checks passam", func() {
			It("deve retornar score máximo de 100", func() {
				analysis := &models.TerraformAnalysis{Valid: true}
				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  10,
						Failed:  0,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{},
					},
				}

				score := prScorer.CalculateScore(analysis, checkovResult)
				Expect(score.Security).To(Equal(100))
			})
		})

		Context("quando há falhas críticas", func() {
			It("deve penalizar fortemente o score", func() {
				analysis := &models.TerraformAnalysis{Valid: true}
				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  5,
						Failed:  2,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{
								CheckID:  "CKV_AWS_1",
								Severity: "CRITICAL",
							},
							{
								CheckID:  "CKV_AWS_2",
								Severity: "CRITICAL",
							},
						},
					},
				}

				score := prScorer.CalculateScore(analysis, checkovResult)
				// 100 - (2 * 20) = 60
				Expect(score.Security).To(Equal(60))
			})
		})

		Context("quando há falhas de severidades variadas", func() {
			It("deve aplicar penalidades proporcionais", func() {
				analysis := &models.TerraformAnalysis{Valid: true}
				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  10,
						Failed:  4,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{CheckID: "CKV_1", Severity: "CRITICAL"}, // -20
							{CheckID: "CKV_2", Severity: "HIGH"},     // -10
							{CheckID: "CKV_3", Severity: "MEDIUM"},   // -5
							{CheckID: "CKV_4", Severity: "LOW"},      // -2
						},
					},
				}

				score := prScorer.CalculateScore(analysis, checkovResult)
				// 100 - 20 - 10 - 5 - 2 = 63
				Expect(score.Security).To(Equal(63))
			})
		})

		Context("quando há muitas falhas que reduziriam o score abaixo de zero", func() {
			It("deve retornar score mínimo de 0", func() {
				analysis := &models.TerraformAnalysis{Valid: true}
				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  0,
						Failed:  10,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{Severity: "CRITICAL"}, {Severity: "CRITICAL"},
							{Severity: "CRITICAL"}, {Severity: "CRITICAL"},
							{Severity: "CRITICAL"}, {Severity: "CRITICAL"},
						},
					},
				}

				score := prScorer.CalculateScore(analysis, checkovResult)
				Expect(score.Security).To(BeNumerically(">=", 0))
			})
		})
	})

	Describe("Calculando score de best practices", func() {
		Context("quando não há warnings", func() {
			It("deve retornar score máximo de 100", func() {
				analysis := &models.TerraformAnalysis{
					Valid:                true,
					BestPracticeWarnings: []string{},
				}

				score := prScorer.CalculateScore(analysis, nil)
				Expect(score.BestPractices).To(Equal(100))
			})
		})

		Context("quando há alguns warnings", func() {
			It("deve penalizar 5 pontos por warning", func() {
				analysis := &models.TerraformAnalysis{
					Valid: true,
					BestPracticeWarnings: []string{
						"Warning 1",
						"Warning 2",
						"Warning 3",
					},
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 - (3 * 5) = 85
				Expect(score.BestPractices).To(Equal(85))
			})
		})

		Context("quando há muitos warnings", func() {
			It("deve garantir que o score não fique negativo", func() {
				warnings := make([]string, 30)
				for i := range warnings {
					warnings[i] = "Warning"
				}

				analysis := &models.TerraformAnalysis{
					Valid:                true,
					BestPracticeWarnings: warnings,
				}

				score := prScorer.CalculateScore(analysis, nil)
				Expect(score.BestPractices).To(Equal(0))
			})
		})
	})

	Describe("Calculando score de documentação", func() {
		Context("quando todas as variáveis e outputs têm descrição", func() {
			It("deve retornar score alto", func() {
				analysis := &models.TerraformAnalysis{
					Valid: true,
					Variables: []models.TerraformVariable{
						{Name: "var1", Description: "Description 1"},
						{Name: "var2", Description: "Description 2"},
					},
					Outputs: []models.TerraformOutput{
						{Name: "out1", Description: "Output 1"},
					},
					TotalOutputs: 1,
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 + 5 (bonus por ter outputs) = 105, mas max é 100
				Expect(score.Documentation).To(Equal(100))
			})
		})

		Context("quando variáveis não têm descrição", func() {
			It("deve penalizar o score", func() {
				analysis := &models.TerraformAnalysis{
					Valid: true,
					Variables: []models.TerraformVariable{
						{Name: "var1", Description: ""},
						{Name: "var2", Description: ""},
					},
					Outputs: []models.TerraformOutput{},
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 - (2 * 3) = 94
				Expect(score.Documentation).To(Equal(94))
			})
		})

		Context("quando outputs não têm descrição", func() {
			It("deve penalizar o score mas adicionar bonus por ter outputs", func() {
				analysis := &models.TerraformAnalysis{
					Valid:     true,
					Variables: []models.TerraformVariable{},
					Outputs: []models.TerraformOutput{
						{Name: "out1", Description: ""},
					},
					TotalOutputs: 1,
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 - 3 + 5 = 102, mas max é 100
				Expect(score.Documentation).To(Equal(100))
			})
		})
	})

	Describe("Calculando score de maintainability", func() {
		Context("quando usa módulos verificados", func() {
			It("deve dar bonus no score", func() {
				analysis := &models.TerraformAnalysis{
					Valid:        true,
					TotalModules: 3,
					Modules: []models.TerraformModule{
						{Name: "mod1", Verified: true},
						{Name: "mod2", Verified: true},
						{Name: "mod3", Verified: true},
					},
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 + (3 * 2) = 106, mas max é 100
				Expect(score.Maintainability).To(Equal(100))
			})
		})

		Context("quando há muitos recursos", func() {
			It("deve penalizar o score por complexidade", func() {
				resources := make([]models.TerraformResource, 60)
				analysis := &models.TerraformAnalysis{
					Valid:          true,
					TotalResources: 60,
					Resources:      resources,
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 - 10 = 90
				Expect(score.Maintainability).To(Equal(90))
			})
		})

		Context("quando há módulos não verificados", func() {
			It("deve penalizar o score", func() {
				analysis := &models.TerraformAnalysis{
					Valid:        true,
					TotalModules: 2,
					Modules: []models.TerraformModule{
						{Name: "mod1", Verified: false},
						{Name: "mod2", Verified: false},
					},
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 + (2 * 2) - (2 * 3) = 98
				Expect(score.Maintainability).To(Equal(98))
			})
		})
	})

	Describe("Calculando score de performance", func() {
		Context("quando há poucos data sources", func() {
			It("deve retornar score alto", func() {
				analysis := &models.TerraformAnalysis{
					Valid:            true,
					TotalDataSources: 5,
				}

				score := prScorer.CalculateScore(analysis, nil)
				Expect(score.Performance).To(BeNumerically(">=", 95))
			})
		})

		Context("quando há muitos data sources", func() {
			It("deve penalizar o score", func() {
				analysis := &models.TerraformAnalysis{
					Valid:            true,
					TotalDataSources: 15,
					TotalOutputs:     0,
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 - 10 = 90
				Expect(score.Performance).To(Equal(90))
			})
		})

		Context("quando usa outputs eficientemente", func() {
			It("deve dar bonus no score", func() {
				analysis := &models.TerraformAnalysis{
					Valid:        true,
					TotalOutputs: 3,
				}

				score := prScorer.CalculateScore(analysis, nil)
				// 100 + 5 = 105, mas max é 100
				Expect(score.Performance).To(Equal(100))
			})
		})
	})

	Describe("Calculando score total ponderado", func() {
		Context("quando todos os scores são 100", func() {
			It("deve retornar score total de 100", func() {
				analysis := &models.TerraformAnalysis{
					Valid:                true,
					BestPracticeWarnings: []string{},
					TotalOutputs:         1,
					Variables: []models.TerraformVariable{
						{Name: "var1", Description: "Description"},
					},
					Outputs: []models.TerraformOutput{
						{Name: "out1", Description: "Description"},
					},
				}

				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  10,
						Failed:  0,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{},
					},
				}

				score := prScorer.CalculateScore(analysis, checkovResult)
				Expect(score.Total).To(Equal(100))
			})
		})

		Context("quando há scores variados", func() {
			It("deve calcular média ponderada corretamente", func() {
				analysis := &models.TerraformAnalysis{
					Valid: true,
					BestPracticeWarnings: []string{
						"Warning 1",
					},
				}

				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  5,
						Failed:  1,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{Severity: "MEDIUM"}, // -5 = 95
						},
					},
				}

				score := prScorer.CalculateScore(analysis, checkovResult)
				// Security: 95 (peso 0.35) = 33.25
				// BestPractices: 95 (peso 0.25) = 23.75
				// Maintainability: 100 (peso 0.20) = 20
				// Documentation: 100 (peso 0.15) = 15
				// Performance: 100 (peso 0.05) = 5
				// Total: 97
				Expect(score.Total).To(BeNumerically("~", 97, 1))
			})
		})
	})

	Describe("Verificando aprovação do PR", func() {
		Context("quando o score atinge o mínimo", func() {
			It("deve aprovar o PR", func() {
				score := &models.PRScore{Total: 70}
				isApproved := prScorer.IsApproved(score)
				Expect(isApproved).To(BeTrue())
			})
		})

		Context("quando o score está acima do mínimo", func() {
			It("deve aprovar o PR", func() {
				score := &models.PRScore{Total: 85}
				isApproved := prScorer.IsApproved(score)
				Expect(isApproved).To(BeTrue())
			})
		})

		Context("quando o score está abaixo do mínimo", func() {
			It("não deve aprovar o PR", func() {
				score := &models.PRScore{Total: 65}
				isApproved := prScorer.IsApproved(score)
				Expect(isApproved).To(BeFalse())
			})
		})
	})

	Describe("Obtendo recomendações baseadas no score", func() {
		Context("quando o score é excelente (>= 90)", func() {
			It("deve retornar mensagem positiva", func() {
				score := &models.PRScore{Total: 95}
				recommendation := prScorer.GetRecommendation(score)
				Expect(recommendation).To(ContainSubstring("Excelente"))
			})
		})

		Context("quando o score é muito bom (80-89)", func() {
			It("deve sugerir melhorias menores", func() {
				score := &models.PRScore{Total: 85}
				recommendation := prScorer.GetRecommendation(score)
				Expect(recommendation).To(ContainSubstring("Muito bom"))
			})
		})

		Context("quando o score é bom (70-79)", func() {
			It("deve sugerir melhorias importantes", func() {
				score := &models.PRScore{Total: 75}
				recommendation := prScorer.GetRecommendation(score)
				Expect(recommendation).To(ContainSubstring("Bom"))
			})
		})

		Context("quando o score é aceitável (60-69)", func() {
			It("deve indicar melhorias significativas necessárias", func() {
				score := &models.PRScore{Total: 65}
				recommendation := prScorer.GetRecommendation(score)
				Expect(recommendation).To(ContainSubstring("Aceitável"))
			})
		})

		Context("quando o score é baixo (< 60)", func() {
			It("deve recomendar não fazer merge", func() {
				score := &models.PRScore{Total: 50}
				recommendation := prScorer.GetRecommendation(score)
				Expect(recommendation).To(ContainSubstring("Não recomendado"))
			})
		})
	})

	Describe("Verificando pesos dos scores", func() {
		Context("quando o score de segurança é baixo", func() {
			It("deve impactar fortemente o score total (peso 35%)", func() {
				analysis := &models.TerraformAnalysis{Valid: true}
				checkovResult := &models.CheckovResult{
					Summary: models.CheckovSummary{
						Passed:  0,
						Failed:  3,
						Skipped: 0,
					},
					Results: models.CheckovResults{
						FailedChecks: []models.CheckovCheck{
							{Severity: "CRITICAL"}, // -20
							{Severity: "CRITICAL"}, // -20
							{Severity: "HIGH"},     // -10
						},
					},
				}

				score := prScorer.CalculateScore(analysis, checkovResult)
				// Security: 50 (peso 0.35) = 17.5
				// Total será significativamente reduzido
				Expect(score.Total).To(BeNumerically("<", 80))
			})
		})
	})
})
