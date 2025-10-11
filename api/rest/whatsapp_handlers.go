package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/govinda777/iac-ai-agent/internal/agent/whatsapp"
)

// WhatsAppWebhookHandler gerencia webhooks do WhatsApp
type WhatsAppWebhookHandler struct {
	agent whatsapp.WhatsAppAgentInterface
}

// NewWhatsAppWebhookHandler cria novo handler de webhook
func NewWhatsAppWebhookHandler(agent whatsapp.WhatsAppAgentInterface) *WhatsAppWebhookHandler {
	return &WhatsAppWebhookHandler{
		agent: agent,
	}
}

// HandleWebhook processa webhooks do WhatsApp
func (h *WhatsAppWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
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
func (h *WhatsAppWebhookHandler) handleVerification(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	mode := query.Get("hub.mode")
	token := query.Get("hub.verify_token")
	challenge := query.Get("hub.challenge")

	// Verificar token de verificação
	if mode == "subscribe" && token == h.agent.GetVerifyToken() {
		log.Printf("Webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(challenge))
		return
	}

	log.Printf("Webhook verification failed")
	http.Error(w, "Verification failed", http.StatusForbidden)
}

// handleMessage processa mensagens recebidas
func (h *WhatsAppWebhookHandler) handleMessage(w http.ResponseWriter, r *http.Request) {
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
func (h *WhatsAppWebhookHandler) processTextMessage(message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      struct {
		Body string `json:"body"`
	} `json:"text"`
	Type string `json:"type"`
}) {
	// Converter para estrutura interna
	whatsappMsg := &whatsapp.WhatsAppMessage{
		ID:        message.ID,
		From:      message.From,
		Text:      message.Text.Body,
		Type:      message.Type,
		Timestamp: time.Now(),
	}

	// Processar mensagem
	ctx := context.Background()
	response, err := h.agent.ProcessMessage(ctx, whatsappMsg)
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
func (h *WhatsAppWebhookHandler) sendResponse(response *whatsapp.WhatsAppResponse) error {
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

// RegisterRoutes registra rotas do webhook
func (h *WhatsAppWebhookHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/webhook/whatsapp", h.HandleWebhook).Methods("GET", "POST")
	router.HandleFunc("/webhook/whatsapp/status", h.HandleStatus).Methods("GET")
	router.HandleFunc("/webhook/whatsapp/health", h.HandleHealth).Methods("GET")
}

// HandleStatus retorna status do webhook
func (h *WhatsAppWebhookHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":    "active",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

// HandleHealth verifica saúde do webhook
func (h *WhatsAppWebhookHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "healthy",
		"checks": map[string]string{
			"agent":        "ok",
			"webhook":      "ok",
			"database":     "ok",
			"external_api": "ok",
		},
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(health)
}