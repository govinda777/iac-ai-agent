package mocks

import (
	"context"
	"math/big"

	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
)

// MockPrivyClient simula o cliente Privy para testes
type MockPrivyClient struct {
	// Configurações de comportamento do mock
	ShouldFailAuth    bool
	ShouldFailWallet  bool
	MockUser          *web3.PrivyUser
	MockWalletAddress string
	MockUserID        string
	MockAccessToken   string
}

// NewMockPrivyClient cria um novo mock do PrivyClient
func NewMockPrivyClient() *MockPrivyClient {
	return &MockPrivyClient{
		MockUser: &web3.PrivyUser{
			ID:    "mock_user_123",
			Email: "test@example.com",
			Wallets: []web3.PrivyWallet{
				{
					Address: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
					Type:    "ethereum",
					ChainID: "84531", // Base Goerli
				},
			},
		},
		MockWalletAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		MockUserID:        "mock_user_123",
		MockAccessToken:   "mock_access_token_123",
	}
}

// Authenticate simula autenticação
func (m *MockPrivyClient) Authenticate(ctx context.Context, accessToken string) (*web3.PrivyUser, error) {
	if m.ShouldFailAuth {
		return nil, &web3.PrivyError{
			Code:    "AUTHENTICATION_FAILED",
			Message: "Mock authentication failed",
		}
	}
	return m.MockUser, nil
}

// GetUser simula busca de usuário
func (m *MockPrivyClient) GetUser(userID string) (*web3.PrivyUser, error) {
	if m.ShouldFailAuth {
		return nil, &web3.PrivyError{
			Code:    "USER_NOT_FOUND",
			Message: "Mock user not found",
		}
	}
	return m.MockUser, nil
}

// ValidateWalletOwnership simula validação de propriedade da wallet
func (m *MockPrivyClient) ValidateWalletOwnership(userID, walletAddress string) (bool, error) {
	if m.ShouldFailWallet {
		return false, &web3.PrivyError{
			Code:    "WALLET_VALIDATION_FAILED",
			Message: "Mock wallet validation failed",
		}
	}
	return walletAddress == m.MockWalletAddress, nil
}

// MockNFTAccessManager simula o gerenciador de NFT
type MockNFTAccessManager struct {
	// Configurações de comportamento
	ShouldFailMint      bool
	ShouldFailTransfer  bool
	MockTokenID         string
	MockTransactionHash string
	MockBalance         int
}

// NewMockNFTAccessManager cria um novo mock do NFTAccessManager
func NewMockNFTAccessManager() *MockNFTAccessManager {
	return &MockNFTAccessManager{
		MockTokenID:         "mock_token_123",
		MockTransactionHash: "0xmock_tx_hash_123",
		MockBalance:         1,
	}
}

// MintNFTAccess simula mint de NFT
func (m *MockNFTAccessManager) MintNFTAccess(ctx context.Context, userID, walletAddress, tier string) (*web3.NFTMintResult, error) {
	if m.ShouldFailMint {
		return nil, &web3.Web3Error{
			Code:    "MINT_FAILED",
			Message: "Mock mint failed",
		}
	}

	return &web3.NFTMintResult{
		TokenID:         m.MockTokenID,
		TransactionHash: m.MockTransactionHash,
		Status:          "success",
		Tier:            tier,
	}, nil
}

// CheckNFTAccess simula verificação de acesso NFT
func (m *MockNFTAccessManager) CheckNFTAccess(ctx context.Context, walletAddress string) (*web3.NFTAccessStatus, error) {
	return &web3.NFTAccessStatus{
		HasAccess:   m.MockBalance > 0,
		TokenID:     m.MockTokenID,
		Tier:        "pro",
		Balance:     m.MockBalance,
		LastUpdated: 1700000000,
	}, nil
}

// MockBotTokenManager simula o gerenciador de tokens
type MockBotTokenManager struct {
	// Configurações de comportamento
	ShouldFailTransfer  bool
	ShouldFailMint      bool
	MockBalance         *big.Int
	MockTransactionHash string
}

// NewMockBotTokenManager cria um novo mock do BotTokenManager
func NewMockBotTokenManager() *MockBotTokenManager {
	return &MockBotTokenManager{
		MockBalance:         big.NewInt(1000), // 1000 tokens
		MockTransactionHash: "0xmock_token_tx_123",
	}
}

