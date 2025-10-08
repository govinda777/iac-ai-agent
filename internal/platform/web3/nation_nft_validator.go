package web3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// NationNFTValidator valida NFT Pass do Nation em tempo de execução
type NationNFTValidator struct {
	config     *config.Config
	logger     *logger.Logger
	httpClient *http.Client
}

// NationNFTResponse representa a resposta da API do Nation.fun
type NationNFTResponse struct {
	Success bool `json:"success"`
	Data    struct {
		HasNFT    bool   `json:"has_nft"`
		TokenID   string `json:"token_id,omitempty"`
		Tier      string `json:"tier,omitempty"`
		IsActive  bool   `json:"is_active,omitempty"`
		ExpiresAt int64  `json:"expires_at,omitempty"`
		Metadata  string `json:"metadata,omitempty"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// NationTestResponse representa a resposta de um teste para o Nation
type NationTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		TestID    string `json:"test_id"`
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
		Response  string `json:"response,omitempty"`
		Error     string `json:"error,omitempty"`
	} `json:"data"`
}

// NewNationNFTValidator cria um novo validador de NFT do Nation
func NewNationNFTValidator(cfg *config.Config, log *logger.Logger) *NationNFTValidator {
	return &NationNFTValidator{
		config: cfg,
		logger: log,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ValidateWalletNFT valida se a carteira possui NFT Pass do Nation
func (nv *NationNFTValidator) ValidateWalletNFT(ctx context.Context, walletAddress string) (*NationNFTResponse, error) {
	nv.logger.Info("Validando NFT Pass do Nation", "wallet", walletAddress)

	// Verificar se é a carteira padrão autorizada
	if !nv.isDefaultWallet(walletAddress) {
		return nil, fmt.Errorf("wallet não autorizada para acesso: %s", walletAddress)
	}

	// Verificar NFT via API do Nation.fun
	nftResponse, err := nv.checkNFTOwnership(ctx, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar NFT: %w", err)
	}

	if !nftResponse.Success || !nftResponse.Data.HasNFT {
		return nil, fmt.Errorf("carteira não possui NFT Pass do Nation válido")
	}

	nv.logger.Info("NFT Pass do Nation validado com sucesso",
		"wallet", walletAddress,
		"token_id", nftResponse.Data.TokenID,
		"tier", nftResponse.Data.Tier)

	return nftResponse, nil
}

// SendTestToNation envia um teste para o Nation.fun e coleta resposta
func (nv *NationNFTValidator) SendTestToNation(ctx context.Context, testMessage string) (*NationTestResponse, error) {
	nv.logger.Info("Enviando teste para Nation.fun", "message", testMessage)

	// Preparar dados do teste
	testData := map[string]interface{}{
		"message":   testMessage,
		"timestamp": time.Now().Unix(),
		"source":    "iac-ai-agent",
		"wallet":    nv.config.Web3.WalletAddress,
	}

	// Enviar teste via API
	response, err := nv.sendTestRequest(ctx, testData)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar teste: %w", err)
	}

	nv.logger.Info("Teste enviado para Nation.fun com sucesso",
		"test_id", response.Data.TestID,
		"status", response.Data.Status)

	return response, nil
}

// ValidateAtStartup valida NFT na inicialização da aplicação
func (nv *NationNFTValidator) ValidateAtStartup(ctx context.Context) error {
	nv.logger.Info("Iniciando validação de NFT Pass do Nation na inicialização")

	walletAddress := nv.config.Web3.WalletAddress
	if walletAddress == "" {
		return fmt.Errorf("WALLET_ADDRESS não configurado")
	}

	// Validar NFT
	nftResponse, err := nv.ValidateWalletNFT(ctx, walletAddress)
	if err != nil {
		nv.logger.Error("Falha na validação de NFT na inicialização", "error", err)
		return fmt.Errorf("validação de NFT falhou na inicialização: %w", err)
	}

	// Enviar teste de conectividade
	testResponse, err := nv.SendTestToNation(ctx, "Teste de conectividade na inicialização")
	if err != nil {
		nv.logger.Warn("Falha no teste de conectividade", "error", err)
		// Não falha a inicialização por causa do teste, apenas loga
	} else {
		nv.logger.Info("Teste de conectividade com Nation.fun bem-sucedido",
			"test_id", testResponse.Data.TestID)
	}

	nv.logger.Info("Validação de NFT Pass do Nation concluída com sucesso",
		"wallet", walletAddress,
		"token_id", nftResponse.Data.TokenID,
		"tier", nftResponse.Data.Tier)

	return nil
}

// isDefaultWallet verifica se é a carteira padrão autorizada
func (nv *NationNFTValidator) isDefaultWallet(walletAddress string) bool {
	defaultWallet := "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"
	return common.HexToAddress(walletAddress) == common.HexToAddress(defaultWallet)
}

// checkNFTOwnership verifica propriedade de NFT via API do Nation.fun
func (nv *NationNFTValidator) checkNFTOwnership(ctx context.Context, walletAddress string) (*NationNFTResponse, error) {
	// URL da API do Nation.fun para verificação de NFT
	apiURL := fmt.Sprintf("https://api.nation.fun/v1/nft/check/%s", walletAddress)

	// Criar requisição
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Adicionar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "iac-ai-agent/1.0")

	// Se temos token de autenticação, adicionar
	if nv.config.Web3.WalletToken != "" {
		req.Header.Set("Authorization", "Bearer "+nv.config.Web3.WalletToken)
	}

	// Fazer requisição
	resp, err := nv.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verificar status HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	// Parsear resposta JSON
	var nftResponse NationNFTResponse
	if err := json.Unmarshal(body, &nftResponse); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta JSON: %w", err)
	}

	return &nftResponse, nil
}

// sendTestRequest envia requisição de teste para o Nation.fun
func (nv *NationNFTValidator) sendTestRequest(ctx context.Context, testData map[string]interface{}) (*NationTestResponse, error) {
	// URL da API do Nation.fun para testes
	apiURL := "https://api.nation.fun/v1/test/send"

	// Converter dados para JSON
	jsonData, err := json.Marshal(testData)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados para JSON: %w", err)
	}

	// Criar requisição POST
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Adicionar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "iac-ai-agent/1.0")

	// Se temos token de autenticação, adicionar
	if nv.config.Web3.WalletToken != "" {
		req.Header.Set("Authorization", "Bearer "+nv.config.Web3.WalletToken)
	}

	// Fazer requisição
	resp, err := nv.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verificar status HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	// Parsear resposta JSON
	var testResponse NationTestResponse
	if err := json.Unmarshal(body, &testResponse); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta JSON: %w", err)
	}

	return &testResponse, nil
}

// GetNFTInfo retorna informações detalhadas do NFT
func (nv *NationNFTValidator) GetNFTInfo(ctx context.Context, walletAddress string) (*NationNFTResponse, error) {
	return nv.ValidateWalletNFT(ctx, walletAddress)
}

// IsNFTAccessRequired verifica se a validação de NFT é obrigatória
func (nv *NationNFTValidator) IsNFTAccessRequired() bool {
	return nv.config.Web3.NationNFTRequired
}

// GetDefaultWalletAddress retorna o endereço da carteira padrão
func (nv *NationNFTValidator) GetDefaultWalletAddress() string {
	return "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"
}
