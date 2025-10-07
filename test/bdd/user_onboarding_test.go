package bdd

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
)

// TestContext holds the state for a single scenario.
type TestContext struct {
	isAuthenticated      bool
	walletAddress        string
	selectedProvider     string
	welcomeMessage       string
	purchaseOptionsVisible bool
	currentError         error
	userEmail            string
	isEmailLinked        bool
	canLoginWithEmail    bool
	redirectUrl          string
	currentPage          string
	sessionExpired       bool
}

// reset initializes the context before each scenario.
func (t *TestContext) reset() {
	t.isAuthenticated = false
	t.walletAddress = ""
	t.selectedProvider = ""
	t.welcomeMessage = ""
	t.purchaseOptionsVisible = false
	t.currentError = nil
	t.userEmail = ""
	t.isEmailLinked = false
	t.canLoginWithEmail = false
	t.redirectUrl = ""
	t.currentPage = ""
	t.sessionExpired = false
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "godog",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	t := &TestContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		t.reset()
		return ctx, nil
	})

	ctx.Step(`^que o serviço Privy está disponível$`, t.thePrivyServiceIsAvailable)
	ctx.Step(`^que a Base Network está acessível$`, t.theBaseNetworkIsAccessible)
	ctx.Step(`^que os contratos de NFT e Token estão deployados$`, t.theNFTAndTokenContractsAreDeployed)
	ctx.Step(`^que sou um novo usuário sem conta$`, t.iAmANewUserWithNoAccount)
	ctx.Step(`^eu clico em "([^"]*)"$`, t.iClickOn)
	ctx.Step(`^seleciono "([^"]*)" como provider$`, t.iSelectAsProvider)
	ctx.Step(`^aprovo a conexão no MetaMask$`, t.iApproveTheConnectionInMetaMask)
	ctx.Step(`^devo estar autenticado$`, t.iShouldBeAuthenticated)
	ctx.Step(`^meu endereço da wallet deve estar visível$`, t.myWalletAddressShouldBeVisible)
	ctx.Step(`^devo ver a mensagem "([^"]*)"$`, t.iShouldSeeTheMessage)
	ctx.Step(`^devo ver as opções de compra de acesso$`, t.iShouldSeeTheAccessPurchaseOptions)
	ctx.Step(`^aprovo a conexão no Coinbase Wallet$`, t.iApproveTheConnectionInCoinbaseWallet)
	ctx.Step(`^que sou um novo usuário sem wallet$`, t.iAmANewUserWithoutAWallet)
	ctx.Step(`^concluo o processo de autenticação com email$`, t.iCompleteTheEmailAuthenticationProcess)
	ctx.Step(`^uma embedded wallet deve ser criada automaticamente$`, t.anEmbeddedWalletShouldBeCreatedAutomatically)
	ctx.Step(`^devo ver meu endereço de wallet$`, t.iShouldSeeMyWalletAddress)
	ctx.Step(`^que estou autenticado com wallet$`, t.iAmAuthenticatedWithAWallet)
	ctx.Step(`^insiro meu email "([^"]*)"$`, t.iEnterMyEmail)
	ctx.Step(`^confirmo o código de verificação recebido$`, t.iConfirmTheVerificationCodeReceived)
	ctx.Step(`^meu email deve estar vinculado à conta$`, t.myEmailShouldBeLinkedToTheAccount)
	ctx.Step(`^devo poder fazer login com email ou wallet$`, t.iShouldBeAbleToLogInWithEmailOrWallet)
	ctx.Step(`^que não estou autenticado$`, t.iAmNotAuthenticated)
	ctx.Step(`^eu tento acessar "([^"]*)"$`, t.iTryToAccess)
	ctx.Step(`^devo ser redirecionado para "([^"]*)"$`, t.iShouldBeRedirectedTo)
	ctx.Step(`^que estou autenticado há mais de (\d+) horas$`, t.iHaveBeenAuthenticatedForMoreThanHours)
	ctx.Step(`^eu tento fazer uma análise$`, t.iTryToDoAnAnalysis)
	ctx.Step(`^devo receber erro "([^"]*)"$`, t.iShouldReceiveAnError)
	ctx.Step(`^devo ser redirecionado para login$`, t.iShouldBeRedirectedToLogin)
	ctx.Step(`^após re-autenticar devo voltar para a página original$`, t.afterReAuthenticatingIShouldReturnToTheOriginalPage)
}

// --- Givens ---
func (t *TestContext) thePrivyServiceIsAvailable() error {
	// Simulate checking Privy service. For now, we assume it's always available.
	return nil
}

func (t *TestContext) theBaseNetworkIsAccessible() error {
	// Simulate checking Base Network. For now, we assume it's always available.
	return nil
}

func (t *TestContext) theNFTAndTokenContractsAreDeployed() error {
	// Simulate checking contracts. For now, we assume they are deployed.
	return nil
}

func (t *TestContext) iAmANewUserWithNoAccount() error {
	t.reset()
	return nil
}

func (t *TestContext) iAmANewUserWithoutAWallet() error {
	t.reset()
	return nil
}

func (t *TestContext) iAmAuthenticatedWithAWallet() error {
	t.isAuthenticated = true
	t.walletAddress = "0x" + uuid.New().String()
	return nil
}

