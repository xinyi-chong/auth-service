package response

import (
	apperrors "auth-service/pkg/error"
	locale "auth-service/pkg/i18n"
	"auth-service/pkg/success"
	"errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, success *success.Success, data interface{}) {
	message := locale.Translate(c, locale.CategorySuccess, success.Code, success.TemplateData)

	c.JSON(success.HTTPStatus, Response{
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	var appErr *apperrors.Error
	ok := errors.As(err, &appErr)
	if !ok {
		appErr = apperrors.ErrInternalServerError
	}

	message := locale.Translate(c, locale.CategoryError, appErr.Code, appErr.TemplateData)

	c.JSON(appErr.HTTPStatus, Response{
		Message: message,
	})
}
