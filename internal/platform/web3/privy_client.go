package web3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// PrivyClient é o cliente para integração com Privy.io
type PrivyClient struct {
	config     *config.Config
	logger     *logger.Logger
	httpClient *http.Client
	appID      string
	baseURL    string
}

// NewPrivyClient cria um novo cliente Privy
func NewPrivyClient(cfg *config.Config, log *logger.Logger) *PrivyClient {
	return &PrivyClient{
		config:     cfg,
		logger:     log,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		appID:      cfg.Web3.PrivyAppID,
		baseURL:    "https://auth.privy.io",
	}
}

// PrivyUser representa um usuário autenticado via Privy
type PrivyUser struct {
	ID               string          `json:"id"`
	CreatedAt        time.Time       `json:"created_at"`
	LinkedAccounts   []LinkedAccount `json:"linked_accounts"`
	WalletAddress    string          `json:"-"` // Extracted from linked_accounts
	Email            string          `json:"-"` // Extracted if available
	MFAEnabled       bool            `json:"mfa_enabled"`
	HasAcceptedTerms bool            `json:"has_accepted_terms"`
}

// LinkedAccount representa uma conta vinculada ao usuário
type LinkedAccount struct {
	Type            string    `json:"type"`                    // wallet, email, discord, twitter, etc
	Address         string    `json:"address,omitempty"`       // For wallet
	Email           string    `json:"email,omitempty"`         // For email
	ChainType       string    `json:"chain_type,omitempty"`    // ethereum, base, polygon
	WalletClient    string    `json:"wallet_client,omitempty"` // metamask, coinbase_wallet, etc
	VerifiedAt      time.Time `json:"verified_at"`
	FirstVerifiedAt time.Time `json:"first_verified_at"`
}

// VerifyTokenRequest é a requisição para verificar um token Privy
type VerifyTokenRequest struct {
	Token string `json:"token"`
}

// VerifyTokenResponse é a resposta da verificação de token
type VerifyTokenResponse struct {
	User  PrivyUser `json:"user"`
	Valid bool      `json:"valid"`
}

// VerifyToken verifica um token de autenticação Privy
func (pc *PrivyClient) VerifyToken(token string) (*PrivyUser, error) {
	url := fmt.Sprintf("%s/api/v1/verification_keys", pc.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Privy usa apenas App ID para autenticação
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("privy-app-id", pc.appID)

	pc.logger.Info("Verificando token Privy")
	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar Privy API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Privy API retornou erro %d: %s", resp.StatusCode, string(body))
	}

	var user PrivyUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	// Extrair wallet address e email dos linked accounts
	for _, account := range user.LinkedAccounts {
		if account.Type == "wallet" && user.WalletAddress == "" {
			user.WalletAddress = account.Address
		}
		if account.Type == "email" && user.Email == "" {
			user.Email = account.Email
		}
	}

	pc.logger.Info("Token verificado com sucesso", "user_id", user.ID, "wallet", user.WalletAddress)
	return &user, nil
}

// GetUser obtém informações de um usuário pelo ID
func (pc *PrivyClient) GetUser(userID string) (*PrivyUser, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s", pc.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("privy-app-id", pc.appID)

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar Privy API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Privy API retornou erro %d: %s", resp.StatusCode, string(body))
	}

	var user PrivyUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	// Extrair wallet address dos linked accounts
	for _, account := range user.LinkedAccounts {
		if account.Type == "wallet" && user.WalletAddress == "" {
			user.WalletAddress = account.Address
		}
		if account.Type == "email" && user.Email == "" {
			user.Email = account.Email
		}
	}

	return &user, nil
}

// LinkWallet vincula uma wallet a um usuário
func (pc *PrivyClient) LinkWallet(userID, walletAddress, signature string) error {
	url := fmt.Sprintf("%s/api/v1/users/%s/wallets", pc.baseURL, userID)

	payload := map[string]string{
		"address":    walletAddress,
		"signature":  signature,
		"chain_type": "base", // Base Network
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("privy-app-id", pc.appID)

	pc.logger.Info("Vinculando wallet", "user_id", userID, "wallet", walletAddress)
	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao chamar Privy API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Privy API retornou erro %d: %s", resp.StatusCode, string(body))
	}

	pc.logger.Info("Wallet vinculada com sucesso")
	return nil
}

// GetWalletsByUser obtém todas as wallets de um usuário
func (pc *PrivyClient) GetWalletsByUser(userID string) ([]LinkedAccount, error) {
	user, err := pc.GetUser(userID)
	if err != nil {
		return nil, err
	}

	wallets := []LinkedAccount{}
	for _, account := range user.LinkedAccounts {
		if account.Type == "wallet" {
			wallets = append(wallets, account)
		}
	}

	return wallets, nil
}

// ValidateWalletOwnership valida que um usuário é dono de uma wallet
func (pc *PrivyClient) ValidateWalletOwnership(userID, walletAddress string) (bool, error) {
	wallets, err := pc.GetWalletsByUser(userID)
	if err != nil {
		return false, err
	}

	for _, wallet := range wallets {
		if wallet.Address == walletAddress {
			return true, nil
		}
	}

	return false, nil
}

// CreateEmbeddedWallet cria uma embedded wallet para um usuário
func (pc *PrivyClient) CreateEmbeddedWallet(userID string) (*LinkedAccount, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/embedded_wallet", pc.baseURL, userID)

	payload := map[string]string{
		"chain_type": "base", // Base Network
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("privy-app-id", pc.appID)

	pc.logger.Info("Criando embedded wallet", "user_id", userID)
	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar Privy API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Privy API retornou erro %d: %s", resp.StatusCode, string(body))
	}

	var wallet LinkedAccount
	if err := json.Unmarshal(body, &wallet); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %w", err)
	}

	pc.logger.Info("Embedded wallet criada", "address", wallet.Address)
	return &wallet, nil
}
