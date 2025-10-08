package services

import (
	"context"
	"fmt"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/platform/notion"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// NotionAgentService gerencia agentes no Notion
type NotionAgentService struct {
	config *config.Config
	logger *logger.Logger
	client *notion.Client
}

// NewNotionAgentService cria uma nova instância do serviço de agentes Notion
func NewNotionAgentService(cfg *config.Config, log *logger.Logger) (*NotionAgentService, error) {
	client, err := notion.NewClient(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente Notion: %w", err)
	}

	return &NotionAgentService{
		config: cfg,
		logger: log,
		client: client,
	}, nil
}

// AgentInfo representa informações do agente
type AgentInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Available   bool      `json:"available"`
}

// CreateAgentRequest representa a requisição para criar um agente
type CreateAgentRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
}

// GetOrCreateDefaultAgent obtém o agente padrão ou cria um novo se não existir
func (s *NotionAgentService) GetOrCreateDefaultAgent(ctx context.Context) (*AgentInfo, error) {
	s.logger.Info("Obtendo ou criando agente padrão do Notion")

	// Verifica se o serviço está disponível
	if !s.client.IsAgentAvailable(ctx) {
		s.logger.Warn("Serviço de agentes Notion não está disponível")
		return nil, fmt.Errorf("serviço de agentes Notion não está disponível")
	}

	// Tenta encontrar o agente pelo nome padrão
	agentName := s.config.Notion.AgentName
	agent, err := s.client.FindAgentByName(ctx, agentName)
	if err == nil {
		s.logger.Info("Agente padrão encontrado", "agent_id", agent.ID, "name", agent.Name)
		return &AgentInfo{
			ID:          agent.ID,
			Name:        agent.Name,
			Description: agent.Description,
			Status:      agent.Status,
			CreatedAt:   agent.CreatedAt,
			UpdatedAt:   agent.UpdatedAt,
			Available:   true,
		}, nil
	}

	// Se não encontrou, verifica se deve criar automaticamente
	if !s.config.Notion.EnableAgentCreation {
		s.logger.Warn("Criação de agente desabilitada e agente não encontrado", "agent_name", agentName)
		return nil, fmt.Errorf("agente '%s' não encontrado e criação automática está desabilitada", agentName)
	}

	// Cria o agente
	s.logger.Info("Criando novo agente padrão", "agent_name", agentName)
	createReq := &notion.CreateAgentRequest{
		Name:         agentName,
		Description:  s.config.Notion.AgentDescription,
		Capabilities: []string{"terraform_analysis", "security_review", "cost_optimization", "iam_analysis"},
	}

	newAgent, err := s.client.CreateAgent(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar agente: %w", err)
	}

	s.logger.Info("Agente padrão criado com sucesso", "agent_id", newAgent.ID, "name", newAgent.Name)
	return &AgentInfo{
		ID:          newAgent.ID,
		Name:        newAgent.Name,
		Description: newAgent.Description,
		Status:      newAgent.Status,
		CreatedAt:   newAgent.CreatedAt,
		UpdatedAt:   newAgent.UpdatedAt,
		Available:   true,
	}, nil
}

// CreateAgent cria um novo agente personalizado
func (s *NotionAgentService) CreateAgent(ctx context.Context, req *CreateAgentRequest) (*AgentInfo, error) {
	s.logger.Info("Criando novo agente personalizado", "name", req.Name)

	// Verifica se o serviço está disponível
	if !s.client.IsAgentAvailable(ctx) {
		return nil, fmt.Errorf("serviço de agentes Notion não está disponível")
	}

	// Verifica se já existe um agente com o mesmo nome
	existingAgent, err := s.client.FindAgentByName(ctx, req.Name)
	if err == nil {
		s.logger.Warn("Agente já existe", "agent_id", existingAgent.ID, "name", req.Name)
		return &AgentInfo{
			ID:          existingAgent.ID,
			Name:        existingAgent.Name,
			Description: existingAgent.Description,
			Status:      existingAgent.Status,
			CreatedAt:   existingAgent.CreatedAt,
			UpdatedAt:   existingAgent.UpdatedAt,
			Available:   true,
		}, nil
	}

	// Cria o agente
	createReq := &notion.CreateAgentRequest{
		Name:         req.Name,
		Description:  req.Description,
		Capabilities: req.Capabilities,
	}

	newAgent, err := s.client.CreateAgent(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar agente: %w", err)
	}

	s.logger.Info("Agente personalizado criado com sucesso", "agent_id", newAgent.ID, "name", newAgent.Name)
	return &AgentInfo{
		ID:          newAgent.ID,
		Name:        newAgent.Name,
		Description: newAgent.Description,
		Status:      newAgent.Status,
		CreatedAt:   newAgent.CreatedAt,
		UpdatedAt:   newAgent.UpdatedAt,
		Available:   true,
	}, nil
}

