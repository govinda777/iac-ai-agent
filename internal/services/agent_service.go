package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// AgentService gerencia agentes de IA
type AgentService struct {
	logger    *logger.Logger
	agents    map[string]*models.Agent // In-memory storage (usar DB em produção)
	templates map[string]*models.AgentTemplate
}

// NewAgentService cria um novo serviço de agentes
func NewAgentService(log *logger.Logger) *AgentService {
	service := &AgentService{
		logger:    log,
		agents:    make(map[string]*models.Agent),
		templates: make(map[string]*models.AgentTemplate),
	}

	// Carrega templates padrão
	service.loadDefaultTemplates()

	return service
}

// loadDefaultTemplates carrega os templates pré-definidos
func (as *AgentService) loadDefaultTemplates() {
	as.templates = map[string]*models.AgentTemplate{
		"general-purpose":      as.createGeneralPurposeTemplate(),
		"security-specialist":  as.createSecuritySpecialistTemplate(),
		"cost-optimizer":       as.createCostOptimizerTemplate(),
		"architecture-advisor": as.createArchitectureAdvisorTemplate(),
	}

	as.logger.Info("Templates de agentes carregados", "count", len(as.templates))
}

// GetOrCreateDefaultAgent obtém o agente padrão ou cria um novo
func (as *AgentService) GetOrCreateDefaultAgent(ctx context.Context, ownerWallet string) (*models.Agent, error) {
	as.logger.Info("Verificando agente padrão", "owner", ownerWallet)

	// Procura agente existente do owner
	for _, agent := range as.agents {
		if agent.Owner == ownerWallet && agent.Status == "active" {
			as.logger.Info("Agente existente encontrado", "agent_id", agent.ID, "name", agent.Name)
			return agent, nil
		}
	}

	// Não encontrou, cria um novo automaticamente
	as.logger.Info("Nenhum agente encontrado, criando novo agente automaticamente...")

	req := &models.CreateAgentRequest{
		TemplateID:  "general-purpose",
		Name:        fmt.Sprintf("IaC Agent - %s", ownerWallet[:8]),
		Description: "Agente de IA criado automaticamente para análise de Infrastructure as Code",
		Owner:       ownerWallet,
	}

	agent, err := as.CreateAgent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar agente automaticamente: %w", err)
	}

	as.logger.Info("✅ Novo agente criado automaticamente",
		"agent_id", agent.ID,
		"name", agent.Name,
		"template", req.TemplateID)

	return agent, nil
}

// CreateAgent cria um novo agente
func (as *AgentService) CreateAgent(ctx context.Context, req *models.CreateAgentRequest) (*models.Agent, error) {
	// Valida requisição
	if req.Name == "" {
		return nil, fmt.Errorf("nome do agente é obrigatório")
	}

	if req.Owner == "" {
		return nil, fmt.Errorf("owner (wallet) é obrigatório")
	}

	// Busca template se especificado
	var template *models.AgentTemplate
	if req.TemplateID != "" {
		var ok bool
		template, ok = as.templates[req.TemplateID]
		if !ok {
			return nil, fmt.Errorf("template não encontrado: %s", req.TemplateID)
		}
	} else {
		// Usa template padrão
		template = as.templates["general-purpose"]
	}

	// Cria agente baseado no template
	agent := &models.Agent{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Version:     "1.0.0",
		Description: req.Description,
		Owner:       req.Owner,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      "active",

		Config:       template.DefaultConfig,
		Capabilities: template.DefaultCapabilities,
		Personality:  template.DefaultPersonality,
		Knowledge:    template.DefaultKnowledge,
		Limits:       template.DefaultLimits,

		Metrics: models.AgentMetrics{
			TotalRequests:      0,
			SuccessfulRequests: 0,
			FailedRequests:     0,
			TotalTokensUsed:    0,
			TotalCostUSD:       0,
			LastUsed:           time.Now(),
		},
	}

	// Aplica overrides se existirem
	// TODO: Implementar aplicação de overrides

	// Salva agente
	as.agents[agent.ID] = agent

	as.logger.Info("Agente criado com sucesso",
		"agent_id", agent.ID,
		"name", agent.Name,
		"owner", agent.Owner,
		"template", req.TemplateID)

	return agent, nil
}

// GetAgent obtém um agente por ID
func (as *AgentService) GetAgent(ctx context.Context, agentID string) (*models.Agent, error) {
	agent, ok := as.agents[agentID]
	if !ok {
		return nil, fmt.Errorf("agente não encontrado: %s", agentID)
	}

	return agent, nil
}

// ListAgents lista agentes de um owner
func (as *AgentService) ListAgents(ctx context.Context, ownerWallet string) ([]*models.Agent, error) {
	agents := []*models.Agent{}

	for _, agent := range as.agents {
		if agent.Owner == ownerWallet {
			agents = append(agents, agent)
		}
	}

	as.logger.Info("Agentes listados", "owner", ownerWallet, "count", len(agents))
	return agents, nil
}

