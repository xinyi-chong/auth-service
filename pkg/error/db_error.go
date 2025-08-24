package apperrors

import (
	"errors"

	"gorm.io/gorm"
)

type Model string

const (
	ModelUser Model = "user"
)

func HandleDbFindErr(err error, model Model) *Error {
	if isDbErrNotFound(err) {
		return getModelNotFoundError(model).Wrap(err)
	}
	return ErrInternalServerError.Wrap(err)
}

func HandleDbCreateErr(err error, model Model) *Error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return getModelConflictError(model).Wrap(err)
	}
	return ErrInternalServerError.Wrap(err)
}

func HandleDbUpdateErr(err error, model Model) *Error {
	if isDbErrNotFound(err) {
		return getModelNotFoundError(model).Wrap(err)
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return getModelConflictError(model).Wrap(err)
	}
	return ErrInternalServerError.Wrap(err)
}

func HandleDbDeleteErr(err error, model Model) *Error {
	if isDbErrNotFound(err) {
		return getModelNotFoundError(model).Wrap(err)
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return getModelConflictError(model).Wrap(err)
	}
	return ErrInternalServerError.Wrap(err)
}

func isDbErrNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func getModelNotFoundError(model Model) *Error {
	switch model {
	case ModelUser:
		return ErrUserNotFound
	default:
		return ErrInternalServerError
	}
}

func getModelConflictError(model Model) *Error {
	switch model {
	case ModelUser:
		return ErrUserConflict
	default:
		return ErrInternalServerError
	}
}
