package integration

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNFTAccessIntegration executa testes integrados de acesso NFT
// Estes testes vão além dos mocks e validam efetivamente o acesso a NFTs Nation Pass
func TestNFTAccessIntegration(t *testing.T) {
	// Skip se não estiver configurado para testes de integração
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests. Set INTEGRATION_TESTS=true to run.")
	}

	// Configuração para testes de integração
	cfg := setupIntegrationConfig(t)
	log := logger.New("debug", "json")

	// Criar cliente Base Network
	baseClient, err := web3.NewBaseClient(cfg, log)
	require.NoError(t, err, "Failed to create Base client")

	// Criar NFTAccessManager
	nftManager := web3.NewNFTAccessManager(cfg, log, baseClient)

	t.Run("RealWalletNationPassValidation", func(t *testing.T) {
		testRealWalletNationPassValidation(t, nftManager, cfg)
	})

	t.Run("ContractInteractionTests", func(t *testing.T) {
		testContractInteraction(t, nftManager, cfg)
	})

	t.Run("AccessTierValidation", func(t *testing.T) {
		testAccessTierValidation(t, nftManager)
	})

	t.Run("NationPassIntegrationFlow", func(t *testing.T) {
		testNationPassIntegrationFlow(t, nftManager, cfg)
	})
}

// testRealWalletNationPassValidation testa validação com wallet real configurada
func testRealWalletNationPassValidation(t *testing.T, nftManager *web3.NFTAccessManager, cfg *config.Config) {
	walletAddress := os.Getenv("WALLET_ADDRESS")
	if walletAddress == "" {
		t.Skip("WALLET_ADDRESS not configured for real wallet tests")
	}

	ctx := context.Background()

	t.Logf("Testing with real wallet: %s", walletAddress)

	// Teste 1: Verificar se a wallet tem acesso
	t.Run("CheckWalletAccess", func(t *testing.T) {
		nft, err := nftManager.CheckAccess(ctx, walletAddress)

		if err != nil {
			t.Logf("Wallet does not have NFT access: %v", err)
			// Isso é esperado se a wallet não tem NFT
			return
		}

		require.NotNil(t, nft, "NFT should not be nil if no error")
		assert.Equal(t, walletAddress, nft.Owner, "NFT owner should match wallet address")
		assert.True(t, nft.IsActive, "NFT should be active")
		assert.Greater(t, nft.Tier.TierID, uint8(0), "Tier ID should be greater than 0")

		t.Logf("Wallet has NFT access - Tier: %s, Token ID: %s",
			nft.Tier.Name, nft.TokenID.String())
	})

	// Teste 2: Listar NFTs da wallet
	t.Run("ListWalletNFTs", func(t *testing.T) {
		nfts, err := nftManager.ListAccessNFTs(ctx, walletAddress)
		require.NoError(t, err, "Should be able to list NFTs")

		if len(nfts) == 0 {
			t.Logf("Wallet %s has no NFTs", walletAddress)
			return
		}

		t.Logf("Wallet has %d NFTs", len(nfts))
		for i, nft := range nfts {
			assert.Equal(t, walletAddress, nft.Owner, "NFT %d owner should match wallet", i)
			assert.True(t, nft.IsActive, "NFT %d should be active", i)
			t.Logf("NFT %d: Tier %s, Token ID %s", i, nft.Tier.Name, nft.TokenID.String())
		}
	})

	// Teste 3: Validar acesso para diferentes tiers
	t.Run("ValidateAccessForTiers", func(t *testing.T) {
		tiers := []uint8{1, 2, 3} // Basic, Pro, Enterprise

		for _, tier := range tiers {
			hasAccess, err := nftManager.ValidateAccess(ctx, walletAddress, tier)
			require.NoError(t, err, "Should be able to validate access for tier %d", tier)

			t.Logf("Tier %d access: %v", tier, hasAccess)
		}
	})
}

