package services

import (
	"encoding/json"
	"strings"

	"github.com/gosouza/iac-ai-agent/internal/models"
)

// enrichSuggestionsWithLLM usa LLM para adicionar contexto inteligente
func (as *AnalysisService) enrichSuggestionsWithLLM(
	analysis *models.AnalysisDetails,
	baseSuggestions []models.Suggestion,
) ([]models.Suggestion, error) {

	// Consulta Knowledge Base
	relevantPractices := as.knowledgeBase.GetRelevantPractices(analysis)
	relevantModules := as.moduleRegistry.FindApplicableModules(analysis.Terraform.Resources)

	// Constrói prompt contextualizado
	prompt, err := as.promptBuilder.BuildEnrichmentPrompt(
		analysis,
		baseSuggestions,
        relevantPractices,
        relevantModules,
	)
	if err != nil {
		as.logger.Warn("Failed to build LLM enrichment prompt", "error", err)
		return nil, err
	}

	// Chama LLM
	llmResp, err := as.llmClient.Generate(&models.LLMRequest{
		Prompt:      prompt,
		Temperature: 0.2,
		MaxTokens:   2000,
	})
	if err != nil {
		as.logger.Warn("LLM enrichment failed, using base suggestions", "error", err)
		return baseSuggestions, nil // Fallback gracioso
	}

	// Parse resposta estruturada
	var enrichedSuggestions []models.EnrichedSuggestion
	if err := as.parseLLMSuggestions(llmResp.Content, &enrichedSuggestions); err != nil {
        as.logger.Warn("Failed to parse LLM suggestions", "error", err)
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
	return json.Unmarshal([]byte(jsonStr), target)
}

func extractJSONFromText(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || start > end {
		return ""
	}
	return text[start : end+1]
}

func (as *AnalysisService) convertEnrichedSuggestions(enriched []models.EnrichedSuggestion) []models.Suggestion {
	suggestions := make([]models.Suggestion, len(enriched))
	for i, e := range enriched {
		suggestion := models.Suggestion{
			Type:           e.Type,
			Severity:       e.Severity,
			Message:        e.Message,
			Recommendation: e.CodeExample,
		}
		if len(e.References) > 0 {
			suggestion.ReferenceLink = e.References[0]
		}
		suggestions[i] = suggestion
	}
	return suggestions
}