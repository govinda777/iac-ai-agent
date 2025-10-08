package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
	"gopkg.in/yaml.v3"
)

// AgentService gerencia a criação e configuração de agentes
type AgentService struct {
	config       *config.Config
	logger       *logger.Logger
	baseClient   *web3.BaseClient
	litClient    *web3.LitProtocolClient
	tokenManager *web3.BotTokenManager
}

// NewAgentService cria um novo serviço de agentes
func NewAgentService(cfg *config.Config, log *logger.Logger, baseClient *web3.BaseClient) *AgentService {
	litClient := web3.NewLitProtocolClient(cfg, log, baseClient)
	tokenManager := web3.NewBotTokenManager(cfg, log, baseClient)

	return &AgentService{
		config:       cfg,
		logger:       log,
		baseClient:   baseClient,
		litClient:    litClient,
		tokenManager: tokenManager,
	}
}

// AgentTemplates contém a lista de templates disponíveis
type AgentTemplates struct {
	Templates           []AgentTemplate        `yaml:"templates"`
	DefaultConfig       map[string]interface{} `yaml:"default_config"`
	DefaultCapabilities map[string]bool        `yaml:"default_capabilities"`
	DefaultPersonality  map[string]interface{} `yaml:"default_personality"`
	DefaultKnowledge    map[string]interface{} `yaml:"default_knowledge"`
	DefaultLimits       map[string]interface{} `yaml:"default_limits"`
}

// AgentTemplate representa um template de agente
type AgentTemplate struct {
	ID            string   `yaml:"id"`
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	Category      string   `yaml:"category"`
	IsRecommended bool     `yaml:"is_recommended"`
	UseCases      []string `yaml:"use_cases"`
	Tags          []string `yaml:"tags"`
}

