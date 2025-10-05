package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/internal/services"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// Handler gerencia requisições HTTP
type Handler struct {
	config          *config.Config
	logger          *logger.Logger
	analysisService *services.AnalysisService
	reviewService   *services.ReviewService
}

// NewHandler cria um novo handler
func NewHandler(cfg *config.Config, log *logger.Logger) *Handler {
	return &Handler{
		config:          cfg,
		logger:          log,
		analysisService: services.NewAnalysisService(cfg, log),
		reviewService:   services.NewReviewService(cfg, log),
	}
}

// SetupRoutes configura as rotas da API
func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", h.HandleHealth).Methods("GET")

	// Analysis endpoints
	r.HandleFunc("/analyze", h.HandleAnalyze).Methods("POST")

	// Review endpoints
	r.HandleFunc("/review", h.HandleReview).Methods("POST")

	// Info endpoint
	r.HandleFunc("/", h.HandleRoot).Methods("GET")

	return r
}

// HandleHealth retorna status de saúde do serviço
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"service": "iac-ai-agent",
		"version": "1.0.0",
	})
}

// HandleRoot retorna informações sobre a API
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": "IaC AI Agent",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"health":  "GET /health",
			"analyze": "POST /analyze",
			"review":  "POST /review",
		},
	})
}

// HandleAnalyze processa requisição de análise
func (h *Handler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Requisição de análise recebida")

	// Parse request
	var req models.AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao fazer parse da requisição", "error", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validações básicas
	if req.Content == "" && req.Path == "" {
		h.respondError(w, "Either 'content' or 'path' must be provided", http.StatusBadRequest)
		return
	}

	// Executa análise
	response, err := h.analysisService.Analyze(&req)
	if err != nil {
		h.logger.Error("Erro ao executar análise", "error", err)
		h.respondError(w, "Analysis failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retorna resultado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	h.logger.Info("Análise concluída",
		"id", response.ID,
		"score", response.Score,
		"suggestions", len(response.Suggestions))
}

// HandleReview processa requisição de review
func (h *Handler) HandleReview(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Requisição de review recebida")

	// Parse request
	var req models.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao fazer parse da requisição", "error", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validações
	if req.Repository == "" {
		h.respondError(w, "Repository is required", http.StatusBadRequest)
		return
	}
	if req.PRNumber == 0 {
		h.respondError(w, "PR number is required", http.StatusBadRequest)
		return
	}
	if req.Owner == "" {
		h.respondError(w, "Owner is required", http.StatusBadRequest)
		return
	}

	// Executa review
	response, err := h.reviewService.ReviewPR(&req)
	if err != nil {
		h.logger.Error("Erro ao executar review", "error", err)
		h.respondError(w, "Review failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retorna resultado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	h.logger.Info("Review concluído",
		"id", response.ID,
		"pr", response.PRNumber,
		"score", response.Score)
}

// respondError envia resposta de erro
func (h *Handler) respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    string(rune(statusCode)),
		Message: message,
	})
}
