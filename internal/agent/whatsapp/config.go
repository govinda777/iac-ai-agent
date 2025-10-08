package whatsapp

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// DefaultConfig retorna configuração padrão
func DefaultConfig() *WhatsAppAgentConfig {
	return &WhatsAppAgentConfig{
		Name:        "IaC AI Agent WhatsApp",
		Description: "Agente para análise de infraestrutura via WhatsApp",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5", // Wallet padrão
		WebhookURL:  "",
		VerifyToken: "",
	}
}

// LoadConfigFromFile carrega configuração de arquivo YAML
func LoadConfigFromFile(filePath string) (*WhatsAppAgentConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config WhatsAppAgentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfigToFile salva configuração em arquivo YAML
func (c *WhatsAppAgentConfig) SaveConfigToFile(filePath string) error {
	// Criar diretório se não existir
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate valida a configuração
func (c *WhatsAppAgentConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("agent name is required")
	}

	if c.Description == "" {
		return fmt.Errorf("agent description is required")
	}

	if c.WalletAddr == "" {
		return fmt.Errorf("wallet address is required")
	}

	// Validar formato do endereço da wallet (básico)
	if len(c.WalletAddr) != 42 || c.WalletAddr[:2] != "0x" {
		return fmt.Errorf("invalid wallet address format")
	}

	return nil
}

// MergeWithDefaults mescla configuração com valores padrão
func (c *WhatsAppAgentConfig) MergeWithDefaults() *WhatsAppAgentConfig {
	defaults := DefaultConfig()

	if c.Name == "" {
		c.Name = defaults.Name
	}

	if c.Description == "" {
		c.Description = defaults.Description
	}

	if c.WalletAddr == "" {
		c.WalletAddr = defaults.WalletAddr
	}

	return c
}

// GetConfigPath retorna caminho padrão para arquivo de configuração
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./whatsapp-agent-config.yaml"
	}

	return filepath.Join(homeDir, ".iac-ai-agent", "whatsapp-agent-config.yaml")
}

// LoadOrCreateConfig carrega configuração existente ou cria nova
func LoadOrCreateConfig() (*WhatsAppAgentConfig, error) {
	configPath := GetConfigPath()

	// Tentar carregar configuração existente
	if _, err := os.Stat(configPath); err == nil {
		config, err := LoadConfigFromFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load existing config: %w", err)
		}
		return config.MergeWithDefaults(), nil
	}

	// Criar nova configuração com valores padrão
	config := DefaultConfig()
	if err := config.SaveConfigToFile(configPath); err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	return config, nil
}
