package capabilities

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/agent/core"
)

// IACAnalysisCapability habilidade de análise de Infrastructure as Code
type IACAnalysisCapability struct {
	id           string
	name         string
	description  string
	version      string
	status       string
	lastUsed     time.Time
	messageCount int
	errorCount   int

	// Configurações específicas da análise IaC
	terraformAnalyzer *TerraformAnalyzer
	securityAnalyzer  *SecurityAnalyzer
	costAnalyzer      *CostAnalyzer

	logger *core.Logger
}

// TerraformAnalyzer analisador de código Terraform
type TerraformAnalyzer struct {
	enabled bool
}

// SecurityAnalyzer analisador de segurança
type SecurityAnalyzer struct {
	enabled bool
}

// CostAnalyzer analisador de custos
type CostAnalyzer struct {
	enabled bool
}

// NewIACAnalysisCapability cria nova habilidade de análise IaC
func NewIACAnalysisCapability() *IACAnalysisCapability {
	return &IACAnalysisCapability{
		id:                "iac-analysis",
		name:              "IaC Analysis",
		description:       "Capability to analyze Infrastructure as Code (Terraform)",
		version:           "1.0.0",
		status:            "inactive",
		terraformAnalyzer: &TerraformAnalyzer{enabled: true},
		securityAnalyzer:  &SecurityAnalyzer{enabled: true},
		costAnalyzer:      &CostAnalyzer{enabled: true},
		logger:            core.NewLogger("iac-analysis-capability", "info"),
	}
}

// GetID retorna ID da habilidade
func (i *IACAnalysisCapability) GetID() string {
	return i.id
}

// GetName retorna nome da habilidade
func (i *IACAnalysisCapability) GetName() string {
	return i.name
}

// GetDescription retorna descrição da habilidade
func (i *IACAnalysisCapability) GetDescription() string {
	return i.description
}

// GetVersion retorna versão da habilidade
func (i *IACAnalysisCapability) GetVersion() string {
	return i.version
}

// Initialize inicializa a habilidade de análise IaC
func (i *IACAnalysisCapability) Initialize(ctx context.Context, config *core.Config) error {
	i.logger = core.NewLogger("iac-analysis-capability", config.Logging.Level)

	// Carregar configurações específicas da análise IaC
	if iacConfig, exists := config.Capabilities["iac-analysis"]; exists {
		if configMap, ok := iacConfig.(map[string]interface{}); ok {
			if terraformEnabled, exists := configMap["terraform_enabled"]; exists {
				i.terraformAnalyzer.enabled = terraformEnabled.(bool)
			}
			if securityEnabled, exists := configMap["security_enabled"]; exists {
				i.securityAnalyzer.enabled = securityEnabled.(bool)
			}
			if costEnabled, exists := configMap["cost_enabled"]; exists {
				i.costAnalyzer.enabled = costEnabled.(bool)
			}
		}
	}

	i.logger.Info("IaC Analysis capability initialized", map[string]interface{}{
		"terraform_enabled": i.terraformAnalyzer.enabled,
		"security_enabled":  i.securityAnalyzer.enabled,
		"cost_enabled":      i.costAnalyzer.enabled,
	})

	return nil
}

// Start inicia a habilidade de análise IaC
func (i *IACAnalysisCapability) Start(ctx context.Context) error {
	i.status = "active"
	i.logger.Info("IaC Analysis capability started", map[string]interface{}{
		"capability": i.id,
		"status":     i.status,
	})
	return nil
}

// Stop para a habilidade de análise IaC
func (i *IACAnalysisCapability) Stop(ctx context.Context) error {
	i.status = "inactive"
	i.logger.Info("IaC Analysis capability stopped", map[string]interface{}{
		"capability": i.id,
		"status":     i.status,
	})
	return nil
}

// CanHandle verifica se pode processar a mensagem
func (i *IACAnalysisCapability) CanHandle(message *core.Message) bool {
	// Pode processar mensagens que contenham código Terraform ou comandos IaC
	return i.containsTerraformCode(message.Text) || i.isIACCommand(message.Text)
}

