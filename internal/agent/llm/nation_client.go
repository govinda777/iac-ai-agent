package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// NationClient implementa a interface LLMProvider para Nation.fun
type NationClient struct {
	config      *config.Config
	logger      *logger.Logger
	httpClient  *http.Client
	apiKey      string
	nftAddress  string
	baseURL     string
	modelName   string
	walletToken string
}

// NewNationClient cria um novo cliente Nation.fun
func NewNationClient(cfg *config.Config, log *logger.Logger) (*NationClient, error) {
	if cfg.LLM.APIKey == "" {
		return nil, errors.New("Nation.fun API key não configurada")
	}

	if cfg.Web3.NFTAccessContractAddress == "" {
		return nil, errors.New("Nation.fun NFT contract address não configurado")
	}

	if cfg.Web3.WalletToken == "" {
		return nil, errors.New("Nation.fun wallet token não configurado")
	}

	// Valida URL base
	baseURL := cfg.LLM.BaseURL
	if baseURL == "" {
		baseURL = "https://api.nation.fun/v1"
		log.Debug("URL base Nation.fun não especificada, usando padrão", "url", baseURL)
	}

	// Valida modelo
	modelName := cfg.LLM.Model
	if modelName == "" {
		modelName = "nation-1" // Modelo padrão
		log.Warn("Modelo Nation.fun não especificado, usando padrão", "model", modelName)
	}

	// Define agente padrão caso não seja especificado
	defaultAgentAddress := "0x147e832418Cc06A501047019E956714271098b89"
	if cfg.Web3.DefaultAgentAddress == "" {
		cfg.Web3.DefaultAgentAddress = defaultAgentAddress
		log.Info("Usando agente Nation.fun padrão", "agent", defaultAgentAddress)
	}

	client := &NationClient{
		config:      cfg,
		logger:      log,
		httpClient:  &http.Client{Timeout: 60 * time.Second},
		apiKey:      cfg.LLM.APIKey,
		nftAddress:  cfg.Web3.NFTAccessContractAddress,
		baseURL:     baseURL,
		modelName:   modelName,
		walletToken: cfg.Web3.WalletToken,
	}

	// Testa conexão
	err := client.ValidateConnection()
	if err != nil {
		// Se a validação falhar e tivermos a chave privada, podemos tentar gerar um token
		if cfg.Web3.WalletToken != "" {
			log.Warn("Falha na conexão, tentando gerar um novo WALLET_TOKEN")
			tokenGenerator := web3.NewWalletTokenGenerator(cfg, log)
			token, genErr := tokenGenerator.GenerateToken()
			if genErr == nil {
				client.walletToken = token
				cfg.Web3.WalletToken = token
				// Tenta validar novamente
				log.Info("Novo token gerado, tentando validar novamente")
				err = client.ValidateConnection()
			}
		}

		// Se ainda houver erro, mas temos a API key, podemos continuar
		if err != nil && cfg.LLM.APIKey != "" {
			log.Warn("Falha ao validar com Nation.fun API, mas temos API key. Continuando com funcionalidade limitada",
				"error", err.Error())
			// Continuamos mesmo com erro, pois temos a API key
		} else if err != nil {
			// Se não temos nem API key nem conexão validada, retornamos erro
			return nil, fmt.Errorf("falha ao conectar com Nation.fun API: %w", err)
		}
	}

	log.Info("Cliente Nation.fun inicializado com sucesso",
		"model", modelName,
		"nft_contract", cfg.Web3.NFTAccessContractAddress)

	return client, nil
}

