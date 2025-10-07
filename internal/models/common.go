package models

import "time"

// AnalysisRequest representa uma requisição de análise de código
type AnalysisRequest struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
	Content    string `json:"content"`
	Branch     string `json:"branch,omitempty"`
	CommitSHA  string `json:"commit_sha,omitempty"`
}

// AnalysisResponse representa o resultado de uma análise
type AnalysisResponse struct {
	ID          string                 `json:"id"`
	Score       int                    `json:"score"`
	Analysis    AnalysisDetails        `json:"analysis"`
	Suggestions []Suggestion           `json:"suggestions"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// AnalysisDetails contém detalhes de todas as análises
type AnalysisDetails struct {
	Terraform TerraformAnalysis `json:"terraform"`
	Security  SecurityAnalysis  `json:"security"`
	IAM       IAMAnalysis       `json:"iam"`
	Cost      CostAnalysis      `json:"cost,omitempty"`
}

// Suggestion representa uma sugestão de melhoria
type Suggestion struct {
	Type             string                 `json:"type"`     // security, cost, best_practice, performance
	Severity         string                 `json:"severity"` // critical, high, medium, low, info
	Message          string                 `json:"message"`
	Recommendation   string                 `json:"recommendation"`
	File             string                 `json:"file"`
	Line             int                    `json:"line,omitempty"`
	Resource         string                 `json:"resource,omitempty"`
	EstimatedSavings string                 `json:"estimated_savings,omitempty"`
	ReferenceLink    string                 `json:"reference_link,omitempty"`
	References       []string               `json:"references,omitempty"`
	AutoFixAvailable bool                   `json:"auto_fix_available"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// SecurityAnalysis contém resultados de análise de segurança
type SecurityAnalysis struct {
	Critical     int               `json:"critical"`
	High         int               `json:"high"`
	Medium       int               `json:"medium"`
	Low          int               `json:"low"`
	Info         int               `json:"info"`
	TotalIssues  int               `json:"total_issues"`
	ChecksPassed int               `json:"checks_passed"`
	ChecksFailed int               `json:"checks_failed"`
	Findings     []SecurityFinding `json:"findings"`
}

// SecurityFinding representa um achado de segurança
type SecurityFinding struct {
	ID          string   `json:"id"`
	CheckID     string   `json:"check_id"`
	CheckName   string   `json:"check_name"`
	Severity    string   `json:"severity"`
	Resource    string   `json:"resource"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	Description string   `json:"description"`
	Guideline   string   `json:"guideline"`
	References  []string `json:"references"`
}

// IAMAnalysis contém análise de políticas IAM
type IAMAnalysis struct {
	TotalPolicies       int             `json:"total_policies"`
	TotalRoles          int             `json:"total_roles"`
	OverlyPermissive    bool            `json:"overly_permissive"`
	WildcardActions     []string        `json:"wildcard_actions"`
	AdminAccessDetected bool            `json:"admin_access_detected"`
	PublicAccess        []string        `json:"public_access"`
	Recommendations     []string        `json:"recommendations"`
	PrincipalRisks      []PrincipalRisk `json:"principal_risks"`
}

// PrincipalRisk representa um risco relacionado a principal IAM
type PrincipalRisk struct {
	Principal   string   `json:"principal"`
	Type        string   `json:"type"` // user, role, service
	RiskLevel   string   `json:"risk_level"`
	Permissions []string `json:"permissions"`
	Reason      string   `json:"reason"`
}

// CostAnalysis contém análise de custos
type CostAnalysis struct {
	EstimatedMonthlyCost  float64              `json:"estimated_monthly_cost"`
	Currency              string               `json:"currency"`
	OptimizationPotential float64              `json:"optimization_potential"`
	Recommendations       []CostRecommendation `json:"recommendations"`
}

// CostRecommendation representa uma recomendação de otimização de custo
type CostRecommendation struct {
	Resource                 string  `json:"resource"`
	CurrentCost              float64 `json:"current_cost"`
	PotentialSavings         float64 `json:"potential_savings"`
	Recommendation           string  `json:"recommendation"`
	ImplementationDifficulty string  `json:"implementation_difficulty"` // easy, medium, hard
}

// LLMRequest representa uma requisição para o LLM
type LLMRequest struct {
	Prompt          string    `json:"prompt"`
	SystemPrompt    string    `json:"system_prompt,omitempty"`
	ContextMessages []Message `json:"context_messages,omitempty"`
	Temperature     float64   `json:"temperature"`
	MaxTokens       int       `json:"max_tokens"`
	ResponseFormat  string    `json:"response_format,omitempty"`
}

// LLMResponse representa a resposta do LLM
type LLMResponse struct {
	Content    string                 `json:"content"`
	TokensUsed int                    `json:"tokens_used"`
	Model      string                 `json:"model"`
	Provider   string                 `json:"provider,omitempty"`
	LatencyMs  int64                  `json:"latency_ms,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Message represents a message in a conversation context
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
