package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/govinda777/iac-ai-agent/api/rest"
	"github.com/govinda777/iac-ai-agent/internal/agent/capabilities"
	"github.com/govinda777/iac-ai-agent/internal/agent/core"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

var (
	router http.Handler
	once   sync.Once
)

// Handler is the main entry point for Vercel Serverless Functions.
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		// This block is executed only once per container instance (on cold start).

		// 1. Create configuration and logger for the REST handler.
		// Vercel manages the host and port, so we use dummy values here.
		// Logging should go to stdout/stderr, which Vercel captures.
		restConfig := &config.Config{
			Server: config.ServerConfig{
				Port: "8080",
				Host: "localhost",
			},
			Logging: config.LoggingConfig{
				Level:  "info",
				Format: "console", // Use "console" for Vercel logs
			},
			// In a real-world scenario, populate other fields like LLM, Web3 from env vars.
			// Example:
			// LLM: config.LLMConfig{
			// 	APIKey: os.Getenv("LLM_API_KEY"),
			// },
		}
		restLogger := logger.New("info", "console")

		// 2. Create the core agent.
		agentConfig, err := loadAgentConfig()
		if err != nil {
			log.Fatalf("FATAL: Failed to load agent configuration: %v", err)
		}

		agent := core.NewAgent(agentConfig)

		// 3. Register agent capabilities.
		if err := registerCapabilities(agent); err != nil {
			log.Fatalf("FATAL: Failed to register capabilities: %v", err)
		}

		// The agent doesn't need to be explicitly started or stopped in a serverless context.
		// Initialization should be sufficient.
		if err := agent.Initialize(context.Background()); err != nil {
			log.Fatalf("FATAL: Failed to initialize agent: %v", err)
		}
		log.Println("Agent initialized successfully.")

		// 4. Create REST handlers.
		handler := rest.NewHandler(restConfig, restLogger)
		agentHandler := rest.NewAgentHandler(agent)

		// 5. Setup and combine routes.
		rtr := handler.SetupRoutes()
		agentHandler.RegisterRoutes(rtr)

		router = rtr
		log.Println("Router initialized successfully.")
	})

	// Serve the request using the initialized router.
	router.ServeHTTP(w, r)
}

// loadAgentConfig creates the agent's configuration, mirroring the logic from main.go.
// It should be adapted to use environment variables provided by Vercel.
func loadAgentConfig() (*core.Config, error) {
	// Default configuration, values can be overridden by environment variables.
	conf := &core.Config{
		AgentID:     "iac-ai-agent-vercel",
		AgentName:   "IaC AI Agent (Vercel)",
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
			File:       "", // No file logging in serverless
			MaxSize:    "10MB",
			MaxBackups: 1,
			MaxAge:     1,
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
	return conf, nil
}

// registerCapabilities registers all agent capabilities, mirroring the logic from main.go.
func registerCapabilities(agent *core.Agent) error {
	// Register WhatsApp capability
	whatsappCapability := capabilities.NewWhatsAppCapability()
	if err := agent.RegisterCapability(whatsappCapability); err != nil {
		return fmt.Errorf("failed to register WhatsApp capability: %w", err)
	}

	// Register IaC Analysis capability
	iacCapability := capabilities.NewIACAnalysisCapability()
	if err := agent.RegisterCapability(iacCapability); err != nil {
		return fmt.Errorf("failed to register IaC Analysis capability: %w", err)
	}

	// Connect capabilities (WhatsApp can use IaC Analysis)
	whatsappCapability.SetIACCapability(iacCapability)

	log.Printf("Capabilities registered successfully")
	return nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}