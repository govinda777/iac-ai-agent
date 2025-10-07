package cloudcontroller

import (
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// ModuleRegistry gerencia informações sobre módulos Terraform
type ModuleRegistry struct {
	logger          *logger.Logger
	modules         map[string][]models.Module
	resourceModules map[string][]models.Module
	useCaseModules  map[string][]models.Module
}

// NewModuleRegistry cria um novo registro de módulos
func NewModuleRegistry(log *logger.Logger) *ModuleRegistry {
	mr := &ModuleRegistry{
		logger:          log,
		modules:         make(map[string][]models.Module),
		resourceModules: make(map[string][]models.Module),
		useCaseModules:  make(map[string][]models.Module),
	}

	// Carrega dados iniciais
	mr.loadModules()

	log.Info("Module Registry inicializado",
		"total_modules", len(mr.modules),
		"resource_mappings", len(mr.resourceModules),
		"use_cases", len(mr.useCaseModules))

	return mr
}

// FindApplicableModules encontra módulos aplicáveis para os recursos
func (mr *ModuleRegistry) FindApplicableModules(
	resources []models.TerraformResource,
) []models.Module {
	applicable := []models.Module{}
	seen := make(map[string]bool)

	// Busca por tipo de recurso
	for _, res := range resources {
		if modules, ok := mr.resourceModules[res.Type]; ok {
			for _, module := range modules {
				if !seen[module.Source] {
					seen[module.Source] = true
					applicable = append(applicable, module)
				}
			}
		}
	}

	// Adiciona módulos recomendados gerais
	for _, module := range mr.modules["recommended"] {
		if !seen[module.Source] {
			seen[module.Source] = true
			applicable = append(applicable, module)
		}
	}

	return applicable
}

// GetModuleBySource retorna um módulo pelo source
func (mr *ModuleRegistry) GetModuleBySource(source string) *models.Module {
	for _, modules := range mr.modules {
		for _, module := range modules {
			if module.Source == source {
				return &module
			}
		}
	}
	return nil
}

// GetModulesByUseCase retorna módulos por caso de uso
func (mr *ModuleRegistry) GetModulesByUseCase(useCase string) []models.Module {
	if modules, ok := mr.useCaseModules[useCase]; ok {
		return modules
	}
	return []models.Module{}
}

// loadModules carrega dados de módulos
func (mr *ModuleRegistry) loadModules() {
	// AWS VPC
	vpcModule := models.Module{
		Name:        "AWS VPC",
		Source:      "terraform-aws-modules/vpc/aws",
		Version:     "~> 5.0",
		Description: "Terraform module which creates VPC resources on AWS",
		UseCase:     "networking",
		Recommended: true,
		Resources:   []string{"aws_vpc", "aws_subnet", "aws_route_table"},
	}
	mr.addModule(vpcModule, "aws", "networking", "recommended")
	mr.mapResourceToModule("aws_vpc", vpcModule)
	mr.mapResourceToModule("aws_subnet", vpcModule)

	// AWS EKS
	eksModule := models.Module{
		Name:        "AWS EKS",
		Source:      "terraform-aws-modules/eks/aws",
		Version:     "~> 19.0",
		Description: "Terraform module which creates EKS resources on AWS",
		UseCase:     "container",
		Recommended: true,
		Resources:   []string{"aws_eks_cluster", "aws_eks_node_group"},
	}
	mr.addModule(eksModule, "aws", "container", "recommended")
	mr.mapResourceToModule("aws_eks_cluster", eksModule)

	// AWS RDS
	rdsModule := models.Module{
		Name:        "AWS RDS",
		Source:      "terraform-aws-modules/rds/aws",
		Version:     "~> 5.0",
		Description: "Terraform module which creates RDS resources on AWS",
		UseCase:     "database",
		Recommended: true,
		Resources:   []string{"aws_db_instance", "aws_db_parameter_group"},
	}
	mr.addModule(rdsModule, "aws", "database")
	mr.mapResourceToModule("aws_db_instance", rdsModule)

	// AWS S3
	s3Module := models.Module{
		Name:        "AWS S3",
		Source:      "terraform-aws-modules/s3-bucket/aws",
		Version:     "~> 3.0",
		Description: "Terraform module which creates S3 bucket on AWS with all/most features provided",
		UseCase:     "storage",
		Recommended: true,
		Resources:   []string{"aws_s3_bucket"},
	}
	mr.addModule(s3Module, "aws", "storage", "recommended")
	mr.mapResourceToModule("aws_s3_bucket", s3Module)

	// AWS Security Group
	sgModule := models.Module{
		Name:        "AWS Security Group",
		Source:      "terraform-aws-modules/security-group/aws",
		Version:     "~> 4.0",
		Description: "Terraform module which creates Security Group resources on AWS",
		UseCase:     "security",
		Recommended: true,
		Resources:   []string{"aws_security_group", "aws_security_group_rule"},
	}
	mr.addModule(sgModule, "aws", "security", "recommended")
	mr.mapResourceToModule("aws_security_group", sgModule)

	// Azure Virtual Network
	azVnetModule := models.Module{
		Name:        "Azure Virtual Network",
		Source:      "Azure/vnet/azurerm",
		Version:     "~> 4.0",
		Description: "Terraform module which creates virtual network resources on Azure",
		UseCase:     "networking",
		Recommended: true,
		Resources:   []string{"azurerm_virtual_network", "azurerm_subnet"},
	}
	mr.addModule(azVnetModule, "azure", "networking")
	mr.mapResourceToModule("azurerm_virtual_network", azVnetModule)

	// GCP Network
	gcpNetworkModule := models.Module{
		Name:        "GCP Network",
		Source:      "terraform-google-modules/network/google",
		Version:     "~> 6.0",
		Description: "Terraform module which creates network resources on GCP",
		UseCase:     "networking",
		Recommended: true,
		Resources:   []string{"google_compute_network", "google_compute_subnetwork"},
	}
	mr.addModule(gcpNetworkModule, "gcp", "networking")
	mr.mapResourceToModule("google_compute_network", gcpNetworkModule)
}

// addModule adiciona um módulo às coleções
func (mr *ModuleRegistry) addModule(module models.Module, categories ...string) {
	for _, category := range categories {
		mr.modules[category] = append(mr.modules[category], module)
	}
	mr.useCaseModules[module.UseCase] = append(mr.useCaseModules[module.UseCase], module)
}

// mapResourceToModule mapeia um recurso a um módulo
func (mr *ModuleRegistry) mapResourceToModule(resourceType string, module models.Module) {
	mr.resourceModules[resourceType] = append(mr.resourceModules[resourceType], module)
}