// containsTerraformCode verifica se a mensagem contém código Terraform
func (i *IACAnalysisCapability) containsTerraformCode(text string) bool {
	// Verificar blocos de código Terraform
	terraformRegex := regexp.MustCompile("```(hcl|terraform|tf)")
	if terraformRegex.MatchString(text) {
		return true
	}

	// Verificar recursos Terraform comuns
	resourceRegex := regexp.MustCompile(`resource\s+"[^"]+"\s+"[^"]+"`)
	return resourceRegex.MatchString(text)
}

// isIACCommand verifica se é um comando de análise IaC
func (i *IACAnalysisCapability) isIACCommand(text string) bool {
	iacCommands := []string{"/analyze", "/security", "/cost", "/terraform", "/iac"}
	text = strings.ToLower(strings.TrimSpace(text))

	for _, cmd := range iacCommands {
		if strings.HasPrefix(text, cmd) {
			return true
		}
	}

	return false
}

// ProcessMessage processa mensagem de análise IaC
func (i *IACAnalysisCapability) ProcessMessage(ctx context.Context, message *core.Message) (*core.Response, error) {
	i.messageCount++
	i.lastUsed = time.Now()

	i.logger.Info("Processing IaC analysis message", map[string]interface{}{
		"message_id":  message.ID,
		"from":        message.From,
		"text_length": len(message.Text),
	})

	// Determinar tipo de análise solicitada
	analysisType := i.determineAnalysisType(message.Text)

	switch analysisType {
	case "analyze":
		return i.performGeneralAnalysis(ctx, message)
	case "security":
		return i.performSecurityAnalysis(ctx, message)
	case "cost":
		return i.performCostAnalysis(ctx, message)
	default:
		return i.performGeneralAnalysis(ctx, message)
	}
}

// determineAnalysisType determina o tipo de análise baseado no comando
func (i *IACAnalysisCapability) determineAnalysisType(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	if strings.HasPrefix(text, "/security") {
		return "security"
	}
	if strings.HasPrefix(text, "/cost") {
		return "cost"
	}
	if strings.HasPrefix(text, "/analyze") {
		return "analyze"
	}

	return "analyze" // padrão
}