// Generate implementa a interface LLMProvider.Generate para Nation.fun
func (c *NationClient) Generate(req *models.LLMRequest) (*models.LLMResponse, error) {
	startTime := time.Now()
	c.logger.Debug("Gerando resposta com Nation.fun",
		"model", c.modelName,
		"prompt_length", len(req.Prompt),
		"max_tokens", req.MaxTokens,
		"temperature", req.Temperature)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Prepara request para Nation.fun
	url := fmt.Sprintf("%s/completions", c.baseURL)

	// Estrutura da requisição
	payload := map[string]interface{}{
		"model":         c.modelName,
		"prompt":        req.Prompt,
		"temperature":   req.Temperature,
		"max_tokens":    req.MaxTokens,
		"system":        req.SystemPrompt,
		"nft_address":   c.nftAddress,
		"wallet_token":  c.walletToken,
		"agent_address": c.config.Web3.DefaultAgentAddress,
	}

	// Adiciona mensagens de contexto se existirem
	if len(req.ContextMessages) > 0 {
		messages := []map[string]string{}

		for _, msg := range req.ContextMessages {
			messages = append(messages, map[string]string{
				"role":    msg.Role,
				"content": msg.Content,
			})
		}

		payload["context"] = messages
	}

	// Configura formato de resposta JSON se solicitado
	if req.ResponseFormat == "json" {
		payload["response_format"] = "json"
	}

	// Serializa payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	// Cria requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Adiciona headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("X-Nation-NFT-Address", c.nftAddress)
	httpReq.Header.Set("X-Nation-Wallet-Token", c.walletToken)

	// Faz a chamada para a API
	c.logger.Debug("Enviando requisição para Nation.fun")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar Nation.fun API: %w", err)
	}
	defer resp.Body.Close()

	// Lê resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verifica status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Nation.fun API retornou status %d: %s", resp.StatusCode, string(body))
	}

	// Parse da resposta
	var nationResp struct {
		Content    string `json:"content"`
		Model      string `json:"model"`
		TokensUsed int    `json:"tokens_used"`
		NFTUsed    string `json:"nft_used"`
	}

	if err := json.Unmarshal(body, &nationResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	// Calcula métricas
	duration := time.Since(startTime)

	// Constrói resposta
	result := &models.LLMResponse{
		Content:    nationResp.Content,
		Model:      nationResp.Model,
		TokensUsed: nationResp.TokensUsed,
		LatencyMs:  duration.Milliseconds(),
		Provider:   "nation.fun",
		Metadata: map[string]interface{}{
			"nft_used": nationResp.NFTUsed,
		},
	}

	c.logger.Debug("Resposta Nation.fun recebida",
		"duration_ms", duration.Milliseconds(),
		"tokens_used", nationResp.TokensUsed,
		"nft_used", nationResp.NFTUsed)

	return result, nil
}

// GenerateStructured gera uma resposta estruturada usando o LLM
func (c *NationClient) GenerateStructured(req *models.LLMRequest, responseStruct interface{}) error {
	// Força formato JSON
	reqCopy := *req
	reqCopy.ResponseFormat = "json"

	// Adiciona instrução para formato JSON
	if reqCopy.SystemPrompt == "" {
		reqCopy.SystemPrompt = "Você é um assistente que responde apenas em formato JSON válido."
	} else {
		reqCopy.SystemPrompt += "\n\nIMPORTANTE: Responda apenas em formato JSON válido."
	}

	// Adiciona estrutura esperada ao prompt
	structType := fmt.Sprintf("%T", responseStruct)
	structExample, _ := json.MarshalIndent(responseStruct, "", "  ")

	reqCopy.Prompt += fmt.Sprintf("\n\nResponda com um JSON válido no formato: %s\n\nExemplo de estrutura:\n```json\n%s\n```",
		structType, string(structExample))

	// Gera resposta
	resp, err := c.Generate(&reqCopy)
	if err != nil {
		return fmt.Errorf("erro ao gerar resposta estruturada: %w", err)
	}

	// Tenta fazer parse do JSON
	err = json.Unmarshal([]byte(resp.Content), responseStruct)
	if err != nil {
		c.logger.Error("Falha ao fazer parse da resposta JSON",
			"error", err,
			"response", resp.Content)
		return fmt.Errorf("falha ao fazer parse da resposta JSON: %w", err)
	}

	return nil
}

// GetCompletion obtém uma resposta para o prompt fornecido
func (c *NationClient) GetCompletion(ctx context.Context, prompt string) (string, error) {
	req := &models.LLMRequest{
		Prompt:       prompt,
		MaxTokens:    4000, // Valor padrão
		Temperature:  0.7,  // Valor padrão
		SystemPrompt: "Você é um especialista em análise de infraestrutura e segurança em nuvem.",
	}

	resp, err := c.Generate(req)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// ValidateConnection testa a conexão com Nation.fun
func (c *NationClient) ValidateConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// URL para validação
	url := fmt.Sprintf("%s/validate", c.baseURL)

	// Cria requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Adiciona headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("X-Nation-NFT-Address", c.nftAddress)
	req.Header.Set("X-Nation-Wallet-Token", c.walletToken)

	// Faz a chamada para a API
	c.logger.Debug("Validando conexão com Nation.fun")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao conectar com Nation.fun API: %w", err)
	}
	defer resp.Body.Close()

	// Verifica status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Nation.fun API retornou status %d: %s", resp.StatusCode, string(body))
	}

	c.logger.Info("Conexão com Nation.fun validada com sucesso")
	return nil
}
