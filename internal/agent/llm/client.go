package llm

import (
	"fmt"
	"time"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// Client é o cliente principal para interagir com LLMs
type Client struct {
	config   *config.Config
	logger   *logger.Logger
	provider LLMProvider
}

// NewClient cria um novo cliente LLM
func NewClient(cfg *config.Config, log *logger.Logger) *Client {
	client := &Client{
		config: cfg,
		logger: log,
	}

	// Inicializa o provedor LLM
	provider, err := NewLLMProvider(cfg, log)
	if err != nil {
		log.Error("Falha ao inicializar provedor LLM", "error", err)
		// Não retornamos erro aqui para permitir inicialização sem LLM
		// O erro será retornado quando Generate for chamado
	} else {
		client.provider = provider
		log.Info("Cliente LLM inicializado com sucesso", 
			"provider", cfg.LLM.Provider, 
			"model", cfg.LLM.Model)
	}

	return client
}

// Generate gera resposta do LLM baseado na requisição
func (c *Client) Generate(req *models.LLMRequest) (*models.LLMResponse, error) {
	// Valida se o provedor está inicializado
	if c.provider == nil {
		return nil, fmt.Errorf("provedor LLM não inicializado")
	}

	// Adiciona timestamp para métricas
	startTime := time.Now()

	// Adiciona defaults se necessário
	if req.SystemPrompt == "" {
		req.SystemPrompt = "Você é um especialista em Infrastructure as Code, segurança e otimização de cloud. Forneça análises detalhadas e recomendações práticas."
	}
	
	if req.MaxTokens == 0 {
		req.MaxTokens = 4000
	}
	
	if req.Temperature == 0 {
		req.Temperature = 0.2
	}

	// Log da requisição
	c.logger.Info("Gerando resposta LLM", 
		"provider", c.config.LLM.Provider,
		"model", c.config.LLM.Model,
		"prompt_length", len(req.Prompt),
		"max_tokens", req.MaxTokens)

	// Chama o provedor
	resp, err := c.provider.Generate(req)
	if err != nil {
		c.logger.Error("Erro ao gerar resposta LLM", "error", err)
		return nil, fmt.Errorf("erro ao gerar resposta LLM: %w", err)
	}

	// Adiciona latência
	resp.LatencyMs = time.Since(startTime).Milliseconds()

	// Log da resposta
	c.logger.Info("Resposta LLM recebida", 
		"latency_ms", resp.LatencyMs,
		"tokens_used", resp.TokensUsed,
		"response_length", len(resp.Content))

	return resp, nil
}

// GenerateStructured gera resposta estruturada parseada como JSON
func (c *Client) GenerateStructured(req *models.LLMRequest, target interface{}) error {
	// Valida se o provedor está inicializado
	if c.provider == nil {
		return fmt.Errorf("provedor LLM não inicializado")
	}

	// Log da requisição
	c.logger.Info("Gerando resposta estruturada", 
		"provider", c.config.LLM.Provider,
		"model", c.config.LLM.Model,
		"target_type", fmt.Sprintf("%T", target))

	// Chama o provedor
	err := c.provider.GenerateStructured(req, target)
	if err != nil {
		c.logger.Error("Erro ao gerar resposta estruturada", "error", err)
		return fmt.Errorf("erro ao gerar resposta estruturada: %w", err)
	}

	return nil
}

// ValidateConnection testa a conexão com o LLM
func (c *Client) ValidateConnection() error {
	// Valida se o provedor está inicializado
	if c.provider == nil {
		return fmt.Errorf("provedor LLM não inicializado")
	}

	// Chama o provedor
	err := c.provider.ValidateConnection()
	if err != nil {
		c.logger.Error("Falha na validação de conexão LLM", "error", err)
		return fmt.Errorf("falha na validação de conexão LLM: %w", err)
	}

	c.logger.Info("Conexão LLM validada com sucesso", 
		"provider", c.config.LLM.Provider,
		"model", c.config.LLM.Model)

	return nil
}
