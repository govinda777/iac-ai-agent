package models

import "time"

// LLMStructuredResponse é o template principal de resposta estruturada da LLM
type LLMStructuredResponse struct {
	ID                    string                `json:"id"`
	Timestamp             time.Time             `json:"timestamp"`
	ModelUsed             string                `json:"model_used"`
	TokensUsed            int                   `json:"tokens_used"`
	AnalysisType          string                `json:"analysis_type"` // full, preview, security, cost, etc
	ExecutiveSummary      ExecutiveSummary      `json:"executive_summary"`
	CriticalIssues        []EnrichedIssue       `json:"critical_issues"`
	Improvements          []EnrichedImprovement `json:"improvements"`
	BestPractices         []BestPracticeCheck   `json:"best_practices"`
	ArchitecturalInsights ArchitecturalInsights `json:"architectural_insights"`
	PriorityActions       []PriorityAction      `json:"priority_actions"`
	QuickWins             []QuickWin            `json:"quick_wins"`
	Metadata              ResponseMetadata      `json:"metadata"`
}

// ExecutiveSummary é o resumo executivo da análise
type ExecutiveSummary struct {
	OverallScore       int      `json:"overall_score"`        // 0-100
	ScoreLevel         string   `json:"score_level"`          // excellent, good, fair, poor, critical
	Summary            string   `json:"summary"`              // 2-3 sentences
	MainConcerns       []string `json:"main_concerns"`        // Top 3-5 concerns
	KeyStrengths       []string `json:"key_strengths"`        // What's good
	Recommendation     string   `json:"recommendation"`       // approve, request_changes, needs_review
	EstimatedRiskLevel string   `json:"estimated_risk_level"` // low, medium, high, critical
}

// EnrichedIssue é um problema identificado com contexto da LLM
type EnrichedIssue struct {
	ID                   string      `json:"id"`
	Title                string      `json:"title"`
	Category             string      `json:"category"` // security, preview_error, iam, drift, timeout
	Severity             string      `json:"severity"` // critical, high, medium, low
	Resource             string      `json:"resource,omitempty"`
	File                 string      `json:"file,omitempty"`
	Line                 int         `json:"line,omitempty"`
	Description          string      `json:"description"`
	BusinessImpact       string      `json:"business_impact"` // Why it matters
	TechnicalExplanation string      `json:"technical_explanation"`
	HowToFix             HowToFix    `json:"how_to_fix"`
	EstimatedEffort      string      `json:"estimated_effort"` // easy, medium, hard
	References           []Reference `json:"references"`
	ComplianceImpact     []string    `json:"compliance_impact,omitempty"` // GDPR, SOC2, HIPAA, etc
	RelatedIssues        []string    `json:"related_issues,omitempty"`    // IDs of related issues
}

// EnrichedImprovement é uma melhoria sugerida com contexto da LLM
type EnrichedImprovement struct {
	ID                       string      `json:"id"`
	Title                    string      `json:"title"`
	Category                 string      `json:"category"` // cost, architecture, modules, secrets, provider
	Priority                 string      `json:"priority"` // high, medium, low
	CurrentState             string      `json:"current_state"`
	ProposedState            string      `json:"proposed_state"`
	Benefits                 []string    `json:"benefits"`
	Implementation           HowToFix    `json:"implementation"`
	EstimatedSavings         *Savings    `json:"estimated_savings,omitempty"`
	EstimatedEffort          string      `json:"estimated_effort"`            // easy, medium, hard
	EstimatedTimeToImplement string      `json:"estimated_time_to_implement"` // 1h, 1d, 1w, etc
	ROI                      string      `json:"roi,omitempty"`               // Return on Investment description
	References               []Reference `json:"references"`
}

// BestPracticeCheck é uma verificação de boa prática
type BestPracticeCheck struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Category       string      `json:"category"` // structure, documentation, testing, size, etc
	Status         string      `json:"status"`   // passed, warning, failed
	CurrentValue   string      `json:"current_value,omitempty"`
	ExpectedValue  string      `json:"expected_value,omitempty"`
	Message        string      `json:"message"`
	Recommendation string      `json:"recommendation,omitempty"`
	Priority       string      `json:"priority"` // high, medium, low
	References     []Reference `json:"references,omitempty"`
}

// ArchitecturalInsights são insights sobre a arquitetura
type ArchitecturalInsights struct {
	DetectedPattern     string               `json:"detected_pattern"` // 3-tier, serverless, microservices, etc
	Confidence          float64              `json:"confidence"`       // 0.0 - 1.0
	Strengths           []string             `json:"strengths"`
	Weaknesses          []string             `json:"weaknesses"`
	AreasForImprovement []string             `json:"areas_for_improvement"`
	RecommendedPatterns []RecommendedPattern `json:"recommended_patterns"`
	ScalabilityConcerns []string             `json:"scalability_concerns,omitempty"`
	ResilienceConcerns  []string             `json:"resilience_concerns,omitempty"`
}

