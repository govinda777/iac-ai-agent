package models

// EnrichedSuggestion represents a single, detailed suggestion from the LLM.
type EnrichedSuggestion struct {
	OriginalID          string   `json:"original_id"`
	Type                string   `json:"type"`
	Severity            string   `json:"severity"`
	Title               string   `json:"title"`
	Message             string   `json:"message"`
	CodeExample         string   `json:"code_example"`
	ImplementationEffort string   `json:"implementation_effort"`
	EstimatedImpact     string   `json:"estimated_impact"`
	WhyItMatters        string   `json:"why_it_matters"`
	References          []string `json:"references"`
}