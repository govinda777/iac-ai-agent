package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/govinda777/iac-ai-agent/api/rest"
	_ "github.com/govinda777/iac-ai-agent/docs" // Swagger docs
	"github.com/govinda777/iac-ai-agent/internal/platform/web3"
	"github.com/govinda777/iac-ai-agent/internal/services"
	"github.com/govinda777/iac-ai-agent/internal/startup"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	version = "1.0.0"
	banner  = `
	██╗ █████╗  ██████╗     █████╗ ██╗     █████╗  ██████╗ ███████╗███╗   ██╗████████╗
	██║██╔══██╗██╔════╝    ██╔══██╗██║    ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝
	██║███████║██║         ███████║██║    ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   
	██║██╔══██║██║         ██╔══██║██║    ██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   
	██║██║  ██║╚██████╗    ██║  ██║██║    ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   
	╚═╝╚═╝  ╚═╝ ╚═════╝    ╚═╝  ╚═╝╚═╝    ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   
	
	Infrastructure as Code AI Agent v%s
	Powered by LLM | Secured by Privy.io | Running on Base Network
	`
)

func main() {
	// Print banner
	fmt.Printf(banner, version)
	fmt.Println()

	// Load configuration
	cfg, err := config.Load("configs/app.yaml")
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	log.Info("🚀 Starting IaC AI Agent", "version", version)

	// Context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ============================================================
	// STARTUP VALIDATION (OBRIGATÓRIO)
	// ============================================================
	log.Info("🔍 Executando validações de startup...")

	// Verificar se estamos no modo de desenvolvimento sem validação
	if os.Getenv("ENABLE_STARTUP_VALIDATION") == "false" {
		log.Info("⚠️ Validações desabilitadas - modo desenvolvimento")
		log.Info("✅ Pulando validação de startup")
	} else {
		validator := startup.NewValidator(cfg, log)

		// Validar tudo - panic se falhar
		validator.MustValidate(ctx)

		log.Info("✅ Todas as validações passaram!")
	}
	log.Info("")

	// ============================================================
	// INITIALIZE SERVICES
	// ============================================================
	log.Info("📦 Inicializando serviços...")

	// Inicializar cliente Base Network
	baseClient, err := web3.NewBaseClient(cfg, log)
	if err != nil {
		log.Error("❌ Falha ao inicializar cliente Base Network", "error", err)
		os.Exit(1)
	}

	// Inicializar serviço de agentes
	agentService := services.NewAgentService(cfg, log, baseClient)

	// Verificar se já existe um agente para a wallet ou criar um novo
	log.Info("🤖 Verificando agente para wallet", "address", cfg.Web3.WalletAddress)
	agent, err := agentService.EnsureAgentExists(ctx)
	if err != nil {
		log.Warn("⚠️ Não foi possível verificar/criar agente", "error", err)
		log.Warn("⚠️ Algumas funcionalidades podem estar limitadas")
	} else {
		log.Info("✅ Agente configurado com sucesso",
			"agent_id", agent.ID,
			"template", agent.TemplateID,
			"contract", agent.ContractAddress)

		// Verificar se o agente tem chave do WhatsApp
		if agent.HasWhatsAppKey {
			log.Info("✅ Agente possui chave de API do WhatsApp configurada")
		} else {
			log.Info("ℹ️ Agente não possui chave de API do WhatsApp configurada")
			log.Info("ℹ️ Use a API para configurar a chave do WhatsApp usando Lit Protocol")
		}
	}

	log.Info("✅ Serviços inicializados com sucesso")

	// ============================================================
	// SETUP HTTP SERVER
	// ============================================================
	log.Info("🌐 Configurando servidor HTTP...")

	// API Handlers
	apiHandlers := rest.NewHandler(cfg, log)
	router := apiHandlers.SetupRoutes()

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+cfg.Server.Port+"/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	// HTTP Server
	server := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ============================================================
	// START SERVER
	// ============================================================
	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Info("🚀 Servidor HTTP iniciado",
			"address", cfg.GetAddress(),
			"environment", os.Getenv("ENVIRONMENT"))
		log.Info("📚 Swagger UI: http://localhost:" + cfg.Server.Port + "/swagger/")
		log.Info("❤️  Health Check: http://localhost:" + cfg.Server.Port + "/health")
		log.Info("")
		log.Info("✨ Aplicação pronta para receber requisições!")
		log.Info("Press Ctrl+C to shutdown gracefully")
		log.Info("")

		serverErrors <- server.ListenAndServe()
	}()

	// ============================================================
	// GRACEFUL SHUTDOWN
	// ============================================================
	// Channel to listen for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown
	select {
	case err := <-serverErrors:
		log.Error("❌ Erro ao iniciar servidor", "error", err)
		os.Exit(1)

	case sig := <-shutdown:
		log.Info("🛑 Shutdown signal recebido", "signal", sig)

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Asking listener to shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Error("❌ Erro no graceful shutdown", "error", err)

			// Force close
			if err := server.Close(); err != nil {
				log.Error("❌ Erro ao forçar fechamento", "error", err)
			}
		}

		log.Info("👋 Aplicação encerrada com sucesso")
	}
}
