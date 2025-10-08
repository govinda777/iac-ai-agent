package analyzer

import (
	"regexp"
	"strings"

	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

type SecretsAnalyzer struct {
	patterns []SecretPattern
	logger   *logger.Logger
}

type SecretPattern struct {
	Name        string
	Regex       *regexp.Regexp
	Severity    string
	Description string
	Suggestion  string
}

func NewSecretsAnalyzer(log *logger.Logger) *SecretsAnalyzer {
	sa := &SecretsAnalyzer{
		logger:   log,
		patterns: []SecretPattern{},
	}
	sa.loadPatterns()
	return sa
}

func (sa *SecretsAnalyzer) loadPatterns() {
	sa.patterns = []SecretPattern{
		{
			Name:        "AWS Access Key",
			Regex:       regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
			Severity:    "critical",
			Description: "AWS Access Key ID detected in plaintext",
			Suggestion:  "Use AWS Secrets Manager or environment variables",
		},
		{
			Name:        "AWS Secret Key",
			Regex:       regexp.MustCompile(`aws_secret_access_key\s*=\s*["']([^"']+)["']`),
			Severity:    "critical",
			Description: "AWS Secret Access Key in plaintext",
			Suggestion:  "Never commit AWS credentials. Use IAM roles or AWS SSO",
		},
		{
			Name:        "Generic Password",
			Regex:       regexp.MustCompile(`password\s*=\s*["']([^"']+)["']`),
			Severity:    "high",
			Description: "Password in plaintext",
			Suggestion:  "Use variable with sensitive=true or secrets manager",
		},
		{
			Name:        "Private Key",
			Regex:       regexp.MustCompile(`-----BEGIN.*PRIVATE KEY-----`),
			Severity:    "critical",
			Description: "Private key detected in code",
			Suggestion:  "Store private keys in secure vault",
		},
		{
			Name:        "API Key",
			Regex:       regexp.MustCompile(`api[_-]?key\s*=\s*["']([A-Za-z0-9]{20,})["']`),
			Severity:    "high",
			Description: "API key in plaintext",
			Suggestion:  "Use environment variables or parameter store",
		},
		{
			Name:        "Database Password",
			Regex:       regexp.MustCompile(`db[_-]?password\s*=\s*["']([^"']+)["']`),
			Severity:    "high",
			Description: "Database password in plaintext",
			Suggestion:  "Use AWS RDS secrets manager or environment variables",
		},
		{
			Name:        "JWT Secret",
			Regex:       regexp.MustCompile(`jwt[_-]?secret\s*=\s*["']([^"']+)["']`),
			Severity:    "critical",
			Description: "JWT secret in plaintext",
			Suggestion:  "Use environment variables or secrets manager",
		},
		{
			Name:        "GitHub Token",
			Regex:       regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`),
			Severity:    "critical",
			Description: "GitHub personal access token detected",
			Suggestion:  "Use GitHub secrets or environment variables",
		},
		{
			Name:        "Slack Token",
			Regex:       regexp.MustCompile(`xoxb-[0-9]{11}-[0-9]{11}-[A-Za-z0-9]{24}`),
			Severity:    "high",
			Description: "Slack bot token detected",
			Suggestion:  "Use Slack secrets or environment variables",
		},
		{
			Name:        "Generic Token",
			Regex:       regexp.MustCompile(`token\s*=\s*["']([A-Za-z0-9]{32,})["']`),
			Severity:    "high",
			Description: "Generic token in plaintext",
			Suggestion:  "Use environment variables or secrets manager",
		},
		{
			Name:        "SSH Private Key",
			Regex:       regexp.MustCompile(`-----BEGIN OPENSSH PRIVATE KEY-----`),
			Severity:    "critical",
			Description: "SSH private key detected",
			Suggestion:  "Store SSH keys in secure vault or use SSH agent",
		},
		{
			Name:        "Azure Storage Key",
			Regex:       regexp.MustCompile(`DefaultEndpointsProtocol=https;AccountName=[^;]+;AccountKey=[^;]+`),
			Severity:    "critical",
			Description: "Azure storage account key in connection string",
			Suggestion:  "Use Azure Key Vault or managed identity",
		},
		{
			Name:        "Google Service Account Key",
			Regex:       regexp.MustCompile(`"type":\s*"service_account"`),
			Severity:    "critical",
			Description: "Google service account key detected",
			Suggestion:  "Use Google Cloud Secret Manager or workload identity",
		},
		{
			Name:        "Docker Registry Password",
			Regex:       regexp.MustCompile(`docker[_-]?password\s*=\s*["']([^"']+)["']`),
			Severity:    "high",
			Description: "Docker registry password in plaintext",
			Suggestion:  "Use Docker secrets or environment variables",
		},
		{
			Name:        "Redis Password",
			Regex:       regexp.MustCompile(`redis[_-]?password\s*=\s*["']([^"']+)["']`),
			Severity:    "high",
			Description: "Redis password in plaintext",
			Suggestion:  "Use AWS ElastiCache auth or environment variables",
		},
	}

	sa.logger.Info("Secrets patterns loaded", "count", len(sa.patterns))
}

// ScanContent escaneia conteúdo em busca de secrets
func (sa *SecretsAnalyzer) ScanContent(
	content string,
	filename string,
) []models.SecretFinding {
	findings := []models.SecretFinding{}
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		for _, pattern := range sa.patterns {
			if pattern.Regex.MatchString(line) {
				findings = append(findings, models.SecretFinding{
					Type:        pattern.Name,
					File:        filename,
					Line:        lineNum + 1,
					Value:       sa.maskSecret(line),
					Severity:    pattern.Severity,
					Description: pattern.Description,
					Suggestion:  pattern.Suggestion,
				})
			}
		}
	}

	sa.logger.Info("Secrets scan completed",
		"file", filename,
		"findings", len(findings))

	return findings
}

// ScanDirectory escaneia um diretório completo
func (sa *SecretsAnalyzer) ScanDirectory(dir string) ([]models.SecretFinding, error) {
	// Esta implementação seria expandida para escanear arquivos
	// Por enquanto, retorna um placeholder
	sa.logger.Info("Directory secrets scan requested", "directory", dir)
	return []models.SecretFinding{}, nil
}

// maskSecret mascara o valor do secret
func (sa *SecretsAnalyzer) maskSecret(line string) string {
	if len(line) > 50 {
		return line[:20] + "***REDACTED***" + line[len(line)-10:]
	}
	return "***REDACTED***"
}

// GetPatterns retorna todos os padrões carregados
func (sa *SecretsAnalyzer) GetPatterns() []SecretPattern {
	return sa.patterns
}

// AddCustomPattern adiciona um padrão customizado
func (sa *SecretsAnalyzer) AddCustomPattern(pattern SecretPattern) error {
	// Valida o regex
	if _, err := regexp.Compile(pattern.Regex.String()); err != nil {
		return err
	}

	sa.patterns = append(sa.patterns, pattern)
	sa.logger.Info("Custom pattern added", "name", pattern.Name)
	return nil
}

// RemovePattern remove um padrão por nome
func (sa *SecretsAnalyzer) RemovePattern(name string) bool {
	for i, pattern := range sa.patterns {
		if pattern.Name == name {
			sa.patterns = append(sa.patterns[:i], sa.patterns[i+1:]...)
			sa.logger.Info("Pattern removed", "name", name)
			return true
		}
	}
	return false
}

// ValidateSecrets valida se secrets encontrados são válidos
func (sa *SecretsAnalyzer) ValidateSecrets(findings []models.SecretFinding) []models.SecretFinding {
	validated := []models.SecretFinding{}

	for _, finding := range findings {
		// Adiciona validações específicas por tipo
		if sa.isValidSecret(finding) {
			validated = append(validated, finding)
		}
	}

	return validated
}

// isValidSecret valida se um secret encontrado é realmente um secret válido
func (sa *SecretsAnalyzer) isValidSecret(finding models.SecretFinding) bool {
	// Implementa validações específicas por tipo
	switch finding.Type {
	case "AWS Access Key":
		// Valida formato AKIA + 16 caracteres alfanuméricos
		return len(finding.Value) >= 20 && strings.HasPrefix(finding.Value, "AKIA")
	case "Private Key":
		// Valida se contém BEGIN e END
		return strings.Contains(finding.Value, "BEGIN") && strings.Contains(finding.Value, "END")
	case "API Key":
		// Valida comprimento mínimo
		return len(finding.Value) >= 20
	default:
		return true // Para outros tipos, assume válido
	}
}

// GetSeverityCount retorna contagem por severidade
func (sa *SecretsAnalyzer) GetSeverityCount(findings []models.SecretFinding) map[string]int {
	counts := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
	}

	for _, finding := range findings {
		counts[finding.Severity]++
	}

	return counts
}

// GenerateReport gera relatório de secrets encontrados
func (sa *SecretsAnalyzer) GenerateReport(findings []models.SecretFinding) *models.SecretsReport {
	severityCounts := sa.GetSeverityCount(findings)

	report := &models.SecretsReport{
		TotalFindings: len(findings),
		CriticalCount: severityCounts["critical"],
		HighCount:     severityCounts["high"],
		MediumCount:   severityCounts["medium"],
		LowCount:      severityCounts["low"],
		Findings:      findings,
		RiskLevel:     sa.calculateOverallRisk(severityCounts),
	}

	return report
}

// calculateOverallRisk calcula risco geral baseado nas severidades
func (sa *SecretsAnalyzer) calculateOverallRisk(counts map[string]int) string {
	if counts["critical"] > 0 {
		return "critical"
	}
	if counts["high"] > 2 {
		return "high"
	}
	if counts["high"] > 0 || counts["medium"] > 3 {
		return "medium"
	}
	return "low"
}
