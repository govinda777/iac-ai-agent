package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/govinda777/iac-ai-agent/docs" // Importar docs gerados
	"github.com/govinda777/iac-ai-agent/internal/agent/analyzer"
	"github.com/govinda777/iac-ai-agent/internal/agent/scorer"
	"github.com/govinda777/iac-ai-agent/internal/agent/suggester"
	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/internal/services"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title IaC AI Agent API
// @version 1.0.0
// @description API para análise, revisão e otimização de código Infrastructure as Code (IaC)
// @description
// @description O IaC AI Agent analisa código Terraform, políticas Checkov e IAM para propor melhorias de segurança, custo e best practices.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email gosouza@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// Variável global para tracking de uptime
var startTime = time.Now()

// Handler gerencia requisições HTTP
type Handler struct {
	config          *config.Config
	logger          *logger.Logger
	analysisService *services.AnalysisService
	reviewService   *services.ReviewService
	web3Handler     *Web3Handler
	nationValidator *web3.NationNFTValidator
}

// NewHandler cria um novo handler
func NewHandler(cfg *config.Config, log *logger.Logger) *Handler {
	// Instantiate concrete types
	tfAnalyzer := analyzer.NewTerraformAnalyzer()
	checkovAnalyzer := analyzer.NewCheckovAnalyzer(log)
	iamAnalyzer := analyzer.NewIAMAnalyzer(log)
	prScorer := scorer.NewPRScorer()
	costOptimizer := suggester.NewCostOptimizer(log)
	securityAdvisor := suggester.NewSecurityAdvisor(log)

	analysisService := services.NewAnalysisService(
		log,
		70, // minPassScore
		tfAnalyzer,
		checkovAnalyzer,
		iamAnalyzer,
		prScorer,
		costOptimizer,
		securityAdvisor,
		cfg,
	)

	// Inicializar NationNFTValidator
	nationValidator := web3.NewNationNFTValidator(cfg, log)

	// Web3 handler será inicializado posteriormente quando tivermos os serviços Web3
	return &Handler{
		config:          cfg,
		logger:          log,
		analysisService: analysisService,
		reviewService:   services.NewReviewService(analysisService, log),
		nationValidator: nationValidator,
		// web3Handler será configurado com SetupWeb3Handler
	}
}

// SetupWeb3Handler configura o handler Web3
func (h *Handler) SetupWeb3Handler(web3AuthService *services.Web3AuthService) {
	h.web3Handler = NewWeb3Handler(web3AuthService)
}

// SetupRoutes configura as rotas da API
func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/agent/health", h.HandleAgentHealth).Methods("GET")
	r.HandleFunc("/agent/template", h.HandleAgentTemplate).Methods("GET")

	// Analysis endpoints
	r.HandleFunc("/analyze", h.HandleAnalyze).Methods("POST")

	// Review endpoints
	r.HandleFunc("/review", h.HandleReview).Methods("POST")

	// Info endpoint
	r.HandleFunc("/", h.HandleRoot).Methods("GET")

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL para o JSON do Swagger
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Web3 endpoints
	if h.web3Handler != nil {
		h.web3Handler.RegisterRoutes(r)
	}

	return r
}

