package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var terraformFunctions = map[string]function.Function{
	"jsonencode": function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name:             "val",
				Type:             cty.DynamicPseudoType,
				AllowNull:        true,
				AllowUnknown:     true,
				AllowDynamicType: true,
			},
		},
		Type: function.StaticReturnType(cty.String),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			return stdlib.JSONEncode(args[0])
		},
	}),
}

// TerraformAnalyzer realiza análise de código Terraform
type TerraformAnalyzer struct {
	parser *hclparse.Parser
}

// NewTerraformAnalyzer cria uma nova instância do analisador
func NewTerraformAnalyzer() *TerraformAnalyzer {
	return &TerraformAnalyzer{
		parser: hclparse.NewParser(),
	}
}

// AnalyzeDirectory analisa todos os arquivos Terraform em um diretório
func (ta *TerraformAnalyzer) AnalyzeDirectory(dir string) (*models.TerraformAnalysis, error) {
	analysis := &models.TerraformAnalysis{
		Valid:                true,
		Resources:            []models.TerraformResource{},
		Modules:              []models.TerraformModule{},
		Variables:            []models.TerraformVariable{},
		Outputs:              []models.TerraformOutput{},
		Providers:            []string{},
		SyntaxErrors:         []models.SyntaxError{},
		BestPracticeWarnings: []string{},
	}

	// Busca arquivos .tf no diretório
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".tf") {
			return nil
		}

		// Analisa o arquivo
		fileAnalysis, err := ta.analyzeFile(path)
		if err != nil {
			// Adiciona erro de sintaxe
			analysis.Valid = false
			analysis.SyntaxErrors = append(analysis.SyntaxErrors, models.SyntaxError{
				File:    path,
				Message: err.Error(),
			})
			return nil // Continua processando outros arquivos
		}

		// Merge resultados
		ta.mergeAnalysis(analysis, fileAnalysis)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao percorrer diretório: %w", err)
	}

	// Calcula totais
	analysis.TotalResources = len(analysis.Resources)
	analysis.TotalModules = len(analysis.Modules)
	analysis.TotalVariables = len(analysis.Variables)
	analysis.TotalOutputs = len(analysis.Outputs)

	// Verifica best practices
	ta.checkBestPractices(analysis)

	return analysis, nil
}

// AnalyzeContent analisa conteúdo Terraform direto
func (ta *TerraformAnalyzer) AnalyzeContent(content string, filename string) (*models.TerraformAnalysis, error) {
	analysis := &models.TerraformAnalysis{
		Valid:                true,
		Resources:            []models.TerraformResource{},
		Modules:              []models.TerraformModule{},
		Variables:            []models.TerraformVariable{},
		Outputs:              []models.TerraformOutput{},
		Providers:            []string{},
		SyntaxErrors:         []models.SyntaxError{},
		BestPracticeWarnings: []string{},
	}

	file, diags := ta.parser.ParseHCL([]byte(content), filename)
	if diags.HasErrors() {
		analysis.Valid = false
		for _, diag := range diags {
			analysis.SyntaxErrors = append(analysis.SyntaxErrors, models.SyntaxError{
				File:    filename,
				Line:    diag.Subject.Start.Line,
				Column:  diag.Subject.Start.Column,
				Message: diag.Summary,
				Snippet: diag.Detail,
			})
		}
		return analysis, nil
	}

	// Parse resources, modules, etc
	ta.parseFile(file, filename, analysis)

	// Calcula totais
	analysis.TotalResources = len(analysis.Resources)
	analysis.TotalModules = len(analysis.Modules)
	analysis.TotalVariables = len(analysis.Variables)
	analysis.TotalOutputs = len(analysis.Outputs)

	ta.checkBestPractices(analysis)

	return analysis, nil
}

// analyzeFile analisa um arquivo individual
func (ta *TerraformAnalyzer) analyzeFile(path string) (*models.TerraformAnalysis, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	return ta.AnalyzeContent(string(content), path)
}

// parseFile faz parsing de um arquivo HCL
func (ta *TerraformAnalyzer) parseFile(file *hcl.File, filename string, analysis *models.TerraformAnalysis) {
	// Extrai o conteúdo do body
	body, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}},
			{Type: "module", LabelNames: []string{"name"}},
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "output", LabelNames: []string{"name"}},
			{Type: "provider", LabelNames: []string{"name"}},
			{Type: "data", LabelNames: []string{"type", "name"}},
		},
	})

	if diags.HasErrors() {
		analysis.Valid = false
		for _, diag := range diags {
			analysis.SyntaxErrors = append(analysis.SyntaxErrors, models.SyntaxError{
				File:    filename,
				Line:    diag.Subject.Start.Line,
				Column:  diag.Subject.Start.Column,
				Message: diag.Summary,
				Snippet: diag.Detail,
			})
		}
	}

	// Parse blocks
	for _, block := range body.Blocks {
		switch block.Type {
		case "resource":
			ta.parseResource(block, filename, analysis)
		case "module":
			ta.parseModule(block, filename, analysis)
		case "variable":
			ta.parseVariable(block, filename, analysis)
		case "output":
			ta.parseOutput(block, filename, analysis)
		case "provider":
			if len(block.Labels) > 0 {
				analysis.Providers = append(analysis.Providers, block.Labels[0])
			}
		case "data":
			analysis.TotalDataSources++
		}
	}
}

