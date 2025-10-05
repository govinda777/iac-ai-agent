package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// Client é o cliente para interagir com LLMs
type Client struct {
	config     *config.Config
	logger     *logger.Logger
	httpClient *http.Client
	provider   string
	apiKey     string
	model      string
}

// NewClient cria um novo cliente LLM
func NewClient(cfg *config.Config, log *logger.Logger) *Client {
	return &Client{
		config:     cfg,
		logger:     log,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		provider:   cfg.LLM.Provider,
		apiKey:     cfg.LLM.APIKey,
		model:      cfg.LLM.Model,
	}
}

// Generate gera resposta do LLM baseado na requisição
func (c *Client) Generate(req *models.LLMRequest) (*models.LLMResponse, error) {
	switch c.provider {
	case "openai":
		return c.generateOpenAI(req)
	case "anthropic":
		return c.generateAnthropic(req)
	default:
		return nil, fmt.Errorf("provedor LLM não suportado: %s", c.provider)
	}
}

// generateOpenAI gera resposta usando OpenAI API
func (c *Client) generateOpenAI(req *models.LLMRequest) (*models.LLMResponse, error) {
	url := "https://api.openai.com/v1/chat/completions"

	payload := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Você é um especialista em Infrastructure as Code, segurança e otimização de cloud. Forneça análises detalhadas e recomendações práticas.",
			},
			{
				"role":    "user",
				"content": req.Prompt,
			},
		},
		"temperature": req.Temperature,
	}

	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	c.logger.Info("Chamando OpenAI API", "model", c.model)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou erro %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("nenhuma resposta retornada pela API")
	}

	return &models.LLMResponse{
		Content:      openAIResp.Choices[0].Message.Content,
		TokensUsed:   openAIResp.Usage.TotalTokens,
		Model:        openAIResp.Model,
		FinishReason: openAIResp.Choices[0].FinishReason,
		Metadata:     make(map[string]interface{}),
	}, nil
}

// generateAnthropic gera resposta usando Anthropic API (Claude)
func (c *Client) generateAnthropic(req *models.LLMRequest) (*models.LLMResponse, error) {
	url := "https://api.anthropic.com/v1/messages"

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	payload := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": req.Prompt,
			},
		},
		"max_tokens":  maxTokens,
		"temperature": req.Temperature,
		"system":      "Você é um especialista em Infrastructure as Code, segurança e otimização de cloud. Forneça análises detalhadas e recomendações práticas.",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	c.logger.Info("Chamando Anthropic API", "model", c.model)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou erro %d: %s", resp.StatusCode, string(body))
	}

	var anthropicResp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		StopReason string `json:"stop_reason"`
		Usage      struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return nil, fmt.Errorf("nenhuma resposta retornada pela API")
	}

	totalTokens := anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens

	return &models.LLMResponse{
		Content:      anthropicResp.Content[0].Text,
		TokensUsed:   totalTokens,
		Model:        anthropicResp.Model,
		FinishReason: anthropicResp.StopReason,
		Metadata:     make(map[string]interface{}),
	}, nil
}

// GenerateStructured gera resposta estruturada parseada como JSON
func (c *Client) GenerateStructured(req *models.LLMRequest, target interface{}) error {
	// Adiciona instrução para responder em JSON
	req.Prompt += "\n\nResponda APENAS com JSON válido, sem texto adicional."

	resp, err := c.Generate(req)
	if err != nil {
		return err
	}

	// Parse JSON da resposta
	if err := json.Unmarshal([]byte(resp.Content), target); err != nil {
		return fmt.Errorf("erro ao fazer parse da resposta estruturada: %w", err)
	}

	return nil
}
