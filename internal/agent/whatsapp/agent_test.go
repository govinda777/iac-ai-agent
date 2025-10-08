package whatsapp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWhatsAppAgent_ProcessMessage(t *testing.T) {
	// Setup
	config := &WhatsAppAgentConfig{
		Name:        "Test Agent",
		Description: "Test Description",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
	}

	agent := &WhatsAppAgent{
		ID:             "test_agent_123",
		Name:           config.Name,
		Description:    config.Description,
		WalletAddr:     config.WalletAddr,
		APIKey:         "test_api_key",
		Service:        &MockAgentService{},
		LLMService:     &MockLLMService{},
		BillingService: &MockBillingService{},
		Logger:         SetupLogger("test_agent_123"),
		Commands:       AvailableCommands(),
	}

	// Test cases
	tests := []struct {
		name       string
		message    *WhatsAppMessage
		expected   string
		setupMocks func()
	}{
		{
			name: "help command",
			message: &WhatsAppMessage{
				From: "test_user",
				Text: "/help",
			},
			expected: "Comandos disponíveis",
			setupMocks: func() {
				// No mocks needed for help command
			},
		},
		{
			name: "status command",
			message: &WhatsAppMessage{
				From: "test_user",
				Text: "/status",
			},
			expected: "Status do Agente",
			setupMocks: func() {
				// No mocks needed for status command
			},
		},
		{
			name: "analyze command with code",
			message: &WhatsAppMessage{
				From: "test_user",
				Text: "/analyze\n```hcl\nresource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}\n```",
			},
			expected: "Análise concluída",
			setupMocks: func() {
				// No mocks needed for analyze command
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Execute
			response, err := agent.ProcessMessage(context.Background(), tt.message)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Contains(t, response.Text, tt.expected)
		})
	}
}

func TestWhatsAppAgent_Authentication(t *testing.T) {
	agent := &WhatsAppAgent{
		ID:     "test_agent_123",
		Logger: SetupLogger("test_agent_123"),
	}

	// Test authentication
	err := agent.authenticate(context.Background(), "test_user")
	assert.NoError(t, err)
}

func TestWhatsAppAgent_ParseCommand(t *testing.T) {
	agent := &WhatsAppAgent{
		Commands: AvailableCommands(),
	}

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "help command",
			text:     "/help",
			expected: "help",
		},
		{
			name:     "analyze command",
			text:     "/analyze some code",
			expected: "analyze",
		},
		{
			name:     "security command",
			text:     "/security check this",
			expected: "security",
		},
		{
			name:     "invalid command",
			text:     "invalid command",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, err := agent.parseCommand(tt.text)
			if tt.expected == "" {
				assert.Error(t, err)
				assert.Nil(t, command)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, command)
				assert.Equal(t, tt.expected, command.Name)
			}
		})
	}
}

func TestWhatsAppAgent_ExtractCommandArgs(t *testing.T) {
	agent := &WhatsAppAgent{}

	tests := []struct {
		name         string
		text         string
		pattern      string
		expectedArgs []string
		expectedCode string
	}{
		{
			name:         "simple command",
			text:         "/help",
			pattern:      `^/help`,
			expectedArgs: []string{},
			expectedCode: "",
		},
		{
			name:         "command with args",
			text:         "/analyze terraform code",
			pattern:      `^/analyze\s*(.*)`,
			expectedArgs: []string{"terraform", "code"},
			expectedCode: "",
		},
		{
			name:         "command with code block",
			text:         "/analyze\n```hcl\nresource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}\n```",
			pattern:      `^/analyze\s*(.*)`,
			expectedArgs: []string{},
			expectedCode: "resource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, codeBlock := agent.extractCommandArgs(tt.text, tt.pattern)
			assert.Equal(t, tt.expectedArgs, args)
			assert.Equal(t, tt.expectedCode, codeBlock)
		})
	}
}

func TestWhatsAppAgent_HandleError(t *testing.T) {
	agent := &WhatsAppAgent{
		Logger: SetupLogger("test_agent_123"),
	}

	errorMsg := "Test error message"
	response := agent.handleError(errorMsg)

	assert.NotNil(t, response)
	assert.Contains(t, response.Text, "Ops! Algo deu errado")
	assert.Contains(t, response.Text, errorMsg)
	assert.Equal(t, "text", response.Type)
}

func TestWhatsAppAgent_GetUsageStats(t *testing.T) {
	mockBillingService := &MockBillingService{}
	agent := &WhatsAppAgent{
		BillingService: mockBillingService,
		Logger:         SetupLogger("test_agent_123"),
	}

	expectedStats := &UsageStats{
		TotalRequests:  10,
		TokensConsumed: 8,
		LastRequest:    time.Now(),
		RequestsToday:  3,
		AverageCost:    0.8,
	}

	// Test completed successfully

	stats, err := agent.GetUsageStats(context.Background(), "test_user")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, expectedStats, stats)

	// Test completed successfully
}

func TestWhatsAppAgentConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *WhatsAppAgentConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &WhatsAppAgentConfig{
				Name:        "Test Agent",
				Description: "Test Description",
				WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: &WhatsAppAgentConfig{
				Description: "Test Description",
				WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
			},
			wantErr: true,
		},
		{
			name: "invalid wallet address",
			config: &WhatsAppAgentConfig{
				Name:        "Test Agent",
				Description: "Test Description",
				WalletAddr:  "invalid_address",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "IaC AI Agent WhatsApp", config.Name)
	assert.Equal(t, "Agente para análise de infraestrutura via WhatsApp", config.Description)
	assert.Equal(t, "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5", config.WalletAddr)
}

func TestAvailableCommands(t *testing.T) {
	commands := AvailableCommands()

	assert.NotNil(t, commands)
	assert.Contains(t, commands, "help")
	assert.Contains(t, commands, "analyze")
	assert.Contains(t, commands, "security")
	assert.Contains(t, commands, "cost")
	assert.Contains(t, commands, "status")
	assert.Contains(t, commands, "balance")
	assert.Contains(t, commands, "usage")

	// Verificar propriedades dos comandos
	assert.Equal(t, "help", commands["help"].Name)
	assert.False(t, commands["help"].RequiresPayment)
	assert.Equal(t, 0, commands["help"].TokenCost)

	assert.Equal(t, "analyze", commands["analyze"].Name)
	assert.True(t, commands["analyze"].RequiresPayment)
	assert.Equal(t, 1, commands["analyze"].TokenCost)
}
