package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gosouza/iac-ai-agent/internal/models"
	"github.com/gosouza/iac-ai-agent/internal/services"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// WebhookHandler processa webhooks do GitHub
type WebhookHandler struct {
	config        *config.Config
	logger        *logger.Logger
	reviewService *services.ReviewService
	secret        string
}

// NewWebhookHandler cria um novo handler de webhooks
func NewWebhookHandler(cfg *config.Config, log *logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		config:        cfg,
		logger:        log,
		reviewService: services.NewReviewService(cfg, log),
		secret:        cfg.GitHub.WebhookSecret,
	}
}

// HandleGitHub processa webhook do GitHub
func (wh *WebhookHandler) HandleGitHub(w http.ResponseWriter, r *http.Request) {
	wh.logger.Info("Webhook recebido", "event", r.Header.Get("X-GitHub-Event"))

	// Verifica assinatura
	if !wh.verifySignature(r) {
		wh.logger.Warn("Assinatura inválida do webhook")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Lê body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		wh.logger.Error("Erro ao ler body", "error", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse payload
	var payload models.GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		wh.logger.Error("Erro ao fazer parse do payload", "error", err)
		http.Error(w, "Error parsing payload", http.StatusBadRequest)
		return
	}

	// Processa baseado no tipo de evento
	event := r.Header.Get("X-GitHub-Event")
	switch event {
	case "pull_request":
		wh.handlePullRequest(&payload, w)
	case "push":
		wh.handlePush(&payload, w)
	default:
		wh.logger.Info("Evento ignorado", "event", event)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Event ignored"})
	}
}

// handlePullRequest processa evento de pull request
func (wh *WebhookHandler) handlePullRequest(payload *models.GitHubWebhookPayload, w http.ResponseWriter) {
	// Verifica ações relevantes
	if payload.Action != "opened" && payload.Action != "synchronize" && payload.Action != "reopened" {
		wh.logger.Info("Ação ignorada", "action", payload.Action)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Action ignored"})
		return
	}

	if payload.PullRequest == nil {
		wh.logger.Error("Pull request não encontrado no payload")
		http.Error(w, "Missing pull request", http.StatusBadRequest)
		return
	}

	wh.logger.Info("Processando pull request",
		"repo", payload.Repository.FullName,
		"pr", payload.PullRequest.Number,
		"action", payload.Action)

	// Extrai owner e repo
	owner := payload.Repository.Owner.Login
	repo := payload.Repository.Name

	// Cria request de review
	reviewReq := &models.ReviewRequest{
		Repository: repo,
		Owner:      owner,
		PRNumber:   payload.PullRequest.Number,
	}

	if payload.Installation != nil {
		reviewReq.InstallationID = payload.Installation.ID
	}

	// Executa review em background (não bloqueia webhook)
	go func() {
		_, err := wh.reviewService.ReviewPR(reviewReq)
		if err != nil {
			wh.logger.Error("Erro ao processar review", "error", err)
		}
	}()

	// Responde imediatamente
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Review started",
		"pr":      fmt.Sprintf("%d", payload.PullRequest.Number),
	})
}

// handlePush processa evento de push
func (wh *WebhookHandler) handlePush(payload *models.GitHubWebhookPayload, w http.ResponseWriter) {
	wh.logger.Info("Push event recebido",
		"repo", payload.Repository.FullName,
		"ref", payload.PullRequest.Head.Ref)

	// Por enquanto, apenas loga
	// Em uma implementação completa, poderia disparar análise de branch

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Push event received"})
}

// verifySignature verifica a assinatura HMAC do webhook
func (wh *WebhookHandler) verifySignature(r *http.Request) bool {
	if wh.secret == "" {
		wh.logger.Warn("Webhook secret não configurado, pulando verificação")
		return true // Sem secret configurado, aceita qualquer requisição
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return false
	}

	// Lê body para calcular HMAC
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}

	// Restaura body para leitura posterior
	r.Body = io.NopCloser(io.Reader(body))

	// Calcula HMAC
	mac := hmac.New(sha256.New, []byte(wh.secret))
	mac.Write(body)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
