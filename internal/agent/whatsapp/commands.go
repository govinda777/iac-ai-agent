package whatsapp

import (
	"context"
	"fmt"
)

// AvailableCommands retorna comandos disponíveis
func AvailableCommands() map[string]*Command {
	return map[string]*Command{
		"analyze": {
			Name:            "analyze",
			Description:     "Analisa código Terraform",
			Pattern:         `(?s)^/analyze\s*(.*)`,
			Handler:         handleAnalyzeCommand,
			RequiresPayment: true,
			TokenCost:       1,
		},
		"security": {
			Name:            "security",
			Description:     "Verifica segurança do código",
			Pattern:         `(?s)^/security\s*(.*)`,
			Handler:         handleSecurityCommand,
			RequiresPayment: true,
			TokenCost:       1,
		},
		"cost": {
			Name:            "cost",
			Description:     "Otimiza custos do código",
			Pattern:         `(?s)^/cost\s*(.*)`,
			Handler:         handleCostCommand,
			RequiresPayment: true,
			TokenCost:       1,
		},
		"help": {
			Name:            "help",
			Description:     "Lista comandos disponíveis",
			Pattern:         `^/help`,
			Handler:         handleHelpCommand,
			RequiresPayment: false,
			TokenCost:       0,
		},
		"status": {
			Name:            "status",
			Description:     "Status do agente",
			Pattern:         `^/status`,
			Handler:         handleStatusCommand,
			RequiresPayment: false,
			TokenCost:       0,
		},
		"balance": {
			Name:            "balance",
			Description:     "Verifica saldo de tokens",
			Pattern:         `^/balance`,
			Handler:         handleBalanceCommand,
			RequiresPayment: false,
			TokenCost:       0,
		},
		"usage": {
			Name:            "usage",
			Description:     "Estatísticas de uso",
			Pattern:         `^/usage`,
			Handler:         handleUsageCommand,
			RequiresPayment: false,
			TokenCost:       0,
		},
	}
}

// handleAnalyzeCommand processa comando de análise
func handleAnalyzeCommand(agent *WhatsAppAgent, ctx *CommandContext) (*WhatsAppResponse, error) {
	if ctx.CodeBlock == "" {
		return &WhatsAppResponse{
			Text: "❌ Por favor, forneça o código Terraform para análise.\n\nExemplo:\n/analyze\n```hcl\nresource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}\n```",
			Type: "text",
		}, nil
	}

	// Executar análise
	analysis, err := agent.Service.AnalyzeCode(ctx.CodeBlock)
	if err != nil {
		return &WhatsAppResponse{
			Text: fmt.Sprintf("❌ Erro na análise: %v", err),
			Type: "text",
		}, nil
	}

	// Gerar resposta
	response := fmt.Sprintf("✅ Análise concluída!\n\n")
	response += fmt.Sprintf("🔍 Problemas encontrados: %d\n", len(analysis.Issues))

	for _, issue := range analysis.Issues {
		response += fmt.Sprintf("• %s\n", issue.Description)
	}

	if len(analysis.Suggestions) > 0 {
		response += "\n💡 Sugestões:\n"
		for _, suggestion := range analysis.Suggestions {
			response += fmt.Sprintf("• %s\n", suggestion)
		}
	}

	response += fmt.Sprintf("\n💰 Custo: %d token(s) IACAI", 1)

	return &WhatsAppResponse{
		Text: response,
		Type: "text",
	}, nil
}

// handleSecurityCommand processa comando de segurança
func handleSecurityCommand(agent *WhatsAppAgent, ctx *CommandContext) (*WhatsAppResponse, error) {
	if ctx.CodeBlock == "" {
		return &WhatsAppResponse{
			Text: "❌ Por favor, forneça o código para verificação de segurança.\n\nExemplo:\n/security\n```hcl\nresource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}\n```",
			Type: "text",
		}, nil
	}

	// Executar análise de segurança
	securityAnalysis, err := agent.Service.AnalyzeSecurity(ctx.CodeBlock)
	if err != nil {
		return &WhatsAppResponse{
			Text: fmt.Sprintf("❌ Erro na análise de segurança: %v", err),
			Type: "text",
		}, nil
	}

	// Gerar resposta
	response := fmt.Sprintf("🔒 Análise de Segurança Concluída!\n\n")
	response += fmt.Sprintf("⚠️ Vulnerabilidades encontradas: %d\n", len(securityAnalysis.Vulnerabilities))

	for _, vuln := range securityAnalysis.Vulnerabilities {
		response += fmt.Sprintf("• %s (Severidade: %s)\n", vuln.Description, vuln.Severity)
	}

	if len(securityAnalysis.Recommendations) > 0 {
		response += "\n🛡️ Recomendações de Segurança:\n"
		for _, rec := range securityAnalysis.Recommendations {
			response += fmt.Sprintf("• %s\n", rec)
		}
	}

	response += fmt.Sprintf("\n💰 Custo: %d token(s) IACAI", 1)

	return &WhatsAppResponse{
		Text: response,
		Type: "text",
	}, nil
}

