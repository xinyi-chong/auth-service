package api

import (
	"auth-service/db"
	"auth-service/internal/config"
	"auth-service/internal/middleware"
	locale "auth-service/pkg/i18n"
	token "auth-service/pkg/jwt"
	"auth-service/pkg/logger"
	redisclient "auth-service/pkg/redis"
	"net/http"
	"os"

	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	redis  *redis.Client
	config *config.Config
	logger *zap.Logger
}

func NewServer() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	gormDB, err := db.Init(cfg.Postgres.Config)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisclient.Init(cfg.Redis.Config)
	if err != nil {
		return nil, err
	}

	err = token.Init(cfg)
	if err != nil {
		return nil, err
	}

	err = locale.Init()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = gormDB.WithContext(ctx).Exec("SELECT 1").Error
	if err != nil {
		return nil, err
	}

	server := &Server{
		router: gin.New(),
		db:     gormDB,
		redis:  redisClient,
		config: cfg,
		logger: logger.Log,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server, nil
}

func (s *Server) Start() error {
	addr := ":" + s.config.Server.Port
	s.logger.Info("Starting server",
		zap.String("address", addr),
		zap.String("environment", os.Getenv("APP_ENV")),
	)

	return s.router.Run(addr)
}

func (s *Server) setupMiddleware() {
	s.router.Use(
		ginzap.Ginzap(logger.Log, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger.Log, true),
		middleware.CORSMiddleware(),
		middleware.RateLimit(100, time.Hour),
		middleware.LocaleMiddleware(),
	)
}

func (s *Server) setupRoutes() {
	s.router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
