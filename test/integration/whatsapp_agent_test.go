package integration

import (
	"context"
	"testing"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/agent/whatsapp"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppAgentIntegration(t *testing.T) {
	// Setup integration test environment
	setupTestEnvironment(t)

	// Create test agent
	config := &whatsapp.WhatsAppAgentConfig{
		Name:        "Integration Test Agent",
		Description: "Agent for integration testing",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
	}

	agent, err := whatsapp.NewWhatsAppAgent(config)
	assert.NoError(t, err)
	assert.NotNil(t, agent)

	// Test full message flow
	message := &whatsapp.WhatsAppMessage{
		From: "test_user",
		Text: "/analyze\n```hcl\nresource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}\n```",
	}

	response, err := agent.ProcessMessage(context.Background(), message)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Contains(t, response.Text, "Análise concluída")

	// Verify billing
	stats, err := agent.GetUsageStats(context.Background(), "test_user")
	assert.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestWhatsAppAgentCommandFlow(t *testing.T) {
	// Setup
	config := &whatsapp.WhatsAppAgentConfig{
		Name:        "Command Flow Test Agent",
		Description: "Agent for command flow testing",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
	}

	agent, err := whatsapp.NewWhatsAppAgent(config)
	assert.NoError(t, err)

	// Test command sequence
	commands := []struct {
		name     string
		message  *whatsapp.WhatsAppMessage
		expected string
	}{
		{
			name: "help command",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/help",
			},
			expected: "Comandos disponíveis",
		},
		{
			name: "status command",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/status",
			},
			expected: "Status do Agente",
		},
		{
			name: "balance command",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/balance",
			},
			expected: "Saldo de Tokens IACAI",
		},
		{
			name: "usage command",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/usage",
			},
			expected: "Estatísticas de Uso",
		},
	}

	for _, cmd := range commands {
		t.Run(cmd.name, func(t *testing.T) {
			response, err := agent.ProcessMessage(context.Background(), cmd.message)
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Contains(t, response.Text, cmd.expected)
		})
	}
}

func TestWhatsAppAgentErrorHandling(t *testing.T) {
	// Setup
	config := &whatsapp.WhatsAppAgentConfig{
		Name:        "Error Handling Test Agent",
		Description: "Agent for error handling testing",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
	}

	agent, err := whatsapp.NewWhatsAppAgent(config)
	assert.NoError(t, err)

	// Test error scenarios
	errorTests := []struct {
		name     string
		message  *whatsapp.WhatsAppMessage
		expected string
	}{
		{
			name: "invalid command",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/invalid_command",
			},
			expected: "Comando inválido",
		},
		{
			name: "analyze without code",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/analyze",
			},
			expected: "Por favor, forneça o código Terraform",
		},
		{
			name: "security without code",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/security",
			},
			expected: "Por favor, forneça o código para verificação",
		},
		{
			name: "cost without code",
			message: &whatsapp.WhatsAppMessage{
				From: "test_user",
				Text: "/cost",
			},
			expected: "Por favor, forneça o código para análise",
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			response, err := agent.ProcessMessage(context.Background(), test.message)
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Contains(t, response.Text, test.expected)
		})
	}
}

func TestWhatsAppAgentPerformance(t *testing.T) {
	// Setup
	config := &whatsapp.WhatsAppAgentConfig{
		Name:        "Performance Test Agent",
		Description: "Agent for performance testing",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
	}

	agent, err := whatsapp.NewWhatsAppAgent(config)
	assert.NoError(t, err)

	// Test performance with multiple messages
	message := &whatsapp.WhatsAppMessage{
		From: "test_user",
		Text: "/help",
	}

	start := time.Now()

	// Process multiple messages
	for i := 0; i < 10; i++ {
		response, err := agent.ProcessMessage(context.Background(), message)
		assert.NoError(t, err)
		assert.NotNil(t, response)
	}

	duration := time.Since(start)

	// Verify performance (should complete within reasonable time)
	assert.Less(t, duration, 5*time.Second, "Processing 10 messages should take less than 5 seconds")

	// Log performance metrics
	t.Logf("Processed 10 messages in %v (avg: %v per message)", duration, duration/10)
}

