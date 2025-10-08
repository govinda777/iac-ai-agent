package mocks

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// MockTestEnvironment configura o ambiente de teste com mocks
type MockTestEnvironment struct {
	// Configuração
	Config *config.Config
	Logger *logger.Logger

	// Mocks dos serviços Web3
	PrivyClient      *MockPrivyClient
	NFTAccessManager *MockNFTAccessManager
	BotTokenManager  *MockBotTokenManager
	OnrampManager    *MockPrivyOnrampManager

	// Mock do serviço de análise
	AnalysisService *MockAnalysisService

	// Estado do teste
	IsMockMode bool
}

// NewMockTestEnvironment cria um novo ambiente de teste com mocks
func NewMockTestEnvironment() *MockTestEnvironment {
	// Configuração mockada
	cfg := &config.Config{
		Web3: config.Web3Config{
			PrivyAppID:               "cmgh6un8w007bl10ci0tgitwp",
			BaseRPCURL:               "https://goerli.base.org",
			BaseChainID:              84531,
			NFTAccessContractAddress: "0x147e832418Cc06A501047019E956714271098b89",
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
	}

	// Logger mockado
	log := logger.New("mock_test", "info")

	return &MockTestEnvironment{
		Config:           cfg,
		Logger:           log,
		PrivyClient:      NewMockPrivyClient(),
		NFTAccessManager: NewMockNFTAccessManager(),
		BotTokenManager:  NewMockBotTokenManager(),
		OnrampManager:    NewMockPrivyOnrampManager(),
		AnalysisService:  NewMockAnalysisService(),
		IsMockMode:       true,
	}
}

// SetupMockMode configura o ambiente para usar mocks
func (mte *MockTestEnvironment) SetupMockMode() {
	mte.IsMockMode = true
	mte.Logger.Info("Configurando ambiente de teste com mocks")
}

// SetupRealMode configura o ambiente para usar serviços reais
func (mte *MockTestEnvironment) SetupRealMode() {
	mte.IsMockMode = false
	mte.Logger.Info("Configurando ambiente de teste com serviços reais")
}

// GetPrivyClient retorna o cliente Privy apropriado (mock ou real)
func (mte *MockTestEnvironment) GetPrivyClient() *MockPrivyClient {
	if mte.IsMockMode {
		return mte.PrivyClient
	}
	// Em modo real, retornaria o cliente real
	// return web3.NewPrivyClient(mte.Config.Web3.PrivyAppID, mte.Logger)
	return mte.PrivyClient // Por enquanto, sempre mock
}

// GetNFTAccessManager retorna o gerenciador NFT apropriado
func (mte *MockTestEnvironment) GetNFTAccessManager() *MockNFTAccessManager {
	if mte.IsMockMode {
		return mte.NFTAccessManager
	}
	// Em modo real, retornaria o gerenciador real
	return mte.NFTAccessManager // Por enquanto, sempre mock
}

// GetBotTokenManager retorna o gerenciador de tokens apropriado
func (mte *MockTestEnvironment) GetBotTokenManager() *MockBotTokenManager {
	if mte.IsMockMode {
		return mte.BotTokenManager
	}
	// Em modo real, retornaria o gerenciador real
	return mte.BotTokenManager // Por enquanto, sempre mock
}

// GetOnrampManager retorna o gerenciador de onramp apropriado
func (mte *MockTestEnvironment) GetOnrampManager() *MockPrivyOnrampManager {
	if mte.IsMockMode {
		return mte.OnrampManager
	}
	// Em modo real, retornaria o gerenciador real
	return mte.OnrampManager // Por enquanto, sempre mock
}

// GetAnalysisService retorna o serviço de análise apropriado
func (mte *MockTestEnvironment) GetAnalysisService() *MockAnalysisService {
	if mte.IsMockMode {
		return mte.AnalysisService
	}
	// Em modo real, retornaria o serviço real
	return mte.AnalysisService // Por enquanto, sempre mock
}

// SimulateUserLogin simula login de usuário
func (mte *MockTestEnvironment) SimulateUserLogin(walletAddress string) (*web3.PrivyUser, error) {
	mte.Logger.Info("Simulando login de usuário", "wallet", walletAddress)

	// Atualiza o mock com o endereço fornecido
	mte.PrivyClient.MockWalletAddress = walletAddress
	mte.PrivyClient.MockUser.Wallets[0].Address = walletAddress

	return mte.PrivyClient.MockUser, nil
}

// SimulateNFTPurchase simula compra de NFT
func (mte *MockTestEnvironment) SimulateNFTPurchase(userID, walletAddress, tier string) (*web3.NFTMintResult, error) {
	mte.Logger.Info("Simulando compra de NFT",
		"user_id", userID,
		"wallet", walletAddress,
		"tier", tier)

	return mte.NFTAccessManager.MintNFTAccess(context.Background(), userID, walletAddress, tier)
}

// SimulateTokenPurchase simula compra de tokens
func (mte *MockTestEnvironment) SimulateTokenPurchase(userID, walletAddress string, amount int64) (*web3.TokenMintResult, error) {
	mte.Logger.Info("Simulando compra de tokens",
		"user_id", userID,
		"wallet", walletAddress,
		"amount", amount)

	bigAmount := big.NewInt(amount)
	return mte.BotTokenManager.MintTokens(context.Background(), userID, walletAddress, bigAmount)
}

// SimulateAnalysis simula análise de código
func (mte *MockTestEnvironment) SimulateAnalysis(code, analysisType string) (map[string]interface{}, error) {
	mte.Logger.Info("Simulando análise de código",
		"analysis_type", analysisType,
		"code_length", len(code))

	return mte.AnalysisService.AnalyzeCode(context.Background(), code, analysisType)
}

// SimulateOnrampSession simula criação de sessão de onramp
func (mte *MockTestEnvironment) SimulateOnrampSession(userID, walletAddress, purpose, itemID string) (*web3.OnrampSession, error) {
	mte.Logger.Info("Simulando sessão de onramp",
		"user_id", userID,
		"wallet", walletAddress,
		"purpose", purpose,
		"item_id", itemID)

	req := &web3.CreateOnrampSessionRequest{
		UserID:         userID,
		WalletAddress:  walletAddress,
		Purpose:        purpose,
		TargetItemID:   itemID,
		SourceCurrency: "USD",
		Metadata: map[string]interface{}{
			"test_mode": true,
			"timestamp": time.Now().Unix(),
		},
	}

	return mte.OnrampManager.CreateOnrampSession(context.Background(), req)
}

// ValidateMockResults valida os resultados dos mocks
func (mte *MockTestEnvironment) ValidateMockResults() error {
	mte.Logger.Info("Validando resultados dos mocks")

	// Validações básicas
	if mte.PrivyClient.MockUser == nil {
		return fmt.Errorf("Mock Privy user não configurado")
	}

	if mte.NFTAccessManager.MockTokenID == "" {
		return fmt.Errorf("Mock NFT token ID não configurado")
	}

	if mte.BotTokenManager.MockBalance == nil {
		return fmt.Errorf("Mock token balance não configurado")
	}

	if mte.OnrampManager.MockSessionID == "" {
		return fmt.Errorf("Mock onramp session ID não configurado")
	}

	if len(mte.AnalysisService.MockAnalysisResult) == 0 {
		return fmt.Errorf("Mock analysis result não configurado")
	}

	mte.Logger.Info("Validação dos mocks concluída com sucesso")
	return nil
}

// ResetMocks reseta todos os mocks para estado inicial
func (mte *MockTestEnvironment) ResetMocks() {
	mte.Logger.Info("Resetando mocks para estado inicial")

	// Reset Privy Client
	mte.PrivyClient.ShouldFailAuth = false
	mte.PrivyClient.ShouldFailWallet = false

	// Reset NFT Access Manager
	mte.NFTAccessManager.ShouldFailMint = false
	mte.NFTAccessManager.ShouldFailTransfer = false
	mte.NFTAccessManager.MockBalance = 1

	// Reset Bot Token Manager
	mte.BotTokenManager.ShouldFailTransfer = false
	mte.BotTokenManager.ShouldFailMint = false
	mte.BotTokenManager.MockBalance = big.NewInt(1000)

	// Reset Onramp Manager
	mte.OnrampManager.ShouldFailSession = false
	mte.OnrampManager.ShouldFailPayment = false

	// Reset Analysis Service
	mte.AnalysisService.ShouldFailAnalysis = false
	mte.AnalysisService.MockProcessingTime = 2
}

// GetTestSummary retorna um resumo do estado atual dos testes
func (mte *MockTestEnvironment) GetTestSummary() map[string]interface{} {
	return map[string]interface{}{
		"mode":              mte.IsMockMode,
		"privy_user_id":     mte.PrivyClient.MockUserID,
		"wallet_address":    mte.PrivyClient.MockWalletAddress,
		"nft_token_id":      mte.NFTAccessManager.MockTokenID,
		"token_balance":     mte.BotTokenManager.MockBalance.String(),
		"onramp_session_id": mte.OnrampManager.MockSessionID,
		"analysis_score":    mte.AnalysisService.MockAnalysisResult["score"],
		"config_app_id":     mte.Config.Web3.PrivyAppID,
	}
}
