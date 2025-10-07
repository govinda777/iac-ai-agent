package models

import "time"

// Agent representa um agente de IA configurado
type Agent struct {
	ID          string    `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Version     string    `json:"version" yaml:"version"`
	Description string    `json:"description" yaml:"description"`
	Owner       string    `json:"owner" yaml:"owner"` // Wallet address
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" yaml:"updated_at"`
	Status      string    `json:"status" yaml:"status"` // active, inactive, training

	// Configurações
	Config AgentConfig `json:"config" yaml:"config"`

	// Habilidades
	Capabilities AgentCapabilities `json:"capabilities" yaml:"capabilities"`

	// Personalidade e Comportamento
	Personality AgentPersonality `json:"personality" yaml:"personality"`

	// Conhecimento Especializado
	Knowledge AgentKnowledge `json:"knowledge" yaml:"knowledge"`

	// Limites e Restrições
	Limits AgentLimits `json:"limits" yaml:"limits"`

	// Métricas de Performance
	Metrics AgentMetrics `json:"metrics" yaml:"metrics"`
}

// AgentConfig são as configurações técnicas do agente
type AgentConfig struct {
	// LLM Configuration
	LLMProvider      string  `json:"llm_provider" yaml:"llm_provider"` // openai, anthropic
	LLMModel         string  `json:"llm_model" yaml:"llm_model"`       // gpt-4, claude-3-opus
	Temperature      float64 `json:"temperature" yaml:"temperature"`   // 0.0 - 1.0
	MaxTokens        int     `json:"max_tokens" yaml:"max_tokens"`     // Token limit per request
	TopP             float64 `json:"top_p" yaml:"top_p"`               // Nucleus sampling
	FrequencyPenalty float64 `json:"frequency_penalty" yaml:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty" yaml:"presence_penalty"`

	// Analysis Configuration
	EnableCheckov         bool `json:"enable_checkov" yaml:"enable_checkov"`
	EnableIAMAnalysis     bool `json:"enable_iam_analysis" yaml:"enable_iam_analysis"`
	EnableCostAnalysis    bool `json:"enable_cost_analysis" yaml:"enable_cost_analysis"`
	EnableDriftDetection  bool `json:"enable_drift_detection" yaml:"enable_drift_detection"`
	EnablePreviewAnalysis bool `json:"enable_preview_analysis" yaml:"enable_preview_analysis"`
	EnableSecretsScanning bool `json:"enable_secrets_scanning" yaml:"enable_secrets_scanning"`

	// Response Configuration
	ResponseFormat      string `json:"response_format" yaml:"response_format"` // json, markdown, text
	IncludeCodeExamples bool   `json:"include_code_examples" yaml:"include_code_examples"`
	IncludeReferences   bool   `json:"include_references" yaml:"include_references"`
	DetailLevel         string `json:"detail_level" yaml:"detail_level"` // brief, standard, detailed

	// Language and Localization
	Language string `json:"language" yaml:"language"` // en, pt-br, es
	Timezone string `json:"timezone" yaml:"timezone"` // UTC, America/Sao_Paulo
}

// AgentCapabilities define o que o agente pode fazer
type AgentCapabilities struct {
	// Analysis Capabilities
	CanAnalyzeTerraform bool `json:"can_analyze_terraform" yaml:"can_analyze_terraform"`
	CanAnalyzeCheckov   bool `json:"can_analyze_checkov" yaml:"can_analyze_checkov"`
	CanAnalyzeIAM       bool `json:"can_analyze_iam" yaml:"can_analyze_iam"`
	CanAnalyzeCosts     bool `json:"can_analyze_costs" yaml:"can_analyze_costs"`
	CanDetectDrift      bool `json:"can_detect_drift" yaml:"can_detect_drift"`
	CanAnalyzePreview   bool `json:"can_analyze_preview" yaml:"can_analyze_preview"`
	CanScanSecrets      bool `json:"can_scan_secrets" yaml:"can_scan_secrets"`

	// Generation Capabilities
	CanGenerateCode          bool `json:"can_generate_code" yaml:"can_generate_code"`
	CanGenerateTests         bool `json:"can_generate_tests" yaml:"can_generate_tests"`
	CanGenerateDocumentation bool `json:"can_generate_documentation" yaml:"can_generate_documentation"`
	CanRefactorCode          bool `json:"can_refactor_code" yaml:"can_refactor_code"`

	// Advisory Capabilities
	CanSuggestArchitecture  bool `json:"can_suggest_architecture" yaml:"can_suggest_architecture"`
	CanSuggestModules       bool `json:"can_suggest_modules" yaml:"can_suggest_modules"`
	CanSuggestOptimizations bool `json:"can_suggest_optimizations" yaml:"can_suggest_optimizations"`
	CanSuggestSecurity      bool `json:"can_suggest_security" yaml:"can_suggest_security"`

	// Integration Capabilities
	CanIntegrateGitHub  bool `json:"can_integrate_github" yaml:"can_integrate_github"`
	CanIntegrateCI      bool `json:"can_integrate_ci" yaml:"can_integrate_ci"`
	CanIntegrateSlack   bool `json:"can_integrate_slack" yaml:"can_integrate_slack"`
	CanIntegrateDiscord bool `json:"can_integrate_discord" yaml:"can_integrate_discord"`

	// Learning Capabilities
	CanLearnFromFeedback   bool `json:"can_learn_from_feedback" yaml:"can_learn_from_feedback"`
	CanAdaptToContext      bool `json:"can_adapt_to_context" yaml:"can_adapt_to_context"`
	CanRememberPreferences bool `json:"can_remember_preferences" yaml:"can_remember_preferences"`
}