// parseResource extrai informações de um resource block
func (ta *TerraformAnalyzer) parseResource(block *hcl.Block, filename string, analysis *models.TerraformAnalysis) {
	if len(block.Labels) < 2 {
		analysis.Valid = false
		analysis.SyntaxErrors = append(analysis.SyntaxErrors, models.SyntaxError{
			File:    filename,
			Line:    block.DefRange.Start.Line,
			Message: "Bloco de recurso malformado: requer tipo e nome",
		})
		return
	}

	resource := models.TerraformResource{
		Type:       block.Labels[0],
		Name:       block.Labels[1],
		Provider:   strings.Split(block.Labels[0], "_")[0],
		File:       filename,
		LineStart:  block.DefRange.Start.Line,
		LineEnd:    block.DefRange.End.Line,
		Attributes: make(map[string]interface{}),
	}

	// Avalia atributos
	evalCtx := &hcl.EvalContext{
		Functions: terraformFunctions,
	}
	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		analysis.Valid = false
		for _, diag := range diags {
			analysis.SyntaxErrors = append(analysis.SyntaxErrors, models.SyntaxError{
				File:    filename,
				Line:    diag.Subject.Start.Line,
				Column:  diag.Subject.Start.Column,
				Message: diag.Summary,
				Snippet: diag.Detail,
			})
		}
	}
	for name, attr := range attrs {
		val, valDiags := attr.Expr.Value(evalCtx)
		if valDiags.HasErrors() {
			continue
		}

		// Tenta converter para tipos Go nativos
		if val.IsKnown() && !val.IsNull() {
			switch val.Type() {
			case cty.String:
				resource.Attributes[name] = val.AsString()
			case cty.Bool:
				resource.Attributes[name] = val.True()
			case cty.Number:
				bf := val.AsBigFloat()
				f64, _ := bf.Float64()
				resource.Attributes[name] = f64
			}
		}
	}

	analysis.Resources = append(analysis.Resources, resource)
}

// parseModule extrai informações de um module block
func (ta *TerraformAnalyzer) parseModule(block *hcl.Block, filename string, analysis *models.TerraformAnalysis) {
	if len(block.Labels) < 1 {
		return
	}

	module := models.TerraformModule{
		Name:      block.Labels[0],
		File:      filename,
		LineStart: block.DefRange.Start.Line,
		Inputs:    make(map[string]interface{}),
	}

	analysis.Modules = append(analysis.Modules, module)
}

// parseVariable extrai informações de um variable block
func (ta *TerraformAnalyzer) parseVariable(block *hcl.Block, filename string, analysis *models.TerraformAnalysis) {
	if len(block.Labels) < 1 {
		return
	}

	variable := models.TerraformVariable{
		Name:     block.Labels[0],
		File:     filename,
		Required: true,
	}

	analysis.Variables = append(analysis.Variables, variable)
}

// parseOutput extrai informações de um output block
func (ta *TerraformAnalyzer) parseOutput(block *hcl.Block, filename string, analysis *models.TerraformAnalysis) {
	if len(block.Labels) < 1 {
		return
	}

	output := models.TerraformOutput{
		Name: block.Labels[0],
		File: filename,
	}

	analysis.Outputs = append(analysis.Outputs, output)
}

// mergeAnalysis combina resultados de análise de múltiplos arquivos
func (ta *TerraformAnalyzer) mergeAnalysis(dest, src *models.TerraformAnalysis) {
	dest.Resources = append(dest.Resources, src.Resources...)
	dest.Modules = append(dest.Modules, src.Modules...)
	dest.Variables = append(dest.Variables, src.Variables...)
	dest.Outputs = append(dest.Outputs, src.Outputs...)
	dest.SyntaxErrors = append(dest.SyntaxErrors, src.SyntaxErrors...)

	// Merge providers (unique)
	providerMap := make(map[string]bool)
	for _, p := range dest.Providers {
		providerMap[p] = true
	}
	for _, p := range src.Providers {
		providerMap[p] = true
	}
	dest.Providers = []string{}
	for p := range providerMap {
		dest.Providers = append(dest.Providers, p)
	}

	dest.Valid = dest.Valid && src.Valid
}

// checkBestPractices verifica best practices
func (ta *TerraformAnalyzer) checkBestPractices(analysis *models.TerraformAnalysis) {
	// Verifica se recursos têm tags
	for _, resource := range analysis.Resources {
		if len(resource.Tags) == 0 && ta.shouldHaveTags(resource.Type) {
			analysis.BestPracticeWarnings = append(analysis.BestPracticeWarnings,
				fmt.Sprintf("Recurso %s.%s não possui tags", resource.Type, resource.Name))
		}
	}

	// Verifica se há outputs
	if len(analysis.Outputs) == 0 && len(analysis.Resources) > 0 {
		analysis.BestPracticeWarnings = append(analysis.BestPracticeWarnings,
			"Considere adicionar outputs para facilitar integração com outros módulos")
	}

	// Verifica se variáveis têm descrição
	for _, variable := range analysis.Variables {
		if variable.Description == "" {
			analysis.BestPracticeWarnings = append(analysis.BestPracticeWarnings,
				fmt.Sprintf("Variável %s não possui descrição", variable.Name))
		}
	}
}

// shouldHaveTags verifica se um tipo de recurso deveria ter tags
func (ta *TerraformAnalyzer) shouldHaveTags(resourceType string) bool {
	taggableTypes := []string{
		"aws_instance", "aws_s3_bucket", "aws_vpc", "aws_subnet",
		"aws_security_group", "aws_rds_instance", "aws_lambda_function",
		"azurerm_resource_group", "azurerm_virtual_machine",
		"google_compute_instance", "google_storage_bucket",
	}

	for _, t := range taggableTypes {
		if t == resourceType {
			return true
		}
	}
	return false
}
