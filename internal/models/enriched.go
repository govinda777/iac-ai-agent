package models

// EnrichedSuggestion representa uma sugestão enriquecida pelo LLM
type EnrichedSuggestion struct {
	OriginalID           string   `json:"original_id,omitempty"`
	Type                 string   `json:"type"`
	Severity             string   `json:"severity"`
	Title                string   `json:"title"`
	Message              string   `json:"message"`
	Resource             string   `json:"resource,omitempty"`
	CodeExample          string   `json:"code_example"`
	ImplementationEffort string   `json:"implementation_effort"`
	EstimatedImpact      string   `json:"estimated_impact"`
	WhyItMatters         string   `json:"why_it_matters"`
	References           []string `json:"references,omitempty"`
}

// ArchitecturalInsights representa insights arquiteturais do LLM
type ArchitecturalInsights struct {
	PatternDetected       string   `json:"pattern_detected"`
	Strengths             []string `json:"strengths"`
	AreasForImprovement   []string `json:"areas_for_improvement"`
	RecommendedModules    []string `json:"recommended_modules,omitempty"`
	ArchitectureScore     int      `json:"architecture_score,omitempty"`
	ScalabilityAssessment string   `json:"scalability_assessment,omitempty"`
	SecurityAssessment    string   `json:"security_assessment,omitempty"`
	CostAssessment        string   `json:"cost_assessment,omitempty"`
}

// EnrichmentResponse representa a resposta completa do LLM para enriquecimento
type EnrichmentResponse struct {
	EnrichedSuggestions   []EnrichedSuggestion  `json:"enriched_suggestions"`
	ArchitecturalInsights ArchitecturalInsights `json:"architectural_insights,omitempty"`
	PriorityActions       []string              `json:"priority_actions,omitempty"`
}

// PreviewAnalysisResponse representa a resposta do LLM para análise de preview
type PreviewAnalysisResponse struct {
	RiskAssessment struct {
		OverallRisk     string `json:"overall_risk"`
		HighRiskChanges []struct {
			Resource    string `json:"resource"`
			Action      string `json:"action"`
			RiskDetails string `json:"risk_details"`
			Mitigation  string `json:"mitigation"`
		} `json:"high_risk_changes"`
		BlastRadius string `json:"blast_radius"`
	} `json:"risk_assessment"`
	ImpactAnalysis struct {
		EstimatedDowntime      string   `json:"estimated_downtime"`
		AffectedDependencies   []string `json:"affected_dependencies"`
		PerformanceImpact      string   `json:"performance_impact"`
		DataLossRisk           string   `json:"data_loss_risk,omitempty"`
		BusinessContinuityRisk string   `json:"business_continuity_risk,omitempty"`
	} `json:"impact_analysis"`
	Recommendations struct {
		ApplyStrategy   string   `json:"apply_strategy"`
		TestingStrategy string   `json:"testing_strategy"`
		RollbackPlan    string   `json:"rollback_plan"`
		PreApplyChecks  []string `json:"pre_apply_checks,omitempty"`
		PostApplyChecks []string `json:"post_apply_checks,omitempty"`
	} `json:"recommendations"`
	BestPractices []struct {
		Issue          string `json:"issue"`
		Recommendation string `json:"recommendation"`
	} `json:"best_practices"`
}

// SecurityAuditResponse representa a resposta do LLM para auditoria de segurança
type SecurityAuditResponse struct {
	SecurityAnalysis struct {
		CriticalFindings []SecurityFindingDetail `json:"critical_findings"`
		HighFindings     []SecurityFindingDetail `json:"high_findings"`
		MediumFindings   []SecurityFindingDetail `json:"medium_findings"`
		LowFindings      []SecurityFindingDetail `json:"low_findings"`
	} `json:"security_analysis"`
	ImplementationPlan struct {
		ImmediateActions []string `json:"immediate_actions"`
		ShortTermActions []string `json:"short_term_actions"`
		LongTermActions  []string `json:"long_term_actions"`
	} `json:"implementation_plan"`
	SecurityPostureRecommendations []struct {
		Area           string `json:"area"`
		Recommendation string `json:"recommendation"`
		Implementation string `json:"implementation"`
	} `json:"security_posture_recommendations"`
}

// SecurityFindingDetail representa um achado de segurança detalhado pelo LLM
type SecurityFindingDetail struct {
	FindingID              string   `json:"finding_id"`
	Title                  string   `json:"title"`
	Description            string   `json:"description"`
	ExploitScenario        string   `json:"exploit_scenario"`
	IsFalsePositive        bool     `json:"is_false_positive"`
	RemediationCode        string   `json:"remediation_code"`
	RemediationExplanation string   `json:"remediation_explanation"`
	AdditionalControls     []string `json:"additional_controls"`
}

// CostOptimizationResponse representa a resposta do LLM para otimização de custos
type CostOptimizationResponse struct {
	CostOptimization struct {
		TotalEstimatedSavings string `json:"total_estimated_savings"`
		Opportunities         []struct {
			Resource             string   `json:"resource"`
			CurrentConfiguration string   `json:"current_configuration"`
			Recommendation       string   `json:"recommendation"`
			EstimatedSavings     string   `json:"estimated_savings"`
			ImplementationEffort string   `json:"implementation_effort"`
			ROI                  string   `json:"roi"`
			CodeExample          string   `json:"code_example"`
			Risks                []string `json:"risks"`
		} `json:"opportunities"`
	} `json:"cost_optimization"`
	ArchitectureRecommendations []struct {
		Area                     string `json:"area"`
		CurrentApproach          string `json:"current_approach"`
		RecommendedApproach      string `json:"recommended_approach"`
		EstimatedSavings         string `json:"estimated_savings"`
		ImplementationComplexity string `json:"implementation_complexity"`
	} `json:"architecture_recommendations"`
	QuickWins []struct {
		Action  string `json:"action"`
		Savings string `json:"savings"`
		Effort  string `json:"effort"`
	} `json:"quick_wins"`
}

// BestPractice representa uma prática recomendada da Knowledge Base
type BestPractice struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Severity    string   `json:"severity"`
	Resources   []string `json:"resources,omitempty"`
	References  []string `json:"references,omitempty"`
}

// Module representa um módulo recomendado da Knowledge Base
type Module struct {
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	UseCase     string   `json:"use_case"`
	Recommended bool     `json:"recommended"`
	Resources   []string `json:"resources,omitempty"`
}

// SecurityFinding representa um achado de segurança
type SecurityFinding struct {
	Type        string `json:"type"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Value       string `json:"value"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
	Resource    string `json:"resource,omitempty"`
}

// SecretFinding representa um achado de secret
type SecretFinding struct {
	Type        string `json:"type"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Value       string `json:"value"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}
