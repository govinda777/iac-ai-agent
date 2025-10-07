package models

// PreviewAnalysis holds the results of analyzing a Terraform plan.
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
}

// TerraformPlan is a struct for unmarshalling the JSON output of a Terraform plan.
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

// RiskWarning represents a high-risk change detected in a plan.
type RiskWarning struct {
    Severity string
    Resource string
    Message  string
    Action   string
}