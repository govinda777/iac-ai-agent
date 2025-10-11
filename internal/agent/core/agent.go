package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Agent representa o agente principal com capacidades modulares
type Agent struct {
	ID           string
	Name         string
	Description  string
	Version      string
	Capabilities map[string]Capability
	Config       *Config
	Logger       *Logger
	mu           sync.RWMutex
}

// Capability interface para habilidades do agente
type Capability interface {
	// Identificação da habilidade
	GetID() string
	GetName() string
	GetDescription() string
	GetVersion() string

	// Ciclo de vida
	Initialize(ctx context.Context, config *Config) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	// Processamento de mensagens
	CanHandle(message *Message) bool
	ProcessMessage(ctx context.Context, message *Message) (*Response, error)

	// Status e saúde
	GetStatus() *CapabilityStatus
	HealthCheck(ctx context.Context) error
}

// Message representa uma mensagem genérica
type Message struct {
	ID        string            `json:"id"`
	Source    string            `json:"source"`   // whatsapp, telegram, slack, etc.
	Channel   string            `json:"channel"`  // canal específico
	From      string            `json:"from"`     // remetente
	To        string            `json:"to"`       // destinatário
	Text      string            `json:"text"`     // conteúdo da mensagem
	Type      string            `json:"type"`     // text, image, document, etc.
	Metadata  map[string]string `json:"metadata"` // metadados adicionais
	Timestamp time.Time         `json:"timestamp"`
}

// Response representa uma resposta genérica
type Response struct {
	ID        string            `json:"id"`
	To        string            `json:"to"`
	Text      string            `json:"text"`
	Type      string            `json:"type"`
	Metadata  map[string]string `json:"metadata"`
	Timestamp time.Time         `json:"timestamp"`
}

// Config configuração do agente
type Config struct {
	AgentID      string                 `yaml:"agent_id"`
	AgentName    string                 `yaml:"agent_name"`
	Description  string                 `yaml:"description"`
	Version      string                 `yaml:"version"`
	Capabilities map[string]interface{} `yaml:"capabilities"`
	Logging      LoggingConfig          `yaml:"logging"`
	Web3         Web3Config             `yaml:"web3"`
	Billing      BillingConfig          `yaml:"billing"`
}

// LoggingConfig configuração de logging
type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    string `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// Web3Config configuração Web3
type Web3Config struct {
	WalletAddress string `yaml:"wallet_address"`
	NFTContract   string `yaml:"nft_contract"`
	LitProtocol   bool   `yaml:"lit_protocol"`
}

// BillingConfig configuração de billing
type BillingConfig struct {
	Enabled      bool           `yaml:"enabled"`
	TokenCosts   map[string]int `yaml:"token_costs"`
	FreeCommands []string       `yaml:"free_commands"`
}

// CapabilityStatus status de uma habilidade
type CapabilityStatus struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"` // active, inactive, error
	LastUsed     time.Time `json:"last_used"`
	MessageCount int       `json:"message_count"`
	ErrorCount   int       `json:"error_count"`
}

// Logger sistema de logging do agente
type Logger struct {
	AgentID string
	Level   string
}

// NewAgent cria um novo agente
func NewAgent(config *Config) *Agent {
	return &Agent{
		ID:           config.AgentID,
		Name:         config.AgentName,
		Description:  config.Description,
		Version:      config.Version,
		Capabilities: make(map[string]Capability),
		Config:       config,
		Logger:       NewLogger(config.AgentID, config.Logging.Level),
	}
}

// RegisterCapability registra uma nova habilidade
func (a *Agent) RegisterCapability(capability Capability) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	capID := capability.GetID()
	if _, exists := a.Capabilities[capID]; exists {
		return fmt.Errorf("capability %s already registered", capID)
	}

	a.Capabilities[capID] = capability
	a.Logger.Info("Capability registered", map[string]interface{}{
		"capability_id": capID,
		"name":          capability.GetName(),
		"version":       capability.GetVersion(),
	})

	return nil
}

// Initialize inicializa o agente e todas as habilidades
func (a *Agent) Initialize(ctx context.Context) error {
	a.Logger.Info("Initializing agent", map[string]interface{}{
		"agent_id":           a.ID,
		"name":               a.Name,
		"version":            a.Version,
		"capabilities_count": len(a.Capabilities),
	})

	// Inicializar cada habilidade
	for capID, capability := range a.Capabilities {
		if err := capability.Initialize(ctx, a.Config); err != nil {
			a.Logger.Error("Failed to initialize capability", map[string]interface{}{
				"capability_id": capID,
				"error":         err.Error(),
			})
			return fmt.Errorf("failed to initialize capability %s: %w", capID, err)
		}

		a.Logger.Info("Capability initialized", map[string]interface{}{
			"capability_id": capID,
		})
	}

	return nil
}