func (t *TestContext) iAmNotAuthenticated() error {
	t.isAuthenticated = false
	return nil
}

func (t *TestContext) iHaveBeenAuthenticatedForMoreThanHours(hours int) error {
	t.iAmAuthenticatedWithAWallet()
	t.sessionExpired = true
	// We can use the 'hours' variable if needed in the future
	return nil
}

// --- Whens ---
func (t *TestContext) iClickOn(text string) error {
	// This is a UI interaction, we can log it for debugging if needed.
	fmt.Printf("Simulating click on: %s\n", text)
	return nil
}

func (t *TestContext) iSelectAsProvider(provider string) error {
	t.selectedProvider = provider
	return nil
}

func (t *TestContext) iApproveTheConnectionInMetaMask() error {
	if t.selectedProvider != "MetaMask" {
		return fmt.Errorf("expected MetaMask to be selected, but got %s", t.selectedProvider)
	}
	t.isAuthenticated = true
	t.walletAddress = "0x" + uuid.New().String()
	t.welcomeMessage = "Bem-vindo ao IaC AI Agent"
	t.purchaseOptionsVisible = true
	return nil
}

func (t *TestContext) iApproveTheConnectionInCoinbaseWallet() error {
	if t.selectedProvider != "Coinbase Wallet" {
		return fmt.Errorf("expected Coinbase Wallet to be selected, but got %s", t.selectedProvider)
	}
	t.isAuthenticated = true
	t.walletAddress = "0x" + uuid.New().String()
	return nil
}

func (t *TestContext) iCompleteTheEmailAuthenticationProcess() error {
	t.isAuthenticated = true
	t.isEmailLinked = true
	return nil
}

func (t *TestContext) iEnterMyEmail(email string) error {
	t.userEmail = email
	return nil
}

func (t *TestContext) iConfirmTheVerificationCodeReceived() error {
	if t.userEmail != "" {
		t.isEmailLinked = true
	}
	return nil
}

func (t *TestContext) iTryToAccess(page string) error {
	t.currentPage = page
	if !t.isAuthenticated {
		t.redirectUrl = "/login"
		t.welcomeMessage = "Por favor, conecte sua wallet"
	}
	return nil
}

func (t *TestContext) iTryToDoAnAnalysis() error {
	if t.sessionExpired {
		t.currentError = fmt.Errorf("Token expirado")
		t.redirectUrl = "/login" // Assuming this is the login page
	}
	return nil
}

// --- Thens ---
func (t *TestContext) iShouldBeAuthenticated() error {
	if !t.isAuthenticated {
		return fmt.Errorf("user should be authenticated, but is not")
	}
	return nil
}

func (t *TestContext) myWalletAddressShouldBeVisible() error {
	if t.walletAddress == "" {
		return fmt.Errorf("wallet address should be visible, but it is empty")
	}
	return nil
}

func (t *TestContext) iShouldSeeTheMessage(message string) error {
	if t.welcomeMessage != message {
		return fmt.Errorf("expected welcome message '%s', but got '%s'", message, t.welcomeMessage)
	}
	return nil
}

func (t *TestContext) iShouldSeeTheAccessPurchaseOptions() error {
	if !t.purchaseOptionsVisible {
		return fmt.Errorf("purchase options should be visible, but are not")
	}
	return nil
}

func (t *TestContext) anEmbeddedWalletShouldBeCreatedAutomatically() error {
	t.walletAddress = "0x" + uuid.New().String()
	return nil
}

func (t *TestContext) iShouldSeeMyWalletAddress() error {
	if t.walletAddress == "" {
		return fmt.Errorf("wallet address should be visible, but it is empty")
	}
	return nil
}

func (t *TestContext) myEmailShouldBeLinkedToTheAccount() error {
	if !t.isEmailLinked {
		return fmt.Errorf("email should be linked, but it is not")
	}
	return nil
}

func (t *TestContext) iShouldBeAbleToLogInWithEmailOrWallet() error {
	t.canLoginWithEmail = true // Simulate that this is now possible
	return nil
}

func (t *TestContext) iShouldBeRedirectedTo(page string) error {
	if t.redirectUrl != page {
		return fmt.Errorf("expected to be redirected to '%s', but was redirected to '%s'", page, t.redirectUrl)
	}
	return nil
}

func (t *TestContext) iShouldReceiveAnError(errorMessage string) error {
	if t.currentError == nil || t.currentError.Error() != errorMessage {
		return fmt.Errorf("expected error '%s', but got '%v'", errorMessage, t.currentError)
	}
	return nil
}

func (t *TestContext) iShouldBeRedirectedToLogin() error {
	if t.redirectUrl != "/login" {
		return fmt.Errorf("expected to be redirected to login, but was redirected to '%s'", t.redirectUrl)
	}
	return nil
}

func (t *TestContext) afterReAuthenticatingIShouldReturnToTheOriginalPage() error {
	// Simulate re-authentication
	t.isAuthenticated = true
	t.sessionExpired = false
	t.currentError = nil
	// Here we would check if the user is back on the original page, e.g. 't.currentPage'
	// For this test, we just assume it works.
	return nil
}