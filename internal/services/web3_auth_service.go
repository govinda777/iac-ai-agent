package services

import (
	"context"
	"fmt"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// Web3AuthService gerencia autenticação e autorização Web3
type Web3AuthService struct {
	config      *config.Config
	logger      *logger.Logger
	privyClient *web3.PrivyClient
	nftAccess   *web3.NFTAccessManager
	botToken    *web3.BotTokenManager
}

// NewWeb3AuthService cria um novo serviço de autenticação Web3
func NewWeb3AuthService(cfg *config.Config, log *logger.Logger,
	privyClient *web3.PrivyClient, nftAccess *web3.NFTAccessManager,
	botToken *web3.BotTokenManager) *Web3AuthService {
	return &Web3AuthService{
		config:      cfg,
		logger:      log,
		privyClient: privyClient,
		nftAccess:   nftAccess,
		botToken:    botToken,
	}
}

// AuthenticatedUser representa um usuário autenticado via Web3
type AuthenticatedUser struct {
	UserID            string   `json:"user_id"`
	WalletAddress     string   `json:"wallet_address"`
	Email             string   `json:"email,omitempty"`
	HasNFTAccess      bool     `json:"has_nft_access"`
	NFTTier           int      `json:"nft_tier,omitempty"`
	TokenBalance      string   `json:"token_balance,omitempty"`
	AllowedOperations []string `json:"allowed_operations,omitempty"`
}

// VerifyToken verifica um token de autenticação e retorna informações do usuário
func (s *Web3AuthService) VerifyToken(ctx context.Context, token string) (*AuthenticatedUser, error) {
	s.logger.Info("Verificando token de autenticação Web3")

	// Verificar token com Privy
	privyUser, err := s.privyClient.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("erro na verificação do token Privy: %w", err)
	}

	// Criar resposta básica
	user := &AuthenticatedUser{
		UserID:        privyUser.ID,
		WalletAddress: privyUser.WalletAddress,
		Email:         privyUser.Email,
	}

	// Se temos wallet address, verificar NFT e tokens
	if user.WalletAddress != "" {
		// Verificar NFT de acesso
		nftAccess, err := s.nftAccess.CheckAccess(ctx, user.WalletAddress)
		if err == nil && nftAccess != nil {
			user.HasNFTAccess = true
			user.NFTTier = int(nftAccess.Tier.TierID)
		} else {
			s.logger.Info("Usuário sem NFT de acesso", "wallet", user.WalletAddress)
		}

		// Verificar saldo de tokens
		if s.config.Web3.EnableTokenPayments {
			tokenBalance, err := s.botToken.GetBalance(ctx, user.WalletAddress)
			if err == nil && tokenBalance != nil {
				user.TokenBalance = tokenBalance.BalanceFormatted
			}
		}

		// Determinar operações permitidas
		user.AllowedOperations = s.determineAllowedOperations(user.NFTTier)
	}

	s.logger.Info("Token verificado com sucesso",
		"user_id", user.UserID,
		"wallet", user.WalletAddress,
		"has_nft", user.HasNFTAccess,
		"tier", user.NFTTier)

	return user, nil
}

// determineAllowedOperations determina quais operações um usuário pode realizar com base no tier
func (s *Web3AuthService) determineAllowedOperations(tier int) []string {
	// Operações básicas que todos podem fazer
	operations := []string{"view_docs", "basic_analysis"}

	// Adicionar operações com base no tier
	switch tier {
	case 1: // Basic Access
		operations = append(operations,
			"terraform_analysis",
			"checkov_scan")
	case 2: // Pro Access
		operations = append(operations,
			"terraform_analysis",
			"checkov_scan",
			"llm_analysis",
			"preview_analysis",
			"security_audit")
	case 3: // Enterprise Access
		operations = append(operations,
			"terraform_analysis",
			"checkov_scan",
			"llm_analysis",
			"preview_analysis",
			"security_audit",
			"cost_optimization",
			"priority_support",
			"full_review")
	}

	return operations
}

// IsOperationAllowed verifica se um usuário pode realizar uma operação
func (s *Web3AuthService) IsOperationAllowed(ctx context.Context, walletAddress string, operation string) (bool, error) {
	// Verificar NFT de acesso
	nftAccess, err := s.nftAccess.CheckAccess(ctx, walletAddress)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar acesso NFT: %w", err)
	}

	// Se não tem NFT, não pode fazer nada além de operações públicas
	if nftAccess == nil {
		return operation == "view_docs" || operation == "basic_analysis", nil
	}

	// Verificar operações com base no tier
	operations := s.determineAllowedOperations(int(nftAccess.Tier.TierID))

	for _, op := range operations {
		if op == operation {
			return true, nil
		}
	}

	return false, nil
}

// GetTokenCost retorna o custo em tokens de uma operação
func (s *Web3AuthService) GetTokenCost(operation string) (string, error) {
	cost, err := s.botToken.CalculateTokenCost(operation)
	if err != nil {
		return "", err
	}

	return s.botToken.FormatTokenAmount(cost), nil
}

// CheckRateLimit verifica se o usuário atingiu o limite de taxa para seu tier
func (s *Web3AuthService) CheckRateLimit(ctx context.Context, walletAddress string) (bool, error) {
	// Para simplificar, vamos apenas verificar o tier do usuário e comparar com os limites
	nftAccess, err := s.nftAccess.CheckAccess(ctx, walletAddress)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar acesso NFT: %w", err)
	}

	// Se não tem NFT, usar limite básico
	tier := 0
	if nftAccess != nil {
		tier = int(nftAccess.Tier.TierID)
	}

	// Pegar limite de requisições com base no tier
	var limit int
	switch tier {
	case 1: // Basic
		limit = s.config.Web3.BasicTierRateLimit
	case 2: // Pro
		limit = s.config.Web3.ProTierRateLimit
	case 3: // Enterprise
		limit = s.config.Web3.EnterpriseTierRateLimit
	default: // Público
		limit = 5 // Valor baixo para usuários sem NFT
	}

	// TODO: Implementar verificação real de rate limiting usando Redis ou similar
	// Por enquanto, simulamos com um valor fixo dentro do limite
	currentUsage := 1 // Simular uso atual

	return currentUsage <= limit, nil
}

// CreateSessionToken cria um token de sessão para o usuário
func (s *Web3AuthService) CreateSessionToken(userID, walletAddress string) (string, time.Time, error) {
	// TODO: Implementar geração de token JWT ou similar
	// Por ora, retornar um placeholder
	expiry := time.Now().Add(24 * time.Hour)
	return fmt.Sprintf("web3_session_%s_%d", userID, time.Now().Unix()), expiry, nil
}