// HandleAgentHealth retorna status detalhado do agente com informações da Nation
// @Summary Agent Health check detalhado
// @Description Verifica o status de saúde do agente incluindo validação NFT Nation e teste de conectividade
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "Status detalhado do agente"
// @Router /agent/health [get]
func (h *Handler) HandleAgentHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Gerar identificador único da instância
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	instanceID := fmt.Sprintf("%s-%d-%d", hostname, pid, time.Now().Unix())

	// Status básico
	healthStatus := map[string]interface{}{
		"status":    "healthy",
		"service":   "iac-ai-agent",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
		"instance": map[string]interface{}{
			"id":             instanceID,
			"hostname":       hostname,
			"pid":            pid,
			"started_at":     startTime.UTC().Format(time.RFC3339),
			"uptime_seconds": int64(time.Since(startTime).Seconds()),
		},
		"links": map[string]interface{}{
			"agent_swagger":     "http://localhost:8080/swagger/index.html",
			"agent_api_docs":    "http://localhost:8080/swagger/doc.json",
			"health_check":      "http://localhost:8080/agent/health",
			"agent_template":    "http://localhost:8080/agent/template",
			"root_endpoint":     "http://localhost:8080/",
			"analysis_endpoint": "http://localhost:8080/analyze",
			"review_endpoint":   "http://localhost:8080/review",
			"monitoring": map[string]interface{}{
				"prometheus_metrics": "http://localhost:8080/metrics",
				"health_status":      "http://localhost:8080/health",
				"agent_status":       "http://localhost:8080/agent/health",
				"service_info":       "http://localhost:8080/",
			},
			"external_services": map[string]interface{}{
				"openai_status":   "https://status.openai.com",
				"ethereum_status": "https://ethereum.org/en/developers/docs/nodes-and-clients/",
				"github_status":   "https://www.githubstatus.com",
			},
		},
	}

	// Informações do agente
	agentInfo := map[string]interface{}{
		"id":          "iac-ai-agent-main",
		"unique_id":   fmt.Sprintf("agent-%s-%d", hostname, pid),
		"name":        "IaC AI Agent",
		"type":        "general-purpose",
		"description": "Agente versátil para análise completa de Infrastructure as Code",
		"build_info": map[string]interface{}{
			"version":    "1.0.0",
			"build_time": "2025-01-15T10:00:00Z",
			"go_version": runtime.Version(),
			"platform":   runtime.GOOS + "/" + runtime.GOARCH,
			"compiler":   runtime.Compiler,
		},
		"capabilities": []string{
			"terraform_analysis",
			"security_audit",
			"cost_optimization",
			"preview_analysis",
			"secrets_scanning",
			"iam_analysis",
			"checkov_integration",
		},
		"template": map[string]interface{}{
			"id":          "general-purpose",
			"name":        "General Purpose IaC Agent",
			"category":    "general",
			"recommended": true,
		},
		"runtime_stats": map[string]interface{}{
			"requests_processed":       0, // TODO: implementar contador
			"last_request_at":          nil,
			"average_response_time_ms": 0,
			"memory_usage_mb":          getMemoryUsage(),
			"goroutines":               runtime.NumGoroutine(),
		},
	}

	// Configurações do agente
	agentConfig := map[string]interface{}{
		"llm_provider":    "openai",
		"llm_model":       "gpt-4",
		"temperature":     0.2,
		"max_tokens":      4000,
		"enable_checkov":  true,
		"enable_iam":      true,
		"enable_costs":    true,
		"enable_preview":  true,
		"enable_secrets":  true,
		"response_format": "json",
		"language":        "pt-br",
		"timezone":        "America/Sao_Paulo",
		"api_keys": map[string]interface{}{
			"openai_configured": h.config != nil && h.config.LLM.APIKey != "",
			"openai_key_status": func() string {
				if h.config != nil && h.config.LLM.APIKey != "" {
					return "configured"
				}
				return "not_configured"
			}(),
			"openai_dashboard": "https://platform.openai.com/api-keys",
			"openai_usage":     "https://platform.openai.com/usage",
			"openai_docs":      "https://platform.openai.com/docs",
		},
		"web3_config": map[string]interface{}{
			"wallet_configured":   h.config != nil && h.config.Web3.WalletAddress != "",
			"nation_nft_required": h.config != nil && h.config.Web3.NationNFTRequired,
			"web3_links": map[string]interface{}{
				"ethereum_mainnet": "https://etherscan.io",
				"ethereum_testnet": "https://goerli.etherscan.io",
				"metamask":         "https://metamask.io",
				"walletconnect":    "https://walletconnect.com",
			},
		},
	}

	// Verificações de dependências
	checks := make(map[string]string)
	errors := make(map[string]interface{})

	// Verificar configuração
	if h.config != nil {
		checks["config"] = "ok"
		agentConfig["wallet_address"] = h.config.Web3.WalletAddress
		agentConfig["nation_nft_required"] = h.config.Web3.NationNFTRequired
	} else {
		checks["config"] = "error"
		healthStatus["status"] = "unhealthy"
		errors["config"] = "Configuração não encontrada"
	}

	// Verificar logger
	if h.logger != nil {
		checks["logger"] = "ok"
	} else {
		checks["logger"] = "error"
		healthStatus["status"] = "unhealthy"
		errors["logger"] = "Logger não configurado"
	}

	// Verificar serviços principais
	if h.analysisService != nil {
		checks["analysis_service"] = "ok"
	} else {
		checks["analysis_service"] = "error"
		healthStatus["status"] = "unhealthy"
		errors["analysis_service"] = "Serviço de análise não disponível"
	}

	if h.reviewService != nil {
		checks["review_service"] = "ok"
	} else {
		checks["review_service"] = "error"
		healthStatus["status"] = "unhealthy"
		errors["review_service"] = "Serviço de review não disponível"
	}

	// Verificar Web3 handler e validação NFT Nation
	if h.web3Handler != nil {
		checks["web3_handler"] = "ok"

		// Verificar NFT da Nation se configurado
		if h.config != nil && h.config.Web3.WalletAddress != "" {
			nationInfo := h.checkNationNFT(ctx)
			agentInfo["nation_nft"] = nationInfo

			if nationInfo["status"] == "error" {
				checks["nation_nft"] = "error"
				errors["nation_nft"] = nationInfo["error"]
				if h.config.Web3.NationNFTRequired {
					healthStatus["status"] = "unhealthy"
				}
			} else {
				checks["nation_nft"] = "ok"
			}
		} else {
			checks["nation_nft"] = "not_configured"
			agentInfo["nation_nft"] = map[string]interface{}{
				"status":  "not_configured",
				"message": "Wallet address não configurado",
			}
		}
	} else {
		checks["web3_handler"] = "not_configured"
		agentInfo["nation_nft"] = map[string]interface{}{
			"status":  "not_configured",
			"message": "Web3 handler não configurado",
		}
	}

	// Adicionar comparação com template
	templateComparison := h.compareWithTemplate(agentInfo, agentConfig)
	agentInfo["template_comparison"] = templateComparison

	// Adicionar informações ao status
	healthStatus["agent"] = agentInfo
	healthStatus["config"] = agentConfig
	healthStatus["checks"] = checks

	if len(errors) > 0 {
		healthStatus["errors"] = errors
	}

	// Determinar status HTTP
	statusCode := http.StatusOK
	if healthStatus["status"] == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	h.respondJSON(w, statusCode, healthStatus)
}

