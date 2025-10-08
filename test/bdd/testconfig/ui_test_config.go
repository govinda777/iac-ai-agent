package testconfig

import (
	"fmt"
	"os"
	"testing"

	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// UITestConfig configuração específica para testes de UI
type UITestConfig struct {
	BaseURL         string
	BrowserType     string
	HeadlessMode    bool
	ScreenshotPath  string
	TestTimeout     int
	RetryAttempts   int
	MockMode        bool
	Web3MockEnabled bool
	APIMockEnabled  bool
}

// LoadUITestConfig carrega configuração para testes de UI
func LoadUITestConfig() *UITestConfig {
	return &UITestConfig{
		BaseURL:         getEnv("UI_TEST_BASE_URL", "http://localhost:8080"),
		BrowserType:     getEnv("UI_TEST_BROWSER", "chrome"),
		HeadlessMode:    getEnvBool("UI_TEST_HEADLESS", true),
		ScreenshotPath:  getEnv("UI_TEST_SCREENSHOTS", "./test/screenshots"),
		TestTimeout:     getEnvInt("UI_TEST_TIMEOUT", 30),
		RetryAttempts:   getEnvInt("UI_TEST_RETRY", 3),
		MockMode:        getEnvBool("UI_TEST_MOCK_MODE", true),
		Web3MockEnabled: getEnvBool("UI_TEST_WEB3_MOCK", true),
		APIMockEnabled:  getEnvBool("UI_TEST_API_MOCK", true),
	}
}

// UIScenarioContext contexto para cenários de teste de UI
type UIScenarioContext struct {
	Config       *config.Config
	Logger       *logger.Logger
	UITestConfig *UITestConfig
	CurrentPage  string
	UserState    map[string]interface{}
	TestData     map[string]interface{}
	Screenshots  []string
	Errors       []error
}

// NewUIScenarioContext cria novo contexto para cenários de UI
func NewUIScenarioContext(cfg *config.Config, log *logger.Logger) *UIScenarioContext {
	return &UIScenarioContext{
		Config:       cfg,
		Logger:       log,
		UITestConfig: LoadUITestConfig(),
		UserState:    make(map[string]interface{}),
		TestData:     make(map[string]interface{}),
		Screenshots:  []string{},
		Errors:       []error{},
	}
}

// SetupUITestEnvironment configura ambiente para testes de UI
func SetupUITestEnvironment(t *testing.T) (*UIScenarioContext, error) {
	// Carregar configuração (simplificado para testes)
	cfg := &config.Config{}

	// Configurar logger (simplificado para testes)
	log := &logger.Logger{}

	// Criar contexto
	ctx := NewUIScenarioContext(cfg, log)

	// Configurar ambiente de teste
	if err := ctx.setupTestEnvironment(); err != nil {
		return nil, err
	}

	return ctx, nil
}

// setupTestEnvironment configura o ambiente de teste
func (ctx *UIScenarioContext) setupTestEnvironment() error {
	ctx.Logger.Info("Configurando ambiente de teste de UI")

	// Criar diretório de screenshots se não existir
	if err := os.MkdirAll(ctx.UITestConfig.ScreenshotPath, 0755); err != nil {
		return err
	}

	// Configurar estado inicial do usuário
	ctx.UserState["authenticated"] = false
	ctx.UserState["wallet_address"] = ""
	ctx.UserState["nft_tier"] = ""
	ctx.UserState["token_balance"] = 0
	ctx.UserState["current_step"] = 1

	// Configurar dados de teste
	ctx.TestData["mock_wallet"] = "0x742d35Cc6634C0532925a3b844Bc9e7595f0bDce"
	ctx.TestData["mock_nft_tier"] = "pro"
	ctx.TestData["mock_token_balance"] = 500
	ctx.TestData["sample_terraform_code"] = `
resource "aws_s3_bucket" "example" {
  bucket = "my-terraform-bucket"
  acl    = "public-read"
}

resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
}`

	ctx.Logger.Info("Ambiente de teste configurado com sucesso")
	return nil
}

// CleanupUITestEnvironment limpa ambiente após testes
func (ctx *UIScenarioContext) CleanupUITestEnvironment() error {
	ctx.Logger.Info("Limpando ambiente de teste de UI")

	// Limpar estado do usuário
	ctx.UserState = make(map[string]interface{})
	ctx.TestData = make(map[string]interface{})
	ctx.Screenshots = []string{}
	ctx.Errors = []error{}

	ctx.Logger.Info("Ambiente de teste limpo")
	return nil
}

// TakeScreenshot captura screenshot durante teste
func (ctx *UIScenarioContext) TakeScreenshot(name string) error {
	// Implementação simplificada para testes
	screenshotPath := ctx.UITestConfig.ScreenshotPath + "/" + name + ".png"
	ctx.Screenshots = append(ctx.Screenshots, screenshotPath)
	ctx.Logger.Info(fmt.Sprintf("Screenshot capturado: %s", screenshotPath))
	return nil
}

// ValidateUIState valida estado atual da UI
func (ctx *UIScenarioContext) ValidateUIState(expectedState map[string]interface{}) error {
	ctx.Logger.Info("Validando estado da UI")

	for key, expectedValue := range expectedState {
		if actualValue, exists := ctx.UserState[key]; !exists {
			return fmt.Errorf("estado '%s' não encontrado", key)
		} else if actualValue != expectedValue {
			return fmt.Errorf("estado '%s' esperado '%v', encontrado '%v'", key, expectedValue, actualValue)
		}
	}

	ctx.Logger.Info("Estado da UI validado com sucesso")
	return nil
}

// SimulateUserAction simula ação do usuário
func (ctx *UIScenarioContext) SimulateUserAction(action string, params map[string]interface{}) error {
	ctx.Logger.Info(fmt.Sprintf("Simulando ação do usuário: %s", action))

	switch action {
	case "login":
		ctx.UserState["authenticated"] = true
		ctx.UserState["wallet_address"] = params["wallet_address"].(string)
		ctx.UserState["current_step"] = 2

	case "purchase_nft":
		ctx.UserState["nft_tier"] = params["tier"].(string)
		ctx.UserState["current_step"] = 3

	case "purchase_tokens":
		amount := params["amount"].(int)
		currentBalance := ctx.UserState["token_balance"].(int)
		ctx.UserState["token_balance"] = currentBalance + amount
		ctx.UserState["current_step"] = 4

	case "submit_analysis":
		cost := params["cost"].(int)
		currentBalance := ctx.UserState["token_balance"].(int)
		if currentBalance < cost {
			return fmt.Errorf("saldo insuficiente: %d < %d", currentBalance, cost)
		}
		ctx.UserState["token_balance"] = currentBalance - cost
		ctx.UserState["current_step"] = 5

	default:
		return fmt.Errorf("ação não reconhecida: %s", action)
	}

	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true"
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simplified - in real implementation, use strconv.Atoi
		return defaultValue
	}
	return defaultValue
}
