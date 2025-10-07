package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config representa a configuração completa da aplicação
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	LLM      LLMConfig      `yaml:"llm"`
	GitHub   GitHubConfig   `yaml:"github"`
	Analysis AnalysisConfig `yaml:"analysis"`
	Scoring  ScoringConfig  `yaml:"scoring"`
	Logging  LoggingConfig  `yaml:"logging"`
	Web3     Web3Config     `yaml:"web3"`
}

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// LLMConfig configurações do LLM
type LLMConfig struct {
	Provider    string  `yaml:"provider"` // openai, anthropic
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

// GitHubConfig configurações do GitHub
type GitHubConfig struct {
	Token         string `yaml:"token"`
	WebhookSecret string `yaml:"webhook_secret"`
	AutoComment   bool   `yaml:"auto_comment"`
}

// AnalysisConfig configurações de análise
type AnalysisConfig struct {
	CheckovEnabled          bool `yaml:"checkov_enabled"`
	IAMAnalysisEnabled      bool `yaml:"iam_analysis_enabled"`
	CostOptimizationEnabled bool `yaml:"cost_optimization_enabled"`
}

// ScoringConfig configurações de scoring
type ScoringConfig struct {
	MinPassScore int `yaml:"min_pass_score"`
}

// LoggingConfig configurações de logging
type LoggingConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error
	Format string `yaml:"format"` // json, text
}

// Web3Config configurações Web3 (Privy, Base Network, etc)
type Web3Config struct {
	// Privy Configuration
	PrivyAppID              string `yaml:"privy_app_id"`
	PrivyVerificationKeyURL string `yaml:"privy_verification_key_url"`

	// Base Network Configuration
	BaseRPCURL  string `yaml:"base_rpc_url"`
	BaseChainID int    `yaml:"base_chain_id"`

	// Smart Contract Addresses
	NFTAccessContractAddress string `yaml:"nft_access_contract_address"`
	BotTokenContractAddress  string `yaml:"bot_token_contract_address"`

	// Token Configuration
	TokenSymbol   string `yaml:"token_symbol"`
	TokenDecimals int    `yaml:"token_decimals"`

	// Features
	EnableNFTAccess     bool `yaml:"enable_nft_access"`
	EnableTokenPayments bool `yaml:"enable_token_payments"`

	// Rate Limiting by Tier (requests per hour)
	BasicTierRateLimit      int `yaml:"basic_tier_rate_limit"`
	ProTierRateLimit        int `yaml:"pro_tier_rate_limit"`
	EnterpriseTierRateLimit int `yaml:"enterprise_tier_rate_limit"`

	// Nation.fun Configuration
	WalletToken         string `yaml:"wallet_token"`          // Token de autenticação da wallet
	WalletAddress       string `yaml:"wallet_address"`        // Endereço da wallet
	DefaultAgentAddress string `yaml:"default_agent_address"` // Endereço do agente padrão
}

