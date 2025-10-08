package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// Client representa o cliente para comunicação com a API do Notion
type Client struct {
	config     *config.Config
	logger     *logger.Logger
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewClient cria uma nova instância do cliente Notion
func NewClient(cfg *config.Config, log *logger.Logger) (*Client, error) {
	// Obtém a API key via Git Secrets
	apiKey, err := cfg.GetNotionAPIKey()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter API key do Notion: %w", err)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key do Notion não encontrada")
	}

	return &Client{
		config:     cfg,
		logger:     log,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     apiKey,
		baseURL:    cfg.Notion.BaseURL,
	}, nil
}

// Agent representa um agente no Notion
type Agent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Status       string    `json:"status"`
	Capabilities []string  `json:"capabilities"`
}

// CreateAgentRequest representa a requisição para criar um agente
type CreateAgentRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
}

// CreateAgentResponse representa a resposta da criação de agente
type CreateAgentResponse struct {
	Agent Agent `json:"agent"`
}

// ListAgentsResponse representa a resposta da listagem de agentes
type ListAgentsResponse struct {
	Agents []Agent `json:"agents"`
	Total  int     `json:"total"`
}

// ErrorResponse representa uma resposta de erro da API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// CreateAgent cria um novo agente no Notion
func (c *Client) CreateAgent(ctx context.Context, req *CreateAgentRequest) (*Agent, error) {
	c.logger.Info("Criando agente no Notion", "name", req.Name)

	// Prepara a requisição
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Faz a requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/agents", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Adiciona headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	// Executa a requisição
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verifica o status code
	if resp.StatusCode != http.StatusCreated {
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return nil, fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("erro ao criar agente: %s - %s", errorResp.Error, errorResp.Message)
	}

	// Parse da resposta
	var createResp CreateAgentResponse
	if err := json.Unmarshal(respBody, &createResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	c.logger.Info("Agente criado com sucesso", "agent_id", createResp.Agent.ID)
	return &createResp.Agent, nil
}

// ListAgents lista todos os agentes disponíveis
func (c *Client) ListAgents(ctx context.Context) ([]Agent, error) {
	c.logger.Info("Listando agentes do Notion")

	// Faz a requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/agents", nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Adiciona headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	// Executa a requisição
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verifica o status code
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return nil, fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("erro ao listar agentes: %s - %s", errorResp.Error, errorResp.Message)
	}

	// Parse da resposta
	var listResp ListAgentsResponse
	if err := json.Unmarshal(respBody, &listResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	c.logger.Info("Agentes listados com sucesso", "total", listResp.Total)
	return listResp.Agents, nil
}

// GetAgent obtém um agente específico por ID
func (c *Client) GetAgent(ctx context.Context, agentID string) (*Agent, error) {
	c.logger.Info("Obtendo agente do Notion", "agent_id", agentID)

	// Faz a requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/agents/"+agentID, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Adiciona headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	// Executa a requisição
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verifica o status code
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("agente não encontrado: %s", agentID)
		}
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return nil, fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("erro ao obter agente: %s - %s", errorResp.Error, errorResp.Message)
	}

	// Parse da resposta
	var agent Agent
	if err := json.Unmarshal(respBody, &agent); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	c.logger.Info("Agente obtido com sucesso", "agent_id", agent.ID)
	return &agent, nil
}

// FindAgentByName busca um agente pelo nome
func (c *Client) FindAgentByName(ctx context.Context, name string) (*Agent, error) {
	c.logger.Info("Buscando agente por nome", "name", name)

	agents, err := c.ListAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar agentes: %w", err)
	}

	for _, agent := range agents {
		if agent.Name == name {
			c.logger.Info("Agente encontrado", "agent_id", agent.ID, "name", agent.Name)
			return &agent, nil
		}
	}

	c.logger.Info("Agente não encontrado", "name", name)
	return nil, fmt.Errorf("agente não encontrado: %s", name)
}

// DeleteAgent remove um agente
func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	c.logger.Info("Removendo agente do Notion", "agent_id", agentID)

	// Faz a requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+"/agents/"+agentID, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Adiciona headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	// Executa a requisição
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Verifica o status code
	if resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("erro ao remover agente: %s - %s", errorResp.Error, errorResp.Message)
	}

	c.logger.Info("Agente removido com sucesso", "agent_id", agentID)
	return nil
}

// IsAgentAvailable verifica se o serviço de agentes está disponível
func (c *Client) IsAgentAvailable(ctx context.Context) bool {
	c.logger.Info("Verificando disponibilidade do serviço de agentes Notion")

	// Faz uma requisição simples para verificar se o serviço está disponível
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/agents", nil)
	if err != nil {
		c.logger.Error("Erro ao criar requisição de verificação", "error", err)
		return false
	}

	// Adiciona headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	// Executa a requisição
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Erro ao executar requisição de verificação", "error", err)
		return false
	}
	defer resp.Body.Close()

	// Considera disponível se retornar qualquer status válido (não erro de conexão)
	available := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
	c.logger.Info("Verificação de disponibilidade concluída", "status_code", resp.StatusCode, "available", available)

	return available
}
