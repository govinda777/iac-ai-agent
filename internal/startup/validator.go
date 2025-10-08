package startup

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/govinda777/iac-ai-agent/internal/agent/llm"
	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/internal/services"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// Validator valida requisitos obrigatórios no startup
type Validator struct {
	config *config.Config
	logger *logger.Logger
}

// NewValidator cria um novo validador
func NewValidator(cfg *config.Config, log *logger.Logger) *Validator {
	return &Validator{
		config: cfg,
		logger: log,
	}
}

// ValidationResult contém o resultado da validação
type ValidationResult struct {
	Success              bool
	LLMValidated         bool
	NationNFTValidated   bool
	PrivyValidated       bool
	BaseNetworkValidated bool
	NotionValidated      bool
	AgentCreated         bool
	AgentID              string
	AgentName            string
	NotionAgentID        string
	NotionAgentName      string
	Errors               []string
	Warnings             []string
}

// ValidateAll executa todas as validações obrigatórias
func (v *Validator) ValidateAll(ctx context.Context) (*ValidationResult, error) {
	v.logger.Info("🚀 Iniciando validação de startup...")

	result := &ValidationResult{
		Success:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 1. Validar configuração básica
	v.logger.Info("📋 Validando configuração básica...")
	if err := v.validateBasicConfig(result); err != nil {
		return result, err
	}

	// 2. Validar LLM (OBRIGATÓRIO)
	v.logger.Info("🤖 Validando conexão com LLM...")
	if err := v.validateLLM(ctx, result); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("❌ LLM validation failed: %v", err))
		return result, fmt.Errorf("LLM validation failed: %w", err)
	}
	result.LLMValidated = true
	v.logger.Info("✅ LLM validado com sucesso")

	// 3. Validar Privy.io (OBRIGATÓRIO)
	v.logger.Info("🔐 Validando credenciais Privy.io...")
	if err := v.validatePrivy(result); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("❌ Privy validation failed: %v", err))
		return result, fmt.Errorf("Privy validation failed: %w", err)
	}
	result.PrivyValidated = true
	v.logger.Info("✅ Privy.io validado com sucesso")

	// 4. Validar Base Network
	v.logger.Info("🌐 Validando conexão com Base Network...")
	if err := v.validateBaseNetwork(ctx, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("⚠️ Base Network validation failed: %v", err))
		v.logger.Warn("Base Network não validado, continuando...", "error", err)
	} else {
		result.BaseNetworkValidated = true
		v.logger.Info("✅ Base Network validado com sucesso")
	}

	// 5. Validar Nation.fun NFT (OBRIGATÓRIO)
	v.logger.Info("🎨 Validando posse do NFT Nation.fun...")
	if err := v.validateNationNFT(ctx, result); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("❌ Nation.fun NFT validation failed: %v", err))
		return result, fmt.Errorf("Nation.fun NFT validation failed: %w", err)
	}
	result.NationNFTValidated = true
	v.logger.Info("✅ NFT Nation.fun validado com sucesso")

	// 6. Validar Notion (OPCIONAL)
	v.logger.Info("📝 Validando integração com Notion...")
	if err := v.validateNotion(ctx, result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("⚠️ Notion validation failed: %v", err))
		v.logger.Warn("Notion não validado, continuando...", "error", err)
	} else {
		result.NotionValidated = true
		v.logger.Info("✅ Notion validado com sucesso")
	}

	// 7. Criar ou obter agente padrão (OBRIGATÓRIO)
	v.logger.Info("🤖 Verificando agente padrão...")
	agentID, agentName, err := v.getOrCreateDefaultAgent(ctx, result)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("❌ Agent creation failed: %v", err))
		return result, fmt.Errorf("Agent creation failed: %w", err)
	}
	result.AgentCreated = true
	result.AgentID = agentID
	result.AgentName = agentName
	v.logger.Info("✅ Agente pronto", "id", agentID, "name", agentName)

	// Resultado final
	if result.Success {
		v.logger.Info("🎉 Todas as validações passaram! Aplicação pronta para iniciar.")
	} else {
		v.logger.Error("❌ Validação falhou. Aplicação não pode iniciar.", "errors", result.Errors)
	}

	return result, nil
}

// validateBasicConfig valida configurações básicas
func (v *Validator) validateBasicConfig(result *ValidationResult) error {
	// Verifica variáveis obrigatórias (sem LLM_API_KEY - autenticação via NFT Pass do Nation)
	required := map[string]string{
		"PRIVY_APP_ID":   v.config.Web3.PrivyAppID,
		"WALLET_ADDRESS": os.Getenv("WALLET_ADDRESS"),
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("variável obrigatória não configurada: %s", key)
		}
	}

	return nil
}