// Start inicia o agente e todas as habilidades
func (a *Agent) Start(ctx context.Context) error {
	a.Logger.Info("Starting agent", map[string]interface{}{
		"agent_id": a.ID,
	})

	// Iniciar cada habilidade
	for capID, capability := range a.Capabilities {
		if err := capability.Start(ctx); err != nil {
			a.Logger.Error("Failed to start capability", map[string]interface{}{
				"capability_id": capID,
				"error":         err.Error(),
			})
			return fmt.Errorf("failed to start capability %s: %w", capID, err)
		}

		a.Logger.Info("Capability started", map[string]interface{}{
			"capability_id": capID,
		})
	}

	return nil
}

// Stop para o agente e todas as habilidades
func (a *Agent) Stop(ctx context.Context) error {
	a.Logger.Info("Stopping agent", map[string]interface{}{
		"agent_id": a.ID,
	})

	// Parar cada habilidade
	for capID, capability := range a.Capabilities {
		if err := capability.Stop(ctx); err != nil {
			a.Logger.Error("Failed to stop capability", map[string]interface{}{
				"capability_id": capID,
				"error":         err.Error(),
			})
		}

		a.Logger.Info("Capability stopped", map[string]interface{}{
			"capability_id": capID,
		})
	}

	return nil
}

// ProcessMessage processa uma mensagem usando a habilidade apropriada
func (a *Agent) ProcessMessage(ctx context.Context, message *Message) (*Response, error) {
	a.Logger.Info("Processing message", map[string]interface{}{
		"message_id":  message.ID,
		"source":      message.Source,
		"from":        message.From,
		"text_length": len(message.Text),
	})

	// Encontrar habilidade capaz de processar a mensagem
	var selectedCapability Capability
	for _, capability := range a.Capabilities {
		if capability.CanHandle(message) {
			selectedCapability = capability
			break
		}
	}

	if selectedCapability == nil {
		return nil, fmt.Errorf("no capability can handle message from source: %s", message.Source)
	}

	// Processar mensagem com a habilidade selecionada
	response, err := selectedCapability.ProcessMessage(ctx, message)
	if err != nil {
		a.Logger.Error("Failed to process message", map[string]interface{}{
			"message_id":    message.ID,
			"capability_id": selectedCapability.GetID(),
			"error":         err.Error(),
		})
		return nil, err
	}

	a.Logger.Info("Message processed successfully", map[string]interface{}{
		"message_id":      message.ID,
		"capability_id":   selectedCapability.GetID(),
		"response_length": len(response.Text),
	})

	return response, nil
}

// GetCapabilities retorna lista de habilidades registradas
func (a *Agent) GetCapabilities() map[string]*CapabilityStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()

	statuses := make(map[string]*CapabilityStatus)
	for capID, capability := range a.Capabilities {
		statuses[capID] = capability.GetStatus()
	}

	return statuses
}

// GetCapability retorna uma habilidade específica
func (a *Agent) GetCapability(capID string) (Capability, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	capability, exists := a.Capabilities[capID]
	if !exists {
		return nil, fmt.Errorf("capability %s not found", capID)
	}

	return capability, nil
}

// HealthCheck verifica saúde de todas as habilidades
func (a *Agent) HealthCheck(ctx context.Context) map[string]error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	errors := make(map[string]error)
	for capID, capability := range a.Capabilities {
		if err := capability.HealthCheck(ctx); err != nil {
			errors[capID] = err
		}
	}

	return errors
}

// NewLogger cria novo logger
func NewLogger(agentID, level string) *Logger {
	return &Logger{
		AgentID: agentID,
		Level:   level,
	}
}

// Info logs mensagem informativa
func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log("INFO", message, fields)
}

// Error logs mensagem de erro
func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log("ERROR", message, fields)
}

// Debug logs mensagem de debug
func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.log("DEBUG", message, fields)
}

// Warn logs mensagem de aviso
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log("WARN", message, fields)
}

// log implementação interna de logging
func (l *Logger) log(level, message string, fields map[string]interface{}) {
	// Implementação simplificada - em produção usar logger estruturado
	fmt.Printf("[%s] [%s] %s", l.AgentID, level, message)
	for key, value := range fields {
		fmt.Printf(" %s=%v", key, value)
	}
	fmt.Println()
}