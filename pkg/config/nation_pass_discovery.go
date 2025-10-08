package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NationPassDiscovery descobre automaticamente informações da Nation Pass
type NationPassDiscovery struct {
	ContractAddress string
	RPCURL          string
	ChainID         int64
	IsValid         bool
	Source          string
}

// discoverNationPassContract tenta descobrir automaticamente o contrato Nation Pass
func discoverNationPassContract() string {
	discovery := &NationPassDiscovery{}

	// Tentar diferentes métodos de descoberta
	if contract := discoverFromAPI(discovery); contract != "" {
		return contract
	}

	if contract := discoverFromKnownContracts(discovery); contract != "" {
		return contract
	}

	// Fallback para contrato conhecido
	return "0x147e832418Cc06A501047019E956714271098b89"
}

// discoverFromAPI tenta descobrir o contrato via API da Nation.fun
func discoverFromAPI(discovery *NationPassDiscovery) string {
	// URLs conhecidas da API Nation.fun
	apiURLs := []string{
		"https://api.nation.fun/v1/contracts",
		"https://nation.fun/api/contracts",
		"https://api.nation.fun/contracts/base",
	}

	for _, url := range apiURLs {
		contract := tryDiscoverFromURL(url)
		if contract != "" {
			discovery.Source = "API"
			return contract
		}
	}

	return ""
}

// tryDiscoverFromURL tenta descobrir contrato de uma URL específica
func tryDiscoverFromURL(url string) string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	// Tentar parsear diferentes formatos de resposta
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}

	// Procurar por endereços de contrato na resposta
	if contracts, ok := data["contracts"].([]interface{}); ok {
		for _, contract := range contracts {
			if contractMap, ok := contract.(map[string]interface{}); ok {
				if address, ok := contractMap["address"].(string); ok {
					if common.IsHexAddress(address) {
						return address
					}
				}
			}
		}
	}

	// Procurar por campo direto
	if address, ok := data["address"].(string); ok {
		if common.IsHexAddress(address) {
			return address
		}
	}

	return ""
}

// discoverFromKnownContracts tenta descobrir usando contratos conhecidos
func discoverFromKnownContracts(discovery *NationPassDiscovery) string {
	// Lista de contratos conhecidos da Nation.fun na Base Network
	knownContracts := []string{
		"0x147e832418Cc06A501047019E956714271098b89", // Contrato principal conhecido
		"0x1234567890123456789012345678901234567890", // Contrato alternativo
	}

	// Tentar validar cada contrato
	for _, contract := range knownContracts {
		if validateContract(contract) {
			discovery.Source = "Known Contracts"
			return contract
		}
	}

	return ""
}

// validateContract valida se um contrato existe e é válido
func validateContract(address string) bool {
	if !common.IsHexAddress(address) {
		return false
	}

	// Tentar conectar com Base Network
	rpcURLs := []string{
		"https://mainnet.base.org",
		"https://base-mainnet.g.alchemy.com/v2/demo",
		"https://base-mainnet.public.blastapi.io",
	}

	for _, rpcURL := range rpcURLs {
		if validateContractWithRPC(address, rpcURL) {
			return true
		}
	}

	return false
}

// validateContractWithRPC valida contrato usando um RPC específico
func validateContractWithRPC(address, rpcURL string) bool {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verificar se o contrato tem código
	code, err := client.CodeAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return false
	}

	// Contrato válido se tem código
	return len(code) > 0
}

// discoverBaseNetworkRPC descobre automaticamente RPCs da Base Network
func discoverBaseNetworkRPC() string {
	// Lista de RPCs públicos conhecidos da Base Network
	rpcURLs := []string{
		"https://mainnet.base.org",
		"https://base-mainnet.g.alchemy.com/v2/demo",
		"https://base-mainnet.public.blastapi.io",
		"https://base.blockpi.network/v1/rpc/public",
		"https://base.meowrpc.com",
	}

	// Testar cada RPC
	for _, rpcURL := range rpcURLs {
		if testRPCConnection(rpcURL) {
			return rpcURL
		}
	}

	// Fallback para RPC principal
	return "https://mainnet.base.org"
}

// testRPCConnection testa se um RPC está funcionando
func testRPCConnection(rpcURL string) bool {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Tentar obter chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return false
	}

	// Base Network tem chain ID 8453
	return chainID.Int64() == 8453
}

// GetDefaultWalletAddress retorna o endereço padrão da wallet
func GetDefaultWalletAddress() string {
	if addr := os.Getenv("WALLET_ADDRESS"); addr != "" {
		return addr
	}
	return "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"
}

// GetDefaultBaseRPC retorna o RPC padrão da Base Network
func GetDefaultBaseRPC() string {
	if rpc := os.Getenv("BASE_RPC_URL"); rpc != "" {
		return rpc
	}
	return discoverBaseNetworkRPC()
}

// GetDefaultNationPassContract retorna o contrato padrão da Nation Pass
func GetDefaultNationPassContract() string {
	if contract := os.Getenv("NFT_CONTRACT_ADDRESS"); contract != "" {
		return contract
	}
	return discoverNationPassContract()
}

// ValidateNationPassSetup valida se a configuração Nation Pass está correta
func ValidateNationPassSetup() (*NationPassDiscovery, error) {
	discovery := &NationPassDiscovery{
		ContractAddress: GetDefaultNationPassContract(),
		RPCURL:          GetDefaultBaseRPC(),
		ChainID:         8453,
	}

	// Validar RPC
	if !testRPCConnection(discovery.RPCURL) {
		return discovery, fmt.Errorf("RPC da Base Network não está funcionando: %s", discovery.RPCURL)
	}

	// Validar contrato
	if !validateContract(discovery.ContractAddress) {
		return discovery, fmt.Errorf("Contrato Nation Pass não é válido: %s", discovery.ContractAddress)
	}

	discovery.IsValid = true
	return discovery, nil
}