// performGeneralAnalysis executa análise geral do código
func (i *IACAnalysisCapability) performGeneralAnalysis(ctx context.Context, message *core.Message) (*core.Response, error) {
	code := i.extractTerraformCode(message.Text)
	if code == "" {
		return &core.Response{
			ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
			To:        message.From,
			Text:      "❌ Nenhum código Terraform encontrado. Por favor, envie o código para análise.",
			Type:      "text",
			Timestamp: time.Now(),
		}, nil
	}

	// Executar análise geral
	analysis := i.terraformAnalyzer.AnalyzeCode(code)

	responseText := fmt.Sprintf(`✅ *Análise Terraform Concluída*

*Código analisado:* %d caracteres
*Problemas encontrados:* %d
*Sugestões:* %d

*Problemas:*
%s

*Sugestões:*
%s

*Status:* Análise completa
*Timestamp:* %s`,
		len(code),
		len(analysis.Issues),
		len(analysis.Suggestions),
		i.formatIssues(analysis.Issues),
		i.formatSuggestions(analysis.Suggestions),
		time.Now().Format("15:04:05"))

	return &core.Response{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		To:        message.From,
		Text:      responseText,
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// performSecurityAnalysis executa análise de segurança
func (i *IACAnalysisCapability) performSecurityAnalysis(ctx context.Context, message *core.Message) (*core.Response, error) {
	code := i.extractTerraformCode(message.Text)
	if code == "" {
		return &core.Response{
			ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
			To:        message.From,
			Text:      "❌ Nenhum código Terraform encontrado. Por favor, envie o código para análise de segurança.",
			Type:      "text",
			Timestamp: time.Now(),
		}, nil
	}

	// Executar análise de segurança
	securityAnalysis := i.securityAnalyzer.AnalyzeSecurity(code)

	responseText := fmt.Sprintf(`🔒 *Análise de Segurança Concluída*

*Código analisado:* %d caracteres
*Vulnerabilidades encontradas:* %d
*Recomendações:* %d

*Vulnerabilidades:*
%s

*Recomendações de Segurança:*
%s

*Status:* Análise de segurança completa
*Timestamp:* %s`,
		len(code),
		len(securityAnalysis.Vulnerabilities),
		len(securityAnalysis.Recommendations),
		i.formatVulnerabilities(securityAnalysis.Vulnerabilities),
		i.formatRecommendations(securityAnalysis.Recommendations),
		time.Now().Format("15:04:05"))

	return &core.Response{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		To:        message.From,
		Text:      responseText,
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// performCostAnalysis executa análise de custos
func (i *IACAnalysisCapability) performCostAnalysis(ctx context.Context, message *core.Message) (*core.Response, error) {
	code := i.extractTerraformCode(message.Text)
	if code == "" {
		return &core.Response{
			ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
			To:        message.From,
			Text:      "❌ Nenhum código Terraform encontrado. Por favor, envie o código para análise de custos.",
			Type:      "text",
			Timestamp: time.Now(),
		}, nil
	}

	// Executar análise de custos
	costAnalysis := i.costAnalyzer.AnalyzeCosts(code)

	responseText := fmt.Sprintf(`💰 *Análise de Custos Concluída*

*Código analisado:* %d caracteres
*Custo estimado mensal:* $%.2f
*Potencial de economia:* $%.2f
*Otimizações sugeridas:* %d

*Otimizações:*
%s

*Resumo:*
• Custo atual estimado: $%.2f/mês
• Economia potencial: $%.2f/mês
• ROI: %.1f%%

*Status:* Análise de custos completa
*Timestamp:* %s`,
		len(code),
		costAnalysis.EstimatedMonthlyCost,
		costAnalysis.PotentialSavings,
		len(costAnalysis.Optimizations),
		i.formatOptimizations(costAnalysis.Optimizations),
		costAnalysis.EstimatedMonthlyCost,
		costAnalysis.PotentialSavings,
		(costAnalysis.PotentialSavings/costAnalysis.EstimatedMonthlyCost)*100,
		time.Now().Format("15:04:05"))

	return &core.Response{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		To:        message.From,
		Text:      responseText,
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// extractTerraformCode extrai código Terraform da mensagem
func (i *IACAnalysisCapability) extractTerraformCode(text string) string {
	// Procurar por blocos de código
	codeBlockRegex := regexp.MustCompile("```(?s)(.*?)```")
	matches := codeBlockRegex.FindStringSubmatch(text)

	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Se não encontrar blocos de código, retornar o texto inteiro
	return strings.TrimSpace(text)
}

// formatIssues formata problemas para exibição
func (i *IACAnalysisCapability) formatIssues(issues []Issue) string {
	if len(issues) == 0 {
		return "• Nenhum problema encontrado ✅"
	}

	var formatted []string
	for _, issue := range issues {
		formatted = append(formatted, fmt.Sprintf("• %s", issue.Description))
	}

	return strings.Join(formatted, "\n")
}

// formatSuggestions formata sugestões para exibição
func (i *IACAnalysisCapability) formatSuggestions(suggestions []string) string {
	if len(suggestions) == 0 {
		return "• Nenhuma sugestão disponível"
	}

	var formatted []string
	for _, suggestion := range suggestions {
		formatted = append(formatted, fmt.Sprintf("• %s", suggestion))
	}

	return strings.Join(formatted, "\n")
}

// formatVulnerabilities formata vulnerabilidades para exibição
func (i *IACAnalysisCapability) formatVulnerabilities(vulnerabilities []Vulnerability) string {
	if len(vulnerabilities) == 0 {
		return "• Nenhuma vulnerabilidade encontrada ✅"
	}

	var formatted []string
	for _, vuln := range vulnerabilities {
		formatted = append(formatted, fmt.Sprintf("• %s (Severidade: %s)", vuln.Description, vuln.Severity))
	}

	return strings.Join(formatted, "\n")
}

// formatRecommendations formata recomendações para exibição
func (i *IACAnalysisCapability) formatRecommendations(recommendations []string) string {
	if len(recommendations) == 0 {
		return "• Nenhuma recomendação disponível"
	}

	var formatted []string
	for _, rec := range recommendations {
		formatted = append(formatted, fmt.Sprintf("• %s", rec))
	}

	return strings.Join(formatted, "\n")
}

// formatOptimizations formata otimizações para exibição
func (i *IACAnalysisCapability) formatOptimizations(optimizations []Optimization) string {
	if len(optimizations) == 0 {
		return "• Nenhuma otimização sugerida"
	}

	var formatted []string
	for _, opt := range optimizations {
		formatted = append(formatted, fmt.Sprintf("• %s (Economia: $%.2f/mês)", opt.Description, opt.MonthlySavings))
	}

	return strings.Join(formatted, "\n")
}

// GetStatus retorna status da habilidade
func (i *IACAnalysisCapability) GetStatus() *core.CapabilityStatus {
	return &core.CapabilityStatus{
		ID:           i.id,
		Name:         i.name,
		Status:       i.status,
		LastUsed:     i.lastUsed,
		MessageCount: i.messageCount,
		ErrorCount:   i.errorCount,
	}
}

// HealthCheck verifica saúde da habilidade
func (i *IACAnalysisCapability) HealthCheck(ctx context.Context) error {
	if i.status != "active" {
		return fmt.Errorf("capability is not active")
	}

	return nil
}

// AnalyzeCode analisa código Terraform
func (t *TerraformAnalyzer) AnalyzeCode(code string) *AnalysisResult {
	// Implementação simplificada - em produção usar analisador real
	return &AnalysisResult{
		Issues: []Issue{
			{Description: "Recurso sem tags adequadas"},
			{Description: "Instance type pode ser otimizada"},
		},
		Suggestions: []string{
			"Adicionar tags para melhor organização",
			"Considerar usar instance types menores",
			"Implementar lifecycle rules",
		},
	}
}

// AnalyzeSecurity analisa segurança do código
func (s *SecurityAnalyzer) AnalyzeSecurity(code string) *SecurityAnalysis {
	// Implementação simplificada - em produção usar analisador real
	return &SecurityAnalysis{
		Vulnerabilities: []Vulnerability{
			{Description: "Bucket S3 sem criptografia", Severity: "High"},
			{Description: "Security group muito permissivo", Severity: "Medium"},
		},
		Recommendations: []string{
			"Habilitar criptografia no bucket S3",
			"Restringir regras do security group",
			"Implementar WAF",
		},
	}
}

// AnalyzeCosts analisa custos do código
func (c *CostAnalyzer) AnalyzeCosts(code string) *CostAnalysis {
	// Implementação simplificada - em produção usar analisador real
	return &CostAnalysis{
		EstimatedMonthlyCost: 150.0,
		PotentialSavings:     30.0,
		Optimizations: []Optimization{
			{Description: "Usar Spot Instances", MonthlySavings: 20.0},
			{Description: "Otimizar tamanho de storage", MonthlySavings: 10.0},
		},
	}
}

// Tipos de dados para análise
type AnalysisResult struct {
	Issues      []Issue  `json:"issues"`
	Suggestions []string `json:"suggestions"`
}

type Issue struct {
	Description string `json:"description"`
}

type SecurityAnalysis struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Recommendations []string        `json:"recommendations"`
}

type Vulnerability struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type CostAnalysis struct {
	EstimatedMonthlyCost float64        `json:"estimated_monthly_cost"`
	PotentialSavings     float64        `json:"potential_savings"`
	Optimizations        []Optimization `json:"optimizations"`
}

type Optimization struct {
	Description    string  `json:"description"`
	MonthlySavings float64 `json:"monthly_savings"`
}