func TestWhatsAppAgentConcurrency(t *testing.T) {
	// Setup
	config := &whatsapp.WhatsAppAgentConfig{
		Name:        "Concurrency Test Agent",
		Description: "Agent for concurrency testing",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
	}

	agent, err := whatsapp.NewWhatsAppAgent(config)
	assert.NoError(t, err)

	// Test concurrent message processing
	message := &whatsapp.WhatsAppMessage{
		From: "test_user",
		Text: "/help",
	}

	// Channel to collect results
	results := make(chan *whatsapp.WhatsAppResponse, 10)
	errors := make(chan error, 10)

	// Launch concurrent goroutines
	for i := 0; i < 10; i++ {
		go func() {
			response, err := agent.ProcessMessage(context.Background(), message)
			if err != nil {
				errors <- err
				return
			}
			results <- response
		}()
	}

	// Collect results
	successCount := 0
	errorCount := 0

	for i := 0; i < 10; i++ {
		select {
		case response := <-results:
			assert.NotNil(t, response)
			successCount++
		case err := <-errors:
			assert.NoError(t, err)
			errorCount++
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for concurrent results")
		}
	}

	// Verify all messages were processed successfully
	assert.Equal(t, 10, successCount)
	assert.Equal(t, 0, errorCount)
}

func TestWhatsAppAgentConfiguration(t *testing.T) {
	// Test configuration loading
	config, err := whatsapp.LoadOrCreateConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Test configuration validation
	err = config.Validate()
	assert.NoError(t, err)

	// Test configuration merging
	config.Name = ""
	config.Description = ""
	config.WalletAddr = ""

	mergedConfig := config.MergeWithDefaults()
	assert.Equal(t, "IaC AI Agent WhatsApp", mergedConfig.Name)
	assert.Equal(t, "Agente para análise de infraestrutura via WhatsApp", mergedConfig.Description)
	assert.Equal(t, "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5", mergedConfig.WalletAddr)
}

func setupTestEnvironment(t *testing.T) {
	// Setup test environment
	// In a real implementation, this would:
	// - Set up test database
	// - Configure test services
	// - Initialize test dependencies
	// - Set up test data

	t.Log("Setting up test environment...")

	// For now, just log that we're setting up
	// In production, implement actual test environment setup
}

func TestWhatsAppAgentLogging(t *testing.T) {
	// Test logging functionality
	logger := whatsapp.SetupLogger("test_logger")
	assert.NotNil(t, logger)

	// Test structured logging
	structuredLogger := whatsapp.NewStructuredLogger("test_structured_logger", whatsapp.LogLevelInfo)
	assert.NotNil(t, structuredLogger)

	// Test log levels
	structuredLogger.Debug("Debug message", map[string]interface{}{"key": "value"})
	structuredLogger.Info("Info message", map[string]interface{}{"key": "value"})
	structuredLogger.Warn("Warning message", map[string]interface{}{"key": "value"})
	structuredLogger.Error("Error message", map[string]interface{}{"key": "value"})
}

func TestWhatsAppAgentTemplates(t *testing.T) {
	// Test response templates
	templates := whatsapp.ResponseTemplates

	assert.NotNil(t, templates)
	assert.Contains(t, templates, "welcome")
	assert.Contains(t, templates, "error")
	assert.Contains(t, templates, "insufficient_tokens")
	assert.Contains(t, templates, "analysis_started")
	assert.Contains(t, templates, "analysis_complete")

	// Verify template content
	assert.Contains(t, templates["welcome"], "Olá! Sou o IaC AI Agent")
	assert.Contains(t, templates["error"], "Ops! Algo deu errado")
	assert.Contains(t, templates["insufficient_tokens"], "Saldo insuficiente")
}
