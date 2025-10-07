package models

// TerraformAnalysis contém resultados da análise do código Terraform
type TerraformAnalysis struct {
	Valid                bool                  `json:"valid"`
	TotalResources       int                   `json:"total_resources"`
	TotalModules         int                   `json:"total_modules"`
	TotalVariables       int                   `json:"total_variables"`
	TotalOutputs         int                   `json:"total_outputs"`
	TotalDataSources     int                   `json:"total_data_sources"`
	Providers            []string              `json:"providers"`
	Resources            []TerraformResource   `json:"resources"`
	Modules              []TerraformModule     `json:"modules"`
	Variables            []TerraformVariable   `json:"variables"`
	Outputs              []TerraformOutput     `json:"outputs"`
	SyntaxErrors         []SyntaxError         `json:"syntax_errors,omitempty"`
	BestPracticeWarnings []string              `json:"best_practice_warnings,omitempty"`
	ResourceGraph        *ResourceGraph        `json:"resource_graph,omitempty"`
}

// TerraformResource representa um recurso Terraform
type TerraformResource struct {
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	Provider     string                 `json:"provider"`
	File         string                 `json:"file"`
	LineStart    int                    `json:"line_start"`
	LineEnd      int                    `json:"line_end"`
	Attributes   map[string]interface{} `json:"attributes"`
	Dependencies []string               `json:"dependencies,omitempty"`
	Tags         map[string]string      `json:"tags,omitempty"`
}

// TerraformModule representa um módulo Terraform
type TerraformModule struct {
	Name         string                 `json:"name"`
	Source       string                 `json:"source"`
	Version      string                 `json:"version,omitempty"`
	File         string                 `json:"file"`
	LineStart    int                    `json:"line_start"`
	Inputs       map[string]interface{} `json:"inputs"`
	Verified     bool                   `json:"verified"`
}

// TerraformVariable representa uma variável Terraform
type TerraformVariable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Required    bool        `json:"required"`
	Sensitive   bool        `json:"sensitive"`
	File        string      `json:"file"`
}

// TerraformOutput representa um output Terraform
type TerraformOutput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value"`
	Sensitive   bool   `json:"sensitive"`
	File        string `json:"file"`
}

// SyntaxError representa um erro de sintaxe
type SyntaxError struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Snippet string `json:"snippet,omitempty"`
}

// ResourceGraph representa o grafo de dependências entre recursos
type ResourceGraph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// GraphNode representa um nó no grafo
type GraphNode struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Resource string `json:"resource"`
}

// GraphEdge representa uma aresta no grafo (dependência)
type GraphEdge struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Type       string `json:"type"` // explicit, implicit
}