// MintTokens simula mint de tokens
func (m *MockBotTokenManager) MintTokens(ctx context.Context, userID, walletAddress string, amount *big.Int) (*web3.TokenMintResult, error) {
	if m.ShouldFailMint {
		return nil, &web3.Web3Error{
			Code:    "MINT_FAILED",
			Message: "Mock token mint failed",
		}
	}

	// Adiciona tokens ao saldo mockado
	m.MockBalance.Add(m.MockBalance, amount)

	return &web3.TokenMintResult{
		Amount:          amount.String(),
		TransactionHash: m.MockTransactionHash,
		Status:          "success",
		NewBalance:      m.MockBalance.String(),
	}, nil
}

// GetTokenBalance simula busca de saldo
func (m *MockBotTokenManager) GetTokenBalance(ctx context.Context, walletAddress string) (*big.Int, error) {
	return m.MockBalance, nil
}

// TransferTokens simula transferência de tokens
func (m *MockBotTokenManager) TransferTokens(ctx context.Context, fromWallet, toWallet string, amount *big.Int) (*web3.TokenTransferResult, error) {
	if m.ShouldFailTransfer {
		return nil, &web3.Web3Error{
			Code:    "TRANSFER_FAILED",
			Message: "Mock transfer failed",
		}
	}

	// Subtrai tokens do saldo mockado
	m.MockBalance.Sub(m.MockBalance, amount)

	return &web3.TokenTransferResult{
		Amount:          amount.String(),
		TransactionHash: m.MockTransactionHash,
		Status:          "success",
		FromBalance:     m.MockBalance.String(),
	}, nil
}

// MockPrivyOnrampManager simula o gerenciador de onramp
type MockPrivyOnrampManager struct {
	// Configurações de comportamento
	ShouldFailSession bool
	ShouldFailPayment bool
	MockSessionID     string
	MockTransactionID string
	MockQuote         *web3.OnrampQuote
}

// NewMockPrivyOnrampManager cria um novo mock do PrivyOnrampManager
func NewMockPrivyOnrampManager() *MockPrivyOnrampManager {
	return &MockPrivyOnrampManager{
		MockSessionID:     "mock_session_123",
		MockTransactionID: "mock_tx_123",
		MockQuote: &web3.OnrampQuote{
			QuoteID:        "mock_quote_123",
			SourceAmount:   "125.00",
			SourceCurrency: "USD",
			TargetAmount:   "0.05",
			TargetCurrency: "ETH",
			NetworkFee:     "2.50",
			ProviderFee:    "5.00",
			TotalCost:      "132.50",
			ExchangeRate:   "2500.00",
			EstimatedTime:  "5-10 minutes",
			Provider:       "moonpay",
			ExpiresAt:      1700001800,
		},
	}
}

// CreateOnrampSession simula criação de sessão de onramp
func (m *MockPrivyOnrampManager) CreateOnrampSession(ctx context.Context, req *web3.CreateOnrampSessionRequest) (*web3.OnrampSession, error) {
	if m.ShouldFailSession {
		return nil, &web3.Web3Error{
			Code:    "SESSION_CREATION_FAILED",
			Message: "Mock session creation failed",
		}
	}

	return &web3.OnrampSession{
		SessionID:      m.MockSessionID,
		UserID:         req.UserID,
		WalletAddress:  req.WalletAddress,
		Purpose:        req.Purpose,
		RequiredAmount: big.NewInt(50000000000000000), // 0.05 ETH
		Currency:       "ETH",
		Status:         "created",
		Quote:          m.MockQuote,
		Metadata:       req.Metadata,
		CreatedAt:      1700000000,
		ExpiresAt:      1700003600,
	}, nil
}

// InitiatePayment simula início de pagamento
func (m *MockPrivyOnrampManager) InitiatePayment(ctx context.Context, sessionID, paymentMethod string) (*web3.OnrampTransaction, error) {
	if m.ShouldFailPayment {
		return nil, &web3.Web3Error{
			Code:    "PAYMENT_INITIATION_FAILED",
			Message: "Mock payment initiation failed",
		}
	}

	return &web3.OnrampTransaction{
		TransactionID:  m.MockTransactionID,
		QuoteID:        m.MockQuote.QuoteID,
		UserID:         "mock_user_123",
		WalletAddress:  "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		Status:         "pending",
		SourceAmount:   m.MockQuote.SourceAmount,
		SourceCurrency: m.MockQuote.SourceCurrency,
		TargetAmount:   m.MockQuote.TargetAmount,
		TargetCurrency: m.MockQuote.TargetCurrency,
		PaymentMethod:  paymentMethod,
		CreatedAt:      1700000000,
	}, nil
}

