package cloudcontroller

import "github.com/gosouza/iac-ai-agent/internal/models"

// TerraformModule representa um módulo Terraform recomendado
type TerraformModule struct {
	Name        string
	Source      string
	Version     string
	Description string
	Provider    string
	Verified    bool
	UseCase     string
}

// ModuleRegistry mantém registro de módulos Terraform recomendados
type ModuleRegistry struct {
	modules map[string]*TerraformModule
}

// NewModuleRegistry cria um novo registry
func NewModuleRegistry() *ModuleRegistry {
	mr := &ModuleRegistry{
		modules: make(map[string]*TerraformModule),
	}

	mr.loadRecommendedModules()

	return mr
}

// loadRecommendedModules carrega módulos recomendados
func (mr *ModuleRegistry) loadRecommendedModules() {
	// AWS VPC
	mr.modules["aws-vpc"] = &TerraformModule{
		Name:        "terraform-aws-modules/vpc/aws",
		Source:      "terraform-aws-modules/vpc/aws",
		Version:     "~> 5.0",
		Description: "Módulo oficial para criação de VPC na AWS",
		Provider:    "aws",
		Verified:    true,
		UseCase:     "networking",
	}

	// AWS ECS
	mr.modules["aws-ecs"] = &TerraformModule{
		Name:        "terraform-aws-modules/ecs/aws",
		Source:      "terraform-aws-modules/ecs/aws",
		Version:     "~> 5.0",
		Description: "Módulo oficial para ECS cluster e serviços",
		Provider:    "aws",
		Verified:    true,
		UseCase:     "compute",
	}

	// AWS RDS
	mr.modules["aws-rds"] = &TerraformModule{
		Name:        "terraform-aws-modules/rds/aws",
		Source:      "terraform-aws-modules/rds/aws",
		Version:     "~> 6.0",
		Description: "Módulo oficial para RDS instances",
		Provider:    "aws",
		Verified:    true,
		UseCase:     "database",
	}

	// AWS S3
	mr.modules["aws-s3"] = &TerraformModule{
		Name:        "terraform-aws-modules/s3-bucket/aws",
		Source:      "terraform-aws-modules/s3-bucket/aws",
		Version:     "~> 3.0",
		Description: "Módulo oficial para S3 buckets com best practices",
		Provider:    "aws",
		Verified:    true,
		UseCase:     "storage",
	}

	// AWS Security Group
	mr.modules["aws-security-group"] = &TerraformModule{
		Name:        "terraform-aws-modules/security-group/aws",
		Source:      "terraform-aws-modules/security-group/aws",
		Version:     "~> 5.0",
		Description: "Módulo oficial para Security Groups",
		Provider:    "aws",
		Verified:    true,
		UseCase:     "security",
	}

	// AWS Lambda
	mr.modules["aws-lambda"] = &TerraformModule{
		Name:        "terraform-aws-modules/lambda/aws",
		Source:      "terraform-aws-modules/lambda/aws",
		Version:     "~> 6.0",
		Description: "Módulo oficial para Lambda functions",
		Provider:    "aws",
		Verified:    true,
		UseCase:     "compute",
	}
}

// GetModule retorna um módulo pelo ID
func (mr *ModuleRegistry) GetModule(id string) *TerraformModule {
	if module, ok := mr.modules[id]; ok {
		return module
	}
	return nil
}

// GetModulesByProvider retorna módulos de um provider
func (mr *ModuleRegistry) GetModulesByProvider(provider string) []*TerraformModule {
	modules := []*TerraformModule{}
	for _, module := range mr.modules {
		if module.Provider == provider {
			modules = append(modules, module)
		}
	}
	return modules
}

// FindApplicableModules is a placeholder implementation to satisfy the interface.
func (mr *ModuleRegistry) FindApplicableModules(resources []models.TerraformResource) []models.ApprovedModule {
	// This is a basic implementation. A real one would have more complex logic
	// to suggest modules based on the resources being used.
	return []models.ApprovedModule{}
}

// GetModulesByUseCase retorna módulos por caso de uso
func (mr *ModuleRegistry) GetModulesByUseCase(useCase string) []*TerraformModule {
	modules := []*TerraformModule{}
	for _, module := range mr.modules {
		if module.UseCase == useCase {
			modules = append(modules, module)
		}
	}
	return modules
}

// GetAllModules retorna todos os módulos
func (mr *ModuleRegistry) GetAllModules() []*TerraformModule {
	modules := []*TerraformModule{}
	for _, module := range mr.modules {
		modules = append(modules, module)
	}
	return modules
}
