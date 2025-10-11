package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/govinda777/iac-ai-agent/internal/agent/core"
)

// AgentHandler handler REST para o agente principal
type AgentHandler struct {
	agent *core.Agent
}

// NewAgentHandler cria novo handler do agente
func NewAgentHandler(agent *core.Agent) *AgentHandler {
	return &AgentHandler{
		agent: agent,
	}
}

// WhatsAppWebhookRequest estrutura da requisição do webhook WhatsApp
type WhatsAppWebhookRequest struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type string `json:"type"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// WhatsAppWebhookResponse estrutura da resposta do webhook WhatsApp
type WhatsAppWebhookResponse struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

// HandleWebhook processa webhooks do WhatsApp
func (h *AgentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleVerification(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.handleMessage(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleVerification processa verificação do webhook
func (h *AgentHandler) handleVerification(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	mode := query.Get("hub.mode")
	token := query.Get("hub.verify_token")
	challenge := query.Get("hub.challenge")

	// Verificar token de verificação
	if mode == "subscribe" && token == "your_verify_token_here" {
		log.Printf("Webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(challenge))
		return
	}

	log.Printf("Webhook verification failed")
	http.Error(w, "Verification failed", http.StatusForbidden)
}

// handleMessage processa mensagens recebidas
func (h *AgentHandler) handleMessage(w http.ResponseWriter, r *http.Request) {
	var webhookReq WhatsAppWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&webhookReq); err != nil {
		log.Printf("Failed to decode webhook request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Processar mensagens
	for _, entry := range webhookReq.Entry {
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {
				if message.Type == "text" {
					h.processTextMessage(message)
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

// processTextMessage processa mensagem de texto
func (h *AgentHandler) processTextMessage(message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      struct {
		Body string `json:"body"`
	} `json:"text"`
	Type string `json:"type"`
}) {
	// Converter para estrutura interna do agente
	agentMessage := &core.Message{
		ID:        message.ID,
		Source:    "whatsapp",
		Channel:   "whatsapp",
		From:      message.From,
		Text:      message.Text.Body,
		Type:      message.Type,
		Timestamp: time.Now(),
	}

	// Processar mensagem usando o agente
	ctx := context.Background()
	response, err := h.agent.ProcessMessage(ctx, agentMessage)
	if err != nil {
		log.Printf("Failed to process message: %v", err)
		return
	}

	// Enviar resposta
	if err := h.sendResponse(response); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

// sendResponse envia resposta via API do WhatsApp
func (h *AgentHandler) sendResponse(response *core.Response) error {
	// Preparar resposta
	webhookResp := WhatsAppWebhookResponse{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               response.To,
		Type:             response.Type,
		Text: struct {
			Body string `json:"body"`
		}{
			Body: response.Text,
		},
	}

	// Serializar JSON
	jsonData, err := json.Marshal(webhookResp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Enviar via API do WhatsApp
	// Em produção, implementar chamada real para API
	log.Printf("Sending WhatsApp response: %s", string(jsonData))

	return nil
}

// RegisterRoutes registra rotas do agente
func (h *AgentHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/webhook/whatsapp", h.HandleWebhook).Methods("GET", "POST")
	router.HandleFunc("/agent/status", h.HandleStatus).Methods("GET")
	router.HandleFunc("/agent/health", h.HandleHealth).Methods("GET")
	router.HandleFunc("/agent/capabilities", h.HandleCapabilities).Methods("GET")
}

// HandleStatus retorna status do agente
func (h *AgentHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	capabilities := h.agent.GetCapabilities()

	status := map[string]interface{}{
		"agent_id":     h.agent.ID,
		"agent_name":   h.agent.Name,
		"version":      h.agent.Version,
		"status":       "active",
		"capabilities": capabilities,
		"timestamp":    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

// HandleHealth verifica saúde do agente
func (h *AgentHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	errors := h.agent.HealthCheck(ctx)

	health := map[string]interface{}{
		"status": "healthy",
		"checks": map[string]string{
			"agent":        "ok",
			"capabilities": "ok",
		},
		"timestamp": time.Now(),
	}

	if len(errors) > 0 {
		health["status"] = "unhealthy"
		health["errors"] = errors
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(health)
}

// HandleCapabilities retorna lista de habilidades
func (h *AgentHandler) HandleCapabilities(w http.ResponseWriter, r *http.Request) {
	capabilities := h.agent.GetCapabilities()

	response := map[string]interface{}{
		"capabilities": capabilities,
		"count":        len(capabilities),
		"timestamp":    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// Middleware para logging de requisições
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log da requisição
		log.Printf("Agent Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Processar requisição
		next.ServeHTTP(w, r)

		// Log da resposta
		duration := time.Since(start)
		log.Printf("Agent Response: %s %s completed in %v", r.Method, r.URL.Path, duration)
	})
}

// Middleware para validação de token
func TokenValidationMiddleware(verifyToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Para requisições GET (verificação), não validar token aqui
			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// Para requisições POST, validar token se necessário
			// Em produção, implementar validação adequada
			next.ServeHTTP(w, r)
		})
	}
}
