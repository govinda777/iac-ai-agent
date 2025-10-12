package handler

import (
	"net/http"

	"github.com/govinda777/iac-ai-agent/api/rest"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Create configuration and logger for the handler REST
	restConfig := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
	restLogger := logger.New("info", "json")

	// Create handler REST
	handler := rest.NewHandler(restConfig, restLogger)

	// Set up routes
	router := handler.SetupRoutes()

	router.ServeHTTP(w, r)
}