package api

import (
	"auth-service/db"
	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/middleware"
	"auth-service/internal/user"
	locale "auth-service/pkg/i18n"
	token "auth-service/pkg/jwt"
	"auth-service/pkg/logger"
	redisclient "auth-service/pkg/redis"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"

	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	redis  *redis.Client
	config *config.Config
	logger *zap.Logger

	authCtrl *auth.Controller
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

	log := logger.Get()
	userRepo := user.NewRepository(gormDB)
	userSvc := user.NewService(userRepo, log)
	authSvc := auth.NewService(userSvc, log)
	authCtrl := auth.NewController(authSvc, log)

	s := &Server{
		router:   gin.New(),
		db:       gormDB,
		redis:    redisClient,
		config:   cfg,
		logger:   log,
		authCtrl: authCtrl,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
}

func (s *Server) Start() error {
	addr := ":" + s.config.Server.Port
	s.logger.Info("Starting server",
		zap.String("addr", addr),
		zap.String("environment", os.Getenv("APP_ENV")))
	return s.router.Run(addr)
}

func (s *Server) setupMiddleware() {
	s.router.Use(
		func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
				ginzap.Ginzap(s.logger, time.RFC3339, true)(c)
			}
			c.Next()
		},
		ginzap.RecoveryWithZap(s.logger, true),
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

	api := s.router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			s.registerAuthRoutes(v1)
		}
	}
}

func (s *Server) registerAuthRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	{
		g.POST("/register", s.authCtrl.Register)
		g.POST("/login", s.authCtrl.Login)
		g.POST("/refresh", s.authCtrl.RefreshToken)
		g.Use(middleware.AuthMiddleware())
		g.PATCH("/change-password", s.authCtrl.ChangePassword)
		g.POST("/logout", s.authCtrl.Logout)
	}
}