// testContractInteraction testa interação com contratos reais
func testContractInteraction(t *testing.T, nftManager *web3.NFTAccessManager, cfg *config.Config) {
	ctx := context.Background()

	// Teste 1: Conectar com Base Network
	t.Run("BaseNetworkConnection", func(t *testing.T) {
		client, err := ethclient.Dial(cfg.Web3.BaseRPCURL)
		require.NoError(t, err, "Should connect to Base Network")
		defer client.Close()

		// Obter chain ID
		chainID, err := client.ChainID(ctx)
		require.NoError(t, err, "Should get chain ID")

		expectedChainID := big.NewInt(int64(cfg.Web3.BaseChainID))
		assert.Equal(t, expectedChainID, chainID, "Chain ID should match expected Base Network ID")

		t.Logf("Connected to Base Network - Chain ID: %s", chainID.String())
	})

	// Teste 2: Verificar contrato NFT (se configurado)
	t.Run("NFTContractValidation", func(t *testing.T) {
		contractAddr := cfg.Web3.NFTAccessContractAddress
		if contractAddr == "" {
			t.Skip("NFT contract address not configured")
		}

		client, err := ethclient.Dial(cfg.Web3.BaseRPCURL)
		require.NoError(t, err, "Should connect to Base Network")
		defer client.Close()

		// Verificar se o contrato existe (código não é vazio)
		code, err := client.CodeAt(ctx, common.HexToAddress(contractAddr), nil)
		require.NoError(t, err, "Should get contract code")

		assert.NotEmpty(t, code, "Contract code should not be empty")
		t.Logf("NFT contract validated at address: %s", contractAddr)
	})

	// Teste 3: Obter tiers disponíveis
	t.Run("GetAvailableTiers", func(t *testing.T) {
		tiers, err := nftManager.GetAccessTiers(ctx)
		require.NoError(t, err, "Should get available tiers")
		require.Len(t, tiers, 3, "Should have 3 tiers")

		for _, tier := range tiers {
			assert.Greater(t, tier.TierID, uint8(0), "Tier ID should be positive")
			assert.NotEmpty(t, tier.Name, "Tier name should not be empty")
			assert.NotEmpty(t, tier.Description, "Tier description should not be empty")
			assert.NotNil(t, tier.Price, "Tier price should not be nil")
			assert.Greater(t, tier.MaxSupply, uint64(0), "Max supply should be positive")
			assert.True(t, tier.IsActive, "Tier should be active")

			t.Logf("Tier %d: %s - %s (Price: %s ETH)",
				tier.TierID, tier.Name, tier.Description, tier.Price.String())
		}
	})
}

// testAccessTierValidation testa validação de diferentes tiers
func testAccessTierValidation(t *testing.T, nftManager *web3.NFTAccessManager) {
	ctx := context.Background()

	// Wallet de teste (sem NFT)
	testWallet := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"

	t.Run("NoAccessWallet", func(t *testing.T) {
		// Teste com wallet que não tem NFT
		hasAccess, err := nftManager.ValidateAccess(ctx, testWallet, 1)
		require.NoError(t, err, "Should validate access without error")
		assert.False(t, hasAccess, "Wallet without NFT should not have access")
	})

	t.Run("TierRequirements", func(t *testing.T) {
		// Teste diferentes requisitos de tier
		tiers := []uint8{1, 2, 3}

		for _, tier := range tiers {
			hasAccess, err := nftManager.ValidateAccess(ctx, testWallet, tier)
			require.NoError(t, err, "Should validate tier %d access", tier)
			assert.False(t, hasAccess, "Test wallet should not have tier %d access", tier)
		}
	})

	t.Run("GasEstimation", func(t *testing.T) {
		// Teste estimativa de gas para mint
		gasEstimate, err := nftManager.EstimateMintGas(ctx, testWallet, 1)
		require.NoError(t, err, "Should estimate gas for mint")
		assert.Greater(t, gasEstimate, uint64(0), "Gas estimate should be positive")

		t.Logf("Gas estimate for minting tier 1: %d", gasEstimate)
	})
}

