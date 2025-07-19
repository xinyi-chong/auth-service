package response

import (
	locale "auth-service/pkg/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func newSuccess(c *gin.Context, data interface{}, messageID string, templateData ...locale.TemplateData) Response {
	message := locale.Translate(c, messageID, templateData...)
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func newError(c *gin.Context, messageID string, templateData ...locale.TemplateData) Response {
	messageID = "errors." + messageID
	message := locale.Translate(c, messageID, templateData...)
	return Response{
		Success: false,
		Message: message,
	}
}

func Success(c *gin.Context, data interface{}, messageID string, templateData ...locale.TemplateData) {
	c.JSON(http.StatusOK, newSuccess(c, data, messageID, templateData...))
}

func Created(c *gin.Context, data interface{}, messageID string, templateData ...locale.TemplateData) {
	c.JSON(http.StatusCreated, newSuccess(c, data, messageID, templateData...))
}

func BadRequest(c *gin.Context, messageID string, templateData ...locale.TemplateData) {
	c.JSON(http.StatusBadRequest, newError(c, messageID, templateData...))
}

func NotFound(c *gin.Context, messageID string, templateData ...locale.TemplateData) {
	c.JSON(http.StatusNotFound, newError(c, messageID, templateData...))
}

func InternalServerError(c *gin.Context, messageID string, templateData ...locale.TemplateData) {
	c.JSON(http.StatusInternalServerError, newError(c, messageID, templateData...))
}

func TooManyRequests(c *gin.Context, messageID string, templateData ...locale.TemplateData) {
	c.AbortWithStatusJSON(http.StatusTooManyRequests, newError(c, messageID, templateData...))
}

func Unauthorized(c *gin.Context, messageID string, templateData ...locale.TemplateData) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, newError(c, messageID, templateData...))
}
