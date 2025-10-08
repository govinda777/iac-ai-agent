package integration

import (
	"context"
	"testing"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/services"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// TestNotionIntegration testa a integração básica com Notion
func TestNotionIntegration(t *testing.T) {
	// Configuração de teste
	cfg := &config.Config{
		Notion: config.NotionConfig{
			BaseURL:              "https://api.notion.com/v1",
			AgentName:            "Test IaC AI Agent",
			AgentDescription:     "Test Agent for IaC AI",
			EnableAgentCreation:  true,
			AutoCreateOnStartup:  true,
			MaxRequestsPerMinute: 60,
		},
	}

	// Logger de teste
	log := logger.New("info", "text")

	// Cria serviço Notion
	notionService, err := services.NewNotionAgentService(cfg, log)
	if err != nil {
		t.Skipf("Notion service não pode ser criado (provavelmente falta API key): %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Testa disponibilidade do serviço
	t.Run("ServiceAvailability", func(t *testing.T) {
		available := notionService.IsServiceAvailable(ctx)
		if !available {
			t.Skip("Serviço Notion não está disponível")
		}
		t.Log("Serviço Notion está disponível")
	})

	// Testa validação de configuração
	t.Run("ConfigurationValidation", func(t *testing.T) {
		err := notionService.ValidateConfiguration()
		if err != nil {
			t.Errorf("Validação de configuração falhou: %v", err)
		}
		t.Log("Configuração validada com sucesso")
	})

	// Testa obtenção de capacidades padrão
	t.Run("DefaultCapabilities", func(t *testing.T) {
		capabilities := notionService.GetDefaultAgentCapabilities()
		if len(capabilities) == 0 {
			t.Error("Nenhuma capacidade padrão encontrada")
		}
		t.Logf("Capacidades padrão: %v", capabilities)
	})

	// Testa criação/obtenção de agente padrão (apenas se serviço estiver disponível)
	t.Run("GetOrCreateDefaultAgent", func(t *testing.T) {
		if !notionService.IsServiceAvailable(ctx) {
			t.Skip("Serviço Notion não está disponível")
		}

		agent, err := notionService.GetOrCreateDefaultAgent(ctx)
		if err != nil {
			t.Errorf("Erro ao obter/criar agente padrão: %v", err)
			return
		}

		if agent.ID == "" {
			t.Error("ID do agente não foi retornado")
		}

		if agent.Name == "" {
			t.Error("Nome do agente não foi retornado")
		}

		t.Logf("Agente obtido/criado: ID=%s, Name=%s", agent.ID, agent.Name)
	})
}

// TestNotionConfiguration testa a configuração do Notion
func TestNotionConfiguration(t *testing.T) {
	// Testa configuração padrão sem validação obrigatória
	cfg := &config.Config{
		Notion: config.NotionConfig{
			BaseURL:              "https://api.notion.com/v1",
			AgentName:            "Test Agent",
			AgentDescription:     "Test Description",
			EnableAgentCreation:  true,
			AutoCreateOnStartup:  true,
			MaxRequestsPerMinute: 60,
		},
	}

	// Verifica se configuração Notion foi inicializada
	if cfg.Notion.BaseURL == "" {
		t.Error("BaseURL do Notion não foi inicializada")
	}

	if cfg.Notion.AgentName == "" {
		t.Error("AgentName do Notion não foi inicializada")
	}

	if cfg.Notion.MaxRequestsPerMinute == 0 {
		t.Error("MaxRequestsPerMinute do Notion não foi inicializada")
	}

	t.Logf("Configuração Notion: BaseURL=%s, AgentName=%s, MaxRequestsPerMinute=%d",
		cfg.Notion.BaseURL, cfg.Notion.AgentName, cfg.Notion.MaxRequestsPerMinute)
}

// TestNotionClientCreation testa a criação do cliente Notion
func TestNotionClientCreation(t *testing.T) {
	cfg := &config.Config{
		Notion: config.NotionConfig{
			BaseURL: "https://api.notion.com/v1",
			// APIKey não definida intencionalmente para testar erro
		},
	}

	log := logger.New("info", "text")

	// Deve falhar sem API key
	_, err := services.NewNotionAgentService(cfg, log)
	if err == nil {
		t.Error("Esperava erro ao criar serviço sem API key")
	}

	t.Logf("Erro esperado ao criar serviço sem API key: %v", err)
}