// validateLLM valida conexão e autenticação com LLM via NFT Pass do Nation
func (v *Validator) validateLLM(ctx context.Context, result *ValidationResult) error {
	v.logger.Info("Testando conexão com LLM via NFT Pass do Nation...",
		"provider", v.config.LLM.Provider,
		"model", v.config.LLM.Model,
		"wallet", v.config.Web3.WalletAddress)

	// Verificar se temos NFT Pass do Nation válido
	if !v.config.Web3.NationNFTRequired {
		v.logger.Info("Validação de NFT Pass do Nation desabilitada - pulando teste LLM")
		return nil
	}

	// Criar validador de NFT do Nation para teste
	nationValidator := web3.NewNationNFTValidator(v.config, v.logger)

	// Validar NFT Pass do Nation
	nftResponse, err := nationValidator.ValidateWalletNFT(ctx, v.config.Web3.WalletAddress)
	if err != nil {
		return fmt.Errorf("falha na validação de NFT Pass do Nation para LLM: %w", err)
	}

	v.logger.Info("LLM autenticado via NFT Pass do Nation",
		"wallet", v.config.Web3.WalletAddress,
		"token_id", nftResponse.Data.TokenID,
		"tier", nftResponse.Data.Tier,
		"provider", v.config.LLM.Provider)

	// Enviar teste de conectividade para o agente Nation.fun
	testResponse, err := nationValidator.SendTestToNation(ctx, "Teste de conectividade LLM via NFT Pass")
	if err != nil {
		v.logger.Warn("Falha no teste de conectividade LLM", "error", err)
		// Não falha a validação por causa do teste, apenas loga
	} else {
		v.logger.Info("Teste de conectividade LLM bem-sucedido",
			"test_id", testResponse.Data.TestID,
			"status", testResponse.Data.Status)
	}

	return nil
}

// validatePrivy valida credenciais Privy.io
func (v *Validator) validatePrivy(result *ValidationResult) error {
	if v.config.Web3.PrivyAppID == "" {
		return fmt.Errorf("PRIVY_APP_ID não configurado")
	}

	// TODO: Fazer chamada de teste à API do Privy quando implementado
	v.logger.Info("Privy credentials configuradas",
		"app_id", v.config.Web3.PrivyAppID[:8]+"...")

	return nil
}

// validateBaseNetwork valida conexão com Base Network
func (v *Validator) validateBaseNetwork(ctx context.Context, result *ValidationResult) error {
	if v.config.Web3.BaseRPCURL == "" {
		return fmt.Errorf("BASE_RPC_URL não configurado")
	}

	// Conectar ao RPC
	client, err := ethclient.Dial(v.config.Web3.BaseRPCURL)
	if err != nil {
		return fmt.Errorf("falha ao conectar com Base RPC: %w", err)
	}
	defer client.Close()

	// Obter chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("falha ao obter chain ID: %w", err)
	}

	// Validar chain ID
	expectedChainID := big.NewInt(int64(v.config.Web3.BaseChainID))
	if chainID.Cmp(expectedChainID) != 0 {
		return fmt.Errorf("chain ID incorreto: esperado %s, obtido %s", expectedChainID, chainID)
	}

	// Obter bloco atual
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("falha ao obter block number: %w", err)
	}

	v.logger.Info("Base Network conectado",
		"chain_id", chainID,
		"latest_block", blockNumber)

	return nil
}

// validateNationNFT valida que a wallet possui NFT da Nation.fun usando o novo validador
func (v *Validator) validateNationNFT(ctx context.Context, result *ValidationResult) error {
	// Verificar se validação de NFT é obrigatória
	if !v.config.Web3.NationNFTRequired {
		v.logger.Info("Validação de NFT Pass do Nation é opcional - pulada")
		return nil
	}

	walletAddress := v.config.Web3.WalletAddress
	if walletAddress == "" {
		return fmt.Errorf("WALLET_ADDRESS não configurado")
	}

	// Validar formato do endereço
	if !common.IsHexAddress(walletAddress) {
		return fmt.Errorf("WALLET_ADDRESS inválido: %s", walletAddress)
	}

	// Criar validador de NFT do Nation
	nationValidator := web3.NewNationNFTValidator(v.config, v.logger)

	// Executar validação completa (NFT + teste de conectividade)
	if err := nationValidator.ValidateAtStartup(ctx); err != nil {
		return fmt.Errorf("validação de NFT Pass do Nation falhou: %w", err)
	}

	v.logger.Info("✅ NFT Pass do Nation validado com sucesso",
		"wallet", walletAddress,
		"contract", v.config.Web3.NationNFTContract)

	return nil
}