// UpdateAgent atualiza um agente
func (as *AgentService) UpdateAgent(ctx context.Context, agentID string, req *models.AgentUpdateRequest) (*models.Agent, error) {
	agent, ok := as.agents[agentID]
	if !ok {
		return nil, fmt.Errorf("agente não encontrado: %s", agentID)
	}

	// Aplica updates
	if req.Name != nil {
		agent.Name = *req.Name
	}
	if req.Description != nil {
		agent.Description = *req.Description
	}
	if req.Config != nil {
		agent.Config = *req.Config
	}
	if req.Capabilities != nil {
		agent.Capabilities = *req.Capabilities
	}
	if req.Personality != nil {
		agent.Personality = *req.Personality
	}
	if req.Knowledge != nil {
		agent.Knowledge = *req.Knowledge
	}
	if req.Limits != nil {
		agent.Limits = *req.Limits
	}
	if req.Status != nil {
		agent.Status = *req.Status
	}

	agent.UpdatedAt = time.Now()

	as.logger.Info("Agente atualizado", "agent_id", agentID)
	return agent, nil
}

// DeleteAgent deleta um agente
func (as *AgentService) DeleteAgent(ctx context.Context, agentID string) error {
	if _, ok := as.agents[agentID]; !ok {
		return fmt.Errorf("agente não encontrado: %s", agentID)
	}

	delete(as.agents, agentID)
	as.logger.Info("Agente deletado", "agent_id", agentID)
	return nil
}

// GetTemplates lista templates disponíveis
func (as *AgentService) GetTemplates(ctx context.Context) []*models.AgentTemplate {
	templates := make([]*models.AgentTemplate, 0, len(as.templates))
	for _, template := range as.templates {
		templates = append(templates, template)
	}
	return templates
}

// ====================================================================
// TEMPLATES PRÉ-DEFINIDOS
// ====================================================================

// createGeneralPurposeTemplate cria template de agente de propósito geral
func (as *AgentService) createGeneralPurposeTemplate() *models.AgentTemplate {
	return &models.AgentTemplate{
		ID:            "general-purpose",
		Name:          "General Purpose IaC Agent",
		Description:   "Agente versátil para análise completa de Infrastructure as Code",
		Category:      "general",
		IsRecommended: true,
		UseCases: []string{
			"Análise completa de Terraform",
			"Review de Pull Requests",
			"Detecção de problemas de segurança",
			"Sugestões de otimização",
			"Best practices",
		},
		Tags: []string{"terraform", "aws", "azure", "gcp", "security", "cost"},

		DefaultConfig: models.AgentConfig{
			LLMProvider:      "openai",
			LLMModel:         "gpt-4",
			Temperature:      0.2,
			MaxTokens:        4000,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,

			EnableCheckov:         true,
			EnableIAMAnalysis:     true,
			EnableCostAnalysis:    true,
			EnableDriftDetection:  true,
			EnablePreviewAnalysis: true,
			EnableSecretsScanning: true,

			ResponseFormat:      "json",
			IncludeCodeExamples: true,
			IncludeReferences:   true,
			DetailLevel:         "standard",
			Language:            "pt-br",
			Timezone:            "America/Sao_Paulo",
		},

		DefaultCapabilities: models.AgentCapabilities{
			CanAnalyzeTerraform: true,
			CanAnalyzeCheckov:   true,
			CanAnalyzeIAM:       true,
			CanAnalyzeCosts:     true,
			CanDetectDrift:      true,
			CanAnalyzePreview:   true,
			CanScanSecrets:      true,

			CanGenerateCode:          true,
			CanGenerateTests:         false,
			CanGenerateDocumentation: true,
			CanRefactorCode:          true,

			CanSuggestArchitecture:  true,
			CanSuggestModules:       true,
			CanSuggestOptimizations: true,
			CanSuggestSecurity:      true,

			CanIntegrateGitHub:  true,
			CanIntegrateCI:      false,
			CanIntegrateSlack:   false,
			CanIntegrateDiscord: false,

			CanLearnFromFeedback:   false,
			CanAdaptToContext:      true,
			CanRememberPreferences: false,
		},

		DefaultPersonality: models.AgentPersonality{
			Style:         "professional",
			Tone:          "encouraging",
			Verbosity:     "balanced",
			UseEmojis:     true,
			UseHumor:      false,
			BeEncouraging: true,
			BeDirective:   false,

			ExplainReasoning:    true,
			ProvideExamples:     true,
			CompareAlternatives: true,
			HighlightRisks:      true,

			AskClarifyingQuestions: false,
			OfferAlternatives:      true,
			SuggestBestPractices:   true,
		},

		DefaultKnowledge: models.AgentKnowledge{
			TerraformExpertise: "expert",
			AWSExpertise:       "expert",
			AzureExpertise:     "intermediate",
			GCPExpertise:       "intermediate",

			SecurityExpertise:   "expert",
			NetworkingExpertise: "intermediate",
			KubernetesExpertise: "intermediate",
			DatabaseExpertise:   "intermediate",

			ComplianceFrameworks: []string{"GDPR", "SOC2", "HIPAA", "PCI-DSS"},
			IndustryFocus:        []string{"general"},
			ArchitecturePatterns: []string{"3-tier", "microservices", "serverless"},

			CustomRules:      []models.CustomRule{},
			PreferredModules: []string{},
			BannedResources:  []string{},
			RequiredTags:     []string{"Environment", "Owner", "ManagedBy"},
		},

		DefaultLimits: models.AgentLimits{
			MaxRequestsPerHour:    100,
			MaxRequestsPerDay:     1000,
			MaxConcurrentRequests: 5,

			MaxTokensPerRequest: 4000,
			MaxTokensPerDay:     100000,

			MaxFilesPerAnalysis: 50,
			MaxFileSizeMB:       10,
			MaxResourcesPerFile: 200,

			MaxCostPerRequest: 0.50,   // USD
			MaxCostPerDay:     10.00,  // USD
			MaxCostPerMonth:   200.00, // USD

			MaxAnalysisTimeSeconds: 300, // 5 minutes
			RequestTimeoutSeconds:  60,  // 1 minute
		},
	}
}

