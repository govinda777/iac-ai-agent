package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// ReviewService orquestra processo de review de PRs
type ReviewService struct {
	analysisService *AnalysisService
	logger          *logger.Logger
}

// NewReviewService cria uma nova instância do serviço de review
func NewReviewService(analysisService *AnalysisService, log *logger.Logger) *ReviewService {
	return &ReviewService{
		analysisService: analysisService,
		logger:          log,
	}
}

// ReviewPR realiza review completo de um PR
func (rs *ReviewService) ReviewPR(request *models.ReviewRequest) (*models.ReviewResponse, error) {
	rs.logger.Info("Iniciando review de PR",
		"repository", request.Repository,
		"pr_number", request.PRNumber)

	// Por simplicidade, esta implementação assume que temos acesso aos arquivos
	// Em produção, isso se conectaria ao GitHub para buscar os arquivos do PR

	review := &models.ReviewResponse{
		ID:               uuid.New().String(),
		Repository:       request.Repository,
		PRNumber:         request.PRNumber,
		Status:           "commented",
		FilesAnalyzed:    0,
		TotalSuggestions: 0,
		FileReviews:      []models.FileReview{},
		Timestamp:        time.Now(),
	}

	// Simula análise básica
	// Em produção, buscaria os arquivos alterados do GitHub
	rs.logger.Info("Review de PR concluído",
		"pr_number", request.PRNumber,
		"status", review.Status)

	return review, nil
}

// ReviewFiles realiza review de arquivos específicos
func (rs *ReviewService) ReviewFiles(files []string) (*models.ReviewResponse, error) {
	rs.logger.Info("Iniciando review de arquivos", "file_count", len(files))

	review := &models.ReviewResponse{
		ID:               uuid.New().String(),
		Status:           "commented",
		FilesAnalyzed:    len(files),
		TotalSuggestions: 0,
		FileReviews:      []models.FileReview{},
		Timestamp:        time.Now(),
	}

	totalScore := 0
	allSuggestions := []models.Suggestion{}

	// Analisa cada arquivo
	for _, file := range files {
		fileReview, err := rs.reviewFile(file)
		if err != nil {
			rs.logger.Warn("Erro ao analisar arquivo", "file", file, "error", err)
			continue
		}

		review.FileReviews = append(review.FileReviews, *fileReview)
		totalScore += fileReview.Score
		allSuggestions = append(allSuggestions, fileReview.Suggestions...)
	}

	// Calcula score médio
	if len(review.FileReviews) > 0 {
		review.Score = totalScore / len(review.FileReviews)
	}

	review.TotalSuggestions = len(allSuggestions)

	// Define status baseado no score
	review.Status = rs.determineStatus(review.Score)
	review.Summary = rs.generateSummary(review)

	rs.logger.Info("Review concluído",
		"files_analyzed", review.FilesAnalyzed,
		"score", review.Score,
		"status", review.Status)

	return review, nil
}

// reviewFile analisa um arquivo individual
func (rs *ReviewService) reviewFile(filename string) (*models.FileReview, error) {
	// Lê o arquivo
	content, err := rs.readFile(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Analisa o conteúdo
	analysis, err := rs.analysisService.AnalyzeContent(string(content), filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar arquivo: %w", err)
	}

	// Cria review do arquivo
	fileReview := &models.FileReview{
		Filename:    filename,
		Status:      "modified",
		Suggestions: analysis.Suggestions,
		Score:       analysis.Score,
		Comments:    rs.generateComments(analysis),
	}

	return fileReview, nil
}

// readFile lê conteúdo de um arquivo
func (rs *ReviewService) readFile(filename string) ([]byte, error) {
	// Em produção, isso buscaria do GitHub ou do sistema de arquivos
	// Por enquanto, retorna erro pois é um placeholder
	return nil, fmt.Errorf("implementação pendente")
}

// determineStatus determina status do review baseado no score
func (rs *ReviewService) determineStatus(score int) string {
	if score >= 90 {
		return "approved"
	}
	if score >= 70 {
		return "commented"
	}
	return "changes_requested"
}

// generateSummary gera sumário do review
func (rs *ReviewService) generateSummary(review *models.ReviewResponse) string {
	if review.Score >= 90 {
		return fmt.Sprintf("✅ Excelente! Código analisado em %d arquivo(s) com score %d/100",
			review.FilesAnalyzed, review.Score)
	}
	if review.Score >= 70 {
		return fmt.Sprintf("✔️ Bom código! %d arquivo(s) analisados com score %d/100. Algumas melhorias sugeridas.",
			review.FilesAnalyzed, review.Score)
	}
	return fmt.Sprintf("⚠️ Atenção! %d arquivo(s) analisados com score %d/100. Melhorias necessárias.",
		review.FilesAnalyzed, review.Score)
}

// generateComments gera comentários baseados na análise
func (rs *ReviewService) generateComments(analysis *models.AnalysisResponse) []models.Comment {
	comments := []models.Comment{}

	// Gera comentários para sugestões críticas e de alta severidade
	for _, suggestion := range analysis.Suggestions {
		if suggestion.Severity == "critical" || suggestion.Severity == "high" {
			comment := models.Comment{
				Path: suggestion.File,
				Line: suggestion.Line,
				Body: fmt.Sprintf("**[%s]** %s\n\n%s",
					suggestion.Severity,
					suggestion.Message,
					suggestion.Recommendation),
				Side: "RIGHT",
			}
			comments = append(comments, comment)
		}
	}

	return comments
}

// ApproveIfScoreIsHigh aprova automaticamente se score for alto
func (rs *ReviewService) ApproveIfScoreIsHigh(review *models.ReviewResponse, threshold int) bool {
	if review.Score >= threshold {
		review.Status = "approved"
		rs.logger.Info("PR aprovado automaticamente",
			"score", review.Score,
			"threshold", threshold)
		return true
	}
	return false
}
