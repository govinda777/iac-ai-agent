package whatsapp

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogMessage logs incoming messages
func (l *WhatsAppLogger) LogMessage(msg *WhatsAppMessage) {
	log.Printf("[%s] WhatsApp Message - From: %s, Text: %s, Time: %s",
		l.AgentID, msg.From, msg.Text, time.Now().Format(time.RFC3339))
}

// LogResponse logs outgoing responses
func (l *WhatsAppLogger) LogResponse(response *WhatsAppResponse) {
	log.Printf("[%s] WhatsApp Response - Text: %s, Time: %s",
		l.AgentID, response.Text, time.Now().Format(time.RFC3339))
}

// LogError logs errors
func (l *WhatsAppLogger) LogError(err error, context string) {
	log.Printf("[%s] WhatsApp Error - Context: %s, Error: %v, Time: %s",
		l.AgentID, context, err, time.Now().Format(time.RFC3339))
}

// LogBilling logs billing events
func (l *WhatsAppLogger) LogBilling(userAddr string, amount int, txHash string) {
	log.Printf("[%s] WhatsApp Billing - User: %s, Amount: %d, TX: %s, Time: %s",
		l.AgentID, userAddr, amount, txHash, time.Now().Format(time.RFC3339))
}

// LogCommand logs command execution
func (l *WhatsAppLogger) LogCommand(command string, userAddr string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	log.Printf("[%s] WhatsApp Command - Command: %s, User: %s, Status: %s, Time: %s",
		l.AgentID, command, userAddr, status, time.Now().Format(time.RFC3339))
}

// LogAuthentication logs authentication events
func (l *WhatsAppLogger) LogAuthentication(userAddr string, success bool, method string) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	log.Printf("[%s] WhatsApp Auth - User: %s, Method: %s, Status: %s, Time: %s",
		l.AgentID, userAddr, method, status, time.Now().Format(time.RFC3339))
}

// LogRateLimit logs rate limiting events
func (l *WhatsAppLogger) LogRateLimit(userAddr string, limit int, current int) {
	log.Printf("[%s] WhatsApp RateLimit - User: %s, Limit: %d, Current: %d, Time: %s",
		l.AgentID, userAddr, limit, current, time.Now().Format(time.RFC3339))
}

// LogPerformance logs performance metrics
func (l *WhatsAppLogger) LogPerformance(operation string, duration time.Duration, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	log.Printf("[%s] WhatsApp Performance - Operation: %s, Duration: %v, Status: %s, Time: %s",
		l.AgentID, operation, duration, status, time.Now().Format(time.RFC3339))
}

// LogWebhook logs webhook events
func (l *WhatsAppLogger) LogWebhook(event string, details string) {
	log.Printf("[%s] WhatsApp Webhook - Event: %s, Details: %s, Time: %s",
		l.AgentID, event, details, time.Now().Format(time.RFC3339))
}

// LogSystem logs system events
func (l *WhatsAppLogger) LogSystem(event string, details string) {
	log.Printf("[%s] WhatsApp System - Event: %s, Details: %s, Time: %s",
		l.AgentID, event, details, time.Now().Format(time.RFC3339))
}

// SetupLogger configura o sistema de logging
func SetupLogger(agentID string) *WhatsAppLogger {
	// Configurar formato de log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Criar arquivo de log se necessário
	logFile := os.Getenv("WHATSAPP_LOG_FILE")
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open log file: %v", err)
		} else {
			log.SetOutput(file)
		}
	}

	return &WhatsAppLogger{
		AgentID: agentID,
	}
}

// LogLevel define níveis de log
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// StructuredLogger logger estruturado para WhatsApp
type StructuredLogger struct {
	AgentID  string
	LogLevel LogLevel
}

// NewStructuredLogger cria novo logger estruturado
func NewStructuredLogger(agentID string, level LogLevel) *StructuredLogger {
	return &StructuredLogger{
		AgentID:  agentID,
		LogLevel: level,
	}
}

// LogWithLevel logs com nível específico
func (sl *StructuredLogger) LogWithLevel(level LogLevel, message string, fields map[string]interface{}) {
	if level < sl.LogLevel {
		return
	}

	// Construir mensagem estruturada
	logMessage := fmt.Sprintf("[%s] [%s] %s", sl.AgentID, level.String(), message)

	// Adicionar campos
	for key, value := range fields {
		logMessage += fmt.Sprintf(" %s=%v", key, value)
	}

	logMessage += fmt.Sprintf(" time=%s", time.Now().Format(time.RFC3339))

	log.Println(logMessage)
}

// Debug logs mensagem de debug
func (sl *StructuredLogger) Debug(message string, fields map[string]interface{}) {
	sl.LogWithLevel(LogLevelDebug, message, fields)
}

// Info logs mensagem informativa
func (sl *StructuredLogger) Info(message string, fields map[string]interface{}) {
	sl.LogWithLevel(LogLevelInfo, message, fields)
}

// Warn logs mensagem de aviso
func (sl *StructuredLogger) Warn(message string, fields map[string]interface{}) {
	sl.LogWithLevel(LogLevelWarn, message, fields)
}

// Error logs mensagem de erro
func (sl *StructuredLogger) Error(message string, fields map[string]interface{}) {
	sl.LogWithLevel(LogLevelError, message, fields)
}

// String converte LogLevel para string
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
