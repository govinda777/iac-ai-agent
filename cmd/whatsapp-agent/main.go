package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/govinda777/iac-ai-agent/api/rest"
	"github.com/govinda777/iac-ai-agent/internal/agent/whatsapp"
)

func main() {
	// Carregar configuração
	config, err := whatsapp.LoadOrCreateConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validar configuração
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Criar agente WhatsApp
	agent, err := whatsapp.NewWhatsAppAgent(config)
	if err != nil {
		log.Fatalf("Failed to create WhatsApp agent: %v", err)
	}

	log.Printf("WhatsApp Agent created successfully: %s", agent.ID)

	// Criar handler de webhook
	webhookHandler := rest.NewWhatsAppWebhookHandler(agent)

	// Configurar rotas
	router := mux.NewRouter()

	// Middleware
	router.Use(rest.LoggingMiddleware)
	router.Use(rest.TokenValidationMiddleware(config.VerifyToken))

	// Registrar rotas
	webhookHandler.RegisterRoutes(router)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Starting WhatsApp webhook server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

