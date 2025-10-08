package web3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// NFTAccessManager gerencia NFTs de acesso ao bot
type NFTAccessManager struct {
	config       *config.Config
	logger       *logger.Logger
	baseClient   *BaseClient
	contractAddr common.Address
}

// NewNFTAccessManager cria um novo gerenciador de NFT de acesso
func NewNFTAccessManager(cfg *config.Config, log *logger.Logger, baseClient *BaseClient) *NFTAccessManager {
	contractAddr := common.HexToAddress(cfg.Web3.NFTAccessContractAddress)

	return &NFTAccessManager{
		config:       cfg,
		logger:       log,
		baseClient:   baseClient,
		contractAddr: contractAddr,
	}
}

// NFTAccessTier representa um tier de acesso
type NFTAccessTier struct {
	TierID      uint8   `json:"tier_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       *big.Int `json:"price_wei"`
	PriceUSD    string  `json:"price_usd"`
	MaxSupply   uint64  `json:"max_supply"`
	CurrentSupply uint64 `json:"current_supply"`
	Benefits    []string `json:"benefits"`
	IsActive    bool    `json:"is_active"`
}

// NFTAccess representa um NFT de acesso
type NFTAccess struct {
	TokenID     *big.Int `json:"token_id"`
	Owner       string   `json:"owner"`
	Tier        NFTAccessTier `json:"tier"`
	MintedAt    uint64   `json:"minted_at"`
	ExpiresAt   uint64   `json:"expires_at,omitempty"` // 0 = never expires
	IsActive    bool     `json:"is_active"`
	Metadata    string   `json:"metadata"` // IPFS URI
}

// GetAccessTiers retorna todos os tiers de acesso disponíveis
func (nam *NFTAccessManager) GetAccessTiers(ctx context.Context) ([]NFTAccessTier, error) {
	// Hardcoded tiers - em produção, viriam do smart contract
	tiers := []NFTAccessTier{
		{
			TierID:      1,
			Name:        "Basic Access",
			Description: "Acesso básico ao IaC AI Agent",
			Price:       big.NewInt(10000000000000000), // 0.01 ETH
			PriceUSD:    "25.00",
			MaxSupply:   10000,
			CurrentSupply: 0,
			Benefits: []string{
				"Análises ilimitadas de Terraform",
				"Detecção de segurança com Checkov",
				"Sugestões básicas de otimização",
				"Suporte via Discord",
			},
			IsActive: true,
		},
		{
			TierID:      2,
			Name:        "Pro Access",
			Description: "Acesso profissional com AI avançada",
			Price:       big.NewInt(50000000000000000), // 0.05 ETH
			PriceUSD:    "125.00",
			MaxSupply:   5000,
			CurrentSupply: 0,
			Benefits: []string{
				"Tudo do Basic Access",
				"Análise com LLM (GPT-4/Claude)",
				"Sugestões contextualizadas inteligentes",
				"Análise de Preview e Drift",
				"Detecção de Secrets",
				"Recomendações de arquitetura",
				"Priority support",
			},
			IsActive: true,
		},
		{
			TierID:      3,
			Name:        "Enterprise Access",
			Description: "Acesso enterprise com features exclusivas",
			Price:       big.NewInt(200000000000000000), // 0.2 ETH
			PriceUSD:    "500.00",
			MaxSupply:   1000,
			CurrentSupply: 0,
			Benefits: []string{
				"Tudo do Pro Access",
				"API dedicada com rate limits maiores",
				"Custom knowledge base",
				"Integração com CI/CD privado",
				"Suporte dedicado 24/7",
				"SLA garantido",
				"Governance tokens inclusos",
			},
			IsActive: true,
		},
	}

	nam.logger.Info("Tiers de acesso carregados", "count", len(tiers))
	return tiers, nil
}

// CheckAccess verifica se uma wallet tem acesso ao bot
func (nam *NFTAccessManager) CheckAccess(ctx context.Context, walletAddress string) (*NFTAccess, error) {
	addr := common.HexToAddress(walletAddress)
	
	// TODO: Chamar smart contract para verificar balance
	// Por ora, simula a verificação
	
	nam.logger.Info("Verificando acesso", "wallet", walletAddress)
	
	// Simulação: verifica se tem algum NFT
	// Em produção, isso seria uma chamada ao contrato
	balance := nam.getSimulatedBalance(addr)
	
	if balance.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("wallet não possui NFT de acesso")
	}

	// Retorna o NFT (simulado)
	tiers, _ := nam.GetAccessTiers(ctx)
	return &NFTAccess{
		TokenID:  big.NewInt(123), // Simulado
		Owner:    walletAddress,
		Tier:     tiers[1], // Pro tier
		MintedAt: 1700000000,
		ExpiresAt: 0, // Never expires
		IsActive: true,
		Metadata: "ipfs://QmExample...",
	}, nil
}

// MintAccessNFT minta um novo NFT de acesso
func (nam *NFTAccessManager) MintAccessNFT(ctx context.Context, walletAddress string, tierID uint8) (*NFTAccess, error) {
	// Get tier info
	tiers, err := nam.GetAccessTiers(ctx)
	if err != nil {
		return nil, err
	}

	var selectedTier *NFTAccessTier
	for _, tier := range tiers {
		if tier.TierID == tierID {
			selectedTier = &tier
			break
		}
	}

	if selectedTier == nil {
		return nil, fmt.Errorf("tier %d não encontrado", tierID)
	}

	if !selectedTier.IsActive {
		return nil, fmt.Errorf("tier %d não está ativo", tierID)
	}

	nam.logger.Info("Mintando NFT de acesso", 
		"wallet", walletAddress, 
		"tier", selectedTier.Name,
		"price", selectedTier.Price.String())

	// TODO: Chamar smart contract para mintar
	// Por ora, simula o mint

	// Em produção, isso seria:
	// 1. Verificar pagamento
	// 2. Chamar contract.mint(walletAddress, tierID)
	// 3. Aguardar confirmação
	// 4. Retornar NFT mintado

	nft := &NFTAccess{
		TokenID:  big.NewInt(12345), // Seria retornado pelo contrato
		Owner:    walletAddress,
		Tier:     *selectedTier,
		MintedAt: uint64(1700000000), // timestamp atual
		ExpiresAt: 0, // Never expires
		IsActive: true,
		Metadata: "ipfs://QmExample...",
	}

	nam.logger.Info("NFT de acesso mintado com sucesso", "token_id", nft.TokenID.String())
	return nft, nil
}

// TransferAccess transfere um NFT de acesso
func (nam *NFTAccessManager) TransferAccess(ctx context.Context, from, to string, tokenID *big.Int) error {
	nam.logger.Info("Transferindo NFT de acesso",
		"from", from,
		"to", to,
		"token_id", tokenID.String())

	// TODO: Chamar smart contract para transferir
	// contract.transferFrom(from, to, tokenID)

	nam.logger.Info("NFT transferido com sucesso")
	return nil
}

// RevokeAccess revoga um NFT de acesso (apenas admin)
func (nam *NFTAccessManager) RevokeAccess(ctx context.Context, tokenID *big.Int) error {
	nam.logger.Warn("Revogando acesso", "token_id", tokenID.String())

	// TODO: Chamar smart contract para revogar
	// contract.revoke(tokenID)

	return nil
}

// GetAccessByTokenID obtém informações de um NFT específico
func (nam *NFTAccessManager) GetAccessByTokenID(ctx context.Context, tokenID *big.Int) (*NFTAccess, error) {
	// TODO: Chamar smart contract
	// contract.tokenURI(tokenID)
	// contract.ownerOf(tokenID)

	tiers, _ := nam.GetAccessTiers(ctx)
	return &NFTAccess{
		TokenID:  tokenID,
		Owner:    "0x...",
		Tier:     tiers[0],
		MintedAt: 1700000000,
		ExpiresAt: 0,
		IsActive: true,
		Metadata: "ipfs://QmExample...",
	}, nil
}

// ListAccessNFTs lista todos os NFTs de uma wallet
func (nam *NFTAccessManager) ListAccessNFTs(ctx context.Context, walletAddress string) ([]*NFTAccess, error) {
	addr := common.HexToAddress(walletAddress)

	nam.logger.Info("Listando NFTs de acesso", "wallet", walletAddress)

	// TODO: Chamar smart contract
	// balance := contract.balanceOf(addr)
	// for i := 0; i < balance; i++ {
	//     tokenId := contract.tokenOfOwnerByIndex(addr, i)
	//     ...
	// }

	// Simulação
	balance := nam.getSimulatedBalance(addr)
	if balance.Cmp(big.NewInt(0)) == 0 {
		return []*NFTAccess{}, nil
	}

	tiers, _ := nam.GetAccessTiers(ctx)
	return []*NFTAccess{
		{
			TokenID:  big.NewInt(123),
			Owner:    walletAddress,
			Tier:     tiers[1],
			MintedAt: 1700000000,
			ExpiresAt: 0,
			IsActive: true,
			Metadata: "ipfs://QmExample...",
		},
	}, nil
}

// getSimulatedBalance simula o balance (apenas para desenvolvimento)
func (nam *NFTAccessManager) getSimulatedBalance(addr common.Address) *big.Int {
	// Em produção, seria: contract.balanceOf(addr)
	return big.NewInt(1) // Simula que tem 1 NFT
}

// UpgradeAccess faz upgrade de um NFT para um tier superior
func (nam *NFTAccessManager) UpgradeAccess(ctx context.Context, tokenID *big.Int, newTierID uint8) error {
	current, err := nam.GetAccessByTokenID(ctx, tokenID)
	if err != nil {
		return err
	}

	if current.Tier.TierID >= newTierID {
		return fmt.Errorf("tier atual já é igual ou superior ao tier solicitado")
	}

	tiers, err := nam.GetAccessTiers(ctx)
	if err != nil {
		return err
	}

	var newTier *NFTAccessTier
	for _, tier := range tiers {
		if tier.TierID == newTierID {
			newTier = &tier
			break
		}
	}

	if newTier == nil {
		return fmt.Errorf("tier %d não encontrado", newTierID)
	}

	// Calcula diferença de preço
	priceDiff := new(big.Int).Sub(newTier.Price, current.Tier.Price)

	nam.logger.Info("Upgrade de tier",
		"token_id", tokenID.String(),
		"from_tier", current.Tier.Name,
		"to_tier", newTier.Name,
		"price_diff", priceDiff.String())

	// TODO: Chamar smart contract para upgrade
	// Requer pagamento da diferença

	return nil
}

// ValidateAccess valida se uma wallet tem acesso válido
func (nam *NFTAccessManager) ValidateAccess(ctx context.Context, walletAddress string, requiredTier uint8) (bool, error) {
	nft, err := nam.CheckAccess(ctx, walletAddress)
	if err != nil {
		return false, nil // Não tem acesso
	}

	// Verifica se está ativo
	if !nft.IsActive {
		return false, fmt.Errorf("NFT de acesso está inativo")
	}

	// Verifica se expirou
	if nft.ExpiresAt > 0 && uint64(1700000000) > nft.ExpiresAt {
		return false, fmt.Errorf("NFT de acesso expirou")
	}

	// Verifica tier
	if nft.Tier.TierID < requiredTier {
		return false, fmt.Errorf("tier de acesso insuficiente: tem %d, requer %d", nft.Tier.TierID, requiredTier)
	}

	return true, nil
}

// GetContractAddress retorna o endereço do contrato de NFT
func (nam *NFTAccessManager) GetContractAddress() string {
	return nam.contractAddr.Hex()
}

// EstimateMintGas estima o gas necessário para mintar um NFT
func (nam *NFTAccessManager) EstimateMintGas(ctx context.Context, walletAddress string, tierID uint8) (uint64, error) {
	// TODO: Estimar gas real do contrato
	// Em Base Network, transações são muito baratas
	
	// Estimativa conservadora
	gasEstimate := uint64(150000) // ~150k gas

	nam.logger.Info("Gas estimado para mint", 
		"wallet", walletAddress,
		"tier", tierID,
		"gas", gasEstimate)

	return gasEstimate, nil
}
