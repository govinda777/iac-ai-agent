package web3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// BotTokenManager gerencia o token do bot (ERC-20)
type BotTokenManager struct {
	config       *config.Config
	logger       *logger.Logger
	baseClient   *BaseClient
	contractAddr common.Address
	tokenSymbol  string
	tokenName    string
	decimals     uint8
}

// NewBotTokenManager cria um novo gerenciador de tokens do bot
func NewBotTokenManager(cfg *config.Config, log *logger.Logger, baseClient *BaseClient) *BotTokenManager {
	contractAddr := common.HexToAddress(cfg.Web3.BotTokenContractAddress)

	return &BotTokenManager{
		config:       cfg,
		logger:       log,
		baseClient:   baseClient,
		contractAddr: contractAddr,
		tokenSymbol:  "IACAI", // IaC AI Token
		tokenName:    "IaC AI Agent Token",
		decimals:     18,
	}
}

// TokenInfo contém informações sobre o token
type TokenInfo struct {
	Address      string   `json:"address"`
	Name         string   `json:"name"`
	Symbol       string   `json:"symbol"`
	Decimals     uint8    `json:"decimals"`
	TotalSupply  *big.Int `json:"total_supply"`
	TotalSupplyFormatted string `json:"total_supply_formatted"`
	Price        *TokenPrice `json:"price,omitempty"`
}

// TokenPrice contém informações de preço
type TokenPrice struct {
	USD         string `json:"usd"`
	ETH         string `json:"eth"`
	LastUpdated int64  `json:"last_updated"`
}

// TokenBalance representa o saldo de tokens de uma wallet
type TokenBalance struct {
	Address         string   `json:"address"`
	Balance         *big.Int `json:"balance"`
	BalanceFormatted string  `json:"balance_formatted"`
	ValueUSD        string   `json:"value_usd,omitempty"`
}

