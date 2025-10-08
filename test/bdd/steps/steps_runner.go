package steps

import (
	"github.com/cucumber/godog"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// RegisterAllSteps registra todos os steps nos testes BDD
func RegisterAllSteps(ctx *godog.ScenarioContext, cfg *config.Config, log *logger.Logger) {
	// Contexto de caminho crítico
	criticalPathCtx := NewCriticalPathContext(cfg, log)
	RegisterCriticalPathSteps(ctx, criticalPathCtx)

	// Contexto de onboarding
	onboardingCtx := NewOnboardingContext(cfg, log)
	RegisterOnboardingSteps(ctx, onboardingCtx)

	// Contexto de NFT
	nftCtx := NewNFTPurchaseContext(cfg, log)
	RegisterNFTPurchaseSteps(ctx, nftCtx)

	// Contexto de tokens
	tokenCtx := NewTokenPurchaseContext(cfg, log)
	RegisterTokenPurchaseSteps(ctx, tokenCtx)

	// Contexto de análise
	analysisCtx := NewBotAnalysisContext(cfg, log)
	RegisterBotAnalysisSteps(ctx, analysisCtx)
}

// NewOnboardingContext cria contexto para testes de onboarding
func NewOnboardingContext(cfg *config.Config, log *logger.Logger) *OnboardingContext {
	return &OnboardingContext{
		config: cfg,
		logger: log,
	}
}

// OnboardingContext mantém o estado para testes de onboarding
type OnboardingContext struct {
	config *config.Config
	logger *logger.Logger
	// outros campos necessários
}

// RegisterOnboardingSteps registra os steps de onboarding
func RegisterOnboardingSteps(ctx *godog.ScenarioContext, oCtx *OnboardingContext) {
	// Implementar registro de steps
}

// NewNFTPurchaseContext cria contexto para testes de compra de NFT
func NewNFTPurchaseContext(cfg *config.Config, log *logger.Logger) *NFTPurchaseContext {
	return &NFTPurchaseContext{
		config: cfg,
		logger: log,
	}
}

// NFTPurchaseContext mantém o estado para testes de compra de NFT
type NFTPurchaseContext struct {
	config *config.Config
	logger *logger.Logger
	// outros campos necessários
}

// RegisterNFTPurchaseSteps registra os steps de compra de NFT
func RegisterNFTPurchaseSteps(ctx *godog.ScenarioContext, nftCtx *NFTPurchaseContext) {
	// Implementar registro de steps
}

// NewTokenPurchaseContext cria contexto para testes de compra de tokens
func NewTokenPurchaseContext(cfg *config.Config, log *logger.Logger) *TokenPurchaseContext {
	return &TokenPurchaseContext{
		config: cfg,
		logger: log,
	}
}

// TokenPurchaseContext mantém o estado para testes de compra de tokens
type TokenPurchaseContext struct {
	config *config.Config
	logger *logger.Logger
	// outros campos necessários
}

// RegisterTokenPurchaseSteps registra os steps de compra de tokens
func RegisterTokenPurchaseSteps(ctx *godog.ScenarioContext, tokenCtx *TokenPurchaseContext) {
	// Implementar registro de steps
}

// NewBotAnalysisContext cria contexto para testes de análise
func NewBotAnalysisContext(cfg *config.Config, log *logger.Logger) *BotAnalysisContext {
	return &BotAnalysisContext{
		config: cfg,
		logger: log,
	}
}

// BotAnalysisContext mantém o estado para testes de análise
type BotAnalysisContext struct {
	config *config.Config
	logger *logger.Logger
	// outros campos necessários
}

// RegisterBotAnalysisSteps registra os steps de análise
func RegisterBotAnalysisSteps(ctx *godog.ScenarioContext, analysisCtx *BotAnalysisContext) {
	// Implementar registro de steps
}