// AgentPersonality define a personalidade e estilo do agente
type AgentPersonality struct {
	Style         string `json:"style" yaml:"style"`         // professional, casual, friendly, technical
	Tone          string `json:"tone" yaml:"tone"`           // formal, informal, encouraging, direct
	Verbosity     string `json:"verbosity" yaml:"verbosity"` // concise, balanced, verbose
	UseEmojis     bool   `json:"use_emojis" yaml:"use_emojis"`
	UseHumor      bool   `json:"use_humor" yaml:"use_humor"`
	BeEncouraging bool   `json:"be_encouraging" yaml:"be_encouraging"`
	BeDirective   bool   `json:"be_directive" yaml:"be_directive"` // Give direct commands vs suggestions

	// Communication Preferences
	ExplainReasoning    bool `json:"explain_reasoning" yaml:"explain_reasoning"`
	ProvideExamples     bool `json:"provide_examples" yaml:"provide_examples"`
	CompareAlternatives bool `json:"compare_alternatives" yaml:"compare_alternatives"`
	HighlightRisks      bool `json:"highlight_risks" yaml:"highlight_risks"`

	// Interaction Style
	AskClarifyingQuestions bool `json:"ask_clarifying_questions" yaml:"ask_clarifying_questions"`
	OfferAlternatives      bool `json:"offer_alternatives" yaml:"offer_alternatives"`
	SuggestBestPractices   bool `json:"suggest_best_practices" yaml:"suggest_best_practices"`
}

// AgentKnowledge define o conhecimento especializado do agente
type AgentKnowledge struct {
	// Infrastructure Expertise
	TerraformExpertise string `json:"terraform_expertise" yaml:"terraform_expertise"` // beginner, intermediate, expert
	AWSExpertise       string `json:"aws_expertise" yaml:"aws_expertise"`
	AzureExpertise     string `json:"azure_expertise" yaml:"azure_expertise"`
	GCPExpertise       string `json:"gcp_expertise" yaml:"gcp_expertise"`

	// Domain Knowledge
	SecurityExpertise   string `json:"security_expertise" yaml:"security_expertise"`
	NetworkingExpertise string `json:"networking_expertise" yaml:"networking_expertise"`
	KubernetesExpertise string `json:"kubernetes_expertise" yaml:"kubernetes_expertise"`
	DatabaseExpertise   string `json:"database_expertise" yaml:"database_expertise"`

	// Specialized Knowledge
	ComplianceFrameworks []string `json:"compliance_frameworks" yaml:"compliance_frameworks"` // GDPR, SOC2, HIPAA, PCI-DSS
	IndustryFocus        []string `json:"industry_focus" yaml:"industry_focus"`               // fintech, healthcare, e-commerce
	ArchitecturePatterns []string `json:"architecture_patterns" yaml:"architecture_patterns"` // microservices, serverless, 3-tier

	// Custom Knowledge Base
	CustomRules      []CustomRule `json:"custom_rules" yaml:"custom_rules"`
	PreferredModules []string     `json:"preferred_modules" yaml:"preferred_modules"`
	BannedResources  []string     `json:"banned_resources" yaml:"banned_resources"`
	RequiredTags     []string     `json:"required_tags" yaml:"required_tags"`
}

