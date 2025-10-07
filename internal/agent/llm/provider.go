package llm

import (
	"context"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// Provider define a interface para provedores de LLM
type Provider interface {
	// Generate gera uma resposta de texto para o prompt fornecido
	Generate(req *models.LLMRequest) (*models.LLMResponse, error)

	// GenerateStructured gera uma resposta estruturada usando o LLM
	GenerateStructured(req *models.LLMRequest, responseStruct interface{}) error

	// GetCompletion obtém uma resposta para o prompt fornecido
	GetCompletion(ctx context.Context, prompt string) (string, error)

	// ValidateConnection testa a conexão com o provedor
	ValidateConnection() error
}

// NewLLMProvider cria um novo provedor de LLM baseado na configuração
func NewLLMProvider(cfg *config.Config, log *logger.Logger) (Provider, error) {
	if cfg.LLM.Provider == "" {
		// Default para Nation.fun se não especificado
		log.Info("Provedor LLM não especificado, usando Nation.fun como padrão")
		return NewNationClient(cfg, log)
	}

	// Cria provedor baseado na configuração
	switch cfg.LLM.Provider {
	case "nation", "nation.fun":
		return NewNationClient(cfg, log)
	case "openai":
		log.Warn("OpenAI não é mais suportado como provedor LLM, usando Nation.fun")
		return NewNationClient(cfg, log)
	case "anthropic":
		log.Warn("Anthropic não é mais suportado como provedor LLM, usando Nation.fun")
		return NewNationClient(cfg, log)
	default:
		log.Warn("Provedor LLM não suportado, usando Nation.fun", "provider", cfg.LLM.Provider)
		return NewNationClient(cfg, log)
	}
}
