package suggester

import (
	"fmt"

	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// SecurityAdvisor gera recomendações de segurança
type SecurityAdvisor struct {
	logger *logger.Logger
}

// NewSecurityAdvisor cria uma nova instância
func NewSecurityAdvisor(log *logger.Logger) *SecurityAdvisor {
	return &SecurityAdvisor{
		logger: log,
	}
}

// GenerateSuggestions gera sugestões de segurança
func (sa *SecurityAdvisor) GenerateSuggestions(
	securityAnalysis *models.SecurityAnalysis,
	iamAnalysis *models.IAMAnalysis,
) []models.Suggestion {
	suggestions := []models.Suggestion{}

	// Sugestões baseadas em security findings
	for _, finding := range securityAnalysis.Findings {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "security",
			Severity:       mapSeverity(finding.Severity),
			Message:        finding.CheckName,
			Recommendation: finding.Guideline,
			File:           finding.File,
			Line:           finding.Line,
			ReferenceLink:  fmt.Sprintf("https://docs.bridgecrew.io/docs/%s", finding.CheckID),
		})
	}

	// Sugestões baseadas em IAM analysis
	if iamAnalysis.AdminAccessDetected {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "security",
			Severity:       "critical",
			Message:        "Acesso administrativo detectado (Action: *)",
			Recommendation: "Use princípio do menor privilégio. Especifique apenas as ações necessárias.",
		})
	}

	for _, publicAccess := range iamAnalysis.PublicAccess {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "security",
			Severity:       "high",
			Message:        publicAccess,
			Recommendation: "Restrinja o acesso público apenas quando absolutamente necessário.",
		})
	}

	for _, risk := range iamAnalysis.PrincipalRisks {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "security",
			Severity:       risk.RiskLevel,
			Message:        fmt.Sprintf("Risco de principal: %s", risk.Principal),
			Recommendation: risk.Reason,
		})
	}

	return suggestions
}

// mapSeverity mapeia severidade para formato padronizado
func mapSeverity(severity string) string {
	switch severity {
	case "CRITICAL":
		return "critical"
	case "HIGH":
		return "high"
	case "MEDIUM":
		return "medium"
	case "LOW":
		return "low"
	default:
		return "info"
	}
}
