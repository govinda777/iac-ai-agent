package models

// PreviewAnalysis representa análise de um terraform plan
type PreviewAnalysis struct {
	PlannedChanges    []PlannedChange `json:"planned_changes"`
	Errors            []string        `json:"errors"`
	Warnings          []string        `json:"warnings"`
	ResourcesAffected int             `json:"resources_affected"`
	CreateCount       int             `json:"create_count"`
	UpdateCount       int             `json:"update_count"`
	DestroyCount      int             `json:"destroy_count"`
	ReplaceCount      int             `json:"replace_count"`
	RiskLevel         string          `json:"risk_level"` // low, medium, high, critical
	RiskWarnings      []RiskWarning   `json:"risk_warnings"`
}

// PlannedChange representa uma mudança planejada em um recurso
type PlannedChange struct {
	Resource  string                 `json:"resource"`
	Action    string                 `json:"action"` // create, update, destroy, replace
	Changes   map[string]interface{} `json:"changes"`
	RiskScore int                    `json:"risk_score"`
	Impact    string                 `json:"impact"`
}

// RiskWarning representa um aviso de risco
type RiskWarning struct {
	Severity string `json:"severity"`
	Resource string `json:"resource"`
	Message  string `json:"message"`
	Action   string `json:"action"`
}

// TerraformPlan representa a estrutura de um terraform plan JSON
type TerraformPlan struct {
	FormatVersion    string                    `json:"format_version"`
	TerraformVersion string                    `json:"terraform_version"`
	ResourceChanges  []TerraformResourceChange `json:"resource_changes"`
}

// TerraformResourceChange representa uma mudança de recurso no plan
type TerraformResourceChange struct {
	Address string `json:"address"`
	Mode    string `json:"mode"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Change  struct {
		Actions []string               `json:"actions"`
		Before  map[string]interface{} `json:"before"`
		After   map[string]interface{} `json:"after"`
	} `json:"change"`
}
