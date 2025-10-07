package models

// PreviewAnalysis representa a análise de um terraform plan
type PreviewAnalysis struct {
	PlannedChanges    []TerraformPlannedChange `json:"planned_changes"`
	Errors            []string                 `json:"errors"`
	Warnings          []string                 `json:"warnings"`
	ResourcesAffected int                      `json:"resources_affected"`
	CreateCount       int                      `json:"create_count"`
	UpdateCount       int                      `json:"update_count"`
	DestroyCount      int                      `json:"destroy_count"`
	ReplaceCount      int                      `json:"replace_count"`
	RiskLevel         string                   `json:"risk_level"` // low, medium, high, critical
}

// TerraformPlannedChange representa uma mudança planejada no terraform plan
type TerraformPlannedChange struct {
	Resource  string                 `json:"resource"`
	Action    string                 `json:"action"` // create, update, destroy, replace
	Changes   map[string]interface{} `json:"changes"`
	RiskScore int                    `json:"risk_score"`
	Impact    string                 `json:"impact"`
}

// TerraformPlan representa o resultado de um terraform plan em formato JSON
type TerraformPlan struct {
	FormatVersion    string `json:"format_version"`
	TerraformVersion string `json:"terraform_version"`
	ResourceChanges  []struct {
		Address string `json:"address"`
		Mode    string `json:"mode"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Change  struct {
			Actions []string               `json:"actions"`
			Before  map[string]interface{} `json:"before"`
			After   map[string]interface{} `json:"after"`
		} `json:"change"`
	} `json:"resource_changes"`
}

// RiskWarning representa um aviso de risco identificado no plan
type RiskWarning struct {
	Severity string `json:"severity"`
	Resource string `json:"resource"`
	Message  string `json:"message"`
	Action   string `json:"action"`
}
