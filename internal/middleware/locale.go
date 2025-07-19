package middleware

import (
	"auth-service/internal/shared/consts"
	locale "auth-service/pkg/i18n"
	"github.com/gin-gonic/gin"
)

func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		accept := c.GetHeader("Accept-Language")

		c.Set(consts.Localizer, locale.GetLocalizer(lang, accept))

		c.Next()
	}
}
