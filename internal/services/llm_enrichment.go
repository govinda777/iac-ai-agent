package services

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gosouza/iac-ai-agent/internal/models"
)

// enrichSuggestionsWithLLM usa LLM para adicionar contexto inteligente às sugestões
func (as *AnalysisService) enrichSuggestionsWithLLM(
	analysis *models.AnalysisDetails,
	baseSuggestions []models.Suggestion,
) ([]models.Suggestion, error) {
	// Verifica se temos um cliente LLM configurado
	if as.llmClient == nil {
		as.logger.Warn("LLM client não configurado, retornando sugestões base")
		return baseSuggestions, nil
	}

	// Consulta Knowledge Base
	relevantPractices := as.knowledgeBase.GetRelevantPractices(analysis)
	relevantModules := as.moduleRegistry.FindApplicableModules(analysis.Terraform.Resources)

	// Constrói prompt contextualizado
	prompt := as.promptBuilder.BuildEnrichmentPrompt(
		analysis,
		baseSuggestions,
		relevantPractices,
		relevantModules,
	)

	// Chama LLM
	llmResp, err := as.llmClient.Generate(&models.LLMRequest{
		Prompt:      prompt,
		Temperature: 0.2,
		MaxTokens:   2000,
	})
	if err != nil {
		as.logger.Warn("LLM enrichment falhou, usando sugestões base", "error", err)
		return baseSuggestions, nil // Fallback gracioso
	}

	// Parse resposta estruturada
	var enrichedSuggestions []models.EnrichedSuggestion
	if err := as.parseLLMSuggestions(llmResp.Content, &enrichedSuggestions); err != nil {
		as.logger.Warn("Falha ao fazer parse da resposta LLM", "error", err)
		return baseSuggestions, nil
	}

	// Converte para formato padrão
	return as.convertEnrichedSuggestions(enrichedSuggestions), nil
}

// parseLLMSuggestions extrai sugestões estruturadas da resposta LLM
func (as *AnalysisService) parseLLMSuggestions(
	content string,
	target interface{},
) error {
	// Extrai JSON da resposta (LLM pode adicionar texto antes/depois)
	jsonStr := extractJSONFromText(content)
	if jsonStr == "" {
		return fmt.Errorf("nenhum JSON encontrado na resposta LLM")
	}
	return json.Unmarshal([]byte(jsonStr), target)
}

// extractJSONFromText extrai o primeiro bloco JSON válido de um texto
func extractJSONFromText(text string) string {
	// Tenta encontrar JSON entre chaves {}
	re := regexp.MustCompile(`(?s)\{.*\}`)
	match := re.FindString(text)
	if match != "" {
		// Verifica se é um JSON válido
		var js map[string]interface{}
		if json.Unmarshal([]byte(match), &js) == nil {
			return match
		}
	}

	// Tenta encontrar JSON entre blocos de código
	codeBlockRe := regexp.MustCompile("```(?:json)?\\s*\\n?(.+?)\\n?```")
	if matches := codeBlockRe.FindStringSubmatch(text); len(matches) > 1 {
		jsonCandidate := matches[1]
		var js interface{}
		if json.Unmarshal([]byte(jsonCandidate), &js) == nil {
			return jsonCandidate
		}
	}

	return ""
}

// convertEnrichedSuggestions converte sugestões enriquecidas para o formato padrão
func (as *AnalysisService) convertEnrichedSuggestions(
	enriched []models.EnrichedSuggestion,
) []models.Suggestion {
	result := []models.Suggestion{}

	for _, e := range enriched {
		suggestion := models.Suggestion{
			Type:           e.Type,
			Severity:       e.Severity,
			Message:        e.Message,
			Recommendation: e.CodeExample,
			Resource:       e.Resource,
			References:     e.References,
			Metadata: map[string]interface{}{
				"implementation_effort": e.ImplementationEffort,
				"estimated_impact":      e.EstimatedImpact,
				"why_it_matters":        e.WhyItMatters,
			},
		}

		result = append(result, suggestion)
	}

	return result
}

// mergeSuggestions combina sugestões baseadas em regras com sugestões do LLM
func (as *AnalysisService) mergeSuggestions(
	ruleBased []models.Suggestion,
	enriched []models.Suggestion,
) []models.Suggestion {
	// Mapa para evitar duplicações
	uniqueSuggestions := make(map[string]models.Suggestion)

	// Adiciona sugestões baseadas em regras
	for _, s := range ruleBased {
		key := fmt.Sprintf("%s:%s:%s", s.Type, s.Severity, s.Resource)
		uniqueSuggestions[key] = s
	}

	// Adiciona sugestões enriquecidas, substituindo as existentes se tiverem a mesma chave
	for _, s := range enriched {
		key := fmt.Sprintf("%s:%s:%s", s.Type, s.Severity, s.Resource)
		uniqueSuggestions[key] = s
	}

	// Converte mapa de volta para slice
	result := []models.Suggestion{}
	for _, s := range uniqueSuggestions {
		result = append(result, s)
	}

	return result
}

// generateSuggestions gera sugestões baseadas na análise
// Versão atualizada que usa LLM para enriquecimento
func (as *AnalysisService) generateSuggestions(
	analysis *models.AnalysisDetails,
) []models.Suggestion {
	// 1. Sugestões baseadas em regras (rápido, determinístico)
	ruleBased := as.generateRuleBasedSuggestions(analysis)

	// 2. Enriquece com LLM (contexto, inteligência)
	enriched, err := as.enrichSuggestionsWithLLM(
		analysis,
		ruleBased,
	)

	// Se houve erro, usa apenas as sugestões baseadas em regras
	if err != nil {
		as.logger.Warn("Falha ao enriquecer sugestões com LLM", "error", err)
		return ruleBased
	}

	// 3. Combina e remove duplicatas
	return as.mergeSuggestions(ruleBased, enriched)
}

// generateRuleBasedSuggestions gera sugestões baseadas apenas em regras
func (as *AnalysisService) generateRuleBasedSuggestions(
	analysis *models.AnalysisDetails,
) []models.Suggestion {
	suggestions := []models.Suggestion{}

	// Adiciona warnings de best practices
	for _, warning := range analysis.Terraform.BestPracticeWarnings {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "best_practice",
			Severity:       "medium",
			Message:        warning.Message,
			Recommendation: warning.Recommendation,
			Resource:       warning.Resource,
		})
	}

	// Adiciona warnings de segurança do Checkov
	for _, finding := range analysis.Checkov.Findings {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "security",
			Severity:       finding.Severity,
			Message:        finding.Description,
			Recommendation: finding.Guideline,
			Resource:       finding.Resource,
			References:     []string{finding.Documentation},
		})
	}

	// Adiciona warnings de IAM
	for _, warning := range analysis.IAM.Warnings {
		suggestions = append(suggestions, models.Suggestion{
			Type:           "iam",
			Severity:       warning.Severity,
			Message:        warning.Message,
			Recommendation: warning.Recommendation,
			Resource:       warning.Resource,
		})
	}

	return suggestions
}
