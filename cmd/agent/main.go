package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gosouza/iac-ai-agent/api/rest"
	"github.com/gosouza/iac-ai-agent/internal/platform/webhook"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

var (
	configPath = flag.String("config", "configs/app.yaml", "Path to config file")
	version    = "1.0.0"
)

func main() {
	flag.Parse()

	// Carrega configuração
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("Erro ao carregar configuração: %v\n", err)
		os.Exit(1)
	}

	// Inicializa logger
	log := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	log.Info("Iniciando IaC AI Agent", "version", version)

	// Cria handlers
	restHandler := rest.NewHandler(cfg, log)
	webhookHandler := webhook.NewWebhookHandler(cfg, log)

	// Setup rotas
	router := restHandler.SetupRoutes()

	// Adiciona rota de webhook
	router.HandleFunc("/webhook/github", webhookHandler.HandleGitHub).Methods("POST")

	// Cria servidor HTTP
	server := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Inicia servidor em goroutine
	go func() {
		log.Info("Servidor HTTP iniciado", "address", cfg.GetAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Erro no servidor HTTP", "error", err)
			os.Exit(1)
		}
	}()

	// Aguarda sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Desligando servidor...")

	// Shutdown graceful
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Erro ao desligar servidor", "error", err)
	}

	log.Info("Servidor desligado com sucesso")
}
