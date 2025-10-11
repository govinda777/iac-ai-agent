package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/agent/capabilities"
	"github.com/govinda777/iac-ai-agent/internal/agent/core"
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// HealthCheckResult represents the result of a health check.
type HealthCheckResult struct {
	Status   string        `json:"status"`
	Message  string        `json:"message"`
	Details  interface{}   `json:"details,omitempty"`
	Duration time.Duration `json:"duration"`
}

// RunComprehensiveHealthCheck executes a comprehensive health check of the agent.
func RunComprehensiveHealthCheck(ctx context.Context, cfg *config.Config, log *logger.Logger) map[string]HealthCheckResult {
	checks := make(map[string]HealthCheckResult)

	// Step 1: Find nation agents in the default wallet
	checks["find_nation_agents"] = runHealthCheckStep(ctx, func(ctx context.Context) HealthCheckResult {
		return findNationAgents(ctx, cfg, log)
	})

	// Step 2: Create a new agent with the template settings
	agent, err := createTestAgent(cfg, log)
	if err != nil {
		checks["create_agent_from_template"] = HealthCheckResult{
			Status:  "error",
			Message: fmt.Sprintf("Failed to create agent: %v", err),
		}
	} else {
		checks["create_agent_from_template"] = HealthCheckResult{
			Status:  "success",
			Message: "Successfully created a new agent from the template.",
			Details: map[string]interface{}{
				"agent_id": agent.ID,
			},
		}
		// Step 3: Send a test and validate the Agent's response
		checks["test_agent_response"] = runHealthCheckStep(ctx, func(ctx context.Context) HealthCheckResult {
			return testAgentResponse(ctx, agent)
		})
	}

	return checks
}

func runHealthCheckStep(ctx context.Context, stepFunc func(context.Context) HealthCheckResult) HealthCheckResult {
	return stepFunc(ctx)
}

func findNationAgents(ctx context.Context, cfg *config.Config, log *logger.Logger) HealthCheckResult {
	start := time.Now()
	nationValidator := web3.NewNationNFTValidator(cfg, log)
	_, err := nationValidator.ValidateWalletNFT(ctx, cfg.Web3.WalletAddress)
	if err != nil {
		return HealthCheckResult{
			Status:   "error",
			Message:  fmt.Sprintf("Failed to find nation agents: %v", err),
			Duration: time.Since(start),
		}
	}

	return HealthCheckResult{
		Status:   "success",
		Message:  "Successfully found Nation agents for the default wallet.",
		Duration: time.Since(start),
	}
}

func createTestAgent(cfg *config.Config, log *logger.Logger) (*core.Agent, error) {
	agentConfig := &core.Config{
		AgentID:     "test-agent",
		AgentName:   "Test Agent",
		Description: "A temporary agent for health check purposes.",
		Version:     "1.0.0",
		Capabilities: map[string]interface{}{
			"whatsapp": map[string]interface{}{
				"webhook_url":  "",
				"verify_token": "test_token",
				"api_key":      "",
			},
			"iac-analysis": map[string]interface{}{
				"terraform_enabled": true,
				"security_enabled":  true,
				"cost_enabled":      true,
			},
		},
		Logging: core.LoggingConfig{
			Level: "info",
		},
	}
	agent := core.NewAgent(agentConfig)

	// Register capabilities
	whatsappCapability := capabilities.NewWhatsAppCapability()
	iacAnalysisCapability := capabilities.NewIACAnalysisCapability()

	if err := agent.RegisterCapability(whatsappCapability); err != nil {
		return nil, fmt.Errorf("failed to register WhatsApp capability: %w", err)
	}
	if err := agent.RegisterCapability(iacAnalysisCapability); err != nil {
		return nil, fmt.Errorf("failed to register IaC Analysis capability: %w", err)
	}

	// Link capabilities
	whatsappCapability.SetIACCapability(iacAnalysisCapability)

	return agent, nil
}

func testAgentResponse(ctx context.Context, agent *core.Agent) HealthCheckResult {
	start := time.Now()

	message := &core.Message{
		ID:        "health-check-message",
		Source:    "whatsapp",
		Channel:   "health-check",
		From:      "health-checker",
		Text:      "/ping",
		Type:      "text",
		Timestamp: time.Now(),
	}

	response, err := agent.ProcessMessage(ctx, message)
	if err != nil {
		return HealthCheckResult{
			Status:   "error",
			Message:  fmt.Sprintf("Failed to process message: %v", err),
			Duration: time.Since(start),
		}
	}

	expectedResponse := "🏓 Pong! Agente WhatsApp ativo e funcionando."
	if response.Text != expectedResponse {
		return HealthCheckResult{
			Status:   "error",
			Message:  fmt.Sprintf("Unexpected response from agent: got %q, want %q", response.Text, expectedResponse),
			Duration: time.Since(start),
		}
	}

	return HealthCheckResult{
		Status:  "success",
		Message: "Agent responded to the test message successfully.",
		Details: map[string]interface{}{
			"response": response.Text,
		},
		Duration: time.Since(start),
	}
}