package models

// CheckovResult representa o resultado completo do Checkov
type CheckovResult struct {
	Summary       CheckovSummary `json:"summary"`
	Results       CheckovResults `json:"results"`
	CheckType     string         `json:"check_type"`
	ExecutionTime float64        `json:"execution_time"`
}

// CheckovSummary contém o resumo da análise Checkov
type CheckovSummary struct {
	Passed         int    `json:"passed"`
	Failed         int    `json:"failed"`
	Skipped        int    `json:"skipped"`
	ParsingErrors  int    `json:"parsing_errors"`
	ResourceCount  int    `json:"resource_count"`
	CheckovVersion string `json:"checkov_version"`
}

// CheckovResults contém os checks que passaram e falharam
type CheckovResults struct {
	PassedChecks  []CheckovCheck `json:"passed_checks"`
	FailedChecks  []CheckovCheck `json:"failed_checks"`
	SkippedChecks []CheckovCheck `json:"skipped_checks"`
}

// CheckovCheck representa um check individual do Checkov
type CheckovCheck struct {
	CheckID       string   `json:"check_id"`
	CheckName     string   `json:"check_name"`
	CheckResult   string   `json:"check_result"` // passed, failed, skipped
	Resource      string   `json:"resource"`
	File          string   `json:"file_path"`
	FileLineRange []int    `json:"file_line_range"`
	Guideline     string   `json:"guideline"`
	Description   string   `json:"description,omitempty"`
	Severity      string   `json:"severity,omitempty"` // CRITICAL, HIGH, MEDIUM, LOW
	CWE           []string `json:"cwe,omitempty"`
	OWASP         []string `json:"owasp,omitempty"`
	FixedCode     string   `json:"fixed_code,omitempty"`
}

// CheckovConfig representa configuração para execução do Checkov
type CheckovConfig struct {
	Directory         string   `json:"directory"`
	Framework         string   `json:"framework"` // terraform, cloudformation, etc
	Checks            []string `json:"checks,omitempty"`
	SkipChecks        []string `json:"skip_checks,omitempty"`
	CompactOutput     bool     `json:"compact_output"`
	Quiet             bool     `json:"quiet"`
	ExternalChecksDir string   `json:"external_checks_dir,omitempty"`
}