// Load carrega configuração de um arquivo YAML
func Load(path string) (*Config, error) {
	// Valores padrão
	config := &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
		LLM: LLMConfig{
			Provider:    "nation.fun",
			Model:       "nation-1",
			Temperature: 0.2,
			MaxTokens:   2000,
		},
		GitHub: GitHubConfig{
			AutoComment: true,
		},
		Analysis: AnalysisConfig{
			CheckovEnabled:          true,
			IAMAnalysisEnabled:      true,
			CostOptimizationEnabled: true,
		},
		Scoring: ScoringConfig{
			MinPassScore: 70,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}

	// Lê arquivo de configuração se existir
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Arquivo não existe, usa defaults
				fmt.Printf("Arquivo de configuração não encontrado, usando defaults\n")
			} else {
				return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
			}
		} else {
			// Parse YAML
			if err := yaml.Unmarshal(data, config); err != nil {
				return nil, fmt.Errorf("erro ao fazer parse do YAML: %w", err)
			}
		}
	}

	// Sobrescreve com variáveis de ambiente
	config.loadFromEnv()

	// Valida configuração
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// loadFromEnv carrega valores de variáveis de ambiente
func (c *Config) loadFromEnv() {
	// Server
	if port := os.Getenv("PORT"); port != "" {
		c.Server.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		c.Server.Host = host
	}

	// LLM
	if provider := os.Getenv("LLM_PROVIDER"); provider != "" {
		c.LLM.Provider = provider
	} else {
		// Default para Nation.fun
		c.LLM.Provider = "nation.fun"
	}
	if apiKey := os.Getenv("LLM_API_KEY"); apiKey != "" {
		c.LLM.APIKey = apiKey
	}
	if model := os.Getenv("LLM_MODEL"); model != "" {
		c.LLM.Model = model
	}

	// GitHub
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		c.GitHub.Token = token
	}
	if secret := os.Getenv("GITHUB_WEBHOOK_SECRET"); secret != "" {
		c.GitHub.WebhookSecret = secret
	}

	// Analysis
	if checkov := os.Getenv("CHECKOV_ENABLED"); checkov == "false" {
		c.Analysis.CheckovEnabled = false
	}
	if iam := os.Getenv("IAM_ANALYSIS_ENABLED"); iam == "false" {
		c.Analysis.IAMAnalysisEnabled = false
	}
	if cost := os.Getenv("COST_OPTIMIZATION_ENABLED"); cost == "false" {
		c.Analysis.CostOptimizationEnabled = false
	}

	// Logging
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		c.Logging.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		c.Logging.Format = format
	}

	// Web3
	if privyAppID := os.Getenv("PRIVY_APP_ID"); privyAppID != "" {
		c.Web3.PrivyAppID = privyAppID
	}
	if baseRPC := os.Getenv("BASE_RPC_URL"); baseRPC != "" {
		c.Web3.BaseRPCURL = baseRPC
	}
	if nftAddr := os.Getenv("NFT_CONTRACT_ADDRESS"); nftAddr != "" {
		c.Web3.NFTAccessContractAddress = nftAddr
	}
	if tokenAddr := os.Getenv("TOKEN_CONTRACT_ADDRESS"); tokenAddr != "" {
		c.Web3.BotTokenContractAddress = tokenAddr
	}

	// Nation.fun
	if walletToken := os.Getenv("WALLET_TOKEN"); walletToken != "" {
		c.Web3.WalletToken = walletToken
	}
	if walletAddr := os.Getenv("WALLET_ADDRESS"); walletAddr != "" {
		c.Web3.WalletAddress = walletAddr
	}
	if agentAddr := os.Getenv("DEFAULT_AGENT_ADDRESS"); agentAddr != "" {
		c.Web3.DefaultAgentAddress = agentAddr
	} else {
		// Agente padrão caso não seja especificado
		c.Web3.DefaultAgentAddress = "0x147e832418Cc06A501047019E956714271098b89"
	}
}

// Validate valida a configuração
func (c *Config) Validate() error {
	// Valida provider LLM
	validProviders := map[string]bool{
		"nation.fun": true,
		"nation":     true,
		"openai":     true, // Mantido para compatibilidade, mas será redirecionado para Nation.fun
		"anthropic":  true, // Mantido para compatibilidade, mas será redirecionado para Nation.fun
	}

	if c.LLM.Provider != "" && !validProviders[c.LLM.Provider] {
		return fmt.Errorf("LLM provider inválido: %s (use 'nation.fun')", c.LLM.Provider)
	}

	// Valida log level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("log level inválido: %s", c.Logging.Level)
	}

	// Valida Nation.fun config quando provider é nation.fun
	if c.LLM.Provider == "nation.fun" || c.LLM.Provider == "nation" {
		if c.Web3.NFTAccessContractAddress == "" {
			return fmt.Errorf("NFT_CONTRACT_ADDRESS é obrigatório para Nation.fun")
		}

		if c.Web3.WalletToken == "" {
			return fmt.Errorf("WALLET_TOKEN é obrigatório para Nation.fun")
		}
	}

	return nil
}

// GetAddress retorna o endereço completo do servidor
func (c *Config) GetAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}