// testNationPassIntegrationFlow testa o fluxo completo de integração Nation Pass
func testNationPassIntegrationFlow(t *testing.T, nftManager *web3.NFTAccessManager, cfg *config.Config) {
	ctx := context.Background()
	walletAddress := os.Getenv("WALLET_ADDRESS")

	if walletAddress == "" {
		t.Skip("WALLET_ADDRESS not configured for integration flow test")
	}

	t.Run("CompleteAccessFlow", func(t *testing.T) {
		// 1. Verificar acesso inicial
		t.Log("Step 1: Checking initial access...")
		nft, err := nftManager.CheckAccess(ctx, walletAddress)

		if err != nil {
			t.Logf("Initial access check failed (expected if no NFT): %v", err)
		} else {
			t.Logf("Wallet has NFT access: Tier %s", nft.Tier.Name)
		}

		// 2. Listar NFTs
		t.Log("Step 2: Listing NFTs...")
		nfts, err := nftManager.ListAccessNFTs(ctx, walletAddress)
		require.NoError(t, err, "Should list NFTs")
		t.Logf("Found %d NFTs", len(nfts))

		// 3. Validar acesso para cada tier
		t.Log("Step 3: Validating access for all tiers...")
		for tier := uint8(1); tier <= 3; tier++ {
			hasAccess, err := nftManager.ValidateAccess(ctx, walletAddress, tier)
			require.NoError(t, err, "Should validate tier %d access", tier)
			t.Logf("Tier %d access: %v", tier, hasAccess)
		}

		// 4. Obter informações do contrato
		t.Log("Step 4: Getting contract information...")
		contractAddr := nftManager.GetContractAddress()
		t.Logf("NFT Contract Address: %s", contractAddr)

		// 5. Testar estimativa de gas
		t.Log("Step 5: Testing gas estimation...")
		for tier := uint8(1); tier <= 3; tier++ {
			gasEstimate, err := nftManager.EstimateMintGas(ctx, walletAddress, tier)
			require.NoError(t, err, "Should estimate gas for tier %d", tier)
			t.Logf("Tier %d mint gas estimate: %d", tier, gasEstimate)
		}
	})

	t.Run("NationPassSpecificValidation", func(t *testing.T) {
		// Validação específica para Nation Pass
		nationContract := os.Getenv("NATION_NFT_CONTRACT")
		if nationContract == "" {
			t.Skip("NATION_NFT_CONTRACT not configured")
		}

		t.Logf("Testing Nation Pass contract: %s", nationContract)

		// Verificar se o contrato está configurado corretamente
		assert.Equal(t, nationContract, cfg.Web3.NFTAccessContractAddress,
			"Nation contract should match config")

		// Testar acesso com Nation Pass
		hasAccess, err := nftManager.ValidateAccess(ctx, walletAddress, 1)
		require.NoError(t, err, "Should validate Nation Pass access")

		if hasAccess {
			t.Log("✅ Wallet has Nation Pass access!")
		} else {
			t.Log("❌ Wallet does not have Nation Pass access")
		}
	})
}

// setupIntegrationConfig configura o ambiente para testes de integração
func setupIntegrationConfig(t *testing.T) *config.Config {
	// Configuração mínima para testes de integração
	cfg := &config.Config{
		Web3: config.Web3Config{
			BaseRPCURL:               config.GetDefaultBaseRPC(),
			BaseChainID:              8453, // Base Mainnet
			NFTAccessContractAddress: config.GetDefaultNationPassContract(),
			WalletAddress:            config.GetDefaultWalletAddress(),
			WalletToken:              getEnvOrDefault("WALLET_TOKEN", ""),
		},
		LLM: config.LLMConfig{
			Provider: "nation.fun",
			Model:    "nation-1",
		},
	}

	// Validar configurações descobertas automaticamente
	t.Logf("Configuração automática descoberta:")
	t.Logf("  Base RPC: %s", cfg.Web3.BaseRPCURL)
	t.Logf("  Nation Contract: %s", cfg.Web3.NFTAccessContractAddress)
	t.Logf("  Wallet Address: %s", cfg.Web3.WalletAddress)

	return cfg
}

// getEnvOrDefault obtém variável de ambiente ou retorna valor padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestNFTAccessManagerRealContract testa interação com contrato real da Nation.fun
func TestNFTAccessManagerRealContract(t *testing.T) {
	if os.Getenv("REAL_CONTRACT_TESTS") != "true" {
		t.Skip("Skipping real contract tests. Set REAL_CONTRACT_TESTS=true to run.")
	}

	cfg := setupIntegrationConfig(t)
	_ = logger.New("debug", "json") // log não usado neste teste

	// Conectar com Base Network
	client, err := ethclient.Dial(cfg.Web3.BaseRPCURL)
	require.NoError(t, err, "Failed to connect to Base Network")
	defer client.Close()

	ctx := context.Background()

	t.Run("NationPassContractValidation", func(t *testing.T) {
		contractAddr := cfg.Web3.NFTAccessContractAddress
		if contractAddr == "" {
			t.Skip("NATION_NFT_CONTRACT not configured")
		}

		// Verificar se é um endereço válido
		addr := common.HexToAddress(contractAddr)
		assert.NotEqual(t, common.Address{}, addr, "Contract address should be valid")

		// Verificar se o contrato tem código
		code, err := client.CodeAt(ctx, addr, nil)
		require.NoError(t, err, "Should get contract code")
		assert.NotEmpty(t, code, "Contract should have code")

		t.Logf("✅ Nation Pass contract validated at: %s", contractAddr)

		// Tentar obter informações básicas do contrato
		// Nota: Isso requer ABI do contrato, que não temos ainda
		// Em uma implementação completa, você faria chamadas específicas do contrato
	})

	t.Run("WalletBalanceCheck", func(t *testing.T) {
		walletAddr := cfg.Web3.WalletAddress
		if walletAddr == "" {
			t.Skip("WALLET_ADDRESS not configured")
		}

		// Verificar saldo ETH da wallet
		addr := common.HexToAddress(walletAddr)
		balance, err := client.BalanceAt(ctx, addr, nil)
		require.NoError(t, err, "Should get wallet balance")

		t.Logf("Wallet %s ETH balance: %s ETH", walletAddr, balance.String())

		// Verificar se tem saldo suficiente para gas
		minBalance := big.NewInt(1000000000000000) // 0.001 ETH
		if balance.Cmp(minBalance) < 0 {
			t.Logf("⚠️ Wallet has low ETH balance: %s ETH (minimum recommended: 0.001 ETH)",
				balance.String())
		} else {
			t.Logf("✅ Wallet has sufficient ETH balance")
		}
	})
}