// handleCostCommand processa comando de otimização de custos
func handleCostCommand(agent *WhatsAppAgent, ctx *CommandContext) (*WhatsAppResponse, error) {
	if ctx.CodeBlock == "" {
		return &WhatsAppResponse{
			Text: "❌ Por favor, forneça o código para análise de custos.\n\nExemplo:\n/cost\n```hcl\nresource \"aws_instance\" \"web\" {\n  instance_type = \"t3.micro\"\n}\n```",
			Type: "text",
		}, nil
	}

	// Executar análise de custos
	costAnalysis, err := agent.Service.AnalyzeCosts(ctx.CodeBlock)
	if err != nil {
		return &WhatsAppResponse{
			Text: fmt.Sprintf("❌ Erro na análise de custos: %v", err),
			Type: "text",
		}, nil
	}

	// Gerar resposta
	response := fmt.Sprintf("💰 Análise de Custos Concluída!\n\n")
	response += fmt.Sprintf("📊 Custo estimado mensal: $%.2f\n", costAnalysis.EstimatedMonthlyCost)
	response += fmt.Sprintf("💡 Potencial de economia: $%.2f\n", costAnalysis.PotentialSavings)

	if len(costAnalysis.Optimizations) > 0 {
		response += "\n🔧 Otimizações sugeridas:\n"
		for _, opt := range costAnalysis.Optimizations {
			response += fmt.Sprintf("• %s (Economia: $%.2f/mês)\n", opt.Description, opt.MonthlySavings)
		}
	}

	response += fmt.Sprintf("\n💰 Custo: %d token(s) IACAI", 1)

	return &WhatsAppResponse{
		Text: response,
		Type: "text",
	}, nil
}

// handleHelpCommand processa comando de ajuda
func handleHelpCommand(agent *WhatsAppAgent, ctx *CommandContext) (*WhatsAppResponse, error) {
	response := ResponseTemplates["welcome"]
	response += "\n\n📋 Comandos Disponíveis:\n\n"

	for _, cmd := range agent.Commands {
		costText := "🆓 Gratuito"
		if cmd.RequiresPayment {
			costText = fmt.Sprintf("💰 %d token(s)", cmd.TokenCost)
		}
		response += fmt.Sprintf("• %s - %s (%s)\n", cmd.Name, cmd.Description, costText)
	}

	response += "\n💡 Dica: Envie seu código Terraform junto com o comando para análise!"

	return &WhatsAppResponse{
		Text: response,
		Type: "text",
	}, nil
}

// handleStatusCommand processa comando de status
func handleStatusCommand(agent *WhatsAppAgent, ctx *CommandContext) (*WhatsAppResponse, error) {
	response := fmt.Sprintf("🤖 Status do Agente\n\n")
	response += fmt.Sprintf("Nome: %s\n", agent.Name)
	response += fmt.Sprintf("ID: %s\n", agent.ID)
	response += fmt.Sprintf("Descrição: %s\n", agent.Description)
	response += fmt.Sprintf("Wallet: %s\n", agent.WalletAddr)
	response += fmt.Sprintf("Status: ✅ Online\n")
	response += fmt.Sprintf("Comandos disponíveis: %d\n", len(agent.Commands))

	return &WhatsAppResponse{
		Text: response,
		Type: "text",
	}, nil
}

// handleBalanceCommand processa comando de saldo
func handleBalanceCommand(agent *WhatsAppAgent, ctx *CommandContext) (*WhatsAppResponse, error) {
	// Por enquanto, retorna saldo simulado
	// Em produção, consultar blockchain
	response := fmt.Sprintf("💰 Saldo de Tokens IACAI\n\n")
	response += fmt.Sprintf("Usuário: %s\n", ctx.Message.From)
	response += fmt.Sprintf("Saldo atual: 100 tokens\n")
	response += fmt.Sprintf("Última atualização: Agora\n\n")
	response += fmt.Sprintf("💡 Use /usage para ver estatísticas detalhadas")

	return &WhatsAppResponse{
		Text: response,
		Type: "text",
	}, nil
}

// handleUsageCommand processa comando de estatísticas de uso
func handleUsageCommand(agent *WhatsAppAgent, ctx *CommandContext) (*WhatsAppResponse, error) {
	stats, err := agent.GetUsageStats(context.Background(), ctx.Message.From)
	if err != nil {
		return &WhatsAppResponse{
			Text: fmt.Sprintf("❌ Erro ao obter estatísticas: %v", err),
			Type: "text",
		}, nil
	}

	response := fmt.Sprintf("📊 Estatísticas de Uso\n\n")
	response += fmt.Sprintf("Total de requisições: %d\n", stats.TotalRequests)
	response += fmt.Sprintf("Tokens consumidos: %d\n", stats.TokensConsumed)
	response += fmt.Sprintf("Requisições hoje: %d\n", stats.RequestsToday)
	response += fmt.Sprintf("Custo médio: %.2f tokens\n", stats.AverageCost)
	response += fmt.Sprintf("Última requisição: %s\n", stats.LastRequest.Format("02/01/2006 15:04"))

	return &WhatsAppResponse{
		Text: response,
		Type: "text",
	}, nil
}
