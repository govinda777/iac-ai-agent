package unit_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/govinda777/iac-ai-agent/internal/agent/analyzer"
	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

var _ = Describe("IAMAnalyzer", func() {
	var (
		iamAnalyzer *analyzer.IAMAnalyzer
		log         *logger.Logger
	)

	BeforeEach(func() {
		log = logger.New("info", "json")
		iamAnalyzer = analyzer.NewIAMAnalyzer(log)
	})

	Describe("Analisando políticas IAM", func() {
		Context("quando detecta políticas com wildcard actions", func() {
			It("deve identificar ações com wildcard e marcar como overly permissive", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_policy",
							Name: "admin_policy",
							Attributes: map[string]interface{}{
								"policy": `{
									"Version": "2012-10-17",
									"Statement": [{
										"Effect": "Allow",
										"Action": "*",
										"Resource": "*"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.OverlyPermissive).To(BeTrue())
				Expect(analysis.AdminAccessDetected).To(BeTrue())
				Expect(analysis.WildcardActions).NotTo(BeEmpty())
				Expect(analysis.TotalPolicies).To(Equal(1))
			})
		})

		Context("quando detecta políticas específicas sem wildcards", func() {
			It("não deve marcar como overly permissive", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_policy",
							Name: "s3_read_policy",
							Attributes: map[string]interface{}{
								"policy": `{
									"Version": "2012-10-17",
									"Statement": [{
										"Effect": "Allow",
										"Action": ["s3:GetObject", "s3:ListBucket"],
										"Resource": "arn:aws:s3:::my-bucket/*"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.OverlyPermissive).To(BeFalse())
				Expect(analysis.AdminAccessDetected).To(BeFalse())
				Expect(analysis.WildcardActions).To(BeEmpty())
			})
		})

		Context("quando detecta wildcard parcial em actions", func() {
			It("deve identificar padrões como s3:*", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_policy",
							Name: "s3_all_policy",
							Attributes: map[string]interface{}{
								"policy": `{
									"Version": "2012-10-17",
									"Statement": [{
										"Effect": "Allow",
										"Action": "s3:*",
										"Resource": "arn:aws:s3:::my-bucket/*"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.OverlyPermissive).To(BeTrue())
				Expect(analysis.WildcardActions).To(ContainElement(ContainSubstring("s3:*")))
			})
		})
	})

	Describe("Analisando roles IAM", func() {
		Context("quando encontra roles AWS", func() {
			It("deve contar corretamente as roles", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_role",
							Name: "lambda_role",
							Attributes: map[string]interface{}{
								"assume_role_policy": `{
									"Version": "2012-10-17",
									"Statement": [{
										"Effect": "Allow",
										"Principal": {
											"Service": "lambda.amazonaws.com"
										},
										"Action": "sts:AssumeRole"
									}]
								}`,
							},
						},
						{
							Type:       "aws_iam_role",
							Name:       "ec2_role",
							Attributes: map[string]interface{}{},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.TotalRoles).To(Equal(2))
			})
		})

		Context("quando role permite serviços de risco", func() {
			It("deve adicionar aos principal risks", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_role",
							Name: "lambda_role",
							Attributes: map[string]interface{}{
								"assume_role_policy": `{
									"Version": "2012-10-17",
									"Statement": [{
										"Effect": "Allow",
										"Principal": {
											"Service": "lambda.amazonaws.com"
										},
										"Action": "sts:AssumeRole"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.PrincipalRisks).To(HaveLen(1))
				Expect(analysis.PrincipalRisks[0].Type).To(Equal("service"))
				Expect(analysis.PrincipalRisks[0].Principal).To(Equal("lambda.amazonaws.com"))
			})
		})
	})

	Describe("Detectando acesso público", func() {
		Context("quando política permite principal público (*)", func() {
			It("deve detectar e adicionar aos riscos", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_policy",
							Name: "public_policy",
							Attributes: map[string]interface{}{
								"policy": `{
									"Version": "2012-10-17",
									"Statement": [{
										"Effect": "Allow",
										"Principal": {
											"AWS": "*"
										},
										"Action": "s3:GetObject",
										"Resource": "*"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.PublicAccess).NotTo(BeEmpty())
				Expect(analysis.PublicAccess[0]).To(ContainSubstring("acesso público"))
				Expect(analysis.PrincipalRisks).To(HaveLen(1))
				Expect(analysis.PrincipalRisks[0].RiskLevel).To(Equal("critical"))
			})
		})

		Context("quando S3 bucket tem ACL pública", func() {
			It("deve detectar acesso público", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_s3_bucket",
							Name: "public_bucket",
							Attributes: map[string]interface{}{
								"acl": "public-read",
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.PublicAccess).NotTo(BeEmpty())
			})
		})

		Context("quando RDS instance está publicly accessible", func() {
			It("deve detectar acesso público", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_db_instance",
							Name: "public_db",
							Attributes: map[string]interface{}{
								"publicly_accessible": true,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.PublicAccess).NotTo(BeEmpty())
			})
		})

		Context("quando security group permite acesso de 0.0.0.0/0", func() {
			It("deve detectar acesso público", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_security_group",
							Name: "public_sg",
							Attributes: map[string]interface{}{
								"ingress": []interface{}{
									map[string]interface{}{
										"cidr_blocks": []interface{}{"0.0.0.0/0"},
										"from_port":   80,
										"to_port":     80,
									},
								},
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.PublicAccess).NotTo(BeEmpty())
			})
		})
	})

	Describe("Gerando recomendações", func() {
		Context("quando detecta acesso admin", func() {
			It("deve recomendar princípio do menor privilégio", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_policy",
							Name: "admin",
							Attributes: map[string]interface{}{
								"policy": `{
									"Statement": [{
										"Effect": "Allow",
										"Action": "*",
										"Resource": "*"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Recommendations).To(ContainElement(
					ContainSubstring("menor privilégio")))
			})
		})

		Context("quando detecta políticas muito permissivas", func() {
			It("deve recomendar especificar ações explicitamente", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_policy",
							Name: "permissive",
							Attributes: map[string]interface{}{
								"policy": `{
									"Statement": [{
										"Effect": "Allow",
										"Action": "s3:*",
										"Resource": "*"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.Recommendations).To(ContainElement(
					ContainSubstring("wildcard")))
			})
		})

		Context("quando há múltiplos recursos com acesso público", func() {
			It("deve recomendar revisar recursos públicos", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_s3_bucket",
							Name: "bucket1",
							Attributes: map[string]interface{}{
								"acl": "public-read",
							},
						},
						{
							Type: "aws_db_instance",
							Name: "db1",
							Attributes: map[string]interface{}{
								"publicly_accessible": true,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.PublicAccess).To(HaveLen(2))
				Expect(analysis.Recommendations).To(ContainElement(
					ContainSubstring("acesso público")))
			})
		})

		Context("quando há muitas ações com wildcard", func() {
			It("deve recomendar revisar permissões", func() {
				resources := []models.TerraformResource{}
				for i := 0; i < 5; i++ {
					resources = append(resources, models.TerraformResource{
						Type: "aws_iam_policy",
						Name: "policy",
						Attributes: map[string]interface{}{
							"policy": `{
								"Statement": [{
									"Effect": "Allow",
									"Action": "ec2:*",
									"Resource": "*"
								}]
							}`,
						},
					})
				}

				tfAnalysis := &models.TerraformAnalysis{
					Valid:     true,
					Resources: resources,
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(analysis.WildcardActions)).To(BeNumerically(">", 3))
				Expect(analysis.Recommendations).To(ContainElement(
					ContainSubstring("Muitas ações com wildcard")))
			})
		})
	})

	Describe("Analisando cenários complexos", func() {
		Context("quando há mix de políticas boas e ruins", func() {
			It("deve identificar apenas os problemas", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_iam_policy",
							Name: "good_policy",
							Attributes: map[string]interface{}{
								"policy": `{
									"Statement": [{
										"Effect": "Allow",
										"Action": ["s3:GetObject"],
										"Resource": "arn:aws:s3:::my-bucket/*"
									}]
								}`,
							},
						},
						{
							Type: "aws_iam_policy",
							Name: "bad_policy",
							Attributes: map[string]interface{}{
								"policy": `{
									"Statement": [{
										"Effect": "Allow",
										"Action": "*",
										"Resource": "*"
									}]
								}`,
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.TotalPolicies).To(Equal(2))
				Expect(analysis.WildcardActions).To(HaveLen(1))
				Expect(analysis.AdminAccessDetected).To(BeTrue())
			})
		})

		Context("quando não há recursos IAM", func() {
			It("deve retornar análise vazia sem erros", func() {
				tfAnalysis := &models.TerraformAnalysis{
					Valid: true,
					Resources: []models.TerraformResource{
						{
							Type: "aws_instance",
							Name: "web",
							Attributes: map[string]interface{}{
								"ami": "ami-12345",
							},
						},
					},
				}

				analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)

				Expect(err).NotTo(HaveOccurred())
				Expect(analysis.TotalPolicies).To(Equal(0))
				Expect(analysis.TotalRoles).To(Equal(0))
				Expect(analysis.WildcardActions).To(BeEmpty())
				Expect(analysis.PublicAccess).To(BeEmpty())
			})
		})
	})

	Describe("Identificando tipos de recursos IAM", func() {
		Context("quando verifica tipos de políticas AWS", func() {
			It("deve reconhecer todos os tipos de política", func() {
				policyTypes := []string{
					"aws_iam_policy",
					"aws_iam_role_policy",
					"aws_iam_user_policy",
					"aws_iam_group_policy",
				}

				for _, policyType := range policyTypes {
					tfAnalysis := &models.TerraformAnalysis{
						Valid: true,
						Resources: []models.TerraformResource{
							{
								Type:       policyType,
								Name:       "test_policy",
								Attributes: map[string]interface{}{},
							},
						},
					}

					analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)
					Expect(err).NotTo(HaveOccurred())
					Expect(analysis.TotalPolicies).To(Equal(1),
						"Should recognize %s as IAM policy", policyType)
				}
			})
		})

		Context("quando verifica tipos de roles multi-cloud", func() {
			It("deve reconhecer roles de diferentes providers", func() {
				roleTypes := []string{
					"aws_iam_role",
					"azurerm_role_assignment",
					"google_project_iam_member",
				}

				for _, roleType := range roleTypes {
					tfAnalysis := &models.TerraformAnalysis{
						Valid: true,
						Resources: []models.TerraformResource{
							{
								Type:       roleType,
								Name:       "test_role",
								Attributes: map[string]interface{}{},
							},
						},
					}

					analysis, err := iamAnalyzer.AnalyzeTerraform(tfAnalysis)
					Expect(err).NotTo(HaveOccurred())
					Expect(analysis.TotalRoles).To(Equal(1),
						"Should recognize %s as IAM role", roleType)
				}
			})
		})
	})
})