// createSecuritySpecialistTemplate cria template focado em segurança
func (as *AgentService) createSecuritySpecialistTemplate() *models.AgentTemplate {
	template := as.createGeneralPurposeTemplate()

	template.ID = "security-specialist"
	template.Name = "Security Specialist Agent"
	template.Description = "Agente especializado em análise de segurança e compliance"
	template.Category = "security"
	template.Tags = []string{"security", "compliance", "checkov", "secrets", "iam"}

	// Ajusta configuração para foco em segurança
	template.DefaultConfig.EnableCheckov = true
	template.DefaultConfig.EnableSecretsScanning = true
	template.DefaultConfig.EnableIAMAnalysis = true
	template.DefaultConfig.EnableCostAnalysis = false
	template.DefaultConfig.DetailLevel = "detailed"

	// Ajusta capabilities
	template.DefaultCapabilities.CanAnalyzeCheckov = true
	template.DefaultCapabilities.CanAnalyzeIAM = true
	template.DefaultCapabilities.CanScanSecrets = true
	template.DefaultCapabilities.CanSuggestSecurity = true

	// Ajusta personalidade para ser mais direto em segurança
	template.DefaultPersonality.BeDirective = true
	template.DefaultPersonality.HighlightRisks = true
	template.DefaultPersonality.Tone = "formal"

	// Ajusta conhecimento
	template.DefaultKnowledge.SecurityExpertise = "expert"
	template.DefaultKnowledge.ComplianceFrameworks = []string{"GDPR", "SOC2", "HIPAA", "PCI-DSS", "ISO27001"}

	return template
}

// createCostOptimizerTemplate cria template focado em otimização de custos
func (as *AgentService) createCostOptimizerTemplate() *models.AgentTemplate {
	template := as.createGeneralPurposeTemplate()

	template.ID = "cost-optimizer"
	template.Name = "Cost Optimizer Agent"
	template.Description = "Agente especializado em análise e otimização de custos"
	template.Category = "cost"
	template.Tags = []string{"cost", "optimization", "savings", "pricing"}

	// Ajusta configuração
	template.DefaultConfig.EnableCostAnalysis = true
	template.DefaultConfig.EnableCheckov = false
	template.DefaultConfig.DetailLevel = "detailed"

	// Ajusta capabilities
	template.DefaultCapabilities.CanAnalyzeCosts = true
	template.DefaultCapabilities.CanSuggestOptimizations = true

	// Ajusta personalidade
	template.DefaultPersonality.ProvideExamples = true
	template.DefaultPersonality.CompareAlternatives = true

	return template
}

// createArchitectureAdvisorTemplate cria template focado em arquitetura
func (as *AgentService) createArchitectureAdvisorTemplate() *models.AgentTemplate {
	template := as.createGeneralPurposeTemplate()

	template.ID = "architecture-advisor"
	template.Name = "Architecture Advisor Agent"
	template.Description = "Agente especializado em análise e sugestões arquiteturais"
	template.Category = "architecture"
	template.Tags = []string{"architecture", "patterns", "best-practices", "design"}

	// Ajusta capabilities
	template.DefaultCapabilities.CanSuggestArchitecture = true
	template.DefaultCapabilities.CanSuggestModules = true
	template.DefaultCapabilities.CanRefactorCode = true

	// Ajusta personalidade
	template.DefaultPersonality.Verbosity = "verbose"
	template.DefaultPersonality.ExplainReasoning = true
	template.DefaultPersonality.CompareAlternatives = true

	// Ajusta conhecimento
	template.DefaultKnowledge.ArchitecturePatterns = []string{
		"3-tier", "microservices", "serverless", "event-driven",
		"cqrs", "saga", "sidecar", "strangler",
	}

	return template
}
