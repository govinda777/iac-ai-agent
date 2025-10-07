package models

// PlatformContext defines the specific context of the platform.
type PlatformContext struct {
	Name              string             `json:"name"`
	Standards         Standards          `json:"standards"`
	SupportedVersions SupportedVersions  `json:"supported_versions"`
	ApprovedModules   []ApprovedModule   `json:"approved_modules"`
}

// Standards defines policies and conventions.
type Standards struct {
	TaggingPolicy    map[string]string `json:"tagging_policy"`
	NamingConvention string            `json:"naming_convention"`
	RequiredOutputs  []string          `json:"required_outputs"`
}

// SupportedVersions defines the supported versions of tools.
type SupportedVersions struct {
	Terraform []string          `json:"terraform"`
	OpenTofu  []string          `json:"open_tofu"`
	Providers map[string]string `json:"providers"`
}

// ApprovedModule defines an approved Terraform module.
type ApprovedModule struct {
	Source      string `json:"source"`
	Version     string `json:"version"`
	UseCase     string `json:"use_case"`
	Recommended bool   `json:"recommended"`
}

// SecurityPolicy defines a security policy.
type SecurityPolicy struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	AutoFix     bool     `json:"auto_fix"`
	FixCode     string   `json:"fix_code"`
	Resources   []string `json:"resources"`
}

// BestPractice defines an IaC best practice.
type BestPractice struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Category    string   `json:"category"`
	Remediation string   `json:"remediation"`
	References  []string `json:"references"`
}