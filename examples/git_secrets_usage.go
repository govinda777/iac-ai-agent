package main

import (
	"fmt"
	"log"

	"github.com/govinda777/iac-ai-agent/pkg/config"
)

func main() {
	// Carregar configuração
	cfg, err := config.Load("configs/app.yaml")
	if err != nil {
		log.Fatal("Erro ao carregar configuração:", err)
	}

	// Exemplo de uso dos secrets via Git Secrets
	fmt.Println("=== Acessando Secrets via Git Secrets ===")

	// GitHub Token
	if token, err := cfg.GetGitHubToken(); err != nil {
		fmt.Printf("GitHub Token: ❌ %v\n", err)
	} else {
		fmt.Printf("GitHub Token: ✅ %s...\n", token[:10])
	}

	// GitHub Webhook Secret
	if secret, err := cfg.GetGitHubWebhookSecret(); err != nil {
		fmt.Printf("GitHub Webhook Secret: ❌ %v\n", err)
	} else {
		fmt.Printf("GitHub Webhook Secret: ✅ %s...\n", secret[:10])
	}

	// Chave privada da carteira
	if privateKey, err := cfg.GetWalletPrivateKey(); err != nil {
		fmt.Printf("Wallet Private Key: ❌ %v\n", err)
	} else {
		fmt.Printf("Wallet Private Key: ✅ %s...\n", privateKey[:10])
	}

	// WhatsApp API Key
	if whatsappKey, err := cfg.GetWhatsAppAPIKey(); err != nil {
		fmt.Printf("WhatsApp API Key: ❌ %v\n", err)
	} else {
		fmt.Printf("WhatsApp API Key: ✅ %s...\n", whatsappKey[:10])
	}

	fmt.Println("\n=== Configuração Carregada ===")
	fmt.Printf("Servidor: %s\n", cfg.GetAddress())
	fmt.Printf("LLM Provider: %s\n", cfg.LLM.Provider)
	fmt.Printf("LLM Model: %s\n", cfg.LLM.Model)
	fmt.Printf("GitHub Auto Comment: %t\n", cfg.GitHub.AutoComment)
}