// RecommendedPattern é um padrão arquitetural recomendado
type RecommendedPattern struct {
	Name      string   `json:"name"`
	Reason    string   `json:"reason"`
	Benefits  []string `json:"benefits"`
	Tradeoffs []string `json:"tradeoffs,omitempty"`
	Effort    string   `json:"effort"` // easy, medium, hard
	Reference string   `json:"reference,omitempty"`
}

// PriorityAction é uma ação prioritária
type PriorityAction struct {
	Order       int      `json:"order"` // 1, 2, 3, etc
	Action      string   `json:"action"`
	Reason      string   `json:"reason"`
	IssuesFixed []string `json:"issues_fixed"` // IDs of issues this fixes
	Effort      string   `json:"effort"`       // easy, medium, hard
	Impact      string   `json:"impact"`       // high, medium, low
}

// QuickWin é uma vitória rápida (high impact, low effort)
type QuickWin struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Implementation HowToFix `json:"implementation"`
	Benefit        string   `json:"benefit"`
	EstimatedTime  string   `json:"estimated_time"` // 5m, 15m, 1h, etc
	ROI            string   `json:"roi"`            // Very High, High, Medium
}

// HowToFix descreve como corrigir um problema
type HowToFix struct {
	Steps           []string       `json:"steps"`
	CodeExample     *CodeExample   `json:"code_example,omitempty"`
	CommandsToRun   []string       `json:"commands_to_run,omitempty"`
	ConfigChanges   []ConfigChange `json:"config_changes,omitempty"`
	ValidationSteps []string       `json:"validation_steps,omitempty"`
	RollbackPlan    string         `json:"rollback_plan,omitempty"`
}

// CodeExample é um exemplo de código
type CodeExample struct {
	Language string `json:"language"` // hcl, go, bash, yaml, etc
	Before   string `json:"before,omitempty"`
	After    string `json:"after"`
	Diff     string `json:"diff,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

// ConfigChange descreve uma mudança de configuração
type ConfigChange struct {
	File     string `json:"file"`
	Key      string `json:"key"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value"`
	Reason   string `json:"reason"`
}

// Reference é uma referência externa
type Reference struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"` // documentation, blog, video, github, etc
}

// Savings descreve economia estimada
type Savings struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Period      string  `json:"period"`               // monthly, yearly
	Percentage  float64 `json:"percentage,omitempty"` // % of current cost
	Confidence  string  `json:"confidence"`           // low, medium, high
	Description string  `json:"description"`
}

// ResponseMetadata contém metadados da resposta
type ResponseMetadata struct {
	AnalysisDuration     string                 `json:"analysis_duration"`
	LLMCallCount         int                    `json:"llm_call_count"`
	KnowledgeBaseEntries int                    `json:"knowledge_base_entries"`
	RulesApplied         int                    `json:"rules_applied"`
	Platform             string                 `json:"platform,omitempty"`    // Nation.fun
	Environment          string                 `json:"environment,omitempty"` // prod, staging, dev
	CustomData           map[string]interface{} `json:"custom_data,omitempty"`
}

// ===== TEMPLATES ESPECÍFICOS POR TIPO DE ANÁLISE =====

// PreviewAnalysisResponse é a resposta específica para análise de preview
type PreviewAnalysisResponse struct {
	LLMStructuredResponse
	PreviewDetails PreviewDetails `json:"preview_details"`
}

// DangerousOperation representa uma operação perigosa detectada
type DangerousOperation struct {
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Mitigation  string `json:"mitigation"`
}

// PreviewDetails contém detalhes do preview analisado
type PreviewDetails struct {
	PlannedChanges       []PlannedChange      `json:"planned_changes"`
	ResourcesAffected    int                  `json:"resources_affected"`
	CreateCount          int                  `json:"create_count"`
	UpdateCount          int                  `json:"update_count"`
	DestroyCount         int                  `json:"destroy_count"`
	ReplaceCount         int                  `json:"replace_count"`
	RiskLevel            string               `json:"risk_level"` // low, medium, high, critical
	DangerousOperations  []DangerousOperation `json:"dangerous_operations"`
	EstimatedApplyTime   string               `json:"estimated_apply_time"`
	RecommendedApprovals []string             `json:"recommended_approvals,omitempty"`
}

// SecurityAuditResponse é a resposta específica para auditoria de segurança
type SecurityAuditResponse struct {
	LLMStructuredResponse
	SecurityDetails SecurityAuditDetails `json:"security_details"`
}

// SecurityAuditDetails contém detalhes da auditoria de segurança
type SecurityAuditDetails struct {
	OverallRiskLevel       string                 `json:"overall_risk_level"` // low, medium, high, critical
	ComplianceStatus       ComplianceStatus       `json:"compliance_status"`
	VulnerabilityBreakdown VulnerabilityBreakdown `json:"vulnerability_breakdown"`
	Top5Vulnerabilities    []EnrichedIssue        `json:"top5_vulnerabilities"`
	SecurityRoadmap        SecurityRoadmap        `json:"security_roadmap"`
	SecretsDetected        []SecretDetection      `json:"secrets_detected,omitempty"`
}

