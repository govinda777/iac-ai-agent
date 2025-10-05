package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// CheckovAnalyzer executa análise de segurança usando Checkov
type CheckovAnalyzer struct {
	checkovPath string
	logger      *logger.Logger
}

// NewCheckovAnalyzer cria uma nova instância do analisador Checkov
func NewCheckovAnalyzer(log *logger.Logger) *CheckovAnalyzer {
	checkovPath, _ := exec.LookPath("checkov")
	return &CheckovAnalyzer{
		checkovPath: checkovPath,
		logger:      log,
	}
}

// IsAvailable verifica se o Checkov está instalado
func (ca *CheckovAnalyzer) IsAvailable() bool {
	return ca.checkovPath != ""
}

// AnalyzeDirectory executa Checkov em um diretório
func (ca *CheckovAnalyzer) AnalyzeDirectory(dir string, config *models.CheckovConfig) (*models.SecurityAnalysis, error) {
	if !ca.IsAvailable() {
		return nil, fmt.Errorf("checkov não está instalado ou não foi encontrado no PATH")
	}

	// Prepara comando
	args := []string{
		"-d", dir,
		"-o", "json",
		"--quiet",
		"--compact",
	}

	if config != nil {
		if config.Framework != "" {
			args = append(args, "--framework", config.Framework)
		}
		if len(config.SkipChecks) > 0 {
			args = append(args, "--skip-check", strings.Join(config.SkipChecks, ","))
		}
		if len(config.Checks) > 0 {
			args = append(args, "--check", strings.Join(config.Checks, ","))
		}
	}

	// Executa Checkov
	ca.logger.Info("Executando Checkov", "directory", dir)
	cmd := exec.Command(ca.checkovPath, args...)
	output, err := cmd.Output()

	// Checkov retorna exit code != 0 quando há falhas, mas isso é esperado
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Se há stderr, algo deu errado
			if len(exitErr.Stderr) > 0 {
				return nil, fmt.Errorf("erro ao executar checkov: %s", string(exitErr.Stderr))
			}
			// Caso contrário, provavelmente há apenas falhas de check
			output = exitErr.Stderr
		} else {
			return nil, fmt.Errorf("erro ao executar checkov: %w", err)
		}
	}

	// Parse resultado
	var result models.CheckovResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do resultado checkov: %w", err)
	}

	return ca.convertToSecurityAnalysis(&result), nil
}

// AnalyzeFiles executa Checkov em arquivos específicos
func (ca *CheckovAnalyzer) AnalyzeFiles(files []string, config *models.CheckovConfig) (*models.SecurityAnalysis, error) {
	if !ca.IsAvailable() {
		return nil, fmt.Errorf("checkov não está instalado")
	}

	// Cria diretório temporário
	tmpDir, err := os.MkdirTemp("", "checkov-analysis-*")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório temporário: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Copia arquivos para o diretório temporário
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		tmpFile := filepath.Join(tmpDir, filepath.Base(file))
		if err := os.WriteFile(tmpFile, content, 0644); err != nil {
			continue
		}
	}

	return ca.AnalyzeDirectory(tmpDir, config)
}

// convertToSecurityAnalysis converte resultado Checkov para SecurityAnalysis
func (ca *CheckovAnalyzer) convertToSecurityAnalysis(result *models.CheckovResult) *models.SecurityAnalysis {
	analysis := &models.SecurityAnalysis{
		ChecksPassed: result.Summary.Passed,
		ChecksFailed: result.Summary.Failed,
		TotalIssues:  result.Summary.Failed,
		Findings:     []models.SecurityFinding{},
	}

	// Processa falhas
	for _, check := range result.Results.FailedChecks {
		finding := models.SecurityFinding{
			ID:          check.CheckID,
			CheckID:     check.CheckID,
			CheckName:   check.CheckName,
			Severity:    ca.determineSeverity(check),
			Resource:    check.Resource,
			File:        check.File,
			Description: check.Description,
			Guideline:   check.Guideline,
		}

		if len(check.FileLineRange) > 0 {
			finding.Line = check.FileLineRange[0]
		}

		analysis.Findings = append(analysis.Findings, finding)

		// Conta por severidade
		switch finding.Severity {
		case "CRITICAL":
			analysis.Critical++
		case "HIGH":
			analysis.High++
		case "MEDIUM":
			analysis.Medium++
		case "LOW":
			analysis.Low++
		default:
			analysis.Info++
		}
	}

	return analysis
}

// determineSeverity determina a severidade de um check
func (ca *CheckovAnalyzer) determineSeverity(check models.CheckovCheck) string {
	if check.Severity != "" {
		return strings.ToUpper(check.Severity)
	}

	// Heurísticas baseadas no check ID ou nome
	checkLower := strings.ToLower(check.CheckID + " " + check.CheckName)

	if strings.Contains(checkLower, "encryption") ||
		strings.Contains(checkLower, "credential") ||
		strings.Contains(checkLower, "secret") ||
		strings.Contains(checkLower, "password") ||
		strings.Contains(checkLower, "public access") {
		return "HIGH"
	}

	if strings.Contains(checkLower, "logging") ||
		strings.Contains(checkLower, "monitoring") ||
		strings.Contains(checkLower, "versioning") {
		return "MEDIUM"
	}

	if strings.Contains(checkLower, "tag") ||
		strings.Contains(checkLower, "description") {
		return "LOW"
	}

	return "MEDIUM"
}

// GetRecommendations gera recomendações baseadas nas findings
func (ca *CheckovAnalyzer) GetRecommendations(analysis *models.SecurityAnalysis) []models.Suggestion {
	suggestions := []models.Suggestion{}

	for _, finding := range analysis.Findings {
		suggestion := models.Suggestion{
			Type:           "security",
			Severity:       strings.ToLower(finding.Severity),
			Message:        finding.CheckName,
			Recommendation: finding.Guideline,
			File:           finding.File,
			Line:           finding.Line,
			ReferenceLink:  ca.getBridgecrewLink(finding.CheckID),
		}

		suggestions = append(suggestions, suggestion)
	}

	return suggestions
}

// getBridgecrewLink retorna link para documentação do check
func (ca *CheckovAnalyzer) getBridgecrewLink(checkID string) string {
	if checkID == "" {
		return ""
	}
	return fmt.Sprintf("https://docs.bridgecrew.io/docs/%s", strings.ToLower(checkID))
}
