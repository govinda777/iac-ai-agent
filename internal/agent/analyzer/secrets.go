package analyzer

import (
	"regexp"
	"strings"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
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
	}
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

	return findings
}

// maskSecret mascara o valor do secret
func (sa *SecretsAnalyzer) maskSecret(line string) string {
	trimmedLine := strings.TrimSpace(line)
	if len(trimmedLine) > 50 {
		return trimmedLine[:20] + "***REDACTED***" + trimmedLine[len(trimmedLine)-10:]
	}
	return "***REDACTED***"
}