package cloudcontroller

import (
	"strings"
)

// KnowledgeBase contém conhecimento sobre best practices e padrões
type KnowledgeBase struct {
	bestPractices map[string][]string
	securityRules map[string]string
}

// NewKnowledgeBase cria uma nova base de conhecimento
func NewKnowledgeBase() *KnowledgeBase {
	kb := &KnowledgeBase{
		bestPractices: make(map[string][]string),
		securityRules: make(map[string]string),
	}

	kb.loadBestPractices()
	kb.loadSecurityRules()

	return kb
}

// loadBestPractices carrega best practices
func (kb *KnowledgeBase) loadBestPractices() {
	kb.bestPractices["aws_s3_bucket"] = []string{
		"Habilite versionamento para proteção contra exclusão acidental",
		"Configure encriptação em repouso (AES-256 ou KMS)",
		"Bloqueie acesso público a menos que seja absolutamente necessário",
		"Configure lifecycle policies para otimizar custos",
		"Habilite logging de acesso",
	}

	kb.bestPractices["aws_instance"] = []string{
		"Use IMDSv2 para metadata service",
		"Habilite monitoring detalhado",
		"Use Security Groups restritivos",
		"Configure backup automático via AWS Backup",
		"Considere usar Auto Scaling Groups",
	}

	kb.bestPractices["aws_rds_instance"] = []string{
		"Habilite encriptação em repouso",
		"Configure automated backups com retenção apropriada",
		"Use Parameter Groups customizados",
		"Habilite Enhanced Monitoring",
		"Configure Multi-AZ para produção",
	}

	kb.bestPractices["aws_security_group"] = []string{
		"Siga princípio do menor privilégio",
		"Evite 0.0.0.0/0 para ingress (exceto portas públicas necessárias)",
		"Use Security Group IDs como source ao invés de CIDRs",
		"Adicione descrições claras em todas as regras",
		"Revise periodicamente regras não utilizadas",
	}
}

// loadSecurityRules carrega regras de segurança
func (kb *KnowledgeBase) loadSecurityRules() {
	kb.securityRules["public_s3_bucket"] = "S3 buckets não devem ser públicos a menos que haja justificativa clara"
	kb.securityRules["unencrypted_storage"] = "Todos os storages devem ser encriptados em repouso"
	kb.securityRules["open_security_group"] = "Security groups não devem permitir acesso de 0.0.0.0/0 exceto para portas públicas"
	kb.securityRules["weak_iam_policy"] = "Políticas IAM não devem usar wildcards (*) em Actions e Resources"
	kb.securityRules["missing_logging"] = "Recursos críticos devem ter logging habilitado"
}

// GetBestPractices retorna best practices para um tipo de recurso
func (kb *KnowledgeBase) GetBestPractices(resourceType string) []string {
	if practices, ok := kb.bestPractices[resourceType]; ok {
		return practices
	}
	return []string{}
}

// GetSecurityRule retorna uma regra de segurança
func (kb *KnowledgeBase) GetSecurityRule(ruleID string) string {
	if rule, ok := kb.securityRules[ruleID]; ok {
		return rule
	}
	return ""
}

// SearchBestPractices busca best practices por palavra-chave
func (kb *KnowledgeBase) SearchBestPractices(keyword string) map[string][]string {
	results := make(map[string][]string)
	keyword = strings.ToLower(keyword)

	for resourceType, practices := range kb.bestPractices {
		matchingPractices := []string{}
		for _, practice := range practices {
			if strings.Contains(strings.ToLower(practice), keyword) {
				matchingPractices = append(matchingPractices, practice)
			}
		}
		if len(matchingPractices) > 0 {
			results[resourceType] = matchingPractices
		}
	}

	return results
}
