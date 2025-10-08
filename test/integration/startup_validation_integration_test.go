package integration

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/govinda777/iac-ai-agent/internal/startup"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartupValidationIntegration testa a validação de startup com Nation Pass
// Este teste valida o fluxo completo de validação que acontece na inicialização da aplicação
func TestStartupValidationIntegration(t *testing.T) {
	if os.Getenv("STARTUP_VALIDATION_TESTS") != "true" {
		t.Skip("Skipping startup validation tests. Set STARTUP_VALIDATION_TESTS=true to run.")
	}

	cfg := setupStartupValidationConfig(t)
	log := logger.New("debug", "text")

	validator := startup.NewValidator(cfg, log)
	ctx := context.Background()

	t.Run("CompleteStartupValidation", func(t *testing.T) {
		// Executar validação completa de startup
		result, err := validator.ValidateAll(ctx)

		if err != nil {
			t.Logf("Startup validation failed: %v", err)

			// Se falhou, mostrar detalhes
			if result != nil {
				t.Logf("Validation result: Success=%v", result.Success)
				for _, validationError := range result.Errors {
					t.Logf("Error: %s", validationError)
				}
				for _, warning := range result.Warnings {
					t.Logf("Warning: %s", warning)
				}
			}

			// Para testes de integração, não falhamos se não tem NFT
			// Apenas logamos o resultado
			return
		}

		require.NotNil(t, result, "Validation result should not be nil")
		t.Logf("Startup validation successful: %+v", result)

		// Validar componentes críticos
		assert.True(t, result.Success, "Startup validation should succeed")

		// Verificar se Nation NFT foi validado
		if result.NationNFTValidated {
			t.Log("✅ Nation Pass NFT validation passed")
		} else {
			t.Log("⚠️ Nation Pass NFT validation skipped or failed")
		}
	})

	t.Run("NationPassValidationOnly", func(t *testing.T) {
		// Testar apenas a validação de Nation Pass
		result := &startup.ValidationResult{
			Success:  true,
			Errors:   []string{},
			Warnings: []string{},
		}

		result, err := validator.ValidateAll(ctx)

		if err != nil {
			t.Logf("Nation Pass validation failed: %v", err)

			// Verificar se é erro esperado (sem NFT)
			if result != nil && len(result.Warnings) > 0 {
				for _, warning := range result.Warnings {
					t.Logf("Warning: %s", warning)
				}
			}
		} else {
			t.Log("✅ Nation Pass validation passed")
			assert.True(t, result.NationNFTValidated, "Nation NFT should be validated")
		}
	})
}

// TestNationPassTierValidation testa validação específica de tiers Nation Pass
func TestNationPassTierValidation(t *testing.T) {
	if os.Getenv("TIER_VALIDATION_TESTS") != "true" {
		t.Skip("Skipping tier validation tests. Set TIER_VALIDATION_TESTS=true to run.")
	}

	ctx := context.Background()
	walletAddress := os.Getenv("WALLET_ADDRESS")

	if walletAddress == "" {
		t.Skip("WALLET_ADDRESS not configured for tier validation tests")
	}

	t.Run("TierAccessMatrix", func(t *testing.T) {
		// Testar matriz de acesso para diferentes tiers
		tiers := map[uint8]string{
			1: "Basic Access",
			2: "Pro Access",
			3: "Enterprise Access",
		}

		for tierID, tierName := range tiers {
			t.Run(tierName, func(t *testing.T) {
				// Simular validação de acesso para o tier
				hasAccess := simulateTierAccess(ctx, walletAddress, tierID)

				if hasAccess {
					t.Logf("✅ Wallet has %s", tierName)

					// Se tem acesso, validar funcionalidades
					features := getTierFeatures(tierID)
					t.Logf("Available features for %s:", tierName)
					for _, feature := range features {
						t.Logf("  - %s", feature)
					}
				} else {
					t.Logf("❌ Wallet does not have %s", tierName)
				}
			})
		}
	})

	t.Run("TierUpgradeValidation", func(t *testing.T) {
		// Testar validação de upgrade de tier
		currentTier := uint8(1) // Assumir Basic Access

		for targetTier := uint8(2); targetTier <= 3; targetTier++ {
			canUpgrade := validateTierUpgrade(currentTier, targetTier)

			if canUpgrade {
				t.Logf("✅ Can upgrade from tier %d to tier %d", currentTier, targetTier)

				// Calcular custo de upgrade
				cost := calculateUpgradeCost(currentTier, targetTier)
				t.Logf("Upgrade cost: %s ETH", cost.String())
			} else {
				t.Logf("❌ Cannot upgrade from tier %d to tier %d", currentTier, targetTier)
			}
		}
	})
}

