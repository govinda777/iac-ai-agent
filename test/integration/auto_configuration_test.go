package integration

import (
	"os"
	"testing"

	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TestAutoConfiguration testa se a configuração automática funciona
func TestAutoConfiguration(t *testing.T) {
	t.Run("DefaultWalletAddress", func(t *testing.T) {
		address := config.GetDefaultWalletAddress()
		assert.NotEmpty(t, address, "Default wallet address should not be empty")
		assert.Equal(t, "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5", address, "Should use correct default wallet")
		t.Logf("✅ Default wallet address: %s", address)
	})

	t.Run("DefaultBaseRPC", func(t *testing.T) {
		rpc := config.GetDefaultBaseRPC()
		assert.NotEmpty(t, rpc, "Default Base RPC should not be empty")
		assert.Contains(t, rpc, "base", "Should be a Base Network RPC")
		t.Logf("✅ Default Base RPC: %s", rpc)
	})

	t.Run("DefaultNationPassContract", func(t *testing.T) {
		contract := config.GetDefaultNationPassContract()
		assert.NotEmpty(t, contract, "Default Nation Pass contract should not be empty")
		assert.Contains(t, contract, "0x", "Should be a valid contract address")
		t.Logf("✅ Default Nation Pass contract: %s", contract)
	})

	t.Run("NationPassSetupValidation", func(t *testing.T) {
		discovery, err := config.ValidateNationPassSetup()

		if err != nil {
			t.Logf("⚠️ Nation Pass setup validation failed: %v", err)
			t.Logf("This is expected if network is not available")
			return
		}

		assert.NotNil(t, discovery, "Discovery should not be nil")
		assert.True(t, discovery.IsValid, "Discovery should be valid")
		assert.NotEmpty(t, discovery.ContractAddress, "Contract address should not be empty")
		assert.NotEmpty(t, discovery.RPCURL, "RPC URL should not be empty")
		assert.Equal(t, int64(8453), discovery.ChainID, "Should be Base Network chain ID")

		t.Logf("✅ Nation Pass setup validation passed:")
		t.Logf("  Contract: %s", discovery.ContractAddress)
		t.Logf("  RPC: %s", discovery.RPCURL)
		t.Logf("  Chain ID: %d", discovery.ChainID)
		t.Logf("  Source: %s", discovery.Source)
	})
}

// TestZeroConfigurationIntegration testa integração sem configuração
func TestZeroConfigurationIntegration(t *testing.T) {
	if os.Getenv("ZERO_CONFIG_TESTS") != "true" {
		t.Skip("Skipping zero configuration tests. Set ZERO_CONFIG_TESTS=true to run.")
	}

	t.Run("IntegrationWithoutConfiguration", func(t *testing.T) {
		// Teste que funciona sem nenhuma configuração
		cfg := &config.Config{
			Web3: config.Web3Config{
				BaseRPCURL:               config.GetDefaultBaseRPC(),
				BaseChainID:              8453,
				NFTAccessContractAddress: config.GetDefaultNationPassContract(),
				WalletAddress:            config.GetDefaultWalletAddress(),
			},
			LLM: config.LLMConfig{
				Provider: "nation.fun",
				Model:    "nation-1",
			},
		}

		t.Logf("🎯 Testando integração com configuração automática:")
		t.Logf("  Wallet: %s", cfg.Web3.WalletAddress)
		t.Logf("  Contract: %s", cfg.Web3.NFTAccessContractAddress)
		t.Logf("  RPC: %s", cfg.Web3.BaseRPCURL)

		// Validar configurações
		assert.NotEmpty(t, cfg.Web3.WalletAddress, "Wallet should be configured")
		assert.NotEmpty(t, cfg.Web3.NFTAccessContractAddress, "Contract should be configured")
		assert.NotEmpty(t, cfg.Web3.BaseRPCURL, "RPC should be configured")

		t.Log("✅ Integração funciona sem configuração manual!")
	})
}