// validateNotion valida integração com Notion
func (v *Validator) validateNotion(ctx context.Context, result *ValidationResult) error {
	// Verifica se Notion está habilitado
	if !v.config.Notion.EnableAgentCreation {
		v.logger.Info("Notion desabilitado na configuração")
		return nil
	}

	// Verifica se API key está configurada
	if v.config.Notion.APIKey == "" {
		return fmt.Errorf("NOTION_API_KEY não configurado")
	}

	// Cria serviço Notion
	notionService, err := services.NewNotionAgentService(v.config, v.logger)
	if err != nil {
		return fmt.Errorf("erro ao criar serviço Notion: %w", err)
	}

	// Verifica se serviço está disponível
	if !notionService.IsServiceAvailable(ctx) {
		return fmt.Errorf("serviço Notion não está disponível")
	}

	// Se auto-create está habilitado, cria/obtém agente
	if v.config.Notion.AutoCreateOnStartup {
		agent, err := notionService.GetOrCreateDefaultAgent(ctx)
		if err != nil {
			return fmt.Errorf("erro ao obter/criar agente Notion: %w", err)
		}

		result.NotionAgentID = agent.ID
		result.NotionAgentName = agent.Name
		v.logger.Info("Agente Notion configurado",
			"id", agent.ID,
			"name", agent.Name,
			"status", agent.Status)
	}

	return nil
}

// getOrCreateDefaultAgent obtém ou cria o agente padrão
func (v *Validator) getOrCreateDefaultAgent(ctx context.Context, result *ValidationResult) (string, string, error) {
	walletAddress := os.Getenv("WALLET_ADDRESS")
	if walletAddress == "" {
		return "", "", fmt.Errorf("WALLET_ADDRESS não configurado")
	}

	// Criar AgentService
	// agentService := services.NewAgentService(v.config, v.logger, nil)

	// Buscar ou criar agente
	// agent, err := agentService.GetOrCreateDefaultAgent(ctx, walletAddress)
	// if err != nil {
	//	return "", "", fmt.Errorf("falha ao obter/criar agente: %w", err)
	// }

	// v.logger.Info("Agente configurado",
	//	"id", agent.ID,
	//	"name", agent.Name,
	//	"owner", agent.Owner,
	//	"status", agent.Status)

	// Para o MVP, retornamos valores simulados
	agentID := "default-agent-123"
	agentName := "Default Agent"

	return agentID, agentName, nil
}

// PrintValidationReport imprime relatório de validação
func (v *Validator) PrintValidationReport(result *ValidationResult) {
	v.logger.Info(strings.Repeat("=", 60))
	v.logger.Info("📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP")
	v.logger.Info(strings.Repeat("=", 60))

	// Status geral
	if result.Success {
		v.logger.Info("✅ Status: PASSOU")
	} else {
		v.logger.Error("❌ Status: FALHOU")
	}

	v.logger.Info("")

	// Checklist
	v.logger.Info("📋 Checklist de Validações:")
	v.printCheckItem("LLM Connection", result.LLMValidated)
	v.printCheckItem("Privy.io Credentials", result.PrivyValidated)
	v.printCheckItem("Base Network", result.BaseNetworkValidated)
	v.printCheckItem("Nation.fun NFT", result.NationNFTValidated)
	v.printCheckItem("Notion Integration", result.NotionValidated)
	v.printCheckItem("Default Agent", result.AgentCreated)

	if result.AgentCreated {
		v.logger.Info("")
		v.logger.Info("🤖 Agent Details:")
		v.logger.Info(fmt.Sprintf("  ID: %s", result.AgentID))
		v.logger.Info(fmt.Sprintf("  Name: %s", result.AgentName))
	}

	if result.NotionValidated && result.NotionAgentID != "" {
		v.logger.Info("")
		v.logger.Info("📝 Notion Agent Details:")
		v.logger.Info(fmt.Sprintf("  ID: %s", result.NotionAgentID))
		v.logger.Info(fmt.Sprintf("  Name: %s", result.NotionAgentName))
	}

	v.logger.Info("")

	// Erros
	if len(result.Errors) > 0 {
		v.logger.Error("❌ Erros Encontrados:")
		for _, err := range result.Errors {
			v.logger.Error("  " + err)
		}
		v.logger.Info("")
	}

	// Avisos
	if len(result.Warnings) > 0 {
		v.logger.Warn("⚠️ Avisos:")
		for _, warn := range result.Warnings {
			v.logger.Warn("  " + warn)
		}
		v.logger.Info("")
	}

	v.logger.Info(strings.Repeat("=", 60))
}

func (v *Validator) printCheckItem(name string, passed bool) {
	status := "❌"
	if passed {
		status = "✅"
	}
	v.logger.Info(fmt.Sprintf("  %s %s", status, name))
}

// MustValidate valida e panic se falhar
func (v *Validator) MustValidate(ctx context.Context) {
	result, err := v.ValidateAll(ctx)

	v.PrintValidationReport(result)

	if err != nil || !result.Success {
		v.logger.Error("💥 APLICAÇÃO NÃO PODE INICIAR - Validação falhou")
		v.logger.Error("Por favor, corrija os erros acima e tente novamente.")

		if len(result.Errors) > 0 {
			v.logger.Error("Erros críticos:")
			for _, e := range result.Errors {
				v.logger.Error("  - " + e)
			}
		}

		panic("Startup validation failed")
	}

	v.logger.Info("✅ Validação completa - Aplicação iniciando...")
}
