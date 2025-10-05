package logger

import (
	"io"
	"log/slog"
	"os"
)

// Logger é um wrapper para o logger estruturado
type Logger struct {
	logger *slog.Logger
}

// New cria um novo logger
func New(level, format string) *Logger {
	var handler slog.Handler
	var logLevel slog.Level

	// Parse level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	// Escolhe formato
	var output io.Writer = os.Stdout
	if format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	return &Logger{
		logger: slog.New(handler),
	}
}

// Debug registra mensagem de debug
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

// Info registra mensagem informativa
func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Warn registra warning
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

// Error registra erro
func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

// With adiciona campos contextuais ao logger
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}
