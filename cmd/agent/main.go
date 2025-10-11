package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/govinda777/iac-ai-agent/api/rest"
	"github.com/govinda777/iac-ai-agent/internal/agent/capabilities"
	"github.com/govinda777/iac-ai-agent/internal/agent/core"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

func main() {
	// Carregar configuração
	agentConfig, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Criar agente principal
	agent := core.NewAgent(agentConfig)

	// Registrar habilidades
	if err := registerCapabilities(agent); err != nil {
		log.Fatalf("Failed to register capabilities: %v", err)
	}

	// Inicializar agente
	ctx := context.Background()
	if err := agent.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// Iniciar agente
	if err := agent.Start(ctx); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	log.Printf("Agent started successfully: %s", agent.ID)

	// Criar configuração e logger para o handler REST
	restConfig := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
	restLogger := logger.New("info", "json")

	// Criar handler REST com Swagger
	handler := rest.NewHandler(restConfig, restLogger)

	// Configurar rotas com Swagger UI
	router := handler.SetupRoutes()

	// Middleware (comentado temporariamente para teste)
	// router.Use(rest.LoggingMiddleware)
	// router.Use(rest.TokenValidationMiddleware("your_verify_token_here"))

	// Registrar rotas do agente também
	agentHandler := rest.NewAgentHandler(agent)
	agentHandler.RegisterRoutes(router)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Starting agent server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down agent...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parar agente
	if err := agent.Stop(ctx); err != nil {
		log.Printf("Failed to stop agent: %v", err)
	}

	// Parar servidor
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Agent stopped successfully")
}

// loadConfig carrega configuração do agente
func loadConfig() (*core.Config, error) {
	// Configuração padrão
	config := &core.Config{
		AgentID:     "iac-ai-agent",
		AgentName:   "IaC AI Agent",
		Description: "Agente inteligente para análise de Infrastructure as Code",
		Version:     "1.0.0",
		Capabilities: map[string]interface{}{
			"whatsapp": map[string]interface{}{
				"webhook_url":  getEnv("WHATSAPP_WEBHOOK_URL", ""),
				"verify_token": getEnv("WHATSAPP_VERIFY_TOKEN", "your_verify_token_here"),
				"api_key":      getEnv("WHATSAPP_API_KEY", ""),
			},
			"iac-analysis": map[string]interface{}{
				"terraform_enabled": true,
				"security_enabled":  true,
				"cost_enabled":      true,
			},
		},
		Logging: core.LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			File:       getEnv("LOG_FILE", "/var/log/iac-ai-agent.log"),
			MaxSize:    "100MB",
			MaxBackups: 3,
			MaxAge:     28,
		},
		Web3: core.Web3Config{
			WalletAddress: getEnv("WALLET_ADDRESS", "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"),
			NFTContract:   getEnv("NFT_CONTRACT", "nation.fun"),
			LitProtocol:   true,
		},
		Billing: core.BillingConfig{
			Enabled: true,
			TokenCosts: map[string]int{
				"analyze":  1,
				"security": 1,
				"cost":     1,
			},
			FreeCommands: []string{"help", "status", "ping"},
		},
	}

	return config, nil
}

// registerCapabilities registra todas as habilidades do agente
func registerCapabilities(agent *core.Agent) error {
	// Registrar habilidade WhatsApp
	whatsappCapability := capabilities.NewWhatsAppCapability()
	if err := agent.RegisterCapability(whatsappCapability); err != nil {
		return fmt.Errorf("failed to register WhatsApp capability: %w", err)
	}

	// Registrar habilidade de análise IaC
	iacCapability := capabilities.NewIACAnalysisCapability()
	if err := agent.RegisterCapability(iacCapability); err != nil {
		return fmt.Errorf("failed to register IaC Analysis capability: %w", err)
	}

	// Conectar habilidades (WhatsApp pode usar IaC Analysis)
	whatsappCapability.SetIACCapability(iacCapability)

	log.Printf("Capabilities registered successfully")
	return nil
}

// getEnv obtém variável de ambiente com valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}