// TestNationPassAccessValidation testa validação específica de acesso Nation Pass
func TestNationPassAccessValidation(t *testing.T) {
	if os.Getenv("NATION_PASS_TESTS") != "true" {
		t.Skip("Skipping Nation Pass tests. Set NATION_PASS_TESTS=true to run.")
	}

	cfg := setupIntegrationConfig(t)
	log := logger.New("debug", "json")

	baseClient, err := web3.NewBaseClient(cfg, log)
	require.NoError(t, err, "Failed to create Base client")

	nftManager := web3.NewNFTAccessManager(cfg, log, baseClient)
	ctx := context.Background()

	walletAddress := os.Getenv("WALLET_ADDRESS")
	if walletAddress == "" {
		t.Fatal("WALLET_ADDRESS is required for Nation Pass tests")
	}

	t.Run("NationPassOwnership", func(t *testing.T) {
		// Verificar se a wallet possui Nation Pass NFT
		nft, err := nftManager.CheckAccess(ctx, walletAddress)

		if err != nil {
			t.Logf("❌ Wallet does not have Nation Pass NFT: %v", err)

			// Se não tem NFT, mostrar como obter
			t.Log("To get Nation Pass access:")
			t.Log("1. Visit https://nation.fun/")
			t.Log("2. Connect your wallet")
			t.Log("3. Purchase a Nation Pass NFT")
			t.Log("4. Update WALLET_ADDRESS in your .env file")

			return
		}

		t.Logf("✅ Wallet has Nation Pass NFT!")
		t.Logf("   Token ID: %s", nft.TokenID.String())
		t.Logf("   Tier: %s", nft.Tier.Name)
		t.Logf("   Owner: %s", nft.Owner)
		t.Logf("   Active: %v", nft.IsActive)

		// Validar que o NFT está ativo
		assert.True(t, nft.IsActive, "Nation Pass NFT should be active")
		assert.Equal(t, walletAddress, nft.Owner, "NFT owner should match wallet")
	})

	t.Run("NationPassTierAccess", func(t *testing.T) {
		// Testar acesso para diferentes tiers
		tiers := map[uint8]string{
			1: "Basic Access",
			2: "Pro Access",
			3: "Enterprise Access",
		}

		for tierID, tierName := range tiers {
			hasAccess, err := nftManager.ValidateAccess(ctx, walletAddress, tierID)
			require.NoError(t, err, "Should validate %s access", tierName)

			if hasAccess {
				t.Logf("✅ Has %s", tierName)
			} else {
				t.Logf("❌ No %s", tierName)
			}
		}
	})

	t.Run("NationPassIntegration", func(t *testing.T) {
		// Teste de integração completa
		t.Log("Testing Nation Pass integration...")

		// 1. Verificar configuração
		contractAddr := nftManager.GetContractAddress()
		assert.NotEmpty(t, contractAddr, "Contract address should be configured")
		t.Logf("Nation Pass contract: %s", contractAddr)

		// 2. Listar todos os NFTs da wallet
		nfts, err := nftManager.ListAccessNFTs(ctx, walletAddress)
		require.NoError(t, err, "Should list NFTs")

		t.Logf("Found %d NFTs in wallet", len(nfts))
		for i, nft := range nfts {
			t.Logf("NFT %d: %s (Tier %d)", i+1, nft.Tier.Name, nft.Tier.TierID)
		}

		// 3. Testar funcionalidades disponíveis
		if len(nfts) > 0 {
			t.Log("✅ Nation Pass integration working!")
			t.Log("Available features:")
			for _, nft := range nfts {
				for _, benefit := range nft.Tier.Benefits {
					t.Logf("  - %s", benefit)
				}
			}
		} else {
			t.Log("❌ No Nation Pass NFTs found")
			t.Log("Please ensure your wallet has a Nation Pass NFT")
		}
	})
}
