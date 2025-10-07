package models

import (
	"encoding/json"
)

// PreviewReview representa uma análise de preview de terraform
type PreviewReview struct {
	RiskAssessment struct {
		OverallRisk    string `json:"overall_risk"`
		HighRiskChanges []struct {
			Resource    string `json:"resource"`
			Action      string `json:"action"`
			RiskDetails string `json:"risk_details"`
			Mitigation  string `json:"mitigation"`
		} `json:"high_risk_changes"`
		BlastRadius string `json:"blast_radius"`
	} `json:"risk_assessment"`
	ImpactAnalysis struct {
		EstimatedDowntime    string   `json:"estimated_downtime"`
		AffectedDependencies []string `json:"affected_dependencies"`
		PerformanceImpact    string   `json:"performance_impact"`
	} `json:"impact_analysis"`
	Recommendations struct {
		ApplyStrategy   string `json:"apply_strategy"`
		TestingStrategy string `json:"testing_strategy"`
		RollbackPlan    string `json:"rollback_plan"`
	} `json:"recommendations"`
	BestPractices []struct {
		Issue         string `json:"issue"`
		Recommendation string `json:"recommendation"`
	} `json:"best_practices"`
}

// SecurityAnalysisResponse representa uma análise de segurança detalhada
type SecurityAnalysisResponse struct {
	SecurityAnalysis struct {
		CriticalFindings []SecurityFindingDetail `json:"critical_findings"`
		HighFindings     []SecurityFindingDetail `json:"high_findings"`
		MediumFindings   []SecurityFindingDetail `json:"medium_findings"`
		LowFindings      []SecurityFindingDetail `json:"low_findings"`
	} `json:"security_analysis"`
	ImplementationPlan struct {
		ImmediateActions  []string `json:"immediate_actions"`
		ShortTermActions  []string `json:"short_term_actions"`
		LongTermActions   []string `json:"long_term_actions"`
	} `json:"implementation_plan"`
	SecurityPostureRecommendations []struct {
		Area           string `json:"area"`
		Recommendation string `json:"recommendation"`
		Implementation string `json:"implementation"`
	} `json:"security_posture_recommendations"`
}

// SecurityFindingDetail representa um detalhe de uma vulnerabilidade de segurança
type SecurityFindingDetail struct {
	FindingID            string   `json:"finding_id"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	ExploitScenario      string   `json:"exploit_scenario"`
	IsFalsePositive      bool     `json:"is_false_positive"`
	RemediationCode      string   `json:"remediation_code"`
	RemediationExplanation string `json:"remediation_explanation"`
	AdditionalControls   []string `json:"additional_controls"`
}

// CostOptimization representa recomendações de otimização de custo
type CostOptimization struct {
	CostOptimization struct {
		TotalEstimatedSavings string `json:"total_estimated_savings"`
		Opportunities []struct {
			Resource            string   `json:"resource"`
			CurrentConfiguration string   `json:"current_configuration"`
			Recommendation      string   `json:"recommendation"`
			EstimatedSavings    string   `json:"estimated_savings"`
			ImplementationEffort string   `json:"implementation_effort"`
			ROI                 string   `json:"roi"`
			CodeExample         string   `json:"code_example"`
			Risks               []string `json:"risks"`
		} `json:"opportunities"`
	} `json:"cost_optimization"`
	ArchitectureRecommendations []struct {
		Area                    string `json:"area"`
		CurrentApproach         string `json:"current_approach"`
		RecommendedApproach     string `json:"recommended_approach"`
		EstimatedSavings        string `json:"estimated_savings"`
		ImplementationComplexity string `json:"implementation_complexity"`
	} `json:"architecture_recommendations"`
	QuickWins []struct {
		Action  string `json:"action"`
		Savings string `json:"savings"`
		Effort  string `json:"effort"`
	} `json:"quick_wins"`
}

// EnrichedSuggestionResponse representa a resposta com sugestões enriquecidas
type EnrichedSuggestionResponse struct {
	EnrichedSuggestions []EnrichedSuggestionLLM `json:"enriched_suggestions"`
	ArchitecturalInsights struct {
		PatternDetected      string   `json:"pattern_detected"`
		Strengths            []string `json:"strengths"`
		AreasForImprovement  []string `json:"areas_for_improvement"`
	} `json:"architectural_insights"`
	PriorityActions []string `json:"priority_actions"`
}

// EnrichedSuggestionLLM representa uma sugestão enriquecida pelo LLM na resposta
type EnrichedSuggestionLLM struct {
	OriginalID         string   `json:"original_id"`
	Type               string   `json:"type"`
	Severity           string   `json:"severity"`
	Title              string   `json:"title"`
	Message            string   `json:"message"`
	CodeExample        string   `json:"code_example"`
	ImplementationEffort string   `json:"implementation_effort"`
	EstimatedImpact    string   `json:"estimated_impact"`
	WhyItMatters       string   `json:"why_it_matters"`
	References         []string `json:"references"`
}

// ParseEnrichedSuggestions analisa a resposta do LLM para extrair sugestões enriquecidas
func ParseEnrichedSuggestions(response string) ([]Suggestion, error) {
	var enrichedResponse EnrichedSuggestionResponse
	err := json.Unmarshal([]byte(response), &enrichedResponse)
	if err != nil {
		return nil, err
	}

	suggestions := make([]Suggestion, 0, len(enrichedResponse.EnrichedSuggestions))
	for _, es := range enrichedResponse.EnrichedSuggestions {
		suggestion := Suggestion{
			Type:           es.Type,
			Severity:       es.Severity,
			Message:        es.Message,
			Recommendation: es.CodeExample,
			Metadata:       map[string]interface{}{
				"implementation_effort": es.ImplementationEffort,
				"estimated_impact":      es.EstimatedImpact,
				"why_it_matters":       es.WhyItMatters,
				"references":           es.References,
			},
			References: es.References,
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// ParsePreviewReview analisa a resposta do LLM para extrair uma análise de preview
func ParsePreviewReview(response string) (*PreviewReview, error) {
	var previewReview PreviewReview
	err := json.Unmarshal([]byte(response), &previewReview)
	if err != nil {
		return nil, err
	}
	return &previewReview, nil
}

// ParseSecurityAnalysis analisa a resposta do LLM para extrair uma análise de segurança
func ParseSecurityAnalysis(response string) (*SecurityAnalysisResponse, error) {
	var securityAnalysis SecurityAnalysisResponse
	err := json.Unmarshal([]byte(response), &securityAnalysis)
	if err != nil {
		return nil, err
	}
	return &securityAnalysis, nil
}

// ParseCostOptimization analisa a resposta do LLM para extrair recomendações de otimização de custo
func ParseCostOptimization(response string) (*CostOptimization, error) {
	var costOptimization CostOptimization
	err := json.Unmarshal([]byte(response), &costOptimization)
	if err != nil {
		return nil, err
	}
	return &costOptimization, nil
}