// Agent representa um agente configurado
type Agent struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	TemplateID      string                 `json:"template_id"`
	OwnerAddress    string                 `json:"owner_address"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Config          map[string]interface{} `json:"config"`
	Capabilities    map[string]bool        `json:"capabilities"`
	Personality     map[string]interface{} `json:"personality"`
	Knowledge       map[string]interface{} `json:"knowledge"`
	Limits          map[string]interface{} `json:"limits"`
	HasWhatsAppKey  bool                   `json:"has_whatsapp_key"`
	ContractAddress string                 `json:"contract_address"`
}

// LoadTemplates carrega os templates de agentes do arquivo de configuração
func (as *AgentService) LoadTemplates() (*AgentTemplates, error) {
	// Caminho para o arquivo de templates
	templatesPath := "configs/agent_templates.yaml"

	// Ler o arquivo
	data, err := ioutil.ReadFile(templatesPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de templates: %w", err)
	}

	// Parse do YAML
	var templates AgentTemplates
	if err := yaml.Unmarshal(data, &templates); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do YAML: %w", err)
	}

	return &templates, nil
}

// CreateAgentFromTemplate cria um novo agente a partir de um template
func (as *AgentService) CreateAgentFromTemplate(ctx context.Context, templateID string) (*Agent, error) {
	// Carregar templates
	templates, err := as.LoadTemplates()
	if err != nil {
		return nil, err
	}

	// Encontrar o template solicitado
	var selectedTemplate *AgentTemplate
	for _, tmpl := range templates.Templates {
		if tmpl.ID == templateID {
			selectedTemplate = &tmpl
			break
		}
	}

	if selectedTemplate == nil {
		// Se não encontrou o template específico, usa o template de WhatsApp por padrão
		for _, tmpl := range templates.Templates {
			if tmpl.ID == "whatsapp-chatbot" {
				selectedTemplate = &tmpl
				break
			}
		}

		// Se ainda não encontrou, usa o primeiro template disponível
		if selectedTemplate == nil && len(templates.Templates) > 0 {
			selectedTemplate = &templates.Templates[0]
		} else {
			return nil, fmt.Errorf("nenhum template disponível")
		}
	}

	// Verificar se a wallet tem um NFT da Nation.fun
	// Na implementação real, verificaríamos isso na blockchain
	// Para o MVP, assumimos que sim

	// Criar o agente
	agent := &Agent{
		ID:           uuid.New().String(),
		Name:         selectedTemplate.Name,
		Description:  selectedTemplate.Description,
		TemplateID:   selectedTemplate.ID,
		OwnerAddress: as.config.Web3.WalletAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Config:       templates.DefaultConfig,
		Capabilities: templates.DefaultCapabilities,
		Personality:  templates.DefaultPersonality,
		Knowledge:    templates.DefaultKnowledge,
		Limits:       templates.DefaultLimits,
	}

	// Verificar se já existe uma chave de WhatsApp armazenada
	hasWhatsAppKey := as.litClient.HasStoredWhatsAppAPIKey()
	agent.HasWhatsAppKey = hasWhatsAppKey

	// Gerar um endereço de contrato para o agente (simulado para o MVP)
	agent.ContractAddress = fmt.Sprintf("0x%s", uuid.New().String()[:40])

	// Salvar o agente
	if err := as.SaveAgent(agent); err != nil {
		return nil, err
	}

	// Atualizar o endereço do agente padrão na configuração
	as.config.Web3.DefaultAgentAddress = agent.ContractAddress

	as.logger.Info("Agente criado com sucesso",
		"agent_id", agent.ID,
		"template", agent.TemplateID,
		"owner", agent.OwnerAddress,
		"contract", agent.ContractAddress)

	return agent, nil
}

// SaveAgent salva um agente no armazenamento local
func (as *AgentService) SaveAgent(agent *Agent) error {
	// Criar diretório de agentes se não existir
	agentsDir := "data/agents"
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de agentes: %w", err)
	}

	// Caminho para o arquivo do agente
	agentPath := filepath.Join(agentsDir, fmt.Sprintf("%s.json", agent.ID))

	// Serializar o agente para JSON
	data, err := json.MarshalIndent(agent, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar agente: %w", err)
	}

	// Salvar o arquivo
	if err := ioutil.WriteFile(agentPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar arquivo do agente: %w", err)
	}

	return nil
}

// GetAgentByAddress busca um agente pelo endereço do contrato
func (as *AgentService) GetAgentByAddress(contractAddress string) (*Agent, error) {
	// Diretório de agentes
	agentsDir := "data/agents"

	// Listar arquivos no diretório
	files, err := ioutil.ReadDir(agentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("nenhum agente encontrado")
		}
		return nil, fmt.Errorf("erro ao listar agentes: %w", err)
	}

	// Procurar o agente pelo endereço do contrato
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Ler o arquivo do agente
		data, err := ioutil.ReadFile(filepath.Join(agentsDir, file.Name()))
		if err != nil {
			continue
		}

		// Deserializar o agente
		var agent Agent
		if err := json.Unmarshal(data, &agent); err != nil {
			continue
		}

		// Verificar se é o agente procurado
		if agent.ContractAddress == contractAddress {
			return &agent, nil
		}
	}

	return nil, fmt.Errorf("agente não encontrado: %s", contractAddress)
}

// GetAgentByOwner busca um agente pelo endereço do proprietário
func (as *AgentService) GetAgentByOwner(ownerAddress string) (*Agent, error) {
	// Diretório de agentes
	agentsDir := "data/agents"

	// Listar arquivos no diretório
	files, err := ioutil.ReadDir(agentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("nenhum agente encontrado")
		}
		return nil, fmt.Errorf("erro ao listar agentes: %w", err)
	}

	// Procurar o agente pelo endereço do proprietário
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Ler o arquivo do agente
		data, err := ioutil.ReadFile(filepath.Join(agentsDir, file.Name()))
		if err != nil {
			continue
		}

		// Deserializar o agente
		var agent Agent
		if err := json.Unmarshal(data, &agent); err != nil {
			continue
		}

		// Verificar se é o agente procurado
		if agent.OwnerAddress == ownerAddress {
			return &agent, nil
		}
	}

	return nil, fmt.Errorf("nenhum agente encontrado para o proprietário: %s", ownerAddress)
}

// StoreWhatsAppAPIKey armazena a chave de API do WhatsApp para o agente
func (as *AgentService) StoreWhatsAppAPIKey(agentID string, apiKey string) error {
	// Armazenar a chave usando o Lit Protocol
	_, err := as.litClient.StoreWhatsAppAPIKey(apiKey)
	if err != nil {
		return fmt.Errorf("erro ao armazenar chave de API do WhatsApp: %w", err)
	}

	// Atualizar o agente
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	agent.HasWhatsAppKey = true
	agent.UpdatedAt = time.Now()

	// Salvar o agente atualizado
	if err := as.SaveAgent(agent); err != nil {
		return err
	}

	as.logger.Info("Chave de API do WhatsApp armazenada com sucesso",
		"agent_id", agent.ID,
		"owner", agent.OwnerAddress)

	return nil
}

// GetWhatsAppAPIKey recupera a chave de API do WhatsApp para o agente
func (as *AgentService) GetWhatsAppAPIKey() (string, error) {
	// Recuperar a chave usando o Lit Protocol
	return as.litClient.GetWhatsAppAPIKey()
}

// GetAgentByID busca um agente pelo ID
func (as *AgentService) GetAgentByID(agentID string) (*Agent, error) {
	// Caminho para o arquivo do agente
	agentPath := filepath.Join("data/agents", fmt.Sprintf("%s.json", agentID))

	// Verificar se o arquivo existe
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("agente não encontrado: %s", agentID)
	}

	// Ler o arquivo do agente
	data, err := ioutil.ReadFile(agentPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo do agente: %w", err)
	}

	// Deserializar o agente
	var agent Agent
	if err := json.Unmarshal(data, &agent); err != nil {
		return nil, fmt.Errorf("erro ao deserializar agente: %w", err)
	}

	return &agent, nil
}

// EnsureAgentExists verifica se já existe um agente para o proprietário
// Se não existir, cria um novo agente com o template padrão
func (as *AgentService) EnsureAgentExists(ctx context.Context) (*Agent, error) {
	// Verificar se já existe um agente para o proprietário
	agent, err := as.GetAgentByOwner(as.config.Web3.WalletAddress)
	if err == nil {
		// Agente já existe
		as.logger.Info("Agente existente encontrado",
			"agent_id", agent.ID,
			"template", agent.TemplateID,
			"contract", agent.ContractAddress)

		// Atualizar o endereço do agente padrão na configuração
		as.config.Web3.DefaultAgentAddress = agent.ContractAddress

		return agent, nil
	}

	// Criar um novo agente com o template de WhatsApp
	return as.CreateAgentFromTemplate(ctx, "whatsapp-chatbot")
}

// DeductTokenForAnalysis debita um token do usuário para uma análise
func (as *AgentService) DeductTokenForAnalysis(ctx context.Context) error {
	// Calcular o custo em tokens para uma análise
	cost, err := as.tokenManager.CalculateTokenCost("terraform_analysis")
	if err != nil {
		return fmt.Errorf("erro ao calcular custo da análise: %w", err)
	}

	// Debitar os tokens
	if err := as.tokenManager.SpendTokens(ctx, as.config.Web3.WalletAddress, cost, "Análise de Terraform"); err != nil {
		return fmt.Errorf("erro ao debitar tokens: %w", err)
	}

	as.logger.Info("Tokens debitados com sucesso para análise",
		"wallet", as.config.Web3.WalletAddress,
		"cost", cost.String())

	return nil
}
