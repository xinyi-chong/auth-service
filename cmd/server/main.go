package main

import (
	_ "auth-service/docs"
	"auth-service/internal/api"
	"auth-service/pkg/logger"
	"go.uber.org/zap"
	"os"
)

// @title Auth Service
// @version 1.0
// @description JWT Authentication
// @securityDefinitions.apikey BearerTokenAuth
// @in header
// @name Authorization
// @description Bearer Token Authentication. Example: `Bearer <your-jwt-token>`
func main() {
	if err := logger.Init(); err != nil {
		logger.Fatal("Failed to initialize logger", zap.Error(err))
	}
	defer logger.Log.Sync()

	server, err := api.NewServer()
	if err != nil {
		logger.Fatal("Failed to initialize server", zap.Error(err))
	}

	if err := server.Start(); err != nil {
		logger.Error("Server exited with error", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
