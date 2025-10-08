package capabilities

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/govinda777/iac-ai-agent/internal/agent/core"
)

// WhatsAppCapability habilidade de comunicação via WhatsApp
type WhatsAppCapability struct {
	id           string
	name         string
	description  string
	version      string
	status       string
	lastUsed     time.Time
	messageCount int
	errorCount   int

	// Configurações específicas do WhatsApp
	webhookURL  string
	verifyToken string
	apiKey      string

	// Dependências
	iacCapability *IACAnalysisCapability
	logger        *core.Logger
}

// NewWhatsAppCapability cria nova habilidade WhatsApp
func NewWhatsAppCapability() *WhatsAppCapability {
	return &WhatsAppCapability{
		id:          "whatsapp",
		name:        "WhatsApp Communication",
		description: "Capability to communicate via WhatsApp Business API",
		version:     "1.0.0",
		status:      "inactive",
		logger:      core.NewLogger("whatsapp-capability", "info"),
	}
}

// GetID retorna ID da habilidade
func (w *WhatsAppCapability) GetID() string {
	return w.id
}

// GetName retorna nome da habilidade
func (w *WhatsAppCapability) GetName() string {
	return w.name
}

// GetDescription retorna descrição da habilidade
func (w *WhatsAppCapability) GetDescription() string {
	return w.description
}

// GetVersion retorna versão da habilidade
func (w *WhatsAppCapability) GetVersion() string {
	return w.version
}

// Initialize inicializa a habilidade WhatsApp
func (w *WhatsAppCapability) Initialize(ctx context.Context, config *core.Config) error {
	w.logger = core.NewLogger("whatsapp-capability", config.Logging.Level)

	// Carregar configurações específicas do WhatsApp
	if whatsappConfig, exists := config.Capabilities["whatsapp"]; exists {
		if configMap, ok := whatsappConfig.(map[string]interface{}); ok {
			if webhookURL, exists := configMap["webhook_url"]; exists {
				w.webhookURL = webhookURL.(string)
			}
			if verifyToken, exists := configMap["verify_token"]; exists {
				w.verifyToken = verifyToken.(string)
			}
			if apiKey, exists := configMap["api_key"]; exists {
				w.apiKey = apiKey.(string)
			}
		}
	}

	w.logger.Info("WhatsApp capability initialized", map[string]interface{}{
		"webhook_url":      w.webhookURL,
		"has_verify_token": w.verifyToken != "",
		"has_api_key":      w.apiKey != "",
	})

	return nil
}

// Start inicia a habilidade WhatsApp
func (w *WhatsAppCapability) Start(ctx context.Context) error {
	w.status = "active"
	w.logger.Info("WhatsApp capability started", map[string]interface{}{
		"capability": w.id,
		"status":     w.status,
	})
	return nil
}

// Stop para a habilidade WhatsApp
func (w *WhatsAppCapability) Stop(ctx context.Context) error {
	w.status = "inactive"
	w.logger.Info("WhatsApp capability stopped", map[string]interface{}{
		"capability": w.id,
		"status":     w.status,
	})
	return nil
}

// CanHandle verifica se pode processar a mensagem
func (w *WhatsAppCapability) CanHandle(message *core.Message) bool {
	return message.Source == "whatsapp"
}

// ProcessMessage processa mensagem do WhatsApp
func (w *WhatsAppCapability) ProcessMessage(ctx context.Context, message *core.Message) (*core.Response, error) {
	w.messageCount++
	w.lastUsed = time.Now()

	w.logger.Info("Processing WhatsApp message", map[string]interface{}{
		"message_id":  message.ID,
		"from":        message.From,
		"text_length": len(message.Text),
	})

	// Verificar se é um comando de análise IaC
	if w.isIACCommand(message.Text) {
		return w.processIACCommand(ctx, message)
	}

	// Processar outros comandos WhatsApp
	return w.processWhatsAppCommand(ctx, message)
}

// isIACCommand verifica se é um comando de análise IaC
func (w *WhatsAppCapability) isIACCommand(text string) bool {
	iacCommands := []string{"/analyze", "/security", "/cost", "/terraform", "/iac"}
	text = strings.ToLower(strings.TrimSpace(text))

	for _, cmd := range iacCommands {
		if strings.HasPrefix(text, cmd) {
			return true
		}
	}

	// Verificar se contém blocos de código Terraform
	terraformRegex := regexp.MustCompile("```(hcl|terraform|tf)")
	return terraformRegex.MatchString(text)
}

