package whatsapp

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

// WhatsAppAgent representa um agente WhatsApp
type WhatsAppAgent struct {
	ID             string
	Name           string
	Description    string
	WalletAddr     string
	APIKey         string
	Service        *MockAgentService
	LLMService     *MockLLMService
	AuthService    *MockAuthService
	BillingService *MockBillingService
	Logger         *WhatsAppLogger
	Commands       map[string]*Command
}

// WhatsAppLogger sistema de logging para WhatsApp
type WhatsAppLogger struct {
	AgentID string
}

// MockAgentService mock para AgentService
type MockAgentService struct{}

func (m *MockAgentService) AnalyzeCode(code string) (*AnalysisResult, error) {
	return &AnalysisResult{
		Issues: []Issue{
			{Description: "Test issue"},
		},
		Suggestions: []string{"Test suggestion"},
	}, nil
}

func (m *MockAgentService) AnalyzeSecurity(code string) (*SecurityAnalysis, error) {
	return &SecurityAnalysis{
		Vulnerabilities: []Vulnerability{
			{Description: "Test vulnerability", Severity: "Medium"},
		},
		Recommendations: []string{"Test recommendation"},
	}, nil
}

func (m *MockAgentService) AnalyzeCosts(code string) (*CostAnalysis, error) {
	return &CostAnalysis{
		EstimatedMonthlyCost: 100.0,
		PotentialSavings:     20.0,
		Optimizations: []Optimization{
			{Description: "Test optimization", MonthlySavings: 10.0},
		},
	}, nil
}

// MockLLMService mock para LLMService
type MockLLMService struct{}

func (m *MockLLMService) GenerateResponse(prompt string) (string, error) {
	return "Mock response", nil
}

// MockAuthService mock para AuthService
type MockAuthService struct{}

func (m *MockAuthService) VerifyWalletNFT(ctx context.Context) error {
	return nil
}

func (m *MockAuthService) RecoverAPIKey(ctx context.Context) (string, error) {
	return "mock_api_key", nil
}

// MockBillingService mock para BillingService
type MockBillingService struct{}

func (m *MockBillingService) ChargeTokens(ctx context.Context, userAddr string, amount int) error {
	return nil
}

func (m *MockBillingService) GetUsageStats(ctx context.Context, userAddr string) (*UsageStats, error) {
	return &UsageStats{
		TotalRequests:  10,
		TokensConsumed: 8,
		LastRequest:    time.Now(),
		RequestsToday:  3,
		AverageCost:    0.8,
	}, nil
}

// AnalysisResult resultado da análise
type AnalysisResult struct {
	Issues      []Issue  `json:"issues"`
	Suggestions []string `json:"suggestions"`
}

// Issue problema encontrado
type Issue struct {
	Description string `json:"description"`
}

// SecurityAnalysis resultado da análise de segurança
type SecurityAnalysis struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Recommendations []string        `json:"recommendations"`
}

// Vulnerability vulnerabilidade encontrada
type Vulnerability struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// CostAnalysis resultado da análise de custos
type CostAnalysis struct {
	EstimatedMonthlyCost float64        `json:"estimated_monthly_cost"`
	PotentialSavings     float64        `json:"potential_savings"`
	Optimizations        []Optimization `json:"optimizations"`
}

// Optimization otimização sugerida
type Optimization struct {
	Description    string  `json:"description"`
	MonthlySavings float64 `json:"monthly_savings"`
}

// NewWhatsAppAgent cria um novo agente WhatsApp
func NewWhatsAppAgent(config *WhatsAppAgentConfig) (*WhatsAppAgent, error) {
	// Verificar wallet e NFT
	authService := &MockAuthService{}
	if err := authService.VerifyWalletNFT(context.Background()); err != nil {
		return nil, fmt.Errorf("wallet verification failed: %w", err)
	}

	// Recuperar chave API via Lit Protocol
	apiKey, err := authService.RecoverAPIKey(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to recover API key: %w", err)
	}

	// Criar serviços
	agentService := &MockAgentService{}
	llmService := &MockLLMService{}
	billingService := &MockBillingService{}

	// Criar logger
	logger := &WhatsAppLogger{
		AgentID: generateAgentID(),
	}

	// Criar agente
	agent := &WhatsAppAgent{
		ID:             logger.AgentID,
		Name:           config.Name,
		Description:    config.Description,
		WalletAddr:     config.WalletAddr,
		APIKey:         apiKey,
		Service:        agentService,
		LLMService:     llmService,
		AuthService:    authService,
		BillingService: billingService,
		Logger:         logger,
		Commands:       AvailableCommands(),
	}

	return agent, nil
}