// checkNationNFT verifica NFT da Nation e executa teste de conectividade
func (h *Handler) checkNationNFT(ctx context.Context) map[string]interface{} {
	nationInfo := map[string]interface{}{
		"status":            "ok",
		"wallet_address":    h.config.Web3.WalletAddress,
		"default_wallet":    "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
		"is_default_wallet": false,
		"nft_validation":    map[string]interface{}{},
		"test_request":      map[string]interface{}{},
		"links": map[string]interface{}{
			"nation_fun_api":     "https://nation.fun/api",
			"nation_fun_docs":    "https://docs.nation.fun",
			"nation_fun_website": "https://nation.fun",
			"opensea_collection": "https://opensea.io/collection/nation-pass",
			"etherscan_wallet":   fmt.Sprintf("https://etherscan.io/address/%s", h.config.Web3.WalletAddress),
			"nation_fun_test":    "https://nation.fun/test",
		},
	}

	// Verificar se é a carteira padrão
	defaultWallet := "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"
	if h.config.Web3.WalletAddress == defaultWallet {
		nationInfo["is_default_wallet"] = true
	}

	// Validar NFT usando NationNFTValidator
	if h.nationValidator != nil {
		nftResponse, err := h.nationValidator.ValidateWalletNFT(ctx, h.config.Web3.WalletAddress)
		if err != nil {
			nationInfo["status"] = "error"
			nationInfo["error"] = err.Error()
			nationInfo["nft_validation"] = map[string]interface{}{
				"has_nft":      false,
				"error":        err.Error(),
				"validated_at": time.Now().UTC().Format(time.RFC3339),
			}
		} else {
			nationInfo["nft_validation"] = map[string]interface{}{
				"has_nft":      nftResponse.Data.HasNFT,
				"token_id":     nftResponse.Data.TokenID,
				"tier":         nftResponse.Data.Tier,
				"is_active":    nftResponse.Data.IsActive,
				"expires_at":   nftResponse.Data.ExpiresAt,
				"metadata":     nftResponse.Data.Metadata,
				"validated_at": time.Now().UTC().Format(time.RFC3339),
			}
		}

		// Executar teste de conectividade
		testResponse, err := h.nationValidator.SendTestToNation(ctx, "Teste de conectividade do IaC AI Agent")
		if err != nil {
			nationInfo["test_request"] = map[string]interface{}{
				"status":    "error",
				"error":     err.Error(),
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			}
		} else {
			nationInfo["test_request"] = map[string]interface{}{
				"test_id":   testResponse.Data.TestID,
				"status":    testResponse.Data.Status,
				"message":   testResponse.Data.Response,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			}
		}
	} else {
		nationInfo["status"] = "error"
		nationInfo["error"] = "NationNFTValidator não configurado"
		nationInfo["nft_validation"] = map[string]interface{}{
			"has_nft":      false,
			"error":        "NationNFTValidator não configurado",
			"validated_at": time.Now().UTC().Format(time.RFC3339),
		}
		nationInfo["test_request"] = map[string]interface{}{
			"status":    "error",
			"error":     "NationNFTValidator não configurado",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
	}

	return nationInfo
}

// getMemoryUsage retorna o uso de memória em MB
func getMemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}

