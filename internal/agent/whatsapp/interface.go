package whatsapp

import (
	"context"
	"time"
)

// WhatsAppAgentInterface defines the interface for the WhatsApp agent.
type WhatsAppAgentInterface interface {
	ProcessMessage(context.Context, *WhatsAppMessage) (*WhatsAppResponse, error)
	GetVerifyToken() string
}

// WhatsAppMessage represents a message received from WhatsApp.
type WhatsAppMessage struct {
	ID        string
	From      string
	To        string
	Text      string
	Timestamp time.Time
	Type      string
}

// WhatsAppResponse represents a message to be sent to WhatsApp.
type WhatsAppResponse struct {
	To      string
	Text    string
	Type    string
	ReplyTo string
}

// WhatsAppAgentConfig configuração para criação de agente WhatsApp
type WhatsAppAgentConfig struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	WalletAddr  string `json:"wallet_address" yaml:"wallet_address"`
	WebhookURL  string `json:"webhook_url" yaml:"webhook_url"`
	VerifyToken string `json:"verify_token" yaml:"verify_token"`
}

// Command representa um comando disponível no agente
type Command struct {
	Name            string
	Description     string
	Pattern         string
	Handler         func(*WhatsAppAgent, *CommandContext) (*WhatsAppResponse, error)
	RequiresPayment bool
	TokenCost       int
}

// CommandContext contexto para execução de comandos
type CommandContext struct {
	Message   *WhatsAppMessage
	Agent     *WhatsAppAgent
	Args      []string
	CodeBlock string
}

// UsageStats estatísticas de uso do usuário
type UsageStats struct {
	TotalRequests  int       `json:"total_requests"`
	TokensConsumed int       `json:"tokens_consumed"`
	LastRequest    time.Time `json:"last_request"`
	RequestsToday  int       `json:"requests_today"`
	AverageCost    float64   `json:"average_cost"`
}