// ProcessMessage processa mensagens recebidas
func (a *WhatsAppAgent) ProcessMessage(ctx context.Context, msg *WhatsAppMessage) (*WhatsAppResponse, error) {
	// Log da mensagem recebida
	a.Logger.LogMessage(msg)

	// Verificar autenticação
	if err := a.authenticate(ctx, msg.From); err != nil {
		a.Logger.LogError(err, "authentication")
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Processar comando
	command, err := a.parseCommand(msg.Text)
	if err != nil {
		return a.handleError("Comando inválido. Use /help para ver comandos disponíveis."), nil
	}

	// Executar comando
	response, err := a.executeCommand(ctx, command, msg)
	if err != nil {
		a.Logger.LogError(err, "command execution")
		return a.handleError(fmt.Sprintf("Erro ao executar comando: %v", err)), nil
	}

	// Cobrar tokens se necessário
	if command.RequiresPayment {
		if err := a.BillingService.ChargeTokens(ctx, msg.From, command.TokenCost); err != nil {
			a.Logger.LogError(err, "billing")
			log.Printf("Failed to charge tokens: %v", err)
		}
	}

	// Log da resposta
	a.Logger.LogResponse(response)

	return response, nil
}

// authenticate verifica autenticação do usuário
func (a *WhatsAppAgent) authenticate(ctx context.Context, userAddr string) error {
	// Por enquanto, aceita qualquer usuário
	// Em produção, implementar verificação de wallet/NFT
	return nil
}

// parseCommand analisa o texto da mensagem e identifica o comando
func (a *WhatsAppAgent) parseCommand(text string) (*Command, error) {
	text = strings.TrimSpace(text)

	for _, cmd := range a.Commands {
		matched, err := regexp.MatchString(cmd.Pattern, text)
		if err != nil {
			continue
		}
		if matched {
			return cmd, nil
		}
	}

	return nil, fmt.Errorf("command not found")
}

// executeCommand executa o comando identificado
func (a *WhatsAppAgent) executeCommand(ctx context.Context, cmd *Command, msg *WhatsAppMessage) (*WhatsAppResponse, error) {
	// Extrair argumentos e código
	args, codeBlock := a.extractCommandArgs(msg.Text, cmd.Pattern)

	context := &CommandContext{
		Message:   msg,
		Agent:     a,
		Args:      args,
		CodeBlock: codeBlock,
	}

	return cmd.Handler(a, context)
}

// extractCommandArgs extrai argumentos e blocos de código do comando
func (a *WhatsAppAgent) extractCommandArgs(text, pattern string) ([]string, string) {
	// Remover o comando do texto
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return []string{}, ""
	}

	remaining := strings.TrimSpace(matches[1])

	// Procurar por blocos de código
	codeBlockRegex := regexp.MustCompile("```(?s)(.*?)```")
	codeMatches := codeBlockRegex.FindStringSubmatch(remaining)

	var codeBlock string
	if len(codeMatches) > 1 {
		codeBlock = strings.TrimSpace(codeMatches[1])
		// Remover o bloco de código do texto restante
		remaining = codeBlockRegex.ReplaceAllString(remaining, "")
	}

	// Dividir argumentos
	args := strings.Fields(remaining)

	return args, codeBlock
}

// handleError cria resposta de erro padronizada
func (a *WhatsAppAgent) handleError(message string) *WhatsAppResponse {
	return &WhatsAppResponse{
		Text: ResponseTemplates["error"] + "\n\n" + message,
		Type: "text",
	}
}

// GetUsageStats retorna estatísticas de uso do usuário
func (a *WhatsAppAgent) GetUsageStats(ctx context.Context, userAddr string) (*UsageStats, error) {
	return a.BillingService.GetUsageStats(ctx, userAddr)
}

// generateAgentID gera ID único para o agente
func generateAgentID() string {
	return fmt.Sprintf("whatsapp_agent_%d", time.Now().UnixNano())
}

// encryptWithAES criptografa dados com AES-256
func encryptWithAES(plaintext string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// decryptWithAES descriptografa dados com AES-256
func decryptWithAES(ciphertext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// generateAESKey gera chave AES-256
func generateAESKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		// Fallback para hash determinístico
		hash := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
		return hash[:]
	}
	return key
}