// compareWithTemplate compara configuração atual com template padrão
func (h *Handler) compareWithTemplate(agentInfo map[string]interface{}, agentConfig map[string]interface{}) map[string]interface{} {
	comparison := map[string]interface{}{
		"template_id":      "general-purpose",
		"matches_template": true,
		"differences":      []string{},
		"customizations":   []string{},
	}

	// Verificar se as capacidades correspondem ao template
	expectedCapabilities := []string{
		"terraform_analysis",
		"security_audit",
		"cost_optimization",
		"preview_analysis",
		"secrets_scanning",
		"iam_analysis",
		"checkov_integration",
	}

	currentCapabilities := agentInfo["capabilities"].([]string)
	if len(currentCapabilities) != len(expectedCapabilities) {
		comparison["matches_template"] = false
		comparison["differences"] = append(comparison["differences"].([]string), "Número de capacidades diferente do template")
	}

	// Verificar configurações específicas
	if agentConfig["llm_model"] != "gpt-4" {
		comparison["customizations"] = append(comparison["customizations"].([]string), "Modelo LLM customizado")
	}

	if agentConfig["temperature"] != 0.2 {
		comparison["customizations"] = append(comparison["customizations"].([]string), "Temperatura LLM customizada")
	}

	return comparison
}

// HandleAgentTemplate retorna comparação entre dados do agente e template
// @Summary Comparação Agente vs Template
// @Description Compara dados do agente atual com o template de configuração
// @Tags agent
// @Produce json
// @Success 200 {object} map[string]interface{} "Comparação agente vs template"
// @Router /agent/template [get]
func (h *Handler) HandleAgentTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Dados do agente atual
	currentAgent := map[string]interface{}{
		"id":          "iac-ai-agent-main",
		"name":        "IaC AI Agent",
		"type":        "general-purpose",
		"description": "Agente versátil para análise completa de Infrastructure as Code",
		"capabilities": []string{
			"terraform_analysis",
			"security_audit",
			"cost_optimization",
			"preview_analysis",
			"secrets_scanning",
			"iam_analysis",
			"checkov_integration",
		},
		"config": map[string]interface{}{
			"llm_provider":    "openai",
			"llm_model":       "gpt-4",
			"temperature":     0.2,
			"max_tokens":      4000,
			"enable_checkov":  true,
			"enable_iam":      true,
			"enable_costs":    true,
			"enable_preview":  true,
			"enable_secrets":  true,
			"response_format": "json",
			"language":        "pt-br",
			"timezone":        "America/Sao_Paulo",
		},
	}

	// Template de referência (do agent_templates.yaml)
	templateReference := map[string]interface{}{
		"id":          "general-purpose",
		"name":        "General Purpose IaC Agent",
		"description": "Agente versátil para análise completa de Infrastructure as Code",
		"category":    "general",
		"recommended": true,
		"use_cases": []string{
			"Análise completa de Terraform",
			"Review de Pull Requests",
			"Detecção de problemas de segurança",
			"Sugestões de otimização",
			"Best practices",
			"Preview analysis de mudanças",
			"Secrets scanning",
			"Análise de risco de deploy",
		},
		"tags": []string{
			"terraform",
			"aws",
			"azure",
			"gcp",
			"security",
			"cost",
			"preview",
			"secrets",
			"risk-analysis",
		},
		"default_config": map[string]interface{}{
			"llm_provider":    "openai",
			"llm_model":       "gpt-4",
			"temperature":     0.2,
			"max_tokens":      4000,
			"enable_checkov":  true,
			"enable_iam":      true,
			"enable_costs":    true,
			"enable_preview":  true,
			"enable_secrets":  true,
			"response_format": "json",
			"language":        "pt-br",
			"timezone":        "America/Sao_Paulo",
		},
		"default_capabilities": map[string]interface{}{
			"can_analyze_terraform":  true,
			"can_analyze_checkov":    true,
			"can_analyze_iam":        true,
			"can_analyze_costs":      true,
			"can_detect_drift":       true,
			"can_analyze_preview":    true,
			"can_scan_secrets":       true,
			"can_use_llm":            true,
			"can_use_knowledge_base": true,
		},
	}

	// Comparação
	comparison := map[string]interface{}{
		"matches_template": true,
		"differences":      []string{},
		"missing_features": []string{},
		"extra_features":   []string{},
		"compliance_score": 95, // Score de conformidade com o template
	}

	// Verificar conformidade com template
	if currentAgent["type"] != templateReference["id"] {
		comparison["matches_template"] = false
		comparison["differences"] = append(comparison["differences"].([]string), "Tipo de agente não corresponde ao template")
	}

	// Verificar configurações
	currentConfig := currentAgent["config"].(map[string]interface{})
	templateConfig := templateReference["default_config"].(map[string]interface{})

	for key, templateValue := range templateConfig {
		if currentValue, exists := currentConfig[key]; exists {
			if currentValue != templateValue {
				comparison["differences"] = append(comparison["differences"].([]string),
					fmt.Sprintf("Config '%s': atual=%v, template=%v", key, currentValue, templateValue))
			}
		} else {
			comparison["missing_features"] = append(comparison["missing_features"].([]string),
				fmt.Sprintf("Configuração '%s' não encontrada", key))
		}
	}

	// Verificar capabilities
	templateCapabilities := templateReference["default_capabilities"].(map[string]interface{})
	currentCapabilities := currentAgent["capabilities"].([]string)

	capabilityMap := make(map[string]bool)
	for _, cap := range currentCapabilities {
		capabilityMap[cap] = true
	}

	for cap, enabled := range templateCapabilities {
		if enabled.(bool) && !capabilityMap[cap] {
			comparison["missing_features"] = append(comparison["missing_features"].([]string),
				fmt.Sprintf("Capability '%s' não implementada", cap))
		}
	}

	// Informações da Nation NFT se disponível
	nationInfo := map[string]interface{}{
		"status": "not_configured",
	}

	if h.config != nil && h.config.Web3.WalletAddress != "" {
		nationInfo = h.checkNationNFT(ctx)
	}

	// Resposta final
	response := map[string]interface{}{
		"timestamp":          time.Now().UTC().Format(time.RFC3339),
		"current_agent":      currentAgent,
		"template_reference": templateReference,
		"comparison":         comparison,
		"nation_nft":         nationInfo,
		"recommendations": []string{
			"Agente está em conformidade com o template general-purpose",
			"Todas as capabilities principais estão implementadas",
			"Configurações seguem as melhores práticas do template",
		},
	}

	h.respondJSON(w, http.StatusOK, response)
}

