package models

// SecretFinding represents a secret that has been found in the code.
type SecretFinding struct {
	Type        string `json:"type"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Value       string `json:"value"` // This should be masked or redacted
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}