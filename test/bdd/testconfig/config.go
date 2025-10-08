package testconfig

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// TestMode define o modo de teste
type TestMode string

const (
	MockMode        TestMode = "mock"
	IntegrationMode TestMode = "integration"
	RealMode        TestMode = "real"
	HybridMode      TestMode = "hybrid"
)

// TestConfig contém configurações específicas para testes
type TestConfig struct {
	Mode            TestMode
	UseRealServices bool
	MockFailures    bool
	PerformanceTest bool
	LoadTestCount   int
	TimeoutSeconds  int

	// Configurações específicas para integração
	RealPrivyAppID   string
	RealBaseRPCURL   string
	RealContractAddr string
	RealNationAPIKey string

	// Configurações específicas para mocks
	MockUserID        string
	MockWalletAddress string
	MockTokenID       string
	MockSessionID     string
}

// LoadTestConfig carrega configuração de teste baseada em variáveis de ambiente
func LoadTestConfig() (*TestConfig, error) {
	cfg := &TestConfig{
		Mode:            getTestMode(),
		UseRealServices: getBoolEnv("TEST_USE_REAL_SERVICES", false),
		MockFailures:    getBoolEnv("TEST_MOCK_FAILURES", false),
		PerformanceTest: getBoolEnv("TEST_PERFORMANCE", false),
		LoadTestCount:   getIntEnv("TEST_LOAD_COUNT", 10),
		TimeoutSeconds:  getIntEnv("TEST_TIMEOUT_SECONDS", 30),

		// Configurações reais
		RealPrivyAppID:   getStringEnv("TEST_REAL_PRIVY_APP_ID", "cmgh6un8w007bl10ci0tgitwp"),
		RealBaseRPCURL:   getStringEnv("TEST_REAL_BASE_RPC_URL", "https://goerli.base.org"),
		RealContractAddr: getStringEnv("TEST_REAL_CONTRACT_ADDR", "0x147e832418Cc06A501047019E956714271098b89"),
		RealNationAPIKey: getStringEnv("TEST_REAL_NATION_API_KEY", ""),

		// Configurações mock
		MockUserID:        getStringEnv("TEST_MOCK_USER_ID", "mock_user_123"),
		MockWalletAddress: getStringEnv("TEST_MOCK_WALLET_ADDRESS", "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"),
		MockTokenID:       getStringEnv("TEST_MOCK_TOKEN_ID", "mock_token_123"),
		MockSessionID:     getStringEnv("TEST_MOCK_SESSION_ID", "mock_session_123"),
	}

	return cfg, nil
}

