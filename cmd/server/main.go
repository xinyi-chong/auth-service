package main

import (
	_ "auth-service/docs"
	"auth-service/internal/api"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/xinyi-chong/common-lib/logger"
	"go.uber.org/zap"
	"os"
)

// @title Auth Service
// @version 1.0
// @host localhost:8080
// @basePath /api/v1
// @description JWT Authentication
// @securityDefinitions.apikey BearerTokenAuth
// @in header
// @name Authorization
// @description Bearer Token Authentication. Example: `Bearer <your-jwt-token>`
func main() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	if appEnv == "local" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}
	}

	if err := logger.Init(); err != nil {
		fmt.Printf("FATAL: logger init failed: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

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
