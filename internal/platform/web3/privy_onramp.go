package web3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// PrivyOnrampManager gerencia compras de crypto via Privy Onramp
type PrivyOnrampManager struct {
	config      *config.Config
	logger      *logger.Logger
	privyClient *PrivyClient
}

// NewPrivyOnrampManager cria um novo gerenciador de onramp
func NewPrivyOnrampManager(cfg *config.Config, log *logger.Logger, privyClient *PrivyClient) *PrivyOnrampManager {
	return &PrivyOnrampManager{
		config:      cfg,
		logger:      log,
		privyClient: privyClient,
	}
}

// OnrampQuote representa uma cotação de onramp
type OnrampQuote struct {
	QuoteID         string   `json:"quote_id"`
	SourceAmount    string   `json:"source_amount"` // Valor em moeda fiat
	SourceCurrency  string   `json:"source_currency"` // USD, EUR, BRL, etc
	TargetAmount    string   `json:"target_amount"` // Valor em crypto
	TargetCurrency  string   `json:"target_currency"` // ETH, USDC, etc
	NetworkFee      string   `json:"network_fee"`
	ProviderFee     string   `json:"provider_fee"`
	TotalCost       string   `json:"total_cost"`
	ExchangeRate    string   `json:"exchange_rate"`
	EstimatedTime   string   `json:"estimated_time"` // "5-10 minutes"
	Provider        string   `json:"provider"` // moonpay, transak, etc
	ExpiresAt       int64    `json:"expires_at"`
}

// OnrampTransaction representa uma transação de onramp
type OnrampTransaction struct {
	TransactionID   string `json:"transaction_id"`
	QuoteID         string `json:"quote_id"`
	UserID          string `json:"user_id"`
	WalletAddress   string `json:"wallet_address"`
	Status          string `json:"status"` // pending, processing, completed, failed
	SourceAmount    string `json:"source_amount"`
	SourceCurrency  string `json:"source_currency"`
	TargetAmount    string `json:"target_amount"`
	TargetCurrency  string `json:"target_currency"`
	PaymentMethod   string `json:"payment_method"` // credit_card, debit_card, bank_transfer, pix
	TxHash          string `json:"tx_hash,omitempty"` // Hash da transação on-chain
	CreatedAt       int64  `json:"created_at"`
	CompletedAt     int64  `json:"completed_at,omitempty"`
	ErrorMessage    string `json:"error_message,omitempty"`
}

