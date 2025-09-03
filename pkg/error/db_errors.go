package dberrors

import (
	"errors"
	apperrors "github.com/xinyi-chong/common-lib/errors"

	"gorm.io/gorm"
)

type Model string

const (
	ModelUser Model = "user"
)

func getModelNotFoundError(model Model) *apperrors.Error {
	switch model {
	case ModelUser:
		return apperrors.ErrUserNotFound
	default:
		return apperrors.ErrInternalServerError
	}
}

func getModelConflictError(model Model) *apperrors.Error {
	switch model {
	case ModelUser:
		return apperrors.ErrUserConflict
	default:
		return apperrors.ErrInternalServerError
	}
}

func WrapDBError(err error, model Model) *apperrors.Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return getModelNotFoundError(model).Wrap(err)
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return getModelConflictError(model).Wrap(err)
	}

	return apperrors.ErrInternalServerError.Wrap(err)
}