// GetOnrampStatus simula busca de status
func (m *MockPrivyOnrampManager) GetOnrampStatus(ctx context.Context, transactionID string) (*web3.OnrampTransaction, error) {
	return &web3.OnrampTransaction{
		TransactionID:  transactionID,
		Status:         "completed",
		SourceAmount:   m.MockQuote.SourceAmount,
		SourceCurrency: m.MockQuote.SourceCurrency,
		TargetAmount:   m.MockQuote.TargetAmount,
		TargetCurrency: m.MockQuote.TargetCurrency,
		TxHash:         "0xmock_completed_tx",
		CompletedAt:    1700001000,
	}, nil
}

// MockAnalysisService simula o serviço de análise
type MockAnalysisService struct {
	// Configurações de comportamento
	ShouldFailAnalysis bool
	MockAnalysisResult map[string]interface{}
	MockProcessingTime int // em segundos
}

// NewMockAnalysisService cria um novo mock do AnalysisService
func NewMockAnalysisService() *MockAnalysisService {
	return &MockAnalysisService{
		MockProcessingTime: 2,
		MockAnalysisResult: map[string]interface{}{
			"score": 85,
			"issues": []map[string]interface{}{
				{
					"severity":   "HIGH",
					"message":    "S3 bucket has public read ACL",
					"resource":   "aws_s3_bucket.example",
					"suggestion": "Remove public-read ACL for security",
					"impact":     "High security risk",
					"code_fix":   "acl = \"private\"",
				},
				{
					"severity":   "MEDIUM",
					"message":    "Instance type could be optimized",
					"resource":   "aws_instance.web",
					"suggestion": "Consider using t3.small for better performance",
					"impact":     "Cost optimization opportunity",
					"code_fix":   "instance_type = \"t3.small\"",
				},
			},
			"sections": []string{
				"Executive Summary",
				"Security Issues",
				"Best Practices",
				"Cost Optimization",
				"Detailed Findings",
			},
			"summary": "Found 2 issues requiring attention",
			"recommendations": []string{
				"Fix S3 bucket ACL security issue",
				"Consider upgrading instance type",
			},
		},
	}
}

// AnalyzeCode simula análise de código
func (m *MockAnalysisService) AnalyzeCode(ctx context.Context, code string, analysisType string) (map[string]interface{}, error) {
	if m.ShouldFailAnalysis {
		return nil, &web3.Web3Error{
			Code:    "ANALYSIS_FAILED",
			Message: "Mock analysis failed",
		}
	}

	// Simula tempo de processamento
	// time.Sleep(time.Duration(m.MockProcessingTime) * time.Second)

	return m.MockAnalysisResult, nil
}

// MockConfigProvider fornece configurações mockadas para testes
type MockConfigProvider struct {
	MockConfig map[string]interface{}
}

// NewMockConfigProvider cria um novo provedor de configuração mockada
func NewMockConfigProvider() *MockConfigProvider {
	return &MockConfigProvider{
		MockConfig: map[string]interface{}{
			"web3": map[string]interface{}{
				"privy_app_id":                "cmgh6un8w007bl10ci0tgitwp",
				"base_rpc_url":                "https://goerli.base.org",
				"base_chain_id":               84531,
				"nft_access_contract_address": "0x147e832418Cc06A501047019E956714271098b89",
				"bot_token_contract_address":  "0xMockTokenContract",
				"enable_nft_access":           true,
				"enable_token_payments":       true,
				"basic_tier_rate_limit":       100,
				"pro_tier_rate_limit":         1000,
				"enterprise_tier_rate_limit":  10000,
			},
			"llm": map[string]interface{}{
				"provider":    "nation.fun",
				"model":       "nation-1",
				"temperature": 0.2,
				"max_tokens":  4000,
			},
			"analysis": map[string]interface{}{
				"checkov_enabled":           true,
				"iam_analysis_enabled":      true,
				"cost_optimization_enabled": true,
			},
		},
	}
}

// GetConfig retorna configuração mockada
func (m *MockConfigProvider) GetConfig() map[string]interface{} {
	return m.MockConfig
}

// GetWeb3Config retorna configuração Web3 mockada
func (m *MockConfigProvider) GetWeb3Config() map[string]interface{} {
	if web3Config, ok := m.MockConfig["web3"].(map[string]interface{}); ok {
		return web3Config
	}
	return make(map[string]interface{})
}
