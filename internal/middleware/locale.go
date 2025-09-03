package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xinyi-chong/common-lib/consts"
	locale "github.com/xinyi-chong/common-lib/i18n"
)

func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		accept := c.GetHeader("Accept-Language")

		c.Set(consts.Localizer, locale.GetLocalizer(lang, accept))

		c.Next()
	}
}