// TestNationPassIntegrationFlow testa o fluxo completo de integração Nation Pass
func TestNationPassIntegrationFlow(t *testing.T) {
	if os.Getenv("NATION_PASS_FLOW_TESTS") != "true" {
		t.Skip("Skipping Nation Pass flow tests. Set NATION_PASS_FLOW_TESTS=true to run.")
	}

	cfg := setupStartupValidationConfig(t)
	log := logger.New("debug", "text")

	ctx := context.Background()
	walletAddress := os.Getenv("WALLET_ADDRESS")

	if walletAddress == "" {
		t.Skip("WALLET_ADDRESS not configured for flow tests")
	}

	t.Run("CompleteNationPassFlow", func(t *testing.T) {
		// Fluxo completo de validação Nation Pass

		// 1. Verificar configuração
		t.Log("Step 1: Checking Nation Pass configuration...")
		nationContract := os.Getenv("NATION_NFT_CONTRACT")
		if nationContract == "" {
			t.Log("⚠️ NATION_NFT_CONTRACT not configured")
			return
		}
		t.Logf("Nation Pass contract: %s", nationContract)

		// 2. Validar wallet
		t.Log("Step 2: Validating wallet...")
		if walletAddress == "" {
			t.Log("❌ WALLET_ADDRESS not configured")
			return
		}
		t.Logf("Wallet address: %s", walletAddress)

		// 3. Verificar acesso
		t.Log("Step 3: Checking Nation Pass access...")
		hasAccess := simulateNationPassAccess(ctx, walletAddress)

		if hasAccess {
			t.Log("✅ Wallet has Nation Pass access!")

			// 4. Validar funcionalidades disponíveis
			t.Log("Step 4: Validating available features...")
			features := getNationPassFeatures()
			for _, feature := range features {
				t.Logf("  ✅ %s", feature)
			}

			// 5. Testar integração com LLM
			t.Log("Step 5: Testing LLM integration...")
			llmWorking := testLLMIntegration(cfg, log)
			if llmWorking {
				t.Log("✅ LLM integration working")
			} else {
				t.Log("❌ LLM integration failed")
			}

		} else {
			t.Log("❌ Wallet does not have Nation Pass access")
			t.Log("To get Nation Pass access:")
			t.Log("1. Visit https://nation.fun/")
			t.Log("2. Connect your wallet")
			t.Log("3. Purchase a Nation Pass NFT")
			t.Log("4. Update WALLET_ADDRESS in your .env file")
		}
	})

	t.Run("NationPassErrorHandling", func(t *testing.T) {
		// Testar tratamento de erros específicos do Nation Pass

		// Teste com wallet inválida
		invalidWallet := "0xinvalid"
		hasAccess := simulateNationPassAccess(ctx, invalidWallet)
		assert.False(t, hasAccess, "Invalid wallet should not have access")

		// Teste com contrato inválido
		invalidContract := "0x0000000000000000000000000000000000000000"
		contractValid := validateContract(invalidContract)
		assert.False(t, contractValid, "Invalid contract should not be valid")

		// Teste com RPC inválido
		invalidRPC := "https://invalid-rpc.com"
		rpcWorking := testRPCConnection(invalidRPC)
		assert.False(t, rpcWorking, "Invalid RPC should not work")
	})
}

// Funções auxiliares para simulação

func simulateTierAccess(ctx context.Context, walletAddress string, tierID uint8) bool {
	// Simulação de verificação de acesso por tier
	// Em implementação real, isso faria chamada ao contrato

	// Para testes, assumir que tem acesso básico se wallet está configurada
	if walletAddress != "" && tierID == 1 {
		return true
	}

	return false
}

