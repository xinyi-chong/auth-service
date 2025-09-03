package middleware

import (
	"auth-service/internal/shared/consts"
	token "auth-service/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	apperrors "github.com/xinyi-chong/common-lib/errors"
	"github.com/xinyi-chong/common-lib/logger"
	"github.com/xinyi-chong/common-lib/response"
	"go.uber.org/zap"
	"strings"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Language"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func RateLimit(requests int, duration time.Duration) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: duration,
		Limit:  int64(requests),
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		context, err := instance.Get(c, c.ClientIP())
		if err != nil {
			response.Error(c, apperrors.ErrInternalServerError)
			return
		}

		if context.Reached {
			response.Error(c, apperrors.ErrTooManyRequests)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var bearerToken string
		if strings.HasPrefix(authHeader, "Bearer ") {
			bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			logger.Debug("missing or invalid Authorization header")
			response.Error(c, apperrors.ErrUnauthorized)
			return
		}

		ctx := c.Request.Context()
		blacklisted, _ := token.IsTokenBlacklisted(ctx, bearerToken)
		if blacklisted {
			logger.Debug("token is blacklisted", zap.String("token", bearerToken))
			response.Error(c, apperrors.ErrUnauthorized)
			return
		}

		accessClaims, err := token.ParseAccessToken(bearerToken)
		if err != nil {
			logger.Debug("invalid access token", zap.Error(err))
			response.Error(c, apperrors.ErrUnauthorized)
			return
		}

		logger.Debug("access claims", zap.Any("claims", accessClaims))

		c.Set(consts.AccessToken, bearerToken)
		c.Set(consts.UserId, accessClaims.UserID)
		c.Set(consts.Username, accessClaims.Username)
		c.Set(consts.Email, accessClaims.Email)

		c.Next()
	}
}
