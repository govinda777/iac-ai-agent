package whatsapp

// ResponseTemplates templates de resposta padronizados
var ResponseTemplates = map[string]string{
	"welcome": `🤖 Olá! Sou o IaC AI Agent.

Posso ajudar você a analisar código Terraform, verificar segurança e otimizar custos.

Comandos disponíveis:
• /analyze - Analisa código Terraform
• /security - Verifica segurança
• /cost - Otimiza custos
• /help - Lista comandos
• /status - Status do agente

Envie seu código e eu farei a análise!`,

	"error": `❌ Ops! Algo deu errado.

Verifique se:
• O comando está correto
• O código está bem formatado
• Você tem tokens suficientes

Use /help para ver comandos disponíveis.`,

	"insufficient_tokens": `💰 Saldo insuficiente!

Você precisa de mais tokens IACAI para usar este serviço.

Comandos gratuitos:
• /help
• /status

Para análises, você precisa de tokens IACAI.`,

	"analysis_started": `🔄 Iniciando análise...

Por favor, aguarde enquanto analiso seu código Terraform.`,

	"analysis_complete": `✅ Análise concluída!

Seu código foi analisado com sucesso.`,

	"security_scan_started": `🔒 Iniciando verificação de segurança...

Analisando vulnerabilidades e problemas de segurança.`,

	"security_scan_complete": `🛡️ Verificação de segurança concluída!

Análise de segurança finalizada.`,

	"cost_analysis_started": `💰 Iniciando análise de custos...

Calculando custos e identificando oportunidades de economia.`,

	"cost_analysis_complete": `📊 Análise de custos concluída!

Relatório de custos gerado com sucesso.`,

	"invalid_command": `❌ Comando inválido!

Use /help para ver todos os comandos disponíveis.`,

	"no_code_provided": `❌ Nenhum código fornecido!

Por favor, envie o código Terraform junto com o comando.

Exemplo:
/analyze
` + "```hcl" + `
resource "aws_instance" "web" {
  instance_type = "t3.micro"
}
` + "```",

	"processing_error": `❌ Erro no processamento!

Ocorreu um erro interno. Tente novamente em alguns instantes.

Se o problema persistir, entre em contato com o suporte.`,

	"rate_limit_exceeded": `⏰ Limite de requisições excedido!

Você atingiu o limite de requisições por hora.

Tente novamente em alguns minutos.`,

	"authentication_failed": `🔐 Falha na autenticação!

Verifique se:
• Sua wallet está conectada
• Você possui o NFT necessário
• Sua assinatura é válida`,

	"billing_error": `💳 Erro na cobrança!

Não foi possível processar o pagamento.

Verifique seu saldo de tokens IACAI.`,

	"maintenance_mode": `🔧 Modo de manutenção!

O sistema está temporariamente indisponível para manutenção.

Tente novamente em alguns minutos.`,

	"feature_unavailable": `🚧 Funcionalidade indisponível!

Esta funcionalidade ainda não está disponível.

Use /help para ver comandos disponíveis.`,
}
