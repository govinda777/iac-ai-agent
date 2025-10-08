package whatsapp

import (
	"time"
)

// WhatsAppMessage representa uma mensagem recebida do WhatsApp
type WhatsAppMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // text, image, document, etc.
}

// WhatsAppResponse representa uma resposta a ser enviada via WhatsApp
type WhatsAppResponse struct {
	To      string `json:"to"`
	Text    string `json:"text"`
	Type    string `json:"type"` // text, image, document, etc.
	ReplyTo string `json:"reply_to,omitempty"`
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

// TokenTransaction representa uma transação de tokens
type TokenTransaction struct {
	UserAddress     string    `json:"user_address"`
	AgentAddress    string    `json:"agent_address"`
	Amount          int       `json:"amount"`
	TransactionHash string    `json:"transaction_hash"`
	Timestamp       time.Time `json:"timestamp"`
	Service         string    `json:"service"`
}

// APIKeyStorage estrutura para armazenamento seguro de chaves API
type APIKeyStorage struct {
	EncryptedAPIKey []byte    `json:"encrypted_api_key"`
	EncryptedAESKey []byte    `json:"encrypted_aes_key"`
	WalletAddress   string    `json:"wallet_address"`
	Timestamp       time.Time `json:"timestamp"`
}