// TokenPackage representa um pacote de tokens para venda
type TokenPackage struct {
	PackageID   uint8    `json:"package_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TokenAmount *big.Int `json:"token_amount"`
	TokenAmountFormatted string `json:"token_amount_formatted"`
	Price       *big.Int `json:"price_wei"`
	PriceUSD    string   `json:"price_usd"`
	Discount    string   `json:"discount,omitempty"`
	IsActive    bool     `json:"is_active"`
}

// GetTokenInfo retorna informações sobre o token
func (btm *BotTokenManager) GetTokenInfo(ctx context.Context) (*TokenInfo, error) {
	// TODO: Buscar do smart contract
	// totalSupply := contract.totalSupply()
	
	totalSupply := new(big.Int)
	totalSupply.SetString("1000000000000000000000000", 10) // 1M tokens

	return &TokenInfo{
		Address:      btm.contractAddr.Hex(),
		Name:         btm.tokenName,
		Symbol:       btm.tokenSymbol,
		Decimals:     btm.decimals,
		TotalSupply:  totalSupply,
		TotalSupplyFormatted: btm.formatTokenAmount(totalSupply),
		Price: &TokenPrice{
			USD:         "0.10",
			ETH:         "0.00004",
			LastUpdated: 1700000000,
		},
	}, nil
}

// GetBalance obtém o saldo de tokens de uma wallet
func (btm *BotTokenManager) GetBalance(ctx context.Context, walletAddress string) (*TokenBalance, error) {
	addr := common.HexToAddress(walletAddress)

	btm.logger.Info("Obtendo saldo de tokens", "wallet", walletAddress)

	// TODO: Chamar smart contract
	// balance := contract.balanceOf(addr)

	// Simulação
	balance := new(big.Int)
	balance.SetString("1000000000000000000000", 10) // 1000 tokens

	return &TokenBalance{
		Address:          walletAddress,
		Balance:          balance,
		BalanceFormatted: btm.formatTokenAmount(balance),
		ValueUSD:         "100.00",
	}, nil
}

// GetTokenPackages retorna os pacotes de tokens disponíveis
func (btm *BotTokenManager) GetTokenPackages(ctx context.Context) ([]TokenPackage, error) {
	packages := []TokenPackage{
		{
			PackageID:   1,
			Name:        "Starter Pack",
			Description: "Pacote inicial para começar",
			TokenAmount: btm.parseTokenAmount("100"),
			TokenAmountFormatted: "100 IACAI",
			Price:       big.NewInt(5000000000000000), // 0.005 ETH
			PriceUSD:    "10.00",
			Discount:    "",
			IsActive:    true,
		},
		{
			PackageID:   2,
			Name:        "Power Pack",
			Description: "Pacote popular com 10% de desconto",
			TokenAmount: btm.parseTokenAmount("500"),
			TokenAmountFormatted: "500 IACAI",
			Price:       big.NewInt(22500000000000000), // 0.0225 ETH
			PriceUSD:    "45.00",
			Discount:    "10%",
			IsActive:    true,
		},
		{
			PackageID:   3,
			Name:        "Pro Pack",
			Description: "Pacote profissional com 15% de desconto",
			TokenAmount: btm.parseTokenAmount("1000"),
			TokenAmountFormatted: "1000 IACAI",
			Price:       big.NewInt(42500000000000000), // 0.0425 ETH
			PriceUSD:    "85.00",
			Discount:    "15%",
			IsActive:    true,
		},
		{
			PackageID:   4,
			Name:        "Enterprise Pack",
			Description: "Pacote enterprise com 25% de desconto",
			TokenAmount: btm.parseTokenAmount("5000"),
			TokenAmountFormatted: "5000 IACAI",
			Price:       big.NewInt(187500000000000000), // 0.1875 ETH
			PriceUSD:    "375.00",
			Discount:    "25%",
			IsActive:    true,
		},
	}

	btm.logger.Info("Pacotes de tokens carregados", "count", len(packages))
	return packages, nil
}

// BuyTokens compra tokens usando ETH
func (btm *BotTokenManager) BuyTokens(ctx context.Context, walletAddress string, packageID uint8) (*Transaction, error) {
	packages, err := btm.GetTokenPackages(ctx)
	if err != nil {
		return nil, err
	}

	var selectedPackage *TokenPackage
	for _, pkg := range packages {
		if pkg.PackageID == packageID {
			selectedPackage = &pkg
			break
		}
	}

	if selectedPackage == nil {
		return nil, fmt.Errorf("pacote %d não encontrado", packageID)
	}

	if !selectedPackage.IsActive {
		return nil, fmt.Errorf("pacote %d não está ativo", packageID)
	}

	btm.logger.Info("Comprando tokens",
		"wallet", walletAddress,
		"package", selectedPackage.Name,
		"amount", selectedPackage.TokenAmountFormatted,
		"price", selectedPackage.Price.String())

	// TODO: Criar transação real
	// 1. Usuário aprova gasto de ETH
	// 2. Chamar contract.buyTokens(packageID) com value = price
	// 3. Aguardar confirmação
	// 4. Tokens são transferidos para a wallet

	// Simulação de transação
	txHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

	btm.logger.Info("Tokens comprados com sucesso", "tx_hash", txHash)

	return &Transaction{
		Hash:     txHash,
		From:     walletAddress,
		To:       btm.contractAddr.Hex(),
		Value:    selectedPackage.Price,
		ValueETH: btm.formatEth(selectedPackage.Price),
		Status:   1,
	}, nil
}

// Transfer transfere tokens entre wallets
func (btm *BotTokenManager) Transfer(ctx context.Context, from, to string, amount *big.Int) (*Transaction, error) {
	fromAddr := common.HexToAddress(from)
	toAddr := common.HexToAddress(to)

	btm.logger.Info("Transferindo tokens",
		"from", from,
		"to", to,
		"amount", btm.formatTokenAmount(amount))

	// TODO: Chamar smart contract
	// contract.transfer(to, amount)

	// Verificar saldo suficiente
	balance, err := btm.GetBalance(ctx, from)
	if err != nil {
		return nil, err
	}

	if balance.Balance.Cmp(amount) < 0 {
		return nil, fmt.Errorf("saldo insuficiente: tem %s, precisa %s",
			btm.formatTokenAmount(balance.Balance),
			btm.formatTokenAmount(amount))
	}

	txHash := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"

	btm.logger.Info("Tokens transferidos com sucesso", "tx_hash", txHash)

	return &Transaction{
		Hash:  txHash,
		From:  from,
		To:    to,
		Value: amount,
		Status: 1,
	}, nil
}

// Approve aprova gasto de tokens por outro endereço
func (btm *BotTokenManager) Approve(ctx context.Context, owner, spender string, amount *big.Int) error {
	ownerAddr := common.HexToAddress(owner)
	spenderAddr := common.HexToAddress(spender)

	btm.logger.Info("Aprovando gasto de tokens",
		"owner", owner,
		"spender", spender,
		"amount", btm.formatTokenAmount(amount))

	// TODO: Chamar smart contract
	// contract.approve(spender, amount)

	return nil
}

// GetAllowance obtém o allowance de um spender
func (btm *BotTokenManager) GetAllowance(ctx context.Context, owner, spender string) (*big.Int, error) {
	ownerAddr := common.HexToAddress(owner)
	spenderAddr := common.HexToAddress(spender)

	btm.logger.Info("Obtendo allowance",
		"owner", owner,
		"spender", spender)

	// TODO: Chamar smart contract
	// allowance := contract.allowance(owner, spender)

	// Simulação
	allowance := big.NewInt(0)
	return allowance, nil
}

// SpendTokens gasta tokens do usuário (usado pelo bot para cobrar por análises)
func (btm *BotTokenManager) SpendTokens(ctx context.Context, userWallet string, amount *big.Int, reason string) error {
	btm.logger.Info("Gastando tokens do usuário",
		"wallet", userWallet,
		"amount", btm.formatTokenAmount(amount),
		"reason", reason)

	// TODO: Implementar lógica de cobrança
	// 1. Verificar saldo
	// 2. Verificar allowance para o contrato do bot
	// 3. Transferir tokens para treasury
	// 4. Registrar gasto no histórico

	balance, err := btm.GetBalance(ctx, userWallet)
	if err != nil {
		return err
	}

	if balance.Balance.Cmp(amount) < 0 {
		return fmt.Errorf("saldo insuficiente de tokens")
	}

	btm.logger.Info("Tokens gastos com sucesso")
	return nil
}

// GetTokenPrice obtém o preço atual do token
func (btm *BotTokenManager) GetTokenPrice(ctx context.Context) (*TokenPrice, error) {
	// TODO: Integrar com oracle de preços (Chainlink, etc)
	// ou com DEX (Uniswap/SushiSwap na Base)

	return &TokenPrice{
		USD:         "0.10",
		ETH:         "0.00004",
		LastUpdated: 1700000000,
	}, nil
}

// CalculateTokenCost calcula o custo em tokens para uma operação
func (btm *BotTokenManager) CalculateTokenCost(operationType string) (*big.Int, error) {
	// Tabela de preços por tipo de operação
	costs := map[string]string{
		"terraform_analysis": "1",    // 1 token
		"checkov_scan":       "2",    // 2 tokens
		"llm_analysis":       "5",    // 5 tokens (mais caro, usa LLM)
		"preview_analysis":   "3",    // 3 tokens
		"security_audit":     "10",   // 10 tokens (auditoria completa)
		"cost_optimization":  "5",    // 5 tokens
		"full_review":        "15",   // 15 tokens (review completo)
	}

	costStr, ok := costs[operationType]
	if !ok {
		return nil, fmt.Errorf("tipo de operação desconhecido: %s", operationType)
	}

	return btm.parseTokenAmount(costStr), nil
}

// Helper functions

func (btm *BotTokenManager) parseTokenAmount(amount string) *big.Int {
	// Converte string para big.Int com decimais
	value := new(big.Int)
	value.SetString(amount, 10)
	
	// Multiplica por 10^decimals
	multiplier := new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(int64(btm.decimals)),
		nil,
	)
	
	return value.Mul(value, multiplier)
}

func (btm *BotTokenManager) formatTokenAmount(amount *big.Int) string {
	// Converte big.Int para string com decimais
	divisor := new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(int64(btm.decimals)),
		nil,
	)

	value := new(big.Float).Quo(
		new(big.Float).SetInt(amount),
		new(big.Float).SetInt(divisor),
	)

	return fmt.Sprintf("%s %s", value.Text('f', 2), btm.tokenSymbol)
}

func (btm *BotTokenManager) formatEth(amount *big.Int) string {
	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(amount),
		big.NewFloat(1e18),
	)
	return fmt.Sprintf("%s ETH", ethValue.Text('f', 6))
}

// GetContractAddress retorna o endereço do contrato de token
func (btm *BotTokenManager) GetContractAddress() string {
	return btm.contractAddr.Hex()
}