// processIACCommand processa comando de análise IaC
func (w *WhatsAppCapability) processIACCommand(ctx context.Context, message *core.Message) (*core.Response, error) {
	w.logger.Info("Processing IAC command", map[string]interface{}{
		"message_id": message.ID,
		"command":    message.Text,
	})

	// Se não temos a habilidade IaC registrada, retornar erro
	if w.iacCapability == nil {
		return &core.Response{
			ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
			To:        message.From,
			Text:      "❌ Habilidade de análise IaC não está disponível no momento.",
			Type:      "text",
			Timestamp: time.Now(),
		}, nil
	}

	// Delegar para a habilidade IaC
	iacMessage := &core.Message{
		ID:        message.ID,
		Source:    "whatsapp",
		Channel:   message.Channel,
		From:      message.From,
		To:        message.To,
		Text:      message.Text,
		Type:      message.Type,
		Metadata:  message.Metadata,
		Timestamp: message.Timestamp,
	}

	iacResponse, err := w.iacCapability.ProcessMessage(ctx, iacMessage)
	if err != nil {
		w.errorCount++
		return &core.Response{
			ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
			To:        message.From,
			Text:      fmt.Sprintf("❌ Erro na análise IaC: %v", err),
			Type:      "text",
			Timestamp: time.Now(),
		}, nil
	}

	// Converter resposta IaC para formato WhatsApp
	return &core.Response{
		ID:        iacResponse.ID,
		To:        message.From,
		Text:      iacResponse.Text,
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// processWhatsAppCommand processa comandos específicos do WhatsApp
func (w *WhatsAppCapability) processWhatsAppCommand(ctx context.Context, message *core.Message) (*core.Response, error) {
	text := strings.TrimSpace(message.Text)

	switch {
	case strings.HasPrefix(text, "/help"):
		return w.handleHelpCommand(message)
	case strings.HasPrefix(text, "/status"):
		return w.handleStatusCommand(message)
	case strings.HasPrefix(text, "/ping"):
		return w.handlePingCommand(message)
	default:
		return w.handleUnknownCommand(message)
	}
}

// handleHelpCommand processa comando de ajuda
func (w *WhatsAppCapability) handleHelpCommand(message *core.Message) (*core.Response, error) {
	helpText := `🤖 *IaC AI Agent - Comandos Disponíveis*

*Análise IaC:*
• /analyze - Analisa código Terraform
• /security - Verifica segurança do código
• /cost - Otimiza custos do código

*Comandos Gerais:*
• /help - Lista comandos disponíveis
• /status - Status do agente
• /ping - Testa conectividade

*Como usar:*
Envie seu código Terraform junto com o comando:

/analyze
` + "```hcl" + `
resource "aws_instance" "web" {
  instance_type = "t3.micro"
}
` + "```" + `

💡 *Dica:* O agente pode analisar código Terraform, verificar segurança e otimizar custos automaticamente!`

	return &core.Response{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		To:        message.From,
		Text:      helpText,
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// handleStatusCommand processa comando de status
func (w *WhatsAppCapability) handleStatusCommand(message *core.Message) (*core.Response, error) {
	statusText := fmt.Sprintf(`🤖 *Status do IaC AI Agent*

*Agente:* %s
*Versão:* %s
*Status:* ✅ Ativo
*Habilidade WhatsApp:* ✅ Ativa
*Mensagens processadas:* %d
*Última atividade:* %s

*Habilidades disponíveis:*
• WhatsApp Communication ✅
• IaC Analysis %s

*Sistema:* Operacional`,
		w.name, w.version, w.messageCount, w.lastUsed.Format("15:04:05"),
		func() string {
			if w.iacCapability != nil {
				return "✅"
			}
			return "❌"
		}())

	return &core.Response{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		To:        message.From,
		Text:      statusText,
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// handlePingCommand processa comando ping
func (w *WhatsAppCapability) handlePingCommand(message *core.Message) (*core.Response, error) {
	return &core.Response{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		To:        message.From,
		Text:      "🏓 Pong! Agente WhatsApp ativo e funcionando.",
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// handleUnknownCommand processa comando desconhecido
func (w *WhatsAppCapability) handleUnknownCommand(message *core.Message) (*core.Response, error) {
	return &core.Response{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		To:        message.From,
		Text:      "❓ Comando não reconhecido. Use /help para ver comandos disponíveis.",
		Type:      "text",
		Timestamp: time.Now(),
	}, nil
}

// GetStatus retorna status da habilidade
func (w *WhatsAppCapability) GetStatus() *core.CapabilityStatus {
	return &core.CapabilityStatus{
		ID:           w.id,
		Name:         w.name,
		Status:       w.status,
		LastUsed:     w.lastUsed,
		MessageCount: w.messageCount,
		ErrorCount:   w.errorCount,
	}
}

// HealthCheck verifica saúde da habilidade
func (w *WhatsAppCapability) HealthCheck(ctx context.Context) error {
	if w.status != "active" {
		return fmt.Errorf("capability is not active")
	}

	if w.apiKey == "" {
		return fmt.Errorf("WhatsApp API key not configured")
	}

	return nil
}

// SetIACCapability define a habilidade IaC para integração
func (w *WhatsAppCapability) SetIACCapability(iacCapability *IACAnalysisCapability) {
	w.iacCapability = iacCapability
	w.logger.Info("IAC capability linked to WhatsApp", map[string]interface{}{
		"iac_capability_id": iacCapability.GetID(),
	})
}