// CustomRule é uma regra customizada do agente
type CustomRule struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Severity    string `json:"severity" yaml:"severity"` // critical, high, medium, low, info
	Pattern     string `json:"pattern" yaml:"pattern"`   // Regex or rule pattern
	Message     string `json:"message" yaml:"message"`
	Suggestion  string `json:"suggestion" yaml:"suggestion"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
}

// AgentLimits define os limites e restrições do agente
type AgentLimits struct {
	// Rate Limits
	MaxRequestsPerHour    int `json:"max_requests_per_hour" yaml:"max_requests_per_hour"`
	MaxRequestsPerDay     int `json:"max_requests_per_day" yaml:"max_requests_per_day"`
	MaxConcurrentRequests int `json:"max_concurrent_requests" yaml:"max_concurrent_requests"`

	// Token Limits
	MaxTokensPerRequest int `json:"max_tokens_per_request" yaml:"max_tokens_per_request"`
	MaxTokensPerDay     int `json:"max_tokens_per_day" yaml:"max_tokens_per_day"`

	// Analysis Limits
	MaxFilesPerAnalysis int `json:"max_files_per_analysis" yaml:"max_files_per_analysis"`
	MaxFileSizeMB       int `json:"max_file_size_mb" yaml:"max_file_size_mb"`
	MaxResourcesPerFile int `json:"max_resources_per_file" yaml:"max_resources_per_file"`

	// Cost Limits
	MaxCostPerRequest float64 `json:"max_cost_per_request" yaml:"max_cost_per_request"` // USD
	MaxCostPerDay     float64 `json:"max_cost_per_day" yaml:"max_cost_per_day"`         // USD
	MaxCostPerMonth   float64 `json:"max_cost_per_month" yaml:"max_cost_per_month"`     // USD

	// Time Limits
	MaxAnalysisTimeSeconds int `json:"max_analysis_time_seconds" yaml:"max_analysis_time_seconds"`
	RequestTimeoutSeconds  int `json:"request_timeout_seconds" yaml:"request_timeout_seconds"`
}

// AgentMetrics contém métricas de performance do agente
type AgentMetrics struct {
	// Usage Metrics
	TotalRequests      int64   `json:"total_requests" yaml:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests" yaml:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests" yaml:"failed_requests"`
	TotalTokensUsed    int64   `json:"total_tokens_used" yaml:"total_tokens_used"`
	TotalCostUSD       float64 `json:"total_cost_usd" yaml:"total_cost_usd"`

	// Performance Metrics
	AverageResponseTime     float64 `json:"average_response_time" yaml:"average_response_time"` // seconds
	AverageTokensPerRequest float64 `json:"average_tokens_per_request" yaml:"average_tokens_per_request"`
	AverageCostPerRequest   float64 `json:"average_cost_per_request" yaml:"average_cost_per_request"`

	// Quality Metrics
	AverageUserRating    float64 `json:"average_user_rating" yaml:"average_user_rating"`       // 0-5
	PositiveFeedbackRate float64 `json:"positive_feedback_rate" yaml:"positive_feedback_rate"` // 0-1
	IssuesDetected       int64   `json:"issues_detected" yaml:"issues_detected"`
	IssuesResolved       int64   `json:"issues_resolved" yaml:"issues_resolved"`

	// Time Metrics
	LastUsed    time.Time `json:"last_used" yaml:"last_used"`
	TotalUptime int64     `json:"total_uptime" yaml:"total_uptime"` // seconds
}

// AgentTemplate é um template pré-configurado de agente
type AgentTemplate struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Category    string `json:"category" yaml:"category"` // general, security, cost, architecture

	// Template Configuration
	DefaultConfig       AgentConfig       `json:"default_config" yaml:"default_config"`
	DefaultCapabilities AgentCapabilities `json:"default_capabilities" yaml:"default_capabilities"`
	DefaultPersonality  AgentPersonality  `json:"default_personality" yaml:"default_personality"`
	DefaultKnowledge    AgentKnowledge    `json:"default_knowledge" yaml:"default_knowledge"`
	DefaultLimits       AgentLimits       `json:"default_limits" yaml:"default_limits"`

	// Metadata
	IsRecommended bool     `json:"is_recommended" yaml:"is_recommended"`
	UseCases      []string `json:"use_cases" yaml:"use_cases"`
	Tags          []string `json:"tags" yaml:"tags"`
}

// CreateAgentRequest é a requisição para criar um novo agente
type CreateAgentRequest struct {
	TemplateID  string                 `json:"template_id" yaml:"template_id"` // Optional: use template
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Owner       string                 `json:"owner" yaml:"owner"`         // Wallet address
	Overrides   map[string]interface{} `json:"overrides" yaml:"overrides"` // Override template values
}

// AgentUpdateRequest é a requisição para atualizar um agente
type AgentUpdateRequest struct {
	Name         *string            `json:"name,omitempty" yaml:"name,omitempty"`
	Description  *string            `json:"description,omitempty" yaml:"description,omitempty"`
	Config       *AgentConfig       `json:"config,omitempty" yaml:"config,omitempty"`
	Capabilities *AgentCapabilities `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Personality  *AgentPersonality  `json:"personality,omitempty" yaml:"personality,omitempty"`
	Knowledge    *AgentKnowledge    `json:"knowledge,omitempty" yaml:"knowledge,omitempty"`
	Limits       *AgentLimits       `json:"limits,omitempty" yaml:"limits,omitempty"`
	Status       *string            `json:"status,omitempty" yaml:"status,omitempty"`
}

// AgentListResponse é a resposta com lista de agentes
type AgentListResponse struct {
	Agents     []Agent `json:"agents" yaml:"agents"`
	TotalCount int     `json:"total_count" yaml:"total_count"`
	Page       int     `json:"page" yaml:"page"`
	PageSize   int     `json:"page_size" yaml:"page_size"`
}

// AgentAnalysisRequest é uma requisição de análise com agente específico
type AgentAnalysisRequest struct {
	AgentID      string                 `json:"agent_id" yaml:"agent_id"`
	Code         string                 `json:"code" yaml:"code"`
	AnalysisType string                 `json:"analysis_type" yaml:"analysis_type"` // full, security, cost, etc
	Context      map[string]interface{} `json:"context,omitempty" yaml:"context,omitempty"`
}