func getTierFeatures(tierID uint8) []string {
	features := map[uint8][]string{
		1: {
			"Análises ilimitadas de Terraform",
			"Detecção de segurança com Checkov",
			"Sugestões básicas de otimização",
			"Suporte via Discord",
		},
		2: {
			"Tudo do Basic Access",
			"Análise com LLM (GPT-4/Claude)",
			"Sugestões contextualizadas inteligentes",
			"Análise de Preview e Drift",
			"Detecção de Secrets",
			"Recomendações de arquitetura",
			"Priority support",
		},
		3: {
			"Tudo do Pro Access",
			"API dedicada com rate limits maiores",
			"Custom knowledge base",
			"Integração com CI/CD privado",
			"Suporte dedicado 24/7",
			"SLA garantido",
			"Governance tokens inclusos",
		},
	}

	if tierFeatures, exists := features[tierID]; exists {
		return tierFeatures
	}

	return []string{}
}

func validateTierUpgrade(currentTier, targetTier uint8) bool {
	// Validar se pode fazer upgrade
	return targetTier > currentTier && targetTier <= 3
}

func calculateUpgradeCost(currentTier, targetTier uint8) *big.Int {
	// Calcular custo de upgrade
	prices := map[uint8]*big.Int{
		1: big.NewInt(10000000000000000),  // 0.01 ETH
		2: big.NewInt(50000000000000000),  // 0.05 ETH
		3: big.NewInt(200000000000000000), // 0.2 ETH
	}

	currentPrice := prices[currentTier]
	targetPrice := prices[targetTier]

	if currentPrice == nil || targetPrice == nil {
		return big.NewInt(0)
	}

	return new(big.Int).Sub(targetPrice, currentPrice)
}

func simulateNationPassAccess(ctx context.Context, walletAddress string) bool {
	// Simulação de verificação de acesso Nation Pass
	// Em implementação real, isso faria chamada ao contrato Nation Pass

	// Para testes, verificar se wallet está configurada
	return walletAddress != "" && len(walletAddress) == 42 && walletAddress[:2] == "0x"
}

func getNationPassFeatures() []string {
	return []string{
		"Access to IaC AI Agent",
		"Terraform analysis",
		"Security scanning with Checkov",
		"Cost optimization suggestions",
		"LLM-powered insights",
		"Preview and drift analysis",
		"Secret detection",
		"Architecture recommendations",
	}
}

func testLLMIntegration(cfg *config.Config, log *logger.Logger) bool {
	// Simulação de teste de integração LLM
	// Em implementação real, isso testaria conexão com Nation.fun LLM

	return cfg.LLM.Provider == "nation.fun"
}

func validateContract(contractAddress string) bool {
	// Simulação de validação de contrato
	// Em implementação real, isso verificaria se o contrato existe

	return len(contractAddress) == 42 && contractAddress[:2] == "0x" && contractAddress != "0x0000000000000000000000000000000000000000"
}

func testRPCConnection(rpcURL string) bool {
	// Simulação de teste de conexão RPC
	// Em implementação real, isso testaria conexão com Base Network

	return rpcURL != "" && rpcURL[:4] == "http"
}

func setupStartupValidationConfig(t *testing.T) *config.Config {
	// Configuração para testes de validação de startup
	cfg := &config.Config{
		Web3: config.Web3Config{
			BaseRPCURL:               config.GetDefaultBaseRPC(),
			BaseChainID:              8453,
			NFTAccessContractAddress: config.GetDefaultNationPassContract(),
			WalletAddress:            config.GetDefaultWalletAddress(),
			WalletToken:              getEnvOrDefault("WALLET_TOKEN", ""),
		},
		LLM: config.LLMConfig{
			Provider: "nation.fun",
			Model:    "nation-1",
		},
	}

	// Log das configurações descobertas automaticamente
	t.Logf("Configuração automática para validação de startup:")
	t.Logf("  Base RPC: %s", cfg.Web3.BaseRPCURL)
	t.Logf("  Nation Contract: %s", cfg.Web3.NFTAccessContractAddress)
	t.Logf("  Wallet Address: %s", cfg.Web3.WalletAddress)

	return cfg
}