// GetAppConfig retorna configuração da aplicação baseada no modo de teste
func (tc *TestConfig) GetAppConfig() *config.Config {
	baseConfig := &config.Config{
		Web3: config.Web3Config{
			PrivyAppID:               tc.RealPrivyAppID,
			BaseRPCURL:               tc.RealBaseRPCURL,
			BaseChainID:              84531,
			NFTAccessContractAddress: tc.RealContractAddr,
			BotTokenContractAddress:  "0xMockTokenContract",
			EnableNFTAccess:          true,
			EnableTokenPayments:      true,
			BasicTierRateLimit:       100,
			ProTierRateLimit:         1000,
			EnterpriseTierRateLimit:  10000,
		},
		LLM: config.LLMConfig{
			Provider:    "nation.fun",
			Model:       "nation-1",
			Temperature: 0.2,
			MaxTokens:   4000,
		},
		Analysis: config.AnalysisConfig{
			CheckovEnabled:          true,
			IAMAnalysisEnabled:      true,
			CostOptimizationEnabled: true,
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}

	// Ajusta configurações baseadas no modo de teste
	if tc.Mode == MockMode {
		baseConfig.Web3.PrivyAppID = "mock_app_id"
		baseConfig.Web3.BaseRPCURL = "mock_rpc_url"
		baseConfig.Web3.NFTAccessContractAddress = "0xMockNFTContract"
	}

	return baseConfig
}

// GetLogger retorna logger configurado para testes
func (tc *TestConfig) GetLogger() *logger.Logger {
	logLevel := "info"
	if tc.PerformanceTest {
		logLevel = "warn" // Reduz logs durante testes de performance
	}

	return logger.New(logLevel, "json")
}

// IsMockMode retorna true se estiver em modo mock
func (tc *TestConfig) IsMockMode() bool {
	return tc.Mode == MockMode
}

// IsIntegrationMode retorna true se estiver em modo integração
func (tc *TestConfig) IsIntegrationMode() bool {
	return tc.Mode == IntegrationMode || tc.Mode == RealMode
}

// ShouldUseRealServices retorna true se deve usar serviços reais
func (tc *TestConfig) ShouldUseRealServices() bool {
	return tc.UseRealServices || tc.IsIntegrationMode()
}

// GetTestTags retorna tags para filtrar testes baseado no modo
func (tc *TestConfig) GetTestTags() []string {
	switch tc.Mode {
	case MockMode:
		return []string{"@mock", "@unit"}
	case IntegrationMode:
		return []string{"@integration", "@real"}
	case RealMode:
		return []string{"@integration", "@real", "@e2e"}
	case HybridMode:
		return []string{"@integration", "@mock", "@real"}
	default:
		return []string{"@mock"}
	}
}

// GetExcludedTags retorna tags a serem excluídas baseado no modo
func (tc *TestConfig) GetExcludedTags() []string {
	switch tc.Mode {
	case MockMode:
		return []string{"@integration", "@real", "@e2e"}
	case IntegrationMode:
		return []string{"@mock", "@unit"}
	case RealMode:
		return []string{"@mock", "@unit"}
	case HybridMode:
		return []string{"@unit"} // Híbrido permite todos exceto unitários
	default:
		return []string{"@integration", "@real", "@e2e"}
	}
}

// ValidateConfig valida a configuração de teste
func (tc *TestConfig) ValidateConfig() error {
	if tc.Mode == IntegrationMode || tc.Mode == RealMode {
		if tc.RealPrivyAppID == "" {
			return fmt.Errorf("RealPrivyAppID é obrigatório para testes de integração")
		}
		if tc.RealBaseRPCURL == "" {
			return fmt.Errorf("RealBaseRPCURL é obrigatório para testes de integração")
		}
		if tc.RealContractAddr == "" {
			return fmt.Errorf("RealContractAddr é obrigatório para testes de integração")
		}
	}

	if tc.PerformanceTest && tc.LoadTestCount <= 0 {
		return fmt.Errorf("LoadTestCount deve ser maior que 0 para testes de performance")
	}

	if tc.TimeoutSeconds <= 0 {
		return fmt.Errorf("TimeoutSeconds deve ser maior que 0")
	}

	return nil
}

// PrintConfig imprime a configuração atual (para debug)
func (tc *TestConfig) PrintConfig() {
	fmt.Printf("=== Configuração de Teste ===\n")
	fmt.Printf("Modo: %s\n", tc.Mode)
	fmt.Printf("Usar Serviços Reais: %t\n", tc.UseRealServices)
	fmt.Printf("Simular Falhas: %t\n", tc.MockFailures)
	fmt.Printf("Teste de Performance: %t\n", tc.PerformanceTest)
	fmt.Printf("Contagem de Carga: %d\n", tc.LoadTestCount)
	fmt.Printf("Timeout (segundos): %d\n", tc.TimeoutSeconds)
	fmt.Printf("Privy App ID: %s\n", tc.RealPrivyAppID)
	fmt.Printf("Base RPC URL: %s\n", tc.RealBaseRPCURL)
	fmt.Printf("Contrato NFT: %s\n", tc.RealContractAddr)
	fmt.Printf("=============================\n")
}

// Funções auxiliares para ler variáveis de ambiente

func getTestMode() TestMode {
	mode := strings.ToLower(getStringEnv("TEST_MODE", "mock"))
	switch mode {
	case "integration":
		return IntegrationMode
	case "real":
		return RealMode
	case "hybrid":
		return HybridMode
	case "mock":
		return MockMode
	default:
		return MockMode
	}
}

func getStringEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// TestEnvironmentSetup configura o ambiente de teste
type TestEnvironmentSetup struct {
	TestConfig *TestConfig
	AppConfig  *config.Config
	Logger     *logger.Logger
}

// NewTestEnvironmentSetup cria uma nova configuração de ambiente de teste
func NewTestEnvironmentSetup() (*TestEnvironmentSetup, error) {
	testConfig, err := LoadTestConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração de teste: %w", err)
	}

	if err := testConfig.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("configuração de teste inválida: %w", err)
	}

	appConfig := testConfig.GetAppConfig()
	logger := testConfig.GetLogger()

	return &TestEnvironmentSetup{
		TestConfig: testConfig,
		AppConfig:  appConfig,
		Logger:     logger,
	}, nil
}

// SetupEnvironment configura o ambiente baseado no modo de teste
func (tes *TestEnvironmentSetup) SetupEnvironment() error {
	tes.Logger.Info("Configurando ambiente de teste",
		"mode", tes.TestConfig.Mode,
		"use_real_services", tes.TestConfig.UseRealServices)

	if tes.TestConfig.IsMockMode() {
		tes.Logger.Info("Ambiente configurado para modo mock")
	} else if tes.TestConfig.IsIntegrationMode() {
		tes.Logger.Info("Ambiente configurado para modo de integração")
	}

	return nil
}

// GetTestSummary retorna um resumo da configuração de teste
func (tes *TestEnvironmentSetup) GetTestSummary() map[string]interface{} {
	return map[string]interface{}{
		"mode":              tes.TestConfig.Mode,
		"use_real_services": tes.TestConfig.UseRealServices,
		"mock_failures":     tes.TestConfig.MockFailures,
		"performance_test":  tes.TestConfig.PerformanceTest,
		"load_test_count":   tes.TestConfig.LoadTestCount,
		"timeout_seconds":   tes.TestConfig.TimeoutSeconds,
		"privy_app_id":      tes.TestConfig.RealPrivyAppID,
		"base_rpc_url":      tes.TestConfig.RealBaseRPCURL,
		"contract_address":  tes.TestConfig.RealContractAddr,
		"test_tags":         tes.TestConfig.GetTestTags(),
		"excluded_tags":     tes.TestConfig.GetExcludedTags(),
	}
}