// GetAgent obtém um agente específico por ID
func (s *NotionAgentService) GetAgent(ctx context.Context, agentID string) (*AgentInfo, error) {
	s.logger.Info("Obtendo agente", "agent_id", agentID)

	// Verifica se o serviço está disponível
	if !s.client.IsAgentAvailable(ctx) {
		return nil, fmt.Errorf("serviço de agentes Notion não está disponível")
	}

	agent, err := s.client.GetAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter agente: %w", err)
	}

	return &AgentInfo{
		ID:          agent.ID,
		Name:        agent.Name,
		Description: agent.Description,
		Status:      agent.Status,
		CreatedAt:   agent.CreatedAt,
		UpdatedAt:   agent.UpdatedAt,
		Available:   true,
	}, nil
}

// ListAgents lista todos os agentes disponíveis
func (s *NotionAgentService) ListAgents(ctx context.Context) ([]AgentInfo, error) {
	s.logger.Info("Listando agentes do Notion")

	// Verifica se o serviço está disponível
	if !s.client.IsAgentAvailable(ctx) {
		return nil, fmt.Errorf("serviço de agentes Notion não está disponível")
	}

	agents, err := s.client.ListAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar agentes: %w", err)
	}

	var agentInfos []AgentInfo
	for _, agent := range agents {
		agentInfos = append(agentInfos, AgentInfo{
			ID:          agent.ID,
			Name:        agent.Name,
			Description: agent.Description,
			Status:      agent.Status,
			CreatedAt:   agent.CreatedAt,
			UpdatedAt:   agent.UpdatedAt,
			Available:   true,
		})
	}

	s.logger.Info("Agentes listados com sucesso", "total", len(agentInfos))
	return agentInfos, nil
}

// DeleteAgent remove um agente
func (s *NotionAgentService) DeleteAgent(ctx context.Context, agentID string) error {
	s.logger.Info("Removendo agente", "agent_id", agentID)

	// Verifica se o serviço está disponível
	if !s.client.IsAgentAvailable(ctx) {
		return fmt.Errorf("serviço de agentes Notion não está disponível")
	}

	err := s.client.DeleteAgent(ctx, agentID)
	if err != nil {
		return fmt.Errorf("erro ao remover agente: %w", err)
	}

	s.logger.Info("Agente removido com sucesso", "agent_id", agentID)
	return nil
}

// IsServiceAvailable verifica se o serviço de agentes está disponível
func (s *NotionAgentService) IsServiceAvailable(ctx context.Context) bool {
	return s.client.IsAgentAvailable(ctx)
}

// ValidateConfiguration valida a configuração do Notion
func (s *NotionAgentService) ValidateConfiguration() error {
	if s.config.Notion.EnableAgentCreation && s.config.Notion.APIKey == "" {
		return fmt.Errorf("NOTION_API_KEY é obrigatório quando enable_agent_creation está habilitado")
	}

	if s.config.Notion.AgentName == "" {
		return fmt.Errorf("NOTION_AGENT_NAME é obrigatório")
	}

	return nil
}

// GetDefaultAgentCapabilities retorna as capacidades padrão do agente
func (s *NotionAgentService) GetDefaultAgentCapabilities() []string {
	return []string{
		"terraform_analysis",
		"security_review",
		"cost_optimization",
		"iam_analysis",
		"infrastructure_review",
		"compliance_check",
		"best_practices",
		"documentation_generation",
	}
}
