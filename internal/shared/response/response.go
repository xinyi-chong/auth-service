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

func Error(c *gin.Context, err error, overrideErr ...*apperrors.Error) {
	var override *apperrors.Error
	if len(overrideErr) > 0 {
		override = overrideErr[0]
	}

	appErr := getDisplayErr(err, override)
	message := locale.Translate(c, locale.CategoryError, appErr.MessageKey, appErr.TemplateData)

	c.AbortWithStatusJSON(appErr.HTTPStatus, Response{Message: message})
}

func getDisplayErr(err error, overrideErr *apperrors.Error) *apperrors.Error {
	var appErr *apperrors.Error
	if overrideErr != nil {
		if !errors.As(err, &appErr) || errors.Is(err, apperrors.ErrInternalServerError) {
			return overrideErr.Wrap(err)
		}
	}

	if errors.As(err, &appErr) {
		return appErr
	}

	return apperrors.ErrInternalServerError.Wrap(err)
}
