package rest

import (
	"encoding/json"
	"net/http"

	"github.com/govinda777/iac-ai-agent/internal/services"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// NotionAgentHandler gerencia endpoints relacionados aos agentes Notion
type NotionAgentHandler struct {
	notionService *services.NotionAgentService
	logger        *logger.Logger
}

// NewNotionAgentHandler cria um novo handler para agentes Notion
func NewNotionAgentHandler(notionService *services.NotionAgentService, log *logger.Logger) *NotionAgentHandler {
	return &NotionAgentHandler{
		notionService: notionService,
		logger:        log,
	}
}

// CreateAgentRequest representa a requisição para criar um agente
type CreateAgentRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
}

// CreateAgentResponse representa a resposta da criação de agente
type CreateAgentResponse struct {
	Agent services.AgentInfo `json:"agent"`
}

// ListAgentsResponse representa a resposta da listagem de agentes
type ListAgentsResponse struct {
	Agents []services.AgentInfo `json:"agents"`
	Total  int                  `json:"total"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// CreateAgent cria um novo agente no Notion
func (h *NotionAgentHandler) CreateAgent(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Criando novo agente Notion")

	var req CreateAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validação básica
	if req.Name == "" {
		h.logger.Error("Nome do agente é obrigatório")
		http.Error(w, "Agent name is required", http.StatusBadRequest)
		return
	}

	// Cria o agente
	agent, err := h.notionService.CreateAgent(r.Context(), &services.CreateAgentRequest{
		Name:         req.Name,
		Description:  req.Description,
		Capabilities: req.Capabilities,
	})
	if err != nil {
		h.logger.Error("Erro ao criar agente", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta
	response := CreateAgentResponse{Agent: *agent}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Erro ao codificar resposta", "error", err)
	}
}

// GetOrCreateDefaultAgent obtém ou cria o agente padrão
func (h *NotionAgentHandler) GetOrCreateDefaultAgent(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Obtendo ou criando agente padrão Notion")

	agent, err := h.notionService.GetOrCreateDefaultAgent(r.Context())
	if err != nil {
		h.logger.Error("Erro ao obter/criar agente padrão", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta
	response := CreateAgentResponse{Agent: *agent}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Erro ao codificar resposta", "error", err)
	}
}

// ListAgents lista todos os agentes disponíveis
func (h *NotionAgentHandler) ListAgents(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Listando agentes Notion")

	agents, err := h.notionService.ListAgents(r.Context())
	if err != nil {
		h.logger.Error("Erro ao listar agentes", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta
	response := ListAgentsResponse{
		Agents: agents,
		Total:  len(agents),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Erro ao codificar resposta", "error", err)
	}
}

// GetAgent obtém um agente específico por ID
func (h *NotionAgentHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL (implementação simplificada)
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		h.logger.Error("ID do agente é obrigatório")
		http.Error(w, "Agent ID is required", http.StatusBadRequest)
		return
	}

	h.logger.Info("Obtendo agente", "agent_id", agentID)

	agent, err := h.notionService.GetAgent(r.Context(), agentID)
	if err != nil {
		h.logger.Error("Erro ao obter agente", "error", err, "agent_id", agentID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta
	response := CreateAgentResponse{Agent: *agent}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Erro ao codificar resposta", "error", err)
	}
}

// DeleteAgent remove um agente
func (h *NotionAgentHandler) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL (implementação simplificada)
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		h.logger.Error("ID do agente é obrigatório")
		http.Error(w, "Agent ID is required", http.StatusBadRequest)
		return
	}

	h.logger.Info("Removendo agente", "agent_id", agentID)

	err := h.notionService.DeleteAgent(r.Context(), agentID)
	if err != nil {
		h.logger.Error("Erro ao remover agente", "error", err, "agent_id", agentID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetServiceStatus verifica o status do serviço Notion
func (h *NotionAgentHandler) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Verificando status do serviço Notion")

	available := h.notionService.IsServiceAvailable(r.Context())

	status := map[string]interface{}{
		"available": available,
		"service":   "notion",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(status); err != nil {
		h.logger.Error("Erro ao codificar resposta", "error", err)
	}
}

// GetDefaultCapabilities retorna as capacidades padrão do agente
func (h *NotionAgentHandler) GetDefaultCapabilities(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Obtendo capacidades padrão do agente")

	capabilities := h.notionService.GetDefaultAgentCapabilities()

	response := map[string]interface{}{
		"capabilities": capabilities,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Erro ao codificar resposta", "error", err)
	}
}
