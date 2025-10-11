package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/govinda777/iac-ai-agent/internal/services"
)

// Web3Handler lida com endpoints relacionados a web3
type Web3Handler struct {
	web3AuthService *services.Web3AuthService
}

// NewWeb3Handler cria um novo manipulador web3
func NewWeb3Handler(web3AuthService *services.Web3AuthService) *Web3Handler {
	return &Web3Handler{
		web3AuthService: web3AuthService,
	}
}

// RegisterRoutes registra as rotas web3
func (h *Web3Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/auth/web3/verify", h.VerifyToken).Methods("POST")
	router.HandleFunc("/api/auth/web3/check-access", h.CheckAccess).Methods("POST")
	router.HandleFunc("/api/auth/web3/token-cost", h.GetTokenCost).Methods("GET")
	router.HandleFunc("/api/auth/web3/spend-tokens", h.SpendTokens).Methods("POST")
}

// VerifyTokenRequest é a requisição para verificar um token Privy
type VerifyTokenRequest struct {
	Token string `json:"token"`
}

// VerifyToken verifica um token de autenticação web3
func (h *Web3Handler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Verificar token
	user, err := h.web3AuthService.VerifyToken(r.Context(), req.Token)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Responder com dados do usuário
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

// CheckAccessRequest é a requisição para verificar acesso a uma operação
type CheckAccessRequest struct {
	WalletAddress string `json:"wallet_address"`
	Operation     string `json:"operation"`
}

// CheckAccess verifica se um usuário tem acesso a uma operação
func (h *Web3Handler) CheckAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CheckAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.WalletAddress == "" || req.Operation == "" {
		http.Error(w, "Wallet address and operation are required", http.StatusBadRequest)
		return
	}

	// Sanitizar endereço da wallet
	req.WalletAddress = strings.TrimSpace(req.WalletAddress)

	// Verificar acesso
	allowed, err := h.web3AuthService.IsOperationAllowed(r.Context(), req.WalletAddress, req.Operation)
	if err != nil {
		http.Error(w, "Error checking access: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar rate limit
	withinLimit, err := h.web3AuthService.CheckRateLimit(r.Context(), req.WalletAddress)
	if err != nil {
		http.Error(w, "Error checking rate limit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta
	response := struct {
		Allowed     bool   `json:"allowed"`
		WithinLimit bool   `json:"within_rate_limit"`
		Message     string `json:"message,omitempty"`
	}{
		Allowed:     allowed,
		WithinLimit: withinLimit,
	}

	if !allowed {
		response.Message = "Operation not allowed for this wallet"
	} else if !withinLimit {
		response.Message = "Rate limit exceeded"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// TokenCostRequest é a requisição para obter o custo de uma operação em tokens
type TokenCostRequest struct {
	Operation string `json:"operation"`
}

// GetTokenCost retorna o custo em tokens de uma operação
func (h *Web3Handler) GetTokenCost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TokenCostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Operation == "" {
		http.Error(w, "Operation is required", http.StatusBadRequest)
		return
	}

	// Obter custo da operação
	cost, err := h.web3AuthService.GetTokenCost(req.Operation)
	if err != nil {
		http.Error(w, "Error getting token cost: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta
	response := struct {
		Operation string `json:"operation"`
		TokenCost string `json:"token_cost"`
	}{
		Operation: req.Operation,
		TokenCost: cost,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// SpendTokensRequest é a requisição para gastar tokens
type SpendTokensRequest struct {
	WalletAddress string `json:"wallet_address"`
	Operation     string `json:"operation"`
}

// SpendTokens gasta tokens para uma operação
func (h *Web3Handler) SpendTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SpendTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.WalletAddress == "" || req.Operation == "" {
		http.Error(w, "Wallet address and operation are required", http.StatusBadRequest)
		return
	}

	// Sanitizar endereço da wallet
	req.WalletAddress = strings.TrimSpace(req.WalletAddress)

	// Verificar acesso
	allowed, err := h.web3AuthService.IsOperationAllowed(r.Context(), req.WalletAddress, req.Operation)
	if err != nil || !allowed {
		http.Error(w, "Operation not allowed: "+err.Error(), http.StatusForbidden)
		return
	}

	// Obter custo da operação
	costStr, err := h.web3AuthService.GetTokenCost(req.Operation)
	if err != nil {
		http.Error(w, "Error getting operation cost: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extrair valor numérico do custo (formato "X IACAI")
	parts := strings.Split(costStr, " ")
	costStr = parts[0]
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		http.Error(w, "Invalid cost format", http.StatusInternalServerError)
		return
	}

	// Simular gasto de tokens (no futuro será implementado no Web3AuthService)
	// Em uma implementação real, chamaríamos:
	// tokenBalance, err := h.web3AuthService.GetUserTokenBalance(r.Context(), req.WalletAddress)
	// newBalance, err := h.web3AuthService.SpendUserTokens(r.Context(), req.WalletAddress, cost, req.Operation)

	// Por enquanto, simulamos o gasto como bem-sucedido
	newBalance := 100 - cost // Valor simulado

	// Resposta
	response := struct {
		Success    bool   `json:"success"`
		Operation  string `json:"operation"`
		Cost       int    `json:"cost"`
		NewBalance int    `json:"new_balance"`
	}{
		Success:    true,
		Operation:  req.Operation,
		Cost:       cost,
		NewBalance: newBalance,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