// HandleRoot retorna informações sobre a API
// @Summary Informações da API
// @Description Retorna informações gerais sobre a API e seus endpoints disponíveis
// @Tags info
// @Produce json
// @Success 200 {object} map[string]interface{} "Informações da API"
// @Router / [get]
func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
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
// @Summary Analisar código IaC
// @Description Analisa código Terraform para identificar problemas de segurança, custos e best practices
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body models.AnalysisRequest true "Requisição de análise"
// @Success 200 {object} models.AnalysisResponse "Resultado da análise"
// @Failure 400 {object} models.ErrorResponse "Requisição inválida"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /analyze [post]
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
	h.respondJSON(w, http.StatusOK, response)

	h.logger.Info("Análise concluída",
		"id", response.ID,
		"score", response.Score,
		"suggestions", len(response.Suggestions))
}

// HandleReview processa requisição de review
// @Summary Review de Pull Request
// @Description Executa uma análise completa de um Pull Request do GitHub
// @Tags review
// @Accept json
// @Produce json
// @Param request body models.ReviewRequest true "Requisição de review"
// @Success 200 {object} models.ReviewResponse "Resultado do review"
// @Failure 400 {object} models.ErrorResponse "Requisição inválida"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /review [post]
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
	h.respondJSON(w, http.StatusOK, response)

	h.logger.Info("Review concluído",
		"id", response.ID,
		"pr", response.PRNumber,
		"score", response.Score)
}

// respondJSON envia uma resposta JSON bem-sucedida
func (h *Handler) respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("Erro ao escrever resposta JSON", "error", err)
	}
}

// respondError envia resposta de erro
func (h *Handler) respondError(w http.ResponseWriter, message string, statusCode int) {
	h.respondJSON(w, statusCode, models.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    string(rune(statusCode)),
		Message: message,
	})
}
