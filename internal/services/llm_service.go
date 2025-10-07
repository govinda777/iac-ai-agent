package services

import (
	"context"

	"github.com/gosouza/iac-ai-agent/internal/agent/llm"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// LLMService representa o serviço para interação com modelos de linguagem
type LLMService struct {
	client  llm.Provider
	logger  *logger.Logger
	builder *llm.PromptBuilder
}

// NewLLMService cria uma nova instância do serviço LLM
func NewLLMService(client llm.Provider, log *logger.Logger) *LLMService {
	return &LLMService{
		client:  client,
		logger:  log,
		builder: llm.NewPromptBuilder(log),
	}
}

// EnrichSuggestions enriquece sugestões com informações do LLM
func (s *LLMService) EnrichSuggestions(
	ctx context.Context,
	analysis *models.AnalysisDetails,
	baseSuggestions []models.Suggestion,
	relevantPractices []models.BestPractice,
	relevantModules []models.Module,
) ([]models.Suggestion, error) {
	// Constrói o prompt para enriquecer sugestões
	prompt := s.builder.BuildEnrichmentPrompt(
		analysis,
		baseSuggestions,
		relevantPractices,
		relevantModules,
	)

	// Chama o LLM para enriquecer sugestões
	response, err := s.client.GetCompletion(ctx, prompt)
	if err != nil {
		s.logger.Error("Erro ao obter resposta do LLM para enriquecimento", "error", err)
		return baseSuggestions, err
	}

	// Processa resposta para extrair sugestões enriquecidas
	enrichedSuggestions, err := models.ParseEnrichedSuggestions(response)
	if err != nil {
		s.logger.Error("Erro ao analisar resposta do LLM para enriquecimento", "error", err)
		return baseSuggestions, err
	}

	return enrichedSuggestions, nil
}

// AnalyzePreview analisa um preview de terraform usando LLM
func (s *LLMService) AnalyzePreview(
	ctx context.Context,
	previewAnalysis *models.PreviewAnalysis,
) (*models.PreviewReview, error) {
	// Constrói o prompt para análise de preview
	prompt := s.builder.BuildPreviewAnalysisPrompt(previewAnalysis)

	// Chama o LLM para analisar o preview
	response, err := s.client.GetCompletion(ctx, prompt)
	if err != nil {
		s.logger.Error("Erro ao obter resposta do LLM para análise de preview", "error", err)
		return nil, err
	}

	// Processa resposta para extrair análise de preview
	previewReview, err := models.ParsePreviewReview(response)
	if err != nil {
		s.logger.Error("Erro ao analisar resposta do LLM para preview", "error", err)
		return nil, err
	}

	return previewReview, nil
}

// GenerateSecurityAnalysis gera uma análise de segurança detalhada
func (s *LLMService) GenerateSecurityAnalysis(
	ctx context.Context,
	analysis *models.AnalysisDetails,
	securityFindings []models.SecurityFinding,
) (*models.SecurityAnalysisResponse, error) {
	// Constrói o prompt para análise de segurança
	prompt := s.builder.BuildSecurityAnalysisPrompt(analysis, securityFindings)

	// Chama o LLM para análise de segurança
	response, err := s.client.GetCompletion(ctx, prompt)
	if err != nil {
		s.logger.Error("Erro ao obter resposta do LLM para segurança", "error", err)
		return nil, err
	}

	// Processa resposta para extrair análise de segurança
	securityAnalysis, err := models.ParseSecurityAnalysis(response)
	if err != nil {
		s.logger.Error("Erro ao analisar resposta do LLM para segurança", "error", err)
		return nil, err
	}

	return securityAnalysis, nil
}

// GenerateCostOptimization gera recomendações de otimização de custo
func (s *LLMService) GenerateCostOptimization(
	ctx context.Context,
	analysis *models.AnalysisDetails,
) (*models.CostOptimization, error) {
	// Constrói o prompt para otimização de custo
	prompt := s.builder.BuildCostOptimizationPrompt(analysis)

	// Chama o LLM para otimização de custo
	response, err := s.client.GetCompletion(ctx, prompt)
	if err != nil {
		s.logger.Error("Erro ao obter resposta do LLM para otimização de custo", "error", err)
		return nil, err
	}

	// Processa resposta para extrair recomendações de otimização de custo
	costOptimization, err := models.ParseCostOptimization(response)
	if err != nil {
		s.logger.Error("Erro ao analisar resposta do LLM para custo", "error", err)
		return nil, err
	}

	return costOptimization, nil
}