// ComplianceStatus representa o status de conformidade
type ComplianceStatus struct {
	Frameworks        map[string]FrameworkCompliance `json:"frameworks"`         // GDPR, SOC2, HIPAA, PCI-DSS
	OverallCompliance float64                        `json:"overall_compliance"` // 0-100%
}

// FrameworkCompliance representa conformidade com um framework
type FrameworkCompliance struct {
	Name            string   `json:"name"`
	ComplianceLevel float64  `json:"compliance_level"` // 0-100%
	Status          string   `json:"status"`           // compliant, partial, non_compliant
	Gaps            []string `json:"gaps"`
	Recommendations []string `json:"recommendations"`
}

// VulnerabilityBreakdown é o breakdown de vulnerabilidades
type VulnerabilityBreakdown struct {
	BySeverity    map[string]int `json:"by_severity"`              // critical, high, medium, low
	ByCategory    map[string]int `json:"by_category"`              // network, iam, encryption, etc
	ByResource    map[string]int `json:"by_resource"`              // aws_s3_bucket, aws_iam_role, etc
	TrendAnalysis string         `json:"trend_analysis,omitempty"` // Better/Worse than previous
}

// SecurityRoadmap é o roadmap de segurança
type SecurityRoadmap struct {
	UrgentActions     []PriorityAction `json:"urgent_actions"`      // Fix now
	ShortTermActions  []PriorityAction `json:"short_term_actions"`  // Fix in 1-2 weeks
	MediumTermActions []PriorityAction `json:"medium_term_actions"` // Fix in 1-3 months
}

// SecretDetection representa um secret detectado
type SecretDetection struct {
	Type        string `json:"type"` // aws_key, password, private_key, api_key
	File        string `json:"file"`
	Line        int    `json:"line"`
	Severity    string `json:"severity"`
	MaskedValue string `json:"masked_value"`
	Suggestion  string `json:"suggestion"`
}

// CostOptimizationResponse é a resposta específica para otimização de custo
type CostOptimizationResponse struct {
	LLMStructuredResponse
	CostDetails CostOptimizationDetails `json:"cost_details"`
}

// CostOptimizationDetails contém detalhes da otimização de custo
type CostOptimizationDetails struct {
	CurrentMonthlyCost            float64                         `json:"current_monthly_cost"`
	OptimizedMonthlyCost          float64                         `json:"optimized_monthly_cost"`
	TotalSavingsPotential         float64                         `json:"total_savings_potential"`
	Currency                      string                          `json:"currency"`
	OptimizationsByCategory       map[string]CategoryOptimization `json:"optimizations_by_category"`
	ReservedInstanceOpportunities []ReservedInstanceSuggestion    `json:"reserved_instance_opportunities,omitempty"`
	UnusedResources               []UnusedResource                `json:"unused_resources,omitempty"`
	RightsizingOpportunities      []RightsizingSuggestion         `json:"rightsizing_opportunities,omitempty"`
}

// CategoryOptimization é otimização por categoria
type CategoryOptimization struct {
	Category         string  `json:"category"` // compute, storage, network, database
	CurrentCost      float64 `json:"current_cost"`
	PotentialSavings float64 `json:"potential_savings"`
	Opportunities    int     `json:"opportunities"`
}

// ReservedInstanceSuggestion sugere Reserved Instances
type ReservedInstanceSuggestion struct {
	ResourceType         string  `json:"resource_type"`
	CurrentOnDemandCost  float64 `json:"current_on_demand_cost"`
	ReservedInstanceCost float64 `json:"reserved_instance_cost"`
	AnnualSavings        float64 `json:"annual_savings"`
	Term                 string  `json:"term"`           // 1yr, 3yr
	PaymentOption        string  `json:"payment_option"` // all_upfront, partial_upfront, no_upfront
}

// UnusedResource representa um recurso não utilizado
type UnusedResource struct {
	Resource       string  `json:"resource"`
	Type           string  `json:"type"`
	MonthlyCost    float64 `json:"monthly_cost"`
	LastUsed       string  `json:"last_used,omitempty"`
	Recommendation string  `json:"recommendation"`
}

// RightsizingSuggestion sugere ajuste de tamanho de recurso
type RightsizingSuggestion struct {
	Resource          string  `json:"resource"`
	CurrentType       string  `json:"current_type"`
	SuggestedType     string  `json:"suggested_type"`
	CurrentCost       float64 `json:"current_cost"`
	NewCost           float64 `json:"new_cost"`
	MonthlySavings    float64 `json:"monthly_savings"`
	PerformanceImpact string  `json:"performance_impact"`
	Utilization       string  `json:"utilization"` // CPU/Memory usage stats
}
