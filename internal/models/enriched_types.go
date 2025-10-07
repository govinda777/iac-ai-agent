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