// OnrampSession representa uma sessão de compra
type OnrampSession struct {
	SessionID      string                 `json:"session_id"`
	UserID         string                 `json:"user_id"`
	WalletAddress  string                 `json:"wallet_address"`
	Purpose        string                 `json:"purpose"` // nft_access, bot_tokens
	RequiredAmount *big.Int               `json:"required_amount"`
	Currency       string                 `json:"currency"`
	Status         string                 `json:"status"` // created, payment_pending, completed, cancelled
	Quote          *OnrampQuote           `json:"quote,omitempty"`
	Transaction    *OnrampTransaction     `json:"transaction,omitempty"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      int64                  `json:"created_at"`
	ExpiresAt      int64                  `json:"expires_at"`
}

// CreateOnrampSessionRequest é a requisição para criar uma sessão de onramp
type CreateOnrampSessionRequest struct {
	UserID         string                 `json:"user_id"`
	WalletAddress  string                 `json:"wallet_address"`
	Purpose        string                 `json:"purpose"` // "nft_access" ou "bot_tokens"
	TargetItemID   string                 `json:"target_item_id"` // tier_id ou package_id
	SourceCurrency string                 `json:"source_currency"` // USD, EUR, BRL
	PaymentMethod  string                 `json:"payment_method,omitempty"` // Opcional
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// CreateOnrampSession cria uma nova sessão de onramp
func (pom *PrivyOnrampManager) CreateOnrampSession(ctx context.Context, req *CreateOnrampSessionRequest) (*OnrampSession, error) {
	pom.logger.Info("Criando sessão de onramp",
		"user_id", req.UserID,
		"wallet", req.WalletAddress,
		"purpose", req.Purpose,
		"item_id", req.TargetItemID)

	// Valida wallet
	user, err := pom.privyClient.GetUser(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado: %w", err)
	}

	// Valida ownership da wallet
	owned, err := pom.privyClient.ValidateWalletOwnership(req.UserID, req.WalletAddress)
	if err != nil || !owned {
		return nil, fmt.Errorf("wallet não pertence ao usuário")
	}

	// Determina valor necessário baseado no purpose
	var requiredAmount *big.Int
	var targetCurrency string
	var itemDescription string

	switch req.Purpose {
	case "nft_access":
		// Buscar preço do NFT tier
		requiredAmount = big.NewInt(50000000000000000) // 0.05 ETH (exemplo)
		targetCurrency = "ETH"
		itemDescription = fmt.Sprintf("NFT Access Tier %s", req.TargetItemID)

	case "bot_tokens":
		// Buscar preço do pacote de tokens
		requiredAmount = big.NewInt(22500000000000000) // 0.0225 ETH (exemplo)
		targetCurrency = "ETH"
		itemDescription = fmt.Sprintf("Bot Token Package %s", req.TargetItemID)

	default:
		return nil, fmt.Errorf("purpose inválido: %s", req.Purpose)
	}

	// Cria cotação
	quote, err := pom.getQuote(ctx, requiredAmount, targetCurrency, req.SourceCurrency)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter cotação: %w", err)
	}

	// Cria sessão
	session := &OnrampSession{
		SessionID:      fmt.Sprintf("onramp_%d", 1700000000), // UUID em produção
		UserID:         req.UserID,
		WalletAddress:  req.WalletAddress,
		Purpose:        req.Purpose,
		RequiredAmount: requiredAmount,
		Currency:       targetCurrency,
		Status:         "created",
		Quote:          quote,
		Metadata: map[string]interface{}{
			"item_id":         req.TargetItemID,
			"item_description": itemDescription,
			"source_currency": req.SourceCurrency,
			"user_email":      user.Email,
		},
		CreatedAt: 1700000000,
		ExpiresAt: 1700003600, // 1 hora
	}

	pom.logger.Info("Sessão de onramp criada",
		"session_id", session.SessionID,
		"quote_id", quote.QuoteID,
		"required_amount", requiredAmount.String())

	return session, nil
}

// InitiatePayment inicia o pagamento via Privy Onramp
func (pom *PrivyOnrampManager) InitiatePayment(ctx context.Context, sessionID string, paymentMethod string) (*OnrampTransaction, error) {
	// TODO: Em produção, buscar sessão do banco de dados
	
	pom.logger.Info("Iniciando pagamento",
		"session_id", sessionID,
		"payment_method", paymentMethod)

	// Valida método de pagamento
	validMethods := map[string]bool{
		"credit_card":   true,
		"debit_card":    true,
		"bank_transfer": true,
		"pix":           true,
		"apple_pay":     true,
		"google_pay":    true,
	}

	if !validMethods[paymentMethod] {
		return nil, fmt.Errorf("método de pagamento inválido: %s", paymentMethod)
	}

	// TODO: Integrar com Privy Onramp API
	// 1. Criar transação no provedor (MoonPay, Transak, etc)
	// 2. Retornar URL/widget para o usuário completar pagamento
	// 3. Receber webhook quando pagamento for completado
	// 4. Processar compra do NFT/Tokens

	transaction := &OnrampTransaction{
		TransactionID:  fmt.Sprintf("tx_onramp_%d", 1700000000),
		QuoteID:        "quote_123",
		UserID:         "user_123",
		WalletAddress:  "0x...",
		Status:         "pending",
		SourceAmount:   "125.00",
		SourceCurrency: "USD",
		TargetAmount:   "0.05",
		TargetCurrency: "ETH",
		PaymentMethod:  paymentMethod,
		CreatedAt:      1700000000,
	}

	pom.logger.Info("Pagamento iniciado",
		"transaction_id", transaction.TransactionID,
		"status", transaction.Status)

	return transaction, nil
}

// GetOnrampStatus obtém o status de uma transação de onramp
func (pom *PrivyOnrampManager) GetOnrampStatus(ctx context.Context, transactionID string) (*OnrampTransaction, error) {
	pom.logger.Info("Verificando status de onramp", "transaction_id", transactionID)

	// TODO: Buscar do banco de dados ou API do provedor

	transaction := &OnrampTransaction{
		TransactionID:  transactionID,
		Status:         "completed",
		SourceAmount:   "125.00",
		SourceCurrency: "USD",
		TargetAmount:   "0.05",
		TargetCurrency: "ETH",
		TxHash:         "0xabcdef...",
		CompletedAt:    1700001000,
	}

	return transaction, nil
}

// ProcessOnrampCompletion processa a conclusão de um onramp
func (pom *PrivyOnrampManager) ProcessOnrampCompletion(ctx context.Context, transactionID string) error {
	pom.logger.Info("Processando conclusão de onramp", "transaction_id", transactionID)

	// TODO: 
	// 1. Verificar que pagamento foi recebido
	// 2. Executar a ação correspondente:
	//    - Se purpose = nft_access: mintar NFT
	//    - Se purpose = bot_tokens: transferir tokens
	// 3. Atualizar status da sessão
	// 4. Notificar usuário

	pom.logger.Info("Onramp processado com sucesso")
	return nil
}

// GetSupportedCurrencies retorna moedas fiat suportadas
func (pom *PrivyOnrampManager) GetSupportedCurrencies() []string {
	return []string{
		"USD", // Dólar Americano
		"EUR", // Euro
		"GBP", // Libra Esterlina
		"BRL", // Real Brasileiro
		"AUD", // Dólar Australiano
		"CAD", // Dólar Canadense
		"JPY", // Iene Japonês
		"MXN", // Peso Mexicano
	}
}

// GetSupportedPaymentMethods retorna métodos de pagamento suportados
func (pom *PrivyOnrampManager) GetSupportedPaymentMethods() []string {
	return []string{
		"credit_card",
		"debit_card",
		"bank_transfer",
		"pix", // Brasil
		"apple_pay",
		"google_pay",
	}
}

// GetSupportedCryptos retorna cryptos suportadas para compra
func (pom *PrivyOnrampManager) GetSupportedCryptos() []string {
	return []string{
		"ETH",  // Ethereum (Base Network)
		"USDC", // USD Coin
		"USDT", // Tether
	}
}

// getQuote obtém uma cotação de compra
func (pom *PrivyOnrampManager) getQuote(ctx context.Context, amount *big.Int, targetCurrency, sourceCurrency string) (*OnrampQuote, error) {
	// TODO: Chamar API do provedor de onramp para cotação real
	
	// Simula cotação
	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(amount),
		big.NewFloat(1e18),
	)

	// ETH price: ~$2500
	sourceAmount := new(big.Float).Mul(ethValue, big.NewFloat(2500))

	quote := &OnrampQuote{
		QuoteID:        fmt.Sprintf("quote_%d", 1700000000),
		SourceAmount:   sourceAmount.Text('f', 2),
		SourceCurrency: sourceCurrency,
		TargetAmount:   ethValue.Text('f', 4),
		TargetCurrency: targetCurrency,
		NetworkFee:     "2.50",
		ProviderFee:    "5.00",
		TotalCost:      new(big.Float).Add(sourceAmount, big.NewFloat(7.50)).Text('f', 2),
		ExchangeRate:   "2500.00",
		EstimatedTime:  "5-10 minutes",
		Provider:       "moonpay",
		ExpiresAt:      1700001800, // 30 minutos
	}

	return quote, nil
}

// CancelOnrampSession cancela uma sessão de onramp
func (pom *PrivyOnrampManager) CancelOnrampSession(ctx context.Context, sessionID string) error {
	pom.logger.Info("Cancelando sessão de onramp", "session_id", sessionID)

	// TODO: Atualizar status no banco de dados
	// TODO: Notificar provedor de onramp se necessário

	return nil
}

// GetOnrampHistory obtém histórico de onramps de um usuário
func (pom *PrivyOnrampManager) GetOnrampHistory(ctx context.Context, userID string, limit int) ([]*OnrampTransaction, error) {
	pom.logger.Info("Buscando histórico de onramp", "user_id", userID, "limit", limit)

	// TODO: Buscar do banco de dados

	// Simulação
	history := []*OnrampTransaction{
		{
			TransactionID:  "tx_1",
			Status:         "completed",
			SourceAmount:   "125.00",
			SourceCurrency: "USD",
			TargetAmount:   "0.05",
			TargetCurrency: "ETH",
			PaymentMethod:  "credit_card",
			TxHash:         "0xabc...",
			CreatedAt:      1700000000,
			CompletedAt:    1700000600,
		},
	}

	return history, nil
}
