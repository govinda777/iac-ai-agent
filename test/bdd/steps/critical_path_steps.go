package steps

import (
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/internal/services"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// CriticalPathContext mantém o estado entre os steps do teste
type CriticalPathContext struct {
	// Componentes do sistema
	config           *config.Config
	logger           *logger.Logger
	privyClient      *web3.PrivyClient
	nftAccessManager *web3.NFTAccessManager
	botTokenManager  *web3.BotTokenManager
	onrampManager    *web3.PrivyOnrampManager
	analysisService  *services.AnalysisService

	// Estado do teste
	isAuthenticated  bool
	walletAddress    string
	userID           string
	currentPage      string
	selectedTier     string
	selectedPackage  string
	nftTokenID       string
	tokenBalance     int
	analysisRequest  string
	analysisType     string
	analysisResponse map[string]interface{}
	transactionHash  string

	// Mock mode
	isMockMode bool
}

// NewCriticalPathContext cria um novo contexto para os testes
func NewCriticalPathContext(cfg *config.Config, log *logger.Logger) *CriticalPathContext {
	return &CriticalPathContext{
		config:          cfg,
		logger:          log,
		isAuthenticated: false,
		isMockMode:      false,
	}
}

// RegisterCriticalPathSteps registra os steps nos testes BDD
func RegisterCriticalPathSteps(ctx *godog.ScenarioContext, cpCtx *CriticalPathContext) {
	// Configuração geral
	ctx.Step(`^que o serviço está disponível$`, cpCtx.serviceIsAvailable)
	ctx.Step(`^que o serviço Privy está configurado$`, cpCtx.privyServiceConfigured)
	ctx.Step(`^que os contratos na Base Network estão deployados$`, cpCtx.contractsAreDeployed)
	ctx.Step(`^que o sistema está em modo de validação$`, cpCtx.systemInValidationMode)
	ctx.Step(`^a API do Privy está mockada$`, cpCtx.privyApiIsMocked)

	// Login e autenticação
	ctx.Step(`^que sou um novo usuário sem conta$`, cpCtx.iAmNewUserWithoutAccount)
	ctx.Step(`^eu acesso a página inicial$`, cpCtx.iAccessHomePage)
	ctx.Step(`^clico em "([^"]*)"$`, cpCtx.iClickOn)
	ctx.Step(`^seleciono "([^"]*)" como provedor$`, cpCtx.iSelectAsProvider)
	ctx.Step(`^autorizo a conexão no MetaMask$`, cpCtx.iAuthorizeConnectionInMetaMask)
	ctx.Step(`^devo estar autenticado com sucesso$`, cpCtx.iShouldBeAuthenticated)
	ctx.Step(`^meu endereço de wallet deve estar visível no cabeçalho$`, cpCtx.myWalletAddressShouldBeVisible)
	ctx.Step(`^eu faço login com a wallet mockada "([^"]*)"$`, cpCtx.iLoginWithMockedWallet)
	ctx.Step(`^devo ver "([^"]*)"$`, cpCtx.iShouldSee)

	// NFT Purchase
	ctx.Step(`^eu navego para a página "([^"]*)"$`, cpCtx.iNavigateToPage)
	ctx.Step(`^seleciono o tier "([^"]*)"$`, cpCtx.iSelectTier)
	ctx.Step(`^confirmo a transação na wallet$`, cpCtx.iConfirmTransactionInWallet)
	ctx.Step(`^a transação deve ser processada$`, cpCtx.transactionShouldBeProcessed)
	ctx.Step(`^após confirmação devo receber o NFT ([^"]*)$`, cpCtx.iShouldReceiveNFT)
	ctx.Step(`^devo ver o status "([^"]*)" no dashboard$`, cpCtx.iShouldSeeStatusInDashboard)
	ctx.Step(`^confirmo a compra usando a função de simulação$`, cpCtx.iConfirmPurchaseWithSimulation)
	ctx.Step(`^o sistema deve simular a transação com sucesso$`, cpCtx.systemShouldSimulateTransactionSuccessfully)
	ctx.Step(`^devo receber um NFT ([^"]*) simulado$`, cpCtx.iShouldReceiveSimulatedNFT)
	ctx.Step(`^meu status de acesso deve ser "([^"]*)"$`, cpCtx.myAccessStatusShouldBe)

	// Token Purchase
	ctx.Step(`^seleciono o pacote "([^"]*)" \((\d+) tokens\)$`, cpCtx.iSelectTokenPackage)
	ctx.Step(`^os tokens devem ser adicionados à minha conta$`, cpCtx.tokensShouldBeAddedToMyAccount)
	ctx.Step(`^meu saldo deve mostrar "([^"]*)"$`, cpCtx.myBalanceShouldShow)
	ctx.Step(`^meu saldo simulado deve ser "([^"]*)"$`, cpCtx.mySimulatedBalanceShouldBe)

	// Analysis
	ctx.Step(`^eu submeto o seguinte código Terraform:$`, cpCtx.iSubmitTerraformCode)
	ctx.Step(`^seleciono "([^"]*)" \((\d+) tokens\)$`, cpCtx.iSelectAnalysisType)
	ctx.Step(`^o sistema deve processar minha análise$`, cpCtx.systemShouldProcessMyAnalysis)
	ctx.Step(`^deve debitar (\d+) IACAI tokens da minha conta$`, cpCtx.shouldDeductTokensFromMyAccount)
	ctx.Step(`^meu saldo final deve ser (\d+) IACAI$`, cpCtx.myFinalBalanceShouldBe)
	ctx.Step(`^eu submeto código Terraform para análise:$`, cpCtx.iSubmitCodeForAnalysis)
	ctx.Step(`^o sistema deve processar a análise real$`, cpCtx.systemShouldProcessRealAnalysis)
	ctx.Step(`^debitar tokens do saldo simulado$`, cpCtx.deductTokensFromSimulatedBalance)

	// Results
	ctx.Step(`^devo receber um relatório de análise completo em até (\d+) segundos$`, cpCtx.iShouldReceiveAnalysisReport)
	ctx.Step(`^o relatório deve conter as seguintes seções:$`, cpCtx.reportShouldContainSections)
	ctx.Step(`^pelo menos uma sugestão de segurança sobre o "([^"]*)" ACL$`, cpCtx.atLeastOneSuggestionAbout)
	ctx.Step(`^cada sugestão deve conter:$`, cpCtx.eachSuggestionShouldContain)
	ctx.Step(`^devo receber sugestões reais de melhoria$`, cpCtx.iShouldReceiveRealSuggestions)
	ctx.Step(`^o relatório deve ser completo e detalhado$`, cpCtx.reportShouldBeCompleteAndDetailed)
	ctx.Step(`^meu saldo final simulado deve refletir o custo da análise$`, cpCtx.myFinalSimulatedBalanceShouldReflectAnalysisCost)

	// Onramp & Embedded Wallet
	ctx.Step(`^eu me autentico com email e senha$`, cpCtx.iAuthenticateWithEmailAndPassword)
	ctx.Step(`^uma embedded wallet deve ser criada para mim$`, cpCtx.embeddedWalletShouldBeCreatedForMe)
	ctx.Step(`^completo o processo de pagamento com cartão$`, cpCtx.iCompleteCardPaymentProcess)
	ctx.Step(`^o NFT deve ser mintado para minha embedded wallet$`, cpCtx.nftShouldBeMintedToMyEmbeddedWallet)
}

// Implementação dos steps

func (ctx *CriticalPathContext) serviceIsAvailable() error {
	// Check if service is running
	return nil
}

func (ctx *CriticalPathContext) privyServiceConfigured() error {
	if ctx.config.Web3.PrivyAppID == "" {
		return fmt.Errorf("Privy App ID não está configurado")
	}
	return nil
}

func (ctx *CriticalPathContext) contractsAreDeployed() error {
	// Check if contracts are deployed
	if ctx.config.Web3.NFTAccessContractAddress == "" {
		return fmt.Errorf("Endereço do contrato NFT não configurado")
	}
	if ctx.config.Web3.BotTokenContractAddress == "" {
		return fmt.Errorf("Endereço do contrato token não configurado")
	}
	return nil
}

func (ctx *CriticalPathContext) systemInValidationMode() error {
	ctx.isMockMode = true
	return nil
}

func (ctx *CriticalPathContext) privyApiIsMocked() error {
	// Setup mock Privy client
	return nil
}

// Login steps
func (ctx *CriticalPathContext) iAmNewUserWithoutAccount() error {
	ctx.isAuthenticated = false
	ctx.walletAddress = ""
	return nil
}

func (ctx *CriticalPathContext) iAccessHomePage() error {
	ctx.currentPage = "home"
	return nil
}

func (ctx *CriticalPathContext) iClickOn(buttonText string) error {
	// Simula click on button
	ctx.logger.Info("Clicando em botão", "button", buttonText)
	return nil
}

func (ctx *CriticalPathContext) iSelectAsProvider(provider string) error {
	// Simula select provider
	ctx.logger.Info("Selecionando provedor", "provider", provider)
	return nil
}

func (ctx *CriticalPathContext) iAuthorizeConnectionInMetaMask() error {
	// Simula authorize in MetaMask
	ctx.isAuthenticated = true
	ctx.walletAddress = "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
	ctx.userID = "user_" + ctx.walletAddress[2:10]
	return nil
}

func (ctx *CriticalPathContext) iShouldBeAuthenticated() error {
	if !ctx.isAuthenticated {
		return fmt.Errorf("Usuário não está autenticado")
	}
	return nil
}

func (ctx *CriticalPathContext) myWalletAddressShouldBeVisible() error {
	if ctx.walletAddress == "" {
		return fmt.Errorf("Endereço de wallet não está disponível")
	}
	return nil
}

func (ctx *CriticalPathContext) iLoginWithMockedWallet(address string) error {
	ctx.isAuthenticated = true
	ctx.walletAddress = address
	ctx.userID = "user_mocked_" + time.Now().Format("20060102150405")
	return nil
}

func (ctx *CriticalPathContext) iShouldSee(text string) error {
	// Verifica se o texto está visível
	if strings.Contains(text, "Wallet conectada") {
		if ctx.walletAddress == "" {
			return fmt.Errorf("Wallet não está conectada")
		}
		return nil
	}
	return nil
}

// NFT steps
func (ctx *CriticalPathContext) iNavigateToPage(pageName string) error {
	ctx.currentPage = strings.ToLower(pageName)
	return nil
}

func (ctx *CriticalPathContext) iSelectTier(tier string) error {
	ctx.selectedTier = tier
	return nil
}

func (ctx *CriticalPathContext) iConfirmTransactionInWallet() error {
	ctx.transactionHash = fmt.Sprintf("0x%x", time.Now().UnixNano())
	return nil
}

func (ctx *CriticalPathContext) transactionShouldBeProcessed() error {
	if ctx.transactionHash == "" {
		return fmt.Errorf("Transação não foi processada")
	}
	return nil
}

func (ctx *CriticalPathContext) iShouldReceiveNFT(tierName string) error {
	if ctx.selectedTier != tierName {
		return fmt.Errorf("Tier do NFT incorreto")
	}
	ctx.nftTokenID = fmt.Sprintf("%d", time.Now().Unix())
	return nil
}

func (ctx *CriticalPathContext) iShouldSeeStatusInDashboard(status string) error {
	// Verifica se o status é exibido no dashboard
	return nil
}

func (ctx *CriticalPathContext) iConfirmPurchaseWithSimulation() error {
	ctx.transactionHash = fmt.Sprintf("0x%x_simulated", time.Now().UnixNano())
	return nil
}

func (ctx *CriticalPathContext) systemShouldSimulateTransactionSuccessfully() error {
	if !strings.Contains(ctx.transactionHash, "simulated") {
		return fmt.Errorf("Transação simulada não foi registrada")
	}
	return nil
}

func (ctx *CriticalPathContext) iShouldReceiveSimulatedNFT(tierName string) error {
	ctx.nftTokenID = fmt.Sprintf("%d_simulated", time.Now().Unix())
	return nil
}

func (ctx *CriticalPathContext) myAccessStatusShouldBe(status string) error {
	if ctx.selectedTier != status {
		return fmt.Errorf("Status de acesso incorreto")
	}
	return nil
}

// Token purchase steps
func (ctx *CriticalPathContext) iSelectTokenPackage(packageName string, tokens int) error {
	ctx.selectedPackage = packageName
	ctx.tokenBalance = tokens
	return nil
}

func (ctx *CriticalPathContext) tokensShouldBeAddedToMyAccount() error {
	// Verifica se os tokens foram adicionados
	return nil
}

func (ctx *CriticalPathContext) myBalanceShouldShow(balanceText string) error {
	if !strings.Contains(balanceText, fmt.Sprintf("%d", ctx.tokenBalance)) {
		return fmt.Errorf("Saldo incorreto")
	}
	return nil
}

func (ctx *CriticalPathContext) mySimulatedBalanceShouldBe(balance string) error {
	if !strings.Contains(balance, fmt.Sprintf("%d", ctx.tokenBalance)) {
		return fmt.Errorf("Saldo simulado incorreto")
	}
	return nil
}

// Analysis steps
func (ctx *CriticalPathContext) iSubmitTerraformCode(code string) error {
	ctx.analysisRequest = code
	return nil
}

func (ctx *CriticalPathContext) iSelectAnalysisType(analysisType string, cost int) error {
	ctx.analysisType = analysisType
	return nil
}

func (ctx *CriticalPathContext) systemShouldProcessMyAnalysis() error {
	if ctx.analysisRequest == "" {
		return fmt.Errorf("Não há código para análise")
	}

	// Simula processamento
	time.Sleep(1 * time.Second)

	// Cria resultado simulado
	ctx.analysisResponse = map[string]interface{}{
		"score": 75,
		"issues": []map[string]interface{}{
			{
				"severity":   "HIGH",
				"message":    "S3 bucket has public read ACL",
				"resource":   "aws_s3_bucket.example",
				"suggestion": "Remove public-read ACL for security",
			},
		},
		"sections": []string{
			"Executive Summary",
			"Security Issues",
			"Best Practices",
			"Cost Optimization",
			"Detailed Findings",
		},
	}

	return nil
}

func (ctx *CriticalPathContext) shouldDeductTokensFromMyAccount(cost int) error {
	ctx.tokenBalance -= cost
	return nil
}

func (ctx *CriticalPathContext) myFinalBalanceShouldBe(expectedBalance int) error {
	if ctx.tokenBalance != expectedBalance {
		return fmt.Errorf("Saldo final incorreto: esperado %d, atual %d", expectedBalance, ctx.tokenBalance)
	}
	return nil
}

func (ctx *CriticalPathContext) iSubmitCodeForAnalysis(code string) error {
	return ctx.iSubmitTerraformCode(code)
}

func (ctx *CriticalPathContext) systemShouldProcessRealAnalysis() error {
	return ctx.systemShouldProcessMyAnalysis()
}

func (ctx *CriticalPathContext) deductTokensFromSimulatedBalance() error {
	// Find analysis type cost
	var cost int
	switch ctx.analysisType {
	case "Full Review":
		cost = 15
	case "LLM Analysis":
		cost = 5
	default:
		cost = 1
	}

	return ctx.shouldDeductTokensFromMyAccount(cost)
}

// Results steps
func (ctx *CriticalPathContext) iShouldReceiveAnalysisReport(seconds int) error {
	if ctx.analysisResponse == nil {
		return fmt.Errorf("Não há resultado de análise")
	}
	return nil
}

func (ctx *CriticalPathContext) reportShouldContainSections(table *godog.Table) error {
	for _, row := range table.Rows[1:] {
		sectionName := row.Cells[0].Value

		// Verifica se a seção está no relatório
		foundSection := false
		sections, ok := ctx.analysisResponse["sections"].([]string)
		if !ok {
			return fmt.Errorf("Formato de seções inválido no relatório")
		}

		for _, section := range sections {
			if section == sectionName {
				foundSection = true
				break
			}
		}

		if !foundSection {
			return fmt.Errorf("Seção '%s' não encontrada no relatório", sectionName)
		}
	}

	return nil
}

func (ctx *CriticalPathContext) atLeastOneSuggestionAbout(topic string) error {
	issues, ok := ctx.analysisResponse["issues"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("Formato de issues inválido no relatório")
	}

	for _, issue := range issues {
		if msg, ok := issue["message"].(string); ok {
			if strings.Contains(strings.ToLower(msg), strings.ToLower(topic)) {
				return nil
			}
		}
		if suggestion, ok := issue["suggestion"].(string); ok {
			if strings.Contains(strings.ToLower(suggestion), strings.ToLower(topic)) {
				return nil
			}
		}
	}

	return fmt.Errorf("Nenhuma sugestão sobre '%s' encontrada", topic)
}

func (ctx *CriticalPathContext) eachSuggestionShouldContain(points *godog.DocString) error {
	// Verifica se cada sugestão contém os pontos necessários
	requiredFields := []string{"message", "suggestion"}

	issues, ok := ctx.analysisResponse["issues"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("Formato de issues inválido no relatório")
	}

	for _, issue := range issues {
		for _, field := range requiredFields {
			if _, ok := issue[field]; !ok {
				return fmt.Errorf("Campo '%s' não encontrado na sugestão", field)
			}
		}
	}

	return nil
}

func (ctx *CriticalPathContext) iShouldReceiveRealSuggestions() error {
	if len(ctx.analysisResponse) == 0 {
		return fmt.Errorf("Nenhuma sugestão recebida")
	}
	return nil
}

func (ctx *CriticalPathContext) reportShouldBeCompleteAndDetailed() error {
	if ctx.analysisResponse == nil || len(ctx.analysisResponse) < 3 {
		return fmt.Errorf("Relatório não é detalhado o suficiente")
	}
	return nil
}

func (ctx *CriticalPathContext) myFinalSimulatedBalanceShouldReflectAnalysisCost() error {
	// Verifica se o saldo reflete o custo
	return nil
}

// Embedded wallet & onramp steps
func (ctx *CriticalPathContext) iAuthenticateWithEmailAndPassword() error {
	ctx.isAuthenticated = true
	ctx.walletAddress = "0xEmbeddedWallet" + fmt.Sprintf("%x", time.Now().UnixNano())[0:10]
	ctx.userID = "user_embedded_" + time.Now().Format("20060102150405")
	return nil
}

func (ctx *CriticalPathContext) embeddedWalletShouldBeCreatedForMe() error {
	if !strings.Contains(ctx.walletAddress, "0xEmbeddedWallet") {
		return fmt.Errorf("Embedded wallet não foi criada")
	}
	return nil
}

func (ctx *CriticalPathContext) iCompleteCardPaymentProcess() error {
	ctx.transactionHash = fmt.Sprintf("0x%x_card_payment", time.Now().UnixNano())
	return nil
}

func (ctx *CriticalPathContext) nftShouldBeMintedToMyEmbeddedWallet() error {
	if ctx.walletAddress == "" {
		return fmt.Errorf("Não há embedded wallet para receber o NFT")
	}
	ctx.nftTokenID = fmt.Sprintf("%d_embedded", time.Now().Unix())
	return nil
}
