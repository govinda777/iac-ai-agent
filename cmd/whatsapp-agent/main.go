package main

import (
	"context"
	"fmt"
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

// Exemplo de uso do agente WhatsApp
func exampleUsage() {
	// Criar configuração
	config := &whatsapp.WhatsAppAgentConfig{
		Name:        "Meu Agente WhatsApp",
		Description: "Agente personalizado para análise de código",
		WalletAddr:  "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5",
		WebhookURL:  "https://meudominio.com/webhook/whatsapp",
		VerifyToken: "meu_token_secreto",
	}

	// Criar agente
	agent, err := whatsapp.NewWhatsAppAgent(config)
	if err != nil {
		log.Fatalf("Erro ao criar agente: %v", err)
	}

	// Exemplo de mensagem recebida
	message := &whatsapp.WhatsAppMessage{
		ID:        "msg_123",
		From:      "5511999999999",
		Text:      "/analyze\n```hcl\nresource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}\n```",
		Type:      "text",
		Timestamp: time.Now(),
	}

	// Processar mensagem
	ctx := context.Background()
	response, err := agent.ProcessMessage(ctx, message)
	if err != nil {
		log.Printf("Erro ao processar mensagem: %v", err)
		return
	}

	// Exibir resposta
	fmt.Printf("Resposta: %s\n", response.Text)

	// Obter estatísticas de uso
	stats, err := agent.GetUsageStats(ctx, message.From)
	if err != nil {
		log.Printf("Erro ao obter estatísticas: %v", err)
		return
	}

	fmt.Printf("Estatísticas: %+v\n", stats)
}

// Exemplo de comandos disponíveis
func exampleCommands() {
	fmt.Println("Comandos disponíveis no agente WhatsApp:")

	commands := whatsapp.AvailableCommands()
	for name, cmd := range commands {
		costText := "Gratuito"
		if cmd.RequiresPayment {
			costText = fmt.Sprintf("%d token(s)", cmd.TokenCost)
		}
		fmt.Printf("• %s - %s (%s)\n", name, cmd.Description, costText)
	}
}

// Exemplo de templates de resposta
func exampleTemplates() {
	fmt.Println("Templates de resposta disponíveis:")

	templates := whatsapp.ResponseTemplates
	for name, template := range templates {
		fmt.Printf("• %s: %s\n", name, template[:50]+"...")
	}
}

// Exemplo de configuração
func exampleConfiguration() {
	// Configuração padrão
	config := whatsapp.DefaultConfig()
	fmt.Printf("Configuração padrão: %+v\n", config)

	// Validar configuração
	if err := config.Validate(); err != nil {
		fmt.Printf("Erro na validação: %v\n", err)
	}

	// Mesclar com valores padrão
	config.Name = ""
	config.Description = ""
	config.WalletAddr = ""

	mergedConfig := config.MergeWithDefaults()
	fmt.Printf("Configuração mesclada: %+v\n", mergedConfig)
}

// Exemplo de logging
func exampleLogging() {
	// Criar logger
	logger := whatsapp.SetupLogger("example_agent")

	// Logger estruturado
	structuredLogger := whatsapp.NewStructuredLogger("example_structured", whatsapp.LogLevelInfo)

	// Exemplo de mensagem
	message := &whatsapp.WhatsAppMessage{
		From: "test_user",
		Text: "Hello, world!",
	}

	// Log da mensagem
	logger.LogMessage(message)

	// Log estruturado
	structuredLogger.Info("Message processed", map[string]interface{}{
		"user":    message.From,
		"message": message.Text,
		"time":    time.Now(),
	})
}
