package middleware

import (
	"auth-service/internal/shared/consts"
	"auth-service/internal/shared/response"
	token "auth-service/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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
			response.InternalServerError(c, "internal_server_error")
			return
		}

		if context.Reached {
			response.TooManyRequests(c, "too_many_request")
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
			response.Unauthorized(c, "unauthorized")
			return
		}

		ctx := c.Request.Context()
		blacklisted, _ := token.IsTokenBlacklisted(ctx, bearerToken)
		if blacklisted {
			response.Unauthorized(c, "unauthorized")
			return
		}

		accessClaims, err := token.ParseAccessToken(bearerToken)
		if err != nil {
			response.Unauthorized(c, "unauthorized")
			return
		}

		c.Set(consts.AccessToken, bearerToken)
		c.Set(consts.UserId, accessClaims.UserID)
		c.Set(consts.Username, accessClaims.Username)
		c.Set(consts.Email, accessClaims.Email)

		c.Next()
	}
}
